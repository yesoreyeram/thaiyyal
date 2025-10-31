package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// VariableExecutor executes Variable nodes
type VariableExecutor struct{}

// Execute runs the Variable node
// Handles variable get/set operations for workflow state.
// Variables are scoped to a single workflow execution and shared across nodes.
func (e *VariableExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	if node.Data.VarName == nil {
		return nil, fmt.Errorf("variable node missing var_name")
	}
	if node.Data.VarOp == nil {
		return nil, fmt.Errorf("variable node missing var_op (get or set)")
	}

	varName := *node.Data.VarName
	varOp := *node.Data.VarOp

	switch varOp {
	case "set":
		// Store value in workflow state
		inputs := ctx.GetNodeInputs(node.ID)
		if len(inputs) == 0 {
			return nil, fmt.Errorf("variable set operation requires input value")
		}
		value := inputs[0]
		if err := ctx.SetVariable(varName, value); err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"var_name":  varName,
			"operation": "set",
			"value":     value,
		}, nil

	case "get":
		// Retrieve value from workflow state
		value, err := ctx.GetVariable(varName)
		if err != nil {
			return nil, fmt.Errorf("variable '%s' not found", varName)
		}
		return map[string]interface{}{
			"var_name":  varName,
			"operation": "get",
			"value":     value,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported variable operation: %s (use 'get' or 'set')", varOp)
	}
}

// NodeType returns the node type this executor handles
func (e *VariableExecutor) NodeType() types.NodeType {
	return types.NodeTypeVariable
}

// Validate checks if node configuration is valid
func (e *VariableExecutor) Validate(node types.Node) error {
	if node.Data.VarName == nil {
		return fmt.Errorf("variable node missing var_name")
	}
	if node.Data.VarOp == nil {
		return fmt.Errorf("variable node missing var_op")
	}
	return nil
}
