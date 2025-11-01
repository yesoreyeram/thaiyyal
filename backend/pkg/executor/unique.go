package executor

import (
	"context"
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// UniqueExecutor removes duplicate elements from an array
type UniqueExecutor struct{}

// Execute removes duplicates from the input array
func (e *UniqueExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get input array
	input, ok := inputs["in"]
	if !ok {
		return nil, fmt.Errorf("unique node missing required input 'in'")
	}

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		return map[string]interface{}{
			"error": "input is not an array",
			"input": input,
		}, nil
	}

	// Get field for uniqueness check (optional)
	field := ""
	if node.Data.Field != nil {
		field = *node.Data.Field
	}

	// Track seen values
	seen := make(map[string]bool)
	var unique []interface{}

	for _, item := range arr {
		var key string
		
		if field != "" {
			// Use specific field for uniqueness
			if obj, ok := item.(map[string]interface{}); ok {
				if val, exists := obj[field]; exists {
					key = fmt.Sprintf("%v", val)
				} else {
					// Item missing field, treat as unique
					key = fmt.Sprintf("%p", item)
				}
			} else {
				key = fmt.Sprintf("%v", item)
			}
		} else {
			// Use whole item for uniqueness
			key = fmt.Sprintf("%v", item)
		}

		if !seen[key] {
			seen[key] = true
			unique = append(unique, item)
		}
	}

	return map[string]interface{}{
		"unique":       unique,
		"input_count":  len(arr),
		"output_count": len(unique),
		"removed":      len(arr) - len(unique),
	}, nil
}

// NodeType returns the node type this executor handles
func (e *UniqueExecutor) NodeType() types.NodeType {
	return types.NodeTypeUnique
}

// Validate checks if the node configuration is valid
func (e *UniqueExecutor) Validate(node types.Node) error {
	// No required fields - field is optional
	return nil
}
