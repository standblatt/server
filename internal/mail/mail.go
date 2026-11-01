// Package mail sends emails with a plain-text and an HTML part over SMTP, or only logs them when no server is
// configured.
package mail

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

// Sender delivers mails through the SMTP server at Host:Port, with implicit TLS on port 465 and STARTTLS
// elsewhere when the server offers it; User and Password are optional. With an empty Host it logs the mail
// instead, which is enough for development.
type Sender struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	Log      *slog.Logger
}

// Message is one mail; its parts render to the plain text and to the HTML alternative.
type Message struct {
	To      string
	Subject string
	Parts   []Part
}

// Part is a block of a mail: a paragraph, a button when Link is set (Text is its label) or a table when Table
// is set (the first row names the columns).
type Part struct {
	Text  string
	Link  string
	Table [][]string
}

// P is a paragraph.
func P(text string) Part { return Part{Text: text} }

// Button is a link shown as a button; the plain text shows only the URL, so the paragraph before introduces it.
func Button(label, link string) Part { return Part{Text: label, Link: link} }

// Table is a table with column names and rows.
func Table(columns []string, rows ...[]string) Part {
	return Part{Table: append([][]string{columns}, rows...)}
}

// Text renders the message as plain text, one block per paragraph and a "- column: value" line per table row.
func (m Message) Text() string {
	blocks := make([]string, 0, len(m.Parts))
	for _, p := range m.Parts {
		switch {
		case p.Link != "":
			blocks = append(blocks, p.Link)
		case len(p.Table) > 0:
			lines := []string{}
			for _, row := range p.Table[1:] {
				cells := make([]string, len(row))
				for i, c := range row {
					cells[i] = p.Table[0][i] + ": " + c
				}
				lines = append(lines, "- "+strings.Join(cells, ", "))
			}
			blocks = append(blocks, strings.Join(lines, "\n"))
		default:
			blocks = append(blocks, p.Text)
		}
	}
	return strings.Join(blocks, "\n\n")
}

//go:embed layout.html
var layoutHTML string

var layout = template.Must(template.New("mail").Parse(layoutHTML))

// HTML renders the message into the shared mail layout.
func (m Message) HTML() string {
	var b strings.Builder
	if err := layout.Execute(&b, m); err != nil {
		return ""
	}
	return b.String()
}

// Send delivers a message.
func (s *Sender) Send(ctx context.Context, m Message) error {
	if s.Host == "" {
		s.Log.Info("mail not sent, no SMTP server configured", "to", m.To, "subject", m.Subject, "body", m.Text())
		return nil
	}
	host, addr, implicitTLS := s.Host, net.JoinHostPort(s.Host, strconv.Itoa(s.Port)), s.Port == 465
	d := net.Dialer{Timeout: 10 * time.Second}
	var conn net.Conn
	var err error
	if implicitTLS {
		conn, err = tls.DialWithDialer(&d, "tcp", addr, &tls.Config{ServerName: host})
	} else {
		conn, err = d.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("mail: %w", err)
	}
	conn.SetDeadline(time.Now().Add(30 * time.Second))
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("mail: %w", err)
	}
	defer c.Close()
	if ok, _ := c.Extension("STARTTLS"); ok && !implicitTLS {
		if err := c.StartTLS(&tls.Config{ServerName: host}); err != nil {
			return fmt.Errorf("mail: %w", err)
		}
	}
	if s.User != "" {
		if err := c.Auth(smtp.PlainAuth("", s.User, s.Password, host)); err != nil {
			return fmt.Errorf("mail: %w", err)
		}
	}
	if err := c.Mail(address(s.From)); err != nil {
		return fmt.Errorf("mail: %w", err)
	}
	if err := c.Rcpt(m.To); err != nil {
		return fmt.Errorf("mail: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("mail: %w", err)
	}
	if _, err := w.Write(Format(s.From, m, time.Now())); err != nil {
		return fmt.Errorf("mail: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mail: %w", err)
	}
	return c.Quit()
}

// address returns the bare address of a From value such as "Name <a@b>".
func address(from string) string {
	if i := strings.LastIndex(from, "<"); i >= 0 {
		return strings.TrimSuffix(from[i+1:], ">")
	}
	return from
}

// Format renders a message with headers as multipart/alternative with a plain-text and an HTML part.
func Format(from string, m Message, now time.Time) []byte {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", m.To)
	fmt.Fprintf(&b, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", m.Subject))
	fmt.Fprintf(&b, "Date: %s\r\n", now.Format(time.RFC1123Z))
	fmt.Fprintf(&b, "MIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=%q\r\n\r\n", mw.Boundary())
	for _, part := range []struct{ kind, body string }{{"text/plain", m.Text()}, {"text/html", m.HTML()}} {
		w, _ := mw.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {part.kind + "; charset=utf-8"},
			"Content-Transfer-Encoding": {"quoted-printable"},
		})
		q := quotedprintable.NewWriter(w)
		q.Write([]byte(strings.ReplaceAll(part.body, "\r\n", "\n")))
		q.Close()
	}
	mw.Close()
	return b.Bytes()
}
