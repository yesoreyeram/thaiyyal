package config

import (
	"time"
)

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

	// Deprecated field for backward compatibility
	BlockInternalIPs bool // DEPRECATED: Use BlockPrivateIPs instead

	// Zero Trust Security - Network Access Control
	// ALL NETWORK ACCESS IS DENIED BY DEFAULT (zero trust)
	// Use Allow* fields to explicitly permit access
	AllowHTTP          bool     // Explicitly allow HTTP requests (default: false for zero trust)
	AllowedDomains     []string // Whitelist of allowed domains for HTTP (empty = allow all domains when AllowHTTP is true)
	AllowPrivateIPs    bool     // Allow private IP ranges (10.x, 172.16.x, 192.168.x) - default: false (BLOCKED)
	AllowLocalhost     bool     // Allow localhost and loopback addresses - default: false (BLOCKED)
	AllowLinkLocal     bool     // Allow link-local addresses (169.254.x.x) - default: false (BLOCKED)
	AllowCloudMetadata bool     // Allow cloud metadata endpoints (169.254.169.254, etc.) - default: false (BLOCKED)

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
		MaxHTTPCallsPerExec: 100,              // Default: 100 calls per execution (changed from unlimited for security)
		AllowedURLPatterns:  nil,              // allow all when AllowHTTP is true
		BlockInternalIPs:    true,             // DEPRECATED: kept for backward compatibility

		// Zero Trust Security - DENY BY DEFAULT
		AllowHTTP:          false, // Require HTTPS by default
		AllowedDomains:     nil,   // allow all domains when AllowHTTP is true
		AllowPrivateIPs:    false, // Block private IPs by default (DENY)
		AllowLocalhost:     false, // Block localhost by default (DENY)
		AllowLinkLocal:     false, // Block link-local by default (DENY)
		AllowCloudMetadata: false, // Block cloud metadata by default (DENY)

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
	cfg.AllowPrivateIPs = true     // Allow private IPs
	cfg.BlockInternalIPs = false   // DEPRECATED: kept for backward compatibility
	cfg.AllowLocalhost = true      // Allow localhost
	cfg.AllowCloudMetadata = false // Still block cloud metadata (security best practice)
	cfg.MaxExecutionTime = 10 * time.Minute
	return cfg
}

// Production returns a Config optimized for production with strict security.
func Production() *Config {
	cfg := Default()
	cfg.AllowHTTP = false         // Require HTTPS
	cfg.AllowPrivateIPs = false   // Block private IPs (DENY)
	cfg.BlockInternalIPs = true   // DEPRECATED: kept for backward compatibility
	cfg.AllowLocalhost = false    // Block localhost (DENY)
	cfg.AllowLinkLocal = false    // Block link-local (DENY)
	cfg.AllowCloudMetadata = false // Block cloud metadata (DENY)
	cfg.MaxExecutionTime = 5 * time.Minute
	return cfg
}

// Testing returns a Config optimized for testing with minimal limits.
func Testing() *Config {
	cfg := Default()
	cfg.AllowHTTP = true           // Allow HTTP for test servers
	cfg.AllowPrivateIPs = true     // Allow private IPs
	cfg.BlockInternalIPs = false   // DEPRECATED: kept for backward compatibility
	cfg.AllowLocalhost = true      // Allow localhost
	cfg.AllowCloudMetadata = false // Still block cloud metadata (security best practice)
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
	return &clone
}
