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
type Engine struct {
graph      *graph.Graph
state      *state.Manager
registry   *executor.Registry
config     types.Config
results    map[string]interface{}
resultsMu  sync.RWMutex
executionID string
workflowID  string

// Node storage for lookups
nodes []types.Node
edges []types.Edge
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
var payload types.Payload
if err := json.Unmarshal(payloadJSON, &payload); err != nil {
return nil, fmt.Errorf("failed to parse payload: %w", err)
}

engine := &Engine{
state:       state.New(),
registry:    defaultRegistry(),
config:      config,
results:     make(map[string]interface{}),
executionID: generateExecutionID(),
workflowID:  payload.WorkflowID,
nodes:       payload.Nodes,
edges:       payload.Edges,
}

// Create graph for topological sorting
engine.graph = graph.New(payload.Nodes, payload.Edges)

return engine, nil
}

// defaultRegistry creates and populates the default executor registry with all node executors
func defaultRegistry() *executor.Registry {
	reg := executor.NewRegistry()

	// Register all 25 node type executors
	// Basic I/O nodes
	reg.MustRegister(&executor.NumberExecutor{})
	reg.MustRegister(&executor.TextInputExecutor{})
	reg.MustRegister(&executor.VisualizationExecutor{})

	// Operation nodes
	reg.MustRegister(&executor.OperationExecutor{})
	reg.MustRegister(&executor.TextOperationExecutor{})
	reg.MustRegister(&executor.HTTPExecutor{})

	// Control flow nodes
	reg.MustRegister(&executor.ConditionExecutor{})
	reg.MustRegister(&executor.ForEachExecutor{})
	reg.MustRegister(&executor.WhileLoopExecutor{})

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
// Public API - Execute
// ============================================================================

// Execute runs the workflow and returns the result.
// The workflow execution is protected by a timeout configured in MaxExecutionTime.
// If the workflow takes longer than the timeout, execution is cancelled and an error is returned.
//
// Each workflow execution is assigned a unique execution ID that is passed through the
// execution context and included in the result. This ID can be used for logging and tracing.
//
// Returns:
//   - *types.Result: Workflow execution results including execution ID, node outputs and final output
//   - error: If execution fails, times out, or encounters an error
func (e *Engine) Execute() (*types.Result, error) {
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
return result, err
}

// Step 3: Create context with timeout and execution metadata for workflow execution
ctx, cancel := context.WithTimeout(context.Background(), e.config.MaxExecutionTime)
defer cancel()

// Add execution ID and workflow ID to context for logging and tracing
ctx = context.WithValue(ctx, types.ContextKeyExecutionID, e.executionID)
ctx = context.WithValue(ctx, types.ContextKeyWorkflowID, e.workflowID)

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
return result, err
}
case <-ctx.Done():
return result, fmt.Errorf("workflow execution timeout: exceeded %v", e.config.MaxExecutionTime)
}

// Step 4: Copy results and set final output
result.NodeResults = e.results
result.FinalOutput = e.getFinalOutput()

return result, nil
}

// ============================================================================
// Node Execution
// ============================================================================

// executeNode dispatches node execution to the appropriate executor via the registry.
// Handles template interpolation before execution (except for context nodes).
//
// Parameters:
//   - ctx: Context with execution metadata (execution ID, workflow ID)
//   - node: Node to execute
//
// Returns:
//   - interface{}: Result of node execution (type depends on node)
//   - error: If node execution fails
func (e *Engine) executeNode(ctx context.Context, node types.Node) (interface{}, error) {
// Interpolate templates in node data before execution (except for context nodes)
if node.Type != types.NodeTypeContextVariable && node.Type != types.NodeTypeContextConstant {
e.interpolateNodeData(&node.Data)
}

// Dispatch to appropriate executor via registry
return e.registry.Execute(e, node)
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

// SetVariable sets a variable value
func (e *Engine) SetVariable(name string, value interface{}) error {
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
func (e *Engine) SetContextVariable(name string, value interface{}) {
e.state.SetContextVariable(name, value)
}

// GetContextConstant retrieves a context constant
func (e *Engine) GetContextConstant(name string) (interface{}, bool) {
return e.state.GetContextConstant(name)
}

// SetContextConstant sets a context constant
func (e *Engine) SetContextConstant(name string, value interface{}) {
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
func (e *Engine) SetNodeResult(nodeID string, result interface{}) {
e.resultsMu.Lock()
defer e.resultsMu.Unlock()

e.results[nodeID] = result
}

// GetConfig returns the engine configuration
func (e *Engine) GetConfig() types.Config {
return e.config
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
