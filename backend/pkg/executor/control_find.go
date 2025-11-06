package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/expression"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// FindExecutor finds the first element matching a condition
type FindExecutor struct{}

// Execute finds the first matching element in the array
func (e *FindExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsFindData(node.Data)
	if err != nil {
		return nil, err
	}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("find node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("find node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	// Get condition
	condition := ""
	if data.Condition != nil {
		condition = *data.Condition
	}
	if condition == "" {
		return nil, fmt.Errorf("find node missing required 'condition' string")
	}

	// Get return_index flag
	returnIndex := false
	if data.ReturnIndex != nil {
		returnIndex = *data.ReturnIndex
	}

	// Search for first match
	for i, item := range arr {
		// Create context with item and index variables
		itemCtx := &expression.Context{
			Variables:   make(map[string]interface{}),
			ContextVars: ctx.GetContextVariables(),
			NodeResults: ctx.GetAllNodeResults(),
		}
		// Copy existing variables
		for k, v := range ctx.GetVariables() {
			itemCtx.Variables[k] = v
		}
		itemCtx.Variables["item"] = item
		itemCtx.Variables["index"] = i
		itemCtx.Variables["items"] = arr

		// Evaluate condition
		result, err := expression.Evaluate(condition, item, itemCtx)
		if err != nil {
			// Continue on error
			continue
		}

		// Check if condition is true
		if result {
			// Found match
			if returnIndex {
				return map[string]interface{}{
					"found":     item,
					"index":     i,
					"condition": condition,
				}, nil
			}
			return map[string]interface{}{
				"found":     item,
				"condition": condition,
			}, nil
		}
	}

	// No match found
	return map[string]interface{}{
		"found":       nil,
		"not_found":   true,
		"condition":   condition,
		"input_count": len(arr),
	}, nil
}

// NodeType returns the node type this executor handles
func (e *FindExecutor) NodeType() types.NodeType {
	return types.NodeTypeFind
}

// Validate checks if the node configuration is valid
func (e *FindExecutor) Validate(node types.Node) error {
	data, err := types.AsFindData(node.Data)
	if err != nil {
		return err
	}
	if data.Condition == nil || *data.Condition == "" {
		return fmt.Errorf("find node requires non-empty 'condition' field")
	}
	return nil
}
