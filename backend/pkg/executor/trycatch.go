package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TryCatchExecutor executes TryCatch nodes
type TryCatchExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the TryCatch node
func (e *TryCatchExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *TryCatchExecutor) NodeType() types.NodeType {
return types.NodeTypeTryCatch
}

// Validate checks if node configuration is valid
func (e *TryCatchExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
