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
	"github.com/lynxlinkage/lynxlinkage/backend/internal/uploads"
)

// Server bundles the dependencies handlers need.
type Server struct {
	Logger       *slog.Logger
	Validate     *validator.Validate
	Research     *store.ResearchRepo
	Jobs         *store.JobRepo
	Partners     *store.PartnerRepo
	Contact      *store.ContactRepo
	Users        *store.UserRepo
	Applications *store.ApplicationRepo
	Statuses     *store.StatusRepo
	Uploads      *uploads.Store
	ContactRL    *middleware.IPRateLimiter
	ApplyRL      *middleware.IPRateLimiter
	Auth         *auth.Manager

	// Upload limits surfaced from config so handlers don't reach back
	// into env-loading code.
	MaxUploadFiles      int
	MaxUploadFileBytes  int64
	MaxUploadTotalBytes int64
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

		// Public: candidate applications. Rate-limited per IP separately
		// from contact so a campaign spam-attempt against one form
		// doesn't lock the other out.
		apply := v1.Group("/jobs/:id/applications")
		if s.ApplyRL != nil {
			apply.Use(s.ApplyRL.Middleware())
		}
		apply.POST("", s.handleSubmitApplication)

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

			admin.GET("/applications", s.handleAdminListApplications)
			admin.GET("/applications/:id", s.handleAdminGetApplication)
			admin.PUT("/applications/:id/status", s.handleAdminUpdateApplicationStatus)
			admin.GET("/applications/:id/files/:fileId", s.handleAdminDownloadApplicationFile)

			admin.GET("/application-statuses", s.handleAdminListStatuses)
			admin.POST("/application-statuses", s.handleAdminCreateStatus)
			admin.PUT("/application-statuses/:id", s.handleAdminUpdateStatus)
			admin.DELETE("/application-statuses/:id", s.handleAdminDeleteStatus)
		}
	}
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
