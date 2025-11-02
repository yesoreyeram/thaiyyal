package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestChunkExecutor_Basic(t *testing.T) {
	tests := []struct {
		name           string
		inputArray     []interface{}
		size           interface{}
		expectedChunks int
		description    string
	}{
		{
			name:           "Chunk array into size 2",
			inputArray:     []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			size:           2.0,
			expectedChunks: 3,
			description:    "Should create 3 chunks of size 2",
		},
		{
			name:           "Chunk array into size 3",
			inputArray:     []interface{}{1.0, 2.0, 3.0, 4.0, 5.0, 6.0},
			size:           3.0,
			expectedChunks: 2,
			description:    "Should create 2 chunks of size 3",
		},
		{
			name:           "Chunk size larger than array",
			inputArray:     []interface{}{1.0, 2.0, 3.0},
			size:           10.0,
			expectedChunks: 1,
			description:    "Should create single chunk when size > array length",
		},
		{
			name:           "Empty array",
			inputArray:     []interface{}{},
			size:           2.0,
			expectedChunks: 0,
			description:    "Should handle empty array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ChunkExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeChunk,
				Data: types.NodeData{
					Size: tt.size,
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

			chunks, ok := resultMap["chunks"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'chunks' to be array, got %T", resultMap["chunks"])
			}

			if len(chunks) != tt.expectedChunks {
				t.Errorf("Expected %d chunks, got %d", tt.expectedChunks, len(chunks))
			}
		})
	}
}

func TestChunkExecutor_Validate(t *testing.T) {
	exec := &ChunkExecutor{}

	tests := []struct {
		name        string
		size        interface{}
		expectError bool
	}{
		{"Valid size", 2.0, false},
		{"Zero size", 0.0, true},
		{"Negative size", -1.0, true},
		{"No size (default)", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeData := types.NodeData{}
			if tt.size != nil {
				nodeData.Size = tt.size
			}

			node := types.Node{
				Type: types.NodeTypeChunk,
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
