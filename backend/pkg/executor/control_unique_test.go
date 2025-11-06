package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestUniqueExecutor_Basic(t *testing.T) {
	tests := []struct {
		name          string
		inputArray    []interface{}
		field         *string
		expectedCount int
		description   string
	}{
		{
			name:          "Remove duplicate numbers",
			inputArray:    []interface{}{1.0, 2.0, 1.0, 3.0, 2.0, 4.0},
			expectedCount: 4,
			description:   "Should remove duplicate numbers",
		},
		{
			name:          "Remove duplicate strings",
			inputArray:    []interface{}{"apple", "banana", "apple", "orange"},
			expectedCount: 3,
			description:   "Should remove duplicate strings",
		},
		{
			name: "Unique by field",
			inputArray: []interface{}{
				map[string]interface{}{"id": 1.0, "name": "Alice"},
				map[string]interface{}{"id": 2.0, "name": "Bob"},
				map[string]interface{}{"id": 1.0, "name": "Charlie"},
			},
			field:         stringPtr("id"),
			expectedCount: 2,
			description:   "Should keep unique by field value",
		},
		{
			name:          "All unique",
			inputArray:    []interface{}{1.0, 2.0, 3.0, 4.0},
			expectedCount: 4,
			description:   "Should return all when all are unique",
		},
		{
			name:          "Empty array",
			inputArray:    []interface{}{},
			expectedCount: 0,
			description:   "Should handle empty array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &UniqueExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			nodeData := types.NodeData{}
			if tt.field != nil {
				nodeData.Field = tt.field
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeUnique,
				Data: nodeData,
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

			unique, ok := resultMap["unique"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'unique' to be array, got %T", resultMap["unique"])
			}

			if len(unique) != tt.expectedCount {
				t.Errorf("Expected %d elements, got %d", tt.expectedCount, len(unique))
			}

			outputCount, _ := resultMap["output_count"].(int)
			if outputCount != tt.expectedCount {
				t.Errorf("Expected output_count %d, got %d", tt.expectedCount, outputCount)
			}
		})
	}
}

func TestUniqueExecutor_Validate(t *testing.T) {
	exec := &UniqueExecutor{}
	node := types.Node{
		Type: types.NodeTypeUnique,
		Data: types.UniqueData{},
	}

	err := exec.Validate(node)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}
