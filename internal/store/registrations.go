package store

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/modules"
)

// Errors for the registration constraints.
var (
	ErrLaneTaken      = errors.New("store: lane is occupied")
	ErrProgrammeTaken = errors.New("store: programme already registered for this shooter")
	ErrHasResults     = errors.New("store: entry has results")
	ErrAttemptsUsed   = errors.New("store: shooters already have more attempts")
)

// Registration statuses: waiting without a lane, shooting on a lane, completed.
const (
	StatusRegistered = "registered"
	StatusShooting   = "shooting"
	StatusCompleted  = "completed"
)

// Entry is one attempt at a programme that a registration shoots.
type Entry struct {
	ID             int64  `json:"id"`
	ProgrammeID    int64  `json:"programmeId"`
	Attempt        int    `json:"attempt" doc:"Number of the shooter's attempt at the programme in this event, from 1"`
	Status         string `json:"status" enum:"pending,completed"`
	AwardTierID    *int64 `json:"awardTierId,omitempty"`
	CombinedSeries []int  `json:"combinedSeries"`
}

// RegisteredShooter is the shooter as shown on a registration.
type RegisteredShooter struct {
	ID            int64  `json:"id"`
	LicenseNumber string `json:"licenseNumber"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Birthday      string `json:"birthday"`
}

// Registration is a visit of a shooter at an event.
type Registration struct {
	ID             int64              `json:"id"`
	EventID        int64              `json:"eventId"`
	Shooter        RegisteredShooter  `json:"shooter"`
	ClubID         *int64             `json:"clubId,omitempty"`
	ClubName       string             `json:"clubName"`
	AgeCategory    domain.AgeCategory `json:"ageCategory" enum:"u10,u13,u15,u17,u19,u21,elite,senior,veteran,senior_veteran"`
	RangeID        *int64             `json:"rangeId,omitempty" doc:"Range of the lane"`
	Lane           *int               `json:"lane,omitempty"`
	Status         string             `json:"status" enum:"registered,shooting,completed"`
	RegisteredAt   time.Time          `json:"registeredAt"`
	StartedAt      *time.Time         `json:"startedAt,omitempty"`
	CompletedAt    *time.Time         `json:"completedAt,omitempty"`
	AmmoIssuedAt   *time.Time         `json:"ammoIssuedAt,omitempty" doc:"When the ammunition desk handed out the ammunition"`
	AwardsIssuedAt *time.Time         `json:"awardsIssuedAt,omitempty" doc:"When the shooter collected the awards (wreath cards, prizes, cash)"`
	Entries        []Entry            `json:"entries"`
}

// NewRegistration is a registration to insert; the caller has validated it.
type NewRegistration struct {
	EventID     int64
	ShooterID   int64
	ClubID      *int64
	AgeCategory domain.AgeCategory
	RangeID     *int64
	Lane        *int
	Programmes  []domain.ProgrammeChoice
}

// RangeRules is what registration checks need to know about a range of the event.
type RangeRules struct {
	Name          string
	LaneCount     int
	DisabledLanes []int32
}

// EventRules is what registration checks need to know about an event.
type EventRules struct {
	StartsOn   string
	Template   bool
	Lanes      bool
	Ranges     map[int64]RangeRules
	RangeOrder []int64
	Programmes map[int64]domain.ProgrammeInfo
}

// registrationErr maps the lane and programme unique violations to ErrLaneTaken and ErrProgrammeTaken, others via mapErr.
func registrationErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "registrations_lane_taken":
			return ErrLaneTaken
		case "entries_programme_once":
			return ErrProgrammeTaken
		}
	}
	return mapErr(err)
}

// LoadEventRules returns the event date, ranges with their lanes and the programmes of an event of a tenant.
func (s *Store) LoadEventRules(ctx context.Context, tenantID, eventID int64) (EventRules, error) {
	r := EventRules{Ranges: map[int64]RangeRules{}, Programmes: map[int64]domain.ProgrammeInfo{}}
	var chosen, used []string
	err := s.pool.QueryRow(ctx, `SELECT e.starts_on::text, e.template, e.modules, `+eventUsed+` FROM events e WHERE e.tenant_id = $1 AND e.id = $2`,
		tenantID, eventID).Scan(&r.StartsOn, &r.Template, &chosen, &used)
	if err != nil {
		return r, mapErr(err)
	}
	rows, err := s.pool.Query(ctx, `SELECT rg.id, rg.name, rg.lane_count, rg.disabled_lanes FROM event_ranges er
		JOIN ranges rg ON rg.id = er.range_id WHERE er.event_id = $1 ORDER BY er.position`, eventID)
	if err != nil {
		return r, err
	}
	var rid int64
	var rr RangeRules
	if _, err := pgx.ForEachRow(rows, []any{&rid, &rr.Name, &rr.LaneCount, &rr.DisabledLanes}, func() error {
		r.Ranges[rid] = RangeRules{Name: rr.Name, LaneCount: rr.LaneCount, DisabledLanes: slices.Clone(rr.DisabledLanes)}
		r.RangeOrder = append(r.RangeOrder, rid)
		return nil
	}); err != nil {
		return r, err
	}
	r.Lanes = slices.Contains(modules.Effective(modules.Event, chosen, used, len(r.RangeOrder) > 0), modules.Lanes)
	rows, err = s.pool.Query(ctx, `SELECT p.id, p.name, p.variant, p.active, p.combinable_with_id, p.range_id, p.attempts,
			(SELECT count(*) FROM programme_series s WHERE s.programme_id = p.id),
			coalesce((SELECT array_agg(age_category::text) FROM programme_prices pp WHERE pp.programme_id = p.id), '{}')
		FROM programmes p
		WHERE p.tenant_id = $1 AND p.event_id = $2`, tenantID, eventID)
	if err != nil {
		return r, err
	}
	var p domain.ProgrammeInfo
	var cats []string
	_, err = pgx.ForEachRow(rows, []any{&p.ID, &p.Name, &p.Variant, &p.Active, &p.CombinableWithID, &p.RangeID, &p.AttemptLimit, &p.Series, &cats}, func() error {
		info := p
		info.Categories = make([]domain.AgeCategory, len(cats))
		for i, c := range cats {
			info.Categories[i] = domain.AgeCategory(c)
		}
		r.Programmes[p.ID] = info
		return nil
	})
	return r, err
}

const registrationCols = `r.id, r.event_id, sh.id, coalesce(sh.license_number, ''), sh.firstname, sh.lastname,
	coalesce(sh.birthday::text, ''), r.club_id, coalesce(c.name, ''), r.age_category, r.range_id, r.lane, r.status,
	r.registered_at, r.started_at, r.completed_at, r.ammo_issued_at, r.awards_issued_at,
	coalesce((SELECT json_agg(json_build_object('id', en.id, 'programmeId', en.programme_id, 'attempt', en.attempt, 'status', en.status,
		'awardTierId', en.award_tier_id, 'combinedSeries', coalesce(en.combined_series, '{}')) ORDER BY en.id)
		FROM entries en WHERE en.registration_id = r.id), '[]')`

const registrationFrom = ` FROM registrations r JOIN shooters sh ON sh.id = r.shooter_id LEFT JOIN clubs c ON c.id = r.club_id `

// scanRegistration scans a row of registrationCols into a Registration.
func scanRegistration(row pgx.CollectableRow) (Registration, error) {
	var r Registration
	err := row.Scan(&r.ID, &r.EventID, &r.Shooter.ID, &r.Shooter.LicenseNumber, &r.Shooter.Firstname, &r.Shooter.Lastname,
		&r.Shooter.Birthday, &r.ClubID, &r.ClubName, &r.AgeCategory, &r.RangeID, &r.Lane, &r.Status,
		&r.RegisteredAt, &r.StartedAt, &r.CompletedAt, &r.AmmoIssuedAt, &r.AwardsIssuedAt, &r.Entries)
	return r, err
}

// ListRegistrations lists the registrations of an event; only waiting and shooting ones when active is set.
func (s *Store) ListRegistrations(ctx context.Context, tenantID, eventID int64, active bool) ([]Registration, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+registrationCols+registrationFrom+`
		WHERE r.tenant_id = $1 AND r.event_id = $2 AND (NOT $3 OR r.status <> 'completed')
		ORDER BY r.lane NULLS FIRST, r.registered_at`, tenantID, eventID, active)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanRegistration)
}

// GetRegistration returns one registration of a tenant.
func (s *Store) GetRegistration(ctx context.Context, tenantID, registrationID int64) (Registration, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+registrationCols+registrationFrom+` WHERE r.tenant_id = $1 AND r.id = $2`,
		tenantID, registrationID)
	if err != nil {
		return Registration{}, err
	}
	r, err := pgx.CollectExactlyOneRow(rows, scanRegistration)
	return r, mapErr(err)
}

// statusFor returns shooting when the registration has a lane, else registered.
func statusFor(lane *int) string {
	if lane != nil {
		return StatusShooting
	}
	return StatusRegistered
}

// insertEntries inserts the attempts of every chosen programme of a registration, numbered after the shooter's
// last attempt so the numbers stay in order; ErrProgrammeTaken when the shooter would have more attempts than allowed.
func insertEntries(ctx context.Context, tx pgx.Tx, tenantID, registrationID, eventID, shooterID int64, choices []domain.ProgrammeChoice) error {
	for _, c := range choices {
		for range c.Attempts() {
			tag, err := tx.Exec(ctx, `INSERT INTO entries (tenant_id, registration_id, event_id, shooter_id, programme_id, combined_series, attempt)
				SELECT $1, $2, $3, $4, p.id, $6, n.next
				FROM programmes p, LATERAL (SELECT coalesce(max(attempt), 0) + 1 AS next, count(*) AS used FROM entries
					WHERE event_id = $3 AND shooter_id = $4 AND programme_id = p.id) n
				WHERE p.id = $5 AND (p.attempts = 0 OR n.used < p.attempts)`,
				tenantID, registrationID, eventID, shooterID, c.ProgrammeID, nilIfEmpty(c.CombinedSeries))
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				return ErrProgrammeTaken
			}
		}
	}
	return nil
}

// nilIfEmpty returns nil for an empty slice, so it is stored as NULL.
func nilIfEmpty(s []int) []int {
	if len(s) == 0 {
		return nil
	}
	return s
}

// CreateRegistration inserts a registration with its entries. ErrLaneTaken or ErrProgrammeTaken when another
// registration holds the lane or the programme.
func (s *Store) CreateRegistration(ctx context.Context, tenantID int64, in NewRegistration) (int64, error) {
	var id int64
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		status := statusFor(in.Lane)
		err := tx.QueryRow(ctx, `INSERT INTO registrations (tenant_id, event_id, shooter_id, club_id, age_category, range_id, lane, status, started_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8::registration_status, CASE WHEN $8::registration_status = 'shooting' THEN now() END) RETURNING id`,
			tenantID, in.EventID, in.ShooterID, in.ClubID, in.AgeCategory, in.RangeID, in.Lane, status).Scan(&id)
		if err != nil {
			return err
		}
		return insertEntries(ctx, tx, tenantID, id, in.EventID, in.ShooterID, in.Programmes)
	})
	return id, registrationErr(err)
}

// UpdateRegistration changes lane and programmes of a registration. Entries of programmes that stay keep
// their status and award; removed ones are deleted, and so are surplus attempts when the count goes down,
// those without results from the last one (ErrHasResults when results would be lost). The status follows the
// lane: shooting with one, registered without.
func (s *Store) UpdateRegistration(ctx context.Context, tenantID int64, current Registration, rangeID *int64, lane *int, choices []domain.ProgrammeChoice) error {
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		status := statusFor(lane)
		_, err := tx.Exec(ctx, `UPDATE registrations SET range_id = $5, lane = $3, status = $4::registration_status, completed_at = NULL,
				started_at = CASE WHEN $4::registration_status = 'registered' THEN NULL
					WHEN lane IS DISTINCT FROM $3 OR range_id IS DISTINCT FROM $5 OR started_at IS NULL THEN now() ELSE started_at END
			WHERE tenant_id = $1 AND id = $2`, tenantID, current.ID, lane, status, rangeID)
		if err != nil {
			return err
		}
		keep := map[int64]domain.ProgrammeChoice{}
		for _, c := range choices {
			keep[c.ProgrammeID] = c
		}
		var added []domain.ProgrammeChoice
		have := map[int64]int{}
		surplus := map[int64]int{}
		for _, e := range current.Entries {
			have[e.ProgrammeID]++
		}
		for id, n := range have {
			surplus[id] = n - keep[id].Attempts()
			if _, ok := keep[id]; !ok {
				surplus[id] = n
			}
		}
		entries := slices.Clone(current.Entries)
		slices.SortStableFunc(entries, func(a, b Entry) int { return b.Attempt - a.Attempt })
		var scoredIDs []int64
		if err := tx.QueryRow(ctx, `SELECT coalesce(array_agg(DISTINCT entry_id), '{}') FROM results WHERE entry_id = ANY($1)`,
			entryIDs(entries)).Scan(&scoredIDs); err != nil {
			return err
		}
		for _, pass := range []bool{false, true} {
			for _, e := range entries {
				if surplus[e.ProgrammeID] <= 0 || slices.Contains(scoredIDs, e.ID) != pass {
					continue
				}
				if pass && len(e.CombinedSeries) == 0 {
					return ErrHasResults
				}
				if _, err := tx.Exec(ctx, `DELETE FROM entries WHERE id = $1`, e.ID); err != nil {
					return err
				}
				surplus[e.ProgrammeID]--
			}
		}
		for _, e := range current.Entries {
			c, ok := keep[e.ProgrammeID]
			if !ok {
				continue
			}
			if _, err := tx.Exec(ctx, `UPDATE entries SET combined_series = $2 WHERE id = $1`, e.ID, nilIfEmpty(c.CombinedSeries)); err != nil {
				return err
			}
			if !slices.Equal(e.CombinedSeries, c.CombinedSeries) && (len(e.CombinedSeries) > 0 || len(c.CombinedSeries) > 0) {
				if err := recombine(ctx, tx, e.ID, current.ID, e.ProgrammeID, c.CombinedSeries); err != nil {
					return err
				}
			}
		}
		for _, c := range choices {
			if n := c.Attempts() - have[c.ProgrammeID]; n > 0 {
				c.Count = n
				added = append(added, c)
			}
		}
		return insertEntries(ctx, tx, tenantID, current.ID, current.EventID, current.Shooter.ID, added)
	})
	return registrationErr(err)
}

// recombine recomputes the copied results of a combined entry after its series choice changed; without a
// choice (no longer combined) the copied results are removed.
func recombine(ctx context.Context, tx pgx.Tx, entryID, registrationID, programmeID int64, picked []int) error {
	if len(picked) == 0 {
		return replaceResults(ctx, tx, entryID, nil)
	}
	rows, err := tx.Query(ctx, `SELECT re.points, re.mouches FROM results re
		JOIN entries t ON t.id = re.entry_id
		JOIN programmes p ON p.combinable_with_id = t.programme_id
		WHERE t.registration_id = $1 AND p.id = $2 ORDER BY re.series_number`, registrationID, programmeID)
	if err != nil {
		return err
	}
	target, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.SeriesResult, error) {
		var r domain.SeriesResult
		err := row.Scan(&r.Points, &r.Mouches)
		return r, err
	})
	if err != nil {
		return err
	}
	if len(target) == 0 {
		return replaceResults(ctx, tx, entryID, nil)
	}
	return replaceResults(ctx, tx, entryID, domain.CombinedResults(target, picked))
}

// FinishRegistration marks a registration completed, which frees its lane.
func (s *Store) FinishRegistration(ctx context.Context, tenantID, registrationID int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE registrations SET status = 'completed', completed_at = now()
		WHERE tenant_id = $1 AND id = $2`, tenantID, registrationID))
}

// DeleteRegistration removes a registration with its entries.
func (s *Store) DeleteRegistration(ctx context.Context, tenantID, registrationID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM registrations WHERE tenant_id = $1 AND id = $2`, tenantID, registrationID))
}

// EventStats counts the registrations of an event by status.
type EventStats struct {
	Registered int `json:"registered"`
	Shooting   int `json:"shooting"`
	Completed  int `json:"completed"`
}

// GetEventStats counts the registrations of an event by status.
func (s *Store) GetEventStats(ctx context.Context, tenantID, eventID int64) (EventStats, error) {
	var st EventStats
	err := s.pool.QueryRow(ctx, `SELECT
			count(*) FILTER (WHERE status = 'registered'),
			count(*) FILTER (WHERE status = 'shooting'),
			count(*) FILTER (WHERE status = 'completed')
		FROM registrations WHERE tenant_id = $1 AND event_id = $2`, tenantID, eventID).
		Scan(&st.Registered, &st.Shooting, &st.Completed)
	return st, err
}

// entryIDs returns the IDs of the entries.
func entryIDs(entries []Entry) []int64 {
	ids := make([]int64, len(entries))
	for i, e := range entries {
		ids[i] = e.ID
	}
	return ids
}
