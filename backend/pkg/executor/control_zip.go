package executor

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ZipExecutor combines multiple arrays element-wise
type ZipExecutor struct{}

// Execute combines arrays element-wise into tuples
func (e *ZipExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsZipData(node.Data)
	if err != nil {
		return nil, err
	}
	
	// Get arrays to zip (can be from inputs or specified as array references)
	var arrays [][]interface{}

	inputs := ctx.GetNodeInputs(node.ID)

	// Check for direct input arrays
	if len(inputs) > 0 {
		if arr, ok := inputs[0].([]interface{}); ok {
			arrays = append(arrays, arr)
		}
	}

	// Check for additional arrays specified in config
	if arraysConfig, ok := data.Arrays.([]interface{}); ok {
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
	fillMissing := data.FillMissing

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
// Validate node data type
	if _, err := types.AsZipData(node.Data); err != nil {
		return err
	}
	// Arrays can be provided via inputs or config
	return nil
}

// CompactExecutor removes null and empty values from an array
type CompactExecutor struct{}

// Execute removes null, undefined, and optionally empty values
func (e *CompactExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsCompactData(node.Data)
	if err != nil {
		return nil, err
	}
	
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("compact node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("compact node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	// Get remove_empty option (default: false)
	removeEmpty := false
	if data.RemoveEmpty != nil {
		removeEmpty = *data.RemoveEmpty
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
	// Validate node data type
	if _, err := types.AsCompactData(node.Data); err != nil {
		return err
	}
	// No required fields
	return nil
}
