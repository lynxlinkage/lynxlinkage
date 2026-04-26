package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
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
}

// Load reads configuration from environment variables, applying defaults
// suitable for local development.
func Load() (Config, error) {
	cfg := Config{
		Env:             getEnv("APP_ENV", "development"),
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "file:./data/lynxlinkage.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGIN", "http://localhost:5173"),
		EmailFrom:       getEnv("EMAIL_FROM", ""),
		EmailTo:         getEnv("EMAIL_TO", ""),
		SMTPHost:        getEnv("SMTP_HOST", ""),
		SMTPUser:        getEnv("SMTP_USER", ""),
		SMTPPass:        getEnv("SMTP_PASS", ""),
	}

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

	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if err != nil {
		return cfg, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}
	cfg.SMTPPort = smtpPort

	cfg.ReadTimeout = mustDuration(getEnv("READ_TIMEOUT", "10s"))
	cfg.WriteTimeout = mustDuration(getEnv("WRITE_TIMEOUT", "15s"))
	cfg.IdleTimeout = mustDuration(getEnv("IDLE_TIMEOUT", "60s"))

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func mustDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("invalid duration %q: %v", s, err))
	}
	return d
}
