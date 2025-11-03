package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
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

// DefaultConfig returns the default engine configuration with zero-trust security model.
// 
// **ZERO TRUST BY DEFAULT**: Network access and privileged operations are DISABLED by default.
// You must explicitly enable them in your configuration.
//
// To enable HTTP requests:
//   config := types.DefaultConfig()
//   config.AllowHTTP = true
//   config.AllowedDomains = []string{"api.example.com", "trusted-domain.com"}
//
// Zero Trust Principles:
//   - Network access disabled by default (AllowHTTP = false)
//   - No environment variable access (not implemented in workflow engine)
//   - No file system access (not implemented in workflow engine)
//   - All security protections enabled
//   - Reasonable resource limits for safe execution
//
// For development/testing, use DevelopmentConfig() which has relaxed settings.
func DefaultConfig() Config {
	return Config{
		// Execution limits - Reasonable for production
		MaxExecutionTime:     5 * time.Minute,
		MaxNodeExecutionTime: 30 * time.Second,
		MaxIterations:        1000,

		// HTTP configuration - DISABLED by default (zero trust)
		HTTPTimeout:         30 * time.Second,
		MaxHTTPRedirects:    10,
		MaxResponseSize:     10 * 1024 * 1024, // 10MB
		MaxHTTPCallsPerExec: 100,              // Limit when HTTP is enabled
		AllowedURLPatterns:  []string{},       // Empty = allow all when HTTP enabled
		BlockInternalIPs:    true,             // Deprecated: use BlockPrivateIPs
		
		// Zero Trust Security - DENY ALL by default (explicit opt-in required)
		AllowHTTP:          false,   // HTTP DISABLED - must explicitly enable
		AllowedDomains:     []string{}, // No domains whitelisted (must configure if enabling HTTP)
		BlockPrivateIPs:    true,    // Block all private IP ranges
		BlockLocalhost:     true,    // Block localhost
		BlockLinkLocal:     true,    // Block link-local
		BlockCloudMetadata: true,    // Block cloud metadata

		// Cache configuration
		DefaultCacheTTL: 1 * time.Hour,
		MaxCacheSize:    1000,

		// Resource limits - Reasonable for production
		MaxInputSize:      10 * 1024 * 1024, // 10MB
		MaxPayloadSize:    1 * 1024 * 1024,  // 1MB
		MaxNodes:          1000,
		MaxEdges:          10000,
		MaxNodeExecutions: 10000,            // Limit total node executions including loop iterations
		MaxStringLength:   1024 * 1024,      // 1MB max string length
		MaxArrayLength:    10000,            // 10k elements max in arrays
		MaxVariables:      1000,             // Max 1000 variables in workflow state
		MaxContextDepth:   32,               // Max 32 levels of nesting

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
	c.MaxStringLength = 100 * 1024  // 100KB
	c.MaxArrayLength = 1000
	c.MaxVariables = 100
	c.MaxContextDepth = 16
	// Keep HTTP enabled but with strict limits for validation
	c.AllowHTTP = true
	c.BlockLocalhost = true // Block localhost in validation mode
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
	c.MaxStringLength = 10 * 1024 * 1024 // 10MB
	c.MaxArrayLength = 100000
	c.MaxVariables = 10000
	c.MaxContextDepth = 64
	// Allow localhost for development
	c.AllowHTTP = true
	c.BlockLocalhost = false
	c.BlockPrivateIPs = false // Allow private IPs in development
	return c
}

// ZeroTrustConfig returns ultra-restrictive configuration following zero-trust principles
// This configuration denies all network access and applies minimal resource limits by default.
// Suitable for running untrusted workflows in sandboxed environments.
//
// Zero Trust Principles:
//   - No network access by default (AllowHTTP = false)
//   - No environment variable access (not implemented in workflow engine)
//   - No file system access (not implemented in workflow engine)
//   - Minimal execution time and resource limits
//   - All security protections enabled at maximum level
//
// To enable specific capabilities, explicitly configure them:
//   config := types.ZeroTrustConfig()
//   config.AllowHTTP = true  // Enable HTTP
//   config.AllowedDomains = []string{"api.example.com"}  // Whitelist specific domains
func ZeroTrustConfig() Config {
	return Config{
		// Execution limits - Minimal
		MaxExecutionTime:     30 * time.Second, // Short execution time
		MaxNodeExecutionTime: 5 * time.Second,  // Short per-node time
		MaxIterations:        50,               // Minimal iterations

		// HTTP configuration - DISABLED by default
		HTTPTimeout:         10 * time.Second, // Short timeout
		MaxHTTPRedirects:    0,                // No redirects
		MaxResponseSize:     1 * 1024 * 1024,  // 1MB response limit
		MaxHTTPCallsPerExec: 0,                // 0 = No HTTP calls allowed
		AllowedURLPatterns:  []string{},       // Empty whitelist
		BlockInternalIPs:    true,             // Block internal IPs
		
		// Zero Trust Security - DENY ALL by default
		AllowHTTP:          false,   // Deny all HTTP access
		AllowedDomains:     []string{}, // No domains allowed (deny all)
		BlockPrivateIPs:    true,    // Block all private IP ranges
		BlockLocalhost:     true,    // Block localhost
		BlockLinkLocal:     true,    // Block link-local
		BlockCloudMetadata: true,    // Block cloud metadata

		// Cache configuration - Minimal
		DefaultCacheTTL: 5 * time.Minute, // Short TTL
		MaxCacheSize:    100,             // Small cache

		// Resource limits - Minimal
		MaxInputSize:      1 * 1024 * 1024,   // 1MB input
		MaxPayloadSize:    512 * 1024,        // 512KB payload
		MaxNodes:          50,                // Few nodes
		MaxEdges:          200,               // Few edges
		MaxNodeExecutions: 500,               // Limited executions
		MaxStringLength:   50 * 1024,         // 50KB strings
		MaxArrayLength:    500,               // Small arrays
		MaxVariables:      50,                // Few variables
		MaxContextDepth:   10,                // Shallow nesting

		// Retry configuration - Minimal
		DefaultMaxAttempts: 1, // No retries
		DefaultBackoff:     0, // No backoff
	}
}

// ValidateValue validates a value against resource limits in the config.
// Returns an error if the value violates any limits.
func ValidateValue(value interface{}, config Config) error {
	if value == nil {
		return nil
	}

	// Check string length
	if config.MaxStringLength > 0 {
		if str, ok := value.(string); ok {
			if len(str) > config.MaxStringLength {
				return fmt.Errorf("string too long: %d bytes (limit: %d)", len(str), config.MaxStringLength)
			}
		}
	}

	// Check array length
	if config.MaxArrayLength > 0 {
		if arr, ok := value.([]interface{}); ok {
			if len(arr) > config.MaxArrayLength {
				return fmt.Errorf("array too large: %d elements (limit: %d)", len(arr), config.MaxArrayLength)
			}
			// Recursively validate array elements
			for i, elem := range arr {
				if err := ValidateValue(elem, config); err != nil {
					return fmt.Errorf("array element %d: %w", i, err)
				}
			}
		}
	}

	// Check nesting depth
	if config.MaxContextDepth > 0 {
		depth := getValueDepth(value)
		if depth > config.MaxContextDepth {
			return fmt.Errorf("value too deeply nested: %d levels (limit: %d)", depth, config.MaxContextDepth)
		}
	}

	return nil
}

// getValueDepth calculates the nesting depth of a value
// with protection against stack overflow
func getValueDepth(value interface{}) int {
	return getValueDepthRecursive(value, 0, 1000) // max 1000 depth for safety
}

func getValueDepthRecursive(value interface{}, currentDepth, maxDepth int) int {
	if value == nil || currentDepth >= maxDepth {
		return currentDepth
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		maxChildDepth := currentDepth
		iter := v.MapRange()
		for iter.Next() {
			depth := getValueDepthRecursive(iter.Value().Interface(), currentDepth+1, maxDepth)
			if depth > maxChildDepth {
				maxChildDepth = depth
			}
		}
		return maxChildDepth
	case reflect.Slice, reflect.Array:
		maxChildDepth := currentDepth
		for i := 0; i < v.Len(); i++ {
			depth := getValueDepthRecursive(v.Index(i).Interface(), currentDepth+1, maxDepth)
			if depth > maxChildDepth {
				maxChildDepth = depth
			}
		}
		return maxChildDepth
	default:
		return currentDepth + 1
	}
}
