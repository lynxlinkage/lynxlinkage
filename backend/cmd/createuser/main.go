// Command createuser is the bootstrap tool for adding HR (or other-role)
// users to the database. Run it once after deploying to create the first
// account.
//
// Usage:
//
//	go run ./cmd/createuser -email hr@example.com -role hr
//	# Then enter the password when prompted (it is not echoed).
//
// You can also pass -password directly for non-interactive scripts, but
// be aware the password will appear in shell history and process lists.
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"syscall"

	"github.com/lynxlinkage/lynxlinkage/backend/internal/auth"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/config"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
	"golang.org/x/term"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	emailFlag := flag.String("email", "", "user email (required)")
	roleFlag := flag.String("role", string(domain.RoleHR), "user role (currently only 'hr' is supported)")
	passwordFlag := flag.String("password", "", "password (omit to be prompted on stdin without echo)")
	flag.Parse()

	email := strings.ToLower(strings.TrimSpace(*emailFlag))
	role := strings.TrimSpace(*roleFlag)

	if email == "" {
		return errors.New("-email is required")
	}
	if role != string(domain.RoleHR) {
		return fmt.Errorf("unsupported role %q (only %q is allowed)", role, domain.RoleHR)
	}

	password := *passwordFlag
	if password == "" {
		var err error
		password, err = readPassword()
		if err != nil {
			return err
		}
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx := context.Background()

	db, err := store.Open(ctx, cfg.DatabaseURL, logger)
	if err != nil {
		return err
	}
	defer db.Close()

	users := store.NewUserRepo(db)

	// We don't need a real session manager for hashing, but NewManager
	// requires a secret. Use a throwaway one — only HashPassword is called.
	mgr := auth.NewManager("createuser-stub", auth.DefaultTTL, users, false)
	hash, err := mgr.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	id, err := users.Create(ctx, email, hash, domain.Role(role))
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	fmt.Printf("created user id=%d email=%s role=%s\n", id, email, role)
	return nil
}

func readPassword() (string, error) {
	if !term.IsTerminal(int(syscall.Stdin)) {
		// Non-terminal stdin: read a single line.
		r := bufio.NewReader(os.Stdin)
		s, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimRight(s, "\r\n"), nil
	}
	fmt.Fprint(os.Stderr, "password: ")
	bytePass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	return string(bytePass), nil
}

