package workflow

import "fmt"

// ============================================================================
// Control Flow Node Executors
// ============================================================================
// This file contains executors for control flow nodes:
// - Condition: Conditional branching based on comparisons
// - ForEach: Array iteration (simplified implementation)
// - WhileLoop: Conditional looping (simplified implementation)
// ============================================================================

// executeConditionNode evaluates a condition and passes through the input value
// with metadata about whether the condition was met.
//
// Required fields:
//   - Data.Condition: Condition expression (e.g., ">100", "<50", "==10", "true")
//
// Inputs:
//   - Requires at least one input value to evaluate
//
// Returns:
//   - map: Contains original value, condition_met boolean, and condition string
//   - error: If condition or input is missing
//
// Supported condition formats:
//   - ">N", "<N", ">=N", "<=N", "==N", "!=N" for numeric comparisons
//   - "true", "false" for boolean constants
func (e *Engine) executeConditionNode(node Node) (interface{}, error) {
	if node.Data.Condition == nil {
		return nil, fmt.Errorf("condition node missing condition")
	}

	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("condition node needs at least 1 input")
	}

	input := inputs[0]
	conditionMet := e.evaluateCondition(*node.Data.Condition, input)

	// Return the input value along with metadata about which path was taken
	return map[string]interface{}{
		"value":         input,
		"condition_met": conditionMet,
		"condition":     *node.Data.Condition,
	}, nil
}

// executeForEachNode iterates over an array input.
// This is a simplified implementation that validates the array and returns metadata.
// A full implementation would execute a sub-workflow for each array element.
//
// Optional fields:
//   - Data.MaxIterations: Maximum iterations allowed (default: 1000)
//
// Inputs:
//   - Requires one array input
//
// Returns:
//   - map: Contains items array, count, and iteration count
//   - error: If input is not an array or exceeds max iterations
func (e *Engine) executeForEachNode(node Node) (interface{}, error) {
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("for_each node needs at least 1 input")
	}

	// Check if input is an array (slice)
	inputArray, ok := inputs[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("for_each node requires array input, got %T", inputs[0])
	}

	// Set default max iterations
	maxIter := 1000
	if node.Data.MaxIterations != nil && *node.Data.MaxIterations > 0 {
		maxIter = *node.Data.MaxIterations
	}

	// Limit iterations to prevent resource exhaustion
	iterCount := len(inputArray)
	if iterCount > maxIter {
		return nil, fmt.Errorf("for_each exceeds max iterations: %d > %d", iterCount, maxIter)
	}

	// TODO: In a full implementation, execute sub-workflow for each element
	// For now, return metadata about the iteration
	return map[string]interface{}{
		"items":      inputArray,
		"count":      len(inputArray),
		"iterations": iterCount,
	}, nil
}

// executeWhileLoopNode executes a loop while a condition remains true.
// This is a simplified implementation that validates the condition.
// A full implementation would execute a sub-workflow on each iteration.
//
// Required fields:
//   - Data.Condition: Loop continuation condition
//
// Optional fields:
//   - Data.MaxIterations: Maximum iterations allowed (default: 100)
//
// Inputs:
//   - Requires at least one input value to evaluate
//
// Returns:
//   - map: Contains final value, iteration count, and condition
//   - error: If condition is missing, input missing, or max iterations exceeded
func (e *Engine) executeWhileLoopNode(node Node) (interface{}, error) {
	if node.Data.Condition == nil {
		return nil, fmt.Errorf("while_loop node missing condition")
	}

	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("while_loop node needs at least 1 input")
	}

	// Set default max iterations (lower than for_each to prevent infinite loops)
	maxIter := 100
	if node.Data.MaxIterations != nil && *node.Data.MaxIterations > 0 {
		maxIter = *node.Data.MaxIterations
	}

	currentValue := inputs[0]
	iterationCount := 0

	// Loop while condition is met (with safety limit)
	for e.evaluateCondition(*node.Data.Condition, currentValue) && iterationCount < maxIter {
		iterationCount++
		// TODO: In a full implementation, execute sub-workflow and update currentValue
		// For now, we just count iterations without modifying the value
	}

	if iterationCount >= maxIter {
		return nil, fmt.Errorf("while_loop exceeded max iterations: %d", maxIter)
	}

	return map[string]interface{}{
		"final_value": currentValue,
		"iterations":  iterationCount,
		"condition":   *node.Data.Condition,
	}, nil
}

// evaluateCondition evaluates a condition string against an input value.
//
// Supported condition formats:
//   - "true" - Always true
//   - "false" - Always false
//   - ">N" - Greater than N
//   - "<N" - Less than N
//   - ">=N" - Greater than or equal to N
//   - "<=N" - Less than or equal to N
//   - "==N" - Equal to N
//   - "!=N" - Not equal to N
//
// The value can be a direct number or a map containing a "value" field.
//
// Returns:
//   - bool: true if condition is met, false otherwise
func (e *Engine) evaluateCondition(condition string, value interface{}) bool {
	// Handle boolean constants
	if condition == "true" {
		return true
	}
	if condition == "false" {
		return false
	}

	// Extract numeric value from input
	numVal, ok := value.(float64)
	if !ok {
		// Try to extract value from map (common in node results)
		if m, isMap := value.(map[string]interface{}); isMap {
			if v, exists := m["value"]; exists {
				numVal, ok = v.(float64)
			}
		}
		if !ok {
			return false
		}
	}

	// Parse condition using a simple state machine
	var threshold float64
	var operator string

	if len(condition) >= 2 {
		// Check two-character operators first
		twoChar := condition[0:2]
		switch twoChar {
		case ">=":
			operator = ">="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		case "<=":
			operator = "<="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		case "==":
			operator = "=="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		case "!=":
			operator = "!="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		default:
			// Single-character operators
			switch condition[0] {
			case '>':
				operator = ">"
				fmt.Sscanf(condition[1:], "%f", &threshold)
			case '<':
				operator = "<"
				fmt.Sscanf(condition[1:], "%f", &threshold)
			}
		}
	}

	// Evaluate comparison using strategy pattern
	switch operator {
	case ">":
		return numVal > threshold
	case "<":
		return numVal < threshold
	case ">=":
		return numVal >= threshold
	case "<=":
		return numVal <= threshold
	case "==":
		return numVal == threshold
	case "!=":
		return numVal != threshold
	default:
		return false
	}
}
