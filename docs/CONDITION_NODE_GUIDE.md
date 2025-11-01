# Condition Node - Advanced Expression Evaluation

## Overview

The Condition (IF) node now supports advanced expressions for production use cases. You can reference node outputs, workflow variables, context values, and use complex boolean logic.

## Syntax (No Template Delimiters)

The expression syntax is simple and direct - **no `{{}}` delimiters needed**.

## Expression Types

### 1. Simple Comparisons (Backward Compatible)

```javascript
">100"    // Greater than 100
"<50"     // Less than 50
">=10"    // Greater than or equal to 10
"<=20"    // Less than or equal to 20
"==5"     // Equal to 5
"!=0"     // Not equal to 0
"true"    // Always true
"false"   // Always false
```

### 2. Node References

Reference outputs from other nodes in the workflow:

```javascript
// Access node output value
"node.nodeId.value > 100"

// Access nested fields
"node.http1.output.status == 200"
"node.api1.output.data.count > 10"

// Compare two nodes
"node.input1.value > node.input2.value"
```

### 3. Variable References

Reference workflow variables:

```javascript
// Check variable value
"variables.counter > 10"
"variables.enabled == true"
"variables.username == 'admin'"

// Compare with node
"node.result.value > variables.threshold"
```

### 4. Context References

Reference context variables and constants:

```javascript
"context.maxRetries < 5"
"context.apiKey == 'valid'"
"node.count.value < context.limit"
```

### 5. Boolean Operators

Combine conditions with AND, OR, NOT:

```javascript
// AND operator
"node.a.value > 100 && node.b.value < 50"
"variables.enabled && node.status.value == 'ready'"

// OR operator
"node.result.value > 1000 || variables.override"
"node.error.value != null || variables.hasError"

// NOT operator
"!variables.disabled"
"!(node.status.value == 'error')"

// Complex combinations
"(node.a.value > 100 || node.b.value > 200) && !variables.paused"
```

### 6. String Operations

Work with text values:

```javascript
// String equality
"node.status.value == 'success'"
"variables.mode == 'production'"

// Contains function
"contains(node.message.value, 'error')"
"contains(node.log.value, 'WARNING')"
```

### 7. Input Reference

Reference the direct input to the condition node:

```javascript
"input > 100"
"input == 'active'"
```

## True/False Paths

The Condition node now has **two separate output handles**:

- **Green Handle (Top)**: True path - triggered when condition is met
- **Red Handle (Bottom)**: False path - triggered when condition is not met

Connect different nodes to each handle to create branching logic.

## Examples

### Example 1: HTTP Status Check

```json
{
  "nodes": [
    {"id": "http1", "type": "http", "data": {"url": "https://api.example.com"}},
    {"id": "check", "type": "condition", "data": {
      "condition": "node.http1.output.status == 200"
    }},
    {"id": "success", "type": "text_input", "data": {"text": "API call succeeded"}},
    {"id": "error", "type": "text_input", "data": {"text": "API call failed"}}
  ],
  "edges": [
    {"source": "http1", "target": "check"},
    {"source": "check", "target": "success", "sourceHandle": "true"},
    {"source": "check", "target": "error", "sourceHandle": "false"}
  ]
}
```

### Example 2: Variable Threshold Check

```json
{
  "nodes": [
    {"id": "var1", "type": "variable", "data": {"var_name": "threshold", "var_op": "set"}},
    {"id": "num1", "type": "number", "data": {"value": 150}},
    {"id": "check", "type": "condition", "data": {
      "condition": "node.num1.value > variables.threshold"
    }}
  ],
  "edges": [
    {"source": "var1", "target": "check"},
    {"source": "num1", "target": "check"}
  ]
}
```

### Example 3: Complex Boolean Logic

```json
{
  "nodes": [
    {"id": "temp", "type": "number", "data": {"value": 75}},
    {"id": "humidity", "type": "number", "data": {"value": 60}},
    {"id": "check", "type": "condition", "data": {
      "condition": "(node.temp.value > 70 && node.humidity.value > 50) || variables.override"
    }},
    {"id": "alert", "type": "text_input", "data": {"text": "Environmental alert"}},
    {"id": "normal", "type": "text_input", "data": {"text": "Conditions normal"}}
  ],
  "edges": [
    {"source": "temp", "target": "check"},
    {"source": "humidity", "target": "check"},
    {"source": "check", "target": "alert", "sourceHandle": "true"},
    {"source": "check", "target": "normal", "sourceHandle": "false"}
  ]
}
```

### Example 4: String Matching

```json
{
  "nodes": [
    {"id": "log", "type": "text_input", "data": {"text": "ERROR: Connection failed"}},
    {"id": "check", "type": "condition", "data": {
      "condition": "contains(node.log.value, 'ERROR')"
    }},
    {"id": "errorHandler", "type": "text_input", "data": {"text": "Handle error"}},
    {"id": "continue", "type": "text_input", "data": {"text": "Continue processing"}}
  ],
  "edges": [
    {"source": "log", "target": "check"},
    {"source": "check", "target": "errorHandler", "sourceHandle": "true"},
    {"source": "check", "target": "continue", "sourceHandle": "false"}
  ]
}
```

## Dependency Graph Integration

When you use node references in conditions (e.g., `node.id.value`), the workflow engine automatically:

1. Extracts the referenced node IDs
2. Adds them as dependencies in the execution graph
3. Ensures proper topological sorting
4. Detects circular references and throws errors

This means you don't need to manually create edges for dependencies - the engine handles it automatically!

## Error Handling

The expression evaluator provides clear error messages:

- **Unknown reference**: `"variable not found: counter"`
- **Invalid node**: `"node result not found: http1"`
- **Missing field**: `"field not found: status in node.http1"`
- **Type mismatch**: Automatic type conversion where possible

If an expression fails to evaluate, the condition node falls back to simple numeric comparison for backward compatibility.

## Performance

- **No parsing overhead**: Direct evaluation without tokenization
- **Type-safe**: Runtime type checking with conversions
- **Efficient**: Single-pass evaluation
- **Cached**: Node results are cached during execution

## Migration from Template Syntax

If you have existing conditions with `{{}}`, simply remove the delimiters:

**Before:**
```javascript
"{{ node.id.value > 100 }}"
```

**After:**
```javascript
"node.id.value > 100"
```

The simpler syntax is cleaner and easier to read.

## Best Practices

1. **Use descriptive node IDs**: `node.httpApiCall.output.status` is clearer than `node.n1.output.status`
2. **Test incrementally**: Build complex conditions step by step
3. **Use variables for thresholds**: Makes workflows more maintainable
4. **Add comments in node labels**: Document what each condition checks
5. **Leverage true/false paths**: Visual clarity in workflow design

## See Also

- [Node Types Reference](/docs/NODES.md)
- [Variable Node](/docs/NODES.md#variable-node)
- [Context Nodes](/docs/NODES.md#context-nodes)
- [Workflow Examples](/docs/EXAMPLES.md)
