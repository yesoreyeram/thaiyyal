package expression

import "errors"

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
