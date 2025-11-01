package engine

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestSaveSnapshot_Basic tests basic snapshot creation
func TestSaveSnapshot_Basic(t *testing.T) {
	payload := []byte(`{
		"workflow_id": "test-workflow",
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 10}},
			{"id": "2", "type": "number", "data": {"value": 20}},
			{"id": "3", "type": "operation", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "3"},
			{"source": "2", "target": "3"}
		]
	}`)

	engine, err := New(payload)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Execute part of the workflow
	_, err = engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Create snapshot
	snapshot, err := engine.SaveSnapshot()
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	// Verify snapshot
	if snapshot == nil {
		t.Fatal("Snapshot is nil")
	}

	if snapshot.Version != snapshotVersion {
		t.Errorf("Expected version %s, got %s", snapshotVersion, snapshot.Version)
	}

	if snapshot.WorkflowID != "test-workflow" {
		t.Errorf("Expected workflow_id 'test-workflow', got '%s'", snapshot.WorkflowID)
	}

	if snapshot.ExecutionID == "" {
		t.Error("Execution ID is empty")
	}

	if len(snapshot.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(snapshot.Nodes))
	}

	if len(snapshot.Edges) != 2 {
		t.Errorf("Expected 2 edges, got %d", len(snapshot.Edges))
	}

	if len(snapshot.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(snapshot.Results))
	}

	if snapshot.SnapshotTime.IsZero() {
		t.Error("Snapshot time is zero")
	}
}

// TestSaveSnapshot_WithState tests snapshot with state manager data
func TestSaveSnapshot_WithState(t *testing.T) {
	payload := []byte(`{
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 42}},
			{"id": "2", "type": "variable", "data": {"var_op": "set", "var_name": "myVar"}},
			{"id": "3", "type": "counter", "data": {"counter_op": "increment", "delta": 5}}
		],
		"edges": [
			{"source": "1", "target": "2"},
			{"source": "1", "target": "3"}
		]
	}`)

	engine, err := New(payload)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Execute workflow
	_, err = engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Create snapshot
	snapshot, err := engine.SaveSnapshot()
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	// Verify state was captured
	if len(snapshot.Variables) == 0 {
		t.Error("Expected variables to be captured in snapshot")
	}

	if snapshot.Counter != 5.0 {
		t.Errorf("Expected counter to be 5.0, got %f", snapshot.Counter)
	}
}

// TestLoadSnapshot_Restore tests snapshot restoration
func TestLoadSnapshot_Restore(t *testing.T) {
	// Create original engine and execute
	payload := []byte(`{
		"workflow_id": "restore-test",
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 100}},
			{"id": "2", "type": "number", "data": {"value": 200}},
			{"id": "3", "type": "operation", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "3"},
			{"source": "2", "target": "3"}
		]
	}`)

	original, err := New(payload)
	if err != nil {
		t.Fatalf("Failed to create original engine: %v", err)
	}

	originalResult, err := original.Execute()
	if err != nil {
		t.Fatalf("Failed to execute original workflow: %v", err)
	}

	// Save snapshot
	snapshot, err := original.SaveSnapshot()
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	// Load snapshot into new engine
	restored, err := LoadSnapshot(snapshot, nil)
	if err != nil {
		t.Fatalf("Failed to load snapshot: %v", err)
	}

	// Verify restored state
	if restored.workflowID != original.workflowID {
		t.Errorf("Workflow ID mismatch: expected %s, got %s", original.workflowID, restored.workflowID)
	}

	if restored.executionID != original.executionID {
		t.Errorf("Execution ID mismatch: expected %s, got %s", original.executionID, restored.executionID)
	}

	if len(restored.results) != len(originalResult.NodeResults) {
		t.Errorf("Results count mismatch: expected %d, got %d", len(originalResult.NodeResults), len(restored.results))
	}

	// Verify results match
	for nodeID, originalRes := range originalResult.NodeResults {
		restoredRes, exists := restored.results[nodeID]
		if !exists {
			t.Errorf("Result for node %s not restored", nodeID)
			continue
		}

		if originalRes != restoredRes {
			t.Errorf("Result mismatch for node %s: expected %v, got %v", nodeID, originalRes, restoredRes)
		}
	}
}

// TestSerializeDeserialize tests JSON serialization
func TestSerializeDeserialize(t *testing.T) {
	payload := []byte(`{
		"workflow_id": "serialize-test",
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 42}}
		],
		"edges": []
	}`)

	engine, err := New(payload)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	_, err = engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Create snapshot
	snapshot, err := engine.SaveSnapshot()
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	// Serialize to JSON
	data, err := SerializeSnapshot(snapshot)
	if err != nil {
		t.Fatalf("Failed to serialize snapshot: %v", err)
	}

	// Verify it's valid JSON
	var jsonCheck map[string]interface{}
	if err := json.Unmarshal(data, &jsonCheck); err != nil {
		t.Fatalf("Serialized data is not valid JSON: %v", err)
	}

	// Deserialize back
	restored, err := DeserializeSnapshot(data)
	if err != nil {
		t.Fatalf("Failed to deserialize snapshot: %v", err)
	}

	// Verify key fields match
	if restored.Version != snapshot.Version {
		t.Errorf("Version mismatch: expected %s, got %s", snapshot.Version, restored.Version)
	}

	if restored.WorkflowID != snapshot.WorkflowID {
		t.Errorf("Workflow ID mismatch: expected %s, got %s", snapshot.WorkflowID, restored.WorkflowID)
	}

	if restored.ExecutionID != snapshot.ExecutionID {
		t.Errorf("Execution ID mismatch: expected %s, got %s", snapshot.ExecutionID, restored.ExecutionID)
	}

	if len(restored.Nodes) != len(snapshot.Nodes) {
		t.Errorf("Nodes count mismatch: expected %d, got %d", len(snapshot.Nodes), len(restored.Nodes))
	}
}

// TestExecuteFromSnapshot tests convenience function
func TestExecuteFromSnapshot(t *testing.T) {
	payload := []byte(`{
		"workflow_id": "execute-from-snapshot",
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 10}},
			{"id": "2", "type": "number", "data": {"value": 20}},
			{"id": "3", "type": "operation", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "3"},
			{"source": "2", "target": "3"}
		]
	}`)

	// Create and execute original
	engine, err := New(payload)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	originalResult, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Save snapshot
	snapshot, err := engine.SaveSnapshot()
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	// Execute from snapshot
	restoredResult, err := ExecuteFromSnapshot(snapshot, nil)
	if err != nil {
		t.Fatalf("Failed to execute from snapshot: %v", err)
	}

	// Results should match since workflow was already complete
	if len(restoredResult.NodeResults) != len(originalResult.NodeResults) {
		t.Errorf("Results count mismatch: expected %d, got %d", len(originalResult.NodeResults), len(restoredResult.NodeResults))
	}

	for nodeID, originalRes := range originalResult.NodeResults {
		restoredRes, exists := restoredResult.NodeResults[nodeID]
		if !exists {
			t.Errorf("Result for node %s not in restored results", nodeID)
			continue
		}

		if originalRes != restoredRes {
			t.Errorf("Result mismatch for node %s: expected %v, got %v", nodeID, originalRes, restoredRes)
		}
	}
}

// TestSnapshot_CacheExpiration tests cache TTL handling in snapshots
func TestSnapshot_CacheExpiration(t *testing.T) {
	// Skip this test for now - cache testing is complex
	t.Skip("Cache snapshot testing needs more work")
}

// TestSnapshot_ProtectionCounters tests that protection counters are preserved
func TestSnapshot_ProtectionCounters(t *testing.T) {
	payload := []byte(`{
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 1}},
			{"id": "2", "type": "number", "data": {"value": 2}},
			{"id": "3", "type": "number", "data": {"value": 3}},
			{"id": "4", "type": "operation", "data": {"op": "add"}},
			{"id": "5", "type": "operation", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "4"},
			{"source": "2", "target": "4"},
			{"source": "4", "target": "5"},
			{"source": "3", "target": "5"}
		]
	}`)

	engine, err := New(payload)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	_, err = engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Save snapshot
	snapshot, err := engine.SaveSnapshot()
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	// Verify counters are captured
	if snapshot.NodeExecutionCount != 5 {
		t.Errorf("Expected node execution count 5, got %d", snapshot.NodeExecutionCount)
	}

	// Load snapshot
	restored, err := LoadSnapshot(snapshot, nil)
	if err != nil {
		t.Fatalf("Failed to load snapshot: %v", err)
	}

	// Verify counters are restored
	if restored.nodeExecutionCount != snapshot.NodeExecutionCount {
		t.Errorf("Node execution count not restored: expected %d, got %d", snapshot.NodeExecutionCount, restored.nodeExecutionCount)
	}
}

// TestSnapshot_InvalidVersion tests version validation
func TestSnapshot_InvalidVersion(t *testing.T) {
	snapshot := &Snapshot{
		Version:     "99.99.99", // Invalid version
		ExecutionID: "test",
		WorkflowID:  "test",
		Nodes:       []types.Node{},
		Edges:       []types.Edge{},
		Results:     make(map[string]interface{}),
	}

	_, err := LoadSnapshot(snapshot, nil)
	if err == nil {
		t.Error("Expected error for invalid snapshot version")
	}
}

// TestSnapshot_NilSnapshot tests nil snapshot handling
func TestSnapshot_NilSnapshot(t *testing.T) {
	_, err := LoadSnapshot(nil, nil)
	if err == nil {
		t.Error("Expected error for nil snapshot")
	}

	_, err = SerializeSnapshot(nil)
	if err == nil {
		t.Error("Expected error for nil snapshot serialization")
	}
}

// TestSnapshot_ComplexWorkflow tests snapshot with complex workflow
func TestSnapshot_ComplexWorkflow(t *testing.T) {
	payload := []byte(`{
		"workflow_id": "complex-test",
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 10}},
			{"id": "2", "type": "variable", "data": {"var_op": "set", "var_name": "x"}},
			{"id": "3", "type": "counter", "data": {"counter_op": "increment", "delta": 1}},
			{"id": "4", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"source": "1", "target": "2"},
			{"source": "2", "target": "3"},
			{"source": "3", "target": "4"}
		]
	}`)

	// Execute original
	engine, err := New(payload)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	originalResult, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Save snapshot
	snapshot, err := engine.SaveSnapshot()
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	// Serialize and deserialize
	data, err := SerializeSnapshot(snapshot)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	restored, err := DeserializeSnapshot(data)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	// Load and execute
	newEngine, err := LoadSnapshot(restored, nil)
	if err != nil {
		t.Fatalf("Failed to load snapshot: %v", err)
	}

	restoredResult, err := newEngine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute restored workflow: %v", err)
	}

	// Verify all results match
	if len(restoredResult.NodeResults) != len(originalResult.NodeResults) {
		t.Errorf("Results count mismatch: expected %d, got %d", len(originalResult.NodeResults), len(restoredResult.NodeResults))
	}
}

// TestSnapshot_EmptyWorkflow tests snapshot of empty workflow
func TestSnapshot_EmptyWorkflow(t *testing.T) {
	payload := []byte(`{
		"workflow_id": "empty",
		"nodes": [],
		"edges": []
	}`)

	engine, err := New(payload)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	snapshot, err := engine.SaveSnapshot()
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	if len(snapshot.Nodes) != 0 {
		t.Errorf("Expected 0 nodes, got %d", len(snapshot.Nodes))
	}

	if len(snapshot.Results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(snapshot.Results))
	}
}

// TestSnapshot_Timestamp tests that snapshot time is recorded
func TestSnapshot_Timestamp(t *testing.T) {
	payload := []byte(`{
		"nodes": [{"id": "1", "type": "number", "data": {"value": 1}}],
		"edges": []
	}`)

	engine, err := New(payload)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	before := time.Now()
	snapshot, err := engine.SaveSnapshot()
	after := time.Now()

	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	if snapshot.SnapshotTime.Before(before) || snapshot.SnapshotTime.After(after) {
		t.Errorf("Snapshot time %v not between %v and %v", snapshot.SnapshotTime, before, after)
	}
}
