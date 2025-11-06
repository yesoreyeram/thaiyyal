package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestWhileLoopExecutor_Basic tests basic WhileLoop functionality
func TestWhileLoopExecutor_Basic(t *testing.T) {
	tests := []struct {
		name              string
		condition         string
		input             interface{}
		maxIterations     *int
		expectedIterCount int
		expectError       bool
		description       string
	}{
		{
			name:              "False condition - no iterations",
			condition:         ">100",
			input:             float64(5),
			expectedIterCount: 0,
			expectError:       false,
			description:       "5 > 100 is false, should iterate 0 times",
		},
		{
			name:              "True condition with limit",
			condition:         ">=0",
			input:             float64(5),
			maxIterations:     intPtr(10),
			expectedIterCount: 10,
			expectError:       true, // Should error due to exceeding max
			description:       "Always true condition should hit max iterations",
		},
		{
			name:              "Custom max iterations",
			condition:         ">=0",
			input:             float64(5),
			maxIterations:     intPtr(5),
			expectedIterCount: 5,
			expectError:       true,
			description:       "Should respect custom max iterations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &WhileLoopExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeWhileLoop,
				Data: types.WhileLoopData{
					Condition:     &tt.condition,
					MaxIterations: tt.maxIterations,
				},
			}

			result, err := exec.Execute(ctx, node)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error for exceeding max iterations, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected result to be map, got %T", result)
			}

			iterations, ok := resultMap["iterations"].(int)
			if !ok {
				t.Fatalf("Expected iterations to be int, got %T", resultMap["iterations"])
			}

			if iterations != tt.expectedIterCount {
				t.Errorf("Expected iterations=%d, got %d. Description: %s",
					tt.expectedIterCount, iterations, tt.description)
			}

			// Verify final value is preserved
			if resultMap["final_value"] != tt.input {
				t.Errorf("Expected final_value=%v, got %v", tt.input, resultMap["final_value"])
			}
		})
	}
}

// TestWhileLoopExecutor_Validation tests node validation
func TestWhileLoopExecutor_Validation(t *testing.T) {
	tests := []struct {
		name        string
		node        types.Node
		expectError bool
	}{
		{
			name: "Valid node",
			node: types.Node{
				Type: types.NodeTypeWhileLoop,
				Data: types.WhileLoopData{
					Condition: strPtr(">0"),
				},
			},
			expectError: false,
		},
		{
			name: "Missing condition",
			node: types.Node{
				Type: types.NodeTypeWhileLoop,
				Data: types.WhileLoopData{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &WhileLoopExecutor{}
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

// TestWhileLoopExecutor_MissingInput tests error handling for missing input
func TestWhileLoopExecutor_MissingInput(t *testing.T) {
	exec := &WhileLoopExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{}, // No inputs
	}

	condition := ">0"
	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeWhileLoop,
		Data: types.WhileLoopData{
			Condition: &condition,
		},
	}

	_, err := exec.Execute(ctx, node)
	if err == nil {
		t.Error("Expected error for missing input")
	}
}

// TestWhileLoopExecutor_NodeType tests NodeType method
func TestWhileLoopExecutor_NodeType(t *testing.T) {
	exec := &WhileLoopExecutor{}
	if exec.NodeType() != types.NodeTypeWhileLoop {
		t.Errorf("Expected NodeType to be %s, got %s", types.NodeTypeWhileLoop, exec.NodeType())
	}
}
