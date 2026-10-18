package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func definition(name, variant string) map[string]any {
	return map[string]any{
		"name": name, "variant": variant, "showOnRankings": true,
		"series": []map[string]any{{"shots": 10, "maxPoints": 100}, {"shots": 10, "maxPoints": 100, "trackMouches": true}},
		"prices": []map[string]any{{"ageCategory": "elite", "price": 15}, {"ageCategory": "u21", "price": 5}},
		"awards": []map[string]any{{"name": "Wreath", "kind": "card", "value": 20, "limits": []map[string]any{{"ageCategory": "elite", "points": 180}}}},
	}
}

func (e *env) event(tenantID int64, name string) int64 {
	e.t.Helper()
	ctx := context.Background()
	ranges, _ := e.st.ListRanges(ctx, tenantID)
	if len(ranges) == 0 {
		r, err := e.st.CreateRange(ctx, tenantID, store.RangeInput{Name: "10m", LaneCount: 10})
		if err != nil {
			e.t.Fatal(err)
		}
		ranges = append(ranges, r)
	}
	id, err := e.st.CreateEvent(ctx, tenantID, store.EventInput{RangeIDs: []int64{ranges[0].ID}, Name: name, StartsOn: "2026-11-20", Timezone: "Europe/Zurich"})
	if err != nil {
		e.t.Fatal(err)
	}
	return id
}

func TestProgrammes(t *testing.T) {
	e := newEnv(t)
	tid, other := e.tenant("Club"), e.tenant("Other")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	e.user("outsider", other, store.RoleAdmin)
	admin, staff, outsider := e.login("boss"), e.login("helper"), e.login("outsider")
	ev, next := e.event(tid, "Indoor 2026"), e.event(tid, "Indoor 2027")
	evPath := fmt.Sprintf("/api/tenants/%d/events/%d/programmes", tid, ev)

	bad := definition("Field match", "")
	bad["series"] = []map[string]any{}
	code, b := e.do(admin, "POST", evPath, bad)
	expect(t, "programme without series", code, 422, b)
	code, b = e.do(staff, "POST", evPath, definition("Field match", "Stgw 90"))
	expect(t, "staff creates programme", code, 403, b)

	code, b = e.do(admin, "POST", evPath, definition("Field match", "Stgw 90"))
	expect(t, "create programme", code, 201, b)
	var p90 store.Programme
	json.Unmarshal(b, &p90)
	if p90.Shots != 20 || p90.MaxPoints != 200 || len(p90.Awards) != 1 || !p90.Active {
		t.Fatalf("created programme: %+v", p90)
	}
	code, b = e.do(admin, "POST", evPath, definition("Field match", "Stgw 90"))
	expect(t, "duplicate programme", code, 409, b)
	code, b = e.do(admin, "POST", evPath, definition("Field match", "Stgw 57"))
	expect(t, "second variant", code, 201, b)
	var p57 store.Programme
	json.Unmarshal(b, &p57)

	progPath := func(id int64) string { return fmt.Sprintf("/api/tenants/%d/programmes/%d", tid, id) }
	code, b = e.do(admin, "PUT", progPath(p57.ID)+"/combination", map[string]any{"combinableWithId": p90.ID})
	expect(t, "set combination", code, 204, b)
	code, b = e.do(admin, "PUT", progPath(p57.ID)+"/active", map[string]any{"active": false})
	expect(t, "deactivate", code, 204, b)

	code, b = e.do(staff, "GET", evPath, nil)
	expect(t, "staff lists programmes", code, 200, b)
	var list []store.Programme
	json.Unmarshal(b, &list)
	if len(list) != 2 || list[0].Variant != "Stgw 57" || list[0].Active || list[0].CombinableWithID == nil {
		t.Fatalf("list: %+v", list)
	}

	resp, err := admin.Get(e.srv.URL + progPath(p90.ID) + "/export")
	if err != nil {
		t.Fatal(err)
	}
	var exp web.ProgrammeExport
	json.NewDecoder(resp.Body).Decode(&exp)
	resp.Body.Close()
	if cd := resp.Header.Get("Content-Disposition"); cd != `attachment; filename="field-match-stgw-90.programme.json"` {
		t.Fatalf("content disposition: %q", cd)
	}
	if exp.Format != web.ProgrammeFormat || exp.Programme.Name != "Field match" || len(exp.Programme.Series) != 2 {
		t.Fatalf("export: %+v", exp)
	}

	nextPath := fmt.Sprintf("/api/tenants/%d/events/%d/programmes", tid, next)
	code, b = e.do(admin, "POST", nextPath+"/import", exp)
	expect(t, "import", code, 201, b)
	src, _ := e.st.ListProgrammes(context.Background(), tid, ev)
	for _, p := range src {
		e.st.SetSpecialPrizes(context.Background(), tid, p.ID, domain.SpecialPrizeRules{Ranks: 3, MinCategory: 10, MinClub: 4})
	}
	code, b = e.do(admin, "POST", nextPath+"/copy", map[string]any{"fromEventId": ev})
	expect(t, "copy", code, 200, b)
	if string(b) != `{"copied":1}`+"\n" && string(b) != `{"copied":1}` {
		t.Fatalf("copy result: %s", b)
	}

	tplID, _ := e.st.SaveTemplate(context.Background(), domain.ProgrammeDefinition{
		Name: "Obligatory", Series: []domain.Series{{Shots: 20, MaxPoints: 100}}, Prices: []domain.Price{}, Awards: []domain.AwardTier{},
	})
	code, b = e.do(staff, "GET", "/api/programme-templates", nil)
	expect(t, "list templates", code, 200, b)
	code, b = e.do(admin, "POST", nextPath+"/from-template", map[string]any{"templateId": tplID})
	expect(t, "from template", code, 201, b)
	code, b = e.do(admin, "GET", nextPath, nil)
	json.Unmarshal(b, &list)
	if len(list) != 3 {
		t.Fatalf("next event programmes: %+v", list)
	}
	withPrizes := 0
	for _, p := range list {
		if p.SpecialPrizes == (domain.SpecialPrizeRules{Ranks: 3, MinCategory: 10, MinClub: 4}) {
			withPrizes++
		}
	}
	if withPrizes != 1 {
		t.Fatalf("the copied programme keeps its special prizes: %+v", list)
	}

	code, b = e.do(outsider, "GET", progPath(p90.ID), nil)
	expect(t, "outsider reads programme", code, 404, b)
	code, b = e.do(outsider, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/programmes", other, ev), definition("x", ""))
	expect(t, "outsider creates programme in foreign event", code, 404, b)
	code, b = e.do(outsider, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/programmes/copy", other, ev), map[string]any{"fromEventId": ev})
	expect(t, "outsider copies into foreign event", code, 404, b)

	code, b = e.do(admin, "DELETE", progPath(p90.ID), nil)
	expect(t, "delete programme", code, 204, b)
}

func TestAdminTemplates(t *testing.T) {
	e := newEnv(t)
	c := e.client()
	resp, err := c.PostForm(e.srv.URL+"/admin/login", url.Values{"token": {adminToken}})
	if err != nil || resp.StatusCode != 303 {
		t.Fatalf("admin login: %v %v", resp.StatusCode, err)
	}
	upload := func(content string) int {
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		part, _ := w.CreateFormFile("file", "p.programme.json")
		part.Write([]byte(content))
		w.Close()
		req, _ := http.NewRequest("POST", e.srv.URL+"/admin/templates", &buf)
		req.Header.Set("Content-Type", w.FormDataContentType())
		resp, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		return resp.StatusCode
	}
	doc, _ := json.Marshal(map[string]any{"format": web.ProgrammeFormat, "version": 1, "programme": definition("Obligatory", "")})
	if code := upload(string(doc)); code != 303 {
		t.Fatalf("upload template: %d", code)
	}
	if code := upload(`{"format":"something else"}`); code != 422 {
		t.Fatalf("upload foreign file: %d", code)
	}
	invalid, _ := json.Marshal(map[string]any{"format": web.ProgrammeFormat, "version": 1, "programme": map[string]any{"name": "x"}})
	if code := upload(string(invalid)); code != 422 {
		t.Fatalf("upload invalid programme: %d", code)
	}
	list, _ := e.st.ListTemplates(context.Background())
	if len(list) != 1 || list[0].Definition.Name != "Obligatory" {
		t.Fatalf("templates: %+v", list)
	}
	resp, err = c.Get(fmt.Sprintf("%s/admin/templates/%d", e.srv.URL, list[0].ID))
	if err != nil {
		t.Fatal(err)
	}
	var exp web.ProgrammeExport
	json.NewDecoder(resp.Body).Decode(&exp)
	resp.Body.Close()
	if exp.Programme.Name != "Obligatory" || len(exp.Programme.Series) != 2 {
		t.Fatalf("download: %+v", exp)
	}
}

func TestCopyEvent(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	other := e.tenant("Other")
	e.user("outsider", other, store.RoleAdmin)
	admin, staff, outsider := e.login("boss"), e.login("helper"), e.login("outsider")
	ev := e.event(tid, "41. Lupi 2026")
	home, _ := e.st.CreateClub(ctx, tid, store.ClubInput{Name: "Home"})
	src, _ := e.st.GetEvent(ctx, tid, ev)
	e.st.UpdateEvent(ctx, tid, ev, store.EventInput{RangeIDs: []int64{src.Ranges[0].ID}, Name: src.Name, StartsOn: src.StartsOn, Organizer: "PS Appenzell",
		Timezone: src.Timezone, Locale: "de", Settings: store.EventSettings{ExcludedClubID: &home.ID, ShowLaneStatus: true}})

	def := func(name, variant string) domain.ProgrammeDefinition {
		return domain.ProgrammeDefinition{Name: name, Variant: variant, Series: []domain.Series{{Shots: 10, MaxPoints: 100}},
			Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 20}}, Awards: []domain.AwardTier{}}
	}
	main, _ := e.st.CreateProgramme(ctx, tid, ev, def("Meisterschaft", "Sport"))
	extra, _ := e.st.CreateProgramme(ctx, tid, ev, def("Biberlistich", "Sport"))
	old, _ := e.st.CreateProgramme(ctx, tid, ev, def("Old", ""))
	e.st.SetCombination(ctx, tid, extra, &main)
	e.st.SetSpecialPrizes(ctx, tid, main, domain.SpecialPrizeRules{Ranks: 3, MinCategory: 10, MinClub: 4})
	e.st.SetProgrammeActive(ctx, tid, old, false)
	anna, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Anna", Lastname: "Adult", Birthday: "1990-01-01"}, nil)
	if _, err := e.st.CreateRegistration(ctx, tid, store.NewRegistration{EventID: ev, ShooterID: anna, AgeCategory: domain.Elite,
		Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}}}); err != nil {
		t.Fatal(err)
	}
	e.st.SetEventArchived(ctx, tid, ev, true)

	path := fmt.Sprintf("/api/tenants/%d/events/%d/copy", tid, ev)
	body := map[string]string{"name": "42. Lupi 2027", "startsOn": "2027-01-08"}
	code, b := e.do(staff, "POST", path, body)
	expect(t, "staff copies an event", code, 403, b)
	code, b = e.do(outsider, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/copy", other, ev), body)
	expect(t, "copy an event of another tenant", code, 404, b)
	code, b = e.do(admin, "POST", path, body)
	expect(t, "copy an archived event", code, 201, b)
	var copied store.Event
	json.Unmarshal(b, &copied)
	if copied.Name != "42. Lupi 2027" || copied.StartsOn != "2027-01-08" || copied.Ranges[0].ID != src.Ranges[0].ID || copied.Organizer != "PS Appenzell" ||
		copied.Locale != "de" || copied.ArchivedAt != nil || copied.Settings.ExcludedClubID == nil || *copied.Settings.ExcludedClubID != home.ID ||
		!copied.Settings.ShowLaneStatus {
		t.Fatalf("copied event: %+v", copied)
	}

	programmes, _ := e.st.ListProgrammes(ctx, tid, copied.ID)
	byName := map[string]store.Programme{}
	for _, p := range programmes {
		byName[p.Name] = p
	}
	m, x, o := byName["Meisterschaft"], byName["Biberlistich"], byName["Old"]
	if len(programmes) != 3 || m.ID == main || m.SpecialPrizes.Ranks != 3 || x.CombinableWithID == nil || *x.CombinableWithID != m.ID || o.Active || !m.Active {
		t.Fatalf("copied programmes: %+v", programmes)
	}
	if regs, _ := e.st.ListRegistrations(ctx, tid, copied.ID, true); len(regs) != 0 {
		t.Fatalf("registrations must not be copied: %+v", regs)
	}
}

// TestEventTemplates checks saving an event as a template, creating an event from it and that a template takes no
// registrations.
func TestEventTemplates(t *testing.T) {
	e := newEnv(t)
	f := e.registrationFixture("Templates")
	e.user("boss", f.tenant, store.RoleAdmin)
	boss := e.login("boss")
	base := fmt.Sprintf("/api/tenants/%d", f.tenant)
	code, b := e.do(boss, "POST", fmt.Sprintf("%s/events/%d/copy", base, f.event), map[string]any{"name": "Indoor", "startsOn": "2026-11-20", "template": true})
	expect(t, "save as template", code, 201, b)
	var tpl store.Event
	json.Unmarshal(b, &tpl)
	if !tpl.Template {
		t.Fatalf("not a template: %s", b)
	}
	code, b = e.do(boss, "POST", fmt.Sprintf("%s/events/%d/registrations", base, tpl.ID), web.RegistrationInput{ShooterID: f.shooter, Programmes: []domain.ProgrammeChoice{}})
	expect(t, "register on a template", code, 409, b)

	code, b = e.do(boss, "POST", fmt.Sprintf("%s/events/%d/copy", base, tpl.ID), map[string]any{"name": "Indoor 2027", "startsOn": "2027-11-19"})
	expect(t, "event from the template", code, 201, b)
	var ev store.Event
	json.Unmarshal(b, &ev)
	if ev.Template || ev.StartsOn != "2027-11-19" {
		t.Fatalf("event from template: %s", b)
	}
	code, b = e.do(boss, "GET", fmt.Sprintf("%s/events/%d/programmes", base, ev.ID), nil)
	expect(t, "programmes of the new event", code, 200, b)
	if !strings.Contains(string(b), `"Main"`) {
		t.Fatalf("programmes not copied: %s", b)
	}

	other := e.tenant("Other")
	e.user("eve", other, store.RoleAdmin)
	code, b = e.do(e.login("eve"), "POST", fmt.Sprintf("/api/tenants/%d/events/%d/copy", other, tpl.ID), map[string]any{"name": "Stolen", "startsOn": "2027-01-01"})
	expect(t, "template of another club", code, 404, b)
}
