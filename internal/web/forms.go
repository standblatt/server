package web

import (
	"context"
	"errors"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/docs"
	"github.com/standblatt/server/internal/i18n"
	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/store"
)

// LayoutPath identifies a print form of a tenant in the URL.
type LayoutPath struct {
	TenantID int64 `path:"tenantId"`
	LayoutID int64 `path:"layoutId"`
}

// LayoutField is a placeholder of print forms with a made-up example value.
type LayoutField struct {
	Key     string `json:"key"`
	Example string `json:"example"`
}

const maxBackgroundBytes = 4 << 20

// formValues are the placeholder values of an entry of a registration, as printed onto its form.
func formValues(l i18n.Lang, ev store.Event, sh store.Shooter, club store.Club, reg store.Registration, er EntryResult) map[string]string {
	p := docPerson(l, reg)
	loc, err := time.LoadLocation(ev.Timezone)
	if err != nil {
		loc = time.UTC
	}
	v := map[string]string{
		"shooter.firstname": p.Firstname, "shooter.lastname": p.Lastname, "shooter.name": p.Lastname + " " + p.Firstname,
		"shooter.license": p.License, "shooter.birthday": p.Birthday, "shooter.birthyear": p.BirthYear,
		"shooter.postalcode": sh.Postalcode, "shooter.city": sh.City,
		"club": p.Club, "club.ssv": club.SSVNumber,
		"event.name": ev.Name, "event.date": swissDate(ev.StartsOn), "event.organizer": ev.Organizer,
		"programme": er.Title, "category": p.Category, "registration": strconv.FormatInt(reg.ID, 10),
		"today": time.Now().In(loc).Format("02.01.2006"),
	}
	if reg.Lane != nil {
		v["lane"] = strconv.Itoa(*reg.Lane)
	}
	if len(er.Results) > 0 {
		v["total"], v["mouches"] = docs.Points(er.Total), strconv.Itoa(er.Mouches)
	}
	if er.Award != nil {
		v["award"] = er.Award.Name
	}
	for i, r := range er.Results {
		n := "s" + strconv.Itoa(i+1)
		v[n] = docs.Points(r.Points)
		if r.Mouches != nil {
			v[n+".mouches"] = strconv.Itoa(*r.Mouches)
		}
		for j, shot := range r.Shots {
			v[n+".shot"+strconv.Itoa(j+1)] = docs.Points(shot)
		}
	}
	return v
}

// entryFormValues loads what the form of an entry prints: the event, the shooter, the club and the results.
func (s *Server) entryFormValues(ctx context.Context, tenantID, entryID int64) (store.EntryRef, store.Event, map[string]string, error) {
	e, err := s.Store.GetEntry(ctx, tenantID, entryID)
	if err != nil {
		return e, store.Event{}, nil, notFound(err, "Programme of the registration not found")
	}
	sheet, err := s.resultSheet(ctx, tenantID, e.RegistrationID)
	if err != nil {
		return e, store.Event{}, nil, err
	}
	ev, err := s.Store.GetEvent(ctx, tenantID, e.EventID)
	if err != nil {
		return e, ev, nil, err
	}
	sh, err := s.Store.GetShooter(ctx, tenantID, sheet.Registration.Shooter.ID)
	if err != nil {
		return e, ev, nil, err
	}
	var club store.Club
	if id := sheet.Registration.ClubID; id != nil {
		clubs, err := s.Store.ListClubs(ctx, tenantID)
		if err != nil {
			return e, ev, nil, err
		}
		if i := slices.IndexFunc(clubs, func(c store.Club) bool { return c.ID == *id }); i >= 0 {
			club = clubs[i]
		}
	}
	i := slices.IndexFunc(sheet.Entries, func(x EntryResult) bool { return x.EntryID == e.ID })
	return e, ev, formValues(i18n.Get(ev.Locale), ev, sh, club, sheet.Registration, sheet.Entries[i]), nil
}

// renderForm prints the results of an entry onto the print form of its programme; 409 when it has none.
func (s *Server) renderForm(ctx context.Context, tenantID, entryID int64) (rendered, error) {
	e, ev, values, err := s.entryFormValues(ctx, tenantID, entryID)
	if err != nil {
		return rendered{}, err
	}
	p, err := s.Store.GetProgramme(ctx, tenantID, e.ProgrammeID)
	if err != nil {
		return rendered{}, err
	}
	if p.LayoutID == nil {
		return rendered{}, conflict("This programme has no print form")
	}
	l, _, err := s.Store.GetLayout(ctx, tenantID, *p.LayoutID)
	if err != nil {
		return rendered{}, err
	}
	pdf, err := docs.RenderLayout(l.Layout, values, nil, false)
	name := values["shooter.name"]
	return rendered{
		pdf: pdf, filename: filename(l.Name, values["shooter.lastname"], values["shooter.firstname"]),
		title: l.Name + " " + name, document: store.DocForm, rangeID: printRange(ev, nil), eventID: ev.ID,
	}, err
}

// registerForms registers the print forms: editing, background scans, previews, the form of a programme and of an entry.
func (s *Server) registerForms(api huma.API) {
	type layoutsOut struct{ Body []store.Layout }
	huma.Register(api, huma.Operation{
		OperationID: "list-layouts", Metadata: modules.Metadata(modules.Forms), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/layouts", Tags: []string{"forms"},
	}, func(ctx context.Context, in *TenantPath) (*layoutsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		list, err := s.Store.ListLayouts(ctx, in.TenantID)
		return &layoutsOut{Body: list}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-layout-fields", Metadata: modules.Metadata(modules.Forms), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/layouts/fields", Tags: []string{"forms"},
		Summary: "The placeholders of print forms with example values",
	}, func(ctx context.Context, in *TenantPath) (*struct{ Body []LayoutField }, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		out := []LayoutField{}
		for k, v := range docs.SampleValues() {
			out = append(out, LayoutField{Key: k, Example: v})
		}
		slices.SortFunc(out, func(a, b LayoutField) int { return strings.Compare(a.Key, b.Key) })
		return &struct{ Body []LayoutField }{Body: out}, nil
	})

	type layoutIn struct {
		TenantPath
		Body store.LayoutInput
	}
	type layoutOut struct{ Body store.Layout }
	huma.Register(api, huma.Operation{
		OperationID: "create-layout", Metadata: modules.Metadata(modules.Forms), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/layouts", Tags: []string{"forms"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *layoutIn) (*layoutOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		in.Body.Name = strings.TrimSpace(in.Body.Name)
		id, err := s.Store.CreateLayout(ctx, in.TenantID, in.Body)
		if err != nil {
			return nil, layoutConflict(err)
		}
		l, _, err := s.Store.GetLayout(ctx, in.TenantID, id)
		return &layoutOut{Body: l}, err
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-layout", Metadata: modules.Metadata(modules.Forms), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/layouts/{layoutId}", Tags: []string{"forms"},
	}, func(ctx context.Context, in *LayoutPath) (*layoutOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		l, _, err := s.Store.GetLayout(ctx, in.TenantID, in.LayoutID)
		return &layoutOut{Body: l}, notFound(err, "Print form not found")
	})

	type layoutUpdateIn struct {
		LayoutPath
		Body store.LayoutInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-layout", Metadata: modules.Metadata(modules.Forms), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/layouts/{layoutId}", Tags: []string{"forms"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *layoutUpdateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		in.Body.Name = strings.TrimSpace(in.Body.Name)
		return nil, layoutConflict(s.Store.UpdateLayout(ctx, in.TenantID, in.LayoutID, in.Body))
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-layout", Metadata: modules.Metadata(modules.Forms), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/layouts/{layoutId}", Tags: []string{"forms"},
		Summary: "Delete a print form; programmes that use it no longer print a form", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *LayoutPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.DeleteLayout(ctx, in.TenantID, in.LayoutID), "Print form not found")
	})

	type backgroundIn struct {
		LayoutPath
		RawBody huma.MultipartFormFiles[struct {
			File huma.FormFile `form:"file" contentType:"image/png,image/jpeg,application/octet-stream" required:"true"`
		}]
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-layout-background", Metadata: modules.Metadata(modules.Forms), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/layouts/{layoutId}/background", Tags: []string{"forms"},
		Summary: "Store a scan of the empty form (PNG or JPEG) for the preview; it is never printed", DefaultStatus: http.StatusNoContent,
		MaxBodyBytes: maxBackgroundBytes + 64<<10,
	}, func(ctx context.Context, in *backgroundIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		b, err := io.ReadAll(io.LimitReader(in.RawBody.Data().File, maxBackgroundBytes+1))
		if err != nil {
			return nil, err
		}
		if len(b) > maxBackgroundBytes {
			return nil, invalid("The scan is larger than 4 MB")
		}
		if t := http.DetectContentType(b); t != "image/png" && t != "image/jpeg" {
			return nil, invalid("The scan must be a PNG or JPEG image")
		}
		return nil, notFound(s.Store.SetLayoutBackground(ctx, in.TenantID, in.LayoutID, b), "Print form not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-layout-background", Metadata: modules.Metadata(modules.Forms), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/layouts/{layoutId}/background", Tags: []string{"forms"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *LayoutPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.SetLayoutBackground(ctx, in.TenantID, in.LayoutID, nil), "Print form not found")
	})

	type imageOut struct {
		ContentType  string `header:"Content-Type"`
		CacheControl string `header:"Cache-Control"`
		Body         []byte
	}
	huma.Register(api, huma.Operation{
		OperationID: "get-layout-background", Metadata: modules.Metadata(modules.Forms), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/layouts/{layoutId}/background", Tags: []string{"forms"},
		Summary:   "The scan of the empty form, for the editor",
		Responses: map[string]*huma.Response{"200": {Description: "PNG or JPEG", Content: map[string]*huma.MediaType{"image/png": {}, "image/jpeg": {}}}},
	}, func(ctx context.Context, in *LayoutPath) (*imageOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		_, bg, err := s.Store.GetLayout(ctx, in.TenantID, in.LayoutID)
		if err != nil {
			return nil, notFound(err, "Print form not found")
		}
		if len(bg) == 0 {
			return nil, notFound(store.ErrNotFound, "The print form has no scan")
		}
		return &imageOut{ContentType: http.DetectContentType(bg), CacheControl: "no-cache", Body: bg}, nil
	})

	pdfOp := func(op, path, summary string) huma.Operation {
		return huma.Operation{
			OperationID: op, Metadata: modules.Metadata(modules.Forms), Method: http.MethodGet, Path: path, Tags: []string{"forms"}, Summary: summary,
			Responses: map[string]*huma.Response{"200": {Description: "PDF", Content: map[string]*huma.MediaType{"application/pdf": {}}}},
		}
	}
	type previewIn struct {
		LayoutPath
		EntryID int64 `query:"entryId" required:"false" doc:"Fill in the results of this entry instead of made-up values"`
		Print   bool  `query:"print" required:"false" doc:"Render as printed: with the printer offset and without the background"`
	}
	huma.Register(api, pdfOp("preview-layout", "/api/tenants/{tenantId}/layouts/{layoutId}/preview", "The form with the scan of the empty form behind it"),
		func(ctx context.Context, in *previewIn) (*PDFOutput, error) {
			if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
				return nil, err
			}
			l, bg, err := s.Store.GetLayout(ctx, in.TenantID, in.LayoutID)
			if err != nil {
				return nil, notFound(err, "Print form not found")
			}
			values := docs.SampleValues()
			if in.EntryID != 0 {
				if _, _, values, err = s.entryFormValues(ctx, in.TenantID, in.EntryID); err != nil {
					return nil, err
				}
			}
			if in.Print {
				bg = nil
			}
			pdf, err := docs.RenderLayout(l.Layout, values, bg, !in.Print)
			return pdfOut(rendered{pdf: pdf, filename: filename(l.Name, "preview")}, err)
		})

	huma.Register(api, pdfOp("entry-form", "/api/tenants/{tenantId}/entries/{entryId}/form", "The results of an entry printed onto the form of its programme"),
		func(ctx context.Context, in *EntryPath) (*PDFOutput, error) {
			if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
				return nil, err
			}
			return pdfOut(s.renderForm(ctx, in.TenantID, in.EntryID))
		})

	type formIn struct {
		ProgrammePath
		Body struct {
			LayoutID *int64 `json:"layoutId,omitempty" required:"false" doc:"A print form of the club; leave out for none"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "set-programme-form", Metadata: modules.Metadata(modules.Forms), Method: http.MethodPut, Path: "/api/tenants/{tenantId}/programmes/{programmeId}/form", Tags: []string{"programmes"},
		Summary: "Choose the print form of a programme; allowed while shooters are registered", DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *formIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		return nil, s.programmeChanged(ctx, in.TenantID, in.ProgrammeID, notFound(s.Store.SetProgrammeLayout(ctx, in.TenantID, in.ProgrammeID, in.Body.LayoutID), "Programme or print form not found"))
	})
}

// layoutConflict maps a taken name to 409 and a missing form to 404.
func layoutConflict(err error) error {
	if errors.Is(err, store.ErrExists) {
		return conflict("A print form with this name already exists")
	}
	return notFound(err, "Print form not found")
}
