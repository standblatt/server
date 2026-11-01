package mail

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"testing"
	"time"
)

// TestFormat checks the plain text, the escaped HTML and that both parts of the multipart mail parse back.
func TestFormat(t *testing.T) {
	m := Message{To: "ida@example.org", Subject: "Bahnen für <SV>", Parts: []Part{
		P("Hallo Ida,"),
		P("Hier:"),
		Button("Anfrage öffnen", "https://example.org/b/manage/bk_1?a=1&b=2"),
		Table([]string{"Time", "Lanes"}, []string{"08:00–08:30", "3 (1-3)"}),
		P("<script>"),
	}}
	text := "Hallo Ida,\n\nHier:\n\nhttps://example.org/b/manage/bk_1?a=1&b=2\n\n- Time: 08:00–08:30, Lanes: 3 (1-3)\n\n<script>"
	if got := m.Text(); got != text {
		t.Fatalf("text:\n%s", got)
	}
	html := m.HTML()
	for _, want := range []string{`href="https://example.org/b/manage/bk_1?a=1&amp;b=2"`, ">Anfrage öffnen</a>", "<th", "3 (1-3)</td>", "&lt;script&gt;"} {
		if !strings.Contains(html, want) {
			t.Fatalf("html lacks %q:\n%s", want, html)
		}
	}

	msg, err := mail.ReadMessage(bytes.NewReader(Format("Standblatt <noreply@example.org>", m, time.Now())))
	if err != nil {
		t.Fatal(err)
	}
	kind, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil || kind != "multipart/alternative" {
		t.Fatalf("content type %q: %v", kind, err)
	}
	r := multipart.NewReader(msg.Body, params["boundary"])
	for _, want := range []struct{ kind, body string }{{"text/plain; charset=utf-8", text}, {"text/html; charset=utf-8", html}} {
		p, err := r.NextPart()
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(p)
		if p.Header.Get("Content-Type") != want.kind || strings.ReplaceAll(string(b), "\r\n", "\n") != want.body {
			t.Fatalf("part %s:\n%s", p.Header.Get("Content-Type"), b)
		}
	}
}
