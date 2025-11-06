package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestSchemaValidatorExecutor_Type(t *testing.T) {
	executor := &SchemaValidatorExecutor{}
	if executor.NodeType() != types.NodeTypeSchemaValidator {
		t.Errorf("Expected type %v, got %v", types.NodeTypeSchemaValidator, executor.NodeType())
	}
}

func TestSchemaValidatorExecutor_Validate(t *testing.T) {
	tests := []struct {
		name    string
		schema  interface{}
		strict  *bool
		wantErr bool
	}{
		{
			name: "Valid schema",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{"type": "string"},
				},
			},
			wantErr: false,
		},
		{
			name: "Schema with strict mode",
			schema: map[string]interface{}{
				"type": "object",
			},
			strict:  boolPtr(true),
			wantErr: false,
		},
		{
			name:    "Missing schema",
			schema:  nil,
			wantErr: true,
		},
		{
			name:    "Schema not an object",
			schema:  "invalid",
			wantErr: true,
		},
		{
			name: "Schema missing type",
			schema: map[string]interface{}{
				"properties": map[string]interface{}{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &SchemaValidatorExecutor{}
			node := types.Node{
				Data: types.SchemaValidatorData{
					Schema: tt.schema,
					Strict: tt.strict,
				},
			}
			err := executor.Validate(node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSchemaValidatorExecutor_Execute_BasicValidation(t *testing.T) {
	tests := []struct {
		name      string
		schema    map[string]interface{}
		strict    *bool
		input     interface{}
		wantValid bool
		wantErr   bool
	}{
		{
			name: "Valid object",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{"type": "string"},
					"age":  map[string]interface{}{"type": "number"},
				},
				"required": []interface{}{"name"},
			},
			input: map[string]interface{}{
				"name": "John",
				"age":  float64(30),
			},
			wantValid: true,
			wantErr:   false,
		},
		{
			name: "Missing required field (lenient mode)",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{"type": "string"},
				},
				"required": []interface{}{"name"},
			},
			strict: boolPtr(false),
			input: map[string]interface{}{
				"age": float64(30),
			},
			wantValid: false,
			wantErr:   false,
		},
		{
			name: "Missing required field (strict mode)",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{"type": "string"},
				},
				"required": []interface{}{"name"},
			},
			strict: boolPtr(true),
			input: map[string]interface{}{
				"age": float64(30),
			},
			wantValid: false,
			wantErr:   true,
		},
		{
			name: "Invalid type",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"age": map[string]interface{}{
						"type": "number",
					},
				},
			},
			input: map[string]interface{}{
				"age": "not a number",
			},
			wantValid: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &SchemaValidatorExecutor{}
			node := types.Node{
				ID:   "validator1",
				Type: types.NodeTypeSchemaValidator,
				Data: types.SchemaValidatorData{
					Schema: tt.schema,
					Strict: tt.strict,
				},
			}

			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"validator1": {tt.input},
				},
			}

			result, err := executor.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Expected result to be map, got %T", result)
					return
				}

				valid, ok := resultMap["valid"].(bool)
				if !ok {
					t.Errorf("Expected valid field to be bool")
					return
				}

				if valid != tt.wantValid {
					t.Errorf("Expected valid=%v, got %v", tt.wantValid, valid)
				}

				// In invalid cases, check for errors field
				if !tt.wantValid && !tt.wantErr {
					if _, ok := resultMap["errors"]; !ok {
						t.Errorf("Expected errors field in invalid result")
					}
				}
			}
		})
	}
}

func TestSchemaValidatorExecutor_Execute_ComplexSchema(t *testing.T) {
	// User profile validation schema
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":      "string",
				"minLength": float64(1),
			},
			"age": map[string]interface{}{
				"type":    "number",
				"minimum": float64(0),
				"maximum": float64(150),
			},
			"email": map[string]interface{}{
				"type":    "string",
				"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
			},
		},
		"required": []interface{}{"name", "email"},
	}

	tests := []struct {
		name      string
		input     interface{}
		wantValid bool
	}{
		{
			name: "Valid complete user profile",
			input: map[string]interface{}{
				"name":  "John Doe",
				"age":   float64(30),
				"email": "john@example.com",
			},
			wantValid: true,
		},
		{
			name: "Valid minimal user profile",
			input: map[string]interface{}{
				"name":  "Jane Doe",
				"email": "jane@example.com",
			},
			wantValid: true,
		},
		{
			name: "Missing required name",
			input: map[string]interface{}{
				"email": "test@example.com",
			},
			wantValid: false,
		},
		{
			name: "Invalid email format",
			input: map[string]interface{}{
				"name":  "Test User",
				"email": "invalid-email",
			},
			wantValid: false,
		},
		{
			name: "Age out of range",
			input: map[string]interface{}{
				"name":  "Test User",
				"email": "test@example.com",
				"age":   float64(200),
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &SchemaValidatorExecutor{}
			node := types.Node{
				ID:   "validator1",
				Type: types.NodeTypeSchemaValidator,
				Data: types.SchemaValidatorData{
					Schema: schema,
				},
			}

			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"validator1": {tt.input},
				},
			}

			result, err := executor.Execute(ctx, node)
			if err != nil {
				t.Fatalf("Execute() unexpected error = %v", err)
			}

			resultMap := result.(map[string]interface{})
			valid := resultMap["valid"].(bool)
			if valid != tt.wantValid {
				t.Errorf("Expected valid=%v, got %v", tt.wantValid, valid)
				if !valid {
					t.Logf("Validation errors: %v", resultMap["errors"])
				}
			}

			// Verify data is preserved
			if data, ok := resultMap["data"]; !ok {
				t.Errorf("Expected data field in result")
			} else if data == nil {
				t.Errorf("Expected data to be non-nil")
			}
		})
	}
}
