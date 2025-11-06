package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ContextVariableExecutor executes ContextVariable nodes
type ContextVariableExecutor struct{}

// Execute runs the ContextVariable node
// Context variable nodes are orphan nodes that execute first and store
// their values in the context for later interpolation.
//
// Supports two formats:
// 1. Legacy: Single value with ContextName and ContextValue
// 2. New: Multiple typed values with ContextValues array
func (e *ContextVariableExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsContextVariableData(node.Data)
if err != nil {
return nil, err
}
	// New format: multiple typed values
	if len(data.ContextValues) > 0 {
		result := make(map[string]interface{})

		for _, cv := range data.ContextValues {
			// Convert value based on type
			convertedValue, err := convertTypedValue(cv.Value, cv.Type)
			if err != nil {
				return nil, fmt.Errorf("context_variable node: error converting %s (%s): %w", cv.Name, cv.Type, err)
			}

			// Store in context variables for interpolation
			ctx.SetContextVariable(cv.Name, convertedValue)
			result[cv.Name] = convertedValue
		}

		return map[string]interface{}{
			"type":      "variable",
			"variables": result,
		}, nil
	}

	// Legacy format: single value (backward compatibility)
	if data.ContextName == nil {
		return nil, fmt.Errorf("context_variable node missing context_name or context_values")
	}
	if data.ContextValue == nil {
		return nil, fmt.Errorf("context_variable node missing context_value")
	}

	varName := *data.ContextName
	varValue := data.ContextValue

	// Store in context variables for interpolation
	ctx.SetContextVariable(varName, varValue)

	return map[string]interface{}{
		"type":  "variable",
		"name":  varName,
		"value": varValue,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *ContextVariableExecutor) NodeType() types.NodeType {
	return types.NodeTypeContextVariable
}

// Validate checks if node configuration is valid
func (e *ContextVariableExecutor) Validate(node types.Node) error {
data, err := types.AsContextVariableData(node.Data)
if err != nil {
return err
}
	// Check if using new format
	if len(data.ContextValues) > 0 {
		return nil
	}
	// Check legacy format
	if data.ContextName == nil {
		return fmt.Errorf("context_variable node missing context_name or context_values")
	}
	if data.ContextValue == nil {
		return fmt.Errorf("context_variable node missing context_value")
	}
	return nil
}
