package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// Account is the email state of a user.
type Account struct {
	UserID        int64
	Email         string
	EmailVerified bool
}

// CreateAccount inserts a self-registered user whose login name is the email; ErrExists if it is taken.
func (s *Store) CreateAccount(ctx context.Context, email, passwordHash, locale string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO users (name, email, password_hash, locale)
		VALUES (lower($1), lower($1), $2, $3) RETURNING id`, email, passwordHash, locale).Scan(&id)
	return id, mapErr(err)
}

// GetAccount returns the email state of a user.
func (s *Store) GetAccount(ctx context.Context, userID int64) (Account, error) {
	a := Account{UserID: userID}
	err := s.pool.QueryRow(ctx, `SELECT coalesce(email, ''), email_verified_at IS NOT NULL FROM users WHERE id = $1`, userID).
		Scan(&a.Email, &a.EmailVerified)
	return a, mapErr(err)
}

// SetEmail changes the email of a user and marks it unverified; ErrExists if another account uses it.
func (s *Store) SetEmail(ctx context.Context, userID int64, email string) error {
	return affected(s.pool.Exec(ctx, `UPDATE users SET email = lower($2), email_verified_at = NULL, updated_at = now()
		WHERE id = $1`, userID, email))
}

// UserByEmail looks up a user by email.
func (s *Store) UserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `SELECT id, name, locale, created_at FROM users WHERE email = lower($1)`, email).
		Scan(&u.ID, &u.Name, &u.Locale, &u.CreatedAt)
	return u, mapErr(err)
}

// CreateUserToken stores a one-time token (by its hash) for email verification or a password reset.
// Earlier tokens of the same kind stop working.
func (s *Store) CreateUserToken(ctx context.Context, id string, userID int64, kind string, expires time.Time) error {
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `DELETE FROM user_tokens WHERE user_id = $1 AND kind = $2`, userID, kind); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `INSERT INTO user_tokens (id, user_id, kind, expires_at) VALUES ($1, $2, $3, $4)`, id, userID, kind, expires)
		return err
	})
}

// UseUserToken deletes a live token and returns its user; ErrNotFound if it is unknown, used or expired.
func (s *Store) UseUserToken(ctx context.Context, id, kind string) (int64, error) {
	var userID int64
	err := s.pool.QueryRow(ctx, `DELETE FROM user_tokens WHERE id = $1 AND kind = $2 AND expires_at > now() RETURNING user_id`,
		id, kind).Scan(&userID)
	return userID, mapErr(err)
}

// VerifyEmail marks the email of a user as verified.
func (s *Store) VerifyEmail(ctx context.Context, userID int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE users SET email_verified_at = coalesce(email_verified_at, now()) WHERE id = $1`, userID))
}

// DeleteExpiredUserTokens removes tokens past their expiry.
func (s *Store) DeleteExpiredUserTokens(ctx context.Context) (int64, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM user_tokens WHERE expires_at <= now()`)
	return tag.RowsAffected(), err
}
