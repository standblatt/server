package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// Classification is a ranking over programmes of an event.
type Classification struct {
	ID int64 `json:"id"`
	domain.ClassificationDefinition
}

// ListClassifications lists the classifications of an event of a tenant in their order.
func (s *Store) ListClassifications(ctx context.Context, tenantID, eventID int64) ([]Classification, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, definition FROM classifications
		WHERE tenant_id = $1 AND event_id = $2 ORDER BY position, id`, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Classification, error) {
		var c Classification
		err := row.Scan(&c.ID, &c.ClassificationDefinition)
		c.Normalize()
		return c, err
	})
}

// CreateClassification adds a classification at the end of an event's list; ErrNotFound when the event is not the tenant's.
func (s *Store) CreateClassification(ctx context.Context, tenantID, eventID int64, c domain.ClassificationDefinition) (int64, error) {
	if err := eventInTenant(ctx, s.pool, tenantID, eventID); err != nil {
		return 0, err
	}
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO classifications (tenant_id, event_id, position, definition)
		VALUES ($1, $2, (SELECT coalesce(max(position), 0) + 1 FROM classifications WHERE event_id = $2), $3) RETURNING id`,
		tenantID, eventID, c).Scan(&id)
	return id, mapErr(err)
}

// UpdateClassification replaces a classification of an event of a tenant.
func (s *Store) UpdateClassification(ctx context.Context, tenantID, eventID, id int64, c domain.ClassificationDefinition) error {
	return affected(s.pool.Exec(ctx, `UPDATE classifications SET definition = $4 WHERE tenant_id = $1 AND event_id = $2 AND id = $3`,
		tenantID, eventID, id, c))
}

// DeleteClassification removes a classification of an event of a tenant.
func (s *Store) DeleteClassification(ctx context.Context, tenantID, eventID, id int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM classifications WHERE tenant_id = $1 AND event_id = $2 AND id = $3`, tenantID, eventID, id))
}

// copyClassifications copies the classifications of one event into another, mapping each programme to the one with
// the same name and variant there; parts without such a programme are left out.
func copyClassifications(ctx context.Context, tx pgx.Tx, tenantID, fromEventID, toEventID int64) error {
	rows, err := tx.Query(ctx, `SELECT s.id, t.id FROM programmes s
		JOIN programmes t ON t.event_id = $3 AND t.name = s.name AND t.variant = s.variant
		WHERE s.tenant_id = $1 AND s.event_id = $2`, tenantID, fromEventID, toEventID)
	if err != nil {
		return err
	}
	newID := map[int64]int64{}
	var from, to int64
	if _, err := pgx.ForEachRow(rows, []any{&from, &to}, func() error {
		newID[from] = to
		return nil
	}); err != nil {
		return err
	}
	rows, err = tx.Query(ctx, `SELECT definition FROM classifications WHERE tenant_id = $1 AND event_id = $2 ORDER BY position, id`, tenantID, fromEventID)
	if err != nil {
		return err
	}
	list, err := pgx.CollectRows(rows, pgx.RowTo[domain.ClassificationDefinition])
	if err != nil {
		return err
	}
	for i, c := range list {
		parts := []domain.ClassificationPart{}
		for _, p := range c.Parts {
			if id, ok := newID[p.ProgrammeID]; ok {
				parts = append(parts, domain.ClassificationPart{ProgrammeID: id, Count: p.Count})
			}
		}
		c.Parts = parts
		if _, err := tx.Exec(ctx, `INSERT INTO classifications (tenant_id, event_id, position, definition) VALUES ($1, $2, $3, $4)`,
			tenantID, toEventID, i+1, c); err != nil {
			return err
		}
	}
	return nil
}
