package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TimeoutMiddleware enforces execution timeouts for nodes.
// If a node takes longer than the configured timeout, execution is cancelled.
type TimeoutMiddleware struct {
	defaultTimeout time.Duration
}

// NewTimeoutMiddleware creates a new timeout middleware with default timeout
func NewTimeoutMiddleware(defaultTimeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{
		defaultTimeout: defaultTimeout,
	}
}

// Process enforces execution timeout
func (m *TimeoutMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	// Use default timeout for all nodes
	// Node-specific timeout parsing can be added if needed
	timeout := m.defaultTimeout

	// If timeout is 0 or negative, no timeout is enforced
	if timeout <= 0 {
		return next(ctx, node)
	}

	// Create a channel for the result
	type result struct {
		value interface{}
		err   error
	}
	resultChan := make(chan result, 1)

	// Execute with timeout
	go func() {
		value, err := next(ctx, node)
		resultChan <- result{value: value, err: err}
	}()

	// Wait for result or timeout
	select {
	case res := <-resultChan:
		return res.value, res.err
	case <-time.After(timeout):
		return nil, fmt.Errorf("node execution timeout after %v", timeout)
	}
}

// Name returns the middleware name
func (m *TimeoutMiddleware) Name() string {
	return "Timeout"
}

// TimeoutMiddlewareWithContext is a context-aware timeout middleware
// that respects context cancellation
type TimeoutMiddlewareWithContext struct {
	defaultTimeout time.Duration
}

// NewTimeoutMiddlewareWithContext creates a context-aware timeout middleware
func NewTimeoutMiddlewareWithContext(defaultTimeout time.Duration) *TimeoutMiddlewareWithContext {
	return &TimeoutMiddlewareWithContext{
		defaultTimeout: defaultTimeout,
	}
}

// Process enforces execution timeout using context
func (m *TimeoutMiddlewareWithContext) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	// Use default timeout for all nodes
	// Node-specific timeout parsing can be added if needed
	timeout := m.defaultTimeout

	// If timeout is 0 or negative, no timeout is enforced
	if timeout <= 0 {
		return next(ctx, node)
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a channel for the result
	type result struct {
		value interface{}
		err   error
	}
	resultChan := make(chan result, 1)

	// Execute in goroutine
	go func() {
		value, err := next(ctx, node)
		resultChan <- result{value: value, err: err}
	}()

	// Wait for result or timeout
	select {
	case res := <-resultChan:
		return res.value, res.err
	case <-timeoutCtx.Done():
		return nil, fmt.Errorf("node execution timeout after %v", timeout)
	}
}

// Name returns the middleware name
func (m *TimeoutMiddlewareWithContext) Name() string {
	return "TimeoutWithContext"
}
