package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/expression"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// MapExecutor executes Map nodes
// Transforms each array element through an expression or child nodes
type MapExecutor struct{}

// Execute runs the Map node
// Maps (transforms) each array element to a new value.
//
// The Map node supports two transformation modes:
// 1. Expression-based: Use an expression to transform each element
// 2. Field extraction: Extract a specific field from each object
//
// The expression has access to:
// - `item` - the current array element
// - `index` - the current index (0-based)
// - `items` - the full input array
// - All workflow variables and context
//
// Examples:
//
//	Transform numbers: [1,2,3] → Map(expr="item * 2") → [2,4,6]
//	Extract field: [{name:"Alice"}] → Map(field="name") → ["Alice"]
//	Complex expression: [users] → Map(expr="item.age * 1.1") → [ages with 10% increase]
func (e *MapExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("map node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array (slice)
	inputArray, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("map node requires array input, got %T", input)
	}

	// Determine transformation mode
	hasExpression := node.Data.Expression != nil && *node.Data.Expression != ""
	hasField := node.Data.Field != nil && *node.Data.Field != ""

	if !hasExpression && !hasField {
		// No transformation specified - pass through with warning
		slog.Warn("map node has no expression or field specified, passing through array unchanged",
			slog.String("node_id", node.ID),
		)
		return map[string]interface{}{
			"results":      inputArray,
			"input_count":  len(inputArray),
			"output_count": len(inputArray),
			"warning":      "no transformation specified",
		}, nil
	}

	slog.Debug("map node starting",
		slog.String("node_id", node.ID),
		slog.Int("input_count", len(inputArray)),
		slog.Bool("has_expression", hasExpression),
		slog.Bool("has_field", hasField),
	)

	results := make([]interface{}, 0, len(inputArray))
	successful := 0
	failed := 0

	for i, item := range inputArray {
		var result interface{}
		var err error

		if hasField {
			// Field extraction mode
			result, err = e.extractField(item, *node.Data.Field)
		} else if hasExpression {
			// Expression transformation mode
			result, err = e.evaluateExpression(ctx, node, item, i, inputArray)
		}

		if err != nil {
			slog.Debug("map transformation error (continuing)",
				slog.String("node_id", node.ID),
				slog.Int("index", i),
				slog.String("error", err.Error()),
			)
			failed++
			// Continue on error - collect nil for failed transformation
			results = append(results, nil)
			continue
		}

		results = append(results, result)
		successful++
	}

	slog.Debug("map node completed",
		slog.String("node_id", node.ID),
		slog.Int("successful", successful),
		slog.Int("failed", failed),
	)

	return map[string]interface{}{
		"results":      results,
		"input_count":  len(inputArray),
		"output_count": len(results),
		"successful":   successful,
		"failed":       failed,
	}, nil
}

// extractField extracts a field from an object
func (e *MapExecutor) extractField(item interface{}, field string) (interface{}, error) {
	itemMap, ok := item.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot extract field from non-object: %T", item)
	}

	value, exists := itemMap[field]
	if !exists {
		return nil, fmt.Errorf("field '%s' not found in object", field)
	}

	return value, nil
}

// evaluateExpression evaluates a transformation expression for an item
func (e *MapExecutor) evaluateExpression(
	ctx ExecutionContext,
	node types.Node,
	item interface{},
	index int,
	items []interface{},
) (interface{}, error) {
	// Create expression context with item, index, and items
	exprCtx := &expression.Context{
		NodeResults: ctx.GetAllNodeResults(),
		Variables:   make(map[string]interface{}),
		ContextVars: ctx.GetContextVariables(),
	}

	// Copy existing variables
	for k, v := range ctx.GetVariables() {
		exprCtx.Variables[k] = v
	}

	// Add iteration variables
	exprCtx.Variables["item"] = item
	exprCtx.Variables["index"] = float64(index)
	exprCtx.Variables["items"] = items

	// Evaluate the expression
	// The expression should return a value, not a boolean
	// We use EvaluateExpression instead of Evaluate (which returns bool)
	result, err := expression.EvaluateExpression(*node.Data.Expression, item, exprCtx)
	if err != nil {
		return nil, fmt.Errorf("expression evaluation failed: %w", err)
	}

	return result, nil
}

// NodeType returns the node type this executor handles
func (e *MapExecutor) NodeType() types.NodeType {
	return types.NodeTypeMap
}

// Validate checks if node configuration is valid
func (e *MapExecutor) Validate(node types.Node) error {
	hasExpression := node.Data.Expression != nil && *node.Data.Expression != ""
	hasField := node.Data.Field != nil && *node.Data.Field != ""

	if !hasExpression && !hasField {
		return fmt.Errorf("map node requires either 'expression' or 'field' to be specified")
	}

	if hasExpression && hasField {
		return fmt.Errorf("map node cannot have both 'expression' and 'field' specified, choose one")
	}

	return nil
}
