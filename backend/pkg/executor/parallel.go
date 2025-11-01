package executor

import (
	"fmt"
	"strings"
	"sync"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ParallelExecutor executes Parallel nodes
type ParallelExecutor struct{}

// Execute runs the Parallel node
// Handles parallel execution of multiple branches
func (e *ParallelExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
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

// NodeType returns the node type this executor handles
func (e *ParallelExecutor) NodeType() types.NodeType {
	return types.NodeTypeParallel
}

// Validate checks if node configuration is valid
func (e *ParallelExecutor) Validate(node types.Node) error {
	// No required fields for parallel
	return nil
}
