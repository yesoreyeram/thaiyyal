package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ForEachExecutor executes ForEach nodes
type ForEachExecutor struct{}

// Execute runs the ForEach node
// Iterates over an array input.
// This is a simplified implementation that validates the array and returns metadata.
// A full implementation would execute a sub-workflow for each array element.
func (e *ForEachExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("for_each node needs at least 1 input")
	}

	// Check if input is an array (slice)
	inputArray, ok := inputs[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("for_each node requires array input, got %T", inputs[0])
	}

	// Set default max iterations
	maxIter := 1000
	if node.Data.MaxIterations != nil && *node.Data.MaxIterations > 0 {
		maxIter = *node.Data.MaxIterations
	}

	// Limit iterations to prevent resource exhaustion
	iterCount := len(inputArray)
	if iterCount > maxIter {
		return nil, fmt.Errorf("for_each exceeds max iterations: %d > %d", iterCount, maxIter)
	}

	// TODO: In a full implementation, execute sub-workflow for each element
	// For now, return metadata about the iteration
	return map[string]interface{}{
		"items":      inputArray,
		"count":      len(inputArray),
		"iterations": iterCount,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *ForEachExecutor) NodeType() types.NodeType {
	return types.NodeTypeForEach
}

// Validate checks if node configuration is valid
func (e *ForEachExecutor) Validate(node types.Node) error {
	// No required fields for foreach
	return nil
}
