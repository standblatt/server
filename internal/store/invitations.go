package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// Invitation is an open invitation of a club; Expired when its link no longer works.
type Invitation struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	RoleID    int64     `json:"roleId"`
	CreatedAt time.Time `json:"createdAt" doc:"When it was last sent"`
	ExpiresAt time.Time `json:"expiresAt"`
	Expired   bool      `json:"expired"`
}

// InvitationInfo is what the link of an invitation shows: the club, the address and the role.
type InvitationInfo struct {
	ID         int64
	TenantID   int64
	TenantName string
	Email      string
	RoleKey    string
	RoleName   string
	ExpiresAt  time.Time
	AcceptedAt *time.Time
}

// SaveInvitation invites an address to a tenant with a role, or renews the open invitation of that address with a
// new token, role and expiry; ErrNotFound when the role is not the tenant's.
func (s *Store) SaveInvitation(ctx context.Context, tenantID int64, email string, roleID, invitedBy int64, tokenHash string, expires time.Time) error {
	return affected(s.pool.Exec(ctx, `INSERT INTO invitations (tenant_id, email, role_id, token_hash, invited_by, expires_at)
		SELECT $1, lower($2), id, $4, $5, $6 FROM roles WHERE tenant_id = $1 AND id = $3
		ON CONFLICT (tenant_id, lower(email)) WHERE accepted_at IS NULL
		DO UPDATE SET role_id = EXCLUDED.role_id, token_hash = EXCLUDED.token_hash, invited_by = EXCLUDED.invited_by,
			created_at = now(), expires_at = EXCLUDED.expires_at`, tenantID, email, roleID, tokenHash, invitedBy, expires))
}

// ListInvitations lists the open invitations of a tenant, newest first.
func (s *Store) ListInvitations(ctx context.Context, tenantID int64) ([]Invitation, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, email, role_id, created_at, expires_at, expires_at < now() FROM invitations
		WHERE tenant_id = $1 AND accepted_at IS NULL ORDER BY created_at DESC, id DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Invitation])
}

// GetInvitation returns an open invitation of a tenant.
func (s *Store) GetInvitation(ctx context.Context, tenantID, id int64) (Invitation, error) {
	var i Invitation
	err := s.pool.QueryRow(ctx, `SELECT id, email, role_id, created_at, expires_at, expires_at < now() FROM invitations
		WHERE tenant_id = $1 AND id = $2 AND accepted_at IS NULL`, tenantID, id).
		Scan(&i.ID, &i.Email, &i.RoleID, &i.CreatedAt, &i.ExpiresAt, &i.Expired)
	return i, mapErr(err)
}

// DeleteInvitation withdraws an open invitation of a tenant.
func (s *Store) DeleteInvitation(ctx context.Context, tenantID, id int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM invitations WHERE tenant_id = $1 AND id = $2 AND accepted_at IS NULL`, tenantID, id))
}

// InvitationByToken looks up an invitation by the hash of its link token, accepted or not.
func (s *Store) InvitationByToken(ctx context.Context, tokenHash string) (InvitationInfo, error) {
	var i InvitationInfo
	err := s.pool.QueryRow(ctx, `SELECT i.id, i.tenant_id, t.name, i.email, coalesce(r.key, ''), coalesce(r.name, ''), i.expires_at, i.accepted_at
		FROM invitations i JOIN tenants t ON t.id = i.tenant_id JOIN roles r ON r.id = i.role_id WHERE i.token_hash = $1`, tokenHash).
		Scan(&i.ID, &i.TenantID, &i.TenantName, &i.Email, &i.RoleKey, &i.RoleName, &i.ExpiresAt, &i.AcceptedAt)
	return i, mapErr(err)
}

// AcceptInvitation makes the user a member with the role of the invitation and marks it accepted; a user who is
// already a member keeps their role. ErrNotFound when the invitation is accepted, expired or gone.
func (s *Store) AcceptInvitation(ctx context.Context, id, userID int64) error {
	return mapErr(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var tenantID, roleID int64
		err := tx.QueryRow(ctx, `UPDATE invitations SET accepted_at = now()
			WHERE id = $1 AND accepted_at IS NULL AND expires_at > now() RETURNING tenant_id, role_id`, id).Scan(&tenantID, &roleID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO memberships (tenant_id, user_id, role_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			tenantID, userID, roleID)
		return err
	}))
}

// DeleteOldInvitations removes invitations accepted or expired more than 30 days ago.
func (s *Store) DeleteOldInvitations(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM invitations WHERE coalesce(accepted_at, expires_at) < now() - interval '30 days'`)
	return err
}
