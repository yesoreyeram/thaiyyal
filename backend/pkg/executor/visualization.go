package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// VisualizationExecutor executes Visualization nodes
type VisualizationExecutor struct{}

// Execute acts as a pass-through, returning the input data directly to the frontend
// The frontend RendererNode will handle the visualization logic
func (e *VisualizationExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("visualization needs at least 1 input")
	}

	// Pass through the input data unchanged
	// The frontend will auto-detect the best rendering mode
	return inputs[0], nil
}

// NodeType returns the node type this executor handles
func (e *VisualizationExecutor) NodeType() types.NodeType {
	return types.NodeTypeVisualization
}

// Validate checks if node configuration is valid
func (e *VisualizationExecutor) Validate(node types.Node) error {
	return nil
}
