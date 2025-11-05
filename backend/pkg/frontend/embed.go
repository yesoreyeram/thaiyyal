//go:build !dev
// +build !dev

// Package frontend provides embedded static files for the frontend application.
package frontend

import (
	"embed"
	"io/fs"
)

// Static contains the embedded frontend files from the Next.js standalone build
//
//go:embed all:static
var Static embed.FS

// GetFS returns a filesystem for the embedded static files
func GetFS() (fs.FS, error) {
	return fs.Sub(Static, "static")
}
