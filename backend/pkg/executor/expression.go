package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/expression"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ExpressionExecutor executes Expression nodes
// Applies a user-provided expression to the input and returns the result
type ExpressionExecutor struct{}

// Execute runs the Expression node
// Applies a custom expression to transform the input data.
//
// The Expression node allows users to apply arbitrary transformations
// to input data using an expression language. This is useful for:
// - Simple calculations: input * 2, input + 10
// - Comparisons: input > 100, input == 5, input <= 0
// - Data transformations: input.field * 1.1
// - Conditional logic: input > 100 ? "high" : "low"
// - Complex operations: (input.price * input.quantity) * 1.08
//
// The expression has access to:
// - `input` - the input data from the predecessor node
// - All workflow variables and context
// - All node results from the workflow
//
// Examples:
//
//	Double a number: 5 → Expression("input * 2") → 10
//	Check threshold: 150 → Expression("input > 100") → true
//	Calculate total: {price:10, qty:3} → Expression("input.price * input.qty") → 30
//	Conditional: 150 → Expression("input > 100 ? 'high' : 'low'") → "high"
//	Field comparison: {age:25} → Expression("input.age > 18") → true
func (e *ExpressionExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsExpressionData(node.Data)
	if err != nil {
		return nil, err
	}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("expression node needs at least 1 input")
	}

	input := inputs[0]

	// Check if expression is provided
	if data.Expression == nil || *data.Expression == "" {
		slog.Warn("expression node has no expression specified, passing through input unchanged",
			slog.String("node_id", node.ID),
		)
		return map[string]interface{}{
			"result":  input,
			"warning": "no expression specified",
		}, nil
	}

	expr := *data.Expression

	slog.Debug("expression node starting",
		slog.String("node_id", node.ID),
		slog.String("expression", expr),
	)

	// Create expression context with input
	exprCtx := &expression.Context{
		NodeResults: ctx.GetAllNodeResults(),
		Variables:   make(map[string]interface{}),
		ContextVars: ctx.GetContextVariables(),
	}

	// Copy existing variables
	for k, v := range ctx.GetVariables() {
		exprCtx.Variables[k] = v
	}

	// Add 'input' to the expression context
	exprCtx.Variables["input"] = input

	// Try to evaluate as a value expression first (arithmetic, field access, etc.)
	result, err := expression.EvaluateExpression(expr, input, exprCtx)
	if err != nil {
		// If value expression fails, try as a boolean expression (comparisons)
		// This handles expressions like "input > 2", "input == 5", etc.
		boolResult, boolErr := expression.Evaluate(expr, input, exprCtx)
		if boolErr != nil {
			slog.Error("expression evaluation failed",
				slog.String("node_id", node.ID),
				slog.String("expression", expr),
				slog.String("value_error", err.Error()),
				slog.String("boolean_error", boolErr.Error()),
			)
			return nil, fmt.Errorf("expression evaluation failed: %w", err)
		}
		// Successfully evaluated as boolean
		result = boolResult
	}

	slog.Debug("expression node completed",
		slog.String("node_id", node.ID),
		slog.Any("result", result),
	)

	return result, nil
}

// NodeType returns the type of this executor
func (e *ExpressionExecutor) NodeType() types.NodeType {
	return types.NodeTypeExpression
}

// Validate validates the Expression node configuration
func (e *ExpressionExecutor) Validate(node types.Node) error {
	data, err := types.AsExpressionData(node.Data)
	if err != nil {
		return err
	}
	if node.Type != types.NodeTypeExpression {
		return fmt.Errorf("invalid node type: expected %s, got %s", types.NodeTypeExpression, node.Type)
	}

	if data.Expression == nil || *data.Expression == "" {
		return fmt.Errorf("expression is required")
	}

	// Basic validation - check expression is not too long
	if len(*data.Expression) > 10000 {
		return fmt.Errorf("expression is too long (max 10000 characters)")
	}

	return nil
}
