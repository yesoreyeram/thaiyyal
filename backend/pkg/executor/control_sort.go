package executor

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SortExecutor sorts an array by field value
type SortExecutor struct{}

// Execute sorts the input array
func (e *SortExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("sort node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("sort node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	// Get field to sort by
	field := ""
	if node.Data.Field != nil {
		field = *node.Data.Field
	}

	// Get order (default: asc)
	order := "asc"
	if node.Data.Order != nil {
		order = strings.ToLower(*node.Data.Order)
	}

	// Make a copy to sort
	sorted := make([]interface{}, len(arr))
	copy(sorted, arr)

	// Sort the array using helper function
	sort.SliceStable(sorted, func(i, j int) bool {
		var vi, vj interface{}

		if field != "" {
			// Sort by field
			if obji, ok := sorted[i].(map[string]interface{}); ok {
				vi = obji[field]
			}
			if objj, ok := sorted[j].(map[string]interface{}); ok {
				vj = objj[field]
			}
		} else {
			// Sort by value itself
			vi = sorted[i]
			vj = sorted[j]
		}

		// Use helper function for comparison
		less := compareValuesHelper(vi, vj)

		if order == "desc" {
			return !less
		}
		return less
	})

	return map[string]interface{}{
		"sorted": sorted,
		"count":  len(sorted),
		"field":  field,
		"order":  order,
	}, nil
}

// compareValuesHelper compares two values with type awareness
func compareValuesHelper(a, b interface{}) bool {
	// Handle nil
	if a == nil && b == nil {
		return false
	}
	if a == nil {
		return true // nil comes first
	}
	if b == nil {
		return false
	}

	// Compare numbers
	if an, ok := a.(float64); ok {
		if bn, ok := b.(float64); ok {
			return an < bn
		}
	}
	if an, ok := a.(int); ok {
		if bn, ok := b.(int); ok {
			return an < bn
		}
	}

	// Compare strings
	if as, ok := a.(string); ok {
		if bs, ok := b.(string); ok {
			return as < bs
		}
	}

	// Compare booleans (false < true)
	if ab, ok := a.(bool); ok {
		if bb, ok := b.(bool); ok {
			return !ab && bb
		}
	}

	// Default: convert to string and compare
	return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
}

// NodeType returns the node type this executor handles
func (e *SortExecutor) NodeType() types.NodeType {
	return types.NodeTypeSort
}

// Validate checks if the node configuration is valid
func (e *SortExecutor) Validate(node types.Node) error {
	if node.Data.Order != nil {
		orderLower := strings.ToLower(*node.Data.Order)
		if orderLower != "asc" && orderLower != "desc" {
			return fmt.Errorf("sort order must be 'asc' or 'desc', got '%s'", *node.Data.Order)
		}
	}
	return nil
}
