# Backend Workflow Engine

A Go library for parsing and executing JSON workflow payloads from the frontend workflow builder.

## Features

- **JSON Parsing**: Parse workflow JSON payloads with validation
- **Workflow Execution**: Execute workflows with topological sorting for correct dependency order
- **Node Types Support**:
  - Number nodes (input values)
  - Operation nodes (add, subtract, multiply, divide)
  - Visualization nodes (text, table output)
- **Error Handling**: Comprehensive error handling with validation
- **Circular Dependency Detection**: Detects and reports circular dependencies

## Installation

```bash
go get github.com/yesoreyeram/thaiyyal/backend/workflow
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/yesoreyeram/thaiyyal/backend/workflow"
)

func main() {
	// Example workflow JSON from frontend
	jsonData := []byte(`{
		"nodes": [
			{"id": "1", "data": {"value": 10, "label": "Node 1"}},
			{"id": "2", "data": {"value": 5, "label": "Node 2"}},
			{"id": "3", "data": {"op": "add", "label": "Node 3 (op)"}},
			{"id": "4", "data": {"mode": "text", "label": "Node 4 (viz)"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-4", "source": "3", "target": "4"}
		]
	}`)

	// Execute the workflow
	result, err := workflow.ExecuteWorkflow(jsonData)
	if err != nil {
		log.Fatalf("Error executing workflow: %v", err)
	}

	// Access results
	fmt.Printf("Operation result: %v\n", result.Results["3"].Value) // 15
	fmt.Printf("Visualization output: %v\n", result.Output)         // "Result: 15"
}
```

## API Usage

### Parse Workflow

```go
parser := workflow.NewParser()
wf, err := parser.Parse(jsonData)
if err != nil {
	// handle error
}
```

### Execute Workflow

```go
engine := workflow.NewEngine(wf)
result, err := engine.Execute()
if err != nil {
	// handle error
}
```

### One-Step Execution

```go
result, err := workflow.ExecuteWorkflow(jsonData)
```

## Workflow JSON Format

The workflow JSON follows this structure:

```json
{
  "nodes": [
    {
      "id": "unique-id",
      "type": "numberNode|opNode|vizNode",
      "data": {
        "value": 10,        // for number nodes
        "op": "add",        // for operation nodes: add, subtract, multiply, divide
        "mode": "text",     // for viz nodes: text, table
        "label": "Node 1"   // optional label
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

## Node Types

### Number Node
Provides input values to the workflow.

```json
{
  "id": "1",
  "type": "numberNode",
  "data": {
    "value": 42.5
  }
}
```

### Operation Node
Performs arithmetic operations on input values.

Supported operations:
- `add`: Sum all inputs
- `subtract`: Subtract subsequent inputs from first
- `multiply`: Multiply all inputs
- `divide`: Divide first input by subsequent inputs

```json
{
  "id": "2",
  "type": "opNode",
  "data": {
    "op": "add"
  }
}
```

### Visualization Node
Formats and outputs results.

Supported modes:
- `text`: Simple text output format
- `table`: Structured table format

```json
{
  "id": "3",
  "type": "vizNode",
  "data": {
    "mode": "text"
  }
}
```

## Execution Result

The `ExecutionResult` contains:

```go
type ExecutionResult struct {
	Results map[string]*NodeResult // Map of nodeID to result
	Output  interface{}            // Final output from visualization node
}

type NodeResult struct {
	NodeID string
	Value  interface{} // Result value (typically float64 or string)
	Error  error       // Any error that occurred
}
```

## Error Handling

The library validates workflows and provides detailed error messages:

- Invalid JSON format
- Empty workflows
- Duplicate node IDs
- Invalid edge references
- Circular dependencies
- Division by zero
- Unknown node types or operations

## Testing

Run the test suite:

```bash
cd backend
go test -v ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Examples

### Simple Addition

```go
jsonData := []byte(`{
	"nodes": [
		{"id": "1", "data": {"value": 10}},
		{"id": "2", "data": {"value": 5}},
		{"id": "3", "data": {"op": "add"}}
	],
	"edges": [
		{"id": "e1-3", "source": "1", "target": "3"},
		{"id": "e2-3", "source": "2", "target": "3"}
	]
}`)

result, _ := workflow.ExecuteWorkflow(jsonData)
fmt.Println(result.Results["3"].Value) // Output: 15
```

### Complex Workflow

```go
jsonData := []byte(`{
	"nodes": [
		{"id": "1", "data": {"value": 10}},
		{"id": "2", "data": {"value": 5}},
		{"id": "3", "data": {"op": "multiply"}},
		{"id": "4", "data": {"value": 2}},
		{"id": "5", "data": {"op": "subtract"}},
		{"id": "6", "data": {"mode": "table"}}
	],
	"edges": [
		{"id": "e1-3", "source": "1", "target": "3"},
		{"id": "e2-3", "source": "2", "target": "3"},
		{"id": "e3-5", "source": "3", "target": "5"},
		{"id": "e4-5", "source": "4", "target": "5"},
		{"id": "e5-6", "source": "5", "target": "6"}
	]
}`)

result, _ := workflow.ExecuteWorkflow(jsonData)
// Node 3: 10 * 5 = 50
// Node 5: 50 - 2 = 48
// Node 6: Table output of 48
```

## Architecture

The workflow engine consists of:

1. **Parser**: Validates and parses JSON into workflow structures
2. **Engine**: Executes workflows with dependency resolution
3. **Topological Sort**: Ensures nodes execute in correct order
4. **Node Executors**: Individual execution logic for each node type

## Future Enhancements

Potential additions for future versions:
- More node types (HTTP requests, database queries, etc.)
- Conditional logic and branching
- Parallel execution of independent nodes
- State persistence
- Workflow versioning
- More visualization formats

## License

Part of the Thaiyyal project.
