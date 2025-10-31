package executor

import (
	"fmt"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// OperationExecutor executes arithmetic Operation nodes
type OperationExecutor struct{}

// Execute performs arithmetic operations on two numeric inputs
func (e *OperationExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.Op == nil {
		return nil, fmt.Errorf("operation node missing op")
	}

	// Get inputs from predecessor nodes
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) < 2 {
		return nil, fmt.Errorf("operation needs 2 inputs, got %d", len(inputs))
	}

	// Convert to numbers
	left, ok1 := inputs[0].(float64)
	right, ok2 := inputs[1].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("operation inputs must be numbers")
	}

	// Perform operation using strategy pattern
	switch *node.Data.Op {
	case "add":
		return left + right, nil
	case "subtract":
		return left - right, nil
	case "multiply":
		return left * right, nil
	case "divide":
		if right == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return left / right, nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", *node.Data.Op)
	}
}

// NodeType returns the node type this executor handles
func (e *OperationExecutor) NodeType() types.NodeType {
	return types.NodeTypeOperation
}

// Validate checks if node configuration is valid
func (e *OperationExecutor) Validate(node types.Node) error {
	if node.Data.Op == nil {
		return fmt.Errorf("operation node missing op")
	}
	return nil
}
