package executor

import (
	"context"
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SliceExecutor extracts a portion of an array
type SliceExecutor struct{}

// Execute extracts array slice based on start/end indices
func (e *SliceExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get input array
	input, ok := inputs["in"]
	if !ok {
		return nil, fmt.Errorf("slice node missing required input 'in'")
	}

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		return map[string]interface{}{
			"error": "input is not an array",
			"input": input,
		}, nil
	}

	arrLen := len(arr)

	// Get start index (default: 0)
	start := 0
	if startVal, ok := node.Data.Start.(float64); ok {
		start = int(startVal)
	} else if startVal, ok := node.Data.Start.(int); ok {
		start = startVal
	}

	// Handle negative start (from end)
	if start < 0 {
		start = arrLen + start
	}
	if start < 0 {
		start = 0
	}
	if start > arrLen {
		start = arrLen
	}

	// Get end index or length
	var end int
	hasEnd := false
	
	if endVal, ok := node.Data.End.(float64); ok {
		end = int(endVal)
		hasEnd = true
	} else if endVal, ok := node.Data.End.(int); ok {
		end = endVal
		hasEnd = true
	}

	// Check for length parameter
	if !hasEnd {
		if lengthVal, ok := node.Data.Length.(float64); ok {
			end = start + int(lengthVal)
			hasEnd = true
		} else if lengthVal, ok := node.Data.Length.(int); ok {
			end = start + lengthVal
			hasEnd = true
		}
	}

	// Default end: to the end of array
	if !hasEnd {
		end = arrLen
	}

	// Handle negative end (from end)
	if end < 0 {
		end = arrLen + end
	}
	if end < 0 {
		end = 0
	}
	if end > arrLen {
		end = arrLen
	}

	// Ensure start <= end
	if start > end {
		start = end
	}

	// Extract slice
	sliced := arr[start:end]

	return map[string]interface{}{
		"sliced":       sliced,
		"input_count":  arrLen,
		"output_count": len(sliced),
		"start":        start,
		"end":          end,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *SliceExecutor) NodeType() types.NodeType {
	return types.NodeTypeSlice
}

// Validate checks if the node configuration is valid
func (e *SliceExecutor) Validate(node types.Node) error {
	// All parameters are optional with sensible defaults
	return nil
}
