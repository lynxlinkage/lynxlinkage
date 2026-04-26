//go:build embed

package static

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

// FrontendFS is the prerendered SvelteKit output, rooted at the dist/
// directory shipped alongside this package. The Makefile copies
// frontend/build into backend/internal/static/dist before running
// `go build -tags=embed`.
var FrontendFS fs.FS = mustSub(distFS, "dist")

func mustSub(f embed.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
