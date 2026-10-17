package web_test

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/web"
)

func TestMailFollowsRequestLanguage(t *testing.T) {
	e, box := personalEnv(t, nil)
	id := e.user("felix", 0, "")
	if err := e.st.SetEmail(context.Background(), id, "felix@example.org"); err != nil {
		t.Fatal(err)
	}
	forgot := func(lang string) string {
		req, _ := http.NewRequest("POST", e.srv.URL+"/api/auth/forgot", bytes.NewReader([]byte(`{"email":"felix@example.org"}`)))
		req.Header.Set("Content-Type", "application/json")
		if lang != "" {
			req.Header.Set("Accept-Language", lang)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != 204 {
			t.Fatalf("forgot: %d", resp.StatusCode)
		}
		return (<-box).Subject
	}
	if got := forgot("de"); got != "Neues Passwort für Standblatt" {
		t.Fatalf("German request: %q", got)
	}
	if got := forgot(""); got != "Reset your Standblatt password" {
		t.Fatalf("no language: %q", got)
	}
	if err := e.st.SetUserLocale(context.Background(), id, "fr"); err != nil {
		t.Fatal(err)
	}
	if got := forgot("de"); got != "Nouveau mot de passe pour Standblatt" {
		t.Fatalf("user locale wins: %q", got)
	}
}

func TestVersionEnvironment(t *testing.T) {
	e := newEnvWith(t, func(s *web.Server) { s.Environment = "staging" })
	code, b := e.do(e.client(), "GET", "/api/version", nil)
	expect(t, "version", code, 200, b)
	if !strings.Contains(string(b), `"environment":"staging"`) {
		t.Fatalf("environment missing: %s", b)
	}
}

func TestChangePasswordAndFreshProfile(t *testing.T) {
	e := newEnv(t)
	tid := e.tenant("Club")
	e.user("carla", tid, "staff")
	c := e.login("carla")
	other := e.login("carla")

	code, b := e.do(c, "GET", "/api/my", nil)
	expect(t, "fresh overview", code, 200, b)
	if !strings.Contains(string(b), `"profile":{"firstname":"","lastname":""`) {
		t.Fatalf("no empty profile for a club user: %s", b)
	}
	code, b = e.do(c, "PUT", "/api/me/password", map[string]string{"currentPassword": "wrong", "newPassword": "a-new-password-1"})
	expect(t, "wrong current", code, 403, b)
	code, b = e.do(c, "PUT", "/api/me/password", map[string]string{"currentPassword": "password-carla", "newPassword": "short"})
	expect(t, "weak new", code, 422, b)
	code, b = e.do(c, "PUT", "/api/me/password", map[string]string{"currentPassword": "password-carla", "newPassword": "a-new-password-1"})
	expect(t, "change", code, 204, b)
	code, b = e.do(c, "GET", "/api/me", nil)
	expect(t, "own session survives", code, 200, b)
	code, b = e.do(other, "GET", "/api/me", nil)
	expect(t, "other session ended", code, 401, b)
	code, b = e.do(e.client(), "POST", "/api/auth/login", map[string]string{"name": "carla", "password": "password-carla"})
	expect(t, "old password", code, 401, b)
	code, b = e.do(e.client(), "POST", "/api/auth/login", map[string]string{"name": "carla", "password": "a-new-password-1"})
	expect(t, "new password", code, 200, b)
	if !strings.Contains(string(b), `"email":""`) {
		t.Fatalf("me without email field: %s", b)
	}
}
