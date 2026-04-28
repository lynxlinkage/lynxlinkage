// Package auth provides session helpers (HMAC-signed cookies),
// password hashing (bcrypt), and Gin middleware that gates handlers
// behind authentication and role checks.
//
// The session token has the form "<userID>.<expUnix>.<sigB64>" where the
// signature is HMAC-SHA256 of "<userID>.<expUnix>" using the configured
// session secret. Rotating the secret invalidates all outstanding
// sessions; that's the simplest revocation story for v1.
package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// CookieName is the HttpOnly session cookie set by the API.
const CookieName = "lynx_session"

// DefaultTTL is how long a freshly-issued session lasts.
const DefaultTTL = 7 * 24 * time.Hour

// BcryptCost is the work factor used when hashing passwords. Reasonable
// default; bumped above the bcrypt default of 10 because we expect tiny
// login volumes and care more about brute-force resistance than throughput.
const BcryptCost = 12

// ErrInvalidSession is returned when a cookie is missing, malformed, or
// fails signature/expiry validation.
var ErrInvalidSession = errors.New("invalid session")

// Manager is the auth helper carried around by the API server.
type Manager struct {
	Secret []byte
	TTL    time.Duration
	Users  *store.UserRepo
	Secure bool // set Secure flag on cookies (true in production)
}

// NewManager constructs an auth manager. Panics if secret is empty so
// misconfigurations fail loudly at startup rather than silently issuing
// signable-by-anyone tokens.
func NewManager(secret string, ttl time.Duration, users *store.UserRepo, secure bool) *Manager {
	if secret == "" {
		panic("auth.NewManager: empty session secret")
	}
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return &Manager{
		Secret: []byte(secret),
		TTL:    ttl,
		Users:  users,
		Secure: secure,
	}
}

// HashPassword returns a bcrypt hash for the given plaintext.
func (m *Manager) HashPassword(plain string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// CheckPassword reports whether plain matches the bcrypt hash.
func (m *Manager) CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// Issue creates a signed session token and writes it to the response as an
// HttpOnly cookie. The cookie path is `/` so it accompanies all API and
// frontend requests.
func (m *Manager) Issue(c *gin.Context, userID int64) {
	exp := time.Now().Add(m.TTL).Unix()
	tok := m.signToken(userID, exp)

	sameSite := http.SameSiteStrictMode
	c.SetSameSite(sameSite)
	c.SetCookie(CookieName, tok, int(m.TTL.Seconds()), "/", "", m.Secure, true)
}

// Clear deletes the session cookie.
func (m *Manager) Clear(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(CookieName, "", -1, "/", "", m.Secure, true)
}

// Verify parses and validates a token, returning the embedded user ID.
func (m *Manager) Verify(token string) (int64, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, ErrInvalidSession
	}
	uid, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || uid <= 0 {
		return 0, ErrInvalidSession
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, ErrInvalidSession
	}
	if time.Now().Unix() > exp {
		return 0, ErrInvalidSession
	}
	gotSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return 0, ErrInvalidSession
	}
	wantSig := m.sign(parts[0] + "." + parts[1])
	if subtle.ConstantTimeCompare(gotSig, wantSig) != 1 {
		return 0, ErrInvalidSession
	}
	return uid, nil
}

func (m *Manager) signToken(userID, exp int64) string {
	payload := fmt.Sprintf("%d.%d", userID, exp)
	sig := m.sign(payload)
	return payload + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func (m *Manager) sign(payload string) []byte {
	mac := hmac.New(sha256.New, m.Secret)
	mac.Write([]byte(payload))
	return mac.Sum(nil)
}

// LoadUser reads the session cookie, verifies it, and loads the user.
func (m *Manager) LoadUser(ctx context.Context, c *gin.Context) (*domain.User, error) {
	tok, err := c.Cookie(CookieName)
	if err != nil || tok == "" {
		return nil, ErrInvalidSession
	}
	uid, err := m.Verify(tok)
	if err != nil {
		return nil, err
	}
	user, err := m.Users.GetByID(ctx, uid)
	if err != nil {
		return nil, ErrInvalidSession
	}
	return user, nil
}

// userKey is the gin.Context key under which an authenticated user is
// stored by the RequireAuth middleware.
const userKey = "auth.user"

// RequireAuth returns a Gin middleware that requires a valid session
// cookie and attaches the user to the request context.
func (m *Manager) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := m.LoadUser(c.Request.Context(), c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set(userKey, user)
		c.Next()
	}
}

// RequireRole returns a Gin middleware that requires the authenticated
// user to hold one of the supplied roles.
func (m *Manager) RequireRole(roles ...domain.Role) gin.HandlerFunc {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		u, ok := UserFrom(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if _, ok := allowed[u.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

// UserFrom retrieves the authenticated user from the Gin context. It
// returns (nil, false) when the request was not authenticated.
func UserFrom(c *gin.Context) (*domain.User, bool) {
	v, ok := c.Get(userKey)
	if !ok {
		return nil, false
	}
	u, ok := v.(*domain.User)
	return u, ok
}
