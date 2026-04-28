// Package middleware contains cross-cutting Gin middlewares: structured
// request logging, panic recovery, CORS, and a per-IP rate limiter.
package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RequestLogger logs every request as a single structured line.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("http",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
			"duration_ms", time.Since(start).Milliseconds(),
			"ip", clientIP(c),
		)
	}
}

// Recover converts panics into a 500 response and logs them.
func Recover(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err any) {
		logger.Error("panic", "err", err, "path", c.Request.URL.Path)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": "internal server error"})
	})
}

// CORS allows the configured origin(s) to call the API. The frontend is
// served from the same origin as the API in production, but this lets the
// Vite dev server talk to the backend on :8080.
func CORS(allowOrigin string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, o := range strings.Split(allowOrigin, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			allowed[o] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if _, ok := allowed[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")
			c.Header("Access-Control-Max-Age", "86400")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// IPRateLimiter is a simple per-IP token bucket rate limiter, suitable for
// the contact endpoint. Buckets are kept in-memory; for multi-instance
// deployments swap in Redis.
type IPRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*rate.Limiter
	rps     rate.Limit
	burst   int
}

func NewIPRateLimiter(rps float64, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		buckets: map[string]*rate.Limiter{},
		rps:     rate.Limit(rps),
		burst:   burst,
	}
}

func (l *IPRateLimiter) limiterFor(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	lim, ok := l.buckets[ip]
	if !ok {
		lim = rate.NewLimiter(l.rps, l.burst)
		l.buckets[ip] = lim
	}
	return lim
}

// Middleware returns a Gin handler that 429s requests over the limit.
func (l *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := clientIP(c)
		if !l.limiterFor(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}

// clientIP prefers X-Forwarded-For (first hop) then the remote addr.
// Gin's c.ClientIP() already does most of this when TrustedProxies is set,
// but we keep it explicit and trust only the first hop.
func clientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		if comma := strings.IndexByte(xff, ','); comma > 0 {
			return strings.TrimSpace(xff[:comma])
		}
		return strings.TrimSpace(xff)
	}
	if real := c.GetHeader("X-Real-IP"); real != "" {
		return real
	}
	return c.ClientIP()
}
