// Application mailout ack HTML body.
package mailout

import (
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"
)

// ApplicationAckHTML returns a self-contained HTML body for the application
// receipt. siteURL should be the public site base (https, no trailing slash)
// so the footer link resolves.
func ApplicationAckHTML(candidateName, jobTitle, brand, siteURL string) string {
	name := html.EscapeString(dispName(candidateName))
	title := html.EscapeString(dispTitle(jobTitle))
	brand = html.EscapeString(strings.TrimSpace(brand))
	if brand == "" {
		brand = "LynxLinkage"
	}
	site := strings.TrimSuffix(strings.TrimSpace(siteURL), "/")
	if site == "" {
		site = "https://lynxlinkage.com"
	}
	u, err := url.Parse(site)
	if err != nil || u.Scheme == "" {
		site = "https://lynxlinkage.com"
		u, _ = url.Parse(site)
	}
	siteText := u.Host
	if siteText == "" {
		siteText = "lynxlinkage.com"
	}
	year := time.Now().UTC().Year()
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width"><meta http-equiv="X-UA-Compatible" content="IE=edge"></head>
<body style="margin:0;padding:0;background:#eef1f6;font-family:Georgia,'Times New Roman',Times,serif;">
<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background:#eef1f6;">
  <tr>
    <td align="center" style="padding:48px 32px;">
      <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="width:100%%;max-width:680px;mso-table-lspace:0pt;mso-table-rspace:0pt;">
        <tr>
          <td style="padding:0 0 26px 0;font-size:17px;line-height:1.65;color:#1e293b;">
            <p style="margin:0 0 20px 0;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">Hi %s,</p>
            <p style="margin:0 0 18px 0;">Thank you for applying for the <strong style="color:#111827;font-weight:600;">%s</strong> role. We&rsquo;ve received your application and will review it carefully.</p>
            <p style="margin:0 0 22px 0;">We aim to get back to you <strong style="font-weight:600;color:#111827;">within 7 days</strong>.</p>
            <p style="margin:0;font-size:16px;line-height:1.6;color:#475569;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">If you have questions, reply to this message or write to <a href="mailto:hr@lynxlinkage.com" style="color:#1d4ed8;text-decoration:underline;">hr@lynxlinkage.com</a>.</p>
          </td>
        </tr>
        <tr>
          <td style="padding:24px 0 0 0;border-top:1px solid #cbd5e1;">
            <p style="margin:0;font-size:16px;line-height:1.55;color:#1e293b;">With appreciation,</p>
            <p style="margin:6px 0 0 0;font-size:16px;line-height:1.55;color:#111827;font-weight:600;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">%s</p>
            <p style="margin:20px 0 0 0;font-size:12px;line-height:1.5;color:#94a3b8;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">&copy; %d &middot; <a href="%s" style="color:#64748b;text-decoration:underline;">%s</a></p>
          </td>
        </tr>
      </table>
    </td>
  </tr>
</table>
</body>
</html>`, name, title, brand, year, u.String(), siteText)
}

func dispName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "there"
	}
	return s
}

func dispTitle(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "the position"
	}
	return s
}
