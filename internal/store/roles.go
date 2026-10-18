package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/permissions"
)

// Keys of the roles every club starts with, for code that assigns one of them.
const (
	RoleAdmin  = permissions.Admin
	RoleBoard  = permissions.Board
	RoleStaff  = permissions.Staff
	RoleMember = permissions.Member
)

// Role is a role of a club: one of the defaults (Key set) or one the club made. The admin role has every permission.
type Role struct {
	ID          int64    `json:"id"`
	Key         string   `json:"key,omitempty" required:"false" enum:"admin,board,staff,member" doc:"Set for the roles every club starts with"`
	Name        string   `json:"name" doc:"Empty for a default role that keeps its translated name"`
	Permissions []string `json:"permissions" enum:"events.view,events.operate,events.configure,reports.view,shooters.view,shooters.edit,shooters.manage,ranges.manage,forms.manage,club.manage"`
	Members     int      `json:"members"`
}

// RoleInput holds the editable fields of a role.
type RoleInput struct {
	Name        string   `json:"name" maxLength:"60" doc:"Empty keeps the translated name of a default role"`
	Permissions []string `json:"permissions" enum:"events.view,events.operate,events.configure,reports.view,shooters.view,shooters.edit,shooters.manage,ranges.manage,forms.manage,club.manage"`
}

// createDefaultRoles adds the roles every club starts with.
func createDefaultRoles(ctx context.Context, q querier, tenantID int64) error {
	for i, d := range permissions.Defaults {
		perms := d.Permissions
		if perms == nil {
			perms = []string{}
		}
		if _, err := q.Exec(ctx, `INSERT INTO roles (tenant_id, key, permissions, position) VALUES ($1, $2, $3, $4)`,
			tenantID, d.Key, perms, i+1); err != nil {
			return err
		}
	}
	return nil
}

// ListRoles lists the roles of a tenant in their order, with their member counts.
func (s *Store) ListRoles(ctx context.Context, tenantID int64) ([]Role, error) {
	rows, err := s.pool.Query(ctx, `SELECT r.id, coalesce(r.key, ''), coalesce(r.name, ''), r.permissions,
			(SELECT count(*) FROM memberships m WHERE m.role_id = r.id)
		FROM roles r WHERE r.tenant_id = $1 ORDER BY r.position, r.id`, tenantID)
	if err != nil {
		return nil, err
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Role])
	for i := range list {
		list[i].Permissions = permissions.Of(list[i].Key, list[i].Permissions)
	}
	return list, err
}

// GetRole returns a role of a tenant.
func (s *Store) GetRole(ctx context.Context, tenantID, roleID int64) (Role, error) {
	list, err := s.ListRoles(ctx, tenantID)
	if err != nil {
		return Role{}, err
	}
	for _, r := range list {
		if r.ID == roleID {
			return r, nil
		}
	}
	return Role{}, ErrNotFound
}

// RoleIDByKey returns the ID of a default role of a tenant.
func (s *Store) RoleIDByKey(ctx context.Context, tenantID int64, key string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `SELECT id FROM roles WHERE tenant_id = $1 AND key = $2`, tenantID, key).Scan(&id)
	return id, mapErr(err)
}

// CreateRole adds a role at the end of a tenant's list.
func (s *Store) CreateRole(ctx context.Context, tenantID int64, in RoleInput) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO roles (tenant_id, name, permissions, position)
		VALUES ($1, $2, $3, (SELECT coalesce(max(position), 0) + 1 FROM roles WHERE tenant_id = $1)) RETURNING id`,
		tenantID, in.Name, permissions.Of("", in.Permissions)).Scan(&id)
	return id, mapErr(err)
}

// UpdateRole renames a role and replaces its permissions; the admin role keeps every permission and only takes a name.
func (s *Store) UpdateRole(ctx context.Context, tenantID, roleID int64, in RoleInput) error {
	return affected(s.pool.Exec(ctx, `UPDATE roles SET name = nullif($3, ''),
			permissions = CASE WHEN key = 'admin' THEN '{}' ELSE $4::text[] END
		WHERE tenant_id = $1 AND id = $2`, tenantID, roleID, in.Name, permissions.Of("", in.Permissions)))
}

// DeleteRole removes a role; ErrInUse while members have it, ErrNotFound for the admin role.
func (s *Store) DeleteRole(ctx context.Context, tenantID, roleID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM roles WHERE tenant_id = $1 AND id = $2 AND key IS DISTINCT FROM 'admin'`, tenantID, roleID))
}

// MemberRole returns the user's role in a tenant; ErrNotFound if they are not a member.
func (s *Store) MemberRole(ctx context.Context, tenantID, userID int64) (Role, error) {
	var r Role
	err := s.pool.QueryRow(ctx, `SELECT r.id, coalesce(r.key, ''), coalesce(r.name, ''), r.permissions FROM memberships m
		JOIN roles r ON r.id = m.role_id WHERE m.tenant_id = $1 AND m.user_id = $2`, tenantID, userID).
		Scan(&r.ID, &r.Key, &r.Name, &r.Permissions)
	r.Permissions = permissions.Of(r.Key, r.Permissions)
	return r, mapErr(err)
}
