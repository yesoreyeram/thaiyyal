package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestZipExecutor_Basic(t *testing.T) {
	tests := []struct {
		name          string
		inputArrays   [][]interface{}
		fillMissing   interface{}
		expectedCount int
		description   string
	}{
		{
			name: "Zip two arrays of same length",
			inputArrays: [][]interface{}{
				{1.0, 2.0, 3.0},
			},
			expectedCount: 3,
			description:   "Should zip arrays element-wise",
		},
		{
			name: "Zip with fill value",
			inputArrays: [][]interface{}{
				{1.0, 2.0, 3.0, 4.0},
			},
			fillMissing:   nil,
			expectedCount: 4,
			description:   "Should handle arrays of different lengths",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ZipExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArrays[0]},
				},
			}

			nodeData := types.NodeData{}
			if tt.fillMissing != nil {
				nodeData.FillMissing = tt.fillMissing
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeZip,
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

			zipped, ok := resultMap["zipped"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'zipped' to be array, got %T", resultMap["zipped"])
			}

			if len(zipped) != tt.expectedCount {
				t.Errorf("Expected %d tuples, got %d", tt.expectedCount, len(zipped))
			}
		})
	}
}

func TestZipExecutor_Validate(t *testing.T) {
	exec := &ZipExecutor{}
	node := types.Node{
		Type: types.NodeTypeZip,
		Data: types.NodeData{},
	}

	err := exec.Validate(node)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

func TestCompactExecutor_Basic(t *testing.T) {
	tests := []struct {
		name          string
		inputArray    []interface{}
		removeEmpty   *bool
		expectedCount int
		description   string
	}{
		{
			name:          "Remove null values",
			inputArray:    []interface{}{1.0, nil, 3.0, nil, 5.0},
			expectedCount: 3,
			description:   "Should remove nil values",
		},
		{
			name:          "Remove empty strings",
			inputArray:    []interface{}{"hello", "", "world", ""},
			removeEmpty:   boolPtr(true),
			expectedCount: 2,
			description:   "Should remove empty strings when removeEmpty=true",
		},
		{
			name:          "Keep empty strings",
			inputArray:    []interface{}{"hello", "", "world"},
			removeEmpty:   boolPtr(false),
			expectedCount: 3,
			description:   "Should keep empty strings when removeEmpty=false",
		},
		{
			name:          "All valid values",
			inputArray:    []interface{}{1.0, 2.0, 3.0},
			expectedCount: 3,
			description:   "Should keep all valid values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &CompactExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			nodeData := types.NodeData{}
			if tt.removeEmpty != nil {
				nodeData.RemoveEmpty = tt.removeEmpty
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeCompact,
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

			compacted, ok := resultMap["compacted"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'compacted' to be array, got %T", resultMap["compacted"])
			}

			if len(compacted) != tt.expectedCount {
				t.Errorf("Expected %d elements, got %d", tt.expectedCount, len(compacted))
			}
		})
	}
}

func TestCompactExecutor_Validate(t *testing.T) {
	exec := &CompactExecutor{}
	node := types.Node{
		Type: types.NodeTypeCompact,
		Data: types.NodeData{},
	}

	err := exec.Validate(node)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}
