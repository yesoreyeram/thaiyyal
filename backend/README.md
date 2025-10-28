# Thaiyyal Backend - Workflow Engine

A Go-based workflow execution engine that parses and executes JSON workflow payloads from the Thaiyyal frontend.

## Features

- **JSON Payload Parsing**: Accepts workflow definitions as JSON
- **Explicit Node Types**: Support for explicit type definition with automatic type inference fallback
- **DAG Execution**: Executes workflows using topological sorting to ensure correct node execution order
- **Node Types**:
  - **Number Nodes**: Input nodes that provide numeric values
  - **Operation Nodes**: Perform arithmetic operations (add, subtract, multiply, divide)
  - **Visualization Nodes**: Format and present final results (text, table modes)
- **Cycle Detection**: Prevents execution of workflows with circular dependencies
- **Error Handling**: Comprehensive error reporting for invalid workflows or execution failures
- **Type Safety**: Strongly typed constants for node types, operations, and visualization modes

## Usage

### As a Library

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yesoreyeram/thaiyyal/backend/workflow"
)

func main() {
    // Define a workflow payload with explicit types
    payload := `{
        "nodes": [
            {"id": "1", "type": "number", "data": {"value": 10}},
            {"id": "2", "type": "number", "data": {"value": 5}},
            {"id": "3", "type": "operation", "data": {"op": "add"}},
            {"id": "4", "type": "visualization", "data": {"mode": "text"}}
        ],
        "edges": [
            {"id": "e1-3", "source": "1", "target": "3"},
            {"id": "e2-3", "source": "2", "target": "3"},
            {"id": "e3-4", "source": "3", "target": "4"}
        ]
    }`
    
    // Create and execute the workflow
    engine, err := workflow.NewEngine([]byte(payload))
    if err != nil {
        log.Fatalf("Failed to create engine: %v", err)
    }
    
    result, err := engine.Execute()
    if err != nil {
        log.Fatalf("Execution failed: %v", err)
    }
    
    fmt.Printf("Final Output: %v\n", result.FinalOutput)
}
```

### Running Examples

The `backend/examples` directory contains several example workflows:

```bash
cd backend/examples
go run main.go
```

This will execute multiple example workflows including:
- Simple addition
- Complete workflow with visualization
- Complex multi-operation workflow
- Division operations

## Payload Format

The workflow engine expects a JSON payload with the following structure:

### With Explicit Types (Recommended)

```json
{
  "nodes": [
    {
      "id": "unique-node-id",
      "type": "number",       // explicit type: "number", "operation", or "visualization"
      "data": {
        "value": 10,
        "label": "Node Name"  // optional label
      }
    }
  ],
  "edges": [
    {
      "id": "unique-edge-id",
      "source": "source-node-id",
      "target": "target-node-id"
    }
  ]
}
```

### With Type Inference (Backward Compatible)

The engine can also infer types from data fields if the `type` field is omitted:

```json
{
  "nodes": [
    {
      "id": "unique-node-id",
      "data": {
        "value": 10           // presence of "value" infers type "number"
      }
    }
  ],
  "edges": [...]
}
```

## Node Types

### Type Constants

The package provides strongly-typed constants for all node types:

```go
// Node types
workflow.NodeTypeNumber         // "number"
workflow.NodeTypeOperation      // "operation"
workflow.NodeTypeVisualization  // "visualization"

// Operation types
workflow.OperationAdd       // "add"
workflow.OperationSubtract  // "subtract"
workflow.OperationMultiply  // "multiply"
workflow.OperationDivide    // "divide"

// Visualization modes
workflow.VisualizationModeText  // "text"
workflow.VisualizationModeTable // "table"
```

### Number Node
Provides a numeric input value.
```json
{"id": "1", "type": "number", "data": {"value": 42}}
```

### Operation Node
Performs arithmetic operations on two input values.

Supported operations:
- `add`: Addition
- `subtract`: Subtraction
- `multiply`: Multiplication
- `divide`: Division

```json
{"id": "2", "type": "operation", "data": {"op": "add"}}
```

### Visualization Node
Formats the output for display.

Supported modes:
- `text`: Plain text output
- `table`: Tabular output

```json
{"id": "3", "type": "visualization", "data": {"mode": "text"}}
```

## Result Format

The execution result contains:

```json
{
  "node_results": {
    "1": 10,
    "2": 5,
    "3": 15,
    "4": {
      "mode": "text",
      "value": 15
    }
  },
  "final_output": {
    "mode": "text",
    "value": 15
  },
  "errors": []
}
```

## Testing

Run the test suite:

```bash
cd backend/workflow
go test -v
```

The test suite includes:
- Payload parsing tests
- All operation types (add, subtract, multiply, divide)
- Complex multi-operation workflows
- Cycle detection
- Error handling (division by zero, missing inputs, invalid operations)
- Visualization modes
- Type inference tests
- Explicit type tests
- Mixed explicit and inferred type tests

## Implementation Details

### Execution Algorithm

1. **Parse JSON**: Unmarshal the payload into Go structs
2. **Infer Types**: Determine node types from data fields if not explicitly set
3. **Build Graph**: Create adjacency list and calculate in-degrees
4. **Topological Sort**: Use Kahn's algorithm to determine execution order
5. **Cycle Detection**: Verify all nodes can be executed (no cycles)
6. **Execute Nodes**: Process nodes in topological order using type-specific executors
7. **Collect Results**: Gather intermediate and final results

### Code Organization

The codebase is organized for readability and maintainability:

- **types.go**: All type definitions, constants, and data structures
- **engine.go**: Core execution engine with separate methods for each concern
- **engine_test.go**: Comprehensive execution tests
- **integration_test.go**: Frontend compatibility tests
- **types_test.go**: Type system and constant tests

### Type Inference

The engine automatically infers node types based on data fields:
- Presence of `value` field → `NodeTypeNumber`
- Presence of `op` field → `NodeTypeOperation`
- Presence of `mode` field → `NodeTypeVisualization`

This provides backward compatibility with payloads that don't specify explicit types.

### Error Handling

The engine handles various error conditions:
- Invalid JSON syntax
- Cyclic dependencies
- Missing node inputs
- Division by zero
- Unknown operations
- Type mismatches
- Invalid node types

All errors are reported in the result structure with descriptive messages.

## Future Enhancements

Potential extensions (not in current MVP):
- HTTP trigger nodes
- Database query nodes
- JSON transformation nodes
- Conditional branching (if/switch nodes)
- Parallel execution support
- State persistence
- Workflow versioning
- Advanced visualization types
