package executor

import (
	"context"
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/expression"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// FindExecutor finds the first element matching a condition
type FindExecutor struct{}

// Execute finds the first matching element in the array
func (e *FindExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get input array
	input, ok := inputs["in"]
	if !ok {
		return nil, fmt.Errorf("find node missing required input 'in'")
	}

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		return map[string]interface{}{
			"error": "input is not an array",
			"input": input,
		}, nil
	}

	// Get condition
	condition := ""
	if node.Data.Condition != nil {
		condition = *node.Data.Condition
	}
	if condition == "" {
		return nil, fmt.Errorf("find node missing required 'condition' string")
	}

	// Get return_index flag
	returnIndex := false
	if node.Data.ReturnIndex != nil {
		returnIndex = *node.Data.ReturnIndex
	}

	// Search for first match
	for i, item := range arr {
		// Create context with item and index variables
		itemCtx := &expression.Context{
			Variables:   make(map[string]interface{}),
			ContextVars: map[string]interface{}{},
			NodeResults: nodeResults,
		}
		// Copy existing variables
		for k, v := range variables {
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
					"found": item,
					"index": i,
					"condition": condition,
				}, nil
			}
			return map[string]interface{}{
				"found": item,
				"condition": condition,
			}, nil
		}
	}

	// No match found
	return map[string]interface{}{
		"found":      nil,
		"not_found":  true,
		"condition":  condition,
		"input_count": len(arr),
	}, nil
}

// NodeType returns the node type this executor handles
func (e *FindExecutor) NodeType() types.NodeType {
	return types.NodeTypeFind
}

// Validate checks if the node configuration is valid
func (e *FindExecutor) Validate(node types.Node) error {
	if node.Data.Condition == nil || *node.Data.Condition == "" {
		return fmt.Errorf("find node requires non-empty 'condition' field")
	}
	return nil
}
