// Package types provides shared type definitions for the workflow engine.
// All core data structures used across packages are defined here to avoid circular dependencies.
package types

import (
	"context"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/config"
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
	NodeTypeBooleanInput  NodeType = "boolean_input"
	NodeTypeDateInput     NodeType = "date_input"
	NodeTypeDateTimeInput NodeType = "datetime_input"
	NodeTypeTextOperation NodeType = "text_operation"
	NodeTypeHTTP          NodeType = "http"
	NodeTypeCondition     NodeType = "condition"
	NodeTypeForEach       NodeType = "for_each"
	NodeTypeWhileLoop     NodeType = "while_loop"
	NodeTypeFilter        NodeType = "filter"     // Filter array elements with expression
	NodeTypeMap           NodeType = "map"        // Transform array elements
	NodeTypeReduce        NodeType = "reduce"     // Reduce array to single value
	NodeTypeExpression    NodeType = "expression" // Apply custom expression to input
	// Array Processing nodes
	NodeTypeSlice     NodeType = "slice"     // Extract portion of array (pagination)
	NodeTypeSort      NodeType = "sort"      // Sort array by field
	NodeTypeFind      NodeType = "find"      // Find first matching element
	NodeTypeFlatMap   NodeType = "flat_map"  // Transform and flatten arrays
	NodeTypeGroupBy   NodeType = "group_by"  // Group and aggregate array elements
	NodeTypeUnique    NodeType = "unique"    // Remove duplicate elements
	NodeTypeChunk     NodeType = "chunk"     // Split array into chunks
	NodeTypeReverse   NodeType = "reverse"   // Reverse array order
	NodeTypePartition NodeType = "partition" // Split array into two groups by condition
	NodeTypeZip       NodeType = "zip"       // Combine arrays element-wise
	NodeTypeSample    NodeType = "sample"    // Get random sample from array
	NodeTypeRange     NodeType = "range"     // Generate array of numbers
	NodeTypeCompact   NodeType = "compact"   // Remove null/empty values
	NodeTypeTranspose NodeType = "transpose" // Transpose 2D array (matrix)
	// State & Memory nodes
	NodeTypeVariable    NodeType = "variable"    // Store/retrieve variables
	NodeTypeExtract     NodeType = "extract"     // Extract fields from objects
	NodeTypeTransform   NodeType = "transform"   // Transform data structures
	NodeTypeAccumulator NodeType = "accumulator" // Accumulate values over time
	NodeTypeCounter     NodeType = "counter"     // Increment/decrement counter
	NodeTypeParse       NodeType = "parse"       // Parse string data to structured formats
	NodeTypeFormat      NodeType = "format"      // Format structured data to string (CSV, JSON, TSV)
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
	// Advanced Nodes (Phase 4)
	NodeTypeRateLimiter     NodeType = "rate_limiter"     // Control request rates
	NodeTypeThrottle        NodeType = "throttle"         // Simple request throttling
	NodeTypeSchemaValidator NodeType = "schema_validator" // Validate against JSON schemas
	NodeTypePaginator       NodeType = "paginator"        // Auto-handle API pagination
	// Context nodes (orphan nodes that define workflow-level values)
	NodeTypeContextVariable NodeType = "context_variable" // Define a mutable variable
	NodeTypeContextConstant NodeType = "context_constant" // Define an immutable constant
	// Visualization nodes
	NodeTypeRenderer NodeType = "renderer" // Render data in various formats
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

// Node represents a workflow node with type-safe data
type Node struct {
	ID   string            `json:"id"`
	Type NodeType          `json:"type,omitempty"`
	Data NodeDataInterface `json:"data"`
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
// Supports conditional execution through sourceHandle and condition fields
type Edge struct {
	ID           string  `json:"id"`
	Source       string  `json:"source"`
	Target       string  `json:"target"`
	SourceHandle *string `json:"sourceHandle,omitempty"` // Output port from source node (e.g., "true", "false", "success", "error")
	TargetHandle *string `json:"targetHandle,omitempty"` // Input port on target node (usually not needed)
	Condition    *string `json:"condition,omitempty"`    // Deprecated: Use sourceHandle instead. Kept for backward compatibility.
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

// Config is a type alias for backward compatibility.
// The actual configuration is now in the config package.
// Deprecated: Use github.com/yesoreyeram/thaiyyal/backend/pkg/config.Config instead.
type Config = config.Config
