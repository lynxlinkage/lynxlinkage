package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	Env             string
	HTTPAddr        string
	DatabaseURL     string
	LogLevel        string
	CORSAllowOrigin string

	// Rate limit for the public contact endpoint (requests per IP per window).
	ContactRPS   float64
	ContactBurst int

	// Rate limit for the public job-application endpoint. Defaults are
	// looser than contact since the form has a higher bar to fill in.
	ApplicationRPS   float64
	ApplicationBurst int

	// Upload limits for application attachments.
	UploadDir           string
	MaxUploadFiles      int
	MaxUploadFileBytes  int64
	MaxUploadTotalBytes int64

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// EmailFrom / EmailTo are optional. When set, contact submissions are
	// also forwarded by email; otherwise they're stored only.
	EmailFrom string
	EmailTo   string
	SMTPHost  string
	SMTPPort  int
	SMTPUser  string
	SMTPPass  string

	// SessionSecret signs the HMAC of session cookies. Rotating this value
	// invalidates every outstanding session.
	SessionSecret string
	SessionTTL    time.Duration

	// AppName is used in user-facing text (e.g. application ack emails).
	AppName string

	// SiteURL is the public https origin (no trailing slash), used in HTML
	// emails for absolute asset and link URLs.
	SiteURL string
}

// Load reads configuration from environment variables, applying defaults
// suitable for local development. A `.env` file in the working directory
// (or in `backend/`) is loaded transparently if present so dev workflows
// don't need to source it manually; existing environment variables
// always take precedence.
func Load() (Config, error) {
	loadDotenv()

	cfg := Config{
		Env:             getEnv("APP_ENV", "development"),
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgresql://lynxlinkage:lynxlinkage@localhost:5432/lynxlinkage?sslmode=disable"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGIN", "http://localhost:5173"),
		EmailFrom: getEnv("EMAIL_FROM", ""),
		EmailTo:   getEnv("EMAIL_TO", ""),
	}

	smtpURL := strings.TrimSpace(getEnv("SMTP_URL", ""))
	urlHost, urlPort := parseSMTPAddress(smtpURL)
	smtpHost := strings.TrimSpace(getEnv("SMTP_HOST", ""))
	if smtpHost == "" {
		smtpHost = urlHost
	}
	cfg.SMTPHost = smtpHost

	smtpUser := firstNonEmpty(getEnv("SMTP_USER", ""), getEnv("SMTP_USERNAME", ""))
	smtpUser = strings.TrimSpace(smtpUser)
	cfg.SMTPUser = smtpUser

	smtpPass := firstNonEmpty(getEnv("SMTP_PASS", ""), getEnv("SMTP_PASSWORD", ""))
	smtpPass = strings.TrimSpace(smtpPass)
	cfg.SMTPPass = smtpPass

	rps, err := strconv.ParseFloat(getEnv("CONTACT_RPS", "0.2"), 64)
	if err != nil {
		return cfg, fmt.Errorf("invalid CONTACT_RPS: %w", err)
	}
	cfg.ContactRPS = rps

	burst, err := strconv.Atoi(getEnv("CONTACT_BURST", "3"))
	if err != nil {
		return cfg, fmt.Errorf("invalid CONTACT_BURST: %w", err)
	}
	cfg.ContactBurst = burst

	smtpPort := 587
	if v, ok := os.LookupEnv("SMTP_PORT"); ok && strings.TrimSpace(v) != "" {
		p, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return cfg, fmt.Errorf("invalid SMTP_PORT: %w", err)
		}
		smtpPort = p
	} else if urlPort > 0 {
		smtpPort = urlPort
	}
	cfg.SMTPPort = smtpPort

	appRPS, err := strconv.ParseFloat(getEnv("APPLICATION_RPS", "0.1"), 64)
	if err != nil {
		return cfg, fmt.Errorf("invalid APPLICATION_RPS: %w", err)
	}
	cfg.ApplicationRPS = appRPS
	appBurst, err := strconv.Atoi(getEnv("APPLICATION_BURST", "3"))
	if err != nil {
		return cfg, fmt.Errorf("invalid APPLICATION_BURST: %w", err)
	}
	cfg.ApplicationBurst = appBurst

	cfg.UploadDir = getEnv("UPLOAD_DIR", "./data/uploads")

	maxFiles, err := strconv.Atoi(getEnv("MAX_UPLOAD_FILES", "3"))
	if err != nil {
		return cfg, fmt.Errorf("invalid MAX_UPLOAD_FILES: %w", err)
	}
	cfg.MaxUploadFiles = maxFiles

	perFile, err := strconv.ParseInt(getEnv("MAX_UPLOAD_FILE_BYTES", "10485760"), 10, 64) // 10 MiB
	if err != nil {
		return cfg, fmt.Errorf("invalid MAX_UPLOAD_FILE_BYTES: %w", err)
	}
	cfg.MaxUploadFileBytes = perFile

	totalDefault := strconv.FormatInt(int64(cfg.MaxUploadFiles)*cfg.MaxUploadFileBytes+(2<<20), 10)
	totalCap, err := strconv.ParseInt(getEnv("MAX_UPLOAD_TOTAL_BYTES", totalDefault), 10, 64)
	if err != nil {
		return cfg, fmt.Errorf("invalid MAX_UPLOAD_TOTAL_BYTES: %w", err)
	}
	cfg.MaxUploadTotalBytes = totalCap

	// Defaults bumped from 10s/15s so multi-MB uploads on slow client
	// connections don't get cut mid-stream. Operators can lower for tiny
	// JSON-only deployments.
	cfg.ReadTimeout = mustDuration(getEnv("READ_TIMEOUT", "60s"))
	cfg.WriteTimeout = mustDuration(getEnv("WRITE_TIMEOUT", "60s"))
	cfg.IdleTimeout = mustDuration(getEnv("IDLE_TIMEOUT", "60s"))

	cfg.SessionSecret = getEnv("SESSION_SECRET", "")
	if cfg.SessionSecret == "" {
		if cfg.Env == "production" {
			return cfg, fmt.Errorf("SESSION_SECRET must be set in production")
		}
		// Stable but obviously-non-secret default for local development so
		// developers don't have to fish a value out of an .env on first run.
		cfg.SessionSecret = "dev-only-not-secret-change-me"
	}
	cfg.SessionTTL = mustDuration(getEnv("SESSION_TTL", "168h"))

	cfg.AppName = getEnv("APP_NAME", "LynxLinkage")

	cfg.SiteURL = strings.TrimSuffix(strings.TrimSpace(getEnv("SITE_URL", "https://lynxlinkage.com")), "/")

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// loadDotenv walks up from the current working directory looking for a
// `.env` file (or `backend/.env`) and loads it into the environment.
// Variables that are already exported in the shell win over the file,
// matching the behaviour of `godotenv.Load` itself.
func loadDotenv() {
	candidates := []string{
		".env",
		filepath.Join("backend", ".env"),
		filepath.Join("..", ".env"),
		filepath.Join("..", "backend", ".env"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			_ = godotenv.Load(p)
			return
		}
	}
}

func mustDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("invalid duration %q: %v", s, err))
	}
	return d
}

// firstNonEmpty returns the first non-empty (after trim) string.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// parseSMTPAddress parses "host:port" from SMTP_URL; if there is no port, host
// is the whole string and the returned port is 0.
func parseSMTPAddress(s string) (host string, port int) {
	if s == "" {
		return "", 0
	}
	h, pstr, err := net.SplitHostPort(s)
	if err != nil {
		return s, 0
	}
	p, err := strconv.Atoi(pstr)
	if err != nil {
		return h, 0
	}
	return h, p
}
