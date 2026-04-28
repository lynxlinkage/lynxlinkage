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
