package store

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/modules"
)

// eventUsage tells per event module when an event has data of it (SQL over the alias e for events); a module with
// data stays on and cannot be switched off. "" means the module never holds data.
var eventUsage = map[string]string{
	modules.Reservations:    `e.booking_code IS NOT NULL`,
	modules.LaneSheet:       "",
	modules.ControlStrip:    "",
	modules.Ammo:            "",
	modules.Attempts:        `EXISTS (SELECT 1 FROM programmes x WHERE x.event_id = e.id AND x.attempts <> 1)`,
	modules.Awards:          `EXISTS (SELECT 1 FROM award_tiers t JOIN programmes x ON x.id = t.programme_id WHERE x.event_id = e.id)`,
	modules.SpecialPrizes:   `EXISTS (SELECT 1 FROM programmes x WHERE x.event_id = e.id AND x.prize_ranks > 0)`,
	modules.ClubCompetition: "",
	modules.Classifications: `EXISTS (SELECT 1 FROM classifications x WHERE x.event_id = e.id)`,
	modules.Kiosk:           `EXISTS (SELECT 1 FROM kiosk_tokens x WHERE x.event_id = e.id)`,
	modules.Helpers:         `EXISTS (SELECT 1 FROM helper_codes x WHERE x.event_id = e.id)`,
}

// clubUsage is eventUsage for club modules, SQL over the alias t for tenants.
var clubUsage = map[string]string{
	modules.Forms: `EXISTS (SELECT 1 FROM layouts x WHERE x.tenant_id = t.id) OR EXISTS (SELECT 1 FROM images x WHERE x.tenant_id = t.id)`,
}

// copiers copy the data of event modules that is not part of the programmes when an event is copied.
var copiers = []func(ctx context.Context, tx pgx.Tx, tenantID, fromEventID, toEventID int64) error{
	copyClassifications,
}

// usedSQL builds the SQL array of the modules of a level that hold data from their usage conditions.
func usedSQL(level modules.Level, usage map[string]string) string {
	parts := []string{"NULL"}
	for _, k := range modules.Keys(level) {
		if cond := usage[k]; cond != "" {
			parts = append(parts, "CASE WHEN "+cond+" THEN '"+k+"' END")
		}
	}
	return "array_remove(ARRAY[" + strings.Join(parts, ", ") + "]::text[], NULL)"
}

// eventUsed and clubUsed are the SQL arrays of the modules with data, for the aliases e and t.
var (
	eventUsed = usedSQL(modules.Event, eventUsage)
	clubUsed  = usedSQL(modules.Club, clubUsage)
)

// ModuleUsedError is returned when modules with data are switched off.
type ModuleUsedError struct{ Keys []string }

// Error lists the modules that are in use.
func (e *ModuleUsedError) Error() string {
	return "store: modules in use: " + strings.Join(e.Keys, ", ")
}

// EventOf returns the event of a programme, registration or entry of a tenant; param is the path parameter name
// (programmeId, registrationId or entryId).
func (s *Store) EventOf(ctx context.Context, tenantID int64, param string, id int64) (int64, error) {
	table := map[string]string{"programmeId": "programmes", "registrationId": "registrations", "entryId": "entries"}[param]
	if table == "" {
		return 0, ErrNotFound
	}
	var eventID int64
	err := s.pool.QueryRow(ctx, `SELECT event_id FROM `+table+` WHERE tenant_id = $1 AND id = $2`, tenantID, id).Scan(&eventID)
	return eventID, mapErr(err)
}

// ClubModules returns the modules that are on for a club and those that hold data.
func (s *Store) ClubModules(ctx context.Context, tenantID int64) (on, used []string, err error) {
	var chosen []string
	err = s.pool.QueryRow(ctx, `SELECT t.modules, `+clubUsed+` FROM tenants t WHERE t.id = $1`, tenantID).Scan(&chosen, &used)
	if err != nil {
		return nil, nil, mapErr(err)
	}
	return modules.Effective(modules.Club, chosen, used, true), used, nil
}

// SetClubModules chooses the modules of a club; ModuleUsedError when modules with data would be switched off.
func (s *Store) SetClubModules(ctx context.Context, tenantID int64, keys []string) error {
	chosen := modules.Effective(modules.Club, keys, nil, true)
	return mapErr(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var used []string
		if err := tx.QueryRow(ctx, `SELECT `+clubUsed+` FROM tenants t WHERE t.id = $1 FOR UPDATE`, tenantID).Scan(&used); err != nil {
			return err
		}
		if off := modules.Missing(used, chosen); len(off) > 0 {
			return &ModuleUsedError{Keys: off}
		}
		_, err := tx.Exec(ctx, `UPDATE tenants SET modules = $2 WHERE id = $1`, tenantID, chosen)
		return err
	}))
}
