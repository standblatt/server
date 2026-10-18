package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/standblatt/server/internal/domain"
)

// LayoutInput holds the editable fields of a print form.
type LayoutInput struct {
	Name string `json:"name" minLength:"1" maxLength:"100"`
	domain.Layout
}

// Layout is a print form of a tenant.
type Layout struct {
	ID            int64     `json:"id"`
	HasBackground bool      `json:"hasBackground" doc:"A scan of the empty form is stored for the preview"`
	UpdatedAt     time.Time `json:"updatedAt"`
	LayoutInput
}

// ListLayouts lists the print forms of a tenant by name.
func (s *Store) ListLayouts(ctx context.Context, tenantID int64) ([]Layout, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, background IS NOT NULL, updated_at, name, definition FROM layouts
		WHERE tenant_id = $1 ORDER BY lower(name)`, tenantID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanLayout)
}

// scanLayout reads a row of id, has-background, updated_at, name and definition.
func scanLayout(row pgx.CollectableRow) (Layout, error) {
	var l Layout
	err := row.Scan(&l.ID, &l.HasBackground, &l.UpdatedAt, &l.Name, &l.Layout)
	if l.Elements == nil {
		l.Elements = []domain.Element{}
	}
	return l, err
}

// GetLayout returns a print form of a tenant, without its background scan.
func (s *Store) GetLayout(ctx context.Context, tenantID, id int64) (Layout, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, background IS NOT NULL, updated_at, name, definition FROM layouts
		WHERE tenant_id = $1 AND id = $2`, tenantID, id)
	if err != nil {
		return Layout{}, err
	}
	l, err := pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (Layout, error) {
		var l Layout
		err := row.Scan(&l.ID, &l.HasBackground, &l.UpdatedAt, &l.Name, &l.Layout)
		return l, err
	})
	if l.Elements == nil {
		l.Elements = []domain.Element{}
	}
	return l, mapErr(err)
}

// LayoutBackground returns the background scan of a print form of a tenant (nil without one).
func (s *Store) LayoutBackground(ctx context.Context, tenantID, id int64) ([]byte, error) {
	var bg []byte
	err := s.pool.QueryRow(ctx, `SELECT background FROM layouts WHERE tenant_id = $1 AND id = $2`, tenantID, id).Scan(&bg)
	return bg, mapErr(err)
}

// CreateLayout adds a print form; ErrExists if the name is taken.
func (s *Store) CreateLayout(ctx context.Context, tenantID int64, in LayoutInput) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `INSERT INTO layouts (tenant_id, name, definition) VALUES ($1, $2, $3) RETURNING id`,
		tenantID, in.Name, in.Layout).Scan(&id)
	return id, mapErr(err)
}

// UpdateLayout changes the name and definition of a print form; the background stays.
func (s *Store) UpdateLayout(ctx context.Context, tenantID, id int64, in LayoutInput) error {
	return mapErr(affected(s.pool.Exec(ctx, `UPDATE layouts SET name = $3, definition = $4, updated_at = now()
		WHERE tenant_id = $1 AND id = $2`, tenantID, id, in.Name, in.Layout)))
}

// SetLayoutBackground stores the scan of the empty form, or removes it with nil.
func (s *Store) SetLayoutBackground(ctx context.Context, tenantID, id int64, image []byte) error {
	return affected(s.pool.Exec(ctx, `UPDATE layouts SET background = $3, updated_at = now() WHERE tenant_id = $1 AND id = $2`,
		tenantID, id, image))
}

// DeleteLayout removes a print form; programmes that used it keep no form.
func (s *Store) DeleteLayout(ctx context.Context, tenantID, id int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM layouts WHERE tenant_id = $1 AND id = $2`, tenantID, id))
}

// SetProgrammeLayout chooses the print form of a programme, or none with nil; allowed while in use. ErrNotFound
// unless both belong to the tenant.
func (s *Store) SetProgrammeLayout(ctx context.Context, tenantID, programmeID int64, layoutID *int64) error {
	return affected(s.pool.Exec(ctx, `UPDATE programmes SET layout_id = $3 WHERE tenant_id = $1 AND id = $2
		AND ($3::bigint IS NULL OR EXISTS (SELECT 1 FROM layouts WHERE tenant_id = $1 AND id = $3))`, tenantID, programmeID, layoutID))
}
