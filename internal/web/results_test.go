package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func TestResults(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	e.user("helper", tid, store.RoleStaff)
	other := e.tenant("Other")
	e.user("outsider", other, store.RoleStaff)
	staff, outsider := e.login("helper"), e.login("outsider")
	ev := e.event(tid, "Indoor")

	main, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Main",
		Series: []domain.Series{
			{Shots: 10, MaxPoints: 100}, {Shots: 10, MaxPoints: 100, TrackMouches: true},
			{Shots: 10, MaxPoints: 100}, {Shots: 10, MaxPoints: 100},
		},
		Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 20}},
		Awards: []domain.AwardTier{
			{Name: "Badge", Kind: domain.AwardNatural, Limits: []domain.AwardLimit{{AgeCategory: domain.Elite, Points: 300}}},
			{Name: "Wreath", Kind: domain.AwardCard, Value: 25, Limits: []domain.AwardLimit{{AgeCategory: domain.Elite, Points: 350}}},
		},
	})
	section, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Section", Series: []domain.Series{{Shots: 10, MaxPoints: 100}, {Shots: 10, MaxPoints: 100, TrackMouches: true}},
		Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 5}}, Awards: []domain.AwardTier{},
	})
	e.st.SetCombination(ctx, tid, section, &main)
	anna, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Anna", Lastname: "Adult", Birthday: "1990-01-01"}, nil)

	lane := 1
	code, b := e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/registrations", tid, ev), web.RegistrationInput{
		ShooterID: anna, Lane: &lane, Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}, {ProgrammeID: section, CombinedSeries: []int{2, 4}}},
	})
	expect(t, "register", code, 201, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)
	resultsPath := fmt.Sprintf("/api/tenants/%d/registrations/%d/results", tid, reg.ID)
	entryPath := func(programme int64) string {
		for _, en := range reg.Entries {
			if en.ProgrammeID == programme {
				return fmt.Sprintf("/api/tenants/%d/entries/%d", tid, en.ID)
			}
		}
		t.Fatalf("no entry for programme %d", programme)
		return ""
	}
	sheetOf := func(b []byte) (web.ResultSheet, map[string]web.EntryResult) {
		var s web.ResultSheet
		json.Unmarshal(b, &s)
		m := map[string]web.EntryResult{}
		for _, er := range s.Entries {
			m[er.Title] = er
		}
		return s, m
	}

	code, b = e.do(staff, "GET", resultsPath, nil)
	expect(t, "result sheet", code, 200, b)
	_, entries := sheetOf(b)
	if entries["Section"].CombinedFrom != "Main" || entries["Main"].Status != "pending" || entries["Main"].MaxPoints != 400 {
		t.Fatalf("initial sheet: %+v", entries)
	}

	series := func(points ...float64) map[string]any {
		out := []map[string]any{}
		for i, p := range points {
			r := map[string]any{"points": p}
			if i == 1 {
				r["mouches"] = 3
			}
			out = append(out, r)
		}
		return map[string]any{"series": out}
	}
	code, b = e.do(staff, "PUT", entryPath(section)+"/results", series(1, 2))
	expect(t, "results for a combined programme", code, 409, b)
	code, b = e.do(staff, "PUT", entryPath(main)+"/results", series(90, 95, 80))
	expect(t, "missing series", code, 422, b)
	code, b = e.do(staff, "PUT", entryPath(main)+"/results", series(90, 95, 80, 101))
	expect(t, "points above maximum", code, 422, b)

	code, b = e.do(staff, "PUT", entryPath(main)+"/results", series(90, 95, 80, 99))
	expect(t, "save results", code, 200, b)
	sheet, entries := sheetOf(b)
	m, s := entries["Main"], entries["Section"]
	if m.Total != 364 || m.Mouches != 3 || m.Status != "completed" || len(m.EligibleAwards) != 2 || m.EligibleAwards[0].Name != "Wreath" {
		t.Fatalf("main after save: %+v", m)
	}
	if s.Total != 194 || s.Mouches != 3 || s.Status != "completed" || len(s.Results) != 2 || s.Results[1].Points != 99 {
		t.Fatalf("section copies series 2 and 4: %+v", s)
	}
	if sheet.Registration.Status != "completed" {
		t.Fatalf("registration not completed: %s", sheet.Registration.Status)
	}

	if m.Award == nil || m.Award.Name != "Wreath" {
		t.Fatalf("most valuable award not given on save: %+v", m)
	}

	code, b = e.do(staff, "PUT", entryPath(section)+"/award", map[string]any{"tierId": m.EligibleAwards[0].ID})
	expect(t, "award of another programme", code, 422, b)
	code, b = e.do(staff, "PUT", entryPath(main)+"/award", map[string]any{"tierId": m.EligibleAwards[1].ID})
	expect(t, "choose award", code, 200, b)
	_, entries = sheetOf(b)
	if entries["Main"].Award == nil || entries["Main"].Award.Name != "Badge" {
		t.Fatalf("award not stored: %+v", entries["Main"])
	}
	code, b = e.do(staff, "PUT", entryPath(main)+"/results", series(90, 95, 80, 98))
	expect(t, "corrected results", code, 200, b)
	if _, entries = sheetOf(b); entries["Main"].Award == nil || entries["Main"].Award.Name != "Badge" {
		t.Fatalf("chosen award should survive a correction: %+v", entries["Main"])
	}
	code, b = e.do(staff, "PUT", entryPath(main)+"/award", map[string]any{"tierId": m.EligibleAwards[0].ID})
	expect(t, "choose wreath again", code, 200, b)

	code, b = e.do(staff, "PUT", entryPath(main)+"/results", series(80, 80, 80, 80))
	expect(t, "lower results", code, 200, b)
	_, entries = sheetOf(b)
	if entries["Main"].Award == nil || entries["Main"].Award.Name != "Badge" || len(entries["Main"].EligibleAwards) != 1 {
		t.Fatalf("award should fall back to Badge, the only one reached: %+v", entries["Main"])
	}
	wreath := int64(0)
	for _, tr := range m.EligibleAwards {
		if tr.Name == "Wreath" {
			wreath = tr.ID
		}
	}
	code, b = e.do(staff, "PUT", entryPath(main)+"/award", map[string]any{"tierId": wreath})
	expect(t, "award not reached", code, 422, b)

	regPath := fmt.Sprintf("/api/tenants/%d/registrations/%d", tid, reg.ID)
	code, b = e.do(staff, "PUT", regPath, web.RegistrationUpdate{Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}, {ProgrammeID: section, CombinedSeries: []int{1, 3}}}})
	expect(t, "change combined series", code, 200, b)
	code, b = e.do(staff, "GET", resultsPath, nil)
	_, entries = sheetOf(b)
	if entries["Section"].Total != 160 || entries["Section"].Mouches != 0 {
		t.Fatalf("section after recombination: %+v", entries["Section"])
	}
	code, b = e.do(staff, "PUT", regPath, web.RegistrationUpdate{Programmes: []domain.ProgrammeChoice{{ProgrammeID: section}}})
	expect(t, "remove scored programme", code, 409, b)

	code, b = e.do(outsider, "GET", resultsPath, nil)
	expect(t, "outsider reads results", code, 404, b)
	code, b = e.do(outsider, "PUT", fmt.Sprintf("/api/tenants/%d/entries/%d/results", other, reg.Entries[0].ID), series(1, 1, 1, 1))
	expect(t, "outsider saves results via own tenant", code, 404, b)
}

// TestConcurrentSaves saves results of the same entry from many goroutines at once; before the per-registration
// lock two saves could both delete and then both insert, which failed on the results primary key.
func TestConcurrentSaves(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	ev := e.event(tid, "Indoor")
	main, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Main", Series: []domain.Series{{Shots: 10, MaxPoints: 100}, {Shots: 10, MaxPoints: 100}},
		Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 20}}, Awards: []domain.AwardTier{},
	})
	sh, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Anna", Lastname: "Adult", Birthday: "1990-01-01"}, nil)
	id, err := e.st.CreateRegistration(ctx, tid, store.NewRegistration{EventID: ev, ShooterID: sh, AgeCategory: domain.Elite,
		Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}}})
	if err != nil {
		t.Fatal(err)
	}
	reg, _ := e.st.GetRegistration(ctx, tid, id)
	entry, _ := e.st.GetEntry(ctx, tid, reg.Entries[0].ID)
	errs := make(chan error, 20)
	for i := range cap(errs) {
		go func() {
			_, err := e.st.SaveResults(ctx, tid, entry, []domain.SeriesResult{{Points: float64(80 + i)}, {Points: 90}})
			errs <- err
		}()
	}
	for range cap(errs) {
		if err := <-errs; err != nil {
			t.Fatalf("concurrent save: %v", err)
		}
	}
	res, _ := e.st.ListResults(ctx, tid, id)
	if len(res[entry.ID]) != 2 {
		t.Fatalf("results after concurrent saves: %+v", res)
	}
}

func TestSingleShotResults(t *testing.T) {
	e := newEnv(t)
	f := e.registrationFixture("Shots")
	e.user("boss", f.tenant, store.RoleAdmin)
	boss := e.login("boss")
	ctx := context.Background()
	p, err := e.st.CreateProgramme(ctx, f.tenant, f.event, domain.ProgrammeDefinition{Name: "Obligatory", Awards: []domain.AwardTier{},
		Series: []domain.Series{{Shots: 3, MaxPoints: 15, TrackShots: true}}, Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 0}}})
	if err != nil {
		t.Fatal(err)
	}
	reg, err := e.st.CreateRegistration(ctx, f.tenant, store.NewRegistration{EventID: f.event, ShooterID: f.shooter, AgeCategory: domain.Elite,
		Programmes: []domain.ProgrammeChoice{{ProgrammeID: p}}})
	if err != nil {
		t.Fatal(err)
	}
	r, _ := e.st.GetRegistration(ctx, f.tenant, reg)
	path := fmt.Sprintf("/api/tenants/%d/entries/%d/results", f.tenant, r.Entries[0].ID)

	code, b := e.do(boss, "PUT", path, map[string]any{"series": []map[string]any{{"points": 0, "shots": []float64{5, 4, 5}}}})
	expect(t, "save shots", code, 200, b)
	var sheet web.ResultSheet
	json.Unmarshal(b, &sheet)
	if res := sheet.Entries[0].Results; len(res) != 1 || res[0].Points != 14 || fmt.Sprint(res[0].Shots) != "[5 4 5]" || sheet.Entries[0].Total != 14 {
		t.Fatalf("sheet with shots: %s", b)
	}
	code, b = e.do(boss, "PUT", path, map[string]any{"series": []map[string]any{{"points": 0, "shots": []float64{5, 4}}}})
	expect(t, "too few shots", code, 422, b)
}

func TestEnterResult(t *testing.T) {
	e := newEnv(t)
	f := e.registrationFixture("Quick")
	e.user("helper", f.tenant, store.RoleStaff)
	other := e.tenant("Other")
	e.user("outsider", other, store.RoleAdmin)
	staff := e.login("helper")
	ctx := context.Background()
	path := fmt.Sprintf("/api/tenants/%d/events/%d/results", f.tenant, f.event)
	line := func(shooter, programme int64, points ...float64) map[string]any {
		series := []map[string]any{}
		for _, p := range points {
			series = append(series, map[string]any{"points": p})
		}
		return map[string]any{"shooterId": shooter, "programmeId": programme, "clubId": f.club, "series": series}
	}

	code, b := e.do(staff, "POST", path, line(f.shooter, f.main, 90, 91))
	expect(t, "too few series", code, 422, b)
	if regs, _ := e.st.ListRegistrations(ctx, f.tenant, f.event, false); len(regs) != 0 {
		t.Fatalf("an invalid line left a registration: %+v", regs)
	}
	code, b = e.do(staff, "POST", path, line(f.shooter, f.main, 90, 91, 92, 93))
	expect(t, "enter a new result", code, 200, b)
	var sheet web.ResultSheet
	json.Unmarshal(b, &sheet)
	if sheet.Registration.Status != "completed" || sheet.Registration.Lane != nil || sheet.Entries[0].Total != 366 || sheet.Registration.ClubID == nil {
		t.Fatalf("new result: %s", b)
	}
	code, b = e.do(staff, "POST", path, line(f.shooter, f.main, 90, 91, 92, 93))
	expect(t, "the same programme twice", code, 409, b)

	desk, _ := e.st.CreateRegistration(ctx, f.tenant, store.NewRegistration{EventID: f.event, ShooterID: f.youth, AgeCategory: domain.U21,
		Programmes: []domain.ProgrammeChoice{{ProgrammeID: f.junior}}})
	code, b = e.do(staff, "POST", path, line(f.youth, f.junior, 50, 60))
	expect(t, "fill a registration from the desk", code, 200, b)
	json.Unmarshal(b, &sheet)
	if sheet.Registration.ID != desk || sheet.Entries[0].Total != 110 {
		t.Fatalf("desk registration: %s", b)
	}

	code, b = e.do(e.login("outsider"), "POST", fmt.Sprintf("/api/tenants/%d/events/%d/results", other, f.event), line(f.shooter, f.section, 50, 50))
	expect(t, "other tenant", code, 404, b)
}
