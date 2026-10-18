package web

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
)

// ImportLineError is a problem with one line of an import file.
type ImportLineError struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

var importColumns = map[string]string{
	"license": "license", "licensenumber": "license",
	"firstname":  "firstname",
	"lastname":   "lastname",
	"birthday":   "birthday",
	"postalcode": "postalcode", "zip": "postalcode",
	"city": "city",
	"club": "club",
}

const maxImportErrors = 50

// openCSV skips a byte order mark, takes the delimiter (; or ,) from the header row, reads that row and returns the
// reader with the header names in lower case without spaces, dots and dashes.
func openCSV(r io.Reader) (*csv.Reader, []string, error) {
	br := bufio.NewReader(r)
	head, err := br.Peek(4096)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, bufio.ErrBufferFull) {
		return nil, nil, err
	}
	head = bytes.TrimPrefix(head, []byte("\xef\xbb\xbf"))
	firstLine, _, _ := bytes.Cut(head, []byte("\n"))
	if len(bytes.TrimSpace(firstLine)) == 0 {
		return nil, nil, errors.New("The file is empty")
	}
	cr := csv.NewReader(br)
	cr.Comma = ','
	if bytes.Count(firstLine, []byte(";")) > bytes.Count(firstLine, []byte(",")) {
		cr.Comma = ';'
	}
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true
	header, err := cr.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot read the header row: %w", err)
	}
	for i, h := range header {
		key := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(h, "\ufeff")))
		header[i] = strings.NewReplacer(" ", "", ".", "", "-", "").Replace(key)
	}
	return cr, header, nil
}

// parseShooterCSV reads a shooter list with a header row. The delimiter (; or ,) is taken from the header.
// Required columns: firstname and lastname; optional: license, birthday, postalcode, city, club.
func parseShooterCSV(r io.Reader) ([]store.ImportRow, []ImportLineError, error) {
	cr, header, err := openCSV(r)
	if err != nil {
		return nil, nil, err
	}
	col := map[string]int{}
	for i, h := range header {
		if f, ok := importColumns[strings.ReplaceAll(h, "_", "")]; ok {
			col[f] = i
		}
	}
	for _, need := range []string{"firstname", "lastname"} {
		if _, ok := col[need]; !ok {
			return nil, nil, errors.New("The header row needs at least the columns firstname and lastname")
		}
	}

	var rows []store.ImportRow
	var errs []ImportLineError
	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		line, _ := cr.FieldPos(0)
		if err != nil {
			errs = append(errs, ImportLineError{Line: line, Message: err.Error()})
			if len(errs) >= maxImportErrors {
				break
			}
			continue
		}
		get := func(f string) string {
			if i, ok := col[f]; ok && i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}
		if strings.Join(rec, "") == "" {
			continue
		}
		in := store.ShooterInput{
			LicenseNumber: domain.NormalizeLicense(get("license")),
			Firstname:     get("firstname"),
			Lastname:      get("lastname"),
			Postalcode:    get("postalcode"),
			City:          get("city"),
		}
		var msgs []string
		if in.Firstname == "" || in.Lastname == "" {
			msgs = append(msgs, "First and last name are required")
		}
		if in.Birthday, err = parseBirthday(get("birthday")); err != nil {
			msgs = append(msgs, err.Error())
		}
		if len(msgs) > 0 {
			errs = append(errs, ImportLineError{Line: line, Message: strings.Join(msgs, "; ")})
			if len(errs) >= maxImportErrors {
				break
			}
			continue
		}
		rows = append(rows, store.ImportRow{ShooterInput: in, Club: get("club")})
	}
	return rows, errs, nil
}

// parseBirthday accepts empty, YYYY-MM-DD and D.M.YYYY and returns YYYY-MM-DD.
func parseBirthday(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	for _, layout := range []string{time.DateOnly, "2.1.2006"} {
		if t, err := time.Parse(layout, s); err == nil {
			if t.Year() < 1900 || t.After(time.Now()) {
				break
			}
			return t.Format(time.DateOnly), nil
		}
	}
	return "", fmt.Errorf("Invalid birthday %q (YYYY-MM-DD or DD.MM.YYYY)", s)
}
