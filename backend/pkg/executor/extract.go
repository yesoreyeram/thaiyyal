package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ExtractExecutor executes Extract nodes
type ExtractExecutor struct{}

// Execute runs the Extract node
// Extracts specific fields from object inputs.
func (e *ExtractExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsExtractData(node.Data)
if err != nil {
return nil, err
}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("extract node requires input")
	}

	input := inputs[0]
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("extract node requires object input, got %T", input)
	}

	// Single field extraction
	if data.Field != nil {
		field := *data.Field
		value, exists := inputMap[field]
		if !exists {
			return nil, fmt.Errorf("field '%s' not found in input object", field)
		}
		return map[string]interface{}{
			"field": field,
			"value": value,
		}, nil
	}

	// Multiple fields extraction
	if len(data.Fields) > 0 {
		result := make(map[string]interface{})
		for _, field := range data.Fields {
			value, exists := inputMap[field]
			if exists {
				result[field] = value
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("extract node requires 'field' or 'fields' configuration")
}

// NodeType returns the node type this executor handles
func (e *ExtractExecutor) NodeType() types.NodeType {
	return types.NodeTypeExtract
}

// Validate checks if node configuration is valid
func (e *ExtractExecutor) Validate(node types.Node) error {
data, err := types.AsExtractData(node.Data)
if err != nil {
return err
}
	if data.Field == nil && len(data.Fields) == 0 {
		return fmt.Errorf("extract node requires 'field' or 'fields' configuration")
	}
	return nil
}
