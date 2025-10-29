package workflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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

// Test HTTP node with successful response
func TestHTTPNodeSuccess(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from server"))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"url": server.URL}},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["1"] != "Hello from server" {
		t.Errorf("Expected 'Hello from server', got %v", result.NodeResults["1"])
	}
}

// Test HTTP node with error status code
func TestHTTPNodeErrorStatus(t *testing.T) {
	// Create test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"url": server.URL}},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	_, err := engine.Execute()
	if err == nil {
		t.Error("Expected error for 404 status code")
	}
	if err != nil && err.Error() != "error executing node 1: HTTP request returned error status: 404" {
		t.Logf("Got error: %v", err)
	}
}

// Test HTTP node with invalid URL
func TestHTTPNodeInvalidURL(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"url": "http://invalid-url-that-does-not-exist.local"}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

// Test HTTP node output passed to text operation
func TestHTTPNodeToTextOperation(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"url": server.URL}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text_op": "uppercase"}},
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

	// Check HTTP node result
	if result.NodeResults["1"] != "hello world" {
		t.Errorf("Expected 'hello world' from HTTP, got %v", result.NodeResults["1"])
	}

	// Check text operation result
	if result.NodeResults["2"] != "HELLO WORLD" {
		t.Errorf("Expected 'HELLO WORLD' from text operation, got %v", result.NodeResults["2"])
	}
}

// Test HTTP node error followed by text operation (should fail)
func TestHTTPNodeErrorToTextOperation(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"url": server.URL}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text_op": "uppercase"}},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "2"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	_, err := engine.Execute()
	if err == nil {
		t.Error("Expected error when HTTP node fails")
	}
}

// Test chained HTTP to multiple text operations
func TestHTTPNodeToChainedTextOperations(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HELLO WORLD"))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"url": server.URL}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text_op": "lowercase"}},
			map[string]interface{}{"id": "3", "data": map[string]interface{}{"text_op": "titlecase"}},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "2"},
			map[string]interface{}{"id": "e2", "source": "2", "target": "3"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// HELLO WORLD -> hello world -> Hello World
	if result.NodeResults["1"] != "HELLO WORLD" {
		t.Errorf("HTTP result: Expected 'HELLO WORLD', got %v", result.NodeResults["1"])
	}

	if result.NodeResults["2"] != "hello world" {
		t.Errorf("Lowercase result: Expected 'hello world', got %v", result.NodeResults["2"])
	}

	if result.NodeResults["3"] != "Hello World" {
		t.Errorf("Titlecase result: Expected 'Hello World', got %v", result.NodeResults["3"])
	}
}

// Test explicit HTTP node type
func TestExplicitHTTPNodeType(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "type": "http", "data": map[string]interface{}{"url": server.URL}},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["1"] != "test response" {
		t.Errorf("Expected 'test response', got %v", result.NodeResults["1"])
	}
}

// Test HTTP node with different status codes
func TestHTTPNodeStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		shouldFail bool
	}{
		{"200 OK", http.StatusOK, false},
		{"201 Created", http.StatusCreated, false},
		{"204 No Content", http.StatusNoContent, false},
		{"400 Bad Request", http.StatusBadRequest, true},
		{"404 Not Found", http.StatusNotFound, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte("response body"))
			}))
			defer server.Close()

			payload := map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{"id": "1", "data": map[string]interface{}{"url": server.URL}},
				},
				"edges": []interface{}{},
			}
			jsonData, _ := json.Marshal(payload)

			engine, _ := NewEngine(jsonData)
			_, err := engine.Execute()

			if tt.shouldFail && err == nil {
				t.Errorf("Expected error for status code %d", tt.statusCode)
			}
			if !tt.shouldFail && err != nil {
				t.Errorf("Unexpected error for status code %d: %v", tt.statusCode, err)
			}
		})
	}
}

// Test concat operation with two inputs
func TestTextOperationConcat(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "Hello"}},
			{"id": "2", "data": {"text": "World"}},
			{"id": "3", "data": {"text_op": "concat"}}
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

	if result.NodeResults["3"] != "HelloWorld" {
		t.Errorf("Expected 'HelloWorld', got %v", result.NodeResults["3"])
	}
}

// Test concat operation with separator
func TestTextOperationConcatWithSeparator(t *testing.T) {
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"text": "Hello"}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text": "World"}},
			map[string]interface{}{"id": "3", "data": map[string]interface{}{"text_op": "concat", "separator": " "}},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "3"},
			map[string]interface{}{"id": "e2", "source": "2", "target": "3"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["3"] != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", result.NodeResults["3"])
	}
}

// Test concat with multiple inputs and custom separator
func TestTextOperationConcatMultiple(t *testing.T) {
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"text": "one"}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text": "two"}},
			map[string]interface{}{"id": "3", "data": map[string]interface{}{"text": "three"}},
			map[string]interface{}{"id": "4", "data": map[string]interface{}{"text_op": "concat", "separator": ", "}},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "4"},
			map[string]interface{}{"id": "e2", "source": "2", "target": "4"},
			map[string]interface{}{"id": "e3", "source": "3", "target": "4"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["4"] != "one, two, three" {
		t.Errorf("Expected 'one, two, three', got %v", result.NodeResults["4"])
	}
}

// Test concat with non-text input (should fail)
func TestTextOperationConcatNonTextInput(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 42}},
			{"id": "2", "data": {"text": "World"}},
			{"id": "3", "data": {"text_op": "concat"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "3"},
			{"id": "e2", "source": "2", "target": "3"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	if err == nil {
		t.Error("Expected error when concat receives non-text input")
	}
}

// Test repeat operation
func TestTextOperationRepeat(t *testing.T) {
	repeatN := 3
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"text": "Ha"}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text_op": "repeat", "repeat_n": repeatN}},
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

	if result.NodeResults["2"] != "HaHaHa" {
		t.Errorf("Expected 'HaHaHa', got %v", result.NodeResults["2"])
	}
}

// Test repeat operation with zero count
func TestTextOperationRepeatZero(t *testing.T) {
	repeatN := 0
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"text": "Hello"}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text_op": "repeat", "repeat_n": repeatN}},
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

	if result.NodeResults["2"] != "" {
		t.Errorf("Expected empty string, got %v", result.NodeResults["2"])
	}
}

// Test repeat operation without repeat_n (should fail)
func TestTextOperationRepeatMissingN(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "Hello"}},
			{"id": "2", "data": {"text_op": "repeat"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	if err == nil {
		t.Error("Expected error when repeat_n is missing")
	}
}

// Test repeat with negative count (should fail)
func TestTextOperationRepeatNegative(t *testing.T) {
	repeatN := -1
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"text": "Hello"}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text_op": "repeat", "repeat_n": repeatN}},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "2"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	_, err := engine.Execute()
	if err == nil {
		t.Error("Expected error for negative repeat_n")
	}
}

// Test chained concat and repeat operations
func TestConcatAndRepeatChained(t *testing.T) {
	repeatN := 2
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"text": "Hello"}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text": "World"}},
			map[string]interface{}{"id": "3", "data": map[string]interface{}{"text_op": "concat", "separator": " "}},
			map[string]interface{}{"id": "4", "data": map[string]interface{}{"text_op": "repeat", "repeat_n": repeatN}},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "3"},
			map[string]interface{}{"id": "e2", "source": "2", "target": "3"},
			map[string]interface{}{"id": "e3", "source": "3", "target": "4"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// "Hello" + " " + "World" = "Hello World"
	// "Hello World" repeated 2 times = "Hello WorldHello World"
	if result.NodeResults["3"] != "Hello World" {
		t.Errorf("Concat result: Expected 'Hello World', got %v", result.NodeResults["3"])
	}

	if result.NodeResults["4"] != "Hello WorldHello World" {
		t.Errorf("Repeat result: Expected 'Hello WorldHello World', got %v", result.NodeResults["4"])
	}
}

// Test HTTP output to concat
func TestHTTPToConcat(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API"))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"url": server.URL}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text": "Response"}},
			map[string]interface{}{"id": "3", "data": map[string]interface{}{"text_op": "concat", "separator": ": "}},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "3"},
			map[string]interface{}{"id": "e2", "source": "2", "target": "3"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.NodeResults["3"] != "API: Response" {
		t.Errorf("Expected 'API: Response', got %v", result.NodeResults["3"])
	}
}

// Test complex workflow: HTTP -> uppercase -> repeat -> concat
func TestComplexTextWorkflow(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hi"))
	}))
	defer server.Close()

	repeatN := 3
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "data": map[string]interface{}{"url": server.URL}},
			map[string]interface{}{"id": "2", "data": map[string]interface{}{"text_op": "uppercase"}},
			map[string]interface{}{"id": "3", "data": map[string]interface{}{"text_op": "repeat", "repeat_n": repeatN}},
			map[string]interface{}{"id": "4", "data": map[string]interface{}{"text": "!!!"}},
			map[string]interface{}{"id": "5", "data": map[string]interface{}{"text_op": "concat"}},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "2"},
			map[string]interface{}{"id": "e2", "source": "2", "target": "3"},
			map[string]interface{}{"id": "e3", "source": "3", "target": "5"},
			map[string]interface{}{"id": "e4", "source": "4", "target": "5"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// "hi" -> "HI" -> "HIHIHI" -> "HIHIHI" + "!!!" = "HIHIHI!!!"
	if result.NodeResults["2"] != "HI" {
		t.Errorf("Uppercase: Expected 'HI', got %v", result.NodeResults["2"])
	}

	if result.NodeResults["3"] != "HIHIHI" {
		t.Errorf("Repeat: Expected 'HIHIHI', got %v", result.NodeResults["3"])
	}

	if result.NodeResults["5"] != "HIHIHI!!!" {
		t.Errorf("Concat: Expected 'HIHIHI!!!', got %v", result.NodeResults["5"])
	}
}

// ====================================
// HTTP Pagination Node Tests
// ====================================

// Test HTTP pagination with max_pages
func TestHTTPPaginationWithMaxPages(t *testing.T) {
	// Track pages requested
	pagesRequested := []int{}
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract page parameter
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}
		pageNum, _ := strconv.Atoi(page)
		pagesRequested = append(pagesRequested, pageNum)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"page":%s,"data":["item1","item2"]}`, page)))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id": "1",
				"type": "http_pagination",
				"data": map[string]interface{}{
					"base_url":  server.URL,
					"max_pages": 5,
				},
			},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Check result
	nodeResult := result.NodeResults["1"].(map[string]interface{})
	
	if !nodeResult["success"].(bool) {
		t.Errorf("Expected success to be true")
	}
	
	if nodeResult["pages_fetched"].(int) != 5 {
		t.Errorf("Expected 5 pages fetched, got %v", nodeResult["pages_fetched"])
	}
	
	results := nodeResult["results"].([]interface{})
	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	if len(pagesRequested) != 5 {
		t.Errorf("Expected 5 pages requested, got %d", len(pagesRequested))
	}
}

// Test HTTP pagination with total_items and page_size
func TestHTTPPaginationWithTotalItems(t *testing.T) {
	pagesRequested := []int{}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}
		pageNum, _ := strconv.Atoi(page)
		pagesRequested = append(pagesRequested, pageNum)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"page":%s,"items":10}`, page)))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id": "1",
				"type": "http_pagination",
				"data": map[string]interface{}{
					"base_url":    server.URL,
					"total_items": 50,
					"page_size":   10,
				},
			},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	nodeResult := result.NodeResults["1"].(map[string]interface{})
	
	// Should fetch 5 pages (50 items / 10 items per page = 5)
	if nodeResult["pages_fetched"].(int) != 5 {
		t.Errorf("Expected 5 pages fetched, got %v", nodeResult["pages_fetched"])
	}
	
	if len(pagesRequested) != 5 {
		t.Errorf("Expected 5 pages requested, got %d", len(pagesRequested))
	}
}

// Test HTTP pagination with URL placeholder
func TestHTTPPaginationWithPlaceholder(t *testing.T) {
	pagesRequested := []int{}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract page from path
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			page := parts[len(parts)-1]
			pageNum, _ := strconv.Atoi(page)
			pagesRequested = append(pagesRequested, pageNum)
			
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"page":%s}`, page)))
		}
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id": "1",
				"type": "http_pagination",
				"data": map[string]interface{}{
					"base_url":  server.URL + "/api/items/{page}",
					"max_pages": 3,
				},
			},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	nodeResult := result.NodeResults["1"].(map[string]interface{})
	
	if nodeResult["pages_fetched"].(int) != 3 {
		t.Errorf("Expected 3 pages fetched, got %v", nodeResult["pages_fetched"])
	}
	
	// Check that pages 1, 2, 3 were requested
	expectedPages := []int{1, 2, 3}
	if !intSlicesEqual(pagesRequested, expectedPages) {
		t.Errorf("Expected pages %v, got %v", expectedPages, pagesRequested)
	}
}

// Test HTTP pagination with custom start page
func TestHTTPPaginationWithStartPage(t *testing.T) {
	pagesRequested := []int{}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		pageNum, _ := strconv.Atoi(page)
		pagesRequested = append(pagesRequested, pageNum)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"page":%s}`, page)))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id": "1",
				"type": "http_pagination",
				"data": map[string]interface{}{
					"base_url":   server.URL,
					"start_page": 5,
					"max_pages":  3,
				},
			},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	nodeResult := result.NodeResults["1"].(map[string]interface{})
	
	if nodeResult["pages_fetched"].(int) != 3 {
		t.Errorf("Expected 3 pages fetched, got %v", nodeResult["pages_fetched"])
	}
	
	// Should request pages 5, 6, 7
	expectedPages := []int{5, 6, 7}
	if !intSlicesEqual(pagesRequested, expectedPages) {
		t.Errorf("Expected pages %v, got %v", expectedPages, pagesRequested)
	}
}

// Test HTTP pagination breaks on error
func TestHTTPPaginationBreakOnError(t *testing.T) {
	requestCount := 0
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		page := r.URL.Query().Get("page")
		pageNum, _ := strconv.Atoi(page)
		
		// Fail on page 3
		if pageNum == 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"page":%s}`, page)))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id": "1",
				"type": "http_pagination",
				"data": map[string]interface{}{
					"base_url":       server.URL,
					"max_pages":      5,
					"break_on_error": true,
				},
			},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	
	// Should return error since break_on_error is true
	if err == nil {
		t.Error("Expected error when page 3 fails")
	}
	
	// Should have made 3 requests (pages 1, 2, 3)
	if requestCount != 3 {
		t.Errorf("Expected 3 requests, got %d", requestCount)
	}
	
	// Check partial results in result - only if result is not nil
	if result != nil && result.NodeResults["1"] != nil {
		nodeResult := result.NodeResults["1"].(map[string]interface{})
		if nodeResult["success"].(bool) {
			t.Error("Expected success to be false")
		}
		
		if nodeResult["pages_fetched"].(int) != 2 {
			t.Errorf("Expected 2 pages fetched before error, got %v", nodeResult["pages_fetched"])
		}
	}
}

// Test HTTP pagination continues on error
func TestHTTPPaginationContinueOnError(t *testing.T) {
	requestCount := 0
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		page := r.URL.Query().Get("page")
		pageNum, _ := strconv.Atoi(page)
		
		// Fail on page 2
		if pageNum == 2 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"page":%s}`, page)))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id": "1",
				"type": "http_pagination",
				"data": map[string]interface{}{
					"base_url":       server.URL,
					"max_pages":      4,
					"break_on_error": false,
				},
			},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	result, err := engine.Execute()
	
	// Should not return error since break_on_error is false
	if err != nil {
		t.Errorf("Expected no error with break_on_error=false, got: %v", err)
	}
	
	// Should have made 4 requests (all pages)
	if requestCount != 4 {
		t.Errorf("Expected 4 requests, got %d", requestCount)
	}
	
	nodeResult := result.NodeResults["1"].(map[string]interface{})
	
	// Should have 3 successful pages (1, 3, 4)
	if nodeResult["pages_fetched"].(int) != 3 {
		t.Errorf("Expected 3 pages fetched, got %v", nodeResult["pages_fetched"])
	}
	
	// Should have 1 error
	if nodeResult["error_count"].(int) != 1 {
		t.Errorf("Expected 1 error, got %v", nodeResult["error_count"])
	}
	
	// success should be false because there was an error
	if nodeResult["success"].(bool) {
		t.Error("Expected success to be false when there are errors")
	}
}

// Test HTTP pagination with custom page parameter name
func TestHTTPPaginationCustomPageParam(t *testing.T) {
	var capturedParam string
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should have "p" parameter instead of "page"
		if r.URL.Query().Has("p") {
			capturedParam = "p"
		}
		
		pageNum := r.URL.Query().Get("p")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"p":%s}`, pageNum)))
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id": "1",
				"type": "http_pagination",
				"data": map[string]interface{}{
					"base_url":   server.URL,
					"max_pages":  2,
					"page_param": "p",
				},
			},
		},
		"edges": []interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	engine, _ := NewEngine(jsonData)
	_, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	if capturedParam != "p" {
		t.Errorf("Expected custom page param 'p', got '%s'", capturedParam)
	}
}

// Test HTTP pagination missing required fields
func TestHTTPPaginationMissingFields(t *testing.T) {
	tests := []struct {
		name string
		data map[string]interface{}
		expectedError string
	}{
		{
			name: "missing base_url",
			data: map[string]interface{}{
				"max_pages": 5,
			},
			expectedError: "http_pagination node missing base_url",
		},
		{
			name: "missing both max_pages and total_items",
			data: map[string]interface{}{
				"base_url": "http://example.com",
			},
			expectedError: "http_pagination node requires either max_pages or both total_items and page_size",
		},
		{
			name: "missing page_size with total_items",
			data: map[string]interface{}{
				"base_url":    "http://example.com",
				"total_items": 50,
			},
			expectedError: "http_pagination node requires either max_pages or both total_items and page_size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"id":   "1",
						"type": "http_pagination",
						"data": tt.data,
					},
				},
				"edges": []interface{}{},
			}
			jsonData, _ := json.Marshal(payload)

			engine, _ := NewEngine(jsonData)
			_, err := engine.Execute()
			
			if err == nil {
				t.Errorf("Expected error for %s", tt.name)
			} else if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing '%s', got: %v", tt.expectedError, err)
			}
		})
	}
}

// Helper function to compare int slices
func intSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
