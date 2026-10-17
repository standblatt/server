package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// APIToken is a personal API token as the owner sees it; the secret itself is never stored.
type APIToken struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Scopes     []string   `json:"scopes"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

// TokenAuth is a token found by its hash, with its owner.
type TokenAuth struct {
	ID        int64
	Scopes    []string
	ExpiresAt *time.Time
	User      User
}

// ListAPITokens returns the tokens of a user, newest first.
func (s *Store) ListAPITokens(ctx context.Context, userID int64) ([]APIToken, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, scopes, expires_at, last_used_at, created_at
		FROM api_tokens WHERE user_id = $1 ORDER BY created_at DESC, id DESC`, userID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[APIToken])
}

// CreateAPIToken stores a token by the sha256 hash of its secret.
func (s *Store) CreateAPIToken(ctx context.Context, userID int64, name, hash string, scopes []string, expiresAt *time.Time) (APIToken, error) {
	t := APIToken{Name: name, Scopes: scopes, ExpiresAt: expiresAt}
	err := s.pool.QueryRow(ctx, `INSERT INTO api_tokens (user_id, name, token_hash, scopes, expires_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`, userID, name, hash, scopes, expiresAt).Scan(&t.ID, &t.CreatedAt)
	return t, err
}

// DeleteAPIToken revokes a token of the user.
func (s *Store) DeleteAPIToken(ctx context.Context, userID, id int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM api_tokens WHERE user_id = $1 AND id = $2`, userID, id))
}

// APITokenByHash finds a token and its owner; ErrNotFound for unknown or revoked tokens.
func (s *Store) APITokenByHash(ctx context.Context, hash string) (TokenAuth, error) {
	var a TokenAuth
	err := s.pool.QueryRow(ctx, `SELECT t.id, t.scopes, t.expires_at, u.id, u.name, u.locale, u.created_at
		FROM api_tokens t JOIN users u ON u.id = t.user_id WHERE t.token_hash = $1`, hash).
		Scan(&a.ID, &a.Scopes, &a.ExpiresAt, &a.User.ID, &a.User.Name, &a.User.Locale, &a.User.CreatedAt)
	return a, mapErr(err)
}

// TouchAPIToken records that a token was used.
func (s *Store) TouchAPIToken(ctx context.Context, id int64) error {
	_, err := s.pool.Exec(ctx, `UPDATE api_tokens SET last_used_at = now() WHERE id = $1`, id)
	return err
}
