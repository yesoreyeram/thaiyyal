package workflow

import (
	"testing"
)

// frontendDefaultPayload is the exact payload format from src/app/page.tsx (initialNodes and initialEdges)
const frontendDefaultPayload = `{
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
}`

// TestFrontendIntegration tests the exact payload format from the frontend
func TestFrontendIntegration(t *testing.T) {
	// Test that the backend can parse the frontend payload
	engine, err := NewEngine([]byte(frontendDefaultPayload))
	if err != nil {
		t.Fatalf("Failed to parse frontend payload: %v", err)
	}

	// Execute the workflow
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute frontend workflow: %v", err)
	}

	// Verify all nodes were executed
	if len(result.NodeResults) != 4 {
		t.Errorf("Expected 4 node results, got %d", len(result.NodeResults))
	}

	// Verify the final visualization output
	vizResult, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected visualization result to be a map, got %T", result.FinalOutput)
	}

	if vizResult["mode"] != "text" {
		t.Errorf("Expected mode 'text', got %v", vizResult["mode"])
	}

	value, ok := vizResult["value"].(float64)
	if !ok {
		t.Fatalf("Expected value to be float64, got %T", vizResult["value"])
	}

	if value != 15.0 {
		t.Errorf("Expected value 15, got %v", value)
	}

	// Verify no errors
	if len(result.Errors) > 0 {
		t.Errorf("Unexpected errors: %v", result.Errors)
	}
}

// TestFrontendDefaultWorkflow tests the initial default workflow from the frontend
func TestFrontendDefaultWorkflow(t *testing.T) {
	engine, err := NewEngine([]byte(frontendDefaultPayload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	// The default workflow: 10 + 5 = 15
	if result.NodeResults["3"] != 15.0 {
		t.Errorf("Expected node 3 result to be 15, got %v", result.NodeResults["3"])
	}
}

// TestFrontendAllNodeTypes tests all node types available in the frontend
func TestFrontendAllNodeTypes(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		expected interface{}
	}{
		{
			name: "number node only",
			payload: `{
				"nodes": [{"id": "1", "data": {"value": 42}}],
				"edges": []
			}`,
			expected: 42.0,
		},
		{
			name: "add operation",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 5}},
					{"id": "3", "data": {"op": "add"}}
				],
				"edges": [
					{"id": "e1-3", "source": "1", "target": "3"},
					{"id": "e2-3", "source": "2", "target": "3"}
				]
			}`,
			expected: 15.0,
		},
		{
			name: "subtract operation",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 5}},
					{"id": "3", "data": {"op": "subtract"}}
				],
				"edges": [
					{"id": "e1-3", "source": "1", "target": "3"},
					{"id": "e2-3", "source": "2", "target": "3"}
				]
			}`,
			expected: 5.0,
		},
		{
			name: "multiply operation",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 5}},
					{"id": "3", "data": {"op": "multiply"}}
				],
				"edges": [
					{"id": "e1-3", "source": "1", "target": "3"},
					{"id": "e2-3", "source": "2", "target": "3"}
				]
			}`,
			expected: 50.0,
		},
		{
			name: "divide operation",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 5}},
					{"id": "3", "data": {"op": "divide"}}
				],
				"edges": [
					{"id": "e1-3", "source": "1", "target": "3"},
					{"id": "e2-3", "source": "2", "target": "3"}
				]
			}`,
			expected: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine([]byte(tt.payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Failed to execute: %v", err)
			}

			if result.FinalOutput != tt.expected {
				t.Errorf("Expected final output %v, got %v", tt.expected, result.FinalOutput)
			}
		})
	}
}
