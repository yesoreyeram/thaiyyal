package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SwitchExecutor executes Switch nodes
type SwitchExecutor struct{}

// Execute runs the Switch node
// Handles switch/case node execution
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

	// Check each case
	for _, switchCase := range data.Cases {
		matched := false

		// If switchCase.Value is set, do value matching
		if switchCase.Value != nil {
			matched = compareValues(inputValue, switchCase.Value)
		} else {
			// Otherwise, evaluate as a condition
			matched = evaluateCondition(switchCase.When, inputValue)
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
	if data.DefaultPath != nil {
		defaultPath = *data.DefaultPath
	}

	return map[string]interface{}{
		"value":       inputValue,
		"matched":     false,
		"output_path": defaultPath,
	}, nil
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
	// Switch node should have at least one case
	if len(data.Cases) == 0 {
		return fmt.Errorf("switch node requires at least one case")
	}
	return nil
}
