package web

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/store"
)

const (
	apiTokenPrefix = "sb_"
	scopeRead      = "training:read"
	scopeWrite     = "training:write"
)

// API tokens only reach the training data of their owner, never the account, the
// token management or club data.
var tokenPaths = []string{"/api/my/trainings", "/api/my/stats", "/api/my/disciplines", "/api/my/weapons", "/api/my/goals", "/api/my/competitions", "/api/my/export"}

// tokenPathAllowed reports whether an API token may be used for path: /api/my and everything under tokenPaths.
func tokenPathAllowed(path string) bool {
	if path == "/api/my" {
		return true
	}
	for _, p := range tokenPaths {
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

// bearerToken returns the Bearer token of the Authorization header, else the X-API-Key header.
func bearerToken(r *http.Request) string {
	if v, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok {
		return strings.TrimSpace(v)
	}
	return strings.TrimSpace(r.Header.Get("X-API-Key"))
}

// writeProblem answers an RFC 9457 problem JSON with status and detail, for errors outside huma.
func writeProblem(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{"title": http.StatusText(status), "status": status, "detail": detail})
}

// tokenUser authenticates a request to /api/my/* that carries an API token. It returns
// ok=false after it has written an error.
func (s *Server) tokenUser(ctx context.Context, w http.ResponseWriter, r *http.Request, token string) (store.User, bool) {
	a, err := s.Store.APITokenByHash(ctx, sessionID(token))
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeProblem(w, http.StatusUnauthorized, "Invalid or revoked API token")
		return store.User{}, false
	case err != nil:
		s.Log.Error("api token lookup", "err", err)
		writeProblem(w, http.StatusInternalServerError, "internal error")
		return store.User{}, false
	case a.ExpiresAt != nil && a.ExpiresAt.Before(time.Now()):
		writeProblem(w, http.StatusUnauthorized, "API token has expired")
		return store.User{}, false
	case !tokenPathAllowed(r.URL.Path):
		writeProblem(w, http.StatusForbidden, "API tokens cannot use this endpoint")
		return store.User{}, false
	}
	need := scopeWrite
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		need = scopeRead
	}
	if !slices.Contains(a.Scopes, need) {
		writeProblem(w, http.StatusForbidden, "Token is missing the required scope: "+need)
		return store.User{}, false
	}
	if err := s.Store.TouchAPIToken(ctx, a.ID); err != nil {
		s.Log.Error("api token touch", "err", err)
	}
	return a.User, true
}

// NewAPIToken is the request to create a personal API token.
type NewAPIToken struct {
	Name          string   `json:"name" minLength:"1" maxLength:"100"`
	Scopes        []string `json:"scopes" minItems:"1" uniqueItems:"true" enum:"training:read,training:write"`
	ExpiresInDays int      `json:"expiresInDays" required:"false" enum:"0,30,90,365" doc:"0 or omitted: never expires"`
}

// CreatedAPIToken is a new token with its secret, which is shown only this once.
type CreatedAPIToken struct {
	store.APIToken
	Token string `json:"token"`
}

// newAPITokenSecret returns a new API token secret: the sb_ prefix and 40 random hex digits.
func newAPITokenSecret() string {
	b := make([]byte, 20)
	rand.Read(b)
	return apiTokenPrefix + hex.EncodeToString(b)
}

// registerAPITokens registers the endpoints that list, create and revoke personal API tokens.
func (s *Server) registerAPITokens(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-my-api-tokens", Method: http.MethodGet, Path: "/api/my/tokens", Tags: []string{"my"},
		Summary: "Personal API tokens",
	}, func(ctx context.Context, _ *struct{}) (*struct{ Body []store.APIToken }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		list, err := s.Store.ListAPITokens(ctx, u.ID)
		return &struct{ Body []store.APIToken }{Body: list}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-my-api-token", Method: http.MethodPost, Path: "/api/my/tokens", Tags: []string{"my"},
		Summary: "Create a personal API token; the secret is only returned here", DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *struct{ Body NewAPIToken }) (*struct{ Body CreatedAPIToken }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		name := strings.TrimSpace(in.Body.Name)
		if name == "" {
			return nil, invalid("Name is empty")
		}
		var expires *time.Time
		if in.Body.ExpiresInDays > 0 {
			t := time.Now().AddDate(0, 0, in.Body.ExpiresInDays)
			expires = &t
		}
		secret := newAPITokenSecret()
		t, err := s.Store.CreateAPIToken(ctx, u.ID, name, sessionID(secret), in.Body.Scopes, expires)
		return &struct{ Body CreatedAPIToken }{Body: CreatedAPIToken{APIToken: t, Token: secret}}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-my-api-token", Method: http.MethodDelete, Path: "/api/my/tokens/{id}", Tags: []string{"my"},
		Summary: "Revoke a personal API token", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *struct {
		ID int64 `path:"id"`
	}) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.DeleteAPIToken(ctx, u.ID, in.ID), "API token not found")
	})
}
