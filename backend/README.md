# Thaiyyal Backend - Workflow Engine

A Go-based workflow execution engine that parses and executes JSON workflow payloads from the Thaiyyal frontend.

## Features

- **JSON Payload Parsing**: Accepts workflow definitions as JSON
- **DAG Execution**: Executes workflows using topological sorting to ensure correct node execution order
- **Node Types**:
  - **Number Nodes**: Input nodes that provide numeric values
  - **Operation Nodes**: Perform arithmetic operations (add, subtract, multiply, divide)
  - **Visualization Nodes**: Format and present final results (text, table modes)
- **Cycle Detection**: Prevents execution of workflows with circular dependencies
- **Error Handling**: Comprehensive error reporting for invalid workflows or execution failures

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
    // Define a workflow payload
    payload := `{
        "nodes": [
            {"id": "1", "data": {"value": 10}},
            {"id": "2", "data": {"value": 5}},
            {"id": "3", "data": {"op": "add"}},
            {"id": "4", "data": {"mode": "text"}}
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

```json
{
  "nodes": [
    {
      "id": "unique-node-id",
      "data": {
        "value": 10,          // for number nodes
        "op": "add",          // for operation nodes (add, subtract, multiply, divide)
        "mode": "text",       // for visualization nodes (text, table)
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

## Node Types

### Number Node
Provides a numeric input value.
```json
{"id": "1", "data": {"value": 42}}
```

### Operation Node
Performs arithmetic operations on two input values.

Supported operations:
- `add`: Addition
- `subtract`: Subtraction
- `multiply`: Multiplication
- `divide`: Division

```json
{"id": "2", "data": {"op": "add"}}
```

### Visualization Node
Formats the output for display.

Supported modes:
- `text`: Plain text output
- `table`: Tabular output

```json
{"id": "3", "data": {"mode": "text"}}
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

## Implementation Details

### Execution Algorithm

1. **Parse JSON**: Unmarshal the payload into Go structs
2. **Build Graph**: Create adjacency list and calculate in-degrees
3. **Topological Sort**: Use Kahn's algorithm to determine execution order
4. **Cycle Detection**: Verify all nodes can be executed (no cycles)
5. **Execute Nodes**: Process nodes in topological order
6. **Collect Results**: Gather intermediate and final results

### Error Handling

The engine handles various error conditions:
- Invalid JSON syntax
- Cyclic dependencies
- Missing node inputs
- Division by zero
- Unknown operations
- Type mismatches

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
