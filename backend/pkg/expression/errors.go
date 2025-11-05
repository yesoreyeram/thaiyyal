package expression

import (
	"errors"
	"fmt"
)

// Sentinel errors for expression evaluation
var (
	// Parsing errors
	ErrSyntaxError          = errors.New("expression syntax error")
	ErrUnexpectedToken      = errors.New("unexpected token")
	ErrUnmatchedParenthesis = errors.New("unmatched parenthesis")
	ErrInvalidOperator      = errors.New("invalid operator")

	// Evaluation errors
	ErrEvaluationFailed  = errors.New("expression evaluation failed")
	ErrUndefinedVariable = errors.New("undefined variable")
	ErrUndefinedFunction = errors.New("undefined function")
	ErrInvalidArgument   = errors.New("invalid function argument")
	ErrTypeMismatch      = errors.New("type mismatch in expression")

	// Function errors
	ErrFunctionNotFound     = errors.New("function not found")
	ErrInvalidArgumentCount = errors.New("invalid number of arguments")
	ErrArgumentTypeMismatch = errors.New("function argument type mismatch")

	// Field access errors
	ErrFieldNotFound      = errors.New("field not found")
	ErrInvalidFieldAccess = errors.New("invalid field access")
	ErrIndexOutOfBounds   = errors.New("index out of bounds")

	// Security errors
	ErrExpressionTooComplex   = errors.New("expression too complex")
	ErrRecursionDepthExceeded = errors.New("maximum recursion depth exceeded")
)

// ExpressionError represents an error that occurred during expression evaluation
// with additional context about where and why the error occurred
type ExpressionError struct {
	Expression string // The full expression being evaluated
	Position   int    // Position in the expression where error occurred (if known)
	Message    string // The error message
	Context    string // Additional context about what was being done
	Cause      error  // The underlying error (if any)
}

// Error implements the error interface
func (e *ExpressionError) Error() string {
	if e.Position >= 0 && e.Position < len(e.Expression) {
		// Show the expression with a pointer to the error position
		pointer := ""
		if e.Position > 0 {
			pointer = fmt.Sprintf("%s^", repeatString(" ", e.Position))
		} else {
			pointer = "^"
		}
		
		msg := fmt.Sprintf("Expression error: %s\n  %s\n  %s", e.Message, e.Expression, pointer)
		if e.Context != "" {
			msg += fmt.Sprintf("\n  Context: %s", e.Context)
		}
		if e.Cause != nil {
			msg += fmt.Sprintf("\n  Caused by: %s", e.Cause.Error())
		}
		return msg
	}
	
	msg := fmt.Sprintf("Expression error: %s", e.Message)
	if e.Expression != "" {
		msg += fmt.Sprintf("\n  Expression: %s", e.Expression)
	}
	if e.Context != "" {
		msg += fmt.Sprintf("\n  Context: %s", e.Context)
	}
	if e.Cause != nil {
		msg += fmt.Sprintf("\n  Caused by: %s", e.Cause.Error())
	}
	return msg
}

// Unwrap returns the underlying cause error
func (e *ExpressionError) Unwrap() error {
	return e.Cause
}

// repeatString repeats a string n times
func repeatString(s string, n int) string {
	if n <= 0 {
		return ""
	}
	result := make([]byte, n*len(s))
	for i := 0; i < n; i++ {
		copy(result[i*len(s):], s)
	}
	return string(result)
}

// newExpressionError creates a new ExpressionError
func newExpressionError(expr string, message string) *ExpressionError {
	return &ExpressionError{
		Expression: expr,
		Position:   -1,
		Message:    message,
	}
}

// newExpressionErrorWithPos creates a new ExpressionError with position
func newExpressionErrorWithPos(expr string, pos int, message string) *ExpressionError {
	return &ExpressionError{
		Expression: expr,
		Position:   pos,
		Message:    message,
	}
}

// newExpressionErrorWithContext creates a new ExpressionError with context
func newExpressionErrorWithContext(expr string, message string, context string) *ExpressionError {
	return &ExpressionError{
		Expression: expr,
		Position:   -1,
		Message:    message,
		Context:    context,
	}
}

// newExpressionErrorWithCause creates a new ExpressionError with an underlying cause
func newExpressionErrorWithCause(expr string, message string, cause error) *ExpressionError {
	return &ExpressionError{
		Expression: expr,
		Position:   -1,
		Message:    message,
		Cause:      cause,
	}
}
