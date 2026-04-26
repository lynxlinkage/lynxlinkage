// Package api wires HTTP handlers for the public landing-page API. All
// endpoints under /api/v1 are read-only with the exception of /contact,
// which accepts a validated submission.
package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/auth"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
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
	Users     *store.UserRepo
	ContactRL *middleware.IPRateLimiter
	Auth      *auth.Manager
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

		// Authentication: open login + protected logout/me.
		v1auth := v1.Group("/auth")
		v1auth.POST("/login", s.handleLogin)
		v1auth.POST("/logout", s.handleLogout)
		v1auth.GET("/me", s.Auth.RequireAuth(), s.handleMe)

		// Admin: HR-only mutations on job postings.
		admin := v1.Group("/admin")
		admin.Use(s.Auth.RequireAuth())
		admin.Use(s.Auth.RequireRole(domain.RoleHR))
		{
			admin.GET("/jobs", s.handleAdminListJobs)
			admin.POST("/jobs", s.handleAdminCreateJob)
			admin.PUT("/jobs/:id", s.handleAdminUpdateJob)
			admin.DELETE("/jobs/:id", s.handleAdminDeleteJob)
		}
	}
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
