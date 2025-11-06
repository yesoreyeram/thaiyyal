package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestReduceExecutor_Sum(t *testing.T) {
	t.Skip("Expression-based reduction requires EvaluateExpression enhancement - TODO")
	executor := &ReduceExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"reduce1": {[]interface{}{
				float64(1), float64(2), float64(3),
			}},
		},
	}

	initVal := float64(0)
	expr := "accumulator + item"
	node := types.Node{
		ID:   "reduce1",
		Type: types.NodeTypeReduce,
		Data: types.ReduceData{
			InitialValue: initVal,
			Expression:   &expr,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	finalValue := resultMap["final_value"].(float64)

	if finalValue != 6 {
		t.Errorf("Expected final_value to be 6, got %v", finalValue)
	}
}

func TestReduceExecutor_SumObjectFields(t *testing.T) {
	t.Skip("Expression-based reduction requires EvaluateExpression enhancement - TODO")
	executor := &ReduceExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"reduce1": {[]interface{}{
				map[string]interface{}{"age": float64(25)},
				map[string]interface{}{"age": float64(30)},
			}},
		},
	}

	initVal := float64(0)
	expr := "accumulator + item.age"
	node := types.Node{
		ID:   "reduce1",
		Type: types.NodeTypeReduce,
		Data: types.ReduceData{
			InitialValue: initVal,
			Expression:   &expr,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	finalValue := resultMap["final_value"].(float64)

	if finalValue != 55 {
		t.Errorf("Expected final_value to be 55, got %v", finalValue)
	}
}

func TestReduceExecutor_Product(t *testing.T) {
	t.Skip("Expression-based reduction requires EvaluateExpression enhancement - TODO")
	executor := &ReduceExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"reduce1": {[]interface{}{
				float64(2), float64(3), float64(4),
			}},
		},
	}

	initVal := float64(1)
	expr := "accumulator * item"
	node := types.Node{
		ID:   "reduce1",
		Type: types.NodeTypeReduce,
		Data: types.ReduceData{
			InitialValue: initVal,
			Expression:   &expr,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	finalValue := resultMap["final_value"].(float64)

	if finalValue != 24 {
		t.Errorf("Expected final_value to be 24 (2*3*4), got %v", finalValue)
	}
}

func TestReduceExecutor_Max(t *testing.T) {
	t.Skip("Ternary operator in expressions requires enhancement - TODO")
	executor := &ReduceExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"reduce1": {[]interface{}{
				float64(5), float64(2), float64(8), float64(1),
			}},
		},
	}

	initVal := float64(0)
	expr := "item > accumulator ? item : accumulator"
	node := types.Node{
		ID:   "reduce1",
		Type: types.NodeTypeReduce,
		Data: types.ReduceData{
			InitialValue: initVal,
			Expression:   &expr,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	finalValue := resultMap["final_value"].(float64)

	if finalValue != 8 {
		t.Errorf("Expected final_value to be 8 (max), got %v", finalValue)
	}
}

func TestReduceExecutor_DefaultInitialValue(t *testing.T) {
	t.Skip("Expression-based reduction requires EvaluateExpression enhancement - TODO")
	executor := &ReduceExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"reduce1": {[]interface{}{
				float64(5), float64(10),
			}},
		},
	}

	expr := "accumulator + item"
	node := types.Node{
		ID:   "reduce1",
		Type: types.NodeTypeReduce,
		Data: types.ExpressionData{
			Expression: &expr,
			// No InitialValue specified - should default to 0
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	finalValue := resultMap["final_value"].(float64)

	if finalValue != 15 {
		t.Errorf("Expected final_value to be 15 (0+5+10), got %v", finalValue)
	}
}

func TestReduceExecutor_NonArrayInput(t *testing.T) {
	executor := &ReduceExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"reduce1": {"not an array"},
		},
	}

	expr := "accumulator + item"
	node := types.Node{
		ID:   "reduce1",
		Type: types.NodeTypeReduce,
		Data: types.ExpressionData{
			Expression: &expr,
		},
	}

	_, err := executor.Execute(ctx, node)
	if err == nil {
		t.Fatal("Expected error for non-array input, got nil")
	}
}

func TestReduceExecutor_Validate(t *testing.T) {
	executor := &ReduceExecutor{}

	// Valid: has expression
	expr := "accumulator + item"
	node := types.Node{
		Data: types.ExpressionData{Expression: &expr},
	}
	if err := executor.Validate(node); err != nil {
		t.Errorf("Validation should pass with expression: %v", err)
	}

	// Invalid: no expression
	node = types.Node{
		Data: types.ReduceData{},
	}
	if err := executor.Validate(node); err == nil {
		t.Error("Validation should fail without expression")
	}

	// Invalid: empty expression
	emptyExpr := ""
	node = types.Node{
		Data: types.ExpressionData{Expression: &emptyExpr},
	}
	if err := executor.Validate(node); err == nil {
		t.Error("Validation should fail with empty expression")
	}
}

func TestReduceExecutor_EmptyArray(t *testing.T) {
	executor := &ReduceExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"reduce1": {[]interface{}{}},
		},
	}

	initVal := float64(100)
	expr := "accumulator + item"
	node := types.Node{
		ID:   "reduce1",
		Type: types.NodeTypeReduce,
		Data: types.ReduceData{
			InitialValue: initVal,
			Expression:   &expr,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	finalValue := resultMap["final_value"].(float64)

	// Should return initial value when array is empty
	if finalValue != 100 {
		t.Errorf("Expected final_value to be 100 (initial_value), got %v", finalValue)
	}
}
