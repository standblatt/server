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

func TestArchivedEvent(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	admin, staff := e.login("boss"), e.login("helper")
	ev := e.event(tid, "Indoor")
	main, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Main", Series: []domain.Series{{Shots: 10, MaxPoints: 100}}, Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 10}},
		Awards: []domain.AwardTier{{Name: "Badge", Kind: domain.AwardNatural, Limits: []domain.AwardLimit{{AgeCategory: domain.Elite, Points: 50}}}},
	})
	anna, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Anna", Lastname: "Adult", Birthday: "1990-01-01"}, nil)
	bert, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Bert", Lastname: "Best", Birthday: "1991-01-01"}, nil)

	base := fmt.Sprintf("/api/tenants/%d", tid)
	evPath := fmt.Sprintf("%s/events/%d", base, ev)
	code, b := e.do(staff, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: anna, Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}}})
	expect(t, "register", code, 201, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)
	entry := fmt.Sprintf("%s/entries/%d", base, reg.Entries[0].ID)
	code, b = e.do(staff, "PUT", entry+"/results", map[string]any{"series": []map[string]any{{"points": 80}}})
	expect(t, "results", code, 200, b)

	archive := func(archived bool) (int, []byte) {
		return e.do(admin, "PUT", evPath+"/archived", map[string]bool{"archived": archived})
	}
	code, b = e.do(staff, "PUT", evPath+"/archived", map[string]bool{"archived": true})
	expect(t, "staff archives", code, 403, b)
	code, b = archive(true)
	expect(t, "archive", code, 204, b)
	code, b = archive(true)
	expect(t, "archive twice", code, 204, b)

	code, b = e.do(staff, "GET", evPath, nil)
	var got store.Event
	json.Unmarshal(b, &got)
	if code != 200 || got.ArchivedAt == nil {
		t.Fatalf("archived event: %d %s", code, b)
	}

	refused := map[string]func() (int, []byte){
		"results": func() (int, []byte) {
			return e.do(staff, "PUT", entry+"/results", map[string]any{"series": []map[string]any{{"points": 90}}})
		},
		"award": func() (int, []byte) { return e.do(staff, "PUT", entry+"/award", map[string]any{}) },
		"register": func() (int, []byte) {
			return e.do(staff, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: bert, Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}}})
		},
		"delete reg": func() (int, []byte) {
			return e.do(staff, "DELETE", fmt.Sprintf("%s/registrations/%d", base, reg.ID), nil)
		},
		"programme": func() (int, []byte) {
			return e.do(admin, "PUT", fmt.Sprintf("%s/programmes/%d/active", base, main), map[string]bool{"active": false})
		},
		"new program": func() (int, []byte) {
			return e.do(admin, "POST", evPath+"/programmes", domain.ProgrammeDefinition{Name: "New", Series: []domain.Series{{Shots: 1, MaxPoints: 10}}, Prices: []domain.Price{}, Awards: []domain.AwardTier{}})
		},
		"update event": func() (int, []byte) {
			return e.do(admin, "PUT", evPath, map[string]any{"rangeIds": []int64{got.Ranges[0].ID}, "name": "Renamed", "startsOn": got.StartsOn})
		},
		"delete event": func() (int, []byte) { return e.do(admin, "DELETE", evPath, nil) },
	}
	for what, call := range refused {
		code, b := call()
		expect(t, what+" on an archived event", code, 409, b)
	}

	for _, path := range []string{evPath + "/rankings", evPath + "/statistics", evPath + "/reports/financial", fmt.Sprintf("%s/registrations/%d/results", base, reg.ID)} {
		code, b := e.do(staff, "GET", path, nil)
		expect(t, "read "+path, code, 200, b)
	}

	code, b = archive(false)
	expect(t, "unarchive", code, 204, b)
	code, b = e.do(staff, "PUT", entry+"/results", map[string]any{"series": []map[string]any{{"points": 90}}})
	expect(t, "results after unarchive", code, 200, b)

	second := e.event(tid, "Second")
	other, _ := e.st.CreateProgramme(ctx, tid, second, domain.ProgrammeDefinition{
		Name: "Other", Series: []domain.Series{{Shots: 10, MaxPoints: 100}}, Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 10}}, Awards: []domain.AwardTier{},
	})
	code, b = e.do(staff, "POST", fmt.Sprintf("%s/events/%d/registrations", base, second), web.RegistrationInput{ShooterID: anna, Programmes: []domain.ProgrammeChoice{{ProgrammeID: other}}})
	expect(t, "register for the second event", code, 201, b)
	code, b = e.do(admin, "DELETE", fmt.Sprintf("%s/events/%d", base, second), nil)
	expect(t, "delete an event with registrations", code, 409, b)

	archive(true)
	if err := e.st.DeleteTenant(ctx, tid); err != nil {
		t.Fatalf("deleting a tenant with an archived event: %v", err)
	}
}
