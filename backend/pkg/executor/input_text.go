package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TextInputExecutor executes TextInput nodes
type TextInputExecutor struct{}

// Execute returns the text value from a text input node
func (e *TextInputExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsTextInputData(node.Data)
	if err != nil {
		return nil, err
	}
	if data.Text == nil {
		return nil, fmt.Errorf("text input node missing text")
	}
	return *data.Text, nil
}

// NodeType returns the node type this executor handles
func (e *TextInputExecutor) NodeType() types.NodeType {
	return types.NodeTypeTextInput
}

// Validate checks if node configuration is valid
func (e *TextInputExecutor) Validate(node types.Node) error {
	data, err := types.AsTextInputData(node.Data)
	if err != nil {
		return err
	}
	if data.Text == nil {
		return fmt.Errorf("text input node missing text")
	}
	return nil
}
