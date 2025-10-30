package workflow

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// ============================================================================
// HTTP Node Executor
// ============================================================================
// This file contains the executor for HTTP nodes that perform HTTP requests.
// Includes security features: SSRF protection, timeouts, response size limits.
// ============================================================================

// executeHTTPNode performs an HTTP GET request and returns the response body.
//
// Security features:
//   - URL validation (blocks internal IPs by default)
//   - Request timeout (30s default, configurable)
//   - Response size limit (10MB default, configurable)
//   - SSRF protection against cloud metadata endpoints
//
// Required fields:
//   - Data.URL: The URL to send the HTTP GET request to
//
// Returns:
//   - string: Response body content
//   - error: If URL is missing, validation fails, request fails, or HTTP status indicates error (non-2xx)
//
// Note: Currently supports only GET requests. Future versions may support
// POST, PUT, DELETE with configurable methods and bodies.
func (e *Engine) executeHTTPNode(node Node) (interface{}, error) {
	if node.Data.URL == nil {
		return nil, fmt.Errorf("HTTP node missing url")
	}

	// Validate URL for security (SSRF protection)
	if err := isAllowedURL(*node.Data.URL, e.config); err != nil {
		return nil, fmt.Errorf("URL validation failed: %w", err)
	}

	// Create HTTP client with timeout and security settings
	client := &http.Client{
		Timeout: e.config.HTTPTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			DisableKeepAlives:   false,
		},
		// Limit redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= e.config.MaxHTTPRedirects {
				return fmt.Errorf("too many redirects (max %d)", e.config.MaxHTTPRedirects)
			}
			// Validate redirect URL as well (prevent redirect-based SSRF)
			if err := isAllowedURL(req.URL.String(), e.config); err != nil {
				return fmt.Errorf("redirect URL validation failed: %w", err)
			}
			return nil
		},
	}

	// Make HTTP GET request
	resp, err := client.Get(*node.Data.URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for error status codes (only 2xx considered success)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request returned error status: %d", resp.StatusCode)
	}

	// Read response body with size limit
	limitedReader := io.LimitReader(resp.Body, e.config.MaxResponseSize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if response was truncated due to size limit
	if int64(len(body)) == e.config.MaxResponseSize {
		// Try to read one more byte to see if there's more data
		oneByte := make([]byte, 1)
		if n, _ := resp.Body.Read(oneByte); n > 0 {
			return nil, fmt.Errorf("response too large (exceeds %d bytes limit)", e.config.MaxResponseSize)
		}
	}

	return string(body), nil
}
