package workflow

import (
	"fmt"
	"regexp"
)

// ============================================================================
// Context Node Executors
// ============================================================================
// This file contains executors for context nodes that define workflow-level
// variables and constants. These nodes are typically orphan nodes (no inputs)
// that execute first and populate the context for use in other nodes.
//
// Context nodes support template interpolation in other nodes using:
// - {{ variable.name }} - for mutable variables
// - {{ const.name }} - for immutable constants
// ============================================================================

// executeContextVariableNode defines a workflow-level mutable variable.
// Context variable nodes are orphan nodes that execute first and store
// their values in the engine's contextVariables map for later interpolation.
//
// Required fields:
//   - Data.ContextName: Name of the variable
//   - Data.ContextValue: Value of the variable
//
// Returns:
//   - interface{}: The variable value
//   - error: If required fields are missing
func (e *Engine) executeContextVariableNode(node Node) (interface{}, error) {
	if node.Data.ContextName == nil {
		return nil, fmt.Errorf("context_variable node missing context_name")
	}
	if node.Data.ContextValue == nil {
		return nil, fmt.Errorf("context_variable node missing context_value")
	}

	varName := *node.Data.ContextName
	varValue := node.Data.ContextValue

	// Store in context variables for interpolation
	e.contextVariables[varName] = varValue

	return map[string]interface{}{
		"type":  "variable",
		"name":  varName,
		"value": varValue,
	}, nil
}

// executeContextConstantNode defines a workflow-level immutable constant.
// Context constant nodes are orphan nodes that execute first and store
// their values in the engine's contextConstants map for later interpolation.
//
// Required fields:
//   - Data.ContextName: Name of the constant
//   - Data.ContextValue: Value of the constant
//
// Returns:
//   - interface{}: The constant value
//   - error: If required fields are missing
func (e *Engine) executeContextConstantNode(node Node) (interface{}, error) {
	if node.Data.ContextName == nil {
		return nil, fmt.Errorf("context_constant node missing context_name")
	}
	if node.Data.ContextValue == nil {
		return nil, fmt.Errorf("context_constant node missing context_value")
	}

	constName := *node.Data.ContextName
	constValue := node.Data.ContextValue

	// Store in context constants for interpolation
	e.contextConstants[constName] = constValue

	return map[string]interface{}{
		"type":  "const",
		"name":  constName,
		"value": constValue,
	}, nil
}

// ============================================================================
// Template Interpolation
// ============================================================================

// templateRegex matches {{ variable.name }} or {{ const.name }}
var templateRegex = regexp.MustCompile(`\{\{\s*(variable|const)\.(\w+)\s*\}\}`)

// interpolateTemplate replaces template placeholders in a string with actual values from context
func (e *Engine) interpolateTemplate(text string) string {
	if len(e.contextVariables) == 0 && len(e.contextConstants) == 0 {
		return text
	}

	// Replace all template placeholders
	result := templateRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the type and name from the match
		parts := templateRegex.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match // Return original if parsing fails
		}

		contextType := parts[1]
		varName := parts[2]

		// Look up the value in the appropriate context map
		var value interface{}
		var exists bool

		if contextType == "variable" {
			value, exists = e.contextVariables[varName]
		} else if contextType == "const" {
			value, exists = e.contextConstants[varName]
		}

		if exists {
			return fmt.Sprintf("%v", value)
		}

		// Return original if not found
		return match
	})

	return result
}

// interpolateValue recursively interpolates templates in various data types
func (e *Engine) interpolateValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return e.interpolateTemplate(v)
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = e.interpolateValue(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = e.interpolateValue(val)
		}
		return result
	default:
		return value
	}
}

// interpolateNodeData interpolates all string fields in NodeData
func (e *Engine) interpolateNodeData(data *NodeData) {
	if len(e.contextVariables) == 0 && len(e.contextConstants) == 0 {
		return
	}

	// Interpolate string pointer fields
	if data.Text != nil {
		interpolated := e.interpolateTemplate(*data.Text)
		data.Text = &interpolated
	}
	if data.URL != nil {
		interpolated := e.interpolateTemplate(*data.URL)
		data.URL = &interpolated
	}
	if data.Label != nil {
		interpolated := e.interpolateTemplate(*data.Label)
		data.Label = &interpolated
	}
	if data.VarName != nil {
		interpolated := e.interpolateTemplate(*data.VarName)
		data.VarName = &interpolated
	}
	if data.Field != nil {
		interpolated := e.interpolateTemplate(*data.Field)
		data.Field = &interpolated
	}
	if data.CacheKey != nil {
		interpolated := e.interpolateTemplate(*data.CacheKey)
		data.CacheKey = &interpolated
	}

	// Interpolate string arrays
	if len(data.Fields) > 0 {
		for i, field := range data.Fields {
			data.Fields[i] = e.interpolateTemplate(field)
		}
	}
	if len(data.Paths) > 0 {
		for i, path := range data.Paths {
			data.Paths[i] = e.interpolateTemplate(path)
		}
	}

	// Interpolate interface{} fields that might contain strings
	if data.InitialValue != nil {
		data.InitialValue = e.interpolateValue(data.InitialValue)
	}
	if data.FallbackValue != nil {
		data.FallbackValue = e.interpolateValue(data.FallbackValue)
	}
}
