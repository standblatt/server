package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/i18n"
	tlmail "github.com/standblatt/server/internal/mail"
	"github.com/standblatt/server/internal/store"
)

const (
	adminCookie     = "standblatt_admin"
	adminLangCookie = "standblatt_admin_lang"
	adminTTL        = 8 * time.Hour
	adminLabel      = "standblatt-admin-v1:"
	adminWindow     = 15 * time.Minute
	adminIPMax      = 10
)

//go:embed templates/*.html
var templateFS embed.FS

// adminFuncs are the template functions; T and Tf are replaced per request with the panel's language.
var adminFuncs = template.FuncMap{
	"local":   func(t *time.Time) time.Time { return t.In(zone("")) },
	"T":       func(s string) string { return s },
	"Tf":      func(f string, a ...any) string { return fmt.Sprintf(f, a...) },
	"problem": func(p domain.Problem) string { return p.String() },
	"langNames": func() map[string]string {
		return map[string]string{"en": "English", "de": "Deutsch", "fr": "Français", "it": "Italiano"}
	},
}

var adminTmpl = template.Must(template.New("").Funcs(adminFuncs).ParseFS(templateFS, "templates/*.html"))

// signAdmin returns the admin cookie value: the expiry in Unix seconds and its HMAC-SHA256 keyed with the admin token.
func signAdmin(token string, exp time.Time) string {
	unix := strconv.FormatInt(exp.Unix(), 10)
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write([]byte(adminLabel + unix))
	return unix + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// validAdmin reports whether an admin cookie value carries a valid signature and an expiry that is in the future but within adminTTL.
func validAdmin(token, value string, now time.Time) bool {
	unix, sig, ok := strings.Cut(value, ".")
	if !ok {
		return false
	}
	exp, err := strconv.ParseInt(unix, 10, 64)
	if err != nil || now.Unix() >= exp || exp-now.Unix() > int64(adminTTL/time.Second) {
		return false
	}
	want, err := base64.RawURLEncoding.Strict().DecodeString(sig)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write([]byte(adminLabel + unix))
	return hmac.Equal(mac.Sum(nil), want)
}

// tokenMatches compares the posted admin token with the configured one in constant time.
func tokenMatches(configured, posted string) bool {
	a, b := sha256.Sum256([]byte(configured)), sha256.Sum256([]byte(posted))
	return subtle.ConstantTimeCompare(a[:], b[:]) == 1
}

// adminData is what the admin templates render; each page fills the part it shows.
type adminData struct {
	Page   string
	Lang   string
	Errors []domain.Problem
	Notice string

	Stats       store.AdminStats
	Tenants     []store.Tenant
	Tenant      store.Tenant
	Members     []store.Member
	Events      []store.Event
	Users       []store.UserStat
	Query       string
	User        store.UserStat
	Memberships []store.Membership
	Templates   []store.ProgrammeTemplate
	Maintenance bool
	Version     string
	Environment string
	Signup      bool
	TOTP        bool
}

// registerAdmin registers the admin panel routes; everything but the index, login, language and logout needs the admin cookie.
func (s *Server) registerAdmin(mux *http.ServeMux) {
	mux.HandleFunc("GET /admin", s.adminIndex)
	mux.HandleFunc("POST /admin/login", s.adminLogin)
	mux.HandleFunc("POST /admin/logout", s.adminLogout)
	mux.HandleFunc("POST /admin/lang", s.adminLang)
	h := func(pattern string, fn http.HandlerFunc) { mux.Handle(pattern, s.requireAdmin(fn)) }
	for _, page := range []string{"clubs", "users", "library", "system"} {
		h("GET /admin/"+page, s.adminPageHandler(page))
	}
	h("GET /admin/clubs/{id}", s.adminPageHandler("club"))
	h("GET /admin/users/{id}", s.adminPageHandler("user"))
	h("POST /admin/clubs", s.adminCreateTenant)
	h("POST /admin/clubs/{id}/delete", s.adminDeleteTenant)
	h("POST /admin/clubs/{id}/members", s.adminSetMembership)
	h("POST /admin/clubs/{id}/members/delete", s.adminRemoveMembership)
	h("POST /admin/users", s.adminCreateUser)
	h("POST /admin/users/{id}/email", s.adminSetEmail)
	h("POST /admin/users/{id}/password", s.adminSetPassword)
	h("POST /admin/users/{id}/password-link", s.adminPasswordLink)
	h("POST /admin/users/{id}/sessions/delete", s.adminRevokeSessions)
	h("POST /admin/users/{id}/totp/delete", s.adminResetTOTP)
	h("POST /admin/users/{id}/memberships/delete", s.adminRemoveMembership)
	h("POST /admin/users/{id}/delete", s.adminDeleteUser)
	h("POST /admin/sessions/delete", s.adminRevokeAllSessions)
	h("POST /admin/maintenance", s.adminMaintenance)
	h("POST /admin/templates", s.adminUploadTemplate)
	h("GET /admin/templates/{id}", s.adminDownloadTemplate)
	h("POST /admin/templates/{id}/delete", s.adminDeleteTemplate)
}

// isAdmin reports whether the request carries a valid admin cookie.
func (s *Server) isAdmin(r *http.Request) bool {
	c, err := r.Cookie(adminCookie)
	return err == nil && validAdmin(s.AdminToken, c.Value, time.Now())
}

// requireAdmin wraps next so that requests without a valid admin cookie get the login page with a 401.
func (s *Server) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdmin(r) {
			s.adminRender(w, r, http.StatusUnauthorized, "admin_login", adminData{Errors: []domain.Problem{domain.Problemf("Session expired. Enter the token again.")}})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// adminLocale is the panel's language: the language cookie, else the browser's Accept-Language.
func adminLocale(r *http.Request) string {
	if c, err := r.Cookie(adminLangCookie); err == nil && (c.Value == "de" || c.Value == "fr" || c.Value == "it") {
		return c.Value
	}
	if l := acceptLang(r.Header.Get("Accept-Language")); l == "de" || l == "fr" || l == "it" {
		return l
	}
	return "en"
}

// adminRender writes the admin template name with status in the panel's language, uncached.
func (s *Server) adminRender(w http.ResponseWriter, r *http.Request, status int, name string, d adminData) {
	d.Lang = adminLocale(r)
	l := i18n.Get(d.Lang)
	tmpl, err := adminTmpl.Clone()
	if err != nil {
		s.adminFail(w, err)
		return
	}
	tmpl.Funcs(template.FuncMap{
		"T":       func(s string) string { return i18n.T(l, s) },
		"Tf":      func(f string, a ...any) string { return i18n.Tf(l, f, a...) },
		"problem": func(p domain.Problem) string { return i18n.Tf(l, p.Format, p.Args...) },
	})
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, name+".html", d); err != nil {
		s.Log.Error("admin template", "err", err)
	}
}

// adminFail logs err and answers a plain 500.
func (s *Server) adminFail(w http.ResponseWriter, err error) {
	s.Log.Error("admin", "err", err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}

// adminIndex shows the login page to visitors and the overview to logged-in admins.
func (s *Server) adminIndex(w http.ResponseWriter, r *http.Request) {
	if !s.isAdmin(r) {
		s.adminRender(w, r, http.StatusOK, "admin_login", adminData{})
		return
	}
	s.adminShow(w, r, http.StatusOK, "overview", adminData{})
}

// adminPageHandler returns the handler that shows page.
func (s *Server) adminPageHandler(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { s.adminShow(w, r, http.StatusOK, page, adminData{}) }
}

// adminShow loads the data of page (the {id} path value names the club or user) into d and renders it with status.
func (s *Server) adminShow(w http.ResponseWriter, r *http.Request, status int, page string, d adminData) {
	ctx := r.Context()
	d.Page = page
	var err error
	switch page {
	case "overview":
		if d.Stats, err = s.Store.AdminStats(ctx); err == nil {
			d.Maintenance, err = s.Store.Maintenance(ctx)
		}
		d.Version, d.Environment = s.Version, s.Environment
	case "clubs":
		d.Tenants, err = s.Store.ListTenants(ctx)
	case "club":
		if d.Tenant, err = s.Store.GetTenant(ctx, pathID(r)); err == nil {
			if d.Members, err = s.Store.Members(ctx, d.Tenant.ID); err == nil {
				d.Events, err = s.Store.ListEvents(ctx, d.Tenant.ID)
			}
		}
	case "users":
		d.Query = strings.TrimSpace(r.URL.Query().Get("q"))
		d.Users, err = s.Store.ListUsersWithStats(ctx, d.Query)
	case "user":
		if d.User, err = s.Store.UserStats(ctx, pathID(r)); err == nil {
			d.Memberships, err = s.Store.Memberships(ctx, d.User.ID)
		}
	case "library":
		d.Templates, err = s.Store.ListTemplates(ctx)
	case "system":
		d.Maintenance, err = s.Store.Maintenance(ctx)
		d.Version, d.Environment, d.Signup, d.TOTP = s.Version, s.Environment, s.Signup, s.TOTP
	}
	if errors.Is(err, store.ErrNotFound) {
		s.adminShow(w, r, http.StatusNotFound, "overview", adminData{Errors: []domain.Problem{domain.Problemf("Not found")}})
		return
	}
	if err != nil {
		s.adminFail(w, err)
		return
	}
	s.adminRender(w, r, status, "admin_"+page, d)
}

// adminLogin checks the posted admin token under a per-IP rate limit and sets the signed admin cookie on success.
func (s *Server) adminLogin(w http.ResponseWriter, r *http.Request) {
	ip := s.clientIP(r)
	key := "admin:ip:" + ip
	n, err := s.Store.RateCount(r.Context(), key, adminWindow)
	if err != nil {
		s.adminFail(w, err)
		return
	}
	if n >= adminIPMax {
		s.Log.Warn("admin login rate limited", "remote", ip)
		s.adminRender(w, r, http.StatusTooManyRequests, "admin_login",
			adminData{Errors: []domain.Problem{domain.Problemf("Too many attempts. Try again in %d minutes.", int(adminWindow.Minutes()))}})
		return
	}
	if !tokenMatches(s.AdminToken, r.PostFormValue("token")) {
		if _, err := s.Store.RateHit(r.Context(), key, adminWindow); err != nil {
			s.adminFail(w, err)
			return
		}
		s.Log.Warn("admin login failed", "remote", ip)
		s.adminRender(w, r, http.StatusUnauthorized, "admin_login", adminData{Errors: []domain.Problem{domain.Problemf("Wrong token.")}})
		return
	}
	exp := time.Now().Add(adminTTL)
	http.SetCookie(w, &http.Cookie{
		Name: adminCookie, Value: signAdmin(s.AdminToken, exp), Path: "/admin",
		HttpOnly: true, Secure: secureRequest(r), SameSite: http.SameSiteStrictMode, MaxAge: int(adminTTL.Seconds()),
	})
	s.Log.Info("admin login", "remote", ip)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// adminLogout clears the admin cookie and redirects to the login page.
func (s *Server) adminLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: adminCookie, Value: "", Path: "/admin", HttpOnly: true, Secure: secureRequest(r), SameSite: http.SameSiteStrictMode, MaxAge: -1})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// adminLang stores the chosen panel language in a cookie and returns to the page the form came from.
func (s *Server) adminLang(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: adminLangCookie, Value: r.PostFormValue("lang"), Path: "/admin",
		HttpOnly: true, Secure: secureRequest(r), SameSite: http.SameSiteStrictMode, MaxAge: 365 * 24 * 3600,
	})
	http.Redirect(w, r, adminBack(r, "/admin"), http.StatusSeeOther)
}

// adminBack is the panel path to return to after an action: the posted "back" field when it is a panel path, else fallback.
func adminBack(r *http.Request, fallback string) string {
	if b := r.PostFormValue("back"); strings.HasPrefix(b, "/admin") && !strings.Contains(b, "//") {
		return b
	}
	return fallback
}

// adminLog logs an admin action with its target and the client IP.
func (s *Server) adminLog(r *http.Request, action, target string) {
	s.Log.Info("admin action", "action", action, "target", target, "remote", s.clientIP(r))
}

// pathID returns the {id} path value as int64, 0 when it is not a number.
func pathID(r *http.Request) int64 {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id
}

// formID returns the form value key as int64, 0 when it is not a number.
func formID(r *http.Request, key string) int64 {
	id, _ := strconv.ParseInt(r.PostFormValue(key), 10, 64)
	return id
}

// adminDone finishes an action: 404 on the overview for ErrNotFound, 500 for other errors, else a redirect to back.
func (s *Server) adminDone(w http.ResponseWriter, r *http.Request, err error, back string) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		s.adminShow(w, r, http.StatusNotFound, "overview", adminData{Errors: []domain.Problem{domain.Problemf("Not found")}})
	case err != nil:
		s.adminFail(w, err)
	default:
		http.Redirect(w, r, back, http.StatusSeeOther)
	}
}

// adminInvalid renders page again with a 422 and the problems.
func (s *Server) adminInvalid(w http.ResponseWriter, r *http.Request, page string, errs ...domain.Problem) {
	s.adminShow(w, r, http.StatusUnprocessableEntity, page, adminData{Errors: errs})
}

// adminCreateTenant creates a club with the posted name, 422 when the name is empty or taken.
func (s *Server) adminCreateTenant(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.PostFormValue("name"))
	if name == "" {
		s.adminInvalid(w, r, "clubs", domain.Problemf("Name is missing"))
		return
	}
	id, err := s.Store.CreateTenant(r.Context(), name)
	if errors.Is(err, store.ErrExists) {
		s.adminInvalid(w, r, "clubs", domain.Problemf("Club already exists"))
		return
	} else if err != nil {
		s.adminFail(w, err)
		return
	}
	s.adminLog(r, "tenant created", name)
	http.Redirect(w, r, fmt.Sprintf("/admin/clubs/%d", id), http.StatusSeeOther)
}

// adminDeleteTenant deletes the club in the path with all its data.
func (s *Server) adminDeleteTenant(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	err := s.Store.DeleteTenant(r.Context(), id)
	if err == nil {
		s.adminLog(r, "tenant deleted", strconv.FormatInt(id, 10))
	}
	s.adminDone(w, r, err, "/admin/clubs")
}

// adminSetMembership gives the posted user (by name or email) the posted role in the club in the path.
func (s *Server) adminSetMembership(w http.ResponseWriter, r *http.Request) {
	tenantID := pathID(r)
	role := store.Role(r.PostFormValue("role"))
	name := strings.ToLower(strings.TrimSpace(r.PostFormValue("user")))
	if !role.Valid() {
		s.adminInvalid(w, r, "club", domain.Problemf("Unknown role"))
		return
	}
	u, err := s.Store.UserByName(r.Context(), name)
	if errors.Is(err, store.ErrNotFound) {
		u, err = s.Store.UserByEmail(r.Context(), name)
	}
	if errors.Is(err, store.ErrNotFound) {
		s.adminInvalid(w, r, "club", domain.Problemf("No such user"))
		return
	} else if err != nil {
		s.adminFail(w, err)
		return
	}
	if err := s.Store.SetMembership(r.Context(), tenantID, u.ID, role); errors.Is(err, store.ErrInUse) {
		s.adminDone(w, r, store.ErrNotFound, "")
		return
	} else if err != nil {
		s.adminFail(w, err)
		return
	}
	s.adminLog(r, "membership set", fmt.Sprintf("%s@%d=%s", name, tenantID, role))
	http.Redirect(w, r, fmt.Sprintf("/admin/clubs/%d", tenantID), http.StatusSeeOther)
}

// adminRemoveMembership removes a user from a club; the club comes from the path on club pages and from the form on user pages.
func (s *Server) adminRemoveMembership(w http.ResponseWriter, r *http.Request) {
	tenantID, userID := pathID(r), formID(r, "user")
	back := fmt.Sprintf("/admin/clubs/%d", tenantID)
	if strings.HasPrefix(r.URL.Path, "/admin/users/") {
		tenantID, userID = formID(r, "tenant"), pathID(r)
		back = fmt.Sprintf("/admin/users/%d", userID)
	}
	err := s.Store.RemoveMembership(r.Context(), tenantID, userID)
	if err == nil {
		s.adminLog(r, "membership removed", fmt.Sprintf("%d@%d", userID, tenantID))
	}
	s.adminDone(w, r, err, back)
}

// adminCreateUser creates a user with a name, an optional email and a password; without a password the user gets a mail with a link to set one.
func (s *Server) adminCreateUser(w http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(strings.TrimSpace(r.PostFormValue("name")))
	email := strings.ToLower(strings.TrimSpace(r.PostFormValue("email")))
	password := r.PostFormValue("password")
	errs := ValidateUserName(name)
	if email != "" {
		if a, err := mail.ParseAddress(email); err != nil || a.Address != email {
			errs = append(errs, domain.Problemf("Not an email address"))
		}
	}
	sendLink := password == "" && email != ""
	if sendLink {
		password = randomToken()
	}
	errs = append(errs, ValidatePassword(password)...)
	if len(errs) > 0 {
		s.adminInvalid(w, r, "users", errs...)
		return
	}
	hash, err := HashPassword(password)
	if err != nil {
		s.adminFail(w, err)
		return
	}
	id, err := s.Store.CreateUser(r.Context(), name, hash)
	if errors.Is(err, store.ErrExists) {
		s.adminInvalid(w, r, "users", domain.Problemf("User name is already taken"))
		return
	} else if err != nil {
		s.adminFail(w, err)
		return
	}
	if email != "" {
		if err := s.Store.SetUserEmail(r.Context(), id, email, false); errors.Is(err, store.ErrExists) {
			s.adminInvalid(w, r, "users", domain.Problemf("Another account uses this email address"))
			return
		} else if err != nil {
			s.adminFail(w, err)
			return
		}
	}
	s.adminLog(r, "user created", name)
	if sendLink {
		if err := s.sendPasswordLink(r, id); err != nil {
			s.adminFail(w, err)
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/admin/users/%d", id), http.StatusSeeOther)
}

// adminSetEmail sets, verifies or clears the email of the user in the path.
func (s *Server) adminSetEmail(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	email := strings.ToLower(strings.TrimSpace(r.PostFormValue("email")))
	if email != "" {
		if a, err := mail.ParseAddress(email); err != nil || a.Address != email {
			s.adminInvalid(w, r, "user", domain.Problemf("Not an email address"))
			return
		}
	}
	err := s.Store.SetUserEmail(r.Context(), id, email, r.PostFormValue("verified") == "on")
	if errors.Is(err, store.ErrExists) {
		s.adminInvalid(w, r, "user", domain.Problemf("Another account uses this email address"))
		return
	}
	if err == nil {
		s.adminLog(r, "email set", strconv.FormatInt(id, 10))
	}
	s.adminDone(w, r, err, fmt.Sprintf("/admin/users/%d", id))
}

// adminSetPassword sets a new validated password for the user in the path.
func (s *Server) adminSetPassword(w http.ResponseWriter, r *http.Request) {
	password := r.PostFormValue("password")
	if errs := ValidatePassword(password); len(errs) > 0 {
		s.adminInvalid(w, r, "user", errs...)
		return
	}
	hash, err := HashPassword(password)
	if err != nil {
		s.adminFail(w, err)
		return
	}
	id := pathID(r)
	err = s.Store.SetPassword(r.Context(), id, hash)
	if err == nil {
		s.adminLog(r, "password set", strconv.FormatInt(id, 10))
	}
	s.adminDone(w, r, err, fmt.Sprintf("/admin/users/%d", id))
}

// adminPasswordLink mails the user in the path a link to set a new password.
func (s *Server) adminPasswordLink(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	err := s.sendPasswordLink(r, id)
	if errors.Is(err, errNoEmail) {
		s.adminInvalid(w, r, "user", domain.Problemf("The user has no email address"))
		return
	}
	if err != nil {
		s.adminDone(w, r, err, "")
		return
	}
	s.adminShow(w, r, http.StatusOK, "user", adminData{Notice: i18n.T(i18n.Get(adminLocale(r)), "The link to set the password is on its way")})
}

// errNoEmail says a mail could not be sent because the account has no address.
var errNoEmail = errors.New("no email")

// sendPasswordLink creates a reset token valid for inviteTTL and mails it to the user, in the user's language or the panel's.
func (s *Server) sendPasswordLink(r *http.Request, userID int64) error {
	ctx := r.Context()
	a, err := s.Store.GetAccount(ctx, userID)
	if err != nil {
		return err
	}
	if a.Email == "" {
		return errNoEmail
	}
	u, err := s.Store.UserStats(ctx, userID)
	if err != nil {
		return err
	}
	token := randomToken()
	if err := s.Store.CreateUserToken(ctx, sessionID(token), userID, "reset", time.Now().Add(inviteTTL)); err != nil {
		return err
	}
	l := i18n.Get(adminLocale(r))
	if u.Locale != "" {
		l = i18n.Get(u.Locale)
	}
	s.mail(tlmail.Message{To: a.Email, Subject: i18n.T(l, "Set your Standblatt password"), Body: strings.Join([]string{
		greeting(l, s.firstname(ctx, userID)), "",
		i18n.T(l, "an administrator created a link for you to set the password of your Standblatt account:"), "",
		s.link(ctx, "/reset-password?token="+token), "",
		i18n.Tf(l, "The link is valid for seven days. Your user name is %s.", u.Name),
	}, "\n")})
	s.adminLog(r, "password link sent", strconv.FormatInt(userID, 10))
	return nil
}

// adminResetTOTP switches off the second factor of the user in the path.
func (s *Server) adminResetTOTP(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	err := s.Store.SetTOTP(r.Context(), id, "")
	if err == nil {
		s.adminLog(r, "2fa reset", strconv.FormatInt(id, 10))
	}
	s.adminDone(w, r, err, fmt.Sprintf("/admin/users/%d", id))
}

// adminRevokeSessions deletes all sessions of the user in the path and shows how many were revoked.
func (s *Server) adminRevokeSessions(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	n, err := s.Store.DeleteUserSessions(r.Context(), id)
	if err != nil {
		s.adminFail(w, err)
		return
	}
	s.adminLog(r, "sessions revoked", strconv.FormatInt(id, 10))
	s.adminShow(w, r, http.StatusOK, "user", adminData{Notice: i18n.Tf(i18n.Get(adminLocale(r)), "%d session(s) revoked", n)})
}

// adminDeleteUser deletes the user in the path.
func (s *Server) adminDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	err := s.Store.DeleteUser(r.Context(), id)
	if err == nil {
		s.adminLog(r, "user deleted", strconv.FormatInt(id, 10))
	}
	s.adminDone(w, r, err, "/admin/users")
}

// adminRevokeAllSessions deletes every session, which logs everyone out.
func (s *Server) adminRevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	n, err := s.Store.DeleteAllSessions(r.Context())
	if err != nil {
		s.adminFail(w, err)
		return
	}
	s.adminLog(r, "all sessions revoked", "*")
	s.adminShow(w, r, http.StatusOK, "system", adminData{Notice: i18n.Tf(i18n.Get(adminLocale(r)), "%d session(s) revoked, everyone is logged out", n)})
}

// adminMaintenance switches maintenance mode on or off.
func (s *Server) adminMaintenance(w http.ResponseWriter, r *http.Request) {
	on := r.PostFormValue("on") == "true"
	if err := s.Store.SetMaintenance(r.Context(), on); err != nil {
		s.adminFail(w, err)
		return
	}
	s.adminLog(r, map[bool]string{true: "maintenance on", false: "maintenance off"}[on], "*")
	http.Redirect(w, r, adminBack(r, "/admin/system"), http.StatusSeeOther)
}

// adminUploadTemplate validates an uploaded programme export file (at most 1 MiB) and saves it to the library.
func (s *Server) adminUploadTemplate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	f, _, err := r.FormFile("file")
	if err != nil {
		s.adminInvalid(w, r, "library", domain.Problemf("Choose a programme file"))
		return
	}
	defer f.Close()
	var doc ProgrammeExport
	if err := json.NewDecoder(f).Decode(&doc); err != nil || doc.Format != ProgrammeFormat && doc.Format != legacyProgrammeFormat {
		s.adminInvalid(w, r, "library", domain.Problemf("Not a programme export file"))
		return
	}
	doc.Programme.Normalize()
	if errs := doc.Programme.Validate(); len(errs) > 0 {
		s.adminInvalid(w, r, "library", errs...)
		return
	}
	if _, err := s.Store.SaveTemplate(r.Context(), doc.Programme); err != nil {
		s.adminFail(w, err)
		return
	}
	s.adminLog(r, "template saved", doc.Programme.Name+" "+doc.Programme.Variant)
	http.Redirect(w, r, "/admin/library", http.StatusSeeOther)
}

// adminDownloadTemplate sends a library template as a programme export file.
func (s *Server) adminDownloadTemplate(w http.ResponseWriter, r *http.Request) {
	tpl, err := s.Store.GetTemplate(r.Context(), pathID(r))
	if err != nil {
		s.adminDone(w, r, err, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="`+exportFilename(tpl.Definition)+`"`)
	if err := json.NewEncoder(w).Encode(ProgrammeExport{Format: ProgrammeFormat, Version: 1, Programme: tpl.Definition}); err != nil {
		s.Log.Warn("template download", "err", err)
	}
}

// adminDeleteTemplate removes the template in the path from the library.
func (s *Server) adminDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	err := s.Store.DeleteTemplate(r.Context(), id)
	if err == nil {
		s.adminLog(r, "template deleted", strconv.FormatInt(id, 10))
	}
	s.adminDone(w, r, err, "/admin/library")
}
