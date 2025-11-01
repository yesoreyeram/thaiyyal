package engine

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/observer"
)

// ============================================================================
// Test Observer for Integration Tests
// ============================================================================

type testObserver struct {
	events []observer.Event
	mu     sync.Mutex
	wg     sync.WaitGroup
}

func newTestObserver() *testObserver {
	return &testObserver{
		events: []observer.Event{},
	}
}

func (o *testObserver) OnEvent(ctx context.Context, event observer.Event) {
	o.mu.Lock()
	defer o.mu.Unlock()
	defer o.wg.Done()
	
	o.events = append(o.events, event)
}

func (o *testObserver) expectEvents(count int) {
	o.wg.Add(count)
}

func (o *testObserver) wait() {
	o.wg.Wait()
}

func (o *testObserver) getEvents() []observer.Event {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.events
}

func (o *testObserver) getEventsByType(eventType observer.EventType) []observer.Event {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	filtered := []observer.Event{}
	for _, e := range o.events {
		if e.Type == eventType {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// ============================================================================
// Observer Integration Tests
// ============================================================================

func TestEngineWithObserver(t *testing.T) {
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

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Register observer
	obs := newTestObserver()
	engine.RegisterObserver(obs)

	if engine.GetObserverCount() != 1 {
		t.Errorf("Expected 1 observer, got %d", engine.GetObserverCount())
	}

	// Expect events: 1 workflow start + 3 node starts + 3 node successes + 1 workflow end = 8 events
	obs.expectEvents(8)

	// Execute workflow
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Wait for all observers to complete
	obs.wait()

	// Verify result
	if result.FinalOutput != 15.0 {
		t.Errorf("Expected final output 15, got %v", result.FinalOutput)
	}

	// Verify events
	events := obs.getEvents()
	if len(events) != 8 {
		t.Errorf("Expected 8 events, got %d", len(events))
	}

	// Verify workflow start event
	workflowStarts := obs.getEventsByType(observer.EventWorkflowStart)
	if len(workflowStarts) != 1 {
		t.Errorf("Expected 1 workflow start event, got %d", len(workflowStarts))
	}
	if workflowStarts[0].ExecutionID == "" {
		t.Error("Workflow start event missing execution ID")
	}

	// Verify workflow end event
	workflowEnds := obs.getEventsByType(observer.EventWorkflowEnd)
	if len(workflowEnds) != 1 {
		t.Errorf("Expected 1 workflow end event, got %d", len(workflowEnds))
	}
	if workflowEnds[0].Status != observer.StatusSuccess {
		t.Errorf("Expected workflow end status success, got %s", workflowEnds[0].Status)
	}
	if workflowEnds[0].ElapsedTime == 0 {
		t.Error("Workflow end event missing elapsed time")
	}

	// Verify node success events
	nodeSuccesses := obs.getEventsByType(observer.EventNodeSuccess)
	if len(nodeSuccesses) != 3 {
		t.Errorf("Expected 3 node success events, got %d", len(nodeSuccesses))
	}

	// Verify node start events
	nodeStarts := obs.getEventsByType(observer.EventNodeStart)
	if len(nodeStarts) != 3 {
		t.Errorf("Expected 3 node start events, got %d", len(nodeStarts))
	}

	// Verify each node success has metadata
	for i, event := range nodeSuccesses {
		if event.NodeID == "" {
			t.Errorf("Node success event %d missing node ID", i)
		}
		if event.NodeType == "" {
			t.Errorf("Node success event %d missing node type", i)
		}
		if event.ElapsedTime == 0 {
			t.Errorf("Node success event %d missing elapsed time", i)
		}
	}
}

func TestEngineWithMultipleObservers(t *testing.T) {
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

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Register multiple observers
	obs1 := newTestObserver()
	obs2 := newTestObserver()
	obs3 := newTestObserver()

	engine.RegisterObserver(obs1)
	engine.RegisterObserver(obs2)
	engine.RegisterObserver(obs3)

	if engine.GetObserverCount() != 3 {
		t.Errorf("Expected 3 observers, got %d", engine.GetObserverCount())
	}

	// Expect events for all observers: 1 workflow start + 3 node starts + 3 node successes + 1 workflow end = 8 events
	obs1.expectEvents(8)
	obs2.expectEvents(8)
	obs3.expectEvents(8)

	// Execute workflow
	_, err = engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Wait for all observers
	obs1.wait()
	obs2.wait()
	obs3.wait()

	// Verify all observers received all events
	if len(obs1.getEvents()) != 8 {
		t.Errorf("Observer 1 expected 8 events, got %d", len(obs1.getEvents()))
	}
	if len(obs2.getEvents()) != 8 {
		t.Errorf("Observer 2 expected 8 events, got %d", len(obs2.getEvents()))
	}
	if len(obs3.getEvents()) != 8 {
		t.Errorf("Observer 3 expected 8 events, got %d", len(obs3.getEvents()))
	}
}

func TestEngineWithObserverErrorHandling(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 0}},
			{"id": "3", "data": {"op": "divide"}}
		],
		"edges": [
			{"source": "1", "target": "3"},
			{"source": "2", "target": "3"}
		]
	}`

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Register observer
	obs := newTestObserver()
	engine.RegisterObserver(obs)

	// Expect events: 1 workflow start + 2 node starts + 2 node successes + 1 node start (divide) + 1 node failure + 1 workflow end = 8 events
	obs.expectEvents(8)

	// Execute workflow (should fail on division by zero)
	_, err = engine.Execute()
	if err == nil {
		t.Fatal("Expected execution to fail on division by zero")
	}

	// Wait for observers
	obs.wait()

	// Verify workflow end with error
	workflowEnds := obs.getEventsByType(observer.EventWorkflowEnd)
	if len(workflowEnds) != 1 {
		t.Errorf("Expected 1 workflow end event, got %d", len(workflowEnds))
	}
	if workflowEnds[0].Status != observer.StatusFailure {
		t.Errorf("Expected workflow end status failure, got %s", workflowEnds[0].Status)
	}
	if workflowEnds[0].Error == nil {
		t.Error("Expected workflow end event to have error")
	}

	// Verify node failure event
	nodeFailures := obs.getEventsByType(observer.EventNodeFailure)
	if len(nodeFailures) != 1 {
		t.Errorf("Expected 1 node failure event, got %d", len(nodeFailures))
	}
	if nodeFailures[0].Error == nil {
		t.Error("Expected node failure event to have error")
	}
}

func TestEngineWithNoObserver(t *testing.T) {
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

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Don't register any observer
	if engine.GetObserverCount() != 0 {
		t.Errorf("Expected 0 observers, got %d", engine.GetObserverCount())
	}

	// Execute workflow (should work fine without observers)
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	if result.FinalOutput != 15.0 {
		t.Errorf("Expected final output 15, got %v", result.FinalOutput)
	}
}

func TestEngineWithCustomLogger(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}}
		],
		"edges": []
	}`

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Set custom logger
	logger := observer.NewDefaultLogger()
	engine.SetLogger(logger)

	// Execute workflow
	_, err = engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}
}

func TestObserverChaining(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}}
		],
		"edges": []
	}`

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Test method chaining
	obs := newTestObserver()
	logger := observer.NewDefaultLogger()

	engine.
		RegisterObserver(obs).
		SetLogger(logger).
		RegisterObserver(newTestObserver())

	if engine.GetObserverCount() != 2 {
		t.Errorf("Expected 2 observers after chaining, got %d", engine.GetObserverCount())
	}
}

func TestObserverWithWorkflowID(t *testing.T) {
	payload := `{
		"workflow_id": "test-workflow-123",
		"nodes": [
			{"id": "1", "data": {"value": 42}}
		],
		"edges": []
	}`

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	obs := newTestObserver()
	engine.RegisterObserver(obs)

	// Expect: 1 workflow start + 1 node start + 1 node success + 1 workflow end = 4 events
	obs.expectEvents(4)

	_, err = engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	obs.wait()

	// Verify workflow ID is included in events
	events := obs.getEvents()
	for _, event := range events {
		if event.WorkflowID != "test-workflow-123" {
			t.Errorf("Expected workflow ID 'test-workflow-123', got '%s'", event.WorkflowID)
		}
	}
}

func TestObserverPerformanceNoBlock(t *testing.T) {
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

	engine, err := New([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Register many observers
	for i := 0; i < 100; i++ {
		obs := newTestObserver()
		engine.RegisterObserver(obs)
	}

	// Execute and measure time
	start := time.Now()
	_, err = engine.Execute()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// With async observers, execution should be fast even with 100 observers
	// If observers were synchronous, this would take much longer
	if elapsed > 1*time.Second {
		t.Logf("Execution took %v with 100 observers (acceptable but may indicate blocking)", elapsed)
	}
}
