package workflow

import (
	"testing"
)

func TestContextVariableNode(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_variable",
				"data": {
					"context_name": "username",
					"context_value": "john_doe"
				}
			}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Check that the variable was stored in context
	if val, ok := engine.contextVariables["username"]; !ok {
		t.Error("Expected 'username' to be in contextVariables")
	} else if val != "john_doe" {
		t.Errorf("Expected username value to be 'john_doe', got %v", val)
	}

	// Check the result
	resultMap, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected final output to be a map, got %T", result.FinalOutput)
	}
	if resultMap["type"] != "variable" {
		t.Errorf("Expected type to be 'variable', got %v", resultMap["type"])
	}
	if resultMap["name"] != "username" {
		t.Errorf("Expected name to be 'username', got %v", resultMap["name"])
	}
}

func TestContextConstantNode(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_constant",
				"data": {
					"context_name": "apiUrl",
					"context_value": "https://api.example.com"
				}
			}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Check that the constant was stored in context
	if val, ok := engine.contextConstants["apiUrl"]; !ok {
		t.Error("Expected 'apiUrl' to be in contextConstants")
	} else if val != "https://api.example.com" {
		t.Errorf("Expected apiUrl value to be 'https://api.example.com', got %v", val)
	}

	// Check the result
	resultMap, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected final output to be a map, got %T", result.FinalOutput)
	}
	if resultMap["type"] != "const" {
		t.Errorf("Expected type to be 'const', got %v", resultMap["type"])
	}
}

func TestTemplateInterpolation_Variable(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_variable",
				"data": {
					"context_name": "greeting",
					"context_value": "Hello"
				}
			},
			{
				"id": "2",
				"type": "text_input",
				"data": {
					"text": "{{ variable.greeting }}, World!"
				}
			},
			{
				"id": "3",
				"type": "visualization",
				"data": {
					"mode": "text"
				}
			}
		],
		"edges": [
			{"id": "e2-3", "source": "2", "target": "3"}
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

	finalOutput, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected final output to be a map, got %T", result.FinalOutput)
	}

	expectedText := "Hello, World!"
	if finalOutput["value"] != expectedText {
		t.Errorf("Expected value to be %q, got %v", expectedText, finalOutput["value"])
	}
}

func TestTemplateInterpolation_Constant(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_constant",
				"data": {
					"context_name": "baseUrl",
					"context_value": "https://api.example.com"
				}
			},
			{
				"id": "2",
				"type": "text_input",
				"data": {
					"text": "{{ const.baseUrl }}/users"
				}
			},
			{
				"id": "3",
				"type": "visualization",
				"data": {
					"mode": "text"
				}
			}
		],
		"edges": [
			{"id": "e2-3", "source": "2", "target": "3"}
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

	finalOutput, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected final output to be a map, got %T", result.FinalOutput)
	}

	expectedText := "https://api.example.com/users"
	if finalOutput["value"] != expectedText {
		t.Errorf("Expected value to be %q, got %v", expectedText, finalOutput["value"])
	}
}

func TestTemplateInterpolation_MixedTypes(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_variable",
				"data": {
					"context_name": "user",
					"context_value": "alice"
				}
			},
			{
				"id": "2",
				"type": "context_constant",
				"data": {
					"context_name": "host",
					"context_value": "api.example.com"
				}
			},
			{
				"id": "3",
				"type": "text_input",
				"data": {
					"text": "User {{ variable.user }} at {{ const.host }}"
				}
			},
			{
				"id": "4",
				"type": "visualization",
				"data": {
					"mode": "text"
				}
			}
		],
		"edges": [
			{"id": "e3-4", "source": "3", "target": "4"}
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

	finalOutput, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected final output to be a map, got %T", result.FinalOutput)
	}

	expectedText := "User alice at api.example.com"
	if finalOutput["value"] != expectedText {
		t.Errorf("Expected value to be %q, got %v", expectedText, finalOutput["value"])
	}
}

func TestTemplateInterpolation_NotFound(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "text_input",
				"data": {
					"text": "Hello {{ variable.missing }}"
				}
			},
			{
				"id": "2",
				"type": "visualization",
				"data": {
					"mode": "text"
				}
			}
		],
		"edges": [
			{"id": "e1-2", "source": "1", "target": "2"}
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

	finalOutput, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected final output to be a map, got %T", result.FinalOutput)
	}

	// Should remain unchanged when variable not found
	expectedText := "Hello {{ variable.missing }}"
	if finalOutput["value"] != expectedText {
		t.Errorf("Expected value to be %q, got %v", expectedText, finalOutput["value"])
	}
}

func TestTemplateInterpolation_VariableNode(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_constant",
				"data": {
					"context_name": "storageKey",
					"context_value": "userdata"
				}
			},
			{
				"id": "2",
				"type": "number",
				"data": {
					"value": 100
				}
			},
			{
				"id": "3",
				"type": "variable",
				"data": {
					"var_name": "{{ const.storageKey }}",
					"var_op": "set"
				}
			},
			{
				"id": "4",
				"type": "variable",
				"data": {
					"var_name": "{{ const.storageKey }}",
					"var_op": "get"
				}
			},
			{
				"id": "5",
				"type": "visualization",
				"data": {
					"mode": "text"
				}
			}
		],
		"edges": [
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-4", "source": "3", "target": "4"},
			{"id": "e4-5", "source": "4", "target": "5"}
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

	finalOutput, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected final output to be a map, got %T", result.FinalOutput)
	}

	// The value should be retrieved from the variable store
	varResult, ok := finalOutput["value"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected value to be a map, got %T", finalOutput["value"])
	}

	if varResult["value"] != float64(100) {
		t.Errorf("Expected value to be 100, got %v", varResult["value"])
	}
}

func TestTemplateInterpolation_NumericValue(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_variable",
				"data": {
					"context_name": "count",
					"context_value": 42
				}
			},
			{
				"id": "2",
				"type": "text_input",
				"data": {
					"text": "Count is {{ variable.count }}"
				}
			},
			{
				"id": "3",
				"type": "visualization",
				"data": {
					"mode": "text"
				}
			}
		],
		"edges": [
			{"id": "e2-3", "source": "2", "target": "3"}
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

	finalOutput, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected final output to be a map, got %T", result.FinalOutput)
	}

	expectedText := "Count is 42"
	if finalOutput["value"] != expectedText {
		t.Errorf("Expected value to be %q, got %v", expectedText, finalOutput["value"])
	}
}

func TestMultipleContextNodes(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_variable",
				"data": {
					"context_name": "var1",
					"context_value": "value1"
				}
			},
			{
				"id": "2",
				"type": "context_variable",
				"data": {
					"context_name": "var2",
					"context_value": "value2"
				}
			},
			{
				"id": "3",
				"type": "context_constant",
				"data": {
					"context_name": "const1",
					"context_value": "constValue"
				}
			},
			{
				"id": "4",
				"type": "text_input",
				"data": {
					"text": "{{ variable.var1 }} {{ variable.var2 }} {{ const.const1 }}"
				}
			},
			{
				"id": "5",
				"type": "visualization",
				"data": {
					"mode": "text"
				}
			}
		],
		"edges": [
			{"id": "e4-5", "source": "4", "target": "5"}
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

	finalOutput, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected final output to be a map, got %T", result.FinalOutput)
	}

	expectedText := "value1 value2 constValue"
	if finalOutput["value"] != expectedText {
		t.Errorf("Expected value to be %q, got %v", expectedText, finalOutput["value"])
	}
}

func TestContextNode_MissingName(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_variable",
				"data": {
					"context_value": "value"
				}
			}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	_, err = engine.Execute()
	if err == nil {
		t.Error("Expected error for missing context_name, got nil")
	}
}

func TestContextNode_MissingValue(t *testing.T) {
	payload := `{
		"nodes": [
			{
				"id": "1",
				"type": "context_constant",
				"data": {
					"context_name": "test"
				}
			}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	_, err = engine.Execute()
	if err == nil {
		t.Error("Expected error for missing context_value, got nil")
	}
}
