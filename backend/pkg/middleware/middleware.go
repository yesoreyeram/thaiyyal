// Package middleware provides the Chain of Responsibility pattern implementation
// for node execution middleware. This enables cross-cutting concerns like logging,
// metrics, validation, and timeouts to be added without modifying executor logic.
package middleware

import (
	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Handler is a function that executes a node and returns a result.
// This is the function signature that both executors and middleware use.
type Handler func(ctx executor.ExecutionContext, node types.Node) (interface{}, error)

// Middleware defines the interface for execution middleware.
// Middleware can inspect, modify, or short-circuit node execution.
//
// Example middleware implementations:
//   - LoggingMiddleware: logs execution start/end
//   - MetricsMiddleware: records performance metrics
//   - ValidationMiddleware: validates inputs before execution
//   - TimeoutMiddleware: enforces execution timeouts
//   - RetryMiddleware: retries failed executions
type Middleware interface {
	// Process handles the node execution, optionally calling next() to continue the chain.
	// The middleware can:
	//   - Pre-process: modify context or node before calling next
	//   - Execute: call next to continue the chain
	//   - Post-process: inspect or modify the result after next returns
	//   - Short-circuit: return without calling next (e.g., cache hit)
	//
	// Parameters:
	//   ctx: Execution context with workflow state
	//   node: Node being executed
	//   next: Next handler in the chain (may be another middleware or the executor)
	//
	// Returns:
	//   result: Node execution result
	//   error: Execution error, if any
	Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error)

	// Name returns the middleware name for logging and debugging
	Name() string
}

// Chain represents an ordered chain of middleware.
// Middleware are executed in the order they were added.
type Chain struct {
	middlewares []Middleware
}

// NewChain creates a new middleware chain
func NewChain() *Chain {
	return &Chain{
		middlewares: make([]Middleware, 0),
	}
}

// Use adds middleware to the chain.
// Middleware are executed in the order they are added.
func (c *Chain) Use(middleware Middleware) *Chain {
	c.middlewares = append(c.middlewares, middleware)
	return c
}

// Execute runs the middleware chain followed by the final handler.
// The chain executes middleware in order, with each middleware able to:
//   - Pre-process before calling next
//   - Call next to continue the chain
//   - Post-process after next returns
//   - Short-circuit by not calling next
//
// Example execution flow with 3 middleware:
//   M1.Process(pre) -> M2.Process(pre) -> M3.Process(pre) -> handler() ->
//   M3.Process(post) -> M2.Process(post) -> M1.Process(post) -> return
func (c *Chain) Execute(ctx executor.ExecutionContext, node types.Node, handler Handler) (interface{}, error) {
	if len(c.middlewares) == 0 {
		return handler(ctx, node)
	}

	// Build the chain from the end to the beginning
	// This creates a nested structure where each middleware wraps the next
	index := 0
	var next Handler
	next = func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		if index >= len(c.middlewares) {
			return handler(ctx, node)
		}
		middleware := c.middlewares[index]
		index++
		return middleware.Process(ctx, node, next)
	}

	return next(ctx, node)
}

// Len returns the number of middleware in the chain
func (c *Chain) Len() int {
	return len(c.middlewares)
}

// Middlewares returns all middleware in the chain
func (c *Chain) Middlewares() []Middleware {
	// Return a copy to prevent external modification
	result := make([]Middleware, len(c.middlewares))
	copy(result, c.middlewares)
	return result
}
