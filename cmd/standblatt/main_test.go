package main

import "testing"

func TestLoadConfigPublicURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
	t.Setenv("STANDBLATT_SMTP_HOST", "mail.example.org")
	t.Setenv("STANDBLATT_PUBLIC_URL", "")
	if _, err := loadConfig(); err == nil {
		t.Fatal("SMTP without a public URL: want an error")
	}
	t.Setenv("STANDBLATT_PUBLIC_URL", "https://app.example.org")
	if _, err := loadConfig(); err != nil {
		t.Fatal(err)
	}
	t.Setenv("STANDBLATT_SMTP_HOST", "")
	t.Setenv("STANDBLATT_PUBLIC_URL", "")
	if _, err := loadConfig(); err != nil {
		t.Fatalf("mails only logged: %v", err)
	}
}

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
