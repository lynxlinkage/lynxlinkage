package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/auth"
)

// handleSession exchanges an Authelia-validated identity for a user object.
// Traefik's Authelia forwardauth middleware injects Remote-Email before this
// handler is reached; RequireAuth reads it and populates the gin context.
func (s *Server) handleSession(c *gin.Context) {
	user, ok := auth.UserFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := s.Users.TouchLogin(c.Request.Context(), user.ID); err != nil {
		s.Logger.Warn("touch login failed", "err", err, "user_id", user.ID)
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// handleLogout clears the Authelia session by redirecting to the Authelia
// logout endpoint. The browser follows the redirect and lands on the
// Authelia-hosted logout page.
func (s *Server) handleLogout(c *gin.Context) {
	c.Redirect(http.StatusFound, "https://auth."+s.SiteDomain+"/logout")
}

func (s *Server) handleMe(c *gin.Context) {
	user, ok := auth.UserFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
