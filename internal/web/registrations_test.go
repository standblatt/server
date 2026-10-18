package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

type regFixture struct {
	tenant, event, main, section, junior, shooter, youth, club int64
}

func (e *env) registrationFixture(name string) regFixture {
	e.t.Helper()
	ctx := context.Background()
	var f regFixture
	f.tenant = e.tenant(name)
	r, _ := e.st.CreateRange(ctx, f.tenant, store.RangeInput{Name: "10m", LaneCount: 5, DisabledLanes: []int32{5}})
	f.event, _ = e.st.CreateEvent(ctx, f.tenant, store.EventInput{RangeIDs: []int64{r.ID}, Name: "Indoor", StartsOn: "2026-11-20", Timezone: "Europe/Zurich"})
	def := func(n string, series int, cats ...domain.AgeCategory) domain.ProgrammeDefinition {
		d := domain.ProgrammeDefinition{Name: n, Awards: []domain.AwardTier{}}
		for i := 0; i < series; i++ {
			d.Series = append(d.Series, domain.Series{Shots: 10, MaxPoints: 100})
		}
		for _, c := range cats {
			d.Prices = append(d.Prices, domain.Price{AgeCategory: c, Price: 10})
		}
		return d
	}
	f.main, _ = e.st.CreateProgramme(ctx, f.tenant, f.event, def("Main", 4, domain.Elite, domain.U21))
	f.section, _ = e.st.CreateProgramme(ctx, f.tenant, f.event, def("Section", 2, domain.Elite))
	f.junior, _ = e.st.CreateProgramme(ctx, f.tenant, f.event, def("Junior", 2, domain.U21))
	e.st.SetCombination(ctx, f.tenant, f.section, &f.main)
	club, _ := e.st.CreateClub(ctx, f.tenant, store.ClubInput{Name: "Home"})
	f.club = club.ID
	f.shooter, _ = e.st.CreateShooter(ctx, f.tenant, store.ShooterInput{LicenseNumber: "100001", Firstname: "Anna", Lastname: "Adult", Birthday: "1990-01-01"}, []int64{f.club})
	f.youth, _ = e.st.CreateShooter(ctx, f.tenant, store.ShooterInput{Firstname: "Yves", Lastname: "Young", Birthday: "2007-03-03"}, nil)
	return f
}

func TestRegistrations(t *testing.T) {
	e := newEnv(t)
	f := e.registrationFixture("Club")
	e.user("helper", f.tenant, store.RoleStaff)
	other := e.tenant("Other")
	e.user("outsider", other, store.RoleAdmin)
	staff, outsider := e.login("helper"), e.login("outsider")
	evPath := fmt.Sprintf("/api/tenants/%d/events/%d", f.tenant, f.event)
	lane := func(n int) *int { return &n }

	body := web.RegistrationInput{ShooterID: f.shooter, ClubID: &f.club, Lane: lane(2), Programmes: []domain.ProgrammeChoice{
		{ProgrammeID: f.main}, {ProgrammeID: f.section, CombinedSeries: []int{2, 4}},
	}}
	code, b := e.do(staff, "POST", evPath+"/registrations", body)
	expect(t, "register", code, 201, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)
	if reg.Status != "shooting" || reg.AgeCategory != domain.Elite || reg.StartedAt == nil || len(reg.Entries) != 2 || reg.ClubName != "Home" {
		t.Fatalf("registration: %+v", reg)
	}
	regPath := fmt.Sprintf("/api/tenants/%d/registrations/%d", f.tenant, reg.ID)

	for name, c := range map[string]struct {
		body web.RegistrationInput
		code int
		msg  string
	}{
		"lane taken":            {web.RegistrationInput{ShooterID: f.youth, Lane: lane(2), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}}, 409, "Lane 2 is occupied"},
		"lane out of service":   {web.RegistrationInput{ShooterID: f.youth, Lane: lane(5), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}}, 422, "out of service"},
		"lane missing":          {web.RegistrationInput{ShooterID: f.youth, Lane: lane(6), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}}, 422, "does not exist"},
		"category not allowed":  {web.RegistrationInput{ShooterID: f.youth, Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.section}}}, 422, "not open to the age category u21"},
		"programme twice":       {web.RegistrationInput{ShooterID: f.shooter, Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.main}}}, 409, "already registered"},
		"club of someone else":  {web.RegistrationInput{ShooterID: f.youth, ClubID: &f.club, Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}}, 422, "not a member"},
		"combination too short": {web.RegistrationInput{ShooterID: f.youth, Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.main}, {ProgrammeID: f.section, CombinedSeries: []int{1}}}}, 422, "needs exactly 2 series"},
	} {
		code, b := e.do(staff, "POST", evPath+"/registrations", c.body)
		if code != c.code || !strings.Contains(string(b), c.msg) {
			t.Errorf("%s: %d %s", name, code, b)
		}
	}

	code, b = e.do(staff, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.youth, Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}})
	expect(t, "register youth without lane", code, 201, b)
	var youth store.Registration
	json.Unmarshal(b, &youth)
	if youth.Status != "registered" || youth.AgeCategory != domain.U21 || youth.Lane != nil {
		t.Fatalf("youth registration: %+v", youth)
	}

	code, b = e.do(staff, "PUT", regPath, web.RegistrationUpdate{Lane: lane(3), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.main}}})
	expect(t, "update registration", code, 200, b)
	var updated store.Registration
	json.Unmarshal(b, &updated)
	if *updated.Lane != 3 || len(updated.Entries) != 1 || updated.Entries[0].ID != reg.Entries[0].ID {
		t.Fatalf("updated registration keeps the main entry: %+v", updated)
	}

	code, b = e.do(staff, "GET", evPath+"/lanes", nil)
	expect(t, "lanes", code, 200, b)
	var byRange []web.RangeLanes
	json.Unmarshal(b, &byRange)
	lanes := byRange[0].Lanes
	if len(lanes) != 5 || lanes[2].State != "occupied" || lanes[2].ShooterName != "Adult Anna" || lanes[1].State != "free" || lanes[4].State != "disabled" {
		t.Fatalf("lanes: %+v", lanes)
	}

	code, b = e.do(staff, "PUT", fmt.Sprintf("/api/tenants/%d/programmes/%d", f.tenant, f.main), map[string]any{
		"name": "Main", "series": []map[string]any{{"shots": 5, "maxPoints": 50}}, "prices": []any{}, "awards": []any{},
	})
	expect(t, "staff edits programme", code, 403, b)
	e.user("boss", f.tenant, store.RoleAdmin)
	code, b = e.do(e.login("boss"), "PUT", fmt.Sprintf("/api/tenants/%d/programmes/%d", f.tenant, f.main), map[string]any{
		"name": "Main", "series": []map[string]any{{"shots": 5, "maxPoints": 50}}, "prices": []any{}, "awards": []any{},
	})
	expect(t, "edit programme in use", code, 409, b)
	boss := e.login("boss")
	kept := func(prices ...string) map[string]any {
		var ps []map[string]any
		for _, c := range prices {
			ps = append(ps, map[string]any{"ageCategory": c, "price": 12})
		}
		return map[string]any{
			"name": "Main cup", "series": []map[string]any{{"shots": 10, "maxPoints": 100, "name": "A"}, {"shots": 10, "maxPoints": 100}, {"shots": 10, "maxPoints": 100}, {"shots": 10, "maxPoints": 100}},
			"prices": ps, "awards": []map[string]any{{"name": "Wreath", "kind": "card", "value": 8, "limits": []map[string]any{{"ageCategory": "elite", "points": 360}}}},
		}
	}
	progPath := fmt.Sprintf("/api/tenants/%d/programmes/%d", f.tenant, f.main)
	code, b = e.do(boss, "PUT", progPath, kept("elite", "u21", "senior"))
	expect(t, "edit programme in use without touching the series", code, 204, b)
	code, b = e.do(boss, "GET", progPath, nil)
	expect(t, "get", code, 200, b)
	var edited store.Programme
	json.Unmarshal(b, &edited)
	if !edited.InUse || edited.Name != "Main cup" || len(edited.Prices) != 3 || edited.Prices[0].Price != 12 || len(edited.Awards) != 1 || edited.Series[0].Name != "A" {
		t.Errorf("edited programme in use: %+v", edited)
	}
	code, b = e.do(boss, "PUT", progPath, kept("elite"))
	expect(t, "remove an age category in use", code, 409, b)
	code, b = e.do(boss, "DELETE", progPath, nil)
	expect(t, "delete programme in use", code, 409, b)

	code, b = e.do(staff, "POST", regPath+"/finish", nil)
	expect(t, "finish", code, 204, b)
	code, b = e.do(staff, "GET", evPath+"/stats", nil)
	expect(t, "stats", code, 200, b)
	var st store.EventStats
	json.Unmarshal(b, &st)
	if st != (store.EventStats{Registered: 1, Completed: 1}) {
		t.Fatalf("stats: %+v", st)
	}
	code, b = e.do(staff, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.youth + 0, Lane: lane(3), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.main}}})
	if code == 409 && strings.Contains(string(b), "Lane 3") {
		t.Fatalf("lane of a completed registration is still taken: %s", b)
	}

	code, b = e.do(outsider, "GET", regPath, nil)
	expect(t, "outsider reads registration", code, 404, b)
	code, b = e.do(outsider, "POST", fmt.Sprintf("/api/tenants/%d/registrations/%d/finish", other, reg.ID), nil)
	expect(t, "outsider finishes via own tenant", code, 404, b)
	code, b = e.do(staff, "DELETE", fmt.Sprintf("/api/tenants/%d/shooters/%d", f.tenant, f.shooter), nil)
	expect(t, "staff deletes shooter", code, 403, b)
	code, b = e.do(e.login("boss"), "DELETE", fmt.Sprintf("/api/tenants/%d/shooters/%d", f.tenant, f.shooter), nil)
	expect(t, "delete shooter with registrations", code, 409, b)

	code, b = e.do(staff, "DELETE", regPath, nil)
	expect(t, "delete registration", code, 204, b)
}

func TestLiveUpdates(t *testing.T) {
	e := newEnv(t)
	f := e.registrationFixture("Club")
	e.user("helper", f.tenant, store.RoleStaff)
	staff := e.login("helper")
	other := e.registrationFixture("Other")

	dial := func(tenant, event int64) (*websocket.Conn, int) {
		u, _ := http.NewRequest("GET", e.srv.URL, nil)
		header := http.Header{"Cookie": {}}
		for _, c := range staff.Jar.Cookies(u.URL) {
			header.Add("Cookie", c.Name+"="+c.Value)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		url := strings.Replace(e.srv.URL, "http", "ws", 1) + fmt.Sprintf("/api/tenants/%d/events/%d/live", tenant, event)
		conn, resp, err := websocket.Dial(ctx, url, &websocket.DialOptions{HTTPHeader: header})
		if err != nil {
			if resp != nil {
				return nil, resp.StatusCode
			}
			t.Fatal(err)
		}
		return conn, 101
	}
	if _, code := dial(other.tenant, other.event); code != 404 {
		t.Fatalf("live feed of a foreign tenant: %d", code)
	}
	conn, _ := dial(f.tenant, f.event)
	defer conn.CloseNow()

	read := func() string {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, b, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("no message: %v", err)
		}
		var m map[string]string
		json.Unmarshal(b, &m)
		return m["topic"]
	}
	time.Sleep(200 * time.Millisecond)
	code, b := e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/registrations", f.tenant, f.event),
		web.RegistrationInput{ShooterID: f.youth, Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}})
	expect(t, "register", code, 201, b)
	topics := map[string]bool{}
	for len(topics) < 2 {
		topic := read()
		if topic == "resync" {
			continue
		}
		topics[topic] = true
	}
	if !topics["registrations"] || !topics["lanes"] {
		t.Fatalf("topics: %v", topics)
	}
}
