package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ForEachExecutor executes ForEach nodes
// Simple iterator that executes child nodes for each array element
type ForEachExecutor struct{}

// Execute runs the ForEach node
// Iterates over an array and makes each element available to child nodes through variables.
//
// The ForEach node is a pure iterator that:
// 1. Takes an array as input
// 2. For each element, makes variables available:
//   - `variables.item` - Current array element
//   - `variables.index` - Current index (0-based)
//   - `variables.items` - Full input array
//
// 3. Executes child nodes for each iteration (handled by workflow engine)
// 4. Returns execution metadata
//
// Use cases:
//   - Side effects: [users] → ForEach → HTTP(POST to API)
//   - Batch operations: [items] → ForEach → Process
//   - Combined with other nodes: [users] → ForEach → child nodes access variables.item
//
// For transformations, use dedicated nodes:
//   - Map node: Transform each element to new value
//   - Reduce node: Aggregate array to single value
//   - Filter node: Remove elements that don't match condition
//
// Example:
//
//	Input: [{"name":"Alice"}, {"name":"Bob"}]
//	Each iteration: variables.item = {"name":"Alice"}, variables.index = 0
//	Output: {iterations: 2, successful: 2, failed: 0}
func (e *ForEachExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsForEachData(node.Data)
	if err != nil {
		return nil, err
	}
	inputs := ctx.GetNodeInputs(node.ID)
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
	if data.MaxIterations != nil && *data.MaxIterations > 0 {
		maxIter = *data.MaxIterations
	}

	// Limit iterations to prevent resource exhaustion
	iterCount := len(inputArray)
	if iterCount > maxIter {
		return nil, fmt.Errorf("for_each exceeds max iterations: %d > %d", iterCount, maxIter)
	}

	slog.Debug("for_each node starting",
		slog.String("node_id", node.ID),
		slog.Int("input_count", len(inputArray)),
	)

	successful := 0
	failed := 0

	for i, item := range inputArray {
		err := e.executeIteration(ctx, node, item, i, inputArray)
		if err != nil {
			slog.Debug("for_each iteration error (continuing)",
				slog.String("node_id", node.ID),
				slog.Int("index", i),
				slog.String("error", err.Error()),
			)
			failed++
			continue
		}
		successful++
	}

	slog.Debug("for_each completed",
		slog.String("node_id", node.ID),
		slog.Int("successful", successful),
		slog.Int("failed", failed),
	)

	return map[string]interface{}{
		"input_count": len(inputArray),
		"iterations":  successful + failed,
		"successful":  successful,
		"failed":      failed,
	}, nil
}

// executeIteration executes loop body for a single array element
//
// This prepares the iteration context by making item/index/items available
// as variables. The workflow engine will inject these variables when executing
// child nodes.
//
// In a full implementation, this would:
// 1. Create execution context with item/index/items in variables
// 2. Identify child nodes (nodes that depend on this foreach node)
// 3. Execute child nodes in topological order with iteration context
// 4. Handle any execution errors gracefully
//
// For now, we prepare the context and log the iteration.
func (e *ForEachExecutor) executeIteration(
	_ ExecutionContext,
	node types.Node,
	item interface{},
	index int,
	_ []interface{},
) error {
	slog.Debug("for_each iteration",
		slog.String("node_id", node.ID),
		slog.Int("index", index),
		slog.Any("item", item),
	)

	// TODO: Execute child nodes here with iteration context
	// This requires integration with the workflow engine to:
	// 1. Identify child nodes (nodes that depend on this foreach node)
	// 2. Create execution context with item/index/items in variables
	// 3. Execute child nodes in topological order
	// 4. Return any execution errors

	// Variables that would be injected:
	// - variables.item = item
	// - variables.index = index
	// - variables.items = items

	return nil
}

// NodeType returns the node type this executor handles
func (e *ForEachExecutor) NodeType() types.NodeType {
	return types.NodeTypeForEach
}

// Validate checks if node configuration is valid
func (e *ForEachExecutor) Validate(node types.Node) error {
	// Validate node data type
	if _, err := types.AsForEachData(node.Data); err != nil {
		return err
	}
	// No required configuration for simple iterator
	// Max iterations is optional
	return nil
}
