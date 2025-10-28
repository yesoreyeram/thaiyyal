package workflow

import (
	"encoding/json"
	"testing"
)

// Test creating engine with valid and invalid payloads
func TestNewEngine(t *testing.T) {
	validPayload := `{"nodes":[{"id":"1","data":{"value":10}}],"edges":[]}`
	_, err := NewEngine([]byte(validPayload))
	if err != nil {
		t.Errorf("NewEngine with valid payload failed: %v", err)
	}

	invalidPayload := `{invalid}`
	_, err = NewEngine([]byte(invalidPayload))
	if err == nil {
		t.Error("NewEngine should fail with invalid JSON")
	}
}

// Test simple addition workflow
func TestSimpleAddition(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "data": {"op": "add"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["3"] != 15.0 {
		t.Errorf("Expected 15, got %v", result.NodeResults["3"])
	}
}

// Test all arithmetic operations
func TestAllOperations(t *testing.T) {
	tests := []struct {
		op       string
		left     float64
		right    float64
		expected float64
		hasError bool
	}{
		{"add", 10, 5, 15, false},
		{"subtract", 10, 5, 5, false},
		{"multiply", 10, 5, 50, false},
		{"divide", 10, 5, 2, false},
		{"divide", 10, 0, 0, true}, // division by zero
	}

	for _, tt := range tests {
		t.Run(tt.op, func(t *testing.T) {
			payload := map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{"id": "1", "data": map[string]interface{}{"value": tt.left}},
					map[string]interface{}{"id": "2", "data": map[string]interface{}{"value": tt.right}},
					map[string]interface{}{"id": "3", "data": map[string]interface{}{"op": tt.op}},
				},
				"edges": []interface{}{
					map[string]interface{}{"id": "e1", "source": "1", "target": "3"},
					map[string]interface{}{"id": "e2", "source": "2", "target": "3"},
				},
			}
			jsonData, _ := json.Marshal(payload)

			engine, _ := NewEngine(jsonData)
			result, err := engine.Execute()

			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.hasError && result.NodeResults["3"] != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result.NodeResults["3"])
			}
		})
	}
}

// Test complete workflow with visualization
func TestCompleteWorkflow(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10, "label": "Node 1"}},
			{"id": "2", "data": {"value": 5, "label": "Node 2"}},
			{"id": "3", "data": {"op": "add", "label": "Add"}},
			{"id": "4", "data": {"mode": "text", "label": "Display"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-4", "source": "3", "target": "4"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Check visualization output
	vizResult, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatal("Final output should be a map")
	}
	if vizResult["mode"] != "text" {
		t.Errorf("Expected mode 'text', got %v", vizResult["mode"])
	}
	if vizResult["value"] != 15.0 {
		t.Errorf("Expected value 15, got %v", vizResult["value"])
	}
}

// Test multiple chained operations
func TestMultipleOperations(t *testing.T) {
	// (10 + 5) * 2 = 30
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "data": {"op": "add"}},
			{"id": "4", "data": {"value": 2}},
			{"id": "5", "data": {"op": "multiply"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-5", "source": "3", "target": "5"},
			{"id": "e4-5", "source": "4", "target": "5"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["5"] != 30.0 {
		t.Errorf("Expected 30, got %v", result.NodeResults["5"])
	}
}

// Test cyclic workflow detection
func TestCyclicWorkflow(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"op": "add"}},
			{"id": "3", "data": {"op": "multiply"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	if err == nil {
		t.Error("Expected error for cyclic workflow")
	}
}

// Test operation with missing inputs
func TestMissingInputs(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"op": "add"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	if err == nil {
		t.Error("Expected error for missing inputs")
	}
}

// Test explicit node types
func TestExplicitNodeTypes(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 10}},
			{"id": "2", "type": "number", "data": {"value": 5}},
			{"id": "3", "type": "operation", "data": {"op": "add"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "3"},
			{"id": "e2", "source": "2", "target": "3"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["3"] != 15.0 {
		t.Errorf("Expected 15, got %v", result.NodeResults["3"])
	}
}

// Test type inference
func TestTypeInference(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 42}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["1"] != 42.0 {
		t.Errorf("Expected 42, got %v", result.NodeResults["1"])
	}
}

// Test frontend default payload
func TestFrontendDefaultPayload(t *testing.T) {
	// This is the exact payload from the frontend
	payload := `{
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

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify all nodes executed
	if len(result.NodeResults) != 4 {
		t.Errorf("Expected 4 results, got %d", len(result.NodeResults))
	}

	// Verify no errors
	if len(result.Errors) > 0 {
		t.Errorf("Unexpected errors: %v", result.Errors)
	}

	// Verify addition result
	if result.NodeResults["3"] != 15.0 {
		t.Errorf("Expected 15, got %v", result.NodeResults["3"])
	}
}

// Test visualization modes
func TestVisualizationModes(t *testing.T) {
	modes := []string{"text", "table"}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			payload := map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{"id": "1", "data": map[string]interface{}{"value": 42.0}},
					map[string]interface{}{"id": "2", "data": map[string]interface{}{"mode": mode}},
				},
				"edges": []interface{}{
					map[string]interface{}{"id": "e1", "source": "1", "target": "2"},
				},
			}
			jsonData, _ := json.Marshal(payload)

			engine, _ := NewEngine(jsonData)
			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Execute failed: %v", err)
			}

			vizResult, ok := result.FinalOutput.(map[string]interface{})
			if !ok {
				t.Fatal("Final output should be a map")
			}

			if vizResult["mode"] != mode {
				t.Errorf("Expected mode %s, got %v", mode, vizResult["mode"])
			}
		})
	}
}
