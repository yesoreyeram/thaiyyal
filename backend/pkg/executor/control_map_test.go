package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestMapExecutor_ExtractField(t *testing.T) {
	executor := &MapExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"map1": {[]interface{}{
				map[string]interface{}{"name": "Alice", "age": float64(25)},
				map[string]interface{}{"name": "Bob", "age": float64(30)},
			}},
		},
	}

	field := "name"
	node := types.Node{
		ID:   "map1",
		Type: types.NodeTypeMap,
		Data: types.MapData{
			Field: &field,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	results := resultMap["results"].([]interface{})

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0] != "Alice" {
		t.Errorf("Expected first result to be 'Alice', got %v", results[0])
	}
	if results[1] != "Bob" {
		t.Errorf("Expected second result to be 'Bob', got %v", results[1])
	}
}

func TestMapExecutor_ExpressionTransform(t *testing.T) {
	t.Skip("Expression-based transformations require EvaluateExpression enhancement - TODO")
	executor := &MapExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"map1": {[]interface{}{
				float64(1), float64(2), float64(3),
			}},
		},
	}

	expr := "item * 2"
	node := types.Node{
		ID:   "map1",
		Type: types.NodeTypeMap,
		Data: types.ExpressionData{
			Expression: &expr,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	results := resultMap["results"].([]interface{})

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	expected := []float64{2, 4, 6}
	for i, exp := range expected {
		if results[i] != exp {
			t.Errorf("Expected result[%d] to be %v, got %v", i, exp, results[i])
		}
	}
}

func TestMapExecutor_ExpressionWithObjectField(t *testing.T) {
	t.Skip("Expression-based transformations require EvaluateExpression enhancement - TODO")
	executor := &MapExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"map1": {[]interface{}{
				map[string]interface{}{"age": float64(25)},
				map[string]interface{}{"age": float64(30)},
			}},
		},
	}

	expr := "item.age * 1.1"
	node := types.Node{
		ID:   "map1",
		Type: types.NodeTypeMap,
		Data: types.ExpressionData{
			Expression: &expr,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	results := resultMap["results"].([]interface{})

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// 25 * 1.1 = 27.5, 30 * 1.1 = 33.0
	if results[0].(float64) < 27.4 || results[0].(float64) > 27.6 {
		t.Errorf("Expected result[0] to be ~27.5, got %v", results[0])
	}
	if results[1].(float64) < 32.9 || results[1].(float64) > 33.1 {
		t.Errorf("Expected result[1] to be ~33.0, got %v", results[1])
	}
}

func TestMapExecutor_NonArrayInput(t *testing.T) {
	executor := &MapExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"map1": {"not an array"},
		},
	}

	field := "name"
	node := types.Node{
		ID:   "map1",
		Type: types.NodeTypeMap,
		Data: types.MapData{
			Field: &field,
		},
	}

	_, err := executor.Execute(ctx, node)
	if err == nil {
		t.Fatal("Expected error for non-array input, got nil")
	}
}

func TestMapExecutor_Validate(t *testing.T) {
	executor := &MapExecutor{}

	// Valid: has expression
	expr := "item * 2"
	node := types.Node{
		Data: types.MapData{Expression: &expr},
	}
	if err := executor.Validate(node); err != nil {
		t.Errorf("Validation should pass with expression: %v", err)
	}

	// Valid: has field
	field := "name"
	node = types.Node{
		Data: types.MapData{Field: &field},
	}
	if err := executor.Validate(node); err != nil {
		t.Errorf("Validation should pass with field: %v", err)
	}

	// Invalid: has both
	node = types.Node{
		Data: types.MapData{Expression: &expr, Field: &field},
	}
	if err := executor.Validate(node); err == nil {
		t.Error("Validation should fail with both expression and field")
	}

	// Invalid: has neither
	node = types.Node{
		Data: types.MapData{},
	}
	if err := executor.Validate(node); err == nil {
		t.Error("Validation should fail without expression or field")
	}
}

func TestMapExecutor_EmptyArray(t *testing.T) {
	executor := &MapExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"map1": {[]interface{}{}},
		},
	}

	field := "name"
	node := types.Node{
		ID:   "map1",
		Type: types.NodeTypeMap,
		Data: types.MapData{
			Field: &field,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	results := resultMap["results"].([]interface{})

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty array, got %d", len(results))
	}
}
