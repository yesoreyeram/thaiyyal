package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestExpressionExecutor_Execute(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		want       interface{}
		wantErr    bool
	}{
		{
			name:       "simple multiplication",
			expression: "input * 2",
			input:      5.0,
			want:       10.0,
			wantErr:    false,
		},
		{
			name:       "simple addition",
			expression: "input + 10",
			input:      15.0,
			want:       25.0,
			wantErr:    false,
		},
		{
			name:       "field access",
			expression: "input.price * input.quantity",
			input: map[string]interface{}{
				"price":    10.0,
				"quantity": 3.0,
			},
			want:    30.0,
			wantErr: false,
		},
		{
			name:       "nested field access",
			expression: "input.user.name",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
				},
			},
			want:    "Alice",
			wantErr: false,
		},
		{
			name:       "comparison greater than (true)",
			expression: "input > 2",
			input:      5.0,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "comparison greater than (false)",
			expression: "input > 10",
			input:      5.0,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "comparison equal",
			expression: "input == 5",
			input:      5.0,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "comparison less than or equal",
			expression: "input <= 10",
			input:      5.0,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "field comparison",
			expression: "input.age > 18",
			input: map[string]interface{}{
				"age": 25.0,
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &ExpressionExecutor{}
			node := types.Node{
				ID:   "expr1",
				Type: types.NodeTypeExpression,
				Data: types.NodeData{
					Expression: &tt.expression,
				},
			}

			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"expr1": {tt.input},
				},
			}

			got, err := executor.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !compareValues(got, tt.want) {
				t.Errorf("Execute() = %v (type: %T), want %v (type: %T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestExpressionExecutor_Execute_NoExpression(t *testing.T) {
	executor := &ExpressionExecutor{}
	emptyExpr := ""
	node := types.Node{
		ID:   "expr1",
		Type: types.NodeTypeExpression,
		Data: types.NodeData{
			Expression: &emptyExpr,
		},
	}

	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"expr1": {42.0},
		},
	}

	got, err := executor.Execute(ctx, node)
	if err != nil {
		t.Errorf("Execute() error = %v, expected no error with warning", err)
		return
	}

	// Should return a map with result and warning
	result, ok := got.(map[string]interface{})
	if !ok {
		t.Errorf("Execute() returned %T, want map[string]interface{}", got)
		return
	}

	if result["result"] != 42.0 {
		t.Errorf("Execute() result = %v, want 42.0", result["result"])
	}

	if result["warning"] == nil {
		t.Error("Execute() should return warning when expression is empty")
	}
}

func TestExpressionExecutor_Execute_NoInput(t *testing.T) {
	executor := &ExpressionExecutor{}
	expr := "input * 2"
	node := types.Node{
		ID:   "expr1",
		Type: types.NodeTypeExpression,
		Data: types.NodeData{
			Expression: &expr,
		},
	}

	// Create context with no inputs for expr1
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{},
	}

	_, err := executor.Execute(ctx, node)
	if err == nil {
		t.Error("Execute() should return error when no input provided")
	}
}

func TestExpressionExecutor_Validate(t *testing.T) {
	tests := []struct {
		name    string
		node    types.Node
		wantErr bool
	}{
		{
			name: "valid expression node",
			node: types.Node{
				Type: types.NodeTypeExpression,
				Data: types.NodeData{
					Expression: stringPtr("input * 2"),
				},
			},
			wantErr: false,
		},
		{
			name: "missing expression",
			node: types.Node{
				Type: types.NodeTypeExpression,
				Data: types.NodeData{
					Expression: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "empty expression",
			node: types.Node{
				Type: types.NodeTypeExpression,
				Data: types.NodeData{
					Expression: stringPtr(""),
				},
			},
			wantErr: true,
		},
		{
			name: "expression too long",
			node: types.Node{
				Type: types.NodeTypeExpression,
				Data: types.NodeData{
					Expression: stringPtr(string(make([]byte, 10001))),
				},
			},
			wantErr: true,
		},
		{
			name: "wrong node type",
			node: types.Node{
				Type: types.NodeTypeNumber,
				Data: types.NodeData{
					Expression: stringPtr("input * 2"),
				},
			},
			wantErr: true,
		},
	}

	executor := &ExpressionExecutor{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.Validate(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExpressionExecutor_NodeType(t *testing.T) {
	executor := &ExpressionExecutor{}
	if executor.NodeType() != types.NodeTypeExpression {
		t.Errorf("NodeType() = %v, want %v", executor.NodeType(), types.NodeTypeExpression)
	}
}
