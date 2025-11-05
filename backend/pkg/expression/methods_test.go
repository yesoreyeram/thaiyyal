package expression

import (
	"testing"
)

// TestStringMethods tests string method calls
func TestStringMethods(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		ctx        *Context
		want       interface{}
		wantErr    bool
	}{
		{
			name:       "toUpperCase",
			expression: "item.name.toUpperCase()",
			input: map[string]interface{}{
				"name": "alice",
			},
			want: "ALICE",
		},
		{
			name:       "toLowerCase",
			expression: "item.name.toLowerCase()",
			input: map[string]interface{}{
				"name": "ALICE",
			},
			want: "alice",
		},
		{
			name:       "includes - found",
			expression: "item.email.includes('@example.com')",
			input: map[string]interface{}{
				"email": "user@example.com",
			},
			want: true,
		},
		{
			name:       "includes - not found",
			expression: "item.email.includes('@test.com')",
			input: map[string]interface{}{
				"email": "user@example.com",
			},
			want: false,
		},
		{
			name:       "startsWith - true",
			expression: "item.filename.startsWith('report')",
			input: map[string]interface{}{
				"filename": "report_2024.pdf",
			},
			want: true,
		},
		{
			name:       "startsWith - false",
			expression: "item.filename.startsWith('data')",
			input: map[string]interface{}{
				"filename": "report_2024.pdf",
			},
			want: false,
		},
		{
			name:       "endsWith - true",
			expression: "item.filename.endsWith('.pdf')",
			input: map[string]interface{}{
				"filename": "document.pdf",
			},
			want: true,
		},
		{
			name:       "endsWith - false",
			expression: "item.filename.endsWith('.doc')",
			input: map[string]interface{}{
				"filename": "document.pdf",
			},
			want: false,
		},
		{
			name:       "trim",
			expression: "item.text.trim()",
			input: map[string]interface{}{
				"text": "  hello world  ",
			},
			want: "hello world",
		},
		{
			name:       "replace",
			expression: "item.text.replace('old', 'new')",
			input: map[string]interface{}{
				"text": "old text with old values",
			},
			want: "new text with new values",
		},
		{
			name:       "split",
			expression: "item.text.split(',')",
			input: map[string]interface{}{
				"text": "a,b,c,d",
			},
			want: []interface{}{"a", "b", "c", "d"},
		},
		{
			name:       "chained methods",
			expression: "item.name.trim().toUpperCase()",
			input: map[string]interface{}{
				"name": "  alice  ",
			},
			want: "ALICE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expression, tt.input, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// For array comparison
				if gotArr, ok := got.([]interface{}); ok {
					wantArr, ok := tt.want.([]interface{})
					if !ok {
						t.Errorf("EvaluateExpression() got array but want is not array: %T", tt.want)
						return
					}
					if len(gotArr) != len(wantArr) {
						t.Errorf("EvaluateExpression() array length = %d, want %d", len(gotArr), len(wantArr))
						return
					}
					for i := range gotArr {
						if gotArr[i] != wantArr[i] {
							t.Errorf("EvaluateExpression() array[%d] = %v, want %v", i, gotArr[i], wantArr[i])
						}
					}
					return
				}
				
				if got != tt.want {
					t.Errorf("EvaluateExpression() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// TestStringMethodsInConditions tests string methods in boolean expressions
func TestStringMethodsInConditions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		want       bool
	}{
		{
			name:       "toUpperCase in comparison",
			expression: "item.role.toUpperCase() == 'ADMIN'",
			input: map[string]interface{}{
				"role": "admin",
			},
			want: true,
		},
		{
			name:       "toLowerCase with includes",
			expression: "item.email.toLowerCase().includes('@example.com')",
			input: map[string]interface{}{
				"email": "USER@EXAMPLE.COM",
			},
			want: true,
		},
		{
			name:       "startsWith check",
			expression: "item.status.startsWith('active')",
			input: map[string]interface{}{
				"status": "active_pending",
			},
			want: true,
		},
		{
			name:       "endsWith check",
			expression: "item.filename.endsWith('.pdf') && item.filename.startsWith('report')",
			input: map[string]interface{}{
				"filename": "report_2024.pdf",
			},
			want: true,
		},
		{
			name:       "includes with negation",
			expression: "!item.tags.includes('deprecated')",
			input: map[string]interface{}{
				"tags": "active,new,featured",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, tt.input, nil)
			if err != nil {
				t.Errorf("Evaluate() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestArrayMethods tests array method calls
func TestArrayMethods(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		ctx        *Context
		want       interface{}
		wantErr    bool
	}{
		{
			name:       "array includes - found",
			expression: "item.tags.includes('admin')",
			input: map[string]interface{}{
				"tags": []interface{}{"user", "admin", "verified"},
			},
			want: true,
		},
		{
			name:       "array includes - not found",
			expression: "item.tags.includes('superuser')",
			input: map[string]interface{}{
				"tags": []interface{}{"user", "admin", "verified"},
			},
			want: false,
		},
		{
			name:       "join",
			expression: "item.tags.join(', ')",
			input: map[string]interface{}{
				"tags": []interface{}{"red", "green", "blue"},
			},
			want: "red, green, blue",
		},
		{
			name:       "reverse",
			expression: "item.numbers.reverse()",
			input: map[string]interface{}{
				"numbers": []interface{}{1.0, 2.0, 3.0, 4.0},
			},
			want: []interface{}{4.0, 3.0, 2.0, 1.0},
		},
		{
			name:       "first",
			expression: "item.items.first()",
			input: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			want: "a",
		},
		{
			name:       "last",
			expression: "item.items.last()",
			input: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			want: "c",
		},
		{
			name:       "first on empty array",
			expression: "item.items.first()",
			input: map[string]interface{}{
				"items": []interface{}{},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expression, tt.input, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// For array comparison
				if gotArr, ok := got.([]interface{}); ok {
					wantArr, ok := tt.want.([]interface{})
					if !ok {
						t.Errorf("EvaluateExpression() got array but want is not array: %T", tt.want)
						return
					}
					if len(gotArr) != len(wantArr) {
						t.Errorf("EvaluateExpression() array length = %d, want %d", len(gotArr), len(wantArr))
						return
					}
					for i := range gotArr {
						if gotArr[i] != wantArr[i] {
							t.Errorf("EvaluateExpression() array[%d] = %v, want %v", i, gotArr[i], wantArr[i])
						}
					}
					return
				}
				
				if got != tt.want {
					t.Errorf("EvaluateExpression() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// TestMethodsWithVariables tests methods on variables
func TestMethodsWithVariables(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		ctx        *Context
		want       interface{}
	}{
		{
			name:       "variable string toUpperCase",
			expression: "variables.name.toUpperCase()",
			ctx: &Context{
				Variables: map[string]interface{}{
					"name": "alice",
				},
			},
			want: "ALICE",
		},
		{
			name:       "variable array includes",
			expression: "variables.items.includes('test')",
			ctx: &Context{
				Variables: map[string]interface{}{
					"items": []interface{}{"test", "prod", "dev"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expression, nil, tt.ctx)
			if err != nil {
				t.Errorf("EvaluateExpression() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
