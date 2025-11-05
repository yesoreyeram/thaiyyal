package expression

import (
	"testing"
)

// Math functions with arrays
func TestMathFunctionsWithArrays(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		input   interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "round single value",
			expr: "round(3.7)",
			want: 4.0,
		},
		{
			name: "round array",
			expr: "round(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{3.2, 4.7, 5.5, 6.1},
			},
			want: []interface{}{3.0, 5.0, 6.0, 6.0},
		},
		{
			name: "floor array",
			expr: "floor(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{3.9, 4.1, 5.8},
			},
			want: []interface{}{3.0, 4.0, 5.0},
		},
		{
			name: "ceil array",
			expr: "ceil(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{3.1, 4.0, 5.9},
			},
			want: []interface{}{4.0, 4.0, 6.0},
		},
		{
			name: "abs array",
			expr: "abs(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{-3.0, 4.0, -5.5},
			},
			want: []interface{}{3.0, 4.0, 5.5},
		},
		{
			name: "min array",
			expr: "min(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{2.0, 3.0, 1.0, 4.0},
			},
			want: 1.0,
		},
		{
			name: "max array",
			expr: "max(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{2.0, 3.0, 1.0, 4.0},
			},
			want: 4.0,
		},
		{
			name: "min multiple args",
			expr: "min(5, 2, 8, 1)",
			want: 1.0,
		},
		{
			name: "max multiple args",
			expr: "max(5, 2, 8, 1)",
			want: 8.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expr, tt.input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if arr, ok := got.([]interface{}); ok {
					wantArr := tt.want.([]interface{})
					if len(arr) != len(wantArr) {
						t.Errorf("EvaluateExpression() array length = %d, want %d", len(arr), len(wantArr))
						return
					}
					for i := range arr {
						if arr[i] != wantArr[i] {
							t.Errorf("EvaluateExpression() array[%d] = %v, want %v", i, arr[i], wantArr[i])
						}
					}
				} else if got != tt.want {
					t.Errorf("EvaluateExpression() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSumFunction(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		input   interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "sum array",
			expr: "sum(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{1.0, 2.0, 3.0, 4.0},
			},
			want: 10.0,
		},
		{
			name: "sum multiple args",
			expr: "sum(1, 2, 3, 4)",
			want: 10.0,
		},
		{
			name: "sum empty array",
			expr: "sum(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{},
			},
			want: 0.0,
		},
		{
			name: "sum with map",
			expr: "sum(map(item.users, item.age))",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"age": 30.0},
					map[string]interface{}{"age": 25.0},
					map[string]interface{}{"age": 35.0},
				},
			},
			want: 90.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expr, tt.input, nil)
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

func TestArraySortFunction(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		input   interface{}
		want    []interface{}
		wantErr bool
	}{
		{
			name: "sort numbers",
			expr: "sort(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{3.0, 1.0, 4.0, 2.0},
			},
			want: []interface{}{1.0, 2.0, 3.0, 4.0},
		},
		{
			name: "sort strings",
			expr: "sort(item.names)",
			input: map[string]interface{}{
				"names": []interface{}{"charlie", "alice", "bob"},
			},
			want: []interface{}{"alice", "bob", "charlie"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expr, tt.input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				gotArr := got.([]interface{})
				if len(gotArr) != len(tt.want) {
					t.Errorf("sort() length = %d, want %d", len(gotArr), len(tt.want))
					return
				}
				for i := range gotArr {
					if gotArr[i] != tt.want[i] {
						t.Errorf("sort()[%d] = %v, want %v", i, gotArr[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestArraySliceFunction(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		input   interface{}
		want    []interface{}
		wantErr bool
	}{
		{
			name: "slice with start only",
			expr: "slice(item.values, 2)",
			input: map[string]interface{}{
				"values": []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			want: []interface{}{3.0, 4.0, 5.0},
		},
		{
			name: "slice with start and end",
			expr: "slice(item.values, 1, 4)",
			input: map[string]interface{}{
				"values": []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			want: []interface{}{2.0, 3.0, 4.0},
		},
		{
			name: "slice with negative start",
			expr: "slice(item.values, -2)",
			input: map[string]interface{}{
				"values": []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			want: []interface{}{4.0, 5.0},
		},
		{
			name: "slice with negative end",
			expr: "slice(item.values, 1, -1)",
			input: map[string]interface{}{
				"values": []interface{}{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			want: []interface{}{2.0, 3.0, 4.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expr, tt.input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				gotArr := got.([]interface{})
				if len(gotArr) != len(tt.want) {
					t.Errorf("slice() length = %d, want %d", len(gotArr), len(tt.want))
					return
				}
				for i := range gotArr {
					if gotArr[i] != tt.want[i] {
						t.Errorf("slice()[%d] = %v, want %v", i, gotArr[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestArrayUniqueFunction(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		input   interface{}
		want    []interface{}
		wantErr bool
	}{
		{
			name: "unique numbers",
			expr: "unique(item.values)",
			input: map[string]interface{}{
				"values": []interface{}{1.0, 2.0, 2.0, 3.0, 1.0, 4.0},
			},
			want: []interface{}{1.0, 2.0, 3.0, 4.0},
		},
		{
			name: "unique strings",
			expr: "unique(item.tags)",
			input: map[string]interface{}{
				"tags": []interface{}{"a", "b", "a", "c", "b"},
			},
			want: []interface{}{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expr, tt.input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				gotArr := got.([]interface{})
				if len(gotArr) != len(tt.want) {
					t.Errorf("unique() length = %d, want %d", len(gotArr), len(tt.want))
					return
				}
				for i := range gotArr {
					if gotArr[i] != tt.want[i] {
						t.Errorf("unique()[%d] = %v, want %v", i, gotArr[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestArrayReverseFunction(t *testing.T) {
	input := map[string]interface{}{
		"values": []interface{}{1.0, 2.0, 3.0, 4.0},
	}

	got, err := EvaluateExpression("reverse(item.values)", input, nil)
	if err != nil {
		t.Fatalf("reverse() error: %v", err)
	}

	want := []interface{}{4.0, 3.0, 2.0, 1.0}
	gotArr := got.([]interface{})
	if len(gotArr) != len(want) {
		t.Errorf("reverse() length = %d, want %d", len(gotArr), len(want))
		return
	}
	for i := range gotArr {
		if gotArr[i] != want[i] {
			t.Errorf("reverse()[%d] = %v, want %v", i, gotArr[i], want[i])
		}
	}
}

func TestArrayFlattenFunction(t *testing.T) {
	input := map[string]interface{}{
		"nested": []interface{}{
			[]interface{}{1.0, 2.0},
			[]interface{}{3.0, 4.0},
			5.0,
		},
	}

	got, err := EvaluateExpression("flatten(item.nested)", input, nil)
	if err != nil {
		t.Fatalf("flatten() error: %v", err)
	}

	want := []interface{}{1.0, 2.0, 3.0, 4.0, 5.0}
	gotArr := got.([]interface{})
	if len(gotArr) != len(want) {
		t.Errorf("flatten() length = %d, want %d", len(gotArr), len(want))
		return
	}
	for i := range gotArr {
		if gotArr[i] != want[i] {
			t.Errorf("flatten()[%d] = %v, want %v", i, gotArr[i], want[i])
		}
	}
}

func TestArrayZipFunction(t *testing.T) {
	input := map[string]interface{}{
		"names": []interface{}{"alice", "bob", "charlie"},
		"ages":  []interface{}{30.0, 25.0, 35.0},
	}

	got, err := EvaluateExpression("zip(item.names, item.ages)", input, nil)
	if err != nil {
		t.Fatalf("zip() error: %v", err)
	}

	gotArr := got.([]interface{})
	if len(gotArr) != 3 {
		t.Errorf("zip() length = %d, want 3", len(gotArr))
		return
	}

	// Check first tuple
	tuple0 := gotArr[0].([]interface{})
	if tuple0[0] != "alice" || tuple0[1] != 30.0 {
		t.Errorf("zip()[0] = %v, want [alice, 30]", tuple0)
	}
}

func TestComplexComposition(t *testing.T) {
	input := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "alice", "age": 30.0, "score": 85.5},
			map[string]interface{}{"name": "bob", "age": 25.0, "score": 92.3},
			map[string]interface{}{"name": "charlie", "age": 35.0, "score": 78.9},
		},
	}

	tests := []struct {
		name    string
		expr    string
		want    interface{}
		wantErr bool
	}{
		{
			name: "avg of mapped ages",
			expr: "avg(map(item.users, item.age))",
			want: 30.0,
		},
		{
			name: "sum of mapped scores",
			expr: "sum(map(item.users, item.score))",
			want: 256.70000000000005, // Floating point precision
		},
		{
			name: "round avg of scores",
			expr: "round(avg(map(item.users, item.score)))",
			want: 86.0,
		},
		{
			name: "max of mapped ages",
			expr: "max(map(item.users, item.age))",
			want: 35.0,
		},
		{
			name: "min of mapped scores",
			expr: "min(map(item.users, item.score))",
			want: 78.9,
		},
		{
			name: "sort mapped ages",
			expr: "sort(map(item.users, item.age))",
			want: []interface{}{25.0, 30.0, 35.0},
		},
		{
			name: "slice first 2 users ages",
			expr: "slice(map(item.users, item.age), 0, 2)",
			want: []interface{}{30.0, 25.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expr, input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if arr, ok := got.([]interface{}); ok {
					wantArr := tt.want.([]interface{})
					if len(arr) != len(wantArr) {
						t.Errorf("array length = %d, want %d", len(arr), len(wantArr))
						return
					}
					for i := range arr {
						if arr[i] != wantArr[i] {
							t.Errorf("array[%d] = %v, want %v", i, arr[i], wantArr[i])
						}
					}
				} else if got != tt.want {
					t.Errorf("EvaluateExpression() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
