package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// DateInputExecutor executes DateInput nodes
type DateInputExecutor struct{}

// Execute returns the date value from a date input node
func (e *DateInputExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.DateValue != nil {
		return *node.Data.DateValue, nil
	}
	return nil, fmt.Errorf("date input node missing date value")
}

// NodeType returns the node type this executor handles
func (e *DateInputExecutor) NodeType() types.NodeType {
	return types.NodeTypeDateInput
}

// Validate checks if node configuration is valid
func (e *DateInputExecutor) Validate(node types.Node) error {
	if node.Data.DateValue == nil {
		return fmt.Errorf("date input node missing date value")
	}
	return nil
}
