package executor

import (
	"fmt"
	"strings"
	"unicode"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TextOperationExecutor executes TextOperation nodes
type TextOperationExecutor struct{}

// Execute performs text transformations on string inputs
func (e *TextOperationExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.TextOp == nil {
		return nil, fmt.Errorf("text operation node missing text_op")
	}

	// Get input from predecessor node(s)
	inputs := ctx.GetNodeInputs(node.ID)
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

// NodeType returns the node type this executor handles
func (e *TextOperationExecutor) NodeType() types.NodeType {
	return types.NodeTypeTextOperation
}

// Validate checks if node configuration is valid
func (e *TextOperationExecutor) Validate(node types.Node) error {
	if node.Data.TextOp == nil {
		return fmt.Errorf("text operation node missing text_op")
	}
	return nil
}

// executeTextConcat concatenates multiple text inputs with an optional separator
func (e *TextOperationExecutor) executeTextConcat(inputs []interface{}, separator *string) (string, error) {
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

// executeTextRepeat repeats a text input a specified number of times
func (e *TextOperationExecutor) executeTextRepeat(input interface{}, repeatN *int) (string, error) {
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

// toTitleCase converts text to Title Case
func toTitleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

// toCamelCase converts text to camelCase
func toCamelCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		word := words[i]
		if len(word) > 0 {
			result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return result
}

// toInverseCase inverts the case of each character
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
