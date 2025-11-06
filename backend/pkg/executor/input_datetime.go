package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// DateTimeInputExecutor executes DateTimeInput nodes
type DateTimeInputExecutor struct{}

// Execute returns the datetime value from a datetime input node
func (e *DateTimeInputExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsDateTimeInputData(node.Data)
if err != nil {
return nil, err
}
	if data.DateTimeValue != nil {
		return *data.DateTimeValue, nil
	}
	return nil, fmt.Errorf("datetime input node missing datetime value")
}

// NodeType returns the node type this executor handles
func (e *DateTimeInputExecutor) NodeType() types.NodeType {
	return types.NodeTypeDateTimeInput
}

// Validate checks if node configuration is valid
func (e *DateTimeInputExecutor) Validate(node types.Node) error {
data, err := types.AsDateTimeInputData(node.Data)
if err != nil {
return err
}
	if data.DateTimeValue == nil {
		return fmt.Errorf("datetime input node missing datetime value")
	}
	return nil
}
