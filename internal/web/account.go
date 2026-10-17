package web

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/i18n"
	"github.com/standblatt/server/internal/mail"
	"github.com/standblatt/server/internal/store"
)

const (
	verifyTTL    = 7 * 24 * time.Hour
	resetTTL     = time.Hour
	signupWindow = time.Hour
	signupIPMax  = 10
	forgotWindow = 15 * time.Minute
	forgotIPMax  = 5
)

// link builds an absolute URL for path from PublicURL, or from the scheme and host of the current request.
func (s *Server) link(ctx context.Context, path string) string {
	if s.PublicURL != "" {
		return strings.TrimSuffix(s.PublicURL, "/") + path
	}
	ri := reqInfo(ctx)
	return map[bool]string{true: "https", false: "http"}[ri.secure] + "://" + ri.host + path
}

// mailInterval spaces outgoing mails so a bulk invitation never sends more than four per second per pod.
const mailInterval = 250 * time.Millisecond

var (
	mailMu   sync.Mutex
	mailNext time.Time
)

// mail sends m in the background with a one-minute timeout and logs failures instead of returning them.
// Sends are spaced by mailInterval, so a burst drains at four mails per second.
func (s *Server) mail(m mail.Message) {
	send := s.SendMail
	if send == nil {
		send = (&mail.Sender{Log: s.Log}).Send
	}
	mailMu.Lock()
	slot := mailNext
	if now := time.Now(); slot.Before(now) {
		slot = now
	}
	mailNext = slot.Add(mailInterval)
	mailMu.Unlock()
	go func() {
		time.Sleep(time.Until(slot))
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := send(ctx, m); err != nil {
			s.Log.Error("mail", "to", m.To, "err", err)
		}
	}()
}

// greeting returns the translated salutation of a mail, with the name when one is known.
func greeting(l i18n.Lang, name string) string {
	if name == "" {
		return i18n.T(l, "Hello,")
	}
	return i18n.Tf(l, "Hello %s,", name)
}

// sendVerification stores a verify token for the user and mails the confirmation link in the user's language.
func (s *Server) sendVerification(ctx context.Context, userID int64, email, name, locale string) error {
	l := mailLang(ctx, locale)
	token := randomToken()
	if err := s.Store.CreateUserToken(ctx, sessionID(token), userID, "verify", time.Now().Add(verifyTTL)); err != nil {
		return err
	}
	s.mail(mail.Message{To: email, Subject: i18n.T(l, "Confirm your email address"), Body: strings.Join([]string{
		greeting(l, name), "",
		i18n.T(l, "please confirm your email address for Standblatt:"), "",
		s.link(ctx, "/verify?token="+token), "",
		i18n.T(l, "The link is valid for 7 days. If you did not create an account, ignore this mail."),
	}, "\n")})
	return nil
}

// AuthOptions tells the login page what this server offers.
type AuthOptions struct {
	Signup bool `json:"signup"`
	TOTP   bool `json:"totp" doc:"Users may set up an authenticator app"`
}

// SignupInput is the request of POST /api/auth/signup.
type SignupInput struct {
	Body struct {
		Email     string `json:"email" format:"email" maxLength:"200"`
		Password  string `json:"password" minLength:"1"`
		Firstname string `json:"firstname" minLength:"1" maxLength:"100"`
		Lastname  string `json:"lastname" minLength:"1" maxLength:"100"`
		Locale    string `json:"locale" required:"false" enum:"en,de,fr,it,"`
	}
}

// TokenInput carries a token from a mail link.
type TokenInput struct {
	Body struct {
		Token string `json:"token" minLength:"1" maxLength:"100"`
	}
}

// ForgotInput is the request of POST /api/auth/forgot.
type ForgotInput struct {
	Body struct {
		Email string `json:"email" maxLength:"200"`
	}
}

// ResetInput is the request of POST /api/auth/reset.
type ResetInput struct {
	Body struct {
		Token    string `json:"token" minLength:"1" maxLength:"100"`
		Password string `json:"password" minLength:"1"`
	}
}

// rateLimited counts a hit on key and returns a 429 once more than max hits fall within window.
func (s *Server) rateLimited(ctx context.Context, key string, window time.Duration, max int) error {
	n, err := s.Store.RateHit(ctx, key, window)
	if err != nil {
		return err
	}
	if n > max {
		return tooMany("Too many attempts. Try again in %d minutes.", int(window.Minutes()))
	}
	return nil
}

// registerAccount registers the sign-up, email verification and password reset endpoints.
func (s *Server) registerAccount(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-auth-options", Method: http.MethodGet, Path: "/api/auth/options", Tags: []string{"auth"},
		Summary: "Whether sign-up and the second factor are enabled",
	}, func(ctx context.Context, _ *struct{}) (*struct{ Body AuthOptions }, error) {
		out := &struct{ Body AuthOptions }{}
		out.Body = AuthOptions{Signup: s.Signup, TOTP: s.TOTP}
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "signup", Method: http.MethodPost, Path: "/api/auth/signup", Tags: []string{"auth"},
		Summary: "Create a personal account and log in", DefaultStatus: http.StatusCreated,
	}, s.signup)

	huma.Register(api, huma.Operation{
		OperationID: "change-password", Method: http.MethodPut, Path: "/api/me/password", Tags: []string{"auth"},
		Summary: "Change the password of the logged-in user; other sessions end", DefaultStatus: http.StatusNoContent,
	}, s.changePassword)

	huma.Register(api, huma.Operation{
		OperationID: "verify-email", Method: http.MethodPost, Path: "/api/auth/verify", Tags: []string{"auth"},
		Summary: "Confirm an email address with the token from the mail", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *TokenInput) (*struct{}, error) {
		id, err := s.Store.UseUserToken(ctx, sessionID(in.Body.Token), "verify")
		if errors.Is(err, store.ErrNotFound) {
			return nil, gone("This link is no longer valid. Request a new one in your settings.")
		}
		if err != nil {
			return nil, err
		}
		return nil, s.Store.VerifyEmail(ctx, id)
	})

	huma.Register(api, huma.Operation{
		OperationID: "forgot-password", Method: http.MethodPost, Path: "/api/auth/forgot", Tags: []string{"auth"},
		Summary: "Send a password reset link if an account has this email", DefaultStatus: http.StatusNoContent,
	}, s.forgot)

	huma.Register(api, huma.Operation{
		OperationID: "reset-password", Method: http.MethodPost, Path: "/api/auth/reset", Tags: []string{"auth"},
		Summary: "Set a new password with the token from the mail; ends all sessions", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *ResetInput) (*struct{}, error) {
		if errs := ValidatePassword(in.Body.Password); len(errs) > 0 {
			return nil, invalidList(errs)
		}
		id, err := s.Store.UseUserToken(ctx, sessionID(in.Body.Token), "reset")
		if errors.Is(err, store.ErrNotFound) {
			return nil, gone("This link is no longer valid. Request a new one.")
		}
		if err != nil {
			return nil, err
		}
		hash, err := HashPassword(in.Body.Password)
		if err != nil {
			return nil, err
		}
		if err := s.Store.SetPassword(ctx, id, hash); err != nil {
			return nil, err
		}
		if _, err := s.Store.DeleteUserSessions(ctx, id); err != nil {
			return nil, err
		}
		return nil, s.Store.VerifyEmail(ctx, id)
	})
}

// PasswordInput is the request of PUT /api/me/password.
type PasswordInput struct {
	Body struct {
		CurrentPassword string `json:"currentPassword" minLength:"1"`
		NewPassword     string `json:"newPassword" minLength:"1" maxLength:"72"`
	}
}

// changePassword sets a new password after checking the current one and ends every other session of the user.
func (s *Server) changePassword(ctx context.Context, in *PasswordInput) (*struct{}, error) {
	u, err := s.checkPassword(ctx, in.Body.CurrentPassword)
	if err != nil {
		return nil, err
	}
	if errs := ValidatePassword(in.Body.NewPassword); len(errs) > 0 {
		return nil, invalidList(errs)
	}
	hash, err := HashPassword(in.Body.NewPassword)
	if err != nil {
		return nil, err
	}
	if err := s.Store.SetPassword(ctx, u.ID, hash); err != nil {
		return nil, err
	}
	keep := ""
	if c, ok := ctx.Value(cookieKey).(string); ok {
		keep = sessionID(c)
	}
	if _, err := s.Store.DeleteOtherSessions(ctx, u.ID, keep); err != nil {
		return nil, err
	}
	s.Log.Info("password changed", "user", u.Name)
	return nil, nil
}

// signup creates a personal account with profile, mails the verification link and logs the new user in.
func (s *Server) signup(ctx context.Context, in *SignupInput) (*LoginOutput, error) {
	if !s.Signup {
		return nil, forbidden("Sign-up is not open on this server.")
	}
	if err := s.rateLimited(ctx, "signup:ip:"+reqInfo(ctx).ip, signupWindow, signupIPMax); err != nil {
		return nil, err
	}
	email := strings.ToLower(strings.TrimSpace(in.Body.Email))
	if errs := ValidatePassword(in.Body.Password); len(errs) > 0 {
		return nil, invalidList(errs)
	}
	hash, err := HashPassword(in.Body.Password)
	if err != nil {
		return nil, err
	}
	id, err := s.Store.CreateAccount(ctx, email, hash, in.Body.Locale)
	if errors.Is(err, store.ErrExists) {
		return nil, conflict("An account with this email exists already. Log in or reset your password.")
	}
	if err != nil {
		return nil, err
	}
	p := store.Profile{Firstname: strings.TrimSpace(in.Body.Firstname), Lastname: strings.TrimSpace(in.Body.Lastname)}
	if err := s.Store.SaveProfile(ctx, id, p); err != nil {
		return nil, err
	}
	if err := s.sendVerification(ctx, id, email, p.Firstname, in.Body.Locale); err != nil {
		return nil, err
	}
	token := randomToken()
	if err := s.Store.CreateSession(ctx, sessionID(token), id, sessionIdle, sessionTTL); err != nil {
		return nil, err
	}
	s.Log.Info("signup", "user", email, "remote", reqInfo(ctx).ip)
	u, err := s.Store.UserByName(ctx, email)
	if err != nil {
		return nil, err
	}
	me, err := s.me(ctx, u)
	return &LoginOutput{SetCookie: sessionCookieFor(ctx, token, int(sessionTTL.Seconds())), Body: me}, err
}

// forgot mails a password reset link when an account has the email and answers the same either way, so emails cannot be probed.
func (s *Server) forgot(ctx context.Context, in *ForgotInput) (*struct{}, error) {
	if err := s.rateLimited(ctx, "forgot:ip:"+reqInfo(ctx).ip, forgotWindow, forgotIPMax); err != nil {
		return nil, err
	}
	u, err := s.Store.UserByEmail(ctx, strings.TrimSpace(in.Body.Email))
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token := randomToken()
	if err := s.Store.CreateUserToken(ctx, sessionID(token), u.ID, "reset", time.Now().Add(resetTTL)); err != nil {
		return nil, err
	}
	l := mailLang(ctx, u.Locale)
	name := ""
	if p, err := s.Store.GetProfile(ctx, u.ID); err == nil {
		name = p.Firstname
	}
	s.mail(mail.Message{To: strings.ToLower(strings.TrimSpace(in.Body.Email)), Subject: i18n.T(l, "Reset your Standblatt password"), Body: strings.Join([]string{
		greeting(l, name), "",
		i18n.T(l, "someone, hopefully you, asked to reset the password of your Standblatt account. Choose a new password here:"), "",
		s.link(ctx, "/reset-password?token="+token), "",
		i18n.T(l, "The link is valid for one hour. If you did not ask for it, ignore this mail; your password stays the same."),
	}, "\n")})
	return nil, nil
}
