package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ContextConstantExecutor executes ContextConstant nodes
type ContextConstantExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the ContextConstant node
func (e *ContextConstantExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *ContextConstantExecutor) NodeType() types.NodeType {
return types.NodeTypeContextConstant
}

// Validate checks if node configuration is valid
func (e *ContextConstantExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
