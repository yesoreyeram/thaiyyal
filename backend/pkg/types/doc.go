// Package types provides shared type definitions for the Thaiyyal workflow engine.
//
// # Overview
//
// This package contains all core data structures and type definitions used across
// the workflow engine. It serves as the foundation for avoiding circular dependencies
// between other packages while providing a consistent type system.
//
// # Key Components
//
// Node Types: Define the various types of workflow nodes (operation, control flow,
// error handling, etc.)
//
// Workflow Structure: Core data structures for workflows, nodes, edges, and execution state
//
// Execution Context: Context keys and helpers for passing execution metadata
//
// Data Types: Type system for runtime values (string, number, boolean, array, object, null)
//
// # Node Categories
//
// The workflow engine supports several categories of nodes:
//
//   - Basic I/O: Number, TextInput, Visualization
//   - Operations: Math operations, text manipulation, data transformation
//   - HTTP: HTTP request execution
//   - Control Flow: Condition, ForEach, WhileLoop, Switch
//   - Parallel Processing: Parallel, Join, Split
//   - State Management: Variable, Accumulator, Counter, Cache
//   - Array Processing: Filter, Map, Reduce, Sort, Slice, etc.
//   - Error Handling: Retry, TryCatch, Timeout
//   - Context: ContextVariable, ContextConstant
//
// # Usage Example
//
//	workflow := &types.Workflow{
//	    Name: "Example Workflow",
//	    Nodes: []types.Node{
//	        {ID: "1", Type: types.NodeTypeNumber, Data: types.NodeData{Value: 42}},
//	        {ID: "2", Type: types.NodeTypeOperation, Data: types.NodeData{Operation: "add", Value: 10}},
//	    },
//	    Edges: []types.Edge{
//	        {Source: "1", Target: "2"},
//	    },
//	}
//
// # Design Principles
//
//   - Minimal dependencies: Types package has no dependencies on other workflow packages
//   - Immutability: Core types are designed to be immutable where possible
//   - Extensibility: Easy to add new node types without breaking changes
//   - Type safety: Strong typing for workflow components and runtime values
//
// # Thread Safety
//
// The types defined in this package are generally not thread-safe for mutation.
// Concurrent access should be coordinated by the caller using appropriate synchronization.
//
// # Backward Compatibility
//
// This package maintains backward compatibility. Breaking changes are avoided, and
// deprecation warnings are provided when types are phased out.
package types
