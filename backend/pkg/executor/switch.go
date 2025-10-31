package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SwitchExecutor executes Switch nodes
type SwitchExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Switch node
func (e *SwitchExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *SwitchExecutor) NodeType() types.NodeType {
return types.NodeTypeSwitch
}

// Validate checks if node configuration is valid
func (e *SwitchExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
