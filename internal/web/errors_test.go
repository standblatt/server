package web_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestErrorsFollowAcceptLanguage(t *testing.T) {
	e := newEnv(t)
	tid := e.tenant("Club")
	e.user("anna", tid, "admin")
	cases := []struct{ lang, path, body, want string }{
		{"", "/api/auth/login", `{"name":"anna","password":"wrong"}`, "Wrong user name or password."},
		{"de-CH,de;q=0.9", "/api/auth/login", `{"name":"anna","password":"wrong"}`, "Benutzername oder Passwort falsch."},
		{"fr", "/api/auth/login", `{"name":"anna","password":"wrong"}`, "Nom d'utilisateur ou mot de passe incorrect."},
	}
	c := e.client()
	for _, tc := range cases {
		req, _ := http.NewRequest("POST", e.srv.URL+tc.path, strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		if tc.lang != "" {
			req.Header.Set("Accept-Language", tc.lang)
		}
		resp, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var got struct{ Detail string }
		json.Unmarshal(b, &got)
		if resp.StatusCode != 401 || got.Detail != tc.want {
			t.Errorf("%q: got %d %q, want %q", tc.lang, resp.StatusCode, got.Detail, tc.want)
		}
	}
	anna := e.login("anna")
	get := func(method, path, lang, body string) (int, string) {
		req, _ := http.NewRequest(method, e.srv.URL+path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Language", lang)
		resp, err := anna.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return resp.StatusCode, string(b)
	}
	if code, b := get("GET", "/api/tenants/999/events", "it", ""); code != 404 || !strings.Contains(b, "Società non trovata") {
		t.Errorf("404 in Italian: got %d %s", code, b)
	}
	ev := e.event(tid, "Cup")
	path := fmt.Sprintf("/api/tenants/%d/events/%d/programmes", tid, ev)
	if code, b := get("POST", path, "de", `{"name":"Cup","series":[],"prices":[{"ageCategory":"elite","price":5},{"ageCategory":"elite","price":6}],"awards":[]}`); code != 422 || !strings.Contains(b, "Mindestens eine Passe ist erforderlich; Die Alterskategorie elite hat mehr als einen Preis") {
		t.Errorf("domain validation in German: got %d %s", code, b)
	}
}
