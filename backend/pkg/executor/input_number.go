package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// NumberExecutor executes Number nodes
type NumberExecutor struct{}

// Execute returns the numeric value from a number node
func (e *NumberExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.Value == nil {
		return nil, fmt.Errorf("number node missing value")
	}
	return *node.Data.Value, nil
}

// NodeType returns the node type this executor handles
func (e *NumberExecutor) NodeType() types.NodeType {
	return types.NodeTypeNumber
}

// Validate checks if node configuration is valid
func (e *NumberExecutor) Validate(node types.Node) error {
	if node.Data.Value == nil {
		return fmt.Errorf("number node missing value")
	}
	return nil
}
