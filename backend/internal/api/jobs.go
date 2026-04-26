package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
)

func (s *Server) handleListJobs(c *gin.Context) {
	jobs, err := s.Jobs.ListActive(c.Request.Context())
	if err != nil {
		s.Logger.Error("list jobs", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": jobs})
}

func (s *Server) handleGetJob(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	job, err := s.Jobs.Get(c.Request.Context(), id)
	if err != nil {
		if store.IsNoRows(err) || errors.Is(err, errNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		s.Logger.Error("get job", "err", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, job)
}

var errNotFound = errors.New("not found")
