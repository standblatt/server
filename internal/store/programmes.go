package store

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// Programme is a programme of an event with its full definition.
type Programme struct {
	ID               int64                    `json:"id"`
	EventID          int64                    `json:"eventId"`
	Active           bool                     `json:"active"`
	CombinableWithID *int64                   `json:"combinableWithId,omitempty" doc:"Programme whose series this one may take at registration"`
	RangeID          *int64                   `json:"rangeId,omitempty" doc:"Range the programme is shot on; unset means any range of the event"`
	LayoutID         *int64                   `json:"layoutId,omitempty" doc:"Print form the results of the programme are printed onto"`
	Shots            int                      `json:"shots"`
	MaxPoints        int                      `json:"maxPoints"`
	SpecialPrizes    domain.SpecialPrizeRules `json:"specialPrizes"`
	InUse            bool                     `json:"inUse" doc:"Shooters are registered: series are fixed, age categories and awards can only be added"`
	domain.ProgrammeDefinition
}

// ProgrammeTemplate is a programme definition in the global library.
type ProgrammeTemplate struct {
	ID         int64                      `json:"id"`
	Definition domain.ProgrammeDefinition `json:"definition"`
}

// loadProgrammes loads the programmes matching the WHERE fragment with their series, prices and award tiers.
func loadProgrammes(ctx context.Context, q querier, where string, args ...any) ([]Programme, error) {
	rows, err := q.Query(ctx, `SELECT id, event_id, active, combinable_with_id, range_id, layout_id, name, variant, description, barcode,
			counts_for_club_competition, show_on_rankings, prize_ranks, prize_min_category, prize_min_club,
			EXISTS (SELECT 1 FROM entries e WHERE e.programme_id = p.id), attempts, attempt_rule, award_counting_only
		FROM programmes p WHERE `+where+` ORDER BY lower(name), lower(variant), id`, args...)
	if err != nil {
		return nil, err
	}
	list, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Programme, error) {
		var p Programme
		err := row.Scan(&p.ID, &p.EventID, &p.Active, &p.CombinableWithID, &p.RangeID, &p.LayoutID, &p.Name, &p.Variant, &p.Description, &p.Barcode,
			&p.CountsForClubCompetition, &p.ShowOnRankings, &p.SpecialPrizes.Ranks, &p.SpecialPrizes.MinCategory, &p.SpecialPrizes.MinClub, &p.InUse,
			&p.Attempts, &p.AttemptRule, &p.AwardCountingOnly)
		p.UnlimitedAttempts = p.Attempts == 0
		return p, err
	})
	if err != nil || len(list) == 0 {
		return list, err
	}
	ids := make([]int64, len(list))
	byID := map[int64]*Programme{}
	for i := range list {
		ids[i] = list[i].ID
		list[i].Series, list[i].Prices, list[i].Awards = []domain.Series{}, []domain.Price{}, []domain.AwardTier{}
		byID[list[i].ID] = &list[i]
	}

	rows, err = q.Query(ctx, `SELECT programme_id, shots, max_points, name, track_mouches, track_shots FROM programme_series
		WHERE programme_id = ANY($1) ORDER BY programme_id, number`, ids)
	if err != nil {
		return nil, err
	}
	var pid int64
	var sr domain.Series
	if _, err := pgx.ForEachRow(rows, []any{&pid, &sr.Shots, &sr.MaxPoints, &sr.Name, &sr.TrackMouches, &sr.TrackShots}, func() error {
		byID[pid].Series = append(byID[pid].Series, sr)
		return nil
	}); err != nil {
		return nil, err
	}

	rows, err = q.Query(ctx, `SELECT programme_id, age_category, price, repeat_price FROM programme_prices
		WHERE programme_id = ANY($1) ORDER BY programme_id, age_category`, ids)
	if err != nil {
		return nil, err
	}
	var pr domain.Price
	if _, err := pgx.ForEachRow(rows, []any{&pid, &pr.AgeCategory, &pr.Price, &pr.RepeatPrice}, func() error {
		byID[pid].Prices = append(byID[pid].Prices, pr)
		pr.RepeatPrice = nil
		return nil
	}); err != nil {
		return nil, err
	}

	rows, err = q.Query(ctx, `SELECT t.programme_id, t.id, t.name, t.kind, t.value,
			coalesce((SELECT json_agg(json_build_object('ageCategory', l.age_category, 'points', l.points) ORDER BY l.age_category)
				FROM award_limits l WHERE l.tier_id = t.id), '[]')
		FROM award_tiers t WHERE t.programme_id = ANY($1) ORDER BY t.programme_id, t.position`, ids)
	if err != nil {
		return nil, err
	}
	var tierID int64
	var tier domain.AwardTier
	if _, err := pgx.ForEachRow(rows, []any{&pid, &tierID, &tier.Name, &tier.Kind, &tier.Value, &tier.Limits}, func() error {
		byID[pid].Awards = append(byID[pid].Awards, tier)
		tier.Limits = nil
		return nil
	}); err != nil {
		return nil, err
	}
	for i := range list {
		list[i].Shots, list[i].MaxPoints = list[i].ProgrammeDefinition.Shots(), list[i].ProgrammeDefinition.MaxPoints()
	}
	return list, nil
}

// ListProgrammes lists the programmes of an event.
func (s *Store) ListProgrammes(ctx context.Context, tenantID, eventID int64) ([]Programme, error) {
	return loadProgrammes(ctx, s.pool, "tenant_id = $1 AND event_id = $2", tenantID, eventID)
}

// GetProgramme returns one programme of a tenant.
func (s *Store) GetProgramme(ctx context.Context, tenantID, programmeID int64) (Programme, error) {
	list, err := loadProgrammes(ctx, s.pool, "tenant_id = $1 AND id = $2", tenantID, programmeID)
	if err != nil {
		return Programme{}, err
	}
	if len(list) == 0 {
		return Programme{}, ErrNotFound
	}
	return list[0], nil
}

// insertDefinition inserts the series, prices, award tiers and limits of a definition for a programme.
func insertDefinition(ctx context.Context, tx pgx.Tx, programmeID int64, d domain.ProgrammeDefinition) error {
	for i, sr := range d.Series {
		if _, err := tx.Exec(ctx, `INSERT INTO programme_series (programme_id, number, shots, max_points, name, track_mouches, track_shots)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`, programmeID, i+1, sr.Shots, sr.MaxPoints, sr.Name, sr.TrackMouches, sr.TrackShots); err != nil {
			return err
		}
	}
	for _, p := range d.Prices {
		if _, err := tx.Exec(ctx, `INSERT INTO programme_prices (programme_id, age_category, price, repeat_price) VALUES ($1, $2, $3, $4)`,
			programmeID, p.AgeCategory, p.Price, p.RepeatPrice); err != nil {
			return err
		}
	}
	for i, a := range d.Awards {
		var tierID int64
		if err := tx.QueryRow(ctx, `INSERT INTO award_tiers (programme_id, position, name, kind, value)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`, programmeID, i+1, a.Name, a.Kind, a.Value).Scan(&tierID); err != nil {
			return err
		}
		for _, l := range a.Limits {
			if _, err := tx.Exec(ctx, `INSERT INTO award_limits (tier_id, age_category, points) VALUES ($1, $2, $3)`,
				tierID, l.AgeCategory, l.Points); err != nil {
				return err
			}
		}
	}
	return nil
}

// createProgramme inserts a programme of an event with its definition and returns its ID.
func createProgramme(ctx context.Context, tx pgx.Tx, tenantID, eventID int64, d domain.ProgrammeDefinition) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `INSERT INTO programmes (tenant_id, event_id, name, variant, description, barcode,
			counts_for_club_competition, show_on_rankings, attempts, attempt_rule, award_counting_only)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		tenantID, eventID, d.Name, d.Variant, d.Description, d.Barcode, d.CountsForClubCompetition, d.ShowOnRankings,
		d.AttemptLimit(), d.Rule(), d.AwardCountingOnly).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, insertDefinition(ctx, tx, id, d)
}

// eventInTenant returns ErrNotFound unless the event belongs to the tenant.
func eventInTenant(ctx context.Context, q querier, tenantID, eventID int64) error {
	var ok bool
	if err := q.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM events WHERE tenant_id = $1 AND id = $2)`,
		tenantID, eventID).Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

// CreateProgramme adds a programme to an event. ErrNotFound if the event does not belong to the tenant,
// ErrExists if it already has a programme with this name and variant.
func (s *Store) CreateProgramme(ctx context.Context, tenantID, eventID int64, d domain.ProgrammeDefinition) (int64, error) {
	var id int64
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := eventInTenant(ctx, tx, tenantID, eventID); err != nil {
			return err
		}
		var err error
		id, err = createProgramme(ctx, tx, tenantID, eventID, d)
		return err
	})
	return id, mapErr(err)
}

// UpdateProgramme replaces the definition of a programme. Once shooters are registered for it, the series must
// stay as they are and age categories and awards can be added but not removed (ErrInUse otherwise), so results
// and given awards keep their meaning. The attempts may change while in use, but not below the most attempts a
// shooter already has (ErrAttemptsUsed).
func (s *Store) UpdateProgramme(ctx context.Context, tenantID, programmeID int64, d domain.ProgrammeDefinition) error {
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var used bool
		var most int
		if err := tx.QueryRow(ctx, `SELECT count(*) > 0, coalesce(max(n), 0) FROM (SELECT count(*) AS n FROM entries
				WHERE tenant_id = $1 AND programme_id = $2 GROUP BY shooter_id) a`,
			tenantID, programmeID).Scan(&used, &most); err != nil {
			return err
		}
		if limit := d.AttemptLimit(); limit > 0 && most > limit {
			return ErrAttemptsUsed
		}
		tag, err := tx.Exec(ctx, `UPDATE programmes SET name = $3, variant = $4, description = $5, barcode = $6,
				counts_for_club_competition = $7, show_on_rankings = $8, attempts = $9, attempt_rule = $10, award_counting_only = $11
			WHERE tenant_id = $1 AND id = $2`,
			tenantID, programmeID, d.Name, d.Variant, d.Description, d.Barcode, d.CountsForClubCompetition, d.ShowOnRankings,
			d.AttemptLimit(), d.Rule(), d.AwardCountingOnly)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		if used {
			return updateDefinitionInUse(ctx, tx, programmeID, d)
		}
		for _, table := range []string{"programme_series", "programme_prices", "award_tiers"} {
			if _, err := tx.Exec(ctx, `DELETE FROM `+table+` WHERE programme_id = $1`, programmeID); err != nil {
				return err
			}
		}
		return insertDefinition(ctx, tx, programmeID, d)
	})
	return mapErr(err)
}

// updateDefinitionInUse applies the changes allowed while a programme has entries (names, prices, award values,
// new tiers); ErrInUse if the series change or prices or tiers would be removed.
func updateDefinitionInUse(ctx context.Context, tx pgx.Tx, programmeID int64, d domain.ProgrammeDefinition) error {
	rows, err := tx.Query(ctx, `SELECT shots, max_points, name, track_mouches, track_shots FROM programme_series WHERE programme_id = $1 ORDER BY number`, programmeID)
	if err != nil {
		return err
	}
	series, err := pgx.CollectRows(rows, pgx.RowToStructByPos[domain.Series])
	if err != nil {
		return err
	}
	if len(series) != len(d.Series) {
		return ErrInUse
	}
	for i, sr := range d.Series {
		if sr.Shots != series[i].Shots || sr.MaxPoints != series[i].MaxPoints || sr.TrackMouches != series[i].TrackMouches || sr.TrackShots != series[i].TrackShots {
			return ErrInUse
		}
		if _, err := tx.Exec(ctx, `UPDATE programme_series SET name = $3 WHERE programme_id = $1 AND number = $2`, programmeID, i+1, sr.Name); err != nil {
			return err
		}
	}

	rows, err = tx.Query(ctx, `SELECT age_category::text FROM programme_prices WHERE programme_id = $1`, programmeID)
	if err != nil {
		return err
	}
	cats, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, c := range cats {
		if !slices.ContainsFunc(d.Prices, func(p domain.Price) bool { return string(p.AgeCategory) == c }) {
			return ErrInUse
		}
	}
	for _, p := range d.Prices {
		if _, err := tx.Exec(ctx, `INSERT INTO programme_prices (programme_id, age_category, price, repeat_price) VALUES ($1, $2, $3, $4)
			ON CONFLICT (programme_id, age_category) DO UPDATE SET price = EXCLUDED.price, repeat_price = EXCLUDED.repeat_price`,
			programmeID, p.AgeCategory, p.Price, p.RepeatPrice); err != nil {
			return err
		}
	}

	rows, err = tx.Query(ctx, `SELECT id FROM award_tiers WHERE programme_id = $1 ORDER BY position`, programmeID)
	if err != nil {
		return err
	}
	tiers, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		return err
	}
	if len(d.Awards) < len(tiers) {
		return ErrInUse
	}
	for i, a := range d.Awards {
		var tierID int64
		if i < len(tiers) {
			tierID = tiers[i]
			if _, err := tx.Exec(ctx, `UPDATE award_tiers SET name = $2, kind = $3, value = $4 WHERE id = $1`, tierID, a.Name, a.Kind, a.Value); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `DELETE FROM award_limits WHERE tier_id = $1`, tierID); err != nil {
				return err
			}
		} else if err := tx.QueryRow(ctx, `INSERT INTO award_tiers (programme_id, position, name, kind, value)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`, programmeID, i+1, a.Name, a.Kind, a.Value).Scan(&tierID); err != nil {
			return err
		}
		for _, l := range a.Limits {
			if _, err := tx.Exec(ctx, `INSERT INTO award_limits (tier_id, age_category, points) VALUES ($1, $2, $3)`,
				tierID, l.AgeCategory, l.Points); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetProgrammeActive shows or hides a programme at registration.
func (s *Store) SetProgrammeActive(ctx context.Context, tenantID, programmeID int64, active bool) error {
	return affected(s.pool.Exec(ctx, `UPDATE programmes SET active = $3 WHERE tenant_id = $1 AND id = $2`,
		tenantID, programmeID, active))
}

// SetCombination lets a programme take series from another programme of the same event, or clears it with nil.
// ErrNotFound if either programme does not exist in the tenant or they belong to different events.
func (s *Store) SetCombination(ctx context.Context, tenantID, programmeID int64, targetID *int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE programmes p SET combinable_with_id = $3
		WHERE p.tenant_id = $1 AND p.id = $2
		AND ($3::bigint IS NULL OR EXISTS (SELECT 1 FROM programmes q
			WHERE q.tenant_id = $1 AND q.id = $3 AND q.event_id = p.event_id AND q.id <> p.id))`,
		tenantID, programmeID, targetID))
}

// SetProgrammeRange ties a programme to one range of its event, or to any with nil; allowed while in use.
// ErrNotFound if the programme does not exist in the tenant or the range is not one of the event's.
func (s *Store) SetProgrammeRange(ctx context.Context, tenantID, programmeID int64, rangeID *int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE programmes p SET range_id = $3
		WHERE p.tenant_id = $1 AND p.id = $2
		AND ($3::bigint IS NULL OR EXISTS (SELECT 1 FROM event_ranges er WHERE er.event_id = p.event_id AND er.range_id = $3))`,
		tenantID, programmeID, rangeID))
}

// SetSpecialPrizes changes the special prize rules of a programme; allowed while shooters are registered.
func (s *Store) SetSpecialPrizes(ctx context.Context, tenantID, programmeID int64, r domain.SpecialPrizeRules) error {
	return affected(s.pool.Exec(ctx, `UPDATE programmes SET prize_ranks = $3, prize_min_category = $4, prize_min_club = $5
		WHERE tenant_id = $1 AND id = $2`, tenantID, programmeID, r.Ranks, r.MinCategory, r.MinClub))
}

// DeleteProgramme removes a programme and its parts of classifications; ErrInUse once shooters are registered for it.
func (s *Store) DeleteProgramme(ctx context.Context, tenantID, programmeID int64) error {
	return mapErr(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE classifications c SET definition = jsonb_set(definition, '{parts}',
				coalesce((SELECT jsonb_agg(p) FROM jsonb_array_elements(definition->'parts') p WHERE (p->>'programmeId')::bigint <> $2), '[]'))
			WHERE tenant_id = $1 AND event_id = (SELECT event_id FROM programmes WHERE id = $2)`, tenantID, programmeID); err != nil {
			return err
		}
		return affected(tx.Exec(ctx, `DELETE FROM programmes WHERE tenant_id = $1 AND id = $2`, tenantID, programmeID))
	}))
}

// CopyProgrammes copies all programmes of one event into another event of the same tenant, including
// combinations between them, special prizes and the active flag. Programmes whose name and variant already
// exist in the target are skipped.
func (s *Store) CopyProgrammes(ctx context.Context, tenantID, fromEventID, toEventID int64) (int, error) {
	copied := 0
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := eventInTenant(ctx, tx, tenantID, toEventID); err != nil {
			return err
		}
		var err error
		copied, err = copyProgrammes(ctx, tx, tenantID, fromEventID, toEventID)
		return err
	})
	return copied, mapErr(err)
}

// CopyEvent creates a new event of the tenant with the ranges, schedule (dates shifted with the main day), organizer, timezone,
// language, settings, chosen modules, booking hours (other days shifted, no deadline, booking off), all programmes and the
// data of the module copiers of another event, but no registrations; the source may be archived.
func (s *Store) CopyEvent(ctx context.Context, tenantID, fromEventID int64, name, startsOn string) (int64, error) {
	var id int64
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var sched domain.Schedule
		var shift int
		if err := tx.QueryRow(ctx, `SELECT schedule, $3::date - starts_on FROM events WHERE tenant_id = $1 AND id = $2`,
			tenantID, fromEventID, startsOn).Scan(&sched, &shift); err != nil {
			return err
		}
		sched = domain.NormalizeSchedule(sched.Shift(shift), startsOn)
		if err := tx.QueryRow(ctx, `INSERT INTO events (tenant_id, name, starts_on, days, schedule, organizer, timezone, locale, settings, modules, booking)
			SELECT tenant_id, $3, $4, $5::date[], $6, organizer, timezone, locale, settings, modules || `+eventUsed+`,
				booking - 'deadline' - 'days' || coalesce((SELECT jsonb_build_object('days', jsonb_object_agg((k::date + ($4::date - starts_on))::text, v))
					FROM jsonb_each(booking->'days') AS x(k, v)), '{}')
			FROM events e WHERE tenant_id = $1 AND id = $2
			RETURNING id`, tenantID, fromEventID, name, startsOn, sched.Dates(), sched).Scan(&id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO event_ranges (tenant_id, event_id, range_id, position)
			SELECT tenant_id, $2, range_id, position FROM event_ranges WHERE event_id = $1`, fromEventID, id); err != nil {
			return err
		}
		if _, err := copyProgrammes(ctx, tx, tenantID, fromEventID, id); err != nil {
			return err
		}
		for _, copy := range copiers {
			if err := copy(ctx, tx, tenantID, fromEventID, id); err != nil {
				return err
			}
		}
		return nil
	})
	return id, mapErr(err)
}

// copyProgrammes copies the programmes of one event into another, skipping names that exist there, with combinations,
// special prizes, active flag, range and print form; it returns how many were copied.
func copyProgrammes(ctx context.Context, tx pgx.Tx, tenantID, fromEventID, toEventID int64) (int, error) {
	src, err := loadProgrammes(ctx, tx, "tenant_id = $1 AND event_id = $2", tenantID, fromEventID)
	if err != nil {
		return 0, err
	}
	copied := 0
	newID := map[int64]int64{}
	for _, p := range src {
		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM programmes WHERE event_id = $1 AND name = $2 AND variant = $3)`,
			toEventID, p.Name, p.Variant).Scan(&exists); err != nil {
			return 0, err
		}
		if exists {
			continue
		}
		id, err := createProgramme(ctx, tx, tenantID, toEventID, p.ProgrammeDefinition)
		if err != nil {
			return 0, err
		}
		if _, err := tx.Exec(ctx, `UPDATE programmes SET active = $2, prize_ranks = $3, prize_min_category = $4, prize_min_club = $5,
				range_id = (SELECT $7::bigint WHERE EXISTS (SELECT 1 FROM event_ranges er WHERE er.event_id = $6 AND er.range_id = $7)),
				layout_id = $8
			WHERE id = $1`,
			id, p.Active, p.SpecialPrizes.Ranks, p.SpecialPrizes.MinCategory, p.SpecialPrizes.MinClub, toEventID, p.RangeID, p.LayoutID); err != nil {
			return 0, err
		}
		newID[p.ID] = id
		copied++
	}
	for _, p := range src {
		if p.CombinableWithID == nil || newID[p.ID] == 0 || newID[*p.CombinableWithID] == 0 {
			continue
		}
		if _, err := tx.Exec(ctx, `UPDATE programmes SET combinable_with_id = $2 WHERE id = $1`,
			newID[p.ID], newID[*p.CombinableWithID]); err != nil {
			return 0, err
		}
	}
	return copied, nil
}

// ListTemplates lists the global programme library.
func (s *Store) ListTemplates(ctx context.Context) ([]ProgrammeTemplate, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, definition FROM programme_templates ORDER BY lower(name), lower(variant)`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[ProgrammeTemplate])
}

// GetTemplate returns one library template.
func (s *Store) GetTemplate(ctx context.Context, id int64) (ProgrammeTemplate, error) {
	t := ProgrammeTemplate{ID: id}
	err := s.pool.QueryRow(ctx, `SELECT definition FROM programme_templates WHERE id = $1`, id).Scan(&t.Definition)
	return t, mapErr(err)
}

// SaveTemplate adds a template to the library, or replaces the one with the same name and variant.
func (s *Store) SaveTemplate(ctx context.Context, d domain.ProgrammeDefinition) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO programme_templates (name, variant, definition) VALUES ($1, $2, $3)
		ON CONFLICT (name, variant) DO UPDATE SET definition = EXCLUDED.definition, updated_at = now()
		RETURNING id`, d.Name, d.Variant, d).Scan(&id)
	return id, err
}

// DeleteTemplate removes a template from the library.
func (s *Store) DeleteTemplate(ctx context.Context, id int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM programme_templates WHERE id = $1`, id))
}
