package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/standblatt/server/internal/store"
)

func TestPairing(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	f := e.registrationFixture("Club")
	e.user("boss", f.tenant, store.RoleAdmin)
	e.user("helper", f.tenant, store.RoleStaff)
	admin, staff, anon := e.login("boss"), e.login("helper"), e.client()
	ev, _ := e.st.GetEvent(ctx, f.tenant, f.event)
	rangeID := ev.Ranges[0].ID
	pairPath := fmt.Sprintf("/api/tenants/%d/ranges/%d/print-agents/pair", f.tenant, rangeID)

	type pairing struct {
		Code, Secret string
		ExpiresAt    time.Time
	}
	start := func() pairing {
		t.Helper()
		code, b := e.do(anon, "POST", "/api/agent/pair", map[string]any{"name": "DESKTOP-RANGE", "os": "windows", "version": "v0.1.0"})
		expect(t, "pair", code, 201, b)
		var p pairing
		json.Unmarshal(b, &p)
		return p
	}
	type status struct{ Status, Token, Club, Range string }
	poll := func(secret string) (int, status) {
		t.Helper()
		code, b := e.do(anon, "POST", "/api/agent/pair/status", map[string]any{"secret": secret})
		var s status
		json.Unmarshal(b, &s)
		return code, s
	}

	p := start()
	if !regexp.MustCompile(`^[A-HJKMNP-Z2-9]{4}-[A-HJKMNP-Z2-9]{4}$`).MatchString(p.Code) || p.Secret == "" || time.Until(p.ExpiresAt) < 9*time.Minute {
		t.Fatalf("pairing: %+v", p)
	}
	if code, s := poll(p.Secret); code != 200 || s.Status != "pending" || s.Token != "" {
		t.Fatalf("pending: %d %+v", code, s)
	}
	code, b := e.do(anon, "POST", "/api/agent/pair", map[string]any{"name": "x", "os": "plan9", "version": "v1"})
	expect(t, "pair from an unknown OS", code, 422, b)

	typed := strings.ToLower(strings.ReplaceAll(p.Code, "-", " "))
	code, b = e.do(anon, "GET", "/api/agent-pairings/"+typed, nil)
	expect(t, "look up without login", code, 401, b)
	code, b = e.do(admin, "GET", "/api/agent-pairings/"+typed, nil)
	expect(t, "look up", code, 200, b)
	var open store.Pairing
	json.Unmarshal(b, &open)
	if open.Name != "DESKTOP-RANGE" || open.OS != "windows" || open.Version != "v0.1.0" {
		t.Fatalf("open pairing: %s", b)
	}
	code, b = e.do(admin, "GET", "/api/agent-pairings/AAAA-AAAA", nil)
	expect(t, "look up an unknown code", code, 404, b)

	code, b = e.do(staff, "POST", pairPath, map[string]any{"code": p.Code, "name": "Laptop"})
	expect(t, "staff confirms", code, 403, b)
	code, b = e.do(admin, "POST", pairPath, map[string]any{"code": "AAAA-AAAA", "name": "Laptop"})
	expect(t, "confirm an unknown code", code, 404, b)
	other := e.registrationFixture("Other club")
	otherEvent, _ := e.st.GetEvent(ctx, other.tenant, other.event)
	code, b = e.do(admin, "POST", fmt.Sprintf("/api/tenants/%d/ranges/%d/print-agents/pair", f.tenant, otherEvent.Ranges[0].ID), map[string]any{"code": p.Code, "name": "Laptop"})
	expect(t, "confirm on a range of another club", code, 404, b)
	code, b = e.do(admin, "POST", pairPath, map[string]any{"code": typed, "name": " Laptop "})
	expect(t, "confirm", code, 201, b)
	var agent store.PrintAgent
	json.Unmarshal(b, &agent)
	if agent.Name != "Laptop" || agent.OS != "windows" || agent.Version != "v0.1.0" || agent.ID == 0 {
		t.Fatalf("agent: %s", b)
	}
	code, b = e.do(admin, "POST", pairPath, map[string]any{"code": p.Code, "name": "Again"})
	expect(t, "confirm twice", code, 404, b)
	code, b = e.do(admin, "GET", "/api/agent-pairings/"+p.Code, nil)
	expect(t, "look up a confirmed code", code, 404, b)

	code, s := poll(p.Secret)
	if code != 200 || s.Status != "connected" || s.Token == "" || s.Club != "Club" || s.Range != ev.Ranges[0].Name {
		t.Fatalf("connected: %d %+v", code, s)
	}
	if code, _ := poll(p.Secret); code != 404 {
		t.Fatalf("second fetch: %d, want 404", code)
	}
	req, _ := http.NewRequest("POST", e.srv.URL+"/api/agent/hello", strings.NewReader(`{"printers":[],"version":"v0.1.0","os":"windows"}`))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("hello with the paired token: %v %v", resp, err)
	}
	resp.Body.Close()
	code, b = e.do(admin, "GET", fmt.Sprintf("/api/tenants/%d/ranges/%d/print-agents", f.tenant, rangeID), nil)
	expect(t, "agents", code, 200, b)
	if !strings.Contains(string(b), `"name":"Laptop"`) || !strings.Contains(string(b), `"lastSeenAt"`) {
		t.Fatalf("agents: %s", b)
	}

	store.PairingTTL = -time.Second
	t.Cleanup(func() { store.PairingTTL = 10 * time.Minute })
	expired := start()
	if code, _ := poll(expired.Secret); code != 410 {
		t.Fatalf("expired code: %d, want 410", code)
	}
	code, b = e.do(admin, "POST", pairPath, map[string]any{"code": expired.Code, "name": "Late"})
	expect(t, "confirm an expired code", code, 404, b)
	if code, _ := poll("unknown"); code != 404 {
		t.Fatalf("unknown secret: %d, want 404", code)
	}
}

func TestPairingRateLimit(t *testing.T) {
	e := newEnv(t)
	anon := e.client()
	for i := range 21 {
		code, b := e.do(anon, "POST", "/api/agent/pair", map[string]any{"name": "PC", "os": "linux", "version": "dev"})
		if want := map[bool]int{true: 201, false: 429}[i < 20]; code != want {
			t.Fatalf("pairing %d: %d %s, want %d", i+1, code, b, want)
		}
	}
}
