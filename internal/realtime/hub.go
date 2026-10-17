// Package realtime fans out change notifications to WebSocket clients. Every pod listens on a PostgreSQL
// channel, so a change made through any replica reaches the clients of all replicas.
package realtime

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

// Channel is the PostgreSQL notification channel.
const Channel = "standblatt_events"

// Topics a client can react to.
const (
	TopicRegistrations = "registrations"
	TopicLanes         = "lanes"
	TopicProgrammes    = "programmes"
	TopicResults       = "results"
	TopicPrint         = "print"
	// TopicHelpers tells helper pages that a link was revoked; they refetch and run into 401.
	TopicHelpers = "helpers"
	// TopicBooking tells the reservations and setup pages that booking requests or settings changed.
	TopicBooking = "booking"
	TopicResync  = "resync"
)

// Message tells clients of one event that data of a topic changed.
type Message struct {
	TenantID int64  `json:"tenantId"`
	EventID  int64  `json:"eventId"`
	Topic    string `json:"topic"`
}

type client struct {
	tenantID, eventID int64
	send              chan Message
}

// Hub keeps the connected clients of this pod and the listener connection.
type Hub struct {
	DatabaseURL string
	Log         *slog.Logger

	mu      sync.Mutex
	clients map[*client]struct{}
}

// Subscribe registers a client for the messages of one event; call the returned function to leave.
func (h *Hub) Subscribe(tenantID, eventID int64) (<-chan Message, func()) {
	c := &client{tenantID: tenantID, eventID: eventID, send: make(chan Message, 32)}
	h.mu.Lock()
	if h.clients == nil {
		h.clients = map[*client]struct{}{}
	}
	h.clients[c] = struct{}{}
	h.mu.Unlock()
	return c.send, func() {
		h.mu.Lock()
		if _, ok := h.clients[c]; ok {
			delete(h.clients, c)
			close(c.send)
		}
		h.mu.Unlock()
	}
}

// broadcast sends m to the clients of its tenant and event, or to all for a resync, dropping it for full buffers.
func (h *Hub) broadcast(m Message) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		if m.Topic != TopicResync && (c.tenantID != m.TenantID || c.eventID != m.EventID) {
			continue
		}
		select {
		case c.send <- m:
		default:
			// a client that does not keep up loses this message; it resyncs on reconnect
		}
	}
}

// Run listens for notifications until ctx ends, reconnecting with backoff.
func (h *Hub) Run(ctx context.Context) {
	backoff := time.Second
	for ctx.Err() == nil {
		start := time.Now()
		err := h.listen(ctx)
		if ctx.Err() != nil {
			return
		}
		if time.Since(start) > time.Minute {
			backoff = time.Second // the connection was healthy for a while, so this is a fresh outage
		}
		h.Log.Warn("realtime listener lost", "err", err, "retry", backoff)
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, 30*time.Second)
	}
}

// listen opens a LISTEN connection, tells clients to resync and forwards notifications until the connection fails.
func (h *Hub) listen(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, h.DatabaseURL)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	if _, err := conn.Exec(ctx, "LISTEN "+Channel); err != nil {
		return err
	}
	h.broadcast(Message{Topic: TopicResync})
	for {
		n, err := conn.WaitForNotification(ctx)
		if err != nil {
			return err
		}
		var m Message
		if err := json.Unmarshal([]byte(n.Payload), &m); err != nil {
			h.Log.Warn("realtime payload", "err", err)
			continue
		}
		h.broadcast(m)
	}
}
