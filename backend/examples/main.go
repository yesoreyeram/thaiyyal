package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/yesoreyeram/thaiyyal/backend/workflow"
)

func main() {
	fmt.Println("=== Thaiyyal Workflow Engine Example ===")

	// Example 1: Simple Addition
	fmt.Println("Example 1: Simple Addition (10 + 5)")
	example1 := []byte(`{
		"nodes": [
			{"id": "1", "data": {"value": 10, "label": "Number 10"}},
			{"id": "2", "data": {"value": 5, "label": "Number 5"}},
			{"id": "3", "data": {"op": "add", "label": "Add"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"}
		]
	}`)
	runExample(example1)

	// Example 2: Complete Workflow with Visualization
	fmt.Println("\nExample 2: Complete Workflow (10 + 5 = 15 with text visualization)")
	example2 := []byte(`{
		"nodes": [
			{"id": "1", "data": {"value": 10, "label": "Node 1"}},
			{"id": "2", "data": {"value": 5, "label": "Node 2"}},
			{"id": "3", "data": {"op": "add", "label": "Node 3 (op)"}},
			{"id": "4", "data": {"mode": "text", "label": "Node 4 (viz)"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-4", "source": "3", "target": "4"}
		]
	}`)
	runExample(example2)

	// Example 3: Multiplication and Division
	fmt.Println("\nExample 3: Complex Operations ((10 * 5) / 2 = 25)")
	example3 := []byte(`{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "data": {"op": "multiply"}},
			{"id": "4", "data": {"value": 2}},
			{"id": "5", "data": {"op": "divide"}},
			{"id": "6", "data": {"mode": "table"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-5", "source": "3", "target": "5"},
			{"id": "e4-5", "source": "4", "target": "5"},
			{"id": "e5-6", "source": "5", "target": "6"}
		]
	}`)
	runExample(example3)

	// Example 4: Subtraction
	fmt.Println("\nExample 4: Subtraction (100 - 30 - 20 = 50)")
	example4 := []byte(`{
		"nodes": [
			{"id": "1", "data": {"value": 100}},
			{"id": "2", "data": {"value": 30}},
			{"id": "3", "data": {"value": 20}},
			{"id": "4", "data": {"op": "subtract"}}
		],
		"edges": [
			{"id": "e1-4", "source": "1", "target": "4"},
			{"id": "e2-4", "source": "2", "target": "4"},
			{"id": "e3-4", "source": "3", "target": "4"}
		]
	}`)
	runExample(example4)
}

func runExample(jsonData []byte) {
	// Execute the workflow
	result, err := workflow.ExecuteWorkflow(jsonData)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	// Print results
	fmt.Println("Execution Results:")
	for nodeID, nodeResult := range result.Results {
		if nodeResult.Error != nil {
			fmt.Printf("  Node %s: ERROR - %v\n", nodeID, nodeResult.Error)
		} else {
			fmt.Printf("  Node %s: %v\n", nodeID, nodeResult.Value)
		}
	}

	// Print final output if available
	if result.Output != nil {
		fmt.Println("\nFinal Output:")
		if outputMap, ok := result.Output.(map[string]interface{}); ok {
			// Pretty print map
			prettyJSON, _ := json.MarshalIndent(outputMap, "  ", "  ")
			fmt.Printf("  %s\n", prettyJSON)
		} else {
			fmt.Printf("  %v\n", result.Output)
		}
	}
}
