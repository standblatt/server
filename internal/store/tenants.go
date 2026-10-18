package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/modules"
)

// Role is a user's role within one tenant.
type Role string

// Tenant roles: admins manage the club's settings and members, staff run competitions.
const (
	RoleAdmin Role = "admin"
	RoleStaff Role = "staff"
	// RoleHelper is what access reports for a helper link; it is never stored in memberships.
	RoleHelper Role = "helper"
)

// Valid reports whether r is a known role.
func (r Role) Valid() bool { return r == RoleAdmin || r == RoleStaff }

// Tenant is a club organising competitions.
type Tenant struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Members   int       `json:"-"`
	// Events and Registrations are counts for the admin panel.
	Events        int `json:"-"`
	Registrations int `json:"-"`
}

// Membership is a tenant the user belongs to, with the user's role there.
type Membership struct {
	TenantID   int64    `json:"tenantId"`
	TenantName string   `json:"tenantName"`
	Role       Role     `json:"role" enum:"admin,staff"`
	Modules    []string `json:"modules" enum:"forms" doc:"The club modules that are on"`
	InUse      []string `json:"modulesInUse" enum:"forms" doc:"The club modules with data; they cannot be switched off"`
}

// Member is a user within a tenant.
type Member struct {
	UserID int64  `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email" doc:"Empty when the account has none"`
	Role   Role   `json:"role" enum:"admin,staff"`
}

// GetTenant returns one tenant; ErrNotFound if it does not exist.
func (s *Store) GetTenant(ctx context.Context, tenantID int64) (Tenant, error) {
	var t Tenant
	err := s.pool.QueryRow(ctx, `SELECT id, name, created_at FROM tenants WHERE id = $1`, tenantID).Scan(&t.ID, &t.Name, &t.CreatedAt)
	return t, mapErr(err)
}

// CreateTenant inserts a tenant; ErrExists if the name is taken.
func (s *Store) CreateTenant(ctx context.Context, name string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO tenants (name, modules) VALUES ($1, $2) RETURNING id`, name, modules.Defaults(modules.Club)).Scan(&id)
	return id, mapErr(err)
}

// ListTenants lists all tenants with their member count.
func (s *Store) ListTenants(ctx context.Context) ([]Tenant, error) {
	rows, err := s.pool.Query(ctx, `SELECT t.id, t.name, t.created_at,
			(SELECT count(*) FROM memberships m WHERE m.tenant_id = t.id),
			(SELECT count(*) FROM events e WHERE e.tenant_id = t.id),
			(SELECT count(*) FROM registrations r WHERE r.tenant_id = t.id)
		FROM tenants t ORDER BY t.name`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Tenant, error) {
		var t Tenant
		err := row.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.Members, &t.Events, &t.Registrations)
		return t, err
	})
}

// DeleteTenant removes a tenant with all its data.
func (s *Store) DeleteTenant(ctx context.Context, tenantID int64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for _, q := range []string{
		`UPDATE events SET archived_at = NULL WHERE tenant_id = $1 AND archived_at IS NOT NULL`,
		`DELETE FROM registrations WHERE tenant_id = $1`,
		`DELETE FROM events WHERE tenant_id = $1`,
	} {
		if _, err := tx.Exec(ctx, q, tenantID); err != nil {
			return mapErr(err)
		}
	}
	if err := affected(tx.Exec(ctx, `DELETE FROM tenants WHERE id = $1`, tenantID)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// SetMembership adds a user to a tenant or changes their role.
func (s *Store) SetMembership(ctx context.Context, tenantID, userID int64, role Role) error {
	_, err := s.pool.Exec(ctx, `INSERT INTO memberships (tenant_id, user_id, role) VALUES ($1, $2, $3)
		ON CONFLICT (tenant_id, user_id) DO UPDATE SET role = EXCLUDED.role`, tenantID, userID, role)
	return mapErr(err)
}

// RemoveMembership removes a user from a tenant.
func (s *Store) RemoveMembership(ctx context.Context, tenantID, userID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM memberships WHERE tenant_id = $1 AND user_id = $2`, tenantID, userID))
}

// Memberships lists the tenants a user belongs to.
func (s *Store) Memberships(ctx context.Context, userID int64) ([]Membership, error) {
	rows, err := s.pool.Query(ctx, `SELECT t.id, t.name, m.role, t.modules, `+clubUsed+` FROM memberships m
		JOIN tenants t ON t.id = m.tenant_id WHERE m.user_id = $1 ORDER BY t.name`, userID)
	if err != nil {
		return nil, err
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Membership])
	for i := range list {
		list[i].Modules = modules.Effective(modules.Club, list[i].Modules, list[i].InUse, true)
	}
	return list, err
}

// MemberRole returns the user's role in a tenant; ErrNotFound if they are not a member.
func (s *Store) MemberRole(ctx context.Context, tenantID, userID int64) (Role, error) {
	var r Role
	err := s.pool.QueryRow(ctx, `SELECT role FROM memberships WHERE tenant_id = $1 AND user_id = $2`,
		tenantID, userID).Scan(&r)
	return r, mapErr(err)
}

// Members lists the users of a tenant.
func (s *Store) Members(ctx context.Context, tenantID int64) ([]Member, error) {
	rows, err := s.pool.Query(ctx, `SELECT u.id, u.name, coalesce(u.email, ''), m.role FROM memberships m
		JOIN users u ON u.id = m.user_id WHERE m.tenant_id = $1 ORDER BY u.name`, tenantID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Member])
}

// AdminCount returns the number of admins of a tenant.
func (s *Store) AdminCount(ctx context.Context, tenantID int64) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `SELECT count(*) FROM memberships WHERE tenant_id = $1 AND role = 'admin'`,
		tenantID).Scan(&n)
	return n, err
}
