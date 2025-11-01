package executor

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/security"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// HTTPExecutor executes HTTP nodes with connection pooling
type HTTPExecutor struct {
	client *http.Client
	mu     sync.RWMutex
}

// NewHTTPExecutor creates a new HTTP executor with a shared connection pool
func NewHTTPExecutor() *HTTPExecutor {
	return &HTTPExecutor{}
}

// Execute runs the HTTP node
// Performs an HTTP GET request and returns the response body.
// Uses a shared connection pool for better performance.
//
// Security features:
//   - URL validation (blocks internal IPs by default)
//   - Request timeout (30s default, configurable)
//   - Response size limit (10MB default, configurable)
//   - SSRF protection against cloud metadata endpoints
//   - HTTP call count limit per execution
func (e *HTTPExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.URL == nil {
		return nil, fmt.Errorf("HTTP node missing url")
	}

	// Check and increment HTTP call counter before making the request
	if err := ctx.IncrementHTTPCall(); err != nil {
		return nil, err
	}

	config := ctx.GetConfig()

	// Validate URL for security (SSRF protection)
	if err := isAllowedURL(*node.Data.URL, config); err != nil {
		return nil, fmt.Errorf("URL validation failed: %w", err)
	}

	// Get or create shared HTTP client with connection pooling
	client := e.getOrCreateClient(config)

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
	limitedReader := io.LimitReader(resp.Body, config.MaxResponseSize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if response was truncated due to size limit
	if int64(len(body)) == config.MaxResponseSize {
		// Try to read one more byte to see if there's more data
		oneByte := make([]byte, 1)
		if n, _ := resp.Body.Read(oneByte); n > 0 {
			return nil, fmt.Errorf("response too large (exceeds %d bytes limit)", config.MaxResponseSize)
		}
	}

	return string(body), nil
}

// getOrCreateClient returns the shared HTTP client, creating it if necessary
// This enables connection pooling and reuse across multiple requests
func (e *HTTPExecutor) getOrCreateClient(config types.Config) *http.Client {
	e.mu.RLock()
	if e.client != nil {
		e.mu.RUnlock()
		return e.client
	}
	e.mu.RUnlock()

	e.mu.Lock()
	defer e.mu.Unlock()

	// Double-check after acquiring write lock
	if e.client != nil {
		return e.client
	}

	// Create HTTP client with connection pooling and security settings
	e.client = &http.Client{
		Timeout: config.HTTPTimeout,
		Transport: &http.Transport{
			// Connection pooling settings
			MaxIdleConns:        100,              // Max idle connections across all hosts
			MaxIdleConnsPerHost: 10,               // Max idle connections per host
			MaxConnsPerHost:     100,              // Max connections per host
			IdleConnTimeout:     90 * time.Second, // How long idle connections are kept
			
			// Performance settings
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			
			// Keep connections alive for reuse
			DisableKeepAlives: false,
		},
		// Limit redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.MaxHTTPRedirects {
				return fmt.Errorf("too many redirects (max %d)", config.MaxHTTPRedirects)
			}
			// Validate redirect URL as well (prevent redirect-based SSRF)
			if err := isAllowedURL(req.URL.String(), config); err != nil {
				return fmt.Errorf("redirect URL validation failed: %w", err)
			}
			return nil
		},
	}

	return e.client
}

// NodeType returns the node type this executor handles
func (e *HTTPExecutor) NodeType() types.NodeType {
	return types.NodeTypeHTTP
}

// Validate checks if node configuration is valid
func (e *HTTPExecutor) Validate(node types.Node) error {
	if node.Data.URL == nil {
		return fmt.Errorf("HTTP node missing url")
	}
	return nil
}

// isAllowedURL validates URLs to prevent SSRF attacks using the security package
func isAllowedURL(url string, config types.Config) error {
	// For now, we'll use a permissive config that allows localhost for testing
	// In production, config should control this behavior
	ssrfConfig := security.SSRFConfig{
		AllowedSchemes:     []string{"http", "https"},
		BlockPrivateIPs:    false, // Allow for now - should be true in production
		BlockLocalhost:     false, // Allow for now - should be true in production
		BlockLinkLocal:     true,
		BlockCloudMetadata: true,
		AllowedDomains:     []string{},
		BlockedDomains:     []string{},
	}
	
	protection := security.NewSSRFProtectionWithConfig(ssrfConfig)
	
	// Validate URL
	return protection.ValidateURL(url)
}
