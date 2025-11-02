package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/expression"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// PartitionExecutor splits an array into two groups based on a condition
type PartitionExecutor struct{}

// Execute splits the input array into passed and failed groups
func (e *PartitionExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("partition node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("partition node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	// Get condition
	condition := ""
	if node.Data.Condition != nil {
		condition = *node.Data.Condition
	}
	if condition == "" {
		return nil, fmt.Errorf("partition node missing required 'condition' string")
	}

	// Partition the array
	var passed []interface{}
	var failed []interface{}
	passedCount := 0
	failedCount := 0

	for i, item := range arr {
		// Create context with item and index variables
		itemCtx := &expression.Context{
			Variables:   make(map[string]interface{}),
			ContextVars: ctx.GetContextVariables(),
			NodeResults: ctx.GetAllNodeResults(),
		}
		// Copy existing variables
		for k, v := range ctx.GetVariables() {
			itemCtx.Variables[k] = v
		}
		itemCtx.Variables["item"] = item
		itemCtx.Variables["index"] = i
		itemCtx.Variables["items"] = arr

		// Evaluate condition
		result, err := expression.Evaluate(condition, item, itemCtx)
		if err != nil {
			// On error, add to failed
			failed = append(failed, item)
			failedCount++
			continue
		}

		// Check if condition is true
		if result {
			passed = append(passed, item)
			passedCount++
		} else {
			failed = append(failed, item)
			failedCount++
		}
	}

	return map[string]interface{}{
		"passed":       passed,
		"failed":       failed,
		"passed_count": passedCount,
		"failed_count": failedCount,
		"input_count":  len(arr),
		"condition":    condition,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *PartitionExecutor) NodeType() types.NodeType {
	return types.NodeTypePartition
}

// Validate checks if the node configuration is valid
func (e *PartitionExecutor) Validate(node types.Node) error {
	if node.Data.Condition == nil || *node.Data.Condition == "" {
		return fmt.Errorf("partition node requires non-empty 'condition' field")
	}
	return nil
}
