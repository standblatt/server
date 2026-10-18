package web_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/store"
)

// TestModules checks the gate, the default of all modules and that modules with data cannot be switched off.
func TestModules(t *testing.T) {
	e := newEnv(t)
	f := e.registrationFixture("Modular")
	e.user("boss", f.tenant, store.RoleAdmin)
	boss := e.login("boss")
	base := fmt.Sprintf("/api/tenants/%d", f.tenant)
	evPath := fmt.Sprintf("%s/events/%d", base, f.event)
	get := func() store.Event {
		var ev store.Event
		_, b := e.do(boss, "GET", evPath, nil)
		json.Unmarshal(b, &ev)
		return ev
	}

	if ev := get(); len(ev.Modules) != 11 || slices.Contains(ev.Modules, "ammo") || len(ev.InUse) != 0 {
		t.Fatalf("a new event has every module but the ammunition desk: %v %v", ev.Modules, ev.InUse)
	}
	code, b := e.do(boss, "GET", "/api/modules", nil)
	expect(t, "list modules", code, 200, b)

	ranges := []int64{get().Ranges[0].ID}
	update := func(mods ...string) (int, []byte) {
		return e.do(boss, "PUT", evPath, map[string]any{"name": "Indoor", "startsOn": "2026-11-20", "rangeIds": ranges, "modules": mods})
	}
	code, b = update("special_prizes")
	expect(t, "switch modules off", code, 204, b)
	if ev := get(); fmt.Sprint(ev.Modules) != "[awards special_prizes]" {
		t.Fatalf("required modules come along: %v", ev.Modules)
	}

	code, b = e.do(boss, "GET", evPath+"/kiosk-tokens", nil)
	expect(t, "kiosk off", code, 404, b)
	if !strings.Contains(string(b), "module is off") {
		t.Fatalf("message: %s", b)
	}
	code, b = e.do(boss, "GET", evPath+"/reports/club_competition", nil)
	expect(t, "report of a module that is off", code, 404, b)
	code, b = e.do(boss, "PUT", fmt.Sprintf("%s/programmes/%d/special-prizes", base, f.main), map[string]any{"ranks": 3, "minCategory": 0, "minClub": 0})
	expect(t, "special prizes are on", code, 204, b)

	code, b = update("kiosk", "special_prizes")
	expect(t, "switch the kiosk on", code, 204, b)
	code, b = e.do(boss, "POST", evPath+"/kiosk-tokens", map[string]string{"name": "Hall"})
	expect(t, "kiosk link", code, 201, b)
	if ev := get(); fmt.Sprint(ev.InUse) != "[special_prizes kiosk]" {
		t.Fatalf("modules in use: %v", ev.InUse)
	}
	code, b = update("special_prizes")
	expect(t, "kiosk with links", code, 409, b)
	if !strings.Contains(string(b), "kiosk links") {
		t.Fatalf("message: %s", b)
	}

	code, b = e.do(boss, "PUT", evPath, map[string]any{"name": "Renamed", "startsOn": "2026-11-20", "rangeIds": ranges})
	expect(t, "update without modules", code, 204, b)
	if ev := get(); fmt.Sprint(ev.Modules) != "[awards special_prizes kiosk]" {
		t.Fatalf("modules unchanged: %v", ev.Modules)
	}
	code, b = e.do(boss, "POST", evPath+"/copy", map[string]string{"name": "Next year", "startsOn": "2027-11-20"})
	expect(t, "copy", code, 201, b)
	var copied store.Event
	json.Unmarshal(b, &copied)
	if fmt.Sprint(copied.Modules) != "[awards special_prizes kiosk]" {
		t.Fatalf("copied modules: %v", copied.Modules)
	}
}

// TestClubModules checks the gate of a club module and that it cannot be switched off while it holds data.
func TestClubModules(t *testing.T) {
	e := newEnv(t)
	tid := e.tenant("Club modules")
	e.user("boss", tid, store.RoleAdmin)
	boss := e.login("boss")
	base := fmt.Sprintf("/api/tenants/%d", tid)
	me := func() store.Membership {
		var m struct{ Tenants []store.Membership }
		_, b := e.do(boss, "GET", "/api/me", nil)
		json.Unmarshal(b, &m)
		return m.Tenants[0]
	}
	if m := me(); fmt.Sprint(m.Modules) != "[forms]" {
		t.Fatalf("a new club has print forms: %v", m.Modules)
	}
	code, b := e.do(boss, "POST", base+"/layouts", map[string]any{"name": "Sheet", "width": 148, "height": 210, "elements": []map[string]any{}})
	expect(t, "create a form", code, 201, b)
	var l store.Layout
	json.Unmarshal(b, &l)
	code, b = e.do(boss, "PUT", base+"/modules", map[string]any{"modules": []string{}})
	expect(t, "forms with a form", code, 409, b)
	if !strings.Contains(string(b), "print forms") {
		t.Fatalf("message: %s", b)
	}
	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/layouts/%d", base, l.ID), nil)
	expect(t, "delete the form", code, 204, b)
	code, b = e.do(boss, "PUT", base+"/modules", map[string]any{"modules": []string{}})
	expect(t, "switch forms off", code, 204, b)
	if m := me(); len(m.Modules) != 0 {
		t.Fatalf("forms still on: %v", m.Modules)
	}
	code, b = e.do(boss, "GET", base+"/layouts", nil)
	expect(t, "forms are off", code, 404, b)
}
