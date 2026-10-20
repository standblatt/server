package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

func (e *env) upload(c *http.Client, path, content string, clubs ...string) (int, []byte) {
	e.t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, _ := w.CreateFormFile("file", "shooters.csv")
	part.Write([]byte(content))
	for _, m := range clubs {
		w.WriteField("clubs", m)
	}
	w.Close()
	req, _ := http.NewRequest("POST", e.srv.URL+path, &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := c.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b
}

func TestShootersAndClubs(t *testing.T) {
	e := newEnv(t)
	tid, other := e.tenant("Club"), e.tenant("Other")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	e.user("outsider", other, store.RoleAdmin)
	admin, staff, outsider := e.login("boss"), e.login("helper"), e.login("outsider")
	base := fmt.Sprintf("/api/tenants/%d", tid)

	code, b := e.do(staff, "POST", base+"/shooters", map[string]any{
		"licenseNumber": "13756", "firstname": " Sandra ", "lastname": "Schwitter", "birthday": "2.6.1969",
		"clubNames": []string{"Sargans Pistol"},
	})
	expect(t, "staff creates shooter", code, 201, b)
	var sh store.Shooter
	json.Unmarshal(b, &sh)
	if sh.LicenseNumber != "013756" || sh.Firstname != "Sandra" || sh.Birthday != "1969-06-02" || len(sh.Clubs) != 1 {
		t.Fatalf("created shooter: %+v", sh)
	}
	shooterPath := fmt.Sprintf("%s/shooters/%d", base, sh.ID)

	code, b = e.do(staff, "POST", base+"/shooters", map[string]any{"licenseNumber": "013756", "firstname": "x", "lastname": "y"})
	expect(t, "duplicate license", code, 409, b)
	code, b = e.do(staff, "POST", base+"/shooters", map[string]any{"firstname": "x", "lastname": "y", "birthday": "31.02.2000"})
	expect(t, "invalid birthday", code, 422, b)

	code, b = e.do(staff, "GET", base+"/shooters?q=sandra+schw", nil)
	expect(t, "search", code, 200, b)
	var hits []store.Shooter
	json.Unmarshal(b, &hits)
	if len(hits) != 1 || hits[0].ID != sh.ID {
		t.Fatalf("search hits: %+v", hits)
	}

	code, b = e.do(staff, "PUT", shooterPath, map[string]any{"licenseNumber": "013756", "firstname": "Sandra", "lastname": "Schwitter-Muster", "city": "Sargans"})
	expect(t, "update shooter", code, 204, b)
	code, b = e.do(staff, "POST", shooterPath+"/clubs", map[string]any{"clubName": "Second Club"})
	expect(t, "add club by name", code, 204, b)
	code, b = e.do(staff, "GET", shooterPath, nil)
	expect(t, "get shooter", code, 200, b)
	json.Unmarshal(b, &sh)
	if sh.Lastname != "Schwitter-Muster" || sh.Birthday != "" || len(sh.Clubs) != 2 {
		t.Fatalf("updated shooter: %+v", sh)
	}

	code, b = e.do(staff, "GET", base+"/clubs", nil)
	expect(t, "list clubs", code, 200, b)
	var clubs []store.Club
	json.Unmarshal(b, &clubs)
	if len(clubs) != 2 {
		t.Fatalf("clubs: %+v", clubs)
	}
	code, b = e.do(staff, "DELETE", fmt.Sprintf("%s/clubs/%d", shooterPath, clubs[0].ID), nil)
	expect(t, "remove club link", code, 204, b)
	code, b = e.do(staff, "PUT", fmt.Sprintf("%s/clubs/%d", base, clubs[0].ID), map[string]any{"name": "Renamed", "ssvNumber": "123"})
	expect(t, "staff renames club", code, 403, b)
	code, b = e.do(admin, "PUT", fmt.Sprintf("%s/clubs/%d", base, clubs[0].ID), map[string]any{"name": "Renamed", "ssvNumber": "123"})
	expect(t, "admin renames club", code, 204, b)

	code, b = e.do(outsider, "GET", shooterPath, nil)
	expect(t, "outsider reads shooter", code, 404, b)
	code, b = e.do(outsider, "POST", fmt.Sprintf("/api/tenants/%d/shooters/%d/clubs", other, sh.ID), map[string]any{"clubName": "x"})
	expect(t, "outsider links shooter via own tenant", code, 404, b)

	code, b = e.do(staff, "DELETE", shooterPath, nil)
	expect(t, "staff deletes shooter", code, 403, b)
	code, b = e.do(admin, "DELETE", shooterPath, nil)
	expect(t, "admin deletes shooter", code, 204, b)

	clubPath := fmt.Sprintf("%s/clubs/%d", base, clubs[0].ID)
	code, b = e.do(staff, "DELETE", clubPath, nil)
	expect(t, "staff deletes club", code, 403, b)
	code, b = e.do(admin, "DELETE", clubPath, nil)
	expect(t, "admin deletes club", code, 204, b)
	code, b = e.do(admin, "DELETE", clubPath, nil)
	expect(t, "delete club again", code, 404, b)
}

func TestCrossOriginUpload(t *testing.T) {
	e := newEnv(t)
	tid := e.tenant("Club")
	e.user("boss", tid, store.RoleAdmin)
	admin := e.login("boss")
	for _, c := range []struct{ header, value string }{
		{"Sec-Fetch-Site", "same-site"},
		{"Sec-Fetch-Site", "cross-site"},
		{"Origin", "https://www.example.org"},
	} {
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		part, _ := w.CreateFormFile("file", "shooters.csv")
		part.Write([]byte("license;firstname;lastname\n13756;Sandra;Schwitter\n"))
		w.Close()
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/tenants/%d/shooters/import", e.srv.URL, tid), &buf)
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.Header.Set(c.header, c.value)
		resp, err := admin.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("%s: %s answered %d, want 403", c.header, c.value, resp.StatusCode)
		}
	}
	if n, err := e.st.SearchShooters(context.Background(), tid, "", 0, 10, 0); err != nil || len(n) != 0 {
		t.Fatalf("cross-origin uploads imported shooters: %v, %v", n, err)
	}
}

func TestShooterImport(t *testing.T) {
	e := newEnv(t)
	tid := e.tenant("Club")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	admin, staff := e.login("boss"), e.login("helper")
	path := fmt.Sprintf("/api/tenants/%d/shooters/import", tid)

	good := "license;firstname;lastname;birthday;club\n13756;Sandra;Schwitter;02.06.1969;Sargans\n;Guest;Person;;Sargans\n"
	code, b := e.upload(staff, path, good)
	expect(t, "staff imports", code, 403, b)

	code, b = e.upload(admin, path, good+"1;;Nofirst;;\n")
	expect(t, "import with errors", code, 200, b)
	var rep web.ImportReport
	json.Unmarshal(b, &rep)
	if len(rep.Errors) != 1 || rep.Errors[0].Line != 4 || rep.Created != 0 {
		t.Fatalf("report with errors: %+v", rep)
	}

	code, b = e.do(admin, "POST", fmt.Sprintf("/api/tenants/%d/clubs", tid), map[string]any{"name": "Sargans-Wartau"})
	expect(t, "create a similar club", code, 201, b)
	code, b = e.upload(admin, path, good)
	expect(t, "import with an unknown club", code, 200, b)
	rep = web.ImportReport{}
	json.Unmarshal(b, &rep)
	if rep.Created != 0 || len(rep.UnknownClubs) != 1 || rep.UnknownClubs[0].Name != "Sargans" || len(rep.UnknownClubs[0].Similar) != 1 {
		t.Fatalf("report with an unknown club: %s", b)
	}

	foreign := e.tenant("Other")
	fc, _ := e.st.CreateClub(context.Background(), foreign, store.ClubInput{Name: "Foreign"})
	code, b = e.upload(admin, path, good, fmt.Sprintf(`{"Sargans": %d}`, fc.ID))
	expect(t, "map to a club of another tenant", code, 422, b)

	code, b = e.upload(admin, path, good, `{"Sargans": 0}`)
	expect(t, "import", code, 200, b)
	rep = web.ImportReport{}
	json.Unmarshal(b, &rep)
	if rep.Created != 2 || rep.ClubsCreated != 1 || rep.Links != 2 || len(rep.Errors) != 0 {
		t.Fatalf("report: %+v", rep)
	}

	code, b = e.upload(admin, path, "firstname;lastname;club\nNew;Member; SARGANS\nOther;Member;Sargans-Wartau\n")
	expect(t, "import into existing clubs", code, 200, b)
	rep = web.ImportReport{}
	json.Unmarshal(b, &rep)
	if rep.Created != 2 || rep.ClubsCreated != 0 || len(rep.UnknownClubs) != 0 {
		t.Fatalf("clubs matched regardless of case: %s", b)
	}

	code, b = e.upload(admin, path, "city\nBern\n")
	expect(t, "bad header", code, 422, b)

	t.Setenv("TMPDIR", "/nonexistent")
	big := "firstname;lastname\n" + strings.Repeat("Big;Import\n", 2000)
	code, b = e.upload(admin, path, big)
	expect(t, "import larger than 8 KB without a writable temp directory", code, 200, b)
}

func TestClubMatchAndMerge(t *testing.T) {
	e := newEnv(t)
	tid, other := e.tenant("Club"), e.tenant("Other")
	e.user("boss", tid, store.RoleAdmin)
	e.user("helper", tid, store.RoleStaff)
	e.user("stranger", other, store.RoleAdmin)
	admin, staff := e.login("boss"), e.login("helper")
	ctx := context.Background()
	hill, _ := e.st.CreateClub(ctx, tid, store.ClubInput{Name: "SG Hilltown"})
	typo, _ := e.st.CreateClub(ctx, tid, store.ClubInput{Name: "SG Hiltown"})
	base := fmt.Sprintf("/api/tenants/%d/clubs", tid)

	code, b := e.do(staff, "GET", base+"/similar?name=sg%20hilltown", nil)
	expect(t, "staff matches a club", code, 200, b)
	var m web.ClubMatch
	json.Unmarshal(b, &m)
	if m.Match == nil || m.Match.ID != hill.ID || len(m.Similar) != 1 || m.Similar[0].ID != typo.ID {
		t.Fatalf("match: %s", b)
	}
	code, b = e.do(admin, "POST", base, map[string]any{"name": "sg HILLTOWN"})
	expect(t, "create a club differing in case", code, 409, b)

	merge := fmt.Sprintf("%s/%d/merge", base, typo.ID)
	code, b = e.do(staff, "POST", merge, map[string]any{"intoId": hill.ID})
	expect(t, "staff merges", code, 403, b)
	code, b = e.do(e.login("stranger"), "POST", fmt.Sprintf("/api/tenants/%d/clubs/%d/merge", other, typo.ID), map[string]any{"intoId": hill.ID})
	expect(t, "other tenant merges", code, 404, b)
	code, b = e.do(admin, "POST", merge, map[string]any{"intoId": hill.ID})
	expect(t, "merge", code, 204, b)
	if clubs, _ := e.st.ListClubs(ctx, tid); len(clubs) != 1 || clubs[0].ID != hill.ID {
		t.Fatalf("clubs after merge: %+v", clubs)
	}
}
