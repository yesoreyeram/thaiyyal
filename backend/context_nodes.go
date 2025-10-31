package workflow

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
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
//
// NEW: Context nodes now support multiple typed values for better type safety
// and conversion (string, number, boolean, time_string, epoch_second, epoch_ms)
// ============================================================================

// executeContextVariableNode defines workflow-level mutable variables.
// Context variable nodes are orphan nodes that execute first and store
// their values in the engine's contextVariables map for later interpolation.
//
// Supports two formats:
// 1. Legacy: Single value with ContextName and ContextValue
// 2. New: Multiple typed values with ContextValues array
//
// Required fields (legacy):
//   - Data.ContextName: Name of the variable
//   - Data.ContextValue: Value of the variable
//
// Required fields (new):
//   - Data.ContextValues: Array of {name, value, type}
//
// Returns:
//   - interface{}: Map of variable names to converted values
//   - error: If required fields are missing or type conversion fails
func (e *Engine) executeContextVariableNode(node Node) (interface{}, error) {
	// New format: multiple typed values
	if len(node.Data.ContextValues) > 0 {
		result := make(map[string]interface{})
		
		for _, cv := range node.Data.ContextValues {
			// Convert value based on type
			convertedValue, err := convertTypedValue(cv.Value, cv.Type)
			if err != nil {
				return nil, fmt.Errorf("context_variable node: error converting %s (%s): %w", cv.Name, cv.Type, err)
			}
			
			// Store in context variables for interpolation
			e.contextVariables[cv.Name] = convertedValue
			result[cv.Name] = convertedValue
		}
		
		return map[string]interface{}{
			"type":      "variable",
			"variables": result,
		}, nil
	}
	
	// Legacy format: single value (backward compatibility)
	if node.Data.ContextName == nil {
		return nil, fmt.Errorf("context_variable node missing context_name or context_values")
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

// executeContextConstantNode defines workflow-level immutable constants.
// Context constant nodes are orphan nodes that execute first and store
// their values in the engine's contextConstants map for later interpolation.
//
// Supports two formats:
// 1. Legacy: Single value with ContextName and ContextValue
// 2. New: Multiple typed values with ContextValues array
//
// Required fields (legacy):
//   - Data.ContextName: Name of the constant
//   - Data.ContextValue: Value of the constant
//
// Required fields (new):
//   - Data.ContextValues: Array of {name, value, type}
//
// Returns:
//   - interface{}: Map of constant names to converted values
//   - error: If required fields are missing or type conversion fails
func (e *Engine) executeContextConstantNode(node Node) (interface{}, error) {
	// New format: multiple typed values
	if len(node.Data.ContextValues) > 0 {
		result := make(map[string]interface{})
		
		for _, cv := range node.Data.ContextValues {
			// Convert value based on type
			convertedValue, err := convertTypedValue(cv.Value, cv.Type)
			if err != nil {
				return nil, fmt.Errorf("context_constant node: error converting %s (%s): %w", cv.Name, cv.Type, err)
			}
			
			// Store in context constants for interpolation
			e.contextConstants[cv.Name] = convertedValue
			result[cv.Name] = convertedValue
		}
		
		return map[string]interface{}{
			"type":      "const",
			"constants": result,
		}, nil
	}
	
	// Legacy format: single value (backward compatibility)
	if node.Data.ContextName == nil {
		return nil, fmt.Errorf("context_constant node missing context_name or context_values")
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
// Type Conversion
// ============================================================================

// convertTypedValue converts a value to the specified type.
// Supported types:
//   - "string": Convert to string
//   - "number": Convert to float64
//   - "boolean": Convert to bool
//   - "time_string": Parse as RFC3339 time string, return as string
//   - "epoch_second": Parse as Unix epoch seconds, return as time.Time
//   - "epoch_ms": Parse as Unix epoch milliseconds, return as time.Time
//   - "null": Return nil
func convertTypedValue(value interface{}, valueType string) (interface{}, error) {
	switch valueType {
	case "string":
		return fmt.Sprintf("%v", value), nil
		
	case "number":
		// Try to convert to float64
		switch v := value.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %q to number: %w", v, err)
			}
			return f, nil
		default:
			return nil, fmt.Errorf("cannot convert type %T to number", value)
		}
		
	case "boolean":
		// Try to convert to bool
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %q to boolean: %w", v, err)
			}
			return b, nil
		case float64:
			return v != 0, nil
		default:
			return nil, fmt.Errorf("cannot convert type %T to boolean", value)
		}
		
	case "time_string":
		// Parse as RFC3339 time string
		var timeStr string
		switch v := value.(type) {
		case string:
			timeStr = v
		default:
			timeStr = fmt.Sprintf("%v", value)
		}
		
		// Validate it's a valid time string
		_, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid time string %q: %w", timeStr, err)
		}
		return timeStr, nil
		
	case "epoch_second":
		// Convert to Unix epoch seconds, return as time.Time
		var seconds int64
		switch v := value.(type) {
		case float64:
			seconds = int64(v)
		case int:
			seconds = int64(v)
		case int64:
			seconds = v
		case string:
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %q to epoch seconds: %w", v, err)
			}
			seconds = i
		default:
			return nil, fmt.Errorf("cannot convert type %T to epoch seconds", value)
		}
		return time.Unix(seconds, 0), nil
		
	case "epoch_ms":
		// Convert to Unix epoch milliseconds, return as time.Time
		var ms int64
		switch v := value.(type) {
		case float64:
			ms = int64(v)
		case int:
			ms = int64(v)
		case int64:
			ms = v
		case string:
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %q to epoch milliseconds: %w", v, err)
			}
			ms = i
		default:
			return nil, fmt.Errorf("cannot convert type %T to epoch milliseconds", value)
		}
		return time.Unix(ms/1000, (ms%1000)*1000000), nil
		
	case "null":
		return nil, nil
		
	default:
		return nil, fmt.Errorf("unsupported type %q", valueType)
	}
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
