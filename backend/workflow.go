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
	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/observer"
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

// Executor types re-exported from pkg/executor
type (
	// NodeExecutor is the interface that custom node executors must implement
	NodeExecutor = executor.NodeExecutor

	// ExecutionContext provides access to workflow state and operations for executors
	ExecutionContext = executor.ExecutionContext

	// Registry manages node executor registration and lookup
	Registry = executor.Registry
)

// Observer types re-exported from pkg/observer
type (
	// Observer receives notifications about workflow execution events
	Observer = observer.Observer

	// Logger is the interface for custom logging
	Logger = observer.Logger

	// Event represents an execution event with metadata
	Event = observer.Event

	// EventType represents the type of execution event
	EventType = observer.EventType

	// ExecutionStatus represents the status of execution
	ExecutionStatus = observer.ExecutionStatus
)

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

// Observer event type constants
const (
	EventWorkflowStart = observer.EventWorkflowStart
	EventWorkflowEnd   = observer.EventWorkflowEnd
	EventNodeStart     = observer.EventNodeStart
	EventNodeEnd       = observer.EventNodeEnd
	EventNodeSuccess   = observer.EventNodeSuccess
	EventNodeFailure   = observer.EventNodeFailure
)

// Execution status constants
const (
	StatusStarted   = observer.StatusStarted
	StatusSuccess   = observer.StatusSuccess
	StatusFailure   = observer.StatusFailure
	StatusCompleted = observer.StatusCompleted
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

	// NewEngineWithRegistry creates a new workflow engine with a custom executor registry.
	// This allows users to register custom node executors.
	NewEngineWithRegistry = engine.NewWithRegistry
)

// Registry functions
var (
	// NewRegistry creates a new empty executor registry
	NewRegistry = executor.NewRegistry

	// DefaultRegistry creates a registry with all built-in node executors registered
	DefaultRegistry = engine.DefaultRegistry
)

// Observer functions
var (
	// NewConsoleObserver creates a console observer with default logger
	NewConsoleObserver = observer.NewConsoleObserver

	// NewConsoleObserverWithLogger creates a console observer with custom logger
	NewConsoleObserverWithLogger = observer.NewConsoleObserverWithLogger

	// NewDefaultLogger creates the default logger implementation
	NewDefaultLogger = observer.NewDefaultLogger
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
