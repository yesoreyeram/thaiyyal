package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// DateInputExecutor executes DateInput nodes
type DateInputExecutor struct{}

// Execute returns the date value from a date input node
func (e *DateInputExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsDateInputData(node.Data)
if err != nil {
return nil, err
}
	if data.DateValue != nil {
		return *data.DateValue, nil
	}
	return nil, fmt.Errorf("date input node missing date value")
}

// NodeType returns the node type this executor handles
func (e *DateInputExecutor) NodeType() types.NodeType {
	return types.NodeTypeDateInput
}

// Validate checks if node configuration is valid
func (e *DateInputExecutor) Validate(node types.Node) error {
data, err := types.AsDateInputData(node.Data)
if err != nil {
return err
}
	if data.DateValue == nil {
		return fmt.Errorf("date input node missing date value")
	}
	return nil
}
