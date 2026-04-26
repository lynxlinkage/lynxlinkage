// Package uploads is a thin disk-backed object store for candidate
// artifacts (CVs, cover letters, samples). Files are stored under
//
//	<root>/applications/<applicationID>/<random>-<safe-original-name>
//
// where <random> is a hex token to defeat collisions and prevent
// guessing, and <safe-original-name> is a sanitised version of the
// uploaded filename (path separators stripped, length-capped). The
// original (unsanitised) name is stored in the DB row for display and
// Content-Disposition headers; the on-disk name is purely an internal
// identifier.
package uploads

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Store writes and reads candidate artifacts on the local filesystem.
// All paths returned by Store methods are *relative* to Root, so the
// values are stable even if Root changes between deploys.
type Store struct {
	Root string
}

// NewStore creates a Store rooted at root. The directory is created
// (with parents) if it doesn't exist.
func NewStore(root string) (*Store, error) {
	if root == "" {
		return nil, errors.New("uploads: empty root")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("uploads: mkdir root: %w", err)
	}
	return &Store{Root: root}, nil
}

// Save writes r to disk for the given application and returns the
// relative storage path plus the number of bytes copied. limit caps the
// number of bytes that will be read from r; Save returns ErrFileTooLarge
// if the source produces more.
func (s *Store) Save(applicationID int64, originalName string, r io.Reader, limit int64) (string, int64, error) {
	if limit <= 0 {
		return "", 0, errors.New("uploads: non-positive limit")
	}

	dir := filepath.Join("applications", fmt.Sprintf("%d", applicationID))
	absDir := filepath.Join(s.Root, dir)
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		return "", 0, fmt.Errorf("uploads: mkdir %s: %w", absDir, err)
	}

	token, err := randomToken(8)
	if err != nil {
		return "", 0, err
	}
	safeName := safeFilename(originalName)
	rel := filepath.Join(dir, token+"-"+safeName)
	abs := filepath.Join(s.Root, rel)

	f, err := os.OpenFile(abs, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", 0, fmt.Errorf("uploads: open: %w", err)
	}

	// Read up to limit+1 to detect oversized inputs without buffering
	// the whole stream first.
	written, copyErr := io.Copy(f, io.LimitReader(r, limit+1))
	if cerr := f.Close(); cerr != nil && copyErr == nil {
		copyErr = cerr
	}
	if copyErr != nil {
		_ = os.Remove(abs)
		return "", 0, fmt.Errorf("uploads: copy: %w", copyErr)
	}
	if written > limit {
		_ = os.Remove(abs)
		return "", 0, ErrFileTooLarge
	}
	return filepath.ToSlash(rel), written, nil
}

// Open returns a read-only handle to a previously-stored file. The
// returned ReadSeekCloser exposes the size via Stat() too if the caller
// needs it.
func (s *Store) Open(rel string) (*os.File, error) {
	if !isSafeRelative(rel) {
		return nil, fmt.Errorf("uploads: refusing unsafe path %q", rel)
	}
	abs := filepath.Join(s.Root, filepath.FromSlash(rel))
	return os.Open(abs)
}

// Remove deletes a previously-stored file. Safe to call with a path
// that no longer exists; the error is logged at the caller's discretion.
func (s *Store) Remove(rel string) error {
	if !isSafeRelative(rel) {
		return fmt.Errorf("uploads: refusing unsafe path %q", rel)
	}
	abs := filepath.Join(s.Root, filepath.FromSlash(rel))
	return os.Remove(abs)
}

// ErrFileTooLarge is returned by Save when the source produces more
// bytes than the configured limit.
var ErrFileTooLarge = errors.New("uploads: file too large")

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// safeFilename strips directory components and keeps only a small
// allow-list of characters. The result is always non-empty.
func safeFilename(in string) string {
	in = filepath.Base(in)
	in = strings.ReplaceAll(in, "\x00", "")
	const maxLen = 80

	var b strings.Builder
	for _, r := range in {
		switch {
		case r == ' ' || r == '_' || r == '-' || r == '.' || r == '(' || r == ')':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
		if b.Len() >= maxLen {
			break
		}
	}
	out := strings.TrimSpace(b.String())
	out = strings.TrimLeft(out, ".") // no leading dots — avoid hidden files
	if out == "" {
		out = "file"
	}
	return out
}

// isSafeRelative rejects absolute paths and any traversal segments. We
// constrain "stored_path" values written by us to be safe by
// construction, but defence-in-depth in case the DB is tampered with.
func isSafeRelative(p string) bool {
	if p == "" || strings.HasPrefix(p, "/") || strings.HasPrefix(p, `\`) {
		return false
	}
	parts := strings.Split(filepath.ToSlash(p), "/")
	for _, seg := range parts {
		if seg == "" || seg == "." || seg == ".." {
			return false
		}
	}
	return true
}
