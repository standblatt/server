package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Range is a shooting range of a tenant with its lanes.
type Range struct {
	ID            int64        `json:"id"`
	Name          string       `json:"name"`
	LaneCount     int          `json:"laneCount"`
	DisabledLanes []int32      `json:"disabledLanes"`
	Groups        []RangeGroup `json:"groups" doc:"Lane groups, e.g. 25 m and 50 m of one pistol range; empty = one group"`
}

// RangeGroup is a named group of lanes of a range, typically one distance.
type RangeGroup struct {
	ID    int64   `json:"id,omitempty" required:"false" doc:"Unset for a new group"`
	Name  string  `json:"name" minLength:"1" maxLength:"50"`
	Lanes []int32 `json:"lanes" minItems:"1"`
}

// RangeInput holds the editable fields of a range.
type RangeInput struct {
	Name          string       `json:"name" minLength:"1" maxLength:"100"`
	LaneCount     int          `json:"laneCount" minimum:"1" maximum:"99"`
	DisabledLanes []int32      `json:"disabledLanes" required:"false"`
	Groups        []RangeGroup `json:"groups" required:"false" doc:"Groups with an ID are updated, without one added, missing ones removed"`
}

const rangeGroupsCol = `coalesce((SELECT json_agg(json_build_object('id', g.id, 'name', g.name, 'lanes', g.lanes) ORDER BY g.position)
	FROM range_groups g WHERE g.range_id = rg.id), '[]')`

// normLanes returns an empty slice instead of nil.
func normLanes(lanes []int32) []int32 {
	if lanes == nil {
		return []int32{}
	}
	return lanes
}

// ListRanges lists the ranges of a tenant.
func (s *Store) ListRanges(ctx context.Context, tenantID int64) ([]Range, error) {
	rows, err := s.pool.Query(ctx, `SELECT rg.id, rg.name, rg.lane_count, rg.disabled_lanes, `+rangeGroupsCol+` FROM ranges rg
		WHERE rg.tenant_id = $1 ORDER BY rg.name`, tenantID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Range])
}

// CreateRange adds a range with its lane groups to a tenant.
func (s *Store) CreateRange(ctx context.Context, tenantID int64, in RangeInput) (Range, error) {
	st := Range{Name: in.Name, LaneCount: in.LaneCount, DisabledLanes: normLanes(in.DisabledLanes)}
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `INSERT INTO ranges (tenant_id, name, lane_count, disabled_lanes)
			VALUES ($1, $2, $3, $4) RETURNING id`, tenantID, st.Name, st.LaneCount, st.DisabledLanes).Scan(&st.ID); err != nil {
			return err
		}
		var err error
		st.Groups, err = setRangeGroups(ctx, tx, tenantID, st.ID, in.Groups)
		return err
	})
	return st, mapErr(err)
}

// UpdateRange changes a range of a tenant and its lane groups; ErrInUse when a removed group has programmes.
func (s *Store) UpdateRange(ctx context.Context, tenantID, rangeID int64, in RangeInput) error {
	return mapErr(pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := affected(tx.Exec(ctx, `UPDATE ranges SET name = $3, lane_count = $4, disabled_lanes = $5
			WHERE tenant_id = $1 AND id = $2`, tenantID, rangeID, in.Name, in.LaneCount, normLanes(in.DisabledLanes))); err != nil {
			return err
		}
		_, err := setRangeGroups(ctx, tx, tenantID, rangeID, in.Groups)
		return err
	}))
}

// setRangeGroups makes the groups of a range the given ones in order: with an ID updated (ErrNotFound if it is
// not a group of the range), without one inserted, others deleted. It returns the groups with their IDs.
func setRangeGroups(ctx context.Context, tx pgx.Tx, tenantID, rangeID int64, groups []RangeGroup) ([]RangeGroup, error) {
	out := []RangeGroup{}
	keep := []int64{}
	for _, g := range groups {
		if g.ID != 0 {
			keep = append(keep, g.ID)
		}
	}
	if _, err := tx.Exec(ctx, `DELETE FROM range_groups WHERE range_id = $1 AND id <> ALL($2)`, rangeID, keep); err != nil {
		return nil, err
	}
	// Free the names first so groups can swap names within one update.
	if _, err := tx.Exec(ctx, `UPDATE range_groups SET name = '#' || id WHERE range_id = $1`, rangeID); err != nil {
		return nil, err
	}
	for i, g := range groups {
		g.Lanes = normLanes(g.Lanes)
		if g.ID != 0 {
			if err := affected(tx.Exec(ctx, `UPDATE range_groups SET name = $4, lanes = $5, position = $6
				WHERE tenant_id = $1 AND range_id = $2 AND id = $3`, tenantID, rangeID, g.ID, g.Name, g.Lanes, i)); err != nil {
				return nil, err
			}
		} else if err := tx.QueryRow(ctx, `INSERT INTO range_groups (tenant_id, range_id, name, lanes, position)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`, tenantID, rangeID, g.Name, g.Lanes, i).Scan(&g.ID); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}

// DeleteRange removes a range; ErrInUse while events use it.
func (s *Store) DeleteRange(ctx context.Context, tenantID, rangeID int64) error {
	return affected(s.pool.Exec(ctx, `DELETE FROM ranges WHERE tenant_id = $1 AND id = $2`, tenantID, rangeID))
}
