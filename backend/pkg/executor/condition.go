package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ConditionExecutor executes Condition nodes
type ConditionExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Condition node
func (e *ConditionExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *ConditionExecutor) NodeType() types.NodeType {
return types.NodeTypeCondition
}

// Validate checks if node configuration is valid
func (e *ConditionExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
