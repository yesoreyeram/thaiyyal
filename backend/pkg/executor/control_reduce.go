package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/expression"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ReduceExecutor executes Reduce nodes
// Reduces an array to a single value through accumulation
type ReduceExecutor struct{}

// Execute runs the Reduce node
// Reduces (aggregates) an array to a single value by accumulating across iterations.
//
// The Reduce node maintains an accumulator that is updated for each array element.
// The expression has access to:
// - `accumulator` - the current accumulated value
// - `item` - the current array element
// - `index` - the current index (0-based)
// - `items` - the full input array
// - All workflow variables and context
//
// Configuration:
// - `initial_value` (optional): Starting value for accumulator (default: 0)
// - `expression` (required): Expression to compute new accumulator value
//
// Examples:
//
//	Sum: [1,2,3] → Reduce(init=0, expr="accumulator + item") → 6
//	Sum ages: [{age:25},{age:30}] → Reduce(init=0, expr="accumulator + item.age") → 55
//	Max: [5,2,8,1] → Reduce(init=0, expr="item > accumulator ? item : accumulator") → 8
//	Concat: ["A","B","C"] → Reduce(init="", expr="accumulator + item") → "ABC"
func (e *ReduceExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsReduceData(node.Data)
if err != nil {
return nil, err
}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("reduce node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array (slice)
	inputArray, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("reduce node requires array input, got %T", input)
	}

	// Get initial value for accumulator
	accumulator := data.InitialValue
	if accumulator == nil {
		accumulator = float64(0) // Default to 0
	}

	if data.Expression == nil || *data.Expression == "" {
		return nil, fmt.Errorf("reduce node requires an expression")
	}

	slog.Debug("reduce node starting",
		slog.String("node_id", node.ID),
		slog.Int("input_count", len(inputArray)),
		slog.Any("initial_value", accumulator),
	)

	successful := 0
	failed := 0

	for i, item := range inputArray {
		result, err := e.evaluateExpression(ctx, node, *data.Expression, item, i, inputArray, accumulator)
		if err != nil {
			slog.Debug("reduce expression evaluation error (continuing)",
				slog.String("node_id", node.ID),
				slog.Int("index", i),
				slog.String("error", err.Error()),
			)
			failed++
			continue
		}

		// Update accumulator with the result
		accumulator = result
		successful++
	}

	slog.Debug("reduce node completed",
		slog.String("node_id", node.ID),
		slog.Int("successful", successful),
		slog.Int("failed", failed),
		slog.Any("final_value", accumulator),
	)

	return map[string]interface{}{
		"result":        accumulator,
		"initial_value": data.InitialValue,
		"final_value":   accumulator,
		"input_count":   len(inputArray),
		"iterations":    successful + failed,
		"successful":    successful,
		"failed":        failed,
	}, nil
}

// evaluateExpression evaluates the reduce expression for an item and accumulator
func (e *ReduceExecutor) evaluateExpression(
	ctx ExecutionContext,
	node types.Node,
	expressionStr string,
	item interface{},
	index int,
	items []interface{},
	accumulator interface{},
) (interface{}, error) {
	// Create expression context with accumulator, item, index, and items
	exprCtx := &expression.Context{
		NodeResults: ctx.GetAllNodeResults(),
		Variables:   make(map[string]interface{}),
		ContextVars: ctx.GetContextVariables(),
	}

	// Copy existing variables
	for k, v := range ctx.GetVariables() {
		exprCtx.Variables[k] = v
	}

	// Add iteration and accumulator variables
	exprCtx.Variables["accumulator"] = accumulator
	exprCtx.Variables["item"] = item
	exprCtx.Variables["index"] = float64(index)
	exprCtx.Variables["items"] = items

	// Evaluate the expression
	result, err := expression.EvaluateExpression(expressionStr, item, exprCtx)
	if err != nil {
		return nil, fmt.Errorf("expression evaluation failed: %w", err)
	}

	return result, nil
}

// NodeType returns the node type this executor handles
func (e *ReduceExecutor) NodeType() types.NodeType {
	return types.NodeTypeReduce
}

// Validate checks if node configuration is valid
func (e *ReduceExecutor) Validate(node types.Node) error {
data, err := types.AsReduceData(node.Data)
if err != nil {
return err
}
	if data.Expression == nil || *data.Expression == "" {
		return fmt.Errorf("reduce node requires an 'expression'")
	}

	return nil
}
