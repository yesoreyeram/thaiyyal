package middleware

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ValidationMiddleware validates node configuration before execution.
// It uses the executor's Validate method to ensure node data is valid.
type ValidationMiddleware struct {
	registry interface {
		Validate(node types.Node) error
	}
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware(registry interface{ Validate(node types.Node) error }) *ValidationMiddleware {
	return &ValidationMiddleware{
		registry: registry,
	}
}

// Process validates node before execution
func (m *ValidationMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	// Validate node configuration
	if m.registry != nil {
		if err := m.registry.Validate(node); err != nil {
			return nil, fmt.Errorf("node validation failed: %w", err)
		}
	}

	// Validation passed, continue execution
	return next(ctx, node)
}

// Name returns the middleware name
func (m *ValidationMiddleware) Name() string {
	return "Validation"
}

// InputValidationMiddleware validates node inputs before execution
type InputValidationMiddleware struct {
	maxInputSize int64 // Maximum size for input data in bytes
}

// NewInputValidationMiddleware creates a new input validation middleware
func NewInputValidationMiddleware(maxInputSize int64) *InputValidationMiddleware {
	return &InputValidationMiddleware{
		maxInputSize: maxInputSize,
	}
}

// Process validates inputs before execution
func (m *InputValidationMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	// Get node inputs
	inputs := ctx.GetNodeInputs(node.ID)

	// Validate input count (basic check)
	if len(inputs) > 100 {
		return nil, fmt.Errorf("too many inputs: %d (max 100)", len(inputs))
	}

	// Validate input sizes (for string inputs)
	for i, input := range inputs {
		if str, ok := input.(string); ok {
			if m.maxInputSize > 0 && int64(len(str)) > m.maxInputSize {
				return nil, fmt.Errorf("input %d too large: %d bytes (max %d)", i, len(str), m.maxInputSize)
			}
		}
	}

	// Validation passed, continue execution
	return next(ctx, node)
}

// Name returns the middleware name
func (m *InputValidationMiddleware) Name() string {
	return "InputValidation"
}
