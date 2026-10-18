package store

import (
	"cmp"
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/permissions"
)

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

// Membership is a tenant the user belongs to, with the user's role and permissions there.
type Membership struct {
	TenantID    int64    `json:"tenantId"`
	TenantName  string   `json:"tenantName"`
	RoleID      int64    `json:"roleId"`
	RoleKey     string   `json:"roleKey,omitempty" required:"false" enum:"admin,board,staff,member" doc:"Set for the roles every club starts with"`
	RoleName    string   `json:"roleName" doc:"Empty for a default role that keeps its translated name"`
	Permissions []string `json:"permissions" enum:"events.view,events.operate,events.configure,reports.view,shooters.view,shooters.edit,shooters.manage,ranges.manage,forms.manage,club.manage"`
	Modules     []string `json:"modules" enum:"forms" doc:"The club modules that are on"`
	InUse       []string `json:"modulesInUse" enum:"forms" doc:"The club modules with data; they cannot be switched off"`
}

// Member is a user within a tenant with their role.
type Member struct {
	UserID      int64    `json:"userId"`
	Name        string   `json:"name"`
	Email       string   `json:"email" doc:"Empty when the account has none"`
	RoleID      int64    `json:"roleId"`
	RoleKey     string   `json:"-"`
	RoleName    string   `json:"-"`
	Permissions []string `json:"-"`
}

// RoleLabel is the role's own name or, for an unnamed default role, its key (the admin panel shows it untranslated).
func (m Member) RoleLabel() string { return cmp.Or(m.RoleName, m.RoleKey) }

// RoleLabel is the role's own name or, for an unnamed default role, its key (the admin panel shows it untranslated).
func (m Membership) RoleLabel() string { return cmp.Or(m.RoleName, m.RoleKey) }

// GetTenant returns one tenant; ErrNotFound if it does not exist.
func (s *Store) GetTenant(ctx context.Context, tenantID int64) (Tenant, error) {
	var t Tenant
	err := s.pool.QueryRow(ctx, `SELECT id, name, created_at FROM tenants WHERE id = $1`, tenantID).Scan(&t.ID, &t.Name, &t.CreatedAt)
	return t, mapErr(err)
}

// CreateTenant inserts a tenant; ErrExists if the name is taken.
func (s *Store) CreateTenant(ctx context.Context, name string) (int64, error) {
	var id int64
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `INSERT INTO tenants (name, modules) VALUES ($1, $2) RETURNING id`, name, modules.Defaults(modules.Club)).Scan(&id); err != nil {
			return err
		}
		return createDefaultRoles(ctx, tx, id)
	})
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

// SetMembership adds a user to a tenant or changes their role; ErrNotFound when the role is not the tenant's.
func (s *Store) SetMembership(ctx context.Context, tenantID, userID, roleID int64) error {
	return affected(s.pool.Exec(ctx, `INSERT INTO memberships (tenant_id, user_id, role_id)
		SELECT $1, $2, id FROM roles WHERE tenant_id = $1 AND id = $3
		ON CONFLICT (tenant_id, user_id) DO UPDATE SET role_id = EXCLUDED.role_id`, tenantID, userID, roleID))
}

// RemoveMembership removes a user from a tenant.
func (s *Store) RemoveMembership(ctx context.Context, tenantID, userID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM memberships WHERE tenant_id = $1 AND user_id = $2`, tenantID, userID))
}

// Memberships lists the tenants a user belongs to.
func (s *Store) Memberships(ctx context.Context, userID int64) ([]Membership, error) {
	rows, err := s.pool.Query(ctx, `SELECT t.id, t.name, r.id, coalesce(r.key, ''), coalesce(r.name, ''), r.permissions,
			t.modules, `+clubUsed+` FROM memberships m
		JOIN tenants t ON t.id = m.tenant_id JOIN roles r ON r.id = m.role_id WHERE m.user_id = $1 ORDER BY t.name`, userID)
	if err != nil {
		return nil, err
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Membership])
	for i := range list {
		list[i].Modules = modules.Effective(modules.Club, list[i].Modules, list[i].InUse, true)
		list[i].Permissions = permissions.Of(list[i].RoleKey, list[i].Permissions)
	}
	return list, err
}

// Members lists the users of a tenant.
func (s *Store) Members(ctx context.Context, tenantID int64) ([]Member, error) {
	rows, err := s.pool.Query(ctx, `SELECT u.id, u.name, coalesce(u.email, ''), r.id, coalesce(r.key, ''), coalesce(r.name, ''), r.permissions FROM memberships m
		JOIN users u ON u.id = m.user_id JOIN roles r ON r.id = m.role_id WHERE m.tenant_id = $1 ORDER BY u.name`, tenantID)
	if err != nil {
		return nil, err
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Member])
	for i := range list {
		list[i].Permissions = permissions.Of(list[i].RoleKey, list[i].Permissions)
	}
	return list, err
}

// AdminCount returns the number of admins of a tenant.
func (s *Store) AdminCount(ctx context.Context, tenantID int64) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `SELECT count(*) FROM memberships m JOIN roles r ON r.id = m.role_id
		WHERE m.tenant_id = $1 AND r.key = 'admin'`, tenantID).Scan(&n)
	return n, err
}
