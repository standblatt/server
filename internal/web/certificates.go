package web

import (
	"context"
	"strconv"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/docs"
	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/i18n"
	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/store"
)

// Certificates chooses the ranked lines printed onto a print form with the rank as data source: of a programme or
// of a classification (optionally one class), ranks From to To.
type Certificates struct {
	LayoutID         int64  `json:"layoutId,omitempty" query:"layout" required:"false" doc:"The print form; required for certificates and registration forms"`
	ProgrammeID      int64  `json:"programmeId,omitempty" query:"programme" required:"false"`
	ClassificationID int64  `json:"classificationId,omitempty" query:"classification" required:"false"`
	Class            string `json:"class,omitempty" query:"class" required:"false" doc:"Only this class of the classification"`
	From             int    `json:"from,omitempty" query:"from" required:"false" minimum:"0" doc:"First rank, 1 when 0"`
	To               int    `json:"to,omitempty" query:"to" required:"false" minimum:"0" doc:"Last rank, 3 when 0"`
}

// sourceForm loads a print form of the club that must print the given data source, with its images.
func (s *Server) sourceForm(ctx context.Context, tenantID, layoutID int64, source string) (store.Layout, map[int64][]byte, error) {
	l, err := s.Store.GetLayout(ctx, tenantID, layoutID)
	if err != nil {
		return l, nil, notFound(err, "Print form not found")
	}
	if l.DataSource() != source {
		return l, nil, invalid("The print form %s prints another data source", l.Name)
	}
	images, err := s.Store.ImagesByID(ctx, tenantID, l.ImageIDs())
	return l, images, err
}

// renderRegistrationForm prints a registration onto a print form with the registration as data source.
func (s *Server) renderRegistrationForm(ctx context.Context, tenantID, registrationID, layoutID int64) (rendered, error) {
	reg, ev, programmes, err := s.registrationContext(ctx, tenantID, registrationID)
	if err != nil {
		return rendered{}, err
	}
	l, images, err := s.sourceForm(ctx, tenantID, layoutID, domain.SourceRegistration)
	if err != nil {
		return rendered{}, err
	}
	sh, club, err := s.shooterClub(ctx, tenantID, reg)
	if err != nil {
		return rendered{}, err
	}
	pdf, err := docs.RenderLayout(l.Layout, []map[string]string{registrationValues(i18n.Get(ev.Locale), ev, sh, club, reg, programmes)}, nil, images, false)
	return rendered{
		pdf: pdf, filename: filename(l.Name, reg.Shooter.Lastname, reg.Shooter.Firstname), title: l.Name + " " + reg.Shooter.Lastname + " " + reg.Shooter.Firstname,
		document: store.DocForm, rangeID: printRange(ev, &reg), eventID: ev.ID,
	}, err
}

// renderCertificates prints one page per ranked shooter in the chosen ranks onto a print form with the rank as data
// source; 409 when nobody is ranked there.
func (s *Server) renderCertificates(ctx context.Context, tenantID, eventID int64, c Certificates) (rendered, error) {
	from, to := c.From, c.To
	if from == 0 {
		from = 1
	}
	if to == 0 {
		to = 3
	}
	ev, err := s.Store.GetEvent(ctx, tenantID, eventID)
	if err != nil {
		return rendered{}, notFound(err, "Event not found")
	}
	l, images, err := s.sourceForm(ctx, tenantID, c.LayoutID, domain.SourceRank)
	if err != nil {
		return rendered{}, err
	}
	lang := i18n.Get(ev.Locale)
	var pages []map[string]string
	add := func(v map[string]string, rank int, total float64, ranking, class string) {
		v["rank"], v["total"], v["ranking"], v["class"] = strconv.Itoa(rank), docs.Points(total), ranking, class
		pages = append(pages, v)
	}
	switch {
	case c.ProgrammeID != 0:
		rankings, err := s.rankings(ctx, tenantID, eventID, false)
		if err != nil {
			return rendered{}, err
		}
		for _, r := range rankings {
			if r.ProgrammeID != c.ProgrammeID {
				continue
			}
			for _, line := range r.Lines {
				if line.Rank < from || line.Rank > to {
					continue
				}
				f, err := s.entryFormValues(ctx, tenantID, line.EntryID)
				if err != nil {
					return rendered{}, err
				}
				add(f.values, line.Rank, line.Total, r.Title, "")
			}
		}
	case c.ClassificationID != 0:
		list, err := s.classificationRankings(ctx, tenantID, eventID, false)
		if err != nil {
			return rendered{}, err
		}
		clubs, err := s.Store.ListClubs(ctx, tenantID)
		if err != nil {
			return rendered{}, err
		}
		clubByName := map[string]store.Club{}
		for _, cl := range clubs {
			clubByName[cl.Name] = cl
		}
		for _, cr := range list {
			if cr.ID != c.ClassificationID {
				continue
			}
			for _, cl := range cr.Classes {
				if c.Class != "" && cl.Name != c.Class {
					continue
				}
				for _, line := range cl.Lines {
					if line.Rank < from || line.Rank > to {
						continue
					}
					sh, err := s.Store.GetShooter(ctx, tenantID, line.ShooterID)
					if err != nil {
						return rendered{}, err
					}
					p := docs.Person{Lastname: line.Lastname, Firstname: line.Firstname, Birthday: swissDate(sh.Birthday), BirthYear: line.BirthYear,
						License: sh.LicenseNumber, Club: line.Club, Category: categoryLabel(lang, line.Category), CategoryShort: categoryShort[line.Category]}
					add(baseValues(ev, sh, clubByName[line.Club], p), line.Rank, line.Total, cr.Name, cl.Name)
				}
			}
		}
	default:
		return rendered{}, invalid("Choose a programme or a classification")
	}
	if len(pages) == 0 {
		return rendered{}, conflict("Nobody is ranked in these ranks")
	}
	pdf, err := docs.RenderLayout(l.Layout, pages, nil, images, false)
	return rendered{
		pdf: pdf, filename: filename(l.Name, ev.Name), title: l.Name + " " + ev.Name,
		document: store.DocCertificate, rangeID: printRange(ev, nil), eventID: ev.ID,
	}, err
}

// registerCertificates registers the PDFs of registration documents and certificates.
func (s *Server) registerCertificates(api huma.API) {
	type regFormIn struct {
		RegistrationPath
		LayoutID int64 `path:"layoutId"`
	}
	huma.Register(api, withModule(pdfOp("registration-form", "forms", "/api/tenants/{tenantId}/registrations/{registrationId}/forms/{layoutId}", "A registration printed onto a print form with the registration as data source"), modules.Forms),
		func(ctx context.Context, in *regFormIn) (*PDFOutput, error) {
			if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
				return nil, err
			}
			return pdfOut(s.renderRegistrationForm(ctx, in.TenantID, in.RegistrationID, in.LayoutID))
		})

	type certificatesIn struct {
		EventPath
		Certificates
	}
	huma.Register(api, withModule(pdfOp("certificates", "forms", "/api/tenants/{tenantId}/events/{eventId}/certificates", "Certificates: one page per ranked shooter of a programme or classification"), modules.Forms),
		func(ctx context.Context, in *certificatesIn) (*PDFOutput, error) {
			if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
				return nil, err
			}
			return pdfOut(s.renderCertificates(ctx, in.TenantID, in.EventID, in.Certificates))
		})
}
