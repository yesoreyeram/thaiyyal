// Package executor provides node execution implementations for the workflow engine.
//
// # Overview
//
// This package implements the execution logic for all built-in node types in the
// Thaiyyal workflow engine. Each node type has a dedicated executor that handles
// the specific logic for that operation.
//
// # Node Executor Architecture
//
// Executors implement the Executor interface:
//
//	type Executor interface {
//	    Execute(ctx context.Context, node *types.Node, inputs map[string]interface{}) (interface{}, error)
//	}
//
// Each executor is responsible for:
//
//   - Input validation: Verify inputs meet requirements
//   - Execution logic: Perform the node-specific operation
//   - Output generation: Return results in expected format
//   - Error handling: Provide clear error messages
//
// # Node Categories
//
// Basic Operations:
//   - Number: Constant numeric values
//   - TextInput: Constant text values
//   - Operation: Mathematical operations (add, subtract, multiply, divide)
//   - TextOperation: String manipulation (concat, uppercase, lowercase, etc.)
//   - Visualization: Display values (no transformation)
//
// HTTP Operations:
//   - HTTP: Execute HTTP requests with connection pooling and timeout management
//
// Control Flow:
//   - Condition: Conditional branching based on expressions
//   - Switch: Multi-way branching with pattern matching
//   - ForEach: Iterate over arrays
//   - WhileLoop: Loop while condition is true
//
// Array Processing:
//   - Filter: Select elements matching criteria
//   - Map: Transform each array element
//   - Reduce: Aggregate array to single value
//   - Sort: Sort arrays by field or expression
//   - Slice: Extract array portion (pagination)
//   - Find: Locate first matching element
//   - FlatMap: Transform and flatten nested arrays
//   - GroupBy: Group and aggregate elements
//   - Unique: Remove duplicates
//   - Chunk: Split into fixed-size batches
//   - Reverse: Reverse array order
//   - Partition: Split into two groups by predicate
//   - Zip: Combine multiple arrays element-wise
//   - Sample: Random sampling
//   - Range: Generate numeric sequences
//   - Compact: Remove null/empty values
//   - Transpose: Matrix transposition
//
// State Management:
//   - Variable: Store and retrieve values
//   - Accumulator: Accumulate values over multiple executions
//   - Counter: Increment/decrement counters
//   - Cache: Cache results with TTL
//
// Parallel Processing:
//   - Parallel: Execute multiple branches concurrently
//   - Join: Merge multiple inputs
//   - Split: Fan out to multiple paths
//
// Error Handling & Resilience:
//   - Retry: Automatic retry with exponential backoff
//   - TryCatch: Exception handling with fallback
//   - Timeout: Enforce time limits on operations
//
// Context Management:
//   - ContextVariable: Define workflow-level mutable variables
//   - ContextConstant: Define workflow-level immutable constants
//
// Data Transformation:
//   - Extract: Extract fields from objects
//   - Transform: Complex data structure transformations
//
// Utilities:
//   - Delay: Pause execution for specified duration
//
// # Registry System
//
// All executors are registered in a central registry:
//
//	registry := executor.NewRegistry()
//	registry.Register("number", &NumberExecutor{})
//	registry.Register("operation", &OperationExecutor{})
//
// The registry provides:
//
//   - Type-safe executor lookup
//   - Custom executor registration
//   - Executor validation
//   - Namespace management
//
// # Expression Integration
//
// Many executors integrate with the expression package for dynamic evaluation:
//
//   - Filter conditions
//   - Map transformations
//   - Reduce operations
//   - Conditional branching
//   - Dynamic field access
//
// # Resource Management
//
// Executors manage resources efficiently:
//
//   - HTTP connection pooling (global pool with limits)
//   - Context cancellation support
//   - Timeout enforcement
//   - Memory-efficient array processing
//   - Lazy evaluation where possible
//
// # Error Handling
//
// Executors provide detailed error information:
//
//   - Input validation errors: Clear messages about invalid inputs
//   - Execution errors: Context about what failed and why
//   - Type errors: Expected vs. actual type information
//   - Expression errors: Details about expression evaluation failures
//
// # Performance Characteristics
//
//   - Array operations use streaming where possible
//   - HTTP requests use connection pooling
//   - Expression evaluation is optimized for common cases
//   - State operations minimize lock contention
//
// # Thread Safety
//
// Most executors are stateless and safe for concurrent use.
// Stateful executors (Variable, Accumulator, Counter, Cache) use
// synchronization to ensure thread safety.
//
// # Testing
//
// Each executor has comprehensive test coverage including:
//
//   - Happy path scenarios
//   - Edge cases (empty arrays, null values, etc.)
//   - Error conditions
//   - Performance benchmarks
//   - Concurrent execution tests
//
// # Custom Executors
//
// To create a custom executor:
//
//	type CustomExecutor struct{}
//
//	func (e *CustomExecutor) Execute(ctx context.Context, node *types.Node, inputs map[string]interface{}) (interface{}, error) {
//	    // Validate inputs
//	    // Perform operation
//	    // Return result
//	}
//
//	// Register
//	registry.Register("custom", &CustomExecutor{})
package executor
