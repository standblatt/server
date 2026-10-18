package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// Club is a shooting club in a tenant's master data.
type Club struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	SSVNumber    string `json:"ssvNumber"`
	ShooterCount int    `json:"shooterCount"`
}

// ClubInput holds the editable fields of a club.
type ClubInput struct {
	Name      string `json:"name" minLength:"1" maxLength:"200"`
	SSVNumber string `json:"ssvNumber" maxLength:"40" required:"false"`
}

// ClubRef is a club as listed on a shooter.
type ClubRef struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ListClubs lists the clubs of a tenant with the number of shooters assigned to each.
func (s *Store) ListClubs(ctx context.Context, tenantID int64) ([]Club, error) {
	rows, err := s.pool.Query(ctx, `SELECT c.id, c.name, c.ssv_number,
			(SELECT count(*) FROM shooter_clubs sc WHERE sc.club_id = c.id)
		FROM clubs c WHERE c.tenant_id = $1 ORDER BY lower(c.name)`, tenantID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Club])
}

// GetClub returns a club of a tenant.
func (s *Store) GetClub(ctx context.Context, tenantID, id int64) (Club, error) {
	var c Club
	err := s.pool.QueryRow(ctx, `SELECT c.id, c.name, c.ssv_number, (SELECT count(*) FROM shooter_clubs sc WHERE sc.club_id = c.id)
		FROM clubs c WHERE c.tenant_id = $1 AND c.id = $2`, tenantID, id).Scan(&c.ID, &c.Name, &c.SSVNumber, &c.ShooterCount)
	return c, mapErr(err)
}

// clubByKey returns the ID of the tenant's oldest club other than except whose name has the key of name, or 0.
func clubByKey(ctx context.Context, q querier, tenantID int64, name string, except int64) (int64, error) {
	rows, err := q.Query(ctx, `SELECT id, name FROM clubs WHERE tenant_id = $1 AND id <> $2 ORDER BY id`, tenantID, except)
	if err != nil {
		return 0, err
	}
	key := domain.ClubKey(name)
	var id, found int64
	var n string
	_, err = pgx.ForEachRow(rows, []any{&id, &n}, func() error {
		if found == 0 && domain.ClubKey(n) == key {
			found = id
		}
		return nil
	})
	return found, err
}

// CreateClub adds a club; ErrExists if a club of the tenant has the same name regardless of case and spaces.
func (s *Store) CreateClub(ctx context.Context, tenantID int64, in ClubInput) (Club, error) {
	c := Club{Name: in.Name, SSVNumber: in.SSVNumber}
	if id, err := clubByKey(ctx, s.pool, tenantID, in.Name, 0); err != nil || id != 0 {
		return c, existsUnless(err)
	}
	err := s.pool.QueryRow(ctx, `INSERT INTO clubs (tenant_id, name, ssv_number) VALUES ($1, $2, $3) RETURNING id`,
		tenantID, in.Name, in.SSVNumber).Scan(&c.ID)
	return c, mapErr(err)
}

// UpdateClub changes a club of a tenant; ErrExists if another club has the same name regardless of case and spaces.
func (s *Store) UpdateClub(ctx context.Context, tenantID, clubID int64, in ClubInput) error {
	if id, err := clubByKey(ctx, s.pool, tenantID, in.Name, clubID); err != nil || id != 0 {
		return existsUnless(err)
	}
	return affected(s.pool.Exec(ctx, `UPDATE clubs SET name = $3, ssv_number = $4 WHERE tenant_id = $1 AND id = $2`,
		tenantID, clubID, in.Name, in.SSVNumber))
}

// DeleteClub removes a club; its shooter assignments go with it, the shooters stay.
func (s *Store) DeleteClub(ctx context.Context, tenantID, clubID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM clubs WHERE tenant_id = $1 AND id = $2`, tenantID, clubID))
}

// existsUnless returns err, or ErrExists when there is none.
func existsUnless(err error) error {
	if err != nil {
		return err
	}
	return ErrExists
}

// ensureClub returns the ID of the tenant's club with this name regardless of case and spaces, creating it if
// needed, and whether it was created.
func ensureClub(ctx context.Context, q querier, tenantID int64, name string) (int64, bool, error) {
	if id, err := clubByKey(ctx, q, tenantID, name, 0); err != nil || id != 0 {
		return id, false, err
	}
	var id int64
	var created bool
	err := q.QueryRow(ctx, `INSERT INTO clubs (tenant_id, name) VALUES ($1, $2)
		ON CONFLICT (tenant_id, name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, xmax = 0`, tenantID, name).Scan(&id, &created)
	return id, created, err
}

// EnsureClub returns the club with this name regardless of case and spaces, creating it if needed.
func (s *Store) EnsureClub(ctx context.Context, tenantID int64, name string) (int64, error) {
	id, _, err := ensureClub(ctx, s.pool, tenantID, name)
	return id, err
}

// MergeClubs moves everything of club from to club into and deletes from: shooter assignments, registrations (of
// closed events too), booking requests and the club excluded from an event's club competition. into keeps its name
// and takes the SSV number of from when it has none. ErrNotFound unless both are distinct clubs of the tenant.
func (s *Store) MergeClubs(ctx context.Context, tenantID, from, into int64) error {
	if from == into {
		return ErrNotFound
	}
	return mapErr(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var n int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM clubs WHERE tenant_id = $1 AND id IN ($2, $3)`, tenantID, from, into).Scan(&n); err != nil {
			return err
		}
		if n != 2 {
			return ErrNotFound
		}
		if _, err := tx.Exec(ctx, `SET LOCAL standblatt.merge = 'on'`); err != nil {
			return err
		}
		for _, sql := range []string{
			`INSERT INTO shooter_clubs (tenant_id, shooter_id, club_id) SELECT tenant_id, shooter_id, $3 FROM shooter_clubs
				WHERE tenant_id = $1 AND club_id = $2 ON CONFLICT DO NOTHING`,
			`UPDATE registrations SET club_id = $3 WHERE tenant_id = $1 AND club_id = $2`,
			`UPDATE booking_requests SET club_id = $3 WHERE tenant_id = $1 AND club_id = $2`,
			`UPDATE events SET settings = jsonb_set(settings, '{excludedClubId}', to_jsonb($3::bigint))
				WHERE tenant_id = $1 AND (settings->>'excludedClubId')::bigint = $2`,
			`UPDATE clubs SET ssv_number = (SELECT ssv_number FROM clubs WHERE id = $2) WHERE tenant_id = $1 AND id = $3 AND ssv_number = ''`,
		} {
			if _, err := tx.Exec(ctx, sql, tenantID, from, into); err != nil {
				return err
			}
		}
		_, err := tx.Exec(ctx, `DELETE FROM clubs WHERE tenant_id = $1 AND id = $2`, tenantID, from)
		return err
	}))
}
