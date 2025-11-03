package executor

import "errors"

// Sentinel errors for executor operations
var (
// Input validation errors
ErrInvalidInput       = errors.New("invalid input")
ErrMissingRequiredInput = errors.New("missing required input")
ErrInputTypeMismatch  = errors.New("input type mismatch")
ErrInvalidInputValue  = errors.New("invalid input value")

// Operation errors
ErrDivisionByZero     = errors.New("division by zero")
ErrInvalidOperation   = errors.New("invalid operation")
ErrOperationFailed    = errors.New("operation failed")
ErrUnsupportedOperation = errors.New("unsupported operation")

// HTTP errors
ErrHTTPRequestFailed  = errors.New("HTTP request failed")
ErrHTTPTimeout        = errors.New("HTTP request timeout")
ErrInvalidURL         = errors.New("invalid URL")
ErrURLNotAllowed      = errors.New("URL not allowed by security policy")
ErrMaxRedirectsExceeded = errors.New("maximum redirects exceeded")
ErrResponseTooLarge   = errors.New("response body too large")

// Array operation errors
ErrInvalidArrayIndex  = errors.New("invalid array index")
ErrArrayEmpty         = errors.New("array is empty")
ErrNotAnArray         = errors.New("value is not an array")
ErrArrayTooLarge      = errors.New("array exceeds maximum size")

// Expression errors
ErrExpressionEvaluation = errors.New("expression evaluation failed")
ErrInvalidExpression  = errors.New("invalid expression")

// Loop errors
ErrMaxLoopIterations  = errors.New("maximum loop iterations exceeded")
ErrInfiniteLoop       = errors.New("infinite loop detected")

// Cache errors
ErrCacheKeyNotFound   = errors.New("cache key not found")
ErrCacheExpired       = errors.New("cache entry expired")

// Retry errors
ErrMaxAttemptsExceeded = errors.New("maximum retry attempts exceeded")
ErrRetryFailed        = errors.New("retry failed")
)
