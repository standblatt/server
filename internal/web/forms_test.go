package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func TestPrintForms(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	tid := e.tenant("Club")
	other := e.tenant("Other")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	e.user("stranger", other, store.RoleAdmin)
	admin, staff := e.login("boss"), e.login("helper")
	ev := e.event(tid, "Annual")
	p, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{Name: "Obligatory", Awards: []domain.AwardTier{},
		Series: []domain.Series{{Shots: 3, MaxPoints: 15, TrackShots: true}}, Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 0}}})
	plain, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{Name: "Plain", Awards: []domain.AwardTier{},
		Series: []domain.Series{{Shots: 3, MaxPoints: 15}}, Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 0}}})
	register := func(first string, lane int, programme int64) store.Registration {
		id, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: first, Lastname: "Form", Birthday: "1990-01-01"}, nil)
		code, b := e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/registrations", tid, ev),
			web.RegistrationInput{ShooterID: id, Lane: &lane, Programmes: []domain.ProgrammeChoice{{ProgrammeID: programme}}})
		expect(t, "register "+first, code, 201, b)
		var reg store.Registration
		json.Unmarshal(b, &reg)
		code, b = e.do(staff, "PUT", fmt.Sprintf("/api/tenants/%d/entries/%d/results", tid, reg.Entries[0].ID),
			map[string]any{"series": []map[string]any{{"points": 0, "shots": []float64{5, 4, 5}}}})
		if programme == plain {
			code, b = e.do(staff, "PUT", fmt.Sprintf("/api/tenants/%d/entries/%d/results", tid, reg.Entries[0].ID), map[string]any{"series": []map[string]any{{"points": 12}}})
		}
		expect(t, "results "+first, code, 200, b)
		return reg
	}
	anna := register("Anna", 2, p)
	bert := register("Bert", 7, p)
	cleo := register("Cleo", 3, plain)

	base := fmt.Sprintf("/api/tenants/%d", tid)
	form := map[string]any{"name": "Army sheet", "width": 148, "height": 210, "offsetX": 1, "elements": []map[string]any{
		{"type": "text", "x": 20, "y": 20, "value": "{{shooter.name}}"},
		{"type": "text", "x": 20, "y": 40, "w": 8, "value": "{{s1.shot#}}", "repeat": 3, "dx": 9},
	}}
	code, b := e.do(staff, "POST", base+"/layouts", form)
	expect(t, "staff creates a form", code, 403, b)
	code, b = e.do(admin, "POST", base+"/layouts", form)
	expect(t, "create a form", code, 201, b)
	var l store.Layout
	json.Unmarshal(b, &l)
	if l.ID == 0 || len(l.Elements) != 2 || l.Elements[1].Repeat != 3 || l.HasBackground {
		t.Fatalf("form: %s", b)
	}
	code, b = e.do(admin, "POST", base+"/layouts", form)
	expect(t, "same name twice", code, 409, b)
	code, b = e.do(staff, "GET", base+"/layouts/fields", nil)
	if code != 200 || !strings.Contains(string(b), `"s1.shot1"`) {
		t.Fatalf("fields: %d %s", code, b)
	}

	pdf := func(c *http.Client, path string) int {
		t.Helper()
		resp, err := c.Get(e.srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 200 && (resp.Header.Get("Content-Type") != "application/pdf" || !bytes.HasPrefix(body, []byte("%PDF"))) {
			t.Fatalf("%s is no PDF", path)
		}
		return resp.StatusCode
	}
	entryForm := func(r store.Registration) string { return fmt.Sprintf("%s/entries/%d/form", base, r.Entries[0].ID) }
	if code := pdf(staff, entryForm(anna)); code != 409 {
		t.Fatalf("form of a programme without one: %d", code)
	}
	code, b = e.do(admin, "PUT", fmt.Sprintf("%s/programmes/%d/form", base, p), map[string]any{"layoutId": l.ID})
	expect(t, "choose the form", code, 204, b)
	if code := pdf(staff, entryForm(anna)); code != 200 {
		t.Fatalf("entry form: %d", code)
	}
	if code := pdf(staff, entryForm(cleo)); code != 409 {
		t.Fatalf("form of a plain programme: %d", code)
	}
	code, b = e.do(staff, "POST", base+"/print", web.PrintRequest{Document: "form", EntryID: anna.Entries[0].ID})
	if code != 200 || !strings.Contains(string(b), `"queued":false`) {
		t.Fatalf("print a form without a printer: %d %s", code, b)
	}
	var sheet web.ResultSheet
	_, b = e.do(staff, "GET", fmt.Sprintf("%s/registrations/%d/results", base, anna.ID), nil)
	json.Unmarshal(b, &sheet)
	if !sheet.Entries[0].Form {
		t.Fatalf("the result sheet should offer the form: %s", b)
	}

	var img bytes.Buffer
	png.Encode(&img, image.NewGray(image.Rect(0, 0, 20, 28)))
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, _ := w.CreateFormFile("file", "scan.png")
	part.Write(img.Bytes())
	w.Close()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s/layouts/%d/background", e.srv.URL, base, l.ID), &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := admin.Do(req)
	if err != nil || resp.StatusCode != 204 {
		t.Fatalf("upload the scan: %v %v", resp.StatusCode, err)
	}
	resp.Body.Close()
	if code := pdf(admin, fmt.Sprintf("%s/layouts/%d/preview?entryId=%d", base, l.ID, anna.Entries[0].ID)); code != 200 {
		t.Fatalf("preview: %d", code)
	}

	scorer := func(role string) string {
		ranges, _ := e.st.ListRanges(ctx, tid)
		c := store.HelperCode{Name: role, Role: store.HelperRole(role)}
		if role == "scorer" {
			c.RangeID, c.Lanes = &ranges[0].ID, []int{1, 2, 3}
		}
		_, b := e.do(admin, "POST", fmt.Sprintf("%s/events/%d/helper-codes", base, ev), map[string]any{"name": c.Name, "role": role, "rangeId": c.RangeID, "lanes": "1-3"})
		var out struct {
			Secret string `json:"secret"`
		}
		json.Unmarshal(b, &out)
		return out.Secret
	}
	sc, ammo := scorer("scorer"), scorer("ammunition")
	code, _ = e.helper("GET", entryForm(anna), sc, nil)
	expect(t, "scorer opens the form of a lane of theirs", code, 200, nil)
	code, _ = e.helper("GET", entryForm(bert), sc, nil)
	if code != 403 && code != 404 {
		t.Fatalf("scorer opens the form of another lane: %d", code)
	}
	code, b = e.helper("POST", base+"/print", sc, web.PrintRequest{Document: "form", EntryID: anna.Entries[0].ID})
	expect(t, "scorer prints a form", code, 200, b)
	code, b = e.helper("POST", base+"/print", sc, web.PrintRequest{Document: "form", EntryID: bert.Entries[0].ID})
	if code != 403 && code != 404 {
		t.Fatalf("scorer prints the form of another lane: %d %s", code, b)
	}
	code, b = e.helper("POST", base+"/print", ammo, web.PrintRequest{Document: "form", EntryID: anna.Entries[0].ID})
	expect(t, "ammunition desk prints a form", code, 403, b)

	stranger := e.login("stranger")
	code, b = e.do(stranger, "GET", fmt.Sprintf("/api/tenants/%d/layouts/%d", other, l.ID), nil)
	expect(t, "form of another tenant", code, 404, b)
	foreign, _ := e.st.CreateLayout(ctx, other, store.LayoutInput{Name: "Foreign", Layout: domain.Layout{Width: 100, Height: 100, Elements: []domain.Element{}}})
	code, b = e.do(admin, "PUT", fmt.Sprintf("%s/programmes/%d/form", base, plain), map[string]any{"layoutId": foreign})
	expect(t, "form of another tenant on a programme", code, 404, b)

	code, b = e.do(admin, "DELETE", fmt.Sprintf("%s/layouts/%d", base, l.ID), nil)
	expect(t, "delete the form", code, 204, b)
	if prog, _ := e.st.GetProgramme(ctx, tid, p); prog.LayoutID != nil {
		t.Fatalf("the programme keeps a deleted form: %+v", prog.LayoutID)
	}
}
