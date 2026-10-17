package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// ErrOverlap means a reservation would overlap another one on a lane.
var ErrOverlap = errors.New("store: reservation overlaps another one")

// Reservation blocks lanes of a range for a time. It shows as reserved from LeadMinutes before the start.
type Reservation struct {
	ID          int64     `json:"id"`
	RangeID     int64     `json:"rangeId"`
	Lanes       []int     `json:"lanes"`
	StartsAt    time.Time `json:"startsAt"`
	EndsAt      time.Time `json:"endsAt"`
	LeadMinutes int       `json:"leadMinutes"`
	Description string    `json:"description"`
}

const reservationCols = `id, range_id, lanes, starts_at, ends_at, lead_minutes, description`

// ListReservations lists the reservations of a range ordered by start; past ones only with all.
func (s *Store) ListReservations(ctx context.Context, tenantID, rangeID int64, all bool) ([]Reservation, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+reservationCols+` FROM reservations
		WHERE tenant_id = $1 AND range_id = $2 AND ($3 OR ends_at > now()) ORDER BY starts_at, id`, tenantID, rangeID, all)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Reservation])
}

// ActiveReservations lists the reservations of a range that are shown as reserved at time t.
func (s *Store) ActiveReservations(ctx context.Context, rangeID int64, t time.Time) ([]Reservation, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+reservationCols+` FROM reservations
		WHERE range_id = $1 AND starts_at - make_interval(mins => lead_minutes) <= $2 AND ends_at > $2
		ORDER BY starts_at`, rangeID, t)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Reservation])
}

// GetReservation returns one reservation of a tenant.
func (s *Store) GetReservation(ctx context.Context, tenantID, id int64) (Reservation, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+reservationCols+` FROM reservations WHERE tenant_id = $1 AND id = $2`, tenantID, id)
	if err != nil {
		return Reservation{}, err
	}
	r, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[Reservation])
	return r, mapErr(err)
}

// SaveReservation inserts (ID 0) or updates a reservation. The range row is locked so that two overlapping
// reservations cannot be stored at the same time; ErrOverlap names none, the caller lists them with Overlaps.
func (s *Store) SaveReservation(ctx context.Context, tenantID int64, r Reservation) (int64, error) {
	id := r.ID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var locked int64
		if err := tx.QueryRow(ctx, `SELECT id FROM ranges WHERE tenant_id = $1 AND id = $2 FOR UPDATE`, tenantID, r.RangeID).Scan(&locked); err != nil {
			return err
		}
		var overlap bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM reservations
			WHERE range_id = $1 AND id <> $2 AND lanes && $3 AND starts_at < $5 AND ends_at > $4)`,
			r.RangeID, r.ID, r.Lanes, r.StartsAt, r.EndsAt).Scan(&overlap); err != nil {
			return err
		}
		if overlap {
			return ErrOverlap
		}
		if r.ID == 0 {
			return tx.QueryRow(ctx, `INSERT INTO reservations (tenant_id, range_id, lanes, starts_at, ends_at, lead_minutes, description)
				VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
				tenantID, r.RangeID, r.Lanes, r.StartsAt, r.EndsAt, r.LeadMinutes, r.Description).Scan(&id)
		}
		tag, err := tx.Exec(ctx, `UPDATE reservations SET lanes = $3, starts_at = $4, ends_at = $5, lead_minutes = $6, description = $7
			WHERE tenant_id = $1 AND id = $2`, tenantID, r.ID, r.Lanes, r.StartsAt, r.EndsAt, r.LeadMinutes, r.Description)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	if errors.Is(err, ErrOverlap) {
		return 0, err
	}
	return id, mapErr(err)
}

// Overlaps lists the reservations that collide with r on a lane.
func (s *Store) Overlaps(ctx context.Context, r Reservation) ([]Reservation, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+reservationCols+` FROM reservations
		WHERE range_id = $1 AND id <> $2 AND lanes && $3 AND starts_at < $5 AND ends_at > $4 ORDER BY starts_at`,
		r.RangeID, r.ID, r.Lanes, r.StartsAt, r.EndsAt)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Reservation])
}

// DeleteReservation removes a reservation and returns its range.
func (s *Store) DeleteReservation(ctx context.Context, tenantID, id int64) (int64, error) {
	var rangeID int64
	err := s.pool.QueryRow(ctx, `DELETE FROM reservations WHERE tenant_id = $1 AND id = $2 RETURNING range_id`, tenantID, id).Scan(&rangeID)
	return rangeID, mapErr(err)
}

// EventRef is an event to notify.
type EventRef struct {
	TenantID int64
	EventID  int64
}

// EventsOfRange lists the events of a range that take place around now; reservation changes are announced to them.
func (s *Store) EventsOfRange(ctx context.Context, rangeID int64) ([]EventRef, error) {
	rows, err := s.pool.Query(ctx, `SELECT e.tenant_id, e.id FROM events e JOIN event_ranges er ON er.event_id = e.id
		WHERE er.range_id = $1 AND EXISTS (SELECT 1 FROM unnest(e.days) d WHERE d BETWEEN current_date - 7 AND current_date + 7)`, rangeID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[EventRef])
}

const reservationTickKey = "reservation_tick"

// ReservationTick finds the ranges whose reservations started showing or ended since the previous tick and
// remembers now as the new tick. Only one replica runs it at a time (transaction advisory lock); others get nil.
func (s *Store) ReservationTick(ctx context.Context, now time.Time) ([]int64, error) {
	var ranges []int64
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var got bool
		if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock(7346264)`).Scan(&got); err != nil || !got {
			return err
		}
		var last time.Time
		err := tx.QueryRow(ctx, `SELECT value::timestamptz FROM settings WHERE key = $1`, reservationTickKey).Scan(&last)
		if errors.Is(err, pgx.ErrNoRows) {
			last = now.Add(-time.Minute)
		} else if err != nil {
			return err
		}
		rows, err := tx.Query(ctx, `SELECT DISTINCT range_id FROM reservations
			WHERE starts_at - make_interval(mins => lead_minutes) BETWEEN $1 AND $2 OR ends_at BETWEEN $1 AND $2`, last, now)
		if err != nil {
			return err
		}
		if ranges, err = pgx.CollectRows(rows, pgx.RowTo[int64]); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO settings (key, value) VALUES ($1, $2)
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = now()`, reservationTickKey, now.Format(time.RFC3339Nano))
		return err
	})
	return ranges, err
}
