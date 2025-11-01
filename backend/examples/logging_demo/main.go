package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
)

func main() {
	fmt.Println("=== Structured Logging Example ===")
	fmt.Println()

	// Create a simple workflow with multiple nodes
	workflow := map[string]interface{}{
		"workflow_id": "demo-workflow",
		"nodes": []map[string]interface{}{
			{
				"id":   "1",
				"type": "number",
				"data": map[string]interface{}{"value": 10},
			},
			{
				"id":   "2",
				"type": "number",
				"data": map[string]interface{}{"value": 20},
			},
			{
				"id":   "3",
				"type": "operation",
				"data": map[string]interface{}{"op": "add"},
			},
			{
				"id":   "4",
				"type": "visualization",
				"data": map[string]interface{}{"mode": "text"},
			},
		},
		"edges": []map[string]interface{}{
			{"source": "1", "target": "3"},
			{"source": "2", "target": "3"},
			{"source": "3", "target": "4"},
		},
	}

	payloadJSON, err := json.Marshal(workflow)
	if err != nil {
		log.Fatalf("Failed to marshal workflow: %v", err)
	}

	// Create engine - it will automatically initialize structured logger
	fmt.Println("Creating workflow engine with structured logging...")
	eng, err := engine.New(payloadJSON)
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}

	// Execute the workflow - observe structured logs
	fmt.Println("\nExecuting workflow (logs will be in JSON format):")
	fmt.Println("============================================================")
	
	result, err := eng.Execute()
	if err != nil {
		log.Fatalf("Failed to execute workflow: %v", err)
	}

	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("Workflow execution completed successfully!")
	fmt.Printf("Final result: %v\n", result.FinalOutput)
}
