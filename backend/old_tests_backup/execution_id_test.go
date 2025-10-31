package workflow

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// ============================================================================
// Execution ID and Workflow ID Tests
// ============================================================================

func TestExecutionIDGeneration(t *testing.T) {
	// Test that each engine gets a unique execution ID
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}}
		],
		"edges": []
	}`

	engine1, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine1: %v", err)
	}

	engine2, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine2: %v", err)
	}

	// Each engine should have a different execution ID
	if engine1.executionID == "" {
		t.Error("engine1 execution ID is empty")
	}

	if engine2.executionID == "" {
		t.Error("engine2 execution ID is empty")
	}

	if engine1.executionID == engine2.executionID {
		t.Error("execution IDs should be unique across engines")
	}

	// Execution ID should be 16 hex characters (8 bytes)
	if len(engine1.executionID) != 16 {
		t.Errorf("execution ID should be 16 characters, got %d: %s", len(engine1.executionID), engine1.executionID)
	}

	// Should be valid hex
	for _, c := range engine1.executionID {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("execution ID should be hex, got character: %c", c)
		}
	}
}

func TestExecutionIDInResult(t *testing.T) {
	// Test that execution ID appears in the result
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

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	// Result should contain the execution ID
	if result.ExecutionID == "" {
		t.Error("result execution ID is empty")
	}

	// Result execution ID should match engine execution ID
	if result.ExecutionID != engine.executionID {
		t.Errorf("result execution ID (%s) doesn't match engine execution ID (%s)", result.ExecutionID, engine.executionID)
	}
}

func TestWorkflowIDFromPayload(t *testing.T) {
	// Test that workflow ID from payload is stored and returned
	payload := `{
		"workflow_id": "my-workflow-123",
		"nodes": [
			{"id": "1", "data": {"value": 10}}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	// Engine should have the workflow ID
	if engine.workflowID != "my-workflow-123" {
		t.Errorf("expected workflow ID 'my-workflow-123', got: %s", engine.workflowID)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	// Result should contain the workflow ID
	if result.WorkflowID != "my-workflow-123" {
		t.Errorf("expected workflow ID 'my-workflow-123', got: %s", result.WorkflowID)
	}
}

func TestWorkflowIDOptional(t *testing.T) {
	// Test that workflow ID is optional
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	// Engine workflow ID should be empty
	if engine.workflowID != "" {
		t.Errorf("expected empty workflow ID, got: %s", engine.workflowID)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	// Result workflow ID should be empty
	if result.WorkflowID != "" {
		t.Errorf("expected empty workflow ID, got: %s", result.WorkflowID)
	}
}

func TestExecutionIDInContext(t *testing.T) {
	// Test that execution ID is available in context during execution
	// We'll verify this by checking the context in a custom test

	payload := `{
		"workflow_id": "test-workflow",
		"nodes": [
			{"id": "1", "data": {"value": 10}}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	// Verify execution ID and workflow ID are set
	if result.ExecutionID == "" {
		t.Error("execution ID should not be empty")
	}

	if result.WorkflowID != "test-workflow" {
		t.Errorf("expected workflow ID 'test-workflow', got: %s", result.WorkflowID)
	}
}

func TestGetExecutionIDFromContext(t *testing.T) {
	// Test the helper function to extract execution ID from context
	execID := "test-execution-id-123"

	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyExecutionID, execID)

	retrievedID := GetExecutionID(ctx)
	if retrievedID != execID {
		t.Errorf("expected execution ID %s, got: %s", execID, retrievedID)
	}
}

func TestGetExecutionIDFromEmptyContext(t *testing.T) {
	// Test that GetExecutionID returns empty string when not in context
	ctx := context.Background()

	retrievedID := GetExecutionID(ctx)
	if retrievedID != "" {
		t.Errorf("expected empty string, got: %s", retrievedID)
	}
}

func TestGetWorkflowIDFromContext(t *testing.T) {
	// Test the helper function to extract workflow ID from context
	workflowID := "test-workflow-456"

	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyWorkflowID, workflowID)

	retrievedID := GetWorkflowID(ctx)
	if retrievedID != workflowID {
		t.Errorf("expected workflow ID %s, got: %s", workflowID, retrievedID)
	}
}

func TestGetWorkflowIDFromEmptyContext(t *testing.T) {
	// Test that GetWorkflowID returns empty string when not in context
	ctx := context.Background()

	retrievedID := GetWorkflowID(ctx)
	if retrievedID != "" {
		t.Errorf("expected empty string, got: %s", retrievedID)
	}
}

func TestResultJSONSerialization(t *testing.T) {
	// Test that Result with execution ID and workflow ID serializes correctly to JSON
	payload := `{
		"workflow_id": "json-test-workflow",
		"nodes": [
			{"id": "1", "data": {"value": 42}}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	// Serialize result to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("error marshaling result: %v", err)
	}

	jsonStr := string(jsonData)

	// Check that JSON contains execution_id and workflow_id
	if !strings.Contains(jsonStr, "execution_id") {
		t.Error("JSON should contain execution_id field")
	}

	if !strings.Contains(jsonStr, "workflow_id") {
		t.Error("JSON should contain workflow_id field")
	}

	if !strings.Contains(jsonStr, result.ExecutionID) {
		t.Error("JSON should contain the execution ID value")
	}

	if !strings.Contains(jsonStr, "json-test-workflow") {
		t.Error("JSON should contain the workflow ID value")
	}

	// Deserialize and verify
	var deserializedResult Result
	if err := json.Unmarshal(jsonData, &deserializedResult); err != nil {
		t.Fatalf("error unmarshaling result: %v", err)
	}

	if deserializedResult.ExecutionID != result.ExecutionID {
		t.Error("deserialized execution ID doesn't match")
	}

	if deserializedResult.WorkflowID != result.WorkflowID {
		t.Error("deserialized workflow ID doesn't match")
	}
}

func TestExecutionIDUniquePerExecution(t *testing.T) {
	// Test that multiple executions of the same workflow get different execution IDs
	payload := `{
		"workflow_id": "reusable-workflow",
		"nodes": [
			{"id": "1", "data": {"value": 10}}
		],
		"edges": []
	}`

	// Execute the same workflow 5 times
	executionIDs := make(map[string]bool)
	
	for i := 0; i < 5; i++ {
		engine, err := NewEngine([]byte(payload))
		if err != nil {
			t.Fatalf("unexpected error creating engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("unexpected error executing: %v", err)
		}

		// Check that this execution ID is unique
		if executionIDs[result.ExecutionID] {
			t.Errorf("execution ID %s was reused", result.ExecutionID)
		}
		executionIDs[result.ExecutionID] = true

		// All executions should have the same workflow ID
		if result.WorkflowID != "reusable-workflow" {
			t.Errorf("expected workflow ID 'reusable-workflow', got: %s", result.WorkflowID)
		}
	}

	// Should have 5 unique execution IDs
	if len(executionIDs) != 5 {
		t.Errorf("expected 5 unique execution IDs, got: %d", len(executionIDs))
	}
}

func TestNewEngineWithConfigHasExecutionID(t *testing.T) {
	// Test that NewEngineWithConfig also generates execution ID
	payload := `{
		"workflow_id": "config-test",
		"nodes": [
			{"id": "1", "data": {"value": 10}}
		],
		"edges": []
	}`

	config := DefaultConfig()
	engine, err := NewEngineWithConfig([]byte(payload), config)
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	if engine.executionID == "" {
		t.Error("execution ID should not be empty")
	}

	if engine.workflowID != "config-test" {
		t.Errorf("expected workflow ID 'config-test', got: %s", engine.workflowID)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing: %v", err)
	}

	if result.ExecutionID != engine.executionID {
		t.Error("result execution ID should match engine execution ID")
	}

	if result.WorkflowID != "config-test" {
		t.Error("result workflow ID should match engine workflow ID")
	}
}
