package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func TestReservations(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	f := e.registrationFixture("Club")
	e.user("helper", f.tenant, store.RoleStaff)
	other := e.tenant("Other")
	e.user("outsider", other, store.RoleAdmin)
	staff, outsider := e.login("helper"), e.login("outsider")
	ev, _ := e.st.GetEvent(ctx, f.tenant, f.event)
	resPath := fmt.Sprintf("/api/tenants/%d/ranges/%d/reservations", f.tenant, ev.Ranges[0].ID)
	now := time.Now().Truncate(time.Minute)
	in := func(lanes, desc string, start, end time.Duration, lead int) web.ReservationInput {
		return web.ReservationInput{Lanes: lanes, StartsAt: now.Add(start), EndsAt: now.Add(end), LeadMinutes: lead, Description: desc}
	}

	code, b := e.do(staff, "POST", resPath, in("2-3, 5", "Training", time.Hour, 3*time.Hour, 90))
	expect(t, "create reservation", code, 201, b)
	var training web.ReservationView
	json.Unmarshal(b, &training)
	if training.LanesText != "2, 3, 5" {
		t.Fatalf("reservation: %+v", training)
	}

	code, b = e.do(staff, "GET", fmt.Sprintf("/api/tenants/%d/events/%d/lanes", f.tenant, f.event), nil)
	var byRange []web.RangeLanes
	json.Unmarshal(b, &byRange)
	lanes := byRange[0].Lanes
	if lanes[1].State != "reserved" || lanes[1].Reservation != "Training" || lanes[0].State != "free" || lanes[4].State != "disabled" {
		t.Fatalf("lanes during lead time: %+v", lanes)
	}

	for name, c := range map[string]struct {
		body web.ReservationInput
		code int
		msg  string
	}{
		"overlap":        {in("3,4", "Club evening", 2*time.Hour, 4*time.Hour, 0), 409, "Training\\\" (lanes 2, 3, 5)"},
		"invalid lanes":  {in("0", "x", time.Hour, 2*time.Hour, 0), 422, "Invalid lane"},
		"lane missing":   {in("6", "x", time.Hour, 2*time.Hour, 0), 422, "lanes 1–5"},
		"end before":     {in("1", "x", 2*time.Hour, time.Hour, 0), 422, "after the start"},
		"ended already":  {in("1", "x", -2*time.Hour, -time.Hour, 0), 422, "in the past"},
		"back to back":   {in("2", "Next", 3*time.Hour, 4*time.Hour, 0), 201, ""},
		"other lanes ok": {in("1, 4", "Parallel", time.Hour, 3*time.Hour, 0), 201, ""},
	} {
		code, b := e.do(staff, "POST", resPath, c.body)
		if code != c.code || !strings.Contains(string(b), c.msg) {
			t.Errorf("%s: %d %s", name, code, b)
		}
	}

	var list []web.ReservationView
	code, b = e.do(staff, "GET", resPath, nil)
	json.Unmarshal(b, &list)
	var parallel web.ReservationView
	for _, r := range list {
		if r.Description == "Parallel" {
			parallel = r
		}
	}
	code, b = e.do(staff, "PUT", fmt.Sprintf("/api/tenants/%d/reservations/%d", f.tenant, parallel.ID), in("1, 4, 5", "Parallel", time.Hour, 3*time.Hour, 0))
	expect(t, "update into an overlap", code, 409, b)
	code, b = e.do(staff, "PUT", fmt.Sprintf("/api/tenants/%d/reservations/%d", f.tenant, parallel.ID), in("1, 4", "Parallel", 5*time.Hour, 6*time.Hour, 0))
	expect(t, "move in time", code, 200, b)

	code, b = e.do(outsider, "GET", resPath, nil)
	expect(t, "outsider lists reservations", code, 404, b)
	code, b = e.do(outsider, "DELETE", fmt.Sprintf("/api/tenants/%d/reservations/%d", other, training.ID), nil)
	expect(t, "outsider deletes via own tenant", code, 404, b)
	code, b = e.do(staff, "DELETE", fmt.Sprintf("/api/tenants/%d/reservations/%d", f.tenant, training.ID), nil)
	expect(t, "delete", code, 204, b)
}

func TestReservationTick(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	f := e.registrationFixture("Club")
	ev, _ := e.st.GetEvent(ctx, f.tenant, f.event)
	now := time.Now()

	if _, err := e.st.ReservationTick(ctx, now.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	_, err := e.st.SaveReservation(ctx, f.tenant, store.Reservation{
		RangeID: ev.Ranges[0].ID, Lanes: []int{1}, StartsAt: now.Add(10 * time.Minute), EndsAt: now.Add(time.Hour), LeadMinutes: 15, Description: "Soon",
	})
	if err != nil {
		t.Fatal(err)
	}
	ranges, err := e.st.ReservationTick(ctx, now)
	if err != nil || len(ranges) != 1 || ranges[0] != ev.Ranges[0].ID {
		t.Fatalf("first tick should announce the range: %v %v", ranges, err)
	}
	if ranges, _ = e.st.ReservationTick(ctx, now.Add(time.Second)); len(ranges) != 0 {
		t.Fatalf("second tick should be quiet: %v", ranges)
	}
	if ranges, _ = e.st.ReservationTick(ctx, now.Add(time.Hour+time.Second)); len(ranges) != 1 {
		t.Fatalf("tick after the end should announce the range again: %v", ranges)
	}
}
