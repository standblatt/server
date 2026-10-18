package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// Image is an image of a club's library; Width and Height are in pixels.
type Image struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Size      int       `json:"size" doc:"Bytes"`
	CreatedAt time.Time `json:"createdAt"`
}

// ListImages lists the images of a tenant by name, without their data.
func (s *Store) ListImages(ctx context.Context, tenantID int64) ([]Image, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, width, height, length(data), created_at FROM images
		WHERE tenant_id = $1 ORDER BY lower(name)`, tenantID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Image])
}

// GetImage returns the data of an image of a tenant.
func (s *Store) GetImage(ctx context.Context, tenantID, id int64) ([]byte, error) {
	var data []byte
	err := s.pool.QueryRow(ctx, `SELECT data FROM images WHERE tenant_id = $1 AND id = $2`, tenantID, id).Scan(&data)
	return data, mapErr(err)
}

// CreateImage adds an image; ErrExists if the name is taken.
func (s *Store) CreateImage(ctx context.Context, tenantID int64, name string, data []byte, width, height int) (Image, error) {
	img := Image{Name: name, Width: width, Height: height, Size: len(data)}
	err := s.pool.QueryRow(ctx, `INSERT INTO images (tenant_id, name, data, width, height) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`, tenantID, name, data, width, height).Scan(&img.ID, &img.CreatedAt)
	return img, mapErr(err)
}

// DeleteImage removes an image; ErrInUse while a print form shows it.
func (s *Store) DeleteImage(ctx context.Context, tenantID, id int64) error {
	var used bool
	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM layouts l, jsonb_array_elements(l.definition->'elements') e
		WHERE l.tenant_id = $1 AND e->>'type' = 'image' AND (e->>'image')::bigint = $2)`, tenantID, id).Scan(&used); err != nil {
		return err
	}
	if used {
		return ErrInUse
	}
	return affected(s.pool.Exec(ctx, `DELETE FROM images WHERE tenant_id = $1 AND id = $2`, tenantID, id))
}

// ImagesByID returns the data of the tenant's images with these IDs; unknown IDs are missing from the map.
func (s *Store) ImagesByID(ctx context.Context, tenantID int64, ids []int64) (map[int64][]byte, error) {
	out := map[int64][]byte{}
	if len(ids) == 0 {
		return out, nil
	}
	rows, err := s.pool.Query(ctx, `SELECT id, data FROM images WHERE tenant_id = $1 AND id = ANY($2)`, tenantID, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var data []byte
		if err := rows.Scan(&id, &data); err != nil {
			return nil, err
		}
		out[id] = data
	}
	return out, rows.Err()
}

// MissingImages returns the IDs that are not images of the tenant.
func (s *Store) MissingImages(ctx context.Context, tenantID int64, ids []int64) ([]int64, error) {
	missing := []int64{}
	if len(ids) == 0 {
		return missing, nil
	}
	rows, err := s.pool.Query(ctx, `SELECT x FROM unnest($2::bigint[]) x WHERE NOT EXISTS (SELECT 1 FROM images WHERE tenant_id = $1 AND id = x)`, tenantID, ids)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[int64])
}
