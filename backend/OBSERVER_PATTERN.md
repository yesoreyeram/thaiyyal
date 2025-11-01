# Observer Pattern Documentation

The Thaiyyal workflow engine supports the Observer pattern, allowing library consumers to track and monitor workflow execution behavior in real-time. Observers are executed **asynchronously** and do not block workflow execution.

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Quick Start](#quick-start)
- [Observer Interface](#observer-interface)
- [Event Types](#event-types)
- [Built-in Observers](#built-in-observers)
- [Custom Observers](#custom-observers)
- [Custom Loggers](#custom-loggers)
- [Multiple Observers](#multiple-observers)
- [Asynchronous Execution](#asynchronous-execution)
- [Best Practices](#best-practices)
- [Complete Examples](#complete-examples)

## Overview

The observer pattern enables you to:

- **Monitor** workflow execution in real-time
- **Track** node-level execution metrics (timing, status, errors)
- **Integrate** with external monitoring systems (DataDog, Prometheus, etc.)
- **Debug** workflow issues with detailed execution logs
- **Audit** workflow execution for compliance and analysis

## Key Features

✅ **Asynchronous Execution** - Observers run in separate goroutines and never block workflow execution  
✅ **Multiple Observers** - Register as many observers as needed  
✅ **Type-Safe Events** - Strongly-typed event structures with complete metadata  
✅ **Built-in Implementations** - Console observer and default logger included  
✅ **Custom Logger Support** - Integrate with your existing logging infrastructure  
✅ **Panic Recovery** - Observer panics are recovered and don't affect workflow execution  
✅ **Zero Configuration** - Works out of the box with sensible defaults  

## Quick Start

### Basic Observer Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yesoreyeram/thaiyyal/backend"
)

func main() {
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
    
    // Create engine
    engine, err := workflow.NewEngine([]byte(payload))
    if err != nil {
        log.Fatal(err)
    }
    
    // Register console observer for debugging
    observer := workflow.NewConsoleObserver()
    engine.RegisterObserver(observer)
    
    // Execute workflow
    result, err := engine.Execute()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %v\n", result.FinalOutput)
}
```

### Output

```
[INFO] 2025/11/01 07:14:49 [workflow_start] started map[execution_id:a1b2c3d4e5f6g7h8 workflow_id: status:started type:workflow_start]
[INFO] 2025/11/01 07:14:49 [DEBUG] [node_start] started map[execution_id:a1b2c3d4e5f6g7h8 node_id:1 node_type:number status:started type:node_start]
[INFO] 2025/11/01 07:14:49 [DEBUG] [node_success] success map[elapsed_time:245µs execution_id:a1b2c3d4e5f6g7h8 node_id:1 node_type:number status:success type:node_success]
[INFO] 2025/11/01 07:14:49 [DEBUG] [node_start] started map[execution_id:a1b2c3d4e5f6g7h8 node_id:2 node_type:number status:started type:node_start]
[INFO] 2025/11/01 07:14:49 [DEBUG] [node_success] success map[elapsed_time:189µs execution_id:a1b2c3d4e5f6g7h8 node_id:2 node_type:number status:success type:node_success]
[INFO] 2025/11/01 07:14:49 [DEBUG] [node_start] started map[execution_id:a1b2c3d4e5f6g7h8 node_id:3 node_type:operation status:started type:node_start]
[INFO] 2025/11/01 07:14:49 [DEBUG] [node_success] success map[elapsed_time:312µs execution_id:a1b2c3d4e5f6g7h8 node_id:3 node_type:operation status:success type:node_success]
[INFO] 2025/11/01 07:14:49 [workflow_end] success map[elapsed_time:1.567ms execution_id:a1b2c3d4e5f6g7h8 status:success type:workflow_end]
Result: 15
```

## Observer Interface

All observers must implement the `Observer` interface:

```go
type Observer interface {
    // OnEvent is called when an execution event occurs
    OnEvent(ctx context.Context, event Event)
}
```

### Event Structure

```go
type Event struct {
    // Event identification
    Type      EventType       // Type of event (workflow_start, node_success, etc.)
    Status    ExecutionStatus // Status (started, success, failure)
    Timestamp time.Time       // When the event occurred
    
    // Execution context
    ExecutionID string // Unique ID for this execution
    WorkflowID  string // Workflow definition ID (if provided)
    
    // Node-specific data (empty for workflow-level events)
    NodeID   string   // Node identifier
    NodeType NodeType // Type of node (number, operation, etc.)
    
    // Timing information
    StartTime   time.Time     // When execution started
    ElapsedTime time.Duration // How long execution took
    
    // Execution results
    Result interface{} // Node or workflow result
    Error  error       // Error if execution failed
    
    // Additional metadata
    Metadata map[string]interface{} // Custom metadata
}
```

## Event Types

The following event types are available:

```go
const (
    EventWorkflowStart EventType = "workflow_start" // Workflow execution started
    EventWorkflowEnd   EventType = "workflow_end"   // Workflow execution completed
    
    EventNodeStart   EventType = "node_start"   // Node execution started
    EventNodeEnd     EventType = "node_end"     // Node execution ended
    EventNodeSuccess EventType = "node_success" // Node executed successfully
    EventNodeFailure EventType = "node_failure" // Node execution failed
)
```

### Execution Status

```go
const (
    StatusStarted   ExecutionStatus = "started"   // Execution has started
    StatusSuccess   ExecutionStatus = "success"   // Execution succeeded
    StatusFailure   ExecutionStatus = "failure"   // Execution failed
    StatusCompleted ExecutionStatus = "completed" // Execution completed
)
```

## Built-in Observers

### NoOpObserver

A no-operation observer that ignores all events. Useful when you want to disable observation:

```go
observer := &workflow.NoOpObserver{}
engine.RegisterObserver(observer)
```

### ConsoleObserver

Prints events to stdout/stderr using the default logger:

```go
observer := workflow.NewConsoleObserver()
engine.RegisterObserver(observer)
```

#### With Custom Logger

```go
logger := workflow.NewDefaultLogger()
observer := workflow.NewConsoleObserverWithLogger(logger)
engine.RegisterObserver(observer)
```

## Custom Observers

Create custom observers for integration with monitoring systems:

### Example: Metrics Observer

```go
package main

import (
    "context"
    "sync"
    "time"
    
    "github.com/yesoreyeram/thaiyyal/backend"
)

// MetricsObserver collects execution metrics
type MetricsObserver struct {
    mu                sync.Mutex
    workflowCount     int
    nodeCount         int
    failureCount      int
    totalElapsedTime  time.Duration
}

func NewMetricsObserver() *MetricsObserver {
    return &MetricsObserver{}
}

func (o *MetricsObserver) OnEvent(ctx context.Context, event workflow.Event) {
    o.mu.Lock()
    defer o.mu.Unlock()
    
    switch event.Type {
    case workflow.EventWorkflowEnd:
        o.workflowCount++
        o.totalElapsedTime += event.ElapsedTime
        if event.Status == workflow.StatusFailure {
            o.failureCount++
        }
    case workflow.EventNodeSuccess:
        o.nodeCount++
    }
}

func (o *MetricsObserver) GetMetrics() map[string]interface{} {
    o.mu.Lock()
    defer o.mu.Unlock()
    
    return map[string]interface{}{
        "workflows_executed": o.workflowCount,
        "nodes_executed":     o.nodeCount,
        "failures":           o.failureCount,
        "avg_elapsed_time":   o.totalElapsedTime / time.Duration(o.workflowCount),
    }
}

func main() {
    // Create and register metrics observer
    metrics := NewMetricsObserver()
    
    engine, _ := workflow.NewEngine([]byte(payload))
    engine.RegisterObserver(metrics)
    
    // Execute workflow
    engine.Execute()
    
    // Get metrics
    fmt.Printf("Metrics: %+v\n", metrics.GetMetrics())
}
```

### Example: Distributed Tracing Observer

```go
package main

import (
    "context"
    
    "github.com/yesoreyeram/thaiyyal/backend"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

// TracingObserver integrates with OpenTelemetry
type TracingObserver struct {
    tracer trace.Tracer
    spans  map[string]trace.Span
    mu     sync.Mutex
}

func NewTracingObserver() *TracingObserver {
    return &TracingObserver{
        tracer: otel.Tracer("workflow-engine"),
        spans:  make(map[string]trace.Span),
    }
}

func (o *TracingObserver) OnEvent(ctx context.Context, event workflow.Event) {
    o.mu.Lock()
    defer o.mu.Unlock()
    
    switch event.Type {
    case workflow.EventWorkflowStart:
        ctx, span := o.tracer.Start(ctx, "workflow.execute")
        span.SetAttributes(
            attribute.String("execution_id", event.ExecutionID),
            attribute.String("workflow_id", event.WorkflowID),
        )
        o.spans[event.ExecutionID] = span
        
    case workflow.EventWorkflowEnd:
        if span, ok := o.spans[event.ExecutionID]; ok {
            if event.Error != nil {
                span.RecordError(event.Error)
            }
            span.End()
            delete(o.spans, event.ExecutionID)
        }
        
    case workflow.EventNodeStart:
        ctx, span := o.tracer.Start(ctx, "node.execute")
        span.SetAttributes(
            attribute.String("node_id", event.NodeID),
            attribute.String("node_type", string(event.NodeType)),
        )
        o.spans[event.NodeID] = span
        
    case workflow.EventNodeSuccess, workflow.EventNodeFailure:
        if span, ok := o.spans[event.NodeID]; ok {
            if event.Error != nil {
                span.RecordError(event.Error)
            }
            span.End()
            delete(o.spans, event.NodeID)
        }
    }
}
```

### Example: Database Audit Logger

```go
package main

import (
    "context"
    "database/sql"
    "encoding/json"
    
    "github.com/yesoreyeram/thaiyyal/backend"
)

// AuditObserver logs all workflow executions to a database
type AuditObserver struct {
    db *sql.DB
}

func NewAuditObserver(db *sql.DB) *AuditObserver {
    return &AuditObserver{db: db}
}

func (o *AuditObserver) OnEvent(ctx context.Context, event workflow.Event) {
    // Only log workflow-level events
    if event.Type != workflow.EventWorkflowEnd {
        return
    }
    
    // Serialize event to JSON
    eventJSON, err := json.Marshal(event)
    if err != nil {
        return
    }
    
    // Insert into database
    _, err = o.db.ExecContext(ctx,
        `INSERT INTO workflow_audit_log (execution_id, workflow_id, status, elapsed_time, event_data, created_at)
         VALUES ($1, $2, $3, $4, $5, $6)`,
        event.ExecutionID,
        event.WorkflowID,
        event.Status,
        event.ElapsedTime.Milliseconds(),
        eventJSON,
        event.Timestamp,
    )
}
```

## Custom Loggers

Integrate with your existing logging infrastructure:

```go
type Logger interface {
    Debug(msg string, fields map[string]interface{})
    Info(msg string, fields map[string]interface{})
    Warn(msg string, fields map[string]interface{})
    Error(msg string, fields map[string]interface{})
}
```

### Example: Logrus Integration

```go
package main

import (
    "github.com/sirupsen/logrus"
    "github.com/yesoreyeram/thaiyyal/backend"
)

// LogrusAdapter adapts logrus to the Logger interface
type LogrusAdapter struct {
    logger *logrus.Logger
}

func NewLogrusAdapter(logger *logrus.Logger) *LogrusAdapter {
    return &LogrusAdapter{logger: logger}
}

func (l *LogrusAdapter) Debug(msg string, fields map[string]interface{}) {
    l.logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l *LogrusAdapter) Info(msg string, fields map[string]interface{}) {
    l.logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l *LogrusAdapter) Warn(msg string, fields map[string]interface{}) {
    l.logger.WithFields(logrus.Fields(fields)).Warn(msg)
}

func (l *LogrusAdapter) Error(msg string, fields map[string]interface{}) {
    l.logger.WithFields(logrus.Fields(fields)).Error(msg)
}

func main() {
    // Setup logrus
    logrusLogger := logrus.New()
    logrusLogger.SetFormatter(&logrus.JSONFormatter{})
    
    // Create adapter
    logger := NewLogrusAdapter(logrusLogger)
    
    // Use with console observer
    observer := workflow.NewConsoleObserverWithLogger(logger)
    
    engine, _ := workflow.NewEngine([]byte(payload))
    engine.SetLogger(logger).RegisterObserver(observer)
    
    engine.Execute()
}
```

### Example: Zap Integration

```go
package main

import (
    "go.uber.org/zap"
    "github.com/yesoreyeram/thaiyyal/backend"
)

// ZapAdapter adapts zap logger to the Logger interface
type ZapAdapter struct {
    logger *zap.Logger
}

func NewZapAdapter(logger *zap.Logger) *ZapAdapter {
    return &ZapAdapter{logger: logger}
}

func (l *ZapAdapter) Debug(msg string, fields map[string]interface{}) {
    l.logger.Debug(msg, l.toZapFields(fields)...)
}

func (l *ZapAdapter) Info(msg string, fields map[string]interface{}) {
    l.logger.Info(msg, l.toZapFields(fields)...)
}

func (l *ZapAdapter) Warn(msg string, fields map[string]interface{}) {
    l.logger.Warn(msg, l.toZapFields(fields)...)
}

func (l *ZapAdapter) Error(msg string, fields map[string]interface{}) {
    l.logger.Error(msg, l.toZapFields(fields)...)
}

func (l *ZapAdapter) toZapFields(fields map[string]interface{}) []zap.Field {
    zapFields := make([]zap.Field, 0, len(fields))
    for k, v := range fields {
        zapFields = append(zapFields, zap.Any(k, v))
    }
    return zapFields
}
```

## Multiple Observers

Register as many observers as needed:

```go
engine, _ := workflow.NewEngine([]byte(payload))

// Register multiple observers
engine.
    RegisterObserver(NewMetricsObserver()).
    RegisterObserver(NewTracingObserver()).
    RegisterObserver(NewAuditObserver(db)).
    RegisterObserver(workflow.NewConsoleObserver())

// All observers will receive events
engine.Execute()
```

## Asynchronous Execution

**Important**: Observers are executed asynchronously in separate goroutines. This means:

✅ **No Performance Impact** - Observers never block workflow execution  
✅ **Parallel Execution** - Multiple observers run concurrently  
✅ **Panic Recovery** - Observer panics are recovered and don't crash the workflow  
✅ **Fire and Forget** - Events are sent to observers without waiting for completion  

### Performance Guarantee

```go
// Even with 100 observers, execution is not blocked
for i := 0; i < 100; i++ {
    engine.RegisterObserver(NewMetricsObserver())
}

start := time.Now()
engine.Execute()
elapsed := time.Since(start)
// elapsed will be the same as without observers
```

### Panic Safety

```go
// If an observer panics, it won't affect the workflow or other observers
type PanickyObserver struct{}

func (o *PanickyObserver) OnEvent(ctx context.Context, event workflow.Event) {
    panic("something went wrong!")  // Recovered automatically
}

engine.RegisterObserver(&PanickyObserver())
engine.RegisterObserver(workflow.NewConsoleObserver())

// Both observers registered, console observer still works
result, _ := engine.Execute()  // Succeeds despite panic
```

## Best Practices

### 1. Use Appropriate Event Types

Only handle events you care about:

```go
func (o *MyObserver) OnEvent(ctx context.Context, event workflow.Event) {
    switch event.Type {
    case workflow.EventWorkflowEnd:
        // Only care about workflow completion
        o.recordMetrics(event)
    }
}
```

### 2. Keep Observers Fast

Although observers are async, keep them lightweight:

```go
func (o *MyObserver) OnEvent(ctx context.Context, event workflow.Event) {
    // Bad: Heavy computation in observer
    // result := expensiveOperation()
    
    // Good: Queue for async processing
    o.eventQueue <- event
}
```

### 3. Handle Context Cancellation

Respect context cancellation in long-running observers:

```go
func (o *MyObserver) OnEvent(ctx context.Context, event workflow.Event) {
    select {
    case <-ctx.Done():
        return  // Context cancelled
    default:
        o.processEvent(event)
    }
}
```

### 4. Use Structured Logging

Leverage the event metadata for structured logging:

```go
func (o *MyObserver) OnEvent(ctx context.Context, event workflow.Event) {
    logger.WithFields(logrus.Fields{
        "execution_id": event.ExecutionID,
        "workflow_id":  event.WorkflowID,
        "node_id":      event.NodeID,
        "node_type":    event.NodeType,
        "elapsed_time": event.ElapsedTime,
    }).Info("Node executed")
}
```

### 5. Thread-Safe Observer State

Protect observer state with mutexes:

```go
type MetricsObserver struct {
    mu     sync.Mutex
    counts map[string]int
}

func (o *MetricsObserver) OnEvent(ctx context.Context, event workflow.Event) {
    o.mu.Lock()
    defer o.mu.Unlock()
    
    o.counts[string(event.Type)]++
}
```

## Complete Examples

### Production Monitoring Stack

```go
package main

import (
    "database/sql"
    "fmt"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/yesoreyeram/thaiyyal/backend"
    "go.uber.org/zap"
)

func main() {
    // Setup logging
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    
    // Setup database for audit logging
    db, _ := sql.Open("postgres", "...")
    
    // Setup Prometheus metrics
    reg := prometheus.NewRegistry()
    
    // Create engine
    engine, _ := workflow.NewEngine([]byte(payload))
    
    // Register production observers
    engine.
        RegisterObserver(NewPrometheusObserver(reg)).
        RegisterObserver(NewAuditObserver(db)).
        RegisterObserver(NewZapObserver(logger)).
        SetLogger(NewZapAdapter(logger))
    
    // Execute with full observability
    result, err := engine.Execute()
    if err != nil {
        logger.Error("Workflow failed", zap.Error(err))
        return
    }
    
    fmt.Printf("Result: %v\n", result.FinalOutput)
}
```

### Development Debugging

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yesoreyeram/thaiyyal/backend"
)

func main() {
    payload := `{...}`
    
    engine, err := workflow.NewEngine([]byte(payload))
    if err != nil {
        log.Fatal(err)
    }
    
    // Use console observer for debugging
    engine.RegisterObserver(workflow.NewConsoleObserver())
    
    result, err := engine.Execute()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %v\n", result.FinalOutput)
}
```

## API Reference

### Engine Methods

```go
// RegisterObserver adds an observer (returns engine for chaining)
func (e *Engine) RegisterObserver(obs Observer) *Engine

// SetLogger sets the logger (returns engine for chaining)
func (e *Engine) SetLogger(logger Logger) *Engine

// GetObserverCount returns the number of registered observers
func (e *Engine) GetObserverCount() int
```

### Observer Factory Functions

```go
// NewConsoleObserver creates a console observer with default logger
func NewConsoleObserver() *ConsoleObserver

// NewConsoleObserverWithLogger creates a console observer with custom logger
func NewConsoleObserverWithLogger(logger Logger) *ConsoleObserver

// NewDefaultLogger creates the default logger
func NewDefaultLogger() *DefaultLogger
```

## Migration Guide

Existing code works without changes. To add observability:

**Before:**
```go
engine, _ := workflow.NewEngine([]byte(payload))
result, _ := engine.Execute()
```

**After:**
```go
engine, _ := workflow.NewEngine([]byte(payload))
engine.RegisterObserver(workflow.NewConsoleObserver())  // Added
result, _ := engine.Execute()
```

## FAQ

**Q: Do observers impact workflow performance?**  
A: No. Observers run asynchronously in separate goroutines and never block execution.

**Q: What happens if an observer panics?**  
A: Panics are recovered automatically. Other observers and workflow execution are not affected.

**Q: Can I use multiple observers?**  
A: Yes. Register as many observers as needed. They all receive events independently.

**Q: Are events guaranteed to be delivered in order?**  
A: Within a single observer, yes. Across multiple observers, no ordering guarantee.

**Q: Can I remove observers after registration?**  
A: Currently no. Observers are registered for the lifetime of the engine instance.

**Q: Do I need to close/cleanup observers?**  
A: Not required by the framework, but implement cleanup in your observers if needed.

## See Also

- [Custom Node Executors](CUSTOM_NODES.md)
- [Architecture Documentation](ARCHITECTURE.md)
- [Node Type Reference](../docs/NODES.md)
- [Integration Guide](INTEGRATION.md)
