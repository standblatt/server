package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Range is a shooting range of a tenant with its lanes.
type Range struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	LaneCount     int     `json:"laneCount"`
	DisabledLanes []int32 `json:"disabledLanes"`
}

// RangeInput holds the editable fields of a range.
type RangeInput struct {
	Name          string  `json:"name" minLength:"1" maxLength:"100"`
	LaneCount     int     `json:"laneCount" minimum:"1" maximum:"99"`
	DisabledLanes []int32 `json:"disabledLanes" required:"false"`
}

// normLanes returns an empty slice instead of nil.
func normLanes(lanes []int32) []int32 {
	if lanes == nil {
		return []int32{}
	}
	return lanes
}

// ListRanges lists the ranges of a tenant.
func (s *Store) ListRanges(ctx context.Context, tenantID int64) ([]Range, error) {
	rows, err := s.pool.Query(ctx, `SELECT rg.id, rg.name, rg.lane_count, rg.disabled_lanes FROM ranges rg
		WHERE rg.tenant_id = $1 ORDER BY rg.name`, tenantID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Range])
}

// CreateRange adds a range to a tenant.
func (s *Store) CreateRange(ctx context.Context, tenantID int64, in RangeInput) (Range, error) {
	st := Range{Name: in.Name, LaneCount: in.LaneCount, DisabledLanes: normLanes(in.DisabledLanes)}
	err := s.pool.QueryRow(ctx, `INSERT INTO ranges (tenant_id, name, lane_count, disabled_lanes)
		VALUES ($1, $2, $3, $4) RETURNING id`, tenantID, st.Name, st.LaneCount, st.DisabledLanes).Scan(&st.ID)
	return st, mapErr(err)
}

// UpdateRange changes a range of a tenant.
func (s *Store) UpdateRange(ctx context.Context, tenantID, rangeID int64, in RangeInput) error {
	return mapErr(affected(s.pool.Exec(ctx, `UPDATE ranges SET name = $3, lane_count = $4, disabled_lanes = $5
		WHERE tenant_id = $1 AND id = $2`, tenantID, rangeID, in.Name, in.LaneCount, normLanes(in.DisabledLanes))))
}

// DeleteRange removes a range; ErrInUse while events use it.
func (s *Store) DeleteRange(ctx context.Context, tenantID, rangeID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM ranges WHERE tenant_id = $1 AND id = $2`, tenantID, rangeID))
}
