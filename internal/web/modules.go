package web

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strconv"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/store"
)

// moduleParams are the path parameters the module gate finds the event by, in this order.
var moduleParams = []string{"eventId", "programmeId", "registrationId", "entryId"}

// documentModules are the documents that belong to a module.
var documentModules = map[string]string{store.DocLaneSheet: modules.LaneSheet, store.DocControlStrip: modules.ControlStrip}

// reportModules are the reports that belong to a module.
var reportModules = map[string]string{
	ReportAwards: modules.Awards, ReportAwardsOpen: modules.Awards, ReportClubCompetition: modules.ClubCompetition,
	ReportSpecialPrizes: modules.SpecialPrizes, ReportBookings: modules.Reservations, ReportClassifications: modules.Classifications,
}

// moduleGate refuses operations of a module (see modules.Metadata) with 404 while the module is off for the event
// in the path, or for the club of the path for club modules. Requests the handler refuses anyway (no access, unknown event) pass through to it.
func (s *Server) moduleGate(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		key := modules.Of(ctx.Operation().Metadata)
		if key == "" {
			next(ctx)
			return
		}
		on, err := s.moduleOn(ctx.Context(), key, ctx.Param)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "internal error")
			return
		}
		if !on {
			huma.WriteErr(api, ctx, http.StatusNotFound, moduleOff().Error())
			return
		}
		next(ctx)
	}
}

// moduleOff is the 404 of an operation whose module is off.
func moduleOff() error {
	return missing("This module is off")
}

// moduleOn tells whether a module is on for the event the path points to; true when the caller has no access or the
// event is unknown, so the handler answers.
func (s *Server) moduleOn(ctx context.Context, key string, param func(string) string) (bool, error) {
	tenantID, err := strconv.ParseInt(param("tenantId"), 10, 64)
	if err != nil {
		return true, nil
	}
	if _, _, err := s.access(ctx, tenantID, false); err != nil {
		return true, nil
	}
	if m, _ := modules.Get(key); m.Level == modules.Club {
		on, _, err := s.Store.ClubModules(ctx, tenantID)
		return slices.Contains(on, key), ignoreNotFound(err)
	}
	for _, p := range moduleParams {
		id, err := strconv.ParseInt(param(p), 10, 64)
		if err != nil {
			continue
		}
		eventID := id
		if p != "eventId" {
			if eventID, err = s.Store.EventOf(ctx, tenantID, p, id); err != nil {
				return errors.Is(err, store.ErrNotFound), ignoreNotFound(err)
			}
		}
		ev, err := s.Store.GetEvent(ctx, tenantID, eventID)
		if err != nil {
			return errors.Is(err, store.ErrNotFound), ignoreNotFound(err)
		}
		return slices.Contains(ev.Modules, key), nil
	}
	return true, nil
}

// ignoreNotFound drops ErrNotFound.
func ignoreNotFound(err error) error {
	if errors.Is(err, store.ErrNotFound) {
		return nil
	}
	return err
}

// requireModule answers 404 when a module is off for an event, for checks the gate cannot make from the path.
func requireModule(ev store.Event, key string) error {
	if !slices.Contains(ev.Modules, key) {
		return missing("This module is off for this event")
	}
	return nil
}

// moduleUsed turns a store.ModuleUsedError into a 409 for its first module, telling what to remove first.
func moduleUsed(e *store.ModuleUsedError) error {
	switch e.Keys[0] {
	case modules.Reservations:
		return conflict("Switch the online booking link off first")
	case modules.Attempts:
		return conflict("Programmes with several attempts use this module")
	case modules.Awards:
		return conflict("Programmes with awards use this module")
	case modules.SpecialPrizes:
		return conflict("Programmes with special prizes use this module")
	case modules.Classifications:
		return conflict("Delete the classifications first")
	case modules.Kiosk:
		return conflict("Delete the kiosk links first")
	case modules.Helpers:
		return conflict("Delete the helper links first")
	case modules.Forms:
		return conflict("Delete the print forms first")
	}
	return conflict("The module is in use")
}

// registerModules registers the list of modules.
func (s *Server) registerModules(api huma.API) {
	type listOut struct{ Body []modules.Module }
	huma.Register(api, huma.Operation{
		OperationID: "list-modules", Method: http.MethodGet, Path: "/api/modules", Tags: []string{"events"},
		Summary: "The modules a club or an event can switch on, with what they require",
	}, func(ctx context.Context, _ *struct{}) (*listOut, error) {
		return &listOut{Body: modules.All}, nil
	})

	type clubModulesIn struct {
		TenantPath
		Body struct {
			Modules []string `json:"modules" enum:"forms"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-club-modules", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/modules", Tags: []string{"tenants"},
		Summary: "Choose the modules of a club; modules that hold data stay on", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *clubModulesIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		err := s.Store.SetClubModules(ctx, in.TenantID, in.Body.Modules)
		var used *store.ModuleUsedError
		if errors.As(err, &used) {
			return nil, moduleUsed(used)
		}
		return nil, notFound(err, "Club not found")
	})
}

// withModule marks an operation as part of a module.
func withModule(op huma.Operation, key string) huma.Operation {
	op.Metadata = modules.Metadata(key)
	return op
}
