package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TextInputExecutor executes TextInput nodes
type TextInputExecutor struct{}

// Execute returns the text value from a text input node
func (e *TextInputExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.Text == nil {
		return nil, fmt.Errorf("text input node missing text")
	}
	return *node.Data.Text, nil
}

// NodeType returns the node type this executor handles
func (e *TextInputExecutor) NodeType() types.NodeType {
	return types.NodeTypeTextInput
}

// Validate checks if node configuration is valid
func (e *TextInputExecutor) Validate(node types.Node) error {
	if node.Data.Text == nil {
		return fmt.Errorf("text input node missing text")
	}
	return nil
}
