package executor

import (
"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// AccumulatorExecutor executes Accumulator nodes
type AccumulatorExecutor struct {
// Delegate to engine method
ExecuteFunc func(ctx ExecutionContext, node types.Node) (interface{}, error)
}

// Execute runs the Accumulator node
func (e *AccumulatorExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
if e.ExecuteFunc != nil {
return e.ExecuteFunc(ctx, node)
}
return nil, nil
}

// NodeType returns the node type this executor handles
func (e *AccumulatorExecutor) NodeType() types.NodeType {
return types.NodeTypeAccumulator
}

// Validate checks if node configuration is valid
func (e *AccumulatorExecutor) Validate(node types.Node) error {
// Basic validation - specific validation in Execute
return nil
}
