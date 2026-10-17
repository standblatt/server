package docs

import (
	"bytes"
	"encoding/csv"
	"time"

	"github.com/standblatt/server/internal/i18n"
)

// Column is a table column; Width is relative to the other columns.
type Column struct {
	Title string
	Width float64
	Align string
}

// Section is one table of a report, optionally with a bold total row.
type Section struct {
	Heading string
	Note    string
	Columns []Column
	Rows    [][]string
	Total   []string
}

// Report is a tabular report that renders as PDF or CSV; Lang translates the texts the renderer adds.
type Report struct {
	Lang      i18n.Lang
	Title     string
	Subtitle  string
	Landscape bool
	Summary   [][2]string
	Sections  []Section
	Generated time.Time
}

// RenderReport renders a report on A4; tables continue on new pages with their header repeated.
func RenderReport(r Report) ([]byte, error) {
	orientation, width, bottom := "P", 190.0, 282.0
	if r.Landscape {
		orientation, width, bottom = "L", 277.0, 195.0
	}
	d := newDoc(orientation, "A4", 0, 0)
	left := 10.0
	d.SetFooterFunc(func() {
		d.SetTextColor(100, 100, 100)
		d.text(left, bottom+6, width, 4, "", 7, "C",
			i18n.Tf(r.Lang, "Generated %s · Standblatt · page %d", r.Generated.Format("02.01.2006 15:04"), d.PageNo()))
		d.SetTextColor(0, 0, 0)
	})
	d.AddPage()
	d.text(left, 10, width, 8, "B", 15, "C", r.Title)
	y := 19.0
	if r.Subtitle != "" {
		d.text(left, y, width, 5, "", 10, "C", r.Subtitle)
		y += 7
	}
	y += 3
	for _, s := range r.Summary {
		d.text(left, y, 70, 5, "", 9, "L", s[0])
		d.text(left+70, y, 60, 5, "B", 9, "L", s[1])
		y += 5
	}
	if len(r.Summary) > 0 {
		y += 4
	}
	for _, sec := range r.Sections {
		total := 0.0
		for _, c := range sec.Columns {
			total += c.Width
		}
		widths := make([]float64, len(sec.Columns))
		for i, c := range sec.Columns {
			widths[i] = c.Width / total * width
		}
		row := func(y float64, cells []string, style string, fill bool) {
			d.SetFont("Helvetica", style, 7.5)
			x := left
			for i, c := range sec.Columns {
				v := ""
				if i < len(cells) {
					v = cells[i]
				}
				d.SetXY(x, y)
				d.CellFormat(widths[i], 5, d.tr(d.fit(v, widths[i]-1.5)), "1", 0, c.Align, fill, 0, "")
				x += widths[i]
			}
		}
		head := func(y float64) float64 {
			titles := make([]string, len(sec.Columns))
			for i, c := range sec.Columns {
				titles[i] = c.Title
			}
			d.SetFillColor(229, 231, 235)
			row(y, titles, "B", true)
			return y + 5
		}
		if y > bottom-25 {
			d.AddPage()
			y = 10
		}
		if sec.Heading != "" {
			d.text(left, y, width, 6, "B", 11, "L", sec.Heading)
			y += 7
		}
		if sec.Note != "" {
			d.text(left, y, width, 4, "I", 8, "L", sec.Note)
			y += 5
		}
		y = head(y)
		for _, cells := range sec.Rows {
			if y > bottom-5 {
				d.AddPage()
				y = head(10)
			}
			row(y, cells, "", false)
			y += 5
		}
		if len(sec.Rows) == 0 {
			d.text(left, y, width, 5, "I", 8, "L", i18n.T(r.Lang, "No data."))
			y += 5
		}
		if sec.Total != nil {
			d.SetFillColor(243, 244, 246)
			row(y, sec.Total, "B", true)
			y += 5
		}
		y += 6
	}
	return d.bytes()
}

// ReportCSV writes a report as CSV for spreadsheets: UTF-8 with BOM, semicolon separated, sections divided by
// an empty line.
func ReportCSV(r Report) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(&buf)
	w.Comma = ';'
	w.Write([]string{r.Title})
	if r.Subtitle != "" {
		w.Write([]string{r.Subtitle})
	}
	for _, s := range r.Summary {
		w.Write([]string{s[0], s[1]})
	}
	for _, sec := range r.Sections {
		w.Write(nil)
		if sec.Heading != "" {
			w.Write([]string{sec.Heading})
		}
		titles := make([]string, len(sec.Columns))
		for i, c := range sec.Columns {
			titles[i] = c.Title
		}
		w.Write(titles)
		for _, row := range sec.Rows {
			w.Write(row)
		}
		if sec.Total != nil {
			w.Write(sec.Total)
		}
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}
