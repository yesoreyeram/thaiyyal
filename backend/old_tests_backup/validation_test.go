package workflow

import (
	"testing"
)

// ============================================================================
// Workflow Validation Tests
// ============================================================================

func TestValidation_EmptyWorkflow(t *testing.T) {
	payload := `{"nodes": [], "edges": []}`
	result, err := ValidatePayload([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}

	if result.Valid {
		t.Error("expected workflow to be invalid (no nodes)")
	}

	if len(result.Errors) == 0 {
		t.Error("expected validation errors but got none")
	}

	// Check for specific error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "nodes" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about missing nodes")
	}
}

func TestValidation_DuplicateNodeIDs(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "1", "data": {"value": 20}}
		],
		"edges": []
	}`

	result, err := ValidatePayload([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}

	if result.Valid {
		t.Error("expected workflow to be invalid (duplicate IDs)")
	}

	// Check for duplicate ID error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "id" && err.NodeID == "1" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about duplicate node ID")
	}
}

func TestValidation_EmptyNodeID(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "", "data": {"value": 10}}
		],
		"edges": []
	}`

	result, err := ValidatePayload([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}

	if result.Valid {
		t.Error("expected workflow to be invalid (empty node ID)")
	}
}

func TestValidation_InvalidEdgeSource(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 20}}
		],
		"edges": [
			{"source": "999", "target": "2"}
		]
	}`

	result, err := ValidatePayload([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}

	if result.Valid {
		t.Error("expected workflow to be invalid (invalid edge source)")
	}

	// Check for edge error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "edges[0].source" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about invalid edge source")
	}
}

func TestValidation_InvalidEdgeTarget(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 20}}
		],
		"edges": [
			{"source": "1", "target": "999"}
		]
	}`

	result, err := ValidatePayload([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}

	if result.Valid {
		t.Error("expected workflow to be invalid (invalid edge target)")
	}
}

func TestValidation_SelfReferencingEdge(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}}
		],
		"edges": [
			{"source": "1", "target": "1"}
		]
	}`

	result, err := ValidatePayload([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}

	if result.Valid {
		t.Error("expected workflow to be invalid (self-referencing edge)")
	}

	// Check for self-referential error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "edges[0]" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about self-referential edge")
	}
}

func TestValidation_CyclicWorkflow(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 20}},
			{"id": "3", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "2"},
			{"source": "2", "target": "3"},
			{"source": "3", "target": "1"}
		]
	}`

	result, err := ValidatePayload([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}

	if result.Valid {
		t.Error("expected workflow to be invalid (contains cycle)")
	}

	// Check for cycle error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "edges" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about cyclic dependencies")
	}
}

func TestValidation_NumberNodeMissingValue(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {}}
		],
		"edges": []
	}`

	// Explicitly set the type to Number since inference won't work without a value
	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeNumber

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (number node missing value)")
	}

	// Check for missing value error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "data.value" && err.NodeID == "1" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about missing value field")
	}
}

func TestValidation_TextInputNodeMissingText(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": ""}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeTextInput

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (text input missing text)")
	}

	// Check for missing text error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "data.text" && err.NodeID == "1" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about missing text field")
	}
}

func TestValidation_OperationNodeInvalidOp(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"op": "invalid_operation"}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeOperation

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (invalid operation)")
	}

	// Check for invalid operation error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "data.op" && err.NodeID == "1" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about invalid operation")
	}
}

func TestValidation_OperationNodeMissingOp(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeOperation

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (operation node missing op)")
	}
}

func TestValidation_HTTPNodeMissingURL(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeHTTP

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (HTTP node missing URL)")
	}

	// Check for missing URL error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "data.url" && err.NodeID == "1" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about missing URL field")
	}
}

func TestValidation_HTTPNodeInvalidMethod(t *testing.T) {
	// Note: HTTP nodes currently only support GET requests and don't have a method field
	// This test validates that URL is still required
	url := "http://example.com"
	payload := `{
		"nodes": [
			{"id": "1", "data": {"url": "http://example.com"}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeHTTP
	engine.nodes[0].Data.URL = &url

	result := engine.Validate()

	if !result.Valid {
		t.Error("expected workflow to be valid (HTTP node with URL)")
	}
}

func TestValidation_ConditionNodeMissingCondition(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeCondition

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (condition node missing condition)")
	}
}

func TestValidation_VariableNodeMissingKey(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeVariable

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (variable node missing var_name)")
	}

	// Check for missing var_name error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "data.var_name" && err.NodeID == "1" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about missing var_name field")
	}
}

func TestValidation_ContextConstantMissingValue(t *testing.T) {
	contextName := "myconst"
	payload := `{
		"nodes": [
			{"id": "1", "data": {"context_name": "myconst"}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeContextConstant
	engine.nodes[0].Data.ContextName = &contextName

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (context constant missing context_value)")
	}

	// Check for missing context_value error
	foundError := false
	for _, err := range result.Errors {
		if err.Field == "data.context_value" && err.NodeID == "1" {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("expected error about missing context_value field")
	}
}

func TestValidation_RetryNodeInvalidMaxAttempts(t *testing.T) {
	maxAttempts := 0
	payload := `{
		"nodes": [
			{"id": "1", "data": {"max_attempts": 0}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeRetry
	engine.nodes[0].Data.MaxAttempts = &maxAttempts

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (retry max_attempts < 1)")
	}
}

func TestValidation_TimeoutNodeMissingTimeout(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeTimeout

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (timeout node missing timeout)")
	}
}

func TestValidation_ValidSimpleWorkflow(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "3"},
			{"source": "2", "target": "3"}
		]
	}`

	result, err := ValidatePayload([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}

	if !result.Valid {
		t.Errorf("expected workflow to be valid but got errors: %v", result.Errors)
	}

	if len(result.Errors) != 0 {
		t.Errorf("expected no validation errors but got %d: %v", len(result.Errors), result.Errors)
	}
}

func TestValidation_ValidComplexWorkflow(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "data": {"op": "add"}},
			{"id": "4", "data": {"op": "multiply"}},
			{"id": "5", "data": {"value": 2}}
		],
		"edges": [
			{"source": "1", "target": "3"},
			{"source": "2", "target": "3"},
			{"source": "3", "target": "4"},
			{"source": "5", "target": "4"}
		]
	}`

	result, err := ValidatePayload([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error parsing payload: %v", err)
	}

	if !result.Valid {
		t.Errorf("expected workflow to be valid but got errors: %v", result.Errors)
	}
}

func TestValidation_TextOperationInvalidOp(t *testing.T) {
	textOp := "invalid_text_op"
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text_op": "invalid_text_op"}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeTextOperation
	engine.nodes[0].Data.TextOp = &textOp

	result := engine.Validate()

	if result.Valid {
		t.Error("expected workflow to be invalid (invalid text operation)")
	}
}

func TestValidation_ValidHTTPMethods(t *testing.T) {
	// Note: HTTP nodes currently only support GET and don't have a method field
	// This test verifies that a valid HTTP node with URL passes validation
	url := "http://example.com"
	payload := `{
		"nodes": [
			{"id": "1", "data": {"url": "http://example.com"}}
		],
		"edges": []
	}`

	engine, _ := NewEngine([]byte(payload))
	engine.nodes[0].Type = NodeTypeHTTP
	engine.nodes[0].Data.URL = &url

	result := engine.Validate()

	if !result.Valid {
		t.Errorf("expected HTTP node with valid URL to be valid but got errors: %v", result.Errors)
	}
}

func TestValidation_ValidOperations(t *testing.T) {
	validOps := []string{"add", "subtract", "multiply", "divide", "modulo", "power"}

	for _, op := range validOps {
		o := op
		payload := `{
			"nodes": [
				{"id": "1", "data": {"op": "` + op + `"}}
			],
			"edges": []
		}`

		engine, _ := NewEngine([]byte(payload))
		engine.nodes[0].Type = NodeTypeOperation
		engine.nodes[0].Data.Op = &o

		result := engine.Validate()

		if !result.Valid {
			t.Errorf("expected operation %s to be valid but got errors: %v", op, result.Errors)
		}
	}
}

func TestValidation_ValidTextOperations(t *testing.T) {
	validOps := []string{"concat", "uppercase", "lowercase", "trim", "split", "replace", "length", "substring"}

	for _, textOp := range validOps {
		to := textOp
		payload := `{
			"nodes": [
				{"id": "1", "data": {"text_op": "` + textOp + `"}}
			],
			"edges": []
		}`

		engine, _ := NewEngine([]byte(payload))
		engine.nodes[0].Type = NodeTypeTextOperation
		engine.nodes[0].Data.TextOp = &to

		result := engine.Validate()

		if !result.Valid {
			t.Errorf("expected text operation %s to be valid but got errors: %v", textOp, result.Errors)
		}
	}
}
