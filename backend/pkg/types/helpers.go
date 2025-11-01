package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// generateExecutionID creates a unique execution identifier.
// Uses crypto/rand for cryptographically secure random IDs.
// Format: 16 hex characters (8 bytes) for balance between uniqueness and readability.
func GenerateExecutionID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("exec_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// DefaultConfig returns the default engine configuration
func DefaultConfig() Config {
	return Config{
		// Execution limits
		MaxExecutionTime:     5 * time.Minute,
		MaxNodeExecutionTime: 30 * time.Second,
		MaxIterations:        1000,

		// HTTP configuration
		HTTPTimeout:         30 * time.Second,
		MaxHTTPRedirects:    10,
		MaxResponseSize:     10 * 1024 * 1024, // 10MB
		MaxHTTPCallsPerExec: 100,              // Limit to 100 HTTP calls per execution
		AllowedURLPatterns:  []string{},       // Empty = allow all external URLs
		BlockInternalIPs:    true,             // Block internal IPs by default

		// Cache configuration
		DefaultCacheTTL: 1 * time.Hour,
		MaxCacheSize:    1000,

		// Resource limits
		MaxInputSize:      10 * 1024 * 1024, // 10MB
		MaxPayloadSize:    1 * 1024 * 1024,  // 1MB
		MaxNodes:          1000,
		MaxEdges:          10000,
		MaxNodeExecutions: 10000, // Limit total node executions including loop iterations

		// Retry configuration
		DefaultMaxAttempts: 3,
		DefaultBackoff:     1 * time.Second,
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
	c.MaxHTTPCallsPerExec = 10
	c.MaxNodeExecutions = 1000
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
	c.MaxHTTPCallsPerExec = 1000
	c.MaxNodeExecutions = 100000
	return c
}
