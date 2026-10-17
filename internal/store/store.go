// Package store is the PostgreSQL persistence layer, including the embedded schema migrations.
package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store wraps the connection pool and provides all database access.
type Store struct {
	pool *pgxpool.Pool
}

// New opens a connection pool for the given PostgreSQL URL.
func New(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return &Store{pool: pool}, nil
}

// Close releases all pooled connections.
func (s *Store) Close() { s.pool.Close() }

// Ping checks that the database is reachable.
func (s *Store) Ping(ctx context.Context) error { return s.pool.Ping(ctx) }

// ConnString returns the connection string the pool was created with.
func (s *Store) ConnString() string { return s.pool.Config().ConnString() }
