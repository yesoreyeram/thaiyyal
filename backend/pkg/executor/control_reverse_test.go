package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestReverseExecutor_Basic(t *testing.T) {
	tests := []struct {
		name        string
		inputArray  []interface{}
		expectFirst interface{}
		expectLast  interface{}
		description string
	}{
		{
			name:        "Reverse numbers",
			inputArray:  []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			expectFirst: 5.0,
			expectLast:  1.0,
			description: "Should reverse number array",
		},
		{
			name:        "Reverse strings",
			inputArray:  []interface{}{"first", "second", "third"},
			expectFirst: "third",
			expectLast:  "first",
			description: "Should reverse string array",
		},
		{
			name:        "Single element",
			inputArray:  []interface{}{42.0},
			expectFirst: 42.0,
			expectLast:  42.0,
			description: "Should handle single element",
		},
		{
			name:        "Empty array",
			inputArray:  []interface{}{},
			description: "Should handle empty array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ReverseExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeReverse,
				Data: types.ReverseData{},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected result to be map, got %T", result)
			}

			reversed, ok := resultMap["reversed"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'reversed' to be array, got %T", resultMap["reversed"])
			}

			if len(reversed) != len(tt.inputArray) {
				t.Errorf("Expected %d elements, got %d", len(tt.inputArray), len(reversed))
			}

			if len(reversed) > 0 && tt.expectFirst != nil {
				if reversed[0] != tt.expectFirst {
					t.Errorf("Expected first element to be %v, got %v", tt.expectFirst, reversed[0])
				}
			}

			if len(reversed) > 0 && tt.expectLast != nil {
				if reversed[len(reversed)-1] != tt.expectLast {
					t.Errorf("Expected last element to be %v, got %v", tt.expectLast, reversed[len(reversed)-1])
				}
			}
		})
	}
}

func TestReverseExecutor_Validate(t *testing.T) {
	exec := &ReverseExecutor{}
	node := types.Node{
		Type: types.NodeTypeReverse,
		Data: types.ReverseData{},
	}

	err := exec.Validate(node)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}
