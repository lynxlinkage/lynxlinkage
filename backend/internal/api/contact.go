package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/mailout"
)

// contactRequest is the JSON payload accepted by POST /api/v1/contact.
//
// Validation rules:
//   - name 2..120 chars
//   - email valid email up to 254 chars
//   - company optional up to 200 chars
//   - message 3..5000 chars
//   - kind one of the recognised ContactKind values
type contactRequest struct {
	Name    string `json:"name"    validate:"required,min=2,max=120"`
	Email   string `json:"email"   validate:"required,email,max=254"`
	Company string `json:"company" validate:"omitempty,max=200"`
	Message string `json:"message" validate:"required,min=3,max=5000"`
	Kind    string `json:"kind"    validate:"omitempty,oneof=general partnership research hiring"`
}

func (s *Server) handleSubmitContact(c *gin.Context) {
	var req contactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if err := s.Validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":  "validation failed",
			"detail": err.Error(),
		})
		return
	}

	kind := domain.ContactKind(strings.TrimSpace(req.Kind))
	if kind == "" {
		kind = domain.ContactGeneral
	}

	sub := &domain.ContactSubmission{
		Name:      strings.TrimSpace(req.Name),
		Email:     strings.TrimSpace(strings.ToLower(req.Email)),
		Company:   strings.TrimSpace(req.Company),
		Message:   strings.TrimSpace(req.Message),
		Kind:      kind,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}
	id, err := s.Contact.Insert(c.Request.Context(), sub)
	if err != nil {
		s.Logger.Error("insert contact", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save submission"})
		return
	}
	s.Logger.Info("contact submitted", "id", id, "kind", kind, "email", sub.Email)
	s.sendContactEmails(id, sub)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// sendContactEmails sends (1) a receipt to the submitter and (2) a staff
// notification. If SMTP is not fully configured, both are skipped and a WARN
// is logged so operators know why inbox is silent.
func (s *Server) sendContactEmails(id int64, sub *domain.ContactSubmission) {
	if s.Mail == nil || !s.Mail.Ready() {
		s.Logger.Warn("contact: outgoing email skipped",
			"reason", "SMTP not configured or incomplete",
			"hint", "Set EMAIL_FROM, SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS — check server logs at startup for smtp warnings",
		)
		return
	}

	staffTo := strings.TrimSpace(s.ContactStaffTo)
	if staffTo == "" {
		staffTo = mailout.ContactStaffRecipient
	}

	go func() {
		ackSubject := mailout.ContactAckSubject()
		ackPlain := mailout.ContactAckBody(sub.Name, s.AppName)
		ackHTML := mailout.ContactAckHTML(sub.Name, s.AppName, s.SiteURL)
		if err := s.Mail.SendAlternative(sub.Email, ackSubject, ackPlain, ackHTML); err != nil {
			s.Logger.Error("contact: ack email failed", "err", err, "to", sub.Email)
		} else {
			s.Logger.Info("contact: ack email sent", "to", sub.Email)
		}

		staffSubject := mailout.ContactStaffSubject(id, string(sub.Kind))
		staffPlain := mailout.ContactStaffPlain(
			id,
			sub.Name,
			sub.Email,
			sub.Company,
			string(sub.Kind),
			sub.Message,
			sub.IPAddress,
			sub.UserAgent,
		)
		staffHTML := mailout.ContactStaffHTML(
			id,
			sub.Name,
			sub.Email,
			sub.Company,
			string(sub.Kind),
			sub.Message,
			sub.IPAddress,
			sub.UserAgent,
		)
		if err := s.Mail.SendAlternative(staffTo, staffSubject, staffPlain, staffHTML); err != nil {
			s.Logger.Error("contact: staff notify email failed", "err", err, "to", staffTo)
		} else {
			s.Logger.Info("contact: staff notify email sent", "to", staffTo, "submission_id", id)
		}
	}()
}
