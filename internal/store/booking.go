package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// BookingStatus is the state of a booking request.
type BookingStatus string

// Booking request states.
const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingDeclined  BookingStatus = "declined"
	BookingCancelled BookingStatus = "cancelled"
)

// BookingRequest is what an external club asked for: lanes in slots and a list of shooters.
type BookingRequest struct {
	ID            int64                `json:"id"`
	EventID       int64                `json:"eventId"`
	Club          string               `json:"club"`
	ClubID        *int64               `json:"clubId,omitempty" doc:"Club of the master data the request was matched to on confirmation"`
	Contact       string               `json:"contact"`
	Email         string               `json:"email"`
	Phone         string               `json:"phone"`
	Note          string               `json:"note"`
	Status        BookingStatus        `json:"status" enum:"pending,confirmed,declined,cancelled"`
	DeclineReason string               `json:"declineReason"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
	DecidedAt     *time.Time           `json:"decidedAt,omitempty"`
	Items         []BookingItemView    `json:"items"`
	Shooters      []BookingShooterView `json:"shooters"`
}

// BookingItemView is an item of a request with the lanes it got on confirmation.
type BookingItemView struct {
	domain.BookingItem
	Assigned []int `json:"assigned" doc:"Lanes reserved for it; empty until confirmed"`
}

// BookingShooterView is an announced shooter with its link to the master data.
type BookingShooterView struct {
	domain.BookingShooter
	ShooterID *int64 `json:"shooterId,omitempty" doc:"Set once the shooter was taken into the master data"`
	Known     bool   `json:"known" doc:"A shooter with this licence, or this name and birthday, is in the master data"`
}

// BookingRequestInput is the editable part of a request.
type BookingRequestInput struct {
	Club     string                  `json:"club" minLength:"1" maxLength:"200"`
	Contact  string                  `json:"contact" minLength:"1" maxLength:"200"`
	Email    string                  `json:"email" format:"email" maxLength:"200"`
	Phone    string                  `json:"phone" maxLength:"50" required:"false"`
	Note     string                  `json:"note" maxLength:"2000" required:"false"`
	Items    []domain.BookingItem    `json:"items" maxItems:"500"`
	Shooters []domain.BookingShooter `json:"shooters" required:"false" maxItems:"200"`
}

// BookingTarget is the event behind a public booking code.
type BookingTarget struct {
	TenantID int64
	Event    Event
	Settings domain.BookingSettings
}

// FullSlotError tells which item of a request found too few free lanes on confirmation.
type FullSlotError struct {
	Item domain.BookingItem
	Free int
}

// Error describes the full slot.
func (e *FullSlotError) Error() string {
	return fmt.Sprintf("store: only %d lanes free on range %d at %s", e.Free, e.Item.RangeID, e.Item.Start)
}

// ErrNotPending means the request was already decided or cancelled.
var ErrNotPending = errors.New("store: booking request is not pending")

// BookingConfig returns the booking code (empty when off) and settings of an event.
func (s *Store) BookingConfig(ctx context.Context, tenantID, eventID int64) (string, domain.BookingSettings, error) {
	var code string
	var set domain.BookingSettings
	err := s.pool.QueryRow(ctx, `SELECT coalesce(booking_code, ''), booking FROM events WHERE tenant_id = $1 AND id = $2`,
		tenantID, eventID).Scan(&code, &set)
	return code, set, mapErr(err)
}

// SetBookingSettings stores the booking settings of an event.
func (s *Store) SetBookingSettings(ctx context.Context, tenantID, eventID int64, set domain.BookingSettings) error {
	return affected(s.pool.Exec(ctx, `UPDATE events SET booking = $3 WHERE tenant_id = $1 AND id = $2`, tenantID, eventID, set))
}

// SetBookingCode switches the online booking of an event on with a new code, or off with an empty one.
func (s *Store) SetBookingCode(ctx context.Context, tenantID, eventID int64, code string) error {
	return affected(s.pool.Exec(ctx, `UPDATE events SET booking_code = nullif($3, '') WHERE tenant_id = $1 AND id = $2`,
		tenantID, eventID, code))
}

// BookingByCode returns the open (not archived) event behind a public booking code.
func (s *Store) BookingByCode(ctx context.Context, code string) (BookingTarget, error) {
	var t BookingTarget
	var eventID int64
	err := s.pool.QueryRow(ctx, `SELECT tenant_id, id, booking FROM events WHERE booking_code = $1 AND archived_at IS NULL`, code).
		Scan(&t.TenantID, &eventID, &t.Settings)
	if err != nil {
		return t, mapErr(err)
	}
	t.Event, err = s.GetEvent(ctx, t.TenantID, eventID)
	return t, err
}

// ReservationsBetween lists the reservations of the given ranges of a tenant that overlap [from, to), leaving out
// those of the booking request exceptRequest (0 = none), which a change of that request would free.
func (s *Store) ReservationsBetween(ctx context.Context, tenantID int64, rangeIDs []int64, from, to time.Time, exceptRequest int64) ([]Reservation, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+reservationCols+` FROM reservations
		WHERE tenant_id = $1 AND range_id = ANY($2) AND starts_at < $4 AND ends_at > $3
			AND (booking_item_id IS NULL OR booking_item_id NOT IN (SELECT id FROM booking_items WHERE request_id = $5))
		ORDER BY starts_at`, tenantID, rangeIDs, from, to, exceptRequest)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Reservation])
}

// CreateBookingRequest stores a new pending request of an event with its items and shooters.
func (s *Store) CreateBookingRequest(ctx context.Context, tenantID, eventID int64, in BookingRequestInput, manageHash string) (int64, error) {
	var id int64
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `INSERT INTO booking_requests (tenant_id, event_id, club, contact, email, phone, note, manage_hash)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			tenantID, eventID, in.Club, in.Contact, in.Email, in.Phone, in.Note, manageHash).Scan(&id); err != nil {
			return err
		}
		if err := insertBookingItems(ctx, tx, tenantID, id, in.Items); err != nil {
			return err
		}
		return syncBookingShooters(ctx, tx, tenantID, id, in.Shooters)
	})
	return id, mapErr(err)
}

// insertBookingItems adds the items of a request.
func insertBookingItems(ctx context.Context, tx pgx.Tx, tenantID, requestID int64, items []domain.BookingItem) error {
	for _, it := range items {
		if _, err := tx.Exec(ctx, `INSERT INTO booking_items (tenant_id, request_id, range_id, starts_at, ends_at, lanes)
			VALUES ($1, $2, $3, $4, $5, $6)`, tenantID, requestID, it.RangeID, it.Start, it.End, it.Lanes); err != nil {
			return err
		}
	}
	return nil
}

// shooterKey identifies an announced shooter for diffing two lists.
func shooterKey(b domain.BookingShooter) string {
	return strings.ToLower(strings.TrimSpace(b.Firstname)) + "\x00" + strings.ToLower(strings.TrimSpace(b.Lastname)) + "\x00" +
		b.Birthday + "\x00" + domain.NormalizeLicense(b.LicenseNumber)
}

// syncBookingShooters makes the shooter list of a request equal to list: unchanged rows keep their master data
// link, removed ones are deleted, new ones are added without a link.
func syncBookingShooters(ctx context.Context, tx pgx.Tx, tenantID, requestID int64, list []domain.BookingShooter) error {
	rows, err := tx.Query(ctx, `SELECT id, firstname, lastname, birthday::text, license_number FROM booking_shooters
		WHERE tenant_id = $1 AND request_id = $2`, tenantID, requestID)
	if err != nil {
		return err
	}
	type row struct {
		id int64
		domain.BookingShooter
	}
	have, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (row, error) {
		var x row
		err := r.Scan(&x.id, &x.Firstname, &x.Lastname, &x.Birthday, &x.LicenseNumber)
		return x, err
	})
	if err != nil {
		return err
	}
	keep := map[string][]int64{}
	for _, h := range have {
		keep[shooterKey(h.BookingShooter)] = append(keep[shooterKey(h.BookingShooter)], h.id)
	}
	for _, b := range list {
		k := shooterKey(b)
		if ids := keep[k]; len(ids) > 0 {
			keep[k] = ids[1:]
			continue
		}
		if _, err := tx.Exec(ctx, `INSERT INTO booking_shooters (tenant_id, request_id, firstname, lastname, birthday, license_number)
			VALUES ($1, $2, $3, $4, $5, $6)`, tenantID, requestID, strings.TrimSpace(b.Firstname), strings.TrimSpace(b.Lastname),
			b.Birthday, domain.NormalizeLicense(b.LicenseNumber)); err != nil {
			return err
		}
	}
	var drop []int64
	for _, ids := range keep {
		drop = append(drop, ids...)
	}
	_, err = tx.Exec(ctx, `DELETE FROM booking_shooters WHERE tenant_id = $1 AND id = ANY($2)`, tenantID, drop)
	return err
}

const bookingRequestCols = `r.id, r.event_id, r.club, r.club_id, r.contact, r.email, r.phone, r.note, r.status, r.decline_reason,
	r.created_at, r.updated_at, r.decided_at,
	coalesce((SELECT json_agg(json_build_object('rangeId', i.range_id, 'start', i.starts_at, 'end', i.ends_at, 'lanes', i.lanes,
		'assigned', coalesce((SELECT array_agg(l ORDER BY l) FROM reservations rv, unnest(rv.lanes) l WHERE rv.booking_item_id = i.id), '{}'))
		ORDER BY i.starts_at, i.range_id) FROM booking_items i WHERE i.request_id = r.id), '[]'),
	coalesce((SELECT json_agg(json_build_object('firstname', b.firstname, 'lastname', b.lastname, 'birthday', b.birthday,
		'licenseNumber', b.license_number, 'shooterId', b.shooter_id,
		'known', b.shooter_id IS NOT NULL OR EXISTS (SELECT 1 FROM shooters s WHERE s.tenant_id = b.tenant_id AND (
			(b.license_number <> '' AND ltrim(s.license_number, '0') = ltrim(b.license_number, '0'))
			OR (lower(s.firstname) = lower(b.firstname) AND lower(s.lastname) = lower(b.lastname) AND s.birthday = b.birthday))))
		ORDER BY b.lastname, b.firstname, b.id) FROM booking_shooters b WHERE b.request_id = r.id), '[]')`

// ListBookingRequests lists the requests of an event, newest first.
func (s *Store) ListBookingRequests(ctx context.Context, tenantID, eventID int64) ([]BookingRequest, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+bookingRequestCols+` FROM booking_requests r
		WHERE r.tenant_id = $1 AND r.event_id = $2 ORDER BY r.created_at DESC, r.id DESC`, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[BookingRequest])
}

// GetBookingRequest returns one request of a tenant.
func (s *Store) GetBookingRequest(ctx context.Context, tenantID, id int64) (BookingRequest, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+bookingRequestCols+` FROM booking_requests r WHERE r.tenant_id = $1 AND r.id = $2`, tenantID, id)
	if err != nil {
		return BookingRequest{}, err
	}
	r, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[BookingRequest])
	return r, mapErr(err)
}

// BookingRequestByHash returns the request behind a manage link with its tenant.
func (s *Store) BookingRequestByHash(ctx context.Context, hash string) (int64, BookingRequest, error) {
	var tenantID, id int64
	if err := s.pool.QueryRow(ctx, `SELECT tenant_id, id FROM booking_requests WHERE manage_hash = $1`, hash).Scan(&tenantID, &id); err != nil {
		return 0, BookingRequest{}, mapErr(err)
	}
	r, err := s.GetBookingRequest(ctx, tenantID, id)
	return tenantID, r, err
}

// UpdateBookingRequest replaces contact data, items and shooters of a pending or confirmed request. When the items
// change, their reservations go with them and a confirmed request is pending again; it returns whether that happened.
func (s *Store) UpdateBookingRequest(ctx context.Context, tenantID, id int64, in BookingRequestInput) (bool, error) {
	reopened := false
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var status BookingStatus
		if err := tx.QueryRow(ctx, `SELECT status FROM booking_requests WHERE tenant_id = $1 AND id = $2 FOR UPDATE`, tenantID, id).Scan(&status); err != nil {
			return err
		}
		if status != BookingPending && status != BookingConfirmed {
			return ErrNotPending
		}
		var items []domain.BookingItem
		rows, err := tx.Query(ctx, `SELECT range_id, starts_at, ends_at, lanes FROM booking_items WHERE tenant_id = $1 AND request_id = $2`, tenantID, id)
		if err != nil {
			return err
		}
		if items, err = pgx.CollectRows(rows, pgx.RowToStructByPos[domain.BookingItem]); err != nil {
			return err
		}
		if !domain.SameBookingItems(items, in.Items) {
			if _, err := tx.Exec(ctx, `DELETE FROM booking_items WHERE tenant_id = $1 AND request_id = $2`, tenantID, id); err != nil {
				return err
			}
			if err := insertBookingItems(ctx, tx, tenantID, id, in.Items); err != nil {
				return err
			}
			reopened = status == BookingConfirmed
			status = BookingPending
		}
		if _, err := tx.Exec(ctx, `UPDATE booking_requests SET club = $3, contact = $4, email = $5, phone = $6, note = $7, status = $8::booking_status,
			decided_at = CASE WHEN $8::booking_status = 'pending' THEN NULL ELSE decided_at END, updated_at = now() WHERE tenant_id = $1 AND id = $2`,
			tenantID, id, in.Club, in.Contact, in.Email, in.Phone, in.Note, status); err != nil {
			return err
		}
		return syncBookingShooters(ctx, tx, tenantID, id, in.Shooters)
	})
	return reopened, mapErr(err)
}

// CancelBookingRequest cancels a pending or confirmed request and frees its reservations.
func (s *Store) CancelBookingRequest(ctx context.Context, tenantID, id int64) error {
	return s.closeBooking(ctx, tenantID, id, BookingCancelled, "")
}

// DeclineBookingRequest declines a pending or confirmed request with a reason and frees its reservations.
func (s *Store) DeclineBookingRequest(ctx context.Context, tenantID, id int64, reason string) error {
	return s.closeBooking(ctx, tenantID, id, BookingDeclined, reason)
}

// closeBooking moves a pending or confirmed request to a final status and deletes its reservations.
func (s *Store) closeBooking(ctx context.Context, tenantID, id int64, status BookingStatus, reason string) error {
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE booking_requests SET status = $3, decline_reason = $4, decided_at = now(), updated_at = now()
			WHERE tenant_id = $1 AND id = $2 AND status IN ('pending', 'confirmed')`, tenantID, id, status, reason)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			var exists bool
			if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM booking_requests WHERE tenant_id = $1 AND id = $2)`, tenantID, id).Scan(&exists); err != nil {
				return err
			}
			if exists {
				return ErrNotPending
			}
			return ErrNotFound
		}
		_, err = tx.Exec(ctx, `DELETE FROM reservations WHERE tenant_id = $1 AND booking_item_id IN
			(SELECT id FROM booking_items WHERE request_id = $2)`, tenantID, id)
		return err
	})
	if errors.Is(err, ErrNotPending) {
		return err
	}
	return mapErr(err)
}

// ConfirmBooking confirms a pending request: every item gets free lanes of its range (neighbours preferred) as a
// reservation described with the club, and the shooters are taken into the master data with clubID, or without
// one with the club of the request's name (regardless of case and spaces, created if missing). The ranges are locked
// like SaveReservation does; a *FullSlotError names the first item that found too few lanes, and nothing is stored
// then. ErrNotFound if the request or the club is not the tenant's.
func (s *Store) ConfirmBooking(ctx context.Context, tenantID, id int64, clubID *int64) error {
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var status BookingStatus
		var club string
		if err := tx.QueryRow(ctx, `SELECT status, club FROM booking_requests WHERE tenant_id = $1 AND id = $2 FOR UPDATE`,
			tenantID, id).Scan(&status, &club); err != nil {
			return err
		}
		if status != BookingPending {
			return ErrNotPending
		}
		rows, err := tx.Query(ctx, `SELECT i.id, i.range_id, i.starts_at, i.ends_at, i.lanes, rg.lane_count, rg.disabled_lanes
			FROM booking_items i JOIN ranges rg ON rg.id = i.range_id
			WHERE i.tenant_id = $1 AND i.request_id = $2 ORDER BY i.range_id, i.starts_at`, tenantID, id)
		if err != nil {
			return err
		}
		type item struct {
			ID       int64
			Item     domain.BookingItem
			Count    int
			Disabled []int32
		}
		items, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (item, error) {
			var x item
			err := r.Scan(&x.ID, &x.Item.RangeID, &x.Item.Start, &x.Item.End, &x.Item.Lanes, &x.Count, &x.Disabled)
			return x, err
		})
		if err != nil {
			return err
		}
		locked := map[int64]bool{}
		for _, it := range items {
			if !locked[it.Item.RangeID] {
				if _, err := tx.Exec(ctx, `SELECT id FROM ranges WHERE tenant_id = $1 AND id = $2 FOR UPDATE`, tenantID, it.Item.RangeID); err != nil {
					return err
				}
				locked[it.Item.RangeID] = true
			}
			var busy []int
			if err := tx.QueryRow(ctx, `SELECT coalesce(array_agg(DISTINCT l), '{}') FROM reservations, unnest(lanes) l
				WHERE range_id = $1 AND starts_at < $3 AND ends_at > $2`, it.Item.RangeID, it.Item.Start, it.Item.End).Scan(&busy); err != nil {
				return err
			}
			free := domain.FreeLanes(it.Count, it.Disabled, busy)
			lanes, ok := domain.AssignLanes(free, it.Item.Lanes)
			if !ok {
				return &FullSlotError{Item: it.Item, Free: len(free)}
			}
			if _, err := tx.Exec(ctx, `INSERT INTO reservations (tenant_id, range_id, lanes, starts_at, ends_at, description, booking_item_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`, tenantID, it.Item.RangeID, lanes, it.Item.Start, it.Item.End, club, it.ID); err != nil {
				return err
			}
		}
		cid, err := requestClub(ctx, tx, tenantID, club, clubID)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE booking_requests SET status = 'confirmed', decline_reason = '', decided_at = now(), updated_at = now(),
			club_id = $3 WHERE tenant_id = $1 AND id = $2`, tenantID, id, cid); err != nil {
			return err
		}
		return adoptBookingShooters(ctx, tx, tenantID, id, cid)
	})
	var full *FullSlotError
	if errors.As(err, &full) || errors.Is(err, ErrNotPending) {
		return err
	}
	return mapErr(err)
}

// AdoptBookingShooters takes the shooters of a confirmed request that are not linked yet into the master data,
// with the club the request was matched to.
func (s *Store) AdoptBookingShooters(ctx context.Context, tenantID, id int64) error {
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var status BookingStatus
		var name string
		var clubID *int64
		if err := tx.QueryRow(ctx, `SELECT status, club, club_id FROM booking_requests WHERE tenant_id = $1 AND id = $2 FOR UPDATE`,
			tenantID, id).Scan(&status, &name, &clubID); err != nil {
			return err
		}
		if status != BookingConfirmed {
			return ErrNotPending
		}
		club, err := requestClub(ctx, tx, tenantID, name, clubID)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE booking_requests SET club_id = $3 WHERE tenant_id = $1 AND id = $2`, tenantID, id, club); err != nil {
			return err
		}
		return adoptBookingShooters(ctx, tx, tenantID, id, club)
	})
	if errors.Is(err, ErrNotPending) {
		return err
	}
	return mapErr(err)
}

// requestClub returns clubID after checking that it is the tenant's (ErrNotFound otherwise), or without one the
// club of the request's name, created if missing.
func requestClub(ctx context.Context, tx pgx.Tx, tenantID int64, name string, clubID *int64) (int64, error) {
	if clubID == nil {
		id, _, err := ensureClub(ctx, tx, tenantID, strings.TrimSpace(name))
		return id, err
	}
	var ok bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM clubs WHERE tenant_id = $1 AND id = $2)`, tenantID, *clubID).Scan(&ok); err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrNotFound
	}
	return *clubID, nil
}

// adoptBookingShooters links every unlinked shooter of a request to the master data and to clubID. A shooter is
// found by licence (leading zeros ignored), else by name and birthday, and created otherwise. Found shooters are
// never overwritten: the list comes from a public form, so only a missing birthday or licence is filled in.
func adoptBookingShooters(ctx context.Context, tx pgx.Tx, tenantID, requestID, clubID int64) error {
	rows, err := tx.Query(ctx, `SELECT id, firstname, lastname, birthday::text, license_number FROM booking_shooters
		WHERE tenant_id = $1 AND request_id = $2 AND shooter_id IS NULL ORDER BY id`, tenantID, requestID)
	if err != nil {
		return err
	}
	type row struct {
		ID int64
		domain.BookingShooter
	}
	list, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (row, error) {
		var x row
		err := r.Scan(&x.ID, &x.Firstname, &x.Lastname, &x.Birthday, &x.LicenseNumber)
		return x, err
	})
	if err != nil || len(list) == 0 {
		return err
	}
	for _, b := range list {
		var shooterID int64
		err := pgx.ErrNoRows
		if b.LicenseNumber != "" {
			err = tx.QueryRow(ctx, `SELECT id FROM shooters WHERE tenant_id = $1 AND ltrim(license_number, '0') = ltrim($2, '0')
				ORDER BY id LIMIT 1`, tenantID, b.LicenseNumber).Scan(&shooterID)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			err = tx.QueryRow(ctx, `SELECT id FROM shooters WHERE tenant_id = $1 AND lower(firstname) = lower($2) AND lower(lastname) = lower($3)
				AND birthday = $4::date ORDER BY id LIMIT 1`, tenantID, b.Firstname, b.Lastname, b.Birthday).Scan(&shooterID)
		}
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			err = tx.QueryRow(ctx, `INSERT INTO shooters (tenant_id, license_number, firstname, lastname, birthday)
				VALUES ($1, nullif($2, ''), $3, $4, $5::date) RETURNING id`, tenantID, b.LicenseNumber, b.Firstname, b.Lastname, b.Birthday).Scan(&shooterID)
		case err == nil:
			_, err = tx.Exec(ctx, `UPDATE shooters SET birthday = coalesce(birthday, $3::date),
				license_number = CASE WHEN license_number IS NULL AND $4 <> '' AND NOT EXISTS
					(SELECT 1 FROM shooters o WHERE o.tenant_id = $1 AND o.license_number = $4) THEN $4 ELSE license_number END,
				updated_at = now() WHERE tenant_id = $1 AND id = $2`, tenantID, shooterID, b.Birthday, b.LicenseNumber)
		}
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO shooter_clubs (tenant_id, shooter_id, club_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			tenantID, shooterID, clubID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE booking_shooters SET shooter_id = $3 WHERE tenant_id = $1 AND id = $2`, tenantID, b.ID, shooterID); err != nil {
			return err
		}
	}
	return nil
}
