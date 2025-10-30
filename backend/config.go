package workflow

import "time"

// Config holds workflow engine configuration
type Config struct {
	// Execution limits
	MaxExecutionTime    time.Duration // Maximum time for entire workflow execution
	MaxNodeExecutionTime time.Duration // Maximum time for single node execution
	MaxIterations       int           // Default max iterations for loops (if not specified)
	
	// HTTP node configuration
	HTTPTimeout         time.Duration // Timeout for HTTP requests
	MaxHTTPRedirects    int           // Maximum number of HTTP redirects to follow
	MaxResponseSize     int64         // Maximum size of HTTP response body (bytes)
	AllowedURLPatterns  []string      // Whitelist of allowed URL patterns (if empty, all external URLs allowed)
	BlockInternalIPs    bool          // Block requests to internal/private IP addresses
	
	// Cache configuration
	DefaultCacheTTL     time.Duration // Default TTL for cache entries if not specified
	MaxCacheSize        int           // Maximum number of cache entries (LRU eviction)
	
	// Resource limits
	MaxInputSize        int           // Maximum size of input data (bytes)
	MaxPayloadSize      int           // Maximum size of workflow payload (bytes)
	MaxNodes            int           // Maximum number of nodes in workflow
	MaxEdges            int           // Maximum number of edges in workflow
	
	// Retry configuration
	DefaultMaxAttempts  int           // Default max retry attempts
	DefaultBackoff      time.Duration // Default initial backoff delay
}

// DefaultConfig returns the default engine configuration
func DefaultConfig() Config {
	return Config{
		// Execution limits
		MaxExecutionTime:     5 * time.Minute,
		MaxNodeExecutionTime: 30 * time.Second,
		MaxIterations:        1000,
		
		// HTTP configuration
		HTTPTimeout:          30 * time.Second,
		MaxHTTPRedirects:     10,
		MaxResponseSize:      10 * 1024 * 1024, // 10MB
		AllowedURLPatterns:   []string{},       // Empty = allow all external URLs
		BlockInternalIPs:     true,             // Block internal IPs by default
		
		// Cache configuration
		DefaultCacheTTL:      1 * time.Hour,
		MaxCacheSize:         1000,
		
		// Resource limits
		MaxInputSize:         10 * 1024 * 1024, // 10MB
		MaxPayloadSize:       1 * 1024 * 1024,  // 1MB
		MaxNodes:             1000,
		MaxEdges:             10000,
		
		// Retry configuration
		DefaultMaxAttempts:   3,
		DefaultBackoff:       1 * time.Second,
	}
}

// ValidationLimits returns limits suitable for strict validation
func ValidationLimits() Config {
	c := DefaultConfig()
	c.MaxExecutionTime = 1 * time.Minute
	c.MaxNodeExecutionTime = 10 * time.Second
	c.MaxIterations = 100
	c.MaxNodes = 100
	c.MaxEdges = 1000
	return c
}

// DevelopmentConfig returns relaxed limits for development/testing
func DevelopmentConfig() Config {
	c := DefaultConfig()
	c.MaxExecutionTime = 30 * time.Minute
	c.MaxNodeExecutionTime = 5 * time.Minute
	c.MaxIterations = 10000
	c.MaxNodes = 10000
	c.MaxEdges = 100000
	return c
}
