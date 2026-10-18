package web

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
)

const (
	sessionCookie    = "standblatt_session"
	sessionIdle      = 7 * 24 * time.Hour
	sessionTTL       = 30 * 24 * time.Hour
	rememberTTL      = 90 * 24 * time.Hour
	maxPasswordBytes = 72

	loginWindow  = 15 * time.Minute
	loginIPMax   = 20
	loginUserMax = 5

	dummyHash = "$2a$12$invalidinvalidinvalidinvalidinvalidinvalidinvalidinval"
)

// bcryptCost is the cost of new password hashes; tests lower it.
var bcryptCost = 12

type ctxKey int

const (
	userKey ctxKey = iota
	requestKey
	cookieKey
)

type requestInfo struct {
	ip     string
	secure bool
	host   string
	// lang is the UI language the SPA sends as Accept-Language; mails fall back to it when the user has no locale yet.
	lang string
}

// HashPassword returns the bcrypt hash of a password.
func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(h), err
}

var userNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{2,31}$`)

// ValidateUserName checks the format of a (lower-cased) user name.
func ValidateUserName(name string) []domain.Problem {
	if !userNameRe.MatchString(name) {
		return []domain.Problem{domain.Problemf("User name: 3 to 32 characters, lower-case letters, digits, . _ -, starting with a letter or digit")}
	}
	return nil
}

// ValidatePassword checks the length rules for a password.
func ValidatePassword(password string) []domain.Problem {
	var errs []domain.Problem
	if len(password) < 10 {
		errs = append(errs, domain.Problemf("Password: at least 10 characters"))
	}
	if len(password) > maxPasswordBytes {
		errs = append(errs, domain.Problemf("Password: at most %d bytes", maxPasswordBytes))
	}
	return errs
}

// randomToken returns 32 random bytes as base64url, used for session, mail and kiosk secrets.
func randomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// sessionID returns the hex sha256 of a token, the form in which tokens are stored.
func sessionID(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// clientIP returns the client address: the last X-Forwarded-For entry behind a trusted proxy, else the remote address.
func (s *Server) clientIP(r *http.Request) string {
	if s.TrustProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			return strings.TrimSpace(parts[len(parts)-1])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// secureRequest reports whether the request came over TLS, directly or through a proxy.
func secureRequest(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

// withSession puts the request info and the caller (API token user, helper or session user) into the context of next.
func (s *Server) withSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestKey, requestInfo{ip: s.clientIP(r), secure: secureRequest(r), host: r.Host, lang: acceptLang(r.Header.Get("Accept-Language"))})
		if token := bearerToken(r); strings.HasPrefix(token, apiTokenPrefix) {
			if a, ok := s.tokenRequest(ctx, w, r, token); ok {
				ctx = context.WithValue(context.WithValue(ctx, userKey, a.User), tokenKey, &a)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
			return
		}
		h, err := s.helperFromRequest(ctx, r)
		if err != nil {
			s.Log.Error("helper lookup", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if h != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, helperKey, h)))
			return
		}
		if c, err := r.Cookie(sessionCookie); err == nil {
			ctx = context.WithValue(ctx, cookieKey, c.Value)
			u, err := s.Store.SessionUser(ctx, sessionID(c.Value))
			switch {
			case err == nil:
				ctx = context.WithValue(ctx, userKey, u)
			case !errors.Is(err, store.ErrNotFound):
				s.Log.Error("session lookup", "err", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// reqInfo returns the request info that withSession stored in ctx.
func reqInfo(ctx context.Context) requestInfo {
	ri, _ := ctx.Value(requestKey).(requestInfo)
	return ri
}

// currentUser returns the logged-in user from ctx, 401 when there is none.
func currentUser(ctx context.Context) (store.User, error) {
	u, ok := ctx.Value(userKey).(store.User)
	if !ok {
		return store.User{}, unauthorized("not logged in")
	}
	return u, nil
}

// sessionCookieFor returns the session cookie for token with maxAge seconds, Secure when the request was.
func sessionCookieFor(ctx context.Context, token string, maxAge int) http.Cookie {
	return http.Cookie{
		Name: sessionCookie, Value: token, Path: "/", HttpOnly: true, Secure: reqInfo(ctx).secure,
		SameSite: http.SameSiteLaxMode, MaxAge: maxAge,
	}
}

// Me is the logged-in user with the tenants they can access.
type Me struct {
	User          store.User         `json:"user"`
	Tenants       []store.Membership `json:"tenants"`
	Email         string             `json:"email" doc:"Empty until the user adds one"`
	EmailVerified bool               `json:"emailVerified"`
	TOTP          bool               `json:"totp" doc:"The user logs in with a second factor"`
}

// MeOutput is the response of GET /api/me.
type MeOutput struct {
	Body Me
}

// LocaleInput is the request of PUT /api/me/locale.
type LocaleInput struct {
	Body struct {
		Locale string `json:"locale" enum:"en,de,fr,it"`
	}
}

// LoginInput is the request of POST /api/auth/login.
type LoginInput struct {
	Body struct {
		Name     string `json:"name" minLength:"1"`
		Password string `json:"password" minLength:"1"`
		Code     string `json:"code,omitempty" required:"false" maxLength:"8" doc:"Authenticator code; the login answers 428 when it is needed"`
		Remember bool   `json:"remember,omitempty" required:"false" doc:"Stay logged in for 90 days instead of 7 days of inactivity (30 at most)"`
	}
}

// LoginOutput sets the session cookie and returns the user.
type LoginOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      Me
}

// LogoutOutput clears the session cookie.
type LogoutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

// me returns the Me view of u: memberships, whether the personal log is set up and whether the second factor is on.
func (s *Server) me(ctx context.Context, u store.User) (Me, error) {
	ms, err := s.Store.Memberships(ctx, u.ID)
	if err != nil {
		return Me{}, err
	}
	a, err := s.Store.GetAccount(ctx, u.ID)
	if err != nil {
		return Me{}, err
	}
	secret, _, err := s.Store.TOTP(ctx, u.ID)
	if err != nil {
		return Me{}, err
	}
	return Me{User: u, Tenants: ms, Email: a.Email, EmailVerified: a.EmailVerified, TOTP: s.TOTP && secret != ""}, nil
}

// registerAuth registers the login, logout, current user and locale endpoints.
func (s *Server) registerAuth(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "login", Method: http.MethodPost, Path: "/api/auth/login", Tags: []string{"auth"},
		Summary: "Log in and start a session",
	}, s.login)

	huma.Register(api, huma.Operation{
		OperationID: "logout", Method: http.MethodPost, Path: "/api/auth/logout", Tags: []string{"auth"},
		Summary: "End the current session", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, _ *struct{}) (*LogoutOutput, error) {
		if c, ok := ctx.Value(cookieKey).(string); ok {
			if err := s.Store.DeleteSession(ctx, sessionID(c)); err != nil {
				s.Log.Warn("logout", "err", err)
			}
		}
		return &LogoutOutput{SetCookie: sessionCookieFor(ctx, "", -1)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-me", Method: http.MethodGet, Path: "/api/me", Tags: []string{"auth"},
		Summary: "The logged-in user and their tenants",
	}, func(ctx context.Context, _ *struct{}) (*MeOutput, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		me, err := s.me(ctx, u)
		return &MeOutput{Body: me}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "set-my-locale", Method: http.MethodPut, Path: "/api/me/locale", Tags: []string{"auth"},
		Summary: "Set the UI language of the logged-in user", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *LocaleInput) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		return nil, s.Store.SetUserLocale(ctx, u.ID, in.Body.Locale)
	})
}

// login checks maintenance mode, rate limits, password and second factor, then starts a session and sets its cookie.
func (s *Server) login(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
	name := strings.ToLower(strings.TrimSpace(in.Body.Name))
	ip := reqInfo(ctx).ip

	on, err := s.Store.Maintenance(ctx)
	if err != nil {
		return nil, err
	}
	if on {
		return nil, &apiError{http.StatusServiceUnavailable, []domain.Problem{domain.Problemf("Maintenance: logging in is not possible right now.")}}
	}

	ipKey, nameKey := "login:ip:"+ip, "login:user:"+name
	for _, l := range []struct {
		key string
		max int
	}{{ipKey, loginIPMax}, {nameKey, loginUserMax}} {
		n, err := s.Store.RateCount(ctx, l.key, loginWindow)
		if err != nil {
			return nil, err
		}
		if n >= l.max {
			s.Log.Warn("login rate limited", "user", name, "remote", ip)
			return nil, tooMany("Too many attempts. Try again in %d minutes.", int(loginWindow.Minutes()))
		}
	}

	u, hash, err := s.Store.Credentials(ctx, name)
	if errors.Is(err, store.ErrNotFound) {
		hash, err = dummyHash, nil
	}
	if err != nil {
		return nil, err
	}
	// unknown users burn a real bcrypt compare so names cannot be probed by timing
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Body.Password)) != nil || u.ID == 0 {
		for _, k := range []string{ipKey, nameKey} {
			if _, err := s.Store.RateHit(ctx, k, loginWindow); err != nil {
				s.Log.Warn("rate limit", "err", err)
			}
		}
		s.Log.Warn("login failed", "user", name, "remote", ip)
		return nil, unauthorized("Wrong user name or password.")
	}
	if required, err := s.checkTOTP(ctx, u.ID, in.Body.Code); err != nil {
		if required && in.Body.Code != "" {
			for _, k := range []string{ipKey, nameKey} {
				if _, err := s.Store.RateHit(ctx, k, loginWindow); err != nil {
					s.Log.Warn("rate limit", "err", err)
				}
			}
			s.Log.Warn("login code failed", "user", name, "remote", ip)
		}
		return nil, err
	}

	idle, ttl := sessionIdle, sessionTTL
	if in.Body.Remember {
		idle, ttl = rememberTTL, rememberTTL
	}
	token := randomToken()
	if err := s.Store.CreateSession(ctx, sessionID(token), u.ID, idle, ttl); err != nil {
		return nil, err
	}
	s.Log.Info("login", "user", name, "remote", ip, "remember", in.Body.Remember)
	me, err := s.me(ctx, u)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{SetCookie: sessionCookieFor(ctx, token, int(ttl.Seconds())), Body: me}, nil
}
