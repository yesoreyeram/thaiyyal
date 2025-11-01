package executor

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SampleExecutor gets a sample from an array
type SampleExecutor struct{}

// Execute gets a sample from the input array
func (e *SampleExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get input array
	input, ok := inputs["in"]
	if !ok {
		return nil, fmt.Errorf("sample node missing required input 'in'")
	}

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		return map[string]interface{}{
			"error": "input is not an array",
			"input": input,
		}, nil
	}

	// Get count
	count := 1 // default
	if countVal, ok := node.Data.Count.(float64); ok {
		count = int(countVal)
	} else if countVal, ok := node.Data.Count.(int); ok {
		count = countVal
	}

	if count <= 0 {
		return nil, fmt.Errorf("sample count must be greater than 0, got %d", count)
	}

	// Get method
	method := "random" // default
	if node.Data.Method != nil {
		method = *node.Data.Method
	}

	var sample []interface{}

	switch method {
	case "first":
		// Take first N items
		end := count
		if end > len(arr) {
			end = len(arr)
		}
		sample = arr[:end]

	case "last":
		// Take last N items
		start := len(arr) - count
		if start < 0 {
			start = 0
		}
		sample = arr[start:]

	case "random":
		// Take random N items
		if count >= len(arr) {
			// Return all items
			sample = make([]interface{}, len(arr))
			copy(sample, arr)
		} else {
			// Fisher-Yates shuffle to get random sample
			indices := make([]int, len(arr))
			for i := range indices {
				indices[i] = i
			}
			
			// Initialize random seed
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			
			// Partial shuffle - only shuffle first 'count' elements
			for i := 0; i < count; i++ {
				j := i + rng.Intn(len(indices)-i)
				indices[i], indices[j] = indices[j], indices[i]
			}

			// Get sampled items
			sample = make([]interface{}, count)
			for i := 0; i < count; i++ {
				sample[i] = arr[indices[i]]
			}
		}

	default:
		return nil, fmt.Errorf("unknown sample method: %s (must be 'random', 'first', or 'last')", method)
	}

	return map[string]interface{}{
		"sample":       sample,
		"input_count":  len(arr),
		"output_count": len(sample),
		"method":       method,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *SampleExecutor) NodeType() types.NodeType {
	return types.NodeTypeSample
}

// Validate checks if the node configuration is valid
func (e *SampleExecutor) Validate(node types.Node) error {
	if node.Data.Count != nil {
		count := 0
		if countVal, ok := node.Data.Count.(float64); ok {
			count = int(countVal)
		} else if countVal, ok := node.Data.Count.(int); ok {
			count = countVal
		}
		if count <= 0 {
			return fmt.Errorf("sample count must be greater than 0, got %d", count)
		}
	}

	if node.Data.Method != nil {
		method := *node.Data.Method
		if method != "random" && method != "first" && method != "last" {
			return fmt.Errorf("sample method must be 'random', 'first', or 'last', got %s", method)
		}
	}

	return nil
}
