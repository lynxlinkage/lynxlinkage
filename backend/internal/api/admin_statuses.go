package api

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
)

// isUniqueViolation reports whether err is a Postgres "23505 unique
// violation" error. The status repo wraps a few of these as 409s so HR
// sees a friendly conflict message rather than a generic 500.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

type statusUpsertRequest struct {
	Slug         string `json:"slug"          validate:"omitempty,max=80"`
	Name         string `json:"name"          validate:"required,min=1,max=80"`
	Kind         string `json:"kind"          validate:"required,oneof=open accept reject"`
	Color        string `json:"color"         validate:"omitempty,max=20"`
	DisplayOrder *int   `json:"displayOrder"`
	IsDefault    *bool  `json:"isDefault"`
}

var slugAllowed = regexp.MustCompile(`[^a-z0-9-]+`)

// slugify produces a URL-safe lowercased identifier from name when the
// HR user didn't supply one explicitly.
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugAllowed.ReplaceAllString(s, "")
	s = strings.Trim(s, "-")
	if len(s) > 80 {
		s = s[:80]
	}
	return s
}

func (s *Server) handleAdminListStatuses(c *gin.Context) {
	items, err := s.Statuses.List(c.Request.Context())
	if err != nil {
		s.Logger.Error("admin list statuses", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) handleAdminCreateStatus(c *gin.Context) {
	body, ok := s.bindStatusBody(c)
	if !ok {
		return
	}

	slug := strings.TrimSpace(body.Slug)
	if slug == "" {
		slug = slugify(body.Name)
	}
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug could not be derived from name"})
		return
	}

	order := 0
	if body.DisplayOrder != nil {
		order = *body.DisplayOrder
	}
	isDefault := false
	if body.IsDefault != nil {
		isDefault = *body.IsDefault
	}

	status := &domain.ApplicationStatus{
		Slug:         slug,
		Name:         strings.TrimSpace(body.Name),
		Kind:         domain.ApplicationStatusKind(body.Kind),
		Color:        strings.TrimSpace(body.Color),
		DisplayOrder: order,
		IsDefault:    isDefault,
	}
	if _, err := s.Statuses.Create(c.Request.Context(), status); err != nil {
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "a status with that slug already exists"})
			return
		}
		s.Logger.Error("admin create status", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, status)
}

func (s *Server) handleAdminUpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	body, ok := s.bindStatusBody(c)
	if !ok {
		return
	}
	slug := strings.TrimSpace(body.Slug)
	if slug == "" {
		slug = slugify(body.Name)
	}
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug could not be derived from name"})
		return
	}

	order := 0
	if body.DisplayOrder != nil {
		order = *body.DisplayOrder
	}
	isDefault := false
	if body.IsDefault != nil {
		isDefault = *body.IsDefault
	}

	status := &domain.ApplicationStatus{
		ID:           id,
		Slug:         slug,
		Name:         strings.TrimSpace(body.Name),
		Kind:         domain.ApplicationStatusKind(body.Kind),
		Color:        strings.TrimSpace(body.Color),
		DisplayOrder: order,
		IsDefault:    isDefault,
	}
	if err := s.Statuses.Update(c.Request.Context(), status); err != nil {
		switch {
		case errors.Is(err, store.ErrStatusNoDefault):
			c.JSON(http.StatusConflict,
				gin.H{"error": "at least one status must be marked default"})
			return
		case store.IsNoRows(err) || errors.Is(err, store.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		case isUniqueViolation(err):
			c.JSON(http.StatusConflict, gin.H{"error": "a status with that slug already exists"})
			return
		}
		s.Logger.Error("admin update status", "err", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (s *Server) handleAdminDeleteStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := s.Statuses.Delete(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrStatusInUse):
			c.JSON(http.StatusConflict, gin.H{
				"error": "this status is still in use; reassign affected applications first",
			})
			return
		case errors.Is(err, store.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		s.Logger.Error("admin delete status", "err", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) bindStatusBody(c *gin.Context) (statusUpsertRequest, bool) {
	var body statusUpsertRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return body, false
	}
	body.Slug = strings.TrimSpace(body.Slug)
	body.Name = strings.TrimSpace(body.Name)
	body.Kind = strings.TrimSpace(body.Kind)
	body.Color = strings.TrimSpace(body.Color)
	if err := s.Validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return body, false
	}
	return body, true
}
