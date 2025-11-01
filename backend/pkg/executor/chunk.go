package executor

import (
	"context"
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ChunkExecutor splits an array into fixed-size chunks
type ChunkExecutor struct{}

// Execute splits the input array into chunks
func (e *ChunkExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get input array
	input, ok := inputs["in"]
	if !ok {
		return nil, fmt.Errorf("chunk node missing required input 'in'")
	}

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		return map[string]interface{}{
			"error": "input is not an array",
			"input": input,
		}, nil
	}

	// Get chunk size
	size := 10 // default
	if sizeVal, ok := node.Data.Size.(float64); ok {
		size = int(sizeVal)
	} else if sizeVal, ok := node.Data.Size.(int); ok {
		size = sizeVal
	}

	if size <= 0 {
		return nil, fmt.Errorf("chunk size must be greater than 0, got %d", size)
	}

	// Create chunks
	var chunks []interface{}
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		chunk := arr[i:end]
		chunks = append(chunks, chunk)
	}

	return map[string]interface{}{
		"chunks":      chunks,
		"input_count": len(arr),
		"chunk_count": len(chunks),
		"chunk_size":  size,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *ChunkExecutor) NodeType() types.NodeType {
	return types.NodeTypeChunk
}

// Validate checks if the node configuration is valid
func (e *ChunkExecutor) Validate(node types.Node) error {
	if node.Data.Size != nil {
		size := 0
		if sizeVal, ok := node.Data.Size.(float64); ok {
			size = int(sizeVal)
		} else if sizeVal, ok := node.Data.Size.(int); ok {
			size = sizeVal
		}
		if size <= 0 {
			return fmt.Errorf("chunk size must be greater than 0, got %d", size)
		}
	}
	return nil
}
