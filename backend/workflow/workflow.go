package workflow

import "fmt"

// ExecuteWorkflow is a convenience function that parses and executes a workflow
// from a JSON payload in one call
func ExecuteWorkflow(jsonData []byte) (*ExecutionResult, error) {
	// Parse the workflow
	parser := NewParser()
	workflow, err := parser.Parse(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}

	// Execute the workflow
	engine := NewEngine(workflow)
	result, err := engine.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to execute workflow: %w", err)
	}

	return result, nil
}
