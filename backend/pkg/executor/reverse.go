package executor

import (
	"context"
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ReverseExecutor reverses the order of an array
type ReverseExecutor struct{}

// Execute reverses the input array
func (e *ReverseExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get input array
	input, ok := inputs["in"]
	if !ok {
		return nil, fmt.Errorf("reverse node missing required input 'in'")
	}

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		return map[string]interface{}{
			"error": "input is not an array",
			"input": input,
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
