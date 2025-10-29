package workflow

import (
	"fmt"
	"io"
	"net/http"
)

// ============================================================================
// HTTP Node Executor
// ============================================================================
// This file contains the executor for HTTP nodes that perform HTTP requests.
// ============================================================================

// executeHTTPNode performs an HTTP GET request and returns the response body.
//
// Required fields:
//   - Data.URL: The URL to send the HTTP GET request to
//
// Returns:
//   - string: Response body content
//   - error: If URL is missing, request fails, or HTTP status indicates error (non-2xx)
//
// Note: Currently supports only GET requests. Future versions may support
// POST, PUT, DELETE with configurable methods and bodies.
func (e *Engine) executeHTTPNode(node Node) (interface{}, error) {
	if node.Data.URL == nil {
		return nil, fmt.Errorf("HTTP node missing url")
	}

	// Make HTTP GET request
	resp, err := http.Get(*node.Data.URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for error status codes (only 2xx considered success)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request returned error status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
