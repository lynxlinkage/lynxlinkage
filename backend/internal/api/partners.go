package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleListPartners(c *gin.Context) {
	partners, err := s.Partners.List(c.Request.Context())
	if err != nil {
		s.Logger.Error("list partners", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": partners})
}
