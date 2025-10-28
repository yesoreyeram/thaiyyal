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

// Test text input node
func TestTextInputNode(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "Hello World"}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["1"] != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", result.NodeResults["1"])
	}
}

// Test text operation node - uppercase
func TestTextOperationUppercase(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "hello world"}},
			{"id": "2", "data": {"text_op": "uppercase"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["2"] != "HELLO WORLD" {
		t.Errorf("Expected 'HELLO WORLD', got %v", result.NodeResults["2"])
	}
}

// Test text operation node - lowercase
func TestTextOperationLowercase(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "HELLO WORLD"}},
			{"id": "2", "data": {"text_op": "lowercase"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["2"] != "hello world" {
		t.Errorf("Expected 'hello world', got %v", result.NodeResults["2"])
	}
}

// Test text operation node - titlecase
func TestTextOperationTitlecase(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "hello world from go"}},
			{"id": "2", "data": {"text_op": "titlecase"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["2"] != "Hello World From Go" {
		t.Errorf("Expected 'Hello World From Go', got %v", result.NodeResults["2"])
	}
}

// Test text operation node - camelcase
func TestTextOperationCamelcase(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "hello world from go"}},
			{"id": "2", "data": {"text_op": "camelcase"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["2"] != "helloWorldFromGo" {
		t.Errorf("Expected 'helloWorldFromGo', got %v", result.NodeResults["2"])
	}
}

// Test text operation node - inversecase
func TestTextOperationInversecase(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "HeLLo WoRLd"}},
			{"id": "2", "data": {"text_op": "inversecase"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["2"] != "hEllO wOrlD" {
		t.Errorf("Expected 'hEllO wOrlD', got %v", result.NodeResults["2"])
	}
}

// Test chained text operations
func TestChainedTextOperations(t *testing.T) {
	// Input: "hello world" -> camelcase -> inversecase
	// Expected: "helloWorld" -> "HELLOwORLD"
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "hello world"}},
			{"id": "2", "data": {"text_op": "camelcase"}},
			{"id": "3", "data": {"text_op": "inversecase"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// First operation: camelcase "hello world" -> "helloWorld"
	if result.NodeResults["2"] != "helloWorld" {
		t.Errorf("Expected 'helloWorld' at node 2, got %v", result.NodeResults["2"])
	}

	// Second operation: inversecase "helloWorld" -> "HELLOwORLD"
	if result.NodeResults["3"] != "HELLOwORLD" {
		t.Errorf("Expected 'HELLOwORLD' at node 3, got %v", result.NodeResults["3"])
	}
}

// Test text operation with non-text input (should fail)
func TestTextOperationNonTextInput(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 42}},
			{"id": "2", "data": {"text_op": "uppercase"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	if err == nil {
		t.Error("Expected error when text operation receives non-text input")
	}
	if err != nil && err.Error() != "error executing node 2: text operation input must be text/string" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

// Test explicit text node types
func TestExplicitTextNodeTypes(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "type": "text_input", "data": {"text": "test"}},
			{"id": "2", "type": "text_operation", "data": {"text_op": "uppercase"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["2"] != "TEST" {
		t.Errorf("Expected 'TEST', got %v", result.NodeResults["2"])
	}
}

// Test complex chained text operations
func TestComplexTextChain(t *testing.T) {
	// "HELLO WORLD" -> lowercase -> titlecase -> camelcase
	// "HELLO WORLD" -> "hello world" -> "Hello World" -> "helloWorld"
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "HELLO WORLD"}},
			{"id": "2", "data": {"text_op": "lowercase"}},
			{"id": "3", "data": {"text_op": "titlecase"}},
			{"id": "4", "data": {"text_op": "camelcase"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["2"] != "hello world" {
		t.Errorf("Step 1: Expected 'hello world', got %v", result.NodeResults["2"])
	}

	if result.NodeResults["3"] != "Hello World" {
		t.Errorf("Step 2: Expected 'Hello World', got %v", result.NodeResults["3"])
	}

	if result.NodeResults["4"] != "helloWorld" {
		t.Errorf("Step 3: Expected 'helloWorld', got %v", result.NodeResults["4"])
	}
}
