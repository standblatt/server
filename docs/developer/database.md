# Database

PostgreSQL, accessed with a handwritten pgx store (`internal/store`). Migrations are numbered SQL files in `internal/store/migrations`, embedded in the binary and applied at startup under an advisory lock. Add a new file for every change; never edit a migration that has been released.

## Guarantees in the database

| Rule | How |
|---|---|
| one active registration per lane and event | partial unique index `registrations_lane_taken` |
| a shooter shoots a programme once per event | unique constraint `entries_programme_once` |
| no cross-tenant references | composite foreign keys `(id, tenant_id)` |
| shooters, clubs and programmes in use cannot be deleted | foreign keys without cascade, mapped to `ErrInUse` |
| archived events cannot change | triggers from `0011_archive.sql`, SQLSTATE `TLARC`, mapped to 409 by the `huma.NewError` hook in `internal/web/server.go` |

Age categories are the enum `age_category`. Reservations overlap check locks the range row (`FOR UPDATE`) before testing `lanes && …`.

## Deleting

The cascade from an event could reach programmes before the entries that reference them, so `store.DeleteTenant` deletes registrations first, then events, then the tenant (after unarchiving its events). Events with registrations cannot be deleted from the app.
