// Package workflow provides a backward-compatible facade for the refactored workflow engine.
//
// DEPRECATED: This package is maintained for backward compatibility only.
// New code should import and use the pkg/* packages directly:
//   - github.com/yesoreyeram/thaiyyal/backend/pkg/types
//   - github.com/yesoreyeram/thaiyyal/backend/pkg/engine
//   - github.com/yesoreyeram/thaiyyal/backend/pkg/executor
//   - github.com/yesoreyeram/thaiyyal/backend/pkg/graph
//   - github.com/yesoreyeram/thaiyyal/backend/pkg/state
//
// # Architecture
//
// The workflow engine consists of:
//   - Type system: 25 node types organized by category (I/O, operations, control flow, state, resilience)
//   - Validation: Comprehensive validation of workflow structure and node data before execution
//   - Execution engine: Parses JSON, infers types, sorts nodes, executes in order
//   - State management: Variables, accumulators, counters, and cache scoped to workflow execution
//
// # Node Categories
//
// Basic I/O: Number, TextInput, Visualization
// Operations: Operation, TextOperation, HTTP
// Control Flow: Condition, ForEach, WhileLoop
// State & Memory: Variable, Extract, Transform, Accumulator, Counter
// Advanced Control: Switch, Parallel, Join, Split, Delay, Cache
// Resilience: Retry, TryCatch, Timeout
// Context: ContextVariable, ContextConstant
//
// # Example Usage
//
//	payload := `{
//	  "nodes": [
//	    {"id": "1", "data": {"value": 10}},
//	    {"id": "2", "data": {"value": 5}},
//	    {"id": "3", "data": {"op": "add"}}
//	  ],
//	  "edges": [
//	    {"source": "1", "target": "3"},
//	    {"source": "2", "target": "3"}
//	  ]
//	}`
//
//	// Create and execute workflow
//	engine, err := workflow.NewEngine([]byte(payload))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	result, err := engine.Execute()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Result: %v\n", result.FinalOutput) // Output: 15
//
// # Design Principles
//
//   - Simplicity: Minimal dependencies, clean interfaces
//   - Testability: Comprehensive test coverage
//   - Type Safety: Strong typing with type inference support
//   - Error Handling: Descriptive errors with context
//   - Validation: Early error detection through comprehensive validation
//   - Modularity: Clear package boundaries and separation of concerns
//
// See ARCHITECTURE.md for detailed architecture documentation.
package workflow

import (
	"github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ============================================================================
// Type Re-exports for Backward Compatibility
// ============================================================================

// Core types re-exported from pkg/types
type (
	NodeType             = types.NodeType
	Node                 = types.Node
	NodeData             = types.NodeData
	Edge                 = types.Edge
	Payload              = types.Payload
	Result               = types.Result
	Config               = types.Config
	SwitchCase           = types.SwitchCase
	ContextVariableValue = types.ContextVariableValue
	CacheEntry           = types.CacheEntry
)

// Engine re-exported from pkg/engine
type Engine = engine.Engine

// ============================================================================
// Constant Re-exports
// ============================================================================

// Context key constants
const (
	ContextKeyExecutionID = types.ContextKeyExecutionID
	ContextKeyWorkflowID  = types.ContextKeyWorkflowID
)

// Node type constants - Basic I/O
const (
	NodeTypeNumber        = types.NodeTypeNumber
	NodeTypeTextInput     = types.NodeTypeTextInput
	NodeTypeVisualization = types.NodeTypeVisualization
)

// Node type constants - Operations
const (
	NodeTypeOperation     = types.NodeTypeOperation
	NodeTypeTextOperation = types.NodeTypeTextOperation
	NodeTypeHTTP          = types.NodeTypeHTTP
)

// Node type constants - Control Flow
const (
	NodeTypeCondition = types.NodeTypeCondition
	NodeTypeForEach   = types.NodeTypeForEach
	NodeTypeWhileLoop = types.NodeTypeWhileLoop
)

// Node type constants - State & Memory
const (
	NodeTypeVariable    = types.NodeTypeVariable
	NodeTypeExtract     = types.NodeTypeExtract
	NodeTypeTransform   = types.NodeTypeTransform
	NodeTypeAccumulator = types.NodeTypeAccumulator
	NodeTypeCounter     = types.NodeTypeCounter
)

// Node type constants - Advanced Control Flow
const (
	NodeTypeSwitch   = types.NodeTypeSwitch
	NodeTypeParallel = types.NodeTypeParallel
	NodeTypeJoin     = types.NodeTypeJoin
	NodeTypeSplit    = types.NodeTypeSplit
	NodeTypeDelay    = types.NodeTypeDelay
	NodeTypeCache    = types.NodeTypeCache
)

// Node type constants - Error Handling & Resilience
const (
	NodeTypeRetry    = types.NodeTypeRetry
	NodeTypeTryCatch = types.NodeTypeTryCatch
	NodeTypeTimeout  = types.NodeTypeTimeout
)

// Node type constants - Context nodes
const (
	NodeTypeContextVariable = types.NodeTypeContextVariable
	NodeTypeContextConstant = types.NodeTypeContextConstant
)

// ============================================================================
// Function Re-exports
// ============================================================================

// Engine constructors
var (
	// NewEngine creates a new workflow engine from JSON payload with default configuration
	NewEngine = engine.New

	// NewEngineWithConfig creates a new workflow engine with custom configuration
	NewEngineWithConfig = engine.NewWithConfig
)

// Context helper functions
var (
	// GetExecutionID extracts the execution ID from context
	GetExecutionID = types.GetExecutionID

	// GetWorkflowID extracts the workflow ID from context
	GetWorkflowID = types.GetWorkflowID
)

// Configuration functions
var (
	// DefaultConfig returns the default engine configuration
	DefaultConfig = types.DefaultConfig

	// ValidationLimits returns limits suitable for strict validation
	ValidationLimits = types.ValidationLimits

	// DevelopmentConfig returns relaxed limits for development/testing
	DevelopmentConfig = types.DevelopmentConfig
)
