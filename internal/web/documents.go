package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/docs"
	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/i18n"
	"github.com/standblatt/server/internal/realtime"
	"github.com/standblatt/server/internal/store"
)

// categoryLabel returns the translated label of an age category; youth categories print as their code, e.g. U21.
func categoryLabel(l i18n.Lang, c domain.AgeCategory) string {
	switch c {
	case domain.Elite:
		return i18n.T(l, "Elite")
	case domain.Senior:
		return i18n.T(l, "Senior")
	case domain.Veteran:
		return i18n.T(l, "Veteran")
	case domain.SeniorVeteran:
		return i18n.T(l, "Senior veteran")
	}
	return strings.ToUpper(string(c))
}

var categoryShort = map[domain.AgeCategory]string{
	domain.U10: "U10", domain.U13: "U13", domain.U15: "U15", domain.U17: "U17", domain.U19: "U19", domain.U21: "U21",
	domain.Elite: "E", domain.Senior: "S", domain.Veteran: "V", domain.SeniorVeteran: "SV",
}

// swissDate formats a YYYY-MM-DD date as DD.MM.YYYY, empty when it does not parse.
func swissDate(iso string) string {
	if t, err := time.Parse(time.DateOnly, iso); err == nil {
		return t.Format("02.01.2006")
	}
	return ""
}

// rendered is a PDF together with where it may be printed.
type rendered struct {
	pdf      []byte
	filename string
	title    string
	document string
	rangeID  int64
	eventID  int64
}

// swissDateSpan writes the first and last of the days, "08.01.2027" or "20.12.2026–10.01.2027".
func swissDateSpan(days []string) string {
	if len(days) == 0 {
		return ""
	}
	first, last := swissDate(days[0]), swissDate(days[len(days)-1])
	if first == last {
		return first
	}
	return first + "–" + last
}

// printRange is the range a registration's documents print on: that of its lane, else the event's first.
func printRange(ev store.Event, reg *store.Registration) int64 {
	if reg != nil && reg.RangeID != nil {
		return *reg.RangeID
	}
	if len(ev.Ranges) > 0 {
		return ev.Ranges[0].ID
	}
	return 0
}

// docEvent returns the event header of documents with its days and language.
func docEvent(ev store.Event) docs.Event {
	return docs.Event{Name: ev.Name, Date: swissDateSpan(ev.Days), Organizer: ev.Organizer, Lang: i18n.Get(ev.Locale)}
}

// docPerson returns the shooter block of documents for a registration.
func docPerson(l i18n.Lang, r store.Registration) docs.Person {
	p := docs.Person{
		Lastname: r.Shooter.Lastname, Firstname: r.Shooter.Firstname, Birthday: swissDate(r.Shooter.Birthday),
		License: r.Shooter.LicenseNumber, Club: r.ClubName, Category: categoryLabel(l, r.AgeCategory), CategoryShort: categoryShort[r.AgeCategory],
	}
	if len(r.Shooter.Birthday) >= 4 {
		p.BirthYear = r.Shooter.Birthday[:4]
	}
	return p
}

// filename joins the ASCII letters and digits of parts with underscores into a PDF file name.
func filename(parts ...string) string {
	var b strings.Builder
	for _, p := range parts {
		for _, r := range p {
			switch {
			case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
				b.WriteRune(r)
			case b.Len() > 0 && !strings.HasSuffix(b.String(), "_"):
				b.WriteByte('_')
			}
		}
		if b.Len() > 0 && !strings.HasSuffix(b.String(), "_") {
			b.WriteByte('_')
		}
	}
	return strings.TrimSuffix(b.String(), "_") + ".pdf"
}

// registrationContext loads a registration with its event and the event's programmes, 404 when the registration is missing.
func (s *Server) registrationContext(ctx context.Context, tenantID, registrationID int64) (store.Registration, store.Event, []store.Programme, error) {
	reg, err := s.Store.GetRegistration(ctx, tenantID, registrationID)
	if err != nil {
		return reg, store.Event{}, nil, notFound(err, "Registration not found")
	}
	ev, err := s.Store.GetEvent(ctx, tenantID, reg.EventID)
	if err != nil {
		return reg, ev, nil, err
	}
	programmes, err := s.Store.ListProgrammes(ctx, tenantID, reg.EventID)
	return reg, ev, programmes, err
}

// renderRegistrationDoc renders the lane sheet, control strip or result sheet of a registration.
func (s *Server) renderRegistrationDoc(ctx context.Context, tenantID, registrationID int64, document string) (rendered, error) {
	reg, ev, programmes, err := s.registrationContext(ctx, tenantID, registrationID)
	if err != nil {
		return rendered{}, err
	}
	byID := map[int64]store.Programme{}
	for _, p := range programmes {
		byID[p.ID] = p
	}
	l := i18n.Get(ev.Locale)
	person := docPerson(l, reg)
	out := rendered{document: document, rangeID: printRange(ev, &reg), eventID: ev.ID}
	name := reg.Shooter.Lastname + " " + reg.Shooter.Firstname
	switch document {
	case store.DocLaneSheet:
		sheet := docs.LaneSheet{Event: docEvent(ev), Person: person, RegistrationID: reg.ID, Lane: reg.Lane}
		if reg.RangeID != nil && len(ev.Ranges) > 1 {
			for _, r := range ev.Ranges {
				if r.ID == *reg.RangeID {
					sheet.Range = r.Name
				}
			}
		}
		for _, e := range reg.Entries {
			p := byID[e.ProgrammeID]
			lp := docs.LaneSheetProgramme{Title: programmeTitle(p.Name, p.Variant), Shots: p.Shots, Barcode: p.Barcode}
			if len(e.CombinedSeries) > 0 && p.CombinableWithID != nil {
				t := byID[*p.CombinableWithID]
				lp.CombinedWith, lp.Shots = programmeTitle(t.Name, t.Variant), 0
			}
			sheet.Programmes = append(sheet.Programmes, lp)
		}
		out.pdf, err = docs.RenderLaneSheet(sheet)
		out.filename, out.title = filename("Lane_sheet", reg.Shooter.LicenseNumber, reg.Shooter.Lastname), i18n.Tf(l, "Lane sheet %s", name)
	case store.DocControlStrip:
		var passes []domain.PassesEntry
		for _, e := range reg.Entries {
			p := byID[e.ProgrammeID]
			passes = append(passes, domain.PassesEntry{ProgrammeID: p.ID, Series: len(p.Series), CombinableWithID: p.CombinableWithID, CombinedSeries: e.CombinedSeries})
		}
		out.pdf, err = docs.RenderControlStrip(docs.ControlStrip{EventName: ev.Name, Person: person, RegistrationID: reg.ID, Passes: domain.PassesText(passes)})
		out.filename, out.title = filename("Control_strip", reg.Shooter.LicenseNumber, reg.Shooter.Lastname), i18n.Tf(l, "Control strip %s", name)
	case store.DocResultSheet:
		sheet, err := s.resultSheet(ctx, tenantID, reg.ID)
		if err != nil {
			return out, err
		}
		rs := docs.ResultSheet{Event: docEvent(ev), Person: person}
		for _, er := range sheet.Entries {
			if er.Status != "completed" {
				return out, conflict("Not all programmes have results yet")
			}
			points := make([]float64, len(er.Results))
			for i, r := range er.Results {
				points[i] = r.Points
			}
			award := ""
			if er.Award != nil {
				award = er.Award.Name
			}
			rs.Programmes = append(rs.Programmes, docs.ResultSheetProgramme{Title: er.Title, Series: points, Total: er.Total, MaxPoints: er.MaxPoints, Award: award})
		}
		out.pdf, err = docs.RenderResultSheet(rs)
		out.filename, out.title = filename("Result_sheet", reg.Shooter.Lastname, reg.Shooter.Firstname), i18n.Tf(l, "Result sheet %s", name)
	default:
		return out, invalid("Unknown document %s", document)
	}
	return out, err
}

// renderWreathCard renders the wreath card of an entry, 409 when its award is not a wreath card.
func (s *Server) renderWreathCard(ctx context.Context, tenantID, entryID int64) (rendered, error) {
	e, err := s.Store.GetEntry(ctx, tenantID, entryID)
	if err != nil {
		return rendered{}, notFound(err, "Programme of the registration not found")
	}
	sheet, err := s.resultSheet(ctx, tenantID, e.RegistrationID)
	if err != nil {
		return rendered{}, err
	}
	var award *domain.Tier
	for _, er := range sheet.Entries {
		if er.EntryID == e.ID {
			award = er.Award
		}
	}
	if award == nil || award.Kind != domain.AwardCard {
		return rendered{}, conflict("This programme has no wreath card award")
	}
	ev, err := s.Store.GetEvent(ctx, tenantID, e.EventID)
	if err != nil {
		return rendered{}, err
	}
	l := i18n.Get(ev.Locale)
	p := docPerson(l, sheet.Registration)
	pdf, err := docs.RenderWreathCard(docs.WreathCard{
		Lastname: p.Lastname, Firstname: p.Firstname, BirthYear: p.BirthYear, Club: p.Club, Date: swissDate(ev.StartsOn),
		Value: award.Value, ShowValue: ev.Settings.PrintAwardValue, ValueWords: domain.AmountInWords(ev.Locale, award.Value), Organizer: ev.Organizer,
	})
	return rendered{
		pdf: pdf, filename: filename("Wreath_card", p.Lastname, p.Firstname), title: i18n.Tf(l, "Wreath card %s", p.Lastname+" "+p.Firstname),
		document: store.DocWreathCard, rangeID: printRange(ev, &sheet.Registration), eventID: ev.ID,
	}, err
}

// renderRanking renders the ranking list PDF of a programme, kiosk-hidden programmes included.
func (s *Server) renderRanking(ctx context.Context, tenantID, programmeID int64) (rendered, error) {
	p, err := s.Store.GetProgramme(ctx, tenantID, programmeID)
	if err != nil {
		return rendered{}, notFound(err, "Programme not found")
	}
	ev, err := s.Store.GetEvent(ctx, tenantID, p.EventID)
	if err != nil {
		return rendered{}, err
	}
	rankings, err := s.rankings(ctx, tenantID, p.EventID, false)
	if err != nil {
		return rendered{}, err
	}
	list := docs.RankingList{Event: docEvent(ev), Title: programmeTitle(p.Name, p.Variant), SeriesCount: len(p.Series), Generated: time.Now()}
	if loc, err := time.LoadLocation(ev.Timezone); err == nil {
		list.Generated = list.Generated.In(loc)
	}
	for _, r := range rankings {
		if r.ProgrammeID != p.ID {
			continue
		}
		for _, l := range r.Lines {
			line := docs.RankingLine{
				Rank: l.Rank, Lastname: l.Lastname, Firstname: l.Firstname, BirthYear: l.BirthYear,
				CategoryShort: categoryShort[l.Category], Club: l.Club, Series: l.Series, Total: l.Total,
			}
			if r.Repeatable {
				line.Note = fmt.Sprintf("(%d/%d)", l.Attempt, l.Attempts)
			}
			list.Lines = append(list.Lines, line)
		}
	}
	pdf, err := docs.RenderRankingList(list)
	return rendered{pdf: pdf, filename: filename("Ranking", p.Name, p.Variant)}, err
}

// PDFOutput is a PDF shown in the browser.
type PDFOutput struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	Body               []byte
}

// pdfOut turns a rendered document into an inline PDF response.
func pdfOut(r rendered, err error) (*PDFOutput, error) {
	if err != nil {
		return nil, err
	}
	return &PDFOutput{ContentType: "application/pdf", ContentDisposition: fmt.Sprintf(`inline; filename="%s"`, r.filename), Body: r.pdf}, nil
}

// PrintRequest names what to print: a registration document, or the wreath card or print form of an entry.
type PrintRequest struct {
	Document       string `json:"document" enum:"lane_sheet,control_strip,result_sheet,wreath_card,form"`
	RegistrationID int64  `json:"registrationId,omitempty" required:"false"`
	EntryID        int64  `json:"entryId,omitempty" required:"false"`
}

// PrintResult tells whether the document went to a printer; if not, the client shows the PDF instead.
type PrintResult struct {
	Queued bool  `json:"queued"`
	JobID  int64 `json:"jobId,omitempty"`
}

// registerDocuments registers the PDF endpoints and the print endpoint.
func (s *Server) registerDocuments(api huma.API) {
	pdfOp := func(op, path string) huma.Operation {
		return huma.Operation{
			OperationID: op, Method: http.MethodGet, Path: path, Tags: []string{"documents"},
			Responses: map[string]*huma.Response{"200": {Description: "PDF", Content: map[string]*huma.MediaType{"application/pdf": {}}}},
		}
	}
	type regDocIn struct {
		RegistrationPath
		Document string `path:"document" enum:"lane_sheet,control_strip,result_sheet"`
	}
	huma.Register(api, pdfOp("registration-document", "/api/tenants/{tenantId}/registrations/{registrationId}/documents/{document}"),
		func(ctx context.Context, in *regDocIn) (*PDFOutput, error) {
			if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
				return nil, err
			}
			if err := s.helperDocument(ctx, in.TenantID, in.RegistrationID, 0, in.Document); err != nil {
				return nil, err
			}
			return pdfOut(s.renderRegistrationDoc(ctx, in.TenantID, in.RegistrationID, in.Document))
		})
	huma.Register(api, pdfOp("wreath-card", "/api/tenants/{tenantId}/entries/{entryId}/wreath-card"),
		func(ctx context.Context, in *EntryPath) (*PDFOutput, error) {
			if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
				return nil, err
			}
			return pdfOut(s.renderWreathCard(ctx, in.TenantID, in.EntryID))
		})
	huma.Register(api, pdfOp("ranking-pdf", "/api/tenants/{tenantId}/programmes/{programmeId}/ranking"),
		func(ctx context.Context, in *ProgrammePath) (*PDFOutput, error) {
			if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
				return nil, err
			}
			return pdfOut(s.renderRanking(ctx, in.TenantID, in.ProgrammeID))
		})

	type printIn struct {
		TenantPath
		Body PrintRequest
	}
	type printOut struct{ Body PrintResult }
	huma.Register(api, huma.Operation{
		OperationID: "print", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/print", Tags: []string{"documents"},
		Summary: "Send a document to the printer of its range; queued is false when no printer is set up or no agent is online",
	}, func(ctx context.Context, in *printIn) (*printOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		if err := s.helperDocument(ctx, in.TenantID, in.Body.RegistrationID, in.Body.EntryID, in.Body.Document); err != nil {
			return nil, err
		}
		var r rendered
		var err error
		switch in.Body.Document {
		case store.DocWreathCard:
			r, err = s.renderWreathCard(ctx, in.TenantID, in.Body.EntryID)
		case store.DocForm:
			r, err = s.renderForm(ctx, in.TenantID, in.Body.EntryID)
		default:
			r, err = s.renderRegistrationDoc(ctx, in.TenantID, in.Body.RegistrationID, in.Body.Document)
		}
		if err != nil {
			return nil, err
		}
		id, err := s.Store.EnqueuePrint(ctx, in.TenantID, r.rangeID, &r.eventID, r.document, r.title, r.pdf)
		if errors.Is(err, store.ErrNoPrinter) {
			return &printOut{Body: PrintResult{}}, nil
		}
		if err != nil {
			return nil, err
		}
		s.notify(ctx, in.TenantID, r.eventID, realtime.TopicPrint)
		return &printOut{Body: PrintResult{Queued: true, JobID: id}}, nil
	})
}
