package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestRangeExecutor_Basic(t *testing.T) {
	tests := []struct {
		name          string
		start         interface{}
		end           interface{}
		step          interface{}
		expectedCount int
		expectError   bool
		description   string
	}{
		{
			name:          "Range 0 to 10 step 1",
			start:         0.0,
			end:           10.0,
			step:          1.0,
			expectedCount: 11,
			description:   "Should generate 0,1,2...10",
		},
		{
			name:          "Range 0 to 10 step 2",
			start:         0.0,
			end:           10.0,
			step:          2.0,
			expectedCount: 6,
			description:   "Should generate 0,2,4,6,8,10",
		},
		{
			name:          "Range 10 to 0 step -1",
			start:         10.0,
			end:           0.0,
			step:          -1.0,
			expectedCount: 11,
			description:   "Should generate 10,9,8...0",
		},
		{
			name:          "Range 5 to 5",
			start:         5.0,
			end:           5.0,
			step:          1.0,
			expectedCount: 1,
			description:   "Should generate single value",
		},
		{
			name:        "Invalid: positive step with start > end",
			start:       10.0,
			end:         0.0,
			step:        1.0,
			expectError: true,
			description: "Should error on mismatched direction",
		},
		{
			name:        "Invalid: zero step",
			start:       0.0,
			end:         10.0,
			step:        0.0,
			expectError: true,
			description: "Should error on zero step",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &RangeExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {},
				},
			}

			nodeData := types.NodeData{}
			if tt.start != nil {
				nodeData.Start = tt.start
			}
			if tt.end != nil {
				nodeData.End = tt.end
			}
			if tt.step != nil {
				nodeData.Step = tt.step
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeRange,
				Data: nodeData,
			}

			result, err := exec.Execute(ctx, node)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected result to be map, got %T", result)
			}

			rangeArr, ok := resultMap["range"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'range' to be array, got %T", resultMap["range"])
			}

			if len(rangeArr) != tt.expectedCount {
				t.Errorf("Expected %d elements, got %d", tt.expectedCount, len(rangeArr))
			}
		})
	}
}

func TestRangeExecutor_Validate(t *testing.T) {
	exec := &RangeExecutor{}

	tests := []struct {
		name        string
		start       interface{}
		end         interface{}
		step        interface{}
		expectError bool
	}{
		{"Valid ascending range", 0.0, 10.0, 1.0, false},
		{"Valid descending range", 10.0, 0.0, -1.0, false},
		{"Invalid: zero step", 0.0, 10.0, 0.0, true},
		{"Invalid: positive step, descending", 10.0, 0.0, 1.0, true},
		{"Invalid: negative step, ascending", 0.0, 10.0, -1.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeData := types.NodeData{}
			if tt.start != nil {
				nodeData.Start = tt.start
			}
			if tt.end != nil {
				nodeData.End = tt.end
			}
			if tt.step != nil {
				nodeData.Step = tt.step
			}

			node := types.Node{
				Type: types.NodeTypeRange,
				Data: nodeData,
			}

			err := exec.Validate(node)
			if tt.expectError && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}
