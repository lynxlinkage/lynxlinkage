// Command server runs the lynxlinkage backend HTTP server. It serves the
// public JSON API under /api/v1 and (when built with -tags=embed) the
// prerendered SvelteKit frontend on every other path.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/api"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/auth"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/config"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/middleware"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/static"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := newLogger(cfg)
	slog.SetDefault(logger)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := ensureSQLiteDir(cfg.DatabaseURL); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := store.Open(ctx, cfg.DatabaseURL, logger)
	if err != nil {
		return err
	}
	defer db.Close()

	contactRL := middleware.NewIPRateLimiter(cfg.ContactRPS, cfg.ContactBurst)
	users := store.NewUserRepo(db)
	authMgr := auth.NewManager(cfg.SessionSecret, cfg.SessionTTL, users, cfg.Env == "production")

	server := &api.Server{
		Logger:    logger,
		Validate:  validator.New(validator.WithRequiredStructEnabled()),
		Research:  store.NewResearchRepo(db),
		Jobs:      store.NewJobRepo(db),
		Partners:  store.NewPartnerRepo(db),
		Contact:   store.NewContactRepo(db),
		Users:     users,
		ContactRL: contactRL,
		Auth:      authMgr,
	}

	r := gin.New()
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.Recover(logger))
	r.Use(middleware.CORS(cfg.CORSAllowOrigin))

	server.Register(r)

	// Catch-all serves the embedded frontend. Registered with NoRoute so it
	// runs only when no /api/v1/* handler matched.
	frontend := static.Handler(static.FrontendFS)
	r.NoRoute(frontend)

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", cfg.HTTPAddr, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case sig := <-stop:
		logger.Info("shutting down", "signal", sig.String())
	case err := <-errCh:
		if err != nil {
			return err
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	return srv.Shutdown(shutdownCtx)
}

func newLogger(cfg config.Config) *slog.Logger {
	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToUpper(cfg.LogLevel))); err != nil {
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: level}
	if cfg.Env == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

// ensureSQLiteDir creates the parent directory of a SQLite DSN like
// "file:./data/foo.db?...". It's a no-op for in-memory or non-file DSNs.
func ensureSQLiteDir(dsn string) error {
	const prefix = "file:"
	if !strings.HasPrefix(dsn, prefix) {
		return nil
	}
	rest := strings.TrimPrefix(dsn, prefix)
	if i := strings.IndexByte(rest, '?'); i >= 0 {
		rest = rest[:i]
	}
	if rest == "" || rest == ":memory:" {
		return nil
	}
	return os.MkdirAll(filepath.Dir(rest), 0o755)
}
