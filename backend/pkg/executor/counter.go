package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// CounterExecutor executes Counter nodes
type CounterExecutor struct{}

// Execute runs the Counter node
// Handles counter operations (increment, decrement, reset, get).
// The counter maintains a single numeric value across workflow execution.
func (e *CounterExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.CounterOp == nil {
		return nil, fmt.Errorf("counter node missing counter_op")
	}

	counterOp := *node.Data.CounterOp
	currentCounter := ctx.GetCounter()

	// Initialize counter if configured
	if node.Data.InitialValue != nil {
		if val, ok := node.Data.InitialValue.(float64); ok {
			currentCounter = val
			ctx.SetCounter(currentCounter)
		}
	}

	// Execute counter operation
	switch counterOp {
	case "increment":
		delta := 1.0
		if node.Data.Delta != nil {
			delta = *node.Data.Delta
		}
		currentCounter += delta
		ctx.SetCounter(currentCounter)

	case "decrement":
		delta := 1.0
		if node.Data.Delta != nil {
			delta = *node.Data.Delta
		}
		currentCounter -= delta
		ctx.SetCounter(currentCounter)

	case "reset":
		resetValue := 0.0
		if node.Data.InitialValue != nil {
			if val, ok := node.Data.InitialValue.(float64); ok {
				resetValue = val
			}
		}
		currentCounter = resetValue
		ctx.SetCounter(currentCounter)

	case "get":
		// Just return current counter value (no modification)

	default:
		return nil, fmt.Errorf("unsupported counter operation: %s (use increment, decrement, reset, or get)", counterOp)
	}

	return map[string]interface{}{
		"operation": counterOp,
		"value":     currentCounter,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *CounterExecutor) NodeType() types.NodeType {
	return types.NodeTypeCounter
}

// Validate checks if node configuration is valid
func (e *CounterExecutor) Validate(node types.Node) error {
	if node.Data.CounterOp == nil {
		return fmt.Errorf("counter node missing counter_op")
	}
	return nil
}
