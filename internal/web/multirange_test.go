package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func TestMultiRangeEvent(t *testing.T) {
	var srv *web.Server
	e := newEnvWith(t, func(s *web.Server) { srv = s })
	f := e.registrationFixture("Fest")
	e.user("boss", f.tenant, store.RoleAdmin)
	boss := e.login("boss")
	ctx := context.Background()
	r50, _ := e.st.CreateRange(ctx, f.tenant, store.RangeInput{Name: "50m", LaneCount: 3})
	ev, _ := e.st.GetEvent(ctx, f.tenant, f.event)
	r10 := ev.Ranges[0].ID
	evPath := fmt.Sprintf("/api/tenants/%d/events/%d", f.tenant, f.event)
	lane := func(n int) *int { return &n }

	code, b := e.do(boss, "PUT", evPath, map[string]any{"rangeIds": []int64{r10, r50.ID}, "name": ev.Name, "startsOn": "2026-11-21",
		"schedule": domain.Schedule{Days: []domain.ScheduleDay{{Date: "2026-11-22"}, {Date: "2026-11-18"}}}})
	expect(t, "add a range and days", code, 204, b)
	code, b = e.do(boss, "PUT", evPath, map[string]any{"rangeIds": []int64{r10}, "name": ev.Name, "startsOn": "2026-11-21",
		"schedule": domain.Schedule{From: "2026-12-01", To: "2026-11-01"}})
	if code != 422 || !strings.Contains(string(b), "ends before it starts") {
		t.Fatalf("period backwards: %d %s", code, b)
	}
	code, b = e.do(boss, "GET", evPath, nil)
	expect(t, "get", code, 200, b)
	json.Unmarshal(b, &ev)
	if len(ev.Ranges) != 2 || ev.Ranges[1].Name != "50m" || strings.Join(ev.Days, ",") != "2026-11-18,2026-11-21,2026-11-22" {
		t.Fatalf("event with two ranges: %+v", ev)
	}

	code, b = e.do(boss, "PUT", fmt.Sprintf("/api/tenants/%d/programmes/%d/range", f.tenant, f.junior), map[string]any{"rangeId": r50.ID})
	expect(t, "tie junior to 50m", code, 204, b)
	other := e.tenant("Other")
	foreign, _ := e.st.CreateRange(ctx, other, store.RangeInput{Name: "25m", LaneCount: 2})
	code, b = e.do(boss, "PUT", fmt.Sprintf("/api/tenants/%d/programmes/%d/range", f.tenant, f.junior), map[string]any{"rangeId": foreign.ID})
	expect(t, "range of another tenant", code, 404, b)

	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.youth, Lane: lane(1), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}})
	expect(t, "lane without range on a two-range event", code, 422, b)
	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.youth, RangeID: &r10, Lane: lane(1), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}})
	if code != 422 || !strings.Contains(string(b), "not shot on this range") {
		t.Fatalf("junior on 10m: %d %s", code, b)
	}
	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.youth, RangeID: &r50.ID, Lane: lane(1), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}})
	expect(t, "junior on 50m lane 1", code, 201, b)
	var junior store.Registration
	json.Unmarshal(b, &junior)
	form, err := e.st.CreateLayout(ctx, f.tenant, store.LayoutInput{Name: "Sheet", Layout: domain.Layout{Width: 100, Height: 100,
		Elements: []domain.Element{{Type: "text", X: 10, Y: 10, Value: "{{shooter.name}}"}}}})
	if err != nil {
		t.Fatal(err)
	}
	e.st.SetProgrammeLayout(ctx, f.tenant, f.junior, &form)
	if id, err := srv.FormRangeForTest(ctx, f.tenant, junior.Entries[0].ID); err != nil || id != r50.ID {
		t.Fatalf("print form of a 50m entry goes to range %d (%v), want %d", id, err, r50.ID)
	}
	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.shooter, ClubID: &f.club, RangeID: &r10, Lane: lane(1), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.main}}})
	expect(t, "main on 10m lane 1 (same lane number, other range)", code, 201, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)
	if reg.RangeID == nil || *reg.RangeID != r10 {
		t.Fatalf("registration range: %+v", reg)
	}

	code, b = e.do(boss, "GET", evPath+"/lanes", nil)
	expect(t, "get", code, 200, b)
	var lanes []web.RangeLanes
	json.Unmarshal(b, &lanes)
	if len(lanes) != 2 || lanes[0].RangeName != "10m" || lanes[0].Lanes[0].State != "occupied" || len(lanes[1].Lanes) != 3 || lanes[1].Lanes[0].ShooterName != "Young Yves" {
		t.Fatalf("lanes per range: %+v", lanes)
	}

	code, b = e.do(boss, "PUT", evPath, map[string]any{"rangeIds": []int64{r10}, "name": ev.Name, "startsOn": "2026-11-21"})
	expect(t, "drop a range in use", code, 409, b)

	code, b = e.do(boss, "POST", evPath+"/copy", map[string]any{"name": "Fest 2027", "startsOn": "2027-11-20"})
	expect(t, "copy", code, 201, b)
	var copied store.Event
	json.Unmarshal(b, &copied)
	if len(copied.Ranges) != 2 || strings.Join(copied.Days, ",") != "2027-11-17,2027-11-20,2027-11-21" {
		t.Fatalf("copied event: %+v", copied)
	}
	progs, _ := e.st.ListProgrammes(ctx, f.tenant, copied.ID)
	for _, p := range progs {
		if p.Name == "Junior" && (p.RangeID == nil || *p.RangeID != r50.ID) {
			t.Fatalf("copied programme keeps its range: %+v", p)
		}
	}

	code, b = e.do(boss, "GET", evPath+"/statistics", nil)
	if code != 200 || !strings.Contains(string(b), `"registrations":2`) {
		t.Fatalf("stats per day: %d %s", code, b)
	}
}

func TestEventWithoutRange(t *testing.T) {
	e := newEnv(t)
	f := e.registrationFixture("Field")
	e.user("boss", f.tenant, store.RoleAdmin)
	boss := e.login("boss")
	base := fmt.Sprintf("/api/tenants/%d", f.tenant)
	evPath := fmt.Sprintf("%s/events/%d", base, f.event)

	code, b := e.do(boss, "POST", base+"/events", map[string]any{"name": "Annual programme", "startsOn": "2027-03-01"})
	expect(t, "create an event without a range", code, 201, b)
	var annual store.Event
	json.Unmarshal(b, &annual)
	code, b = e.do(boss, "PUT", evPath, map[string]any{"name": "Indoor", "startsOn": "2026-11-20", "modules": []string{"control_strip", "reservations"}})
	expect(t, "drop the range", code, 204, b)
	var ev store.Event
	_, b = e.do(boss, "GET", evPath, nil)
	json.Unmarshal(b, &ev)
	if len(ev.Ranges) != 0 || fmt.Sprint(ev.Modules) != "[control_strip]" {
		t.Fatalf("event without a range: %s", b)
	}
	code, b = e.do(boss, "GET", evPath+"/lanes", nil)
	if code != 200 || strings.TrimSpace(string(b)) != "[]" {
		t.Fatalf("lanes of an event without a range: %d %s", code, b)
	}
	lane := 1
	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.shooter, Lane: &lane, Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.main}}})
	expect(t, "lane on an event without a range", code, 422, b)
	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.shooter, Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.main}}})
	expect(t, "register without a lane", code, 201, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)
	code, b = e.do(boss, "GET", fmt.Sprintf("%s/registrations/%d/documents/lane_sheet", base, reg.ID), nil)
	expect(t, "lane sheet switched off", code, 404, b)
	code, _ = e.do(boss, "GET", fmt.Sprintf("%s/registrations/%d/documents/control_strip", base, reg.ID), nil)
	expect(t, "control strip", code, 200, nil)

	code, b = e.do(boss, "DELETE", evPath, nil)
	expect(t, "delete an event with registrations", code, 409, b)
	annualPath := fmt.Sprintf("%s/events/%d", base, annual.ID)
	code, b = e.do(boss, "DELETE", annualPath, nil)
	expect(t, "delete an event", code, 204, b)
	code, b = e.do(boss, "GET", annualPath, nil)
	expect(t, "deleted event", code, 404, b)
	ranges, _ := e.st.ListRanges(context.Background(), f.tenant)
	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/ranges/%d", base, ranges[0].ID), nil)
	expect(t, "delete a range no event uses", code, 204, b)
	code, b = e.do(boss, "GET", base+"/ranges", nil)
	if code != 200 || strings.TrimSpace(string(b)) != "[]" {
		t.Fatalf("ranges after the delete: %d %s", code, b)
	}
}
