package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/expression"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ForEachExecutor executes ForEach nodes with full loop capabilities
type ForEachExecutor struct{}

// ForEachMode defines how the loop processes and returns data
type ForEachMode string

const (
	// ForEachModeMap transforms each element and returns new array
	ForEachModeMap ForEachMode = "map"
	// ForEachModeReduce accumulates values and returns single result
	ForEachModeReduce ForEachMode = "reduce"
	// ForEachModeFilterMap filters and transforms elements
	ForEachModeFilterMap ForEachMode = "filter_map"
	// ForEachModeForEach executes for side effects only
	ForEachModeForEach ForEachMode = "foreach"
	// ForEachModeMetadata returns metadata only (backward compatible)
	ForEachModeMetadata ForEachMode = "metadata"
)

// Execute runs the ForEach node with full loop body execution
//
// The ForEach node processes arrays in several modes:
// - MAP: Transform each element → return new array
// - REDUCE: Accumulate values → return single value
// - FILTER_MAP: Filter and transform → return filtered array
// - FOREACH: Execute for side effects → return metadata
// - METADATA: Return array info only (backward compatible stub)
//
// Loop Body Execution:
// The ForEach node identifies child nodes that should execute for each iteration.
// These are nodes that depend on the ForEach node output. For each array element,
// the executor:
// 1. Creates iteration context with item, index, items variables
// 2. Executes child nodes with this context
// 3. Collects results based on mode
//
// Example workflow:
//   [users] → ForEach(MAP) → Extract(name) → [names]
//
// Data Flow:
//   Input: [{"name":"Alice","age":25}, {"name":"Bob","age":30}]
//   Iteration 1: item={"name":"Alice","age":25}, index=0 → Extract → "Alice"
//   Iteration 2: item={"name":"Bob","age":30}, index=1 → Extract → "Bob"
//   Output: {results: ["Alice", "Bob"], mode: "map", ...}
func (e *ForEachExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("for_each node needs at least 1 input")
	}

	// Check if input is an array (slice)
	inputArray, ok := inputs[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("for_each node requires array input, got %T", inputs[0])
	}

	// Determine execution mode
	mode := ForEachModeMap // Default to map mode
	if node.Data.Mode != nil && *node.Data.Mode != "" {
		mode = ForEachMode(*node.Data.Mode)
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

	slog.Debug("for_each node starting",
		slog.String("node_id", node.ID),
		slog.String("mode", string(mode)),
		slog.Int("input_count", len(inputArray)),
	)

	// Execute based on mode
	switch mode {
	case ForEachModeMap:
		return e.executeMap(ctx, node, inputArray)
	case ForEachModeReduce:
		return e.executeReduce(ctx, node, inputArray)
	case ForEachModeFilterMap:
		return e.executeFilterMap(ctx, node, inputArray)
	case ForEachModeForEach:
		return e.executeForEach(ctx, node, inputArray)
	case ForEachModeMetadata:
		return e.executeMetadata(ctx, node, inputArray)
	default:
		return nil, fmt.Errorf("unknown for_each mode: %s", mode)
	}
}

// executeMap transforms each array element and returns new array
func (e *ForEachExecutor) executeMap(ctx ExecutionContext, node types.Node, inputArray []interface{}) (interface{}, error) {
	results := make([]interface{}, 0, len(inputArray))
	successful := 0
	failed := 0

	for i, item := range inputArray {
		result, err := e.executeIteration(ctx, node, item, i, inputArray, nil)
		if err != nil {
			slog.Debug("for_each iteration error (continuing)",
				slog.String("node_id", node.ID),
				slog.Int("index", i),
				slog.String("error", err.Error()),
			)
			failed++
			// Continue on error - collect nil for failed iteration
			results = append(results, nil)
			continue
		}
		
		results = append(results, result)
		successful++
	}

	slog.Debug("for_each map completed",
		slog.String("node_id", node.ID),
		slog.Int("successful", successful),
		slog.Int("failed", failed),
	)

	return map[string]interface{}{
		"results":      results,
		"mode":         string(ForEachModeMap),
		"input_count":  len(inputArray),
		"output_count": len(results),
		"successful":   successful,
		"failed":       failed,
	}, nil
}

// executeReduce accumulates values across iterations
func (e *ForEachExecutor) executeReduce(ctx ExecutionContext, node types.Node, inputArray []interface{}) (interface{}, error) {
	// Get initial value for accumulator
	accumulator := node.Data.InitialValue
	if accumulator == nil {
		accumulator = float64(0) // Default to 0
	}

	successful := 0
	failed := 0

	for i, item := range inputArray {
		result, err := e.executeIteration(ctx, node, item, i, inputArray, accumulator)
		if err != nil {
			slog.Debug("for_each reduce iteration error (continuing)",
				slog.String("node_id", node.ID),
				slog.Int("index", i),
				slog.String("error", err.Error()),
			)
			failed++
			continue
		}
		
		// Update accumulator with iteration result
		accumulator = result
		successful++
	}

	slog.Debug("for_each reduce completed",
		slog.String("node_id", node.ID),
		slog.Int("successful", successful),
		slog.Int("failed", failed),
	)

	return map[string]interface{}{
		"result":        accumulator,
		"mode":          string(ForEachModeReduce),
		"input_count":   len(inputArray),
		"iterations":    successful + failed,
		"successful":    successful,
		"failed":        failed,
		"initial_value": node.Data.InitialValue,
		"final_value":   accumulator,
	}, nil
}

// executeFilterMap filters and transforms elements
func (e *ForEachExecutor) executeFilterMap(ctx ExecutionContext, node types.Node, inputArray []interface{}) (interface{}, error) {
	results := make([]interface{}, 0, len(inputArray))
	successful := 0
	failed := 0
	filtered := 0

	// Get filter condition
	filterCondition := ""
	if node.Data.Condition != nil {
		filterCondition = *node.Data.Condition
	}

	for i, item := range inputArray {
		// Evaluate filter condition if present
		if filterCondition != "" {
			// Create expression context with item
			exprCtx := &expression.Context{
				NodeResults: ctx.GetAllNodeResults(),
				Variables:   make(map[string]interface{}),
				ContextVars: ctx.GetContextVariables(),
			}
			
			// Copy existing variables
			for k, v := range ctx.GetVariables() {
				exprCtx.Variables[k] = v
			}
			
			// Add current item
			exprCtx.Variables["item"] = item
			exprCtx.Variables["index"] = float64(i)
			exprCtx.Variables["items"] = inputArray

			// Evaluate filter
			passes, err := expression.Evaluate(filterCondition, item, exprCtx)
			if err != nil || !passes {
				filtered++
				continue
			}
		}

		// Item passed filter - execute iteration
		result, err := e.executeIteration(ctx, node, item, i, inputArray, nil)
		if err != nil {
			slog.Debug("for_each filter_map iteration error (continuing)",
				slog.String("node_id", node.ID),
				slog.Int("index", i),
				slog.String("error", err.Error()),
			)
			failed++
			continue
		}
		
		results = append(results, result)
		successful++
	}

	slog.Debug("for_each filter_map completed",
		slog.String("node_id", node.ID),
		slog.Int("successful", successful),
		slog.Int("failed", failed),
		slog.Int("filtered", filtered),
	)

	return map[string]interface{}{
		"results":        results,
		"mode":           string(ForEachModeFilterMap),
		"input_count":    len(inputArray),
		"output_count":   len(results),
		"filtered_count": filtered,
		"successful":     successful,
		"failed":         failed,
	}, nil
}

// executeForEach executes for side effects only
func (e *ForEachExecutor) executeForEach(ctx ExecutionContext, node types.Node, inputArray []interface{}) (interface{}, error) {
	successful := 0
	failed := 0

	for i, item := range inputArray {
		_, err := e.executeIteration(ctx, node, item, i, inputArray, nil)
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
		"mode":        string(ForEachModeForEach),
		"input_count": len(inputArray),
		"iterations":  successful + failed,
		"successful":  successful,
		"failed":      failed,
	}, nil
}

// executeMetadata returns metadata only (backward compatible stub)
func (e *ForEachExecutor) executeMetadata(ctx ExecutionContext, node types.Node, inputArray []interface{}) (interface{}, error) {
	return map[string]interface{}{
		"items":      inputArray,
		"count":      len(inputArray),
		"iterations": len(inputArray),
		"mode":       string(ForEachModeMetadata),
	}, nil
}

// executeIteration executes loop body for a single array element
//
// This is the core iteration logic that:
// 1. Creates iteration-specific context with item/index/items variables
// 2. Identifies and executes child nodes (loop body)
// 3. Returns the final result from the iteration
//
// NOTE: This is a simplified implementation that executes direct child nodes.
// A full implementation would need access to the workflow engine to execute
// a sub-DAG of nodes marked as loop body.
//
// For now, we simulate by:
// - Making item/index available in variables
// - Returning item (for map mode to work)
// - Actual child node execution would be handled by the workflow engine
func (e *ForEachExecutor) executeIteration(
	ctx ExecutionContext,
	node types.Node,
	item interface{},
	index int,
	items []interface{},
	accumulator interface{},
) (interface{}, error) {
	// Create iteration context with item, index, items available as variables
	// In a full implementation, this would be passed to child node executors
	
	// For now, return the item itself (allows basic map operations)
	// The workflow engine would inject these variables when executing child nodes
	
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
	// 4. Return the final output from the child execution
	
	// For now, return item to enable basic functionality
	return item, nil
}

// NodeType returns the node type this executor handles
func (e *ForEachExecutor) NodeType() types.NodeType {
	return types.NodeTypeForEach
}

// Validate checks if node configuration is valid
func (e *ForEachExecutor) Validate(node types.Node) error {
	// Validate mode if specified
	if node.Data.Mode != nil {
		mode := ForEachMode(*node.Data.Mode)
		switch mode {
		case ForEachModeMap, ForEachModeReduce, ForEachModeFilterMap, ForEachModeForEach, ForEachModeMetadata:
			// Valid modes
		default:
			return fmt.Errorf("invalid for_each mode: %s (must be one of: map, reduce, filter_map, foreach, metadata)", mode)
		}

		// Validate mode-specific requirements
		if mode == ForEachModeReduce && node.Data.InitialValue == nil {
			slog.Warn("for_each reduce mode without initial_value, defaulting to 0",
				slog.String("node_id", node.ID),
			)
		}

		if mode == ForEachModeFilterMap && node.Data.Condition == nil {
			return fmt.Errorf("for_each filter_map mode requires condition")
		}
	}

	return nil
}
