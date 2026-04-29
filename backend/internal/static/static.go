// Package static serves the prerendered SvelteKit build. In production the
// SvelteKit static export is embedded into the binary at build time (see
// `embed` build tag below). In development the embed.FS is empty and the
// frontend is served by the Vite dev server on its own port; the Go server
// then returns 404 for non-API routes, which is fine because the user
// hits Vite directly.
//
// To produce a production binary that contains the frontend run:
//
//	cd frontend && pnpm build
//	cd backend  && go build -tags=embed ./cmd/server
package static

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler returns a Gin handler that serves files from fsys with SPA-style
// fallback: any request that doesn't match a real file falls back to
// index.html so client-side routes still resolve. /api/* requests are
// expected to be handled before this catch-all.
//
// When fsys is empty (development build without -tags=embed), the handler
// responds with 404 and a hint to use the Vite dev server.
func Handler(fsys fs.FS) gin.HandlerFunc {
	if fsys == nil {
		return devFallback
	}
	if entries, err := fs.ReadDir(fsys, "."); err != nil || len(entries) == 0 {
		return devFallback
	}

	fileServer := http.FileServer(http.FS(fsys))

	return func(c *gin.Context) {
		req := c.Request
		urlPath := strings.TrimPrefix(req.URL.Path, "/")

		// Serve root as "/" so FileServer resolves index.html internally.
		// Rewriting to "/index.html" triggers a redirect loop because
		// FileServer always redirects /index.html → ./.
		if urlPath == "" {
			c.Header("Cache-Control", "no-cache")
			req2 := req.Clone(req.Context())
			req2.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, req2)
			return
		}

		if !fileExists(fsys, urlPath) {
			// Try /<path>.html for SvelteKit's prerendered routes.
			switch {
			case fileExists(fsys, urlPath+".html"):
				urlPath = urlPath + ".html"
			case fileExists(fsys, path.Join(urlPath, "index.html")):
				urlPath = path.Join(urlPath, "index.html")
			case fileExists(fsys, "200.html"):
				// SPA fallback emitted by adapter-static for non-prerendered
				// routes such as /admin and /login.
				urlPath = "200.html"
			default:
				urlPath = "index.html"
			}
		}

		// Cache hashed assets aggressively, keep HTML uncached.
		if strings.HasPrefix(urlPath, "_app/") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		} else if strings.HasSuffix(urlPath, ".html") {
			c.Header("Cache-Control", "no-cache")
		}

		req2 := req.Clone(req.Context())
		req2.URL.Path = "/" + urlPath
		fileServer.ServeHTTP(c.Writer, req2)
	}
}

func fileExists(fsys fs.FS, name string) bool {
	f, err := fsys.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return false
	}
	return !st.IsDir()
}

func devFallback(c *gin.Context) {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusNotFound,
		"frontend not embedded; run `pnpm dev` in frontend/ for development "+
			"or build with `-tags=embed` for production.\n")
}
