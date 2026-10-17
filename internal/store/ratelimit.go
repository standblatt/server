package store

import (
	"context"
	"time"
)

// RateHit counts one attempt for key in a fixed window and returns the count including it.
// The upsert is atomic, so the limit holds across replicas.
func (s *Store) RateHit(ctx context.Context, key string, window time.Duration) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx, `INSERT INTO rate_limits (key, count, window_start) VALUES ($1, 1, now())
		ON CONFLICT (key) DO UPDATE SET
			count = CASE WHEN rate_limits.window_start < now() - $2::interval THEN 1 ELSE rate_limits.count + 1 END,
			window_start = CASE WHEN rate_limits.window_start < now() - $2::interval THEN now() ELSE rate_limits.window_start END
		RETURNING count`, key, window.String()).Scan(&count)
	return count, err
}

// RateCount returns the attempts for key in the current window without counting a new one.
func (s *Store) RateCount(ctx context.Context, key string, window time.Duration) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx, `SELECT coalesce((SELECT count FROM rate_limits
		WHERE key = $1 AND window_start >= now() - $2::interval), 0)`, key, window.String()).Scan(&count)
	return count, err
}

// DeleteOldRateLimits drops counters older than a day.
func (s *Store) DeleteOldRateLimits(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM rate_limits WHERE window_start < now() - interval '1 day'`)
	return err
}
