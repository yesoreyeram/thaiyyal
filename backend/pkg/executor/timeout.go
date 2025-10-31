package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TimeoutExecutor executes Timeout nodes
type TimeoutExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Timeout node
func (e *TimeoutExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *TimeoutExecutor) NodeType() types.NodeType {
return types.NodeTypeTimeout
}

// Validate checks if node configuration is valid
func (e *TimeoutExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
