package web

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

// helperPrefix marks helper link secrets (sh_…); personal API tokens are sb_….
const helperPrefix = "sh_"

const helperKey ctxKey = 10

// helperOps lists the API operations each helper role may call; everything else is 403.
var helperOps = func() map[store.HelperRole]map[string]bool {
	set := func(ops ...string) map[string]bool {
		m := map[string]bool{}
		for _, o := range append(ops, "get-event", "get-registration", "search-shooters", "get-shooter", "list-lanes", "list-programmes", "get-event-stats", "get-results", "live", "helper") {
			m[o] = true
		}
		return m
	}
	return map[store.HelperRole]map[string]bool{
		store.HelperScorer:     set("save-results", "set-award", "finish-registration", "list-registrations", "print", "entry-form", "registration-form"),
		store.HelperResults:    set("save-results", "set-award", "list-registrations", "issue-awards", "print", "entry-form", "registration-form"),
		store.HelperAmmunition: set("list-registrations", "update-registration", "issue-ammo", "print", "registration-document"),
	}
}()

// authorize is the huma middleware that keeps helper links in their place before any handler runs: the
// operation must be in helperOps for the role, the tenant and event in the path must be the link's, and a
// registration or entry in the path must be on the scorer's lanes. Logged-in users pass straight through;
// access then checks their membership.
func (s *Server) authorize(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		h := helperOf(ctx.Context())
		if h == nil {
			next(ctx)
			return
		}
		if err := s.helperAllowed(ctx.Context(), h, ctx.Operation().OperationID, ctx.Param); err != nil {
			var se huma.StatusError
			if errors.As(err, &se) {
				huma.WriteErr(api, ctx, se.GetStatus(), se.Error())
				return
			}
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "internal error")
			return
		}
		next(ctx)
	}
}

// helperAllowed decides a helper link's request from the operation and the path parameters.
func (s *Server) helperAllowed(ctx context.Context, h *store.Helper, op string, param func(string) string) error {
	if !helperOps[h.Role][op] {
		return forbidden("Not allowed for this helper link")
	}
	if t := param("tenantId"); t != "" && t != strconv.FormatInt(h.TenantID, 10) {
		return missing("Club not found")
	}
	if e := param("eventId"); e != "" && e != strconv.FormatInt(h.EventID, 10) {
		return missing("Event not found")
	}
	registrationID, _ := strconv.ParseInt(param("registrationId"), 10, 64)
	if entryID, _ := strconv.ParseInt(param("entryId"), 10, 64); entryID != 0 {
		e, err := s.Store.GetEntry(ctx, h.TenantID, entryID)
		if err != nil {
			return notFound(err, "Programme of the registration not found")
		}
		registrationID = e.RegistrationID
	}
	if registrationID == 0 {
		return nil
	}
	reg, err := s.Store.GetRegistration(ctx, h.TenantID, registrationID)
	if err != nil {
		return notFound(err, "Registration not found")
	}
	return helperLane(h, reg)
}

// helperLane refuses a registration outside the link's event and, for scorers, off their lanes.
func helperLane(h *store.Helper, reg store.Registration) error {
	if reg.EventID != h.EventID {
		return missing("Registration not found")
	}
	if h.Role == store.HelperScorer && (reg.RangeID == nil || reg.Lane == nil || *reg.RangeID != *h.RangeID || !slices.Contains(h.Lanes, *reg.Lane)) {
		return forbidden("This lane is not yours")
	}
	return nil
}

// helperOf returns the helper that withSession stored in ctx, nil when the request has none.
func helperOf(ctx context.Context) *store.Helper {
	h, _ := ctx.Value(helperKey).(*store.Helper)
	return h
}

// helperDocument allows helper links their documents from a request body: the ammunition desk the lane sheet,
// scorers and results clerks the print forms of an entry or a registration; the registration is checked the way the
// middleware checks path parameters.
func (s *Server) helperDocument(ctx context.Context, tenantID, registrationID, entryID int64, document string) error {
	h := helperOf(ctx)
	if h == nil {
		return nil
	}
	allowed := document == store.DocLaneSheet && h.Role == store.HelperAmmunition || document == store.DocForm && h.Role != store.HelperAmmunition
	if !allowed {
		return forbidden("Not allowed for this helper link")
	}
	if document == store.DocForm && entryID != 0 {
		e, err := s.Store.GetEntry(ctx, tenantID, entryID)
		if err != nil {
			return notFound(err, "Programme of the registration not found")
		}
		registrationID = e.RegistrationID
	}
	reg, err := s.Store.GetRegistration(ctx, tenantID, registrationID)
	if err != nil {
		return notFound(err, "Registration not found")
	}
	return helperLane(h, reg)
}

// sameProgrammes reports whether an update keeps the programmes of a registration, their attempts and combined series (helpers may only move lanes).
func sameProgrammes(entries []store.Entry, choices []domain.ProgrammeChoice) bool {
	have := map[int64][]int{}
	count := map[int64]int{}
	for _, e := range entries {
		have[e.ProgrammeID] = e.CombinedSeries
		count[e.ProgrammeID]++
	}
	if len(count) != len(choices) {
		return false
	}
	for _, c := range choices {
		series, ok := have[c.ProgrammeID]
		if !ok || !slices.Equal(series, c.CombinedSeries) || count[c.ProgrammeID] != c.Attempts() {
			return false
		}
	}
	return true
}

// registerHelpers registers the helper code endpoints and the helper endpoint.
func (s *Server) registerHelpers(api huma.API) {
	type listOut struct{ Body []store.HelperCode }
	huma.Register(api, huma.Operation{
		OperationID: "list-helper-codes", Metadata: modules.Metadata(modules.Helpers), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/helper-codes", Tags: []string{"helpers"},
	}, func(ctx context.Context, in *EventPath) (*listOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, permissions.EventsConfigure); err != nil {
			return nil, err
		}
		list, err := s.Store.ListHelperCodes(ctx, in.TenantID, in.EventID)
		return &listOut{Body: list}, err
	})

	type createIn struct {
		EventPath
		Body struct {
			Name    string           `json:"name" minLength:"1" maxLength:"100" doc:"Who gets the link, e.g. Scorer range A"`
			Role    store.HelperRole `json:"role" enum:"scorer,ammunition,results"`
			RangeID *int64           `json:"rangeId,omitempty" required:"false" doc:"Scorers: the range of their lanes"`
			Lanes   string           `json:"lanes,omitempty" required:"false" doc:"Scorers: their lanes, e.g. 1-5, 8"`
		}
	}
	type createOut struct {
		Body struct {
			Code   store.HelperCode `json:"code"`
			Secret string           `json:"secret" doc:"Shown once; the helper opens /h/<secret>"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-helper-code", Metadata: modules.Metadata(modules.Helpers), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/helper-codes", Tags: []string{"helpers"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createIn) (*createOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, permissions.EventsConfigure); err != nil {
			return nil, err
		}
		ev, err := s.Store.GetEvent(ctx, in.TenantID, in.EventID)
		if err != nil {
			return nil, notFound(err, "Event not found")
		}
		c := store.HelperCode{Name: strings.TrimSpace(in.Body.Name), Role: in.Body.Role, Lanes: []int{}}
		if c.Role == store.HelperScorer {
			if !slices.Contains(ev.Modules, modules.Lanes) {
				return nil, invalid("Scorer links need the lane assignment")
			}
			i := slices.IndexFunc(ev.Ranges, func(r store.EventRange) bool { return in.Body.RangeID != nil && r.ID == *in.Body.RangeID })
			if i < 0 {
				return nil, invalid("Choose a range of the event for the scorer")
			}
			lanes, err := domain.ParseLanes(in.Body.Lanes)
			if err != nil || len(lanes) == 0 {
				return nil, invalid("Lanes: numbers and ranges such as 1-5, 8")
			}
			if lanes[len(lanes)-1] > ev.Ranges[i].LaneCount {
				return nil, invalid("Lane %d does not exist on %s", lanes[len(lanes)-1], ev.Ranges[i].Name)
			}
			c.RangeID = in.Body.RangeID
			c.Lanes = lanes
		}
		secret := helperPrefix + randomToken()
		id, err := s.Store.CreateHelperCode(ctx, in.TenantID, in.EventID, c, sessionID(secret))
		if err != nil {
			return nil, notFound(err, "Event not found")
		}
		c.ID = id
		out := &createOut{}
		out.Body.Code = c
		out.Body.Secret = secret
		return out, nil
	})

	type codePath struct {
		TenantID int64 `path:"tenantId"`
		CodeID   int64 `path:"codeId"`
	}
	huma.Register(api, huma.Operation{
		OperationID: "delete-helper-code", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/helper-codes/{codeId}", Tags: []string{"helpers"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *codePath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, permissions.EventsConfigure); err != nil {
			return nil, err
		}
		eventID, err := s.Store.DeleteHelperCode(ctx, in.TenantID, in.CodeID)
		if err != nil {
			return nil, notFound(err, "Helper link not found")
		}
		// The revoked helper's page refetches the event on this topic, gets 401 and reloads into the "not valid" screen.
		s.notify(ctx, in.TenantID, eventID, realtime.TopicHelpers)
		return nil, nil
	})

	type meOut struct{ Body store.Helper }
	huma.Register(api, huma.Operation{
		OperationID: "helper", Method: http.MethodGet, Path: "/api/helper", Tags: []string{"helpers"},
		Summary: "The event and role behind the helper link in the Authorization header (Bearer sh_…)",
	}, func(ctx context.Context, _ *struct{}) (*meOut, error) {
		h := helperOf(ctx)
		if h == nil {
			return nil, unauthorized("Unknown or revoked helper link")
		}
		return &meOut{Body: *h}, nil
	})
}

// helperFromRequest resolves a helper link secret from the Authorization header, or from ?token= on the live feed.
func (s *Server) helperFromRequest(ctx context.Context, r *http.Request) (*store.Helper, error) {
	token := bearerToken(r)
	if token == "" && strings.HasSuffix(r.URL.Path, "/live") {
		token = r.URL.Query().Get("token")
	}
	if !strings.HasPrefix(token, helperPrefix) {
		return nil, nil
	}
	h, err := s.Store.HelperByHash(ctx, sessionID(token))
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &h, nil
}
