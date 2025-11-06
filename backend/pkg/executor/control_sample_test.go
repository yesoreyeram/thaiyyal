package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestSampleExecutor_Basic(t *testing.T) {
	tests := []struct {
		name          string
		inputArray    []interface{}
		count         interface{}
		method        *string
		expectedCount int
		description   string
	}{
		{
			name:          "Sample first 2 elements",
			inputArray:    []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			count:         2.0,
			method:        stringPtr("first"),
			expectedCount: 2,
			description:   "Should take first N elements",
		},
		{
			name:          "Sample last 3 elements",
			inputArray:    []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			count:         3.0,
			method:        stringPtr("last"),
			expectedCount: 3,
			description:   "Should take last N elements",
		},
		{
			name:          "Random sample",
			inputArray:    []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			count:         2.0,
			method:        stringPtr("random"),
			expectedCount: 2,
			description:   "Should take random N elements",
		},
		{
			name:          "Sample more than array length",
			inputArray:    []interface{}{1.0, 2.0, 3.0},
			count:         10.0,
			method:        stringPtr("first"),
			expectedCount: 3,
			description:   "Should return all when count > length",
		},
		{
			name:          "Default method (random) and count (1)",
			inputArray:    []interface{}{1.0, 2.0, 3.0},
			expectedCount: 1,
			description:   "Should use defaults",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &SampleExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			nodeData := types.SampleData{}
			if tt.count != nil {
				nodeData.Count = tt.count
			}
			if tt.method != nil {
				nodeData.Method = tt.method
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeSample,
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

			sample, ok := resultMap["sample"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'sample' to be array, got %T", resultMap["sample"])
			}

			if len(sample) != tt.expectedCount {
				t.Errorf("Expected %d elements, got %d", tt.expectedCount, len(sample))
			}
		})
	}
}

func TestSampleExecutor_Validate(t *testing.T) {
	exec := &SampleExecutor{}

	tests := []struct {
		name        string
		count       interface{}
		method      *string
		expectError bool
	}{
		{"Valid random method", 2.0, stringPtr("random"), false},
		{"Valid first method", 3.0, stringPtr("first"), false},
		{"Valid last method", 1.0, stringPtr("last"), false},
		{"Invalid method", 2.0, stringPtr("invalid"), true},
		{"Zero count", 0.0, stringPtr("random"), true},
		{"Negative count", -1.0, stringPtr("random"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeData := types.SampleData{}
			if tt.count != nil {
				nodeData.Count = tt.count
			}
			if tt.method != nil {
				nodeData.Method = tt.method
			}

			node := types.Node{
				Type: types.NodeTypeSample,
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
