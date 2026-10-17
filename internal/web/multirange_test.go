package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func TestMultiRangeEvent(t *testing.T) {
	e := newEnv(t)
	f := e.registrationFixture("Fest")
	e.user("boss", f.tenant, store.RoleAdmin)
	boss := e.login("boss")
	ctx := context.Background()
	r50, _ := e.st.CreateRange(ctx, f.tenant, store.RangeInput{Name: "50m", LaneCount: 3})
	ev, _ := e.st.GetEvent(ctx, f.tenant, f.event)
	r10 := ev.Ranges[0].ID
	evPath := fmt.Sprintf("/api/tenants/%d/events/%d", f.tenant, f.event)
	lane := func(n int) *int { return &n }

	code, b := e.do(boss, "PUT", evPath, map[string]any{"rangeIds": []int64{r10, r50.ID}, "name": ev.Name, "startsOn": "2026-11-21",
		"schedule": domain.Schedule{Days: []domain.ScheduleDay{{Date: "2026-11-22"}, {Date: "2026-11-18"}}}})
	expect(t, "add a range and days", code, 204, b)
	code, b = e.do(boss, "PUT", evPath, map[string]any{"rangeIds": []int64{r10}, "name": ev.Name, "startsOn": "2026-11-21",
		"schedule": domain.Schedule{From: "2026-12-01", To: "2026-11-01"}})
	if code != 422 || !strings.Contains(string(b), "ends before it starts") {
		t.Fatalf("period backwards: %d %s", code, b)
	}
	code, b = e.do(boss, "GET", evPath, nil)
	json.Unmarshal(b, &ev)
	if len(ev.Ranges) != 2 || ev.Ranges[1].Name != "50m" || strings.Join(ev.Days, ",") != "2026-11-18,2026-11-21,2026-11-22" {
		t.Fatalf("event with two ranges: %+v", ev)
	}

	code, b = e.do(boss, "PUT", fmt.Sprintf("/api/tenants/%d/programmes/%d/range", f.tenant, f.junior), map[string]any{"rangeId": r50.ID})
	expect(t, "tie junior to 50m", code, 204, b)
	other := e.tenant("Other")
	foreign, _ := e.st.CreateRange(ctx, other, store.RangeInput{Name: "25m", LaneCount: 2})
	code, b = e.do(boss, "PUT", fmt.Sprintf("/api/tenants/%d/programmes/%d/range", f.tenant, f.junior), map[string]any{"rangeId": foreign.ID})
	expect(t, "range of another tenant", code, 404, b)

	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.youth, Lane: lane(1), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}})
	expect(t, "lane without range on a two-range event", code, 422, b)
	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.youth, RangeID: &r10, Lane: lane(1), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}})
	if code != 422 || !strings.Contains(string(b), "not shot on this range") {
		t.Fatalf("junior on 10m: %d %s", code, b)
	}
	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.youth, RangeID: &r50.ID, Lane: lane(1), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}})
	expect(t, "junior on 50m lane 1", code, 201, b)
	code, b = e.do(boss, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: f.shooter, ClubID: &f.club, RangeID: &r10, Lane: lane(1), Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.main}}})
	expect(t, "main on 10m lane 1 (same lane number, other range)", code, 201, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)
	if reg.RangeID == nil || *reg.RangeID != r10 {
		t.Fatalf("registration range: %+v", reg)
	}

	code, b = e.do(boss, "GET", evPath+"/lanes", nil)
	var lanes []web.RangeLanes
	json.Unmarshal(b, &lanes)
	if len(lanes) != 2 || lanes[0].RangeName != "10m" || lanes[0].Lanes[0].State != "occupied" || len(lanes[1].Lanes) != 3 || lanes[1].Lanes[0].ShooterName != "Young Yves" {
		t.Fatalf("lanes per range: %+v", lanes)
	}

	code, b = e.do(boss, "PUT", evPath, map[string]any{"rangeIds": []int64{r10}, "name": ev.Name, "startsOn": "2026-11-21"})
	expect(t, "drop a range in use", code, 409, b)

	code, b = e.do(boss, "POST", evPath+"/copy", map[string]any{"name": "Fest 2027", "startsOn": "2027-11-20"})
	expect(t, "copy", code, 201, b)
	var copied store.Event
	json.Unmarshal(b, &copied)
	if len(copied.Ranges) != 2 || strings.Join(copied.Days, ",") != "2027-11-17,2027-11-20,2027-11-21" {
		t.Fatalf("copied event: %+v", copied)
	}
	progs, _ := e.st.ListProgrammes(ctx, f.tenant, copied.ID)
	for _, p := range progs {
		if p.Name == "Junior" && (p.RangeID == nil || *p.RangeID != r50.ID) {
			t.Fatalf("copied programme keeps its range: %+v", p)
		}
	}

	code, b = e.do(boss, "GET", evPath+"/statistics", nil)
	if code != 200 || !strings.Contains(string(b), `"registrations":2`) {
		t.Fatalf("stats per day: %d %s", code, b)
	}
}
