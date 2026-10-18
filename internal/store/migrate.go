package store

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
)

//go:embed migrations/*.sql
var migrations embed.FS

const migrateLockKey = 5184227

// MigrationsHash returns a short hash over the names and contents of the embedded migrations.
func MigrationsHash() string {
	h := sha256.New()
	fs.WalkDir(migrations, "migrations", func(path string, _ fs.DirEntry, _ error) error {
		b, _ := migrations.ReadFile(path)
		h.Write([]byte(path))
		h.Write(b)
		return nil
	})
	return hex.EncodeToString(h.Sum(nil))[:12]
}

// Migrate applies all pending embedded migrations in name order. A session-level advisory lock
// serialises concurrent starts, so several replicas can boot at the same time.
func (s *Store) Migrate(ctx context.Context, log *slog.Logger) error {
	return s.migrate(ctx, log, "")
}

// migrate applies the pending migrations up to and including until, or all of them when until is empty.
func (s *Store) migrate(ctx context.Context, log *slog.Logger, until string) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, `SELECT pg_advisory_lock($1)`, migrateLockKey); err != nil {
		return fmt.Errorf("migrate: lock: %w", err)
	}
	defer conn.Exec(context.Background(), `SELECT pg_advisory_unlock($1)`, migrateLockKey)

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version text PRIMARY KEY,
		applied_at timestamptz NOT NULL DEFAULT now())`)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		if until != "" && name > until {
			break
		}
		var applied bool
		if err := conn.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)`, name).Scan(&applied); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
		if applied {
			continue
		}
		sql, err := migrations.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}
		tx, err := conn.Begin(ctx)
		if err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("migrate %s: %w", name, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, name); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("migrate %s: %w", name, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("migrate %s: %w", name, err)
		}
		log.Info("applied migration", "version", name)
	}
	return nil
}
