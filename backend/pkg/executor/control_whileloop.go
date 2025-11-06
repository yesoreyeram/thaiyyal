package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// WhileLoopExecutor executes WhileLoop nodes
type WhileLoopExecutor struct{}

// Execute runs the WhileLoop node
// Executes a loop while a condition remains true.
// This is a simplified implementation that validates the condition.
// A full implementation would execute a sub-workflow on each iteration.
func (e *WhileLoopExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsWhileLoopData(node.Data)
if err != nil {
return nil, err
}
	if data.Condition == nil {
		return nil, fmt.Errorf("while_loop node missing condition")
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("while_loop node needs at least 1 input")
	}

	// Set default max iterations (lower than for_each to prevent infinite loops)
	maxIter := 100
	if data.MaxIterations != nil && *data.MaxIterations > 0 {
		maxIter = *data.MaxIterations
	}

	currentValue := inputs[0]
	iterationCount := 0

	// Loop while condition is met (with safety limit)
	for evaluateCondition(*data.Condition, currentValue) && iterationCount < maxIter {
		iterationCount++
		// TODO: In a full implementation, execute sub-workflow and update currentValue
		// For now, we just count iterations without modifying the value
	}

	if iterationCount >= maxIter {
		return nil, fmt.Errorf("while_loop exceeded max iterations: %d", maxIter)
	}

	return map[string]interface{}{
		"final_value": currentValue,
		"iterations":  iterationCount,
		"condition":   *data.Condition,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *WhileLoopExecutor) NodeType() types.NodeType {
	return types.NodeTypeWhileLoop
}

// Validate checks if node configuration is valid
func (e *WhileLoopExecutor) Validate(node types.Node) error {
data, err := types.AsWhileLoopData(node.Data)
if err != nil {
return err
}
	if data.Condition == nil {
		return fmt.Errorf("while_loop node missing condition")
	}
	return nil
}
