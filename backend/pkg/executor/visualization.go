package executor

import (
	"fmt"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// VisualizationExecutor executes Visualization nodes
type VisualizationExecutor struct{}

// Execute formats output for display
func (e *VisualizationExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.Mode == nil {
		return nil, fmt.Errorf("visualization node missing mode")
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("visualization needs at least 1 input")
	}

	return map[string]interface{}{
		"mode":  *node.Data.Mode,
		"value": inputs[0],
	}, nil
}

// NodeType returns the node type this executor handles
func (e *VisualizationExecutor) NodeType() types.NodeType {
	return types.NodeTypeVisualization
}

// Validate checks if node configuration is valid
func (e *VisualizationExecutor) Validate(node types.Node) error {
	if node.Data.Mode == nil {
		return fmt.Errorf("visualization node missing mode")
	}
	return nil
}
