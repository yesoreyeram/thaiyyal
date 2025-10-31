package workflow

import (
	"testing"
)

// TestVariableNodeSet tests storing a value in a variable
func TestVariableNodeSet(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 42}},
			{"id": "2", "data": {"var_name": "myvar", "var_op": "set"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Check that variable was set
	varResult, ok := result.NodeResults["2"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for variable set, got %T", result.NodeResults["2"])
	}

	if varResult["operation"] != "set" {
		t.Errorf("Expected operation 'set', got %v", varResult["operation"])
	}

	if varResult["value"] != 42.0 {
		t.Errorf("Expected value 42, got %v", varResult["value"])
	}

	// Check internal variable storage
	if engine.variables["myvar"] != 42.0 {
		t.Errorf("Expected variable 'myvar' to be 42, got %v", engine.variables["myvar"])
	}
}

// TestVariableNodeGetAndSet tests setting and retrieving a variable
func TestVariableNodeGetAndSet(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "Hello World"}},
			{"id": "2", "data": {"var_name": "greeting", "var_op": "set"}},
			{"id": "3", "data": {"var_name": "greeting", "var_op": "get"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Check get result
	getResult, ok := result.NodeResults["3"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for variable get, got %T", result.NodeResults["3"])
	}

	if getResult["operation"] != "get" {
		t.Errorf("Expected operation 'get', got %v", getResult["operation"])
	}

	if getResult["value"] != "Hello World" {
		t.Errorf("Expected value 'Hello World', got %v", getResult["value"])
	}
}

// TestVariableNodeGetNonExistent tests retrieving a non-existent variable
func TestVariableNodeGetNonExistent(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"var_name": "nonexistent", "var_op": "get"}}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	_, err = engine.Execute()
	if err == nil {
		t.Fatal("Expected error for non-existent variable, got nil")
	}
}

// TestVariableNodeMissingConfig tests variable node with missing configuration
func TestVariableNodeMissingConfig(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "type": "variable", "data": {"var_name": "test"}}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	_, err = engine.Execute()
	if err == nil {
		t.Fatal("Expected error for missing var_op, got nil")
	}
}

// TestTransformNodeToArray tests converting inputs to an array
func TestTransformNodeToArray(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 20}},
			{"id": "3", "data": {"value": 30}},
			{"id": "4", "type": "transform", "data": {"transform_type": "to_array"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "4"},
			{"id": "e2", "source": "2", "target": "4"},
			{"id": "e3", "source": "3", "target": "4"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	arrResult, ok := result.NodeResults["4"].([]interface{})
	if !ok {
		t.Fatalf("Expected array result for to_array transform, got %T", result.NodeResults["4"])
	}

	if len(arrResult) != 3 {
		t.Errorf("Expected array length 3, got %d", len(arrResult))
	}
}

// TestAccumulatorNodeSum tests accumulating numbers with sum operation
func TestAccumulatorNodeSum(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "type": "accumulator", "data": {"accum_op": "sum"}},
			{"id": "3", "data": {"value": 20}},
			{"id": "4", "type": "accumulator", "data": {"accum_op": "sum"}},
			{"id": "5", "data": {"value": 30}},
			{"id": "6", "type": "accumulator", "data": {"accum_op": "sum"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"},
			{"id": "e5", "source": "5", "target": "6"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Check final accumulator result
	finalResult, ok := result.NodeResults["6"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for accumulator, got %T", result.NodeResults["6"])
	}

	// Should have accumulated 10 + 20 + 30 = 60
	if finalResult["value"] != 60.0 {
		t.Errorf("Expected accumulated value 60, got %v", finalResult["value"])
	}
}

// TestAccumulatorNodeProduct tests accumulating numbers with product operation
func TestAccumulatorNodeProduct(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 2}},
			{"id": "2", "type": "accumulator", "data": {"accum_op": "product"}},
			{"id": "3", "data": {"value": 3}},
			{"id": "4", "type": "accumulator", "data": {"accum_op": "product"}},
			{"id": "5", "data": {"value": 4}},
			{"id": "6", "type": "accumulator", "data": {"accum_op": "product"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"},
			{"id": "e5", "source": "5", "target": "6"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	finalResult, ok := result.NodeResults["6"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for accumulator, got %T", result.NodeResults["6"])
	}

	// Should have accumulated 1 * 2 * 3 * 4 = 24
	if finalResult["value"] != 24.0 {
		t.Errorf("Expected accumulated value 24, got %v", finalResult["value"])
	}
}

// TestAccumulatorNodeConcat tests accumulating strings with concat operation
func TestAccumulatorNodeConcat(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "Hello"}},
			{"id": "2", "type": "accumulator", "data": {"accum_op": "concat"}},
			{"id": "3", "data": {"text": " "}},
			{"id": "4", "type": "accumulator", "data": {"accum_op": "concat"}},
			{"id": "5", "data": {"text": "World"}},
			{"id": "6", "type": "accumulator", "data": {"accum_op": "concat"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"},
			{"id": "e5", "source": "5", "target": "6"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	finalResult, ok := result.NodeResults["6"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for accumulator, got %T", result.NodeResults["6"])
	}

	if finalResult["value"] != "Hello World" {
		t.Errorf("Expected accumulated value 'Hello World', got %v", finalResult["value"])
	}
}

// TestAccumulatorNodeArray tests accumulating values into an array
func TestAccumulatorNodeArray(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "type": "accumulator", "data": {"accum_op": "array"}},
			{"id": "3", "data": {"value": 20}},
			{"id": "4", "type": "accumulator", "data": {"accum_op": "array"}},
			{"id": "5", "data": {"value": 30}},
			{"id": "6", "type": "accumulator", "data": {"accum_op": "array"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"},
			{"id": "e5", "source": "5", "target": "6"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	finalResult, ok := result.NodeResults["6"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for accumulator, got %T", result.NodeResults["6"])
	}

	arr, ok := finalResult["value"].([]interface{})
	if !ok {
		t.Fatalf("Expected array value in accumulator, got %T", finalResult["value"])
	}

	expected := []float64{10.0, 20.0, 30.0}
	if len(arr) != len(expected) {
		t.Errorf("Expected array length %d, got %d", len(expected), len(arr))
	}

	for i, v := range expected {
		if arr[i] != v {
			t.Errorf("Expected arr[%d] = %v, got %v", i, v, arr[i])
		}
	}
}

// TestAccumulatorNodeCount tests counting with accumulator
func TestAccumulatorNodeCount(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "item1"}},
			{"id": "2", "type": "accumulator", "data": {"accum_op": "count"}},
			{"id": "3", "data": {"text": "item2"}},
			{"id": "4", "type": "accumulator", "data": {"accum_op": "count"}},
			{"id": "5", "data": {"text": "item3"}},
			{"id": "6", "type": "accumulator", "data": {"accum_op": "count"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"},
			{"id": "e5", "source": "5", "target": "6"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	finalResult, ok := result.NodeResults["6"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for accumulator, got %T", result.NodeResults["6"])
	}

	if finalResult["value"] != 3.0 {
		t.Errorf("Expected count 3, got %v", finalResult["value"])
	}
}

// TestCounterNodeIncrement tests incrementing a counter
func TestCounterNodeIncrement(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "type": "counter", "data": {"counter_op": "increment"}},
			{"id": "2", "type": "counter", "data": {"counter_op": "increment"}},
			{"id": "3", "type": "counter", "data": {"counter_op": "increment"}},
			{"id": "4", "type": "counter", "data": {"counter_op": "get"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	finalResult, ok := result.NodeResults["4"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for counter, got %T", result.NodeResults["4"])
	}

	if finalResult["value"] != 3.0 {
		t.Errorf("Expected counter value 3, got %v", finalResult["value"])
	}
}

// TestCounterNodeDecrement tests decrementing a counter
func TestCounterNodeDecrement(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "type": "counter", "data": {"counter_op": "increment", "delta": 10}},
			{"id": "2", "type": "counter", "data": {"counter_op": "decrement", "delta": 3}},
			{"id": "3", "type": "counter", "data": {"counter_op": "decrement", "delta": 2}},
			{"id": "4", "type": "counter", "data": {"counter_op": "get"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	finalResult, ok := result.NodeResults["4"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for counter, got %T", result.NodeResults["4"])
	}

	// 10 - 3 - 2 = 5
	if finalResult["value"] != 5.0 {
		t.Errorf("Expected counter value 5, got %v", finalResult["value"])
	}
}

// TestCounterNodeReset tests resetting a counter
func TestCounterNodeReset(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "type": "counter", "data": {"counter_op": "increment"}},
			{"id": "2", "type": "counter", "data": {"counter_op": "increment"}},
			{"id": "3", "type": "counter", "data": {"counter_op": "reset"}},
			{"id": "4", "type": "counter", "data": {"counter_op": "increment"}},
			{"id": "5", "type": "counter", "data": {"counter_op": "get"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	finalResult, ok := result.NodeResults["5"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for counter, got %T", result.NodeResults["5"])
	}

	// Incremented twice, reset to 0, then incremented once = 1
	if finalResult["value"] != 1.0 {
		t.Errorf("Expected counter value 1 after reset, got %v", finalResult["value"])
	}
}

// TestVariableWithArithmetic tests variable storage and retrieval with arithmetic operations
func TestVariableWithArithmetic(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 20}},
			{"id": "3", "data": {"op": "add"}},
			{"id": "4", "data": {"var_name": "sum", "var_op": "set"}},
			{"id": "5", "data": {"var_name": "sum", "var_op": "get"}},
			{"id": "6", "type": "extract", "data": {"field": "value"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "3"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"},
			{"id": "e5", "source": "5", "target": "6"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// 10 + 20 = 30, stored, retrieved, extracted
	extractResult, ok := result.NodeResults["6"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result.NodeResults["6"])
	}

	if extractResult["value"] != 30.0 {
		t.Errorf("Expected extracted value 30, got %v", extractResult["value"])
	}
}

// TestExtractFieldFromVariableResult tests extracting field from variable get result
func TestExtractFieldFromVariableResult(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 99}},
			{"id": "2", "data": {"var_name": "number", "var_op": "set"}},
			{"id": "3", "data": {"var_name": "number", "var_op": "get"}},
			{"id": "4", "type": "extract", "data": {"field": "value"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	extractResult, ok := result.NodeResults["4"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for extract, got %T", result.NodeResults["4"])
	}

	if extractResult["value"] != 99.0 {
		t.Errorf("Expected extracted value 99, got %v", extractResult["value"])
	}
}

// TestTransformToObjectAndExtract tests creating an object and extracting values
func TestTransformToObjectAndExtract(t *testing.T) {
	// This test demonstrates creating an object from array and extracting fields
	// Skipped for now as it requires more complex setup
	t.Skip("Complex transform test - requires more sophisticated workflow")
}

// TestAccumulatorWithCondition tests accumulator with conditional logic
func TestAccumulatorWithCondition(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 5}},
			{"id": "2", "type": "accumulator", "data": {"accum_op": "sum"}},
			{"id": "3", "data": {"value": 10}},
			{"id": "4", "type": "accumulator", "data": {"accum_op": "sum"}},
			{"id": "5", "type": "extract", "data": {"field": "value"}},
			{"id": "6", "type": "condition", "data": {"condition": ">=15"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"},
			{"id": "e5", "source": "5", "target": "6"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	condResult, ok := result.NodeResults["6"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result for condition, got %T", result.NodeResults["6"])
	}

	// Accumulated sum is 5 + 10 = 15, condition is >= 15, should be true
	if condResult["condition_met"] != true {
		t.Errorf("Expected condition_met to be true, got %v", condResult["condition_met"])
	}

	// The value will be the extracted result which is a map
	valueMap, ok := condResult["value"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected value to be a map, got %T: %v", condResult["value"], condResult["value"])
	} else if valueMap["value"] != 15.0 {
		t.Errorf("Expected extracted value 15, got %v", valueMap["value"])
	}
}
