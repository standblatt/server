// Command print-agent prints Standblatt documents on the local CUPS printers of a range.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/standblatt/server/internal/agent"
)

var version string

// envOr returns the environment variable key, or def when it is empty.
func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// main configures the agent from the environment and runs it until SIGINT or SIGTERM, one loop per token
// (STANDBLATT_AGENT_TOKEN may list several, comma separated, to serve several ranges from one computer).
func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	url := os.Getenv("STANDBLATT_URL")
	var tokens []string
	for t := range strings.SplitSeq(os.Getenv("STANDBLATT_AGENT_TOKEN"), ",") {
		if t = strings.TrimSpace(t); t != "" {
			tokens = append(tokens, t)
		}
	}
	if url == "" || len(tokens) == 0 {
		log.Error("STANDBLATT_URL and STANDBLATT_AGENT_TOKEN are required")
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	log.Info("print agent started", "version", version, "server", url, "ranges", len(tokens))
	var wg sync.WaitGroup
	for i, token := range tokens {
		a := &agent.Agent{
			URL:       url,
			Token:     token,
			LP:        envOr("STANDBLATT_LP", "lp"),
			LPStat:    envOr("STANDBLATT_LPSTAT", "lpstat"),
			LPOptions: envOr("STANDBLATT_LPOPTIONS", "lpoptions"),
			Version:   version,
			Client:    &http.Client{Timeout: 60 * time.Second},
			Log:       log.With("token", i+1),
		}
		wg.Go(func() { a.Run(ctx) })
	}
	wg.Wait()
}
