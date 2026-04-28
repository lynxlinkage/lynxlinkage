package mailout

import (
	"strings"
	"testing"
)

func TestApplicationAckBody(t *testing.T) {
	body := ApplicationAckBody("Alex", "Quant Researcher", "LynxLinkage")
	if !strings.Contains(body, "Alex") || !strings.Contains(body, "Quant Researcher") {
		t.Fatalf("body missing name or title: %q", body)
	}
	if !strings.Contains(body, "7 days") {
		t.Fatalf("body should mention reply window: %q", body)
	}
	if !strings.Contains(body, "hr@lynxlinkage.com") {
		t.Fatalf("body should mention HR contact: %q", body)
	}
}

func TestApplicationAckSubject(t *testing.T) {
	s := ApplicationAckSubject("Engineer")
	if s == "" || !strings.Contains(s, "Engineer") {
		t.Fatalf("subject: %q", s)
	}
}

func TestContactAckBody(t *testing.T) {
	body := ContactAckBody("Jordan", "LynxLinkage")
	if !strings.Contains(body, "Jordan") || !strings.Contains(body, "received your message") {
		t.Fatalf("body: %q", body)
	}
	if !strings.Contains(body, "No reply to this message") {
		t.Fatalf("body should state noreply: %q", body)
	}
	if !strings.Contains(body, "eddy@lynxlinkage.com") {
		t.Fatalf("body should mention contact email: %q", body)
	}
}

func TestContactStaffPlain(t *testing.T) {
	body := ContactStaffPlain(7, "Ada", "ada@example.com", "Co", "research", "Hello\nthere", "1.2.3.4", "Mozilla/5.0")
	if !strings.Contains(body, "Ada") || !strings.Contains(body, "ada@example.com") || !strings.Contains(body, "Hello\nthere") {
		t.Fatalf("body: %q", body)
	}
	if !strings.Contains(body, "7") || !strings.Contains(body, "research") {
		t.Fatalf("body should include id and kind: %q", body)
	}
}

func TestContactStaffSubjectLine(t *testing.T) {
	s := ContactStaffSubject(12, "partnership")
	if s == "" || !strings.Contains(s, "12") || !strings.Contains(s, "partnership") {
		t.Fatalf("subject: %q", s)
	}
}
