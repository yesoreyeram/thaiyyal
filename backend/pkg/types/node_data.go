package types

import "fmt"

// ============================================================================
// NodeData Interface - Type-safe node data
// ============================================================================

// NodeData is an interface that all node-specific data types must implement.
// This ensures type safety while allowing extensibility for custom node types.
type NodeDataInterface interface {
	// Validate checks if the node data is valid
	Validate() error
	// GetLabel returns the optional label for the node
	GetLabel() string
}

// ============================================================================
// Common Data - Shared fields across node types
// ============================================================================

// CommonData contains fields that are common across many node types
type CommonData struct {
	Label *string `json:"label,omitempty"`
}

func (c CommonData) GetLabel() string {
	if c.Label != nil {
		return *c.Label
	}
	return ""
}

// ============================================================================
// Input Node Data Types
// ============================================================================

// NumberData contains data for number input nodes
type NumberData struct {
	CommonData
	Value *float64 `json:"value,omitempty"`
}

func (d NumberData) Validate() error {
	if d.Value == nil {
		return ErrMissingRequiredField("value")
	}
	return nil
}

// TextInputData contains data for text input nodes
type TextInputData struct {
	CommonData
	Text *string `json:"text,omitempty"`
}

func (d TextInputData) Validate() error {
	if d.Text == nil {
		return ErrMissingRequiredField("text")
	}
	return nil
}

// BooleanInputData contains data for boolean input nodes
type BooleanInputData struct {
	CommonData
	BooleanValue *bool `json:"boolean_value,omitempty"`
}

func (d BooleanInputData) Validate() error {
	if d.BooleanValue == nil {
		return ErrMissingRequiredField("boolean_value")
	}
	return nil
}

// DateInputData contains data for date input nodes
type DateInputData struct {
	CommonData
	DateValue *string `json:"date_value,omitempty"` // YYYY-MM-DD format
}

func (d DateInputData) Validate() error {
	if d.DateValue == nil {
		return ErrMissingRequiredField("date_value")
	}
	return nil
}

// DateTimeInputData contains data for datetime input nodes
type DateTimeInputData struct {
	CommonData
	DateTimeValue *string `json:"datetime_value,omitempty"` // ISO 8601 format
}

func (d DateTimeInputData) Validate() error {
	if d.DateTimeValue == nil {
		return ErrMissingRequiredField("datetime_value")
	}
	return nil
}

// ============================================================================
// Operation Node Data Types
// ============================================================================

// OperationData contains data for operation nodes
type OperationData struct {
	CommonData
	Op *string `json:"op,omitempty"` // add, subtract, multiply, divide, etc.
}

func (d OperationData) Validate() error {
	if d.Op == nil {
		return ErrMissingRequiredField("op")
	}
	return nil
}

// TextOperationData contains data for text operation nodes
type TextOperationData struct {
	CommonData
	TextOp    *string `json:"text_op,omitempty"`   // concat, uppercase, lowercase, etc.
	Separator *string `json:"separator,omitempty"` // for concat operation
	RepeatN   *int    `json:"repeat_n,omitempty"`  // for repeat operation
}

func (d TextOperationData) Validate() error {
	if d.TextOp == nil {
		return ErrMissingRequiredField("text_op")
	}
	return nil
}

// HTTPData contains data for HTTP request nodes
type HTTPData struct {
	CommonData
	URL           *string `json:"url,omitempty"`
	HTTPClientUID *string `json:"http_client_uid,omitempty"` // Optional named client
}

func (d HTTPData) Validate() error {
	if d.URL == nil {
		return ErrMissingRequiredField("url")
	}
	return nil
}

// ExpressionData contains data for expression nodes
type ExpressionData struct {
	CommonData
	Expression *string `json:"expression,omitempty"`
}

func (d ExpressionData) Validate() error {
	if d.Expression == nil {
		return ErrMissingRequiredField("expression")
	}
	return nil
}

// ============================================================================
// Control Flow Node Data Types
// ============================================================================

// ConditionData contains data for condition nodes
type ConditionData struct {
	CommonData
	Condition *string `json:"condition,omitempty"`
	TruePath  *string `json:"true_path,omitempty"`
	FalsePath *string `json:"false_path,omitempty"`
}

func (d ConditionData) Validate() error {
	if d.Condition == nil {
		return ErrMissingRequiredField("condition")
	}
	return nil
}

// ForEachData contains data for for_each nodes
type ForEachData struct {
	CommonData
	MaxIterations *int `json:"max_iterations,omitempty"`
}

func (d ForEachData) Validate() error {
	return nil // MaxIterations is optional
}

// WhileLoopData contains data for while_loop nodes
type WhileLoopData struct {
	CommonData
	Condition     *string `json:"condition,omitempty"`
	MaxIterations *int    `json:"max_iterations,omitempty"`
}

func (d WhileLoopData) Validate() error {
	if d.Condition == nil {
		return ErrMissingRequiredField("condition")
	}
	return nil
}

// FilterData contains data for filter nodes
type FilterData struct {
	CommonData
	Condition *string `json:"condition,omitempty"`
}

func (d FilterData) Validate() error {
	if d.Condition == nil {
		return ErrMissingRequiredField("condition")
	}
	return nil
}

// MapData contains data for map transformation nodes
type MapData struct {
	CommonData
	Expression *string `json:"expression,omitempty"`
	Field      *string `json:"field,omitempty"` // Field to extract from each element
}

func (d MapData) Validate() error {
	// Either expression or field should be provided, but both are optional
	return nil
}

// ReduceData contains data for reduce nodes
type ReduceData struct {
	CommonData
	Expression   *string     `json:"expression,omitempty"`
	InitialValue interface{} `json:"initial_value,omitempty"`
}

func (d ReduceData) Validate() error {
	if d.Expression == nil {
		return ErrMissingRequiredField("expression")
	}
	return nil
}

// ============================================================================
// Array Processing Node Data Types
// ============================================================================

// SliceData contains data for slice nodes
type SliceData struct {
	CommonData
	Start  interface{} `json:"start,omitempty"`
	End    interface{} `json:"end,omitempty"`
	Length interface{} `json:"length,omitempty"`
}

func (d SliceData) Validate() error {
	return nil // All fields are optional
}

// SortData contains data for sort nodes
type SortData struct {
	CommonData
	Field *string `json:"field,omitempty"` // Field to sort by
	Order *string `json:"order,omitempty"` // asc/desc
}

func (d SortData) Validate() error {
	return nil // Fields are optional (can sort primitive arrays)
}

// FindData contains data for find nodes
type FindData struct {
	CommonData
	Condition   *string `json:"condition,omitempty"`
	ReturnIndex *bool   `json:"return_index,omitempty"`
}

func (d FindData) Validate() error {
	if d.Condition == nil {
		return ErrMissingRequiredField("condition")
	}
	return nil
}

// FlatMapData contains data for flat_map nodes
type FlatMapData struct {
	CommonData
	Expression *string `json:"expression,omitempty"`
	Field      *string `json:"field,omitempty"` // Field to extract from each element before flattening
}

func (d FlatMapData) Validate() error {
	// Either expression or field should be provided, but both are optional
	return nil
}

// GroupByData contains data for group_by nodes
type GroupByData struct {
	CommonData
	Field      *string `json:"field,omitempty"`       // Field to group by
	Aggregate  *string `json:"aggregate,omitempty"`   // count, sum, avg, min, max, values
	ValueField *string `json:"value_field,omitempty"` // Field to aggregate
}

func (d GroupByData) Validate() error {
	if d.Field == nil {
		return ErrMissingRequiredField("field")
	}
	return nil
}

// UniqueData contains data for unique nodes
type UniqueData struct {
	CommonData
	Field *string `json:"field,omitempty"` // Optional field to check uniqueness
}

func (d UniqueData) Validate() error {
	return nil // Field is optional
}

// ChunkData contains data for chunk nodes
type ChunkData struct {
	CommonData
	Size interface{} `json:"size,omitempty"`
}

func (d ChunkData) Validate() error {
	if d.Size == nil {
		return ErrMissingRequiredField("size")
	}
	return nil
}

// ReverseData contains data for reverse nodes
type ReverseData struct {
	CommonData
}

func (d ReverseData) Validate() error {
	return nil // No required fields
}

// PartitionData contains data for partition nodes
type PartitionData struct {
	CommonData
	Condition *string `json:"condition,omitempty"`
}

func (d PartitionData) Validate() error {
	if d.Condition == nil {
		return ErrMissingRequiredField("condition")
	}
	return nil
}

// ZipData contains data for zip nodes
type ZipData struct {
	CommonData
	Arrays      interface{} `json:"arrays,omitempty"`
	FillMissing interface{} `json:"fill_missing,omitempty"`
	RemoveEmpty *bool       `json:"remove_empty,omitempty"` // Remove entries where any array is missing
}

func (d ZipData) Validate() error {
	return nil // Arrays are optional (can zip with inputs)
}

// SampleData contains data for sample nodes
type SampleData struct {
	CommonData
	Count  interface{} `json:"count,omitempty"`  // Number of samples
	Method *string     `json:"method,omitempty"` // random, first, last
}

func (d SampleData) Validate() error {
	if d.Count == nil {
		return ErrMissingRequiredField("count")
	}
	return nil
}

// RangeData contains data for range nodes
type RangeData struct {
	CommonData
	Start interface{} `json:"start,omitempty"`
	End   interface{} `json:"end,omitempty"`
	Step  interface{} `json:"step,omitempty"`
}

func (d RangeData) Validate() error {
	if d.Start == nil || d.End == nil {
		return ErrMissingRequiredField("start and end")
	}
	return nil
}

// CompactData contains data for compact nodes
type CompactData struct {
	CommonData
	RemoveEmpty *bool `json:"remove_empty,omitempty"`
}

func (d CompactData) Validate() error {
	return nil // RemoveEmpty is optional
}

// TransposeData contains data for transpose nodes
type TransposeData struct {
	CommonData
}

func (d TransposeData) Validate() error {
	return nil // No required fields
}

// ============================================================================
// State & Memory Node Data Types
// ============================================================================

// VariableData contains data for variable nodes
type VariableData struct {
	CommonData
	VarName *string `json:"var_name,omitempty"`
	VarOp   *string `json:"var_op,omitempty"` // get, set
	Scope   *string `json:"scope,omitempty"`  // global, workflow, local
}

func (d VariableData) Validate() error {
	if d.VarName == nil {
		return ErrMissingRequiredField("var_name")
	}
	if d.VarOp == nil {
		return ErrMissingRequiredField("var_op")
	}
	return nil
}

// ExtractData contains data for extract nodes
type ExtractData struct {
	CommonData
	Field  *string  `json:"field,omitempty"`  // Single field path
	Fields []string `json:"fields,omitempty"` // Multiple field paths
}

func (d ExtractData) Validate() error {
	if d.Field == nil && len(d.Fields) == 0 {
		return ErrMissingRequiredField("field or fields")
	}
	return nil
}

// TransformData contains data for transform nodes
type TransformData struct {
	CommonData
	TransformType *string `json:"transform_type,omitempty"` // to_array, to_object, etc.
}

func (d TransformData) Validate() error {
	if d.TransformType == nil {
		return ErrMissingRequiredField("transform_type")
	}
	return nil
}

// AccumulatorData contains data for accumulator nodes
type AccumulatorData struct {
	CommonData
	AccumOp      *string     `json:"accum_op,omitempty"`      // sum, product, concat, etc.
	InitialValue interface{} `json:"initial_value,omitempty"` // Initial accumulator value
}

func (d AccumulatorData) Validate() error {
	if d.AccumOp == nil {
		return ErrMissingRequiredField("accum_op")
	}
	return nil
}

// CounterData contains data for counter nodes
type CounterData struct {
	CommonData
	CounterOp    *string     `json:"counter_op,omitempty"`    // increment, decrement, reset
	Delta        *float64    `json:"delta,omitempty"`         // Delta value for increment/decrement
	InitialValue interface{} `json:"initial_value,omitempty"` // Initial counter value
}

func (d CounterData) Validate() error {
	if d.CounterOp == nil {
		return ErrMissingRequiredField("counter_op")
	}
	return nil
}

// ParseData contains data for parse nodes
type ParseData struct {
	CommonData
	InputType *string `json:"input_type,omitempty"` // AUTO, JSON, CSV, TSV, YAML, XML
}

func (d ParseData) Validate() error {
	return nil // InputType is optional (AUTO detection)
}

// FormatData contains data for format nodes
type FormatData struct {
	CommonData
	OutputType     *string `json:"output_type,omitempty"`     // JSON, CSV, TSV
	PrettyPrint    *bool   `json:"pretty_print,omitempty"`    // For JSON
	IncludeHeaders *bool   `json:"include_headers,omitempty"` // For CSV/TSV
	Delimiter      *string `json:"delimiter,omitempty"`       // For CSV
}

func (d FormatData) Validate() error {
	if d.OutputType == nil {
		return ErrMissingRequiredField("output_type")
	}
	return nil
}

// ============================================================================
// Advanced Control Flow Node Data Types
// ============================================================================

// SwitchData contains data for switch nodes
type SwitchData struct {
	CommonData
	Cases []SwitchCase `json:"cases,omitempty"`
}

func (d SwitchData) Validate() error {
	if len(d.Cases) == 0 {
		return ErrMissingRequiredField("cases")
	}
	
	// Must have exactly one default case, and it must be last
	defaultCount := 0
	lastIndex := len(d.Cases) - 1
	
	for i, c := range d.Cases {
		if c.IsDefault {
			defaultCount++
			if i != lastIndex {
				return fmt.Errorf("default case must be the last case in the array (found at index %d, but last index is %d)", i, lastIndex)
			}
			// Default case doesn't need output_path
		} else {
			// Non-default cases must have output_path
			if c.OutputPath == nil || *c.OutputPath == "" {
				return fmt.Errorf("case at index %d must have output_path (non-default cases require output_path)", i)
			}
		}
	}
	
	if defaultCount == 0 {
		return fmt.Errorf("switch node must have exactly one default case (found 0)")
	}
	if defaultCount > 1 {
		return fmt.Errorf("switch node must have exactly one default case (found %d)", defaultCount)
	}
	
	return nil
}

// ParallelData contains data for parallel execution nodes
type ParallelData struct {
	CommonData
	MaxConcurrency *int    `json:"max_concurrency,omitempty"`
	Timeout        *string `json:"timeout,omitempty"`
}

func (d ParallelData) Validate() error {
	return nil // All fields are optional
}

// JoinData contains data for join nodes
type JoinData struct {
	CommonData
	JoinStrategy *string `json:"join_strategy,omitempty"` // all, any, first
	Timeout      *string `json:"timeout,omitempty"`
}

func (d JoinData) Validate() error {
	return nil // Fields are optional
}

// SplitData contains data for split nodes
type SplitData struct {
	CommonData
	Paths []string `json:"paths,omitempty"`
}

func (d SplitData) Validate() error {
	if len(d.Paths) == 0 {
		return ErrMissingRequiredField("paths")
	}
	return nil
}

// DelayData contains data for delay nodes
type DelayData struct {
	CommonData
	Duration *string `json:"duration,omitempty"`
}

func (d DelayData) Validate() error {
	if d.Duration == nil {
		return ErrMissingRequiredField("duration")
	}
	return nil
}

// CacheData contains data for cache nodes
type CacheData struct {
	CommonData
	CacheOp  *string `json:"cache_op,omitempty"` // get, set
	CacheKey *string `json:"cache_key,omitempty"`
	TTL      *string `json:"ttl,omitempty"`
	MaxSize  *int    `json:"max_size,omitempty"`
	Eviction *string `json:"eviction,omitempty"` // lru, lfu, ttl
	Storage  *string `json:"storage,omitempty"`  // memory, redis
}

func (d CacheData) Validate() error {
	if d.CacheOp == nil {
		return ErrMissingRequiredField("cache_op")
	}
	if d.CacheKey == nil {
		return ErrMissingRequiredField("cache_key")
	}
	return nil
}

// ============================================================================
// Error Handling & Resilience Node Data Types
// ============================================================================

// RetryData contains data for retry nodes
type RetryData struct {
	CommonData
	MaxAttempts     *int     `json:"max_attempts,omitempty"`
	BackoffStrategy *string  `json:"backoff_strategy,omitempty"` // exponential, linear, constant
	InitialDelay    *string  `json:"initial_delay,omitempty"`
	MaxDelay        *string  `json:"max_delay,omitempty"`
	Multiplier      *float64 `json:"multiplier,omitempty"`
	RetryOnErrors   []string `json:"retry_on_errors,omitempty"`
}

func (d RetryData) Validate() error {
	return nil // All fields are optional with defaults
}

// TryCatchData contains data for try-catch nodes
type TryCatchData struct {
	CommonData
	FallbackValue   interface{} `json:"fallback_value,omitempty"`
	ContinueOnError *bool       `json:"continue_on_error,omitempty"`
	ErrorOutputPath *string     `json:"error_output_path,omitempty"`
}

func (d TryCatchData) Validate() error {
	return nil // All fields are optional
}

// TimeoutData contains data for timeout nodes
type TimeoutData struct {
	CommonData
	Timeout       *string `json:"timeout,omitempty"`
	TimeoutAction *string `json:"timeout_action,omitempty"` // error, continue_with_partial
}

func (d TimeoutData) Validate() error {
	if d.Timeout == nil {
		return ErrMissingRequiredField("timeout")
	}
	return nil
}

// ============================================================================
// Advanced Node Data Types (Phase 4)
// ============================================================================

// RateLimiterData contains data for rate limiter nodes
type RateLimiterData struct {
	CommonData
	MaxRequests       *int    `json:"max_requests,omitempty"`
	PerDuration       *string `json:"per_duration,omitempty"` // 1s, 1m, 1h
	RateLimitStrategy *string `json:"strategy,omitempty"`     // fixed_window, sliding_window, token_bucket
}

func (d RateLimiterData) Validate() error {
	if d.MaxRequests == nil {
		return ErrMissingRequiredField("max_requests")
	}
	if d.PerDuration == nil {
		return ErrMissingRequiredField("per_duration")
	}
	return nil
}

// ThrottleData contains data for throttle nodes
type ThrottleData struct {
	CommonData
	RequestsPerSecond *float64 `json:"requests_per_second,omitempty"`
}

func (d ThrottleData) Validate() error {
	if d.RequestsPerSecond == nil {
		return ErrMissingRequiredField("requests_per_second")
	}
	return nil
}

// SchemaValidatorData contains data for schema validator nodes
type SchemaValidatorData struct {
	CommonData
	Schema interface{} `json:"schema,omitempty"` // JSON schema
	Strict *bool       `json:"strict,omitempty"`
}

func (d SchemaValidatorData) Validate() error {
	if d.Schema == nil {
		return ErrMissingRequiredField("schema")
	}
	return nil
}

// PaginatorData contains data for paginator nodes
type PaginatorData struct {
	CommonData
	PaginationStrategy *string `json:"pagination_strategy,omitempty"` // offset_limit, page_number, cursor, link_header
	OffsetParam        *string `json:"offset_param,omitempty"`
	LimitParam         *string `json:"limit_param,omitempty"`
	PageSize           *int    `json:"page_size,omitempty"`
	MaxPages           *int    `json:"max_pages,omitempty"`
	PageParam          *string `json:"page_param,omitempty"`
	PerPageParam       *string `json:"per_page_param,omitempty"`
	CursorParam        *string `json:"cursor_param,omitempty"`
	NextCursorPath     *string `json:"next_cursor_path,omitempty"`
	LinkHeader         *string `json:"link_header,omitempty"`
	TotalCountPath     *string `json:"total_count_path,omitempty"`
	ResultsPath        *string `json:"results_path,omitempty"`
}

func (d PaginatorData) Validate() error {
	if d.PaginationStrategy == nil {
		return ErrMissingRequiredField("pagination_strategy")
	}
	return nil
}

// ============================================================================
// Context Node Data Types
// ============================================================================

// ContextVariableData contains data for context variable nodes
type ContextVariableData struct {
	CommonData
	ContextName   *string                `json:"context_name,omitempty"`   // DEPRECATED
	ContextValue  interface{}            `json:"context_value,omitempty"`  // DEPRECATED
	ContextValues []ContextVariableValue `json:"context_values,omitempty"` // Preferred
}

func (d ContextVariableData) Validate() error {
	if len(d.ContextValues) == 0 && d.ContextName == nil {
		return ErrMissingRequiredField("context_values or context_name")
	}
	return nil
}

// ContextConstantData contains data for context constant nodes
type ContextConstantData struct {
	CommonData
	ContextName   *string                `json:"context_name,omitempty"`   // DEPRECATED
	ContextValue  interface{}            `json:"context_value,omitempty"`  // DEPRECATED
	ContextValues []ContextVariableValue `json:"context_values,omitempty"` // Preferred
}

func (d ContextConstantData) Validate() error {
	if len(d.ContextValues) == 0 && d.ContextName == nil {
		return ErrMissingRequiredField("context_values or context_name")
	}
	return nil
}

// ============================================================================
// Visualization Node Data Types
// ============================================================================

// VisualizationData contains data for visualization nodes
type VisualizationData struct {
	CommonData
	Mode *string `json:"mode,omitempty"` // table, chart, json, etc.
}

func (d VisualizationData) Validate() error {
	return nil // Mode is optional
}

// RendererData contains data for renderer nodes
type RendererData struct {
	CommonData
	Mode *string `json:"mode,omitempty"` // Render mode
}

func (d RendererData) Validate() error {
	return nil // Mode is optional
}

// ============================================================================
// Custom Executor Data Type
// ============================================================================

// CustomExecutorData is a generic container for custom node executors
// It allows any fields to support extensibility
type CustomExecutorData struct {
	CommonData
	// Additional fields are stored as a map for maximum flexibility
	Fields map[string]interface{} `json:"-"`
}

func (d CustomExecutorData) Validate() error {
	return nil // Custom executors handle their own validation
}
