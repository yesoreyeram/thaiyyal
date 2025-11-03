package httpclient

import (
	"fmt"
	"net/http"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/security"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Client wraps an HTTP client with its configuration
type Client struct {
	*http.Client
	config *ClientConfig
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *ClientConfig {
	return c.config
}

// Builder creates configured HTTP clients
type Builder struct {
	engineConfig types.Config // Main engine config for security settings
}

// NewBuilder creates a new HTTP client builder
func NewBuilder(engineConfig types.Config) *Builder {
	return &Builder{
		engineConfig: engineConfig,
	}
}

// Build creates an HTTP client from the given configuration
func (b *Builder) Build(config *ClientConfig) (*Client, error) {
	// Apply defaults
	config.ApplyDefaults()

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client config: %w", err)
	}

	// Create transport with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
		TLSHandshakeTimeout: config.TLSHandshakeTimeout,
		DisableKeepAlives:   config.DisableKeepAlives,
	}

	// Create base HTTP client
	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Transport: &authTransport{
			base:   transport,
			config: config,
		},
	}

	// Configure redirect behavior
	if !config.FollowRedirects {
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.MaxRedirects {
				return fmt.Errorf("too many redirects (max %d)", config.MaxRedirects)
			}
			// Validate redirect URL for SSRF protection
			if err := b.validateURL(req.URL.String()); err != nil {
				return fmt.Errorf("redirect URL validation failed: %w", err)
			}
			return nil
		}
	}

	return &Client{
		Client: httpClient,
		config: config,
	}, nil
}

// validateURL validates URLs to prevent SSRF attacks
func (b *Builder) validateURL(url string) error {
	// Build SSRF protection config from workflow engine config
	ssrfConfig := security.SSRFConfig{
		AllowedSchemes:     []string{"http", "https"},
		BlockPrivateIPs:    b.engineConfig.BlockPrivateIPs,
		BlockLocalhost:     b.engineConfig.BlockLocalhost,
		BlockLinkLocal:     b.engineConfig.BlockLinkLocal,
		BlockCloudMetadata: b.engineConfig.BlockCloudMetadata,
		AllowedDomains:     b.engineConfig.AllowedDomains,
		BlockedDomains:     []string{},
	}

	protection := security.NewSSRFProtectionWithConfig(ssrfConfig)
	return protection.ValidateURL(url)
}

// authTransport is an http.RoundTripper that adds authentication headers
type authTransport struct {
	base   http.RoundTripper
	config *ClientConfig
}

// RoundTrip implements http.RoundTripper interface
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Add authentication headers
	switch t.config.AuthType {
	case AuthTypeBasic:
		clonedReq.SetBasicAuth(t.config.Username, t.config.Password)
	case AuthTypeBearer:
		clonedReq.Header.Set("Authorization", "Bearer "+t.config.Token)
	}

	// Add default headers
	for key, value := range t.config.DefaultHeaders {
		// Don't override headers that are already set
		if clonedReq.Header.Get(key) == "" {
			clonedReq.Header.Set(key, value)
		}
	}

	// Add default query parameters
	if len(t.config.DefaultQueryParams) > 0 {
		q := clonedReq.URL.Query()
		for key, value := range t.config.DefaultQueryParams {
			// Don't override query params that are already set
			if !q.Has(key) {
				q.Set(key, value)
			}
		}
		clonedReq.URL.RawQuery = q.Encode()
	}

	// Execute the request
	return t.base.RoundTrip(clonedReq)
}
