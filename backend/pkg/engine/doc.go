// Package engine provides the core workflow execution engine for Thaiyyal.
//
// # Overview
//
// The engine package implements a DAG-based workflow execution system that processes
// workflows by executing nodes in topological order. It handles dependency resolution,
// parallel execution, state management, and comprehensive error handling.
//
// # Key Features
//
//   - Topological execution: Nodes execute in dependency order using Kahn's algorithm
//   - Parallel processing: Independent nodes execute concurrently
//   - State management: Maintains execution state and node outputs
//   - Error recovery: Comprehensive error handling with retry and fallback mechanisms
//   - Observer pattern: Extensible event system for monitoring execution
//   - Custom executors: Support for custom node type implementations
//   - Context propagation: Execution metadata flows through the workflow
//
// # Architecture
//
// The engine uses a multi-phase execution model:
//
//  1. Validation: Verify workflow structure, detect cycles, validate node configurations
//  2. Preparation: Build execution graph, initialize state, set up observers
//  3. Execution: Process nodes in topological order with parallel optimization
//  4. Cleanup: Release resources, finalize state, emit completion events
//
// # Basic Usage
//
//	import (
//	    "context"
//	    "github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
//	    "github.com/yesoreyeram/thaiyyal/backend/pkg/types"
//	)
//
//	// Create engine with default configuration
//	eng := engine.New()
//
//	// Execute workflow
//	result, err := eng.Execute(context.Background(), workflow)
//	if err != nil {
//	    log.Fatalf("Execution failed: %v", err)
//	}
//
//	// Access node outputs
//	for nodeID, output := range result.Outputs {
//	    fmt.Printf("Node %s: %v\n", nodeID, output)
//	}
//
// # Advanced Usage
//
//	// Create engine with custom configuration
//	eng := engine.New(
//	    engine.WithMaxParallel(10),
//	    engine.WithTimeout(5 * time.Minute),
//	    engine.WithObserver(myObserver),
//	)
//
//	// Register custom executor
//	eng.RegisterExecutor("custom_node", &CustomExecutor{})
//
//	// Execute with context metadata
//	ctx := context.WithValue(context.Background(), types.ContextKeyExecutionID, "exec-123")
//	result, err := eng.Execute(ctx, workflow)
//
// # Error Handling
//
// The engine provides detailed error information:
//
//   - Validation errors: Issues with workflow structure or configuration
//   - Execution errors: Runtime failures during node execution
//   - Timeout errors: Execution exceeded time limits
//   - Cycle detection: Circular dependencies in workflow graph
//
// All errors include context about the failing node and execution state.
//
// # Concurrency
//
// The engine is designed for concurrent execution:
//
//   - Nodes without dependencies execute in parallel
//   - Each node execution is isolated with its own goroutine
//   - State updates are synchronized to prevent race conditions
//   - Context cancellation propagates to all running nodes
//
// # Performance Considerations
//
//   - Use context cancellation for long-running workflows
//   - Configure max parallel workers based on workload
//   - Large workflows benefit from parallel execution
//   - Observer overhead is minimal but measurable
//
// # Extensibility
//
// The engine supports multiple extension points:
//
//   - Custom executors: Implement custom node types
//   - Observers: Monitor and react to execution events
//   - Middleware: Intercept and modify execution flow
//   - State providers: Custom state storage backends
//
// # Thread Safety
//
// The Engine type is safe for concurrent use by multiple goroutines.
// A single engine instance can execute multiple workflows concurrently.
package engine
