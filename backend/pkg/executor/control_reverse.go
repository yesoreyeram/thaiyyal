package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ReverseExecutor reverses the order of an array
type ReverseExecutor struct{}

// Execute reverses the input array
func (e *ReverseExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("reverse node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("reverse node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	// Reverse the array
	reversed := make([]interface{}, len(arr))
	for i, v := range arr {
		reversed[len(arr)-1-i] = v
	}

	return map[string]interface{}{
		"reversed": reversed,
		"count":    len(reversed),
	}, nil
}

// NodeType returns the node type this executor handles
func (e *ReverseExecutor) NodeType() types.NodeType {
	return types.NodeTypeReverse
}

// Validate checks if the node configuration is valid
func (e *ReverseExecutor) Validate(node types.Node) error {
	// No configuration needed
	return nil
}
