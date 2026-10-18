package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

// RegistrationPath identifies a registration of a tenant in the URL.
type RegistrationPath struct {
	TenantID       int64 `path:"tenantId"`
	RegistrationID int64 `path:"registrationId"`
}

// RegistrationInput registers a shooter for an event.
type RegistrationInput struct {
	ShooterID  int64                    `json:"shooterId"`
	ClubID     *int64                   `json:"clubId,omitempty" required:"false" doc:"One of the shooter's clubs"`
	RangeID    *int64                   `json:"rangeId,omitempty" required:"false" doc:"Range of the lane; may be left out when the event has one range"`
	Lane       *int                     `json:"lane,omitempty" required:"false"`
	Programmes []domain.ProgrammeChoice `json:"programmes"`
}

// RegistrationUpdate changes range, lane and programmes of a registration.
type RegistrationUpdate struct {
	RangeID    *int64                   `json:"rangeId,omitempty" required:"false"`
	Lane       *int                     `json:"lane,omitempty" required:"false"`
	Programmes []domain.ProgrammeChoice `json:"programmes"`
}

// RangeLanes are the lanes of one range of an event.
type RangeLanes struct {
	RangeID   int64       `json:"rangeId"`
	RangeName string      `json:"rangeName"`
	Lanes     []LaneState `json:"lanes"`
}

// LaneState is the state of one lane of an event's range. Out of service beats occupied beats reserved.
type LaneState struct {
	Lane           int        `json:"lane"`
	State          string     `json:"state" enum:"free,disabled,occupied,reserved"`
	RegistrationID int64      `json:"registrationId,omitempty"`
	ShooterName    string     `json:"shooterName,omitempty"`
	Since          *time.Time `json:"since,omitempty"`
	Reservation    string     `json:"reservation,omitempty" doc:"Description of the reservation that holds the lane now"`
	ReservedUntil  *time.Time `json:"reservedUntil,omitempty"`
}

// notify announces a change of topics to the live clients of an event and only logs failures.
func (s *Server) notify(ctx context.Context, tenantID, eventID int64, topics ...string) {
	if err := s.Store.Notify(ctx, realtime.Channel, tenantID, eventID, topics...); err != nil {
		s.Log.Warn("notify", "err", err)
	}
}

// checkLane validates a lane choice and returns the range it belongs to: the given one, or the event's only range.
func checkLane(rangeID *int64, lane *int, rules store.EventRules) (*int64, error) {
	if lane == nil {
		return nil, nil
	}
	if !rules.Lanes {
		return nil, invalid("Lanes are not assigned at this event")
	}
	if rangeID == nil && len(rules.RangeOrder) == 1 {
		rangeID = &rules.RangeOrder[0]
	}
	if rangeID == nil {
		return nil, invalid("Choose the range of the lane")
	}
	rr, ok := rules.Ranges[*rangeID]
	if !ok {
		return nil, invalid("Range not found")
	}
	if *lane < 1 || *lane > rr.LaneCount {
		return nil, invalid("Lane %d does not exist (1–%d)", *lane, rr.LaneCount)
	}
	if slices.Contains(rr.DisabledLanes, int32(*lane)) {
		return nil, invalid("Lane %d is out of service", *lane)
	}
	return rangeID, nil
}

// register checks and stores a new registration: the club must be one of the shooter's, the age category follows
// the event date, the lane must exist and be free and the programmes must suit the shooter and the lane's range.
func (s *Server) register(ctx context.Context, tenantID, eventID int64, in RegistrationInput) (int64, error) {
	rules, err := s.Store.LoadEventRules(ctx, tenantID, eventID)
	if err != nil {
		return 0, notFound(err, "Event not found")
	}
	if rules.Template {
		return 0, conflict("A template takes no registrations")
	}
	sh, err := s.Store.GetShooter(ctx, tenantID, in.ShooterID)
	if err != nil {
		return 0, notFound(err, "Shooter not found")
	}
	if in.ClubID != nil && !slices.ContainsFunc(sh.Clubs, func(c store.ClubRef) bool { return c.ID == *in.ClubID }) {
		return 0, invalid("The shooter is not a member of this club")
	}
	category := domain.AgeCategoryFor(sh.Birthday, rules.StartsOn)
	rangeID, err := checkLane(in.RangeID, in.Lane, rules)
	if err != nil {
		return 0, err
	}
	if errs := domain.CheckProgrammes(in.Programmes, rules.Programmes, category, nil, rangeID); len(errs) > 0 {
		return 0, invalidList(errs)
	}
	id, err := s.Store.CreateRegistration(ctx, tenantID, store.NewRegistration{
		EventID: eventID, ShooterID: sh.ID, ClubID: in.ClubID, AgeCategory: category,
		RangeID: rangeID, Lane: in.Lane, Programmes: in.Programmes,
	})
	if err != nil {
		return 0, registrationConflict(err, in.Lane)
	}
	s.notify(ctx, tenantID, eventID, realtime.TopicRegistrations, realtime.TopicLanes)
	return id, nil
}

// registrationConflict maps store errors of registration writes to 409 or 404; lane must be set for ErrLaneTaken.
func registrationConflict(err error, lane *int) error {
	switch {
	case errors.Is(err, store.ErrLaneTaken) && lane != nil:
		return conflict("Lane %d is occupied", *lane)
	case errors.Is(err, store.ErrLaneTaken):
		return conflict("The lane is occupied")
	case errors.Is(err, store.ErrProgrammeTaken):
		return conflict("The shooter is already registered for one of these programmes at this event and has no attempts left")
	case errors.Is(err, store.ErrHasResults):
		return conflict("A programme with results cannot be removed from the registration")
	}
	return notFound(err, "Registration not found")
}

// registerRegistrations registers the registration, lane and event stats endpoints.
func (s *Server) registerRegistrations(api huma.API) {
	type listIn struct {
		EventPath
		All bool `query:"all" doc:"Include completed registrations"`
	}
	type listOut struct{ Body []store.Registration }
	huma.Register(api, huma.Operation{
		OperationID: "list-registrations", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/registrations", Tags: []string{"registrations"},
	}, func(ctx context.Context, in *listIn) (*listOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, permissions.EventsOperate); err != nil {
			return nil, err
		}
		list, err := s.Store.ListRegistrations(ctx, in.TenantID, in.EventID, !in.All)
		return &listOut{Body: list}, err
	})

	type regOut struct{ Body store.Registration }
	type createIn struct {
		EventPath
		Body RegistrationInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-registration", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/registrations", Tags: []string{"registrations"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createIn) (*regOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, permissions.EventsOperate); err != nil {
			return nil, err
		}
		id, err := s.register(ctx, in.TenantID, in.EventID, in.Body)
		if err != nil {
			return nil, err
		}
		r, err := s.Store.GetRegistration(ctx, in.TenantID, id)
		return &regOut{Body: r}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-registration", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/registrations/{registrationId}", Tags: []string{"registrations"},
	}, func(ctx context.Context, in *RegistrationPath) (*regOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, permissions.EventsOperate); err != nil {
			return nil, err
		}
		r, err := s.Store.GetRegistration(ctx, in.TenantID, in.RegistrationID)
		return &regOut{Body: r}, notFound(err, "Registration not found")
	})

	type updateIn struct {
		RegistrationPath
		Body RegistrationUpdate
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-registration", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/registrations/{registrationId}", Tags: []string{"registrations"},
	}, func(ctx context.Context, in *updateIn) (*regOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, permissions.EventsOperate); err != nil {
			return nil, err
		}
		current, err := s.Store.GetRegistration(ctx, in.TenantID, in.RegistrationID)
		if err != nil {
			return nil, notFound(err, "Registration not found")
		}
		if helperOf(ctx) != nil && !sameProgrammes(current.Entries, in.Body.Programmes) {
			return nil, forbidden("Not allowed for this helper link")
		}
		rules, err := s.Store.LoadEventRules(ctx, in.TenantID, current.EventID)
		if err != nil {
			return nil, err
		}
		keep := map[int64]bool{}
		for _, e := range current.Entries {
			keep[e.ProgrammeID] = true
		}
		rangeID, err := checkLane(in.Body.RangeID, in.Body.Lane, rules)
		if err != nil {
			return nil, err
		}
		if errs := domain.CheckProgrammes(in.Body.Programmes, rules.Programmes, current.AgeCategory, keep, rangeID); len(errs) > 0 {
			return nil, invalidList(errs)
		}
		if err := s.Store.UpdateRegistration(ctx, in.TenantID, current, rangeID, in.Body.Lane, in.Body.Programmes); err != nil {
			return nil, registrationConflict(err, in.Body.Lane)
		}
		s.notify(ctx, in.TenantID, current.EventID, realtime.TopicRegistrations, realtime.TopicLanes)
		r, err := s.Store.GetRegistration(ctx, in.TenantID, current.ID)
		return &regOut{Body: r}, err
	})

	change := func(op, method, suffix, module string, run func(ctx context.Context, tenantID, id int64) error) {
		huma.Register(api, huma.Operation{
			OperationID: op, Metadata: modules.Metadata(module), Method: method, Path: "/api/tenants/{tenantId}/registrations/{registrationId}" + suffix, Tags: []string{"registrations"},
			DefaultStatus: http.StatusNoContent,
		}, func(ctx context.Context, in *RegistrationPath) (*struct{}, error) {
			if _, _, err := s.access(ctx, in.TenantID, permissions.EventsOperate); err != nil {
				return nil, err
			}
			current, err := s.Store.GetRegistration(ctx, in.TenantID, in.RegistrationID)
			if err != nil {
				return nil, notFound(err, "Registration not found")
			}
			if err := run(ctx, in.TenantID, current.ID); err != nil {
				return nil, notFound(err, "Registration not found")
			}
			s.notify(ctx, in.TenantID, current.EventID, realtime.TopicRegistrations, realtime.TopicLanes)
			return nil, nil
		})
	}
	change("finish-registration", http.MethodPost, "/finish", "", s.Store.FinishRegistration)
	change("delete-registration", http.MethodDelete, "", "", s.Store.DeleteRegistration)
	change("issue-ammo", http.MethodPost, "/ammo", modules.Ammo, s.Store.IssueAmmo)
	change("issue-awards", http.MethodPost, "/awards-issued", modules.Awards, s.Store.IssueAwards)
	change("revoke-ammo", http.MethodDelete, "/ammo", modules.Ammo, s.Store.RevokeAmmo)
	change("revoke-awards", http.MethodDelete, "/awards-issued", modules.Awards, s.Store.RevokeAwards)

	type lanesOut struct{ Body []RangeLanes }
	huma.Register(api, huma.Operation{
		OperationID: "list-lanes", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/lanes", Tags: []string{"registrations"},
	}, func(ctx context.Context, in *EventPath) (*lanesOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, permissions.EventsOperate); err != nil {
			return nil, err
		}
		lanes, err := s.laneStates(ctx, in.TenantID, in.EventID)
		return &lanesOut{Body: lanes}, err
	})

	type statsOut struct{ Body store.EventStats }
	huma.Register(api, huma.Operation{
		OperationID: "get-event-stats", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/stats", Tags: []string{"registrations"},
	}, func(ctx context.Context, in *EventPath) (*statsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, permissions.EventsOperate); err != nil {
			return nil, err
		}
		st, err := s.Store.GetEventStats(ctx, in.TenantID, in.EventID)
		return &statsOut{Body: st}, err
	})
}

// laneStates returns the lanes of every range of the event as free, reserved, occupied or disabled.
func (s *Server) laneStates(ctx context.Context, tenantID, eventID int64) ([]RangeLanes, error) {
	rules, err := s.Store.LoadEventRules(ctx, tenantID, eventID)
	if err != nil {
		return nil, notFound(err, "Event not found")
	}
	out := []RangeLanes{}
	if !rules.Lanes {
		return out, nil
	}
	regs, err := s.Store.ListRegistrations(ctx, tenantID, eventID, true)
	if err != nil {
		return nil, err
	}
	for _, rangeID := range rules.RangeOrder {
		rr := rules.Ranges[rangeID]
		reservations, err := s.Store.ActiveReservations(ctx, rangeID, time.Now())
		if err != nil {
			return nil, err
		}
		lanes := make([]LaneState, rr.LaneCount)
		for i := range lanes {
			lanes[i] = LaneState{Lane: i + 1, State: "free"}
			for _, res := range reservations {
				if slices.Contains(res.Lanes, i+1) && lanes[i].Reservation == "" {
					until := res.EndsAt
					lanes[i].State, lanes[i].Reservation, lanes[i].ReservedUntil = "reserved", res.Description, &until
				}
			}
			if slices.Contains(rr.DisabledLanes, int32(i+1)) {
				// out of service wins over a reservation; the reservation text would only confuse the grid
				lanes[i] = LaneState{Lane: i + 1, State: "disabled"}
			}
		}
		for _, r := range regs {
			if r.Lane == nil || r.RangeID == nil || *r.RangeID != rangeID || *r.Lane > len(lanes) {
				continue
			}
			l := &lanes[*r.Lane-1]
			if l.State == "disabled" {
				continue
			}
			l.State, l.RegistrationID, l.Since = "occupied", r.ID, r.StartedAt
			l.ShooterName = r.Shooter.Lastname + " " + r.Shooter.Firstname
		}
		out = append(out, RangeLanes{RangeID: rangeID, RangeName: rr.Name, Lanes: lanes})
	}
	return out, nil
}

// live streams change notifications of one event over a WebSocket: {"topic": "..."} per change.
func (s *Server) live(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := strconv.ParseInt(r.PathValue("tenantId"), 10, 64)
	eventID, _ := strconv.ParseInt(r.PathValue("eventId"), 10, 64)
	ctx := r.Context()
	_, _, err := s.access(ctx, tenantID, permissions.EventsView)
	if h := helperOf(ctx); err == nil && h != nil {
		err = s.helperAllowed(ctx, h, "live", func(string) string { return "" })
		if err == nil && h.EventID != eventID {
			err = missing("Event not found")
		}
	}
	if err != nil {
		var se huma.StatusError
		if errors.As(err, &se) {
			http.Error(w, se.Error(), se.GetStatus())
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if _, err := s.Store.GetEvent(r.Context(), tenantID, eventID); err != nil {
		http.NotFound(w, r)
		return
	}
	s.serveLive(w, r, tenantID, eventID)
}

// serveLive streams the change topics of an event over a WebSocket until the client leaves, pinging every 30 s.
func (s *Server) serveLive(w http.ResponseWriter, r *http.Request, tenantID, eventID int64) {
	// Subscribe before the handshake: the client resyncs on open, so nothing may fall between the two.
	msgs, leave := s.Hub.Subscribe(tenantID, eventID)
	defer leave()
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()
	ctx := conn.CloseRead(r.Context())
	ping := time.NewTicker(30 * time.Second)
	defer ping.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case m, ok := <-msgs:
			if !ok {
				return
			}
			b, _ := json.Marshal(map[string]string{"topic": m.Topic})
			wctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := conn.Write(wctx, websocket.MessageText, b)
			cancel()
			if err != nil {
				return
			}
		case <-ping.C:
			pctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := conn.Ping(pctx)
			cancel()
			if err != nil {
				return
			}
		}
	}
}
