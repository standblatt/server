package docs

import (
	"sort"
	"strings"

	"github.com/standblatt/server/internal/i18n"
)

// TestPage is the test page of the printing setup: printer, options and where it was sent from.
type TestPage struct {
	Club, Range, Printer string
	Options              map[string]string
	Time                 string
	Lang                 i18n.Lang
}

// RenderTestPage renders the test page as a 90 × 29 mm label in landscape, the size of the control strip: a label
// printer prints it as it is, a sheet printer in the top-left corner (or scaled up when the options say so).
func RenderTestPage(p TestPage) ([]byte, error) {
	d := newDoc("L", "", 29, 90)
	d.AddPage()
	left, w := 4.0, 82.0
	d.Rect(1, 1, 88, 27, "D")
	d.text(left, 2.5, w, 5, "B", 10, "L", i18n.T(p.Lang, "Standblatt test page"))
	d.text(left, 8, w, 4, "", 8, "L", d.fit(p.Club+" · "+p.Range, w))
	d.text(left, 12.5, w, 4, "", 8, "L", d.fit(i18n.Tf(p.Lang, "Printer: %s", p.Printer), w))
	keys := make([]string, 0, len(p.Options))
	for k := range p.Options {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	opts := make([]string, 0, len(keys))
	for _, k := range keys {
		if v := p.Options[k]; v != "" {
			opts = append(opts, k+"="+v)
		} else {
			opts = append(opts, k)
		}
	}
	if len(opts) == 0 {
		opts = append(opts, i18n.T(p.Lang, "default"))
	}
	d.text(left, 17, w, 4, "", 8, "L", d.fit(i18n.Tf(p.Lang, "Options: %s", strings.Join(opts, ", ")), w))
	d.text(left, 22, w, 4, "", 7, "L", p.Time)
	return d.bytes()
}
