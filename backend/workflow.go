// Package workflow provides a workflow execution engine for visual workflow builders.
//
// The engine accepts JSON payloads containing node and edge definitions, performs
// topological sorting to determine execution order, and executes nodes sequentially
// in a directed acyclic graph (DAG) structure.
//
// # Architecture
//
// The workflow engine consists of:
//   - Type system: 23 node types organized by category (I/O, operations, control flow, state, resilience)
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
//	// Validate workflow before execution (recommended)
//	validationResult, err := workflow.ValidatePayload([]byte(payload))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if !validationResult.Valid {
//	    log.Fatalf("Workflow validation failed: %v", validationResult.Errors)
//	}
//
//	// Execute workflow
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
//   - Simplicity: Single package, no external dependencies
//   - Testability: Comprehensive test coverage (165+ tests)
//   - Type Safety: Strong typing with type inference support
//   - Error Handling: Descriptive errors with context
//   - Validation: Early error detection through comprehensive validation
//
// See ARCHITECTURE.md for detailed architecture documentation.
package workflow

import (
"context"
"crypto/rand"
"encoding/hex"
"encoding/json"
"fmt"
"sync"
"time"
)

// ============================================================================
// Core Types and Constants
// ============================================================================

// Context keys for passing execution metadata through the execution chain
type contextKey string

const (
	// ContextKeyExecutionID is the context key for the unique execution ID
	ContextKeyExecutionID contextKey = "execution_id"
	
	// ContextKeyWorkflowID is the context key for the workflow ID
	ContextKeyWorkflowID contextKey = "workflow_id"
)

// NodeType represents the type of a workflow node
type NodeType string

const (
NodeTypeNumber        NodeType = "number"
NodeTypeOperation     NodeType = "operation"
NodeTypeVisualization NodeType = "visualization"
NodeTypeTextInput     NodeType = "text_input"
NodeTypeTextOperation NodeType = "text_operation"
NodeTypeHTTP          NodeType = "http"
NodeTypeCondition     NodeType = "condition"
NodeTypeForEach       NodeType = "for_each"
NodeTypeWhileLoop     NodeType = "while_loop"
// State & Memory nodes
NodeTypeVariable    NodeType = "variable"    // Store/retrieve variables
NodeTypeExtract     NodeType = "extract"     // Extract fields from objects
NodeTypeTransform   NodeType = "transform"   // Transform data structures
NodeTypeAccumulator NodeType = "accumulator" // Accumulate values over time
NodeTypeCounter     NodeType = "counter"     // Increment/decrement counter
// Advanced Control Flow nodes
NodeTypeSwitch   NodeType = "switch"   // Multi-way branching
NodeTypeParallel NodeType = "parallel" // Parallel execution
NodeTypeJoin     NodeType = "join"     // Combine multiple inputs
NodeTypeSplit    NodeType = "split"    // Split to multiple paths
NodeTypeDelay    NodeType = "delay"    // Delay execution
NodeTypeCache    NodeType = "cache"    // Cache get/set operations
// Error Handling & Resilience nodes
NodeTypeRetry    NodeType = "retry"     // Retry with backoff
NodeTypeTryCatch NodeType = "try_catch" // Error handling with fallback
NodeTypeTimeout  NodeType = "timeout"   // Enforce time limits
// Context nodes (orphan nodes that define workflow-level values)
NodeTypeContextVariable NodeType = "context_variable" // Define a mutable variable
NodeTypeContextConstant NodeType = "context_constant" // Define an immutable constant
)

// ============================================================================
// Helper Functions
// ============================================================================

// generateExecutionID creates a unique execution identifier.
// Uses crypto/rand for cryptographically secure random IDs.
// Format: 16 hex characters (8 bytes) for balance between uniqueness and readability.
// Example: "a1b2c3d4e5f6g7h8"
func generateExecutionID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("exec_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// GetExecutionID extracts the execution ID from context.
// Returns empty string if not found in context.
//
// This can be used by node executors for logging:
//
//	executionID := workflow.GetExecutionID(ctx)
//	log.Printf("[%s] Processing node", executionID)
func GetExecutionID(ctx context.Context) string {
	if id, ok := ctx.Value(ContextKeyExecutionID).(string); ok {
		return id
	}
	return ""
}

// GetWorkflowID extracts the workflow ID from context.
// Returns empty string if not found in context.
//
// This can be used by node executors for logging:
//
//	workflowID := workflow.GetWorkflowID(ctx)
//	log.Printf("[workflow:%s] Processing", workflowID)
func GetWorkflowID(ctx context.Context) string {
	if id, ok := ctx.Value(ContextKeyWorkflowID).(string); ok {
		return id
	}
	return ""
}

// ============================================================================
// Payload and Node Types
// ============================================================================

// Payload represents the JSON payload from the frontend
type Payload struct {
	WorkflowID string `json:"workflow_id,omitempty"` // Optional workflow identifier
	Nodes      []Node `json:"nodes"`
	Edges      []Edge `json:"edges"`
}

// Node represents a workflow node
type Node struct {
ID   string   `json:"id"`
Type NodeType `json:"type,omitempty"`
Data NodeData `json:"data"`
}

// NodeData contains the node-specific configuration
// This uses the Composite Pattern to support multiple node types
type NodeData struct {
Value         *float64 `json:"value,omitempty"`          // for number nodes
Op            *string  `json:"op,omitempty"`             // for operation nodes
Mode          *string  `json:"mode,omitempty"`           // for visualization nodes
Label         *string  `json:"label,omitempty"`          // optional label
Text          *string  `json:"text,omitempty"`           // for text input nodes
TextOp        *string  `json:"text_op,omitempty"`        // for text operation nodes
URL           *string  `json:"url,omitempty"`            // for HTTP nodes
Separator     *string  `json:"separator,omitempty"`      // for concat text operation
RepeatN       *int     `json:"repeat_n,omitempty"`       // for repeat text operation
Condition     *string  `json:"condition,omitempty"`      // for condition nodes
TruePath      *string  `json:"true_path,omitempty"`      // for condition nodes (output port name)
FalsePath     *string  `json:"false_path,omitempty"`     // for condition nodes (output port name)
MaxIterations *int     `json:"max_iterations,omitempty"` // for for_each and while_loop nodes
// State & Memory fields
VarName       *string     `json:"var_name,omitempty"`       // for variable nodes (variable name)
VarOp         *string     `json:"var_op,omitempty"`         // for variable nodes (get/set)
Field         *string     `json:"field,omitempty"`          // for extract nodes (field path)
Fields        []string    `json:"fields,omitempty"`         // for extract nodes (multiple fields)
TransformType *string     `json:"transform_type,omitempty"` // for transform nodes (to_array, to_object, etc.)
InitialValue  interface{} `json:"initial_value,omitempty"`  // for accumulator/counter initial value
AccumOp       *string     `json:"accum_op,omitempty"`       // for accumulator operation (sum, product, concat, etc.)
CounterOp     *string     `json:"counter_op,omitempty"`     // for counter operation (increment, decrement, reset)
Delta         *float64    `json:"delta,omitempty"`          // for counter delta value
// Advanced Control Flow fields
Cases          []SwitchCase `json:"cases,omitempty"`           // for switch node (case definitions)
DefaultPath    *string      `json:"default_path,omitempty"`    // for switch node (default case)
MaxConcurrency *int         `json:"max_concurrency,omitempty"` // for parallel node
JoinStrategy   *string      `json:"join_strategy,omitempty"`   // for join node (all/any/first)
Timeout        *string      `json:"timeout,omitempty"`         // for join/parallel/timeout nodes
Paths          []string     `json:"paths,omitempty"`           // for split node
Duration       *string      `json:"duration,omitempty"`        // for delay node
CacheOp        *string      `json:"cache_op,omitempty"`        // for cache node (get/set)
CacheKey       *string      `json:"cache_key,omitempty"`       // for cache node
TTL            *string      `json:"ttl,omitempty"`             // for cache node
// Error Handling & Resilience fields
MaxAttempts      *int        `json:"max_attempts,omitempty"`      // for retry node
BackoffStrategy  *string     `json:"backoff_strategy,omitempty"`  // for retry node (exponential/linear/constant)
InitialDelay     *string     `json:"initial_delay,omitempty"`     // for retry node
MaxDelay         *string     `json:"max_delay,omitempty"`         // for retry node
Multiplier       *float64    `json:"multiplier,omitempty"`        // for retry node (backoff multiplier)
RetryOnErrors    []string    `json:"retry_on_errors,omitempty"`   // for retry node (error patterns to retry on)
FallbackValue    interface{} `json:"fallback_value,omitempty"`    // for try-catch node
ContinueOnError  *bool       `json:"continue_on_error,omitempty"` // for try-catch node
ErrorOutputPath  *string     `json:"error_output_path,omitempty"` // for try-catch node
TimeoutAction    *string     `json:"timeout_action,omitempty"`    // for timeout node (error/continue_with_partial)
// Context node fields
ContextName  *string     `json:"context_name,omitempty"`  // DEPRECATED: Use ContextValues for multiple values
ContextValue interface{} `json:"context_value,omitempty"` // DEPRECATED: Use ContextValues for multiple values
ContextValues []ContextVariableValue `json:"context_values,omitempty"` // for context nodes (multiple typed values)
}

// SwitchCase represents a case in a switch node
type SwitchCase struct {
	When       string      `json:"when"`                  // condition or value to match
	Value      interface{} `json:"value,omitempty"`       // value to match (for value matching)
	OutputPath *string     `json:"output_path,omitempty"` // output port name
}

// ContextVariableValue represents a typed value in a context variable/constant node.
// Supports multiple values with explicit type information for proper type conversion.
type ContextVariableValue struct {
	Name  string      `json:"name"`            // Variable name
	Value interface{} `json:"value"`           // The actual value
	Type  string      `json:"type"`            // Type: "string", "number", "boolean", "time_string", "epoch_second", "epoch_ms", "null"
}

// Edge represents a connection between nodes
type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// Result represents the execution result of the workflow
type Result struct {
	ExecutionID string                 `json:"execution_id"`         // Unique execution identifier
	WorkflowID  string                 `json:"workflow_id,omitempty"` // Optional workflow identifier
	NodeResults map[string]interface{} `json:"node_results"`
	FinalOutput interface{}            `json:"final_output"`
	Errors      []string               `json:"errors,omitempty"`
}

// CacheEntry represents a cached value with expiration
type CacheEntry struct {
Value      interface{}
Expiration time.Time
}

// ============================================================================
// Engine Definition
// ============================================================================

// Engine is the workflow execution engine.
// It manages workflow state and coordinates node execution in topological order.
//
// The Engine uses the following design patterns:
//   - Strategy Pattern: Different execution strategies for different node types
//   - State Pattern: Manages workflow state (variables, accumulator, counter, cache)
//   - Template Method: Execute() defines the workflow execution algorithm
type Engine struct {
	nodes       []Node
	edges       []Edge
	nodeResults map[string]interface{}
	resultsMutex sync.RWMutex // protects nodeResults for parallel execution
	config      Config // configuration for execution limits and security
	
	// Execution metadata
	executionID string // unique identifier for this execution
	workflowID  string // optional identifier for the workflow definition
	
	// State management
	variables   map[string]interface{} // stores variables across nodes
	accumulator interface{}            // stores accumulated value
	counter     float64                // stores counter value
	
	// Cache management
	cache      map[string]*CacheEntry // stores cached values with TTL
	cacheMutex sync.RWMutex           // protects cache access
	
	// Context for template interpolation (populated by context nodes)
	contextVariables map[string]interface{} // workflow-level variables from context_variable nodes
	contextConstants map[string]interface{} // workflow-level constants from context_constant nodes
}

// ============================================================================
// Public API
// ============================================================================

// NewEngine creates a new workflow engine from JSON payload.
//
// The payload should contain:
//   - workflow_id (optional): Identifier for the workflow definition
//   - nodes: Array of node definitions with id, type (optional), and data
//   - edges: Array of edge definitions connecting nodes
//
// An execution ID is automatically generated for this execution.
//
// Returns:
//   - *Engine: Initialized engine ready for execution
//   - error: If JSON parsing fails
func NewEngine(payloadJSON []byte) (*Engine, error) {
	var payload Payload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	return &Engine{
		nodes:            payload.Nodes,
		edges:            payload.Edges,
		nodeResults:      make(map[string]interface{}),
		config:           DefaultConfig(),
		executionID:      generateExecutionID(),
		workflowID:       payload.WorkflowID,
		variables:        make(map[string]interface{}),
		accumulator:      nil,
		counter:          0,
		cache:            make(map[string]*CacheEntry),
		contextVariables: make(map[string]interface{}),
		contextConstants: make(map[string]interface{}),
	}, nil
}

// NewEngineWithConfig creates a new workflow engine with custom configuration.
// This is useful for testing or when you need non-default security settings.
//
// An execution ID is automatically generated for this execution.
func NewEngineWithConfig(payloadJSON []byte, config Config) (*Engine, error) {
	var payload Payload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	return &Engine{
		nodes:            payload.Nodes,
		edges:            payload.Edges,
		nodeResults:      make(map[string]interface{}),
		config:           config,
		executionID:      generateExecutionID(),
		workflowID:       payload.WorkflowID,
		variables:        make(map[string]interface{}),
		accumulator:      nil,
		counter:          0,
		cache:            make(map[string]*CacheEntry),
		contextVariables: make(map[string]interface{}),
		contextConstants: make(map[string]interface{}),
	}, nil
}


// Execute runs the workflow and returns the result.
// The workflow execution is protected by a timeout configured in MaxExecutionTime.
// If the workflow takes longer than the timeout, execution is cancelled and an error is returned.
//
// Each workflow execution is assigned a unique execution ID that is passed through the
// execution context and included in the result. This ID can be used for logging and tracing.
//
// Returns:
//   - *Result: Workflow execution results including execution ID, node outputs and final output
//   - error: If execution fails, times out, or encounters an error
func (e *Engine) Execute() (*Result, error) {
	result := &Result{
		ExecutionID: e.executionID,
		WorkflowID:  e.workflowID,
		NodeResults: make(map[string]interface{}),
		Errors:      []string{},
	}

	// Step 1: Infer node types if not set
	e.inferNodeTypes()

	// Step 2: Get execution order using topological sort
	executionOrder, err := e.topologicalSort()
	if err != nil {
		return result, err
	}

	// Step 3: Create context with timeout and execution metadata for workflow execution
	ctx, cancel := context.WithTimeout(context.Background(), e.config.MaxExecutionTime)
	defer cancel()
	
	// Add execution ID and workflow ID to context for logging and tracing
	ctx = context.WithValue(ctx, ContextKeyExecutionID, e.executionID)
	ctx = context.WithValue(ctx, ContextKeyWorkflowID, e.workflowID)

	// Use a channel to communicate execution completion
	done := make(chan error, 1)

	// Execute workflow in a goroutine
	go func() {
		// Execute each node in order
		for _, nodeID := range executionOrder {
			// Check if context was cancelled (timeout or parent cancellation)
			select {
			case <-ctx.Done():
				done <- ctx.Err()
				return
			default:
			}

			node := e.getNode(nodeID)
			value, err := e.executeNodeWithContext(ctx, node)
			if err != nil {
				errMsg := fmt.Sprintf("error executing node %s: %v", nodeID, err)
				result.Errors = append(result.Errors, errMsg)
				done <- fmt.Errorf("%s", errMsg)
				return
			}
			e.nodeResults[nodeID] = value
		}
		done <- nil
	}()

	// Wait for execution to complete or timeout
	select {
	case err := <-done:
		if err != nil {
			return result, err
		}
	case <-ctx.Done():
		return result, fmt.Errorf("workflow execution timeout: exceeded %v", e.config.MaxExecutionTime)
	}

	// Step 4: Copy results and set final output
	result.NodeResults = e.nodeResults
	result.FinalOutput = e.getFinalOutput()

	return result, nil
}
