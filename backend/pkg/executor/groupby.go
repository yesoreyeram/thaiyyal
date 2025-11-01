package executor

import (
	"context"
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// GroupByExecutor groups array elements by field value with aggregation
type GroupByExecutor struct{}

// Execute groups the input array and aggregates
func (e *GroupByExecutor) Execute(ctx context.Context, node types.Node, inputs map[string]interface{}, nodeResults map[string]interface{}, variables map[string]interface{}) (interface{}, error) {
	// Get input array
	input, ok := inputs["in"]
	if !ok {
		return nil, fmt.Errorf("group_by node missing required input 'in'")
	}

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		return map[string]interface{}{
			"error": "input is not an array",
			"input": input,
		}, nil
	}

	// Get grouping field
	field := ""
	if node.Data.Field != nil {
		field = *node.Data.Field
	}
	if field == "" {
		return nil, fmt.Errorf("group_by node requires 'field' string")
	}

	// Get aggregate function (default: count)
	aggregate := "count"
	if node.Data.Aggregate != nil {
		aggregate = *node.Data.Aggregate
	}

	// Get value field for numeric aggregations
	valueField := ""
	if node.Data.ValueField != nil {
		valueField = *node.Data.ValueField
	}

	// Group items
	groups := make(map[string][]interface{})
	for _, item := range arr {
		var key string
		if obj, ok := item.(map[string]interface{}); ok {
			if val, exists := obj[field]; exists {
				key = fmt.Sprintf("%v", val)
			} else {
				key = "<missing>"
			}
		} else {
			key = fmt.Sprintf("%v", item)
		}

		groups[key] = append(groups[key], item)
	}

	// Perform aggregation
	aggregated := make(map[string]interface{})

	switch aggregate {
	case "count":
		for key, items := range groups {
			aggregated[key] = len(items)
		}

	case "sum":
		for key, items := range groups {
			sum := 0.0
			for _, item := range items {
				if obj, ok := item.(map[string]interface{}); ok {
					if val, exists := obj[valueField]; exists {
						if num, ok := val.(float64); ok {
							sum += num
						} else if num, ok := val.(int); ok {
							sum += float64(num)
						}
					}
				}
			}
			aggregated[key] = sum
		}

	case "avg":
		for key, items := range groups {
			sum := 0.0
			count := 0
			for _, item := range items {
				if obj, ok := item.(map[string]interface{}); ok {
					if val, exists := obj[valueField]; exists {
						if num, ok := val.(float64); ok {
							sum += num
							count++
						} else if num, ok := val.(int); ok {
							sum += float64(num)
							count++
						}
					}
				}
			}
			if count > 0 {
				aggregated[key] = sum / float64(count)
			} else {
				aggregated[key] = 0.0
			}
		}

	case "min":
		for key, items := range groups {
			var min *float64
			for _, item := range items {
				if obj, ok := item.(map[string]interface{}); ok {
					if val, exists := obj[valueField]; exists {
						var num float64
						if n, ok := val.(float64); ok {
							num = n
						} else if n, ok := val.(int); ok {
							num = float64(n)
						} else {
							continue
						}
						if min == nil || num < *min {
							min = &num
						}
					}
				}
			}
			if min != nil {
				aggregated[key] = *min
			} else {
				aggregated[key] = nil
			}
		}

	case "max":
		for key, items := range groups {
			var max *float64
			for _, item := range items {
				if obj, ok := item.(map[string]interface{}); ok {
					if val, exists := obj[valueField]; exists {
						var num float64
						if n, ok := val.(float64); ok {
							num = n
						} else if n, ok := val.(int); ok {
							num = float64(n)
						} else {
							continue
						}
						if max == nil || num > *max {
							max = &num
						}
					}
				}
			}
			if max != nil {
				aggregated[key] = *max
			} else {
				aggregated[key] = nil
			}
		}

	case "values":
		// Return all items in each group
		for key, items := range groups {
			aggregated[key] = items
		}

	default:
		return nil, fmt.Errorf("unknown aggregate function: %s", aggregate)
	}

	// Build result
	result := map[string]interface{}{
		"groups":      groups,
		"group_count": len(groups),
		"input_count": len(arr),
		"field":       field,
		"aggregate":   aggregate,
	}

	// Add aggregated results with appropriate key
	switch aggregate {
	case "count":
		result["counts"] = aggregated
	case "sum":
		result["sums"] = aggregated
	case "avg":
		result["averages"] = aggregated
	case "min":
		result["minimums"] = aggregated
	case "max":
		result["maximums"] = aggregated
	case "values":
		result["values"] = aggregated
	}

	return result, nil
}

// NodeType returns the node type this executor handles
func (e *GroupByExecutor) NodeType() types.NodeType {
	return types.NodeTypeGroupBy
}

// Validate checks if the node configuration is valid
func (e *GroupByExecutor) Validate(node types.Node) error {
	if node.Data.Field == nil || *node.Data.Field == "" {
		return fmt.Errorf("group_by node requires non-empty 'field'")
	}

	if node.Data.Aggregate != nil {
		aggregate := *node.Data.Aggregate
		validAgg := []string{"count", "sum", "avg", "min", "max", "values"}
		found := false
		for _, valid := range validAgg {
			if aggregate == valid {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("group_by aggregate must be one of: count, sum, avg, min, max, values")
		}

		// Check value_field for numeric aggregations
		if (aggregate == "sum" || aggregate == "avg" || aggregate == "min" || aggregate == "max") && (node.Data.ValueField == nil || *node.Data.ValueField == "") {
			return fmt.Errorf("group_by with '%s' aggregate requires 'value_field'", aggregate)
		}
	}

	return nil
}
