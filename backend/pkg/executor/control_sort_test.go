package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestSortExecutor_Basic(t *testing.T) {
	tests := []struct {
		name        string
		inputArray  []interface{}
		field       *string
		order       *string
		expectFirst interface{}
		expectLast  interface{}
		description string
	}{
		{
			name:        "Sort numbers ascending",
			inputArray:  []interface{}{5.0, 2.0, 8.0, 1.0, 9.0},
			order:       stringPtr("asc"),
			expectFirst: 1.0,
			expectLast:  9.0,
			description: "Should sort numbers in ascending order",
		},
		{
			name:        "Sort numbers descending",
			inputArray:  []interface{}{5.0, 2.0, 8.0, 1.0, 9.0},
			order:       stringPtr("desc"),
			expectFirst: 9.0,
			expectLast:  1.0,
			description: "Should sort numbers in descending order",
		},
		{
			name:        "Sort strings ascending",
			inputArray:  []interface{}{"zebra", "apple", "banana"},
			order:       stringPtr("asc"),
			expectFirst: "apple",
			expectLast:  "zebra",
			description: "Should sort strings alphabetically",
		},
		{
			name: "Sort objects by field",
			inputArray: []interface{}{
				map[string]interface{}{"name": "Charlie", "age": 30.0},
				map[string]interface{}{"name": "Alice", "age": 25.0},
				map[string]interface{}{"name": "Bob", "age": 35.0},
			},
			field:       stringPtr("name"),
			order:       stringPtr("asc"),
			expectFirst: map[string]interface{}{"name": "Alice", "age": 25.0},
			description: "Should sort objects by field value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &SortExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			nodeData := types.SortData{}
			if tt.field != nil {
				nodeData.Field = tt.field
			}
			if tt.order != nil {
				nodeData.Order = tt.order
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeSort,
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

			sorted, ok := resultMap["sorted"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'sorted' to be array, got %T", resultMap["sorted"])
			}

			if len(sorted) != len(tt.inputArray) {
				t.Errorf("Expected %d elements, got %d", len(tt.inputArray), len(sorted))
			}

			if len(sorted) > 0 {
				if tt.expectFirst != nil && !deepEqual(sorted[0], tt.expectFirst) {
					t.Errorf("Expected first element to be %v, got %v", tt.expectFirst, sorted[0])
				}
				if tt.expectLast != nil && !deepEqual(sorted[len(sorted)-1], tt.expectLast) {
					t.Errorf("Expected last element to be %v, got %v", tt.expectLast, sorted[len(sorted)-1])
				}
			}
		})
	}
}

func deepEqual(a, b interface{}) bool {
	switch av := a.(type) {
	case float64:
		if bv, ok := b.(float64); ok {
			return av == bv
		}
	case string:
		if bv, ok := b.(string); ok {
			return av == bv
		}
	case map[string]interface{}:
		if bv, ok := b.(map[string]interface{}); ok {
			if len(av) != len(bv) {
				return false
			}
			for k, v := range av {
				if !deepEqual(v, bv[k]) {
					return false
				}
			}
			return true
		}
	}
	return false
}
