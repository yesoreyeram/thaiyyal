package workflow

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// executeSwitchNode handles switch/case node execution
func (e *Engine) executeSwitchNode(node Node) (interface{}, error) {
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("switch node requires at least one input")
	}

	// Get the input value to switch on
	inputValue := inputs[0]

	// Check each case
	for _, switchCase := range node.Data.Cases {
		matched := false

		// If switchCase.Value is set, do value matching
		if switchCase.Value != nil {
			matched = compareValues(inputValue, switchCase.Value)
		} else {
			// Otherwise, evaluate as a condition
			matched = e.evaluateCondition(switchCase.When, inputValue)
		}

		if matched {
			outputPath := "matched"
			if switchCase.OutputPath != nil {
				outputPath = *switchCase.OutputPath
			}
			return map[string]interface{}{
				"value":       inputValue,
				"matched":     true,
				"case":        switchCase.When,
				"output_path": outputPath,
			}, nil
		}
	}

	// No case matched, use default
	defaultPath := "default"
	if node.Data.DefaultPath != nil {
		defaultPath = *node.Data.DefaultPath
	}

	return map[string]interface{}{
		"value":       inputValue,
		"matched":     false,
		"output_path": defaultPath,
	}, nil
}

// executeParallelNode handles parallel execution of multiple branches
func (e *Engine) executeParallelNode(node Node) (interface{}, error) {
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("parallel node requires at least one input")
	}

	maxConcurrency := 10 // default
	if node.Data.MaxConcurrency != nil {
		maxConcurrency = *node.Data.MaxConcurrency
	}

	// Create semaphore for concurrency control
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]interface{}, len(inputs))
	errors := make([]error, len(inputs))

	for i, input := range inputs {
		wg.Add(1)
		go func(index int, value interface{}) {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release

			// Process the input (in real implementation, this would execute a sub-workflow)
			mu.Lock()
			results[index] = value
			mu.Unlock()
		}(i, input)
	}

	wg.Wait()

	// Check for errors
	var errorMsgs []string
	for _, err := range errors {
		if err != nil {
			errorMsgs = append(errorMsgs, err.Error())
		}
	}

	if len(errorMsgs) > 0 {
		return nil, fmt.Errorf("parallel execution errors: %s", strings.Join(errorMsgs, "; "))
	}

	return map[string]interface{}{
		"results":     results,
		"count":       len(results),
		"concurrency": maxConcurrency,
	}, nil
}

// executeJoinNode handles joining/merging multiple inputs
func (e *Engine) executeJoinNode(node Node) (interface{}, error) {
	inputs := e.getNodeInputs(node.ID)

	strategy := "all" // default strategy
	if node.Data.JoinStrategy != nil {
		strategy = *node.Data.JoinStrategy
	}

	switch strategy {
	case "all":
		// Wait for all inputs and combine them
		if len(inputs) == 0 {
			return nil, fmt.Errorf("join node with 'all' strategy requires at least one input")
		}
		return map[string]interface{}{
			"strategy": "all",
			"values":   inputs,
			"count":    len(inputs),
		}, nil

	case "any":
		// Return as soon as any input is available
		if len(inputs) > 0 {
			return map[string]interface{}{
				"strategy": "any",
				"value":    inputs[0],
				"count":    len(inputs),
			}, nil
		}
		return nil, fmt.Errorf("join node with 'any' strategy has no inputs")

	case "first":
		// Return only the first input
		if len(inputs) > 0 {
			return map[string]interface{}{
				"strategy": "first",
				"value":    inputs[0],
			}, nil
		}
		return nil, fmt.Errorf("join node with 'first' strategy has no inputs")

	default:
		return nil, fmt.Errorf("unsupported join strategy: %s (use all, any, or first)", strategy)
	}
}

// executeSplitNode handles splitting single input to multiple paths
func (e *Engine) executeSplitNode(node Node) (interface{}, error) {
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("split node requires at least one input")
	}

	inputValue := inputs[0]
	paths := node.Data.Paths
	if len(paths) == 0 {
		// Default to splitting to 2 paths
		paths = []string{"path1", "path2"}
	}

	// Create a copy of the input for each path
	outputs := make(map[string]interface{})
	for _, path := range paths {
		outputs[path] = inputValue
	}

	return map[string]interface{}{
		"value":  inputValue,
		"paths":  paths,
		"outputs": outputs,
	}, nil
}

// executeDelayNode handles execution delay
func (e *Engine) executeDelayNode(node Node) (interface{}, error) {
	inputs := e.getNodeInputs(node.ID)
	var inputValue interface{}
	if len(inputs) > 0 {
		inputValue = inputs[0]
	}

	if node.Data.Duration == nil {
		return nil, fmt.Errorf("delay node requires duration field")
	}

	duration, err := parseDuration(*node.Data.Duration)
	if err != nil {
		return nil, fmt.Errorf("invalid duration format: %w", err)
	}

	// Perform the delay
	time.Sleep(duration)

	return map[string]interface{}{
		"value":    inputValue,
		"duration": *node.Data.Duration,
		"delayed":  true,
	}, nil
}

// executeCacheNode handles cache get/set operations
func (e *Engine) executeCacheNode(node Node) (interface{}, error) {
	if node.Data.CacheOp == nil {
		return nil, fmt.Errorf("cache node requires cache_op field")
	}
	if node.Data.CacheKey == nil {
		return nil, fmt.Errorf("cache node requires cache_key field")
	}

	cacheOp := *node.Data.CacheOp
	cacheKey := *node.Data.CacheKey

	switch cacheOp {
	case "set":
		inputs := e.getNodeInputs(node.ID)
		if len(inputs) == 0 {
			return nil, fmt.Errorf("cache set operation requires an input value")
		}

		value := inputs[0]

		// Parse TTL
		var expiration time.Time
		if node.Data.TTL != nil {
			ttlDuration, err := parseDuration(*node.Data.TTL)
			if err != nil {
				return nil, fmt.Errorf("invalid TTL format: %w", err)
			}
			expiration = time.Now().Add(ttlDuration)
		} else {
			// Default TTL: 5 minutes
			expiration = time.Now().Add(5 * time.Minute)
		}

		e.cacheMutex.Lock()
		e.cache[cacheKey] = &CacheEntry{
			Value:      value,
			Expiration: expiration,
		}
		e.cacheMutex.Unlock()

		return map[string]interface{}{
			"operation": "set",
			"key":       cacheKey,
			"value":     value,
			"ttl":       node.Data.TTL,
		}, nil

	case "get":
		e.cacheMutex.RLock()
		entry, exists := e.cache[cacheKey]
		e.cacheMutex.RUnlock()

		if !exists {
			return map[string]interface{}{
				"operation": "get",
				"key":       cacheKey,
				"found":     false,
				"value":     nil,
			}, nil
		}

		// Check if expired
		if time.Now().After(entry.Expiration) {
			// Remove expired entry
			e.cacheMutex.Lock()
			delete(e.cache, cacheKey)
			e.cacheMutex.Unlock()

			return map[string]interface{}{
				"operation": "get",
				"key":       cacheKey,
				"found":     false,
				"expired":   true,
				"value":     nil,
			}, nil
		}

		return map[string]interface{}{
			"operation": "get",
			"key":       cacheKey,
			"found":     true,
			"value":     entry.Value,
		}, nil

	case "delete":
		e.cacheMutex.Lock()
		_, existed := e.cache[cacheKey]
		delete(e.cache, cacheKey)
		e.cacheMutex.Unlock()

		return map[string]interface{}{
			"operation": "delete",
			"key":       cacheKey,
			"deleted":   existed,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported cache operation: %s (use get, set, or delete)", cacheOp)
	}
}

// Helper function to parse duration strings
func parseDuration(durationStr string) (time.Duration, error) {
	// Support formats like "5s", "10m", "1h", "100ms"
	if duration, err := time.ParseDuration(durationStr); err == nil {
		return duration, nil
	}

	// Also support integer milliseconds
	if ms, err := strconv.Atoi(durationStr); err == nil {
		return time.Duration(ms) * time.Millisecond, nil
	}

	return 0, fmt.Errorf("invalid duration format: %s (use formats like '5s', '10m', '1h')", durationStr)
}

// Helper function to compare values for switch cases
func compareValues(a, b interface{}) bool {
	// Simple equality check
	switch aVal := a.(type) {
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal == bVal
		}
	case string:
		if bVal, ok := b.(string); ok {
			return aVal == bVal
		}
	case bool:
		if bVal, ok := b.(bool); ok {
			return aVal == bVal
		}
	}
	return false
}
