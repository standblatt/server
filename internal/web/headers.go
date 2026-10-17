package web

import (
	"crypto/sha256"
	"encoding/base64"
	"io/fs"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var inlineScript = regexp.MustCompile(`(?s)<script>(.*?)</script>`)

// scriptHashes returns the CSP hashes of the inline scripts in the given HTML documents.
func scriptHashes(docs ...[]byte) []string {
	var out []string
	for _, doc := range docs {
		for _, m := range inlineScript.FindAllSubmatch(doc, -1) {
			sum := sha256.Sum256(m[1])
			out = append(out, "'sha256-"+base64.StdEncoding.EncodeToString(sum[:])+"'")
		}
	}
	return out
}

// framable reports whether a path may be shown in a frame of the app itself: the preview PDF of print forms,
// which the form editor shows next to the elements.
func framable(path string) bool {
	return strings.Contains(path, "/layouts/") && strings.HasSuffix(path, "/preview")
}

// secureHeaders adds the browser security headers to every response. Inline scripts are
// allowed by hash only: the SvelteKit start script and the admin panel's scripts. No page may be framed,
// except framable ones by the app itself.
func (s *Server) secureHeaders(next http.Handler) http.Handler {
	index, _ := fs.ReadFile(s.UI, "index.html")
	head, _ := templateFS.ReadFile("templates/head.html")
	analytics := ""
	if u, err := url.Parse(s.Matomo.URL); err == nil && u.Host != "" {
		analytics = " " + u.Scheme + "://" + u.Host
	}
	csp := strings.Join([]string{
		"default-src 'self'",
		"script-src 'self' " + strings.Join(scriptHashes(index, head), " ") + analytics,
		"connect-src 'self'" + analytics,
		"style-src 'self' 'unsafe-inline'",
		"img-src 'self' data: blob:" + analytics,
		"object-src 'none'",
		"base-uri 'self'",
		"form-action 'self'",
		"frame-ancestors 'none'",
	}, "; ")
	sameOrigin := strings.Replace(csp, "frame-ancestors 'none'", "frame-ancestors 'self'", 1)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		if framable(r.URL.Path) {
			h.Set("Content-Security-Policy", sameOrigin)
			h.Set("X-Frame-Options", "SAMEORIGIN")
		} else {
			h.Set("Content-Security-Policy", csp)
			h.Set("X-Frame-Options", "DENY")
		}
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		h.Set("Cross-Origin-Resource-Policy", "same-origin")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}
