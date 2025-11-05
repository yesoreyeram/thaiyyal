//go:build !dev
// +build !dev

package server

import (
	"embed"
	"io/fs"
	"net/http/httputil"
)

// staticFiles contains the embedded frontend files from the Next.js build
// These files should be built and copied before building the Go binary
//
//go:embed static
//go:embed static/_next
//go:embed static/_not-found
var staticFiles embed.FS

// getStaticFS returns a filesystem for the embedded static files
func getStaticFS() (fs.FS, error) {
	return fs.Sub(staticFiles, "static")
}

// getDevProxy returns nil in production mode
func getDevProxy() (*httputil.ReverseProxy, error) {
	return nil, nil
}

// isDevMode returns false in production builds
func isDevMode() bool {
	return false
}
