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

// Node represents a workflow node
type Node struct {
	ID   string   `json:"id"`
	Type NodeType `json:"type,omitempty"`
	Data NodeData `json:"data"`
}

// NodeData contains the node-specific configuration
type NodeData struct {
	Value         *float64 `json:"value,omitempty"`           // for number nodes
	Op            *string  `json:"op,omitempty"`              // for operation nodes
	Mode          *string  `json:"mode,omitempty"`            // for visualization nodes
	Label         *string  `json:"label,omitempty"`           // optional label
	Text          *string  `json:"text,omitempty"`            // for text input nodes
	TextOp        *string  `json:"text_op,omitempty"`         // for text operation nodes
	URL           *string  `json:"url,omitempty"`             // for HTTP nodes
	HTTPClientUID *string  `json:"http_client_uid,omitempty"` // for HTTP nodes - immutable UID of HTTP client from registry
	Separator     *string  `json:"separator,omitempty"`       // for concat text operation
	RepeatN       *int     `json:"repeat_n,omitempty"`        // for repeat text operation
	Condition     *string  `json:"condition,omitempty"`       // for condition, filter, partition, find nodes
	Expression    *string  `json:"expression,omitempty"`      // for map, reduce nodes (transformation expression)
	TruePath      *string  `json:"true_path,omitempty"`       // for condition nodes (output port name)
	FalsePath     *string  `json:"false_path,omitempty"`      // for condition nodes (output port name)
	MaxIterations *int     `json:"max_iterations,omitempty"`  // for for_each and while_loop nodes
	// Array Processing fields
	Start       interface{} `json:"start,omitempty"`        // for slice, range nodes (start index/value)
	End         interface{} `json:"end,omitempty"`          // for slice, range nodes (end index/value)
	Length      interface{} `json:"length,omitempty"`       // for slice node (length instead of end)
	Step        interface{} `json:"step,omitempty"`         // for range node (step value)
	Order       *string     `json:"order,omitempty"`        // for sort node (asc/desc)
	ReturnIndex *bool       `json:"return_index,omitempty"` // for find node
	Size        interface{} `json:"size,omitempty"`         // for chunk node
	Count       interface{} `json:"count,omitempty"`        // for sample node
	Method      *string     `json:"method,omitempty"`       // for sample node (random/first/last)
	Aggregate   *string     `json:"aggregate,omitempty"`    // for group_by node (count/sum/avg/min/max/values)
	ValueField  *string     `json:"value_field,omitempty"`  // for group_by node (field to aggregate)
	Arrays      interface{} `json:"arrays,omitempty"`       // for zip node (array references)
	FillMissing interface{} `json:"fill_missing,omitempty"` // for zip node (value for shorter arrays)
	RemoveEmpty *bool       `json:"remove_empty,omitempty"` // for compact node
	// State & Memory fields
	VarName        *string     `json:"var_name,omitempty"`        // for variable nodes (variable name)
	VarOp          *string     `json:"var_op,omitempty"`          // for variable nodes (get/set)
	Field          *string     `json:"field,omitempty"`           // for extract nodes (field path)
	Fields         []string    `json:"fields,omitempty"`          // for extract nodes (multiple fields)
	TransformType  *string     `json:"transform_type,omitempty"`  // for transform nodes (to_array, to_object, etc.)
	InitialValue   interface{} `json:"initial_value,omitempty"`   // for accumulator/counter initial value
	AccumOp        *string     `json:"accum_op,omitempty"`        // for accumulator operation (sum, product, concat, etc.)
	CounterOp      *string     `json:"counter_op,omitempty"`      // for counter operation (increment, decrement, reset)
	Delta          *float64    `json:"delta,omitempty"`           // for counter delta value
	InputType      *string     `json:"input_type,omitempty"`      // for parse node (AUTO, JSON, CSV, TSV, YAML, XML)
	OutputType     *string     `json:"output_type,omitempty"`     // for format node (JSON, CSV, TSV)
	PrettyPrint    *bool       `json:"pretty_print,omitempty"`    // for format node (JSON pretty printing)
	IncludeHeaders *bool       `json:"include_headers,omitempty"` // for format node (CSV/TSV headers)
	Delimiter      *string     `json:"delimiter,omitempty"`       // for format node (CSV delimiter)
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
	MaxAttempts     *int        `json:"max_attempts,omitempty"`      // for retry node
	BackoffStrategy *string     `json:"backoff_strategy,omitempty"`  // for retry node (exponential/linear/constant)
	InitialDelay    *string     `json:"initial_delay,omitempty"`     // for retry node
	MaxDelay        *string     `json:"max_delay,omitempty"`         // for retry node
	Multiplier      *float64    `json:"multiplier,omitempty"`        // for retry node (backoff multiplier)
	RetryOnErrors   []string    `json:"retry_on_errors,omitempty"`   // for retry node (error patterns to retry on)
	FallbackValue   interface{} `json:"fallback_value,omitempty"`    // for try-catch node
	ContinueOnError *bool       `json:"continue_on_error,omitempty"` // for try-catch node
	ErrorOutputPath *string     `json:"error_output_path,omitempty"` // for try-catch node
	TimeoutAction   *string     `json:"timeout_action,omitempty"`    // for timeout node (error/continue_with_partial)
	// Phase 4: Advanced Node fields
	MaxRequests        *int        `json:"max_requests,omitempty"`        // for rate_limiter node
	PerDuration        *string     `json:"per_duration,omitempty"`        // for rate_limiter node (1s, 1m, 1h)
	RateLimitStrategy  *string     `json:"strategy,omitempty"`            // for rate_limiter node (fixed_window, sliding_window, token_bucket)
	RequestsPerSecond  *float64    `json:"requests_per_second,omitempty"` // for throttle node
	Schema             interface{} `json:"schema,omitempty"`              // for schema_validator node (JSON schema)
	Strict             *bool       `json:"strict,omitempty"`              // for schema_validator node
	PaginationStrategy *string     `json:"pagination_strategy,omitempty"` // for paginator node (offset_limit, page_number, cursor, link_header)
	OffsetParam        *string     `json:"offset_param,omitempty"`        // for paginator node
	LimitParam         *string     `json:"limit_param,omitempty"`         // for paginator node
	PageSize           *int        `json:"page_size,omitempty"`           // for paginator node
	MaxPages           *int        `json:"max_pages,omitempty"`           // for paginator node
	PageParam          *string     `json:"page_param,omitempty"`          // for paginator node
	PerPageParam       *string     `json:"per_page_param,omitempty"`      // for paginator node
	CursorParam        *string     `json:"cursor_param,omitempty"`        // for paginator node
	NextCursorPath     *string     `json:"next_cursor_path,omitempty"`    // for paginator node
	LinkHeader         *string     `json:"link_header,omitempty"`         // for paginator node
	TotalCountPath     *string     `json:"total_count_path,omitempty"`    // for paginator node
	ResultsPath        *string     `json:"results_path,omitempty"`        // for paginator node
	MaxSize            *int        `json:"max_size,omitempty"`            // for enhanced cache node
	Eviction           *string     `json:"eviction,omitempty"`            // for enhanced cache node (lru, lfu, ttl)
	Storage            *string     `json:"storage,omitempty"`             // for enhanced cache node (memory, redis)
	Scope              *string     `json:"scope,omitempty"`               // for variable node (global, workflow, local)
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

// Config is a type alias for backward compatibility.
// The actual configuration is now in the config package.
// Deprecated: Use github.com/yesoreyeram/thaiyyal/backend/pkg/config.Config instead.
type Config = config.Config
