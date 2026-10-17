package docs

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var person = Person{
	Lastname: "Böckli", Firstname: "Fritz", Birthday: "02.05.1930", BirthYear: "1930", License: "114870",
	Club: "Frauenfeld Stadtschützen Gesellschaft mit einem sehr langen Namen", Category: "Senior veteran", CategoryShort: "SV",
}

func pdfText(t *testing.T, pdf []byte) string {
	t.Helper()
	if dir := os.Getenv("DOCS_OUT"); dir != "" {
		os.WriteFile(dir+"/"+t.Name()+".pdf", pdf, 0o644)
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Fatalf("not a PDF: %q", pdf[:min(20, len(pdf))])
	}
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return ""
	}
	f, _ := os.CreateTemp(t.TempDir(), "*.pdf")
	f.Write(pdf)
	f.Close()
	out, err := exec.Command("pdftotext", "-layout", f.Name(), "-").Output()
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func pageSize(t *testing.T, pdf []byte) string {
	t.Helper()
	if _, err := exec.LookPath("pdfinfo"); err != nil {
		return ""
	}
	f, _ := os.CreateTemp(t.TempDir(), "*.pdf")
	f.Write(pdf)
	f.Close()
	out, _ := exec.Command("pdfinfo", f.Name()).Output()
	for _, l := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(l, "Page size:") {
			return strings.TrimSpace(strings.TrimPrefix(l, "Page size:"))
		}
	}
	return ""
}

func expectText(t *testing.T, text string, want ...string) {
	t.Helper()
	if text == "" {
		return
	}
	for _, w := range want {
		if !strings.Contains(text, w) {
			t.Errorf("missing %q in:\n%s", w, text)
		}
	}
}

func TestCode(t *testing.T) {
	if got := Code(42, ""); got != "REGCODEAAAA42AAAAL0" {
		t.Fatal(got)
	}
	if got := Code(42, "013756"); got != "REGCODEAAAA42AAAAL013756" {
		t.Fatal(got)
	}
}

func TestLaneSheet(t *testing.T) {
	lane := 7
	pdf, err := RenderLaneSheet(LaneSheet{
		Event: Event{Name: "Indoor 2026", Date: "20.11.2026", Organizer: "SV Muster"}, Person: person, RegistrationID: 42, Lane: &lane,
		Programmes: []LaneSheetProgramme{
			{Title: "Field match (Stgw 90)", Shots: 18, Barcode: "706625002584250(75)"},
			{Title: "Section", Shots: 0, CombinedWith: "Main"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectText(t, pdfText(t, pdf), "LANE SHEET", "Indoor 2026", "Böckli Fritz, 02.05.1930", "LANE 7", "Field match (Stgw 90) – 18 shots", "070662500258425075", "Combined with Main")
	if s := pageSize(t, pdf); s != "" && !strings.HasSuffix(s, "(A5)") {
		t.Errorf("lane sheet should be A5, got %s", s)
	}
}

func TestControlStrip(t *testing.T) {
	pdf, err := RenderControlStrip(ControlStrip{EventName: "Indoor 2026", Person: person, RegistrationID: 42, Passes: "P1 P2 + P0"})
	if err != nil {
		t.Fatal(err)
	}
	expectText(t, pdfText(t, pdf), "Böckli Fritz", "SV", "P1 P2 + P0", "ID 42")
	if s := pageSize(t, pdf); s != "" && !strings.HasPrefix(s, "255.9") && !strings.HasPrefix(s, "256") {
		t.Errorf("label should be 90.3 × 29 mm landscape, got %s", s)
	}
}

func TestResultSheet(t *testing.T) {
	pdf, err := RenderResultSheet(ResultSheet{
		Event: Event{Name: "Indoor 2026", Date: "20.11.2026"}, Person: person,
		Programmes: []ResultSheetProgramme{{Title: "Main", Series: []float64{95, 88.5, 100, 97}, Total: 380.5, MaxPoints: 400, Award: "Wreath"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectText(t, pdfText(t, pdf), "RESULT SHEET", "Böckli Fritz, 1930", "88.5", "380.5 / 400", "Award: Wreath")
}

func TestWreathCard(t *testing.T) {
	pdf, err := RenderWreathCard(WreathCard{Lastname: "Böckli", Firstname: "Fritz", BirthYear: "1930", Club: "SV", Date: "20.11.2026",
		Value: 25, ShowValue: true, ValueWords: "Twenty-five francs", Organizer: "SV Muster"})
	if err != nil {
		t.Fatal(err)
	}
	expectText(t, pdfText(t, pdf), "Böckli", "CHF 25.00", "Twenty-five francs")
	if s := pageSize(t, pdf); s != "" && !strings.HasSuffix(s, "(A6)") {
		t.Errorf("wreath card should be A6, got %s", s)
	}
}

func TestRankingList(t *testing.T) {
	var lines []RankingLine
	for i := 1; i <= 120; i++ {
		lines = append(lines, RankingLine{Rank: i, Lastname: "Shooter", Firstname: "Number", BirthYear: "1980", CategoryShort: "E", Club: person.Club, Series: []float64{95, 90}, Total: 185})
	}
	pdf, err := RenderRankingList(RankingList{Event: Event{Name: "Indoor 2026", Date: "20.11.2026"}, Title: "Main", SeriesCount: 2, Lines: lines, Generated: time.Date(2026, 11, 20, 18, 5, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	text := pdfText(t, pdf)
	expectText(t, text, "RANKING", "Main", "Shooters: 120", "Generated 20.11.2026 18:05")
	if text != "" && strings.Count(text, "Rank") < 2 {
		t.Errorf("column header should repeat on every page")
	}
}

func TestReport(t *testing.T) {
	var rows [][]string
	for i := 0; i < 70; i++ {
		rows = append(rows, []string{"Field match", "Stgw 90", "12", "CHF 240.00"})
	}
	r := Report{
		Title: "Financial report", Subtitle: "Indoor 2026", Summary: [][2]string{{"Revenue", "CHF 1'234.50"}},
		Sections: []Section{{
			Heading: "Revenue by programme",
			Columns: []Column{{"Programme", 4, "L"}, {"Variant", 2, "L"}, {"Participants", 1, "C"}, {"Revenue", 1.5, "R"}},
			Rows:    rows, Total: []string{"Total", "", "840", "CHF 16'800.00"},
		}},
		Generated: time.Date(2026, 11, 20, 18, 5, 0, 0, time.UTC),
	}
	pdf, err := RenderReport(r)
	if err != nil {
		t.Fatal(err)
	}
	text := pdfText(t, pdf)
	expectText(t, text, "Financial report", "Revenue by programme", "CHF 16'800.00", "page 2")
	if text != "" && strings.Count(text, "Participants") < 2 {
		t.Errorf("table header should repeat on the second page")
	}

	csvOut, err := ReportCSV(r)
	if err != nil {
		t.Fatal(err)
	}
	s := string(csvOut)
	if !strings.HasPrefix(s, "\xef\xbb\xbfFinancial report\n") || !strings.Contains(s, "\nProgramme;Variant;Participants;Revenue\n") ||
		!strings.Contains(s, "\nTotal;;840;CHF 16'800.00\n") || !strings.Contains(s, "Revenue;CHF 1'234.50") {
		t.Fatalf("csv:\n%s", s)
	}
}
