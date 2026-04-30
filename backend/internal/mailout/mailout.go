// Package mailout sends simple transactional email over SMTP (e.g. Hostinger).
package mailout

import (
	"errors"
	"fmt"
	"strings"
)

// Config holds SMTP credentials; zero value is not ready to send.
type Config struct {
	From string
	Host string
	Port int
	User string
	Pass string
}

// Ready reports whether outbound mail can be attempted.
func (c *Config) Ready() bool {
	return c != nil && c.Host != "" && c.From != "" && c.User != "" && c.Pass != "" && c.Port > 0
}

var ErrNotConfigured = errors.New("mailout: smtp not fully configured")

// SendText sends a plain UTF-8 message. "To" and From in Config should be
// valid single addresses (optionally in Name <addr> form in From).
func (c *Config) SendText(to, subject, body string) error {
	if !c.Ready() {
		return ErrNotConfigured
	}
	return c.sendSMTP(to, subject, body, "")
}

// SendAlternative sends multipart/alternative text + HTML (same subject).
func (c *Config) SendAlternative(to, subject, textPlain, html string) error {
	if !c.Ready() {
		return ErrNotConfigured
	}
	return c.sendSMTP(to, subject, textPlain, html)
}

func stripHeaderField(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

// ApplicationAckBody returns the acknowledgment copy for a job application.
func ApplicationAckBody(candidateName, jobTitle, brand string) string {
	name := strings.TrimSpace(candidateName)
	if name == "" {
		name = "there"
	}
	title := strings.TrimSpace(jobTitle)
	if title == "" {
		title = "the position"
	}
	brand = strings.TrimSpace(brand)
	if brand == "" {
		brand = "LynxLinkage"
	}
	return fmt.Sprintf(`Hi %s,

Thank you for applying for the "%s" role. We have received your application and will review it carefully.

We aim to reply within 7 days.

If you have questions, please reach out to hr@lynxlinkage.com.

Best regards,
%s
`, name, title, brand)
}

// ApplicationAckSubject returns a short subject line for the ack email.
func ApplicationAckSubject(jobTitle string) string {
	t := strings.TrimSpace(jobTitle)
	if t == "" {
		return "We received your application"
	}
	return fmt.Sprintf("We received your application — %s", stripHeaderField(t))
}
