package middleware

import (
	"strings"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestSizeLimitMiddleware_InputSizeLimit tests input size limiting
func TestSizeLimitMiddleware_InputSizeLimit(t *testing.T) {
	config := SizeLimitConfig{
		MaxInputSize:     100, // 100 bytes
		EnforceInputSize: true,
	}

	m := NewSizeLimitMiddlewareWithConfig(config)
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}

	// Create mock context with large input
	largeInput := strings.Repeat("x", 200) // 200 bytes
	ctx := &mockExecutionContextWithInputs{
		inputs: []interface{}{largeInput},
	}

	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "ok", nil
	}

	_, err := m.Process(ctx, node, handler)
	if err == nil {
		t.Error("expected error for large input, got nil")
	}

	if !strings.Contains(err.Error(), "input size limit exceeded") {
		t.Errorf("expected size limit error, got: %v", err)
	}
}

// TestSizeLimitMiddleware_ResultSizeLimit tests result size limiting
func TestSizeLimitMiddleware_ResultSizeLimit(t *testing.T) {
	config := SizeLimitConfig{
		MaxResultSize:     100, // 100 bytes
		EnforceResultSize: true,
	}

	m := NewSizeLimitMiddlewareWithConfig(config)
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	ctx := &mockExecutionContextWithInputs{inputs: []interface{}{}}

	// Handler returns large result
	largeResult := strings.Repeat("x", 200)
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return largeResult, nil
	}

	_, err := m.Process(ctx, node, handler)
	if err == nil {
		t.Error("expected error for large result, got nil")
	}

	if !strings.Contains(err.Error(), "result size limit exceeded") {
		t.Errorf("expected result size limit error, got: %v", err)
	}
}

// TestSizeLimitMiddleware_StringLengthLimit tests string length limiting
func TestSizeLimitMiddleware_StringLengthLimit(t *testing.T) {
	config := SizeLimitConfig{
		MaxInputSize:     1000, // Set high enough to not trigger first
		MaxStringLength:  50,
		EnforceInputSize: true,
	}

	m := NewSizeLimitMiddlewareWithConfig(config)
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}

	longString := strings.Repeat("x", 100)
	ctx := &mockExecutionContextWithInputs{
		inputs: []interface{}{longString},
	}

	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "ok", nil
	}

	_, err := m.Process(ctx, node, handler)
	if err == nil {
		t.Error("expected error for long string, got nil")
	}

	if !strings.Contains(err.Error(), "string length") {
		t.Errorf("expected string length error, got: %v", err)
	}
}

// TestSizeLimitMiddleware_ArrayLengthLimit tests array length limiting
func TestSizeLimitMiddleware_ArrayLengthLimit(t *testing.T) {
	config := SizeLimitConfig{
		MaxInputSize:     10000, // Set high enough to not trigger first
		MaxArrayLength:   10,
		EnforceInputSize: true,
	}

	m := NewSizeLimitMiddlewareWithConfig(config)
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}

	// Create array with 20 elements
	longArray := make([]interface{}, 20)
	for i := 0; i < 20; i++ {
		longArray[i] = i
	}

	ctx := &mockExecutionContextWithInputs{
		inputs: []interface{}{longArray},
	}

	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "ok", nil
	}

	_, err := m.Process(ctx, node, handler)
	if err == nil {
		t.Error("expected error for long array, got nil")
	}

	if !strings.Contains(err.Error(), "array length") {
		t.Errorf("expected array length error, got: %v", err)
	}
}

// TestSizeLimitMiddleware_AllowedInputs tests that allowed inputs pass
func TestSizeLimitMiddleware_AllowedInputs(t *testing.T) {
	m := NewSizeLimitMiddleware()
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}

	// Small, valid inputs
	ctx := &mockExecutionContextWithInputs{
		inputs: []interface{}{"hello", 42, true},
	}

	executionCount := 0
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		executionCount++
		return "ok", nil
	}

	result, err := m.Process(ctx, node, handler)
	if err != nil {
		t.Errorf("expected no error for valid inputs, got: %v", err)
	}

	if result != "ok" {
		t.Errorf("expected 'ok', got %v", result)
	}

	if executionCount != 1 {
		t.Errorf("expected handler to be called once, got %d", executionCount)
	}
}

// TestSizeLimitMiddleware_DisabledLimits tests with limits disabled
func TestSizeLimitMiddleware_DisabledLimits(t *testing.T) {
	config := SizeLimitConfig{
		MaxInputSize:      10,
		MaxResultSize:     10,
		EnforceInputSize:  false,
		EnforceResultSize: false,
	}

	m := NewSizeLimitMiddlewareWithConfig(config)
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}

	// Large input and result
	largeInput := strings.Repeat("x", 100)
	ctx := &mockExecutionContextWithInputs{
		inputs: []interface{}{largeInput},
	}

	largeResult := strings.Repeat("y", 100)
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return largeResult, nil
	}

	result, err := m.Process(ctx, node, handler)
	if err != nil {
		t.Errorf("expected no error with disabled limits, got: %v", err)
	}

	if result != largeResult {
		t.Error("result should be returned even if large when limits disabled")
	}
}

// TestSizeLimitMiddleware_Name tests the Name method
func TestSizeLimitMiddleware_Name(t *testing.T) {
	m := NewSizeLimitMiddleware()

	if m.Name() != "SizeLimit" {
		t.Errorf("expected 'SizeLimit', got %s", m.Name())
	}
}

// TestValidateWorkflowSize_NodeCount tests node count validation
func TestValidateWorkflowSize_NodeCount(t *testing.T) {
	config := SizeLimitConfig{
		MaxNodeCount: 5,
	}

	// Create 10 nodes
	nodes := make([]types.Node, 10)
	for i := 0; i < 10; i++ {
		nodes[i] = types.Node{ID: string(rune('a' + i)), Type: types.NodeTypeNumber}
	}

	err := ValidateWorkflowSize(nodes, []types.Edge{}, config)
	if err == nil {
		t.Error("expected error for too many nodes, got nil")
	}

	if !strings.Contains(err.Error(), "nodes") {
		t.Errorf("expected node count error, got: %v", err)
	}
}

// TestValidateWorkflowSize_EdgeCount tests edge count validation
func TestValidateWorkflowSize_EdgeCount(t *testing.T) {
	config := SizeLimitConfig{
		MaxEdgeCount: 5,
	}

	nodes := []types.Node{
		{ID: "1", Type: types.NodeTypeNumber},
		{ID: "2", Type: types.NodeTypeNumber},
	}

	// Create 10 edges
	edges := make([]types.Edge, 10)
	for i := 0; i < 10; i++ {
		edges[i] = types.Edge{Source: "1", Target: "2"}
	}

	err := ValidateWorkflowSize(nodes, edges, config)
	if err == nil {
		t.Error("expected error for too many edges, got nil")
	}

	if !strings.Contains(err.Error(), "edges") {
		t.Errorf("expected edge count error, got: %v", err)
	}
}

// TestValidateWorkflowSize_ValidWorkflow tests valid workflow passes
func TestValidateWorkflowSize_ValidWorkflow(t *testing.T) {
	config := DefaultSizeLimitConfig()

	nodes := []types.Node{
		{ID: "1", Type: types.NodeTypeNumber},
		{ID: "2", Type: types.NodeTypeNumber},
		{ID: "3", Type: types.NodeTypeNumber},
	}

	edges := []types.Edge{
		{Source: "1", Target: "2"},
		{Source: "2", Target: "3"},
	}

	err := ValidateWorkflowSize(nodes, edges, config)
	if err != nil {
		t.Errorf("expected no error for valid workflow, got: %v", err)
	}
}

// TestSizeLimitMiddleware_NestedStructures tests nested data validation
func TestSizeLimitMiddleware_NestedStructures(t *testing.T) {
	config := SizeLimitConfig{
		MaxStringLength:  20,
		EnforceInputSize: true,
	}

	m := NewSizeLimitMiddlewareWithConfig(config)
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}

	// Nested structure with long string
	nestedData := map[string]interface{}{
		"outer": map[string]interface{}{
			"inner": strings.Repeat("x", 50), // Exceeds limit
		},
	}

	ctx := &mockExecutionContextWithInputs{
		inputs: []interface{}{nestedData},
	}

	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "ok", nil
	}

	_, err := m.Process(ctx, node, handler)
	if err == nil {
		t.Error("expected error for nested string exceeding limit, got nil")
	}
}

// mockExecutionContextWithInputs for testing with custom inputs
type mockExecutionContextWithInputs struct {
	inputs []interface{}
}

func (m *mockExecutionContextWithInputs) GetHTTPClientRegistry() interface{} {
	return nil
}

func (m *mockExecutionContextWithInputs) GetNodeInputs(nodeID string) []interface{} {
	return m.inputs
}

func (m *mockExecutionContextWithInputs) GetNode(nodeID string) *types.Node {
	return nil
}

func (m *mockExecutionContextWithInputs) GetVariable(name string) (interface{}, error) {
	return nil, nil
}

func (m *mockExecutionContextWithInputs) SetVariable(name string, value interface{}) error {
	return nil
}

func (m *mockExecutionContextWithInputs) GetAccumulator() interface{} {
	return nil
}

func (m *mockExecutionContextWithInputs) SetAccumulator(value interface{}) {
}

func (m *mockExecutionContextWithInputs) GetCounter() float64 {
	return 0
}

func (m *mockExecutionContextWithInputs) SetCounter(value float64) {
}

func (m *mockExecutionContextWithInputs) GetCache(key string) (interface{}, bool) {
	return nil, false
}

func (m *mockExecutionContextWithInputs) SetCache(key string, value interface{}, ttl time.Duration) {
}

func (m *mockExecutionContextWithInputs) GetWorkflowContext() map[string]interface{} {
	return nil
}

func (m *mockExecutionContextWithInputs) GetContextVariable(name string) (interface{}, bool) {
	return nil, false
}

func (m *mockExecutionContextWithInputs) SetContextVariable(name string, value interface{}) {
}

func (m *mockExecutionContextWithInputs) GetContextConstant(name string) (interface{}, bool) {
	return nil, false
}

func (m *mockExecutionContextWithInputs) SetContextConstant(name string, value interface{}) {
}

func (m *mockExecutionContextWithInputs) InterpolateTemplate(template string) string {
	return template
}

func (m *mockExecutionContextWithInputs) GetNodeResult(nodeID string) (interface{}, bool) {
	return nil, false
}

func (m *mockExecutionContextWithInputs) SetNodeResult(nodeID string, result interface{}) {
}

func (m *mockExecutionContextWithInputs) GetAllNodeResults() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *mockExecutionContextWithInputs) GetVariables() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *mockExecutionContextWithInputs) GetContextVariables() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *mockExecutionContextWithInputs) GetConfig() types.Config {
	return types.Config{}
}

func (m *mockExecutionContextWithInputs) IncrementNodeExecution() error {
	return nil
}

func (m *mockExecutionContextWithInputs) IncrementHTTPCall() error {
	return nil
}

func (m *mockExecutionContextWithInputs) GetNodeExecutionCount() int {
	return 0
}

func (m *mockExecutionContextWithInputs) GetHTTPCallCount() int {
	return 0
}
