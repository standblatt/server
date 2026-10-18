package web

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
)

// RankingLine is one ranked result.
type RankingLine struct {
	Rank      int                `json:"rank"`
	EntryID   int64              `json:"entryId"`
	Firstname string             `json:"firstname"`
	Lastname  string             `json:"lastname"`
	BirthYear string             `json:"birthYear"`
	Category  domain.AgeCategory `json:"category" enum:"u10,u13,u15,u17,u19,u21,elite,senior,veteran,senior_veteran"`
	Club      string             `json:"club"`
	Series    []float64          `json:"series"`
	Total     float64            `json:"total"`
	Mouches   int                `json:"mouches"`
	Award     string             `json:"award,omitempty"`
	Attempt   int                `json:"attempt" doc:"The counting attempt"`
	Attempts  int                `json:"attempts" doc:"Completed attempts of the shooter"`
}

// ProgrammeRanking is the ranking of one programme.
type ProgrammeRanking struct {
	ProgrammeID   int64         `json:"programmeId"`
	Title         string        `json:"title"`
	Shots         int           `json:"shots"`
	MaxPoints     int           `json:"maxPoints"`
	SeriesCount   int           `json:"seriesCount"`
	TracksMouches bool          `json:"tracksMouches"`
	Public        bool          `json:"public" doc:"Shown on kiosk displays (the programme's showOnRankings)"`
	RangeID       *int64        `json:"rangeId,omitempty"`
	RangeName     string        `json:"rangeName,omitempty" doc:"Range the programme is shot on, when tied to one"`
	Repeatable    bool          `json:"repeatable" doc:"Shooters may shoot several attempts; each line shows its counting one"`
	Lines         []RankingLine `json:"lines"`
}

// KioskRange is a range of the kiosk's event.
type KioskRange struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// KioskEvent is what a kiosk display shows in its header.
type KioskEvent struct {
	Name           string       `json:"name"`
	StartsOn       string       `json:"startsOn"`
	Days           []string     `json:"days"`
	Ranges         []KioskRange `json:"ranges" doc:"Ranges in order; ?range=<id> on the kiosk link limits the display to one"`
	Organizer      string       `json:"organizer"`
	Locale         string       `json:"locale"`
	ShowLaneStatus bool         `json:"showLaneStatus"`
}

// rankings ranks every programme of an event; publicOnly leaves out programmes not shown on rankings (kiosk).
func (s *Server) rankings(ctx context.Context, tenantID, eventID int64, publicOnly bool) ([]ProgrammeRanking, error) {
	programmes, err := s.Store.ListProgrammes(ctx, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	ev, err := s.Store.GetEvent(ctx, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	rangeNames := map[int64]string{}
	for _, r := range ev.Ranges {
		rangeNames[r.ID] = r.Name
	}
	rows, err := s.Store.ListRankingRows(ctx, tenantID, eventID)
	if err != nil {
		return nil, err
	}
	byProgramme := map[int64][]store.RankingRow{}
	for _, r := range store.CountingRows(rows) {
		byProgramme[r.ProgrammeID] = append(byProgramme[r.ProgrammeID], r)
	}
	out := []ProgrammeRanking{}
	for _, p := range programmes {
		if publicOnly && !p.ShowOnRankings {
			continue
		}
		pr := ProgrammeRanking{
			ProgrammeID: p.ID, Title: programmeTitle(p.Name, p.Variant), Shots: p.Shots, MaxPoints: p.MaxPoints,
			SeriesCount: len(p.Series), Public: p.ShowOnRankings, RangeID: p.RangeID, Repeatable: p.AttemptLimit() != 1, Lines: []RankingLine{},
		}
		if p.RangeID != nil {
			pr.RangeName = rangeNames[*p.RangeID]
		}
		for _, sr := range p.Series {
			pr.TracksMouches = pr.TracksMouches || sr.TrackMouches
		}
		list := byProgramme[p.ID]
		entries := make([]domain.RankingEntry, len(list))
		for i, r := range list {
			entries[i] = r.Entry()
		}
		for rank, i := range domain.Rank(entries) {
			r := list[i]
			line := RankingLine{
				Rank: rank + 1, EntryID: r.EntryID, Firstname: r.Firstname, Lastname: r.Lastname, Category: r.Category,
				Club: r.ClubName, Series: r.Series, Total: entries[i].Total, Mouches: r.Mouches, Award: r.AwardName,
				Attempt: r.Attempt, Attempts: r.Attempts,
			}
			if len(r.Birthday) >= 4 {
				line.BirthYear = r.Birthday[:4]
			}
			pr.Lines = append(pr.Lines, line)
		}
		out = append(out, pr)
	}
	return out, nil
}

// kioskAuth returns the tenant and event of a kiosk token from the Bearer header or the query, 401 when it is missing or unknown.
func (s *Server) kioskAuth(ctx context.Context, header, query string) (int64, int64, error) {
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok || token == "" {
		token = query
	}
	if token == "" {
		return 0, 0, unauthorized("Kiosk token missing")
	}
	tenantID, eventID, err := s.Store.KioskEvent(ctx, sessionID(token))
	if errors.Is(err, store.ErrNotFound) {
		return 0, 0, unauthorized("Unknown or revoked kiosk token")
	}
	return tenantID, eventID, err
}

// KioskAuth carries the kiosk token of a request.
type KioskAuth struct {
	Authorization string `header:"Authorization" doc:"Bearer <kiosk token>"`
}

// registerRankings registers the ranking, kiosk token and kiosk endpoints.
func (s *Server) registerRankings(api huma.API) {
	type rankingsOut struct{ Body []ProgrammeRanking }
	huma.Register(api, huma.Operation{
		OperationID: "get-rankings", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/rankings", Tags: []string{"rankings"},
	}, func(ctx context.Context, in *EventPath) (*rankingsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		list, err := s.rankings(ctx, in.TenantID, in.EventID, false)
		return &rankingsOut{Body: list}, err
	})

	type tokensOut struct{ Body []store.KioskToken }
	huma.Register(api, huma.Operation{
		OperationID: "list-kiosk-tokens", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/kiosk-tokens", Tags: []string{"kiosk"},
	}, func(ctx context.Context, in *EventPath) (*tokensOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		list, err := s.Store.ListKioskTokens(ctx, in.TenantID, in.EventID)
		return &tokensOut{Body: list}, err
	})

	type createIn struct {
		EventPath
		Body struct {
			Name string `json:"name" minLength:"1" maxLength:"100" doc:"Where the display hangs, e.g. Hall 1"`
		}
	}
	type createOut struct {
		Body struct {
			Token  store.KioskToken `json:"token"`
			Secret string           `json:"secret" doc:"Shown once; put it in the kiosk URL /kiosk/<secret>"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-kiosk-token", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/kiosk-tokens", Tags: []string{"kiosk"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createIn) (*createOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		secret := randomToken()
		name := strings.TrimSpace(in.Body.Name)
		id, err := s.Store.CreateKioskToken(ctx, in.TenantID, in.EventID, name, sessionID(secret))
		if err != nil {
			return nil, notFound(err, "Event not found")
		}
		out := &createOut{}
		out.Body.Token = store.KioskToken{ID: id, Name: name}
		out.Body.Secret = secret
		return out, nil
	})

	type tokenPath struct {
		TenantID int64 `path:"tenantId"`
		TokenID  int64 `path:"tokenId"`
	}
	huma.Register(api, huma.Operation{
		OperationID: "delete-kiosk-token", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/kiosk-tokens/{tokenId}", Tags: []string{"kiosk"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *tokenPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.DeleteKioskToken(ctx, in.TenantID, in.TokenID), "Kiosk token not found")
	})

	type kioskEventOut struct{ Body KioskEvent }
	huma.Register(api, huma.Operation{
		OperationID: "kiosk-event", Method: http.MethodGet, Path: "/api/kiosk/event", Tags: []string{"kiosk"},
	}, func(ctx context.Context, in *KioskAuth) (*kioskEventOut, error) {
		tenantID, eventID, err := s.kioskAuth(ctx, in.Authorization, "")
		if err != nil {
			return nil, err
		}
		ev, err := s.Store.GetEvent(ctx, tenantID, eventID)
		if err != nil {
			return nil, err
		}
		ranges := make([]KioskRange, len(ev.Ranges))
		for i, r := range ev.Ranges {
			ranges[i] = KioskRange{ID: r.ID, Name: r.Name}
		}
		return &kioskEventOut{Body: KioskEvent{
			Name: ev.Name, StartsOn: ev.StartsOn, Days: ev.Days, Ranges: ranges, Organizer: ev.Organizer, Locale: ev.Locale, ShowLaneStatus: ev.Settings.ShowLaneStatus,
		}}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "kiosk-rankings", Method: http.MethodGet, Path: "/api/kiosk/rankings", Tags: []string{"kiosk"},
	}, func(ctx context.Context, in *KioskAuth) (*rankingsOut, error) {
		tenantID, eventID, err := s.kioskAuth(ctx, in.Authorization, "")
		if err != nil {
			return nil, err
		}
		list, err := s.rankings(ctx, tenantID, eventID, true)
		return &rankingsOut{Body: list}, err
	})

	type lanesOut struct{ Body []RangeLanes }
	huma.Register(api, huma.Operation{
		OperationID: "kiosk-lanes", Method: http.MethodGet, Path: "/api/kiosk/lanes", Tags: []string{"kiosk"},
		Summary: "Lane states, empty unless the event shows the lane status on rankings",
	}, func(ctx context.Context, in *KioskAuth) (*lanesOut, error) {
		tenantID, eventID, err := s.kioskAuth(ctx, in.Authorization, "")
		if err != nil {
			return nil, err
		}
		ev, err := s.Store.GetEvent(ctx, tenantID, eventID)
		if err != nil {
			return nil, err
		}
		if !ev.Settings.ShowLaneStatus {
			return &lanesOut{Body: []RangeLanes{}}, nil
		}
		lanes, err := s.laneStates(ctx, tenantID, eventID)
		return &lanesOut{Body: lanes}, err
	})
}

// kioskLive is the WebSocket feed for kiosk displays; the token comes as ?token= since browsers cannot set headers.
func (s *Server) kioskLive(w http.ResponseWriter, r *http.Request) {
	tenantID, eventID, err := s.kioskAuth(r.Context(), r.Header.Get("Authorization"), r.URL.Query().Get("token"))
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	s.serveLive(w, r, tenantID, eventID)
}
