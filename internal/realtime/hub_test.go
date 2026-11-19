package realtime

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

// TestListenerWithPoolParameters checks that the listener connects with a pool_* URL parameter, which PostgreSQL
// rejects when it is sent as a runtime parameter.
func TestListenerWithPoolParameters(t *testing.T) {
	base := os.Getenv("TEST_DATABASE_URL")
	if base == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	h := &Hub{DatabaseURL: base + sep + "pool_max_conns=2", Log: slog.New(slog.NewTextHandler(io.Discard, nil))}
	msgs, leave := h.Subscribe(1, 1)
	defer leave()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go h.Run(ctx)

	select {
	case m := <-msgs:
		if m.Topic != TopicResync {
			t.Fatalf("got %+v, want a resync", m)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("listener did not connect")
	}
}
