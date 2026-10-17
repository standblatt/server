package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// EventSettings are per-event options.
type EventSettings struct {
	ExcludedClubID  *int64 `json:"excludedClubId,omitempty" required:"false"`
	ShowLaneStatus  bool   `json:"showLaneStatus" required:"false"`
	PrintAwardValue bool   `json:"printAwardValue" required:"false"`
	AmmoIssue       bool   `json:"ammoIssue" required:"false" doc:"The organizer hands out the ammunition; shows the Ammo button on the dashboard"`
}

// EventRange is a range an event is shot on.
type EventRange struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	LaneCount     int     `json:"laneCount"`
	DisabledLanes []int32 `json:"disabledLanes"`
}

// Event is a competition of a tenant, shot on one or more ranges on one or more days.
type Event struct {
	ID         int64           `json:"id"`
	Ranges     []EventRange    `json:"ranges" doc:"In the order shown; documents print on the first one"`
	Name       string          `json:"name"`
	StartsOn   string          `json:"startsOn" format:"date" doc:"The main day: decides the age categories and sorts the events"`
	Days       []string        `json:"days" doc:"All shooting days including the main day, sorted"`
	Schedule   domain.Schedule `json:"schedule"`
	Organizer  string          `json:"organizer"`
	Timezone   string          `json:"timezone"`
	Locale     string          `json:"locale" enum:"en,de,fr,it" doc:"Language of documents and kiosk displays"`
	Settings   EventSettings   `json:"settings"`
	ArchivedAt *time.Time      `json:"archivedAt,omitempty" doc:"Set when the event is archived: everything stays readable, nothing can be changed"`
}

// EventInput holds the editable fields of an event.
type EventInput struct {
	RangeIDs  []int64         `json:"rangeIds" minItems:"1"`
	Name      string          `json:"name" minLength:"1" maxLength:"200"`
	StartsOn  string          `json:"startsOn" format:"date"`
	Schedule  domain.Schedule `json:"schedule" required:"false" doc:"Single days and/or a period; the main day is always a shooting day"`
	Days      []string        `json:"-"`
	Organizer string          `json:"organizer" maxLength:"200" required:"false"`
	Timezone  string          `json:"timezone" required:"false"`
	Locale    string          `json:"locale" required:"false" enum:"en,de,fr,it"`
	Settings  EventSettings   `json:"settings" required:"false"`
}

const eventCols = `e.id,
	coalesce((SELECT json_agg(json_build_object('id', rg.id, 'name', rg.name, 'laneCount', rg.lane_count, 'disabledLanes', rg.disabled_lanes) ORDER BY er.position)
		FROM event_ranges er JOIN ranges rg ON rg.id = er.range_id WHERE er.event_id = e.id), '[]'),
	e.name, e.starts_on::text, e.days::text[], e.schedule, e.organizer, e.timezone, e.locale, e.settings, e.archived_at`

// ErrRangeUsed is returned when a range with registrations or programmes is removed from an event.
var ErrRangeUsed = errors.New("store: range has registrations or programmes")

// ListEvents lists the events of a tenant, newest first.
func (s *Store) ListEvents(ctx context.Context, tenantID int64) ([]Event, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+eventCols+` FROM events e
		WHERE e.tenant_id = $1 ORDER BY e.starts_on DESC, e.id DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Event])
}

// GetEvent returns one event of a tenant.
func (s *Store) GetEvent(ctx context.Context, tenantID, eventID int64) (Event, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+eventCols+` FROM events e WHERE e.tenant_id = $1 AND e.id = $2`, tenantID, eventID)
	if err != nil {
		return Event{}, err
	}
	ev, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[Event])
	return ev, mapErr(err)
}

// CreateEvent adds an event; every range must belong to the same tenant (ErrInUse otherwise).
func (s *Store) CreateEvent(ctx context.Context, tenantID int64, in EventInput) (int64, error) {
	var id int64
	in.Schedule = domain.NormalizeSchedule(in.Schedule, in.StartsOn)
	in.Days = in.Schedule.Dates()
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `INSERT INTO events (tenant_id, name, starts_on, days, schedule, organizer, timezone, locale, settings)
			VALUES ($1, $2, $3, $4::date[], $5, $6, $7, $8, $9) RETURNING id`,
			tenantID, in.Name, in.StartsOn, in.Days, in.Schedule, in.Organizer, in.Timezone, in.Locale, in.Settings).Scan(&id); err != nil {
			return err
		}
		return setEventRanges(ctx, tx, tenantID, id, in.RangeIDs)
	})
	return id, mapErr(err)
}

// UpdateEvent changes an event of a tenant. ErrInUse when a range is unknown, ErrRangeUsed when a range
// with registrations or programmes is dropped.
func (s *Store) UpdateEvent(ctx context.Context, tenantID, eventID int64, in EventInput) error {
	in.Schedule = domain.NormalizeSchedule(in.Schedule, in.StartsOn)
	in.Days = in.Schedule.Dates()
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE events SET name = $3, starts_on = $4, days = $5::date[], schedule = $6, organizer = $7,
			timezone = $8, locale = $9, settings = $10 WHERE tenant_id = $1 AND id = $2`,
			tenantID, eventID, in.Name, in.StartsOn, in.Days, in.Schedule, in.Organizer, in.Timezone, in.Locale, in.Settings)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		var used bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM registrations WHERE event_id = $1 AND range_id <> ALL($2))
			OR EXISTS (SELECT 1 FROM programmes WHERE event_id = $1 AND range_id <> ALL($2))`, eventID, in.RangeIDs).Scan(&used); err != nil {
			return err
		}
		if used {
			return ErrRangeUsed
		}
		if _, err := tx.Exec(ctx, `DELETE FROM event_ranges WHERE event_id = $1 AND range_id <> ALL($2)`, eventID, in.RangeIDs); err != nil {
			return err
		}
		return setEventRanges(ctx, tx, tenantID, eventID, in.RangeIDs)
	})
	return mapErr(err)
}

// setEventRanges upserts the ranges of an event in the given order; ErrInUse when one is not a range of the tenant.
func setEventRanges(ctx context.Context, tx pgx.Tx, tenantID, eventID int64, rangeIDs []int64) error {
	tag, err := tx.Exec(ctx, `INSERT INTO event_ranges (tenant_id, event_id, range_id, position)
		SELECT $1, $2, x.id, x.n FROM unnest($3::bigint[]) WITH ORDINALITY AS x(id, n)
		JOIN ranges r ON r.id = x.id AND r.tenant_id = $1
		ON CONFLICT (event_id, range_id) DO UPDATE SET position = EXCLUDED.position`, tenantID, eventID, rangeIDs)
	if err != nil {
		return err
	}
	if int(tag.RowsAffected()) != len(rangeIDs) {
		return ErrInUse
	}
	return nil
}

// SetEventArchived archives or unarchives an event of a tenant; doing it twice is not an error.
func (s *Store) SetEventArchived(ctx context.Context, tenantID, eventID int64, archived bool) error {
	tag, err := s.pool.Exec(ctx, `UPDATE events SET archived_at = CASE WHEN $3 THEN now() END
		WHERE tenant_id = $1 AND id = $2 AND (archived_at IS NOT NULL) <> $3`, tenantID, eventID, archived)
	if err != nil || tag.RowsAffected() > 0 {
		return mapErr(err)
	}
	_, err = s.GetEvent(ctx, tenantID, eventID)
	return err
}

// DeleteEvent removes an event of a tenant; ErrInUse if it has registrations, ErrArchived if it is archived.
func (s *Store) DeleteEvent(ctx context.Context, tenantID, eventID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM events WHERE tenant_id = $1 AND id = $2`, tenantID, eventID))
}
