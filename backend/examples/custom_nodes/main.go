package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/yesoreyeram/thaiyyal/backend"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ============================================================================
// Example 1: Simple Custom Executor - String Reversal
// ============================================================================

// ReverseStringExecutor reverses the input string
type ReverseStringExecutor struct{}

func (e *ReverseStringExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
	// IMPORTANT: Always increment the node execution counter for protection
	if err := ctx.IncrementNodeExecution(); err != nil {
		return nil, err
	}

	// Get inputs from connected nodes
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("reverse_string requires at least one input")
	}

	// Type assertion to ensure input is a string
	str, ok := inputs[0].(string)
	if !ok {
		return nil, fmt.Errorf("reverse_string requires string input, got %T", inputs[0])
	}

	// Reverse the string using runes (handles Unicode correctly)
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}

func (e *ReverseStringExecutor) NodeType() types.NodeType {
	return types.NodeType("reverse_string")
}

func (e *ReverseStringExecutor) Validate(node workflow.Node) error {
	// No specific validation needed - any configuration is valid
	return nil
}

// ============================================================================
// Example 2: Custom Executor with Configuration - JSON Parser
// ============================================================================

// JSONPathExecutor extracts values from JSON using a path
type JSONPathExecutor struct{}

func (e *JSONPathExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
	if err := ctx.IncrementNodeExecution(); err != nil {
		return nil, err
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("json_path requires input")
	}

	// Get the JSON path from node configuration
	path := ""
	if node.Data.Field != nil {
		path = *node.Data.Field
	}

	// Try to parse input as JSON if it's a string
	var data interface{}
	if strInput, ok := inputs[0].(string); ok {
		if err := json.Unmarshal([]byte(strInput), &data); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	} else {
		data = inputs[0]
	}

	// Simple path extraction (supports dot notation like "user.name")
	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		if part == "" {
			continue
		}

		switch v := current.(type) {
		case map[string]interface{}:
			if val, exists := v[part]; exists {
				current = val
			} else {
				return nil, fmt.Errorf("path not found: %s", path)
			}
		default:
			return nil, fmt.Errorf("cannot traverse non-object at path: %s", part)
		}
	}

	return current, nil
}

func (e *JSONPathExecutor) NodeType() types.NodeType {
	return types.NodeType("json_path")
}

func (e *JSONPathExecutor) Validate(node workflow.Node) error {
	if node.Data.Field == nil || *node.Data.Field == "" {
		return fmt.Errorf("json_path node requires 'field' configuration")
	}
	return nil
}

// ============================================================================
// Example 3: Custom Executor with HTTP - External API Call
// ============================================================================

// WeatherAPIExecutor calls a weather API (demonstrates HTTP in custom nodes)
type WeatherAPIExecutor struct{}

func (e *WeatherAPIExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
	if err := ctx.IncrementNodeExecution(); err != nil {
		return nil, err
	}

	// IMPORTANT: If your custom executor makes HTTP calls, increment the HTTP counter
	if err := ctx.IncrementHTTPCall(); err != nil {
		return nil, err
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("weather_api requires city input")
	}

	city, ok := inputs[0].(string)
	if !ok {
		return nil, fmt.Errorf("weather_api requires string city name")
	}

	// In a real implementation, you would make an HTTP call here
	// For this example, we'll return mock data
	mockWeather := map[string]interface{}{
		"city":        city,
		"temperature": 22.5,
		"condition":   "Sunny",
		"humidity":    65,
	}

	return mockWeather, nil
}

func (e *WeatherAPIExecutor) NodeType() types.NodeType {
	return types.NodeType("weather_api")
}

func (e *WeatherAPIExecutor) Validate(node workflow.Node) error {
	return nil
}

// ============================================================================
// Example 4: Custom Executor with Iteration
// ============================================================================

// BatchProcessExecutor processes items in batches (demonstrates iteration)
type BatchProcessExecutor struct{}

func (e *BatchProcessExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
	if err := ctx.IncrementNodeExecution(); err != nil {
		return nil, err
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("batch_process requires array input")
	}

	// Convert input to array
	items, ok := inputs[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("batch_process requires array input")
	}

	// Get batch size from configuration
	batchSize := 10
	if node.Data.Delta != nil {
		batchSize = int(*node.Data.Delta)
	}

	// Process items in batches
	var batches [][]interface{}
	for i := 0; i < len(items); i += batchSize {
		// Increment counter for each batch iteration
		if err := ctx.IncrementNodeExecution(); err != nil {
			return nil, err
		}

		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}

	return map[string]interface{}{
		"total_items": len(items),
		"batch_size":  batchSize,
		"num_batches": len(batches),
		"batches":     batches,
	}, nil
}

func (e *BatchProcessExecutor) NodeType() types.NodeType {
	return types.NodeType("batch_process")
}

func (e *BatchProcessExecutor) Validate(node workflow.Node) error {
	if node.Data.Delta != nil && *node.Data.Delta <= 0 {
		return fmt.Errorf("batch_size must be positive")
	}
	return nil
}

// ============================================================================
// Main Examples
// ============================================================================

func main() {
	fmt.Println("=== Thaiyyal Custom Node Examples ===")
	fmt.Println()

	// Example 1: Simple String Reversal
	runExample1()

	// Example 2: JSON Path Extraction
	runExample2()

	// Example 3: Multiple Custom Nodes Combined
	runExample3()

	// Example 4: Custom Nodes with Built-in Nodes
	runExample4()

	// Example 5: Security - Protection Limits Apply
	runExample5()
}

func runExample1() {
	fmt.Println("Example 1: Simple String Reversal")
	fmt.Println("----------------------------------")

	// Create registry and register custom executor
	registry := workflow.DefaultRegistry()
	registry.MustRegister(&ReverseStringExecutor{})

	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "Hello, World!"}},
			{"id": "2", "type": "reverse_string", "data": {}}
		],
		"edges": [
			{"source": "1", "target": "2"}
		]
	}`

	engine, err := workflow.NewEngineWithRegistry([]byte(payload), workflow.DefaultConfig(), registry)
	if err != nil {
		log.Fatal(err)
	}

	result, err := engine.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Input: Hello, World!\n")
	fmt.Printf("Output: %v\n", result.FinalOutput)
	fmt.Println()
}

func runExample2() {
	fmt.Println("Example 2: JSON Path Extraction")
	fmt.Println("--------------------------------")

	registry := workflow.DefaultRegistry()
	registry.MustRegister(&JSONPathExecutor{})

	// Create a workflow that generates JSON and extracts a field
	field := "user.name"
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "1", "data": {"text": "{\"user\": {\"name\": \"Alice\", \"age\": 30}}"}},
			{"id": "2", "type": "json_path", "data": {"field": "%s"}}
		],
		"edges": [
			{"source": "1", "target": "2"}
		]
	}`, field)

	engine, err := workflow.NewEngineWithRegistry([]byte(payload), workflow.DefaultConfig(), registry)
	if err != nil {
		log.Fatal(err)
	}

	result, err := engine.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("JSON Path: %s\n", field)
	fmt.Printf("Extracted Value: %v\n", result.FinalOutput)
	fmt.Println()
}

func runExample3() {
	fmt.Println("Example 3: Multiple Custom Nodes")
	fmt.Println("---------------------------------")

	registry := workflow.DefaultRegistry()
	registry.MustRegister(&ReverseStringExecutor{})
	registry.MustRegister(&JSONPathExecutor{})
	registry.MustRegister(&WeatherAPIExecutor{})

	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "Paris"}},
			{"id": "2", "type": "weather_api", "data": {}},
			{"id": "3", "type": "json_path", "data": {"field": "condition"}},
			{"id": "4", "type": "reverse_string", "data": {}}
		],
		"edges": [
			{"source": "1", "target": "2"},
			{"source": "2", "target": "3"},
			{"source": "3", "target": "4"}
		]
	}`

	engine, err := workflow.NewEngineWithRegistry([]byte(payload), workflow.DefaultConfig(), registry)
	if err != nil {
		log.Fatal(err)
	}

	result, err := engine.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("City: Paris\n")
	fmt.Printf("Weather Condition (reversed): %v\n", result.FinalOutput)
	fmt.Println()
}

func runExample4() {
	fmt.Println("Example 4: Custom + Built-in Nodes")
	fmt.Println("-----------------------------------")

	registry := workflow.DefaultRegistry()
	registry.MustRegister(&ReverseStringExecutor{})

	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "hello"}},
			{"id": "2", "data": {"text_op": "uppercase"}},
			{"id": "3", "type": "reverse_string", "data": {}},
			{"id": "4", "data": {"text_op": "lowercase"}}
		],
		"edges": [
			{"source": "1", "target": "2"},
			{"source": "2", "target": "3"},
			{"source": "3", "target": "4"}
		]
	}`

	engine, err := workflow.NewEngineWithRegistry([]byte(payload), workflow.DefaultConfig(), registry)
	if err != nil {
		log.Fatal(err)
	}

	result, err := engine.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Workflow: hello → UPPERCASE → REVERSE → lowercase\n")
	fmt.Printf("Result: %v\n", result.FinalOutput)
	fmt.Println()
}

func runExample5() {
	fmt.Println("Example 5: Protection Limits Apply to Custom Nodes")
	fmt.Println("---------------------------------------------------")

	registry := workflow.DefaultRegistry()
	registry.MustRegister(&BatchProcessExecutor{})

	// Create config with low limits
	config := workflow.DefaultConfig()
	config.MaxNodeExecutions = 5 // Very low limit

	batchSize := 2.0
	// Using transform node to create an array, then batch process it
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "1", "data": {"value": 1}},
			{"id": "2", "data": {"value": 2}},
			{"id": "3", "data": {"value": 3}},
			{"id": "4", "data": {"value": 4}},
			{"id": "5", "data": {"value": 5}},
			{"id": "6", "type": "transform", "data": {"transform_type": "to_array"}},
			{"id": "7", "type": "batch_process", "data": {"delta": %f}}
		],
		"edges": [
			{"source": "1", "target": "6"},
			{"source": "2", "target": "6"},
			{"source": "3", "target": "6"},
			{"source": "4", "target": "6"},
			{"source": "5", "target": "6"},
			{"source": "6", "target": "7"}
		]
	}`, batchSize)

	engine, err := workflow.NewEngineWithRegistry([]byte(payload), config, registry)
	if err != nil {
		log.Fatal(err)
	}

	result, err := engine.Execute()
	if err != nil {
		fmt.Printf("✓ Protection limit enforced: %v\n", err)
	} else {
		fmt.Printf("Processed items in batches: %+v\n", result.FinalOutput)
	}
	fmt.Println()
}
