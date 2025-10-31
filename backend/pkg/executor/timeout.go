package executor

import (
	"errors"
	"fmt"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TimeoutExecutor executes Timeout nodes
type TimeoutExecutor struct{}

// Execute runs the Timeout node
// Enforces time limits on operations
// Returns partial results or error if operation exceeds timeout
func (e *TimeoutExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	// Get timeout configuration
	timeoutDuration := 30 * time.Second
	if node.Data.Timeout != nil {
		if d, err := parseDuration(*node.Data.Timeout); err == nil {
			timeoutDuration = d
		} else {
			return nil, fmt.Errorf("invalid timeout duration: %s", *node.Data.Timeout)
		}
	}

	timeoutAction := "error" // "error" or "continue_with_partial"
	if node.Data.TimeoutAction != nil {
		timeoutAction = *node.Data.TimeoutAction
	}

	// Validate inputs
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, errors.New("timeout node requires at least one input")
	}

	input := inputs[0]

	// Simulate timeout check - in real implementation, this would wrap
	// the execution of a previous node with a timeout context

	// Check if input contains a duration field to simulate execution time
	executionTime := 0 * time.Second
	timedOut := false

	if inputMap, ok := input.(map[string]interface{}); ok {
		if et, ok := inputMap["execution_time"].(string); ok {
			if d, err := parseDuration(et); err == nil {
				executionTime = d
				if executionTime > timeoutDuration {
					timedOut = true
				}
			}
		}
	}

	// Build result
	result := map[string]interface{}{
		"value":            input,
		"timeout_duration": timeoutDuration.String(),
		"execution_time":   executionTime.String(),
		"timed_out":        timedOut,
	}

	if timedOut {
		result["timeout_exceeded"] = true

		if timeoutAction == "error" {
			return result, fmt.Errorf("operation timed out after %s (limit: %s)", executionTime, timeoutDuration)
		}

		// continue_with_partial
		result["partial_result"] = true
		return result, nil
	}

	// No timeout
	return result, nil
}

// NodeType returns the node type this executor handles
func (e *TimeoutExecutor) NodeType() types.NodeType {
	return types.NodeTypeTimeout
}

// Validate checks if node configuration is valid
func (e *TimeoutExecutor) Validate(node types.Node) error {
	// No required fields for timeout - all have defaults
	return nil
}
