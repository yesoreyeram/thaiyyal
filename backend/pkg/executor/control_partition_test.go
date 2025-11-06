package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestPartitionExecutor_Basic(t *testing.T) {
	tests := []struct {
		name           string
		inputArray     []interface{}
		condition      string
		expectedPassed int
		expectedFailed int
		description    string
	}{
		{
			name:           "Partition numbers by > 10",
			inputArray:     []interface{}{5.0, 15.0, 8.0, 20.0, 3.0, 12.0},
			condition:      "variables.item > 10",
			expectedPassed: 3,
			expectedFailed: 3,
			description:    "Should partition into two groups",
		},
		{
			name:           "All pass",
			inputArray:     []interface{}{15.0, 20.0, 30.0},
			condition:      "variables.item > 10",
			expectedPassed: 3,
			expectedFailed: 0,
			description:    "Should have all in passed group",
		},
		{
			name:           "All fail",
			inputArray:     []interface{}{1.0, 2.0, 3.0},
			condition:      "variables.item > 10",
			expectedPassed: 0,
			expectedFailed: 3,
			description:    "Should have all in failed group",
		},
		{
			name:           "Empty array",
			inputArray:     []interface{}{},
			condition:      "variables.item > 10",
			expectedPassed: 0,
			expectedFailed: 0,
			description:    "Should handle empty array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &PartitionExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
				variables: make(map[string]interface{}),
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypePartition,
				Data: types.ConditionData{
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

			passed, ok := resultMap["passed"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'passed' to be array, got %T", resultMap["passed"])
			}

			failed, ok := resultMap["failed"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'failed' to be array, got %T", resultMap["failed"])
			}

			if len(passed) != tt.expectedPassed {
				t.Errorf("Expected %d passed, got %d", tt.expectedPassed, len(passed))
			}

			if len(failed) != tt.expectedFailed {
				t.Errorf("Expected %d failed, got %d", tt.expectedFailed, len(failed))
			}
		})
	}
}

func TestPartitionExecutor_Validate(t *testing.T) {
	exec := &PartitionExecutor{}

	// Missing condition
	node := types.Node{
		Type: types.NodeTypePartition,
		Data: types.PartitionData{},
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
