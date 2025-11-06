package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// GroupByExecutor groups array elements by field value with aggregation
type GroupByExecutor struct{}

// Execute groups the input array and aggregates
func (e *GroupByExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsGroupByData(node.Data)
	if err != nil {
		return nil, err
	}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("group_by node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("group_by node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	// Get grouping field
	field := ""
	if data.Field != nil {
		field = *data.Field
	}
	if field == "" {
		return nil, fmt.Errorf("group_by node requires 'field' string")
	}

	// Get aggregate function (default: count)
	aggregate := "count"
	if data.Aggregate != nil {
		aggregate = *data.Aggregate
	}

	// Get value field for numeric aggregations
	valueField := ""
	if data.ValueField != nil {
		valueField = *data.ValueField
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
	data, err := types.AsGroupByData(node.Data)
	if err != nil {
		return err
	}
	if data.Field == nil || *data.Field == "" {
		return fmt.Errorf("group_by node requires non-empty 'field'")
	}

	if data.Aggregate != nil {
		aggregate := *data.Aggregate
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
		if (aggregate == "sum" || aggregate == "avg" || aggregate == "min" || aggregate == "max") && (data.ValueField == nil || *data.ValueField == "") {
			return fmt.Errorf("group_by with '%s' aggregate requires 'value_field'", aggregate)
		}
	}

	return nil
}
