package store

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// User is an account. Accounts are global; access to tenants comes from memberships.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Locale    string    `json:"locale" doc:"UI language; empty means the browser's"`
	CreatedAt time.Time `json:"createdAt"`
}

// UserStat is a user with session and membership counts for the admin panel.
type UserStat struct {
	User
	Sessions      int
	LastLogin     *time.Time
	Tenants       []string
	Email         string
	EmailVerified bool
	Trainings     int
	TOTP          bool
}

// CreateUser inserts a user; ErrExists if the name is taken.
func (s *Store) CreateUser(ctx context.Context, name, passwordHash string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO users (name, password_hash) VALUES (lower($1), $2) RETURNING id`,
		name, passwordHash).Scan(&id)
	return id, mapErr(err)
}

// UpsertUser creates a user or sets a new password for an existing one.
func (s *Store) UpsertUser(ctx context.Context, name, passwordHash string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO users (name, password_hash) VALUES (lower($1), $2)
		ON CONFLICT (name) DO UPDATE SET password_hash = EXCLUDED.password_hash, updated_at = now()
		RETURNING id`, name, passwordHash).Scan(&id)
	return id, mapErr(err)
}

// UserByName looks up a user by name.
func (s *Store) UserByName(ctx context.Context, name string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `SELECT id, name, locale, created_at FROM users WHERE name = lower($1)`, name).
		Scan(&u.ID, &u.Name, &u.Locale, &u.CreatedAt)
	return u, mapErr(err)
}

// Credentials returns the user and password hash for a login attempt by user name or email.
func (s *Store) Credentials(ctx context.Context, name string) (User, string, error) {
	var u User
	var hash string
	err := s.pool.QueryRow(ctx, `SELECT id, name, locale, created_at, password_hash FROM users
		WHERE name = lower($1) OR email = lower($1) ORDER BY name = lower($1) DESC LIMIT 1`, name).
		Scan(&u.ID, &u.Name, &u.Locale, &u.CreatedAt, &hash)
	return u, hash, mapErr(err)
}

// SetPassword replaces the password hash of a user.
func (s *Store) SetPassword(ctx context.Context, userID int64, passwordHash string) error {
	return affected(s.pool.Exec(ctx, `UPDATE users SET password_hash = $2, updated_at = now() WHERE id = $1`,
		userID, passwordHash))
}

// SetUserLocale stores the UI language of a user.
func (s *Store) SetUserLocale(ctx context.Context, userID int64, locale string) error {
	return affected(s.pool.Exec(ctx, `UPDATE users SET locale = $2, updated_at = now() WHERE id = $1`, userID, locale))
}

// DeleteUser removes a user with their sessions and memberships.
func (s *Store) DeleteUser(ctx context.Context, userID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID))
}

// userStatCols selects the columns of UserStat from users u.
const userStatCols = `u.id, u.name, u.locale, u.created_at,
	(SELECT count(*) FROM sessions s WHERE s.user_id = u.id AND s.expires_at > now()),
	(SELECT max(created_at) FROM sessions s WHERE s.user_id = u.id),
	coalesce((SELECT array_agg(t.name || ' (' || m.role || ')' ORDER BY t.name)
		FROM memberships m JOIN tenants t ON t.id = m.tenant_id WHERE m.user_id = u.id), '{}'),
	coalesce(u.email, ''), u.email_verified_at IS NOT NULL,
	(SELECT count(*) FROM training_sessions ts WHERE ts.user_id = u.id), u.totp_secret IS NOT NULL`

// scanUserStat reads one UserStat row.
func scanUserStat(row pgx.Row) (UserStat, error) {
	var u UserStat
	err := row.Scan(&u.ID, &u.Name, &u.Locale, &u.CreatedAt, &u.Sessions, &u.LastLogin, &u.Tenants,
		&u.Email, &u.EmailVerified, &u.Trainings, &u.TOTP)
	return u, err
}

// ListUsersWithStats lists users with active sessions, last login and tenant names; q filters by a part of the name or email.
func (s *Store) ListUsersWithStats(ctx context.Context, q string) ([]UserStat, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+userStatCols+` FROM users u
		WHERE $1 = '' OR u.name ILIKE '%' || $1 || '%' OR u.email ILIKE '%' || $1 || '%'
		ORDER BY u.name LIMIT 500`, strings.TrimSpace(q))
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (UserStat, error) { return scanUserStat(row) })
}

// UserStats returns one user with the counts of the admin panel.
func (s *Store) UserStats(ctx context.Context, userID int64) (UserStat, error) {
	u, err := scanUserStat(s.pool.QueryRow(ctx, `SELECT `+userStatCols+` FROM users u WHERE u.id = $1`, userID))
	return u, mapErr(err)
}

// SetUserEmail sets or clears (empty) the email of a user from the admin panel, verified or not; ErrExists if another account uses it.
func (s *Store) SetUserEmail(ctx context.Context, userID int64, email string, verified bool) error {
	return affected(s.pool.Exec(ctx, `UPDATE users SET email = nullif(lower(trim($2)), ''),
		email_verified_at = CASE WHEN $3 AND nullif(trim($2), '') IS NOT NULL THEN now() END, updated_at = now()
		WHERE id = $1`, userID, email, verified))
}

// AdminStats are the platform counts on the admin overview.
type AdminStats struct {
	Tenants, Users, Events, Registrations, Sessions, Templates int
}

// AdminStats counts clubs, users, events, registrations, live sessions and library templates.
func (s *Store) AdminStats(ctx context.Context) (AdminStats, error) {
	var a AdminStats
	err := s.pool.QueryRow(ctx, `SELECT (SELECT count(*) FROM tenants), (SELECT count(*) FROM users), (SELECT count(*) FROM events),
		(SELECT count(*) FROM registrations), (SELECT count(*) FROM sessions WHERE expires_at > now()), (SELECT count(*) FROM programme_templates)`).
		Scan(&a.Tenants, &a.Users, &a.Events, &a.Registrations, &a.Sessions, &a.Templates)
	return a, err
}

// CreateSession stores a session by the hash of its token; it ends after idle without a request, and after ttl at the latest.
func (s *Store) CreateSession(ctx context.Context, id string, userID int64, idle, ttl time.Duration) error {
	_, err := s.pool.Exec(ctx, `INSERT INTO sessions (id, user_id, expires_at, idle, ttl) VALUES ($1, $2, now() + $3::interval, $3::interval, $4::interval)`,
		id, userID, idle.String(), ttl.String())
	return err
}

// SessionUser returns the user of a live session and slides its expiry by its idle time, capped at its ttl after creation.
func (s *Store) SessionUser(ctx context.Context, id string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `UPDATE sessions s
		SET expires_at = least(now() + s.idle, s.created_at + s.ttl)
		FROM users u
		WHERE s.id = $1 AND s.expires_at > now() AND u.id = s.user_id
		RETURNING u.id, u.name, u.locale, u.created_at`, id).Scan(&u.ID, &u.Name, &u.Locale, &u.CreatedAt)
	return u, mapErr(err)
}

// DeleteSession ends one session.
func (s *Store) DeleteSession(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}

// DeleteOtherSessions ends every session of a user except the one with id keep (after a password change).
func (s *Store) DeleteOtherSessions(ctx context.Context, userID int64, keep string) (int64, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1 AND id <> $2`, userID, keep)
	return tag.RowsAffected(), err
}

// DeleteUserSessions ends all sessions of a user.
func (s *Store) DeleteUserSessions(ctx context.Context, userID int64) (int64, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return tag.RowsAffected(), err
}

// DeleteAllSessions logs everyone out.
func (s *Store) DeleteAllSessions(ctx context.Context) (int64, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM sessions`)
	return tag.RowsAffected(), err
}

// DeleteExpiredSessions removes sessions past their expiry.
func (s *Store) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at <= now()`)
	return tag.RowsAffected(), err
}

// TOTP returns the authenticator secret of a user ("" when off) and the last time step a code was accepted for.
func (s *Store) TOTP(ctx context.Context, userID int64) (secret string, counter int64, err error) {
	var sec *string
	err = s.pool.QueryRow(ctx, `SELECT totp_secret, totp_counter FROM users WHERE id = $1`, userID).Scan(&sec, &counter)
	if sec != nil {
		secret = *sec
	}
	return secret, counter, mapErr(err)
}

// SetTOTP enables the authenticator with a secret, or disables it with ""; the replay counter starts over.
func (s *Store) SetTOTP(ctx context.Context, userID int64, secret string) error {
	var sec *string
	if secret != "" {
		sec = &secret
	}
	return affected(s.pool.Exec(ctx, `UPDATE users SET totp_secret = $2, totp_counter = 0, updated_at = now() WHERE id = $1`, userID, sec))
}

// UseTOTP records the time step of an accepted code; ErrNotFound when that step (or a later one) was used already.
func (s *Store) UseTOTP(ctx context.Context, userID, counter int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE users SET totp_counter = $2 WHERE id = $1 AND totp_counter < $2`, userID, counter))
}
