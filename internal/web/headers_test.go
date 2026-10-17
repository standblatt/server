package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestSecureHeaders(t *testing.T) {
	s := &Server{UI: fstest.MapFS{"index.html": {Data: []byte("<html><script>alert(1)</script></html>")}}}
	h := s.secureHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	h.ServeHTTP(rec, req)
	csp := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "'sha256-bhHHL3z2vDgxUt0W3dWQOrprscmda2Y5pLsLg4GF+pI='") {
		t.Fatalf("index.html script hash missing: %s", csp)
	}
	if strings.Count(csp, "sha256-") != 3 {
		t.Fatalf("expected 3 inline script hashes: %s", csp)
	}
	for k, want := range map[string]string{
		"Strict-Transport-Security":    "max-age=31536000; includeSubDomains",
		"X-Content-Type-Options":       "nosniff",
		"X-Frame-Options":              "DENY",
		"Cross-Origin-Opener-Policy":   "same-origin",
		"Cross-Origin-Resource-Policy": "same-origin",
		"Referrer-Policy":              "strict-origin-when-cross-origin",
		"Permissions-Policy":           "camera=(), microphone=(), geolocation=()",
	} {
		if got := rec.Header().Get(k); got != want {
			t.Errorf("%s = %q, want %q", k, got, want)
		}
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Header().Get("Strict-Transport-Security") != "" {
		t.Error("HSTS must only be sent over https")
	}
}

func TestSecureHeadersMatomo(t *testing.T) {
	s := &Server{UI: fstest.MapFS{}, Matomo: Matomo{URL: "https://analytics.example.org/matomo/", SiteID: 3}}
	rec := httptest.NewRecorder()
	s.secureHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	csp := rec.Header().Get("Content-Security-Policy")
	for _, d := range []string{"script-src 'self' ", "connect-src 'self' https://analytics.example.org;", "img-src 'self' data: blob: https://analytics.example.org;"} {
		if !strings.Contains(csp, d) {
			t.Errorf("CSP lacks %q: %s", d, csp)
		}
	}
	if !strings.Contains(csp, "https://analytics.example.org; connect-src") {
		t.Errorf("script-src lacks the Matomo origin: %s", csp)
	}
}
