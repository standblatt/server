package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func TestReports(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	e.user("helper", tid, store.RoleStaff)
	e.user("boss", tid, store.RoleAdmin)
	other := e.tenant("Other")
	e.user("outsider", other, store.RoleAdmin)
	staff, outsider := e.login("helper"), e.login("outsider")
	ev := e.event(tid, "Indoor")

	main, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Main", CountsForClubCompetition: true, ShowOnRankings: true,
		Series: []domain.Series{{Shots: 10, MaxPoints: 100}, {Shots: 10, MaxPoints: 100}, {Shots: 10, MaxPoints: 100}},
		Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 20}, {AgeCategory: domain.U21, Price: 10}},
		Awards: []domain.AwardTier{{Name: "Wreath", Kind: domain.AwardCard, Value: 15, Limits: []domain.AwardLimit{{AgeCategory: domain.Elite, Points: 270}, {AgeCategory: domain.U21, Points: 260}}}},
	})
	home, _ := e.st.CreateClub(ctx, tid, store.ClubInput{Name: "Home", SSVNumber: "1001"})
	small, _ := e.st.CreateClub(ctx, tid, store.ClubInput{Name: "Small", SSVNumber: "1002"})
	tiers, _ := e.st.ProgrammeTiers(ctx, tid, []int64{main})

	shoot := func(first, birthday string, clubID int64, points ...float64) {
		sh, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: first, Lastname: "Test", Birthday: birthday}, []int64{clubID})
		cat := domain.AgeCategoryFor(birthday, "2026-11-20")
		id, err := e.st.CreateRegistration(ctx, tid, store.NewRegistration{EventID: ev, ShooterID: sh, ClubID: &clubID, AgeCategory: cat,
			Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}}})
		if err != nil {
			t.Fatal(err)
		}
		if len(points) == 0 {
			return
		}
		reg, _ := e.st.GetRegistration(ctx, tid, id)
		entry, _ := e.st.GetEntry(ctx, tid, reg.Entries[0].ID)
		results := []domain.SeriesResult{}
		for _, p := range points {
			results = append(results, domain.SeriesResult{Points: p})
		}
		if _, err := e.st.SaveResults(ctx, tid, entry, results); err != nil {
			t.Fatal(err)
		}
		total, _ := domain.Totals(results)
		if eligible := domain.EligibleAwards(tiers[main], cat, total); len(eligible) > 0 {
			e.st.SetAward(ctx, tid, entry.ID, &eligible[0].ID)
		}
	}
	for i, s := range []float64{95, 94, 93, 92, 91, 90} {
		shoot(fmt.Sprintf("Home%d", i), "1990-01-01", home.ID, s, s, s)
	}
	shoot("Youth", "2006-01-01", small.ID, 88, 88, 88)
	shoot("Pending", "1990-01-01", small.ID)

	base := fmt.Sprintf("/api/tenants/%d/events/%d", tid, ev)
	code, b := e.do(staff, "GET", base+"/statistics", nil)
	expect(t, "statistics", code, 200, b)
	var st web.Statistics
	json.Unmarshal(b, &st)
	if st.Entries != 8 || st.Revenue != 150 || st.Completed != 7 || st.Registered != 1 || st.CompletionRate != 87.5 {
		t.Fatalf("statistics: %+v", st)
	}
	if len(st.Awards) != 1 || st.Awards[0].Count != 7 || st.AwardsValue != 105 || len(st.Categories) != 2 || st.Categories[0].Category != domain.U21 {
		t.Fatalf("awards and categories: %+v", st)
	}

	code, b = e.do(staff, "GET", base+"/club-competition", nil)
	expect(t, "club competition", code, 200, b)
	var cc web.ClubCompetitionView
	json.Unmarshal(b, &cc)
	// six scores of the first two series: 190 188 186 184 182 180, required 5 → (930 + 0.02×180)/5 = 186.72
	if len(cc.Ranked) != 1 || cc.Ranked[0].Name != "Home" || cc.Ranked[0].Result != 186.72 || cc.Ranked[0].Rank != 1 {
		t.Fatalf("ranked: %+v", cc.Ranked)
	}
	if len(cc.Unranked) != 1 || cc.Unranked[0].Name != "Small" || cc.Unranked[0].Participants != 1 {
		t.Fatalf("unranked: %+v", cc.Unranked)
	}

	evData, _ := e.st.GetEvent(ctx, tid, ev)
	e.st.UpdateEvent(ctx, tid, ev, store.EventInput{RangeIDs: []int64{evData.Ranges[0].ID}, Name: evData.Name, StartsOn: evData.StartsOn, Timezone: evData.Timezone,
		Settings: store.EventSettings{ExcludedClubID: &home.ID}})
	code, b = e.do(staff, "GET", base+"/club-competition", nil)
	cc = web.ClubCompetitionView{}
	json.Unmarshal(b, &cc)
	if !cc.Ranked[0].Excluded || cc.Ranked[0].Rank != 0 {
		t.Fatalf("excluded club: %+v", cc.Ranked)
	}

	get := func(path string) (int, string, string) {
		resp, err := staff.Get(e.srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, resp.Header.Get("Content-Type"), string(body)
	}
	for report, want := range map[string]string{
		"financial":         "Revenue;CHF 150.00",
		"revenue":           "Main;U21;1;CHF 10.00;CHF 10.00",
		"awards":            "Test;Home0;Elite;Home;Main;Wreath card;Wreath;CHF 15.00;285;0",
		"club_competition":  "Out of competition",
		"club_participants": "Home;1001;Elite;6",
	} {
		code, ct, body := get(base + "/reports/" + report)
		if code != 200 || ct != "application/pdf" || !strings.HasPrefix(body, "%PDF") {
			t.Errorf("%s pdf: %d %s", report, code, ct)
		}
		code, ct, body = get(base + "/reports/" + report + "?format=csv")
		if code != 200 || !strings.HasPrefix(ct, "text/csv") || !strings.Contains(body, want) {
			t.Errorf("%s csv: %d %s\n%s", report, code, ct, body)
		}
	}
	prizes := fmt.Sprintf("/api/tenants/%d/programmes/%d/special-prizes", tid, main)
	rules := domain.SpecialPrizeRules{Ranks: 2, MinCategory: 2, MinClub: 4}
	code, b = e.do(staff, "PUT", prizes, rules)
	expect(t, "staff sets special prizes", code, 403, b)
	code, b = e.do(e.login("boss"), "PUT", prizes, rules)
	expect(t, "set special prizes on a programme in use", code, 204, b)
	code, _, body := get(base + "/reports/special_prizes?format=csv")
	for _, want := range []string{"Special prizes", "\nElite;1;1;Test Home0;1990;Home;285\n", "\nElite;2;2;Test Home1;1990;Home;282\n", "\nU21;;;;Participants in the category: 1\n"} {
		if code != 200 || !strings.Contains(body, want) {
			t.Fatalf("special prizes: missing %q in\n%s", want, body)
		}
	}

	code, _, body = get(fmt.Sprintf("/api/tenants/%d/programmes/%d/ranking.csv", tid, main))
	if code != 200 || !strings.Contains(body, "Rank;Last name;First name;Year;Category;Club;S1;S2;S3;Total;Mouches") || !strings.Contains(body, "\n1;Test;Home0;1990;Elite;Home;95;95;95;285;0\n") {
		t.Fatalf("ranking csv: %d\n%s", code, body)
	}

	evData, _ = e.st.GetEvent(ctx, tid, ev)
	e.st.UpdateEvent(ctx, tid, ev, store.EventInput{RangeIDs: []int64{evData.Ranges[0].ID}, Name: evData.Name, StartsOn: evData.StartsOn, Timezone: evData.Timezone,
		Locale: "de", Settings: evData.Settings})
	resp, err := staff.Get(e.srv.URL + base + "/reports/financial?format=csv")
	if err != nil {
		t.Fatal(err)
	}
	de, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if !strings.Contains(string(de), "Finanzbericht") || !strings.Contains(string(de), "Einnahmen;CHF 150.00") ||
		!strings.Contains(resp.Header.Get("Content-Disposition"), `filename="financial_`) {
		t.Fatalf("german report: %s\n%s", resp.Header.Get("Content-Disposition"), de)
	}

	code, b = e.do(outsider, "GET", base+"/statistics", nil)
	expect(t, "outsider statistics", code, 404, b)
	code, b = e.do(outsider, "GET", fmt.Sprintf("/api/tenants/%d/events/%d/reports/financial", other, ev), nil)
	expect(t, "outsider report via own tenant", code, 404, b)
}
