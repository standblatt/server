// Package docs renders the printed documents as PDF: lane sheet, control strip label, result sheet,
// wreath card and ranking list. It works on plain data and knows nothing about storage or HTTP.
package docs

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boombuler/barcode/qr"
	"github.com/go-pdf/fpdf"
	"github.com/go-pdf/fpdf/contrib/barcode"

	"github.com/standblatt/server/internal/i18n"
)

// Event is the event as printed in document headers; Lang is the language of the documents.
type Event struct {
	Name      string
	Date      string
	Organizer string
	Lang      i18n.Lang
}

// Person is the shooter as printed on a document.
type Person struct {
	Lastname      string
	Firstname     string
	Birthday      string
	BirthYear     string
	License       string
	Club          string
	Category      string
	CategoryShort string
}

// Code is the content of the registration QR code that scanners type into the search field.
func Code(registrationID int64, license string) string {
	if license == "" {
		license = "0"
	}
	return fmt.Sprintf("REGCODEAAAA%dAAAAL%s", registrationID, license)
}

type doc struct {
	*fpdf.Fpdf
	tr func(string) string
}

// newDoc starts a PDF in millimetres, in a standard size or, when size is empty, a custom w×h page.
func newDoc(orientation, size string, w, h float64) *doc {
	var pdf *fpdf.Fpdf
	if size == "" {
		pdf = fpdf.NewCustom(&fpdf.InitType{UnitStr: "mm", Size: fpdf.SizeType{Wd: w, Ht: h}, OrientationStr: orientation})
	} else {
		pdf = fpdf.New(orientation, "mm", size, "")
	}
	pdf.SetCreator("Standblatt", true)
	pdf.SetAutoPageBreak(false, 0)
	return &doc{Fpdf: pdf, tr: pdf.UnicodeTranslatorFromDescriptor("")}
}

// text writes s in a cell at x, y in Helvetica, translated to cp1252.
func (d *doc) text(x, y, w, h float64, style string, size float64, align, s string) {
	d.SetFont("Helvetica", style, size)
	d.SetXY(x, y)
	d.CellFormat(w, h, d.tr(s), "", 0, align, false, 0, "")
}

// fit shortens s with an ellipsis until it fits into width w.
func (d *doc) fit(s string, w float64) string {
	if d.GetStringWidth(d.tr(s)) <= w {
		return s
	}
	r := []rune(s)
	for len(r) > 0 && d.GetStringWidth(d.tr(string(r)+"…")) > w {
		r = r[:len(r)-1]
	}
	return string(r) + "…"
}

// qr draws a QR code of code as a square of the given size at x, y.
func (d *doc) qr(code string, x, y, size float64) {
	key := barcode.RegisterQR(d.Fpdf, code, qr.H, qr.Unicode)
	barcode.Barcode(d.Fpdf, key, x, y, size, size, false)
}

// bytes returns the finished PDF.
func (d *doc) bytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := d.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Points formats a score: whole numbers without decimals, otherwise one decimal.
func Points(v float64) string {
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', 1, 64)
}

// digits keeps only the digits of s and pads them to an even length with a leading zero, as I25 needs.
func digits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	if b.Len()%2 == 1 {
		return "0" + b.String()
	}
	return b.String()
}

// header draws the title and the event name, date and organizer with a rule below; it returns the next y.
func header(d *doc, title string, ev Event, left, right, y float64) float64 {
	d.text(left, y, 70, 9, "B", 16, "L", title)
	d.text(left+60, y, right-left-60, 6, "B", 11, "R", ev.Name)
	d.text(left+60, y+6, right-left-60, 5, "", 9, "R", ev.Date)
	if ev.Organizer != "" {
		d.text(left+60, y+11, right-left-60, 5, "", 9, "R", ev.Organizer)
	}
	y += 18
	d.Line(left, y, right, y)
	return y + 4
}

// personRows draws the name, club, category and licence rows of a shooter; it returns the next y.
func personRows(d *doc, l i18n.Lang, p Person, left, y, width float64, nameLine string) float64 {
	rows := [][2]string{{i18n.T(l, "Name"), nameLine}, {i18n.T(l, "Club"), p.Club}, {i18n.T(l, "Category"), p.Category}, {i18n.T(l, "License"), orDash(p.License)}}
	for i, r := range rows {
		d.text(left, y, 25, 5, "", 9, "L", r[0]+":")
		style := ""
		if i == 0 {
			style = "B"
		}
		d.SetFont("Helvetica", style, 9)
		d.text(left+25, y, width-25, 5, style, 9, "L", d.fit(r[1], width-26))
		y += 5.5
	}
	return y
}

// orDash returns s, or a dash when s is empty.
func orDash(s string) string {
	if s == "" {
		return "–"
	}
	return s
}

// LaneSheetProgramme is a programme printed on the lane sheet.
type LaneSheetProgramme struct {
	Title        string
	Shots        int
	Barcode      string
	CombinedWith string
}

// LaneSheet is the A5 sheet that goes with the shooter to the lane.
type LaneSheet struct {
	Event          Event
	Person         Person
	RegistrationID int64
	Range          string
	Lane           *int
	Programmes     []LaneSheetProgramme
}

// RenderLaneSheet renders the lane sheet on A5 portrait.
func RenderLaneSheet(s LaneSheet) ([]byte, error) {
	d := newDoc("P", "A5", 0, 0)
	d.AddPage()
	left, right := 10.0, 138.0
	l := s.Event.Lang
	y := header(d, i18n.T(l, "LANE SHEET"), s.Event, left, right, 10)
	code := Code(s.RegistrationID, s.Person.License)
	d.qr(code, right-26, y, 26)
	d.text(right-28, y+26.5, 30, 3, "", 6, "C", fmt.Sprintf("ID %d", s.RegistrationID))
	y = personRows(d, l, s.Person, left, y+2, right-left-30, s.Person.Lastname+" "+s.Person.Firstname+", "+s.Person.Birthday)
	y += 8
	d.Line(left, y, right, y)
	y += 4
	if s.Range != "" {
		d.text(left, y, right-left, 6, "B", 11, "L", d.fit(s.Range, right-left))
		y += 7
	}
	if s.Lane != nil {
		d.text(left, y, right-left, 16, "B", 36, "L", i18n.T(l, "LANE")+fmt.Sprintf("  %d", *s.Lane))
		y += 20
	} else {
		d.text(left, y, right-left, 10, "B", 18, "L", i18n.T(l, "LANE")+"  ________")
		y += 14
	}
	d.Line(left, y, right, y)
	y += 4
	d.text(left, y, 60, 6, "B", 11, "L", i18n.T(l, "PROGRAMMES"))
	y += 9
	for _, p := range s.Programmes {
		d.Rect(left, y+1, 3.5, 3.5, "D")
		d.text(left+6, y, 70, 6, "", 10, "L", d.fit(i18n.Tf(l, "%s – %d shots", p.Title, p.Shots), 70))
		switch {
		case p.CombinedWith != "":
			d.text(left+70, y, right-left-70, 6, "I", 9, "R", i18n.Tf(l, "Combined with %s", p.CombinedWith))
			y += 10
		case digits(p.Barcode) != "":
			key := barcode.RegisterTwoOfFive(d.Fpdf, digits(p.Barcode), true)
			barcode.Barcode(d.Fpdf, key, left+72, y, right-left-72, 12, false)
			d.text(left+72, y+12.5, right-left-72, 3, "", 7, "C", digits(p.Barcode))
			y += 19
		default:
			y += 10
		}
	}
	return d.bytes()
}

// ControlStrip is the 90 × 29 mm label stuck on the target card.
type ControlStrip struct {
	EventName      string
	Person         Person
	RegistrationID int64
	Passes         string
}

// RenderControlStrip renders the control strip label.
func RenderControlStrip(s ControlStrip) ([]byte, error) {
	d := newDoc("P", "", 90.3, 29)
	d.AddPage()
	w := 57.0
	d.SetFont("Helvetica", "I", 7)
	d.text(2, 1.5, w, 3, "I", 7, "L", d.fit(s.EventName, w))
	d.SetFont("Helvetica", "B", 9)
	d.text(2, 5, w, 4, "B", 9, "L", d.fit(s.Person.Lastname+" "+s.Person.Firstname+", "+s.Person.Birthday, w))
	d.text(2, 9.5, w, 3.5, "", 8, "L", orDash(s.Person.License))
	d.SetFont("Helvetica", "", 8)
	d.text(2, 13.5, w, 3.5, "", 8, "L", d.fit(s.Person.Club, w))
	d.SetFont("Helvetica", "B", 13)
	catW := d.GetStringWidth(s.Person.CategoryShort) + 2
	d.text(2, 20, catW, 6, "B", 13, "L", s.Person.CategoryShort)
	d.SetFont("Helvetica", "", 7)
	d.SetXY(2+catW+2, 20)
	d.MultiCell(w-catW-2, 3, d.tr(s.Passes), "", "L", false)
	d.qr(Code(s.RegistrationID, s.Person.License), 63, 1.5, 24)
	d.text(62, 25.5, 26, 2.5, "B", 5, "C", fmt.Sprintf("ID %d", s.RegistrationID))
	return d.bytes()
}

// ResultSheetProgramme is a programme with its results on the result sheet.
type ResultSheetProgramme struct {
	Title     string
	Series    []float64
	Total     float64
	MaxPoints int
	Award     string
}

// ResultSheet is the A5 result sheet handed to the shooter.
type ResultSheet struct {
	Event      Event
	Person     Person
	Programmes []ResultSheetProgramme
}

// RenderResultSheet renders the result sheet on A5 portrait.
func RenderResultSheet(s ResultSheet) ([]byte, error) {
	d := newDoc("P", "A5", 0, 0)
	left, right := 10.0, 138.0
	d.AddPage()
	l := s.Event.Lang
	y := header(d, i18n.T(l, "RESULT SHEET"), s.Event, left, right, 10)
	name := s.Person.Lastname + " " + s.Person.Firstname
	if s.Person.BirthYear != "" {
		name += ", " + s.Person.BirthYear
	}
	y = personRows(d, l, s.Person, left, y+2, right-left, name)
	y += 4
	d.Line(left, y, right, y)
	y += 5
	for _, p := range s.Programmes {
		if y > 180 {
			d.AddPage()
			y = 10
		}
		d.text(left, y, right-left, 6, "B", 10, "L", p.Title)
		y += 7
		x := left
		for i, v := range p.Series {
			if x+14 > right-30 {
				x = left
				y += 13
			}
			d.Rect(x, y, 13, 11, "D")
			d.text(x, y+0.5, 13, 4, "", 7, "C", i18n.Tf(l, "S%d", i+1))
			d.text(x, y+4.5, 13, 6, "B", 10, "C", Points(v))
			x += 14.5
		}
		d.Rect(right-28, y, 28, 11, "D")
		d.text(right-28, y+0.5, 28, 4, "", 7, "C", i18n.T(l, "TOTAL"))
		d.text(right-28, y+4.5, 28, 6, "B", 10, "C", fmt.Sprintf("%s / %d", Points(p.Total), p.MaxPoints))
		y += 12
		if p.Award != "" {
			d.text(right-60, y, 60, 4, "", 8, "R", i18n.Tf(l, "Award: %s", p.Award))
			y += 5
		}
		y += 3
	}
	d.Line(left, y, right, y)
	d.text(left, y+3, right-left, 6, "I", 10, "C", i18n.T(l, "Thank you for taking part!"))
	d.text(left, 202, right-left, 4, "", 7, "C", "www.standblatt.ch")
	return d.bytes()
}

// WreathCard holds the fields printed onto a pre-printed A6 wreath card.
type WreathCard struct {
	Lastname   string
	Firstname  string
	BirthYear  string
	Club       string
	Date       string
	Value      float64
	ShowValue  bool
	ValueWords string
	Organizer  string
}

// RenderWreathCard prints the field values at the positions of the pre-printed card (A6 portrait).
func RenderWreathCard(c WreathCard) ([]byte, error) {
	d := newDoc("P", "A6", 0, 0)
	d.AddPage()
	d.text(30, 94.5, 65, 6, "B", 10, "L", c.Lastname)
	d.text(30, 102.7, 50, 6, "", 10, "L", c.Firstname)
	d.text(85, 102.7, 20, 6, "", 10, "L", c.BirthYear)
	d.text(30, 111, 65, 5, "", 10, "L", c.Club)
	d.text(30, 119.3, 65, 5, "", 10, "L", c.Date)
	d.text(15, 131, 80, 4, "", 8, "R", c.Organizer)
	if c.ShowValue && c.Value > 0 {
		d.text(30, 135.3, 65, 5, "B", 12, "L", fmt.Sprintf("CHF %.2f", c.Value))
		d.text(30, 140.3, 65, 4, "", 8, "L", c.ValueWords)
	}
	return d.bytes()
}

// RankingLine is one line of a printed ranking.
type RankingLine struct {
	Rank          int
	Lastname      string
	Firstname     string
	BirthYear     string
	CategoryShort string
	Club          string
	Series        []float64
	Total         float64
}

// RankingList is the ranking of one programme for printing.
type RankingList struct {
	Event       Event
	Title       string
	SeriesCount int
	Lines       []RankingLine
	Generated   time.Time
}

// RenderRankingList renders a ranking on A4 portrait, repeating the column header on every page.
func RenderRankingList(r RankingList) ([]byte, error) {
	d := newDoc("P", "A4", 0, 0)
	l := r.Event.Lang
	left, width := 10.0, 190.0
	fixed := []float64{10, 12, 10}
	seriesW, totalW := 11.0, 14.0
	rest := width - fixed[0] - fixed[1] - fixed[2] - totalW - seriesW*float64(r.SeriesCount)
	nameW := rest * 0.5
	clubW := rest - nameW
	cols := []struct {
		title string
		w     float64
		align string
	}{{i18n.T(l, "Rank"), fixed[0], "C"}, {i18n.T(l, "Name"), nameW, "L"}, {i18n.T(l, "Year"), fixed[1], "C"}, {i18n.T(l, "Cat."), fixed[2], "C"}, {i18n.T(l, "Club"), clubW, "L"}}
	for i := 0; i < r.SeriesCount; i++ {
		cols = append(cols, struct {
			title string
			w     float64
			align string
		}{i18n.Tf(l, "S%d", i+1), seriesW, "C"})
	}
	cols = append(cols, struct {
		title string
		w     float64
		align string
	}{i18n.T(l, "Total"), totalW, "R"})

	head := func(y float64) float64 {
		d.SetFillColor(229, 231, 235)
		x := left
		for _, c := range cols {
			d.SetFont("Helvetica", "B", 7)
			d.SetXY(x, y)
			d.CellFormat(c.w, 5, d.tr(c.title), "1", 0, c.align, true, 0, "")
			x += c.w
		}
		return y + 5
	}
	podium := [][3]int{{254, 243, 199}, {229, 231, 235}, {254, 215, 170}}

	d.AddPage()
	d.text(left, 10, width, 6, "B", 12, "C", i18n.T(l, "RANKING"))
	d.text(left, 17, width, 5, "B", 10, "C", r.Event.Name)
	d.text(left, 23, width, 5, "B", 9, "C", r.Title)
	d.SetTextColor(80, 80, 80)
	d.text(left, 29, width, 4, "", 8, "C", i18n.Tf(l, "Shooters: %d", len(r.Lines))+" · "+r.Event.Date)
	d.SetTextColor(0, 0, 0)
	y := head(38)
	for _, l := range r.Lines {
		if y > 280 {
			d.AddPage()
			y = head(10)
		}
		fill := l.Rank <= 3
		if fill {
			c := podium[l.Rank-1]
			d.SetFillColor(c[0], c[1], c[2])
		}
		d.SetFont("Helvetica", "", 7)
		values := []string{strconv.Itoa(l.Rank), "", orDash(l.BirthYear), l.CategoryShort, d.fit(l.Club, clubW-2)}
		for i := 0; i < r.SeriesCount; i++ {
			v := "–"
			if i < len(l.Series) {
				v = Points(l.Series[i])
			}
			values = append(values, v)
		}
		values = append(values, Points(l.Total))
		x := left
		for i, c := range cols {
			style := ""
			if i == 0 || i == len(cols)-1 && l.Rank <= 3 {
				style = "B"
			}
			d.SetFont("Helvetica", style, 7)
			d.SetXY(x, y)
			d.CellFormat(c.w, 4.5, d.tr(values[i]), "1", 0, c.align, fill, 0, "")
			if i == 1 {
				d.SetFont("Helvetica", "B", 7)
				d.SetXY(x+0.5, y)
				last := d.tr(l.Lastname)
				d.CellFormat(d.GetStringWidth(last)+0.5, 4.5, last, "", 0, "L", false, 0, "")
				d.SetFont("Helvetica", "", 7)
				d.CellFormat(c.w-d.GetStringWidth(last)-1.5, 4.5, d.tr(" "+l.Firstname), "", 0, "L", false, 0, "")
			}
			x += c.w
		}
		y += 4.5
	}
	d.SetTextColor(100, 100, 100)
	d.text(left, 287, width, 4, "", 7, "C", i18n.Tf(l, "Generated %s", r.Generated.Format("02.01.2006 15:04"))+" · www.standblatt.ch")
	return d.bytes()
}
