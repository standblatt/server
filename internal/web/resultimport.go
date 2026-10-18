package web

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/permissions"
	"github.com/standblatt/server/internal/store"
)

// ResultImportReport is the outcome of a result import. Nothing is imported when Errors is not empty.
type ResultImportReport struct {
	Imported        int               `json:"imported"`
	CreatedShooters int               `json:"createdShooters"`
	Errors          []ImportLineError `json:"errors"`
}

var (
	seriesColumn  = regexp.MustCompile(`^s(\d+)$`)
	mouchesColumn = regexp.MustCompile(`^m(\d+)$`)
	shotColumn    = regexp.MustCompile(`^s(\d+)_(\d+)$`)
)

// resultRow is one line of a result import as read from the file.
type resultRow struct {
	line                                                int
	license, firstname, lastname, birthday, club, title string
	points, mouches                                     map[int]string
	shots                                               map[[2]int]string
}

// parseResultCSV reads a result file: the shooter (license, or firstname and lastname with an optional birthday),
// an optional club, the programme and per series sN, optionally mN for mouches and sN_M for single shots.
func parseResultCSV(r io.Reader) ([]resultRow, []ImportLineError, error) {
	cr, header, err := openCSV(r)
	if err != nil {
		return nil, nil, err
	}
	col := map[string]int{}
	for i, h := range header {
		if f, ok := importColumns[strings.ReplaceAll(h, "_", "")]; ok {
			col[f] = i
		}
		if h == "programme" || h == "program" {
			col["programme"] = i
		}
	}
	if _, ok := col["programme"]; !ok {
		return nil, nil, errors.New("The header row needs the column programme")
	}
	_, lic := col["license"]
	_, first := col["firstname"]
	_, last := col["lastname"]
	if !lic && !(first && last) {
		return nil, nil, errors.New("The header row needs the column license or the columns firstname and lastname")
	}
	var rows []resultRow
	var errs []ImportLineError
	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		line, _ := cr.FieldPos(0)
		if err != nil {
			errs = append(errs, ImportLineError{Line: line, Message: err.Error()})
			continue
		}
		if strings.Join(rec, "") == "" {
			continue
		}
		get := func(f string) string {
			if i, ok := col[f]; ok && i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}
		row := resultRow{line: line, license: domain.NormalizeLicense(get("license")), firstname: get("firstname"), lastname: get("lastname"),
			club: get("club"), title: get("programme"), points: map[int]string{}, mouches: map[int]string{}, shots: map[[2]int]string{}}
		if row.birthday, err = parseBirthday(get("birthday")); err != nil {
			errs = append(errs, ImportLineError{Line: line, Message: err.Error()})
			continue
		}
		for i, h := range header {
			if i >= len(rec) || strings.TrimSpace(rec[i]) == "" {
				continue
			}
			v := strings.TrimSpace(rec[i])
			if m := seriesColumn.FindStringSubmatch(h); m != nil {
				n, _ := strconv.Atoi(m[1])
				row.points[n] = v
			} else if m := mouchesColumn.FindStringSubmatch(h); m != nil {
				n, _ := strconv.Atoi(m[1])
				row.mouches[n] = v
			} else if m := shotColumn.FindStringSubmatch(h); m != nil {
				n, _ := strconv.Atoi(m[1])
				k, _ := strconv.Atoi(m[2])
				row.shots[[2]int{n, k}] = v
			}
		}
		rows = append(rows, row)
	}
	return rows, errs, nil
}

// number reads a decimal with a dot or a comma.
func number(s string) (float64, error) {
	return strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
}

// resultSeries turns the columns of a row into the series of a programme.
func resultSeries(row resultRow, series []domain.Series) ([]domain.SeriesResult, error) {
	out := make([]domain.SeriesResult, len(series))
	for i, s := range series {
		n := i + 1
		var shots []float64
		for k := 1; s.TrackShots && k <= s.Shots; k++ {
			if v, ok := row.shots[[2]int{n, k}]; ok {
				f, err := number(v)
				if err != nil {
					return nil, fmt.Errorf("Series %d, shot %d: %q is not a number", n, k, v)
				}
				shots = append(shots, f)
			}
		}
		out[i].Shots = shots
		if v, ok := row.points[n]; ok {
			f, err := number(v)
			if err != nil {
				return nil, fmt.Errorf("Series %d: %q is not a number", n, v)
			}
			out[i].Points = f
		} else if len(shots) == 0 {
			return nil, fmt.Errorf("Series %d is missing", n)
		}
		if v, ok := row.mouches[n]; ok {
			m, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("Series %d: mouches %q are not a number", n, v)
			}
			out[i].Mouches = &m
		}
	}
	return out, nil
}

// findProgramme finds a programme of the event by barcode, by title ("Name (Variant)") or by its name alone.
func findProgramme(list []store.Programme, title string) (store.Programme, bool) {
	var byName []store.Programme
	for _, p := range list {
		if title == p.Barcode && p.Barcode != "" || strings.EqualFold(title, programmeTitle(p.Name, p.Variant)) {
			return p, true
		}
		if strings.EqualFold(title, p.Name) {
			byName = append(byName, p)
		}
	}
	if len(byName) == 1 {
		return byName[0], true
	}
	return store.Programme{}, false
}

// importShooter finds the shooter of a row by licence, else by name and birthday, and returns the ID; for a shooter
// that is not in the master data yet it returns 0 and the input to create one.
func (s *Server) importShooter(ctx context.Context, tenantID int64, row resultRow) (int64, store.ShooterInput, error) {
	q := row.license
	if q == "" {
		q = row.firstname + " " + row.lastname
	}
	list, err := s.Store.SearchShooters(ctx, tenantID, q, 0, 50, 0)
	if err != nil {
		return 0, store.ShooterInput{}, err
	}
	var hits []store.Shooter
	for _, c := range list {
		switch {
		case row.license != "":
			if c.LicenseNumber == row.license {
				hits = append(hits, c)
			}
		case strings.EqualFold(c.Firstname, row.firstname) && strings.EqualFold(c.Lastname, row.lastname) && (row.birthday == "" || c.Birthday == row.birthday):
			hits = append(hits, c)
		}
	}
	switch {
	case len(hits) == 1:
		return hits[0].ID, store.ShooterInput{}, nil
	case len(hits) > 1:
		return 0, store.ShooterInput{}, fmt.Errorf("Several shooters are called %s %s; add the birthday or the licence", row.firstname, row.lastname)
	case row.firstname == "" || row.lastname == "":
		return 0, store.ShooterInput{}, fmt.Errorf("No shooter with the licence %s; add first and last name to create one", row.license)
	}
	return 0, store.ShooterInput{LicenseNumber: row.license, Firstname: row.firstname, Lastname: row.lastname, Birthday: row.birthday}, nil
}

// newcomerKey identifies a shooter to create, so two lines of one new shooter create one.
func newcomerKey(in store.ShooterInput) string {
	return strings.ToLower(in.LicenseNumber + "|" + in.Firstname + "|" + in.Lastname + "|" + in.Birthday)
}

// registerResultImport registers the import of results from a CSV file.
func (s *Server) registerResultImport(api huma.API) {
	type importIn struct {
		EventPath
		RawBody huma.MultipartFormFiles[struct {
			File huma.FormFile `form:"file" contentType:"text/csv,text/plain,application/vnd.ms-excel,application/octet-stream" required:"true"`
		}]
	}
	type importOut struct{ Body ResultImportReport }
	huma.Register(api, huma.Operation{
		OperationID: "import-results", Metadata: needs("", permissions.EventsConfigure), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/events/{eventId}/results/import", Tags: []string{"results"},
		Summary:      "Import results from CSV (license or firstname/lastname/birthday, optional club, programme, s1…sN, optional m1…mN and sN_M shots)",
		Description:  "Every line is checked first; with any error nothing is imported. Shooters that are not in the master data are created; unknown clubs are an error.",
		MaxBodyBytes: maxImportBytes,
	}, func(ctx context.Context, in *importIn) (*importOut, error) {
		if _, _, err := s.access(ctx, in.TenantID); err != nil {
			return nil, err
		}
		if _, err := s.Store.GetEvent(ctx, in.TenantID, in.EventID); err != nil {
			return nil, notFound(err, "Event not found")
		}
		rows, errs, err := parseResultCSV(in.RawBody.Data().File)
		if err != nil {
			return nil, invalid(err.Error())
		}
		programmes, err := s.Store.ListProgrammes(ctx, in.TenantID, in.EventID)
		if err != nil {
			return nil, err
		}
		clubs, err := s.Store.ListClubs(ctx, in.TenantID)
		if err != nil {
			return nil, err
		}
		type plan struct {
			row       resultRow
			shooterID int64
			newcomer  store.ShooterInput
			clubID    *int64
			line      ResultLine
		}
		var plans []plan
		seen := map[string]int{}
		fail := func(line int, msg string) { errs = append(errs, ImportLineError{Line: line, Message: msg}) }
		for _, row := range rows {
			p, ok := findProgramme(programmes, row.title)
			if !ok {
				fail(row.line, fmt.Sprintf("Programme %q not found", row.title))
				continue
			}
			series, err := resultSeries(row, p.Series)
			if err != nil {
				fail(row.line, err.Error())
				continue
			}
			domain.NormalizeResults(series)
			if problems := domain.ValidateResults(p.Series, series); len(problems) > 0 {
				fail(row.line, domain.Join(problems, "; "))
				continue
			}
			shooterID, newcomer, err := s.importShooter(ctx, in.TenantID, row)
			if err != nil {
				fail(row.line, err.Error())
				continue
			}
			var clubID *int64
			if row.club != "" {
				m := matchClub(row.club, clubs)
				if m.Match == nil {
					fail(row.line, fmt.Sprintf("Club %q is not in the club list", row.club))
					continue
				}
				clubID = &m.Match.ID
			}
			key := fmt.Sprintf("%s|%d", newcomerKey(newcomer), p.ID)
			completed := 0
			if shooterID != 0 {
				key = fmt.Sprintf("#%d|%d", shooterID, p.ID)
				if completed, err = s.Store.CompletedAttempts(ctx, in.TenantID, in.EventID, shooterID, p.ID); err != nil {
					return nil, err
				}
			}
			limit := p.AttemptLimit()
			switch {
			case limit == 1 && completed > 0:
				fail(row.line, fmt.Sprintf("The shooter already has results in %s", programmeTitle(p.Name, p.Variant)))
				continue
			case limit == 1 && seen[key] > 0:
				fail(row.line, "The same shooter and programme appear twice in the file")
				continue
			case limit > 0 && completed+seen[key] >= limit:
				fail(row.line, fmt.Sprintf("The shooter has no attempts left in %s", programmeTitle(p.Name, p.Variant)))
				continue
			}
			seen[key]++
			plans = append(plans, plan{row: row, shooterID: shooterID, newcomer: newcomer, clubID: clubID, line: ResultLine{ProgrammeID: p.ID, ClubID: clubID, Series: series}})
		}
		out := &importOut{Body: ResultImportReport{Errors: errs}}
		if out.Body.Errors == nil {
			out.Body.Errors = []ImportLineError{}
		}
		if len(errs) > 0 {
			return out, nil
		}
		// ponytail: no transaction around the whole file; every line was checked above, a failure now keeps the lines before it.
		created := map[string]int64{}
		for _, pl := range plans {
			id := pl.shooterID
			if id == 0 {
				key := newcomerKey(pl.newcomer)
				if id = created[key]; id == 0 {
					var clubIDs []int64
					if pl.clubID != nil {
						clubIDs = []int64{*pl.clubID}
					}
					if id, err = s.Store.CreateShooter(ctx, in.TenantID, pl.newcomer, clubIDs); err != nil {
						return nil, fmt.Errorf("line %d: %w", pl.row.line, err)
					}
					created[key] = id
					out.Body.CreatedShooters++
				}
			}
			if pl.clubID != nil {
				if err := s.Store.AddShooterClub(ctx, in.TenantID, id, *pl.clubID); err != nil {
					return nil, fmt.Errorf("line %d: %w", pl.row.line, err)
				}
			}
			pl.line.ShooterID = id
			if _, err := s.enterResult(ctx, in.TenantID, in.EventID, pl.line); err != nil {
				return nil, fmt.Errorf("line %d: %w", pl.row.line, err)
			}
			out.Body.Imported++
		}
		s.Log.Info("results imported", "tenant", in.TenantID, "event", in.EventID, "lines", out.Body.Imported, "shooters", out.Body.CreatedShooters)
		return out, nil
	})
}
