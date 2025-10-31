package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SplitExecutor executes Split nodes
type SplitExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Split node
func (e *SplitExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *SplitExecutor) NodeType() types.NodeType {
return types.NodeTypeSplit
}

// Validate checks if node configuration is valid
func (e *SplitExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
