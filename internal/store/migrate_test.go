package store_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/testdb"
)

func TestRangeGroupsBecomeRanges(t *testing.T) {
	st := testdb.Bare(t)
	ctx := context.Background()
	if err := st.MigrateUntil(ctx, "0025_ranking_not_printed.sql"); err != nil {
		t.Fatal(err)
	}
	db := st.Pool()
	var tid, rid, g25, g50, open, closed, shooter, item int64
	scan := func(dst *int64, sql string, args ...any) {
		t.Helper()
		if err := db.QueryRow(ctx, sql, args...).Scan(dst); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
	}
	exec := func(sql string, args ...any) {
		t.Helper()
		if _, err := db.Exec(ctx, sql, args...); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
	}
	scan(&tid, `INSERT INTO tenants (name) VALUES ('Club') RETURNING id`)
	scan(&rid, `INSERT INTO ranges (tenant_id, name, lane_count, disabled_lanes) VALUES ($1, 'Pistol', 15, '{3,12}') RETURNING id`, tid)
	scan(&g25, `INSERT INTO range_groups (tenant_id, range_id, name, lanes, position) VALUES ($1, $2, '25 m', '{1,2,3,4,5,6,7,8,9,10}', 0) RETURNING id`, tid, rid)
	scan(&g50, `INSERT INTO range_groups (tenant_id, range_id, name, lanes, position) VALUES ($1, $2, '50 m', '{11,12,13,14,15}', 1) RETURNING id`, tid, rid)
	scan(&open, `INSERT INTO events (tenant_id, name, starts_on, days) VALUES ($1, 'Open', '2026-11-20', '{2026-11-20}') RETURNING id`, tid)
	scan(&closed, `INSERT INTO events (tenant_id, name, starts_on, days) VALUES ($1, 'Closed', '2025-11-20', '{2025-11-20}') RETURNING id`, tid)
	exec(`INSERT INTO event_ranges (tenant_id, event_id, range_id, position) VALUES ($1, $2, $3, 0), ($1, $4, $3, 0)`, tid, open, rid, closed)
	scan(&shooter, `INSERT INTO shooters (tenant_id, firstname, lastname) VALUES ($1, 'Anna', 'Muster') RETURNING id`, tid)
	exec(`INSERT INTO registrations (tenant_id, event_id, shooter_id, age_category, status, range_id, lane)
		VALUES ($1, $2, $3, 'elite', 'registered', $4, 4), ($1, $5, $3, 'elite', 'completed', $4, 12)`, tid, open, shooter, rid, closed)
	var p50, pAny int64
	scan(&p50, `INSERT INTO programmes (tenant_id, event_id, name, range_id, range_group_id) VALUES ($1, $2, 'Precision', $3, $4) RETURNING id`, tid, open, rid, g50)
	scan(&pAny, `INSERT INTO programmes (tenant_id, event_id, name, range_id) VALUES ($1, $2, 'Any', $3) RETURNING id`, tid, open, rid)
	exec(`UPDATE events SET archived_at = now() WHERE id = $1`, closed)
	var req int64
	scan(&req, `INSERT INTO booking_requests (tenant_id, event_id, club, contact, email, manage_hash) VALUES ($1, $2, 'Guests', 'Bea', 'b@example.org', 'h') RETURNING id`, tid, open)
	scan(&item, `INSERT INTO booking_items (tenant_id, request_id, range_id, starts_at, ends_at, lanes) VALUES ($1, $2, $3, '2026-11-20 08:00Z', '2026-11-20 09:00Z', 3) RETURNING id`, tid, req, rid)
	exec(`INSERT INTO reservations (tenant_id, range_id, lanes, starts_at, ends_at, description, booking_item_id)
		VALUES ($1, $2, '{9,10,11}', '2026-11-20 08:00Z', '2026-11-20 09:00Z', 'Guests', $3)`, tid, rid, item)
	exec(`INSERT INTO helper_codes (tenant_id, event_id, name, role, token_hash, range_id, lanes) VALUES ($1, $2, 'Scorer', 'scorer', 'x', $3, '{11,12}')`, tid, open, rid)
	exec(`INSERT INTO print_settings (tenant_id, range_id, document, printer) VALUES ($1, $2, 'lane_sheet', 'A5')`, tid, rid)

	if err := st.MigrateUntil(ctx, ""); err != nil {
		t.Fatal(err)
	}

	var r50 int64
	scan(&r50, `SELECT id FROM ranges WHERE name = 'Pistol 50 m'`)
	query := func(sql string) string {
		t.Helper()
		var s string
		if err := db.QueryRow(ctx, sql, pgx.NamedArgs{"r50": r50}).Scan(&s); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
		return s
	}
	for _, c := range []struct{ name, sql, want string }{
		{"ranges", `SELECT string_agg(format('%s:%s:%s', name, lane_count, disabled_lanes), ' ' ORDER BY id) FROM ranges`, "Pistol 25 m:10:{3} Pistol 50 m:5:{2}"},
		{"old range keeps its id", `SELECT name FROM ranges WHERE id = ` + fmt.Sprint(rid), "Pistol 25 m"},
		{"registrations", `SELECT string_agg(format('%s:%s', range_id = @r50, lane), ' ' ORDER BY id) FROM registrations`, "f:4 t:2"},
		{"reservations", `SELECT string_agg(format('%s:%s:%s', range_id = @r50, lanes, booking_item_id IS NOT NULL), ' ' ORDER BY range_id = @r50) FROM reservations`, "f:{9,10}:t t:{1}:t"},
		{"booking item", `SELECT (range_id = @r50)::text FROM booking_items`, "true"},
		{"scorer", `SELECT format('%s:%s', range_id = @r50, lanes) FROM helper_codes`, "t:{1,2}"},
		{"programmes", `SELECT string_agg(format('%s:%s', name, coalesce((range_id = @r50)::text, 'any')), ' ' ORDER BY name) FROM programmes`, "Any:any Precision:true"},
		{"event ranges", `SELECT string_agg(format('%s:%s:%s', event_id, range_id = @r50, position), ' ' ORDER BY event_id, position) FROM event_ranges`, fmt.Sprintf("%d:f:1 %d:t:2 %d:f:1 %d:t:2", open, open, closed, closed)},
		{"print settings", `SELECT count(*)::text FROM print_settings WHERE printer = 'A5'`, "2"},
		{"groups gone", `SELECT (to_regclass('range_groups') IS NULL)::text`, "true"},
	} {
		if got := query(c.sql); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

func TestEntriesBecomeFirstAttempts(t *testing.T) {
	st := testdb.Bare(t)
	ctx := context.Background()
	if err := st.MigrateUntil(ctx, "0029_layouts.sql"); err != nil {
		t.Fatal(err)
	}
	db := st.Pool()
	var tid, ev, shooter, prog, reg int64
	scan := func(dst *int64, sql string, args ...any) {
		t.Helper()
		if err := db.QueryRow(ctx, sql, args...).Scan(dst); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
	}
	scan(&tid, `INSERT INTO tenants (name) VALUES ('Club') RETURNING id`)
	scan(&ev, `INSERT INTO events (tenant_id, name, starts_on, days) VALUES ($1, 'Open', '2026-11-20', '{2026-11-20}') RETURNING id`, tid)
	scan(&shooter, `INSERT INTO shooters (tenant_id, firstname, lastname) VALUES ($1, 'Anna', 'Muster') RETURNING id`, tid)
	scan(&prog, `INSERT INTO programmes (tenant_id, event_id, name) VALUES ($1, $2, 'Main') RETURNING id`, tid, ev)
	scan(&reg, `INSERT INTO registrations (tenant_id, event_id, shooter_id, age_category, status) VALUES ($1, $2, $3, 'elite', 'registered') RETURNING id`, tid, ev, shooter)
	insert := `INSERT INTO entries (tenant_id, registration_id, event_id, shooter_id, programme_id) VALUES ($1, $2, $3, $4, $5)`
	if _, err := db.Exec(ctx, insert, tid, reg, ev, shooter, prog); err != nil {
		t.Fatal(err)
	}
	if err := st.MigrateUntil(ctx, ""); err != nil {
		t.Fatal(err)
	}
	var attempt, limit int
	if err := db.QueryRow(ctx, `SELECT e.attempt, p.attempts FROM entries e JOIN programmes p ON p.id = e.programme_id`).Scan(&attempt, &limit); err != nil {
		t.Fatal(err)
	}
	if attempt != 1 || limit != 1 {
		t.Fatalf("attempt %d, limit %d", attempt, limit)
	}
	if _, err := db.Exec(ctx, insert+` ON CONFLICT DO NOTHING`, tid, reg, ev, shooter, prog); err != nil {
		t.Fatal(err)
	}
	var n int
	db.QueryRow(ctx, `SELECT count(*) FROM entries`).Scan(&n)
	if n != 1 {
		t.Fatalf("a second first attempt must collide, got %d entries", n)
	}
}

func TestSettingsBecomeModules(t *testing.T) {
	st := testdb.Bare(t)
	ctx := context.Background()
	if err := st.MigrateUntil(ctx, "0031_classifications.sql"); err != nil {
		t.Fatal(err)
	}
	db := st.Pool()
	var tid int64
	if err := db.QueryRow(ctx, `INSERT INTO tenants (name) VALUES ('Club') RETURNING id`).Scan(&tid); err != nil {
		t.Fatal(err)
	}
	insert := `INSERT INTO events (tenant_id, name, starts_on, days, settings, archived_at) VALUES ($1, $2, '2026-11-20', '{2026-11-20}', $3, $4)`
	if _, err := db.Exec(ctx, insert, tid, "Plain", `{"showLaneStatus": true}`, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, insert, tid, "Closed", `{"noLaneSheet": true, "noControlStrip": true, "ammoIssue": true}`, "2026-12-01T00:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if err := st.MigrateUntil(ctx, "0032_modules.sql"); err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	rows, _ := db.Query(ctx, `SELECT name, array_to_string(modules, ','), settings::text FROM events`)
	for rows.Next() {
		var name, mods, settings string
		rows.Scan(&name, &mods, &settings)
		got[name] = mods + " " + settings
	}
	want := map[string]string{
		"Plain":  `reservations,booking,attempts,awards,special_prizes,club_competition,classifications,kiosk,helpers,lane_sheet,control_strip {"showLaneStatus": true}`,
		"Closed": `reservations,booking,attempts,awards,special_prizes,club_competition,classifications,kiosk,helpers,ammo {}`,
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s: %s, want %s", k, got[k], v)
		}
	}
}

func TestClubsKeepForms(t *testing.T) {
	st := testdb.Bare(t)
	ctx := context.Background()
	if err := st.MigrateUntil(ctx, "0032_modules.sql"); err != nil {
		t.Fatal(err)
	}
	if _, err := st.Pool().Exec(ctx, `INSERT INTO tenants (name) VALUES ('Club')`); err != nil {
		t.Fatal(err)
	}
	if err := st.MigrateUntil(ctx, ""); err != nil {
		t.Fatal(err)
	}
	var mods string
	st.Pool().QueryRow(ctx, `SELECT array_to_string(modules, ',') FROM tenants`).Scan(&mods)
	if mods != "forms" {
		t.Fatalf("modules %q", mods)
	}
}

func TestBookingJoinsReservations(t *testing.T) {
	st := testdb.Bare(t)
	ctx := context.Background()
	if err := st.MigrateUntil(ctx, "0035_certificates.sql"); err != nil {
		t.Fatal(err)
	}
	db := st.Pool()
	var tid int64
	db.QueryRow(ctx, `INSERT INTO tenants (name) VALUES ('Club') RETURNING id`).Scan(&tid)
	for name, mods := range map[string]string{"Both": "{reservations,booking,kiosk}", "Only booking": "{booking}", "None": "{kiosk}"} {
		if _, err := db.Exec(ctx, `INSERT INTO events (tenant_id, name, starts_on, days, modules) VALUES ($1, $2, '2026-11-20', '{2026-11-20}', $3)`, tid, name, mods); err != nil {
			t.Fatal(err)
		}
	}
	if err := st.MigrateUntil(ctx, "0036_booking_in_reservations.sql"); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"Both": "reservations,kiosk", "Only booking": "reservations", "None": "kiosk"}
	rows, _ := db.Query(ctx, `SELECT name, array_to_string(modules, ',') FROM events`)
	for rows.Next() {
		var name, mods string
		rows.Scan(&name, &mods)
		if want[name] != mods {
			t.Errorf("%s: %s, want %s", name, mods, want[name])
		}
	}
}

func TestEventsKeepLanes(t *testing.T) {
	st := testdb.Bare(t)
	ctx := context.Background()
	if err := st.MigrateUntil(ctx, "0036_booking_in_reservations.sql"); err != nil {
		t.Fatal(err)
	}
	db := st.Pool()
	var tid int64
	db.QueryRow(ctx, `INSERT INTO tenants (name) VALUES ('Club') RETURNING id`).Scan(&tid)
	if _, err := db.Exec(ctx, `INSERT INTO events (tenant_id, name, starts_on, days, modules) VALUES ($1, 'Fest', '2026-11-20', '{2026-11-20}', '{reservations,kiosk}')`, tid); err != nil {
		t.Fatal(err)
	}
	if err := st.MigrateUntil(ctx, ""); err != nil {
		t.Fatal(err)
	}
	var mods string
	db.QueryRow(ctx, `SELECT array_to_string(modules, ',') FROM events`).Scan(&mods)
	if mods != "lanes,reservations,kiosk" {
		t.Fatalf("modules %q", mods)
	}
}
