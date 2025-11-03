package config

import "errors"

// Sentinel errors for configuration validation
var (
	// Execution time errors
	ErrInvalidExecutionTime     = errors.New("invalid max execution time: must be non-negative")
	ErrInvalidNodeExecutionTime = errors.New("invalid max node execution time: must be non-negative")
	ErrInvalidMaxIterations     = errors.New("invalid max iterations: must be non-negative")

	// HTTP configuration errors
	ErrInvalidHTTPTimeout     = errors.New("invalid HTTP timeout: must be non-negative")
	ErrInvalidMaxRedirects    = errors.New("invalid max redirects: must be non-negative")
	ErrInvalidMaxResponseSize = errors.New("invalid max response size: must be non-negative")
	ErrInvalidURLPattern      = errors.New("invalid URL pattern")
	ErrInvalidDomain          = errors.New("invalid domain")

	// Cache configuration errors
	ErrInvalidCacheTTL     = errors.New("invalid cache TTL: must be non-negative")
	ErrInvalidMaxCacheSize = errors.New("invalid max cache size: must be non-negative")

	// Resource limit errors
	ErrInvalidInputSize    = errors.New("invalid max input size: must be non-negative")
	ErrInvalidPayloadSize  = errors.New("invalid max payload size: must be non-negative")
	ErrInvalidMaxNodes     = errors.New("invalid max nodes: must be non-negative")
	ErrInvalidMaxEdges     = errors.New("invalid max edges: must be non-negative")
	ErrInvalidStringLength = errors.New("invalid max string length: must be non-negative")
	ErrInvalidArrayLength  = errors.New("invalid max array length: must be non-negative")

	// Retry configuration errors
	ErrInvalidMaxAttempts = errors.New("invalid max attempts: must be positive")
	ErrInvalidBackoff     = errors.New("invalid backoff duration: must be non-negative")

	// File loading errors
	ErrConfigFileNotFound = errors.New("configuration file not found")
	ErrInvalidConfigFile  = errors.New("invalid configuration file format")
	ErrConfigParseFailed  = errors.New("failed to parse configuration file")
)
