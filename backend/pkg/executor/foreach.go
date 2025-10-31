package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ForEachExecutor executes ForEach nodes
type ForEachExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the ForEach node
func (e *ForEachExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *ForEachExecutor) NodeType() types.NodeType {
return types.NodeTypeForEach
}

// Validate checks if node configuration is valid
func (e *ForEachExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
