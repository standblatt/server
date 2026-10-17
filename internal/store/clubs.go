package store

import (
	"context"

	"github.com/jackc/pgx/v5"
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

// CreateClub adds a club; ErrExists if the name is taken in the tenant.
func (s *Store) CreateClub(ctx context.Context, tenantID int64, in ClubInput) (Club, error) {
	c := Club{Name: in.Name, SSVNumber: in.SSVNumber}
	err := s.pool.QueryRow(ctx, `INSERT INTO clubs (tenant_id, name, ssv_number) VALUES ($1, $2, $3) RETURNING id`,
		tenantID, in.Name, in.SSVNumber).Scan(&c.ID)
	return c, mapErr(err)
}

// UpdateClub changes a club of a tenant.
func (s *Store) UpdateClub(ctx context.Context, tenantID, clubID int64, in ClubInput) error {
	return affected(s.pool.Exec(ctx, `UPDATE clubs SET name = $3, ssv_number = $4 WHERE tenant_id = $1 AND id = $2`,
		tenantID, clubID, in.Name, in.SSVNumber))
}

// DeleteClub removes a club; its shooter assignments go with it, the shooters stay.
func (s *Store) DeleteClub(ctx context.Context, tenantID, clubID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM clubs WHERE tenant_id = $1 AND id = $2`, tenantID, clubID))
}

// ensureClub returns the ID of the tenant's club with this name, creating it if needed, and whether it was created.
func ensureClub(ctx context.Context, q querier, tenantID int64, name string) (int64, bool, error) {
	var id int64
	var created bool
	err := q.QueryRow(ctx, `INSERT INTO clubs (tenant_id, name) VALUES ($1, $2)
		ON CONFLICT (tenant_id, name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, xmax = 0`, tenantID, name).Scan(&id, &created)
	return id, created, err
}

// EnsureClub returns the club with this exact name, creating it if needed.
func (s *Store) EnsureClub(ctx context.Context, tenantID int64, name string) (int64, error) {
	id, _, err := ensureClub(ctx, s.pool, tenantID, name)
	return id, err
}
