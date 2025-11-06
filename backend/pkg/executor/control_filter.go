package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/expression"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// FilterExecutor executes Filter nodes
// Filters JSON array elements based on an expression
type FilterExecutor struct{}

// Execute runs the Filter node
// Filters array elements where the condition expression evaluates to true.
// If input is not an array, passes through the original input with a warning.
// The expression has access to the 'item' variable representing the current array element.
func (e *FilterExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsFilterData(node.Data)
if err != nil {
return nil, err
}
	if data.Condition == nil || *data.Condition == "" {
		return nil, fmt.Errorf("filter node missing condition expression")
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("filter node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array (slice)
	inputArray, ok := input.([]interface{})
	if !ok {
		// Try to extract array from common map structures (e.g., range node output)
		if inputMap, isMap := input.(map[string]interface{}); isMap {
			// Try common keys for arrays
			for _, key := range []string{"range", "array", "items", "data", "values"} {
				if arr, found := inputMap[key]; found {
					if arrSlice, isSlice := arr.([]interface{}); isSlice {
						inputArray = arrSlice
						ok = true
						break
					}
				}
			}
		}

		// If still not an array - pass through with warning
		if !ok {
			slog.Warn("filter node received non-array input, passing through unchanged",
				slog.String("node_id", node.ID),
				slog.String("input_type", fmt.Sprintf("%T", input)),
			)

			return map[string]interface{}{
				"input":         input,
				"filtered":      input,
				"is_array":      false,
				"warning":       "input is not an array, passed through unchanged",
				"original_type": fmt.Sprintf("%T", input),
			}, nil
		}
	}

	// Build expression context with access to node results and variables
	exprCtx := &expression.Context{
		NodeResults: ctx.GetAllNodeResults(),
		Variables:   ctx.GetVariables(),
		ContextVars: ctx.GetContextVariables(),
	}

	// Filter array elements
	filtered := make([]interface{}, 0, len(inputArray))
	skippedCount := 0
	errorCount := 0

	for i, item := range inputArray {
		// Create a temporary context with the current item
		// Make the item available as 'item' variable for the expression
		itemCtx := &expression.Context{
			NodeResults: exprCtx.NodeResults,
			Variables:   make(map[string]interface{}),
			ContextVars: exprCtx.ContextVars,
		}

		// Copy existing variables
		for k, v := range exprCtx.Variables {
			itemCtx.Variables[k] = v
		}

		// Add the current item as 'item' variable
		itemCtx.Variables["item"] = item

		// Also add index for potential use
		itemCtx.Variables["index"] = float64(i)

		// Evaluate the condition for this item
		conditionMet, err := expression.Evaluate(*data.Condition, item, itemCtx)
		if err != nil {
			// Log evaluation error but continue processing
			slog.Debug("filter expression evaluation error",
				slog.String("node_id", node.ID),
				slog.Int("item_index", i),
				slog.String("error", err.Error()),
			)
			errorCount++
			skippedCount++
			continue
		}

		// Include item if condition is met
		if conditionMet {
			filtered = append(filtered, item)
		} else {
			skippedCount++
		}
	}

	// Log filtering summary
	slog.Debug("filter node completed",
		slog.String("node_id", node.ID),
		slog.Int("input_count", len(inputArray)),
		slog.Int("output_count", len(filtered)),
		slog.Int("skipped_count", skippedCount),
		slog.Int("error_count", errorCount),
	)

	// Return filtered array with metadata
	return map[string]interface{}{
		"filtered":      filtered,
		"input_count":   len(inputArray),
		"output_count":  len(filtered),
		"skipped_count": skippedCount,
		"error_count":   errorCount,
		"condition":     *data.Condition,
		"is_array":      true,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *FilterExecutor) NodeType() types.NodeType {
	return types.NodeTypeFilter
}

// Validate checks if node configuration is valid
func (e *FilterExecutor) Validate(node types.Node) error {
data, err := types.AsFilterData(node.Data)
if err != nil {
return err
}
	if data.Condition == nil || *data.Condition == "" {
		return fmt.Errorf("filter node requires a condition expression")
	}
	return nil
}
