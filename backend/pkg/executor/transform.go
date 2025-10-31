package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TransformExecutor executes Transform nodes
type TransformExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Transform node
func (e *TransformExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *TransformExecutor) NodeType() types.NodeType {
return types.NodeTypeTransform
}

// Validate checks if node configuration is valid
func (e *TransformExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
