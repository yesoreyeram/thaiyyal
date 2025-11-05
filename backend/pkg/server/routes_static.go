package server

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// handleStaticFiles serves static frontend files from the embedded filesystem
// or proxies to Next.js dev server in dev mode
// It automatically handles .html extension for routes and falls back to index.html for SPA routing
func (s *Server) handleStaticFiles(w http.ResponseWriter, r *http.Request) {
	// In dev mode, proxy to Next.js dev server
	if isDevMode() {
		s.handleDevProxy(w, r)
		return
	}

	// Get the embedded filesystem
	staticFS, err := getStaticFS()
	if err != nil {
		s.logger.WithError(err).Error("failed to get static filesystem")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Clean the path
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	upath = path.Clean(upath)

	// Remove leading slash for fs operations
	filePath := strings.TrimPrefix(upath, "/")

	// Default to index.html for root path
	if filePath == "" || filePath == "/" {
		filePath = "index.html"
	}

	// Try to serve the file as requested
	content, stat, err := tryServeFile(staticFS, filePath)
	if err != nil {
		// If not found and doesn't have an extension, try with .html
		if !strings.Contains(filePath, ".") {
			htmlPath := filePath + ".html"
			content, stat, err = tryServeFile(staticFS, htmlPath)
		}

		// If still not found, only fall back to index.html for HTML routes (not for static assets)
		// Static assets like .js, .css, images should return 404 if not found
		if err != nil {
			// Check if this is a static asset request (has file extension like .js, .css, .png, etc.)
			hasExtension := strings.Contains(filePath, ".")
			isStaticAsset := hasExtension && (strings.HasPrefix(filePath, "_next/") ||
				strings.HasSuffix(filePath, ".js") ||
				strings.HasSuffix(filePath, ".css") ||
				strings.HasSuffix(filePath, ".png") ||
				strings.HasSuffix(filePath, ".jpg") ||
				strings.HasSuffix(filePath, ".svg") ||
				strings.HasSuffix(filePath, ".ico") ||
				strings.HasSuffix(filePath, ".woff") ||
				strings.HasSuffix(filePath, ".woff2"))

			// Only fall back to index.html for potential SPA routes (not static assets)
			if !isStaticAsset {
				content, stat, err = tryServeFile(staticFS, "index.html")
			}

			if err != nil {
				http.NotFound(w, r)
				return
			}
		}
	}

	// Set appropriate content type
	contentType := getContentType(filePath)
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// Serve the content
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), strings.NewReader(string(content)))
}

// handleDevProxy proxies requests to Next.js dev server in development mode
func (s *Server) handleDevProxy(w http.ResponseWriter, r *http.Request) {
	proxy, err := getDevProxy()
	if err != nil {
		s.logger.WithError(err).Error("failed to create dev proxy")
		http.Error(w, "Failed to proxy to dev server", http.StatusInternalServerError)
		return
	}

	s.logger.WithFields(map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Debug("proxying request to Next.js dev server")

	proxy.ServeHTTP(w, r)
}

// tryServeFile attempts to read a file from the filesystem and returns its content and stat
func tryServeFile(fsys fs.FS, filePath string) ([]byte, fs.FileInfo, error) {
	file, err := fsys.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	// If it's a directory, try index.html in that directory
	if stat.IsDir() {
		indexPath := path.Join(filePath, "index.html")
		return tryServeFile(fsys, indexPath)
	}

	// Read the file content
	content, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		return nil, nil, err
	}

	return content, stat, nil
}

// getContentType returns the appropriate Content-Type header for a file
func getContentType(filePath string) string {
	// Get file extension
	ext := path.Ext(filePath)

	// Map common extensions to content types
	contentTypes := map[string]string{
		".html":  "text/html; charset=utf-8",
		".css":   "text/css; charset=utf-8",
		".js":    "application/javascript; charset=utf-8",
		".json":  "application/json; charset=utf-8",
		".png":   "image/png",
		".jpg":   "image/jpeg",
		".jpeg":  "image/jpeg",
		".gif":   "image/gif",
		".svg":   "image/svg+xml",
		".ico":   "image/x-icon",
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".ttf":   "font/ttf",
		".eot":   "application/vnd.ms-fontobject",
	}

	if contentType, ok := contentTypes[ext]; ok {
		return contentType
	}

	return ""
}
