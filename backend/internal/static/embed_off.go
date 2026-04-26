//go:build !embed

package static

import "io/fs"

// FrontendFS is empty unless the binary is built with `-tags=embed`.
// At runtime the static handler will report a friendly 404 for non-API
// routes, instructing developers to use the Vite dev server.
var FrontendFS fs.FS
