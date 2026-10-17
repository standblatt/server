package store

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Errors returned by store methods instead of driver errors.
var (
	ErrNotFound = errors.New("store: not found")
	ErrExists   = errors.New("store: already exists")
	ErrInUse    = errors.New("store: referenced by other rows")
	ErrArchived = errors.New("store: event is archived")
)

// Archived reports whether err comes from a change to an archived event, which the database refuses.
func Archived(err error) bool {
	var pgErr *pgconn.PgError
	return errors.Is(err, ErrArchived) || errors.As(err, &pgErr) && pgErr.Code == "TLARC"
}

// mapErr maps no rows, foreign key, unique and archive violations to ErrNotFound, ErrInUse, ErrExists and ErrArchived.
func mapErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return ErrInUse
		case "23505":
			return ErrExists
		case "TLARC":
			return ErrArchived
		}
	}
	return err
}

// affected maps err, or returns ErrNotFound when the command changed no row.
func affected(tag pgconn.CommandTag, err error) error {
	if err != nil {
		return mapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
