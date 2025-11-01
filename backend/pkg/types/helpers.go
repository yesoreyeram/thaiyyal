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
	return c
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
func getValueDepth(value interface{}) int {
	if value == nil {
		return 0
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		maxDepth := 0
		iter := v.MapRange()
		for iter.Next() {
			depth := getValueDepth(iter.Value().Interface())
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		return 1 + maxDepth
	case reflect.Slice, reflect.Array:
		maxDepth := 0
		for i := 0; i < v.Len(); i++ {
			depth := getValueDepth(v.Index(i).Interface())
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		return 1 + maxDepth
	default:
		return 1
	}
}
