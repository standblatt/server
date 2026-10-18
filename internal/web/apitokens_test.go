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
		code, b := e.do(cl, "POST", "/api/my/tokens", web.NewAPIToken{Name: "script", Scopes: scopes, ExpiresInDays: 30})
		expect(t, "create token", code, 201, b)
		var tok web.CreatedAPIToken
		json.Unmarshal(b, &tok)
		return tok
	}
	ro, rw := create(c, "training:read"), create(c, "training:read", "training:write")
	if !strings.HasPrefix(ro.Token, "sb_") || len(ro.Token) != 43 || ro.ExpiresAt == nil {
		t.Fatalf("token: %+v", ro)
	}
	code, b := e.do(c, "POST", "/api/my/tokens", web.NewAPIToken{Name: "x", Scopes: []string{"club:write"}})
	expect(t, "unknown scope", code, 422, b)

	code, b = e.do(c, "GET", "/api/my/tokens", nil)
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

	for _, p := range []struct{ method, path string }{{"GET", "/api/my/tokens"}, {"DELETE", "/api/my/account"}, {"PUT", "/api/my/email"}, {"PUT", "/api/my/profile"}} {
		code, b = e.do(withToken(rw.Token), p.method, p.path, map[string]any{})
		expect(t, p.method+" "+p.path, code, 403, b)
	}
	code, b = e.do(withToken(rw.Token), "GET", "/api/me", nil)
	expect(t, "token outside /api/my", code, 401, b)
	code, b = e.do(withToken("sb_0000000000000000000000000000000000000000"), "GET", "/api/my", nil)
	expect(t, "unknown token", code, 401, b)

	code, b = e.do(c, "GET", "/api/my/tokens", nil)
	expect(t, "get", code, 200, b)
	var list []store.APIToken
	json.Unmarshal(b, &list)
	for _, tok := range list {
		if tok.LastUsedAt == nil {
			t.Fatalf("last use not recorded: %+v", tok)
		}
	}

	code, b = e.do(other, "DELETE", fmt.Sprintf("/api/my/tokens/%d", ro.ID), nil)
	expect(t, "revoke a foreign token", code, 404, b)
	code, b = e.do(c, "DELETE", fmt.Sprintf("/api/my/tokens/%d", ro.ID), nil)
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
	if _, err := e.st.CreateAPIToken(context.Background(), me.User.ID, "old", hex.EncodeToString(sum[:]), []string{"training:read"}, &past); err != nil {
		t.Fatal(err)
	}
	code, b = e.do(withToken(secret), "GET", "/api/my", nil)
	expect(t, "expired token", code, 401, b)
	if !strings.Contains(string(b), "expired") {
		t.Fatalf("expired token message: %s", b)
	}
}
