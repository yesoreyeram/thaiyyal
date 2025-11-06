package executor

import (
	"strings"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestFormatExecutor_JSON tests JSON formatting
func TestFormatExecutor_JSON(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		prettyPrint bool
		want        string
	}{
		{
			name:  "Simple object",
			input: map[string]interface{}{"name": "John", "age": float64(30)},
			want:  `{"age":30,"name":"John"}`,
		},
		{
			name:  "Array",
			input: []interface{}{1.0, 2.0, 3.0},
			want:  `[1,2,3]`,
		},
		{
			name:  "String",
			input: "hello world",
			want:  `"hello world"`,
		},
		{
			name:  "Number",
			input: float64(42.5),
			want:  `42.5`,
		},
		{
			name:  "Boolean",
			input: true,
			want:  `true`,
		},
		{
			name:  "Null",
			input: nil,
			want:  `null`,
		},
		{
			name:        "Pretty print object",
			input:       map[string]interface{}{"name": "Alice", "active": true},
			prettyPrint: true,
			want: `{
  "active": true,
  "name": "Alice"
}`,
		},
		{
			name:        "Pretty print array",
			input:       []interface{}{map[string]interface{}{"id": 1.0}, map[string]interface{}{"id": 2.0}},
			prettyPrint: true,
			want: `[
  {
    "id": 1
  },
  {
    "id": 2
  }
]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FormatExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			outputType := "JSON"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeFormat,
				Data: types.FormatData{
					OutputType:  &outputType,
					PrettyPrint: &tt.prettyPrint,
				},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Errorf("Execute() error = %v", err)
				return
			}

			got, ok := result.(string)
			if !ok {
				t.Errorf("Execute() result type = %T, want string", result)
				return
			}

			if got != tt.want {
				t.Errorf("Execute() result = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestFormatExecutor_CSV tests CSV formatting
func TestFormatExecutor_CSV(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		includeHeaders bool
		want           string
	}{
		{
			name: "Simple CSV with headers",
			input: []interface{}{
				map[string]interface{}{"name": "Alice", "age": float64(30)},
				map[string]interface{}{"name": "Bob", "age": float64(25)},
			},
			includeHeaders: true,
			want:           "age,name\n30,Alice\n25,Bob\n",
		},
		{
			name: "CSV without headers",
			input: []interface{}{
				map[string]interface{}{"name": "Alice", "score": float64(95)},
			},
			includeHeaders: false,
			want:           "Alice,95\n",
		},
		{
			name: "CSV with mixed types",
			input: []interface{}{
				map[string]interface{}{"id": float64(1), "active": true, "name": "Test"},
				map[string]interface{}{"id": float64(2), "active": false, "name": "Demo"},
			},
			includeHeaders: true,
			want:           "active,id,name\ntrue,1,Test\nfalse,2,Demo\n",
		},
		{
			name: "CSV with null values",
			input: []interface{}{
				map[string]interface{}{"a": "x", "b": nil},
				map[string]interface{}{"a": "y", "b": "z"},
			},
			includeHeaders: true,
			want:           "a,b\nx,\ny,z\n",
		},
		{
			name: "Single object",
			input: map[string]interface{}{
				"name":  "Single",
				"value": float64(100),
			},
			includeHeaders: true,
			want:           "name,value\nSingle,100\n",
		},
		{
			name:           "Empty array",
			input:          []interface{}{},
			includeHeaders: true,
			want:           "",
		},
		{
			name: "CSV with numbers",
			input: []interface{}{
				map[string]interface{}{"price": float64(10.5), "quantity": float64(2)},
				map[string]interface{}{"price": float64(20), "quantity": float64(1)},
			},
			includeHeaders: true,
			want:           "price,quantity\n10.5,2\n20,1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FormatExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			outputType := "CSV"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeFormat,
				Data: types.FormatData{
					OutputType:     &outputType,
					IncludeHeaders: &tt.includeHeaders,
				},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Errorf("Execute() error = %v", err)
				return
			}

			got, ok := result.(string)
			if !ok {
				t.Errorf("Execute() result type = %T, want string", result)
				return
			}

			if got != tt.want {
				t.Errorf("Execute() result = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestFormatExecutor_TSV tests TSV formatting
func TestFormatExecutor_TSV(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		includeHeaders bool
		want           string
	}{
		{
			name: "Simple TSV with headers",
			input: []interface{}{
				map[string]interface{}{"name": "Alice", "score": float64(95)},
				map[string]interface{}{"name": "Bob", "score": float64(88)},
			},
			includeHeaders: true,
			want:           "name\tscore\nAlice\t95\nBob\t88\n",
		},
		{
			name: "TSV without headers",
			input: []interface{}{
				map[string]interface{}{"col1": "a", "col2": "b"},
			},
			includeHeaders: false,
			want:           "a\tb\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FormatExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			outputType := "TSV"
			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeFormat,
				Data: types.FormatData{
					OutputType:     &outputType,
					IncludeHeaders: &tt.includeHeaders,
				},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Errorf("Execute() error = %v", err)
				return
			}

			got, ok := result.(string)
			if !ok {
				t.Errorf("Execute() result type = %T, want string", result)
				return
			}

			if got != tt.want {
				t.Errorf("Execute() result = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestFormatExecutor_CustomDelimiter tests CSV with custom delimiter
func TestFormatExecutor_CustomDelimiter(t *testing.T) {
	exec := &FormatExecutor{}
	input := []interface{}{
		map[string]interface{}{"a": "1", "b": "2"},
		map[string]interface{}{"a": "3", "b": "4"},
	}

	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {input},
		},
	}

	outputType := "CSV"
	delimiter := "|"
	includeHeaders := true

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeFormat,
		Data: types.FormatData{
			OutputType:     &outputType,
			Delimiter:      &delimiter,
			IncludeHeaders: &includeHeaders,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got, ok := result.(string)
	if !ok {
		t.Fatalf("Execute() result type = %T, want string", result)
	}

	want := "a|b\n1|2\n3|4\n"
	if got != want {
		t.Errorf("Execute() result = %q, want %q", got, want)
	}
}

// TestFormatExecutor_Errors tests error cases
func TestFormatExecutor_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		outType string
		wantErr string
	}{
		{
			name:    "Invalid output type",
			input:   map[string]interface{}{"test": "data"},
			outType: "INVALID",
			wantErr: "unsupported output type",
		},
		{
			name:    "CSV with non-object array",
			input:   []interface{}{"string1", "string2"},
			outType: "CSV",
			wantErr: "CSV formatting requires array of objects",
		},
		{
			name:    "CSV with non-array non-object",
			input:   "just a string",
			outType: "CSV",
			wantErr: "CSV formatting requires array of objects or single object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FormatExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeFormat,
				Data: types.FormatData{
					OutputType: &tt.outType,
				},
			}

			_, err := exec.Execute(ctx, node)
			if err == nil {
				t.Errorf("Execute() expected error containing %q, got nil", tt.wantErr)
				return
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Execute() error = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}

// TestFormatExecutor_Validate tests validation
func TestFormatExecutor_Validate(t *testing.T) {
	tests := []struct {
		name     string
		outType  string
		wantErr  bool
		errMatch string
	}{
		{
			name:    "Valid JSON",
			outType: "JSON",
			wantErr: false,
		},
		{
			name:    "Valid CSV",
			outType: "CSV",
			wantErr: false,
		},
		{
			name:    "Valid TSV",
			outType: "TSV",
			wantErr: false,
		},
		{
			name:     "Invalid type",
			outType:  "INVALID",
			wantErr:  true,
			errMatch: "invalid output_type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FormatExecutor{}
			node := types.Node{
				Type: types.NodeTypeFormat,
				Data: types.FormatData{
					OutputType: &tt.outType,
				},
			}

			err := exec.Validate(node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMatch) {
				t.Errorf("Validate() error = %v, want error containing %q", err, tt.errMatch)
			}
		})
	}
}

// TestFormatExecutor_NoInput tests error when no input provided
func TestFormatExecutor_NoInput(t *testing.T) {
	exec := &FormatExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {},
		},
	}

	outputType := "JSON"
	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeFormat,
		Data: types.FormatData{
			OutputType: &outputType,
		},
	}

	_, err := exec.Execute(ctx, node)
	if err == nil {
		t.Error("Execute() expected error for no input, got nil")
	}

	if !strings.Contains(err.Error(), "requires input") {
		t.Errorf("Execute() error = %v, want error containing 'requires input'", err)
	}
}

// TestFormatExecutor_DefaultValues tests default configuration values
func TestFormatExecutor_DefaultValues(t *testing.T) {
	exec := &FormatExecutor{}
	input := map[string]interface{}{"test": "value"}

	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {input},
		},
	}

	// No OutputType specified - should default to JSON
	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeFormat,
		Data: types.FormatData{},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got, ok := result.(string)
	if !ok {
		t.Fatalf("Execute() result type = %T, want string", result)
	}

	// Should produce JSON by default
	if !strings.Contains(got, `"test"`) || !strings.Contains(got, `"value"`) {
		t.Errorf("Execute() with default output type didn't produce JSON: %q", got)
	}
}

// TestRoundTrip_ParseAndFormat tests that parse and format are complementary
func TestRoundTrip_ParseAndFormat(t *testing.T) {
	tests := []struct {
		name       string
		data       []interface{}
		formatType string
		parseType  string
	}{
		{
			name: "CSV round trip",
			data: []interface{}{
				map[string]interface{}{"name": "Alice", "age": float64(30)},
				map[string]interface{}{"name": "Bob", "age": float64(25)},
			},
			formatType: "CSV",
			parseType:  "CSV",
		},
		{
			name: "JSON round trip",
			data: []interface{}{
				map[string]interface{}{"id": float64(1), "active": true},
				map[string]interface{}{"id": float64(2), "active": false},
			},
			formatType: "JSON",
			parseType:  "JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Format to string
			formatExec := &FormatExecutor{}
			formatCtx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"format-node": {tt.data},
				},
			}

			formatNode := types.Node{
				ID:   "format-node",
				Type: types.NodeTypeFormat,
				Data: types.FormatData{
					OutputType: &tt.formatType,
				},
			}

			formatted, err := formatExec.Execute(formatCtx, formatNode)
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			formattedStr, ok := formatted.(string)
			if !ok {
				t.Fatalf("Format() result type = %T, want string", formatted)
			}

			// Parse back to structured data
			parseExec := &ParseExecutor{}
			parseCtx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"parse-node": {formattedStr},
				},
			}

			parseNode := types.Node{
				ID:   "parse-node",
				Type: types.NodeTypeParse,
				Data: types.ParseData{
					InputType: &tt.parseType,
				},
			}

			parsed, err := parseExec.Execute(parseCtx, parseNode)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			// For CSV, the parsed result will be an array of maps
			// For JSON, it should match closely (though number types might differ)
			if tt.formatType == "CSV" {
				// CSV parser may return []map[string]interface{} or []interface{}
				switch parsedData := parsed.(type) {
				case []interface{}:
					if len(parsedData) != len(tt.data) {
						t.Errorf("Parse() array length = %d, want %d", len(parsedData), len(tt.data))
					}
				case []map[string]interface{}:
					if len(parsedData) != len(tt.data) {
						t.Errorf("Parse() array length = %d, want %d", len(parsedData), len(tt.data))
					}
				default:
					t.Fatalf("Parse() result type = %T, want []interface{} or []map[string]interface{}", parsed)
				}
			}

			// We can't do exact DeepEqual because CSV loses type information
			// but we can verify structure
			t.Logf("Round trip successful: %v -> %s -> %v", tt.data, formattedStr, parsed)
		})
	}
}
