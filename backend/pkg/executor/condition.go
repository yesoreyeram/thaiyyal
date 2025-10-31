package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ConditionExecutor executes Condition nodes
type ConditionExecutor struct{}

// Execute runs the Condition node
// Evaluates a condition and passes through the input value
// with metadata about whether the condition was met.
func (e *ConditionExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.Condition == nil {
		return nil, fmt.Errorf("condition node missing condition")
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("condition node needs at least 1 input")
	}

	input := inputs[0]
	conditionMet := evaluateCondition(*node.Data.Condition, input)

	// Return the input value along with metadata about which path was taken
	return map[string]interface{}{
		"value":         input,
		"condition_met": conditionMet,
		"condition":     *node.Data.Condition,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *ConditionExecutor) NodeType() types.NodeType {
	return types.NodeTypeCondition
}

// Validate checks if node configuration is valid
func (e *ConditionExecutor) Validate(node types.Node) error {
	if node.Data.Condition == nil {
		return fmt.Errorf("condition node missing condition")
	}
	return nil
}
