package config

import (
	"time"
)

// HTTPClientConfig represents the configuration for a named HTTP client.
// This is defined here to avoid circular dependencies.
// The actual client building happens in pkg/httpclient package.
type HTTPClientConfig struct {
	// Name is the unique identifier for this HTTP client
	Name string `json:"name" yaml:"name"`

	// Description provides human-readable documentation for this client
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Authentication configuration
	AuthType string `json:"auth_type,omitempty" yaml:"auth_type,omitempty"` // "none", "basic", "bearer"
	Username string `json:"username,omitempty" yaml:"username,omitempty"`   // For basic auth
	Password string `json:"password,omitempty" yaml:"password,omitempty"`   // For basic auth
	Token    string `json:"token,omitempty" yaml:"token,omitempty"`         // For bearer token

	// Network configuration
	Timeout             time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	MaxIdleConns        int           `json:"max_idle_conns,omitempty" yaml:"max_idle_conns,omitempty"`
	MaxIdleConnsPerHost int           `json:"max_idle_conns_per_host,omitempty" yaml:"max_idle_conns_per_host,omitempty"`
	MaxConnsPerHost     int           `json:"max_conns_per_host,omitempty" yaml:"max_conns_per_host,omitempty"`
	IdleConnTimeout     time.Duration `json:"idle_conn_timeout,omitempty" yaml:"idle_conn_timeout,omitempty"`
	TLSHandshakeTimeout time.Duration `json:"tls_handshake_timeout,omitempty" yaml:"tls_handshake_timeout,omitempty"`
	DisableKeepAlives   bool          `json:"disable_keep_alives,omitempty" yaml:"disable_keep_alives,omitempty"`

	// Security configuration
	MaxRedirects    int   `json:"max_redirects,omitempty" yaml:"max_redirects,omitempty"`
	MaxResponseSize int64 `json:"max_response_size,omitempty" yaml:"max_response_size,omitempty"`
	FollowRedirects bool  `json:"follow_redirects,omitempty" yaml:"follow_redirects,omitempty"`

	// Default headers to include in all requests
	DefaultHeaders map[string]string `json:"default_headers,omitempty" yaml:"default_headers,omitempty"`

	// Default query parameters to include in all requests
	DefaultQueryParams map[string]string `json:"default_query_params,omitempty" yaml:"default_query_params,omitempty"`

	// BaseURL is the base URL for all requests (optional)
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty"`
}

// Config holds workflow engine configuration.
// All configuration options are centralized here for easy management and validation.
type Config struct {
	// Execution limits
	MaxExecutionTime     time.Duration // Maximum time for entire workflow execution
	MaxNodeExecutionTime time.Duration // Maximum time for single node execution
	MaxIterations        int           // Default max iterations for loops (if not specified)

	// HTTP node configuration
	HTTPTimeout         time.Duration // Timeout for HTTP requests
	MaxHTTPRedirects    int           // Maximum number of HTTP redirects to follow
	MaxResponseSize     int64         // Maximum size of HTTP response body (bytes)
	MaxHTTPCallsPerExec int           // Maximum HTTP calls allowed per workflow execution (0 = unlimited)
	AllowedURLPatterns  []string      // Whitelist of allowed URL patterns (if empty, all external URLs allowed)

	// Named HTTP clients configuration
	// SDK consumers can define multiple HTTP clients with different authentication and settings
	HTTPClients []HTTPClientConfig `json:"http_clients,omitempty" yaml:"http_clients,omitempty"`

	// Deprecated field for backward compatibility
	BlockInternalIPs bool // DEPRECATED: Use BlockPrivateIPs instead

	// Zero Trust Security - Network Access Control
	AllowHTTP       bool     // Explicitly allow HTTP requests (default: false for zero trust)
	AllowedDomains  []string // Whitelist of allowed domains for HTTP (empty = allow all domains when AllowHTTP is true)
	BlockPrivateIPs bool     // Block private IP ranges (10.x, 172.16.x, 192.168.x)
	BlockLocalhost  bool     // Block localhost and loopback addresses
	BlockLinkLocal  bool     // Block link-local addresses (169.254.x.x)
	BlockCloudMetadata bool  // Block cloud metadata endpoints (169.254.169.254, etc.)

	// Cache configuration
	DefaultCacheTTL time.Duration // Default TTL for cache entries if not specified
	MaxCacheSize    int           // Maximum number of cache entries (LRU eviction)

	// Resource limits
	MaxInputSize      int // Maximum size of input data (bytes)
	MaxPayloadSize    int // Maximum size of workflow payload (bytes)
	MaxNodes          int // Maximum number of nodes in workflow
	MaxEdges          int // Maximum number of edges in workflow
	MaxNodeExecutions int // Maximum total node executions (including loop iterations, 0 = unlimited)
	MaxStringLength   int // Maximum length of string values (0 = unlimited)
	MaxArrayLength    int // Maximum length of array values (0 = unlimited)
	MaxVariables      int // Maximum number of variables in workflow state (0 = unlimited)
	MaxContextDepth   int // Maximum depth of nested objects/arrays (0 = unlimited)

	// Retry configuration
	DefaultMaxAttempts int           // Default max retry attempts
	DefaultBackoff     time.Duration // Default initial backoff delay
}

// Default returns a Config with secure, production-ready default values.
func Default() *Config {
return &Config{
// Execution limits
MaxExecutionTime:     5 * time.Minute,
MaxNodeExecutionTime: 30 * time.Second,
MaxIterations:        10000,

// HTTP configuration
HTTPTimeout:         30 * time.Second,
MaxHTTPRedirects:    10,
MaxResponseSize:     10 * 1024 * 1024, // 10MB
MaxHTTPCallsPerExec: 0,                // unlimited
AllowedURLPatterns:  nil,              // allow all when AllowHTTP is true
BlockInternalIPs:    true,             // DEPRECATED: kept for backward compatibility

// Zero Trust Security
AllowHTTP:          false, // Require HTTPS by default
AllowedDomains:     nil,   // allow all domains when AllowHTTP is true
BlockPrivateIPs:    true,  // Block private IPs by default
BlockLocalhost:     true,  // Block localhost by default
BlockLinkLocal:     true,  // Block link-local by default
BlockCloudMetadata: true,  // Block cloud metadata by default

// Cache configuration
DefaultCacheTTL: 1 * time.Hour,
MaxCacheSize:    1000,

// Resource limits
MaxInputSize:      1024 * 1024,      // 1MB
MaxPayloadSize:    10 * 1024 * 1024, // 10MB
MaxNodes:          1000,
MaxEdges:          5000,
MaxNodeExecutions: 0, // unlimited
MaxStringLength:   0, // unlimited
MaxArrayLength:    0, // unlimited
MaxVariables:      0, // unlimited
MaxContextDepth:   0, // unlimited

// Retry configuration
DefaultMaxAttempts: 3,
DefaultBackoff:     1 * time.Second,
}
}

// Development returns a Config optimized for development with relaxed limits.
func Development() *Config {
cfg := Default()
cfg.AllowHTTP = true           // Allow HTTP in development
cfg.BlockPrivateIPs = false    // Allow private IPs
cfg.BlockInternalIPs = false   // DEPRECATED: kept for backward compatibility
cfg.BlockLocalhost = false     // Allow localhost
cfg.BlockCloudMetadata = false // Allow cloud metadata
cfg.MaxExecutionTime = 10 * time.Minute
return cfg
}

// Production returns a Config optimized for production with strict security.
func Production() *Config {
cfg := Default()
cfg.AllowHTTP = false          // Require HTTPS
cfg.BlockPrivateIPs = true     // Block private IPs
cfg.BlockInternalIPs = true    // DEPRECATED: kept for backward compatibility
cfg.BlockLocalhost = true      // Block localhost
cfg.BlockLinkLocal = true      // Block link-local
cfg.BlockCloudMetadata = true  // Block cloud metadata
cfg.MaxExecutionTime = 5 * time.Minute
return cfg
}

// Testing returns a Config optimized for testing with minimal limits.
func Testing() *Config {
cfg := Default()
cfg.AllowHTTP = true            // Allow HTTP for test servers
cfg.BlockPrivateIPs = false     // Allow private IPs
cfg.BlockInternalIPs = false    // DEPRECATED: kept for backward compatibility
cfg.BlockLocalhost = false      // Allow localhost
cfg.BlockCloudMetadata = false  // Allow cloud metadata
cfg.MaxExecutionTime = 1 * time.Minute
cfg.HTTPTimeout = 5 * time.Second
return cfg
}

// Validate checks if the configuration values are valid.
func (c *Config) Validate() error {
if c.MaxExecutionTime < 0 {
return ErrInvalidExecutionTime
}
if c.MaxNodeExecutionTime < 0 {
return ErrInvalidNodeExecutionTime
}
if c.MaxIterations < 0 {
return ErrInvalidMaxIterations
}
if c.HTTPTimeout < 0 {
return ErrInvalidHTTPTimeout
}
if c.MaxHTTPRedirects < 0 {
return ErrInvalidMaxRedirects
}
if c.MaxResponseSize < 0 {
return ErrInvalidMaxResponseSize
}
if c.DefaultCacheTTL < 0 {
return ErrInvalidCacheTTL
}
if c.MaxCacheSize < 0 {
return ErrInvalidMaxCacheSize
}
if c.DefaultBackoff < 0 {
return ErrInvalidBackoff
}
return nil
}

// Clone creates a deep copy of the configuration.
func (c *Config) Clone() *Config {
	clone := *c
	if c.AllowedURLPatterns != nil {
		clone.AllowedURLPatterns = make([]string, len(c.AllowedURLPatterns))
		copy(clone.AllowedURLPatterns, c.AllowedURLPatterns)
	}
	if c.AllowedDomains != nil {
		clone.AllowedDomains = make([]string, len(c.AllowedDomains))
		copy(clone.AllowedDomains, c.AllowedDomains)
	}
	if c.HTTPClients != nil {
		clone.HTTPClients = make([]HTTPClientConfig, len(c.HTTPClients))
		for i, client := range c.HTTPClients {
			clone.HTTPClients[i] = client
			// Deep copy maps
			if client.DefaultHeaders != nil {
				clone.HTTPClients[i].DefaultHeaders = make(map[string]string, len(client.DefaultHeaders))
				for k, v := range client.DefaultHeaders {
					clone.HTTPClients[i].DefaultHeaders[k] = v
				}
			}
			if client.DefaultQueryParams != nil {
				clone.HTTPClients[i].DefaultQueryParams = make(map[string]string, len(client.DefaultQueryParams))
				for k, v := range client.DefaultQueryParams {
					clone.HTTPClients[i].DefaultQueryParams[k] = v
				}
			}
		}
	}
	return &clone
}
