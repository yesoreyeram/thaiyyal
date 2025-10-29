package workflow

import (
	"encoding/json"
	"testing"
)

// Test condition node with greater than condition - passes input through when condition is true
func TestConditionNodeGreaterThan(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 150}},
			{"id": "2", "type": "condition", "data": {"condition": ">100"}}
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

	condResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}

	if condResult["condition_met"] != true {
		t.Errorf("Expected condition_met to be true, got %v", condResult["condition_met"])
	}

	if condResult["value"] != 150.0 {
		t.Errorf("Expected value to be 150, got %v", condResult["value"])
	}
}

// Test condition node with less than condition - evaluates false when value is greater
func TestConditionNodeLessThan(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 50}},
			{"id": "2", "type": "condition", "data": {"condition": "<100"}}
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

	condResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}

	if condResult["condition_met"] != true {
		t.Errorf("Expected condition_met to be true, got %v", condResult["condition_met"])
	}
}

// Test condition node with greater than or equal condition
func TestConditionNodeGreaterThanOrEqual(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 100}},
			{"id": "2", "type": "condition", "data": {"condition": ">=100"}}
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

	condResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}

	if condResult["condition_met"] != true {
		t.Errorf("Expected condition_met to be true for value >= 100, got %v", condResult["condition_met"])
	}
}

// Test condition node with equality condition
func TestConditionNodeEquals(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 42}},
			{"id": "2", "type": "condition", "data": {"condition": "==42"}}
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

	condResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}

	if condResult["condition_met"] != true {
		t.Errorf("Expected condition_met to be true for value == 42, got %v", condResult["condition_met"])
	}
}

// Test condition node with not equal condition
func TestConditionNodeNotEquals(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 50}},
			{"id": "2", "type": "condition", "data": {"condition": "!=42"}}
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

	condResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}

	if condResult["condition_met"] != true {
		t.Errorf("Expected condition_met to be true for value != 42, got %v", condResult["condition_met"])
	}
}

// Test condition node with boolean true condition
func TestConditionNodeBooleanTrue(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "type": "condition", "data": {"condition": "true"}}
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

	condResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}

	if condResult["condition_met"] != true {
		t.Errorf("Expected condition_met to be true for 'true' condition, got %v", condResult["condition_met"])
	}
}

// Test condition node with boolean false condition
func TestConditionNodeBooleanFalse(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "type": "condition", "data": {"condition": "false"}}
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

	condResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}

	if condResult["condition_met"] != false {
		t.Errorf("Expected condition_met to be false for 'false' condition, got %v", condResult["condition_met"])
	}
}

// Test condition node integrated with arithmetic operation - validates condition on computed result
func TestConditionWithArithmeticOperation(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 50}},
			{"id": "2", "data": {"value": 60}},
			{"id": "3", "data": {"op": "add"}},
			{"id": "4", "type": "condition", "data": {"condition": ">100"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "3"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Check that addition worked
	if result.NodeResults["3"] != 110.0 {
		t.Errorf("Expected addition result to be 110, got %v", result.NodeResults["3"])
	}

	// Check condition evaluation
	condResult, ok := result.NodeResults["4"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}

	if condResult["condition_met"] != true {
		t.Errorf("Expected condition_met to be true for 110 > 100, got %v", condResult["condition_met"])
	}
}

// Test condition node with text input - validates that non-numeric inputs fail gracefully
func TestConditionWithTextInput(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "hello"}},
			{"id": "2", "type": "condition", "data": {"condition": ">100"}}
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

	// Non-numeric input should evaluate to false
	condResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}

	if condResult["condition_met"] != false {
		t.Errorf("Expected condition_met to be false for non-numeric input, got %v", condResult["condition_met"])
	}
}

// Test for_each node with array input - processes array and returns metadata
func TestForEachNodeWithArray(t *testing.T) {
	// Create payload with an array
	payloadMap := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id":   "1",
				"type": "number",
				"data": map[string]interface{}{"value": 0}, // Placeholder
			},
			map[string]interface{}{
				"id":   "2",
				"type": "for_each",
				"data": map[string]interface{}{},
			},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "2"},
		},
	}
	
	jsonData, _ := json.Marshal(payloadMap)
	engine, _ := NewEngine(jsonData)
	
	// Manually set node 1's result to an array
	engine.nodeResults["1"] = []interface{}{1.0, 2.0, 3.0, 4.0, 5.0}
	
	// Execute node 2 directly
	node2 := engine.getNode("2")
	result, err := engine.executeNode(node2)
	if err != nil {
		t.Fatalf("Execute for_each node failed: %v", err)
	}

	forEachResult, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected for_each result to be a map")
	}

	if forEachResult["count"] != 5 {
		t.Errorf("Expected count to be 5, got %v", forEachResult["count"])
	}

	items, ok := forEachResult["items"].([]interface{})
	if !ok || len(items) != 5 {
		t.Errorf("Expected 5 items in result, got %v", items)
	}
}

// Test for_each node with max_iterations limit - validates iteration limit enforcement
func TestForEachNodeMaxIterations(t *testing.T) {
	payloadMap := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id":   "1",
				"type": "number",
				"data": map[string]interface{}{"value": 0},
			},
			map[string]interface{}{
				"id":   "2",
				"type": "for_each",
				"data": map[string]interface{}{"max_iterations": 3},
			},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "2"},
		},
	}
	
	jsonData, _ := json.Marshal(payloadMap)
	engine, _ := NewEngine(jsonData)
	
	// Set array with 10 items, but max_iterations is 3
	largeArray := make([]interface{}, 10)
	for i := range largeArray {
		largeArray[i] = float64(i)
	}
	engine.nodeResults["1"] = largeArray
	
	node2 := engine.getNode("2")
	_, err := engine.executeNode(node2)
	
	// Should fail because array size exceeds max_iterations
	if err == nil {
		t.Error("Expected error for exceeding max_iterations")
	}
}

// Test for_each node with non-array input - validates type checking
func TestForEachNodeWithNonArray(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 42}},
			{"id": "2", "type": "for_each", "data": {}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	
	// Should fail because input is not an array
	if err == nil {
		t.Error("Expected error for non-array input to for_each")
	}
}

// Test for_each integrated with text operations - processes array of text values
func TestForEachWithTextArray(t *testing.T) {
	payloadMap := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"id":   "1",
				"type": "text_input",
				"data": map[string]interface{}{"text": "placeholder"},
			},
			map[string]interface{}{
				"id":   "2",
				"type": "for_each",
				"data": map[string]interface{}{},
			},
		},
		"edges": []interface{}{
			map[string]interface{}{"id": "e1", "source": "1", "target": "2"},
		},
	}
	
	jsonData, _ := json.Marshal(payloadMap)
	engine, _ := NewEngine(jsonData)
	
	// Set text array manually
	engine.nodeResults["1"] = []interface{}{"hello", "world", "test"}
	
	node2 := engine.getNode("2")
	result, err := engine.executeNode(node2)
	if err != nil {
		t.Fatalf("Execute for_each with text array failed: %v", err)
	}

	forEachResult, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected for_each result to be a map")
	}

	if forEachResult["count"] != 3 {
		t.Errorf("Expected count to be 3, got %v", forEachResult["count"])
	}
}

// Test while_loop node with simple condition - validates basic loop execution
// Note: Since the current implementation doesn't modify values in loop iterations,
// a condition that starts true will loop until max_iterations is reached
func TestWhileLoopNode(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 5}},
			{"id": "2", "type": "while_loop", "data": {"condition": "<10", "max_iterations": 10}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	
	// Should fail because condition is true (5 < 10) and value doesn't change,
	// so it will loop until hitting max_iterations
	if err == nil {
		t.Error("Expected error for while loop hitting max iterations with unchanging value")
	}
}

// Test while_loop with max_iterations limit - validates infinite loop prevention
func TestWhileLoopMaxIterations(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 5}},
			{"id": "2", "type": "while_loop", "data": {"condition": "true", "max_iterations": 5}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	
	// Should fail because condition is always true and will hit max_iterations
	if err == nil {
		t.Error("Expected error for exceeding max_iterations in while loop")
	}
}

// Test while_loop with false condition - validates immediate termination
func TestWhileLoopFalseCondition(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 150}},
			{"id": "2", "type": "while_loop", "data": {"condition": "<100", "max_iterations": 10}}
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

	whileResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected while_loop result to be a map")
	}

	// Should execute 0 iterations because condition is false from the start
	if whileResult["iterations"] != 0 {
		t.Errorf("Expected 0 iterations for false condition, got %v", whileResult["iterations"])
	}
}

// Test complex workflow: arithmetic -> condition -> visualization
// Validates integration of control flow with existing node types
func TestComplexWorkflowWithCondition(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 30}},
			{"id": "2", "data": {"value": 70}},
			{"id": "3", "data": {"op": "add"}},
			{"id": "4", "type": "condition", "data": {"condition": ">=100"}},
			{"id": "5", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "3"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify addition
	if result.NodeResults["3"] != 100.0 {
		t.Errorf("Expected addition to be 100, got %v", result.NodeResults["3"])
	}

	// Verify condition
	condResult, ok := result.NodeResults["4"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}
	if condResult["condition_met"] != true {
		t.Errorf("Expected condition to be true for 100 >= 100, got %v", condResult["condition_met"])
	}

	// Verify visualization receives condition result
	vizResult, ok := result.NodeResults["5"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected visualization result to be a map")
	}
	if vizResult["mode"] != "text" {
		t.Errorf("Expected visualization mode to be 'text', got %v", vizResult["mode"])
	}
}

// Test condition node chained with text operations
// Validates condition evaluation with text transformation pipeline
func TestConditionWithTextOperations(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "hello"}},
			{"id": "2", "data": {"text_op": "uppercase"}},
			{"id": "3", "type": "condition", "data": {"condition": "true"}}
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

	// Verify text transformation
	if result.NodeResults["2"] != "HELLO" {
		t.Errorf("Expected uppercase result to be 'HELLO', got %v", result.NodeResults["2"])
	}

	// Verify condition passes through the text
	condResult, ok := result.NodeResults["3"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition result to be a map")
	}
	if condResult["value"] != "HELLO" {
		t.Errorf("Expected condition to pass through 'HELLO', got %v", condResult["value"])
	}
}

// Test HTTP node with condition checking status
// Validates condition can be used to validate HTTP responses
func TestHTTPWithConditionValidation(t *testing.T) {
	// This test would require mocking HTTP, skipping detailed implementation
	// but demonstrates the pattern of HTTP -> condition validation
	t.Skip("HTTP with condition validation - requires HTTP test server setup")
}

// Test missing condition in condition node - validates error handling
func TestConditionNodeMissingCondition(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "type": "condition", "data": {}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	
	if err == nil {
		t.Error("Expected error for condition node without condition")
	}
}

// Test missing condition in while_loop node - validates error handling
func TestWhileLoopMissingCondition(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "type": "while_loop", "data": {}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, _ := NewEngine([]byte(payload))
	_, err := engine.Execute()
	
	if err == nil {
		t.Error("Expected error for while_loop node without condition")
	}
}

// Test control flow nodes with no inputs - validates error handling
func TestControlFlowNodesWithoutInputs(t *testing.T) {
	tests := []struct {
		name     string
		nodeType string
		data     map[string]interface{}
	}{
		{"condition", "condition", map[string]interface{}{"condition": ">10"}},
		{"for_each", "for_each", map[string]interface{}{}},
		{"while_loop", "while_loop", map[string]interface{}{"condition": ">10"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadMap := map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"id":   "1",
						"type": tt.nodeType,
						"data": tt.data,
					},
				},
				"edges": []interface{}{},
			}
			
			jsonData, _ := json.Marshal(payloadMap)
			engine, _ := NewEngine(jsonData)
			_, err := engine.Execute()
			
			if err == nil {
				t.Errorf("Expected error for %s node without inputs", tt.nodeType)
			}
		})
	}
}

// Test multiple conditions in series - validates chaining multiple condition nodes
func TestMultipleConditionsInSeries(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 75}},
			{"id": "2", "type": "condition", "data": {"condition": ">50"}},
			{"id": "3", "type": "condition", "data": {"condition": "<100"}}
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

	// Both conditions should pass
	cond2, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition 2 result to be a map")
	}
	if cond2["condition_met"] != true {
		t.Errorf("Expected first condition (>50) to be true, got %v", cond2["condition_met"])
	}

	cond3, ok := result.NodeResults["3"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected condition 3 result to be a map")
	}
	if cond3["condition_met"] != true {
		t.Errorf("Expected second condition (<100) to be true, got %v", cond3["condition_met"])
	}
}
