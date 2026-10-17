package main

import "testing"

func TestLoadMatomo(t *testing.T) {
	m, err := loadMatomo("https://analytics.example.org", "33")
	if err != nil || m.URL != "https://analytics.example.org/" || m.SiteID != 33 {
		t.Fatalf("got %+v, %v", m, err)
	}
	if m, err := loadMatomo("", ""); err != nil || m.URL != "" {
		t.Fatalf("unset: got %+v, %v", m, err)
	}
	for _, c := range [][2]string{{"https://analytics.example.org/", ""}, {"", "33"}, {"analytics.example.org", "33"}, {"https://analytics.example.org/", "x"}, {"https://analytics.example.org/", "0"}} {
		if _, err := loadMatomo(c[0], c[1]); err == nil {
			t.Errorf("%q, %q: want an error", c[0], c[1])
		}
	}
}
