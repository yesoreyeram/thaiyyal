// Package engine provides the workflow execution engine.
// This orchestrates workflow parsing, validation, and execution using the refactored packages.
package engine

import (
"context"
"crypto/rand"
"encoding/hex"
"encoding/json"
"fmt"
"regexp"
"sync"
"time"

"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
"github.com/yesoreyeram/thaiyyal/backend/pkg/graph"
"github.com/yesoreyeram/thaiyyal/backend/pkg/logging"
"github.com/yesoreyeram/thaiyyal/backend/pkg/observer"
"github.com/yesoreyeram/thaiyyal/backend/pkg/state"
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ============================================================================
// Engine Definition
// ============================================================================

// Engine is the workflow execution engine.
// It manages workflow state and coordinates node execution in topological order.
//
// The Engine uses the following design patterns:
//   - Strategy Pattern: Different execution strategies for different node types (via Registry)
//   - State Pattern: Manages workflow state (variables, accumulator, counter, cache)
//   - Template Method: Execute() defines the workflow execution algorithm
//   - Observer Pattern: Notifies observers of execution events (optional)
type Engine struct {
	graph      *graph.Graph
	state      *state.Manager
	registry   *executor.Registry
	config     types.Config
	results    map[string]interface{}
	resultsMu  sync.RWMutex
	executionID string
	workflowID  string

	// Runtime protection counters
	nodeExecutionCount int
	httpCallCount      int
	countersMu         sync.RWMutex

	// Node storage for lookups
	nodes []types.Node
	edges []types.Edge

	// Observer support
	observerMgr *observer.Manager
	logger      observer.Logger
	
	// Structured logging
	structuredLogger *logging.Logger
	
	// HTTP client registry for named HTTP clients (uses standalone httpclient.Registry)
	httpClientRegistry interface{}
}

// ============================================================================
// Constructor Functions
// ============================================================================

// New creates a new workflow engine from JSON payload.
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
func New(payloadJSON []byte) (*Engine, error) {
return NewWithConfig(payloadJSON, types.DefaultConfig())
}

// NewWithConfig creates a new workflow engine with custom configuration.
// This is useful for testing or when you need non-default security settings.
//
// An execution ID is automatically generated for this execution.
func NewWithConfig(payloadJSON []byte, config types.Config) (*Engine, error) {
	return NewWithRegistry(payloadJSON, config, DefaultRegistry())
}

// NewWithRegistry creates a new workflow engine with a custom executor registry.
// This allows users to register custom node executors while maintaining all
// security protections and workflow execution capabilities.
//
// Example usage:
//
//	// Start with default registry and add custom nodes
//	registry := engine.DefaultRegistry()
//	registry.MustRegister(&MyCustomExecutor{})
//	engine, err := engine.NewWithRegistry(payload, config, registry)
//
//	// Or create a completely custom registry
//	registry := executor.NewRegistry()
//	registry.MustRegister(&MyExecutor{})
//	engine, err := engine.NewWithRegistry(payload, config, registry)
//
// Security Considerations:
//   - Custom executors must implement NodeExecutor interface properly
//   - All protection limits (MaxNodeExecutions, MaxHTTPCallsPerExec, etc.) apply to custom nodes
//   - Custom executors should call ctx.IncrementNodeExecution() if they perform iterations
//   - Custom executors making HTTP calls should call ctx.IncrementHTTPCall()
//   - Custom executors should validate all inputs and handle errors appropriately
//
// An execution ID is automatically generated for this execution.
func NewWithRegistry(payloadJSON []byte, config types.Config, registry *executor.Registry) (*Engine, error) {
	if registry == nil {
		return nil, fmt.Errorf("registry cannot be nil")
	}

	var payload types.Payload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	// Generate execution ID
	executionID := generateExecutionID()

	// Create structured logger with workflow and execution context
	structuredLogger := logging.New(logging.DefaultConfig()).
		WithWorkflowID(payload.WorkflowID).
		WithExecutionID(executionID)

	engine := &Engine{
		state:            state.New(),
		registry:         registry,
		config:           config,
		results:          make(map[string]interface{}),
		executionID:      executionID,
		workflowID:       payload.WorkflowID,
		nodes:            payload.Nodes,
		edges:            payload.Edges,
		observerMgr:      observer.NewManager(),
		logger:           &observer.NoOpLogger{},
		structuredLogger: structuredLogger,
	}

	// Create graph for topological sorting
	engine.graph = graph.New(payload.Nodes, payload.Edges)

	return engine, nil
}

// DefaultRegistry creates and populates the default executor registry with all built-in node executors.
// This function is exported to allow users to start with the default set and add custom executors.
//
// Example usage:
//
//	// Get default registry and add custom nodes
//	registry := engine.DefaultRegistry()
//	registry.MustRegister(&MyCustomExecutor{})
//	engine, err := engine.NewWithRegistry(payload, config, registry)
//
// Returns a registry with all 40 built-in node types registered:
//   - Basic I/O: Number, TextInput, Visualization
//   - Operations: Operation, TextOperation, HTTP
//   - Control Flow: Condition, ForEach, WhileLoop, Filter, Map, Reduce
//   - Array Processing: Slice, Sort, Find, FlatMap, GroupBy, Unique, Chunk, Reverse, Partition, Zip, Sample, Range, Compact, Transpose
//   - State & Memory: Variable, Extract, Transform, Accumulator, Counter
//   - Advanced Control: Switch, Parallel, Join, Split, Delay, Cache
//   - Error Handling: Retry, TryCatch, Timeout
//   - Context: ContextVariable, ContextConstant
func DefaultRegistry() *executor.Registry {
	reg := executor.NewRegistry()

	// Register all 40 node type executors
	// Basic I/O nodes
	reg.MustRegister(&executor.NumberExecutor{})
	reg.MustRegister(&executor.TextInputExecutor{})
	reg.MustRegister(&executor.VisualizationExecutor{})

	// Operation nodes
	reg.MustRegister(&executor.OperationExecutor{})
	reg.MustRegister(&executor.TextOperationExecutor{})
	reg.MustRegister(executor.NewHTTPExecutor())

	// Control flow nodes
	reg.MustRegister(&executor.ConditionExecutor{})
	reg.MustRegister(&executor.ForEachExecutor{})
	reg.MustRegister(&executor.WhileLoopExecutor{})
	reg.MustRegister(&executor.FilterExecutor{})
	reg.MustRegister(&executor.MapExecutor{})
	reg.MustRegister(&executor.ReduceExecutor{})

	// Array processing nodes (14 new nodes)
	reg.MustRegister(&executor.SliceExecutor{})      // Pagination, windowing
	reg.MustRegister(&executor.SortExecutor{})       // Sort by field
	reg.MustRegister(&executor.FindExecutor{})       // Find first match
	reg.MustRegister(&executor.FlatMapExecutor{})    // Transform and flatten
	reg.MustRegister(&executor.GroupByExecutor{})    // Group and aggregate
	reg.MustRegister(&executor.UniqueExecutor{})     // Remove duplicates
	reg.MustRegister(&executor.ChunkExecutor{})      // Split into chunks
	reg.MustRegister(&executor.ReverseExecutor{})    // Reverse array
	reg.MustRegister(&executor.PartitionExecutor{})  // Split by condition
	reg.MustRegister(&executor.ZipExecutor{})        // Combine arrays
	reg.MustRegister(&executor.SampleExecutor{})     // Random sampling
	reg.MustRegister(&executor.RangeExecutor{})      // Generate sequences
	reg.MustRegister(&executor.CompactExecutor{})    // Remove null/empty
	reg.MustRegister(&executor.TransposeExecutor{})  // Transpose matrix

	// State & memory nodes
	reg.MustRegister(&executor.VariableExecutor{})
	reg.MustRegister(&executor.ExtractExecutor{})
	reg.MustRegister(&executor.TransformExecutor{})
	reg.MustRegister(&executor.AccumulatorExecutor{})
	reg.MustRegister(&executor.CounterExecutor{})

	// Advanced control flow nodes
	reg.MustRegister(&executor.SwitchExecutor{})
	reg.MustRegister(&executor.ParallelExecutor{})
	reg.MustRegister(&executor.JoinExecutor{})
	reg.MustRegister(&executor.SplitExecutor{})
	reg.MustRegister(&executor.DelayExecutor{})
	reg.MustRegister(&executor.CacheExecutor{})

	// Error handling & resilience nodes
	reg.MustRegister(&executor.RetryExecutor{})
	reg.MustRegister(&executor.TryCatchExecutor{})
	reg.MustRegister(&executor.TimeoutExecutor{})

	// Context nodes
	reg.MustRegister(&executor.ContextVariableExecutor{})
	reg.MustRegister(&executor.ContextConstantExecutor{})

	return reg
}

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

// ============================================================================
// Observer and Logger Configuration
// ============================================================================

// RegisterObserver adds an observer to receive execution events.
// Multiple observers can be registered and will all receive events.
// Returns the engine for method chaining.
func (e *Engine) RegisterObserver(obs observer.Observer) *Engine {
	if obs != nil {
		e.observerMgr.Register(obs)
	}
	return e
}

// SetLogger sets the logger for the engine.
// If no logger is set, a NoOpLogger is used by default.
// Returns the engine for method chaining.
func (e *Engine) SetLogger(logger observer.Logger) *Engine {
	if logger != nil {
		e.logger = logger
	}
	return e
}

// SetHTTPClientRegistry sets the HTTP client registry for named HTTP clients.
// The registry should be of type *httpclient.Registry from the standalone httpclient package.
// Returns the engine for method chaining.
func (e *Engine) SetHTTPClientRegistry(registry interface{}) *Engine {
	e.httpClientRegistry = registry
	return e
}

// GetObserverCount returns the number of registered observers
func (e *Engine) GetObserverCount() int {
	return e.observerMgr.Count()
}

// ============================================================================
// Public API - Execute
// ============================================================================

// Execute runs the workflow and returns the result.
// The workflow execution is protected by a timeout configured in MaxExecutionTime.
// If the workflow takes longer than the timeout, execution is cancelled and an error is returned.
//
// Each workflow execution is assigned a unique execution ID that is passed through the
// execution context and included in the result. This ID can be used for logging and tracing.
//
// Observers will be notified of workflow and node execution events if registered.
//
// Returns:
//   - *types.Result: Workflow execution results including execution ID, node outputs and final output
//   - error: If execution fails, times out, or encounters an error
func (e *Engine) Execute() (*types.Result, error) {
	workflowStartTime := time.Now()
	
	// Log workflow execution start
	e.structuredLogger.Info("workflow execution started")
	
	result := &types.Result{
		ExecutionID: e.executionID,
		WorkflowID:  e.workflowID,
		NodeResults: make(map[string]interface{}),
		Errors:      []string{},
	}

	// Step 1: Infer node types if not set
	e.inferNodeTypes()

	// Step 2: Get execution order using topological sort
	executionOrder, err := e.graph.TopologicalSort()
	if err != nil {
		e.structuredLogger.WithError(err).Error("topological sort failed")
		return result, err
	}
	
	e.structuredLogger.
		WithField("execution_order", executionOrder).
		WithField("node_count", len(executionOrder)).
		Debug("execution order determined")

	// Step 3: Create context with timeout and execution metadata for workflow execution
	ctx, cancel := context.WithTimeout(context.Background(), e.config.MaxExecutionTime)
	defer cancel()

	// Add execution ID and workflow ID to context for logging and tracing
	ctx = context.WithValue(ctx, types.ContextKeyExecutionID, e.executionID)
	ctx = context.WithValue(ctx, types.ContextKeyWorkflowID, e.workflowID)

	// Notify observers: Workflow start
	e.notifyWorkflowStart(ctx, workflowStartTime)

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
			value, err := e.executeNode(ctx, node)
			if err != nil {
				errMsg := fmt.Sprintf("error executing node %s: %v", nodeID, err)
				result.Errors = append(result.Errors, errMsg)
				e.structuredLogger.WithNodeID(nodeID).WithError(err).Error("node execution failed")
				done <- fmt.Errorf("%s", errMsg)
				return
			}
			e.SetNodeResult(nodeID, value)
		}
		done <- nil
	}()

	// Wait for execution to complete or timeout
	select {
	case err := <-done:
		if err != nil {
			e.structuredLogger.WithError(err).Error("workflow execution failed")
			// Notify observers: Workflow end with error
			e.notifyWorkflowEnd(ctx, workflowStartTime, nil, err)
			return result, err
		}
	case <-ctx.Done():
		timeoutErr := fmt.Errorf("workflow execution timeout: exceeded %v", e.config.MaxExecutionTime)
		e.structuredLogger.WithField("timeout", e.config.MaxExecutionTime).Error("workflow execution timeout")
		// Notify observers: Workflow end with timeout
		e.notifyWorkflowEnd(ctx, workflowStartTime, nil, timeoutErr)
		return result, timeoutErr
	}

	// Step 4: Copy results and set final output
	result.NodeResults = e.results
	result.FinalOutput = e.getFinalOutput()
	
	e.structuredLogger.
		WithField("duration_ms", time.Since(workflowStartTime).Milliseconds()).
		WithField("nodes_executed", len(executionOrder)).
		Info("workflow execution completed successfully")

	// Notify observers: Workflow end (success)
	e.notifyWorkflowEnd(ctx, workflowStartTime, result.FinalOutput, nil)

	return result, nil
}

// ============================================================================
// Node Execution
// ============================================================================

// executeNode dispatches node execution to the appropriate executor via the registry.
// Handles template interpolation before execution (except for context nodes).
// Notifies observers of node execution events.
//
// Parameters:
//   - ctx: Context with execution metadata (execution ID, workflow ID)
//   - node: Node to execute
//
// Returns:
//   - interface{}: Result of node execution (type depends on node)
//   - error: If node execution fails
func (e *Engine) executeNode(ctx context.Context, node types.Node) (interface{}, error) {
	nodeStartTime := time.Now()
	
	// Create node-specific logger
	nodeLogger := e.structuredLogger.
		WithNodeID(node.ID).
		WithNodeType(node.Type)
	
	nodeLogger.Debug("node execution started")
	
	// Notify observers: Node start
	e.notifyNodeStart(ctx, node, nodeStartTime)

	// Check and increment node execution counter
	if err := e.IncrementNodeExecution(); err != nil {
		nodeLogger.WithError(err).Error("node execution counter limit exceeded")
		// This is a protection limit error, not a node execution failure
		// Still notify as failure since execution couldn't proceed
		e.notifyNodeFailure(ctx, node, nodeStartTime, nil, err)
		return nil, err
	}

	// Interpolate templates in node data before execution (except for context nodes)
	if node.Type != types.NodeTypeContextVariable && node.Type != types.NodeTypeContextConstant {
		e.interpolateNodeData(&node.Data)
	}

	// Dispatch to appropriate executor via registry
	result, err := e.registry.Execute(e, node)
	
	if err != nil {
		nodeLogger.WithError(err).Error("node execution failed")
		// Notify observers: Node execution failure
		e.notifyNodeFailure(ctx, node, nodeStartTime, result, err)
		return nil, err
	}
	
	nodeLogger.
		WithField("duration_ms", time.Since(nodeStartTime).Milliseconds()).
		Info("node execution completed successfully")
	
	// Notify observers: Node success
	e.notifyNodeSuccess(ctx, node, nodeStartTime, result)
	
	return result, nil
}

// ============================================================================
// Type Inference
// ============================================================================

// inferNodeTypes determines node types from data if not explicitly set.
// This allows the frontend to omit node types and have them automatically detected.
//
// Type inference is based on the presence of specific fields in NodeData.
// Some nodes (for_each, while_loop, parallel) require explicit types as they
// have ambiguous fields.
func (e *Engine) inferNodeTypes() {
for i := range e.nodes {
if e.nodes[i].Type != "" {
// Type already set, skip inference
continue
}

// Infer type from data fields
e.nodes[i].Type = inferNodeTypeFromData(e.nodes[i].Data)
}
}

// inferNodeTypeFromData infers a node's type from its data fields.
// This implements a simple decision tree based on field presence.
//
// Returns:
//   - types.NodeType: Inferred type, or empty string if cannot infer
func inferNodeTypeFromData(data types.NodeData) types.NodeType {
// Basic I/O nodes (checked first as they're most common)
if data.Value != nil {
return types.NodeTypeNumber
}
if data.Text != nil {
return types.NodeTypeTextInput
}
if data.Mode != nil {
return types.NodeTypeVisualization
}

// Operation nodes
if data.Op != nil {
return types.NodeTypeOperation
}
if data.TextOp != nil {
return types.NodeTypeTextOperation
}
if data.URL != nil {
return types.NodeTypeHTTP
}

// Control flow nodes
if data.Condition != nil {
return types.NodeTypeCondition
}

// State & memory nodes
if data.VarName != nil && data.VarOp != nil {
return types.NodeTypeVariable
}
if data.Field != nil || len(data.Fields) > 0 {
return types.NodeTypeExtract
}
if data.TransformType != nil {
return types.NodeTypeTransform
}
if data.AccumOp != nil {
return types.NodeTypeAccumulator
}
if data.CounterOp != nil {
return types.NodeTypeCounter
}

// Advanced control flow nodes
if len(data.Cases) > 0 {
return types.NodeTypeSwitch
}
if data.JoinStrategy != nil {
return types.NodeTypeJoin
}
if len(data.Paths) > 0 {
return types.NodeTypeSplit
}
if data.Duration != nil {
return types.NodeTypeDelay
}
if data.CacheOp != nil && data.CacheKey != nil {
return types.NodeTypeCache
}

// Context nodes
if data.ContextName != nil && data.ContextValue != nil {
// Default to variable, frontend should specify explicitly
// This is a best-effort inference
return types.NodeTypeContextVariable
}

// Error handling & resilience nodes
if data.MaxAttempts != nil || data.BackoffStrategy != nil {
return types.NodeTypeRetry
}
if data.FallbackValue != nil || data.ContinueOnError != nil {
return types.NodeTypeTryCatch
}
if data.Timeout != nil && data.TimeoutAction != nil {
return types.NodeTypeTimeout
}

// Cannot infer type
// Note: for_each, while_loop, and parallel require explicit type
// as they have ambiguous fields
return ""
}

// ============================================================================
// ExecutionContext Interface Implementation
// ============================================================================

// GetNodeInputs retrieves all input values for a node from its predecessor nodes.
func (e *Engine) GetNodeInputs(nodeID string) []interface{} {
inputs := []interface{}{}
e.resultsMu.RLock()
defer e.resultsMu.RUnlock()

for _, edge := range e.edges {
if edge.Target == nodeID {
if result, ok := e.results[edge.Source]; ok {
inputs = append(inputs, result)
}
}
}
return inputs
}

// GetNode retrieves a node by its ID
func (e *Engine) GetNode(nodeID string) *types.Node {
for i := range e.nodes {
if e.nodes[i].ID == nodeID {
return &e.nodes[i]
}
}
return nil
}

// GetVariable retrieves a variable value
func (e *Engine) GetVariable(name string) (interface{}, error) {
	return e.state.GetVariable(name)
}

// SetVariable sets a variable value with validation
func (e *Engine) SetVariable(name string, value interface{}) error {
	// Validate value against resource limits
	if err := types.ValidateValue(value, e.config); err != nil {
		return fmt.Errorf("variable validation failed: %w", err)
	}

	// Check variable count limit (only for new variables)
	if e.config.MaxVariables > 0 {
		_, err := e.state.GetVariable(name)
		if err != nil {
			// Variable doesn't exist, check count limit
			vars := e.state.GetAllVariables()
			if len(vars) >= e.config.MaxVariables {
				return fmt.Errorf("maximum variables exceeded: %d (limit: %d)", len(vars), e.config.MaxVariables)
			}
		}
	}

	return e.state.SetVariable(name, value)
}

// GetAccumulator returns the current accumulator value
func (e *Engine) GetAccumulator() interface{} {
return e.state.GetAccumulator()
}

// SetAccumulator sets the accumulator value
func (e *Engine) SetAccumulator(value interface{}) {
e.state.SetAccumulator(value)
}

// GetCounter returns the current counter value
func (e *Engine) GetCounter() float64 {
return e.state.GetCounter()
}

// SetCounter sets the counter value
func (e *Engine) SetCounter(value float64) {
e.state.SetCounter(value)
}

// GetCache retrieves a cached value
func (e *Engine) GetCache(key string) (interface{}, bool) {
val, found, _ := e.state.GetCache(key)
return val, found
}

// SetCache sets a cached value with TTL
func (e *Engine) SetCache(key string, value interface{}, ttl time.Duration) {
e.state.SetCache(key, value, ttl)
}

// GetWorkflowContext returns all context variables and constants
func (e *Engine) GetWorkflowContext() map[string]interface{} {
return e.state.GetAllContext()
}

// GetContextVariable retrieves a context variable
func (e *Engine) GetContextVariable(name string) (interface{}, bool) {
	return e.state.GetContextVariable(name)
}

// SetContextVariable sets a context variable
// Note: Validation failures are logged and ignored to maintain backward compatibility.
// Context variables are typically set during workflow initialization.
func (e *Engine) SetContextVariable(name string, value interface{}) {
	// Validate value against resource limits (best effort)
	if err := types.ValidateValue(value, e.config); err != nil {
		// Log validation error but continue to maintain backward compatibility
		e.structuredLogger.
			WithField("variable_name", name).
			WithError(err).
			Warn("context variable validation failed, storing anyway")
		// For now, validation errors are logged but ignored to avoid breaking
		// workflow initialization with large context values
	}
	e.state.SetContextVariable(name, value)
}

// GetContextConstant retrieves a context constant
func (e *Engine) GetContextConstant(name string) (interface{}, bool) {
	return e.state.GetContextConstant(name)
}

// SetContextConstant sets a context constant
// Note: Validation failures are logged and ignored to maintain backward compatibility.
// Context constants are typically set during workflow initialization.
func (e *Engine) SetContextConstant(name string, value interface{}) {
	// Validate value against resource limits (best effort)
	if err := types.ValidateValue(value, e.config); err != nil {
		// Log validation error but continue to maintain backward compatibility
		e.structuredLogger.
			WithField("constant_name", name).
			WithError(err).
			Warn("context constant validation failed, storing anyway")
		// For now, validation errors are logged but ignored to avoid breaking
		// workflow initialization with large context values
	}
	e.state.SetContextConstant(name, value)
}

// InterpolateTemplate replaces template placeholders in a string with actual values from context
func (e *Engine) InterpolateTemplate(template string) string {
return e.interpolateTemplate(template)
}

// GetNodeResult retrieves a node's execution result
func (e *Engine) GetNodeResult(nodeID string) (interface{}, bool) {
e.resultsMu.RLock()
defer e.resultsMu.RUnlock()

result, ok := e.results[nodeID]
return result, ok
}

// SetNodeResult stores a node's execution result
// Note: Validation is best-effort to avoid breaking valid executions.
// Results that exceed limits may still be stored but could cause issues downstream.
func (e *Engine) SetNodeResult(nodeID string, result interface{}) {
	// Validate result against resource limits (best effort)
	// We don't fail here to avoid breaking workflows that produce large intermediate results
	if err := types.ValidateValue(result, e.config); err != nil {
		// Log validation warning but store the result to maintain workflow execution
		e.structuredLogger.
			WithNodeID(nodeID).
			WithError(err).
			Warn("node result validation failed, storing anyway")
		// For now, the validation error is logged but the result is still stored
		// to maintain backward compatibility and avoid breaking workflows
	}

	e.resultsMu.Lock()
	defer e.resultsMu.Unlock()

	e.results[nodeID] = result
}

// GetAllNodeResults returns all node execution results
func (e *Engine) GetAllNodeResults() map[string]interface{} {
e.resultsMu.RLock()
defer e.resultsMu.RUnlock()

// Return a copy to avoid concurrent modification
resultsCopy := make(map[string]interface{}, len(e.results))
for k, v := range e.results {
resultsCopy[k] = v
}
return resultsCopy
}

// GetVariables returns all workflow variables
func (e *Engine) GetVariables() map[string]interface{} {
return e.state.GetAllVariables()
}

// GetContextVariables returns all context variables and constants
func (e *Engine) GetContextVariables() map[string]interface{} {
return e.state.GetAllContext()
}

// GetConfig returns the engine configuration
func (e *Engine) GetConfig() types.Config {
	return e.config
}

// GetHTTPClientRegistry returns the HTTP client registry if configured.
// Returns nil if no registry is set. The caller should type assert to
// *httpclient.Registry if needed.
func (e *Engine) GetHTTPClientRegistry() interface{} {
	return e.httpClientRegistry
}

// IncrementNodeExecution increments the node execution counter and checks limits.
// Returns an error if the limit is exceeded.
func (e *Engine) IncrementNodeExecution() error {
	e.countersMu.Lock()
	defer e.countersMu.Unlock()
	
	e.nodeExecutionCount++
	
	// Check if limit is configured and enforced (0 means unlimited)
	if e.config.MaxNodeExecutions > 0 && e.nodeExecutionCount > e.config.MaxNodeExecutions {
		return fmt.Errorf("maximum node executions exceeded: %d (limit: %d)", e.nodeExecutionCount, e.config.MaxNodeExecutions)
	}
	
	return nil
}

// IncrementHTTPCall increments the HTTP call counter and checks limits.
// Returns an error if the limit is exceeded.
func (e *Engine) IncrementHTTPCall() error {
	e.countersMu.Lock()
	defer e.countersMu.Unlock()
	
	e.httpCallCount++
	
	// Check if limit is configured and enforced (0 means unlimited)
	if e.config.MaxHTTPCallsPerExec > 0 && e.httpCallCount > e.config.MaxHTTPCallsPerExec {
		return fmt.Errorf("maximum HTTP calls per execution exceeded: %d (limit: %d)", e.httpCallCount, e.config.MaxHTTPCallsPerExec)
	}
	
	return nil
}

// GetNodeExecutionCount returns the current node execution count
func (e *Engine) GetNodeExecutionCount() int {
	e.countersMu.RLock()
	defer e.countersMu.RUnlock()
	return e.nodeExecutionCount
}

// GetHTTPCallCount returns the current HTTP call count
func (e *Engine) GetHTTPCallCount() int {
	e.countersMu.RLock()
	defer e.countersMu.RUnlock()
	return e.httpCallCount
}

// ============================================================================
// Template Interpolation
// ============================================================================

// templateRegex matches {{ variable.name }} or {{ const.name }}
var templateRegex = regexp.MustCompile(`\{\{\s*(variable|const)\.(\w+)\s*\}\}`)

// interpolateTemplate replaces template placeholders in a string with actual values from context
func (e *Engine) interpolateTemplate(text string) string {
// Check if we have any context to interpolate
contextVars := e.state.GetAllContext()
if len(contextVars) == 0 {
return text
}

// Replace all template placeholders
result := templateRegex.ReplaceAllStringFunc(text, func(match string) string {
// Extract the type and name from the match
parts := templateRegex.FindStringSubmatch(match)
if len(parts) != 3 {
return match // Return original if parsing fails
}

contextType := parts[1]
varName := parts[2]

// Look up the value in the appropriate context map
var value interface{}
var exists bool

if contextType == "variable" {
value, exists = e.state.GetContextVariable(varName)
} else if contextType == "const" {
value, exists = e.state.GetContextConstant(varName)
}

if exists {
return fmt.Sprintf("%v", value)
}

// Return original if not found
return match
})

return result
}

// interpolateValue recursively interpolates templates in various data types
func (e *Engine) interpolateValue(value interface{}) interface{} {
switch v := value.(type) {
case string:
return e.interpolateTemplate(v)
case map[string]interface{}:
result := make(map[string]interface{})
for key, val := range v {
result[key] = e.interpolateValue(val)
}
return result
case []interface{}:
result := make([]interface{}, len(v))
for i, val := range v {
result[i] = e.interpolateValue(val)
}
return result
default:
return value
}
}

// interpolateNodeData interpolates all string fields in NodeData
func (e *Engine) interpolateNodeData(data *types.NodeData) {
// Check if we have any context to interpolate
contextVars := e.state.GetAllContext()
if len(contextVars) == 0 {
return
}

// Interpolate string pointer fields
if data.Text != nil {
interpolated := e.interpolateTemplate(*data.Text)
data.Text = &interpolated
}
if data.URL != nil {
interpolated := e.interpolateTemplate(*data.URL)
data.URL = &interpolated
}
if data.Label != nil {
interpolated := e.interpolateTemplate(*data.Label)
data.Label = &interpolated
}
if data.VarName != nil {
interpolated := e.interpolateTemplate(*data.VarName)
data.VarName = &interpolated
}
if data.Field != nil {
interpolated := e.interpolateTemplate(*data.Field)
data.Field = &interpolated
}
if data.CacheKey != nil {
interpolated := e.interpolateTemplate(*data.CacheKey)
data.CacheKey = &interpolated
}

// Interpolate string arrays
if len(data.Fields) > 0 {
for i, field := range data.Fields {
data.Fields[i] = e.interpolateTemplate(field)
}
}
if len(data.Paths) > 0 {
for i, path := range data.Paths {
data.Paths[i] = e.interpolateTemplate(path)
}
}

// Interpolate interface{} fields that might contain strings
if data.InitialValue != nil {
data.InitialValue = e.interpolateValue(data.InitialValue)
}
if data.FallbackValue != nil {
data.FallbackValue = e.interpolateValue(data.FallbackValue)
}
}

// ============================================================================
// Helper Methods
// ============================================================================

// getNode retrieves a node by its ID (internal helper)
func (e *Engine) getNode(nodeID string) types.Node {
for _, node := range e.nodes {
if node.ID == nodeID {
return node
}
}
return types.Node{}
}

// getFinalOutput determines the final output of the workflow.
// The final output is the result of a terminal node (node with no outgoing edges).
// Context nodes (context_variable, context_constant) are excluded from being final output
// unless they are the ONLY nodes in the workflow.
//
// If multiple terminal nodes exist, returns the first non-context one found.
// If no terminal nodes exist (all nodes have outgoing edges), returns nil.
//
// Returns:
//   - interface{}: The result value from a terminal node, or nil if none found
func (e *Engine) getFinalOutput() interface{} {
// Build a set of all terminal nodes (nodes with no outgoing edges)
terminalNodes := make(map[string]bool)

// Initially, all nodes are considered terminal
for _, node := range e.nodes {
terminalNodes[node.ID] = true
}

// Remove nodes that have outgoing edges
for _, edge := range e.edges {
terminalNodes[edge.Source] = false
}

e.resultsMu.RLock()
defer e.resultsMu.RUnlock()

// First pass: Try to find a non-context terminal node
for nodeID, isTerminal := range terminalNodes {
if isTerminal {
node := e.getNode(nodeID)
if node.Type != types.NodeTypeContextVariable && node.Type != types.NodeTypeContextConstant {
if result, ok := e.results[nodeID]; ok {
return result
}
}
}
}

// Second pass: If no non-context terminal found, return any terminal (including context)
// This handles the case where workflow contains only context nodes
for nodeID, isTerminal := range terminalNodes {
if isTerminal {
if result, ok := e.results[nodeID]; ok {
return result
}
}
}

// No terminal node found
return nil
}

// ============================================================================
// Observer Notification Helpers
// ============================================================================

// notifyWorkflowStart notifies observers that workflow execution has started
func (e *Engine) notifyWorkflowStart(ctx context.Context, startTime time.Time) {
	if !e.observerMgr.HasObservers() {
		return
	}

	event := observer.Event{
		Type:        observer.EventWorkflowStart,
		Status:      observer.StatusStarted,
		Timestamp:   startTime,
		ExecutionID: e.executionID,
		WorkflowID:  e.workflowID,
		StartTime:   startTime,
	}

	e.observerMgr.Notify(ctx, event)
}

// notifyWorkflowEnd notifies observers that workflow execution has ended
func (e *Engine) notifyWorkflowEnd(ctx context.Context, startTime time.Time, result interface{}, err error) {
	if !e.observerMgr.HasObservers() {
		return
	}

	status := observer.StatusSuccess
	if err != nil {
		status = observer.StatusFailure
	}

	event := observer.Event{
		Type:        observer.EventWorkflowEnd,
		Status:      status,
		Timestamp:   time.Now(),
		ExecutionID: e.executionID,
		WorkflowID:  e.workflowID,
		StartTime:   startTime,
		ElapsedTime: time.Since(startTime),
		Result:      result,
		Error:       err,
	}

	e.observerMgr.Notify(ctx, event)
}

// notifyNodeStart notifies observers that a node execution has started
func (e *Engine) notifyNodeStart(ctx context.Context, node types.Node, startTime time.Time) {
	if !e.observerMgr.HasObservers() {
		return
	}

	event := observer.Event{
		Type:        observer.EventNodeStart,
		Status:      observer.StatusStarted,
		Timestamp:   startTime,
		ExecutionID: e.executionID,
		WorkflowID:  e.workflowID,
		NodeID:      node.ID,
		NodeType:    node.Type,
		StartTime:   startTime,
	}

	e.observerMgr.Notify(ctx, event)
}

// notifyNodeSuccess notifies observers that a node execution succeeded
func (e *Engine) notifyNodeSuccess(ctx context.Context, node types.Node, startTime time.Time, result interface{}) {
	if !e.observerMgr.HasObservers() {
		return
	}

	event := observer.Event{
		Type:        observer.EventNodeSuccess,
		Status:      observer.StatusSuccess,
		Timestamp:   time.Now(),
		ExecutionID: e.executionID,
		WorkflowID:  e.workflowID,
		NodeID:      node.ID,
		NodeType:    node.Type,
		StartTime:   startTime,
		ElapsedTime: time.Since(startTime),
		Result:      result,
	}

	e.observerMgr.Notify(ctx, event)
}

// notifyNodeFailure notifies observers that a node execution failed
func (e *Engine) notifyNodeFailure(ctx context.Context, node types.Node, startTime time.Time, result interface{}, err error) {
	if !e.observerMgr.HasObservers() {
		return
	}

	event := observer.Event{
		Type:        observer.EventNodeFailure,
		Status:      observer.StatusFailure,
		Timestamp:   time.Now(),
		ExecutionID: e.executionID,
		WorkflowID:  e.workflowID,
		NodeID:      node.ID,
		NodeType:    node.Type,
		StartTime:   startTime,
		ElapsedTime: time.Since(startTime),
		Result:      result,
		Error:       err,
	}

	e.observerMgr.Notify(ctx, event)
}
