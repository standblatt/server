// Command standblatt runs the Standblatt server. "standblatt openapi" prints the API description.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/standblatt/server/internal/mail"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
	"github.com/standblatt/server/internal/webui"
)

var version string

const minAdminToken = 32

const janitorInterval = time.Hour

type config struct {
	databaseURL   string
	listen        string
	metricsListen string
	logLevel      slog.Level
	adminToken    string
	trustProxy    bool
	publicURL     string
	environment   string
	smtp          mail.Sender
	signup        bool
	totp          bool
	matomo        web.Matomo
	apiRateLimit  int
}

// loadConfig reads and validates the configuration from the environment.
func loadConfig() (config, error) {
	c := config{
		databaseURL:   os.Getenv("DATABASE_URL"),
		listen:        envOr("STANDBLATT_LISTEN", ":8080"),
		metricsListen: ":9090",
		adminToken:    os.Getenv("STANDBLATT_ADMIN_TOKEN"),
	}
	if v, ok := os.LookupEnv("STANDBLATT_METRICS_LISTEN"); ok {
		c.metricsListen = v
	}
	if c.databaseURL == "" {
		return c, errors.New("DATABASE_URL is required")
	}
	if c.adminToken != "" && len(c.adminToken) < minAdminToken {
		return c, fmt.Errorf("STANDBLATT_ADMIN_TOKEN must be at least %d characters", minAdminToken)
	}
	var err error
	if c.trustProxy, err = envBool("STANDBLATT_TRUST_PROXY", false); err != nil {
		return c, err
	}
	c.publicURL = os.Getenv("STANDBLATT_PUBLIC_URL")
	c.environment = os.Getenv("STANDBLATT_ENVIRONMENT")
	c.smtp = mail.Sender{
		Host:     os.Getenv("STANDBLATT_SMTP_HOST"),
		User:     os.Getenv("STANDBLATT_SMTP_USER"),
		Password: os.Getenv("STANDBLATT_SMTP_PASSWORD"),
		From:     envOr("STANDBLATT_SMTP_FROM", "Standblatt <noreply@standblatt.ch>"),
	}
	if c.smtp.Port, err = strconv.Atoi(envOr("STANDBLATT_SMTP_PORT", "587")); err != nil || c.smtp.Port < 1 || c.smtp.Port > 65535 {
		return c, errors.New("STANDBLATT_SMTP_PORT must be a port number")
	}
	if c.signup, err = envBool("STANDBLATT_SIGNUP", false); err != nil {
		return c, err
	}
	if c.totp, err = envBool("STANDBLATT_TOTP", true); err != nil {
		return c, err
	}
	if c.matomo, err = loadMatomo(os.Getenv("STANDBLATT_MATOMO_URL"), os.Getenv("STANDBLATT_MATOMO_SITE_ID")); err != nil {
		return c, err
	}
	if c.apiRateLimit, err = strconv.Atoi(envOr("STANDBLATT_API_RATE_LIMIT", "600")); err != nil || c.apiRateLimit < 0 {
		return c, errors.New("STANDBLATT_API_RATE_LIMIT must be a number of requests, 0 for no limit")
	}
	if err := c.logLevel.UnmarshalText([]byte(envOr("STANDBLATT_LOG_LEVEL", "info"))); err != nil {
		return c, errors.New("STANDBLATT_LOG_LEVEL must be debug, info, warn or error")
	}
	return c, nil
}

// loadMatomo validates the optional Matomo URL and site ID; both empty disables analytics.
func loadMatomo(rawURL, siteID string) (web.Matomo, error) {
	if rawURL == "" && siteID == "" {
		return web.Matomo{}, nil
	}
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		return web.Matomo{}, errors.New("STANDBLATT_MATOMO_URL must be an http(s) URL, e.g. https://analytics.example.org/")
	}
	id, err := strconv.Atoi(siteID)
	if err != nil || id < 1 {
		return web.Matomo{}, errors.New("STANDBLATT_MATOMO_SITE_ID must be the numeric site ID")
	}
	return web.Matomo{URL: strings.TrimSuffix(u.String(), "/") + "/", SiteID: id}, nil
}

// envBool parses the environment variable key as true or false, with def when it is empty.
func envBool(key string, def bool) (bool, error) {
	switch envOr(key, map[bool]string{true: "true", false: "false"}[def]) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, errors.New(key + " must be true or false")
}

// envOr returns the environment variable key, or def when it is empty.
func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// newLogger returns a JSON logger on stderr at the configured level.
func newLogger(c config) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: c.logLevel}))
}

// main prints the OpenAPI spec for "openapi", otherwise loads the configuration and runs the server until a signal.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "openapi" {
		spec, err := (&web.Server{Version: version}).OpenAPI()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Stdout.Write(spec)
		return
	}
	cfg, err := loadConfig()
	if err != nil {
		newLogger(config{}).Error("config", "err", err)
		os.Exit(2)
	}
	log := newLogger(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg, log); err != nil {
		log.Error("exit", "err", err)
		os.Exit(1)
	}
}

// run migrates the database, starts the realtime hub, janitors and HTTP servers and shuts down when ctx ends.
func run(ctx context.Context, cfg config, log *slog.Logger) error {
	log.Info("starting", "version", version, "listen", cfg.listen, "metrics", cfg.metricsListen)
	st, err := store.New(ctx, cfg.databaseURL)
	if err != nil {
		return err
	}
	defer st.Close()
	if err := st.Migrate(ctx, log); err != nil {
		return err
	}

	hub := &realtime.Hub{DatabaseURL: cfg.databaseURL, Log: log}
	go hub.Run(ctx)
	cfg.smtp.Log = log
	srv := &web.Server{
		Store:       st,
		Hub:         hub,
		UI:          webui.FS(),
		Version:     version,
		TrustProxy:  cfg.trustProxy,
		AdminToken:  cfg.adminToken,
		Log:         log,
		PublicURL:   cfg.publicURL,
		Environment: cfg.environment,
		SendMail:    cfg.smtp.Send,
		Signup:      cfg.signup,
		TOTP:        cfg.totp,
		Matomo:      cfg.matomo,

		APIRateLimit: cfg.apiRateLimit,
	}
	if cfg.adminToken == "" {
		log.Info("STANDBLATT_ADMIN_TOKEN not set, admin panel disabled")
	}
	go janitor(ctx, st, log)
	go printJanitor(ctx, st, log)
	go srv.RunReservationTicker(ctx)
	httpSrv := &http.Server{
		Addr:              cfg.listen,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	var metricsSrv *http.Server
	if cfg.metricsListen != "" {
		mux := http.NewServeMux()
		mux.Handle("GET /metrics", promhttp.Handler())
		metricsSrv = &http.Server{Addr: cfg.metricsListen, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
		go func() {
			if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("metrics server", "err", err)
			}
		}()
	}

	errc := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", cfg.listen)
		errc <- httpSrv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}

	log.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = httpSrv.Shutdown(shutdownCtx)
	if metricsSrv != nil {
		metricsSrv.Close()
	}
	return err
}

// janitor removes expired sessions, rate limits, print jobs, user tokens, invitations and pairings now and every
// janitorInterval.
func janitor(ctx context.Context, st *store.Store, log *slog.Logger) {
	tick := time.NewTicker(janitorInterval)
	defer tick.Stop()
	for {
		if n, err := st.DeleteExpiredSessions(ctx); err != nil {
			log.Warn("session cleanup", "err", err)
		} else if n > 0 {
			log.Info("expired sessions removed", "count", n)
		}
		if err := st.DeleteOldRateLimits(ctx); err != nil {
			log.Warn("rate limit cleanup", "err", err)
		}
		if err := st.DeleteOldPrintJobs(ctx); err != nil {
			log.Warn("print job cleanup", "err", err)
		}
		if _, err := st.DeleteExpiredUserTokens(ctx); err != nil {
			log.Warn("user token cleanup", "err", err)
		}
		if err := st.DeleteOldInvitations(ctx); err != nil {
			log.Warn("invitation cleanup", "err", err)
		}
		if err := st.DeleteOldPairings(ctx); err != nil {
			log.Warn("pairing cleanup", "err", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
	}
}

// printJanitor requeues print jobs that agents did not confirm within two minutes, every 30 seconds.
func printJanitor(ctx context.Context, st *store.Store, log *slog.Logger) {
	tick := time.NewTicker(30 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		if err := st.RequeueStalePrintJobs(ctx, 2*time.Minute); err != nil {
			log.Warn("print job requeue", "err", err)
		}
	}
}
