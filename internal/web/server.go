// Package web wires the HTTP surface: health checks, the JSON API, the admin panel and the embedded SPA.
package web

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/standblatt/server/internal/mail"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

// init makes arrays non-nullable and turns archive trigger errors into a 409 in every huma error.
func init() {
	huma.DefaultArrayNullable = false
	newError := huma.NewError
	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		for _, err := range errs {
			if store.Archived(err) {
				return conflict("The event is archived; unarchive it to make changes").(*apiError)
			}
		}
		return newError(status, msg, errs...)
	}
}

// Server holds the dependencies of the HTTP handlers.
type Server struct {
	Store      *store.Store
	Hub        *realtime.Hub
	UI         fs.FS
	Version    string
	TrustProxy bool
	AdminToken string
	Log        *slog.Logger
	// PublicURL is the base of links in mails; empty means the scheme and host of the request.
	PublicURL string
	// SendMail delivers a mail; nil only logs it.
	SendMail func(context.Context, mail.Message) error
	Signup   bool
	// Environment names the deployment (e.g. staging); anything but production is shown as a banner in the SPA.
	Environment string
	// TOTP lets users add an authenticator app as second factor; off ignores stored secrets.
	TOTP bool
	// Matomo is the web analytics the SPA reports page views to; a zero value turns it off.
	Matomo Matomo
	// APIRateLimit is how many requests one API token may make per five minutes; 0 means no limit.
	APIRateLimit int
}

// Handler returns the root handler: /healthz, /readyz, /api/*, /admin and the SPA for every other path.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	mux.HandleFunc("GET /readyz", s.readyz)

	mux.HandleFunc("GET /api/tenants/{tenantId}/events/{eventId}/live", s.live)
	mux.HandleFunc("GET /api/kiosk/live", s.kioskLive)
	api := humago.New(mux, s.apiConfig())
	s.registerAPI(api)

	if s.AdminToken != "" {
		s.registerAdmin(mux)
	}

	mux.Handle("/", s.spa())
	withSession := s.withSession(mux)
	return s.secureHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			withSession.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	}))
}

// OpenAPI returns the API description as JSON without needing a database.
func (s *Server) OpenAPI() ([]byte, error) {
	api := humago.New(http.NewServeMux(), s.apiConfig())
	s.registerAPI(api)
	return api.OpenAPI().MarshalJSON()
}

// apiConfig returns the huma configuration with the API paths, error translation and the API token scheme.
func (s *Server) apiConfig() huma.Config {
	c := huma.DefaultConfig("Standblatt API", s.Version)
	c.OpenAPIPath = "/api/openapi"
	c.DocsPath = "/api/docs"
	c.SchemasPath = "/api/schemas"
	c.CreateHooks = nil
	c.Transformers = append(c.Transformers, translateErrors)
	c.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"apiToken": {Type: "http", Scheme: "bearer", Description: "Personal API token (sb_…) for the training log under /api/my, created in My training → Settings. Scopes: training:read for GET, training:write for everything else. X-API-Key works too."},
	}
	return c
}

// registerAPI registers all API operations.
func (s *Server) registerAPI(api huma.API) {
	api.UseMiddleware(s.authorize(api), s.permissionGate(api), s.moduleGate(api))
	s.registerModules(api)
	s.registerSystem(api)
	s.registerAuth(api)
	s.registerTenants(api)
	s.registerMasterData(api)
	s.registerResultImport(api)
	s.registerForms(api)
	s.registerImages(api)
	s.registerCertificates(api)
	s.registerProgrammes(api)
	s.registerRegistrations(api)
	s.registerResults(api)
	s.registerRankings(api)
	s.registerClassifications(api)
	s.registerDocuments(api)
	s.registerPrinting(api)
	s.registerReports(api)
	s.registerReservations(api)
	s.registerAccount(api)
	s.registerMy(api)
	s.registerAPITokens(api)
	s.registerHelpers(api)
	s.registerInvites(api)
	s.registerTOTP(api)
	s.registerBooking(api)
}

// readyz answers 200 when the database answers within 2 s, else 503.
func (s *Server) readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.Store.Ping(ctx); err != nil {
		s.Log.Warn("not ready", "err", err)
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	w.Write([]byte("ok"))
}

// spa serves the embedded SPA: files as they are, index.html for unknown routes, 404 for /api, missing assets and other methods.
func (s *Server) spa() http.Handler {
	files := http.FileServerFS(s.UI)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || (r.Method != http.MethodGet && r.Method != http.MethodHead) {
			http.NotFound(w, r)
			return
		}
		name := strings.TrimPrefix(r.URL.Path, "/")
		if name == "" {
			name = "index.html"
		}
		if _, err := fs.Stat(s.UI, name); err != nil {
			if strings.HasPrefix(name, "_app/") || strings.Contains(path.Base(name), ".") {
				http.NotFound(w, r)
				return
			}
			r = r.Clone(r.Context())
			r.URL.Path = "/"
		}
		if strings.HasPrefix(name, "_app/immutable/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		files.ServeHTTP(w, r)
	})
}
