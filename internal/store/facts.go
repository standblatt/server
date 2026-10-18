package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// EntryFact is one attempt at a programme a shooter registered for at an event, with everything the reports need.
type EntryFact struct {
	EntryID            int64
	ProgrammeID        int64
	Attempt            int
	RegistrationID     int64
	RegistrationStatus string
	RegisteredAt       time.Time
	EntryStatus        string
	ShooterID          int64
	Firstname          string
	Lastname           string
	Category           domain.AgeCategory
	ClubID             *int64
	ClubName           string
	ClubSSVNumber      string
	Price              float64
	AwardTierID        *int64
	AwardName          string
	AwardKind          string
	AwardValue         float64
	Total              float64
	Mouches            int
	FirstTwoSeries     float64
	AwardsIssuedAt     *time.Time
}

// ListEntryFacts returns every entry of an event with price, award and result totals.
func (s *Store) ListEntryFacts(ctx context.Context, tenantID, eventID int64) ([]EntryFact, error) {
	rows, err := s.pool.Query(ctx, `SELECT en.id, en.programme_id, en.attempt, r.id, r.status::text, r.registered_at, en.status::text, sh.id, sh.firstname, sh.lastname,
			r.age_category, r.club_id, coalesce(c.name, ''), coalesce(c.ssv_number, ''),
			coalesce(CASE WHEN en.attempt > 1 THEN coalesce(pp.repeat_price, pp.price) ELSE pp.price END, 0)::float8, en.award_tier_id, coalesce(t.name, ''), coalesce(t.kind::text, ''), coalesce(t.value, 0)::float8,
			coalesce(sum(re.points), 0)::float8, coalesce(sum(re.mouches), 0),
			coalesce(sum(re.points) FILTER (WHERE re.series_number <= 2), 0)::float8, r.awards_issued_at
		FROM entries en
		JOIN registrations r ON r.id = en.registration_id
		JOIN shooters sh ON sh.id = r.shooter_id
		LEFT JOIN clubs c ON c.id = r.club_id
		LEFT JOIN programme_prices pp ON pp.programme_id = en.programme_id AND pp.age_category = r.age_category
		LEFT JOIN award_tiers t ON t.id = en.award_tier_id
		LEFT JOIN results re ON re.entry_id = en.id
		WHERE en.tenant_id = $1 AND en.event_id = $2
		GROUP BY en.id, r.id, sh.id, c.id, pp.price, pp.repeat_price, t.id
		ORDER BY sh.lastname, sh.firstname, en.id`, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[EntryFact])
}
