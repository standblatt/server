package docs

import (
	"bytes"
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

// RenderLayout prints the values onto a page of the layout's size. A preview draws background (a PNG or JPEG of
// the empty form) under the elements and leaves out the printer offset; a print never shows the background.
func RenderLayout(l domain.Layout, values map[string]string, background []byte, preview bool) ([]byte, error) {
	d := newDoc("P", "", l.Width, l.Height)
	if l.Width > l.Height {
		d = newDoc("L", "", l.Height, l.Width)
	}
	d.AddPage()
	ox, oy := l.OffsetX, l.OffsetY
	if preview {
		ox, oy = 0, 0
		if len(background) > 0 {
			kind := "PNG"
			if http.DetectContentType(background) == "image/jpeg" {
				kind = "JPG"
			}
			d.RegisterImageOptionsReader("background", fpdf.ImageOptions{ImageType: kind}, bytes.NewReader(background))
			d.ImageOptions("background", 0, 0, l.Width, l.Height, false, fpdf.ImageOptions{ImageType: kind}, 0, "")
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
				h := e.H
				if h == 0 {
					h = size * 0.45
				}
				w := e.W
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
			}
		}
	}
	return d.bytes()
}

// SampleValues are made-up values for every placeholder, used for previews without a real result.
func SampleValues() map[string]string {
	v := map[string]string{
		"shooter.firstname": "Anna", "shooter.lastname": "Example", "shooter.name": "Example Anna",
		"shooter.license": "123456", "shooter.birthday": "02.06.1985", "shooter.birthyear": "1985",
		"shooter.postalcode": "9050", "shooter.city": "Appenzell",
		"club": "SG Hilltown", "club.ssv": "1.23.4.56",
		"event.name": "Annual programme 2027", "event.date": "12.05.2027", "event.organizer": "SG Hilltown",
		"programme": "Obligatory", "category": "Elite", "lane": "7", "registration": "1042", "today": "12.05.2027",
		"total": "85", "mouches": "3", "award": "Card",
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
