package workflow

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Multiple Typed Values Tests for Context Nodes
// ============================================================================

func TestContextVariable_MultipleValues(t *testing.T) {
	// Test new format with multiple typed values
	payload := `{
		"nodes": [
			{
				"id": "ctx1",
				"type": "context_variable",
				"data": {
					"context_values": [
						{"name": "username", "value": "john_doe", "type": "string"},
						{"name": "user_id", "value": 12345, "type": "number"},
						{"name": "is_admin", "value": true, "type": "boolean"}
					]
				}
			}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	// Check that variables were stored
	if engine.contextVariables["username"] != "john_doe" {
		t.Errorf("expected username='john_doe', got: %v", engine.contextVariables["username"])
	}

	if engine.contextVariables["user_id"] != float64(12345) {
		t.Errorf("expected user_id=12345, got: %v", engine.contextVariables["user_id"])
	}

	if engine.contextVariables["is_admin"] != true {
		t.Errorf("expected is_admin=true, got: %v", engine.contextVariables["is_admin"])
	}

	// Check result contains all variables
	ctx1Result := result.NodeResults["ctx1"].(map[string]interface{})
	if ctx1Result["type"] != "variable" {
		t.Errorf("expected type='variable', got: %v", ctx1Result["type"])
	}

	variables := ctx1Result["variables"].(map[string]interface{})
	if len(variables) != 3 {
		t.Errorf("expected 3 variables, got: %d", len(variables))
	}
}

func TestContextConstant_MultipleValues(t *testing.T) {
	// Test new format with multiple typed values for constants
	payload := `{
		"nodes": [
			{
				"id": "ctx1",
				"type": "context_constant",
				"data": {
					"context_values": [
						{"name": "api_url", "value": "https://api.example.com", "type": "string"},
						{"name": "max_retries", "value": 3, "type": "number"},
						{"name": "debug_mode", "value": false, "type": "boolean"}
					]
				}
			}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	// Check that constants were stored
	if engine.contextConstants["api_url"] != "https://api.example.com" {
		t.Errorf("expected api_url, got: %v", engine.contextConstants["api_url"])
	}

	if engine.contextConstants["max_retries"] != float64(3) {
		t.Errorf("expected max_retries=3, got: %v", engine.contextConstants["max_retries"])
	}

	if engine.contextConstants["debug_mode"] != false {
		t.Errorf("expected debug_mode=false, got: %v", engine.contextConstants["debug_mode"])
	}

	// Check result
	ctx1Result := result.NodeResults["ctx1"].(map[string]interface{})
	constants := ctx1Result["constants"].(map[string]interface{})
	if len(constants) != 3 {
		t.Errorf("expected 3 constants, got: %d", len(constants))
	}
}

func TestTypeConversion_String(t *testing.T) {
	tests := []struct{
		name     string
		value    interface{}
		expected string
	}{
		{"number to string", 123, "123"},
		{"bool to string", true, "true"},
		{"already string", "hello", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertTypedValue(tt.value, "string")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTypeConversion_Number(t *testing.T) {
	tests := []struct{
		name     string
		value    interface{}
		expected float64
	}{
		{"int to number", 123, 123.0},
		{"float to number", 45.67, 45.67},
		{"string to number", "89.5", 89.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertTypedValue(tt.value, "number")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTypeConversion_Boolean(t *testing.T) {
	tests := []struct{
		name     string
		value    interface{}
		expected bool
	}{
		{"true bool", true, true},
		{"false bool", false, false},
		{"string true", "true", true},
		{"string false", "false", false},
		{"number non-zero", 1.0, true},
		{"number zero", 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertTypedValue(tt.value, "boolean")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTypeConversion_TimeString(t *testing.T) {
	// Test valid time string
	timeStr := "2025-10-30T12:00:00Z"
	result, err := convertTypedValue(timeStr, "time_string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != timeStr {
		t.Errorf("expected %q, got %q", timeStr, result)
	}

	// Test invalid time string
	_, err = convertTypedValue("not-a-time", "time_string")
	if err == nil {
		t.Error("expected error for invalid time string")
	}
}

func TestTypeConversion_EpochSecond(t *testing.T) {
	// Test epoch seconds conversion
	epochSec := int64(1698667200) // 2023-10-30 12:00:00 UTC
	result, err := convertTypedValue(epochSec, "epoch_second")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	timeResult, ok := result.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time, got %T", result)
	}

	expectedTime := time.Unix(epochSec, 0)
	if !timeResult.Equal(expectedTime) {
		t.Errorf("expected %v, got %v", expectedTime, timeResult)
	}

	// Test with string
	result, err = convertTypedValue("1698667200", "epoch_second")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	timeResult = result.(time.Time)
	if !timeResult.Equal(expectedTime) {
		t.Errorf("expected %v, got %v", expectedTime, timeResult)
	}
}

func TestTypeConversion_EpochMs(t *testing.T) {
	// Test epoch milliseconds conversion
	epochMs := int64(1698667200500) // 2023-10-30 12:00:00.500 UTC
	result, err := convertTypedValue(epochMs, "epoch_ms")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	timeResult, ok := result.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time, got %T", result)
	}

	expectedTime := time.Unix(epochMs/1000, (epochMs%1000)*1000000)
	if !timeResult.Equal(expectedTime) {
		t.Errorf("expected %v, got %v", expectedTime, timeResult)
	}

	// Test with float
	result, err = convertTypedValue(float64(epochMs), "epoch_ms")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	timeResult = result.(time.Time)
	if !timeResult.Equal(expectedTime) {
		t.Errorf("expected %v, got %v", expectedTime, timeResult)
	}
}

func TestTypeConversion_Null(t *testing.T) {
	result, err := convertTypedValue("anything", "null")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestTypeConversion_InvalidType(t *testing.T) {
	_, err := convertTypedValue("value", "invalid_type")
	if err == nil {
		t.Error("expected error for invalid type")
	}
	if !strings.Contains(err.Error(), "unsupported type") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestBackwardCompatibility_SingleValue(t *testing.T) {
	// Test that old format still works
	payload := `{
		"nodes": [
			{
				"id": "ctx1",
				"type": "context_variable",
				"data": {
					"context_name": "user",
					"context_value": "Alice"
				}
			}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	if engine.contextVariables["user"] != "Alice" {
		t.Errorf("expected user='Alice', got: %v", engine.contextVariables["user"])
	}

	// Check result format
	ctx1Result := result.NodeResults["ctx1"].(map[string]interface{})
	if ctx1Result["type"] != "variable" {
		t.Errorf("expected type='variable'")
	}
	if ctx1Result["name"] != "user" {
		t.Errorf("expected name='user'")
	}
	if ctx1Result["value"] != "Alice" {
		t.Errorf("expected value='Alice'")
	}
}

func TestMixedTypes_InSingleNode(t *testing.T) {
	// Test node with various types
	payload := `{
		"nodes": [
			{
				"id": "ctx1",
				"type": "context_variable",
				"data": {
					"context_values": [
						{"name": "str_val", "value": "hello", "type": "string"},
						{"name": "num_val", "value": 42, "type": "number"},
						{"name": "bool_val", "value": true, "type": "boolean"},
						{"name": "time_val", "value": "2025-10-30T12:00:00Z", "type": "time_string"},
						{"name": "epoch_val", "value": 1698667200, "type": "epoch_second"},
						{"name": "null_val", "value": "ignored", "type": "null"}
					]
				}
			}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	// Verify all types were converted correctly
	if engine.contextVariables["str_val"] != "hello" {
		t.Error("str_val conversion failed")
	}
	if engine.contextVariables["num_val"] != float64(42) {
		t.Error("num_val conversion failed")
	}
	if engine.contextVariables["bool_val"] != true {
		t.Error("bool_val conversion failed")
	}
	if engine.contextVariables["time_val"] != "2025-10-30T12:00:00Z" {
		t.Error("time_val conversion failed")
	}
	
	epochTime, ok := engine.contextVariables["epoch_val"].(time.Time)
	if !ok {
		t.Error("epoch_val should be time.Time")
	}
	expected := time.Unix(1698667200, 0)
	if !epochTime.Equal(expected) {
		t.Errorf("epoch_val: expected %v, got %v", expected, epochTime)
	}
	
	if engine.contextVariables["null_val"] != nil {
		t.Error("null_val should be nil")
	}

	// Check result
	ctx1Result := result.NodeResults["ctx1"].(map[string]interface{})
	variables := ctx1Result["variables"].(map[string]interface{})
	if len(variables) != 6 {
		t.Errorf("expected 6 variables, got %d", len(variables))
	}
}

func TestJSONSerialization_TypedValues(t *testing.T) {
	// Test that typed values can be serialized to JSON and back
	values := []ContextVariableValue{
		{Name: "username", Value: "john", Type: "string"},
		{Name: "age", Value: 30, Type: "number"},
		{Name: "active", Value: true, Type: "boolean"},
	}

	data, err := json.Marshal(values)
	if err != nil {
		t.Fatalf("error marshaling: %v", err)
	}

	var decoded []ContextVariableValue
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("error unmarshaling: %v", err)
	}

	if len(decoded) != 3 {
		t.Errorf("expected 3 values, got %d", len(decoded))
	}

	if decoded[0].Name != "username" || decoded[0].Type != "string" {
		t.Error("username not decoded correctly")
	}
}

func TestTypeConversionError_InvalidNumber(t *testing.T) {
	_, err := convertTypedValue("not-a-number", "number")
	if err == nil {
		t.Error("expected error for invalid number")
	}
}

func TestTypeConversionError_InvalidBoolean(t *testing.T) {
	_, err := convertTypedValue("not-a-bool", "boolean")
	if err == nil {
		t.Error("expected error for invalid boolean")
	}
}

func TestTypeConversionError_InvalidEpochSecond(t *testing.T) {
	_, err := convertTypedValue("not-an-epoch", "epoch_second")
	if err == nil {
		t.Error("expected error for invalid epoch second")
	}
}

func TestTypeConversionError_InvalidEpochMs(t *testing.T) {
	_, err := convertTypedValue("not-an-epoch", "epoch_ms")
	if err == nil {
		t.Error("expected error for invalid epoch milliseconds")
	}
}
