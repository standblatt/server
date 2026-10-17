package docs

import (
	"bytes"
	"image"
	"image/png"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/domain"
)

func TestFill(t *testing.T) {
	v := map[string]string{"shooter.lastname": "Example", "s1.shot3": "9"}
	if got := Fill("Name: {{ shooter.lastname }}, shot {{s1.shot#}}, {{unknown}}!", v, 3); got != "Name: Example, shot 9, !" {
		t.Fatalf("got %q", got)
	}
}

func TestRenderLayout(t *testing.T) {
	l := domain.Layout{Width: 148, Height: 210, OffsetX: 1.5, OffsetY: -1, Elements: []domain.Element{
		{Type: "text", X: 20, Y: 20, W: 100, Size: 12, Bold: true, Value: "{{shooter.name}} ({{shooter.birthyear}})"},
		{Type: "text", X: 20, Y: 40, W: 8, Align: "C", Value: "{{s1.shot#}}", Repeat: 10, DX: 9},
		{Type: "text", X: 20, Y: 50, Value: "Total {{total}}"},
		{Type: "qr", X: 120, Y: 180, W: 20, Value: "{{registration}}"},
		{Type: "barcode", X: 20, Y: 180, W: 60, H: 10, Value: "{{shooter.license}}"},
		{Type: "line", X: 20, Y: 60, W: 100},
	}}
	pdf, err := RenderLayout(l, SampleValues(), nil, false)
	if err != nil {
		t.Fatal(err)
	}
	expectText(t, pdfText(t, pdf), "Example Anna (1985)", "Total 85")
	if s := pageSize(t, pdf); s != "" && !strings.HasPrefix(s, "419.5") {
		t.Errorf("page should be 148×210 mm, got %s", s)
	}

	var bg bytes.Buffer
	png.Encode(&bg, image.NewGray(image.Rect(0, 0, 10, 14)))
	preview, err := RenderLayout(l, SampleValues(), bg.Bytes(), true)
	if err != nil {
		t.Fatal(err)
	}
	if len(preview) <= len(pdf) {
		t.Errorf("the preview should carry the background image")
	}
}
