package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// HelperRole is what a helper link may do at an event.
type HelperRole string

// The helper roles: scorers enter results on their lanes, the ammunition desk issues ammunition
// and assigns lanes, results clerks enter results on every lane.
const (
	HelperScorer     HelperRole = "scorer"
	HelperAmmunition HelperRole = "ammunition"
	HelperResults    HelperRole = "results"
)

// Valid reports whether r is a known helper role.
func (r HelperRole) Valid() bool {
	return r == HelperScorer || r == HelperAmmunition || r == HelperResults
}

// HelperCode is a helper link of an event as listed on the configuration page of the helpers module.
type HelperCode struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Role       HelperRole `json:"role" enum:"scorer,ammunition,results"`
	RangeID    *int64     `json:"rangeId,omitempty"`
	Lanes      []int      `json:"lanes"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
}

// Helper is the principal behind a helper link.
type Helper struct {
	ID       int64      `json:"id"`
	TenantID int64      `json:"tenantId"`
	EventID  int64      `json:"eventId"`
	Name     string     `json:"name"`
	Role     HelperRole `json:"role" enum:"scorer,ammunition,results"`
	RangeID  *int64     `json:"rangeId,omitempty"`
	Lanes    []int      `json:"lanes"`
}

// CreateHelperCode stores the hash of a new helper link; ErrNotFound if the event is not the tenant's.
func (s *Store) CreateHelperCode(ctx context.Context, tenantID, eventID int64, c HelperCode, tokenHash string) (int64, error) {
	if err := eventInTenant(ctx, s.pool, tenantID, eventID); err != nil {
		return 0, err
	}
	if c.Lanes == nil {
		c.Lanes = []int{}
	}
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO helper_codes (tenant_id, event_id, name, role, range_id, lanes, token_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		tenantID, eventID, c.Name, c.Role, c.RangeID, c.Lanes, tokenHash).Scan(&id)
	return id, mapErr(err)
}

// ListHelperCodes lists the helper links of an event.
func (s *Store) ListHelperCodes(ctx context.Context, tenantID, eventID int64) ([]HelperCode, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, role, range_id, lanes, created_at, last_used_at FROM helper_codes
		WHERE tenant_id = $1 AND event_id = $2 ORDER BY created_at`, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[HelperCode])
	if list == nil {
		list = []HelperCode{}
	}
	return list, err
}

// DeleteHelperCode revokes a helper link and returns the event it belonged to.
func (s *Store) DeleteHelperCode(ctx context.Context, tenantID, id int64) (eventID int64, err error) {
	err = s.pool.QueryRow(ctx, `DELETE FROM helper_codes WHERE tenant_id = $1 AND id = $2 RETURNING event_id`, tenantID, id).Scan(&eventID)
	return eventID, mapErr(err)
}

// HelperByHash resolves a helper link hash to its principal and records the use.
func (s *Store) HelperByHash(ctx context.Context, tokenHash string) (Helper, error) {
	var h Helper
	err := s.pool.QueryRow(ctx, `UPDATE helper_codes SET last_used_at = now() WHERE token_hash = $1
		RETURNING id, tenant_id, event_id, name, role, range_id, lanes`, tokenHash).
		Scan(&h.ID, &h.TenantID, &h.EventID, &h.Name, &h.Role, &h.RangeID, &h.Lanes)
	return h, mapErr(err)
}

// IssueAwards records that the shooter collected the awards of a registration.
func (s *Store) IssueAwards(ctx context.Context, tenantID, registrationID int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE registrations SET awards_issued_at = coalesce(awards_issued_at, now())
		WHERE tenant_id = $1 AND id = $2`, tenantID, registrationID))
}

// RevokeAwards takes the handout of the awards back, e.g. after a wrong click.
func (s *Store) RevokeAwards(ctx context.Context, tenantID, registrationID int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE registrations SET awards_issued_at = NULL WHERE tenant_id = $1 AND id = $2`, tenantID, registrationID))
}

// RevokeAmmo takes the ammunition issue back.
func (s *Store) RevokeAmmo(ctx context.Context, tenantID, registrationID int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE registrations SET ammo_issued_at = NULL WHERE tenant_id = $1 AND id = $2`, tenantID, registrationID))
}

// IssueAmmo records that the ammunition of a registration was handed out.
func (s *Store) IssueAmmo(ctx context.Context, tenantID, registrationID int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE registrations SET ammo_issued_at = coalesce(ammo_issued_at, now())
		WHERE tenant_id = $1 AND id = $2`, tenantID, registrationID))
}
