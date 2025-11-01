# Custom Node Types - Developer Guide

This guide explains how to extend Thaiyyal's workflow engine with custom node executors. Custom nodes allow you to add domain-specific functionality while maintaining all the security protections and workflow orchestration capabilities of the engine.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Implementing Custom Executors](#implementing-custom-executors)
- [Registration API](#registration-api)
- [Security and Protection](#security-and-protection)
- [Best Practices](#best-practices)
- [Examples](#examples)
- [What Could Go Wrong](#what-could-go-wrong)
- [Testing Custom Executors](#testing-custom-executors)

## Overview

Thaiyyal provides 25 built-in node types for common workflow operations. However, you may need custom functionality for your specific domain. Custom node executors allow you to:

- **Extend functionality**: Add domain-specific operations (e.g., call proprietary APIs, custom data transformations)
- **Integrate systems**: Connect to databases, message queues, or external services
- **Maintain security**: All protection limits (execution limits, HTTP call limits, timeout, etc.) automatically apply
- **Reuse orchestration**: Leverage Thaiyyal's workflow engine without implementing control flow yourself

### What You Can Do

✅ Add custom data processing nodes  
✅ Integrate with external APIs and services  
✅ Implement domain-specific transformations  
✅ Mix custom nodes with built-in nodes  
✅ Apply all protection limits to custom nodes  

### What You Cannot Do (By Design)

❌ Create custom orchestration nodes (use built-in: `for_each`, `while_loop`, `switch`, `parallel`)  
❌ Bypass security and protection limits  
❌ Access workflow engine internals directly  
❌ Modify workflow state outside the provided API  

## Quick Start

Here's a minimal example of creating and using a custom executor:

```go
package main

import (
    "fmt"
    "github.com/yesoreyeram/thaiyyal/backend"
    "github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Step 1: Implement the NodeExecutor interface
type ReverseStringExecutor struct{}

func (e *ReverseStringExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    // Increment execution counter for protection
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    // Get inputs from connected nodes
    inputs := ctx.GetNodeInputs(node.ID)
    if len(inputs) == 0 {
        return nil, fmt.Errorf("reverse_string requires input")
    }
    
    // Process the input
    str := inputs[0].(string)
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
    return nil // No special validation needed
}

// Step 2: Register and use your custom executor
func main() {
    // Get default registry (includes all built-in nodes)
    registry := workflow.DefaultRegistry()
    
    // Register your custom executor
    registry.MustRegister(&ReverseStringExecutor{})
    
    // Create workflow using your custom node
    payload := `{
        "nodes": [
            {"id": "1", "data": {"text": "Hello"}},
            {"id": "2", "type": "reverse_string", "data": {}}
        ],
        "edges": [
            {"source": "1", "target": "2"}
        ]
    }`
    
    // Create engine with custom registry
    engine, _ := workflow.NewEngineWithRegistry(
        []byte(payload),
        workflow.DefaultConfig(),
        registry,
    )
    
    // Execute workflow
    result, _ := engine.Execute()
    fmt.Println(result.FinalOutput) // Output: "olleH"
}
```

## Implementing Custom Executors

### The NodeExecutor Interface

All custom executors must implement three methods:

```go
type NodeExecutor interface {
    // Execute runs the node with given context
    Execute(ctx ExecutionContext, node types.Node) (interface{}, error)
    
    // NodeType returns the type identifier for this executor
    NodeType() types.NodeType
    
    // Validate checks if node configuration is valid
    Validate(node types.Node) error
}
```

### Execute Method

The `Execute` method is where your custom logic lives:

```go
func (e *MyExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    // 1. ALWAYS increment execution counter (for protection limits)
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    // 2. Get inputs from connected nodes
    inputs := ctx.GetNodeInputs(node.ID)
    
    // 3. Access node configuration from node.Data
    config := node.Data.SomeField
    
    // 4. Perform your custom logic
    result := processData(inputs, config)
    
    // 5. Return the result (can be any type)
    return result, nil
}
```

### NodeType Method

Returns a unique identifier for your node type:

```go
func (e *MyExecutor) NodeType() types.NodeType {
    return types.NodeType("my_custom_node")
}
```

**Important**: Choose unique names that won't conflict with built-in types or other custom types.

### Validate Method

Validates the node configuration before execution:

```go
func (e *MyExecutor) Validate(node workflow.Node) error {
    // Check required fields
    if node.Data.RequiredField == nil {
        return fmt.Errorf("my_custom_node requires 'required_field'")
    }
    
    // Check constraints
    if node.Data.Count != nil && *node.Data.Count < 0 {
        return fmt.Errorf("count must be non-negative")
    }
    
    return nil
}
```

## Registration API

### Creating a Registry

**Option 1: Start with defaults (recommended)**
```go
registry := workflow.DefaultRegistry()  // Includes all 25 built-in nodes
registry.MustRegister(&MyCustomExecutor{})
```

**Option 2: Build from scratch**
```go
registry := workflow.NewRegistry()  // Empty registry
registry.MustRegister(&MyCustomExecutor{})
// Note: You'll need to register built-in executors if you want them
```

### Registration Methods

**MustRegister** - Panics on error (use during initialization):
```go
registry.MustRegister(&MyExecutor{})
```

**Register** - Returns error (use when error handling is needed):
```go
if err := registry.Register(&MyExecutor{}); err != nil {
    return fmt.Errorf("failed to register executor: %w", err)
}
```

### Creating Engine with Custom Registry

```go
engine, err := workflow.NewEngineWithRegistry(
    payloadJSON,           // Workflow definition
    workflow.DefaultConfig(),  // Configuration
    registry,              // Your custom registry
)
```

## Security and Protection

### Automatic Protection Limits

All protection limits automatically apply to custom nodes:

| Protection | Description | How It Applies |
|------------|-------------|----------------|
| **MaxNodeExecutions** | Limit total node executions | Engine tracks all nodes |
| **MaxExecutionTime** | Workflow timeout | Entire workflow including custom nodes |
| **MaxHTTPCallsPerExec** | HTTP call limit | Custom nodes must call `IncrementHTTPCall()` |
| **MaxStringLength** | String size limit | Applied to custom node outputs |
| **MaxArrayLength** | Array size limit | Applied to custom node outputs |
| **MaxContextDepth** | Nesting depth limit | Applied to custom node outputs |
| **MaxVariables** | Variable count limit | Custom nodes using variables |

### Calling Protection Methods

**Node Execution Tracking (REQUIRED)**

```go
func (e *MyExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    // ALWAYS call this at the start
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    // For nodes that iterate internally, call again per iteration
    for i := 0; i < iterations; i++ {
        if err := ctx.IncrementNodeExecution(); err != nil {
            return nil, err // Limit exceeded
        }
        // Process iteration
    }
    
    return result, nil
}
```

**HTTP Call Tracking (if making HTTP calls)**

```go
func (e *APIExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    // Increment HTTP counter before making request
    if err := ctx.IncrementHTTPCall(); err != nil {
        return nil, err // HTTP limit exceeded
    }
    
    // Make HTTP request
    resp, err := http.Get(url)
    // ...
    
    return result, nil
}
```

### Input Validation

Always validate inputs to prevent errors:

```go
func (e *MyExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    inputs := ctx.GetNodeInputs(node.ID)
    
    // Check input count
    if len(inputs) == 0 {
        return nil, fmt.Errorf("my_node requires at least one input")
    }
    
    // Check input type
    value, ok := inputs[0].(float64)
    if !ok {
        return nil, fmt.Errorf("expected number input, got %T", inputs[0])
    }
    
    // Check input range/constraints
    if value < 0 {
        return nil, fmt.Errorf("input must be non-negative")
    }
    
    return processValue(value), nil
}
```

## Best Practices

### 1. Always Increment Execution Counter

```go
// ✅ GOOD
func (e *MyExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    // ... rest of logic
}

// ❌ BAD - Missing execution counter
func (e *MyExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    // ... logic without incrementing counter
}
```

### 2. Use Type Assertions Safely

```go
// ✅ GOOD - Check type before using
value, ok := input.(string)
if !ok {
    return nil, fmt.Errorf("expected string, got %T", input)
}

// ❌ BAD - Assumes type without checking
value := input.(string)  // Panics if wrong type
```

### 3. Return Descriptive Errors

```go
// ✅ GOOD - Clear error messages
return nil, fmt.Errorf("json_parser: invalid JSON at line %d: %w", line, err)

// ❌ BAD - Vague error
return nil, fmt.Errorf("error")
```

### 4. Validate in Both Methods

```go
// Validate method catches errors before execution
func (e *MyExecutor) Validate(node workflow.Node) error {
    if node.Data.URL == nil {
        return fmt.Errorf("url is required")
    }
    return nil
}

// Execute method handles runtime errors
func (e *MyExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    // URL is guaranteed to exist (validated), but may be invalid
    resp, err := http.Get(*node.Data.URL)
    if err != nil {
        return nil, fmt.Errorf("HTTP request failed: %w", err)
    }
    // ...
}
```

### 5. Use Workflow State Appropriately

```go
// Get variable
value, err := ctx.GetVariable("myvar")
if err != nil {
    return nil, err
}

// Set variable
if err := ctx.SetVariable("result", computedValue); err != nil {
    return nil, err
}

// Use cache for expensive operations
if cached, ok := ctx.GetCache("expensive_key"); ok {
    return cached, nil
}
result := expensiveOperation()
ctx.SetCache("expensive_key", result, 5*time.Minute)
return result, nil
```

### 6. Document Your Custom Nodes

```go
// WeatherAPIExecutor fetches current weather for a city.
//
// Inputs:
//   - City name (string)
//
// Configuration:
//   - api_key (string): Weather API key
//   - units (string): "celsius" or "fahrenheit" (default: "celsius")
//
// Output:
//   - Weather data (map[string]interface{})
//
// Example:
//   {"id": "1", "type": "weather_api", "data": {"api_key": "xxx", "units": "celsius"}}
type WeatherAPIExecutor struct{}
```

## Examples

### Example 1: String Manipulation

```go
type Base64EncodeExecutor struct{}

func (e *Base64EncodeExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    inputs := ctx.GetNodeInputs(node.ID)
    if len(inputs) == 0 {
        return nil, fmt.Errorf("base64_encode requires input")
    }
    
    str, ok := inputs[0].(string)
    if !ok {
        return nil, fmt.Errorf("base64_encode requires string input")
    }
    
    encoded := base64.StdEncoding.EncodeToString([]byte(str))
    return encoded, nil
}

func (e *Base64EncodeExecutor) NodeType() types.NodeType {
    return types.NodeType("base64_encode")
}

func (e *Base64EncodeExecutor) Validate(node workflow.Node) error {
    return nil
}
```

### Example 2: Database Query

```go
type DatabaseQueryExecutor struct {
    db *sql.DB
}

func (e *DatabaseQueryExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    // Get query from node configuration
    query := ""
    if node.Data.URL != nil {  // Reusing URL field for query
        query = *node.Data.URL
    }
    
    // Execute query
    rows, err := e.db.QueryContext(context.Background(), query)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()
    
    // Convert rows to result
    var results []map[string]interface{}
    // ... row processing logic
    
    return results, nil
}

func (e *DatabaseQueryExecutor) NodeType() types.NodeType {
    return types.NodeType("db_query")
}

func (e *DatabaseQueryExecutor) Validate(node workflow.Node) error {
    if node.Data.URL == nil {
        return fmt.Errorf("query is required")
    }
    return nil
}
```

### Example 3: Multi-Input Aggregation

```go
type AverageExecutor struct{}

func (e *AverageExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    inputs := ctx.GetNodeInputs(node.ID)
    if len(inputs) == 0 {
        return nil, fmt.Errorf("average requires at least one input")
    }
    
    var sum float64
    var count int
    
    for _, input := range inputs {
        num, ok := input.(float64)
        if !ok {
            return nil, fmt.Errorf("average requires numeric inputs, got %T", input)
        }
        sum += num
        count++
    }
    
    return sum / float64(count), nil
}

func (e *AverageExecutor) NodeType() types.NodeType {
    return types.NodeType("average")
}

func (e *AverageExecutor) Validate(node workflow.Node) error {
    return nil
}
```

## What Could Go Wrong

### Security Risks

| Risk | Description | Mitigation |
|------|-------------|------------|
| **Resource Exhaustion** | Custom node runs forever or uses too much memory | Always call `IncrementNodeExecution()`, engine enforces limits |
| **SSRF Attacks** | Custom node makes requests to internal services | Validate URLs, use allowlists, call `IncrementHTTPCall()` |
| **Code Injection** | Custom node evaluates user input as code | Never use `eval()` or similar, validate all inputs |
| **Data Leakage** | Custom node logs sensitive data | Be careful with logging, use structured logging |
| **DoS via Loops** | Custom node has infinite loop | Use iteration limits, check execution counter |

### Common Mistakes

**1. Not Incrementing Execution Counter**
```go
// ❌ BAD - Protection limits won't work
func (e *BadExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    return "result", nil  // Missing IncrementNodeExecution()
}

// ✅ GOOD
func (e *GoodExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    return "result", nil
}
```

**2. Type Assertion Panics**
```go
// ❌ BAD - Panics if wrong type
value := input.(string)

// ✅ GOOD - Handles type mismatch gracefully
value, ok := input.(string)
if !ok {
    return nil, fmt.Errorf("expected string, got %T", input)
}
```

**3. Not Validating Inputs**
```go
// ❌ BAD - Assumes inputs exist and are correct type
func (e *BadExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    value := ctx.GetNodeInputs(node.ID)[0].(float64)  // Can panic!
    return value * 2, nil
}

// ✅ GOOD - Validates everything
func (e *GoodExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    inputs := ctx.GetNodeInputs(node.ID)
    if len(inputs) == 0 {
        return nil, fmt.Errorf("requires input")
    }
    
    value, ok := inputs[0].(float64)
    if !ok {
        return nil, fmt.Errorf("requires numeric input, got %T", inputs[0])
    }
    
    return value * 2, nil
}
```

**4. Modifying Input Data**
```go
// ❌ BAD - Modifies input (side effects)
func (e *BadExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    arr := ctx.GetNodeInputs(node.ID)[0].([]interface{})
    arr[0] = "modified"  // Modifies original!
    return arr, nil
}

// ✅ GOOD - Creates new data
func (e *GoodExecutor) Execute(ctx workflow.ExecutionContext, node workflow.Node) (interface{}, error) {
    if err := ctx.IncrementNodeExecution(); err != nil {
        return nil, err
    }
    
    arr := ctx.GetNodeInputs(node.ID)[0].([]interface{})
    newArr := make([]interface{}, len(arr))
    copy(newArr, arr)
    newArr[0] = "modified"
    return newArr, nil
}
```

### Debugging Tips

1. **Add logging**: Use structured logging to track execution
```go
log.Printf("[%s] Processing node %s with %d inputs", e.NodeType(), node.ID, len(inputs))
```

2. **Return detailed errors**: Include context in error messages
```go
return nil, fmt.Errorf("node %s (%s): failed to process input: %w", node.ID, e.NodeType(), err)
```

3. **Test with ValidationLimits**: Use strict limits during development
```go
config := workflow.ValidationLimits()  // Strict limits for testing
engine, _ := workflow.NewEngineWithRegistry(payload, config, registry)
```

## Testing Custom Executors

### Unit Tests

```go
func TestMyCustomExecutor(t *testing.T) {
    t.Run("basic execution", func(t *testing.T) {
        registry := workflow.DefaultRegistry()
        registry.MustRegister(&MyCustomExecutor{})
        
        payload := `{
            "nodes": [
                {"id": "1", "data": {"value": 10}},
                {"id": "2", "type": "my_custom", "data": {}}
            ],
            "edges": [
                {"source": "1", "target": "2"}
            ]
        }`
        
        engine, err := workflow.NewEngineWithRegistry(
            []byte(payload),
            workflow.DefaultConfig(),
            registry,
        )
        if err != nil {
            t.Fatal(err)
        }
        
        result, err := engine.Execute()
        if err != nil {
            t.Fatal(err)
        }
        
        expected := 20.0
        if result.FinalOutput != expected {
            t.Fatalf("expected %v, got %v", expected, result.FinalOutput)
        }
    })
    
    t.Run("protection limits", func(t *testing.T) {
        registry := workflow.DefaultRegistry()
        registry.MustRegister(&MyCustomExecutor{})
        
        config := workflow.DefaultConfig()
        config.MaxNodeExecutions = 1  // Very low limit
        
        payload := `{...}`  // Workflow that should exceed limit
        
        engine, _ := workflow.NewEngineWithRegistry([]byte(payload), config, registry)
        _, err := engine.Execute()
        
        if err == nil {
            t.Fatal("expected error for exceeding execution limit")
        }
    })
}
```

### Integration Tests

Test custom nodes with built-in nodes:

```go
func TestCustomWithBuiltIn(t *testing.T) {
    registry := workflow.DefaultRegistry()
    registry.MustRegister(&MyCustomExecutor{})
    
    payload := `{
        "nodes": [
            {"id": "1", "data": {"value": 5}},
            {"id": "2", "type": "my_custom", "data": {}},
            {"id": "3", "data": {"op": "add"}},
            {"id": "4", "data": {"value": 3}}
        ],
        "edges": [
            {"source": "1", "target": "2"},
            {"source": "2", "target": "3"},
            {"source": "4", "target": "3"}
        ]
    }`
    
    engine, _ := workflow.NewEngineWithRegistry(
        []byte(payload),
        workflow.DefaultConfig(),
        registry,
    )
    
    result, err := engine.Execute()
    if err != nil {
        t.Fatal(err)
    }
    
    // Verify result
    // ...
}
```

## Complete Working Example

See `backend/examples/custom_nodes/main.go` for a complete working example with:
- Simple string manipulation executor
- JSON path extraction executor
- HTTP API executor (with protection)
- Batch processing executor (with iteration tracking)
- Examples mixing custom and built-in nodes

Run the example:
```bash
cd backend/examples/custom_nodes
go run main.go
```

## Summary

Custom node executors provide a powerful way to extend Thaiyyal while maintaining security and consistency. Key takeaways:

1. **Implement three methods**: `Execute()`, `NodeType()`, `Validate()`
2. **Always call `IncrementNodeExecution()`** at the start of Execute()
3. **Validate all inputs** - never assume input types or values
4. **Use protection methods** - call `IncrementHTTPCall()` if making HTTP requests
5. **Start with `DefaultRegistry()`** - get all built-in nodes + add yours
6. **Test thoroughly** - write unit and integration tests
7. **Document your nodes** - explain inputs, outputs, and configuration

For more examples and details, see:
- `backend/examples/custom_nodes/main.go` - Working examples
- `backend/pkg/engine/custom_executor_test.go` - Comprehensive test suite
- `backend/PROTECTION.md` - Protection limits documentation
- `backend/README.md` - Workflow engine documentation
