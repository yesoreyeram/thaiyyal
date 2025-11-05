package expression

import (
	"testing"
)

// TestArrayLength tests the .length property on arrays
func TestArrayLength(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		ctx        *Context
		want       bool
		wantErr    bool
	}{
		{
			name:       "array length equals",
			expression: "item.tags.length == 3",
			input: map[string]interface{}{
				"tags": []interface{}{"go", "rust", "python"},
			},
			want: true,
		},
		{
			name:       "array length greater than",
			expression: "item.tags.length > 2",
			input: map[string]interface{}{
				"tags": []interface{}{"go", "rust", "python"},
			},
			want: true,
		},
		{
			name:       "variable array length",
			expression: "variables.items.length == 5",
			input:      nil,
			ctx: &Context{
				Variables: map[string]interface{}{
					"items": []interface{}{1, 2, 3, 4, 5},
				},
			},
			want: true,
		},
		{
			name:       "string length",
			expression: "item.name.length > 5",
			input: map[string]interface{}{
				"name": "Alice Smith",
			},
			want: true,
		},
		{
			name:       "empty array length",
			expression: "item.tags.length == 0",
			input: map[string]interface{}{
				"tags": []interface{}{},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, tt.input, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestArrayIndexing tests array index access like items[0]
func TestArrayIndexing(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		ctx        *Context
		want       bool
		wantErr    bool
	}{
		{
			name:       "simple array index",
			expression: "item.tags[0] == 'first'",
			input: map[string]interface{}{
				"tags": []interface{}{"first", "second", "third"},
			},
			want: true,
		},
		{
			name:       "nested object in array",
			expression: "item.users[1].name == 'Bob'",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"name": "Alice"},
					map[string]interface{}{"name": "Bob"},
					map[string]interface{}{"name": "Charlie"},
				},
			},
			want: true,
		},
		{
			name:       "variable array index",
			expression: "variables.items[0] == 'hello'",
			input:      nil,
			ctx: &Context{
				Variables: map[string]interface{}{
					"items": []interface{}{"hello", "world"},
				},
			},
			want: true,
		},
		{
			name:       "index out of bounds",
			expression: "item.tags[5] == 'test'",
			input: map[string]interface{}{
				"tags": []interface{}{"a", "b"},
			},
			want: false, // Out of bounds returns false, doesn't error
		},
		{
			name:       "array first element comparison",
			expression: "item.scores[0] > 90",
			input: map[string]interface{}{
				"scores": []interface{}{95.0, 80.0, 88.0},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, tt.input, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEvaluateExpression_ArrayLength tests EvaluateExpression with array length
func TestEvaluateExpression_ArrayLength(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		ctx        *Context
		want       interface{}
		wantErr    bool
	}{
		{
			name:       "get array length",
			expression: "item.tags.length",
			input: map[string]interface{}{
				"tags": []interface{}{"go", "rust", "python"},
			},
			want: float64(3),
		},
		{
			name:       "array length in arithmetic",
			expression: "item.tags.length * 2",
			input: map[string]interface{}{
				"tags": []interface{}{"a", "b", "c"},
			},
			want: float64(6),
		},
		{
			name:       "string length",
			expression: "item.name.length",
			input: map[string]interface{}{
				"name": "test",
			},
			want: float64(4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expression, tt.input, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("EvaluateExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEvaluateExpression_ArrayIndexing tests EvaluateExpression with array indexing
func TestEvaluateExpression_ArrayIndexing(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		ctx        *Context
		want       interface{}
		wantErr    bool
	}{
		{
			name:       "get array element",
			expression: "item.tags[0]",
			input: map[string]interface{}{
				"tags": []interface{}{"first", "second", "third"},
			},
			want: "first",
		},
		{
			name:       "get nested object field from array",
			expression: "item.users[1].name",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"name": "Alice", "age": float64(25)},
					map[string]interface{}{"name": "Bob", "age": float64(30)},
				},
			},
			want: "Bob",
		},
		{
			name:       "use array element in arithmetic",
			expression: "item.scores[0] * 2",
			input: map[string]interface{}{
				"scores": []interface{}{10.0, 20.0, 30.0},
			},
			want: float64(20),
		},
		{
			name:       "variable array indexing",
			expression: "variables.items[2]",
			input:      nil,
			ctx: &Context{
				Variables: map[string]interface{}{
					"items": []interface{}{"a", "b", "c", "d"},
				},
			},
			want: "c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expression, tt.input, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("EvaluateExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestComplexExpressions tests complex combinations of features
func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		ctx        *Context
		want       bool
		wantErr    bool
	}{
		{
			name:       "array length and comparison",
			expression: "item.tags.length > 0 && item.tags[0] == 'important'",
			input: map[string]interface{}{
				"tags": []interface{}{"important", "urgent"},
			},
			want: true,
		},
		{
			name:       "nested array indexing",
			expression: "item.users[0].tags[1] == 'verified'",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name": "Alice",
						"tags": []interface{}{"admin", "verified"},
					},
				},
			},
			want: true,
		},
		{
			name:       "arithmetic with array length",
			expression: "item.items.length * 2 > 10",
			input: map[string]interface{}{
				"items": []interface{}{1, 2, 3, 4, 5, 6},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, tt.input, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
