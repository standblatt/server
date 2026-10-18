package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// PairingTTL is how long a pairing code can be confirmed; a variable so tests can let codes expire.
var PairingTTL = 10 * time.Minute

// Pairing errors of FetchPairing.
var (
	ErrPairingPending = errors.New("store: pairing not confirmed yet")
	ErrPairingExpired = errors.New("store: pairing expired")
)

// Pairing is a computer waiting to be connected to a range.
type Pairing struct {
	Name      string    `json:"name"`
	OS        string    `json:"os"`
	Version   string    `json:"version"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// CreatePairing stores a pairing under its code and the hash of the computer's poll secret; ErrExists when the
// code is taken.
func (s *Store) CreatePairing(ctx context.Context, code, secretHash string, p Pairing) error {
	_, err := s.pool.Exec(ctx, `INSERT INTO agent_pairings (code, secret_hash, name, os, version, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)`, code, secretHash, p.Name, p.OS, p.Version, p.ExpiresAt)
	return mapErr(err)
}

// OpenPairing returns the pairing of a code that is neither confirmed nor expired.
func (s *Store) OpenPairing(ctx context.Context, code string) (Pairing, error) {
	var p Pairing
	err := s.pool.QueryRow(ctx, `SELECT name, os, version, expires_at FROM agent_pairings
		WHERE code = $1 AND agent_id IS NULL AND expires_at > now()`, code).Scan(&p.Name, &p.OS, &p.Version, &p.ExpiresAt)
	return p, mapErr(err)
}

// ConfirmPairing creates a print agent named name on the range for the open pairing of code. The agent gets
// placeholderHash as token hash until the computer fetches its token.
func (s *Store) ConfirmPairing(ctx context.Context, tenantID, rangeID int64, code, name, placeholderHash string) (PrintAgent, error) {
	a := PrintAgent{Name: name, Printers: []Printer{}}
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := rangeInTenant(ctx, tx, tenantID, rangeID); err != nil {
			return err
		}
		var id int64
		if err := tx.QueryRow(ctx, `SELECT id, os, version FROM agent_pairings
			WHERE code = $1 AND agent_id IS NULL AND expires_at > now() FOR UPDATE`, code).Scan(&id, &a.OS, &a.Version); err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, `INSERT INTO print_agents (tenant_id, range_id, name, token_hash, version, os)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
			tenantID, rangeID, name, placeholderHash, a.Version, a.OS).Scan(&a.ID, &a.CreatedAt); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `UPDATE agent_pairings SET tenant_id = $2, agent_id = $3 WHERE id = $1`, id, tenantID, a.ID)
		return err
	})
	return a, mapErr(err)
}

// FetchPairing hands the agent to the computer behind secretHash once its pairing is confirmed: the agent's token
// hash becomes tokenHash, the pairing is deleted, and the club and range are returned. ErrPairingPending while it
// waits, ErrPairingExpired once the code ran out unconfirmed, ErrNotFound for an unknown or fetched pairing.
func (s *Store) FetchPairing(ctx context.Context, secretHash, tokenHash string) (club, rangeName string, err error) {
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var id int64
		var agentID *int64
		var expired bool
		if err := tx.QueryRow(ctx, `SELECT id, agent_id, expires_at <= now() FROM agent_pairings WHERE secret_hash = $1 FOR UPDATE`,
			secretHash).Scan(&id, &agentID, &expired); err != nil {
			return err
		}
		switch {
		case agentID == nil && expired:
			return ErrPairingExpired
		case agentID == nil:
			return ErrPairingPending
		}
		if err := tx.QueryRow(ctx, `WITH a AS (UPDATE print_agents SET token_hash = $2 WHERE id = $1 RETURNING tenant_id, range_id)
			SELECT t.name, r.name FROM a JOIN tenants t ON t.id = a.tenant_id JOIN ranges r ON r.id = a.range_id`,
			*agentID, tokenHash).Scan(&club, &rangeName); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `DELETE FROM agent_pairings WHERE id = $1`, id)
		return err
	})
	return club, rangeName, mapErr(err)
}

// DeleteOldPairings removes pairings a day after they expired, confirmed or not.
func (s *Store) DeleteOldPairings(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM agent_pairings WHERE expires_at < now() - interval '1 day'`)
	return err
}
