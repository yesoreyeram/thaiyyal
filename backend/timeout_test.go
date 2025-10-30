package workflow

import (
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Workflow Execution Timeout Tests
// ============================================================================

func TestWorkflowExecutionTimeout(t *testing.T) {
	// Create a workflow that would take longer than the timeout
	// We'll use a delay node with a long duration
	payload := `{
		"nodes": [
			{"id": "1", "data": {"duration": "10s"}}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	// Set the node type explicitly to Delay
	engine.nodes[0].Type = NodeTypeDelay

	// Set a very short timeout (1 second)
	engine.config.MaxExecutionTime = 1 * time.Second

	// Execute should timeout
	start := time.Now()
	_, err = engine.Execute()
	duration := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error but got none")
	}

	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected timeout error, got: %v", err)
	}

	// Verify that execution was stopped around the timeout duration
	// Allow some margin for goroutine scheduling
	if duration > 2*time.Second {
		t.Errorf("timeout took too long: expected ~1s, got %v", duration)
	}

	if duration < 500*time.Millisecond {
		t.Errorf("timeout happened too quickly: expected ~1s, got %v", duration)
	}
}

func TestWorkflowExecutionWithinTimeout(t *testing.T) {
	// Create a simple fast workflow
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

	// Set a reasonable timeout
	engine.config.MaxExecutionTime = 5 * time.Second

	// Execute should complete successfully
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.FinalOutput != 15.0 {
		t.Errorf("expected 15, got %v", result.FinalOutput)
	}
}

func TestWorkflowTimeoutWithLongLoop(t *testing.T) {
	// This test creates multiple delay nodes that together exceed the timeout
	// The workflow should timeout before all nodes complete
	payload := `{
		"nodes": [
			{"id": "1", "data": {"duration": "500ms"}},
			{"id": "2", "data": {"duration": "500ms"}},
			{"id": "3", "data": {"duration": "500ms"}}
		],
		"edges": [
			{"source": "1", "target": "2"},
			{"source": "2", "target": "3"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error creating engine: %v", err)
	}

	// Set all nodes to Delay type
	for i := range engine.nodes {
		engine.nodes[i].Type = NodeTypeDelay
	}

	// Set a short timeout (should timeout during execution)
	engine.config.MaxExecutionTime = 1 * time.Second

	// Execute should timeout (3 delays of 500ms = 1.5s > 1s timeout)
	start := time.Now()
	_, err = engine.Execute()
	duration := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error but got none")
	}

	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected timeout error, got: %v", err)
	}

	// Should stop around 1 second (allow margin for goroutine scheduling)
	if duration > 2*time.Second {
		t.Errorf("timeout took too long: expected ~1s, got %v", duration)
	}

	// Verify it ran for at least close to the timeout
	if duration < 800*time.Millisecond {
		t.Errorf("timeout happened too quickly: expected ~1s, got %v", duration)
	}
}

func TestDefaultTimeoutConfiguration(t *testing.T) {
	// Verify default config has reasonable timeout
	config := DefaultConfig()

	if config.MaxExecutionTime != 5*time.Minute {
		t.Errorf("expected default MaxExecutionTime to be 5 minutes, got %v", config.MaxExecutionTime)
	}

	if config.MaxNodeExecutionTime != 30*time.Second {
		t.Errorf("expected default MaxNodeExecutionTime to be 30 seconds, got %v", config.MaxNodeExecutionTime)
	}
}

func TestValidationConfigTimeouts(t *testing.T) {
	// Verify validation config has stricter timeouts
	config := ValidationLimits()

	if config.MaxExecutionTime != 1*time.Minute {
		t.Errorf("expected validation MaxExecutionTime to be 1 minute, got %v", config.MaxExecutionTime)
	}

	if config.MaxNodeExecutionTime != 10*time.Second {
		t.Errorf("expected validation MaxNodeExecutionTime to be 10 seconds, got %v", config.MaxNodeExecutionTime)
	}
}

func TestDevelopmentConfigTimeouts(t *testing.T) {
	// Verify development config has relaxed timeouts
	config := DevelopmentConfig()

	if config.MaxExecutionTime != 30*time.Minute {
		t.Errorf("expected development MaxExecutionTime to be 30 minutes, got %v", config.MaxExecutionTime)
	}

	if config.MaxNodeExecutionTime != 5*time.Minute {
		t.Errorf("expected development MaxNodeExecutionTime to be 5 minutes, got %v", config.MaxNodeExecutionTime)
	}
}
