package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/expression"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SwitchExecutor executes Switch nodes
type SwitchExecutor struct{}

// Execute runs the Switch node
// Evaluates cases in order using expression engine
// Returns the first matching case's output path, or default if no match
func (e *SwitchExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	data, err := types.AsSwitchData(node.Data)
	if err != nil {
		return nil, err
	}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("switch node requires at least one input")
	}

	// Get the input value to switch on
	inputValue := inputs[0]

	// Build expression context with access to node results and variables
	exprCtx := &expression.Context{
		NodeResults: ctx.GetAllNodeResults(),
		Variables:   ctx.GetVariables(),
		ContextVars: ctx.GetContextVariables(),
	}

	// Check each case in order (last case is default)
	for i, switchCase := range data.Cases {
		// If this is the default case (last one), it always matches
		if switchCase.IsDefault {
			outputPath := "default"
			if switchCase.OutputPath != nil && *switchCase.OutputPath != "" {
				outputPath = *switchCase.OutputPath
			}
			return map[string]interface{}{
				"value":       inputValue,
				"matched":     false, // No explicit match, using default
				"case":        "default",
				"output_path": outputPath,
				"case_index":  i,
			}, nil
		}

		// Evaluate the expression
		matched, err := expression.Evaluate(switchCase.When, inputValue, exprCtx)
		if err != nil {
			// If expression evaluation fails, skip this case and continue
			// This allows graceful degradation for malformed expressions
			continue
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
				"case_index":  i,
			}, nil
		}
	}

	// This should never happen if validation is correct (default case exists)
	// But handle gracefully just in case
	return map[string]interface{}{
		"value":       inputValue,
		"matched":     false,
		"output_path": "default",
		"case":        "fallback",
	}, fmt.Errorf("switch node reached end without matching any case or default (validation error)")
}

// NodeType returns the node type this executor handles
func (e *SwitchExecutor) NodeType() types.NodeType {
	return types.NodeTypeSwitch
}

// Validate checks if node configuration is valid
func (e *SwitchExecutor) Validate(node types.Node) error {
	data, err := types.AsSwitchData(node.Data)
	if err != nil {
		return err
	}
	// SwitchData.Validate() already checks:
	// - At least one case exists
	// - Exactly one default case exists
	// - Default case is last
	// - Non-default cases have output_path
	return data.Validate()
}
