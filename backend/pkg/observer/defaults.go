package observer

import (
	"context"
	"fmt"
	"log"
	"os"
)

// ============================================================================
// Default Observer Implementations
// ============================================================================

// NoOpObserver is a no-operation observer that ignores all events.
// This is useful as a default when no observer is configured.
type NoOpObserver struct{}

// OnEvent implements Observer interface (does nothing)
func (o *NoOpObserver) OnEvent(ctx context.Context, event Event) {
	// No operation
}

// ConsoleObserver is a simple observer that prints events to stdout.
// This is useful for development and debugging.
type ConsoleObserver struct {
	logger Logger
}

// NewConsoleObserver creates a new console observer with the default logger
func NewConsoleObserver() *ConsoleObserver {
	return &ConsoleObserver{
		logger: NewDefaultLogger(),
	}
}

// NewConsoleObserverWithLogger creates a new console observer with a custom logger
func NewConsoleObserverWithLogger(logger Logger) *ConsoleObserver {
	return &ConsoleObserver{
		logger: logger,
	}
}

// OnEvent implements Observer interface
func (o *ConsoleObserver) OnEvent(ctx context.Context, event Event) {
	fields := map[string]interface{}{
		"type":         event.Type,
		"status":       event.Status,
		"execution_id": event.ExecutionID,
	}

	if event.WorkflowID != "" {
		fields["workflow_id"] = event.WorkflowID
	}

	if event.NodeID != "" {
		fields["node_id"] = event.NodeID
		fields["node_type"] = event.NodeType
	}

	if event.ElapsedTime > 0 {
		fields["elapsed_time"] = event.ElapsedTime.String()
	}

	msg := fmt.Sprintf("[%s] %s", event.Type, event.Status)

	switch event.Type {
	case EventWorkflowStart:
		o.logger.Info(msg, fields)
	case EventWorkflowEnd:
		if event.Error != nil {
			fields["error"] = event.Error.Error()
			o.logger.Error(msg, fields)
		} else {
			o.logger.Info(msg, fields)
		}
	case EventNodeStart:
		o.logger.Debug(msg, fields)
	case EventNodeSuccess:
		o.logger.Debug(msg, fields)
	case EventNodeFailure:
		if event.Error != nil {
			fields["error"] = event.Error.Error()
		}
		o.logger.Warn(msg, fields)
	case EventNodeEnd:
		o.logger.Debug(msg, fields)
	default:
		o.logger.Info(msg, fields)
	}
}

// ============================================================================
// Default Logger Implementations
// ============================================================================

// NoOpLogger is a no-operation logger that ignores all log messages.
type NoOpLogger struct{}

// Debug implements Logger interface (does nothing)
func (l *NoOpLogger) Debug(msg string, fields map[string]interface{}) {}

// Info implements Logger interface (does nothing)
func (l *NoOpLogger) Info(msg string, fields map[string]interface{}) {}

// Warn implements Logger interface (does nothing)
func (l *NoOpLogger) Warn(msg string, fields map[string]interface{}) {}

// Error implements Logger interface (does nothing)
func (l *NoOpLogger) Error(msg string, fields map[string]interface{}) {}

// DefaultLogger is a simple logger that writes to stdout/stderr.
// This uses the standard library's log package.
type DefaultLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

// Debug implements Logger interface
func (l *DefaultLogger) Debug(msg string, fields map[string]interface{}) {
	l.infoLogger.Printf("[DEBUG] %s %v", msg, fields)
}

// Info implements Logger interface
func (l *DefaultLogger) Info(msg string, fields map[string]interface{}) {
	l.infoLogger.Printf("%s %v", msg, fields)
}

// Warn implements Logger interface
func (l *DefaultLogger) Warn(msg string, fields map[string]interface{}) {
	l.infoLogger.Printf("[WARN] %s %v", msg, fields)
}

// Error implements Logger interface
func (l *DefaultLogger) Error(msg string, fields map[string]interface{}) {
	l.errorLogger.Printf("%s %v", msg, fields)
}

// ============================================================================
// Observer Manager
// ============================================================================

// Manager manages multiple observers and provides a unified notification interface.
// It supports registering multiple observers and notifying them all of events asynchronously.
// Observers are executed in separate goroutines to avoid blocking the main execution flow.
type Manager struct {
	observers []Observer
}

// NewManager creates a new observer manager with no observers
func NewManager() *Manager {
	return &Manager{
		observers: []Observer{},
	}
}

// NewManagerWithObservers creates a new observer manager with initial observers
func NewManagerWithObservers(observers ...Observer) *Manager {
	return &Manager{
		observers: observers,
	}
}

// Register adds an observer to the manager
func (m *Manager) Register(observer Observer) {
	if observer != nil {
		m.observers = append(m.observers, observer)
	}
}

// Notify sends an event to all registered observers asynchronously.
// Each observer is called in a separate goroutine to prevent blocking.
// If an observer panics, it will be recovered and not affect other observers or the main execution.
func (m *Manager) Notify(ctx context.Context, event Event) {
	for _, observer := range m.observers {
		// Create a local copy to avoid closure issues
		obs := observer
		
		// Execute observer asynchronously in a goroutine
		go func() {
			// Recover from any panics in observer code
			defer func() {
				if r := recover(); r != nil {
					// Observer panicked, but we don't propagate it
					// In production, this could be logged to a system logger
					// For now, we silently recover to maintain system stability
				}
			}()
			
			// Call the observer with the event
			obs.OnEvent(ctx, event)
		}()
	}
}

// HasObservers returns true if any observers are registered
func (m *Manager) HasObservers() bool {
	return len(m.observers) > 0
}

// Count returns the number of registered observers
func (m *Manager) Count() int {
	return len(m.observers)
}
