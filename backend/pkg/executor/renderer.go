package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// RendererExecutor executes Renderer nodes
// This node acts as a pass-through in the backend - it just forwards the input data
// The actual rendering happens in the frontend
type RendererExecutor struct{}

// Execute passes through the input data unchanged
func (e *RendererExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
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
func (e *RendererExecutor) Validate(node types.Node) error {
	// Validate render_mode if provided
	if node.Data.RenderMode != nil {
		mode := *node.Data.RenderMode
		validModes := map[string]bool{
			"text":      true,
			"json":      true,
			"csv":       true,
			"tsv":       true,
			"xml":       true,
			"table":     true,
			"bar_chart": true,
		}
		if !validModes[mode] {
			return fmt.Errorf("invalid render_mode: %s. Valid modes: text, json, csv, tsv, xml, table, bar_chart", mode)
		}
	}
	return nil
}
