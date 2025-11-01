package observer

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ============================================================================
// Test Observer Implementation
// ============================================================================

// TestObserver is a test observer that records all events
// It includes synchronization primitives for testing asynchronous behavior
type TestObserver struct {
	events   []Event
	mu       sync.Mutex
	wg       sync.WaitGroup
	expected int // Track expected event count
}

func NewTestObserver() *TestObserver {
	return &TestObserver{
		events:   []Event{},
		expected: 0,
	}
}

func (o *TestObserver) OnEvent(ctx context.Context, event Event) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	o.events = append(o.events, event)
	
	// Only call Done if we're expecting events
	if o.expected > 0 {
		o.wg.Done()
		o.expected--
	}
}

func (o *TestObserver) GetEvents() []Event {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.events
}

func (o *TestObserver) GetEventCount() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return len(o.events)
}

func (o *TestObserver) GetEventsByType(eventType EventType) []Event {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	filtered := []Event{}
	for _, e := range o.events {
		if e.Type == eventType {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func (o *TestObserver) Clear() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.events = []Event{}
}

// ExpectEvents prepares the observer to wait for a specific number of events
func (o *TestObserver) ExpectEvents(count int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.expected += count
	o.wg.Add(count)
}

// Wait waits for all expected events to be received
func (o *TestObserver) Wait() {
	o.wg.Wait()
}

// ============================================================================
// NoOpObserver Tests
// ============================================================================

func TestNoOpObserver(t *testing.T) {
	observer := &NoOpObserver{}
	ctx := context.Background()

	event := Event{
		Type:        EventWorkflowStart,
		Status:      StatusStarted,
		Timestamp:   time.Now(),
		ExecutionID: "test-exec-123",
	}

	// Should not panic
	observer.OnEvent(ctx, event)
}

// ============================================================================
// ConsoleObserver Tests
// ============================================================================

func TestConsoleObserver(t *testing.T) {
	observer := NewConsoleObserver()

	if observer == nil {
		t.Fatal("NewConsoleObserver returned nil")
	}

	ctx := context.Background()
	event := Event{
		Type:        EventWorkflowStart,
		Status:      StatusStarted,
		Timestamp:   time.Now(),
		ExecutionID: "test-exec-123",
		WorkflowID:  "test-workflow-456",
	}

	// Should not panic
	observer.OnEvent(ctx, event)
}

func TestConsoleObserverWithCustomLogger(t *testing.T) {
	logger := NewDefaultLogger()
	observer := NewConsoleObserverWithLogger(logger)

	if observer == nil {
		t.Fatal("NewConsoleObserverWithLogger returned nil")
	}

	ctx := context.Background()

	// Test different event types
	events := []Event{
		{
			Type:        EventWorkflowStart,
			Status:      StatusStarted,
			Timestamp:   time.Now(),
			ExecutionID: "test-exec-123",
		},
		{
			Type:        EventNodeStart,
			Status:      StatusStarted,
			Timestamp:   time.Now(),
			ExecutionID: "test-exec-123",
			NodeID:      "node-1",
			NodeType:    types.NodeTypeNumber,
		},
		{
			Type:        EventNodeSuccess,
			Status:      StatusSuccess,
			Timestamp:   time.Now(),
			ExecutionID: "test-exec-123",
			NodeID:      "node-1",
			ElapsedTime: 100 * time.Millisecond,
		},
		{
			Type:        EventWorkflowEnd,
			Status:      StatusSuccess,
			Timestamp:   time.Now(),
			ExecutionID: "test-exec-123",
			ElapsedTime: 500 * time.Millisecond,
		},
	}

	// Should not panic
	for _, event := range events {
		observer.OnEvent(ctx, event)
	}
}

// ============================================================================
// NoOpLogger Tests
// ============================================================================

func TestNoOpLogger(t *testing.T) {
	logger := &NoOpLogger{}
	fields := map[string]interface{}{
		"key": "value",
	}

	// Should not panic
	logger.Debug("debug message", fields)
	logger.Info("info message", fields)
	logger.Warn("warn message", fields)
	logger.Error("error message", fields)
}

// ============================================================================
// DefaultLogger Tests
// ============================================================================

func TestDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()

	if logger == nil {
		t.Fatal("NewDefaultLogger returned nil")
	}

	fields := map[string]interface{}{
		"execution_id": "test-123",
		"node_id":      "node-1",
	}

	// Should not panic
	logger.Debug("debug message", fields)
	logger.Info("info message", fields)
	logger.Warn("warn message", fields)
	logger.Error("error message", fields)
}

// ============================================================================
// Observer Manager Tests
// ============================================================================

func TestNewManager(t *testing.T) {
	mgr := NewManager()

	if mgr == nil {
		t.Fatal("NewManager returned nil")
	}

	if mgr.Count() != 0 {
		t.Errorf("Expected 0 observers, got %d", mgr.Count())
	}

	if mgr.HasObservers() {
		t.Error("Expected HasObservers to return false")
	}
}

func TestManagerRegister(t *testing.T) {
	mgr := NewManager()
	obs1 := NewTestObserver()
	obs2 := NewTestObserver()

	mgr.Register(obs1)
	if mgr.Count() != 1 {
		t.Errorf("Expected 1 observer, got %d", mgr.Count())
	}

	mgr.Register(obs2)
	if mgr.Count() != 2 {
		t.Errorf("Expected 2 observers, got %d", mgr.Count())
	}

	if !mgr.HasObservers() {
		t.Error("Expected HasObservers to return true")
	}
}

func TestManagerRegisterNil(t *testing.T) {
	mgr := NewManager()
	mgr.Register(nil)

	if mgr.Count() != 0 {
		t.Errorf("Expected 0 observers after registering nil, got %d", mgr.Count())
	}
}

func TestManagerNotify(t *testing.T) {
	mgr := NewManager()
	obs1 := NewTestObserver()
	obs2 := NewTestObserver()

	mgr.Register(obs1)
	mgr.Register(obs2)

	ctx := context.Background()
	event := Event{
		Type:        EventWorkflowStart,
		Status:      StatusStarted,
		Timestamp:   time.Now(),
		ExecutionID: "test-exec-123",
	}

	// Prepare observers to wait for events (asynchronous execution)
	obs1.ExpectEvents(1)
	obs2.ExpectEvents(1)

	mgr.Notify(ctx, event)

	// Wait for async observers to complete
	obs1.Wait()
	obs2.Wait()

	if obs1.GetEventCount() != 1 {
		t.Errorf("Observer 1 expected 1 event, got %d", obs1.GetEventCount())
	}

	if obs2.GetEventCount() != 1 {
		t.Errorf("Observer 2 expected 1 event, got %d", obs2.GetEventCount())
	}

	// Verify event content
	events1 := obs1.GetEvents()
	if events1[0].Type != EventWorkflowStart {
		t.Errorf("Expected event type %s, got %s", EventWorkflowStart, events1[0].Type)
	}
}

func TestManagerNotifyMultipleEvents(t *testing.T) {
	mgr := NewManager()
	obs := NewTestObserver()
	mgr.Register(obs)

	ctx := context.Background()

	events := []Event{
		{Type: EventWorkflowStart, Status: StatusStarted, Timestamp: time.Now(), ExecutionID: "exec-1"},
		{Type: EventNodeStart, Status: StatusStarted, Timestamp: time.Now(), ExecutionID: "exec-1", NodeID: "node-1"},
		{Type: EventNodeSuccess, Status: StatusSuccess, Timestamp: time.Now(), ExecutionID: "exec-1", NodeID: "node-1"},
		{Type: EventWorkflowEnd, Status: StatusSuccess, Timestamp: time.Now(), ExecutionID: "exec-1"},
	}

	// Prepare observer to wait for all events
	obs.ExpectEvents(len(events))

	for _, event := range events {
		mgr.Notify(ctx, event)
	}

	// Wait for async observers to complete
	obs.Wait()

	if obs.GetEventCount() != 4 {
		t.Errorf("Expected 4 events, got %d", obs.GetEventCount())
	}

	// Verify event types
	workflowStarts := obs.GetEventsByType(EventWorkflowStart)
	if len(workflowStarts) != 1 {
		t.Errorf("Expected 1 workflow start event, got %d", len(workflowStarts))
	}

	nodeSuccesses := obs.GetEventsByType(EventNodeSuccess)
	if len(nodeSuccesses) != 1 {
		t.Errorf("Expected 1 node success event, got %d", len(nodeSuccesses))
	}
}

func TestNewManagerWithObservers(t *testing.T) {
	obs1 := NewTestObserver()
	obs2 := NewTestObserver()

	mgr := NewManagerWithObservers(obs1, obs2)

	if mgr.Count() != 2 {
		t.Errorf("Expected 2 observers, got %d", mgr.Count())
	}

	ctx := context.Background()
	event := Event{
		Type:        EventWorkflowStart,
		Status:      StatusStarted,
		Timestamp:   time.Now(),
		ExecutionID: "test-exec-123",
	}

	// Prepare observers to wait for events
	obs1.ExpectEvents(1)
	obs2.ExpectEvents(1)

	mgr.Notify(ctx, event)

	// Wait for async observers to complete
	obs1.Wait()
	obs2.Wait()

	if obs1.GetEventCount() != 1 {
		t.Errorf("Observer 1 expected 1 event, got %d", obs1.GetEventCount())
	}

	if obs2.GetEventCount() != 1 {
		t.Errorf("Observer 2 expected 1 event, got %d", obs2.GetEventCount())
	}
}

// ============================================================================
// Event Tests
// ============================================================================

func TestEventStructure(t *testing.T) {
	now := time.Now()
	event := Event{
		Type:        EventNodeSuccess,
		Status:      StatusSuccess,
		Timestamp:   now,
		ExecutionID: "exec-123",
		WorkflowID:  "workflow-456",
		NodeID:      "node-789",
		NodeType:    types.NodeTypeOperation,
		StartTime:   now.Add(-100 * time.Millisecond),
		ElapsedTime: 100 * time.Millisecond,
		Result:      42,
		Error:       nil,
		Metadata: map[string]interface{}{
			"custom": "data",
		},
	}

	if event.Type != EventNodeSuccess {
		t.Errorf("Expected type %s, got %s", EventNodeSuccess, event.Type)
	}

	if event.Status != StatusSuccess {
		t.Errorf("Expected status %s, got %s", StatusSuccess, event.Status)
	}

	if event.ExecutionID != "exec-123" {
		t.Errorf("Expected execution ID 'exec-123', got '%s'", event.ExecutionID)
	}

	if event.WorkflowID != "workflow-456" {
		t.Errorf("Expected workflow ID 'workflow-456', got '%s'", event.WorkflowID)
	}

	if event.NodeID != "node-789" {
		t.Errorf("Expected node ID 'node-789', got '%s'", event.NodeID)
	}

	if event.Result != 42 {
		t.Errorf("Expected result 42, got %v", event.Result)
	}

	if event.Metadata["custom"] != "data" {
		t.Errorf("Expected metadata custom='data', got %v", event.Metadata["custom"])
	}
}

// ============================================================================
// Asynchronous Execution Tests
// ============================================================================

func TestObserverAsynchronousExecution(t *testing.T) {
	mgr := NewManager()
	
	// Create an observer that sleeps for a bit
	slowObserver := NewTestObserver()
	mgr.Register(slowObserver)

	ctx := context.Background()
	event := Event{
		Type:        EventWorkflowStart,
		Status:      StatusStarted,
		Timestamp:   time.Now(),
		ExecutionID: "test-exec-123",
	}

	// Prepare observer
	slowObserver.ExpectEvents(1)

	// Measure time for notification (should be nearly instant)
	start := time.Now()
	mgr.Notify(ctx, event)
	elapsed := time.Since(start)

	// Notification should return immediately (asynchronous)
	if elapsed > 10*time.Millisecond {
		t.Errorf("Notify blocked for %v, expected to be asynchronous", elapsed)
	}

	// Wait for observer to finish
	slowObserver.Wait()

	// Verify event was received
	if slowObserver.GetEventCount() != 1 {
		t.Errorf("Expected 1 event, got %d", slowObserver.GetEventCount())
	}
}

func TestObserverPanicRecovery(t *testing.T) {
	mgr := NewManager()
	
	// Create a panicking observer
	panicObserver := &PanicObserver{}
	normalObserver := NewTestObserver()
	
	mgr.Register(panicObserver)
	mgr.Register(normalObserver)

	ctx := context.Background()
	event := Event{
		Type:        EventWorkflowStart,
		Status:      StatusStarted,
		Timestamp:   time.Now(),
		ExecutionID: "test-exec-123",
	}

	// Prepare normal observer
	normalObserver.ExpectEvents(1)

	// Should not panic even though one observer panics
	mgr.Notify(ctx, event)

	// Wait for normal observer
	normalObserver.Wait()

	// Normal observer should still receive the event
	if normalObserver.GetEventCount() != 1 {
		t.Errorf("Expected 1 event in normal observer, got %d", normalObserver.GetEventCount())
	}
}

// PanicObserver always panics when OnEvent is called
type PanicObserver struct{}

func (o *PanicObserver) OnEvent(ctx context.Context, event Event) {
	panic("observer panic test")
}

func TestMultipleObserversParallelExecution(t *testing.T) {
	mgr := NewManager()
	
	// Create multiple observers
	observers := make([]*TestObserver, 10)
	for i := 0; i < 10; i++ {
		observers[i] = NewTestObserver()
		mgr.Register(observers[i])
	}

	ctx := context.Background()
	event := Event{
		Type:        EventWorkflowStart,
		Status:      StatusStarted,
		Timestamp:   time.Now(),
		ExecutionID: "test-exec-123",
	}

	// Prepare all observers
	for _, obs := range observers {
		obs.ExpectEvents(1)
	}

	// Notify all observers
	start := time.Now()
	mgr.Notify(ctx, event)
	elapsed := time.Since(start)

	// Should return immediately even with many observers
	if elapsed > 10*time.Millisecond {
		t.Errorf("Notify with 10 observers blocked for %v, expected to be asynchronous", elapsed)
	}

	// Wait for all observers
	for _, obs := range observers {
		obs.Wait()
	}

	// Verify all observers received the event
	for i, obs := range observers {
		if obs.GetEventCount() != 1 {
			t.Errorf("Observer %d expected 1 event, got %d", i, obs.GetEventCount())
		}
	}
}
