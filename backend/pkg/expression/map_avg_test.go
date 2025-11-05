package expression

import (
	"testing"
)

func TestMapFunction_ProjectAges(t *testing.T) {
	input := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "age": 30.0},
			map[string]interface{}{"name": "Bob", "age": 25.0},
		},
	}

	expr := "map(item.users, item.age)"
	got, err := EvaluateExpression(expr, input, nil)
	if err != nil {
		t.Fatalf("EvaluateExpression(map) error: %v", err)
	}

	arr, ok := got.([]interface{})
	if !ok {
		t.Fatalf("map() should return []interface{}, got %T", got)
	}
	if len(arr) != 2 {
		t.Fatalf("map() result length = %d, want 2", len(arr))
	}
	if arr[0] != 30.0 || arr[1] != 25.0 {
		t.Fatalf("map() result = %v, want [30,25]", arr)
	}
}

func TestAvgOfMap(t *testing.T) {
	input := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"age": 30.0},
			map[string]interface{}{"age": 25.0},
		},
	}

	expr := "avg(map(item.users, item.age))"
	got, err := EvaluateExpression(expr, input, nil)
	if err != nil {
		t.Fatalf("EvaluateExpression(avg(map())) error: %v", err)
	}

	if got != 27.5 {
		t.Fatalf("avg(map()) = %v, want 27.5", got)
	}
}

func TestAvgVariants(t *testing.T) {
	// avg on array variable
	ctx := &Context{
		Variables: map[string]interface{}{
			"nums": []interface{}{10.0, 20.0, 30.0},
		},
	}
	got, err := EvaluateExpression("avg(variables.nums)", nil, ctx)
	if err != nil {
		t.Fatalf("avg(array) error: %v", err)
	}
	if got != 20.0 {
		t.Fatalf("avg(array) = %v, want 20", got)
	}

	// avg with multiple args
	got, err = EvaluateExpression("avg(10, 20, 30)", nil, nil)
	if err != nil {
		t.Fatalf("avg(args) error: %v", err)
	}
	if got != 20.0 {
		t.Fatalf("avg(args) = %v, want 20", got)
	}
}

func TestRoundAvgMapUsingInputAlias(t *testing.T) {
	// Top-level input is an array, accessible via 'input'
	input := []interface{}{
		map[string]interface{}{"age": 31.0},
		map[string]interface{}{"age": 29.0},
		map[string]interface{}{"age": 40.0},
	}

	expr := "round(avg(map(input, item.age)))"
	got, err := EvaluateExpression(expr, input, nil)
	if err != nil {
		t.Fatalf("EvaluateExpression(round(avg(map(input,...)))) error: %v", err)
	}

	// Average = (31 + 29 + 40)/3 = 33.333..., round -> 33
	if got != 33.0 {
		t.Fatalf("round(avg(map(input,item.age))) = %v, want 33", got)
	}
}

func TestRoundAvgMapPlusTwo(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{"age": 31.0},
		map[string]interface{}{"age": 29.0},
		map[string]interface{}{"age": 40.0},
	}

	expr := "round(avg(map(input,item.age))) + 2"
	gotAny, err := EvaluateExpression(expr, input, nil)
	if err != nil {
		t.Fatalf("EvaluateExpression(round(avg(map(input,...))) + 2) error: %v", err)
	}

	got, ok := toFloat64(gotAny)
	if !ok {
		t.Fatalf("expected numeric result, got %T (%v)", gotAny, gotAny)
	}

	// Average = 33.333..., round -> 33, then + 2 -> 35
	if got != 35.0 {
		t.Fatalf("round(avg(map(input,item.age))) + 2 = %v, want 35", got)
	}
}

func TestEvaluateArithmeticWithEmbeddedValueFunction(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{"age": 31.0},
		map[string]interface{}{"age": 29.0},
		map[string]interface{}{"age": 40.0},
	}
	expr := "round(avg(map(input,item.age))) + 2"

	// Build a ctx with item/input variables like EvaluateExpression would
	ctx := &Context{NodeResults: map[string]interface{}{}, Variables: map[string]interface{}{}, ContextVars: map[string]interface{}{}}
	ctx.Variables["item"] = input
	ctx.Variables["input"] = input

	got, err := EvaluateArithmetic(expr, ctx)
	if err != nil {
		t.Fatalf("EvaluateArithmetic error: %v", err)
	}
	if got != 35.0 {
		t.Fatalf("EvaluateArithmetic result = %v, want 35", got)
	}
}
