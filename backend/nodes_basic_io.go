package workflow

import "fmt"

// ============================================================================
// Basic I/O Node Executors
// ============================================================================
// This file contains executors for basic input/output nodes:
// - Number: Returns numeric values
// - TextInput: Returns text values
// - Visualization: Formats output for display
// ============================================================================

// executeNumberNode returns the numeric value from a number node.
// Number nodes provide constant numeric input values to workflows.
//
// Required fields:
//   - Data.Value: The numeric value to return
//
// Returns:
//   - float64: The node's numeric value
//   - error: If value field is missing
func (e *Engine) executeNumberNode(node Node) (interface{}, error) {
	if node.Data.Value == nil {
		return nil, fmt.Errorf("number node missing value")
	}
	return *node.Data.Value, nil
}

// executeTextInputNode returns the text value from a text input node.
// Text input nodes provide constant string input values to workflows.
//
// Required fields:
//   - Data.Text: The text value to return
//
// Returns:
//   - string: The node's text value
//   - error: If text field is missing
func (e *Engine) executeTextInputNode(node Node) (interface{}, error) {
	if node.Data.Text == nil {
		return nil, fmt.Errorf("text input node missing text")
	}
	return *node.Data.Text, nil
}

// executeVisualizationNode formats output for display.
// Visualization nodes are terminal nodes that format the final output.
//
// Required fields:
//   - Data.Mode: Display mode (e.g., "text", "table")
//
// Inputs:
//   - Requires at least one input from predecessor nodes
//
// Returns:
//   - map: Contains mode and the first input value
//   - error: If mode is missing or no inputs provided
func (e *Engine) executeVisualizationNode(node Node) (interface{}, error) {
	if node.Data.Mode == nil {
		return nil, fmt.Errorf("visualization node missing mode")
	}

	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("visualization needs at least 1 input")
	}

	return map[string]interface{}{
		"mode":  *node.Data.Mode,
		"value": inputs[0],
	}, nil
}
