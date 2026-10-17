package web

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
)

const defaultZone = "Europe/Zurich"

// WeekSummary is the activity of the current week and the streak of weeks with training.
type WeekSummary struct {
	Sessions int `json:"sessions"`
	Shots    int `json:"shots"`
	Streak   int `json:"streak" doc:"Consecutive weeks with at least one session"`
}

// MyOverview is everything the home screen of the personal log needs.
type MyOverview struct {
	Email         string                  `json:"email"`
	EmailVerified bool                    `json:"emailVerified"`
	Profile       store.Profile           `json:"profile" doc:"Created empty on the first visit; names, birthday and license are optional"`
	Disciplines   []store.Discipline      `json:"disciplines"`
	Goals         []store.Goal            `json:"goals"`
	Tags          []string                `json:"tags"`
	Weapons       []store.Weapon          `json:"weapons"`
	Recent        []store.TrainingSession `json:"recent"`
	Week          WeekSummary             `json:"week"`
	Competitions  int                     `json:"competitions" doc:"Competition results waiting to be added to the log"`
}

// TrainingPath addresses one session of the user.
type TrainingPath struct {
	ID int64 `path:"id"`
}

// DisciplinePath addresses a discipline.
type DisciplinePath struct {
	ID int64 `path:"id"`
}

// WeaponPath addresses a weapon of the user.
type WeaponPath struct {
	ID int64 `path:"id"`
}

// TrainingSaved is the saved session with the feedback for it.
type TrainingSaved struct {
	Session  store.TrainingSession `json:"session"`
	Feedback domain.Feedback       `json:"feedback"`
}

// MyStats is the analysis of one discipline.
type MyStats struct {
	Discipline store.Discipline     `json:"discipline"`
	Stats      domain.TrainingStats `json:"stats"`
	Goal       *store.Goal          `json:"goal"`
}

// zone loads the named time zone, falling back to defaultZone and then UTC.
func zone(name string) *time.Location {
	if loc, err := time.LoadLocation(name); err == nil && name != "" {
		return loc
	}
	loc, err := time.LoadLocation(defaultZone)
	if err != nil {
		return time.UTC
	}
	return loc
}

// cleanTags trims tags, cuts them to 40 bytes and drops empty ones and duplicates.
func cleanTags(tags []string) []string {
	out := []string{}
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if r := []rune(t); len(r) > 40 {
			t = string(r[:40])
		}
		if t != "" && !slices.Contains(out, t) {
			out = append(out, t)
		}
	}
	return out
}

// checkTraining checks that discipline and weapon belong to the user, validates the series and cleans the tags of in.
func (s *Server) checkTraining(ctx context.Context, userID int64, in *store.TrainingInput) error {
	d, err := s.Store.GetDiscipline(ctx, userID, in.DisciplineID)
	if err != nil {
		return notFound(err, "Discipline not found")
	}
	if in.WeaponID != nil {
		if _, err := s.Store.GetWeapon(ctx, userID, *in.WeaponID); err != nil {
			return notFound(err, "Weapon not found")
		}
	}
	series, errs := domain.ValidateTraining(d.Rules(), in.Series)
	if len(errs) > 0 {
		return invalidList(errs)
	}
	in.Series, in.Tags = series, cleanTags(in.Tags)
	return nil
}

// saved returns a saved session of the user together with its feedback.
func (s *Server) saved(ctx context.Context, userID, id int64) (*struct{ Body TrainingSaved }, error) {
	t, err := s.Store.GetTraining(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	facts, err := s.Store.TrainingFacts(ctx, userID, t.DisciplineID)
	if err != nil {
		return nil, err
	}
	return &struct{ Body TrainingSaved }{Body: TrainingSaved{Session: t, Feedback: domain.FeedbackFor(facts, id)}}, nil
}

// overview builds the start page of the personal log: account, profile, disciplines, goals, weapons, recent sessions, the week and competitions.
func (s *Server) overview(ctx context.Context, u store.User, tz string) (MyOverview, error) {
	var o MyOverview
	a, err := s.Store.GetAccount(ctx, u.ID)
	if err != nil {
		return o, err
	}
	o.Email, o.EmailVerified = a.Email, a.EmailVerified
	if o.Profile, err = s.Store.GetProfile(ctx, u.ID); errors.Is(err, store.ErrNotFound) {
		// Every user has the training log: the first visit creates an empty profile.
		if err = s.Store.SaveProfile(ctx, u.ID, store.Profile{}); err != nil {
			return o, err
		}
	} else if err != nil {
		return o, err
	}
	if o.Disciplines, err = s.Store.ListDisciplines(ctx, u.ID); err != nil {
		return o, err
	}
	if o.Goals, err = s.Store.ListGoals(ctx, u.ID); err != nil {
		return o, err
	}
	if o.Tags, err = s.Store.Tags(ctx, u.ID); err != nil {
		return o, err
	}
	if o.Weapons, err = s.Store.ListWeapons(ctx, u.ID); err != nil {
		return o, err
	}
	if o.Recent, err = s.Store.ListTrainings(ctx, u.ID, store.TrainingFilter{Limit: 5}); err != nil {
		return o, err
	}
	activity, err := s.Store.TrainingActivity(ctx, u.ID)
	if err != nil {
		return o, err
	}
	loc, now := zone(tz), time.Now()
	times := make([]time.Time, len(activity))
	monday := now.In(loc).AddDate(0, 0, -((int(now.In(loc).Weekday()) + 6) % 7))
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, loc)
	for i, x := range activity {
		times[i] = x.At
		if !x.At.Before(monday) {
			o.Week.Sessions++
			o.Week.Shots += x.Shots
		}
	}
	o.Week.Streak = domain.WeekStreak(times, now, loc)
	if o.Profile.LicenseNumber != "" && o.Profile.Birthday != "" {
		c, err := s.Store.CompetitionResults(ctx, u.ID)
		if err != nil {
			return o, err
		}
		o.Competitions = len(c)
	}
	return o, nil
}

// registerMy registers the personal log endpoints under /api/my.
func (s *Server) registerMy(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-my-overview", Method: http.MethodGet, Path: "/api/my", Tags: []string{"my"},
		Summary: "Profile, access, disciplines, goals and recent activity of the personal log",
	}, func(ctx context.Context, in *struct {
		TZ string `query:"tz" doc:"IANA time zone for the week summary, default Europe/Zurich"`
	}) (*struct{ Body MyOverview }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		o, err := s.overview(ctx, u, in.TZ)
		return &struct{ Body MyOverview }{Body: o}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "set-my-profile", Method: http.MethodPut, Path: "/api/my/profile", Tags: []string{"my"},
		Summary: "Set up or change the profile", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *struct{ Body store.Profile }) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		p := in.Body
		p.Firstname, p.Lastname, p.LicenseNumber = strings.TrimSpace(p.Firstname), strings.TrimSpace(p.Lastname), strings.TrimSpace(p.LicenseNumber)
		if p.Birthday != "" {
			if _, err := time.Parse("2006-01-02", p.Birthday); err != nil {
				return nil, invalid("Birthday: use YYYY-MM-DD")
			}
		}
		return nil, s.Store.SaveProfile(ctx, u.ID, p)
	})

	huma.Register(api, huma.Operation{
		OperationID: "set-my-email", Method: http.MethodPut, Path: "/api/my/email", Tags: []string{"my"},
		Summary: "Change the email address and send a confirmation link", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *struct {
		Body struct {
			Email    string `json:"email" format:"email" maxLength:"200"`
			Password string `json:"password" minLength:"1"`
		}
	}) (*struct{}, error) {
		u, err := s.checkPassword(ctx, in.Body.Password)
		if err != nil {
			return nil, err
		}
		email := strings.ToLower(strings.TrimSpace(in.Body.Email))
		if err := s.Store.SetEmail(ctx, u.ID, email); errors.Is(err, store.ErrExists) {
			return nil, conflict("Another account uses this email address.")
		} else if err != nil {
			return nil, err
		}
		return nil, s.sendVerification(ctx, u.ID, email, s.firstname(ctx, u.ID), u.Locale)
	})

	huma.Register(api, huma.Operation{
		OperationID: "resend-verification", Method: http.MethodPost, Path: "/api/my/email/verify", Tags: []string{"my"},
		Summary: "Send the confirmation link again", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if err := s.rateLimited(ctx, fmt.Sprintf("verify:user:%d", u.ID), forgotWindow, forgotIPMax); err != nil {
			return nil, err
		}
		a, err := s.Store.GetAccount(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		if a.Email == "" || a.EmailVerified {
			return nil, conflict("There is no email address to confirm.")
		}
		return nil, s.sendVerification(ctx, u.ID, a.Email, s.firstname(ctx, u.ID), u.Locale)
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-my-trainings", Method: http.MethodGet, Path: "/api/my/trainings", Tags: []string{"my"},
		Summary: "Sessions of the log, newest first",
	}, func(ctx context.Context, in *struct {
		Discipline int64     `query:"discipline"`
		Kind       string    `query:"kind" enum:"training,competition,"`
		Before     time.Time `query:"before" doc:"Only sessions that started earlier, for paging"`
		Limit      int       `query:"limit" minimum:"1" maximum:"200" default:"50"`
	}) (*struct{ Body []store.TrainingSession }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		f := store.TrainingFilter{DisciplineID: in.Discipline, Kind: in.Kind, Limit: in.Limit}
		if !in.Before.IsZero() {
			f.Before = &in.Before
		}
		list, err := s.Store.ListTrainings(ctx, u.ID, f)
		return &struct{ Body []store.TrainingSession }{Body: list}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-my-training", Method: http.MethodPost, Path: "/api/my/trainings", Tags: []string{"my"},
		Summary: "Log a session; a retry with the same clientId updates it", DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *struct{ Body store.TrainingInput }) (*struct{ Body TrainingSaved }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if err := s.checkTraining(ctx, u.ID, &in.Body); err != nil {
			return nil, err
		}
		id, err := s.Store.SaveTraining(ctx, u.ID, 0, in.Body)
		if err != nil {
			return nil, err
		}
		return s.saved(ctx, u.ID, id)
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-my-training", Method: http.MethodGet, Path: "/api/my/trainings/{id}", Tags: []string{"my"},
		Summary: "One session with its feedback",
	}, func(ctx context.Context, in *TrainingPath) (*struct{ Body TrainingSaved }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		out, err := s.saved(ctx, u.ID, in.ID)
		return out, notFound(err, "Session not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-my-training", Method: http.MethodPut, Path: "/api/my/trainings/{id}", Tags: []string{"my"},
		Summary: "Change a session",
	}, func(ctx context.Context, in *struct {
		TrainingPath
		Body store.TrainingInput
	}) (*struct{ Body TrainingSaved }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if err := s.checkTraining(ctx, u.ID, &in.Body); err != nil {
			return nil, err
		}
		if _, err := s.Store.SaveTraining(ctx, u.ID, in.ID, in.Body); err != nil {
			return nil, notFound(err, "Session not found")
		}
		return s.saved(ctx, u.ID, in.ID)
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-my-training", Method: http.MethodDelete, Path: "/api/my/trainings/{id}", Tags: []string{"my"},
		Summary: "Delete a session", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *TrainingPath) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.DeleteTraining(ctx, u.ID, in.ID), "Session not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-my-stats", Method: http.MethodGet, Path: "/api/my/stats/{id}", Tags: []string{"my"},
		Summary: "Analysis of all sessions of a discipline, or of those shot with one weapon",
	}, func(ctx context.Context, in *struct {
		DisciplinePath
		Weapon int64 `query:"weapon" doc:"Only sessions shot with this weapon"`
	}) (*struct{ Body MyStats }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		d, err := s.Store.GetDiscipline(ctx, u.ID, in.ID)
		if err != nil {
			return nil, notFound(err, "Discipline not found")
		}
		facts, err := s.Store.TrainingFacts(ctx, u.ID, d.ID)
		if err != nil {
			return nil, err
		}
		now := time.Now()
		weapons := domain.CompareWeapons(d.Rules(), facts, now)
		if in.Weapon != 0 {
			facts = slices.DeleteFunc(facts, func(f domain.SessionFacts) bool { return f.WeaponID != in.Weapon })
		}
		out := &struct{ Body MyStats }{Body: MyStats{Discipline: d, Stats: domain.Stats(d.Rules(), facts, now)}}
		out.Body.Stats.Weapons = weapons
		goals, err := s.Store.ListGoals(ctx, u.ID)
		for _, g := range goals {
			if g.DisciplineID == d.ID {
				out.Body.Goal = &g
			}
		}
		return out, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-my-discipline", Method: http.MethodPost, Path: "/api/my/disciplines", Tags: []string{"my"},
		Summary: "Add an own discipline", DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *struct{ Body store.DisciplineInput }) (*struct{ Body store.Discipline }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if in.Body.Decimal && in.Body.MaxShot != 10 {
			return nil, invalid("Decimal scoring needs 10 as the highest value of a shot")
		}
		id, err := s.Store.CreateDiscipline(ctx, u.ID, in.Body)
		if err != nil {
			return nil, err
		}
		d, err := s.Store.GetDiscipline(ctx, u.ID, id)
		return &struct{ Body store.Discipline }{Body: d}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-my-discipline", Method: http.MethodPut, Path: "/api/my/disciplines/{id}", Tags: []string{"my"},
		Summary: "Change an own discipline", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *struct {
		DisciplinePath
		Body store.DisciplineInput
	}) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if in.Body.Decimal && in.Body.MaxShot != 10 {
			return nil, invalid("Decimal scoring needs 10 as the highest value of a shot")
		}
		return nil, notFound(s.Store.UpdateDiscipline(ctx, u.ID, in.ID, in.Body), "Discipline not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-my-discipline", Method: http.MethodDelete, Path: "/api/my/disciplines/{id}", Tags: []string{"my"},
		Summary: "Delete an own discipline without sessions", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *DisciplinePath) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		err = s.Store.DeleteDiscipline(ctx, u.ID, in.ID)
		if errors.Is(err, store.ErrInUse) {
			return nil, conflict("Sessions use this discipline")
		}
		return nil, notFound(err, "Discipline not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-my-weapon", Method: http.MethodPost, Path: "/api/my/weapons", Tags: []string{"my"},
		Summary: "Add a weapon", DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *struct{ Body store.WeaponInput }) (*struct{ Body store.Weapon }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if in.Body.Name = strings.TrimSpace(in.Body.Name); in.Body.Name == "" {
			return nil, invalid("Name is empty")
		}
		w, err := s.Store.CreateWeapon(ctx, u.ID, in.Body)
		if errors.Is(err, store.ErrExists) {
			return nil, conflict("You already have a weapon with this name")
		}
		return &struct{ Body store.Weapon }{Body: w}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "rename-my-weapon", Method: http.MethodPut, Path: "/api/my/weapons/{id}", Tags: []string{"my"},
		Summary: "Rename a weapon", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *struct {
		WeaponPath
		Body store.WeaponInput
	}) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if in.Body.Name = strings.TrimSpace(in.Body.Name); in.Body.Name == "" {
			return nil, invalid("Name is empty")
		}
		err = s.Store.RenameWeapon(ctx, u.ID, in.ID, in.Body)
		if errors.Is(err, store.ErrExists) {
			return nil, conflict("You already have a weapon with this name")
		}
		return nil, notFound(err, "Weapon not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-my-weapon", Method: http.MethodDelete, Path: "/api/my/weapons/{id}", Tags: []string{"my"},
		Summary: "Delete a weapon; its sessions stay without a weapon", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *WeaponPath) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.DeleteWeapon(ctx, u.ID, in.ID), "Weapon not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "set-my-goal", Method: http.MethodPut, Path: "/api/my/goals/{id}", Tags: []string{"my"},
		Summary: "Set the goal of a discipline", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *struct {
		DisciplinePath
		Body struct {
			TargetTotal float64 `json:"targetTotal" exclusiveMinimum:"0"`
			DueOn       string  `json:"dueOn" required:"false"`
		}
	}) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		d, err := s.Store.GetDiscipline(ctx, u.ID, in.ID)
		if err != nil {
			return nil, notFound(err, "Discipline not found")
		}
		if in.Body.TargetTotal > float64(d.Rules().StandardShots())*d.Rules().ShotMax() {
			return nil, invalid("The goal is higher than the maximum")
		}
		if in.Body.DueOn != "" {
			if _, err := time.Parse("2006-01-02", in.Body.DueOn); err != nil {
				return nil, invalid("Date: use YYYY-MM-DD")
			}
		}
		return nil, s.Store.SetGoal(ctx, u.ID, store.Goal{DisciplineID: d.ID, TargetTotal: in.Body.TargetTotal, DueOn: in.Body.DueOn})
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-my-goal", Method: http.MethodDelete, Path: "/api/my/goals/{id}", Tags: []string{"my"},
		Summary: "Remove the goal of a discipline", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *DisciplinePath) (*struct{}, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.DeleteGoal(ctx, u.ID, in.ID), "Goal not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-my-competitions", Method: http.MethodGet, Path: "/api/my/competitions", Tags: []string{"my"},
		Summary: "Results of Standblatt events with the user's license number and birthday that are not in the log yet",
	}, func(ctx context.Context, _ *struct{}) (*struct{ Body []store.CompetitionResult }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		list, err := s.Store.CompetitionResults(ctx, u.ID)
		return &struct{ Body []store.CompetitionResult }{Body: list}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "import-my-competition", Method: http.MethodPost, Path: "/api/my/competitions/{id}", Tags: []string{"my"},
		Summary: "Add a competition result to the log", DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *struct {
		ID   int64 `path:"id"`
		Body struct {
			DisciplineID int64 `json:"disciplineId"`
		}
	}) (*struct{ Body TrainingSaved }, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		if _, err := s.Store.GetDiscipline(ctx, u.ID, in.Body.DisciplineID); err != nil {
			return nil, notFound(err, "Discipline not found")
		}
		id, err := s.Store.ImportCompetition(ctx, u.ID, in.ID, in.Body.DisciplineID, zone(""))
		if err != nil {
			return nil, notFound(err, "Result not found")
		}
		return s.saved(ctx, u.ID, id)
	})

	huma.Register(api, huma.Operation{
		OperationID: "export-my-trainings", Method: http.MethodGet, Path: "/api/my/export", Tags: []string{"my"},
		Summary: "All sessions as CSV, one line per series; always available",
	}, func(ctx context.Context, _ *struct{}) (*FileOutput, error) {
		u, err := currentUser(ctx)
		if err != nil {
			return nil, err
		}
		b, err := s.exportTrainings(ctx, u.ID)
		return &FileOutput{ContentType: "text/csv; charset=utf-8", ContentDisposition: `attachment; filename="standblatt-trainings.csv"`, Body: b}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-my-account", Method: http.MethodDelete, Path: "/api/my/account", Tags: []string{"my"},
		Summary: "Delete the account with all trainings; not possible while it belongs to a club",
	}, func(ctx context.Context, in *struct {
		Body struct {
			Password string `json:"password" minLength:"1"`
		}
	}) (*LogoutOutput, error) {
		u, err := s.checkPassword(ctx, in.Body.Password)
		if err != nil {
			return nil, err
		}
		ms, err := s.Store.Memberships(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		if len(ms) > 0 {
			return nil, conflict("Your account belongs to a club. Ask the club admin to remove you first.")
		}
		if err := s.Store.DeleteUser(ctx, u.ID); err != nil {
			return nil, err
		}
		s.Log.Info("account deleted", "user", u.Name)
		return &LogoutOutput{SetCookie: sessionCookieFor(ctx, "", -1)}, nil
	})
}

// firstname returns the user's first name from the profile, empty when there is none.
func (s *Server) firstname(ctx context.Context, userID int64) string {
	p, _ := s.Store.GetProfile(ctx, userID)
	return p.Firstname
}

// checkPassword returns the logged-in user if password is theirs, 403 otherwise.
func (s *Server) checkPassword(ctx context.Context, password string) (store.User, error) {
	u, err := currentUser(ctx)
	if err != nil {
		return u, err
	}
	_, hash, err := s.Store.Credentials(ctx, u.Name)
	if err != nil {
		return u, err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return u, forbidden("Wrong password.")
	}
	return u, nil
}

// exportTrainings writes all sessions of the user as semicolon CSV with BOM, one line per series.
func (s *Server) exportTrainings(ctx context.Context, userID int64) ([]byte, error) {
	list, err := s.Store.ListTrainings(ctx, userID, store.TrainingFilter{})
	if err != nil {
		return nil, err
	}
	ds, err := s.Store.ListDisciplines(ctx, userID)
	if err != nil {
		return nil, err
	}
	names := map[int64]string{}
	for _, d := range ds {
		names[d.ID] = d.Name
	}
	ws, err := s.Store.ListWeapons(ctx, userID)
	if err != nil {
		return nil, err
	}
	weapons := map[int64]string{}
	for _, x := range ws {
		weapons[x.ID] = x.Name
	}
	var buf bytes.Buffer
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(&buf)
	w.Comma = ';'
	w.Write([]string{"date", "discipline", "weapon", "kind", "title", "place", "rating", "tags", "notes", "series", "shots", "total", "values"})
	num := func(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }
	for _, t := range list {
		rating, weapon := "", ""
		if t.Rating != nil {
			rating = strconv.Itoa(*t.Rating)
		}
		if t.WeaponID != nil {
			weapon = weapons[*t.WeaponID]
		}
		for i, x := range t.Series {
			values := make([]string, len(x.Values))
			for j, v := range x.Values {
				values[j] = num(v)
			}
			w.Write([]string{t.StartedAt.Format(time.RFC3339), names[t.DisciplineID], weapon, t.Kind, t.Title, t.Place, rating,
				strings.Join(t.Tags, ", "), t.Notes, strconv.Itoa(i + 1), strconv.Itoa(x.Shots), num(x.Total), strings.Join(values, " ")})
		}
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}
