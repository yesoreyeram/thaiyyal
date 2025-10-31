package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// DelayExecutor executes Delay nodes
type DelayExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Delay node
func (e *DelayExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *DelayExecutor) NodeType() types.NodeType {
return types.NodeTypeDelay
}

// Validate checks if node configuration is valid
func (e *DelayExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
