// Package testdb gives integration tests a freshly migrated database of their own.
package testdb

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/store"
)

// templateLockKey serializes building the template database across test processes.
const templateLockKey = 5184228

var (
	templateOnce sync.Once
	templateName string
	templateErr  error
)

// New creates a temporary database from TEST_DATABASE_URL as a copy of a migrated template and drops it when the
// test ends. The template is built once per set of migrations and kept for later runs.
// The test is skipped when TEST_DATABASE_URL is not set.
func New(t testing.TB) *store.Store {
	t.Helper()
	base := baseURL(t)
	templateOnce.Do(func() { templateName, templateErr = buildTemplate(base) })
	if templateErr != nil {
		t.Fatal(templateErr)
	}
	return create(t, base, templateName)
}

// Bare creates a temporary database like New but leaves it empty, for tests of the migrations themselves.
func Bare(t testing.TB) *store.Store {
	t.Helper()
	return create(t, baseURL(t), "")
}

// baseURL returns TEST_DATABASE_URL or skips the test.
func baseURL(t testing.TB) string {
	base := os.Getenv("TEST_DATABASE_URL")
	if base == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	return base
}

// buildTemplate makes sure the template database of the current migrations exists and drops those of older ones.
func buildTemplate(base string) (string, error) {
	ctx := context.Background()
	name := "standblatt_tpl_" + store.MigrationsHash()
	admin, err := pgx.Connect(ctx, base)
	if err != nil {
		return "", err
	}
	defer admin.Close(ctx)
	if _, err := admin.Exec(ctx, `SELECT pg_advisory_lock($1)`, templateLockKey); err != nil {
		return "", err
	}
	defer admin.Exec(ctx, `SELECT pg_advisory_unlock($1)`, templateLockKey)

	rows, err := admin.Query(ctx, `SELECT datname FROM pg_database WHERE datname LIKE 'standblatt\_tpl\_%'`)
	if err != nil {
		return "", err
	}
	existing, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return "", err
	}
	found := false
	for _, db := range existing {
		if db == name {
			found = true
		} else if _, err := admin.Exec(ctx, "DROP DATABASE "+db+" WITH (FORCE)"); err != nil {
			return "", err
		}
	}
	if found {
		return name, nil
	}

	if _, err := admin.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		return "", err
	}
	st, err := store.New(ctx, dbURL(base, name))
	if err == nil {
		err = st.Migrate(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)))
		st.Close()
	}
	if err != nil {
		admin.Exec(ctx, "DROP DATABASE "+name+" WITH (FORCE)")
		return "", err
	}
	return name, nil
}

// create makes a database with a random name, copied from template unless it is empty, and drops it when the test ends.
func create(t testing.TB, base, template string) *store.Store {
	t.Helper()
	ctx := context.Background()
	admin, err := pgx.Connect(ctx, base)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, 6)
	rand.Read(b)
	name := "standblatt_test_" + hex.EncodeToString(b)
	sql := "CREATE DATABASE " + name
	if template != "" {
		sql += " TEMPLATE " + template
	}
	if _, err := admin.Exec(ctx, sql); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		admin.Exec(context.Background(), "DROP DATABASE "+name+" WITH (FORCE)")
		admin.Close(context.Background())
	})

	st, err := store.New(ctx, dbURL(base, name))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(st.Close)
	return st
}

// dbURL returns base with another database name.
func dbURL(base, name string) string {
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	u.Path = "/" + name
	return u.String()
}
