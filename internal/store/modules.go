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
	modules.Reservations:    "",
	modules.Booking:         `e.booking_code IS NOT NULL`,
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

// copiers copy the data of event modules that is not part of the programmes when an event is copied.
var copiers = []func(ctx context.Context, tx pgx.Tx, tenantID, fromEventID, toEventID int64) error{
	copyClassifications,
}

// eventUsed is the SQL array of the event modules with data, for the alias e.
var eventUsed = func() string {
	var parts []string
	for _, k := range modules.Keys(modules.Event) {
		if cond := eventUsage[k]; cond != "" {
			parts = append(parts, "CASE WHEN "+cond+" THEN '"+k+"' END")
		}
	}
	return "array_remove(ARRAY[" + strings.Join(parts, ", ") + "]::text[], NULL)"
}()

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
