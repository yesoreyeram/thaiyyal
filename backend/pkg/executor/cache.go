package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// CacheExecutor executes Cache nodes
type CacheExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Cache node
func (e *CacheExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *CacheExecutor) NodeType() types.NodeType {
return types.NodeTypeCache
}

// Validate checks if node configuration is valid
func (e *CacheExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
