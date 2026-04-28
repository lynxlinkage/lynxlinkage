package mailout

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// smtpAuth returns PLAIN auth. Many servers (including Hostinger) expect the
// authorization identity to match the mailbox username; empty identity also
// works on most hosts, but we set it explicitly.
func (c *Config) smtpAuth() smtp.Auth {
	return smtp.PlainAuth(c.User, c.User, c.Pass, c.Host)
}

func (c *Config) sendSMTP(to, subject, textPlain, html string) error {
	to = strings.TrimSpace(to)
	if to == "" {
		return fmt.Errorf("mailout: empty recipient")
	}
	subject = stripHeaderField(subject)
	msg, err := buildRFC822Message(c.From, to, subject, textPlain, html)
	if err != nil {
		return err
	}
	if c.Port == 465 {
		return c.sendMailImplicitTLS(c.smtpAuth(), to, []byte(msg))
	}
	return c.sendMailStartTLS(c.smtpAuth(), to, []byte(msg))
}

// sendMailStartTLS dials the SMTP server, negotiates STARTTLS, then AUTH.
func (c *Config) sendMailStartTLS(a smtp.Auth, to string, msg []byte) error {
	if err := validateLine(c.From); err != nil {
		return err
	}
	if err := validateLine(to); err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()
	if err = client.Hello("localhost"); err != nil {
		return err
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: c.Host,
			MinVersion: tls.VersionTLS12,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return err
		}
	} else if c.Port == 587 || c.Port == 25 || c.Port == 2525 {
		return fmt.Errorf("mailout: server did not offer STARTTLS on port %d", c.Port)
	}
	if a != nil {
		if err = client.Auth(a); err != nil {
			return err
		}
	}
	if err = client.Mail(extractAddr(c.From)); err != nil {
		return err
	}
	if err = client.Rcpt(extractAddr(to)); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	return w.Close()
}

// sendMailImplicitTLS is for port 465 (SSL/TLS from the first byte).
func (c *Config) sendMailImplicitTLS(a smtp.Auth, to string, msg []byte) error {
	if err := validateLine(c.From); err != nil {
		return err
	}
	if err := validateLine(to); err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	d := tls.Dialer{NetDialer: nil, Config: &tls.Config{
		ServerName: c.Host,
		MinVersion: tls.VersionTLS12,
	}}
	conn, err := d.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, c.Host)
	if err != nil {
		return err
	}
	defer client.Close()
	if a != nil {
		if err = client.Auth(a); err != nil {
			return err
		}
	}
	if err = client.Mail(extractAddr(c.From)); err != nil {
		return err
	}
	if err = client.Rcpt(extractAddr(to)); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return client.Quit()
}

func validateLine(s string) error {
	if strings.ContainsAny(s, "\r\n") {
		return fmt.Errorf("mailout: invalid line break in address")
	}
	return nil
}

// extractAddr takes "Name <a@b.com>" or "a@b.com" and returns the email part.
func extractAddr(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "<"); i >= 0 {
		if j := strings.Index(s[i:], ">"); j > 0 {
			return strings.TrimSpace(s[i+1 : i+j])
		}
	}
	return s
}
