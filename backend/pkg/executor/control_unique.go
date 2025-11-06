package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// UniqueExecutor removes duplicate elements from an array
type UniqueExecutor struct{}

// Execute removes duplicates from the input array
func (e *UniqueExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsUniqueData(node.Data)
	if err != nil {
		return nil, err
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("unique node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("unique node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	// Get field for uniqueness check (optional)
	field := ""
	if data.Field != nil {
		field = *data.Field
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
	// Validate node data type
	if _, err := types.AsUniqueData(node.Data); err != nil {
		return err
	}
	// No required fields - field is optional
	return nil
}
