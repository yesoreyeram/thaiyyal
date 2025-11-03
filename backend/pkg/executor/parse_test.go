package executor

import (
	"reflect"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestParseExecutor_JSON tests JSON parsing
func TestParseExecutor_JSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "JSON object",
			input:    `{"name": "John", "age": 30, "active": true}`,
			expected: map[string]interface{}{"name": "John", "age": float64(30), "active": true},
		},
		{
			name:     "JSON array",
			input:    `[1, 2, 3, 4, 5]`,
			expected: []interface{}{float64(1), float64(2), float64(3), float64(4), float64(5)},
		},
		{
			name:     "JSON string",
			input:    `"hello world"`,
			expected: "hello world",
		},
		{
			name:     "JSON number",
			input:    `42.5`,
			expected: float64(42.5),
		},
		{
			name:     "JSON boolean true",
			input:    `true`,
			expected: true,
		},
		{
			name:     "JSON boolean false",
			input:    `false`,
			expected: false,
		},
		{
			name:     "JSON null",
			input:    `null`,
			expected: nil,
		},
		{
			name:     "Nested JSON",
			input:    `{"user": {"name": "Alice", "roles": ["admin", "user"]}}`,
			expected: map[string]interface{}{"user": map[string]interface{}{"name": "Alice", "roles": []interface{}{"admin", "user"}}},
		},
		{
			name:    "Invalid JSON",
			input:   `{invalid json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ParseExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			inputType := "JSON"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeParse,
				Data: types.NodeData{
					InputType: &inputType,
				},
			}

			result, err := exec.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Execute() result = %#v, expected %#v", result, tt.expected)
			}
		})
	}
}

// TestParseExecutor_CSV tests CSV parsing
func TestParseExecutor_CSV(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:  "Simple CSV",
			input: "name,age,city\nAlice,30,NYC\nBob,25,LA",
			expected: []map[string]interface{}{
				{"name": "Alice", "age": float64(30), "city": "NYC"},
				{"name": "Bob", "age": float64(25), "city": "LA"},
			},
		},
		{
			name:  "CSV with booleans",
			input: "name,active,score\nJohn,true,95.5\nJane,false,87.2",
			expected: []map[string]interface{}{
				{"name": "John", "active": true, "score": float64(95.5)},
				{"name": "Jane", "active": false, "score": float64(87.2)},
			},
		},
		{
			name:     "CSV with headers only",
			input:    "name,age,city",
			expected: []map[string]interface{}{},
		},
		{
			name:     "Empty CSV",
			input:    "",
			expected: []map[string]interface{}{},
		},
		{
			name:  "CSV with null values",
			input: "name,value\nAlice,null\nBob,42",
			expected: []map[string]interface{}{
				{"name": "Alice", "value": nil},
				{"name": "Bob", "value": float64(42)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ParseExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			inputType := "CSV"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeParse,
				Data: types.NodeData{
					InputType: &inputType,
				},
			}

			result, err := exec.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Execute() result = %#v, expected %#v", result, tt.expected)
			}
		})
	}
}

// TestParseExecutor_TSV tests TSV parsing
func TestParseExecutor_TSV(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:  "Simple TSV",
			input: "name\tage\tcity\nAlice\t30\tNYC\nBob\t25\tLA",
			expected: []map[string]interface{}{
				{"name": "Alice", "age": float64(30), "city": "NYC"},
				{"name": "Bob", "age": float64(25), "city": "LA"},
			},
		},
		{
			name:  "TSV with mixed types",
			input: "product\tprice\tin_stock\niPhone\t999.99\ttrue\niPad\t799.99\tfalse",
			expected: []map[string]interface{}{
				{"product": "iPhone", "price": float64(999.99), "in_stock": true},
				{"product": "iPad", "price": float64(799.99), "in_stock": false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ParseExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			inputType := "TSV"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeParse,
				Data: types.NodeData{
					InputType: &inputType,
				},
			}

			result, err := exec.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Execute() result = %#v, expected %#v", result, tt.expected)
			}
		})
	}
}

// TestParseExecutor_YAML tests YAML parsing
func TestParseExecutor_YAML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name: "Simple YAML",
			input: `name: John
age: 30
city: NYC`,
			expected: map[string]interface{}{
				"name": "John",
				"age":  float64(30),
				"city": "NYC",
			},
		},
		{
			name: "YAML with booleans",
			input: `active: true
verified: false
count: 42`,
			expected: map[string]interface{}{
				"active":   true,
				"verified": false,
				"count":    float64(42),
			},
		},
		{
			name: "YAML with null",
			input: `name: Alice
value: null
score: 95.5`,
			expected: map[string]interface{}{
				"name":  "Alice",
				"value": nil,
				"score": float64(95.5),
			},
		},
		{
			name:    "Invalid YAML",
			input:   "no colon separator here",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ParseExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			inputType := "YAML"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeParse,
				Data: types.NodeData{
					InputType: &inputType,
				},
			}

			result, err := exec.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Execute() result = %#v, expected %#v", result, tt.expected)
			}
		})
	}
}

// TestParseExecutor_XML tests XML parsing
func TestParseExecutor_XML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
	}{
		{
			name:    "Simple XML",
			input:   `<name>John</name>`,
			wantErr: false,
		},
		{
			name:    "Invalid XML",
			input:   `<unclosed>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ParseExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			inputType := "XML"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeParse,
				Data: types.NodeData{
					InputType: &inputType,
				},
			}

			_, err := exec.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParseExecutor_AUTO tests automatic format detection
func TestParseExecutor_AUTO(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType string
	}{
		{
			name:         "Detect JSON object",
			input:        `{"name": "John"}`,
			expectedType: "JSON",
		},
		{
			name:         "Detect JSON array",
			input:        `[1, 2, 3]`,
			expectedType: "JSON",
		},
		{
			name:         "Detect CSV",
			input:        "name,age\nJohn,30",
			expectedType: "CSV",
		},
		{
			name:         "Detect TSV",
			input:        "name\tage\nJohn\t30",
			expectedType: "TSV",
		},
		{
			name:         "Detect YAML",
			input:        "name: John\nage: 30",
			expectedType: "YAML",
		},
		{
			name:         "Detect XML",
			input:        "<root>test</root>",
			expectedType: "XML",
		},
		{
			name:         "Detect primitive number",
			input:        "42",
			expectedType: "JSON",
		},
		{
			name:         "Detect boolean",
			input:        "true",
			expectedType: "JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detected := detectFormat(tt.input)
			if detected != tt.expectedType {
				t.Errorf("detectFormat() = %v, expected %v", detected, tt.expectedType)
			}
		})
	}
}

// TestParseExecutor_AUTO_Execute tests AUTO mode execution
func TestParseExecutor_AUTO_Execute(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "AUTO - JSON object",
			input:   `{"name": "John", "age": 30}`,
			wantErr: false,
		},
		{
			name:    "AUTO - CSV",
			input:   "name,age\nJohn,30\nJane,25",
			wantErr: false,
		},
		{
			name:    "AUTO - TSV",
			input:   "name\tage\nJohn\t30",
			wantErr: false,
		},
		{
			name:    "AUTO - YAML",
			input:   "name: John\nage: 30",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ParseExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			inputType := "AUTO"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeParse,
				Data: types.NodeData{
					InputType: &inputType,
				},
			}

			_, err := exec.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParseExecutor_NoInput tests error handling for missing input
func TestParseExecutor_NoInput(t *testing.T) {
	exec := &ParseExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeParse,
		Data: types.NodeData{},
	}

	_, err := exec.Execute(ctx, node)
	if err == nil {
		t.Error("Execute() expected error for missing input, got nil")
	}
}

// TestParseExecutor_DefaultInputType tests default AUTO behavior
func TestParseExecutor_DefaultInputType(t *testing.T) {
	exec := &ParseExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {`{"name": "John"}`},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeParse,
		Data: types.NodeData{
			// InputType not set, should default to AUTO
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Errorf("Execute() unexpected error: %v", err)
	}

	expected := map[string]interface{}{"name": "John"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Execute() result = %#v, expected %#v", result, expected)
	}
}

// TestParseExecutor_Validate tests validation
func TestParseExecutor_Validate(t *testing.T) {
	tests := []struct {
		name      string
		inputType string
		wantErr   bool
	}{
		{"Valid AUTO", "AUTO", false},
		{"Valid JSON", "JSON", false},
		{"Valid CSV", "CSV", false},
		{"Valid TSV", "TSV", false},
		{"Valid YAML", "YAML", false},
		{"Valid XML", "XML", false},
		{"Valid lowercase", "json", false},
		{"Invalid type", "INVALID", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ParseExecutor{}
			node := types.Node{
				Type: types.NodeTypeParse,
				Data: types.NodeData{
					InputType: &tt.inputType,
				},
			}

			err := exec.Validate(node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParseExecutor_NonStringInput tests conversion of non-string inputs
func TestParseExecutor_NonStringInput(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"Number input", float64(42), false},
		{"Boolean input", true, false},
		{"Map input", map[string]interface{}{"test": "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ParseExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			inputType := "JSON"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeParse,
				Data: types.NodeData{
					InputType: &inputType,
				},
			}

			_, err := exec.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
