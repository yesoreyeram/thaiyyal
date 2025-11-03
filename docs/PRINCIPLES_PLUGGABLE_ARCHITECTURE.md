# Principles: Pluggable Architecture

This document describes the pluggable architecture that allows extending Thaiyyal with custom node types and functionality.

## Overview

The pluggable architecture enables:
- **Custom node executors** without modifying core code
- **Middleware injection** for cross-cutting concerns
- **Observer registration** for monitoring
- **HTTP client customization** for special requirements

## Plugin System Design

### 1. Node Executor Plugin Interface

**Core Interface:**

```go
// NodeExecutor is the interface all node executors must implement
type NodeExecutor interface {
    // Type returns the node type this executor handles
    Type() types.NodeType
    
    // Execute runs the node logic
    Execute(ctx ExecutionContext, node types.Node) (interface{}, error)
    
    // Validate checks if node configuration is valid
    Validate(node types.Node) error
}
```

**Creating a Custom Executor:**

```go
package mypackage

import (
    "github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
    "github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// CustomExecutor implements a custom node type
type CustomExecutor struct {
    // Configuration or dependencies
}

func (e *CustomExecutor) Type() types.NodeType {
    return "my_custom_type"
}

func (e *CustomExecutor) Execute(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
    // Get inputs from predecessor nodes
    inputs := ctx.GetNodeInputs(node.ID)
    
    // Access workflow state
    value, _ := ctx.GetVariable("myvar")
    
    // Perform custom logic
    result := myCustomLogic(inputs, value)
    
    // Store results
    ctx.SetVariable("output", result)
    
    return result, nil
}

func (e *CustomExecutor) Validate(node types.Node) error {
    // Validate node configuration
    if node.Data.CustomField == nil {
        return fmt.Errorf("custom field is required")
    }
    return nil
}
```

**Registering Custom Executor:**

```go
// Get default registry with all built-in executors
registry := engine.DefaultRegistry()

// Register custom executor
registry.MustRegister(&mypackage.CustomExecutor{})

// Create engine with custom registry
engine, err := engine.NewWithRegistry(payload, config, registry)
```

### 2. Registry Pattern

**Implementation:**

```go
type Registry struct {
    executors map[types.NodeType]NodeExecutor
    mu        sync.RWMutex
}

func NewRegistry() *Registry {
    return &Registry{
        executors: make(map[types.NodeType]NodeExecutor),
    }
}

func (r *Registry) Register(executor NodeExecutor) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    nodeType := executor.Type()
    if _, exists := r.executors[nodeType]; exists {
        return fmt.Errorf("executor already registered: %s", nodeType)
    }
    
    r.executors[nodeType] = executor
    return nil
}

func (r *Registry) MustRegister(executor NodeExecutor) {
    if err := r.Register(executor); err != nil {
        panic(err)
    }
}

func (r *Registry) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    executor, err := r.Get(node.Type)
    if err != nil {
        return nil, err
    }
    return executor.Execute(ctx, node)
}
```

## Middleware System

### 1. Middleware Interface

```go
type Middleware interface {
    Process(ctx ExecutionContext, node types.Node, next Handler) (interface{}, error)
    Name() string
}

type Handler func(ctx ExecutionContext, node types.Node) (interface{}, error)
```

### 2. Custom Middleware Example

```go
// Logging middleware
type LoggingMiddleware struct {
    logger *logging.Logger
}

func (m *LoggingMiddleware) Process(ctx ExecutionContext, node types.Node, next Handler) (interface{}, error) {
    // Pre-execution logging
    m.logger.WithField("node_id", node.ID).Info("Executing node")
    
    // Execute next in chain
    result, err := next(ctx, node)
    
    // Post-execution logging
    if err != nil {
        m.logger.WithError(err).Error("Node execution failed")
    } else {
        m.logger.Info("Node execution succeeded")
    }
    
    return result, err
}

func (m *LoggingMiddleware) Name() string {
    return "Logging"
}
```

### 3. Metrics Middleware

```go
type MetricsMiddleware struct {
    recorder MetricsRecorder
}

func (m *MetricsMiddleware) Process(ctx ExecutionContext, node types.Node, next Handler) (interface{}, error) {
    start := time.Now()
    
    result, err := next(ctx, node)
    
    duration := time.Since(start)
    m.recorder.RecordDuration(string(node.Type), duration)
    
    if err != nil {
        m.recorder.RecordError(string(node.Type))
    }
    
    return result, err
}
```

## Observer Pattern

### 1. Observer Interface

```go
type Observer interface {
    OnEvent(ctx context.Context, event Event)
}

type Event struct {
    Type        EventType
    Status      ExecutionStatus
    Timestamp   time.Time
    ExecutionID string
    WorkflowID  string
    NodeID      string
    NodeType    NodeType
    Result      interface{}
    Error       error
    ElapsedTime time.Duration
}
```

### 2. Custom Observer Example

```go
// Metrics observer
type MetricsObserver struct {
    collector MetricsCollector
}

func (o *MetricsObserver) OnEvent(ctx context.Context, event Event) {
    switch event.Type {
    case EventWorkflowStart:
        o.collector.IncrementWorkflowCount()
        
    case EventWorkflowEnd:
        o.collector.RecordWorkflowDuration(event.ElapsedTime)
        if event.Error != nil {
            o.collector.IncrementWorkflowFailures()
        }
        
    case EventNodeStart:
        o.collector.IncrementNodeCount(string(event.NodeType))
        
    case EventNodeEnd:
        o.collector.RecordNodeDuration(
            string(event.NodeType),
            event.ElapsedTime,
        )
    }
}

// Register observer
engine.RegisterObserver(&MetricsObserver{collector: myCollector})
```

## HTTP Client Customization

### 1. Named HTTP Clients

```go
// Create HTTP client registry
httpRegistry := httpclient.NewRegistry()

// Add named client with custom config
httpRegistry.Register("api-client", &httpclient.Config{
    Timeout:     30 * time.Second,
    MaxRetries:  3,
    UserAgent:   "Thaiyyal/1.0",
})

// Set registry on engine
engine.SetHTTPClientRegistry(httpRegistry)
```

### 2. HTTP Middleware

```go
// Custom HTTP middleware
type AuthMiddleware struct {
    token string
}

func (m *AuthMiddleware) Process(req *http.Request, next httpclient.Handler) (*http.Response, error) {
    // Add authorization header
    req.Header.Set("Authorization", "Bearer "+m.token)
    
    return next(req)
}
```

## Complete Plugin Example

### Example: Database Query Node

```go
package dbplugin

import (
    "database/sql"
    "fmt"
    
    "github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
    "github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Define custom node type
const NodeTypeDBQuery types.NodeType = "db_query"

// DBQueryExecutor executes database queries
type DBQueryExecutor struct {
    db *sql.DB
}

func NewDBQueryExecutor(db *sql.DB) *DBQueryExecutor {
    return &DBQueryExecutor{db: db}
}

func (e *DBQueryExecutor) Type() types.NodeType {
    return NodeTypeDBQuery
}

func (e *DBQueryExecutor) Execute(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
    // Get query from node configuration
    query := node.Data.Query
    if query == nil {
        return nil, fmt.Errorf("query is required")
    }
    
    // Get parameters from inputs
    inputs := ctx.GetNodeInputs(node.ID)
    
    // Execute query
    rows, err := e.db.QueryContext(ctx, *query, inputs...)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()
    
    // Parse results
    results := []map[string]interface{}{}
    cols, _ := rows.Columns()
    
    for rows.Next() {
        values := make([]interface{}, len(cols))
        valuePtrs := make([]interface{}, len(cols))
        for i := range values {
            valuePtrs[i] = &values[i]
        }
        
        rows.Scan(valuePtrs...)
        
        row := make(map[string]interface{})
        for i, col := range cols {
            row[col] = values[i]
        }
        results = append(results, row)
    }
    
    return results, nil
}

func (e *DBQueryExecutor) Validate(node types.Node) error {
    if node.Data.Query == nil {
        return fmt.Errorf("query is required")
    }
    
    // Validate query syntax
    if err := validateSQL(*node.Data.Query); err != nil {
        return fmt.Errorf("invalid query: %w", err)
    }
    
    return nil
}

// Usage:
// db, _ := sql.Open("postgres", connectionString)
// registry := engine.DefaultRegistry()
// registry.MustRegister(dbplugin.NewDBQueryExecutor(db))
// engine, _ := engine.NewWithRegistry(payload, config, registry)
```

## Best Practices

### 1. Follow Interface Contracts

```go
// Always implement all interface methods
type MyExecutor struct{}

func (e *MyExecutor) Type() types.NodeType { ... }
func (e *MyExecutor) Execute(...) { ... }
func (e *MyExecutor) Validate(...) { ... }
```

### 2. Validate Thoroughly

```go
func (e *MyExecutor) Validate(node types.Node) error {
    // Check required fields
    // Validate data types
    // Check business rules
    // Validate against security constraints
    return nil
}
```

### 3. Handle Errors Properly

```go
func (e *MyExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    result, err := operation()
    if err != nil {
        return nil, fmt.Errorf("operation failed: %w", err)
    }
    return result, nil
}
```

### 4. Respect Resource Limits

```go
func (e *MyExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    // Check execution limits
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    // Respect timeout
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        return e.doWork(ctx, node)
    }
}
```

### 5. Thread Safety

```go
// Ensure executor is thread-safe if it has state
type MyExecutor struct {
    cache map[string]interface{}
    mu    sync.RWMutex
}

func (e *MyExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    // Safe access to shared state
    return e.cache[node.ID], nil
}
```

## Testing Plugins

```go
func TestMyExecutor(t *testing.T) {
    executor := &MyExecutor{}
    
    // Test type
    assert.Equal(t, "my_type", executor.Type())
    
    // Test validation
    node := types.Node{
        ID:   "test",
        Type: "my_type",
        Data: types.NodeData{},
    }
    err := executor.Validate(node)
    assert.Error(t, err)
    
    // Test execution
    ctx := &mockContext{}
    result, err := executor.Execute(ctx, node)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## Related Documentation

- [Architecture Overview](ARCHITECTURE.md)
- [Design Patterns](ARCHITECTURE_DESIGN_PATTERNS.md)
- [API Reference](API_REFERENCE.md)
- [Examples](EXAMPLES.md)

---

**Last Updated:** 2025-11-03
**Version:** 1.0
