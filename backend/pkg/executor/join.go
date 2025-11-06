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
data, err := types.AsJoinData(node.Data)
if err != nil {
return nil, err
}
	inputs := ctx.GetNodeInputs(node.ID)

	strategy := "all" // default strategy
	if data.JoinStrategy != nil {
		strategy = *data.JoinStrategy
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
	// Validate node data type
	if _, err := types.AsJoinData(node.Data); err != nil {
		return err
	}
	// No required fields for join
	return nil
}
