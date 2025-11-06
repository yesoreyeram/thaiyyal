package executor

import (
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// RendererExecutor executes Renderer nodes
// This node acts as a pass-through in the backend - it just forwards the input data
// The actual rendering happens in the frontend, which auto-detects the appropriate format
type RendererExecutor struct{}

// Execute passes through the input data unchanged
func (e *RendererExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsRendererData(node.Data)
if err != nil {
return nil, err
}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		// Return nil if no input - frontend will show "No data"
		return nil, nil
	}

	// Return the first input unchanged - this is a pass-through node
	return inputs[0], nil
}

// NodeType returns the node type this executor handles
func (e *RendererExecutor) NodeType() types.NodeType {
	return types.NodeTypeRenderer
}

// Validate checks if node configuration is valid
// Renderer nodes don't require any specific configuration - they auto-detect format
func (e *RendererExecutor) Validate(node types.Node) error {
data, err := types.AsRendererData(node.Data)
if err != nil {
return err
}
	// No validation needed - renderer auto-detects format from data
	return nil
}
