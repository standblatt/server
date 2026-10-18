package web

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

// ProgrammeFormat identifies programme export files.
const ProgrammeFormat = "standblatt.programme"

// legacyProgrammeFormat is the id files exported before the rename carry; import still accepts it.
const legacyProgrammeFormat = "targetlog.programme"

// ProgrammeExport is the file format for exporting and importing a programme, also used for the library.
type ProgrammeExport struct {
	Format    string                     `json:"format" enum:"standblatt.programme"`
	Version   int                        `json:"version" enum:"1"`
	Programme domain.ProgrammeDefinition `json:"programme"`
}

// ProgrammePath identifies a programme of a tenant in the URL.
type ProgrammePath struct {
	TenantID    int64 `path:"tenantId"`
	ProgrammeID int64 `path:"programmeId"`
}

// checkDefinition normalizes a programme definition and returns its validation problems as one 422.
func checkDefinition(d *domain.ProgrammeDefinition) error {
	d.Normalize()
	if errs := d.Validate(); len(errs) > 0 {
		return invalidList(errs)
	}
	return nil
}

// programmeErr maps store errors of programme writes to 409 or 404.
func programmeErr(err error) error {
	switch {
	case errors.Is(err, store.ErrExists):
		return conflict("The event already has a programme with this name and variant")
	case errors.Is(err, store.ErrNotFound):
		return missing("Programme or event not found")
	case errors.Is(err, store.ErrAttemptsUsed):
		return conflict("Shooters already have more attempts at this programme than the new number allows")
	case errors.Is(err, store.ErrInUse):
		return conflict("Shooters are registered for this programme: the series stay as they are, and age categories and awards can be added but not removed")
	}
	return err
}

// notCombined returns a 409 when the programme takes series from another one or another one takes series from it,
// because such programmes are shot once.
func (s *Server) notCombined(ctx context.Context, tenantID, programmeID int64) error {
	p, err := s.Store.GetProgramme(ctx, tenantID, programmeID)
	if err != nil {
		return notFound(err, "Programme not found")
	}
	list, err := s.Store.ListProgrammes(ctx, tenantID, p.EventID)
	if err != nil {
		return err
	}
	for _, q := range list {
		if p.CombinableWithID != nil || (q.CombinableWithID != nil && *q.CombinableWithID == p.ID) {
			return conflict("Programmes with several attempts cannot be combined")
		}
	}
	return nil
}

// registerProgrammes registers the programme endpoints of events and the library.
func (s *Server) registerProgrammes(api huma.API) {
	type listOut struct{ Body []store.Programme }
	huma.Register(api, huma.Operation{
		OperationID: "list-programmes", Metadata: needs("", permissions.EventsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/events/{eventId}/programmes", Tags: []string{"programmes"},
	}, func(ctx context.Context, in *EventPath) (*listOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		list, err := s.Store.ListProgrammes(ctx, in.TenantID, in.EventID)
		return &listOut{Body: list}, err
	})

	type programmeOut struct{ Body store.Programme }
	createAndReturn := func(ctx context.Context, tenantID, eventID int64, d domain.ProgrammeDefinition) (*programmeOut, error) {
		if err := checkDefinition(&d); err != nil {
			return nil, err
		}
		id, err := s.Store.CreateProgramme(ctx, tenantID, eventID, d)
		if err != nil {
			return nil, programmeErr(err)
		}
		s.notify(ctx, tenantID, eventID, realtime.TopicProgrammes)
		p, err := s.Store.GetProgramme(ctx, tenantID, id)
		return &programmeOut{Body: p}, err
	}

	type createIn struct {
		EventPath
		Body domain.ProgrammeDefinition
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-programme", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/programmes", Tags: []string{"programmes"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createIn) (*programmeOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		return createAndReturn(ctx, in.TenantID, in.EventID, in.Body)
	})

	type importIn struct {
		EventPath
		Body ProgrammeExport
	}
	huma.Register(api, huma.Operation{
		OperationID: "import-programme", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/programmes/import", Tags: []string{"programmes"},
		Summary: "Add a programme from an export file", DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *importIn) (*programmeOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		return createAndReturn(ctx, in.TenantID, in.EventID, in.Body.Programme)
	})

	type fromTemplateIn struct {
		EventPath
		Body struct {
			TemplateID int64 `json:"templateId"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "programme-from-template", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/programmes/from-template", Tags: []string{"programmes"},
		Summary: "Add a copy of a library template", DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *fromTemplateIn) (*programmeOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		tpl, err := s.Store.GetTemplate(ctx, in.Body.TemplateID)
		if err != nil {
			return nil, notFound(err, "Template not found")
		}
		return createAndReturn(ctx, in.TenantID, in.EventID, tpl.Definition)
	})

	type copyIn struct {
		EventPath
		Body struct {
			FromEventID int64 `json:"fromEventId"`
		}
	}
	type copyOut struct {
		Body struct {
			Copied int `json:"copied"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "copy-programmes", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/programmes/copy", Tags: []string{"programmes"},
		Summary: "Copy all programmes of another event of the club; existing name/variant pairs are skipped",
	}, func(ctx context.Context, in *copyIn) (*copyOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		n, err := s.Store.CopyProgrammes(ctx, in.TenantID, in.Body.FromEventID, in.EventID)
		if err == nil {
			s.notify(ctx, in.TenantID, in.EventID, realtime.TopicProgrammes)
		}
		out := &copyOut{}
		out.Body.Copied = n
		return out, notFound(err, "Event not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-programme", Metadata: needs("", permissions.EventsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/programmes/{programmeId}", Tags: []string{"programmes"},
	}, func(ctx context.Context, in *ProgrammePath) (*programmeOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		p, err := s.Store.GetProgramme(ctx, in.TenantID, in.ProgrammeID)
		return &programmeOut{Body: p}, notFound(err, "Programme not found")
	})

	type exportOut struct {
		ContentDisposition string `header:"Content-Disposition"`
		Body               ProgrammeExport
	}
	huma.Register(api, huma.Operation{
		OperationID: "export-programme", Metadata: needs("", permissions.EventsView), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/programmes/{programmeId}/export", Tags: []string{"programmes"},
	}, func(ctx context.Context, in *ProgrammePath) (*exportOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		p, err := s.Store.GetProgramme(ctx, in.TenantID, in.ProgrammeID)
		if err != nil {
			return nil, notFound(err, "Programme not found")
		}
		return &exportOut{
			ContentDisposition: `attachment; filename="` + exportFilename(p.ProgrammeDefinition) + `"`,
			Body:               ProgrammeExport{Format: ProgrammeFormat, Version: 1, Programme: p.ProgrammeDefinition},
		}, nil
	})

	type updateIn struct {
		ProgrammePath
		Body domain.ProgrammeDefinition
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-programme", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/programmes/{programmeId}", Tags: []string{"programmes"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *updateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		if err := checkDefinition(&in.Body); err != nil {
			return nil, err
		}
		if in.Body.AttemptLimit() != 1 {
			if err := s.notCombined(ctx, in.TenantID, in.ProgrammeID); err != nil {
				return nil, err
			}
		}
		return nil, s.programmeChanged(ctx, in.TenantID, in.ProgrammeID, programmeErr(s.Store.UpdateProgramme(ctx, in.TenantID, in.ProgrammeID, in.Body)))
	})

	type activeIn struct {
		ProgrammePath
		Body struct {
			Active bool `json:"active"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-programme-active", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/programmes/{programmeId}/active", Tags: []string{"programmes"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *activeIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		return nil, s.programmeChanged(ctx, in.TenantID, in.ProgrammeID, notFound(s.Store.SetProgrammeActive(ctx, in.TenantID, in.ProgrammeID, in.Body.Active), "Programme not found"))
	})

	type combinationIn struct {
		ProgrammePath
		Body struct {
			CombinableWithID *int64 `json:"combinableWithId,omitempty" required:"false" doc:"Programme of the same event, or omitted to clear"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-programme-combination", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/programmes/{programmeId}/combination", Tags: []string{"programmes"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *combinationIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		if in.Body.CombinableWithID != nil {
			for _, id := range []int64{in.ProgrammeID, *in.Body.CombinableWithID} {
				if p, err := s.Store.GetProgramme(ctx, in.TenantID, id); err == nil && p.AttemptLimit() != 1 {
					return nil, conflict("Programmes with several attempts cannot be combined")
				}
			}
		}
		err := s.Store.SetCombination(ctx, in.TenantID, in.ProgrammeID, in.Body.CombinableWithID)
		return nil, s.programmeChanged(ctx, in.TenantID, in.ProgrammeID, notFound(err, "Programme not found, or the other programme belongs to a different event"))
	})

	type prizesIn struct {
		ProgrammePath
		Body domain.SpecialPrizeRules
	}
	type rangeIn struct {
		ProgrammePath
		Body struct {
			RangeID *int64 `json:"rangeId,omitempty" required:"false" doc:"A range of the event; leave out for any range"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-programme-range", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/programmes/{programmeId}/range", Tags: []string{"programmes"},
		Summary:       "Tie the programme to one range of its event; allowed while shooters are registered",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *rangeIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		return nil, s.programmeChanged(ctx, in.TenantID, in.ProgrammeID, notFound(s.Store.SetProgrammeRange(ctx, in.TenantID, in.ProgrammeID, in.Body.RangeID), "Programme not found, or the range is not one of the event"))
	})

	huma.Register(api, huma.Operation{
		OperationID: "set-special-prizes", Metadata: needs(modules.SpecialPrizes, permissions.EventsConfigure), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/programmes/{programmeId}/special-prizes", Tags: []string{"programmes"},
		Summary:       "Set the special prize rules of a programme; allowed while shooters are registered",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *prizesIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		return nil, s.programmeChanged(ctx, in.TenantID, in.ProgrammeID, notFound(s.Store.SetSpecialPrizes(ctx, in.TenantID, in.ProgrammeID, in.Body), "Programme not found"))
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-programme", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/programmes/{programmeId}", Tags: []string{"programmes"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *ProgrammePath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		p, err := s.Store.GetProgramme(ctx, in.TenantID, in.ProgrammeID)
		if err != nil {
			return nil, notFound(err, "Programme not found")
		}
		if err := s.Store.DeleteProgramme(ctx, in.TenantID, in.ProgrammeID); err != nil {
			if errors.Is(err, store.ErrInUse) {
				return nil, conflict("Shooters are registered for this programme; it cannot be deleted")
			}
			return nil, programmeErr(err)
		}
		s.notify(ctx, in.TenantID, p.EventID, realtime.TopicProgrammes)
		return nil, nil
	})

	type templatesOut struct{ Body []store.ProgrammeTemplate }
	huma.Register(api, huma.Operation{
		OperationID: "list-programme-templates", Method: http.MethodGet, Path: "/api/programme-templates", Tags: []string{"programmes"},
		Summary: "The global programme library",
	}, func(ctx context.Context, _ *struct{}) (*templatesOut, error) {
		if _, err := currentUser(ctx); err != nil {
			return nil, err
		}
		list, err := s.Store.ListTemplates(ctx)
		return &templatesOut{Body: list}, err
	})
}

// programmeChanged tells the programme's event about a successful change and passes err through.
func (s *Server) programmeChanged(ctx context.Context, tenantID, programmeID int64, err error) error {
	if err != nil {
		return err
	}
	if p, err := s.Store.GetProgramme(ctx, tenantID, programmeID); err == nil {
		s.notify(ctx, tenantID, p.EventID, realtime.TopicProgrammes)
	}
	return nil
}

// exportFilename returns the file name of a programme export: name and variant in lower-case words joined by dashes.
func exportFilename(d domain.ProgrammeDefinition) string {
	name := d.Name
	if d.Variant != "" {
		name += " " + d.Variant
	}
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case b.Len() > 0 && !strings.HasSuffix(b.String(), "-"):
			b.WriteByte('-')
		}
	}
	return strings.TrimSuffix(b.String(), "-") + ".programme.json"
}
