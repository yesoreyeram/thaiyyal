// Package observer provides an event-driven observer pattern for workflow execution.
//
// # Overview
//
// The observer package implements the observer pattern to enable monitoring,
// logging, and reacting to workflow execution events. Observers can track
// workflow lifecycle, node execution, state changes, and errors without
// coupling to the engine implementation.
//
// # Features
//
//   - Event-driven: React to workflow and node events
//   - Multiple observers: Register multiple observers simultaneously
//   - Lifecycle hooks: Observe all stages of execution
//   - Error events: Track failures and exceptions
//   - State events: Monitor state changes and updates
//   - Performance events: Track timing and metrics
//   - Thread-safe: Concurrent event emission
//
// # Observer Interface
//
//	type Observer interface {
//	    OnWorkflowStart(ctx context.Context, workflow *types.Workflow)
//	    OnWorkflowComplete(ctx context.Context, workflow *types.Workflow, result *types.WorkflowResult)
//	    OnWorkflowError(ctx context.Context, workflow *types.Workflow, err error)
//	    OnNodeStart(ctx context.Context, node *types.Node)
//	    OnNodeComplete(ctx context.Context, node *types.Node, output interface{})
//	    OnNodeError(ctx context.Context, node *types.Node, err error)
//	}
//
// # Workflow Events
//
// WorkflowStart:
//   - Emitted when workflow execution begins
//   - Before validation and preparation
//   - Access to workflow definition
//
// WorkflowComplete:
//   - Emitted when workflow completes successfully
//   - Includes all node outputs
//   - After all nodes executed
//
// WorkflowError:
//   - Emitted when workflow fails
//   - Includes error details
//   - May occur during validation or execution
//
// # Node Events
//
// NodeStart:
//   - Emitted before node execution
//   - Access to node definition
//   - Before inputs are evaluated
//
// NodeComplete:
//   - Emitted after successful node execution
//   - Includes node output
//   - After output validation
//
// NodeError:
//   - Emitted when node execution fails
//   - Includes error details
//   - May trigger workflow failure
//
// # Basic Usage
//
//	import "github.com/yesoreyeram/thaiyyal/backend/pkg/observer"
//
//	// Create observer
//	obs := observer.NewLoggingObserver(logger)
//
//	// Register with engine
//	engine.RegisterObserver(obs)
//
//	// Execute workflow - observer receives events
//	result, err := engine.Execute(ctx, workflow)
//
// # Custom Observer Example
//
//	type MetricsObserver struct {
//	    metrics MetricsCollector
//	}
//
//	func (o *MetricsObserver) OnWorkflowStart(ctx context.Context, workflow *types.Workflow) {
//	    o.metrics.Increment("workflow.started")
//	}
//
//	func (o *MetricsObserver) OnWorkflowComplete(ctx context.Context, workflow *types.Workflow, result *types.WorkflowResult) {
//	    o.metrics.Increment("workflow.completed")
//	    o.metrics.Histogram("workflow.duration", result.Duration)
//	}
//
//	func (o *MetricsObserver) OnWorkflowError(ctx context.Context, workflow *types.Workflow, err error) {
//	    o.metrics.Increment("workflow.failed")
//	}
//
//	func (o *MetricsObserver) OnNodeStart(ctx context.Context, node *types.Node) {
//	    o.metrics.Increment("node.started", map[string]string{"type": string(node.Type)})
//	}
//
//	func (o *MetricsObserver) OnNodeComplete(ctx context.Context, node *types.Node, output interface{}) {
//	    o.metrics.Increment("node.completed", map[string]string{"type": string(node.Type)})
//	}
//
//	func (o *MetricsObserver) OnNodeError(ctx context.Context, node *types.Node, err error) {
//	    o.metrics.Increment("node.failed", map[string]string{"type": string(node.Type)})
//	}
//
// # Built-in Observers
//
// LoggingObserver:
//   - Logs all workflow and node events
//   - Includes timing information
//   - Structured logging with context
//
// MetricsObserver:
//   - Collects execution metrics
//   - Tracks success/failure rates
//   - Records execution duration
//
// DebugObserver:
//   - Detailed debug output
//   - Includes node inputs/outputs
//   - Useful for troubleshooting
//
// EventStreamObserver:
//   - Streams events to external system
//   - Real-time monitoring
//   - Integration with monitoring tools
//
// # Observer Composition
//
// Multiple observers can be registered:
//
//	engine.RegisterObserver(loggingObserver)
//	engine.RegisterObserver(metricsObserver)
//	engine.RegisterObserver(debugObserver)
//
// All observers receive all events in registration order.
//
// # Event Timing
//
// Events are emitted at specific points in execution:
//
//	Workflow Lifecycle:
//	  OnWorkflowStart
//	    → Node Execution (for each node)
//	       OnNodeStart
//	         → Execute Node
//	       OnNodeComplete or OnNodeError
//	  OnWorkflowComplete or OnWorkflowError
//
// # Context Propagation
//
// Events receive the execution context containing:
//
//   - Execution ID
//   - Workflow ID
//   - User information
//   - Request metadata
//   - Cancellation signals
//
// Observers can extract context values:
//
//	executionID := types.GetExecutionID(ctx)
//	workflowID := types.GetWorkflowID(ctx)
//
// # Performance Considerations
//
//   - Observers should not block
//   - Use buffered channels for async processing
//   - Minimize allocations in hot paths
//   - Consider observer overhead for high-throughput
//
// # Error Handling
//
// Observer errors are logged but don't stop execution:
//
//   - Observer panics are recovered
//   - Observer errors are logged
//   - Execution continues normally
//   - Other observers still receive events
//
// # Use Cases
//
//   - Logging and auditing
//   - Metrics collection and monitoring
//   - Real-time dashboards
//   - Debugging and troubleshooting
//   - Performance profiling
//   - Event streaming to external systems
//   - Workflow analytics
//   - Alerting and notifications
//
// # Testing
//
// For testing, use a mock observer:
//
//	type MockObserver struct {
//	    WorkflowStartCalled bool
//	    NodeStartCount int
//	}
//
//	func (o *MockObserver) OnWorkflowStart(ctx context.Context, workflow *types.Workflow) {
//	    o.WorkflowStartCalled = true
//	}
//
//	func (o *MockObserver) OnNodeStart(ctx context.Context, node *types.Node) {
//	    o.NodeStartCount++
//	}
//
// # Best Practices
//
//   - Keep observer logic simple and fast
//   - Use async processing for expensive operations
//   - Don't modify workflow or node state
//   - Handle errors gracefully (don't panic)
//   - Use structured logging for consistency
//   - Consider observer overhead in production
//
// # Thread Safety
//
// Observer methods may be called concurrently from multiple goroutines.
// Implementations must be thread-safe using appropriate synchronization.
package observer
