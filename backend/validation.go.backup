package workflow

import (
	"fmt"
	"strings"
)

// ============================================================================
// Workflow Validation
// ============================================================================
// This file contains validation functions to ensure workflow payloads are
// well-formed before execution. Early validation prevents runtime errors and
// provides better user feedback.
// ============================================================================

// ValidationError represents a validation error with context.
type ValidationError struct {
	Field   string // The field that failed validation
	Message string // Human-readable error message
	NodeID  string // Optional: The node ID related to the error
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	if e.NodeID != "" {
		return fmt.Sprintf("validation error on node '%s' field '%s': %s", e.NodeID, e.Field, e.Message)
	}
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationResult holds the outcome of workflow validation.
type ValidationResult struct {
	Valid  bool              // Whether the workflow is valid
	Errors []ValidationError // List of validation errors
}

// Validate performs comprehensive validation on a workflow payload.
// It checks for:
//   - Structural issues (empty nodes/edges, duplicate IDs)
//   - Graph issues (invalid edges, cycles, orphaned nodes)
//   - Node-specific issues (missing required fields, invalid data)
//
// Returns:
//   - *ValidationResult: Contains validation status and any errors found
func (e *Engine) Validate() *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	// Run all validation checks
	e.validateStructure(result)
	e.validateGraph(result)
	e.validateNodes(result)

	// Update overall validity
	result.Valid = len(result.Errors) == 0

	return result
}

// validateStructure checks basic structural requirements.
func (e *Engine) validateStructure(result *ValidationResult) {
	// Check if workflow has nodes
	if len(e.nodes) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "nodes",
			Message: "workflow must contain at least one node",
		})
		return
	}

	// Check for duplicate node IDs
	nodeIDs := make(map[string]bool)
	for _, node := range e.nodes {
		if node.ID == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "id",
				Message: "node ID cannot be empty",
			})
			continue
		}
		if nodeIDs[node.ID] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "id",
				Message: fmt.Sprintf("duplicate node ID: %s", node.ID),
				NodeID:  node.ID,
			})
		}
		nodeIDs[node.ID] = true
	}
}

// validateGraph checks graph-related issues.
func (e *Engine) validateGraph(result *ValidationResult) {
	// Build node ID set for quick lookup
	nodeIDs := make(map[string]bool)
	for _, node := range e.nodes {
		nodeIDs[node.ID] = true
	}

	// Validate edges
	for i, edge := range e.edges {
		// Check if source node exists
		if edge.Source == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("edges[%d].source", i),
				Message: "edge source cannot be empty",
			})
			continue
		}
		if !nodeIDs[edge.Source] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("edges[%d].source", i),
				Message: fmt.Sprintf("edge source '%s' does not exist", edge.Source),
			})
		}

		// Check if target node exists
		if edge.Target == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("edges[%d].target", i),
				Message: "edge target cannot be empty",
			})
			continue
		}
		if !nodeIDs[edge.Target] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("edges[%d].target", i),
				Message: fmt.Sprintf("edge target '%s' does not exist", edge.Target),
			})
		}

		// Check for self-referential edges
		if edge.Source == edge.Target {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("edges[%d]", i),
				Message: fmt.Sprintf("self-referential edge detected: node '%s' cannot connect to itself", edge.Source),
			})
		}
	}

	// Check for cycles using topological sort
	// Only check if there are no edge validation errors
	if len(result.Errors) == 0 {
		_, err := e.topologicalSort()
		if err != nil {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "edges",
				Message: "workflow contains cycles (circular dependencies)",
			})
		}
	}
}

// validateNodes checks node-specific requirements.
func (e *Engine) validateNodes(result *ValidationResult) {
	// Infer types first if not set
	e.inferNodeTypes()

	for _, node := range e.nodes {
		// Check if node type is supported
		if node.Type == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "type",
				Message: "node type is required but could not be inferred",
				NodeID:  node.ID,
			})
			continue
		}

		// Validate node-specific requirements
		e.validateNodeData(node, result)
	}
}

// validateNodeData validates node-specific data requirements.
func (e *Engine) validateNodeData(node Node, result *ValidationResult) {
	// NodeData is a struct, not a pointer, so we can access it directly
	data := node.Data

	switch node.Type {
	case NodeTypeNumber:
		// Number nodes require a value
		if data.Value == nil {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.value",
				Message: "number node requires a 'value' field",
				NodeID:  node.ID,
			})
		}

	case NodeTypeTextInput:
		// Text input nodes require text
		if data.Text == nil || *data.Text == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.text",
				Message: "text input node requires a non-empty 'text' field",
				NodeID:  node.ID,
			})
		}

	case NodeTypeOperation:
		// Operation nodes require an operation type
		if data.Op == nil || *data.Op == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.op",
				Message: "operation node requires an 'op' field",
				NodeID:  node.ID,
			})
		} else {
			// Validate operation type
			validOps := []string{"add", "subtract", "multiply", "divide", "modulo", "power"}
			valid := false
			for _, validOp := range validOps {
				if *data.Op == validOp {
					valid = true
					break
				}
			}
			if !valid {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "data.op",
					Message: fmt.Sprintf("invalid operation '%s', must be one of: %s", *data.Op, strings.Join(validOps, ", ")),
					NodeID:  node.ID,
				})
			}
		}

	case NodeTypeTextOperation:
		// Text operation nodes require an operation
		if data.TextOp == nil || *data.TextOp == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.text_op",
				Message: "text operation node requires a 'text_op' field",
				NodeID:  node.ID,
			})
		} else {
			// Validate text operation type
			validOps := []string{"concat", "uppercase", "lowercase", "trim", "split", "replace", "length", "substring"}
			valid := false
			for _, validOp := range validOps {
				if *data.TextOp == validOp {
					valid = true
					break
				}
			}
			if !valid {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "data.text_op",
					Message: fmt.Sprintf("invalid text operation '%s', must be one of: %s", *data.TextOp, strings.Join(validOps, ", ")),
					NodeID:  node.ID,
				})
			}
		}

	case NodeTypeHTTP:
		// HTTP nodes require a URL
		if data.URL == nil || *data.URL == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.url",
				Message: "HTTP node requires a non-empty 'url' field",
				NodeID:  node.ID,
			})
		}
		// Note: HTTP method validation would go here if the Method field existed in NodeData
		// Currently, HTTP nodes only support GET requests

	case NodeTypeCondition:
		// Condition nodes require a condition field
		if data.Condition == nil || *data.Condition == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.condition",
				Message: "condition node requires a 'condition' field",
				NodeID:  node.ID,
			})
		}

	case NodeTypeVariable, NodeTypeExtract:
		// Variable nodes require a var_name
		if node.Type == NodeTypeVariable {
			if data.VarName == nil || *data.VarName == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "data.var_name",
					Message: "variable node requires a 'var_name' field",
					NodeID:  node.ID,
				})
			}
		}
		// Extract nodes require a field or fields
		if node.Type == NodeTypeExtract {
			if (data.Field == nil || *data.Field == "") && len(data.Fields) == 0 {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "data.field",
					Message: "extract node requires either 'field' or 'fields' to be specified",
					NodeID:  node.ID,
				})
			}
		}

	case NodeTypeContextVariable, NodeTypeContextConstant:
		// Context nodes require a context_name
		if data.ContextName == nil || *data.ContextName == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.context_name",
				Message: fmt.Sprintf("%s node requires a 'context_name' field", node.Type),
				NodeID:  node.ID,
			})
		}
		// Context constant also requires a context_value
		if node.Type == NodeTypeContextConstant && data.ContextValue == nil {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.context_value",
				Message: "context constant node requires a 'context_value' field",
				NodeID:  node.ID,
			})
		}

	case NodeTypeRetry:
		// Retry nodes should have valid max_attempts if specified
		if data.MaxAttempts != nil && *data.MaxAttempts < 1 {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.max_attempts",
				Message: "retry node max_attempts must be at least 1",
				NodeID:  node.ID,
			})
		}

	case NodeTypeTimeout:
		// Timeout nodes require a timeout value
		if data.Timeout == nil || *data.Timeout == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "data.timeout",
				Message: "timeout node requires a 'timeout' field",
				NodeID:  node.ID,
			})
		}
	}
}

// ValidatePayload is a convenience function to validate a JSON payload before creating an engine.
// This allows validation without creating an Engine instance.
//
// Returns:
//   - *ValidationResult: Contains validation status and any errors found
//   - error: If the payload cannot be parsed as JSON
func ValidatePayload(payloadJSON []byte) (*ValidationResult, error) {
	// Create a temporary engine for validation
	engine, err := NewEngine(payloadJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	// Run validation
	return engine.Validate(), nil
}
