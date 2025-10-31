package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ExtractExecutor executes Extract nodes
type ExtractExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Extract node
func (e *ExtractExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *ExtractExecutor) NodeType() types.NodeType {
return types.NodeTypeExtract
}

// Validate checks if node configuration is valid
func (e *ExtractExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
