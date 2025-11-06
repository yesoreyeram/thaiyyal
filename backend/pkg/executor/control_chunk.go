package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ChunkExecutor splits an array into fixed-size chunks
type ChunkExecutor struct{}

// Execute splits the input array into chunks
func (e *ChunkExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsChunkData(node.Data)
	if err != nil {
		return nil, err
	}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("chunk node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("chunk node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	// Get chunk size
	size := 10 // default
	if sizeVal, ok := data.Size.(float64); ok {
		size = int(sizeVal)
	} else if sizeVal, ok := data.Size.(int); ok {
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
	data, err := types.AsChunkData(node.Data)
	if err != nil {
		return err
	}
	if data.Size != nil {
		size := 0
		if sizeVal, ok := data.Size.(float64); ok {
			size = int(sizeVal)
		} else if sizeVal, ok := data.Size.(int); ok {
			size = sizeVal
		}
		if size <= 0 {
			return fmt.Errorf("chunk size must be greater than 0, got %d", size)
		}
	}
	return nil
}
