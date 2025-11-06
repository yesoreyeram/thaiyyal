package executor

import (
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// BooleanInputExecutor executes BooleanInput nodes
type BooleanInputExecutor struct{}

// Execute returns the boolean value from a boolean input node
func (e *BooleanInputExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsBooleanInputData(node.Data)
if err != nil {
return nil, err
}
	// Get the boolean value from the node data
	if data.BooleanValue != nil {
		return *data.BooleanValue, nil
	}

	// Default to false if no value is set
	return false, nil
}

// NodeType returns the node type this executor handles
func (e *BooleanInputExecutor) NodeType() types.NodeType {
	return types.NodeTypeBooleanInput
}

// Validate checks if node configuration is valid
func (e *BooleanInputExecutor) Validate(node types.Node) error {
	// Validate node data type
	if _, err := types.AsBooleanInputData(node.Data); err != nil {
		return err
	}
	// Boolean value is optional, defaults to false
	return nil
}
