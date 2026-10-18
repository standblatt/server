package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/testdb"
	"github.com/standblatt/server/internal/web"
)

const adminToken = "test-admin-token-test-admin-token-test"

type env struct {
	t   *testing.T
	st  *store.Store
	srv *httptest.Server
}

func newEnv(t *testing.T) *env { return newEnvWith(t, nil) }

func newEnvWith(t *testing.T, configure func(*web.Server)) *env {
	st := testdb.New(t)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	hub := &realtime.Hub{DatabaseURL: st.ConnString(), Log: log}
	ctx, cancel := context.WithCancel(context.Background())
	go hub.Run(ctx)
	t.Cleanup(cancel)
	s := &web.Server{
		Store:      st,
		Hub:        hub,
		UI:         fstest.MapFS{"index.html": {Data: []byte("<html>spa</html>")}},
		Version:    "test",
		AdminToken: adminToken,
		TOTP:       true,
		Log:        log,
	}
	if configure != nil {
		configure(s)
	}
	srv := httptest.NewServer(s.Handler())
	t.Cleanup(srv.Close)
	return &env{t: t, st: st, srv: srv}
}

func (e *env) client() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{Jar: jar, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
}

func (e *env) user(name string, tenantID int64, role string) int64 {
	e.t.Helper()
	hash, _ := web.HashPassword("password-" + name)
	id, err := e.st.CreateUser(context.Background(), name, hash)
	if err != nil {
		e.t.Fatal(err)
	}
	if tenantID != 0 {
		roleID, err := e.st.RoleIDByKey(context.Background(), tenantID, role)
		if err != nil {
			e.t.Fatal(err)
		}
		if err := e.st.SetMembership(context.Background(), tenantID, id, roleID); err != nil {
			e.t.Fatal(err)
		}
	}
	return id
}

// role returns the ID of a default role of a tenant.
func (e *env) role(tenantID int64, key string) int64 {
	e.t.Helper()
	id, err := e.st.RoleIDByKey(context.Background(), tenantID, key)
	if err != nil {
		e.t.Fatal(err)
	}
	return id
}

func (e *env) tenant(name string) int64 {
	e.t.Helper()
	id, err := e.st.CreateTenant(context.Background(), name)
	if err != nil {
		e.t.Fatal(err)
	}
	return id
}

func (e *env) do(c *http.Client, method, path string, body any) (int, []byte) {
	e.t.Helper()
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, e.srv.URL+path, r)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b
}

func (e *env) login(name string) *http.Client {
	e.t.Helper()
	c := e.client()
	if code, b := e.do(c, "POST", "/api/auth/login", map[string]string{"name": name, "password": "password-" + strings.ToLower(name)}); code != 200 {
		e.t.Fatalf("login %s: %d %s", name, code, b)
	}
	return c
}

func expect(t *testing.T, what string, code, want int, body []byte) {
	t.Helper()
	if code != want {
		t.Fatalf("%s: got %d, want %d: %s", what, code, want, body)
	}
}

func TestLoginAndSession(t *testing.T) {
	e := newEnv(t)
	tid := e.tenant("Test club")
	e.user("anna", tid, store.RoleAdmin)

	anon := e.client()
	code, b := e.do(anon, "GET", "/api/me", nil)
	expect(t, "me anonymous", code, 401, b)
	code, b = e.do(anon, "POST", "/api/auth/login", map[string]string{"name": "anna", "password": "wrong-password"})
	expect(t, "wrong password", code, 401, b)
	code, b = e.do(anon, "POST", "/api/auth/login", map[string]string{"name": "nobody", "password": "whatever-123"})
	expect(t, "unknown user", code, 401, b)

	c := e.login("Anna")
	code, b = e.do(c, "GET", "/api/me", nil)
	expect(t, "me", code, 200, b)
	var me web.Me
	json.Unmarshal(b, &me)
	if me.User.Name != "anna" || len(me.Tenants) != 1 || me.Tenants[0].TenantID != tid || me.Tenants[0].RoleKey != store.RoleAdmin || len(me.Tenants[0].Permissions) != len(permissions.All) {
		t.Fatalf("unexpected me: %+v", me)
	}

	code, b = e.do(c, "PUT", "/api/me/locale", map[string]string{"locale": "xx"})
	expect(t, "unknown locale", code, 422, b)
	code, b = e.do(c, "PUT", "/api/me/locale", map[string]string{"locale": "fr"})
	expect(t, "set locale", code, 204, b)
	_, b = e.do(c, "GET", "/api/me", nil)
	json.Unmarshal(b, &me)
	if me.User.Locale != "fr" {
		t.Fatalf("locale not stored: %+v", me.User)
	}

	code, b = e.do(c, "POST", "/api/auth/logout", nil)
	expect(t, "logout", code, 204, b)
	code, b = e.do(c, "GET", "/api/me", nil)
	expect(t, "me after logout", code, 401, b)
}

func TestLoginRateLimitAndMaintenance(t *testing.T) {
	e := newEnv(t)
	e.user("bert", 0, "")
	c := e.client()
	for i := 0; i < 5; i++ {
		code, b := e.do(c, "POST", "/api/auth/login", map[string]string{"name": "bert", "password": "wrong-password"})
		expect(t, "wrong password", code, 401, b)
	}
	code, b := e.do(c, "POST", "/api/auth/login", map[string]string{"name": "bert", "password": "password-bert"})
	expect(t, "rate limited", code, 429, b)

	e.user("carla", 0, "")
	e.st.SetMaintenance(context.Background(), true)
	code, b = e.do(e.client(), "POST", "/api/auth/login", map[string]string{"name": "carla", "password": "password-carla"})
	expect(t, "maintenance", code, 503, b)
}

func TestTenantAccess(t *testing.T) {
	e := newEnv(t)
	a, b := e.tenant("Club A"), e.tenant("Club B")
	e.user("admin-a", a, store.RoleAdmin)
	e.user("staff-a", a, store.RoleStaff)
	e.user("admin-b", b, store.RoleAdmin)

	admin := e.login("admin-a")
	staff := e.login("staff-a")
	other := e.login("admin-b")
	base := fmt.Sprintf("/api/tenants/%d", a)

	code, body := e.do(admin, "POST", base+"/ranges", map[string]any{"name": "10m", "laneCount": 10, "disabledLanes": []int{2}})
	expect(t, "create range", code, 201, body)
	var rng store.Range
	json.Unmarshal(body, &rng)

	code, body = e.do(admin, "POST", base+"/ranges", map[string]any{"name": "x", "laneCount": 5, "disabledLanes": []int{9}})
	expect(t, "disabled lane out of range", code, 422, body)
	code, body = e.do(admin, "POST", base+"/ranges", map[string]any{"name": "10m", "laneCount": 5})
	expect(t, "duplicate range", code, 409, body)

	code, body = e.do(admin, "POST", base+"/events", map[string]any{"rangeIds": []int64{rng.ID}, "name": "Indoor", "startsOn": "2026-11-20"})
	expect(t, "create event", code, 201, body)
	var ev store.Event
	json.Unmarshal(body, &ev)
	if ev.Timezone != "Europe/Zurich" {
		t.Fatalf("default timezone not applied: %+v", ev)
	}

	code, body = e.do(admin, "PUT", fmt.Sprintf("%s/ranges/%d", base, rng.ID), map[string]any{"name": "10m hall", "laneCount": 12})
	expect(t, "update range", code, 204, body)
	code, body = e.do(admin, "PUT", fmt.Sprintf("%s/events/%d", base, ev.ID), map[string]any{"rangeIds": []int64{rng.ID}, "name": "Indoor 2026", "startsOn": "2026-11-21", "timezone": "Mars/Base"})
	expect(t, "invalid timezone", code, 422, body)
	code, body = e.do(admin, "PUT", fmt.Sprintf("%s/events/%d", base, ev.ID), map[string]any{"rangeIds": []int64{rng.ID}, "name": "Indoor 2026", "startsOn": "2026-11-21"})
	expect(t, "update event", code, 204, body)
	code, body = e.do(staff, "GET", fmt.Sprintf("%s/events/%d", base, ev.ID), nil)
	expect(t, "staff reads event", code, 200, body)
	json.Unmarshal(body, &ev)
	if ev.Name != "Indoor 2026" || ev.StartsOn != "2026-11-21" || ev.Ranges[0].Name != "10m hall" {
		t.Fatalf("update not applied: %+v", ev)
	}

	code, body = e.do(staff, "GET", base+"/events", nil)
	expect(t, "staff lists events", code, 200, body)
	code, body = e.do(staff, "POST", base+"/ranges", map[string]any{"name": "25m", "laneCount": 5})
	expect(t, "staff creates range", code, 403, body)
	code, body = e.do(staff, "GET", base+"/users", nil)
	expect(t, "staff lists users", code, 403, body)

	code, body = e.do(other, "GET", base+"/events", nil)
	expect(t, "other tenant lists events", code, 404, body)
	code, body = e.do(other, "GET", fmt.Sprintf("%s/events/%d", base, ev.ID), nil)
	expect(t, "other tenant reads event", code, 404, body)
	code, body = e.do(other, "DELETE", fmt.Sprintf("%s/events/%d", base, ev.ID), nil)
	expect(t, "other tenant deletes event", code, 404, body)
	code, body = e.do(other, "GET", fmt.Sprintf("/api/tenants/%d/events/%d", b, ev.ID), nil)
	expect(t, "event via own tenant id", code, 404, body)

	code, body = e.do(e.client(), "GET", base+"/events", nil)
	expect(t, "anonymous", code, 401, body)
}

func TestMembers(t *testing.T) {
	e := newEnv(t)
	tid := e.tenant("Club")
	adminID := e.user("boss", tid, store.RoleAdmin)
	c := e.login("boss")
	base := fmt.Sprintf("/api/tenants/%d/users", tid)

	code, b := e.do(c, "POST", base, map[string]any{"name": "Helper1", "roleId": e.role(tid, "staff")})
	expect(t, "new user without password", code, 422, b)
	code, b = e.do(c, "POST", base, map[string]any{"name": "Helper1", "roleId": e.role(tid, "staff"), "password": "helper-password"})
	expect(t, "create member", code, 204, b)
	h := e.client()
	code, b = e.do(h, "POST", "/api/auth/login", map[string]string{"name": "helper1", "password": "helper-password"})
	expect(t, "new member logs in", code, 200, b)

	code, b = e.do(c, "PUT", fmt.Sprintf("%s/%d", base, adminID), map[string]any{"roleId": e.role(tid, "staff")})
	expect(t, "demote last admin", code, 409, b)
	code, b = e.do(c, "DELETE", fmt.Sprintf("%s/%d", base, adminID), nil)
	expect(t, "remove last admin", code, 409, b)

	code, b = e.do(c, "GET", base, nil)
	expect(t, "list members", code, 200, b)
	var members []store.Member
	json.Unmarshal(b, &members)
	if len(members) != 2 {
		t.Fatalf("unexpected members: %+v", members)
	}
}

func TestAdminLoginRateLimit(t *testing.T) {
	e := newEnv(t)
	c := e.client()
	login := func(token string) int {
		resp, err := c.PostForm(e.srv.URL+"/admin/login", url.Values{"token": {token}})
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		return resp.StatusCode
	}

	for i := 0; i < 15; i++ {
		if code := login(adminToken); code != 303 {
			t.Fatalf("successful login %d counted against the limit: %d", i, code)
		}
	}
	for i := 0; i < 10; i++ {
		if code := login("wrong"); code != 401 {
			t.Fatalf("failed login %d: %d", i, code)
		}
	}
	if code := login(adminToken); code != 429 {
		t.Fatalf("login after 10 failures: %d", code)
	}
}

func TestAdminPanel(t *testing.T) {
	e, box := personalEnv(t, nil)
	c := e.client()
	form := func(path string, v url.Values) (int, string) {
		resp, err := c.PostForm(e.srv.URL+path, v)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, string(b)
	}
	get := func(path string) string {
		resp, err := c.Get(e.srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return string(b)
	}

	var body string
	code, _ := form("/admin/clubs", url.Values{"name": {"Club X"}})
	if code != 401 {
		t.Fatalf("admin action without cookie: %d", code)
	}
	if code, _ = form("/admin/login", url.Values{"token": {"wrong"}}); code != 401 {
		t.Fatalf("wrong token: %d", code)
	}
	if code, _ = form("/admin/login", url.Values{"token": {adminToken}}); code != 303 {
		t.Fatalf("admin login: %d", code)
	}
	if code, _ = form("/admin/clubs", url.Values{"name": {"Club X"}}); code != 303 {
		t.Fatalf("create tenant: %d", code)
	}
	if code, _ = form("/admin/users", url.Values{"name": {"xaver"}, "password": {"xaver-password"}}); code != 303 {
		t.Fatalf("create user: %d", code)
	}
	tenants, _ := e.st.ListTenants(context.Background())
	club := fmt.Sprintf("/admin/clubs/%d", tenants[0].ID)
	if code, body = form(club+"/members", url.Values{"user": {"xaver"}, "role": {"admin"}}); code != 303 {
		t.Fatalf("assign membership: %d %s", code, body)
	}
	if b := get("/admin/clubs"); !strings.Contains(b, "Club X") {
		t.Fatalf("clubs page misses the club: %s", b)
	}
	if b := get(club); !strings.Contains(b, `<strong>xaver</strong></a>`) || !strings.Contains(b, `<span class="badge admin">admin</span>`) {
		t.Fatalf("club page misses the member: %s", b)
	}

	u, _ := e.st.UserByName(context.Background(), "xaver")
	user := fmt.Sprintf("/admin/users/%d", u.ID)
	if code, body = form(user+"/password-link", nil); code != 422 || !strings.Contains(body, "no email address") {
		t.Fatalf("password link without email: %d %s", code, body)
	}
	if code, _ = form(user+"/email", url.Values{"email": {"not-an-address"}}); code != 422 {
		t.Fatalf("invalid email: %d", code)
	}
	if code, _ = form(user+"/email", url.Values{"email": {"Xaver@Example.org"}, "verified": {"on"}}); code != 303 {
		t.Fatalf("set email: %d", code)
	}
	if a, _ := e.st.GetAccount(context.Background(), u.ID); a.Email != "xaver@example.org" || !a.EmailVerified {
		t.Fatalf("account after email: %+v", a)
	}
	if b := get("/admin/users?q=xample.org"); !strings.Contains(b, "xaver@example.org") {
		t.Fatalf("user search by email: %s", b)
	}
	if b := get("/admin/users?q=nobody"); strings.Contains(b, "xaver") {
		t.Fatalf("user search should be empty: %s", b)
	}
	if code, body = form(user+"/password-link", nil); code != 200 || !strings.Contains(body, "on its way") {
		t.Fatalf("password link: %d %s", code, body)
	}
	m := box.next(t)
	if m.To != "xaver@example.org" || !strings.Contains(m.Body, "/reset-password?token=") {
		t.Fatalf("password link mail: %+v", m)
	}
	code, b2 := e.do(e.client(), "POST", "/api/auth/reset", map[string]string{"token": tokenRe.FindStringSubmatch(m.Body)[1], "password": "xaver-new-password"})
	expect(t, "reset with the admin link", code, 204, b2)

	code, b2 = e.do(e.client(), "POST", "/api/auth/login", map[string]string{"name": "xaver", "password": "xaver-new-password"})
	expect(t, "login of admin-created user", code, 200, b2)

	// a user created with an email but no password gets the link right away
	if code, body = form("/admin/users", url.Values{"name": {"yves"}, "email": {"yves@example.org"}}); code != 303 {
		t.Fatalf("create user by mail: %d %s", code, body)
	}
	if m := box.next(t); m.To != "yves@example.org" {
		t.Fatalf("mail for the new user: %+v", m)
	}

	req, _ := http.NewRequest("GET", e.srv.URL+"/admin/system", nil)
	req.Header.Set("Accept-Language", "de-CH")
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if !strings.Contains(string(b), "Wartungsmodus") {
		t.Fatalf("german panel: %s", b)
	}
}

func TestSPAFallback(t *testing.T) {
	e := newEnv(t)
	c := e.client()
	if code, b := e.do(c, "GET", "/t/1/events", nil); code != 200 || !strings.Contains(string(b), "spa") {
		t.Errorf("app route: %d %s", code, b)
	}
	if code, _ := e.do(c, "GET", "/_app/immutable/entry/start.old.js", nil); code != 404 {
		t.Errorf("missing asset: got %d, want 404", code)
	}
}
