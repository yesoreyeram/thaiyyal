# Engine Package

The `engine` package provides the main workflow execution engine for Thaiyyal. It orchestrates workflow parsing, validation, and execution using the refactored modular architecture.

## Overview

The Engine is the central component that:
1. Parses JSON workflow definitions
2. Infers node types if not explicitly specified
3. Constructs a DAG (Directed Acyclic Graph)
4. Performs topological sorting to determine execution order
5. Executes nodes using registered executors
6. Manages workflow state (variables, accumulators, counters, cache)
7. Handles template interpolation for context variables

## Architecture

```
Engine
├── graph.Graph        - DAG operations and topological sorting
├── state.Manager      - Workflow state management
├── executor.Registry  - Node executor registry (Strategy Pattern)
└── types.Config       - Execution configuration
```

## Design Patterns

- **Strategy Pattern**: Different execution strategies for different node types via Registry
- **State Pattern**: Manages workflow state across execution
- **Template Method**: Execute() defines the workflow execution algorithm
- **Dependency Injection**: Components (graph, state, registry) are injected

## Usage

### Basic Usage

```go
import "github.com/yesoreyeram/thaiyyal/backend/pkg/engine"

// Create engine from JSON
payload := []byte(`{
    "nodes": [
        {"id": "1", "type": "number", "data": {"value": 10}},
        {"id": "2", "type": "number", "data": {"value": 5}},
        {"id": "3", "type": "operation", "data": {"op": "add"}}
    ],
    "edges": [
        {"id": "e1", "source": "1", "target": "3"},
        {"id": "e2", "source": "2", "target": "3"}
    ]
}`)

eng, err := engine.New(payload)
if err != nil {
    log.Fatal(err)
}

// Execute workflow
result, err := eng.Execute()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %v\n", result.FinalOutput) // Output: 15
```

### With Custom Configuration

```go
import (
    "time"
    "github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
    "github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

config := types.Config{
    MaxExecutionTime:    5 * time.Second,
    MaxHTTPTimeout:      2 * time.Second,
    MaxForEachIterations: 100,
}

eng, err := engine.NewWithConfig(payload, config)
```

## Features

### Type Inference

The engine automatically infers node types from data fields:

```go
// No explicit type needed - inferred as "number" from "value" field
{"id": "1", "data": {"value": 10}}

// Inferred as "text_input" from "text" field
{"id": "2", "data": {"text": "hello"}}

// Inferred as "operation" from "op" field
{"id": "3", "data": {"op": "add"}}
```

### Template Interpolation

Context nodes define variables that can be interpolated in other nodes:

```go
{
    "nodes": [
        {
            "id": "ctx1",
            "type": "context_variable",
            "data": {
                "context_values": [
                    {"name": "api_url", "value": "https://api.example.com", "type": "string"}
                ]
            }
        },
        {
            "id": "http1",
            "type": "http",
            "data": {
                "url": "{{ variable.api_url }}/users"
            }
        }
    ]
}
```

### State Management

The engine implements the ExecutionContext interface for state access:

```go
// Variables
eng.SetVariable("user_id", 123)
val, _ := eng.GetVariable("user_id")

// Accumulator
eng.SetAccumulator([]int{1, 2, 3})
acc := eng.GetAccumulator()

// Counter
eng.SetCounter(42.5)
count := eng.GetCounter()

// Cache with TTL
eng.SetCache("key", "value", 5*time.Minute)
val, found := eng.GetCache("key")

// Context variables
eng.SetContextVariable("api_key", "secret")
key, exists := eng.GetContextVariable("api_key")
```

### Execution Metadata

Each execution has a unique ID for tracking:

```go
result, _ := eng.Execute()
fmt.Printf("Execution ID: %s\n", result.ExecutionID)
fmt.Printf("Workflow ID: %s\n", result.WorkflowID)
```

## ExecutionContext Interface

The Engine implements the `executor.ExecutionContext` interface, providing executors with:

- **Input Retrieval**: `GetNodeInputs(nodeID)`, `GetNode(nodeID)`
- **State Management**: Variables, Accumulator, Counter, Cache
- **Context Operations**: Context variables and constants
- **Result Management**: `GetNodeResult(nodeID)`, `SetNodeResult(nodeID, result)`
- **Configuration**: `GetConfig()`
- **Template Interpolation**: `InterpolateTemplate(template)`

## Node Type Registry

The engine automatically registers all 25 node type executors:

**Basic I/O**: Number, TextInput, Visualization  
**Operations**: Operation, TextOperation, HTTP  
**Control Flow**: Condition, ForEach, WhileLoop  
**State**: Variable, Extract, Transform, Accumulator, Counter  
**Advanced**: Switch, Parallel, Join, Split, Delay, Cache  
**Resilience**: Retry, TryCatch, Timeout  
**Context**: ContextVariable, ContextConstant  

## Error Handling

The engine provides comprehensive error handling:

```go
result, err := eng.Execute()
if err != nil {
    log.Printf("Execution failed: %v", err)
    for _, e := range result.Errors {
        log.Printf("  - %s", e)
    }
}
```

### Timeout Protection

All executions are protected by a timeout:

```go
// Default: 30 seconds
eng, _ := engine.New(payload)

// Custom timeout
config := types.DefaultConfig()
config.MaxExecutionTime = 10 * time.Second
eng, _ := engine.NewWithConfig(payload, config)
```

## Thread Safety

The engine is thread-safe for:
- Result storage (protected by `resultsMu`)
- State management (state.Manager has internal locking)
- Cache operations (protected internally)

However, a single Engine instance should execute one workflow at a time. Create new instances for concurrent executions.

## Testing

```bash
# Run all engine tests
go test ./pkg/engine

# Run specific test
go test ./pkg/engine -run TestExecute

# Run with verbose output
go test -v ./pkg/engine
```

## Performance Considerations

- **Topological Sort**: O(V + E) complexity where V=nodes, E=edges
- **Memory**: O(V) for node results
- **Execution**: Sequential by default (parallel execution via Parallel node)
- **Type Inference**: One-time cost at engine creation

## Migration from Legacy Code

The engine package replaces the monolithic `workflow.Engine`:

### Before (workflow.go)
```go
import "github.com/yesoreyeram/thaiyyal/backend/workflow"

engine, _ := workflow.NewEngine(payload)
result, _ := engine.Execute()
```

### After (pkg/engine)
```go
import "github.com/yesoreyeram/thaiyyal/backend/pkg/engine"

eng, _ := engine.New(payload)
result, _ := eng.Execute()
```

The API is intentionally similar for easy migration.

## Dependencies

```
pkg/engine
├── pkg/types      - Core type definitions
├── pkg/graph      - DAG operations
├── pkg/state      - State management
└── pkg/executor   - Node executors
```

## Future Enhancements

- [ ] Parallel execution optimization
- [ ] Workflow validation before execution
- [ ] Execution metrics and observability
- [ ] Workflow versioning support
- [ ] Distributed execution support

## See Also

- [pkg/executor](../executor/README.md) - Node executor implementation
- [pkg/types](../types/README.md) - Type definitions
- [pkg/graph](../graph/README.md) - Graph operations
- [pkg/state](../state/README.md) - State management
- [ARCHITECTURE.md](../../ARCHITECTURE.md) - Overall architecture

## License

MIT
