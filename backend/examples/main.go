package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/yesoreyeram/thaiyyal/backend/workflow"
)

func main() {
	// Example 1: Simple addition workflow
	fmt.Println("=== Example 1: Simple Addition ===")
	simplePayload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10, "label": "First Number"}},
			{"id": "2", "data": {"value": 5, "label": "Second Number"}},
			{"id": "3", "data": {"op": "add", "label": "Add Operation"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"}
		]
	}`
	executeWorkflow(simplePayload)

	// Example 2: Complete workflow with visualization
	fmt.Println("\n=== Example 2: Complete Workflow with Visualization ===")
	completePayload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10, "label": "Node 1"}},
			{"id": "2", "data": {"value": 5, "label": "Node 2"}},
			{"id": "3", "data": {"op": "multiply", "label": "Node 3 (op)"}},
			{"id": "4", "data": {"mode": "text", "label": "Node 4 (viz)"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-4", "source": "3", "target": "4"}
		]
	}`
	executeWorkflow(completePayload)

	// Example 3: Complex workflow with multiple operations
	fmt.Println("\n=== Example 3: Complex Workflow (10 + 5) * 2 - 3 ===")
	complexPayload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "data": {"op": "add"}},
			{"id": "4", "data": {"value": 2}},
			{"id": "5", "data": {"op": "multiply"}},
			{"id": "6", "data": {"value": 3}},
			{"id": "7", "data": {"op": "subtract"}},
			{"id": "8", "data": {"mode": "table"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-5", "source": "3", "target": "5"},
			{"id": "e4-5", "source": "4", "target": "5"},
			{"id": "e5-7", "source": "5", "target": "7"},
			{"id": "e6-7", "source": "6", "target": "7"},
			{"id": "e7-8", "source": "7", "target": "8"}
		]
	}`
	executeWorkflow(complexPayload)

	// Example 4: Division workflow
	fmt.Println("\n=== Example 4: Division Workflow ===")
	divisionPayload := `{
		"nodes": [
			{"id": "1", "data": {"value": 100}},
			{"id": "2", "data": {"value": 4}},
			{"id": "3", "data": {"op": "divide"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"}
		]
	}`
	executeWorkflow(divisionPayload)
}

func executeWorkflow(payloadJSON string) {
	// Create the workflow engine
	engine, err := workflow.NewEngine([]byte(payloadJSON))
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}

	// Execute the workflow
	result, err := engine.Execute()
	if err != nil {
		log.Printf("Workflow execution failed: %v", err)
		return
	}

	// Print the results
	fmt.Println("Node Results:")
	for nodeID, value := range result.NodeResults {
		fmt.Printf("  Node %s: %v\n", nodeID, value)
	}

	fmt.Println("Final Output:", result.FinalOutput)

	if len(result.Errors) > 0 {
		fmt.Println("Errors:", result.Errors)
	}

	// Pretty print the full result as JSON
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("\nFull Result JSON:")
	fmt.Println(string(resultJSON))
}
