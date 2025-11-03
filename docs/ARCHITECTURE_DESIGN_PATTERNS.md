# Architecture Design Patterns

This document describes the design patterns used in the Thaiyyal workflow engine and explains why they were chosen.

## Table of Contents

- [Creational Patterns](#creational-patterns)
- [Structural Patterns](#structural-patterns)
- [Behavioral Patterns](#behavioral-patterns)
- [Concurrency Patterns](#concurrency-patterns)
- [Pattern Catalog](#pattern-catalog)

## Overview

Thaiyyal uses proven design patterns to ensure:
- **Maintainability**: Easy to understand and modify
- **Extensibility**: Simple to add new features
- **Testability**: Easy to write comprehensive tests
- **Performance**: Efficient execution

## Creational Patterns

### 1. Factory Method

**Location:** `backend/pkg/engine/engine.go`

**Purpose:** Create workflow engine instances with different configurations.

**Implementation:**

```go
// Basic factory
func New(payloadJSON []byte) (*Engine, error)

// Factory with custom configuration
func NewWithConfig(payloadJSON []byte, config types.Config) (*Engine, error)

// Factory with custom registry
func NewWithRegistry(payloadJSON []byte, config types.Config, registry *executor.Registry) (*Engine, error)
```

**Benefits:**
- Consistent initialization
- Configuration flexibility
- Easy testing with different setups

**Usage Example:**

```go
// Development: Use default config
engine, err := engine.New(payload)

// Production: Use strict validation
config := types.ValidationLimits()
engine, err := engine.NewWithConfig(payload, config)

// Custom: Add custom executors
registry := engine.DefaultRegistry()
registry.MustRegister(&MyCustomExecutor{})
engine, err := engine.NewWithRegistry(payload, config, registry)
```

### 2. Builder Pattern (Implicit)

**Location:** `backend/pkg/logging/logger.go`, `backend/pkg/security/ssrf.go`

**Purpose:** Construct complex objects step-by-step.

**Implementation:**

```go
// Logger builder pattern
logger := logging.New(cfg)
logger = logger.WithWorkflowID(id)
logger = logger.WithExecutionID(execID)
logger = logger.WithField("key", value)

// SSRF config builder
config := security.SSRFConfig{
    BlockPrivateIPs: true,
    BlockLocalhost: true,
}
protection := security.NewSSRFProtectionWithConfig(config)
```

**Benefits:**
- Fluent API
- Immutable loggers (new instance per modification)
- Clear configuration

### 3. Singleton (Avoided)

**Why NOT used:**
- Global state complicates testing
- Reduces flexibility
- Makes parallelization harder

**Alternative:** Dependency injection

```go
// Instead of singleton
type Engine struct {
    logger *logging.Logger  // Injected, not global
    state  *state.Manager   // Injected
}
```

## Structural Patterns

### 1. Facade Pattern

**Location:** `backend/workflow.go`

**Purpose:** Provide simple interface to complex subsystem.

**Implementation:**

```go
// Facade provides backward compatibility
package workflow

// Simple API hides complexity
type NodeType = types.NodeType
type Engine = engine.Engine

var NewEngine = engine.New
var DefaultRegistry = engine.DefaultRegistry
```

**Benefits:**
- Backward compatibility
- Simple API for common use cases
- Hides refactoring from users

**Usage:**

```go
// Users use simple facade
import "github.com/yesoreyeram/thaiyyal/backend/workflow"

engine, err := workflow.NewEngine(payload)
```

### 2. Adapter Pattern

**Location:** `backend/pkg/observer/observer.go`

**Purpose:** Convert one interface to another.

**Implementation:**

```go
// Adapter: Convert Logger interface to Observer
type ConsoleObserver struct {
    logger Logger
}

func (o *ConsoleObserver) OnEvent(ctx context.Context, event Event) {
    // Adapt Event to Logger calls
    switch event.Type {
    case EventNodeStart:
        o.logger.Infof("Node %s started", event.NodeID)
    case EventNodeSuccess:
        o.logger.Infof("Node %s succeeded", event.NodeID)
    }
}
```

**Benefits:**
- Reuse existing loggers as observers
- Interface compatibility
- Flexibility

### 3. Decorator Pattern (Middleware)

**Location:** `backend/pkg/middleware/middleware.go`

**Purpose:** Add responsibilities to objects dynamically.

**Implementation:**

```go
// Base handler
type Handler func(ctx ExecutionContext, node types.Node) (interface{}, error)

// Middleware wraps handler
type Middleware interface {
    Process(ctx ExecutionContext, node types.Node, next Handler) (interface{}, error)
    Name() string
}

// Example: Logging middleware
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) Process(ctx ExecutionContext, node types.Node, next Handler) (interface{}, error) {
    // Before
    start := time.Now()
    
    // Execute
    result, err := next(ctx, node)
    
    // After
    duration := time.Since(start)
    log.Printf("Node %s executed in %v", node.ID, duration)
    
    return result, err
}
```

**Middleware Chain:**

```
Request
  ↓
ValidationMiddleware
  ↓
LoggingMiddleware
  ↓
RateLimitMiddleware
  ↓
TimeoutMiddleware
  ↓
Actual Executor
  ↓
Response
```

**Benefits:**
- Add features without modifying executors
- Composable behaviors
- Easy to test individually
- Single Responsibility Principle

## Behavioral Patterns

### 1. Strategy Pattern

**Location:** `backend/pkg/executor/registry.go`

**Purpose:** Select algorithm at runtime.

**Implementation:**

```go
// Strategy interface
type NodeExecutor interface {
    Type() types.NodeType
    Execute(ctx ExecutionContext, node types.Node) (interface{}, error)
    Validate(node types.Node) error
}

// Concrete strategies
type NumberExecutor struct{}
type OperationExecutor struct{}
type HTTPExecutor struct{}

// Context uses strategies
type Registry struct {
    executors map[types.NodeType]NodeExecutor
}

func (r *Registry) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    executor := r.executors[node.Type]
    return executor.Execute(ctx, node)
}
```

**Strategy Selection:**

```
Node Type → Registry → Select Executor → Execute
   ↓
"number" → NumberExecutor
"operation" → OperationExecutor
"http" → HTTPExecutor
"custom" → CustomExecutor
```

**Benefits:**
- Easy to add new node types
- Encapsulates algorithms
- Runtime selection
- Testable in isolation

### 2. Template Method Pattern

**Location:** `backend/pkg/engine/engine.go`

**Purpose:** Define algorithm skeleton, defer steps to subclasses.

**Implementation:**

```go
// Template method defines algorithm
func (e *Engine) Execute() (*types.Result, error) {
    // Step 1: Infer types (hook)
    e.inferNodeTypes()
    
    // Step 2: Sort (algorithm)
    order, err := e.graph.TopologicalSort()
    
    // Step 3: Execute (hook point)
    for _, nodeID := range order {
        result, err := e.executeNode(ctx, node)  // Hook
    }
    
    // Step 4: Return result
    return result, nil
}

// Hook: Can be customized via Registry
func (e *Engine) executeNode(ctx context.Context, node types.Node) (interface{}, error) {
    return e.registry.Execute(e, node)  // Delegates to strategy
}
```

**Execution Template:**

```
Execute() [Template Method]
  │
  ├─→ Parse JSON [Fixed Step]
  ├─→ Infer Types [Hook - can override]
  ├─→ Validate [Fixed Step]
  ├─→ Topological Sort [Fixed Step]
  ├─→ Execute Nodes [Hook - customizable]
  │    └─→ Registry dispatches to executors
  └─→ Return Result [Fixed Step]
```

**Benefits:**
- Consistent algorithm structure
- Extensible at specific points
- Code reuse
- Control inversion

### 3. Observer Pattern

**Location:** `backend/pkg/observer/`

**Purpose:** Notify multiple objects about state changes.

**Implementation:**

```go
// Subject
type Manager struct {
    observers []Observer
    mu        sync.RWMutex
}

// Observer interface
type Observer interface {
    OnEvent(ctx context.Context, event Event)
}

// Register observer
func (m *Manager) Register(obs Observer) {
    m.observers = append(m.observers, obs)
}

// Notify all observers
func (m *Manager) Notify(ctx context.Context, event Event) {
    for _, obs := range m.observers {
        obs.OnEvent(ctx, event)
    }
}
```

**Event Flow:**

```
Engine Execution
  ↓
Event: WorkflowStart
  ↓
Manager.Notify()
  ↓
  ├─→ ConsoleObserver → Logs to console
  ├─→ MetricsObserver → Records metrics
  └─→ CustomObserver → User-defined logic
```

**Events:**
- `WorkflowStart`: Workflow execution begins
- `WorkflowEnd`: Workflow execution completes
- `NodeStart`: Node execution begins
- `NodeSuccess`: Node execution succeeds
- `NodeFailure`: Node execution fails

**Benefits:**
- Loose coupling
- Multiple subscribers
- Easy to add monitoring
- Event-driven architecture

### 4. State Pattern

**Location:** `backend/pkg/state/manager.go`

**Purpose:** Object changes behavior based on internal state.

**Implementation:**

```go
// State manager maintains workflow state
type Manager struct {
    variables map[string]interface{}
    context   map[string]interface{}
    accumulator interface{}
    counter   float64
    cache     map[string]CacheEntry
    mu        sync.RWMutex
}

// State transitions
func (m *Manager) SetVariable(name string, value interface{}) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.variables[name] = value  // State change
    return nil
}
```

**State Lifecycle:**

```
Initial State
  ↓
Variables: {}
Accumulator: nil
Counter: 0
Cache: {}
  ↓
[Execution]
  ↓
Variables: {user: "john", count: 5}
Accumulator: [1, 2, 3]
Counter: 10
Cache: {api_result: {...}}
  ↓
Final State
```

**Benefits:**
- Encapsulates state management
- Thread-safe with mutex
- Clear state transitions
- Isolated from execution logic

### 5. Chain of Responsibility

**Location:** `backend/pkg/middleware/`

**Purpose:** Pass request along chain of handlers.

**Implementation:**

```go
// Handler type
type Handler func(ctx ExecutionContext, node types.Node) (interface{}, error)

// Chain builder
func BuildChain(middlewares []Middleware, executor NodeExecutor) Handler {
    handler := func(ctx ExecutionContext, node types.Node) (interface{}, error) {
        return executor.Execute(ctx, node)
    }
    
    // Build chain in reverse
    for i := len(middlewares) - 1; i >= 0; i-- {
        middleware := middlewares[i]
        next := handler
        handler = func(ctx ExecutionContext, node types.Node) (interface{}, error) {
            return middleware.Process(ctx, node, next)
        }
    }
    
    return handler
}
```

**Chain Execution:**

```
Request
  ↓
Middleware 1 (Validation)
  ├─ Pre-process
  ↓
Middleware 2 (Logging)
  ├─ Pre-process
  ↓
Middleware 3 (Rate Limit)
  ├─ Pre-process
  ↓
Executor (Core Logic)
  ↓
Middleware 3 (Rate Limit)
  ├─ Post-process
  ↓
Middleware 2 (Logging)
  ├─ Post-process
  ↓
Middleware 1 (Validation)
  ├─ Post-process
  ↓
Response
```

**Benefits:**
- Decouple senders/receivers
- Add/remove handlers dynamically
- Order-independent handlers
- Single Responsibility

### 6. Registry Pattern

**Location:** `backend/pkg/executor/registry.go`

**Purpose:** Central repository for node executors.

**Implementation:**

```go
type Registry struct {
    executors map[types.NodeType]NodeExecutor
    mu        sync.RWMutex
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

func (r *Registry) Get(nodeType types.NodeType) (NodeExecutor, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    executor, exists := r.executors[nodeType]
    if !exists {
        return nil, fmt.Errorf("executor not found: %s", nodeType)
    }
    
    return executor, nil
}
```

**Benefits:**
- Central registration
- Type-safe lookup
- Easy to extend
- Thread-safe

## Concurrency Patterns

### 1. Mutex Pattern

**Location:** Throughout codebase

**Purpose:** Protect shared resources.

**Implementation:**

```go
type Manager struct {
    variables map[string]interface{}
    mu        sync.RWMutex
}

func (m *Manager) GetVariable(name string) (interface{}, error) {
    m.mu.RLock()  // Read lock
    defer m.mu.RUnlock()
    
    value, exists := m.variables[name]
    if !exists {
        return nil, fmt.Errorf("variable not found: %s", name)
    }
    return value, nil
}

func (m *Manager) SetVariable(name string, value interface{}) error {
    m.mu.Lock()  // Write lock
    defer m.mu.Unlock()
    
    m.variables[name] = value
    return nil
}
```

**RWMutex Usage:**
- Multiple readers OR single writer
- Readers don't block readers
- Writers block everything

### 2. Context Pattern

**Location:** Throughout execution

**Purpose:** Cancellation, deadlines, values.

**Implementation:**

```go
// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), config.MaxExecutionTime)
defer cancel()

// Add values
ctx = context.WithValue(ctx, types.ContextKeyExecutionID, executionID)

// Check cancellation
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // Continue execution
}
```

**Context Flow:**

```
Root Context
  ↓
WithTimeout (30s)
  ↓
WithValue (executionID)
  ↓
WithValue (workflowID)
  ↓
Passed to all executors
  ↓
Check ctx.Done() in loops
```

**Benefits:**
- Graceful cancellation
- Request-scoped values
- Timeout enforcement
- Propagation tree

### 3. Worker Pool (Future)

**Location:** Planned for parallel execution

**Purpose:** Limit concurrent executions.

**Design:**

```go
type WorkerPool struct {
    workers   int
    taskQueue chan Task
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        go p.worker()
    }
}

func (p *WorkerPool) worker() {
    for task := range p.taskQueue {
        task.Execute()
    }
}
```

## Pattern Catalog

### Pattern Selection Guide

| Need | Pattern | Location |
|------|---------|----------|
| Create objects | Factory Method | `engine.New*()` |
| Configure objects | Builder | `logger.With*()` |
| Simple API | Facade | `workflow.go` |
| Interface conversion | Adapter | `observer/` |
| Add features | Decorator | `middleware/` |
| Select algorithm | Strategy | `executor/registry.go` |
| Algorithm skeleton | Template Method | `engine.Execute()` |
| Event notification | Observer | `observer/` |
| Manage state | State | `state/manager.go` |
| Request handling | Chain of Responsibility | `middleware/` |
| Service lookup | Registry | `executor/registry.go` |
| Thread safety | Mutex | Throughout |
| Cancellation | Context | Throughout |

### Anti-Patterns Avoided

1. **God Object**: Large classes doing everything
   - **Avoided by:** Small, focused packages (SRP)

2. **Singleton**: Global state
   - **Avoided by:** Dependency injection

3. **Magic Numbers**: Hardcoded values
   - **Avoided by:** Configuration constants

4. **Copy-Paste**: Duplicated code
   - **Avoided by:** Shared packages (types, state)

5. **Tight Coupling**: Direct dependencies
   - **Avoided by:** Interfaces, Registry, DI

## Pattern Benefits Summary

### Maintainability
- **Single Responsibility**: Each pattern/package has one job
- **Open/Closed**: Extend without modifying (Strategy, Registry)
- **Clear Structure**: Template Method defines flow

### Extensibility
- **Strategy Pattern**: Add node types easily
- **Registry Pattern**: Register custom executors
- **Observer Pattern**: Add monitoring easily
- **Middleware Pattern**: Add cross-cutting concerns

### Testability
- **Dependency Injection**: Easy to mock dependencies
- **Interfaces**: Test against contracts
- **Isolated Components**: Test in isolation
- **Strategy Pattern**: Test executors independently

### Performance
- **RWMutex**: Concurrent reads
- **Context**: Efficient cancellation
- **Registry**: Fast lookup (map)
- **Caching**: Template compilation, results

## Related Documentation

- [Architecture Overview](ARCHITECTURE.md)
- [Pluggable Architecture](PRINCIPLES_PLUGGABLE_ARCHITECTURE.md)
- [Modular Design](PRINCIPLES_MODULAR_DESIGN.md)
- [Single Responsibility](PRINCIPLES_SINGLE_RESPONSIBILITY.md)

---

**Last Updated:** 2025-11-03
**Version:** 1.0
