package docs

import (
	"bytes"
	"cmp"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-pdf/fpdf"
	"github.com/go-pdf/fpdf/contrib/barcode"

	"github.com/standblatt/server/internal/domain"
)

var placeholder = regexp.MustCompile(`\{\{\s*([A-Za-z0-9_.#]+)\s*\}\}`)

// Fill replaces the placeholders of s with values; # inside a placeholder becomes n. Unknown placeholders vanish.
func Fill(s string, values map[string]string, n int) string {
	return placeholder.ReplaceAllStringFunc(s, func(m string) string {
		key := strings.ReplaceAll(placeholder.FindStringSubmatch(m)[1], "#", strconv.Itoa(n))
		return values[key]
	})
}

// RenderLayout prints every set of values onto a page of the layout's size. A preview draws background (a PNG or
// JPEG of the empty form) under the elements and leaves out the printer offset; a print never shows the background.
// Images holds the PNG or JPEG of every image element by ID; an image element without one is left out.
func RenderLayout(l domain.Layout, pages []map[string]string, background []byte, images map[int64][]byte, preview bool) ([]byte, error) {
	d := newDoc("P", "", l.Width, l.Height)
	if l.Width > l.Height {
		d = newDoc("L", "", l.Height, l.Width)
	}
	for _, values := range pages {
		d.AddPage()
		d.page(l, values, background, images, preview)
	}
	return d.bytes()
}

// page draws the elements of a layout with one set of values onto the current page.
func (d *doc) page(l domain.Layout, values map[string]string, background []byte, images map[int64][]byte, preview bool) {
	ox, oy := l.OffsetX, l.OffsetY
	if preview {
		ox, oy = 0, 0
		if len(background) > 0 {
			d.image("background", background, 0, 0, l.Width, l.Height)
		}
	}
	d.SetLineWidth(0.2)
	for _, e := range l.Elements {
		for k := range max(e.Repeat, 1) {
			x, y := e.X+float64(k)*e.DX+ox, e.Y+float64(k)*e.DY+oy
			v := Fill(e.Value, values, k+1)
			size := e.Size
			if size == 0 {
				size = 10
			}
			switch e.Type {
			case "text":
				style := ""
				if e.Bold {
					style = "B"
				}
				if e.Italic {
					style += "I"
				}
				h := e.H
				if h == 0 {
					h = size * 0.45
				}
				w := e.W
				if e.Wrap && w > 0 {
					d.SetFont("Helvetica", style, size)
					d.SetXY(x, y)
					d.MultiCell(w, h, d.tr(v), "", cmp.Or(e.Align, "L"), false)
					continue
				}
				if w > 0 {
					d.SetFont("Helvetica", style, size)
					v = d.fit(v, w)
				}
				d.text(x, y, w, h, style, size, e.Align, v)
			case "qr":
				if v != "" && e.W > 0 {
					d.qr(v, x, y, e.W)
				}
			case "barcode":
				if v != "" && e.W > 0 && e.H > 0 {
					barcode.Barcode(d.Fpdf, barcode.RegisterCode128(d.Fpdf, v), x, y, e.W, e.H, false)
				}
			case "line":
				d.Line(x, y, x+e.W, y+e.H)
			case "image":
				if img := images[e.Image]; len(img) > 0 && e.W > 0 {
					d.image("image"+strconv.FormatInt(e.Image, 10), img, x, y, e.W, e.H)
				}
			}
		}
	}
}

// image draws a PNG or JPEG in w×h (h 0 keeps the aspect ratio), registering it once under name.
func (d *doc) image(name string, data []byte, x, y, w, h float64) {
	opt := fpdf.ImageOptions{ImageType: "PNG"}
	if http.DetectContentType(data) == "image/jpeg" {
		opt.ImageType = "JPG"
	}
	if d.GetImageInfo(name) == nil {
		d.RegisterImageOptionsReader(name, opt, bytes.NewReader(data))
	}
	d.ImageOptions(name, x, y, w, h, false, opt, 0, "")
}

// SampleValues are made-up values for every placeholder of a data source, used for previews without real data.
func SampleValues(source string) map[string]string {
	v := map[string]string{
		"shooter.firstname": "Anna", "shooter.lastname": "Example", "shooter.name": "Example Anna",
		"shooter.license": "123456", "shooter.birthday": "02.06.1985", "shooter.birthyear": "1985",
		"shooter.postalcode": "9999", "shooter.city": "Hilltown",
		"club": "SG Hilltown", "club.ssv": "1.23.4.56",
		"event.name": "Annual programme 2027", "event.date": "12.05.2027", "event.organizer": "SG Hilltown",
		"category": "Elite", "today": "12.05.2027",
	}
	switch source {
	case domain.SourceRegistration:
		for k, s := range map[string]string{"category.short": "E", "lane": "7", "registration": "1042", "registration.code": Code(1042, "123456"),
			"programmes": "Obligatory, Field match", "passes": "P1 P2 + P0", "price": "CHF 25.00"} {
			v[k] = s
		}
		return v
	case domain.SourceRank:
		for k, s := range map[string]string{"rank": "1", "total": "98", "ranking": "Section championship", "class": "Juniors", "programme": "Obligatory"} {
			v[k] = s
		}
		return v
	}
	for k, s := range map[string]string{"programme": "Obligatory", "lane": "7", "registration": "1042", "total": "85", "mouches": "3", "award": "Card",
		"award.value": "CHF 12.00", "award.value.words": "twelve francs"} {
		v[k] = s
	}
	shots := [][]float64{{5, 4, 5, 5, 4, 5, 5, 4, 5, 5}, {9, 10, 8, 9, 10, 9}}
	for i, ss := range shots {
		n := strconv.Itoa(i + 1)
		sum := 0.0
		for j, s := range ss {
			v["s"+n+".shot"+strconv.Itoa(j+1)] = Points(s)
			sum += s
		}
		v["s"+n] = Points(sum)
		v["s"+n+".mouches"] = "1"
	}
	return v
}
