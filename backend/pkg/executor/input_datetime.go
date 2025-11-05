package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// DateTimeInputExecutor executes DateTimeInput nodes
type DateTimeInputExecutor struct{}

// Execute returns the datetime value from a datetime input node
func (e *DateTimeInputExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.DateTimeValue != nil {
		return *node.Data.DateTimeValue, nil
	}
	return nil, fmt.Errorf("datetime input node missing datetime value")
}

// NodeType returns the node type this executor handles
func (e *DateTimeInputExecutor) NodeType() types.NodeType {
	return types.NodeTypeDateTimeInput
}

// Validate checks if node configuration is valid
func (e *DateTimeInputExecutor) Validate(node types.Node) error {
	if node.Data.DateTimeValue == nil {
		return fmt.Errorf("datetime input node missing datetime value")
	}
	return nil
}
