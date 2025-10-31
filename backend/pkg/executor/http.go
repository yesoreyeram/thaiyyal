package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// HTTPExecutor executes HTTP nodes
type HTTPExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the HTTP node
func (e *HTTPExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *HTTPExecutor) NodeType() types.NodeType {
return types.NodeTypeHTTP
}

// Validate checks if node configuration is valid
func (e *HTTPExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
