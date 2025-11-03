package security

import "errors"

// Sentinel errors for security operations
var (
	// Validation errors
	ErrValidationFailed = errors.New("security validation failed")
	ErrInvalidInput     = errors.New("invalid input detected")
	ErrInputTooLarge    = errors.New("input exceeds maximum size")
	ErrInvalidFormat    = errors.New("invalid input format")

	// URL security errors
	ErrURLNotAllowed    = errors.New("URL not allowed by security policy")
	ErrPrivateIPBlocked = errors.New("access to private IP blocked")
	ErrLocalhostBlocked = errors.New("access to localhost blocked")
	ErrMetadataBlocked  = errors.New("access to cloud metadata blocked")
	ErrInvalidProtocol  = errors.New("invalid or disallowed protocol")

	// Resource limit errors
	ErrResourceLimitExceeded = errors.New("resource limit exceeded")
	ErrTimeoutExceeded       = errors.New("timeout exceeded")
	ErrMemoryLimitExceeded   = errors.New("memory limit exceeded")

	// Access control errors
	ErrPermissionDenied = errors.New("permission denied")
	ErrUnauthorized     = errors.New("unauthorized access")
	ErrAccessForbidden  = errors.New("access forbidden")

	// Rate limiting errors
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrTooManyRequests   = errors.New("too many requests")
)
