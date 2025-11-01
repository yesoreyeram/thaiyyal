package executor

import (
	"context"
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// FlatMapExecutor transforms elements to arrays and flattens the result
type FlatMapExecutor struct{}

// Execute flattens nested arrays
func (e *FlatMapExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get input array
	input, ok := inputs["in"]
	if !ok {
		return nil, fmt.Errorf("flat_map node missing required input 'in'")
	}

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		return map[string]interface{}{
			"error": "input is not an array",
			"input": input,
		}, nil
	}

	// Get field to expand
	field := ""
	if node.Data.Field != nil {
		field = *node.Data.Field
	}

	// FlatMap the array
	var flattened []interface{}

	for _, item := range arr {
		var toFlatten []interface{}

		if field != "" && field != "." {
			// Extract field value
			if obj, ok := item.(map[string]interface{}); ok {
				if val, exists := obj[field]; exists {
					// Check if field value is an array
					if arr, ok := val.([]interface{}); ok {
						toFlatten = arr
					} else {
						// Single value - add as is
						toFlatten = []interface{}{val}
					}
				}
				// If field doesn't exist, skip this item
			} else {
				// Item is not an object, skip
				continue
			}
		} else {
			// No field or "." field - flatten item itself if it's an array
			if arr, ok := item.([]interface{}); ok {
				toFlatten = arr
			} else {
				// Single value
				toFlatten = []interface{}{item}
			}
		}

		// Add all elements from toFlatten
		flattened = append(flattened, toFlatten...)
	}

	return map[string]interface{}{
		"flattened":    flattened,
		"input_count":  len(arr),
		"output_count": len(flattened),
		"field":        field,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *FlatMapExecutor) NodeType() types.NodeType {
	return types.NodeTypeFlatMap
}

// Validate checks if the node configuration is valid
func (e *FlatMapExecutor) Validate(node types.Node) error {
	// Field is optional
	return nil
}
