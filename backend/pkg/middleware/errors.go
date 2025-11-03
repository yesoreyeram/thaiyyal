package middleware

import "errors"

// Sentinel errors for middleware operations
var (
	// Middleware errors
	ErrMiddlewareExecutionFailed = errors.New("middleware execution failed")
	ErrMiddlewareChainBroken     = errors.New("middleware chain broken")
	ErrHandlerNotCalled          = errors.New("next handler not called")

	// Configuration errors
	ErrInvalidMiddleware  = errors.New("invalid middleware")
	ErrMiddlewareNotFound = errors.New("middleware not found")

	// Context errors
	ErrContextValueNotFound = errors.New("context value not found")
	ErrInvalidContextValue  = errors.New("invalid context value")
)
