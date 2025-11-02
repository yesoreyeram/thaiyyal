package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestFindExecutor_Basic(t *testing.T) {
	tests := []struct {
		name        string
		inputArray  []interface{}
		condition   string
		expectFound bool
		expectIndex int
		description string
	}{
		{
			name:        "Find number greater than 10",
			inputArray:  []interface{}{5.0, 15.0, 8.0, 20.0},
			condition:   "variables.item > 10",
			expectFound: true,
			expectIndex: 1,
			description: "Should find first number > 10",
		},
		{
			name:        "Find string equals 'apple'",
			inputArray:  []interface{}{"banana", "apple", "orange"},
			condition:   `variables.item == "apple"`,
			expectFound: true,
			expectIndex: 1,
			description: "Should find matching string",
		},
		{
			name:        "Find with index condition",
			inputArray:  []interface{}{1.0, 2.0, 3.0, 4.0},
			condition:   "variables.index == 2",
			expectFound: true,
			expectIndex: 2,
			description: "Should find by index",
		},
		{
			name:        "No match found",
			inputArray:  []interface{}{1.0, 2.0, 3.0},
			condition:   "variables.item > 100",
			expectFound: false,
			description: "Should return not_found when no match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FindExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
				variables: make(map[string]interface{}),
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeFind,
				Data: types.NodeData{
					Condition: &tt.condition,
				},
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

			if tt.expectFound {
				if resultMap["not_found"] != nil && resultMap["not_found"].(bool) {
					t.Error("Expected to find a match but got not_found")
				}
				if resultMap["found"] == nil {
					t.Error("Expected found field to be present")
				}
			} else {
				if notFound, ok := resultMap["not_found"].(bool); !ok || !notFound {
					t.Error("Expected not_found to be true")
				}
			}
		})
	}
}

func TestFindExecutor_Validate(t *testing.T) {
	exec := &FindExecutor{}

	// Missing condition
	node := types.Node{
		Type: types.NodeTypeFind,
		Data: types.NodeData{},
	}

	err := exec.Validate(node)
	if err == nil {
		t.Error("Expected validation error for missing condition")
	}

	// Valid condition
	condition := "variables.item > 5"
	node.Data.Condition = &condition
	err = exec.Validate(node)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}
