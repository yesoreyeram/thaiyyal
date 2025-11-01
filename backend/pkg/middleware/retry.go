package middleware

import (
	"fmt"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// RetryMiddleware automatically retries failed node executions.
// It implements exponential backoff between retry attempts.
type RetryMiddleware struct {
	maxRetries     int
	initialBackoff time.Duration
	maxBackoff     time.Duration
	backoffFactor  float64
}

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxRetries     int           // Maximum number of retry attempts
	InitialBackoff time.Duration // Initial backoff duration
	MaxBackoff     time.Duration // Maximum backoff duration
	BackoffFactor  float64       // Backoff multiplier (e.g., 2.0 for exponential)
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     10 * time.Second,
		BackoffFactor:  2.0,
	}
}

// NewRetryMiddleware creates a new retry middleware with default config
func NewRetryMiddleware() *RetryMiddleware {
	config := DefaultRetryConfig()
	return &RetryMiddleware{
		maxRetries:     config.MaxRetries,
		initialBackoff: config.InitialBackoff,
		maxBackoff:     config.MaxBackoff,
		backoffFactor:  config.BackoffFactor,
	}
}

// NewRetryMiddlewareWithConfig creates a new retry middleware with custom config
func NewRetryMiddlewareWithConfig(config RetryConfig) *RetryMiddleware {
	return &RetryMiddleware{
		maxRetries:     config.MaxRetries,
		initialBackoff: config.InitialBackoff,
		maxBackoff:     config.MaxBackoff,
		backoffFactor:  config.BackoffFactor,
	}
}

// Process retries failed executions with exponential backoff
func (m *RetryMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	var lastErr error
	backoff := m.initialBackoff

	for attempt := 0; attempt <= m.maxRetries; attempt++ {
		// Execute the node
		result, err := next(ctx, node)

		// Success - return immediately
		if err == nil {
			return result, nil
		}

		// Store the error
		lastErr = err

		// If this was the last attempt, return the error
		if attempt == m.maxRetries {
			break
		}

		// Wait before retrying (exponential backoff)
		if backoff > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * m.backoffFactor)
			if backoff > m.maxBackoff {
				backoff = m.maxBackoff
			}
		}
	}

	// All retries failed
	return nil, fmt.Errorf("node execution failed after %d retries: %w", m.maxRetries, lastErr)
}

// Name returns the middleware name
func (m *RetryMiddleware) Name() string {
	return "Retry"
}

// ConditionalRetryMiddleware retries only for specific error types
type ConditionalRetryMiddleware struct {
	maxRetries      int
	initialBackoff  time.Duration
	maxBackoff      time.Duration
	backoffFactor   float64
	retryableErrors []string // List of error message substrings that should trigger retry
}

// NewConditionalRetryMiddleware creates a retry middleware for specific errors
func NewConditionalRetryMiddleware(retryableErrors []string) *ConditionalRetryMiddleware {
	config := DefaultRetryConfig()
	return &ConditionalRetryMiddleware{
		maxRetries:      config.MaxRetries,
		initialBackoff:  config.InitialBackoff,
		maxBackoff:      config.MaxBackoff,
		backoffFactor:   config.BackoffFactor,
		retryableErrors: retryableErrors,
	}
}

// Process retries only for specific error types
func (m *ConditionalRetryMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	var lastErr error
	backoff := m.initialBackoff

	for attempt := 0; attempt <= m.maxRetries; attempt++ {
		result, err := next(ctx, node)

		if err == nil {
			return result, nil
		}

		lastErr = err

		// Check if error is retryable
		if !m.isRetryable(err) {
			return nil, err
		}

		if attempt == m.maxRetries {
			break
		}

		if backoff > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * m.backoffFactor)
			if backoff > m.maxBackoff {
				backoff = m.maxBackoff
			}
		}
	}

	return nil, fmt.Errorf("node execution failed after %d retries: %w", m.maxRetries, lastErr)
}

// isRetryable checks if an error should trigger a retry
func (m *ConditionalRetryMiddleware) isRetryable(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	for _, retryableErr := range m.retryableErrors {
		if contains(errMsg, retryableErr) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Name returns the middleware name
func (m *ConditionalRetryMiddleware) Name() string {
	return "ConditionalRetry"
}
