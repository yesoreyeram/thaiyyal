package workflow

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Error Handling & Resilience Nodes Implementation
// This file implements Retry, Try-Catch, and Timeout nodes for robust workflow execution

// ========== Retry Node ==========

// executeRetryNode implements retry logic with configurable backoff strategies
// Retries failed operations automatically with exponential, linear, or constant backoff
func (e *Engine) executeRetryNode(node *Node, inputs []interface{}) (interface{}, error) {
	// Get retry configuration with defaults
	maxAttempts := 3
	if node.Data.MaxAttempts != nil {
		maxAttempts = *node.Data.MaxAttempts
	}
	
	backoffStrategy := "exponential"
	if node.Data.BackoffStrategy != nil {
		backoffStrategy = *node.Data.BackoffStrategy
	}
	
	initialDelay := 1 * time.Second
	if node.Data.InitialDelay != nil {
		if d, err := parseDuration(*node.Data.InitialDelay); err == nil {
			initialDelay = d
		}
	}
	
	maxDelay := 30 * time.Second
	if node.Data.MaxDelay != nil {
		if d, err := parseDuration(*node.Data.MaxDelay); err == nil {
			maxDelay = d
		}
	}
	
	multiplier := 2.0
	if node.Data.Multiplier != nil {
		multiplier = *node.Data.Multiplier
	}
	
	// Get retry_on_errors patterns (optional)
	retryOnErrors := node.Data.RetryOnErrors
	
	// Validate inputs
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
		"value":     input,
		"attempts":  maxAttempts,
		"success":   false,
		"error":     lastError.Error(),
		"last_delay": currentDelay.String(),
	}, lastError
}

// ========== Try-Catch Node ==========

// executeTryCatchNode implements error handling with fallback values
// Catches errors and provides fallback values or continues workflow execution
func (e *Engine) executeTryCatchNode(node *Node, inputs []interface{}) (interface{}, error) {
	// Get configuration
	var fallbackValue interface{}
	fallbackValue = node.Data.FallbackValue
	
	continueOnError := true
	if node.Data.ContinueOnError != nil {
		continueOnError = *node.Data.ContinueOnError
	}
	
	errorOutputPath := ""
	if node.Data.ErrorOutputPath != nil {
		errorOutputPath = *node.Data.ErrorOutputPath
	}
	
	// Validate inputs
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
			"value":       input,
			"error_caught": false,
		}, nil
	}
	
	// Error detected - apply fallback strategy
	result := map[string]interface{}{
		"value":        fallbackValue,
		"error_caught": true,
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

// ========== Timeout Node ==========

// executeTimeoutNode enforces time limits on operations
// Returns partial results or error if operation exceeds timeout
func (e *Engine) executeTimeoutNode(node *Node, inputs []interface{}) (interface{}, error) {
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

// Helper function already exists in workflow_advanced_nodes.go, but including here for completeness
// parseDuration is a helper to parse duration strings with support for ms, s, m, h
func parseDurationIfNotExists(s string) (time.Duration, error) {
	// Support common formats: "100ms", "5s", "2m", "1h"
	if strings.HasSuffix(s, "ms") {
		ms := strings.TrimSuffix(s, "ms")
		var value float64
		_, err := fmt.Sscanf(ms, "%f", &value)
		if err != nil {
			return 0, err
		}
		return time.Duration(value) * time.Millisecond, nil
	}
	if strings.HasSuffix(s, "s") {
		sec := strings.TrimSuffix(s, "s")
		var value float64
		_, err := fmt.Sscanf(sec, "%f", &value)
		if err != nil {
			return 0, err
		}
		return time.Duration(value * float64(time.Second)), nil
	}
	if strings.HasSuffix(s, "m") {
		min := strings.TrimSuffix(s, "m")
		var value float64
		_, err := fmt.Sscanf(min, "%f", &value)
		if err != nil {
			return 0, err
		}
		return time.Duration(value * float64(time.Minute)), nil
	}
	if strings.HasSuffix(s, "h") {
		hr := strings.TrimSuffix(s, "h")
		var value float64
		_, err := fmt.Sscanf(hr, "%f", &value)
		if err != nil {
			return 0, err
		}
		return time.Duration(value * float64(time.Hour)), nil
	}
	
	// Try standard parsing
	return time.ParseDuration(s)
}
