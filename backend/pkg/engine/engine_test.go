package engine

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestNew(t *testing.T) {
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

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	if engine == nil {
		t.Fatal("Engine is nil")
	}

	if len(engine.nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(engine.nodes))
	}

	if len(engine.edges) != 2 {
		t.Errorf("Expected 2 edges, got %d", len(engine.edges))
	}
}

func TestExecute(t *testing.T) {
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

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	if result.FinalOutput != 15.0 {
		t.Errorf("Expected final output 15, got %v", result.FinalOutput)
	}
}

func TestExecutionContext(t *testing.T) {
	payload := `{
"nodes": [
{"id": "1", "type": "number", "data": {"value": 42}}
],
"edges": []
}`

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Test GetConfig
	config := engine.GetConfig()
	if config.MaxExecutionTime == 0 {
		t.Error("Config has zero MaxExecutionTime")
	}

	// Test state operations
	engine.SetVariable("test", "value")
	val, err := engine.GetVariable("test")
	if err != nil {
		t.Errorf("Failed to get variable: %v", err)
	}
	if val != "value" {
		t.Errorf("Expected 'value', got %v", val)
	}

	// Test counter
	engine.SetCounter(10.5)
	if engine.GetCounter() != 10.5 {
		t.Errorf("Expected counter 10.5, got %f", engine.GetCounter())
	}

	// Test context variables
	engine.SetContextVariable("ctx_var", "test_value")
	ctxVal, exists := engine.GetContextVariable("ctx_var")
	if !exists {
		t.Error("Context variable not found")
	}
	if ctxVal != "test_value" {
		t.Errorf("Expected 'test_value', got %v", ctxVal)
	}
}

func TestInferNodeTypes(t *testing.T) {
	payload := `{
"nodes": [
{"id": "1", "data": {"value": 10}},
{"id": "2", "data": {"text": "hello"}},
{"id": "3", "data": {"op": "add"}}
],
"edges": []
}`

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Type inference now happens automatically during JSON unmarshaling
	// Just verify the types were inferred correctly

	if engine.nodes[0].Type != types.NodeTypeNumber {
		t.Errorf("Expected type 'number', got %s", engine.nodes[0].Type)
	}
	if engine.nodes[1].Type != types.NodeTypeTextInput {
		t.Errorf("Expected type 'text_input', got %s", engine.nodes[1].Type)
	}
	if engine.nodes[2].Type != types.NodeTypeOperation {
		t.Errorf("Expected type 'operation', got %s", engine.nodes[2].Type)
	}
}
