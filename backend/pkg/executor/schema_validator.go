package executor

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SchemaValidatorExecutor validates data against JSON schemas
type SchemaValidatorExecutor struct{}

// NodeType returns the node type
func (e *SchemaValidatorExecutor) NodeType() types.NodeType {
	return types.NodeTypeSchemaValidator
}

// Execute validates input data against the provided JSON schema
func (e *SchemaValidatorExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	// Get input from node inputs
	inputs := ctx.GetNodeInputs(node.ID)
	var input interface{}
	if len(inputs) > 0 {
		input = inputs[0]
	}
	if input == nil {
		return nil, fmt.Errorf("no input provided for validation")
	}

	// Get schema from node data
	schemaData := node.Data.Schema
	if schemaData == nil {
		return nil, fmt.Errorf("schema not provided")
	}

	// Get strict mode (default: false - return validation errors as metadata)
	strict := false
	if node.Data.Strict != nil {
		strict = *node.Data.Strict
	}

	// Convert schema to JSON for validation
	schemaBytes, err := json.Marshal(schemaData)
	if err != nil {
		return nil, fmt.Errorf("invalid schema format: %w", err)
	}

	// Create schema loader
	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)

	// Convert input to JSON for validation
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize input: %w", err)
	}

	// Create document loader
	documentLoader := gojsonschema.NewBytesLoader(inputBytes)

	// Validate
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	// If validation passed, return input with metadata
	if result.Valid() {
		return map[string]interface{}{
			"valid": true,
			"data":  input,
		}, nil
	}

	// Validation failed - collect errors
	errors := make([]map[string]interface{}, 0, len(result.Errors()))
	for _, err := range result.Errors() {
		errors = append(errors, map[string]interface{}{
			"field":       err.Field(),
			"type":        err.Type(),
			"description": err.Description(),
			"value":       err.Value(),
		})
	}

	// In strict mode, return error
	if strict {
		return nil, fmt.Errorf("validation failed: %d errors found", len(errors))
	}

	// In lenient mode, return validation result with errors
	return map[string]interface{}{
		"valid":  false,
		"data":   input,
		"errors": errors,
	}, nil
}

// Validate checks if the node configuration is valid
func (e *SchemaValidatorExecutor) Validate(node types.Node) error {
	// Check if schema is provided
	if node.Data.Schema == nil {
		return fmt.Errorf("schema is required")
	}

	// Validate schema is a valid object
	schema, ok := node.Data.Schema.(map[string]interface{})
	if !ok {
		return fmt.Errorf("schema must be an object")
	}

	// Check if schema has type field
	if _, ok := schema["type"]; !ok {
		return fmt.Errorf("schema must have a 'type' field")
	}

	// Validate strict mode if provided (already boolean type from NodeData)
	// No validation needed as it's typed correctly

	return nil
}
