package store

import (
	"context"
	"encoding/json"
)

// Notify publishes change topics of an event on the given PostgreSQL channel.
func (s *Store) Notify(ctx context.Context, channel string, tenantID, eventID int64, topics ...string) error {
	for _, t := range topics {
		payload, _ := json.Marshal(map[string]any{"tenantId": tenantID, "eventId": eventID, "topic": t})
		if _, err := s.pool.Exec(ctx, `SELECT pg_notify($1, $2)`, channel, string(payload)); err != nil {
			return err
		}
	}
	return nil
}
