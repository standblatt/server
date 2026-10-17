// Package mail sends plain-text emails over SMTP, or only logs them when no server is configured.
package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"mime"
	"net"
	"net/smtp"
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

// Message is one plain-text mail.
type Message struct {
	To      string
	Subject string
	Body    string
}

// Send delivers a message.
func (s *Sender) Send(ctx context.Context, m Message) error {
	if s.Host == "" {
		s.Log.Info("mail not sent, no SMTP server configured", "to", m.To, "subject", m.Subject, "body", m.Body)
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

// Format renders a message with headers as UTF-8 text.
func Format(from string, m Message, now time.Time) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", m.To)
	fmt.Fprintf(&b, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", m.Subject))
	fmt.Fprintf(&b, "Date: %s\r\n", now.Format(time.RFC1123Z))
	b.WriteString("MIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n")
	b.WriteString(strings.ReplaceAll(strings.ReplaceAll(m.Body, "\r\n", "\n"), "\n", "\r\n"))
	return []byte(b.String())
}
