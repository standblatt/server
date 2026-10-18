package web_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

type headerTransport struct{ key, value string }

func (h headerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set(h.key, h.value)
	return http.DefaultTransport.RoundTrip(r)
}

func withToken(token string) *http.Client {
	return &http.Client{Transport: headerTransport{"Authorization", "Bearer " + token}}
}

func TestAPITokens(t *testing.T) {
	e, _ := personalEnv(t, nil)
	c := e.signup("tina@example.org")
	other := e.signup("otto@example.org")

	create := func(cl *http.Client, scopes ...string) web.CreatedAPIToken {
		code, b := e.do(cl, "POST", "/api/me/tokens", web.NewAPIToken{Name: "script", Scopes: scopes, ExpiresInDays: 30})
		expect(t, "create token", code, 201, b)
		var tok web.CreatedAPIToken
		json.Unmarshal(b, &tok)
		return tok
	}
	ro, rw := create(c, "training:read"), create(c, "training:read", "training:write")
	if !strings.HasPrefix(ro.Token, "sb_") || len(ro.Token) != 43 || ro.ExpiresAt == nil {
		t.Fatalf("token: %+v", ro)
	}
	code, b := e.do(c, "POST", "/api/me/tokens", web.NewAPIToken{Name: "x", Scopes: []string{"club:write"}})
	expect(t, "unknown scope", code, 422, b)

	code, b = e.do(c, "GET", "/api/me/tokens", nil)
	expect(t, "list", code, 200, b)
	if strings.Contains(string(b), "sb_") || strings.Count(string(b), `"id"`) != 2 {
		t.Fatalf("list must show two tokens without secrets: %s", b)
	}

	o := overview(t, e, withToken(ro.Token))
	in := store.TrainingInput{DisciplineID: o.Disciplines[0].ID, Kind: "training", StartedAt: time.Now(), Series: []domain.TrainingSeries{{Shots: 10, Total: 90}}}
	code, b = e.do(withToken(ro.Token), "POST", "/api/my/trainings", in)
	expect(t, "write with a read token", code, 403, b)
	code, b = e.do(withToken(rw.Token), "POST", "/api/my/trainings", in)
	expect(t, "write with a write token", code, 201, b)
	keyClient := &http.Client{Transport: headerTransport{"X-API-Key", ro.Token}}
	code, b = e.do(keyClient, "GET", "/api/my/trainings", nil)
	expect(t, "X-API-Key", code, 200, b)

	for _, p := range []struct{ method, path string }{{"GET", "/api/me/tokens"}, {"GET", "/api/me"}, {"PUT", "/api/me/password"}, {"DELETE", "/api/my/account"}, {"PUT", "/api/my/email"}, {"PUT", "/api/my/profile"}} {
		code, b = e.do(withToken(rw.Token), p.method, p.path, map[string]any{})
		expect(t, p.method+" "+p.path, code, 403, b)
	}
	code, b = e.do(withToken("sb_0000000000000000000000000000000000000000"), "GET", "/api/my", nil)
	expect(t, "unknown token", code, 401, b)

	code, b = e.do(c, "GET", "/api/me/tokens", nil)
	expect(t, "get", code, 200, b)
	var list []store.APIToken
	json.Unmarshal(b, &list)
	for _, tok := range list {
		if tok.LastUsedAt == nil {
			t.Fatalf("last use not recorded: %+v", tok)
		}
	}

	code, b = e.do(other, "DELETE", fmt.Sprintf("/api/me/tokens/%d", ro.ID), nil)
	expect(t, "revoke a foreign token", code, 404, b)
	code, b = e.do(c, "DELETE", fmt.Sprintf("/api/me/tokens/%d", ro.ID), nil)
	expect(t, "revoke", code, 204, b)
	code, b = e.do(withToken(ro.Token), "GET", "/api/my", nil)
	expect(t, "revoked token", code, 401, b)

	code, b = e.do(c, "GET", "/api/me", nil)
	expect(t, "get", code, 200, b)
	var me web.Me
	json.Unmarshal(b, &me)
	secret := "sb_1111111111111111111111111111111111111111"
	sum := sha256.Sum256([]byte(secret))
	past := time.Now().Add(-time.Hour)
	if _, err := e.st.CreateAPIToken(context.Background(), me.User.ID, "old", hex.EncodeToString(sum[:]), []string{"training:read"}, []int64{}, &past); err != nil {
		t.Fatal(err)
	}
	code, b = e.do(withToken(secret), "GET", "/api/my", nil)
	expect(t, "expired token", code, 401, b)
	if !strings.Contains(string(b), "expired") {
		t.Fatalf("expired token message: %s", b)
	}
}

// TestAPITokenClubs checks tokens on club data: only the chosen clubs, only scopes the owner's role has, the role and
// the scopes both limit what the token may do, and the token info.
func TestAPITokenClubs(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	a, b, c := e.tenant("Club A"), e.tenant("Club B"), e.tenant("Club C")
	owner := e.user("owner", a, store.RoleAdmin)
	if err := e.st.SetMembership(ctx, b, owner, e.role(b, store.RoleStaff)); err != nil {
		t.Fatal(err)
	}
	e.event(a, "A event")
	e.event(b, "B event")
	session := e.login("owner")
	create := func(scopes []string, tenants []int64) (int, []byte) {
		return e.do(session, "POST", "/api/me/tokens", web.NewAPIToken{Name: "script", Scopes: scopes, Tenants: tenants})
	}
	token := func(scopes []string, tenants []int64) *http.Client {
		t.Helper()
		code, body := create(scopes, tenants)
		expect(t, fmt.Sprint("token ", scopes, tenants), code, 201, body)
		var tok web.CreatedAPIToken
		json.Unmarshal(body, &tok)
		return withToken(tok.Token)
	}
	for what, call := range map[string][]any{
		"a club the owner is not in":     {[]string{"events.view"}, []int64{c}},
		"a club scope without a club":    {[]string{"events.view"}, []int64{}},
		"a scope the role does not have": {[]string{"club.manage"}, []int64{b}},
	} {
		code, body := create(call[0].([]string), call[1].([]int64))
		expect(t, what, code, 422, body)
	}

	path := func(tenant int64, rest string) string { return fmt.Sprintf("/api/tenants/%d%s", tenant, rest) }
	view := token([]string{"events.view"}, []int64{a})
	code, body := e.do(view, "GET", path(a, "/events"), nil)
	expect(t, "events of a chosen club", code, 200, body)
	code, body = e.do(view, "GET", path(b, "/events"), nil)
	expect(t, "events of a club not chosen", code, 404, body)
	code, body = e.do(view, "GET", path(a, "/shooters"), nil)
	expect(t, "a scope the token lacks although the role has it", code, 403, body)
	if !strings.Contains(string(body), "Not allowed for this API token") {
		t.Errorf("a missing scope names the token: %s", body)
	}

	both := token([]string{"events.view", "shooters.view", "club.manage"}, []int64{a, b})
	code, body = e.do(both, "GET", path(b, "/shooters"), nil)
	expect(t, "helper role in B sees shooters", code, 200, body)
	code, body = e.do(both, "GET", path(b, "/users"), nil)
	expect(t, "a scope the role in B lacks", code, 403, body)
	if !strings.Contains(string(body), "Not allowed for your role") {
		t.Errorf("a permission the role lacks names the role: %s", body)
	}
	code, body = e.do(both, "GET", path(a, "/users"), nil)
	expect(t, "the same scope where the role has it", code, 200, body)
	if err := e.st.SetMembership(ctx, b, owner, e.role(b, store.RoleMember)); err != nil {
		t.Fatal(err)
	}
	code, body = e.do(both, "GET", path(b, "/shooters"), nil)
	expect(t, "a lower role takes effect at once", code, 403, body)
	if err := e.st.RemoveMembership(ctx, b, owner); err != nil {
		t.Fatal(err)
	}
	code, body = e.do(both, "GET", path(b, "/events"), nil)
	expect(t, "no membership any more", code, 404, body)

	code, body = e.do(both, "GET", "/api/token", nil)
	expect(t, "token info", code, 200, body)
	var info web.TokenInfo
	json.Unmarshal(body, &info)
	if info.Owner != "owner" || len(info.Clubs) != 2 || info.Clubs[0].Name != "Club A" || len(info.Scopes) != 3 {
		t.Fatalf("token info: %+v", info)
	}
	code, body = e.do(session, "GET", "/api/token", nil)
	expect(t, "token info without a token", code, 401, body)
	code, body = e.do(both, "POST", path(a, "/roles"), map[string]any{"name": "Wide", "permissions": []string{"ranges.manage"}})
	expect(t, "a token gives no permission it lacks, even an admin's", code, 403, body)
}

// TestAPITokenRateLimits checks the limit per token and the limit of unknown tokens per address.
func TestAPITokenRateLimits(t *testing.T) {
	e := newEnvWith(t, func(s *web.Server) { s.APIRateLimit = 3 })
	tid := e.tenant("Busy")
	e.user("owner", tid, store.RoleAdmin)
	code, body := e.do(e.login("owner"), "POST", "/api/me/tokens", web.NewAPIToken{Name: "busy", Scopes: []string{"events.view"}, Tenants: []int64{tid}})
	expect(t, "token", code, 201, body)
	var tok web.CreatedAPIToken
	json.Unmarshal(body, &tok)
	path := fmt.Sprintf("/api/tenants/%d/events", tid)
	for i := 0; i < 3; i++ {
		code, body = e.do(withToken(tok.Token), "GET", path, nil)
		expect(t, "within the limit", code, 200, body)
	}
	req, _ := http.NewRequest("GET", e.srv.URL+path, nil)
	req.Header.Set("Authorization", "Bearer "+tok.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 429 || resp.Header.Get("Retry-After") == "" {
		t.Fatalf("over the limit: %d %q", resp.StatusCode, resp.Header.Get("Retry-After"))
	}

	bad := withToken("sb_2222222222222222222222222222222222222222")
	for i := 0; i < 20; i++ {
		code, body = e.do(bad, "GET", path, nil)
		expect(t, "unknown token", code, 401, body)
	}
	code, body = e.do(bad, "GET", path, nil)
	expect(t, "too many unknown tokens", code, 429, body)
}
