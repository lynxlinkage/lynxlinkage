package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
)

func (s *Server) handleListResearches(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit < 0 || limit > 100 {
		limit = 0
	}
	cards, err := s.Research.List(c.Request.Context(), store.ListResearchOpts{
		Tag:   c.Query("tag"),
		Limit: limit,
	})
	if err != nil {
		s.Logger.Error("list researches", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": cards})
}
