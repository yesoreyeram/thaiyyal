package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SplitExecutor executes Split nodes
type SplitExecutor struct{}

// Execute runs the Split node
// Handles splitting single input to multiple paths
func (e *SplitExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsSplitData(node.Data)
if err != nil {
return nil, err
}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("split node requires at least one input")
	}

	inputValue := inputs[0]
	paths := data.Paths
	if len(paths) == 0 {
		// Default to splitting to 2 paths
		paths = []string{"path1", "path2"}
	}

	// Create a copy of the input for each path
	outputs := make(map[string]interface{})
	for _, path := range paths {
		outputs[path] = inputValue
	}

	return map[string]interface{}{
		"value":   inputValue,
		"paths":   paths,
		"outputs": outputs,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *SplitExecutor) NodeType() types.NodeType {
	return types.NodeTypeSplit
}

// Validate checks if node configuration is valid
func (e *SplitExecutor) Validate(node types.Node) error {
	// Validate node data type
	if _, err := types.AsSplitData(node.Data); err != nil {
		return err
	}
	// No required fields for split
	return nil
}
