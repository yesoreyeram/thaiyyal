package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// AccumulatorExecutor executes Accumulator nodes
type AccumulatorExecutor struct{}

// Execute runs the Accumulator node
// Accumulates values over successive calls.
// The accumulator maintains state across multiple node executions in a workflow.
func (e *AccumulatorExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.AccumOp == nil {
		return nil, fmt.Errorf("accumulator node missing accum_op")
	}

	accumOp := *node.Data.AccumOp
	inputs := ctx.GetNodeInputs(node.ID)

	// Initialize accumulator with appropriate default or configured initial value
	currentAccum := ctx.GetAccumulator()
	if currentAccum == nil {
		currentAccum = getAccumulatorInitialValue(accumOp, node.Data.InitialValue)
		ctx.SetAccumulator(currentAccum)
	}

	// If no inputs, return current accumulator state
	if len(inputs) == 0 {
		return map[string]interface{}{
			"operation": accumOp,
			"value":     currentAccum,
		}, nil
	}

	// Accumulate the input value
	newAccum, err := accumulateValue(accumOp, currentAccum, inputs[0])
	if err != nil {
		return nil, err
	}
	ctx.SetAccumulator(newAccum)

	return map[string]interface{}{
		"operation": accumOp,
		"value":     newAccum,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *AccumulatorExecutor) NodeType() types.NodeType {
	return types.NodeTypeAccumulator
}

// Validate checks if node configuration is valid
func (e *AccumulatorExecutor) Validate(node types.Node) error {
	if node.Data.AccumOp == nil {
		return fmt.Errorf("accumulator node missing accum_op")
	}
	return nil
}

// getAccumulatorInitialValue returns the appropriate initial value for an accumulator.
func getAccumulatorInitialValue(accumOp string, configuredValue interface{}) interface{} {
	if configuredValue != nil {
		return configuredValue
	}

	// Default initial values based on operation
	switch accumOp {
	case "sum", "count":
		return 0.0
	case "product":
		return 1.0
	case "concat":
		return ""
	case "array":
		return []interface{}{}
	default:
		return nil
	}
}

// accumulateValue adds an input value to the accumulator based on the operation.
func accumulateValue(accumOp string, accum interface{}, input interface{}) (interface{}, error) {
	switch accumOp {
	case "sum":
		return accumulateSum(accum, input)
	case "product":
		return accumulateProduct(accum, input)
	case "concat":
		return accumulateConcat(accum, input)
	case "array":
		return accumulateArray(accum, input)
	case "count":
		return accumulateCount(accum)
	default:
		return nil, fmt.Errorf("unsupported accumulator operation: %s", accumOp)
	}
}

// accumulateSum adds a numeric value to the accumulator.
func accumulateSum(accum interface{}, input interface{}) (interface{}, error) {
	accumVal, ok := accum.(float64)
	if !ok {
		return nil, fmt.Errorf("accumulator value is not a number")
	}
	num, ok := input.(float64)
	if !ok {
		return nil, fmt.Errorf("sum accumulator requires numeric input, got %T", input)
	}
	return accumVal + num, nil
}

// accumulateProduct multiplies the accumulator by a numeric value.
func accumulateProduct(accum interface{}, input interface{}) (interface{}, error) {
	accumVal, ok := accum.(float64)
	if !ok {
		return nil, fmt.Errorf("accumulator value is not a number")
	}
	num, ok := input.(float64)
	if !ok {
		return nil, fmt.Errorf("product accumulator requires numeric input, got %T", input)
	}
	return accumVal * num, nil
}

// accumulateConcat concatenates a string to the accumulator.
func accumulateConcat(accum interface{}, input interface{}) (interface{}, error) {
	accumVal, ok := accum.(string)
	if !ok {
		return nil, fmt.Errorf("accumulator value is not a string")
	}
	str, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("concat accumulator requires string input, got %T", input)
	}
	return accumVal + str, nil
}

// accumulateArray appends a value to the accumulator array.
func accumulateArray(accum interface{}, input interface{}) (interface{}, error) {
	accumVal, ok := accum.([]interface{})
	if !ok {
		return nil, fmt.Errorf("accumulator value is not an array")
	}
	return append(accumVal, input), nil
}

// accumulateCount increments the counter.
func accumulateCount(accum interface{}) (interface{}, error) {
	accumVal, ok := accum.(float64)
	if !ok {
		return nil, fmt.Errorf("accumulator value is not a number")
	}
	return accumVal + 1, nil
}
