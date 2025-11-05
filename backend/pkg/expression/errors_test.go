package expression

import (
	"strings"
	"testing"
)

func TestExpressionError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ExpressionError
		contains []string // Strings that should be in the error message
	}{
		{
			name: "basic error",
			err: &ExpressionError{
				Expression: "item.age > 18",
				Message:    "field not found",
			},
			contains: []string{
				"Expression error",
				"field not found",
				"item.age > 18",
			},
		},
		{
			name: "error with position",
			err: &ExpressionError{
				Expression: "item.unknown > 10",
				Position:   5,
				Message:    "field 'unknown' does not exist",
			},
			contains: []string{
				"field 'unknown' does not exist",
				"item.unknown > 10",
				"^", // pointer
			},
		},
		{
			name: "error with context",
			err: &ExpressionError{
				Expression: "item.tags.includes",
				Message:    "method call missing parentheses",
				Context:    "evaluating field path",
			},
			contains: []string{
				"method call missing parentheses",
				"item.tags.includes",
				"Context: evaluating field path",
			},
		},
		{
			name: "error with cause",
			err: &ExpressionError{
				Expression: "variables.count / 0",
				Message:    "division by zero",
				Cause:      ErrEvaluationFailed,
			},
			contains: []string{
				"division by zero",
				"variables.count / 0",
				"Caused by",
				"evaluation failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			for _, substr := range tt.contains {
				if !strings.Contains(errMsg, substr) {
					t.Errorf("Error message should contain '%s', got:\n%s", substr, errMsg)
				}
			}
		})
	}
}

func TestExpressionErrorHelpers(t *testing.T) {
	t.Run("newExpressionError", func(t *testing.T) {
		err := newExpressionError("test expr", "test message")
		if err.Expression != "test expr" {
			t.Errorf("Expression = %s, want 'test expr'", err.Expression)
		}
		if err.Message != "test message" {
			t.Errorf("Message = %s, want 'test message'", err.Message)
		}
		if err.Position != -1 {
			t.Errorf("Position = %d, want -1", err.Position)
		}
	})

	t.Run("newExpressionErrorWithPos", func(t *testing.T) {
		err := newExpressionErrorWithPos("test expr", 5, "test message")
		if err.Position != 5 {
			t.Errorf("Position = %d, want 5", err.Position)
		}
	})

	t.Run("newExpressionErrorWithContext", func(t *testing.T) {
		err := newExpressionErrorWithContext("test expr", "test message", "test context")
		if err.Context != "test context" {
			t.Errorf("Context = %s, want 'test context'", err.Context)
		}
	})

	t.Run("newExpressionErrorWithCause", func(t *testing.T) {
		cause := ErrFieldNotFound
		err := newExpressionErrorWithCause("test expr", "test message", cause)
		if err.Cause != cause {
			t.Errorf("Cause = %v, want %v", err.Cause, cause)
		}
		if err.Unwrap() != cause {
			t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
		}
	})
}

func TestRepeatString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want string
	}{
		{"empty string", "", 5, ""},
		{"zero times", "x", 0, ""},
		{"negative times", "x", -1, ""},
		{"once", "ab", 1, "ab"},
		{"multiple times", "x", 3, "xxx"},
		{"spaces", " ", 5, "     "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repeatString(tt.s, tt.n)
			if got != tt.want {
				t.Errorf("repeatString(%q, %d) = %q, want %q", tt.s, tt.n, got, tt.want)
			}
		})
	}
}

// TestErrorMessageQuality tests that error messages provide helpful information
func TestErrorMessageQuality(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      interface{}
		ctx        *Context
		wantErr    bool
	}{
		{
			name:       "undefined field",
			expression: "item.nonexistent",
			input:      map[string]interface{}{"name": "test"},
			wantErr:    true,
		},
		{
			name:       "array index out of bounds",
			expression: "item.items[10]",
			input:      map[string]interface{}{"items": []interface{}{1, 2, 3}},
			wantErr:    true,
		},
		{
			name:       "method on wrong type",
			expression: "item.age.toUpperCase()",
			input:      map[string]interface{}{"age": 25},
			wantErr:    true,
		},
		{
			name:       "undefined variable",
			expression: "variables.missing",
			input:      nil,
			ctx:        &Context{Variables: make(map[string]interface{})},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EvaluateExpression(tt.expression, tt.input, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Just verify that we get an error - the specific message can be improved later
			if tt.wantErr && err == nil {
				t.Error("Expected an error but got nil")
			}
		})
	}
}
