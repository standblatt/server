package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

// EntryPath identifies an entry (a programme of a registration) in the URL.
type EntryPath struct {
	TenantID int64 `path:"tenantId"`
	EntryID  int64 `path:"entryId"`
}

// EntryResult is one programme of a registration with its results and awards.
type EntryResult struct {
	EntryID        int64                 `json:"entryId"`
	ProgrammeID    int64                 `json:"programmeId"`
	Title          string                `json:"title"`
	Status         string                `json:"status" enum:"pending,completed"`
	Series         []domain.Series       `json:"series"`
	MaxPoints      int                   `json:"maxPoints"`
	CombinedFrom   string                `json:"combinedFrom,omitempty" doc:"Programme the series are taken from"`
	CombinedSeries []int                 `json:"combinedSeries"`
	Results        []domain.SeriesResult `json:"results"`
	Total          float64               `json:"total"`
	Mouches        int                   `json:"mouches"`
	Award          *domain.Tier          `json:"award,omitempty"`
	EligibleAwards []domain.Tier         `json:"eligibleAwards"`
	Form           bool                  `json:"form" doc:"The programme has a print form for its results"`
}

// ResultSheet is a registration with the results of all its programmes.
type ResultSheet struct {
	Registration store.Registration `json:"registration"`
	Entries      []EntryResult      `json:"entries"`
}

// programmeTitle returns the programme name with the variant in parentheses when there is one.
func programmeTitle(name, variant string) string {
	if variant == "" {
		return name
	}
	return name + " (" + variant + ")"
}

// resultSheet builds the results of every entry of a registration with totals, awards and the awards it can choose from.
func (s *Server) resultSheet(ctx context.Context, tenantID, registrationID int64) (ResultSheet, error) {
	reg, err := s.Store.GetRegistration(ctx, tenantID, registrationID)
	if err != nil {
		return ResultSheet{}, notFound(err, "Registration not found")
	}
	programmes, err := s.Store.ListProgrammes(ctx, tenantID, reg.EventID)
	if err != nil {
		return ResultSheet{}, err
	}
	byID := map[int64]store.Programme{}
	ids := make([]int64, 0, len(programmes))
	for _, p := range programmes {
		byID[p.ID] = p
		ids = append(ids, p.ID)
	}
	tiers, err := s.Store.ProgrammeTiers(ctx, tenantID, ids)
	if err != nil {
		return ResultSheet{}, err
	}
	results, err := s.Store.ListResults(ctx, tenantID, reg.ID)
	if err != nil {
		return ResultSheet{}, err
	}
	sheet := ResultSheet{Registration: reg, Entries: []EntryResult{}}
	for _, e := range reg.Entries {
		p := byID[e.ProgrammeID]
		er := EntryResult{
			EntryID: e.ID, ProgrammeID: p.ID, Title: programmeTitle(p.Name, p.Variant), Status: e.Status,
			Series: p.Series, MaxPoints: p.MaxPoints, CombinedSeries: e.CombinedSeries,
			Results: results[e.ID], EligibleAwards: []domain.Tier{}, Form: p.LayoutID != nil,
		}
		if er.Results == nil {
			er.Results = []domain.SeriesResult{}
		}
		if len(e.CombinedSeries) > 0 && p.CombinableWithID != nil {
			t := byID[*p.CombinableWithID]
			er.CombinedFrom = programmeTitle(t.Name, t.Variant)
		}
		er.Total, er.Mouches = domain.Totals(er.Results)
		if e.Status == "completed" {
			er.EligibleAwards = domain.EligibleAwards(tiers[p.ID], reg.AgeCategory, er.Total)
		}
		for _, t := range tiers[p.ID] {
			if e.AwardTierID != nil && t.ID == *e.AwardTierID {
				award := t
				er.Award = &award
			}
		}
		if er.EligibleAwards == nil {
			er.EligibleAwards = []domain.Tier{}
		}
		sheet.Entries = append(sheet.Entries, er)
	}
	return sheet, nil
}

// ResultLine is the result of a shooter in a programme, entered without the desk.
type ResultLine struct {
	ShooterID   int64                 `json:"shooterId"`
	ClubID      *int64                `json:"clubId,omitempty" required:"false" doc:"One of the shooter's clubs, used when a registration is created"`
	ProgrammeID int64                 `json:"programmeId"`
	Series      []domain.SeriesResult `json:"series"`
}

// enterResult saves a result line into the shooter's entry for the programme when it has no results yet, or into a
// new registration without a lane; 409 when the entry already has results.
func (s *Server) enterResult(ctx context.Context, tenantID, eventID int64, line ResultLine) (ResultSheet, error) {
	e, done, err := s.Store.FindEntry(ctx, tenantID, eventID, line.ShooterID, line.ProgrammeID)
	switch {
	case err == nil && done:
		return ResultSheet{}, conflict("The shooter already has results in this programme")
	case errors.Is(err, store.ErrNotFound):
		p, err := s.Store.GetProgramme(ctx, tenantID, line.ProgrammeID)
		if err != nil || p.EventID != eventID {
			return ResultSheet{}, notFound(store.ErrNotFound, "Programme not found")
		}
		domain.NormalizeResults(line.Series)
		if errs := domain.ValidateResults(p.Series, line.Series); len(errs) > 0 {
			return ResultSheet{}, invalidList(errs)
		}
		if _, err := s.register(ctx, tenantID, eventID, RegistrationInput{ShooterID: line.ShooterID, ClubID: line.ClubID,
			Programmes: []domain.ProgrammeChoice{{ProgrammeID: line.ProgrammeID}}}); err != nil {
			return ResultSheet{}, err
		}
		if e, _, err = s.Store.FindEntry(ctx, tenantID, eventID, line.ShooterID, line.ProgrammeID); err != nil {
			return ResultSheet{}, err
		}
	case err != nil:
		return ResultSheet{}, err
	}
	return s.saveResults(ctx, tenantID, e, line.Series)
}

// saveResults replaces the results of an entry, keeps or picks the default award and tells the event's clients;
// it returns the result sheet of the registration.
func (s *Server) saveResults(ctx context.Context, tenantID int64, e store.EntryRef, series []domain.SeriesResult) (ResultSheet, error) {
	if len(e.CombinedSeries) > 0 {
		return ResultSheet{}, conflict("This programme takes its results from another programme")
	}
	p, err := s.Store.GetProgramme(ctx, tenantID, e.ProgrammeID)
	if err != nil {
		return ResultSheet{}, err
	}
	domain.NormalizeResults(series)
	if errs := domain.ValidateResults(p.Series, series); len(errs) > 0 {
		return ResultSheet{}, invalidList(errs)
	}
	before, err := s.resultSheet(ctx, tenantID, e.RegistrationID)
	if err != nil {
		return ResultSheet{}, err
	}
	completed, err := s.Store.SaveResults(ctx, tenantID, e, series)
	if err != nil {
		return ResultSheet{}, err
	}
	if err := s.defaultAwards(ctx, tenantID, e, before); err != nil {
		return ResultSheet{}, err
	}
	topics := []string{realtime.TopicResults, realtime.TopicRegistrations}
	if completed {
		topics = append(topics, realtime.TopicLanes)
	}
	s.notify(ctx, tenantID, e.EventID, topics...)
	return s.resultSheet(ctx, tenantID, e.RegistrationID)
}

// defaultAwards gives the saved entry and the entries that take series from it an award again: the one
// they had before if it is still reached, otherwise the most valuable one.
func (s *Server) defaultAwards(ctx context.Context, tenantID int64, e store.EntryRef, before ResultSheet) error {
	prior := map[int64]int64{}
	for _, er := range before.Entries {
		if er.Award != nil {
			prior[er.EntryID] = er.Award.ID
		}
	}
	sheet, err := s.resultSheet(ctx, tenantID, e.RegistrationID)
	if err != nil {
		return err
	}
	for _, er := range sheet.Entries {
		if (er.EntryID != e.ID && len(er.CombinedSeries) == 0) || er.Award != nil || len(er.EligibleAwards) == 0 {
			continue
		}
		tier := er.EligibleAwards[0].ID
		for _, t := range er.EligibleAwards {
			if t.ID == prior[er.EntryID] {
				tier = t.ID
			}
		}
		if err := s.Store.SetAward(ctx, tenantID, er.EntryID, &tier); err != nil {
			return err
		}
	}
	return nil
}

// registerResults registers the result and award endpoints.
func (s *Server) registerResults(api huma.API) {
	type sheetOut struct{ Body ResultSheet }
	huma.Register(api, huma.Operation{
		OperationID: "get-results", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/registrations/{registrationId}/results", Tags: []string{"results"},
	}, func(ctx context.Context, in *RegistrationPath) (*sheetOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		sheet, err := s.resultSheet(ctx, in.TenantID, in.RegistrationID)
		return &sheetOut{Body: sheet}, err
	})

	type saveIn struct {
		EntryPath
		Body struct {
			Series []domain.SeriesResult `json:"series"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "save-results", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/entries/{entryId}/results", Tags: []string{"results"},
		Summary:     "Replace the results of a programme and choose its award",
		Description: "Keeps the award of the programme and of the programmes that take series from it while it is still reached, otherwise gives the most valuable award reached.",
	}, func(ctx context.Context, in *saveIn) (*sheetOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		e, err := s.Store.GetEntry(ctx, in.TenantID, in.EntryID)
		if err != nil {
			return nil, notFound(err, "Programme of the registration not found")
		}
		sheet, err := s.saveResults(ctx, in.TenantID, e, in.Body.Series)
		return &sheetOut{Body: sheet}, err
	})

	type enterIn struct {
		EventPath
		Body ResultLine
	}
	huma.Register(api, huma.Operation{
		OperationID: "enter-result", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/results", Tags: []string{"results"},
		Summary:     "Enter the results of a shooter in a programme without the desk",
		Description: "Fills the shooter's entry for the programme when it has no results yet, otherwise registers the shooter without a lane for that programme. 409 when the entry already has results.",
	}, func(ctx context.Context, in *enterIn) (*sheetOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		sheet, err := s.enterResult(ctx, in.TenantID, in.EventID, in.Body)
		return &sheetOut{Body: sheet}, err
	})

	type awardIn struct {
		EntryPath
		Body struct {
			TierID *int64 `json:"tierId,omitempty" required:"false" doc:"One of the eligible awards, or omitted for none"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-award", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/entries/{entryId}/award", Tags: []string{"results"},
	}, func(ctx context.Context, in *awardIn) (*sheetOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		e, err := s.Store.GetEntry(ctx, in.TenantID, in.EntryID)
		if err != nil {
			return nil, notFound(err, "Programme of the registration not found")
		}
		if in.Body.TierID != nil {
			sheet, err := s.resultSheet(ctx, in.TenantID, e.RegistrationID)
			if err != nil {
				return nil, err
			}
			ok := false
			for _, er := range sheet.Entries {
				if er.EntryID != e.ID {
					continue
				}
				for _, t := range er.EligibleAwards {
					ok = ok || t.ID == *in.Body.TierID
				}
			}
			if !ok {
				return nil, invalid("This award has not been reached")
			}
		}
		if err := s.Store.SetAward(ctx, in.TenantID, e.ID, in.Body.TierID); err != nil {
			return nil, err
		}
		s.notify(ctx, in.TenantID, e.EventID, realtime.TopicResults)
		sheet, err := s.resultSheet(ctx, in.TenantID, e.RegistrationID)
		return &sheetOut{Body: sheet}, err
	})
}
