package server

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// handleStaticFiles serves static frontend files from the embedded filesystem
// It automatically handles .html extension for routes and falls back to index.html for SPA routing
func (s *Server) handleStaticFiles(w http.ResponseWriter, r *http.Request) {
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
		
		// If still not found, fall back to index.html for SPA routing
		if err != nil {
			content, stat, err = tryServeFile(staticFS, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
		}
	}

	// Serve the content
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), strings.NewReader(string(content)))
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
