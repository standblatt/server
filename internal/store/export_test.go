package store

import (
	"context"
	"io"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MigrateUntil applies the migrations up to and including until, so a test can fill the old schema.
func (s *Store) MigrateUntil(ctx context.Context, until string) error {
	return s.migrate(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), until)
}

// Pool exposes the connection pool for raw fixtures and checks.
func (s *Store) Pool() *pgxpool.Pool { return s.pool }
