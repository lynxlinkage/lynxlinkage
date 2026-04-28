package mailout

import (
	"fmt"
	"html"
	"strings"
)

// ContactStaffRecipient is where internal contact-form notifications are delivered.
const ContactStaffRecipient = "eddy@lynxlinkage.com"

// ContactStaffSubject is the subject line for the staff notification email.
func ContactStaffSubject(id int64, kind string) string {
	k := stripHeaderField(kind)
	if k == "" {
		k = "general"
	}
	return stripHeaderField(fmt.Sprintf("[LynxLinkage] New contact #%d · %s", id, k))
}

// ContactStaffPlain is the plain-text staff notification body.
func ContactStaffPlain(id int64, name, email, company, kind, message, ip, userAgent string) string {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	company = strings.TrimSpace(company)
	if company == "" {
		company = "(not provided)"
	}
	k := strings.TrimSpace(kind)
	if k == "" {
		k = "general"
	}
	ua := strings.TrimSpace(userAgent)
	if len(ua) > 500 {
		ua = ua[:500] + "…"
	}
	if ua == "" {
		ua = "(not provided)"
	}
	return fmt.Sprintf(`New contact form submission

ID:           %d
Kind:         %s
Name:         %s
Email:        %s
Company:      %s

Message:
%s

IP:           %s
User-Agent:   %s
`, id, k, name, email, company, message, strings.TrimSpace(ip), ua)
}

// ContactStaffHTML is the HTML staff notification (same information, escaped).
func ContactStaffHTML(id int64, name, email, company, kind, message, ip, userAgent string) string {
	email = strings.TrimSpace(email)
	comp := strings.TrimSpace(company)
	if comp == "" {
		comp = "(not provided)"
	}
	k := strings.TrimSpace(kind)
	if k == "" {
		k = "general"
	}
	ua := strings.TrimSpace(userAgent)
	if len(ua) > 500 {
		ua = ua[:500] + "…"
	}
	if ua == "" {
		ua = "(not provided)"
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width"></head>
<body style="margin:0;padding:24px;background:#f8fafc;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;font-size:15px;line-height:1.5;color:#0f172a;">
<p style="margin:0 0 16px 0;font-weight:600;">New contact form submission</p>
<table role="presentation" cellspacing="0" cellpadding="0" style="border-collapse:collapse;background:#fff;border:1px solid #e2e8f0;border-radius:8px;max-width:640px;">
<tr><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;"><strong>ID</strong></td><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;">%d</td></tr>
<tr><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;"><strong>Kind</strong></td><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;">%s</td></tr>
<tr><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;"><strong>Name</strong></td><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;">%s</td></tr>
<tr><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;"><strong>Email</strong></td><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;"><a href="mailto:%s" style="color:#1d4ed8;">%s</a></td></tr>
<tr><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;vertical-align:top;"><strong>Company</strong></td><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;">%s</td></tr>
<tr><td colspan="2" style="padding:14px 18px;border-bottom:1px solid #e2e8f0;vertical-align:top;"><strong>Message</strong></td></tr>
<tr><td colspan="2" style="padding:14px 18px;border-bottom:1px solid #e2e8f0;"><pre style="margin:0;white-space:pre-wrap;word-break:break-word;font-family:inherit;font-size:14px;">%s</pre></td></tr>
<tr><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;vertical-align:top;"><strong>IP</strong></td><td style="padding:14px 18px;border-bottom:1px solid #e2e8f0;font-size:13px;">%s</td></tr>
<tr><td style="padding:14px 18px;vertical-align:top;"><strong>User-Agent</strong></td><td style="padding:14px 18px;font-size:12px;color:#475569;">%s</td></tr>
</table>
</body></html>`,
		id,
		html.EscapeString(k),
		html.EscapeString(strings.TrimSpace(name)),
		email,
		html.EscapeString(email),
		html.EscapeString(comp),
		html.EscapeString(message),
		html.EscapeString(strings.TrimSpace(ip)),
		html.EscapeString(ua),
	)
}
