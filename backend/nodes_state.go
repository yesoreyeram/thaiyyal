package workflow

import "fmt"

// ============================================================================
// State & Memory Node Executors
// ============================================================================
// This file contains executors for state management and data transformation nodes:
// - Variable: Store and retrieve values across the workflow
// - Extract: Extract fields from objects
// - Transform: Transform data structures (to_array, to_object, flatten, keys, values)
// - Accumulator: Accumulate values over time (sum, product, concat, array, count)
// - Counter: Simple counter with increment/decrement/reset/get operations
// ============================================================================

// executeVariableNode handles variable get/set operations for workflow state.
// Variables are scoped to a single workflow execution and shared across nodes.
//
// Required fields:
//   - Data.VarName: Variable name
//   - Data.VarOp: Operation ("get" or "set")
//
// For "set" operation:
//   - Requires one input value to store
//
// For "get" operation:
//   - Retrieves previously stored value
//
// Returns:
//   - map: Contains var_name, operation, and value
//   - error: If required fields missing, input missing (for set), or variable not found (for get)
func (e *Engine) executeVariableNode(node Node) (interface{}, error) {
	if node.Data.VarName == nil {
		return nil, fmt.Errorf("variable node missing var_name")
	}
	if node.Data.VarOp == nil {
		return nil, fmt.Errorf("variable node missing var_op (get or set)")
	}

	varName := *node.Data.VarName
	varOp := *node.Data.VarOp

	switch varOp {
	case "set":
		// Store value in workflow state
		inputs := e.getNodeInputs(node.ID)
		if len(inputs) == 0 {
			return nil, fmt.Errorf("variable set operation requires input value")
		}
		value := inputs[0]
		e.variables[varName] = value
		return map[string]interface{}{
			"var_name":  varName,
			"operation": "set",
			"value":     value,
		}, nil

	case "get":
		// Retrieve value from workflow state
		value, exists := e.variables[varName]
		if !exists {
			return nil, fmt.Errorf("variable '%s' not found", varName)
		}
		return map[string]interface{}{
			"var_name":  varName,
			"operation": "get",
			"value":     value,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported variable operation: %s (use 'get' or 'set')", varOp)
	}
}

// executeExtractNode extracts specific fields from object inputs.
//
// Required fields (one of):
//   - Data.Field: Single field to extract
//   - Data.Fields: Multiple fields to extract
//
// Inputs:
//   - Requires one object (map) input
//
// Returns:
//   - For single field: map with field name and value
//   - For multiple fields: map with all extracted fields
//   - error: If input is not an object or fields don't exist
func (e *Engine) executeExtractNode(node Node) (interface{}, error) {
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("extract node requires input")
	}

	input := inputs[0]
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("extract node requires object input, got %T", input)
	}

	// Single field extraction
	if node.Data.Field != nil {
		field := *node.Data.Field
		value, exists := inputMap[field]
		if !exists {
			return nil, fmt.Errorf("field '%s' not found in input object", field)
		}
		return map[string]interface{}{
			"field": field,
			"value": value,
		}, nil
	}

	// Multiple fields extraction
	if len(node.Data.Fields) > 0 {
		result := make(map[string]interface{})
		for _, field := range node.Data.Fields {
			value, exists := inputMap[field]
			if exists {
				result[field] = value
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("extract node requires 'field' or 'fields' configuration")
}

// executeTransformNode transforms data structures between different formats.
//
// Required fields:
//   - Data.TransformType: Type of transformation
//
// Supported transformations:
//   - "to_array": Convert inputs to array
//   - "to_object": Convert array of key-value pairs to object
//   - "flatten": Recursively flatten nested arrays
//   - "keys": Extract keys from object
//   - "values": Extract values from object
//
// Returns:
//   - interface{}: Transformed data in target format
//   - error: If transform type unsupported or input incompatible
func (e *Engine) executeTransformNode(node Node) (interface{}, error) {
	if node.Data.TransformType == nil {
		return nil, fmt.Errorf("transform node missing transform_type")
	}

	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("transform node requires input")
	}

	transformType := *node.Data.TransformType

	switch transformType {
	case "to_array":
		// Convert all inputs to array
		return inputs, nil

	case "to_object":
		// Convert array of alternating keys and values to object
		return e.transformToObject(inputs[0])

	case "flatten":
		// Recursively flatten nested arrays
		return e.transformFlatten(inputs[0])

	case "keys":
		// Extract all keys from object
		return e.transformKeys(inputs[0])

	case "values":
		// Extract all values from object
		return e.transformValues(inputs[0])

	default:
		return nil, fmt.Errorf("unsupported transform type: %s", transformType)
	}
}

// transformToObject converts an array to an object.
// Array should contain alternating keys (strings) and values.
func (e *Engine) transformToObject(input interface{}) (map[string]interface{}, error) {
	arr, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("to_object requires array input, got %T", input)
	}

	result := make(map[string]interface{})
	for i := 0; i < len(arr)-1; i += 2 {
		key, ok := arr[i].(string)
		if !ok {
			return nil, fmt.Errorf("to_object requires string keys at index %d", i)
		}
		result[key] = arr[i+1]
	}
	return result, nil
}

// transformFlatten recursively flattens nested arrays into a single-level array.
func (e *Engine) transformFlatten(input interface{}) ([]interface{}, error) {
	arr, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("flatten requires array input, got %T", input)
	}

	var flattened []interface{}
	var flatten func(interface{})
	flatten = func(item interface{}) {
		if subArr, ok := item.([]interface{}); ok {
			// Recursively flatten nested arrays
			for _, sub := range subArr {
				flatten(sub)
			}
		} else {
			// Append non-array items directly
			flattened = append(flattened, item)
		}
	}

	for _, item := range arr {
		flatten(item)
	}
	return flattened, nil
}

// transformKeys extracts all keys from an object.
func (e *Engine) transformKeys(input interface{}) ([]interface{}, error) {
	obj, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("keys transform requires object input, got %T", input)
	}

	keys := make([]interface{}, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	return keys, nil
}

// transformValues extracts all values from an object.
func (e *Engine) transformValues(input interface{}) ([]interface{}, error) {
	obj, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("values transform requires object input, got %T", input)
	}

	values := make([]interface{}, 0, len(obj))
	for _, v := range obj {
		values = append(values, v)
	}
	return values, nil
}

// executeAccumulatorNode accumulates values over successive calls.
// The accumulator maintains state across multiple node executions in a workflow.
//
// Required fields:
//   - Data.AccumOp: Operation type (sum, product, concat, array, count)
//
// Optional fields:
//   - Data.InitialValue: Starting value (defaults vary by operation)
//
// Inputs:
//   - Optional: If provided, adds input to accumulator
//   - If no input, returns current accumulator value
//
// Returns:
//   - map: Contains operation type and accumulated value
//   - error: If operation unsupported or input type incompatible
func (e *Engine) executeAccumulatorNode(node Node) (interface{}, error) {
	if node.Data.AccumOp == nil {
		return nil, fmt.Errorf("accumulator node missing accum_op")
	}

	accumOp := *node.Data.AccumOp
	inputs := e.getNodeInputs(node.ID)

	// Initialize accumulator with appropriate default or configured initial value
	if e.accumulator == nil {
		e.accumulator = e.getAccumulatorInitialValue(accumOp, node.Data.InitialValue)
	}

	// If no inputs, return current accumulator state
	if len(inputs) == 0 {
		return map[string]interface{}{
			"operation": accumOp,
			"value":     e.accumulator,
		}, nil
	}

	// Accumulate the input value
	input := inputs[0]
	if err := e.accumulateValue(accumOp, input); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"operation": accumOp,
		"value":     e.accumulator,
	}, nil
}

// getAccumulatorInitialValue returns the appropriate initial value for an accumulator.
func (e *Engine) getAccumulatorInitialValue(accumOp string, configuredValue interface{}) interface{} {
	if configuredValue != nil {
		return configuredValue
	}

	// Default initial values based on operation
	switch accumOp {
	case "sum", "count":
		return 0.0
	case "product":
		return 1.0
	case "concat":
		return ""
	case "array":
		return []interface{}{}
	default:
		return nil
	}
}

// accumulateValue adds an input value to the accumulator based on the operation.
func (e *Engine) accumulateValue(accumOp string, input interface{}) error {
	switch accumOp {
	case "sum":
		return e.accumulateSum(input)
	case "product":
		return e.accumulateProduct(input)
	case "concat":
		return e.accumulateConcat(input)
	case "array":
		return e.accumulateArray(input)
	case "count":
		return e.accumulateCount()
	default:
		return fmt.Errorf("unsupported accumulator operation: %s", accumOp)
	}
}

// accumulateSum adds a numeric value to the accumulator.
func (e *Engine) accumulateSum(input interface{}) error {
	accum, ok := e.accumulator.(float64)
	if !ok {
		return fmt.Errorf("accumulator value is not a number")
	}
	num, ok := input.(float64)
	if !ok {
		return fmt.Errorf("sum accumulator requires numeric input, got %T", input)
	}
	e.accumulator = accum + num
	return nil
}

// accumulateProduct multiplies the accumulator by a numeric value.
func (e *Engine) accumulateProduct(input interface{}) error {
	accum, ok := e.accumulator.(float64)
	if !ok {
		return fmt.Errorf("accumulator value is not a number")
	}
	num, ok := input.(float64)
	if !ok {
		return fmt.Errorf("product accumulator requires numeric input, got %T", input)
	}
	e.accumulator = accum * num
	return nil
}

// accumulateConcat concatenates a string to the accumulator.
func (e *Engine) accumulateConcat(input interface{}) error {
	accum, ok := e.accumulator.(string)
	if !ok {
		return fmt.Errorf("accumulator value is not a string")
	}
	str, ok := input.(string)
	if !ok {
		return fmt.Errorf("concat accumulator requires string input, got %T", input)
	}
	e.accumulator = accum + str
	return nil
}

// accumulateArray appends a value to the accumulator array.
func (e *Engine) accumulateArray(input interface{}) error {
	accum, ok := e.accumulator.([]interface{})
	if !ok {
		return fmt.Errorf("accumulator value is not an array")
	}
	e.accumulator = append(accum, input)
	return nil
}

// accumulateCount increments the counter.
func (e *Engine) accumulateCount() error {
	accum, ok := e.accumulator.(float64)
	if !ok {
		return fmt.Errorf("accumulator value is not a number")
	}
	e.accumulator = accum + 1
	return nil
}

// executeCounterNode handles counter operations (increment, decrement, reset, get).
// The counter maintains a single numeric value across workflow execution.
//
// Required fields:
//   - Data.CounterOp: Operation type (increment, decrement, reset, get)
//
// Optional fields:
//   - Data.Delta: Amount to increment/decrement (default: 1.0)
//   - Data.InitialValue: Value for reset or initialization
//
// Returns:
//   - map: Contains operation type and counter value
//   - error: If operation unsupported
func (e *Engine) executeCounterNode(node Node) (interface{}, error) {
	if node.Data.CounterOp == nil {
		return nil, fmt.Errorf("counter node missing counter_op")
	}

	counterOp := *node.Data.CounterOp

	// Initialize counter if configured
	if node.Data.InitialValue != nil {
		if val, ok := node.Data.InitialValue.(float64); ok {
			e.counter = val
		}
	}

	// Execute counter operation
	switch counterOp {
	case "increment":
		delta := 1.0
		if node.Data.Delta != nil {
			delta = *node.Data.Delta
		}
		e.counter += delta

	case "decrement":
		delta := 1.0
		if node.Data.Delta != nil {
			delta = *node.Data.Delta
		}
		e.counter -= delta

	case "reset":
		resetValue := 0.0
		if node.Data.InitialValue != nil {
			if val, ok := node.Data.InitialValue.(float64); ok {
				resetValue = val
			}
		}
		e.counter = resetValue

	case "get":
		// Just return current counter value (no modification)

	default:
		return nil, fmt.Errorf("unsupported counter operation: %s (use increment, decrement, reset, or get)", counterOp)
	}

	return map[string]interface{}{
		"operation": counterOp,
		"value":     e.counter,
	}, nil
}
