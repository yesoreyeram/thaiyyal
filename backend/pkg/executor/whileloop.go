package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// WhileLoopExecutor executes WhileLoop nodes
type WhileLoopExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the WhileLoop node
func (e *WhileLoopExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *WhileLoopExecutor) NodeType() types.NodeType {
return types.NodeTypeWhileLoop
}

// Validate checks if node configuration is valid
func (e *WhileLoopExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
