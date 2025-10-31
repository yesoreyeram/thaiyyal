package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// CounterExecutor executes Counter nodes
type CounterExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Counter node
func (e *CounterExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *CounterExecutor) NodeType() types.NodeType {
return types.NodeTypeCounter
}

// Validate checks if node configuration is valid
func (e *CounterExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
