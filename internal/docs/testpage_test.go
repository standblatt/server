package docs

import (
	"strings"
	"testing"

	"github.com/standblatt/server/internal/i18n"
)

func TestRenderTestPage(t *testing.T) {
	pdf, err := RenderTestPage(TestPage{Club: "SG Testdorf", Range: "Pistol 25 m", Printer: "Brother_QL",
		Options: map[string]string{"PageSize": "29x90mm", "fit-to-page": ""}, Time: "18.10.2026 21:30", Lang: i18n.Get("en")})
	if err != nil {
		t.Fatal(err)
	}
	if s := pageSize(t, pdf); s != "" && !strings.HasPrefix(s, "255.1") {
		t.Fatalf("page size %s, want 90 × 29 mm", s)
	}
	expectText(t, pdfText(t, pdf), "Standblatt test page", "SG Testdorf · Pistol 25 m", "Printer: Brother_QL",
		"Options: PageSize=29x90mm, fit-to-page", "18.10.2026 21:30")

	pdf, err = RenderTestPage(TestPage{Printer: "Office", Lang: i18n.Get("en")})
	if err != nil {
		t.Fatal(err)
	}
	expectText(t, pdfText(t, pdf), "Options: default")
}
