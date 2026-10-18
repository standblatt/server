package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

// pdfWords returns the text of a PDF via pdftotext, or "" (and skips the check) when it is not installed.
func pdfWords(t *testing.T, pdf []byte) string {
	t.Helper()
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return ""
	}
	f := filepath.Join(t.TempDir(), "doc.pdf")
	os.WriteFile(f, pdf, 0o600)
	out, err := exec.Command("pdftotext", f, "-").Output()
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// TestEventForms checks that an event's wreath card and control strip print through its chosen print forms and fall
// back to the built-in ones.
func TestEventForms(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	e.user("boss", tid, store.RoleAdmin)
	c := e.login("boss")
	ev := e.event(tid, "Match")
	main, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Main", Series: []domain.Series{{Shots: 10, MaxPoints: 100}}, Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 10}},
		Awards: []domain.AwardTier{{Name: "Wreath", Kind: domain.AwardCard, Value: 15, Limits: []domain.AwardLimit{{AgeCategory: domain.Elite, Points: 90}}}},
	})
	sh, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Anna", Lastname: "Winner", Birthday: "1990-01-01"}, nil)
	base := fmt.Sprintf("/api/tenants/%d", tid)
	evPath := fmt.Sprintf("%s/events/%d", base, ev)
	code, b := e.do(c, "POST", evPath+"/registrations", web.RegistrationInput{ShooterID: sh, Programmes: []domain.ProgrammeChoice{{ProgrammeID: main}}})
	expect(t, "register", code, 201, b)
	var reg store.Registration
	json.Unmarshal(b, &reg)
	code, b = e.do(c, "PUT", fmt.Sprintf("%s/entries/%d/results", base, reg.Entries[0].ID), map[string]any{"series": []map[string]any{{"points": 95}}})
	expect(t, "results", code, 200, b)
	wreath := fmt.Sprintf("%s/entries/%d/wreath-card", base, reg.Entries[0].ID)
	strip := fmt.Sprintf("%s/registrations/%d/documents/control_strip", base, reg.ID)

	code, b = e.do(c, "GET", wreath, nil)
	expect(t, "built-in wreath card", code, 200, b)
	if txt := pdfWords(t, b); txt != "" && !strings.Contains(txt, "Winner") {
		t.Fatalf("built-in wreath card: %s", txt)
	}
	code, b = e.do(c, "GET", strip, nil)
	expect(t, "built-in control strip", code, 200, b)
	if txt := pdfWords(t, b); txt != "" && !strings.Contains(txt, fmt.Sprintf("ID %d", reg.ID)) {
		t.Fatalf("built-in control strip: %s", txt)
	}

	form := func(name, source, text string) int64 {
		code, b := e.do(c, "POST", base+"/layouts", map[string]any{"name": name, "width": 105, "height": 148, "source": source,
			"elements": []map[string]any{{"type": "text", "x": 10, "y": 10, "value": text}}})
		expect(t, "form "+name, code, 201, b)
		var l store.Layout
		json.Unmarshal(b, &l)
		return l.ID
	}
	card := form("Own card", "results", "CARD {{shooter.lastname}} {{award.value}}")
	label := form("Own strip", "registration", "STRIP {{passes}}")
	update := func(set map[string]any) (int, []byte) {
		return e.do(c, "PUT", evPath, map[string]any{"name": "Match", "startsOn": "2026-11-20", "settings": set})
	}
	code, b = update(map[string]any{"wreathCardFormId": label})
	expect(t, "registration form as wreath card", code, 422, b)
	code, b = update(map[string]any{"wreathCardFormId": card, "controlStripFormId": label, "printAwardValue": true})
	expect(t, "choose the forms", code, 204, b)

	code, b = e.do(c, "GET", wreath, nil)
	expect(t, "own wreath card", code, 200, b)
	if txt := pdfWords(t, b); txt != "" && !strings.Contains(txt, "CARD Winner CHF 15.00") {
		t.Fatalf("own wreath card: %s", txt)
	}
	code, b = e.do(c, "GET", strip, nil)
	expect(t, "own control strip", code, 200, b)
	if txt := pdfWords(t, b); txt != "" && !strings.Contains(txt, "STRIP P") {
		t.Fatalf("own control strip: %s", txt)
	}

	code, b = e.do(c, "DELETE", fmt.Sprintf("%s/layouts/%d", base, card), nil)
	expect(t, "delete the card", code, 204, b)
	code, b = e.do(c, "GET", wreath, nil)
	expect(t, "back to the built-in card", code, 200, b)
	if txt := pdfWords(t, b); txt != "" && strings.Contains(txt, "CARD") {
		t.Fatalf("deleted form still prints: %s", txt)
	}

	code, b = e.do(c, "GET", base+"/layouts/fields?source=registration", nil)
	expect(t, "registration fields", code, 200, b)
	if !strings.Contains(string(b), "registration.code") || strings.Contains(string(b), `"total"`) {
		t.Fatalf("registration fields: %s", b)
	}
	code, b = e.do(c, "GET", base+"/layouts/builtin", nil)
	expect(t, "built-in forms", code, 200, b)
	code, b = e.do(c, "PUT", fmt.Sprintf("%s/programmes/%d/form", base, main), map[string]any{"layoutId": label})
	expect(t, "registration form on a programme", code, 422, b)
}
