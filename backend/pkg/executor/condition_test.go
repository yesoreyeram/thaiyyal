package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestConditionExecutor_Basic tests basic condition evaluation
func TestConditionExecutor_Basic(t *testing.T) {
	tests := []struct {
		name          string
		condition     string
		input         interface{}
		expectedMet   bool
		description   string
	}{
		{
			name:        "Greater than - true",
			condition:   ">10",
			input:       float64(15),
			expectedMet: true,
			description: "15 > 10 should be true",
		},
		{
			name:        "Greater than - false",
			condition:   ">10",
			input:       float64(5),
			expectedMet: false,
			description: "5 > 10 should be false",
		},
		{
			name:        "Less than - true",
			condition:   "<10",
			input:       float64(5),
			expectedMet: true,
			description: "5 < 10 should be true",
		},
		{
			name:        "Greater than or equal - boundary",
			condition:   ">=10",
			input:       float64(10),
			expectedMet: true,
			description: "10 >= 10 should be true",
		},
		{
			name:        "Less than or equal - boundary",
			condition:   "<=10",
			input:       float64(10),
			expectedMet: true,
			description: "10 <= 10 should be true",
		},
		{
			name:        "Equality - true",
			condition:   "==5",
			input:       float64(5),
			expectedMet: true,
			description: "5 == 5 should be true",
		},
		{
			name:        "Inequality - true",
			condition:   "!=10",
			input:       float64(5),
			expectedMet: true,
			description: "5 != 10 should be true",
		},
		{
			name:        "Boolean literal - true",
			condition:   "true",
			input:       float64(100),
			expectedMet: true,
			description: "Literal true should always be true",
		},
		{
			name:        "Boolean literal - false",
			condition:   "false",
			input:       float64(100),
			expectedMet: false,
			description: "Literal false should always be false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ConditionExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeCondition,
				Data: types.NodeData{
					Condition: &tt.condition,
				},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected result to be map, got %T", result)
			}

			conditionMet, ok := resultMap["condition_met"].(bool)
			if !ok {
				t.Fatalf("Expected condition_met to be bool, got %T", resultMap["condition_met"])
			}

			if conditionMet != tt.expectedMet {
				t.Errorf("Expected condition_met=%v, got %v. Description: %s",
					tt.expectedMet, conditionMet, tt.description)
			}

			// Verify path indicators
			pathTaken := resultMap["path"].(string)
			if tt.expectedMet && pathTaken != "true" {
				t.Errorf("Expected path='true', got '%s'", pathTaken)
			}
			if !tt.expectedMet && pathTaken != "false" {
				t.Errorf("Expected path='false', got '%s'", pathTaken)
			}

			// Verify value is passed through
			if resultMap["value"] != tt.input {
				t.Errorf("Expected value to be passed through: %v, got %v", tt.input, resultMap["value"])
			}
		})
	}
}

// TestConditionExecutor_ComplexExpressions tests advanced expression evaluation
func TestConditionExecutor_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name        string
		condition   string
		input       interface{}
		expectedMet bool
		description string
	}{
		{
			name:        "Object field access",
			condition:   "input.age > 18",
			input:       map[string]interface{}{"age": float64(25), "name": "Alice"},
			expectedMet: true,
			description: "Should access input.age field",
		},
		{
			name:        "Nested field access",
			condition:   "input.profile.verified == true",
			input:       map[string]interface{}{"profile": map[string]interface{}{"verified": true}},
			expectedMet: true,
			description: "Should access nested fields",
		},
		{
			name:        "Boolean logic - AND",
			condition:   "input > 10 && input < 20",
			input:       float64(15),
			expectedMet: true,
			description: "Should evaluate AND correctly",
		},
		{
			name:        "Boolean logic - OR",
			condition:   "input < 5 || input > 20",
			input:       float64(3),
			expectedMet: true,
			description: "Should evaluate OR correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ConditionExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeCondition,
				Data: types.NodeData{
					Condition: &tt.condition,
				},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resultMap := result.(map[string]interface{})
			conditionMet := resultMap["condition_met"].(bool)

			if conditionMet != tt.expectedMet {
				t.Errorf("Expected condition_met=%v, got %v. Description: %s",
					tt.expectedMet, conditionMet, tt.description)
			}
		})
	}
}

// TestConditionExecutor_WithVariables tests condition with workflow variables
func TestConditionExecutor_WithVariables(t *testing.T) {
	exec := &ConditionExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {float64(25)},
		},
		variables: map[string]interface{}{
			"threshold": float64(18),
		},
	}

	condition := "input >= variables.threshold"
	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeCondition,
		Data: types.NodeData{
			Condition: &condition,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	if !resultMap["condition_met"].(bool) {
		t.Error("Expected condition to be met with variable reference")
	}
}

// TestConditionExecutor_WithContextVariables tests condition with context variables
func TestConditionExecutor_WithContextVariables(t *testing.T) {
	exec := &ConditionExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {float64(85)},
		},
		contextVars: map[string]interface{}{
			"passingScore": float64(70),
		},
	}

	condition := "input > context.passingScore"
	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeCondition,
		Data: types.NodeData{
			Condition: &condition,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	if !resultMap["condition_met"].(bool) {
		t.Error("Expected condition to be met with context variable reference")
	}
}

// TestConditionExecutor_WithNodeReferences tests condition with node output references
func TestConditionExecutor_WithNodeReferences(t *testing.T) {
	exec := &ConditionExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {float64(25)},
		},
		nodeResults: map[string]interface{}{
			"threshold-node": map[string]interface{}{"value": float64(10)},
		},
	}

	condition := "input > node.threshold-node.value"
	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeCondition,
		Data: types.NodeData{
			Condition: &condition,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	if !resultMap["condition_met"].(bool) {
		t.Error("Expected condition to be met with node reference")
	}
}

// TestConditionExecutor_Validation tests node validation
func TestConditionExecutor_Validation(t *testing.T) {
	tests := []struct {
		name        string
		node        types.Node
		expectError bool
	}{
		{
			name: "Valid node",
			node: types.Node{
				Type: types.NodeTypeCondition,
				Data: types.NodeData{
					Condition: strPtr(">0"),
				},
			},
			expectError: false,
		},
		{
			name: "Missing condition",
			node: types.Node{
				Type: types.NodeTypeCondition,
				Data: types.NodeData{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ConditionExecutor{}
			err := exec.Validate(tt.node)
			
			if tt.expectError && err == nil {
				t.Error("Expected validation error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
		})
	}
}

// TestConditionExecutor_MissingInput tests error handling for missing input
func TestConditionExecutor_MissingInput(t *testing.T) {
	exec := &ConditionExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{}, // No inputs
	}

	condition := ">10"
	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeCondition,
		Data: types.NodeData{
			Condition: &condition,
		},
	}

	_, err := exec.Execute(ctx, node)
	if err == nil {
		t.Error("Expected error for missing input")
	}
}

// TestConditionExecutor_NodeType tests NodeType method
func TestConditionExecutor_NodeType(t *testing.T) {
	exec := &ConditionExecutor{}
	if exec.NodeType() != types.NodeTypeCondition {
		t.Errorf("Expected NodeType to be %s, got %s", types.NodeTypeCondition, exec.NodeType())
	}
}

// Helper function for string pointers
func strPtr(s string) *string {
	return &s
}
