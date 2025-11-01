package executor

import (
	"fmt"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// DelayExecutor executes Delay nodes
type DelayExecutor struct{}

// Execute runs the Delay node
// Handles execution delay
func (e *DelayExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	var inputValue interface{}
	if len(inputs) > 0 {
		inputValue = inputs[0]
	}

	if node.Data.Duration == nil {
		return nil, fmt.Errorf("delay node requires duration field")
	}

	duration, err := parseDuration(*node.Data.Duration)
	if err != nil {
		return nil, fmt.Errorf("invalid duration format: %w", err)
	}

	// Perform the delay
	time.Sleep(duration)

	return map[string]interface{}{
		"value":    inputValue,
		"duration": *node.Data.Duration,
		"delayed":  true,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *DelayExecutor) NodeType() types.NodeType {
	return types.NodeTypeDelay
}

// Validate checks if node configuration is valid
func (e *DelayExecutor) Validate(node types.Node) error {
	if node.Data.Duration == nil {
		return fmt.Errorf("delay node requires duration field")
	}
	return nil
}
