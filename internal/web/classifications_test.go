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

func TestClassifications(t *testing.T) {
	e := newEnv(t)
	f := e.registrationFixture("Classes")
	e.user("boss", f.tenant, store.RoleAdmin)
	e.user("helper", f.tenant, store.RoleStaff)
	boss, staff := e.login("boss"), e.login("helper")
	ctx := context.Background()
	shoot := func(shooter int64, cat domain.AgeCategory, programme int64, points ...float64) {
		t.Helper()
		id, err := e.st.CreateRegistration(ctx, f.tenant, store.NewRegistration{EventID: f.event, ShooterID: shooter, AgeCategory: cat,
			Programmes: []domain.ProgrammeChoice{{ProgrammeID: programme}}})
		if err != nil {
			t.Fatal(err)
		}
		reg, _ := e.st.GetRegistration(ctx, f.tenant, id)
		series := []domain.SeriesResult{}
		for _, p := range points {
			series = append(series, domain.SeriesResult{Points: p})
		}
		if _, err := e.st.SaveResults(ctx, f.tenant, store.EntryRef{ID: reg.Entries[0].ID, RegistrationID: id, EventID: f.event, ProgrammeID: programme}, series); err != nil {
			t.Fatal(err)
		}
	}
	shoot(f.shooter, domain.Elite, f.main, 90, 90, 90, 90)
	shoot(f.youth, domain.U21, f.main, 95, 95, 95, 95)
	shoot(f.youth, domain.U21, f.junior, 50, 50)

	base := fmt.Sprintf("/api/tenants/%d/events/%d/classifications", f.tenant, f.event)
	c := domain.ClassificationDefinition{Name: "Championship", TieBreak: domain.SSVTieBreak, ShowOnKiosk: true,
		Parts:   []domain.ClassificationPart{{ProgrammeID: f.main, Count: 1}, {ProgrammeID: f.junior, Count: 1}},
		Classes: []domain.ClassificationClass{{Name: "Adults", Categories: []domain.AgeCategory{domain.Elite}}, {Name: "Juniors", Categories: []domain.AgeCategory{domain.U21}}}}
	code, b := e.do(staff, "POST", base, c)
	expect(t, "staff creates a classification", code, 403, b)
	bad := c
	bad.Parts = []domain.ClassificationPart{{ProgrammeID: 999, Count: 1}}
	code, b = e.do(boss, "POST", base, bad)
	expect(t, "programme of another event", code, 422, b)
	code, b = e.do(boss, "POST", base, c)
	expect(t, "create a classification", code, 201, b)
	var created store.Classification
	json.Unmarshal(b, &created)

	code, b = e.do(staff, "GET", fmt.Sprintf("/api/tenants/%d/events/%d/classification-rankings", f.tenant, f.event), nil)
	expect(t, "classification rankings", code, 200, b)
	var ranked []web.ClassificationRanking
	json.Unmarshal(b, &ranked)
	if len(ranked) != 1 || len(ranked[0].Parts) != 2 || ranked[0].Parts[0].Title != "Main" || len(ranked[0].Classes) != 2 {
		t.Fatalf("classification: %s", b)
	}
	adults, juniors := ranked[0].Classes[0], ranked[0].Classes[1]
	if len(adults.Lines) != 1 || adults.Lines[0].Total != 360 || adults.Lines[0].Parts[1].Points != nil {
		t.Fatalf("adults: %+v", adults)
	}
	if len(juniors.Lines) != 1 || juniors.Lines[0].Total != 480 || juniors.Lines[0].Rank != 1 || juniors.Lines[0].Lastname != "Young" {
		t.Fatalf("juniors: %+v", juniors)
	}

	code, b = e.do(boss, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/kiosk-tokens", f.tenant, f.event), map[string]any{"name": "Hall"})
	expect(t, "kiosk token", code, 201, b)
	var token struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(b, &token)
	kioskCount := func() int {
		t.Helper()
		code, b := e.kiosk("GET", "/api/kiosk/classifications", token.Secret)
		expect(t, "kiosk classifications", code, 200, b)
		var list []web.ClassificationRanking
		json.Unmarshal(b, &list)
		return len(list)
	}
	if n := kioskCount(); n != 1 {
		t.Fatalf("kiosk shows %d classifications", n)
	}

	c.Complete, c.ShowOnKiosk = true, false
	code, b = e.do(boss, "PUT", fmt.Sprintf("%s/%d", base, created.ID), c)
	expect(t, "update", code, 204, b)
	code, b = e.do(staff, "GET", fmt.Sprintf("/api/tenants/%d/events/%d/classification-rankings", f.tenant, f.event), nil)
	expect(t, "classification rankings", code, 200, b)
	json.Unmarshal(b, &ranked)
	if a := ranked[0].Classes[0]; len(a.Lines) != 0 || len(a.Unranked) != 1 || a.Unranked[0].Rank != 0 {
		t.Fatalf("incomplete adult is unranked: %+v", a)
	}

	if n := kioskCount(); n != 0 {
		t.Fatalf("hidden from kiosks, but %d shown", n)
	}
	code, b = e.do(boss, "GET", fmt.Sprintf("/api/tenants/%d/events/%d/reports/classifications?format=csv", f.tenant, f.event), nil)
	expect(t, "report", code, 200, b)
	if got := string(b); !strings.Contains(got, "Championship – Juniors") || !strings.Contains(got, "Young Yves") {
		t.Fatalf("report: %s", got)
	}

	evPath := fmt.Sprintf("/api/tenants/%d/events/%d", f.tenant, f.event)
	code, b = e.do(boss, "PUT", evPath+"/archived", map[string]bool{"archived": true})
	expect(t, "close the event", code, 204, b)
	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/%d", base, created.ID), nil)
	expect(t, "delete in a closed event", code, 409, b)
	code, b = e.do(boss, "PUT", evPath+"/archived", map[string]bool{"archived": false})
	expect(t, "reopen", code, 204, b)
	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/%d", base, created.ID), nil)
	expect(t, "delete", code, 204, b)
	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/%d", base, created.ID), nil)
	expect(t, "delete again", code, 404, b)
}
