package logging

import "errors"

// Sentinel errors for logging operations
var (
	// Configuration errors
	ErrInvalidLogLevel  = errors.New("invalid log level")
	ErrInvalidLogFormat = errors.New("invalid log format")
	ErrInvalidOutput    = errors.New("invalid log output")

	// Logging errors
	ErrLogWriteFailed       = errors.New("failed to write log")
	ErrLoggerNotInitialized = errors.New("logger not initialized")
	ErrLogFlushFailed       = errors.New("failed to flush log buffer")
)
