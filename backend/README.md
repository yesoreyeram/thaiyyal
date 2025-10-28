# Thaiyyal Backend - Workflow Engine

A simple, easy-to-understand Go workflow execution engine that parses and executes JSON workflow payloads from the Thaiyyal frontend.

## Features

- **Simple & Readable**: Single file implementation (~700 lines)
- **Easy to Understand**: Straightforward code flow without complex patterns
- **JSON Payload Parsing**: Accepts workflow definitions as JSON
- **DAG Execution**: Uses topological sorting to execute nodes in correct order
- **Type Inference**: Automatically determines node types from data
- **Node Types**:
  - **Number Nodes**: Provide numeric input values
  - **Operation Nodes**: Perform arithmetic (add, subtract, multiply, divide)
  - **Visualization Nodes**: Format output for display (text, table)
  - **Text Input Nodes**: Provide text string inputs
  - **Text Operation Nodes**: Transform text (uppercase, lowercase, titlecase, camelcase, inversecase, concat, repeat)
  - **HTTP Nodes**: Execute HTTP GET requests and return response body
  - **Condition Nodes**: Evaluate conditions and pass through values (NEW ✨)
  - **For Each Nodes**: Iterate over array elements (NEW ✨)
  - **While Loop Nodes**: Loop while conditions are true (NEW ✨)
- **Cycle Detection**: Prevents execution of workflows with circular dependencies
- **Comprehensive Tests**: 79 test cases covering all functionality including 39 control flow tests

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
├── workflow.go                    # Main workflow engine (single file, ~700 lines)
├── workflow_test.go              # Standard tests (40 tests)
├── workflow_controlflow_test.go  # Control flow tests (39 tests)
├── examples/
│   ├── main.go                   # Example usage
│   └── looping_poc.go           # Looping patterns POC
├── README.md                     # This file
└── INTEGRATION.md                # Frontend integration guide
```

## Usage Examples

### Basic Arithmetic

```go
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
// Result: 15
```

### Text Transformations

```go
payload := `{
  "nodes": [
    {"id": "1", "data": {"text": "hello world"}},
    {"id": "2", "data": {"text_op": "camelcase"}},
    {"id": "3", "data": {"text_op": "inversecase"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"}
  ]
}`
// Result: "HELLOwORLD" (hello world → helloWorld → HELLOwORLD)
```

### Condition Node

```go
payload := `{
  "nodes": [
    {"id": "1", "data": {"value": 150}},
    {"id": "2", "type": "condition", "data": {"condition": ">100"}}
  ],
  "edges": [
    {"source": "1", "target": "2"}
  ]
}`
// Result: {"value": 150, "condition_met": true, "condition": ">100"}
```

### For Each Node

```go
// Process an array of items
payload := `{
  "nodes": [
    {"id": "1", "data": {"value": [1, 2, 3, 4, 5]}},
    {"id": "2", "type": "for_each", "data": {"max_iterations": 1000}}
  ],
  "edges": [
    {"source": "1", "target": "2"}
  ]
}`
// Result: {"items": [1,2,3,4,5], "count": 5, "iterations": 5}
```

### While Loop Node

```go
payload := `{
  "nodes": [
    {"id": "1", "data": {"value": 150}},
    {"id": "2", "type": "while_loop", "data": {
      "condition": "<100",
      "max_iterations": 10
    }}
  ],
  "edges": [
    {"source": "1", "target": "2"}
  ]
}`
// Result: {"final_value": 150, "iterations": 0, "condition": "<100"}
// (0 iterations because 150 is not < 100)
```

### Complex Workflow with Condition

```go
payload := `{
  "nodes": [
    {"id": "1", "data": {"value": 30}},
    {"id": "2", "data": {"value": 70}},
    {"id": "3", "data": {"op": "add"}},
    {"id": "4", "type": "condition", "data": {"condition": ">=100"}},
    {"id": "5", "data": {"mode": "text"}}
  ],
  "edges": [
    {"source": "1", "target": "3"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "4"},
    {"source": "4", "target": "5"}
  ]
}`
// Flow: 30 + 70 = 100 → condition (100 >= 100 = true) → visualization
```

## Payload Format

The engine accepts JSON payloads with this structure:

```json
{
  "nodes": [
    {
      "id": "unique-id",
      "type": "number",          // optional: auto-detected from data
      "data": {
        // Number node:
        "value": 10,
        
        // Operation node:
        "op": "add",             // add, subtract, multiply, divide
        
        // Visualization node:
        "mode": "text",          // text, table
        
        // Text input node:
        "text": "Hello",
        
        // Text operation node:
        "text_op": "uppercase",  // uppercase, lowercase, titlecase, camelcase, 
                                 // inversecase, concat, repeat
        "separator": " ",        // for concat operation (optional)
        "repeat_n": 3,          // for repeat operation (required)
        
        // HTTP node:
        "url": "https://api.example.com/data",
        
        // Condition node:
        "condition": ">100",    // >N, <N, >=N, <=N, ==N, !=N, true, false
        
        // For Each node:
        "max_iterations": 1000, // optional, default: 1000
        
        // While Loop node:
        "condition": "<10",     // required
        "max_iterations": 100,  // optional, default: 100
        
        "label": "My Node"      // optional label for all nodes
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
3. **Topological Sort**: Use Kahn's algorithm to find execution order (DAG)
4. **Execute Nodes**: Process each node in order:
   - **Number nodes**: Return their numeric value
   - **Operation nodes**: Compute arithmetic results from inputs
   - **Visualization nodes**: Format outputs for display
   - **Text input nodes**: Return their text value
   - **Text operation nodes**: Transform text inputs
   - **HTTP nodes**: Execute HTTP GET requests
   - **Condition nodes**: Evaluate conditions and pass through values
   - **For Each nodes**: Process array elements
   - **While Loop nodes**: Loop while condition is true
5. **Return Results**: Collect all node results and determine final output

## Testing

Run all tests:
```bash
cd backend
go test -v
```

Run specific test suite:
```bash
# Standard tests (40 tests)
go test -v -run TestSimple
go test -v -run TestText
go test -v -run TestHTTP

# Control flow tests (39 tests)
go test -v -run TestCondition
go test -v -run TestForEach
go test -v -run TestWhileLoop
go test -v -run TestComplex
```

### Test Coverage (79 Total Tests)

**Basic Functionality** (11 tests):
- ✅ Engine creation and JSON parsing
- ✅ Simple addition workflow
- ✅ All arithmetic operations (add, subtract, multiply, divide)
- ✅ Division by zero error handling
- ✅ Complete workflow with visualization
- ✅ Multiple chained operations
- ✅ Cycle detection
- ✅ Missing input error handling
- ✅ Explicit node types
- ✅ Type inference
- ✅ Frontend payload compatibility

**Text Operations** (10 tests):
- ✅ Text input nodes
- ✅ Uppercase, lowercase, titlecase transformations
- ✅ Camelcase and inversecase transformations
- ✅ Chained text operations
- ✅ Non-text input error handling
- ✅ Explicit text node types
- ✅ Complex text transformation chains

**HTTP Nodes** (8 tests):
- ✅ Successful HTTP GET requests
- ✅ Error status code handling (404, 500, etc.)
- ✅ Invalid URL handling
- ✅ HTTP to text operation chaining
- ✅ HTTP error propagation
- ✅ Multiple chained text operations from HTTP
- ✅ Various HTTP status codes (200, 201, 204, 400, 404, 500)

**Text Operations - Concat & Repeat** (11 tests):
- ✅ Concat with 2 inputs
- ✅ Concat with custom separator
- ✅ Concat with multiple inputs
- ✅ Concat with non-text input error
- ✅ Repeat with positive count
- ✅ Repeat with zero count
- ✅ Repeat missing repeat_n error
- ✅ Repeat with negative count error
- ✅ Chained concat and repeat
- ✅ HTTP to concat integration
- ✅ Complex HTTP → text operation workflows

**Control Flow - Condition Nodes** (11 tests):
- ✅ Greater than condition (>N)
- ✅ Less than condition (<N)
- ✅ Greater than or equal (>=N)
- ✅ Equality condition (==N)
- ✅ Not equal condition (!=N)
- ✅ Boolean true/false conditions
- ✅ Condition with arithmetic operation integration
- ✅ Condition with text input
- ✅ Condition with text operations
- ✅ Missing condition error handling
- ✅ Multiple conditions in series

**Control Flow - For Each Nodes** (4 tests):
- ✅ For each with array input
- ✅ Max iterations limit enforcement
- ✅ Non-array input error handling
- ✅ For each with text array

**Control Flow - While Loop Nodes** (4 tests):
- ✅ While loop with condition
- ✅ Max iterations limit enforcement
- ✅ False condition immediate termination
- ✅ Missing condition error handling

**Control Flow - Integration Tests** (9 tests):
- ✅ Complex workflow: arithmetic → condition → visualization
- ✅ Condition chained with text operations
- ✅ HTTP with condition validation
- ✅ Control flow nodes without inputs (error cases)
- ✅ Multiple conditions in series
- ✅ Various node type integrations

## Code Overview

### Main Functions

**NewEngine(payloadJSON []byte) (*Engine, error)**
- Creates a new engine from JSON payload
- Returns error if JSON is invalid

**Execute() (*Result, error)**
- Executes the workflow
- Returns results or error

### Internal Functions

**Core Execution**:
- `inferNodeTypes()` - Determines types from node data
- `topologicalSort()` - Orders nodes using Kahn's algorithm (DAG)
- `executeNode(node)` - Dispatches to specific node executor
- `getNodeInputs(nodeID)` - Gets inputs from predecessor nodes
- `getFinalOutput()` - Finds terminal node output

**Node Executors**:
- `executeNumberNode(node)` - Returns number value
- `executeOperationNode(node)` - Performs arithmetic operations
- `executeVisualizationNode(node)` - Formats output for display
- `executeTextInputNode(node)` - Returns text value
- `executeTextOperationNode(node)` - Transforms text (7 operations)
- `executeHTTPNode(node)` - Executes HTTP GET requests
- `executeConditionNode(node)` - Evaluates conditions (NEW ✨)
- `executeForEachNode(node)` - Processes array iterations (NEW ✨)
- `executeWhileLoopNode(node)` - Loops with condition (NEW ✨)

**Helpers**:
- `evaluateCondition(condition, value)` - Evaluates condition expressions
- `toTitleCase(s)` - Converts to title case
- `toCamelCase(s)` - Converts to camelCase
- `inverseCase(s)` - Swaps character case

## Error Handling

The engine handles:
- Invalid JSON syntax
- Cyclic workflows
- Missing inputs for nodes
- Division by zero
- Unknown operations
- Missing required fields (e.g., condition, url)
- Type mismatches (e.g., non-text input to text operations)
- HTTP errors (4xx, 5xx status codes)
- Invalid URLs
- Iteration limit exceeded (for_each, while_loop)

All errors include descriptive messages.

## Design Principles

- **Simplicity**: Single file implementation, straightforward logic
- **Readability**: Clear function names, simple flow
- **Testability**: Pure functions, comprehensive test coverage
- **Maintainability**: Focused functions, clear separation of concerns

## Limitations & Future Work

**Current Limitations**:
- Control flow nodes (for_each, while_loop) don't yet execute sub-workflows
- Single-threaded execution only
- No persistence or state management
- No HTTP endpoints (library only)

**Future Enhancements**:
- Full sub-workflow execution in loops
- HTTP API endpoints
- Switch/case node for multi-way branching
- Parallel execution node
- Join/merge node for combining multiple paths
- More data transformation nodes
- State persistence
- Workflow validation API
- Real-time streaming support

## Learn More

- See [`docs/NODES.md`](../docs/NODES.md) for complete node type reference
- See [`INTEGRATION.md`](INTEGRATION.md) for frontend integration guide
- See [`examples/looping_poc.go`](examples/looping_poc.go) for looping patterns

