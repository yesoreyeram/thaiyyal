# API Reference

## Engine API

### Creating an Engine

```go
import "github.com/yesoreyeram/thaiyyal/backend/workflow"

// Basic creation
engine, err := workflow.NewEngine(payloadJSON)

// With custom config
config := workflow.ValidationLimits()
engine, err := workflow.NewWithConfig(payloadJSON, config)

// With custom registry
registry := workflow.DefaultRegistry()
registry.MustRegister(&MyExecutor{})
engine, err := workflow.NewWithRegistry(payloadJSON, config, registry)
```

### Executing Workflows

```go
result, err := engine.Execute()
if err != nil {
    // Handle error
}

// Access results
finalOutput := result.FinalOutput
nodeResults := result.NodeResults
executionID := result.ExecutionID
```

### Observers

```go
// Register observer
engine.RegisterObserver(&MyObserver{})

// Observer interface
type Observer interface {
    OnEvent(ctx context.Context, event Event)
}
```

## Type Definitions

### Payload

```go
type Payload struct {
    WorkflowID string `json:"workflow_id"`
    Nodes      []Node `json:"nodes"`
    Edges      []Edge `json:"edges"`
}
```

### Node

```go
type Node struct {
    ID   string   `json:"id"`
    Type NodeType `json:"type"`
    Data NodeData `json:"data"`
}
```

### Result

```go
type Result struct {
    ExecutionID string                 `json:"execution_id"`
    WorkflowID  string                 `json:"workflow_id"`
    NodeResults map[string]interface{} `json:"node_results"`
    FinalOutput interface{}            `json:"final_output"`
    Errors      []string               `json:"errors"`
}
```

---

**Last Updated:** 2025-11-03
**Version:** 1.0
