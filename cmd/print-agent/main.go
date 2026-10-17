// Command print-agent prints Standblatt documents on the local CUPS printers of a range.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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

// main configures the agent from the environment and runs it until SIGINT or SIGTERM.
func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	a := &agent.Agent{
		URL:       os.Getenv("STANDBLATT_URL"),
		Token:     os.Getenv("STANDBLATT_AGENT_TOKEN"),
		LP:        envOr("STANDBLATT_LP", "lp"),
		LPStat:    envOr("STANDBLATT_LPSTAT", "lpstat"),
		LPOptions: envOr("STANDBLATT_LPOPTIONS", "lpoptions"),
		Client:    &http.Client{Timeout: 60 * time.Second},
		Log:       log,
	}
	if a.URL == "" || a.Token == "" {
		log.Error("STANDBLATT_URL and STANDBLATT_AGENT_TOKEN are required")
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	log.Info("print agent started", "version", version, "server", a.URL)
	a.Run(ctx)
}
