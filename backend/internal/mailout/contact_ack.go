package mailout

import (
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"
)

// ContactAckSubject is the acknowledgment subject line for contact form submits.
func ContactAckSubject() string {
	return "We received your message"
}

// ContactAckBody is the plain-text body for contact form acknowledgment.
func ContactAckBody(candidateName, brand string) string {
	name := strings.TrimSpace(candidateName)
	if name == "" {
		name = "there"
	}
	brand = strings.TrimSpace(brand)
	if brand == "" {
		brand = "LynxLinkage"
	}
	return fmt.Sprintf(`Hi %s,

Thank you for contacting us. We've received your message and will review it soon.

No reply to this message.
Reach us at eddy@lynxlinkage.com.

Best regards,
%s
`, name, brand)
}

// ContactAckHTML is the HTML counterpart (same styling family as ApplicationAckHTML).
func ContactAckHTML(candidateName, brand, siteURL string) string {
	name := html.EscapeString(dispName(candidateName))
	brandEscaped := html.EscapeString(strings.TrimSpace(brand))
	if brandEscaped == "" {
		brandEscaped = "LynxLinkage"
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
            <p style="margin:0 0 18px 0;">Thank you for contacting us. We&rsquo;ve received your message and will review it soon.</p>
            <p style="margin:0 0 12px 0;font-size:16px;line-height:1.6;color:#475569;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">No reply to this message.</p>
            <p style="margin:0;font-size:16px;line-height:1.6;color:#475569;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">Reach us at <a href="mailto:eddy@lynxlinkage.com" style="color:#1d4ed8;text-decoration:underline;">eddy@lynxlinkage.com</a>.</p>
          </td>
        </tr>
        <tr>
          <td style="padding:24px 0 0 0;border-top:1px solid #cbd5e1;">
            <p style="margin:0;font-size:16px;line-height:1.55;color:#1e293b;">Best,</p>
            <p style="margin:6px 0 0 0;font-size:16px;line-height:1.55;color:#111827;font-weight:600;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">%s</p>
            <p style="margin:20px 0 0 0;font-size:12px;line-height:1.5;color:#94a3b8;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">&copy; %d &middot; <a href="%s" style="color:#64748b;text-decoration:underline;">%s</a></p>
          </td>
        </tr>
      </table>
    </td>
  </tr>
</table>
</body>
</html>`, name, brandEscaped, year, u.String(), siteText)
}
