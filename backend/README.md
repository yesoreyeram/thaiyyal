# Thaiyyal Backend - Workflow Engine

A simple, easy-to-understand Go workflow execution engine that parses and executes JSON workflow payloads from the Thaiyyal frontend.

## Features

- **Simple & Readable**: Single file implementation (~250 lines)
- **Easy to Understand**: Straightforward code flow without complex patterns
- **JSON Payload Parsing**: Accepts workflow definitions as JSON
- **DAG Execution**: Uses topological sorting to execute nodes in correct order
- **Type Inference**: Automatically determines node types from data
- **Node Types**:
  - **Number Nodes**: Provide numeric input values
  - **Operation Nodes**: Perform arithmetic (add, subtract, multiply, divide)
  - **Visualization Nodes**: Format output for display (text, table)
- **Cycle Detection**: Prevents execution of workflows with circular dependencies
- **Comprehensive Tests**: 11 test cases covering all functionality

## Quick Start

### As a Library

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
            {"id": "e1-3", "source": "1", "target": "3"},
            {"id": "e2-3", "source": "2", "target": "3"}
        ]
    }`
    
    engine, err := workflow.NewEngine([]byte(payload))
    if err != nil {
        log.Fatal(err)
    }
    
    result, err := engine.Execute()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %v\n", result.FinalOutput)
}
```

### Running Examples

```bash
cd backend/examples
go run main.go
```

## File Structure

```
backend/
├── workflow.go          # Main workflow engine (single file)
├── workflow_test.go     # All tests (single file)
├── examples/
│   └── main.go         # Example usage
├── README.md           # This file
└── INTEGRATION.md      # Frontend integration guide
```

## Payload Format

The engine accepts JSON payloads with this structure:

```json
{
  "nodes": [
    {
      "id": "unique-id",
      "type": "number",          // optional: "number", "operation", or "visualization"
      "data": {
        "value": 10,             // for number nodes
        "op": "add",             // for operation nodes: add, subtract, multiply, divide
        "mode": "text",          // for visualization nodes: text, table
        "label": "My Node"       // optional label
      }
    }
  ],
  "edges": [
    {
      "id": "edge-id",
      "source": "source-node-id",
      "target": "target-node-id"
    }
  ]
}
```

## How It Works

1. **Parse JSON**: Unmarshal the payload into Go structs
2. **Infer Types**: Determine node types from data if not explicitly set
3. **Topological Sort**: Use Kahn's algorithm to find execution order
4. **Execute Nodes**: Process each node in order:
   - Number nodes return their value
   - Operation nodes compute results from inputs
   - Visualization nodes format outputs
5. **Return Results**: Collect all node results and determine final output

## Testing

Run tests:
```bash
cd backend
go test -v
```

All tests:
- ✅ Engine creation
- ✅ Simple addition
- ✅ All operations (add, subtract, multiply, divide)
- ✅ Division by zero error
- ✅ Complete workflow with visualization
- ✅ Multiple chained operations
- ✅ Cycle detection
- ✅ Missing input error
- ✅ Explicit node types
- ✅ Type inference
- ✅ Frontend payload compatibility

## Code Overview

### Main Functions

**NewEngine(payloadJSON []byte) (*Engine, error)**
- Creates a new engine from JSON payload
- Returns error if JSON is invalid

**Execute() (*Result, error)**
- Executes the workflow
- Returns results or error

### Internal Functions

- `inferNodeTypes()` - Determines types from node data
- `topologicalSort()` - Orders nodes using Kahn's algorithm
- `executeNode(node)` - Executes a single node
- `executeNumberNode(node)` - Returns number value
- `executeOperationNode(node)` - Performs arithmetic
- `executeVisualizationNode(node)` - Formats output
- `getNodeInputs(nodeID)` - Gets inputs from predecessors
- `getFinalOutput()` - Finds terminal node output

## Error Handling

The engine handles:
- Invalid JSON syntax
- Cyclic workflows
- Missing inputs
- Division by zero
- Unknown operations
- Missing required fields

All errors include descriptive messages.

## Design Principles

- **Simplicity**: Single file, straightforward logic
- **Readability**: Clear function names, simple flow
- **Testability**: Pure functions, no hidden state
- **Maintainability**: Short functions, focused responsibilities

## Limitations

This is an MVP implementation focused on simplicity:
- Single-threaded execution only
- No persistence or state management
- No HTTP endpoints (library only)
- Limited to basic arithmetic operations

## Future Enhancements

- HTTP API endpoints
- More operation types
- Conditional branching
- Parallel execution
- State persistence
- Workflow validation API
