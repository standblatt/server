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
	"github.com/standblatt/server/internal/permissions"
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

// access checks that the logged-in user belongs to the tenant and that their role has one of perms (none: any member).
// It returns the user and the role's permissions. Non-members get 404 so tenant IDs cannot be probed. Helper links
// pass only for viewing and operating events; authorize already checked the operation and the path.
func (s *Server) access(ctx context.Context, tenantID int64, perms ...string) (store.User, []string, error) {
	if h := helperOf(ctx); h != nil {
		helper := []string{permissions.EventsView, permissions.EventsOperate}
		if h.TenantID != tenantID {
			return store.User{Name: h.Name}, nil, missing("Club not found")
		}
		if !permissions.Any(helper, perms...) {
			return store.User{Name: h.Name}, nil, forbidden("Not allowed for this helper link")
		}
		return store.User{Name: h.Name}, helper, nil
	}
	u, err := currentUser(ctx)
	if err != nil {
		return u, nil, err
	}
	_, have, err := s.callerPermissions(ctx, tenantID, u.ID)
	if errors.Is(err, store.ErrNotFound) {
		return u, nil, missing("Club not found")
	}
	if err != nil {
		return u, nil, err
	}
	if !permissions.Any(have, perms...) {
		return u, have, forbidden("Not allowed for your role")
	}
	return u, have, nil
}

// callerPermissions returns the user's role in the tenant and the permissions the request may use: the role's,
// cut down to the scopes of an API token. ErrNotFound when the user is no member or the token is not for this club.
func (s *Server) callerPermissions(ctx context.Context, tenantID, userID int64) (store.Role, []string, error) {
	role, err := s.Store.MemberRole(ctx, tenantID, userID)
	if err != nil {
		return role, nil, err
	}
	t := tokenOf(ctx)
	if t == nil {
		return role, role.Permissions, nil
	}
	if !slices.Contains(t.Tenants, tenantID) {
		return role, nil, store.ErrNotFound
	}
	have := slices.DeleteFunc(slices.Clone(role.Permissions), func(p string) bool { return !slices.Contains(t.Scopes, p) })
	return role, have, nil
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
		OperationID: "list-ranges", Metadata: needs("", permissions.EventsView, permissions.RangesManage), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/ranges", Tags: []string{"ranges"},
	}, func(ctx context.Context, in *TenantPath) (*rangesOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
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
		OperationID: "create-range", Metadata: needs("", permissions.RangesManage), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/ranges", Tags: []string{"ranges"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *rangeIn) (*rangeOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
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
		OperationID: "update-range", Metadata: needs("", permissions.RangesManage), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/ranges/{rangeId}", Tags: []string{"ranges"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *rangeUpdateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
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
		OperationID: "delete-range", Metadata: needs("", permissions.RangesManage), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/ranges/{rangeId}", Tags: []string{"ranges"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *RangePath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
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
		OperationID: "list-events", Metadata: needs("", permissions.EventsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events", Tags: []string{"events"},
	}, func(ctx context.Context, in *TenantPath) (*eventsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.Store.ListEvents(ctx, in.TenantID)
		return &eventsOut{Body: list}, err
	})

	type eventOut struct{ Body store.Event }
	huma.Register(api, huma.Operation{
		OperationID: "get-event", Metadata: needs("", permissions.EventsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}", Tags: []string{"events"},
	}, func(ctx context.Context, in *EventPath) (*eventOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
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
		OperationID: "create-event", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events", Tags: []string{"events"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *eventIn) (*eventOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
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
		OperationID: "copy-event", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/copy", Tags: []string{"events"},
		Summary:       "Create a new event with the settings and programmes of this one, e.g. for next year; registrations are not copied",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *copyEventIn) (*eventOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
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
		OperationID: "update-event", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/events/{eventId}", Tags: []string{"events"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *eventUpdateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
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
		OperationID: "delete-event", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/events/{eventId}", Tags: []string{"events"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *EventPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
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
		OperationID: "archive-event", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/events/{eventId}/archived", Tags: []string{"events"},
		Summary:       "Archive or unarchive an event; an archived event stays readable but refuses every change",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *archiveIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.SetEventArchived(ctx, in.TenantID, in.EventID, in.Body.Archived), "Event not found")
	})

	s.registerMembers(api)
	s.registerRoles(api)
}

// registerMembers registers the endpoints that manage the users of a tenant.
func (s *Server) registerMembers(api huma.API) {
	type membersOut struct{ Body []store.Member }
	huma.Register(api, huma.Operation{
		OperationID: "list-users", Metadata: needs("", permissions.ClubManage), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/users", Tags: []string{"users"},
	}, func(ctx context.Context, in *TenantPath) (*membersOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.Store.Members(ctx, in.TenantID)
		return &membersOut{Body: list}, err
	})

	type addIn struct {
		TenantPath
		Body struct {
			Name     string `json:"name" minLength:"1"`
			RoleID   int64  `json:"roleId"`
			Password string `json:"password,omitempty" required:"false" doc:"Required when the user does not exist yet"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "add-user", Metadata: needs("", permissions.ClubManage), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/users", Tags: []string{"users"},
		Summary:       "Add an existing user to the tenant, or create the user with the given password",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *addIn) (*struct{}, error) {
		u, _, err := s.access(ctx, in.TenantID)
		if err != nil {
			return nil, err
		}
		if err := s.roleWithin(ctx, in.TenantID, in.Body.RoleID); err != nil {
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
		if err := s.Store.SetMembership(ctx, in.TenantID, target.ID, in.Body.RoleID); err != nil {
			return nil, err
		}
		s.Log.Info("member added", "tenant", in.TenantID, "user", name, "role", in.Body.RoleID, "by", u.Name)
		return nil, nil
	})

	type roleIn struct {
		MemberPath
		Body struct {
			RoleID int64 `json:"roleId"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-user-role", Metadata: needs("", permissions.ClubManage), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/users/{userId}", Tags: []string{"users"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *roleIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		if err := s.memberWithin(ctx, in.TenantID, in.UserID); err != nil {
			return nil, err
		}
		if err := s.roleWithin(ctx, in.TenantID, in.Body.RoleID); err != nil {
			return nil, err
		}
		if err := s.keepOneAdmin(ctx, in.TenantID, in.UserID, in.Body.RoleID); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.SetMembership(ctx, in.TenantID, in.UserID, in.Body.RoleID), "Role not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "remove-user", Metadata: needs("", permissions.ClubManage), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/users/{userId}", Tags: []string{"users"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *MemberPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		if err := s.memberWithin(ctx, in.TenantID, in.UserID); err != nil {
			return nil, err
		}
		if err := s.keepOneAdmin(ctx, in.TenantID, in.UserID, 0); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.RemoveMembership(ctx, in.TenantID, in.UserID), "User not found")
	})
}

// RolePath identifies a role of a tenant in the URL.
type RolePath struct {
	TenantID int64 `path:"tenantId"`
	RoleID   int64 `path:"roleId"`
}

// registerRoles registers the endpoints that manage the roles of a tenant and their permissions.
func (s *Server) registerRoles(api huma.API) {
	type listOut struct{ Body []store.Role }
	huma.Register(api, huma.Operation{
		OperationID: "list-roles", Metadata: needs("", permissions.ClubManage), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/roles", Tags: []string{"users"},
	}, func(ctx context.Context, in *TenantPath) (*listOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.Store.ListRoles(ctx, in.TenantID)
		return &listOut{Body: list}, err
	})

	type createIn struct {
		TenantPath
		Body store.RoleInput
	}
	type createOut struct{ Body store.Role }
	huma.Register(api, huma.Operation{
		OperationID: "create-role", Metadata: needs("", permissions.ClubManage), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/roles", Tags: []string{"users"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createIn) (*createOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		in.Body.Name = strings.TrimSpace(in.Body.Name)
		if in.Body.Name == "" {
			return nil, invalid("Name is required")
		}
		if err := s.within(ctx, in.TenantID, in.Body.Permissions); err != nil {
			return nil, err
		}
		id, err := s.Store.CreateRole(ctx, in.TenantID, in.Body)
		if err != nil {
			return nil, err
		}
		r, err := s.Store.GetRole(ctx, in.TenantID, id)
		return &createOut{Body: r}, err
	})

	type updateIn struct {
		RolePath
		Body store.RoleInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-role", Metadata: needs("", permissions.ClubManage), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/roles/{roleId}", Tags: []string{"users"},
		Summary: "Rename a role and set its permissions; the admin role keeps every permission", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *updateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		r, err := s.Store.GetRole(ctx, in.TenantID, in.RoleID)
		if err != nil {
			return nil, notFound(err, "Role not found")
		}
		in.Body.Name = strings.TrimSpace(in.Body.Name)
		if in.Body.Name == "" && r.Key == "" {
			return nil, invalid("Name is required")
		}
		if err := s.within(ctx, in.TenantID, append(r.Permissions, in.Body.Permissions...)); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.UpdateRole(ctx, in.TenantID, in.RoleID, in.Body), "Role not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-role", Metadata: needs("", permissions.ClubManage), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/roles/{roleId}", Tags: []string{"users"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *RolePath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		r, err := s.Store.GetRole(ctx, in.TenantID, in.RoleID)
		if err != nil {
			return nil, notFound(err, "Role not found")
		}
		if r.Key == store.RoleAdmin {
			return nil, conflict("The admin role cannot be deleted")
		}
		if err := s.within(ctx, in.TenantID, r.Permissions); err != nil {
			return nil, err
		}
		if r.Members > 0 {
			return nil, conflict("Give the members of this role another role first")
		}
		return nil, notFound(s.Store.DeleteRole(ctx, in.TenantID, in.RoleID), "Role not found")
	})
}

// within refuses with 403 when the caller is no admin and perms go beyond the caller's own permissions: members who
// manage the club without the admin role give, change and take away only what their own role has. The admin role
// has every permission, so only admins make, change or remove admins.
func (s *Server) within(ctx context.Context, tenantID int64, perms []string) error {
	u, err := currentUser(ctx)
	if err != nil {
		return err
	}
	me, have, err := s.callerPermissions(ctx, tenantID, u.ID)
	if err != nil {
		return notFound(err, "Club not found")
	}
	if me.Key == store.RoleAdmin && tokenOf(ctx) == nil {
		return nil
	}
	for _, p := range perms {
		if !slices.Contains(have, p) {
			return forbidden("Only admins give or change permissions your role does not have")
		}
	}
	return nil
}

// roleWithin checks a role of the tenant with within; 404 for an unknown role.
func (s *Server) roleWithin(ctx context.Context, tenantID, roleID int64) error {
	r, err := s.Store.GetRole(ctx, tenantID, roleID)
	if err != nil {
		return notFound(err, "Role not found")
	}
	return s.within(ctx, tenantID, r.Permissions)
}

// memberWithin checks the role of a member of the tenant with within; 404 for a user who is no member.
func (s *Server) memberWithin(ctx context.Context, tenantID, userID int64) error {
	r, err := s.Store.MemberRole(ctx, tenantID, userID)
	if err != nil {
		return notFound(err, "User not found")
	}
	return s.within(ctx, tenantID, r.Permissions)
}

// keepOneAdmin returns a 409 when userID is the last admin of the tenant and is about to get the role newRoleID
// (0 = removed from the club).
func (s *Server) keepOneAdmin(ctx context.Context, tenantID, userID, newRoleID int64) error {
	role, err := s.Store.MemberRole(ctx, tenantID, userID)
	if err != nil || role.Key != store.RoleAdmin || role.ID == newRoleID {
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
