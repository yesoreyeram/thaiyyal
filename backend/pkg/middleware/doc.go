// Package middleware provides request/response middleware for the workflow engine.
//
// # Overview
//
// The middleware package implements an interceptor pattern for workflow execution,
// allowing pre-processing, post-processing, and wrapping of workflow and node
// execution. This enables cross-cutting concerns like logging, metrics, caching,
// and security checks.
//
// # Features
//
//   - Workflow middleware: Intercept entire workflow execution
//   - Node middleware: Intercept individual node execution
//   - Chain composition: Stack multiple middleware
//   - Order control: Explicit middleware ordering
//   - Context propagation: Pass data through middleware chain
//   - Error handling: Intercept and transform errors
//
// # Middleware Types
//
// Workflow Middleware:
//
//	Wraps entire workflow execution, can:
//	- Add execution metadata
//	- Implement caching
//	- Add authentication/authorization
//	- Collect workflow-level metrics
//	- Transform workflow before execution
//
// Node Middleware:
//
//	Wraps individual node execution, can:
//	- Add node-level logging
//	- Implement retry logic
//	- Add timeout enforcement
//	- Collect node-level metrics
//	- Transform node inputs/outputs
//
// # Middleware Interface
//
//	type WorkflowMiddleware interface {
//	    Process(ctx context.Context, workflow *types.Workflow, next WorkflowHandler) (*types.WorkflowResult, error)
//	}
//
//	type NodeMiddleware interface {
//	    Process(ctx context.Context, node *types.Node, inputs map[string]interface{}, next NodeHandler) (interface{}, error)
//	}
//
// # Basic Usage
//
//	import "github.com/yesoreyeram/thaiyyal/backend/pkg/middleware"
//
//	// Create logging middleware
//	loggingMW := middleware.NewLogging(logger)
//
//	// Create metrics middleware
//	metricsMW := middleware.NewMetrics(metricsCollector)
//
//	// Apply to engine
//	engine := engine.New(
//	    engine.WithWorkflowMiddleware(loggingMW, metricsMW),
//	)
//
// # Custom Middleware Example
//
// Workflow middleware:
//
//	type TimingMiddleware struct{}
//
//	func (m *TimingMiddleware) Process(ctx context.Context, workflow *types.Workflow, next WorkflowHandler) (*types.WorkflowResult, error) {
//	    start := time.Now()
//	    result, err := next(ctx, workflow)
//	    duration := time.Since(start)
//	    log.Printf("Workflow took %v", duration)
//	    return result, err
//	}
//
// Node middleware:
//
//	type ValidationMiddleware struct{}
//
//	func (m *ValidationMiddleware) Process(ctx context.Context, node *types.Node, inputs map[string]interface{}, next NodeHandler) (interface{}, error) {
//	    // Validate inputs before execution
//	    if err := validateInputs(inputs); err != nil {
//	        return nil, err
//	    }
//	    return next(ctx, node, inputs)
//	}
//
// # Built-in Middleware
//
// Logging Middleware:
//   - Logs workflow and node execution
//   - Includes timing information
//   - Captures errors and results
//
// Metrics Middleware:
//   - Collects execution metrics
//   - Tracks success/failure rates
//   - Measures execution duration
//
// Retry Middleware:
//   - Automatic retry on failure
//   - Exponential backoff
//   - Configurable retry limits
//
// Timeout Middleware:
//   - Enforces execution time limits
//   - Cancels long-running operations
//   - Returns timeout errors
//
// Caching Middleware:
//   - Caches workflow results
//   - Configurable TTL
//   - Cache key generation
//
// Security Middleware:
//   - Input validation
//   - Output sanitization
//   - Permission checks
//
// # Middleware Chain
//
// Middleware executes in order (last registered executes first on the way in):
//
//	Chain:  [Auth] → [Logging] → [Metrics] → [Handler]
//	        ↓         ↓           ↓            ↓
//	Request →→→→→→→→→→→→→→→→→→→→→→→→→→→→→→→→→→ Execute
//	        ←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←← Response
//	        ↑         ↑           ↑            ↑
//	        [Auth]   [Logging]   [Metrics]   [Handler]
//
// # Context Enrichment
//
// Middleware can add data to context:
//
//	func (m *AuthMiddleware) Process(ctx context.Context, workflow *types.Workflow, next WorkflowHandler) (*types.WorkflowResult, error) {
//	    user := authenticateUser(ctx)
//	    ctx = context.WithValue(ctx, "user", user)
//	    return next(ctx, workflow)
//	}
//
// # Error Handling
//
// Middleware can intercept and transform errors:
//
//	func (m *ErrorMiddleware) Process(ctx context.Context, node *types.Node, inputs map[string]interface{}, next NodeHandler) (interface{}, error) {
//	    result, err := next(ctx, node, inputs)
//	    if err != nil {
//	        // Log error
//	        logger.Error("Node failed", err)
//	        // Transform error
//	        return nil, fmt.Errorf("node %s failed: %w", node.ID, err)
//	    }
//	    return result, nil
//	}
//
// # Performance Considerations
//
//   - Minimize allocations in hot paths
//   - Use context for request-scoped data
//   - Avoid blocking operations in middleware
//   - Consider middleware overhead for high-throughput scenarios
//
// # Use Cases
//
//   - Authentication and authorization
//   - Request/response logging
//   - Metrics collection and monitoring
//   - Caching and memoization
//   - Rate limiting and throttling
//   - Input validation and sanitization
//   - Error handling and recovery
//   - Request tracing and correlation
//
// # Best Practices
//
//   - Keep middleware focused on a single concern
//   - Avoid modifying workflow/node state in middleware
//   - Use context for passing request-scoped data
//   - Always call next() unless explicitly stopping the chain
//   - Handle errors appropriately (wrap, transform, or log)
//   - Document middleware ordering requirements
//
// # Thread Safety
//
// Middleware implementations should be stateless and thread-safe.
// The same middleware instance may be used concurrently by multiple
// goroutines.
package middleware
