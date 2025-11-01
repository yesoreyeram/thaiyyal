package executor

import (
	"context"
	"fmt"
	"math"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ZipExecutor combines multiple arrays element-wise
type ZipExecutor struct{}

// Execute combines arrays element-wise into tuples
func (e *ZipExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get arrays to zip (can be from inputs or specified as array references)
	var arrays [][]interface{}
	
	// Check for direct input arrays
	if in, ok := inputs["in"]; ok {
		if arr, ok := in.([]interface{}); ok {
			arrays = append(arrays, arr)
		}
	}
	
	// Check for additional arrays specified in config
	if arraysConfig, ok := node.Data.Arrays.([]interface{}); ok {
		for _, arrRef := range arraysConfig {
			// TODO: Resolve array references from node results/variables
			// For now, we expect arrays to be passed directly
			if arr, ok := arrRef.([]interface{}); ok {
				arrays = append(arrays, arr)
			}
		}
	}

	if len(arrays) == 0 {
		return nil, fmt.Errorf("zip node requires at least one array")
	}

	// Get fill value for shorter arrays
	fillMissing := node.Data.FillMissing

	// Find maximum length
	maxLen := 0
	for _, arr := range arrays {
		if len(arr) > maxLen {
			maxLen = len(arr)
		}
	}

	// Zip arrays
	var zipped []interface{}
	for i := 0; i < maxLen; i++ {
		tuple := make([]interface{}, len(arrays))
		for j, arr := range arrays {
			if i < len(arr) {
				tuple[j] = arr[i]
			} else {
				tuple[j] = fillMissing
			}
		}
		zipped = append(zipped, tuple)
	}

	return map[string]interface{}{
		"zipped":      zipped,
		"array_count": len(arrays),
		"tuple_count": len(zipped),
		"max_length":  maxLen,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *ZipExecutor) NodeType() types.NodeType {
	return types.NodeTypeZip
}

// Validate checks if the node configuration is valid
func (e *ZipExecutor) Validate(node types.Node) error {
	// Arrays can be provided via inputs or config
	return nil
}

// CompactExecutor removes null and empty values from an array
type CompactExecutor struct{}

// Execute removes null, undefined, and optionally empty values
func (e *CompactExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get input array
	input, ok := inputs["in"]
	if !ok {
		return nil, fmt.Errorf("compact node missing required input 'in'")
	}

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		return map[string]interface{}{
			"error": "input is not an array",
			"input": input,
		}, nil
	}

	// Get remove_empty option (default: false)
	removeEmpty := false
	if node.Data.RemoveEmpty != nil {
		removeEmpty = *node.Data.RemoveEmpty
	}

	// Compact the array
	var compacted []interface{}
	removedCount := 0

	for _, item := range arr {
		shouldRemove := false

		// Remove null/nil
		if item == nil {
			shouldRemove = true
		}

		// Remove NaN
		if f, ok := item.(float64); ok && math.IsNaN(f) {
			shouldRemove = true
		}

		// Optionally remove empty strings
		if removeEmpty {
			if str, ok := item.(string); ok && str == "" {
				shouldRemove = true
			}
		}

		if !shouldRemove {
			compacted = append(compacted, item)
		} else {
			removedCount++
		}
	}

	return map[string]interface{}{
		"compacted":    compacted,
		"input_count":  len(arr),
		"output_count": len(compacted),
		"removed":      removedCount,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *CompactExecutor) NodeType() types.NodeType {
	return types.NodeTypeCompact
}

// Validate checks if the node configuration is valid
func (e *CompactExecutor) Validate(node types.Node) error {
	// No required fields
	return nil
}
