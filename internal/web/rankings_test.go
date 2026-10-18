package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func (e *env) kiosk(method, path, token string) (int, []byte) {
	e.t.Helper()
	req, _ := http.NewRequest(method, e.srv.URL+path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b
}

func TestRankingsAndKiosk(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	admin, staff := e.login("boss"), e.login("helper")
	ev := e.event(tid, "Indoor")
	main, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Main", ShowOnRankings: true,
		Series: []domain.Series{{Shots: 10, MaxPoints: 100}, {Shots: 10, MaxPoints: 100}},
		Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 10}, {AgeCategory: domain.Senior, Price: 10}, {AgeCategory: domain.U21, Price: 5}},
		Awards: []domain.AwardTier{},
	})
	e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Hidden", Series: []domain.Series{{Shots: 1, MaxPoints: 10}}, Prices: []domain.Price{}, Awards: []domain.AwardTier{},
	})

	shoot := func(first, last, birthday string, points ...float64) {
		id, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: first, Lastname: last, Birthday: birthday}, nil)
		code, b := e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/registrations", tid, ev),
			web.RegistrationInput{ShooterID: id, Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}}})
		expect(t, "register "+first, code, 201, b)
		var reg store.Registration
		json.Unmarshal(b, &reg)
		series := []map[string]any{}
		for _, p := range points {
			series = append(series, map[string]any{"points": p})
		}
		code, b = e.do(staff, "PUT", fmt.Sprintf("/api/tenants/%d/entries/%d/results", tid, reg.Entries[0].ID), map[string]any{"series": series})
		expect(t, "results "+first, code, 200, b)
	}
	shoot("Anna", "Alpha", "1990-01-01", 90, 95)
	shoot("Bert", "Beta", "1985-01-01", 95, 90)
	shoot("Cleo", "Gamma", "2006-01-01", 100, 99)
	shoot("Dora", "Delta", "1970-01-01", 95, 90)

	code, b := e.do(staff, "GET", fmt.Sprintf("/api/tenants/%d/events/%d/rankings", tid, ev), nil)
	expect(t, "rankings", code, 200, b)
	var rankings []web.ProgrammeRanking
	json.Unmarshal(b, &rankings)
	byTitle := map[string]web.ProgrammeRanking{}
	for _, r := range rankings {
		byTitle[r.Title] = r
	}
	if len(rankings) != 2 || !byTitle["Main"].Public || byTitle["Hidden"].Public {
		t.Fatalf("staff see every programme: %+v", rankings)
	}
	var order []string
	for _, l := range byTitle["Main"].Lines {
		order = append(order, fmt.Sprintf("%d %s %.0f", l.Rank, l.Lastname, l.Total))
	}
	if got := strings.Join(order, ", "); got != "1 Gamma 199, 2 Alpha 185, 3 Delta 185, 4 Beta 185" {
		t.Fatalf("ranking order: %s", got)
	}
	code, b = e.do(staff, "GET", fmt.Sprintf("/api/tenants/%d/programmes/%d/ranking", tid, byTitle["Main"].ProgrammeID), nil)
	if code != 200 || !bytes.HasPrefix(b, []byte("%PDF")) {
		t.Fatalf("ranking PDF: %d %.100s", code, b)
	}

	tokensPath := fmt.Sprintf("/api/tenants/%d/events/%d/kiosk-tokens", tid, ev)
	code, b = e.do(staff, "POST", tokensPath, map[string]any{"name": "Hall"})
	expect(t, "staff creates kiosk token", code, 403, b)
	code, b = e.do(admin, "POST", tokensPath, map[string]any{"name": "Hall"})
	expect(t, "create kiosk token", code, 201, b)
	var created struct {
		Token  store.KioskToken `json:"token"`
		Secret string           `json:"secret"`
	}
	json.Unmarshal(b, &created)

	code, _ = e.kiosk("GET", "/api/kiosk/rankings", "")
	if code != 401 {
		t.Fatalf("kiosk without token: %d", code)
	}
	code, _ = e.kiosk("GET", "/api/kiosk/rankings", "wrong")
	if code != 401 {
		t.Fatalf("kiosk with wrong token: %d", code)
	}
	code, b = e.kiosk("GET", "/api/kiosk/rankings", created.Secret)
	if code != 200 || !strings.Contains(string(b), "Gamma") || strings.Contains(string(b), "Hidden") {
		t.Fatalf("kiosk rankings without hidden programmes: %d %s", code, b)
	}
	code, b = e.kiosk("GET", "/api/kiosk/lanes", created.Secret)
	if code != 200 || string(b) != "[]\n" && string(b) != "[]" {
		t.Fatalf("lanes hidden unless enabled: %d %s", code, b)
	}
	code, b = e.kiosk("GET", "/api/kiosk/event", created.Secret)
	if code != 200 || !strings.Contains(string(b), `"name":"Indoor"`) {
		t.Fatalf("kiosk event: %d %s", code, b)
	}
	code, _ = e.kiosk("GET", fmt.Sprintf("/api/tenants/%d/events/%d/rankings", tid, ev), created.Secret)
	if code != 401 {
		t.Fatalf("kiosk token opens the staff API: %d", code)
	}

	ctxWS, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	wsURL := strings.Replace(e.srv.URL, "http", "ws", 1) + "/api/kiosk/live?token=" + created.Secret
	conn, _, err := websocket.Dial(ctxWS, wsURL, nil)
	if err != nil {
		t.Fatalf("kiosk live: %v", err)
	}
	conn.CloseNow()
	if _, resp, err := websocket.Dial(ctxWS, strings.Replace(wsURL, created.Secret, "wrong", 1), nil); err == nil || resp.StatusCode != 401 {
		t.Fatalf("kiosk live with wrong token: %v", err)
	}

	code, b = e.do(admin, "GET", tokensPath, nil)
	expect(t, "get", code, 200, b)
	var tokens []store.KioskToken
	json.Unmarshal(b, &tokens)
	if len(tokens) != 1 || tokens[0].LastUsedAt == nil {
		t.Fatalf("tokens: %+v", tokens)
	}
	code, b = e.do(admin, "DELETE", fmt.Sprintf("/api/tenants/%d/kiosk-tokens/%d", tid, created.Token.ID), nil)
	expect(t, "revoke token", code, 204, b)
	if code, _ = e.kiosk("GET", "/api/kiosk/rankings", created.Secret); code != 401 {
		t.Fatalf("revoked token still works: %d", code)
	}
}
