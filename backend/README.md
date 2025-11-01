# Thaiyyal Backend - Workflow Engine

A simple, easy-to-understand Go workflow execution engine that parses and executes JSON workflow payloads from the Thaiyyal frontend.

## Features

- **Simple & Readable**: Distributed across focused files (~1500 lines total)
- **Easy to Understand**: Straightforward code flow without complex patterns
- **JSON Payload Parsing**: Accepts workflow definitions as JSON
- **DAG Execution**: Uses topological sorting to execute nodes in correct order
- **⚡ Parallel Execution**: Execute independent nodes concurrently for 2-10x speedup
- **🔌 Extensible**: Add custom node types with your own executors
- **👀 Observable**: Monitor workflow execution with the Observer pattern (NEW!)
- **🛡️ Secure by Default**: Comprehensive protection limits prevent resource exhaustion
- **Type Inference**: Automatically determines node types from data
- **Node Types** (25+ types):
  - **Number Nodes**: Provide numeric input values
  - **Operation Nodes**: Perform arithmetic (add, subtract, multiply, divide)
  - **Visualization Nodes**: Format output for display (text, table)
  - **Text Input Nodes**: Provide text string inputs
  - **Text Operation Nodes**: Transform text (uppercase, lowercase, titlecase, camelcase, inversecase, concat, repeat)
  - **HTTP Nodes**: Execute HTTP GET requests and return response body
  - **Condition Nodes**: Evaluate conditions and pass through values
  - **For Each Nodes**: Iterate over array elements
  - **While Loop Nodes**: Loop while conditions are true
  - **Variable Nodes**: Store and retrieve values across workflow
  - **Extract Nodes**: Extract fields from objects
  - **Transform Nodes**: Transform data structures (to_array, to_object, flatten, keys, values)
  - **Accumulator Nodes**: Accumulate values over time (sum, product, concat, array, count)
  - **Counter Nodes**: Simple counter with increment/decrement/reset
  - **Switch Nodes**: Multi-way branching based on value or condition
  - **Parallel Nodes**: Execute multiple branches concurrently with concurrency control
  - **Join Nodes**: Combine outputs from multiple nodes with strategies (all/any/first)
  - **Split Nodes**: Split single input to multiple output paths
  - **Delay Nodes**: Pause execution for specified duration
  - **Cache Nodes**: Get/set cached values with TTL and LRU eviction
  - **Retry Nodes**: Retry with exponential/linear/constant backoff
  - **Try-Catch Nodes**: Error handling with fallback logic
  - **Timeout Nodes**: Enforce time limits on operations
  - **Custom Nodes**: Add your own node types with custom executors
- **State Management**: Variables, accumulators, counters, and cache for stateful workflows
- **Cycle Detection**: Prevents execution of workflows with circular dependencies
- **Thread Safety**: Full mutex protection for concurrent execution
- **Comprehensive Tests**: 190+ test cases (including custom executor tests)
  - 40 standard tests
  - 39 control flow tests
  - 17 state/memory tests
  - 46+ advanced control flow tests
  - 11 parallel execution tests
  - 13 context and validation tests
  - 14 custom executor tests (NEW!)

## Quick Start

### Sequential Execution (Standard)

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
    
    // Sequential execution
    result, err := engine.Execute()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %v\n", result.FinalOutput)
}
```

### Parallel Execution (NEW ⚡)

```go
// Enable parallel execution for 2-10x speedup on branching workflows
config := workflow.DefaultParallelConfig()
result, err := engine.ExecuteWithParallelism(config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %v\n", result.FinalOutput)

// Custom configuration with concurrency limit
config := workflow.ParallelExecutionConfig{
    MaxConcurrency: 4,    // Limit to 4 concurrent nodes
    EnableParallel: true,
}
result, err := engine.ExecuteWithParallelism(config)
```

**Learn More**: See [PARALLEL_EXECUTION.md](PARALLEL_EXECUTION.md) for detailed documentation.

### Custom Node Types (NEW 🔌)

Extend the workflow engine with your own custom node executors:

```go
package main

import (
    "fmt"
    "github.com/yesoreyeram/thaiyyal/backend"
)

// Define a custom executor
type ReverseStringExecutor struct{}

func (e *ReverseStringExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    // Increment execution counter for protection
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    inputs := ctx.GetNodeInputs(node.ID)
    str := inputs[0].(string)
    
    // Reverse the string
    runes := []rune(str)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    
    return string(runes), nil
}

func (e *ReverseStringExecutor) NodeType() types.NodeType {
    return types.NodeType("reverse_string")
}

func (e *ReverseStringExecutor) Validate(node workflow.Node) error {
    return nil
}

func main() {
    // Register your custom executor
    registry := workflow.DefaultRegistry()
    registry.MustRegister(&ReverseStringExecutor{})
    
    payload := `{
        "nodes": [
            {"id": "1", "data": {"text": "Hello"}},
            {"id": "2", "type": "reverse_string", "data": {}}
        ],
        "edges": [{"source": "1", "target": "2"}]
    }`
    
    // Create engine with custom registry
    engine, _ := workflow.NewEngineWithRegistry(
        []byte(payload),
        workflow.DefaultConfig(),
        registry,
    )
    
    result, _ := engine.Execute()
    fmt.Println(result.FinalOutput)  // Output: "olleH"
}
```

**Features:**
- ✅ All protection limits automatically apply to custom nodes
- ✅ Mix custom and built-in nodes in the same workflow
- ✅ Full access to workflow state (variables, cache, context)
- ✅ HTTP call tracking and limits
- ✅ Comprehensive validation support

**Learn More**: See [CUSTOM_NODES.md](CUSTOM_NODES.md) for the complete guide with examples.

### Observer Pattern (NEW 👀)

Monitor workflow execution with custom observers:

```go
package main

import (
    "context"
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
    
    // Register console observer for debugging
    observer := workflow.NewConsoleObserver()
    engine.RegisterObserver(observer)
    
    // Or create a custom observer
    customObserver := &MyCustomObserver{}
    engine.RegisterObserver(customObserver)
    
    // Execute workflow (observers receive events asynchronously)
    result, err := engine.Execute()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %v\n", result.FinalOutput)
}

// Custom observer implementation
type MyCustomObserver struct{}

func (o *MyCustomObserver) OnEvent(ctx context.Context, event workflow.Event) {
    // Handle workflow and node events
    fmt.Printf("[%s] %s - Node: %s, Status: %s\n", 
        event.Type, event.ExecutionID, event.NodeID, event.Status)
}
```

**Features:**
- ✅ Asynchronous execution (observers never block workflow)
- ✅ Multiple observers support
- ✅ Built-in console observer and custom logger
- ✅ Complete event metadata (timing, status, errors)
- ✅ Integration with monitoring systems (DataDog, Prometheus, etc.)
- ✅ Panic recovery (observer failures don't affect workflow)

**Learn More**: See [OBSERVER_PATTERN.md](OBSERVER_PATTERN.md) for the complete guide with examples.

### Running Examples

```bash
cd backend/examples
go run main.go

# Run custom nodes examples
cd backend/examples/custom_nodes
go run main.go
```

## File Structure

```
backend/
├── workflow.go                    # Main workflow engine (single file, ~1100 lines)
├── workflow_test.go              # Standard tests (40 tests)
├── workflow_controlflow_test.go  # Control flow tests (39 tests)
├── workflow_state_test.go        # State/memory tests (17 tests)
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
        
        // Variable node:
        "var_name": "myvar",    // required - name of the variable
        "var_op": "set",        // required - "set" or "get"
        
        // Extract node:
        "field": "name",        // extract single field
        // OR
        "fields": ["name", "email"],  // extract multiple fields
        
        // Transform node:
        "transform_type": "to_array",  // required - to_array, to_object, flatten, keys, values
        
        // Accumulator node:
        "accum_op": "sum",      // required - sum, product, concat, array, count
        "initial_value": 0,     // optional - starting value
        
        // Counter node:
        "counter_op": "increment",  // required - increment, decrement, reset, get
        "delta": 1,             // optional - amount to increment/decrement
        "initial_value": 0,     // optional - reset value
        
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

### State & Memory Node Usage Examples

#### Variable Node Example
```json
{
  "nodes": [
    {"id": "1", "data": {"value": 100}},
    {"id": "2", "data": {"var_name": "result", "var_op": "set"}},
    {"id": "3", "data": {"var_name": "result", "var_op": "get"}},
    {"id": "4", "type": "extract", "data": {"field": "value"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "4"}
  ]
}
```

#### Accumulator Example
```json
{
  "nodes": [
    {"id": "1", "data": {"value": 10}},
    {"id": "2", "type": "accumulator", "data": {"accum_op": "sum"}},
    {"id": "3", "data": {"value": 20}},
    {"id": "4", "type": "accumulator", "data": {"accum_op": "sum"}},
    {"id": "5", "data": {"value": 30}},
    {"id": "6", "type": "accumulator", "data": {"accum_op": "sum"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "4"},
    {"source": "4", "target": "5"},
    {"source": "5", "target": "6"}
  ]
}
// Result: {"operation": "sum", "value": 60}
```

#### Counter Example
```json
{
  "nodes": [
    {"id": "1", "type": "counter", "data": {"counter_op": "increment", "delta": 5}},
    {"id": "2", "type": "counter", "data": {"counter_op": "increment"}},
    {"id": "3", "type": "counter", "data": {"counter_op": "get"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"}
  ]
}
// Result: {"operation": "get", "value": 6}
```

#### Transform Example
```json
{
  "nodes": [
    {"id": "1", "data": {"value": 10}},
    {"id": "2", "data": {"value": 20}},
    {"id": "3", "data": {"value": 30}},
    {"id": "4", "type": "transform", "data": {"transform_type": "to_array"}}
  ],
  "edges": [
    {"source": "1", "target": "4"},
    {"source": "2", "target": "4"},
    {"source": "3", "target": "4"}
  ]
}
// Result: [10, 20, 30]
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
   - **Variable nodes**: Store or retrieve values from workflow state
   - **Extract nodes**: Extract specific fields from object inputs
   - **Transform nodes**: Transform data structures (arrays, objects, etc.)
   - **Accumulator nodes**: Accumulate values across multiple executions
   - **Counter nodes**: Maintain a counter with increment/decrement operations
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

- See [`OBSERVER_PATTERN.md`](OBSERVER_PATTERN.md) for workflow monitoring and observability guide (NEW!)
- See [`CUSTOM_NODES.md`](CUSTOM_NODES.md) for custom node executor guide
- See [`PROTECTION.md`](PROTECTION.md) for security and protection limits
- See [`docs/NODES.md`](../docs/NODES.md) for complete node type reference
- See [`INTEGRATION.md`](INTEGRATION.md) for frontend integration guide
- See [`PARALLEL_EXECUTION.md`](PARALLEL_EXECUTION.md) for parallel execution guide
- See [`examples/custom_nodes/`](examples/custom_nodes/) for custom node examples
- See [`examples/looping_poc.go`](examples/looping_poc.go) for looping patterns

