// Package types provides shared type definitions for the workflow engine.
// All core data structures used across packages are defined here to avoid circular dependencies.
package types

import (
	"context"
	"time"
)

// ============================================================================
// Context Keys
// ============================================================================

// contextKey is used for context keys to avoid collisions
type contextKey string

const (
	// ContextKeyExecutionID is the context key for the unique execution ID
	ContextKeyExecutionID contextKey = "execution_id"

	// ContextKeyWorkflowID is the context key for the workflow ID
	ContextKeyWorkflowID contextKey = "workflow_id"
)

// GetExecutionID extracts the execution ID from context.
// Returns empty string if not found in context.
func GetExecutionID(ctx context.Context) string {
	if id, ok := ctx.Value(ContextKeyExecutionID).(string); ok {
		return id
	}
	return ""
}

// GetWorkflowID extracts the workflow ID from context.
// Returns empty string if not found in context.
func GetWorkflowID(ctx context.Context) string {
	if id, ok := ctx.Value(ContextKeyWorkflowID).(string); ok {
		return id
	}
	return ""
}

// ============================================================================
// Node Types
// ============================================================================

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
// Core Data Structures
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
	ContextName   *string                `json:"context_name,omitempty"`   // DEPRECATED: Use ContextValues for multiple values
	ContextValue  interface{}            `json:"context_value,omitempty"`  // DEPRECATED: Use ContextValues for multiple values
	ContextValues []ContextVariableValue `json:"context_values,omitempty"` // for context nodes (multiple typed values)

	// Custom executor fields - extensible data for user-defined node types
	Factor *float64 `json:"factor,omitempty"` // for multiply_by_n custom executor (example)
	Prefix *string  `json:"prefix,omitempty"` // for concat_prefix custom executor (example)
}

// SwitchCase represents a case in a switch node
type SwitchCase struct {
	When       string      `json:"when"`                  // condition or value to match
	Value      interface{} `json:"value,omitempty"`       // value to match (for value matching)
	OutputPath *string     `json:"output_path,omitempty"` // output port name
}

// ContextVariableValue represents a typed value in a context variable/constant node
type ContextVariableValue struct {
	Name  string      `json:"name"`  // Variable name
	Value interface{} `json:"value"` // The actual value
	Type  string      `json:"type"`  // Type: "string", "number", "boolean", "time_string", "epoch_second", "epoch_ms", "null"
}

// Edge represents a connection between nodes
type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// Result represents the execution result of the workflow
type Result struct {
	ExecutionID string                 `json:"execution_id"`          // Unique execution identifier
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

// Config holds workflow engine configuration
type Config struct {
	// Execution limits
	MaxExecutionTime     time.Duration // Maximum time for entire workflow execution
	MaxNodeExecutionTime time.Duration // Maximum time for single node execution
	MaxIterations        int           // Default max iterations for loops (if not specified)

	// HTTP node configuration
	HTTPTimeout           time.Duration // Timeout for HTTP requests
	MaxHTTPRedirects      int           // Maximum number of HTTP redirects to follow
	MaxResponseSize       int64         // Maximum size of HTTP response body (bytes)
	MaxHTTPCallsPerExec   int           // Maximum HTTP calls allowed per workflow execution (0 = unlimited)
	AllowedURLPatterns    []string      // Whitelist of allowed URL patterns (if empty, all external URLs allowed)
	BlockInternalIPs      bool          // Block requests to internal/private IP addresses

	// Cache configuration
	DefaultCacheTTL time.Duration // Default TTL for cache entries if not specified
	MaxCacheSize    int           // Maximum number of cache entries (LRU eviction)

	// Resource limits
	MaxInputSize      int // Maximum size of input data (bytes)
	MaxPayloadSize    int // Maximum size of workflow payload (bytes)
	MaxNodes          int // Maximum number of nodes in workflow
	MaxEdges          int // Maximum number of edges in workflow
	MaxNodeExecutions int // Maximum total node executions (including loop iterations, 0 = unlimited)
	MaxStringLength   int // Maximum length of string values (0 = unlimited)
	MaxArrayLength    int // Maximum length of array values (0 = unlimited)
	MaxVariables      int // Maximum number of variables in workflow state (0 = unlimited)
	MaxContextDepth   int // Maximum depth of nested objects/arrays (0 = unlimited)

	// Retry configuration
	DefaultMaxAttempts int           // Default max retry attempts
	DefaultBackoff     time.Duration // Default initial backoff delay
}
