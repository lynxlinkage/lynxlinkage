// Package api wires HTTP handlers for the public landing-page API. All
// endpoints under /api/v1 are read-only with the exception of /contact,
// which accepts a validated submission.
package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/middleware"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
)

// Server bundles the dependencies handlers need.
type Server struct {
	Logger    *slog.Logger
	Validate  *validator.Validate
	Research  *store.ResearchRepo
	Jobs      *store.JobRepo
	Partners  *store.PartnerRepo
	Contact   *store.ContactRepo
	ContactRL *middleware.IPRateLimiter
}

// Register mounts all routes under /api/v1 on the provided router.
func (s *Server) Register(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", s.handleHealth)
		v1.GET("/researches", s.handleListResearches)
		v1.GET("/jobs", s.handleListJobs)
		v1.GET("/jobs/:id", s.handleGetJob)
		v1.GET("/partners", s.handleListPartners)

		contact := v1.Group("/contact")
		if s.ContactRL != nil {
			contact.Use(s.ContactRL.Middleware())
		}
		contact.POST("", s.handleSubmitContact)
	}
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
