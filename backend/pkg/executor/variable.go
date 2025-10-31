package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// VariableExecutor executes Variable nodes
type VariableExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Variable node
func (e *VariableExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *VariableExecutor) NodeType() types.NodeType {
return types.NodeTypeVariable
}

// Validate checks if node configuration is valid
func (e *VariableExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
