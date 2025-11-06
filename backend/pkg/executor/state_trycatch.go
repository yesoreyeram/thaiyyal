package executor

import (
	"errors"
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TryCatchExecutor executes TryCatch nodes
type TryCatchExecutor struct{}

// Execute runs the TryCatch node
// Implements error handling with fallback values
// Catches errors and provides fallback values or continues workflow execution
func (e *TryCatchExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsTryCatchData(node.Data)
	if err != nil {
		return nil, err
	}

	// Get configuration
	fallbackValue := data.FallbackValue

	continueOnError := true
	if data.ContinueOnError != nil {
		continueOnError = *data.ContinueOnError
	}

	errorOutputPath := ""
	if data.ErrorOutputPath != nil {
		errorOutputPath = *data.ErrorOutputPath
	}

	// Validate inputs
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, errors.New("try-catch node requires at least one input")
	}

	input := inputs[0]

	// Check if input indicates an error
	isError := false
	errorMsg := ""

	if errMap, ok := input.(map[string]interface{}); ok {
		if err, ok := errMap["error"]; ok {
			isError = true
			if errStr, ok := err.(string); ok {
				errorMsg = errStr
			} else {
				errorMsg = fmt.Sprintf("%v", err)
			}
		}
	}

	// If no error, pass through
	if !isError {
		return map[string]interface{}{
			"value":        input,
			"error_caught": false,
		}, nil
	}

	// Error detected - apply fallback strategy
	result := map[string]interface{}{
		"value":         fallbackValue,
		"error_caught":  true,
		"error_message": errorMsg,
	}

	if errorOutputPath != "" {
		result["error_output_path"] = errorOutputPath
	}

	// Decide whether to return error or continue
	if continueOnError {
		return result, nil
	}

	return result, fmt.Errorf("error caught: %s", errorMsg)
}

// NodeType returns the node type this executor handles
func (e *TryCatchExecutor) NodeType() types.NodeType {
	return types.NodeTypeTryCatch
}

// Validate checks if node configuration is valid
func (e *TryCatchExecutor) Validate(node types.Node) error {
	// Validate node data type
	if _, err := types.AsTryCatchData(node.Data); err != nil {
		return err
	}
	// No required fields for try-catch - all have defaults
	return nil
}
