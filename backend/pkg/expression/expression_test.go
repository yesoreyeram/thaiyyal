package expression

import (
	"math"
	"testing"
)

func TestEvaluate_SimpleComparisons(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		want       bool
		wantErr    bool
	}{
		{"greater than true", ">100", 150.0, true, false},
		{"greater than false", ">100", 50.0, false, false},
		{"less than true", "<100", 50.0, true, false},
		{"less than false", "<100", 150.0, false, false},
		{"equal true", "==100", 100.0, true, false},
		{"equal false", "==100", 50.0, false, false},
		{"not equal true", "!=100", 50.0, true, false},
		{"not equal false", "!=100", 100.0, false, false},
		{"gte true", ">=100", 100.0, true, false},
		{"gte false", ">=100", 50.0, false, false},
		{"lte true", "<=100", 100.0, true, false},
		{"lte false", "<=100", 150.0, false, false},
		{"boolean true", "true", nil, true, false},
		{"boolean false", "false", nil, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, tt.input, nil)
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

func TestEvaluate_BooleanOperators(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		want       bool
	}{
		{"AND both true", "true && true", nil, true},
		{"AND one false", "true && false", nil, false},
		{"AND both false", "false && false", nil, false},
		{"OR both true", "true || true", nil, true},
		{"OR one true", "true || false", nil, true},
		{"OR both false", "false || false", nil, false},
		{"NOT true", "!true", nil, false},
		{"NOT false", "!false", nil, true},
		{"complex AND", ">100 && <200", 150.0, true},
		{"complex OR", ">100 || <50", 75.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Evaluate(tt.expression, tt.input, nil)
			if got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluate_NodeReferences(t *testing.T) {
	ctx := &Context{
		NodeResults: map[string]interface{}{
			"node1": map[string]interface{}{
				"value": 150.0,
				"output": map[string]interface{}{
					"status": 200.0,
					"data":   "success",
				},
			},
			"node2": map[string]interface{}{
				"value": 50.0,
			},
		},
		Variables:   make(map[string]interface{}),
		ContextVars: make(map[string]interface{}),
	}

	tests := []struct {
		name       string
		expression string
		want       bool
	}{
		{"node simple value", "node.node1.value > 100", true},
		{"node nested field", "node.node1.output.status == 200", true},
		{"node comparison", "node.node1.value > node.node2.value", true},
		{"node string", "node.node1.output.data == 'success'", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, nil, ctx)
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

func TestEvaluate_VariableReferences(t *testing.T) {
	ctx := &Context{
		NodeResults: make(map[string]interface{}),
		Variables: map[string]interface{}{
			"counter": 150.0,
			"enabled": true,
			"name":    "test",
		},
		ContextVars: make(map[string]interface{}),
	}

	tests := []struct {
		name       string
		expression string
		want       bool
	}{
		{"variable number", "variables.counter > 100", true},
		{"variable boolean", "variables.enabled == true", true},
		{"variable string", "variables.name == 'test'", true},
		{"variable with AND", "variables.counter > 100 && variables.enabled", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, nil, ctx)
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

func TestEvaluate_StringOperations(t *testing.T) {
	ctx := &Context{
		NodeResults: map[string]interface{}{
			"log": map[string]interface{}{
				"value": "ERROR: Connection failed",
			},
		},
		Variables:   make(map[string]interface{}),
		ContextVars: make(map[string]interface{}),
	}

	tests := []struct {
		name       string
		expression string
		want       bool
	}{
		{"contains true", "contains(node.log.value, 'ERROR')", true},
		{"contains false", "contains(node.log.value, 'SUCCESS')", false},
		{"string equality", "node.log.value == 'ERROR: Connection failed'", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, nil, ctx)
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

func TestEvaluateArithmetic_BasicOperations(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       float64
		wantErr    bool
	}{
		{"addition", "5 + 3", 8, false},
		{"subtraction", "10 - 3", 7, false},
		{"multiplication", "4 * 5", 20, false},
		{"division", "20 / 4", 5, false},
		{"modulo", "10 % 3", 1, false},
		{"negative", "-5", -5, false},
		{"positive", "+5", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateArithmetic(tt.expression, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateArithmetic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("EvaluateArithmetic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateArithmetic_NestedExpressions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       float64
	}{
		{"parentheses", "(5 + 3) * 2", 16},
		{"nested parentheses", "((5 + 3) * 2) / 4", 4},
		{"complex nested", "2 * (3 + (4 * 5))", 46},
		{"multiple operations", "10 + 5 * 2 - 3", 17}, // 10 + 10 - 3
		{"deep nesting", "(((10)))", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateArithmetic(tt.expression, nil)
			if err != nil {
				t.Errorf("EvaluateArithmetic() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateArithmetic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateArithmetic_WithVariables(t *testing.T) {
	ctx := &Context{
		NodeResults: make(map[string]interface{}),
		Variables: map[string]interface{}{
			"a": 10.0,
			"b": 5.0,
			"c": 2.0,
		},
		ContextVars: make(map[string]interface{}),
	}

	tests := []struct {
		name       string
		expression string
		want       float64
	}{
		{"variable addition", "variables.a + variables.b", 15},
		{"variable with constant", "variables.a + 5", 15},
		{"complex with variables", "variables.a + (variables.b * variables.c)", 20},
		{"nested with variables", "(variables.a + variables.b) * variables.c", 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateArithmetic(tt.expression, ctx)
			if err != nil {
				t.Errorf("EvaluateArithmetic() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateArithmetic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateArithmetic_MathFunctions(t *testing.T) {
	ctx := &Context{
		NodeResults: make(map[string]interface{}),
		Variables: map[string]interface{}{
			"foo": 5.0,
		},
		ContextVars: make(map[string]interface{}),
	}

	tests := []struct {
		name       string
		expression string
		want       float64
		tolerance  float64
	}{
		{"pow constant", "pow(2, 3)", 8, 0.001},
		{"pow variable", "pow(variables.foo, 2)", 25, 0.001},
		{"sqrt", "sqrt(16)", 4, 0.001},
		{"abs positive", "abs(5)", 5, 0.001},
		{"abs negative", "abs(-5)", 5, 0.001},
		{"floor", "floor(3.7)", 3, 0.001},
		{"ceil", "ceil(3.2)", 4, 0.001},
		{"round", "round(3.5)", 4, 0.001},
		{"min", "min(5, 3)", 3, 0.001},
		{"max", "max(5, 3)", 5, 0.001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateArithmetic(tt.expression, ctx)
			if err != nil {
				t.Errorf("EvaluateArithmetic() error = %v", err)
				return
			}
			if math.Abs(got-tt.want) > tt.tolerance {
				t.Errorf("EvaluateArithmetic() = %v, want %v (tolerance %v)", got, tt.want, tt.tolerance)
			}
		})
	}
}

func TestEvaluate_ComplexNestedConditions(t *testing.T) {
	ctx := &Context{
		NodeResults: map[string]interface{}{
			"a": map[string]interface{}{"value": 10.0},
			"b": map[string]interface{}{"value": 5.0},
		},
		Variables: map[string]interface{}{
			"foo": 3.0,
		},
		ContextVars: make(map[string]interface{}),
	}

	tests := []struct {
		name       string
		expression string
		want       bool
	}{
		{
			"nested arithmetic in condition",
			"(node.a.value + (node.b.value * 5)) > pow(variables.foo, 2)",
			true, // (10 + 25) > 9 = 35 > 9 = true
		},
		{
			"complex nested with parentheses",
			"(node.a.value + 5) > 10 && node.b.value < 10",
			true,
		},
		{
			"arithmetic with pow",
			"pow(node.a.value, 2) > 50",
			true, // 100 > 50
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, nil, ctx)
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

func TestExtractDependencies(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       []string
	}{
		{"single node", "node.http1.value > 100", []string{"http1"}},
		{"multiple nodes", "node.a.value > node.b.value", []string{"a", "b"}},
		{"with variables", "node.x.value + variables.y > 100", []string{"x"}},
		{"complex expression", "pow(node.n1.value, 2) + node.n2.value > 100", []string{"n1", "n2"}},
		{"no nodes", "variables.x > 100", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractDependencies(tt.expression)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractDependencies() = %v, want %v", got, tt.want)
				return
			}
			// Convert to map for easier comparison
			gotMap := make(map[string]bool)
			for _, id := range got {
				gotMap[id] = true
			}
			for _, id := range tt.want {
				if !gotMap[id] {
					t.Errorf("ExtractDependencies() missing %v", id)
				}
			}
		})
	}
}

func TestEvaluateArithmetic_ErrorCases(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantErr    bool
	}{
		{"division by zero", "10 / 0", true},
		{"unmatched parentheses open", "(5 + 3", true},
		{"unmatched parentheses close", "5 + 3)", true},
		{"invalid operator", "5 # 3", true},
		{"empty expression", "", true},
		{"only operator", "+", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EvaluateArithmetic(tt.expression, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateArithmetic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkEvaluate_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Evaluate(">100", 150.0, nil)
	}
}

func BenchmarkEvaluate_Complex(b *testing.B) {
	ctx := &Context{
		NodeResults: map[string]interface{}{
			"a": map[string]interface{}{"value": 10.0},
			"b": map[string]interface{}{"value": 5.0},
		},
		Variables: map[string]interface{}{
			"foo": 3.0,
		},
		ContextVars: make(map[string]interface{}),
	}
	
	for i := 0; i < b.N; i++ {
		Evaluate("(node.a.value + (node.b.value * 5)) > pow(variables.foo, 2)", nil, ctx)
	}
}

func BenchmarkEvaluateArithmetic(b *testing.B) {
	ctx := &Context{
		Variables: map[string]interface{}{
			"a": 10.0,
			"b": 5.0,
		},
	}
	
	for i := 0; i < b.N; i++ {
		EvaluateArithmetic("(variables.a + variables.b) * 2", ctx)
	}
}
