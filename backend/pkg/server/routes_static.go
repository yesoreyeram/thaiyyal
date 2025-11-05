package server

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/frontend"
)

// handleStaticFiles serves static frontend files from the embedded filesystem
func (s *Server) handleStaticFiles(w http.ResponseWriter, r *http.Request) {
	// Get the embedded filesystem
	staticFS, err := frontend.GetFS()
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

	// Try to open the file
	file, err := staticFS.Open(filePath)
	if err != nil {
		// If file not found, try with .html extension
		if filePath == "workflow" {
			filePath = "workflow.html"
			file, err = staticFS.Open(filePath)
		}
		
		// If still not found, serve index.html for client-side routing
		if err != nil {
			filePath = "index.html"
			file, err = staticFS.Open(filePath)
			if err != nil {
				http.NotFound(w, r)
				return
			}
		}
	}
	defer file.Close()

	// Get file info for proper content type detection
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// If it's a directory, serve index.html from that directory
	if stat.IsDir() {
		indexPath := path.Join(filePath, "index.html")
		file.Close()
		file, err = staticFS.Open(indexPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()
		stat, err = file.Stat()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Serve the file
	// Create a ReadSeeker from the file
	content, err := fs.ReadFile(staticFS, filePath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), strings.NewReader(string(content)))
}
