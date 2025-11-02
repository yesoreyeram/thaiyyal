package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestSliceExecutor_Basic(t *testing.T) {
	tests := []struct {
		name          string
		inputArray    []interface{}
		start         interface{}
		end           interface{}
		length        interface{}
		expectedCount int
		description   string
	}{
		{
			name:          "Slice first 3 elements",
			inputArray:    []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			start:         0.0,
			end:           3.0,
			expectedCount: 3,
			description:   "Should extract first 3 elements",
		},
		{
			name:          "Slice using length parameter",
			inputArray:    []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			start:         1.0,
			length:        2.0,
			expectedCount: 2,
			description:   "Should extract 2 elements starting from index 1",
		},
		{
			name:          "Slice with negative start",
			inputArray:    []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			start:         -2.0,
			expectedCount: 2,
			description:   "Should extract last 2 elements",
		},
		{
			name:          "Slice with negative end",
			inputArray:    []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			start:         0.0,
			end:           -1.0,
			expectedCount: 4,
			description:   "Should extract all but last element",
		},
		{
			name:          "Empty array",
			inputArray:    []interface{}{},
			start:         0.0,
			end:           3.0,
			expectedCount: 0,
			description:   "Should handle empty array gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &SliceExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			nodeData := types.NodeData{}
			if tt.start != nil {
				nodeData.Start = tt.start
			}
			if tt.end != nil {
				nodeData.End = tt.end
			}
			if tt.length != nil {
				nodeData.Length = tt.length
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeSlice,
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

			sliced, ok := resultMap["sliced"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'sliced' to be array, got %T", resultMap["sliced"])
			}

			if len(sliced) != tt.expectedCount {
				t.Errorf("Expected %d elements, got %d", tt.expectedCount, len(sliced))
			}
		})
	}
}

func TestSliceExecutor_NonArrayInput(t *testing.T) {
	exec := &SliceExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {"not an array"},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeSlice,
		Data: types.NodeData{
			Start: 0.0,
			End:   2.0,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Errorf("Should not error on non-array input: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result to be map, got %T", result)
	}

	if resultMap["error"] == nil {
		t.Error("Expected error field in result for non-array input")
	}
}

func TestSliceExecutor_Validate(t *testing.T) {
	exec := &SliceExecutor{}
	node := types.Node{
		Type: types.NodeTypeSlice,
		Data: types.NodeData{},
	}

	err := exec.Validate(node)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}
