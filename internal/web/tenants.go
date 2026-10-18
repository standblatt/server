package web

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
)

const defaultTimezone = "Europe/Zurich"

// TenantPath identifies a tenant in the URL.
type TenantPath struct {
	TenantID int64 `path:"tenantId"`
}

// RangePath identifies a range of a tenant in the URL.
type RangePath struct {
	TenantID int64 `path:"tenantId"`
	RangeID  int64 `path:"rangeId"`
}

// EventPath identifies an event of a tenant in the URL.
type EventPath struct {
	TenantID int64 `path:"tenantId"`
	EventID  int64 `path:"eventId"`
}

// MemberPath identifies a member of a tenant in the URL.
type MemberPath struct {
	TenantID int64 `path:"tenantId"`
	UserID   int64 `path:"userId"`
}

// access checks that the logged-in user belongs to the tenant, and is an admin there if admin is set.
// Non-members get 404 so tenant IDs cannot be probed.
func (s *Server) access(ctx context.Context, tenantID int64, admin bool) (store.User, store.Role, error) {
	if h := helperOf(ctx); h != nil {
		// authorize already checked the operation and the path; admin pages stay closed to links.
		if h.TenantID != tenantID {
			return store.User{Name: h.Name}, store.RoleHelper, missing("Club not found")
		}
		if admin {
			return store.User{Name: h.Name}, store.RoleHelper, forbidden("Not allowed for this helper link")
		}
		return store.User{Name: h.Name}, store.RoleHelper, nil
	}
	u, err := currentUser(ctx)
	if err != nil {
		return u, "", err
	}
	role, err := s.Store.MemberRole(ctx, tenantID, u.ID)
	if errors.Is(err, store.ErrNotFound) {
		return u, "", missing("Club not found")
	}
	if err != nil {
		return u, "", err
	}
	if admin && role != store.RoleAdmin {
		return u, role, forbidden("Club admins only")
	}
	return u, role, nil
}

// notFound turns ErrNotFound into a 404 with msg and returns other errors unchanged.
func notFound(err error, msg string) error {
	if errors.Is(err, store.ErrNotFound) {
		return &apiError{status: http.StatusNotFound, parts: []domain.Problem{{Format: msg}}}
	}
	return err
}

// validateRange checks that the disabled lanes lie within the lane count of the range.
func validateRange(in store.RangeInput) error {
	for _, l := range in.DisabledLanes {
		if l < 1 || int(l) > in.LaneCount {
			return invalid("Lane %d is outside 1–%d", l, in.LaneCount)
		}
	}
	return nil
}

// normalizeEvent validates and normalizes the event input: date, schedule, ranges, time zone and locale.
func normalizeEvent(in *store.EventInput) error {
	in.Name = strings.TrimSpace(in.Name)
	if _, err := time.Parse(time.DateOnly, in.StartsOn); err != nil {
		return invalid("Date must be YYYY-MM-DD")
	}
	if errs := domain.ValidateSchedule(in.Schedule); len(errs) > 0 {
		return invalidList(errs)
	}
	in.RangeIDs = append([]int64{}, slices.Compact(slices.Sorted(slices.Values(in.RangeIDs)))...)
	if in.Timezone == "" {
		in.Timezone = defaultTimezone
	}
	if in.Locale == "" {
		in.Locale = "en"
	}
	if _, err := time.LoadLocation(in.Timezone); err != nil {
		return invalid("Unknown time zone: %s", in.Timezone)
	}
	return nil
}

// registerTenants registers the range and event endpoints.
func (s *Server) registerTenants(api huma.API) {
	type rangesOut struct{ Body []store.Range }
	huma.Register(api, huma.Operation{
		OperationID: "list-ranges", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/ranges", Tags: []string{"ranges"},
	}, func(ctx context.Context, in *TenantPath) (*rangesOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		list, err := s.Store.ListRanges(ctx, in.TenantID)
		return &rangesOut{Body: list}, err
	})

	type rangeIn struct {
		TenantPath
		Body store.RangeInput
	}
	type rangeOut struct{ Body store.Range }
	huma.Register(api, huma.Operation{
		OperationID: "create-range", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/ranges", Tags: []string{"ranges"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *rangeIn) (*rangeOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		if err := validateRange(in.Body); err != nil {
			return nil, err
		}
		st, err := s.Store.CreateRange(ctx, in.TenantID, in.Body)
		if errors.Is(err, store.ErrExists) {
			return nil, conflict("A range with this name already exists")
		}
		return &rangeOut{Body: st}, err
	})

	type rangeUpdateIn struct {
		RangePath
		Body store.RangeInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-range", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/ranges/{rangeId}", Tags: []string{"ranges"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *rangeUpdateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		if err := validateRange(in.Body); err != nil {
			return nil, err
		}
		err := s.Store.UpdateRange(ctx, in.TenantID, in.RangeID, in.Body)
		if errors.Is(err, store.ErrExists) {
			return nil, conflict("A range with this name already exists")
		}
		return nil, notFound(err, "Range not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-range", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/ranges/{rangeId}", Tags: []string{"ranges"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *RangePath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		err := s.Store.DeleteRange(ctx, in.TenantID, in.RangeID)
		if errors.Is(err, store.ErrInUse) {
			return nil, conflict("The range is used by events")
		}
		return nil, notFound(err, "Range not found")
	})

	type eventsOut struct{ Body []store.Event }
	huma.Register(api, huma.Operation{
		OperationID: "list-events", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events", Tags: []string{"events"},
	}, func(ctx context.Context, in *TenantPath) (*eventsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		list, err := s.Store.ListEvents(ctx, in.TenantID)
		return &eventsOut{Body: list}, err
	})

	type eventOut struct{ Body store.Event }
	huma.Register(api, huma.Operation{
		OperationID: "get-event", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}", Tags: []string{"events"},
	}, func(ctx context.Context, in *EventPath) (*eventOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		ev, err := s.Store.GetEvent(ctx, in.TenantID, in.EventID)
		return &eventOut{Body: ev}, notFound(err, "Event not found")
	})

	type eventIn struct {
		TenantPath
		Body store.EventInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-event", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events", Tags: []string{"events"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *eventIn) (*eventOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		if err := normalizeEvent(&in.Body); err != nil {
			return nil, err
		}
		if err := s.checkEventForms(ctx, in.TenantID, in.Body.Settings); err != nil {
			return nil, err
		}
		id, err := s.Store.CreateEvent(ctx, in.TenantID, in.Body)
		if errors.Is(err, store.ErrInUse) {
			return nil, invalid("Range not found")
		}
		if err != nil {
			return nil, err
		}
		ev, err := s.Store.GetEvent(ctx, in.TenantID, id)
		return &eventOut{Body: ev}, err
	})

	type copyEventIn struct {
		EventPath
		Body struct {
			Name     string `json:"name" minLength:"1" maxLength:"200"`
			StartsOn string `json:"startsOn" format:"date"`
			Template bool   `json:"template,omitempty" required:"false" doc:"Save the copy as a template of the club"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "copy-event", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/copy", Tags: []string{"events"},
		Summary:       "Create a new event with the settings and programmes of this one, e.g. for next year; registrations are not copied",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *copyEventIn) (*eventOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		id, err := s.Store.CopyEvent(ctx, in.TenantID, in.EventID, strings.TrimSpace(in.Body.Name), in.Body.StartsOn, in.Body.Template)
		if err != nil {
			return nil, notFound(err, "Event not found")
		}
		ev, err := s.Store.GetEvent(ctx, in.TenantID, id)
		return &eventOut{Body: ev}, err
	})

	type eventUpdateIn struct {
		EventPath
		Body store.EventInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-event", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/events/{eventId}", Tags: []string{"events"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *eventUpdateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		if err := normalizeEvent(&in.Body); err != nil {
			return nil, err
		}
		if err := s.checkEventForms(ctx, in.TenantID, in.Body.Settings); err != nil {
			return nil, err
		}
		err := s.Store.UpdateEvent(ctx, in.TenantID, in.EventID, in.Body)
		var used *store.ModuleUsedError
		switch {
		case errors.As(err, &used):
			return nil, moduleUsed(used)
		case errors.Is(err, store.ErrInUse):
			return nil, invalid("Range not found")
		case errors.Is(err, store.ErrRangeUsed):
			return nil, conflict("Shooters or programmes use a range that was removed from the event")
		}
		return nil, notFound(err, "Event not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-event", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/events/{eventId}", Tags: []string{"events"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *EventPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		err := s.Store.DeleteEvent(ctx, in.TenantID, in.EventID)
		if errors.Is(err, store.ErrInUse) {
			return nil, conflict("Events with registrations cannot be deleted")
		}
		return nil, notFound(err, "Event not found")
	})

	type archiveIn struct {
		EventPath
		Body struct {
			Archived bool `json:"archived"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "archive-event", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/events/{eventId}/archived", Tags: []string{"events"},
		Summary:       "Archive or unarchive an event; an archived event stays readable but refuses every change",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *archiveIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.SetEventArchived(ctx, in.TenantID, in.EventID, in.Body.Archived), "Event not found")
	})

	s.registerMembers(api)
}

// registerMembers registers the endpoints that manage the users of a tenant.
func (s *Server) registerMembers(api huma.API) {
	type membersOut struct{ Body []store.Member }
	huma.Register(api, huma.Operation{
		OperationID: "list-users", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/users", Tags: []string{"users"},
	}, func(ctx context.Context, in *TenantPath) (*membersOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		list, err := s.Store.Members(ctx, in.TenantID)
		return &membersOut{Body: list}, err
	})

	type addIn struct {
		TenantPath
		Body struct {
			Name     string     `json:"name" minLength:"1"`
			Role     store.Role `json:"role" enum:"admin,staff"`
			Password string     `json:"password,omitempty" required:"false" doc:"Required when the user does not exist yet"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "add-user", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/users", Tags: []string{"users"},
		Summary:       "Add an existing user to the tenant, or create the user with the given password",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *addIn) (*struct{}, error) {
		u, _, err := s.access(ctx, in.TenantID, true)
		if err != nil {
			return nil, err
		}
		name := strings.ToLower(strings.TrimSpace(in.Body.Name))
		target, err := s.Store.UserByName(ctx, name)
		switch {
		case errors.Is(err, store.ErrNotFound):
			if errs := append(ValidateUserName(name), ValidatePassword(in.Body.Password)...); len(errs) > 0 {
				return nil, invalidList(errs)
			}
			hash, err := HashPassword(in.Body.Password)
			if err != nil {
				return nil, err
			}
			if target.ID, err = s.Store.CreateUser(ctx, name, hash); err != nil {
				return nil, err
			}
		case err != nil:
			return nil, err
		}
		if err := s.Store.SetMembership(ctx, in.TenantID, target.ID, in.Body.Role); err != nil {
			return nil, err
		}
		s.Log.Info("member added", "tenant", in.TenantID, "user", name, "role", in.Body.Role, "by", u.Name)
		return nil, nil
	})

	type roleIn struct {
		MemberPath
		Body struct {
			Role store.Role `json:"role" enum:"admin,staff"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-user-role", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/users/{userId}", Tags: []string{"users"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *roleIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		if in.Body.Role != store.RoleAdmin {
			if err := s.keepOneAdmin(ctx, in.TenantID, in.UserID); err != nil {
				return nil, err
			}
		}
		if _, err := s.Store.MemberRole(ctx, in.TenantID, in.UserID); err != nil {
			return nil, notFound(err, "User not found")
		}
		return nil, s.Store.SetMembership(ctx, in.TenantID, in.UserID, in.Body.Role)
	})

	huma.Register(api, huma.Operation{
		OperationID: "remove-user", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/users/{userId}", Tags: []string{"users"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *MemberPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		if err := s.keepOneAdmin(ctx, in.TenantID, in.UserID); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.RemoveMembership(ctx, in.TenantID, in.UserID), "User not found")
	})
}

// keepOneAdmin returns a 409 when userID is the last admin of the tenant.
func (s *Server) keepOneAdmin(ctx context.Context, tenantID, userID int64) error {
	role, err := s.Store.MemberRole(ctx, tenantID, userID)
	if err != nil || role != store.RoleAdmin {
		return nil
	}
	n, err := s.Store.AdminCount(ctx, tenantID)
	if err != nil {
		return err
	}
	if n <= 1 {
		return conflict("A club needs at least one admin")
	}
	return nil
}
