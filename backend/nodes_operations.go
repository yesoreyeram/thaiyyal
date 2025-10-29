package workflow

import (
	"fmt"
	"strings"
	"unicode"
)

// ============================================================================
// Operation Node Executors
// ============================================================================
// This file contains executors for operation nodes:
// - Operation: Arithmetic operations (add, subtract, multiply, divide)
// - TextOperation: Text transformations (uppercase, lowercase, etc.)
// ============================================================================

// executeOperationNode performs arithmetic operations on two numeric inputs.
// Supports: add, subtract, multiply, divide
//
// Required fields:
//   - Data.Op: Operation type ("add", "subtract", "multiply", "divide")
//
// Inputs:
//   - Requires exactly 2 numeric inputs from predecessor nodes
//
// Returns:
//   - float64: Result of the arithmetic operation
//   - error: If operation is invalid, inputs are missing/non-numeric, or division by zero
func (e *Engine) executeOperationNode(node Node) (interface{}, error) {
	if node.Data.Op == nil {
		return nil, fmt.Errorf("operation node missing op")
	}

	// Get inputs from predecessor nodes
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) < 2 {
		return nil, fmt.Errorf("operation needs 2 inputs, got %d", len(inputs))
	}

	// Convert to numbers
	left, ok1 := inputs[0].(float64)
	right, ok2 := inputs[1].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("operation inputs must be numbers")
	}

	// Perform operation using strategy pattern
	switch *node.Data.Op {
	case "add":
		return left + right, nil
	case "subtract":
		return left - right, nil
	case "multiply":
		return left * right, nil
	case "divide":
		if right == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return left / right, nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", *node.Data.Op)
	}
}

// executeTextOperationNode performs text transformations on string inputs.
// Supports: uppercase, lowercase, titlecase, camelcase, inversecase, concat, repeat
//
// Required fields:
//   - Data.TextOp: Operation type
//
// Optional fields:
//   - Data.Separator: For concat operation (default: "")
//   - Data.RepeatN: For repeat operation (required for repeat)
//
// Inputs:
//   - Most operations: 1 text input
//   - Concat: 1 or more text inputs
//
// Returns:
//   - string: Transformed text
//   - error: If operation is invalid or inputs are missing/non-text
func (e *Engine) executeTextOperationNode(node Node) (interface{}, error) {
	if node.Data.TextOp == nil {
		return nil, fmt.Errorf("text operation node missing text_op")
	}

	// Get input from predecessor node(s)
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("text operation needs at least 1 input")
	}

	// Handle concat operation (can accept multiple inputs)
	if *node.Data.TextOp == "concat" {
		return e.executeTextConcat(inputs, node.Data.Separator)
	}

	// Handle repeat operation
	if *node.Data.TextOp == "repeat" {
		return e.executeTextRepeat(inputs[0], node.Data.RepeatN)
	}

	// For other operations, validate single input is a string
	inputText, ok := inputs[0].(string)
	if !ok {
		return nil, fmt.Errorf("text operation input must be text/string")
	}

	// Perform text transformation using strategy pattern
	switch *node.Data.TextOp {
	case "uppercase":
		return strings.ToUpper(inputText), nil
	case "lowercase":
		return strings.ToLower(inputText), nil
	case "titlecase":
		return toTitleCase(inputText), nil
	case "camelcase":
		return toCamelCase(inputText), nil
	case "inversecase":
		return toInverseCase(inputText), nil
	default:
		return nil, fmt.Errorf("unknown text operation: %s", *node.Data.TextOp)
	}
}

// executeTextConcat concatenates multiple text inputs with an optional separator.
func (e *Engine) executeTextConcat(inputs []interface{}, separator *string) (string, error) {
	// Validate all inputs are strings
	textInputs := make([]string, 0, len(inputs))
	for i, input := range inputs {
		text, ok := input.(string)
		if !ok {
			return "", fmt.Errorf("concat operation input %d must be text/string", i)
		}
		textInputs = append(textInputs, text)
	}

	// Get separator (default to empty string)
	sep := ""
	if separator != nil {
		sep = *separator
	}

	// Concatenate all inputs
	result := strings.Join(textInputs, sep)
	return result, nil
}

// executeTextRepeat repeats a text input a specified number of times.
func (e *Engine) executeTextRepeat(input interface{}, repeatN *int) (string, error) {
	// Validate input is a string
	inputText, ok := input.(string)
	if !ok {
		return "", fmt.Errorf("repeat operation input must be text/string")
	}

	// Get repeat count (required)
	if repeatN == nil {
		return "", fmt.Errorf("repeat operation requires repeat_n field")
	}

	repeatCount := *repeatN
	if repeatCount < 0 {
		return "", fmt.Errorf("repeat_n must be non-negative, got %d", repeatCount)
	}

	// Repeat the text efficiently
	return strings.Repeat(inputText, repeatCount), nil
}

// ============================================================================
// Text Transformation Helper Functions
// ============================================================================

// toTitleCase converts text to Title Case (first letter of each word capitalized).
func toTitleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

// toCamelCase converts text to camelCase.
// Example: "hello world" → "helloWorld"
func toCamelCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		word := words[i]
		if len(word) > 0 {
			// Capitalize first letter, lowercase rest
			result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return result
}

// toInverseCase inverts the case of each character.
// Example: "Hello" → "hELLO"
func toInverseCase(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			runes[i] = unicode.ToLower(r)
		} else if unicode.IsLower(r) {
			runes[i] = unicode.ToUpper(r)
		}
	}
	return string(runes)
}
