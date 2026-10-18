package web

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/store"
)

const apiTokenPrefix = "sb_"

// tokenKey holds the API token of a request in its context.
const tokenKey ctxKey = 11

// apiRateWindow is the window of Server.APIRateLimit.
const apiRateWindow = 5 * time.Minute

// tokenPaths are the training endpoints a token with a training scope reaches, besides /api/my itself.
var tokenPaths = []string{"/api/my/trainings", "/api/my/stats", "/api/my/disciplines", "/api/my/weapons", "/api/my/goals", "/api/my/competitions", "/api/my/export"}

// trainingPath reports whether path is /api/my or under tokenPaths.
func trainingPath(path string) bool {
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

// tokenOf returns the API token a request came with, nil for sessions and helper links.
func tokenOf(ctx context.Context) *store.TokenAuth {
	t, _ := ctx.Value(tokenKey).(*store.TokenAuth)
	return t
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

// tokenRequest authenticates a request that carries an API token and checks where the token may go: club paths
// (the permission gate decides), the training paths with their scope and /api/token. It answers 429 for an IP with
// too many unknown tokens or a token over Server.APIRateLimit, and returns ok=false after it has written an error.
func (s *Server) tokenRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, secret string) (store.TokenAuth, bool) {
	ipKey := "token:ip:" + reqInfo(ctx).ip
	if n, err := s.Store.RateCount(ctx, ipKey, loginWindow); err == nil && n >= loginIPMax {
		w.Header().Set("Retry-After", strconv.Itoa(int(loginWindow.Seconds())))
		writeProblem(w, http.StatusTooManyRequests, "Too many invalid API tokens")
		return store.TokenAuth{}, false
	}
	a, err := s.Store.APITokenByHash(ctx, sessionID(secret))
	switch {
	case errors.Is(err, store.ErrNotFound):
		if _, err := s.Store.RateHit(ctx, ipKey, loginWindow); err != nil {
			s.Log.Warn("rate limit", "err", err)
		}
		writeProblem(w, http.StatusUnauthorized, "Invalid or revoked API token")
		return a, false
	case err != nil:
		s.Log.Error("api token lookup", "err", err)
		writeProblem(w, http.StatusInternalServerError, "internal error")
		return a, false
	case a.ExpiresAt != nil && a.ExpiresAt.Before(time.Now()):
		writeProblem(w, http.StatusUnauthorized, "API token has expired")
		return a, false
	}
	if s.APIRateLimit > 0 {
		if n, err := s.Store.RateHit(ctx, fmt.Sprintf("token:%d", a.ID), apiRateWindow); err != nil {
			s.Log.Warn("rate limit", "err", err)
		} else if n > s.APIRateLimit {
			w.Header().Set("Retry-After", strconv.Itoa(int(apiRateWindow.Seconds())))
			writeProblem(w, http.StatusTooManyRequests, "Too many requests for this API token")
			return a, false
		}
	}
	path := r.URL.Path
	switch {
	case strings.HasPrefix(path, "/api/tenants/"), path == "/api/token":
	case trainingPath(path):
		need := permissions.TrainingWrite
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			need = permissions.TrainingRead
		}
		if !slices.Contains(a.Scopes, need) {
			writeProblem(w, http.StatusForbidden, "Token is missing the required scope: "+need)
			return a, false
		}
	default:
		writeProblem(w, http.StatusForbidden, "API tokens cannot use this endpoint")
		return a, false
	}
	if err := s.Store.TouchAPIToken(ctx, a.ID); err != nil {
		s.Log.Error("api token touch", "err", err)
	}
	return a, true
}

// NewAPIToken is the request to create a personal API token.
type NewAPIToken struct {
	Name          string   `json:"name" minLength:"1" maxLength:"100"`
	Scopes        []string `json:"scopes" minItems:"1" uniqueItems:"true" enum:"events.view,events.operate,events.configure,reports.view,shooters.view,shooters.edit,shooters.manage,ranges.manage,forms.manage,club.manage,training:read,training:write" doc:"Permissions of the owner's role in the chosen clubs, and training:read/write for the own training log"`
	Tenants       []int64  `json:"tenants,omitempty" required:"false" uniqueItems:"true" doc:"The clubs the token may use; the owner must be a member"`
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

// TokenClub is a club an API token may use.
type TokenClub struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TokenInfo describes the API token of a request: what it is called, whose it is and what it may do.
type TokenInfo struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Owner     string      `json:"owner"`
	Scopes    []string    `json:"scopes"`
	Clubs     []TokenClub `json:"clubs"`
	ExpiresAt *time.Time  `json:"expiresAt,omitempty"`
}

// checkTokenScopes refuses clubs the user is no member of and club scopes that the user's role has in none of the
// chosen clubs, so a token never carries a permission its owner could not use.
func (s *Server) checkTokenScopes(ctx context.Context, userID int64, scopes []string, tenants []int64) error {
	ms, err := s.Store.Memberships(ctx, userID)
	if err != nil {
		return err
	}
	perms := map[string]bool{}
	for _, t := range tenants {
		i := slices.IndexFunc(ms, func(m store.Membership) bool { return m.TenantID == t })
		if i < 0 {
			return invalid("You are no member of club %d", t)
		}
		for _, p := range ms[i].Permissions {
			perms[p] = true
		}
	}
	for _, sc := range scopes {
		switch {
		case sc == permissions.TrainingRead || sc == permissions.TrainingWrite:
		case len(tenants) == 0:
			return invalid("Choose the clubs for %s", sc)
		case !perms[sc]:
			return invalid("Your role in the chosen clubs does not have %s", sc)
		}
	}
	return nil
}

// registerAPITokens registers the endpoints that list, create and revoke personal API tokens, and the token info.
func (s *Server) registerAPITokens(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-api-tokens", Method: http.MethodGet, Path: "/api/me/tokens", Tags: []string{"account"},
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
		OperationID: "create-api-token", Method: http.MethodPost, Path: "/api/me/tokens", Tags: []string{"account"},
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
		tenants := in.Body.Tenants
		if tenants == nil {
			tenants = []int64{}
		}
		if err := s.checkTokenScopes(ctx, u.ID, in.Body.Scopes, tenants); err != nil {
			return nil, err
		}
		var expires *time.Time
		if in.Body.ExpiresInDays > 0 {
			t := time.Now().AddDate(0, 0, in.Body.ExpiresInDays)
			expires = &t
		}
		secret := newAPITokenSecret()
		t, err := s.Store.CreateAPIToken(ctx, u.ID, name, sessionID(secret), in.Body.Scopes, tenants, expires)
		return &struct{ Body CreatedAPIToken }{Body: CreatedAPIToken{APIToken: t, Token: secret}}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-api-token", Method: http.MethodDelete, Path: "/api/me/tokens/{id}", Tags: []string{"account"},
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

	huma.Register(api, huma.Operation{
		OperationID: "token-info", Method: http.MethodGet, Path: "/api/token", Tags: []string{"account"},
		Summary: "The API token of the request: name, owner, scopes, clubs and expiry",
	}, func(ctx context.Context, _ *struct{}) (*struct{ Body TokenInfo }, error) {
		t := tokenOf(ctx)
		if t == nil {
			return nil, unauthorized("Send an API token")
		}
		info := TokenInfo{ID: t.ID, Name: t.Name, Owner: t.User.Name, Scopes: t.Scopes, Clubs: []TokenClub{}, ExpiresAt: t.ExpiresAt}
		for _, id := range t.Tenants {
			if tenant, err := s.Store.GetTenant(ctx, id); err == nil {
				info.Clubs = append(info.Clubs, TokenClub{ID: id, Name: tenant.Name})
			}
		}
		return &struct{ Body TokenInfo }{Body: info}, nil
	})
}
