package types

import "fmt"

// ErrMissingRequiredField creates an error for missing required field
func ErrMissingRequiredField(fieldName string) error {
	return fmt.Errorf("missing required field: %s", fieldName)
}

// ErrInvalidFieldValue creates an error for invalid field value
func ErrInvalidFieldValue(fieldName string, value interface{}, reason string) error {
	return fmt.Errorf("invalid value for field %s: %v (%s)", fieldName, value, reason)
}

// ErrUnknownNodeType creates an error for unknown node type
func ErrUnknownNodeType(nodeType NodeType) error {
	return fmt.Errorf("unknown node type: %s", nodeType)
}
