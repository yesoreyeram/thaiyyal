package executor

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// RetryExecutor executes Retry nodes
type RetryExecutor struct{}

// Execute runs the Retry node
// Implements retry logic with configurable backoff strategies
// Retries failed operations automatically with exponential, linear, or constant backoff
func (e *RetryExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsRetryData(node.Data)
if err != nil {
return nil, err
}
	// Get retry configuration with defaults
	maxAttempts := 3
	if data.MaxAttempts != nil {
		maxAttempts = *data.MaxAttempts
	}

	backoffStrategy := "exponential"
	if data.BackoffStrategy != nil {
		backoffStrategy = *data.BackoffStrategy
	}

	initialDelay := 1 * time.Second
	if data.InitialDelay != nil {
		if d, err := parseDuration(*data.InitialDelay); err == nil {
			initialDelay = d
		}
	}

	maxDelay := 30 * time.Second
	if data.MaxDelay != nil {
		if d, err := parseDuration(*data.MaxDelay); err == nil {
			maxDelay = d
		}
	}

	multiplier := 2.0
	if data.Multiplier != nil {
		multiplier = *data.Multiplier
	}

	// Get retry_on_errors patterns (optional)
	retryOnErrors := data.RetryOnErrors

	// Validate inputs
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, errors.New("retry node requires at least one input")
	}

	input := inputs[0]
	var lastError error
	currentDelay := initialDelay

	// Retry loop
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// For this implementation, we treat the input as the result
		// In a real workflow, this would re-execute a previous node
		// Here we simulate by checking if input is an error type or contains error field

		// Check if input indicates an error
		isError := false
		errorMsg := ""

		if errMap, ok := input.(map[string]interface{}); ok {
			if err, ok := errMap["error"]; ok {
				isError = true
				if errStr, ok := err.(string); ok {
					errorMsg = errStr
				}
			}
		}

		// If not an error, return successfully
		if !isError {
			return map[string]interface{}{
				"value":    input,
				"attempts": attempt,
				"success":  true,
			}, nil
		}

		// Check if this error should be retried
		shouldRetry := len(retryOnErrors) == 0 // If no patterns, retry all errors
		for _, pattern := range retryOnErrors {
			if strings.Contains(errorMsg, pattern) {
				shouldRetry = true
				break
			}
		}

		if !shouldRetry {
			lastError = fmt.Errorf("error not in retry list: %s", errorMsg)
			break
		}

		// If last attempt, don't delay
		if attempt == maxAttempts {
			lastError = fmt.Errorf("max retry attempts (%d) reached: %s", maxAttempts, errorMsg)
			break
		}

		// Calculate delay based on strategy
		var delay time.Duration
		switch backoffStrategy {
		case "exponential":
			delay = time.Duration(float64(initialDelay) * math.Pow(multiplier, float64(attempt-1)))
		case "linear":
			delay = initialDelay * time.Duration(attempt)
		case "constant":
			delay = initialDelay
		default:
			delay = initialDelay
		}

		// Cap at max delay
		if delay > maxDelay {
			delay = maxDelay
		}

		currentDelay = delay

		// Sleep before retry (in real implementation, would be async)
		time.Sleep(delay)

		lastError = fmt.Errorf("%s", errorMsg)
	}

	// All retries failed
	return map[string]interface{}{
		"value":      input,
		"attempts":   maxAttempts,
		"success":    false,
		"error":      lastError.Error(),
		"last_delay": currentDelay.String(),
	}, lastError
}

// NodeType returns the node type this executor handles
func (e *RetryExecutor) NodeType() types.NodeType {
	return types.NodeTypeRetry
}

// Validate checks if node configuration is valid
func (e *RetryExecutor) Validate(node types.Node) error {
	// Validate node data type
	if _, err := types.AsRetryData(node.Data); err != nil {
		return err
	}
	// No required fields for retry - all have defaults
	return nil
}
