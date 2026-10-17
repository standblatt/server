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

func TestAwardsHandout(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	e.user("boss", tid, store.RoleAdmin)
	c := e.login("boss")
	ev := e.event(tid, "Match")
	main, err := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Main", Series: []domain.Series{{Shots: 10, MaxPoints: 100}},
		Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 10}},
		Awards: []domain.AwardTier{{Name: "Payout 20", Kind: domain.AwardCash, Value: 20, Limits: []domain.AwardLimit{{AgeCategory: domain.Elite, Points: 90}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	id, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Anna", Lastname: "Winner", Birthday: "1990-01-01"}, nil)
	code, b := e.do(c, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/registrations", tid, ev), web.RegistrationInput{ShooterID: id, Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}}})
	expect(t, "register", code, 201, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)
	code, b = e.do(c, "PUT", fmt.Sprintf("/api/tenants/%d/entries/%d/results", tid, reg.Entries[0].ID), map[string]any{"series": []map[string]any{{"points": 95}}})
	expect(t, "results", code, 200, b)
	var sheet web.ResultSheet
	json.Unmarshal(b, &sheet)
	if sheet.Entries[0].Award == nil || sheet.Entries[0].Award.Kind != domain.AwardCash {
		t.Fatalf("cash award not given: %s", b)
	}

	report := fmt.Sprintf("/api/tenants/%d/events/%d/reports/awards_open?format=csv", tid, ev)
	code, b = e.do(c, "GET", report, nil)
	expect(t, "open report", code, 200, b)
	if !strings.Contains(string(b), "Winner") || !strings.Contains(string(b), "20.00") {
		t.Fatalf("open report misses the payout: %s", b)
	}
	code, b = e.do(c, "POST", fmt.Sprintf("/api/tenants/%d/registrations/%d/awards-issued", tid, reg.ID), nil)
	expect(t, "hand out", code, 204, b)
	code, b = e.do(c, "GET", fmt.Sprintf("/api/tenants/%d/registrations/%d", tid, reg.ID), nil)
	expect(t, "get", code, 200, b)
	json.Unmarshal(b, &reg)
	if reg.AwardsIssuedAt == nil {
		t.Fatalf("handout not recorded: %s", b)
	}
	code, b = e.do(c, "GET", report, nil)
	expect(t, "open report after", code, 200, b)
	if strings.Contains(string(b), "Winner") {
		t.Fatalf("handed-out award still listed: %s", b)
	}
	code, b = e.do(c, "DELETE", fmt.Sprintf("/api/tenants/%d/registrations/%d/awards-issued", tid, reg.ID), nil)
	expect(t, "take back", code, 204, b)
	code, b = e.do(c, "GET", fmt.Sprintf("/api/tenants/%d/registrations/%d", tid, reg.ID), nil)
	var after store.Registration
	json.Unmarshal(b, &after)
	if after.AwardsIssuedAt != nil {
		t.Fatalf("handout still recorded: %s", b)
	}
	code, b = e.do(c, "GET", fmt.Sprintf("/api/tenants/%d/events/%d/reports/awards?format=csv", tid, ev), nil)
	expect(t, "awards report", code, 200, b)
	if !strings.Contains(string(b), "Cash prize") {
		t.Fatalf("awards report without the cash kind: %s", b)
	}
}
