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
	s.sendContactAck(sub.Name, sub.Email)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// sendContactAck emails the visitor a receipt after storing the submission.
// Failures are logged only; HTTP success is already assured.
func (s *Server) sendContactAck(name, to string) {
	if s.Mail == nil || !s.Mail.Ready() {
		return
	}
	subject := mailout.ContactAckSubject()
	plain := mailout.ContactAckBody(name, s.AppName)
	html := mailout.ContactAckHTML(name, s.AppName, s.SiteURL)
	recipient := to
	go func() {
		if err := s.Mail.SendAlternative(recipient, subject, plain, html); err != nil {
			s.Logger.Error("contact: ack email", "err", err, "to", recipient)
			return
		}
		s.Logger.Info("contact: ack email sent", "to", recipient)
	}()
}
