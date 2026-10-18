package web_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/mail"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/web"
)

type mailbox chan mail.Message

func (m mailbox) next(t *testing.T) mail.Message {
	t.Helper()
	select {
	case msg := <-m:
		return msg
	case <-time.After(5 * time.Second):
		t.Fatal("no mail")
		return mail.Message{}
	}
}

var tokenRe = regexp.MustCompile(`token=([A-Za-z0-9_-]+)`)

func (m mailbox) token(t *testing.T) string {
	t.Helper()
	msg := m.next(t)
	match := tokenRe.FindStringSubmatch(msg.Body)
	if match == nil {
		t.Fatalf("no token in %q", msg.Body)
	}
	return match[1]
}

func personalEnv(t *testing.T, configure func(*web.Server)) (*env, mailbox) {
	box := make(mailbox, 10)
	e := newEnvWith(t, func(s *web.Server) {
		s.Signup = true
		s.SendMail = func(_ context.Context, m mail.Message) error { box <- m; return nil }
		if configure != nil {
			configure(s)
		}
	})
	return e, box
}

func (e *env) signup(email string) *http.Client {
	e.t.Helper()
	c := e.client()
	code, b := e.do(c, "POST", "/api/auth/signup", map[string]string{
		"email": email, "password": "password-" + email, "firstname": "Example", "lastname": "Felix", "locale": "de",
	})
	expect(e.t, "signup "+email, code, 201, b)
	return e.withTokens(c, email)
}

func overview(t *testing.T, e *env, c *http.Client) web.MyOverview {
	t.Helper()
	code, b := e.do(c, "GET", "/api/my", nil)
	expect(t, "overview", code, 200, b)
	var o web.MyOverview
	json.Unmarshal(b, &o)
	return o
}

func TestSignupVerifyReset(t *testing.T) {
	e, _ := personalEnv(t, func(s *web.Server) { s.Signup = false })
	code, b := e.do(e.client(), "POST", "/api/auth/signup", map[string]string{"email": "a@example.org", "password": "long-enough-1", "firstname": "A", "lastname": "B"})
	expect(t, "closed signup", code, 403, b)

	e, box := personalEnv(t, nil)
	code, b = e.do(e.client(), "GET", "/api/auth/options", nil)
	expect(t, "options", code, 200, b)
	if !strings.Contains(string(b), `"signup":true`) {
		t.Fatalf("options: %s", b)
	}
	c := e.signup("Felix@Example.org")
	o := overview(t, e, c)
	if o.Email != "felix@example.org" || o.EmailVerified || o.Profile.Firstname != "Example" || len(o.Disciplines) != 6 {
		t.Fatalf("overview after signup: %+v", o)
	}
	code, b = e.do(e.client(), "POST", "/api/auth/signup", map[string]string{"email": "felix@example.org", "password": "long-enough-1", "firstname": "A", "lastname": "B"})
	expect(t, "duplicate signup", code, 409, b)

	code, b = e.do(c, "POST", "/api/my/email/verify", nil)
	expect(t, "send the link again", code, 204, b)
	box.token(t)
	verify := box.token(t)
	code, b = e.do(c, "POST", "/api/auth/verify", map[string]string{"token": verify})
	expect(t, "verify", code, 204, b)
	code, b = e.do(c, "POST", "/api/my/email/verify", nil)
	expect(t, "send the link of a confirmed address", code, 409, b)
	code, b = e.do(c, "POST", "/api/auth/verify", map[string]string{"token": verify})
	expect(t, "verify twice", code, 410, b)
	if !overview(t, e, c).EmailVerified {
		t.Fatal("email not verified")
	}

	code, b = e.do(e.client(), "POST", "/api/auth/forgot", map[string]string{"email": "nobody@example.org"})
	expect(t, "forgot unknown", code, 204, b)
	code, b = e.do(e.client(), "POST", "/api/auth/forgot", map[string]string{"email": "felix@example.org"})
	expect(t, "forgot", code, 204, b)
	reset := box.token(t)
	code, b = e.do(e.client(), "POST", "/api/auth/reset", map[string]string{"token": reset, "password": "short"})
	expect(t, "weak password", code, 422, b)
	code, b = e.do(e.client(), "POST", "/api/auth/reset", map[string]string{"token": reset, "password": "a-new-password-1"})
	expect(t, "reset", code, 204, b)
	code, b = e.do(c, "GET", "/api/me", nil)
	expect(t, "old session after reset", code, 401, b)
	n := e.client()
	code, b = e.do(n, "POST", "/api/auth/login", map[string]string{"name": "FELIX@example.org", "password": "a-new-password-1"})
	expect(t, "login with email", code, 200, b)
	code, b = e.do(n, "PUT", "/api/my/email", map[string]string{"email": "Felix.New@example.org", "password": "a-new-password-1"})
	expect(t, "change email", code, 204, b)
	box.token(t)
	if o := overview(t, e, n); o.Email != "felix.new@example.org" || o.EmailVerified {
		t.Fatalf("overview after the change: %+v", o)
	}
}

func TestTrainings(t *testing.T) {
	e, _ := personalEnv(t, nil)
	c := e.signup("anna@example.org")
	other := e.signup("bert@example.org")
	o := overview(t, e, c)
	var pistol store.Discipline
	for _, d := range o.Disciplines {
		if d.Key == "air_pistol" {
			pistol = d
		}
	}

	series := []domain.TrainingSeries{{Shots: 10, Total: 91}, {Shots: 10, Total: 93}}
	in := store.TrainingInput{DisciplineID: pistol.ID, Kind: "training", StartedAt: time.Now().Add(-48 * time.Hour), Series: series,
		Tags: []string{" tired ", "tired", ""}, ClientID: "7b0f3c1e-1111-4a6e-9a51-3d1f0a1b2c3d"}
	code, b := e.do(c, "POST", "/api/my/trainings", in)
	expect(t, "create", code, 201, b)
	var saved web.TrainingSaved
	json.Unmarshal(b, &saved)
	if saved.Session.Total != 184 || saved.Session.Shots != 20 || len(saved.Session.Tags) != 1 || saved.Feedback.Of != 1 {
		t.Fatalf("saved: %+v", saved)
	}
	first := saved.Session.ID
	code, b = e.do(c, "POST", "/api/my/trainings", in)
	expect(t, "retry", code, 201, b)
	json.Unmarshal(b, &saved)
	if saved.Session.ID != first {
		t.Fatalf("retry created a second session: %d", saved.Session.ID)
	}

	in.ClientID, in.StartedAt = "", time.Now()
	in.Series = []domain.TrainingSeries{{Shots: 10, Values: []float64{10, 10, 9, 9, 9, 9, 9, 9, 9, 10}, X: make([]float64, 10), Y: make([]float64, 10)}, {Shots: 10, Total: 95}}
	code, b = e.do(c, "POST", "/api/my/trainings", in)
	expect(t, "shots", code, 201, b)
	json.Unmarshal(b, &saved)
	if saved.Session.Total != 188 || !saved.Feedback.PersonalBest || saved.Feedback.Delta == nil || *saved.Feedback.Delta != 0.2 {
		t.Fatalf("feedback: %+v", saved.Feedback)
	}
	second := saved.Session.ID

	in.Series = []domain.TrainingSeries{{Shots: 10, Total: 101}}
	code, b = e.do(c, "POST", "/api/my/trainings", in)
	expect(t, "impossible total", code, 422, b)
	in.DisciplineID = 999999
	code, b = e.do(c, "POST", "/api/my/trainings", in)
	expect(t, "unknown discipline", code, 404, b)

	for _, path := range []string{fmt.Sprintf("/api/my/trainings/%d", first)} {
		code, b = e.do(other, "GET", path, nil)
		expect(t, "foreign get", code, 404, b)
		code, b = e.do(other, "PUT", path, store.TrainingInput{DisciplineID: pistol.ID, Kind: "training", StartedAt: time.Now(), Series: series})
		expect(t, "foreign put", code, 404, b)
		code, b = e.do(other, "DELETE", path, nil)
		expect(t, "foreign delete", code, 404, b)
	}
	code, b = e.do(other, "GET", "/api/my/trainings", nil)
	if code != 200 || strings.TrimSpace(string(b)) != "[]" {
		t.Fatalf("foreign list: %d %s", code, b)
	}

	code, b = e.do(c, "POST", "/api/my/disciplines", store.DisciplineInput{Name: "Field 300 m", SeriesCount: 3, ShotsPerSeries: 6, MaxShot: 5})
	expect(t, "own discipline", code, 201, b)
	var own store.Discipline
	json.Unmarshal(b, &own)
	code, b = e.do(other, "POST", "/api/my/trainings", store.TrainingInput{DisciplineID: own.ID, Kind: "training", StartedAt: time.Now(), Series: series})
	expect(t, "foreign discipline", code, 404, b)
	code, b = e.do(c, "POST", "/api/my/disciplines", store.DisciplineInput{Name: "x", SeriesCount: 1, ShotsPerSeries: 1, MaxShot: 5, Decimal: true})
	expect(t, "decimal needs 10", code, 422, b)

	code, b = e.do(c, "PUT", fmt.Sprintf("/api/my/goals/%d", pistol.ID), map[string]any{"targetTotal": 700})
	expect(t, "goal above max", code, 422, b)
	code, b = e.do(c, "PUT", fmt.Sprintf("/api/my/goals/%d", pistol.ID), map[string]any{"targetTotal": 470, "dueOn": "2027-03-01"})
	expect(t, "goal", code, 204, b)

	code, b = e.do(c, "GET", fmt.Sprintf("/api/my/stats/%d", pistol.ID), nil)
	expect(t, "stats", code, 200, b)
	var st web.MyStats
	json.Unmarshal(b, &st)
	if st.Stats.Sessions != 2 || st.Stats.Distribution[10] != 3 || st.Stats.Group == nil || st.Goal == nil || st.Goal.TargetTotal != 470 {
		t.Fatalf("stats: %+v", st)
	}

	o = overview(t, e, c)
	if o.Week.Streak < 1 || len(o.Recent) != 2 || len(o.Tags) != 1 || o.Tags[0] != "tired" {
		t.Fatalf("overview: %+v", o)
	}

	code, b = e.do(c, "PUT", fmt.Sprintf("/api/my/trainings/%d", second), store.TrainingInput{DisciplineID: pistol.ID, Kind: "competition", StartedAt: time.Now(), Series: series, Notes: "Club match"})
	expect(t, "update", code, 200, b)
	saved = web.TrainingSaved{}
	json.Unmarshal(b, &saved)
	if saved.Session.Kind != "competition" || saved.Session.Total != 184 || saved.Session.Series[0].Values != nil {
		t.Fatalf("updated: %+v", saved.Session)
	}

	code, b = e.do(c, "GET", "/api/my/export", nil)
	if code != 200 || !strings.Contains(string(b), "Air pistol 10 m;;competition") || !strings.Contains(string(b), "Club match") {
		t.Fatalf("export: %d %s", code, b)
	}

	code, b = e.do(c, "GET", fmt.Sprintf("/api/my/trainings/%d", second), nil)
	if code != 200 || !strings.Contains(string(b), "Club match") {
		t.Fatalf("get session: %d %s", code, b)
	}
	code, b = e.do(c, "DELETE", fmt.Sprintf("/api/my/trainings/%d", first), nil)
	expect(t, "delete session", code, 204, b)
	code, b = e.do(c, "GET", fmt.Sprintf("/api/my/trainings/%d", first), nil)
	expect(t, "deleted session", code, 404, b)
	code, b = e.do(c, "DELETE", fmt.Sprintf("/api/my/goals/%d", pistol.ID), nil)
	expect(t, "remove goal", code, 204, b)
	code, b = e.do(c, "GET", fmt.Sprintf("/api/my/stats/%d", pistol.ID), nil)
	if code != 200 || strings.Contains(string(b), `"targetTotal"`) {
		t.Fatalf("stats without goal: %d %s", code, b)
	}
	code, b = e.do(c, "PUT", fmt.Sprintf("/api/my/disciplines/%d", own.ID), store.DisciplineInput{Name: "Field 300 m B", SeriesCount: 3, ShotsPerSeries: 6, MaxShot: 5})
	expect(t, "change discipline", code, 204, b)
	if o := overview(t, e, c); !slices.ContainsFunc(o.Disciplines, func(d store.Discipline) bool { return d.Name == "Field 300 m B" }) {
		t.Fatalf("changed discipline missing: %+v", o.Disciplines)
	}
	code, b = e.do(c, "DELETE", fmt.Sprintf("/api/my/disciplines/%d", own.ID), nil)
	expect(t, "delete unused discipline", code, 204, b)
	code, b = e.do(c, "DELETE", "/api/my/account", map[string]string{"password": "wrong"})
	expect(t, "delete with wrong password", code, 403, b)
	code, b = e.do(c, "DELETE", "/api/my/account", map[string]string{"password": "password-anna@example.org"})
	expect(t, "delete account", code, 204, b)
	code, b = e.do(e.client(), "POST", "/api/auth/login", map[string]string{"name": "anna@example.org", "password": "password-anna@example.org"})
	expect(t, "login after delete", code, 401, b)
}

func TestWeapons(t *testing.T) {
	e, _ := personalEnv(t, nil)
	c := e.signup("carl@example.org")
	other := e.signup("dora@example.org")
	var pistol store.Discipline
	for _, d := range overview(t, e, c).Disciplines {
		if d.Key == "air_pistol" {
			pistol = d
		}
	}
	weapon := func(cl *http.Client, name string) store.Weapon {
		code, b := e.do(cl, "POST", "/api/my/weapons", store.WeaponInput{Name: name})
		expect(t, "weapon "+name, code, 201, b)
		var w store.Weapon
		json.Unmarshal(b, &w)
		return w
	}
	steyr, walther, foreign := weapon(c, " Steyr EVO 10 "), weapon(c, "Walther LP500"), weapon(other, "Morini")
	if steyr.Name != "Steyr EVO 10" {
		t.Fatalf("name not trimmed: %q", steyr.Name)
	}
	code, b := e.do(c, "POST", "/api/my/weapons", store.WeaponInput{Name: "steyr evo 10"})
	expect(t, "duplicate name", code, 409, b)
	code, b = e.do(c, "POST", "/api/my/weapons", store.WeaponInput{Name: "  "})
	expect(t, "blank name", code, 422, b)

	log := func(w *store.Weapon, totals ...float64) {
		in := store.TrainingInput{DisciplineID: pistol.ID, Kind: "training", StartedAt: time.Now()}
		if w != nil {
			in.WeaponID = &w.ID
		}
		for _, x := range totals {
			in.Series = append(in.Series, domain.TrainingSeries{Shots: 10, Total: x})
		}
		code, b := e.do(c, "POST", "/api/my/trainings", in)
		expect(t, "log", code, 201, b)
	}
	log(&steyr, 95, 93)
	log(&steyr, 94, 92)
	log(&walther, 90, 88)
	log(nil, 80)
	in := store.TrainingInput{DisciplineID: pistol.ID, Kind: "training", StartedAt: time.Now(), WeaponID: &foreign.ID, Series: []domain.TrainingSeries{{Shots: 10, Total: 90}}}
	code, b = e.do(c, "POST", "/api/my/trainings", in)
	expect(t, "foreign weapon", code, 404, b)

	var st web.MyStats
	code, b = e.do(c, "GET", fmt.Sprintf("/api/my/stats/%d", pistol.ID), nil)
	expect(t, "stats", code, 200, b)
	json.Unmarshal(b, &st)
	w := st.Stats.Weapons
	if st.Stats.Sessions != 4 || len(w) != 2 || w[0].WeaponID != steyr.ID || w[0].Sessions != 2 || w[0].PerShot != 9.35 || w[1].WeaponID != walther.ID {
		t.Fatalf("compare: %+v", st.Stats)
	}
	code, b = e.do(c, "GET", fmt.Sprintf("/api/my/stats/%d?weapon=%d", pistol.ID, walther.ID), nil)
	expect(t, "stats of one weapon", code, 200, b)
	st = web.MyStats{}
	json.Unmarshal(b, &st)
	if st.Stats.Sessions != 1 || st.Stats.PerShot != 8.9 || len(st.Stats.Weapons) != 2 {
		t.Fatalf("filtered: %+v", st.Stats)
	}

	code, b = e.do(c, "GET", "/api/my/export", nil)
	if code != 200 || !strings.Contains(string(b), "Air pistol 10 m;Walther LP500;training") {
		t.Fatalf("export: %d %s", code, b)
	}

	for _, req := range []struct{ method, path string }{{"PUT", fmt.Sprintf("/api/my/weapons/%d", steyr.ID)}, {"DELETE", fmt.Sprintf("/api/my/weapons/%d", steyr.ID)}} {
		code, b = e.do(other, req.method, req.path, store.WeaponInput{Name: "Mine"})
		expect(t, "foreign "+req.method, code, 404, b)
	}
	if o := overview(t, e, other); len(o.Weapons) != 1 || o.Weapons[0].Name != "Morini" {
		t.Fatalf("foreign overview: %+v", o.Weapons)
	}
	code, b = e.do(c, "PUT", fmt.Sprintf("/api/my/weapons/%d", walther.ID), store.WeaponInput{Name: "Steyr EVO 10"})
	expect(t, "rename to taken name", code, 409, b)
	code, b = e.do(c, "PUT", fmt.Sprintf("/api/my/weapons/%d", walther.ID), store.WeaponInput{Name: "Walther LP400"})
	expect(t, "rename", code, 204, b)
	code, b = e.do(c, "DELETE", fmt.Sprintf("/api/my/weapons/%d", steyr.ID), nil)
	expect(t, "delete", code, 204, b)
	o := overview(t, e, c)
	if len(o.Weapons) != 1 || o.Weapons[0].Name != "Walther LP400" || len(o.Recent) != 4 {
		t.Fatalf("after delete: %+v", o)
	}
	for _, s := range o.Recent {
		if s.WeaponID != nil && *s.WeaponID == steyr.ID {
			t.Fatalf("session still has the deleted weapon: %+v", s)
		}
	}
}

func TestCompetitionImport(t *testing.T) {
	e, _ := personalEnv(t, nil)
	ctx := context.Background()
	tid := e.tenant("PS Example")
	e.user("staff", tid, store.RoleStaff)
	staff := e.login("staff")
	ev := e.event(tid, "Example Championship")
	prog, _ := e.st.CreateProgramme(ctx, tid, ev, domain.ProgrammeDefinition{
		Name: "Championship", Series: []domain.Series{{Shots: 10, MaxPoints: 100}, {Shots: 10, MaxPoints: 100}},
		Prices: []domain.Price{{AgeCategory: domain.Elite, Price: 20}}, Awards: []domain.AwardTier{},
	})
	shooter, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Felix", Lastname: "Example", Birthday: "1990-05-01", LicenseNumber: "012345"}, nil)
	rival, _ := e.st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Rita", Lastname: "Rival", Birthday: "1985-03-03"}, nil)
	for sh, points := range map[int64][]int{shooter: {95, 92}, rival: {96, 94}} {
		code, b := e.do(staff, "POST", fmt.Sprintf("/api/tenants/%d/events/%d/registrations", tid, ev), web.RegistrationInput{
			ShooterID: sh, Programmes: []domain.ProgrammeChoice{{ProgrammeID: prog}},
		})
		expect(t, "register", code, 201, b)
		var r store.Registration
		json.Unmarshal(b, &r)
		code, b = e.do(staff, "PUT", fmt.Sprintf("/api/tenants/%d/entries/%d/results", tid, r.Entries[0].ID),
			map[string]any{"series": []map[string]any{{"points": points[0]}, {"points": points[1]}}})
		expect(t, "results", code, 200, b)
	}
	c := e.signup("felix@example.org")
	code, b := e.do(c, "GET", "/api/my/competitions", nil)
	if code != 200 || strings.TrimSpace(string(b)) != "[]" {
		t.Fatalf("without license: %d %s", code, b)
	}
	code, b = e.do(c, "PUT", "/api/my/profile", store.Profile{Firstname: "Felix", Lastname: "Example", Birthday: "1990-05-01", LicenseNumber: "12345"})
	expect(t, "profile", code, 204, b)
	code, b = e.do(c, "GET", "/api/my/competitions", nil)
	if code != 200 || strings.TrimSpace(string(b)) != "[]" {
		t.Fatalf("before the event is closed: %d %s", code, b)
	}
	if err := e.st.SetEventArchived(ctx, tid, ev, true); err != nil {
		t.Fatal(err)
	}
	code, b = e.do(c, "GET", "/api/my/competitions", nil)
	expect(t, "competitions", code, 200, b)
	var list []store.CompetitionResult
	json.Unmarshal(b, &list)
	if len(list) != 1 || list[0].Total != 187 || list[0].Club != "PS Example" || list[0].MaxPerShot != 10 || list[0].Rank != 2 || list[0].Participants != 2 {
		t.Fatalf("competitions: %+v", list)
	}
	if o := overview(t, e, c); o.Competitions != 1 {
		t.Fatalf("overview competitions: %d", o.Competitions)
	}
	other := e.signup("other@example.org")
	code, b = e.do(other, "POST", fmt.Sprintf("/api/my/competitions/%d", list[0].EntryID), map[string]any{"disciplineId": o2discipline(t, e, other)})
	expect(t, "foreign import", code, 404, b)

	code, b = e.do(c, "POST", fmt.Sprintf("/api/my/competitions/%d", list[0].EntryID), map[string]any{"disciplineId": o2discipline(t, e, c)})
	expect(t, "import", code, 201, b)
	var saved web.TrainingSaved
	json.Unmarshal(b, &saved)
	if saved.Session.Kind != "competition" || saved.Session.Total != 187 || !saved.Session.Imported || saved.Session.Place != "PS Example" ||
		saved.Session.Rank == nil || *saved.Session.Rank != 2 || *saved.Session.Participants != 2 {
		t.Fatalf("imported: %+v", saved.Session)
	}
	code, b = e.do(c, "GET", "/api/my/competitions", nil)
	if code != 200 || strings.TrimSpace(string(b)) != "[]" {
		t.Fatalf("after import: %d %s", code, b)
	}
}

func o2discipline(t *testing.T, e *env, c *http.Client) int64 {
	return overview(t, e, c).Disciplines[0].ID
}
