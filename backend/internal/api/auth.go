package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/auth"
)

// loginRequest is the JSON body posted to /api/v1/auth/login.
type loginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

func (s *Server) handleLogin(c *gin.Context) {
	var body loginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if err := s.Validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(body.Email))
	user, err := s.Users.GetByEmail(c.Request.Context(), email)
	if err != nil {
		// Constant-time-ish: still run a bcrypt compare against a dummy hash
		// to avoid leaking that the email was unknown via response time.
		_ = s.Auth.CheckPassword("$2a$12$invalidsaltinvalidsaltOXyfcmwL5L8WfVpYkHa4u5GfH7ek2", body.Password)
		s.Logger.Info("login failed: unknown email", "email", email)
		time.Sleep(120 * time.Millisecond)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if !s.Auth.CheckPassword(user.PasswordHash, body.Password) {
		s.Logger.Info("login failed: bad password", "email", email)
		time.Sleep(120 * time.Millisecond)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Best-effort: failing to record the login timestamp shouldn't block
	// the user from signing in.
	if err := s.Users.TouchLogin(c.Request.Context(), user.ID); err != nil {
		s.Logger.Warn("touch login failed", "err", err, "user_id", user.ID)
	}

	s.Auth.Issue(c, user.ID)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (s *Server) handleLogout(c *gin.Context) {
	s.Auth.Clear(c)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (s *Server) handleMe(c *gin.Context) {
	user, ok := auth.UserFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
