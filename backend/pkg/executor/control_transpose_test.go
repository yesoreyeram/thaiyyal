package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestTransposeExecutor_Basic(t *testing.T) {
	tests := []struct {
		name         string
		inputArray   []interface{}
		expectedRows int
		expectedCols int
		description  string
	}{
		{
			name: "Transpose 2x3 matrix",
			inputArray: []interface{}{
				[]interface{}{1.0, 2.0, 3.0},
				[]interface{}{4.0, 5.0, 6.0},
			},
			expectedRows: 3,
			expectedCols: 2,
			description:  "Should transpose to 3x2 matrix",
		},
		{
			name: "Transpose square matrix",
			inputArray: []interface{}{
				[]interface{}{1.0, 2.0},
				[]interface{}{3.0, 4.0},
			},
			expectedRows: 2,
			expectedCols: 2,
			description:  "Should transpose square matrix",
		},
		{
			name: "Transpose single row",
			inputArray: []interface{}{
				[]interface{}{1.0, 2.0, 3.0},
			},
			expectedRows: 3,
			expectedCols: 1,
			description:  "Should transpose row to column",
		},
		{
			name:         "Empty array",
			inputArray:   []interface{}{},
			expectedRows: 0,
			expectedCols: 0,
			description:  "Should handle empty array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &TransposeExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeTranspose,
				Data: types.NodeData{},
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

			transposed, ok := resultMap["transposed"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'transposed' to be array, got %T", resultMap["transposed"])
			}

			if len(transposed) != tt.expectedRows {
				t.Errorf("Expected %d rows, got %d", tt.expectedRows, len(transposed))
			}

			if tt.expectedRows > 0 {
				firstRow, ok := transposed[0].([]interface{})
				if !ok {
					t.Fatalf("Expected first row to be array, got %T", transposed[0])
				}
				if len(firstRow) != tt.expectedCols {
					t.Errorf("Expected %d columns, got %d", tt.expectedCols, len(firstRow))
				}
			}
		})
	}
}

func TestTransposeExecutor_InvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		inputArray  []interface{}
		description string
	}{
		{
			name:        "Non-2D array (numbers)",
			inputArray:  []interface{}{1.0, 2.0, 3.0},
			description: "Should handle non-2D array gracefully",
		},
		{
			name: "Jagged array (different row lengths)",
			inputArray: []interface{}{
				[]interface{}{1.0, 2.0, 3.0},
				[]interface{}{4.0, 5.0},
			},
			description: "Should detect inconsistent row lengths",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &TransposeExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeTranspose,
				Data: types.NodeData{},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Errorf("Should not error, but return error in result: %v", err)
				return
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected result to be map, got %T", result)
			}

			// Should have error field in result
			if resultMap["error"] == nil {
				t.Error("Expected error field in result for invalid input")
			}
		})
	}
}

func TestTransposeExecutor_Validate(t *testing.T) {
	exec := &TransposeExecutor{}
	node := types.Node{
		Type: types.NodeTypeTranspose,
		Data: types.NodeData{},
	}

	err := exec.Validate(node)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}
