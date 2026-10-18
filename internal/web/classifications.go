package web

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

// ClassificationPartView is a programme of a classification as shown in its ranking.
type ClassificationPartView struct {
	ProgrammeID int64  `json:"programmeId"`
	Title       string `json:"title"`
	Count       int    `json:"count"`
}

// PartScore is a shooter's points in one part; Points is left out without a counting attempt.
type PartScore struct {
	Points  *float64 `json:"points,omitempty"`
	Counted bool     `json:"counted" doc:"The part counts towards the total (see best)"`
}

// ClassificationLine is a shooter in a classification ranking; Rank is 0 for unranked shooters.
type ClassificationLine struct {
	Rank      int                `json:"rank"`
	ShooterID int64              `json:"shooterId"`
	Firstname string             `json:"firstname"`
	Lastname  string             `json:"lastname"`
	BirthYear string             `json:"birthYear"`
	Category  domain.AgeCategory `json:"category" enum:"u10,u13,u15,u17,u19,u21,elite,senior,veteran,senior_veteran"`
	Club      string             `json:"club"`
	Parts     []PartScore        `json:"parts"`
	Total     float64            `json:"total"`
	Complete  bool               `json:"complete"`
}

// ClassificationClassView is the ranking of one class; Name is empty when the classification has no classes.
type ClassificationClassView struct {
	Name     string               `json:"name"`
	Lines    []ClassificationLine `json:"lines"`
	Unranked []ClassificationLine `json:"unranked" doc:"Shooters without every counting part, when the classification requires it"`
}

// ClassificationRanking is the ranking of a classification.
type ClassificationRanking struct {
	ID      int64                     `json:"id"`
	Name    string                    `json:"name"`
	Public  bool                      `json:"public" doc:"Shown on kiosk displays"`
	Parts   []ClassificationPartView  `json:"parts"`
	Classes []ClassificationClassView `json:"classes"`
}

// classificationRankings ranks every classification of an event; publicOnly leaves out those not shown on kiosks.
func (s *Server) classificationRankings(ctx context.Context, tenantID, eventID int64, publicOnly bool) ([]ClassificationRanking, error) {
	list, err := s.Store.ListClassifications(ctx, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	out := []ClassificationRanking{}
	if len(list) == 0 {
		return out, nil
	}
	programmes, err := s.Store.ListProgrammes(ctx, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	titles := map[int64]string{}
	for _, p := range programmes {
		titles[p.ID] = programmeTitle(p.Name, p.Variant)
	}
	rows, err := s.Store.ListRankingRows(ctx, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	results := make([]domain.AttemptResult, len(rows))
	people := map[int64]store.RankingRow{}
	for i, r := range rows {
		results[i] = domain.AttemptResult{ShooterID: r.ShooterID, ProgrammeID: r.ProgrammeID, Rule: r.Rule, Entry: r.Entry()}
		people[r.ShooterID] = r
	}
	line := func(st domain.Standing) ClassificationLine {
		p := people[st.ShooterID]
		l := ClassificationLine{Rank: st.Rank, BirthYear: birthYear(p.Birthday), ShooterID: st.ShooterID, Firstname: p.Firstname, Lastname: p.Lastname, Category: st.Entry.Category,
			Club: p.ClubName, Total: st.Entry.Total, Complete: st.Complete, Parts: make([]PartScore, len(st.Parts))}
		for i := range st.Parts {
			l.Parts[i] = PartScore{Points: st.Parts[i], Counted: st.Counted[i]}
		}
		return l
	}
	for _, c := range list {
		if publicOnly && !c.ShowOnKiosk {
			continue
		}
		cr := ClassificationRanking{ID: c.ID, Name: c.Name, Public: c.ShowOnKiosk, Parts: []ClassificationPartView{}, Classes: []ClassificationClassView{}}
		for _, p := range c.Parts {
			cr.Parts = append(cr.Parts, ClassificationPartView{ProgrammeID: p.ProgrammeID, Title: titles[p.ProgrammeID], Count: p.Count})
		}
		for _, cl := range domain.Classify(c.ClassificationDefinition, results) {
			v := ClassificationClassView{Name: cl.Name, Lines: []ClassificationLine{}, Unranked: []ClassificationLine{}}
			for _, st := range cl.Ranked {
				v.Lines = append(v.Lines, line(st))
			}
			for _, st := range cl.Unranked {
				v.Unranked = append(v.Unranked, line(st))
			}
			cr.Classes = append(cr.Classes, v)
		}
		out = append(out, cr)
	}
	return out, nil
}

// ClassificationPath identifies a classification of an event of a tenant in the URL.
type ClassificationPath struct {
	TenantID         int64 `path:"tenantId"`
	EventID          int64 `path:"eventId"`
	ClassificationID int64 `path:"classificationId"`
}

// checkClassification normalizes a classification and returns its problems as one 422.
func (s *Server) checkClassification(ctx context.Context, tenantID, eventID int64, c *domain.ClassificationDefinition) error {
	c.Normalize()
	programmes, err := s.Store.ListProgrammes(ctx, tenantID, eventID)
	if err != nil {
		return err
	}
	ids := map[int64]bool{}
	for _, p := range programmes {
		ids[p.ID] = true
	}
	if errs := c.Validate(ids); len(errs) > 0 {
		return invalidList(errs)
	}
	return nil
}

// registerClassifications registers the classification and classification ranking endpoints.
func (s *Server) registerClassifications(api huma.API) {
	type listOut struct{ Body []store.Classification }
	huma.Register(api, huma.Operation{
		OperationID: "list-classifications", Metadata: needs(modules.Classifications, permissions.EventsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/classifications", Tags: []string{"classifications"},
	}, func(ctx context.Context, in *EventPath) (*listOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.Store.ListClassifications(ctx, in.TenantID, in.EventID)
		return &listOut{Body: list}, err
	})

	type createIn struct {
		EventPath
		Body domain.ClassificationDefinition
	}
	type createOut struct{ Body store.Classification }
	huma.Register(api, huma.Operation{
		OperationID: "create-classification", Metadata: needs(modules.Classifications, permissions.EventsConfigure), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/classifications", Tags: []string{"classifications"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createIn) (*createOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		if err := s.checkClassification(ctx, in.TenantID, in.EventID, &in.Body); err != nil {
			return nil, err
		}
		id, err := s.Store.CreateClassification(ctx, in.TenantID, in.EventID, in.Body)
		if err != nil {
			return nil, notFound(err, "Event not found")
		}
		s.notify(ctx, in.TenantID, in.EventID, realtime.TopicProgrammes)
		return &createOut{Body: store.Classification{ID: id, ClassificationDefinition: in.Body}}, nil
	})

	type updateIn struct {
		ClassificationPath
		Body domain.ClassificationDefinition
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-classification", Metadata: needs(modules.Classifications, permissions.EventsConfigure), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/events/{eventId}/classifications/{classificationId}", Tags: []string{"classifications"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *updateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		if err := s.checkClassification(ctx, in.TenantID, in.EventID, &in.Body); err != nil {
			return nil, err
		}
		if err := s.Store.UpdateClassification(ctx, in.TenantID, in.EventID, in.ClassificationID, in.Body); err != nil {
			return nil, notFound(err, "Classification not found")
		}
		s.notify(ctx, in.TenantID, in.EventID, realtime.TopicProgrammes)
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-classification", Metadata: needs(modules.Classifications, permissions.EventsConfigure), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/events/{eventId}/classifications/{classificationId}", Tags: []string{"classifications"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *ClassificationPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		if err := s.Store.DeleteClassification(ctx, in.TenantID, in.EventID, in.ClassificationID); err != nil {
			return nil, notFound(err, "Classification not found")
		}
		s.notify(ctx, in.TenantID, in.EventID, realtime.TopicProgrammes)
		return nil, nil
	})

	type rankingsOut struct{ Body []ClassificationRanking }
	huma.Register(api, huma.Operation{
		OperationID: "get-classification-rankings", Metadata: needs(modules.Classifications, permissions.EventsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/classification-rankings", Tags: []string{"classifications"},
	}, func(ctx context.Context, in *EventPath) (*rankingsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.classificationRankings(ctx, in.TenantID, in.EventID, false)
		return &rankingsOut{Body: list}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "kiosk-classifications", Method: http.MethodGet, Path: "/api/kiosk/classifications", Tags: []string{"kiosk"},
	}, func(ctx context.Context, in *KioskAuth) (*rankingsOut, error) {
		tenantID, eventID, err := s.kioskAuth(ctx, in.Authorization, "")
		if err != nil {
			return nil, err
		}
		list, err := s.classificationRankings(ctx, tenantID, eventID, true)
		return &rankingsOut{Body: list}, err
	})
}
