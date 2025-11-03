package httpclient

import (
	"context"
	"fmt"
	"net/http"
)

// New creates a new HTTP client from the given configuration.
// This is the main entry point for creating HTTP clients.
//
// The context parameter is currently unused but included for future extensibility
// (e.g., for context-based timeout configuration or tracing).
func New(ctx context.Context, config *Config) (*http.Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Apply defaults
	config.ApplyDefaults()

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create base transport with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
		TLSHandshakeTimeout: config.TLSHandshakeTimeout,
		DisableKeepAlives:   config.DisableKeepAlives,
	}

	// Build middleware chain
	var middlewares []Middleware

	// Add SSRF protection middleware if any protection is enabled
	if config.BlockPrivateIPs || config.BlockLocalhost || config.BlockLinkLocal ||
		config.BlockCloudMetadata || len(config.AllowedDomains) > 0 {
		middlewares = append(middlewares, ssrfProtectionMiddleware(config))
	}

	// Add query params middleware if configured
	if len(config.QueryParams) > 0 {
		middlewares = append(middlewares, queryParamsMiddleware(config.QueryParams))
	}

	// Add headers middleware if configured
	if len(config.Headers) > 0 {
		middlewares = append(middlewares, headersMiddleware(config.Headers))
	}

	// Add authentication middleware if configured
	if config.AuthType != AuthTypeNone {
		middlewares = append(middlewares, authMiddleware(config))
	}

	// Apply middleware chain to transport
	var finalTransport http.RoundTripper = transport
	if len(middlewares) > 0 {
		chainedMiddleware := Chain(middlewares...)
		finalTransport = chainedMiddleware(transport)
	}

	// Create HTTP client
	client := &http.Client{
		Timeout:   config.Timeout,
		Transport: finalTransport,
	}

	// Configure redirect behavior
	if !config.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.MaxRedirects {
				return fmt.Errorf("too many redirects (max %d)", config.MaxRedirects)
			}
			// Validate redirect URL for SSRF protection
			if err := validateURL(req.URL.String(), config); err != nil {
				return fmt.Errorf("redirect URL validation failed: %w", err)
			}
			return nil
		}
	}

	return client, nil
}
