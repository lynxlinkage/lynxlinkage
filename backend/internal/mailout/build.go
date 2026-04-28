package mailout

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"
	"unicode/utf8"
)

// buildRFC822Message returns either a plain or multipart/alternative body.
// If html is empty, a single text/plain message is built.
func buildRFC822Message(from, to, subject, textPlain, html string) (string, error) {
	subj, err := encodeSubject(subject)
	if err != nil {
		return "", err
	}
	plain := strings.ReplaceAll(strings.ReplaceAll(textPlain, "\r\n", "\n"), "\n", "\r\n")
	if !strings.HasSuffix(plain, "\r\n") {
		plain += "\r\n"
	}
	if strings.TrimSpace(html) == "" {
		return buildPlain(from, to, subj, plain), nil
	}
	return buildMultipart(from, to, subj, plain, html)
}

func buildPlain(from, to, subject, body string) string {
	var b strings.Builder
	b.WriteString("From: ")
	b.WriteString(from)
	b.WriteString("\r\nTo: ")
	b.WriteString(to)
	b.WriteString("\r\nSubject: ")
	b.WriteString(subject)
	b.WriteString("\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n")
	b.WriteString(body)
	return b.String()
}

func buildMultipart(from, to, subject, textPlain, html string) (string, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	h1 := textproto.MIMEHeader{}
	h1.Set("Content-Type", "text/plain; charset=UTF-8")
	h1.Set("Content-Transfer-Encoding", "8bit")
	p1, err := w.CreatePart(h1)
	if err != nil {
		return "", err
	}
	if _, err = p1.Write([]byte(textPlain)); err != nil {
		return "", err
	}
	h2 := textproto.MIMEHeader{}
	h2.Set("Content-Type", "text/html; charset=UTF-8")
	h2.Set("Content-Transfer-Encoding", "8bit")
	p2, err := w.CreatePart(h2)
	if err != nil {
		return "", err
	}
	if _, err = p2.Write([]byte(html)); err != nil {
		return "", err
	}
	if err = w.Close(); err != nil {
		return "", err
	}
	boundary := w.Boundary()
	var out strings.Builder
	out.WriteString("From: ")
	out.WriteString(from)
	out.WriteString("\r\nTo: ")
	out.WriteString(to)
	out.WriteString("\r\nSubject: ")
	out.WriteString(subject)
	out.WriteString("\r\nMIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=")
	out.WriteString(boundary)
	out.WriteString("\r\n\r\n")
	out.Write(body.Bytes())
	return out.String(), nil
}

func encodeSubject(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("mailout: empty subject")
	}
	ascii := true
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			ascii = false
			break
		}
		if r > 127 {
			ascii = false
			break
		}
		i += size
	}
	if ascii {
		if strings.ContainsAny(s, "\r\n") {
			return "", fmt.Errorf("mailout: invalid subject")
		}
		return s, nil
	}
	enc := base64.StdEncoding.EncodeToString([]byte(s))
	// Per RFC 2047, split long encoded words (here keep one chunk; subjects are short)
	return fmt.Sprintf("=?utf-8?B?%s?=", enc), nil
}
