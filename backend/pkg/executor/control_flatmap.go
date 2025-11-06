package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// FlatMapExecutor transforms elements to arrays and flattens the result
type FlatMapExecutor struct{}

// Execute flattens nested arrays
func (e *FlatMapExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsFlatMapData(node.Data)
	if err != nil {
		return nil, err
	}
	
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("flat_map node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("flat_map node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	// Get field to expand
	field := ""
	if data.Field != nil {
		field = *data.Field
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
	// Validate node data type
	if _, err := types.AsFlatMapData(node.Data); err != nil {
		return err
	}
	// Field is optional
	return nil
}
