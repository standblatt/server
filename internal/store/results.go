package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// EntryRef is an entry with what result handling needs to know about it.
type EntryRef struct {
	ID             int64
	RegistrationID int64
	EventID        int64
	ProgrammeID    int64
	CombinedSeries []int
	AgeCategory    domain.AgeCategory
}

// GetEntry returns an entry of a tenant.
func (s *Store) GetEntry(ctx context.Context, tenantID, entryID int64) (EntryRef, error) {
	var e EntryRef
	err := s.pool.QueryRow(ctx, `SELECT en.id, en.registration_id, en.event_id, en.programme_id,
			coalesce(en.combined_series, '{}'), r.age_category
		FROM entries en JOIN registrations r ON r.id = en.registration_id
		WHERE en.tenant_id = $1 AND en.id = $2`, tenantID, entryID).
		Scan(&e.ID, &e.RegistrationID, &e.EventID, &e.ProgrammeID, &e.CombinedSeries, &e.AgeCategory)
	return e, mapErr(err)
}

// FindEntry returns a shooter's first attempt without results at a programme of an event of a tenant, or else the
// last attempt, and whether that one has results; ErrNotFound when the shooter has none.
func (s *Store) FindEntry(ctx context.Context, tenantID, eventID, shooterID, programmeID int64) (EntryRef, bool, error) {
	var e EntryRef
	var done bool
	err := s.pool.QueryRow(ctx, `SELECT en.id, en.registration_id, en.event_id, en.programme_id,
			coalesce(en.combined_series, '{}'), r.age_category, en.status = 'completed'
		FROM entries en JOIN registrations r ON r.id = en.registration_id
		WHERE en.tenant_id = $1 AND en.event_id = $2 AND en.shooter_id = $3 AND en.programme_id = $4
		ORDER BY en.status = 'completed', CASE WHEN en.status = 'completed' THEN -en.attempt ELSE en.attempt END LIMIT 1`, tenantID, eventID, shooterID, programmeID).
		Scan(&e.ID, &e.RegistrationID, &e.EventID, &e.ProgrammeID, &e.CombinedSeries, &e.AgeCategory, &done)
	return e, done, mapErr(err)
}

// CompletedAttempts counts a shooter's attempts with results at a programme of an event of a tenant.
func (s *Store) CompletedAttempts(ctx context.Context, tenantID, eventID, shooterID, programmeID int64) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `SELECT count(*) FROM entries WHERE tenant_id = $1 AND event_id = $2 AND shooter_id = $3
		AND programme_id = $4 AND status = 'completed'`, tenantID, eventID, shooterID, programmeID).Scan(&n)
	return n, err
}

// Attempts returns every attempt of the shooter at the programme of the given entry, in attempt order.
func (s *Store) Attempts(ctx context.Context, tenantID, entryID int64) ([]EntryRef, error) {
	rows, err := s.pool.Query(ctx, `SELECT o.id, o.registration_id, o.event_id, o.programme_id, coalesce(o.combined_series, '{}'), r.age_category
		FROM entries e
		JOIN entries o ON o.event_id = e.event_id AND o.shooter_id = e.shooter_id AND o.programme_id = e.programme_id
		JOIN registrations r ON r.id = o.registration_id
		WHERE e.tenant_id = $1 AND e.id = $2 ORDER BY o.attempt`, tenantID, entryID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[EntryRef])
}

// ListResults returns the series results of all entries of a registration, keyed by entry.
func (s *Store) ListResults(ctx context.Context, tenantID, registrationID int64) (map[int64][]domain.SeriesResult, error) {
	rows, err := s.pool.Query(ctx, `SELECT re.entry_id, re.points, re.mouches, re.shots FROM results re
		JOIN entries en ON en.id = re.entry_id
		WHERE en.tenant_id = $1 AND en.registration_id = $2 ORDER BY re.entry_id, re.series_number`, tenantID, registrationID)
	if err != nil {
		return nil, err
	}
	out := map[int64][]domain.SeriesResult{}
	var id int64
	var r domain.SeriesResult
	_, err = pgx.ForEachRow(rows, []any{&id, &r.Points, &r.Mouches, &r.Shots}, func() error {
		out[id] = append(out[id], r)
		r = domain.SeriesResult{}
		return nil
	})
	return out, err
}

// ProgrammeTiers returns the award tiers of the given programmes of a tenant, keyed by programme.
func (s *Store) ProgrammeTiers(ctx context.Context, tenantID int64, programmeIDs []int64) (map[int64][]domain.Tier, error) {
	rows, err := s.pool.Query(ctx, `SELECT t.programme_id, t.id, t.position, t.name, t.kind, t.value,
			coalesce((SELECT json_agg(json_build_object('ageCategory', l.age_category, 'points', l.points) ORDER BY l.age_category)
				FROM award_limits l WHERE l.tier_id = t.id), '[]')
		FROM award_tiers t JOIN programmes p ON p.id = t.programme_id
		WHERE p.tenant_id = $1 AND t.programme_id = ANY($2) ORDER BY t.programme_id, t.position`, tenantID, programmeIDs)
	if err != nil {
		return nil, err
	}
	out := map[int64][]domain.Tier{}
	var pid int64
	var t domain.Tier
	_, err = pgx.ForEachRow(rows, []any{&pid, &t.ID, &t.Position, &t.Name, &t.Kind, &t.Value, &t.Limits}, func() error {
		out[pid] = append(out[pid], t)
		t = domain.Tier{}
		return nil
	})
	return out, err
}

// replaceResults replaces the results of an entry, clears its award and marks it completed, or pending without results.
func replaceResults(ctx context.Context, tx pgx.Tx, entryID int64, results []domain.SeriesResult) error {
	if _, err := tx.Exec(ctx, `DELETE FROM results WHERE entry_id = $1`, entryID); err != nil {
		return err
	}
	for i, r := range results {
		if _, err := tx.Exec(ctx, `INSERT INTO results (entry_id, series_number, points, mouches, shots) VALUES ($1, $2, $3, $4, $5)`,
			entryID, i+1, r.Points, r.Mouches, r.Shots); err != nil {
			return err
		}
	}
	status := "completed"
	if len(results) == 0 {
		status = "pending"
	}
	_, err := tx.Exec(ctx, `UPDATE entries SET status = $2::entry_status, award_tier_id = NULL WHERE id = $1`, entryID, status)
	return err
}

// SaveResults replaces the results of an entry, clears its award and marks it completed. Entries of the same
// registration that take series from this programme get the picked series copied and their award cleared.
// When every entry of the registration is completed, the registration completes too. It returns whether
// that happened.
func (s *Store) SaveResults(ctx context.Context, tenantID int64, e EntryRef, results []domain.SeriesResult) (bool, error) {
	var completed bool
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		// Saves of the same registration run one after the other: two at once would both delete and then both insert.
		if _, err := tx.Exec(ctx, `SELECT 1 FROM registrations WHERE id = $1 AND tenant_id = $2 FOR UPDATE`, e.RegistrationID, tenantID); err != nil {
			return err
		}
		if err := replaceResults(ctx, tx, e.ID, results); err != nil {
			return err
		}
		rows, err := tx.Query(ctx, `SELECT en.id, en.combined_series FROM entries en
			JOIN programmes p ON p.id = en.programme_id
			WHERE en.tenant_id = $1 AND en.registration_id = $2 AND p.combinable_with_id = $3
				AND cardinality(en.combined_series) > 0`, tenantID, e.RegistrationID, e.ProgrammeID)
		if err != nil {
			return err
		}
		type source struct {
			id     int64
			series []int
		}
		sources, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (source, error) {
			var src source
			err := row.Scan(&src.id, &src.series)
			return src, err
		})
		if err != nil {
			return err
		}
		for _, src := range sources {
			if err := replaceResults(ctx, tx, src.id, domain.CombinedResults(results, src.series)); err != nil {
				return err
			}
		}
		tag, err := tx.Exec(ctx, `UPDATE registrations SET status = 'completed', completed_at = now()
			WHERE id = $1 AND status <> 'completed'
			AND NOT EXISTS (SELECT 1 FROM entries WHERE registration_id = $1 AND status <> 'completed')`, e.RegistrationID)
		completed = tag.RowsAffected() == 1
		return err
	})
	return completed, err
}

// SetAward stores the award tier an entry received, or clears it with nil.
func (s *Store) SetAward(ctx context.Context, tenantID, entryID int64, tierID *int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE entries SET award_tier_id = $3 WHERE tenant_id = $1 AND id = $2`,
		tenantID, entryID, tierID))
}
