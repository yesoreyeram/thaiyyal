package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// BooleanInputExecutor executes BooleanInput nodes
type BooleanInputExecutor struct{}

// Execute returns the boolean value from a boolean input node
func (e *BooleanInputExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	// The frontend sends the boolean value in the "value" field as interface{}
	// We need to handle it appropriately
	if node.Data.BooleanValue != nil {
		return *node.Data.BooleanValue, nil
	}
	
	// Fallback: try to get from generic value field (frontend compatibility)
	// The frontend might send it as a generic value
	return false, nil // Default to false if no value is set
}

// NodeType returns the node type this executor handles
func (e *BooleanInputExecutor) NodeType() types.NodeType {
	return types.NodeTypeBooleanInput
}

// Validate checks if node configuration is valid
func (e *BooleanInputExecutor) Validate(node types.Node) error {
	// Boolean value is optional, defaults to false
	return nil
}
