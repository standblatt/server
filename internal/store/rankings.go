package store

import (
	"context"
	"math"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// RankingRow is a completed programme result of a shooter, before ranking.
type RankingRow struct {
	ProgrammeID int64
	EntryID     int64
	Firstname   string
	Lastname    string
	Birthday    string
	ClubName    string
	Category    domain.AgeCategory
	AwardName   string
	Series      []float64
	Mouches     int
}

// ListRankingRows returns the completed results of all programmes of an event.
func (s *Store) ListRankingRows(ctx context.Context, tenantID, eventID int64) ([]RankingRow, error) {
	rows, err := s.pool.Query(ctx, `SELECT en.programme_id, en.id, sh.firstname, sh.lastname, coalesce(sh.birthday::text, ''),
			coalesce(c.name, ''), r.age_category, coalesce(t.name, ''),
			array_agg(re.points::float8 ORDER BY re.series_number), coalesce(sum(re.mouches), 0)
		FROM entries en
		JOIN registrations r ON r.id = en.registration_id
		JOIN shooters sh ON sh.id = r.shooter_id
		LEFT JOIN clubs c ON c.id = r.club_id
		LEFT JOIN award_tiers t ON t.id = en.award_tier_id
		JOIN results re ON re.entry_id = en.id
		WHERE en.tenant_id = $1 AND en.event_id = $2 AND en.status = 'completed'
		GROUP BY en.id, sh.id, c.name, r.age_category, t.name`, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[RankingRow])
}

// KioskToken is a read-only display token for the rankings of one event.
type KioskToken struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
}

// CreateKioskToken stores the hash of a new kiosk token for an event; ErrNotFound if the event is not the tenant's.
func (s *Store) CreateKioskToken(ctx context.Context, tenantID, eventID int64, name, tokenHash string) (int64, error) {
	if err := eventInTenant(ctx, s.pool, tenantID, eventID); err != nil {
		return 0, err
	}
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO kiosk_tokens (tenant_id, event_id, name, token_hash) VALUES ($1, $2, $3, $4) RETURNING id`,
		tenantID, eventID, name, tokenHash).Scan(&id)
	return id, mapErr(err)
}

// ListKioskTokens lists the kiosk tokens of an event.
func (s *Store) ListKioskTokens(ctx context.Context, tenantID, eventID int64) ([]KioskToken, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, created_at, last_used_at FROM kiosk_tokens
		WHERE tenant_id = $1 AND event_id = $2 ORDER BY created_at`, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[KioskToken])
}

// DeleteKioskToken revokes a kiosk token.
func (s *Store) DeleteKioskToken(ctx context.Context, tenantID, tokenID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM kiosk_tokens WHERE tenant_id = $1 AND id = $2`, tenantID, tokenID))
}

// KioskEvent resolves a kiosk token hash to its tenant and event and records the use.
func (s *Store) KioskEvent(ctx context.Context, tokenHash string) (tenantID, eventID int64, err error) {
	err = s.pool.QueryRow(ctx, `UPDATE kiosk_tokens SET last_used_at = now()
		WHERE token_hash = $1 RETURNING tenant_id, event_id`, tokenHash).Scan(&tenantID, &eventID)
	return tenantID, eventID, mapErr(err)
}

// Placing is an entry's place in its programme's ranking.
type Placing struct {
	Rank         int
	Participants int
}

// Placings ranks every programme of an event with domain.Rank and returns the place of each completed entry.
func (s *Store) Placings(ctx context.Context, tenantID, eventID int64) (map[int64]Placing, error) {
	rows, err := s.ListRankingRows(ctx, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	byProgramme := map[int64][]RankingRow{}
	for _, r := range rows {
		byProgramme[r.ProgrammeID] = append(byProgramme[r.ProgrammeID], r)
	}
	out := map[int64]Placing{}
	for _, list := range byProgramme {
		entries := make([]domain.RankingEntry, len(list))
		for i, r := range list {
			total := 0.0
			for _, v := range r.Series {
				total += v
			}
			entries[i] = domain.RankingEntry{Total: math.Round(total*100) / 100, Mouches: r.Mouches, Series: r.Series, Birthday: r.Birthday, Category: r.Category}
		}
		for rank, i := range domain.Rank(entries) {
			out[list[i].EntryID] = Placing{Rank: rank + 1, Participants: len(list)}
		}
	}
	return out, nil
}
