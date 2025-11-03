# Principles: Single Responsibility

## Overview

Each package, type, and function in Thaiyyal has a single, well-defined responsibility.

## Package-Level SRP

| Package | Single Responsibility |
|---------|----------------------|
| `types` | Define shared types |
| `engine` | Orchestrate workflow execution |
| `executor` | Execute individual nodes |
| `graph` | Analyze workflow structure |
| `state` | Manage workflow state |
| `security` | Provide security utilities |
| `logging` | Structured logging |
| `observer` | Event notification |

## Type-Level SRP

```go
// Engine: Orchestrate execution only
type Engine struct {
    graph    *graph.Graph       // Delegates to graph
    state    *state.Manager     // Delegates to state
    registry *executor.Registry // Delegates to registry
}

// Graph: Analyze structure only
type Graph struct {
    nodes []Node
    edges []Edge
}

// StateManager: Manage state only  
type Manager struct {
    variables map[string]interface{}
}
```

## Function-Level SRP

```go
// Parse: Only parse JSON
func Parse(data []byte) (*Payload, error)

// Validate: Only validate structure
func Validate(payload Payload) error

// Execute: Only execute workflow
func Execute(payload Payload) (*Result, error)
```

## Benefits

- **Easy to Test**: Single responsibility = simple tests
- **Easy to Understand**: Clear purpose
- **Easy to Modify**: Changes affect single concern
- **Easy to Reuse**: Focused functionality

---

**Last Updated:** 2025-11-03
**Version:** 1.0
