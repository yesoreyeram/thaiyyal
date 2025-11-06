package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestFlatMapExecutor_Basic(t *testing.T) {
	tests := []struct {
		name          string
		inputArray    []interface{}
		field         *string
		expectedCount int
		description   string
	}{
		{
			name: "Flatten nested arrays",
			inputArray: []interface{}{
				[]interface{}{1.0, 2.0},
				[]interface{}{3.0, 4.0},
				[]interface{}{5.0},
			},
			expectedCount: 5,
			description:   "Should flatten nested arrays",
		},
		{
			name: "Flatten object field arrays",
			inputArray: []interface{}{
				map[string]interface{}{"items": []interface{}{1.0, 2.0}},
				map[string]interface{}{"items": []interface{}{3.0, 4.0, 5.0}},
			},
			field:         stringPtr("items"),
			expectedCount: 5,
			description:   "Should extract and flatten field arrays",
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
			exec := &FlatMapExecutor{}
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
				Type: types.NodeTypeFlatMap,
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

			flattened, ok := resultMap["flattened"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'flattened' to be array, got %T", resultMap["flattened"])
			}

			if len(flattened) != tt.expectedCount {
				t.Errorf("Expected %d elements, got %d", tt.expectedCount, len(flattened))
			}
		})
	}
}

func TestFlatMapExecutor_Validate(t *testing.T) {
	exec := &FlatMapExecutor{}
	node := types.Node{
		Type: types.NodeTypeFlatMap,
		Data: types.FlatMapData{},
	}

	err := exec.Validate(node)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}
