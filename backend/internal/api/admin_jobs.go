package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/auth"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
)

// actorID returns a pointer to the authenticated user's id, or nil when
// no user is on the context (which shouldn't happen behind RequireAuth).
func actorID(c *gin.Context) *int64 {
	u, ok := auth.UserFrom(c)
	if !ok {
		return nil
	}
	id := u.ID
	return &id
}

// jobUpsertRequest is the JSON shape posted to admin job create/update
// endpoints. PostedAt is intentionally absent — it is set to NOW() on
// creation and never changed on subsequent edits.
type jobUpsertRequest struct {
	Title           string `json:"title"             validate:"required,min=1,max=200"`
	Team            string `json:"team"              validate:"max=80"`
	Location        string `json:"location"          validate:"max=120"`
	EmploymentType  string `json:"employmentType"    validate:"required,oneof=full_time part_time contract internship"`
	DescriptionMd   string `json:"descriptionMd"     validate:"max=20000"`
	ApplyUrlOrEmail string `json:"applyUrlOrEmail"   validate:"required,max=500"`
	IsActive        *bool  `json:"isActive"`
}

func (s *Server) handleAdminListJobs(c *gin.Context) {
	jobs, err := s.Jobs.ListAll(c.Request.Context())
	if err != nil {
		s.Logger.Error("admin list jobs", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": jobs})
}

func (s *Server) handleAdminCreateJob(c *gin.Context) {
	body, ok := s.bindJobBody(c)
	if !ok {
		return
	}
	posting := s.toJobPosting(0, body)
	if _, err := s.Jobs.Create(c.Request.Context(), posting, actorID(c)); err != nil {
		s.Logger.Error("admin create job", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, posting)
}

func (s *Server) handleAdminDeleteJob(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := s.Jobs.Delete(c.Request.Context(), id); err != nil {
		if store.IsNoRows(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		s.Logger.Error("admin delete job", "err", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) handleAdminUpdateJob(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	body, ok := s.bindJobBody(c)
	if !ok {
		return
	}
	posting := s.toJobPosting(id, body)
	if err := s.Jobs.Update(c.Request.Context(), posting, actorID(c)); err != nil {
		if store.IsNoRows(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		s.Logger.Error("admin update job", "err", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, posting)
}

func (s *Server) bindJobBody(c *gin.Context) (jobUpsertRequest, bool) {
	var body jobUpsertRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return body, false
	}
	body.Title = strings.TrimSpace(body.Title)
	body.Team = strings.TrimSpace(body.Team)
	body.Location = strings.TrimSpace(body.Location)
	body.ApplyUrlOrEmail = strings.TrimSpace(body.ApplyUrlOrEmail)
	if err := s.Validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return body, false
	}
	return body, true
}

func (s *Server) toJobPosting(id int64, body jobUpsertRequest) *domain.JobPosting {
	active := true
	if body.IsActive != nil {
		active = *body.IsActive
	}
	return &domain.JobPosting{
		ID:              id,
		Title:           body.Title,
		Team:            body.Team,
		Location:        body.Location,
		EmploymentType:  domain.EmploymentType(body.EmploymentType),
		DescriptionMD:   body.DescriptionMd,
		ApplyURLOrEmail: body.ApplyUrlOrEmail,
		PostedAt:        time.Now().UTC(),
		IsActive:        active,
	}
}
