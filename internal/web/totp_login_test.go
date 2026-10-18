package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/standblatt/server/internal/web"
)

func TestTOTPLogin(t *testing.T) {
	e := newEnvWith(t, func(s *web.Server) { s.TOTP = true })
	tid := e.tenant("Club")
	e.user("anna", tid, "admin")
	c := e.login("anna")

	code, b := e.do(c, "POST", "/api/me/totp/setup", nil)
	expect(t, "setup", code, 200, b)
	var setup struct{ Secret, URL, QR string }
	json.Unmarshal(b, &setup)
	if len(setup.Secret) != 32 || setup.QR[:22] != "data:image/png;base64," || setup.URL[:15] != "otpauth://totp/" {
		t.Fatalf("setup: %s", b)
	}
	code, b = e.do(c, "PUT", "/api/me/totp", map[string]string{"secret": setup.Secret, "code": "000000"})
	expect(t, "enable with wrong code", code, 401, b)
	now := func() string { v, _ := web.TOTPCodeForTest(setup.Secret, time.Now()); return v }
	enabled := now()
	code, b = e.do(c, "PUT", "/api/me/totp", map[string]string{"secret": setup.Secret, "code": enabled})
	expect(t, "enable", code, 204, b)
	code, b = e.do(c, "GET", "/api/me", nil)
	expect(t, "me", code, 200, b)
	var me web.Me
	json.Unmarshal(b, &me)
	if !me.TOTP {
		t.Fatalf("me.totp false: %s", b)
	}

	login := func(body map[string]string) (int, []byte) { return e.do(e.client(), "POST", "/api/auth/login", body) }
	code, b = login(map[string]string{"name": "anna", "password": "password-anna"})
	expect(t, "login without code", code, http.StatusPreconditionRequired, b)
	code, b = login(map[string]string{"name": "anna", "password": "wrong", "code": now()})
	expect(t, "login wrong password", code, 401, b)
	code, b = login(map[string]string{"name": "anna", "password": "password-anna", "code": "123456"})
	expect(t, "login wrong code", code, 401, b)
	code, b = login(map[string]string{"name": "anna", "password": "password-anna", "code": enabled})
	expect(t, "login replayed code", code, 401, b) // the enable step burned this time step
	next, _ := web.TOTPCodeForTest(setup.Secret, time.Now().Add(30*time.Second))
	code, b = login(map[string]string{"name": "anna", "password": "password-anna", "code": next})
	expect(t, "login with the next code", code, 200, b) // one step of clock drift is accepted and newer than the burned step

	code, b = e.do(c, "POST", "/api/me/totp/disable", map[string]string{"password": "nope"})
	expect(t, "disable wrong password", code, 401, b)
	code, b = e.do(c, "POST", "/api/me/totp/disable", map[string]string{"password": "password-anna"})
	expect(t, "disable", code, 204, b)
	code, b = login(map[string]string{"name": "anna", "password": "password-anna"})
	expect(t, "login after disable", code, 200, b)
}

func TestRememberMe(t *testing.T) {
	e := newEnv(t)
	e.user("bea", 0, "")
	for _, c := range []struct {
		remember bool
		maxAge   string
	}{{false, "Max-Age=2592000"}, {true, "Max-Age=7776000"}} {
		body, _ := json.Marshal(map[string]any{"name": "bea", "password": "password-bea", "remember": c.remember})
		resp, err := http.Post(e.srv.URL+"/api/auth/login", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != 200 || !strings.Contains(resp.Header.Get("Set-Cookie"), c.maxAge) {
			t.Fatalf("remember=%v: %d %q", c.remember, resp.StatusCode, resp.Header.Get("Set-Cookie"))
		}
	}
}
