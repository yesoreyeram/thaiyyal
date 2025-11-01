package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// JoinExecutor executes Join nodes
type JoinExecutor struct{}

// Execute runs the Join node
// Handles joining/merging multiple inputs
func (e *JoinExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)

	strategy := "all" // default strategy
	if node.Data.JoinStrategy != nil {
		strategy = *node.Data.JoinStrategy
	}

	switch strategy {
	case "all":
		// Wait for all inputs and combine them
		if len(inputs) == 0 {
			return nil, fmt.Errorf("join node with 'all' strategy requires at least one input")
		}
		return map[string]interface{}{
			"strategy": "all",
			"values":   inputs,
			"count":    len(inputs),
		}, nil

	case "any":
		// Return as soon as any input is available
		if len(inputs) > 0 {
			return map[string]interface{}{
				"strategy": "any",
				"value":    inputs[0],
				"count":    len(inputs),
			}, nil
		}
		return nil, fmt.Errorf("join node with 'any' strategy has no inputs")

	case "first":
		// Return only the first input
		if len(inputs) > 0 {
			return map[string]interface{}{
				"strategy": "first",
				"value":    inputs[0],
			}, nil
		}
		return nil, fmt.Errorf("join node with 'first' strategy has no inputs")

	default:
		return nil, fmt.Errorf("unsupported join strategy: %s (use all, any, or first)", strategy)
	}
}

// NodeType returns the node type this executor handles
func (e *JoinExecutor) NodeType() types.NodeType {
	return types.NodeTypeJoin
}

// Validate checks if node configuration is valid
func (e *JoinExecutor) Validate(node types.Node) error {
	// No required fields for join
	return nil
}
