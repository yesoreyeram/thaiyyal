package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TransformExecutor executes Transform nodes
type TransformExecutor struct{}

// Execute runs the Transform node
// Transforms data structures between different formats.
func (e *TransformExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.TransformType == nil {
		return nil, fmt.Errorf("transform node missing transform_type")
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("transform node requires input")
	}

	transformType := *node.Data.TransformType

	switch transformType {
	case "to_array":
		// Convert all inputs to array
		return inputs, nil

	case "to_object":
		// Convert array of alternating keys and values to object
		return transformToObject(inputs[0])

	case "flatten":
		// Recursively flatten nested arrays
		return transformFlatten(inputs[0])

	case "keys":
		// Extract all keys from object
		return transformKeys(inputs[0])

	case "values":
		// Extract all values from object
		return transformValues(inputs[0])

	default:
		return nil, fmt.Errorf("unsupported transform type: %s", transformType)
	}
}

// NodeType returns the node type this executor handles
func (e *TransformExecutor) NodeType() types.NodeType {
	return types.NodeTypeTransform
}

// Validate checks if node configuration is valid
func (e *TransformExecutor) Validate(node types.Node) error {
	if node.Data.TransformType == nil {
		return fmt.Errorf("transform node missing transform_type")
	}
	return nil
}

// transformToObject converts an array to an object.
// Array should contain alternating keys (strings) and values.
func transformToObject(input interface{}) (map[string]interface{}, error) {
	arr, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("to_object requires array input, got %T", input)
	}

	result := make(map[string]interface{})
	for i := 0; i < len(arr)-1; i += 2 {
		key, ok := arr[i].(string)
		if !ok {
			return nil, fmt.Errorf("to_object requires string keys at index %d", i)
		}
		result[key] = arr[i+1]
	}
	return result, nil
}

// transformFlatten recursively flattens nested arrays into a single-level array.
func transformFlatten(input interface{}) ([]interface{}, error) {
	arr, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("flatten requires array input, got %T", input)
	}

	var flattened []interface{}
	var flatten func(interface{})
	flatten = func(item interface{}) {
		if subArr, ok := item.([]interface{}); ok {
			// Recursively flatten nested arrays
			for _, sub := range subArr {
				flatten(sub)
			}
		} else {
			// Append non-array items directly
			flattened = append(flattened, item)
		}
	}

	for _, item := range arr {
		flatten(item)
	}
	return flattened, nil
}

// transformKeys extracts all keys from an object.
func transformKeys(input interface{}) ([]interface{}, error) {
	obj, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("keys transform requires object input, got %T", input)
	}

	keys := make([]interface{}, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	return keys, nil
}

// transformValues extracts all values from an object.
func transformValues(input interface{}) ([]interface{}, error) {
	obj, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("values transform requires object input, got %T", input)
	}

	values := make([]interface{}, 0, len(obj))
	for _, v := range obj {
		values = append(values, v)
	}
	return values, nil
}
