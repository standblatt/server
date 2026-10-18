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

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func (e *env) helper(method, path, token string, body any) (int, []byte) {
	e.t.Helper()
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, e.srv.URL+path, r)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b
}

func TestHelperLinks(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	admin, staff := e.login("boss"), e.login("helper")
	ev := e.event(tid, "Match")
	got, _ := e.st.GetEvent(ctx, tid, ev)
	if err := e.st.UpdateEvent(ctx, tid, ev, store.EventInput{Name: got.Name, StartsOn: got.StartsOn, RangeIDs: []int64{got.Ranges[0].ID}, Timezone: got.Timezone, Modules: append(got.Modules, "ammo")}); err != nil {
		t.Fatal(err)
	}
	other := e.event(tid, "Other match")
	ranges, _ := e.st.ListRanges(ctx, tid)
	rangeID := ranges[0].ID
	def := domain.ProgrammeDefinition{
		Name: "Main", Series: []domain.Series{{Shots: 10, MaxPoints: 100}, {Shots: 10, MaxPoints: 100}},
		Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 10}}, Awards: []domain.AwardTier{},
	}
	main, _ := e.st.CreateProgramme(ctx, tid, ev, def)
	side, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{Name: "Side", Series: def.Series, Prices: def.Prices, Awards: def.Awards})
	otherMain, _ := e.st.CreateProgramme(ctx, tid, other, def)

	register := func(event, programme int64, first string, lane int) store.Registration {
		id, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: first, Lastname: "Tester", Birthday: "1990-01-01"}, nil)
		var l *int
		if lane > 0 {
			l = &lane
		}
		code, b := e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/registrations", tid, event),
			web.RegistrationInput{ShooterID: id, Lane: l, Programmes: []domain.ProgrammeChoice{{ProgrammeID: programme}}})
		expect(t, "register "+first, code, 201, b)
		var reg store.Registration
		json.Unmarshal(b, &reg)
		return reg
	}
	anna := register(ev, main, "Anna", 2)
	bert := register(ev, main, "Bert", 7)
	cleo := register(ev, main, "Cleo", 0)
	dora := register(other, otherMain, "Dora", 1)

	codes := fmt.Sprintf("/api/tenants/%d/events/%d/helper-codes", tid, ev)
	code, b := e.do(staff, "POST", codes, map[string]any{"name": "x", "role": "results"})
	expect(t, "staff creates helper link", code, 403, b)
	code, b = e.do(admin, "POST", codes, map[string]any{"name": "Scorer A", "role": "scorer", "rangeId": rangeID, "lanes": "1-3, 11"})
	expect(t, "lane 11 on a 10-lane range", code, 422, b)
	code, b = e.do(admin, "POST", codes, map[string]any{"name": "Scorer A", "role": "scorer", "lanes": "1-3"})
	expect(t, "scorer without range", code, 422, b)
	create := func(body map[string]any) string {
		code, b := e.do(admin, "POST", codes, body)
		expect(t, "create "+body["name"].(string), code, 201, b)
		var out struct {
			Code   store.HelperCode `json:"code"`
			Secret string           `json:"secret"`
		}
		json.Unmarshal(b, &out)
		if !strings.HasPrefix(out.Secret, "sh_") {
			t.Fatalf("secret %q", out.Secret)
		}
		return out.Secret
	}
	scorer := create(map[string]any{"name": "Scorer A", "role": "scorer", "rangeId": rangeID, "lanes": "1-3"})
	ammo := create(map[string]any{"name": "Ammo desk", "role": "ammunition"})
	results := create(map[string]any{"name": "Office", "role": "results"})

	code, b = e.do(admin, "GET", codes, nil)
	expect(t, "list codes", code, 200, b)
	var list []store.HelperCode
	json.Unmarshal(b, &list)
	if len(list) != 3 || list[0].Role != store.HelperScorer || fmt.Sprint(list[0].Lanes) != "[1 2 3]" || list[1].RangeID != nil {
		t.Fatalf("codes: %s", b)
	}

	code, b = e.helper("GET", "/api/helper", scorer, nil)
	expect(t, "helper me", code, 200, b)
	var me store.Helper
	json.Unmarshal(b, &me)
	if me.EventID != ev || me.Role != store.HelperScorer || me.Name != "Scorer A" {
		t.Fatalf("me: %s", b)
	}
	code, b = e.helper("GET", "/api/helper", "sh_nope", nil)
	expect(t, "unknown link", code, 401, b)

	tenants := fmt.Sprintf("/api/tenants/%d", tid)
	results2 := func(points ...float64) map[string]any {
		series := []map[string]any{}
		for _, p := range points {
			series = append(series, map[string]any{"points": p})
		}
		return map[string]any{"series": series}
	}
	// Scorer: own lane yes, other lane no, nothing beyond results and finish.
	code, b = e.helper("GET", fmt.Sprintf("%s/registrations/%d/results", tenants, anna.ID), scorer, nil)
	expect(t, "scorer reads lane 2", code, 200, b)
	code, b = e.helper("PUT", fmt.Sprintf("%s/entries/%d/results", tenants, anna.Entries[0].ID), scorer, results2(90, 95))
	expect(t, "scorer saves lane 2", code, 200, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/registrations/%d/results", tenants, bert.ID), scorer, nil)
	expect(t, "scorer reads lane 7", code, 403, b)
	code, b = e.helper("PUT", fmt.Sprintf("%s/entries/%d/results", tenants, bert.Entries[0].ID), scorer, results2(90, 95))
	expect(t, "scorer saves lane 7", code, 403, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/registrations/%d", tenants, cleo.ID), scorer, nil)
	expect(t, "scorer reads waiting shooter", code, 403, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/registrations/%d", tenants, dora.ID), scorer, nil)
	expect(t, "scorer reads other event", code, 404, b)
	code, b = e.helper("DELETE", fmt.Sprintf("%s/registrations/%d", tenants, anna.ID), scorer, nil)
	expect(t, "scorer deletes", code, 403, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/events/%d/registrations", tenants, ev), scorer, nil)
	expect(t, "scorer lists registrations", code, 200, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/events/%d/registrations", tenants, other), scorer, nil)
	expect(t, "scorer lists other event", code, 404, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/events/%d/lanes", tenants, ev), scorer, nil)
	expect(t, "scorer lane states", code, 200, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/events/%d/programmes", tenants, ev), scorer, nil)
	expect(t, "scorer programmes", code, 200, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/events/%d/stats", tenants, ev), ammo, nil)
	expect(t, "ammo stats", code, 200, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/events/%d/lanes", tenants, other), scorer, nil)
	expect(t, "scorer lane states of other event", code, 404, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/events/%d/rankings", tenants, ev), scorer, nil)
	expect(t, "scorer rankings", code, 403, b)
	code, b = e.helper("POST", fmt.Sprintf("%s/registrations/%d/finish", tenants, anna.ID), scorer, nil)
	expect(t, "scorer finishes lane 2", code, 204, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/events/%d/helper-codes", tenants, ev), scorer, nil)
	expect(t, "scorer lists codes", code, 403, b)

	// Ammunition desk: issue, move to a lane, print the lane sheet, but no programme changes.
	code, b = e.helper("POST", fmt.Sprintf("%s/registrations/%d/ammo", tenants, cleo.ID), ammo, nil)
	expect(t, "ammo issued", code, 204, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/registrations/%d", tenants, cleo.ID), ammo, nil)
	expect(t, "ammo reads", code, 200, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)
	if reg.AmmoIssuedAt == nil {
		t.Fatalf("ammo not recorded: %s", b)
	}
	code, b = e.helper("PUT", fmt.Sprintf("%s/registrations/%d", tenants, cleo.ID), ammo,
		web.RegistrationUpdate{Lane: &[]int{4}[0], Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}, {ProgrammeID: side}}})
	expect(t, "ammo adds a programme", code, 403, b)
	code, b = e.helper("PUT", fmt.Sprintf("%s/registrations/%d", tenants, cleo.ID), ammo,
		web.RegistrationUpdate{Lane: &[]int{4}[0], Programmes: []domain.ProgrammeChoice{{ProgrammeID: main, CombinedSeries: []int{1}}}})
	expect(t, "ammo changes the combined series", code, 403, b)
	code, b = e.helper("PUT", fmt.Sprintf("%s/registrations/%d", tenants, cleo.ID), ammo,
		web.RegistrationUpdate{Lane: &[]int{4}[0], Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}}})
	expect(t, "ammo assigns lane 4", code, 200, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/registrations/%d/documents/lane_sheet", tenants, cleo.ID), ammo, nil)
	expect(t, "ammo lane sheet", code, 200, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/registrations/%d/documents/result_sheet", tenants, cleo.ID), ammo, nil)
	expect(t, "ammo result sheet", code, 403, b)
	code, b = e.helper("POST", tenants+"/print", ammo, web.PrintRequest{Document: "lane_sheet", RegistrationID: cleo.ID})
	expect(t, "ammo prints lane sheet", code, 200, b)
	code, b = e.helper("POST", tenants+"/print", ammo, web.PrintRequest{Document: "result_sheet", RegistrationID: cleo.ID})
	expect(t, "ammo prints result sheet", code, 403, b)
	code, b = e.helper("POST", tenants+"/print", ammo, web.PrintRequest{Document: "wreath_card", EntryID: dora.Entries[0].ID})
	expect(t, "ammo prints a wreath card", code, 403, b)
	code, b = e.helper("PUT", fmt.Sprintf("%s/entries/%d/results", tenants, cleo.Entries[0].ID), ammo, results2(90, 95))
	expect(t, "ammo saves results", code, 403, b)

	// Results clerk: every lane, no finish.
	code, b = e.helper("PUT", fmt.Sprintf("%s/entries/%d/results", tenants, bert.Entries[0].ID), results, results2(80, 85))
	expect(t, "clerk saves lane 7", code, 200, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/events/%d/registrations?all=true", tenants, ev), results, nil)
	expect(t, "clerk lists registrations", code, 200, b)
	code, b = e.helper("POST", fmt.Sprintf("%s/registrations/%d/ammo", tenants, bert.ID), results, nil)
	expect(t, "clerk issues ammo", code, 403, b)
	code, b = e.helper("DELETE", fmt.Sprintf("%s/registrations/%d/awards-issued", tenants, bert.ID), results, nil)
	expect(t, "clerk takes awards back", code, 403, b)
	code, b = e.helper("DELETE", fmt.Sprintf("%s/registrations/%d/ammo", tenants, cleo.ID), ammo, nil)
	expect(t, "ammo desk takes ammo back", code, 403, b)

	// Another club never sees the link.
	tid2 := e.tenant("Other club")
	code, b = e.helper("GET", fmt.Sprintf("/api/tenants/%d/events/%d", tid2, ev), scorer, nil)
	expect(t, "other tenant", code, 404, b)

	// Revoked links stop at once.
	code, b = e.do(admin, "DELETE", fmt.Sprintf("%s/helper-codes/%d", tenants, list[0].ID), nil)
	expect(t, "revoke", code, 204, b)
	code, b = e.helper("GET", "/api/helper", scorer, nil)
	expect(t, "revoked link", code, 401, b)
	code, b = e.helper("GET", fmt.Sprintf("%s/registrations/%d/results", tenants, anna.ID), scorer, nil)
	expect(t, "revoked link on the API", code, 401, b)
}
