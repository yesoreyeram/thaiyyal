# Workflow Validation

This document describes the workflow validation feature added to the Thaiyyal workflow engine.

## Overview

Workflow validation provides early error detection by checking the structure and content of workflow payloads before execution. This improves user experience by catching errors upfront with clear, actionable error messages.

## Quick Start

### Validating a Workflow

```go
import "github.com/yesoreyeram/thaiyyal/backend/workflow"

// Option 1: Validate without creating an engine
payload := []byte(`{...}`)
result, err := workflow.ValidatePayload(payload)
if err != nil {
    // Payload is not valid JSON
    log.Fatal(err)
}

if !result.Valid {
    // Workflow has validation errors
    for _, validationErr := range result.Errors {
        log.Printf("Validation error: %v", validationErr)
    }
    return
}

// Option 2: Validate after creating an engine
engine, err := workflow.NewEngine(payload)
if err != nil {
    log.Fatal(err)
}

result := engine.Validate()
if !result.Valid {
    log.Fatalf("Workflow validation failed: %v", result.Errors)
}
```

## Validation Checks

The validation layer performs three categories of checks:

### 1. Structural Validation

- **Empty workflow**: At least one node is required
- **Empty node IDs**: All nodes must have a non-empty ID
- **Duplicate node IDs**: Node IDs must be unique
- **Empty edges**: Edge source and target must be non-empty

### 2. Graph Validation

- **Invalid edge sources**: Edge source must reference an existing node
- **Invalid edge targets**: Edge target must reference an existing node
- **Self-referencing edges**: Nodes cannot connect to themselves
- **Cyclic dependencies**: Workflow must be a directed acyclic graph (DAG)

### 3. Node Data Validation

Each node type has specific data requirements that are validated:

#### Number Node
- **Required**: `value` field must be present
- Example error: `validation error on node '1' field 'data.value': number node requires a 'value' field`

#### Text Input Node
- **Required**: `text` field must be non-empty
- Example error: `validation error on node '1' field 'data.text': text input node requires a non-empty 'text' field`

#### Operation Node
- **Required**: `op` field must be present
- **Valid operations**: `add`, `subtract`, `multiply`, `divide`, `modulo`, `power`
- Example error: `validation error on node '1' field 'data.op': invalid operation 'invalid_op', must be one of: add, subtract, multiply, divide, modulo, power`

#### Text Operation Node
- **Required**: `text_op` field must be present
- **Valid operations**: `concat`, `uppercase`, `lowercase`, `trim`, `split`, `replace`, `length`, `substring`
- Example error: `validation error on node '1' field 'data.text_op': invalid text operation 'invalid', must be one of: concat, uppercase, lowercase, trim, split, replace, length, substring`

#### HTTP Node
- **Required**: `url` field must be non-empty
- Example error: `validation error on node '1' field 'data.url': HTTP node requires a non-empty 'url' field`

#### Condition Node
- **Required**: `condition` field must be non-empty
- Example error: `validation error on node '1' field 'data.condition': condition node requires a 'condition' field`

#### Variable Node
- **Required**: `var_name` field must be non-empty
- Example error: `validation error on node '1' field 'data.var_name': variable node requires a 'var_name' field`

#### Extract Node
- **Required**: Either `field` or `fields` must be specified
- Example error: `validation error on node '1' field 'data.field': extract node requires either 'field' or 'fields' to be specified`

#### Context Variable/Constant Node
- **Required**: `context_name` field must be non-empty
- **Context Constant**: Also requires `context_value` field
- Example error: `validation error on node '1' field 'data.context_name': context_variable node requires a 'context_name' field`

#### Retry Node
- **Required**: `max_attempts` must be at least 1 if specified
- Example error: `validation error on node '1' field 'data.max_attempts': retry node max_attempts must be at least 1`

#### Timeout Node
- **Required**: `timeout` field must be present
- Example error: `validation error on node '1' field 'data.timeout': timeout node requires a 'timeout' field`

## Validation Result

The `ValidationResult` struct contains:

```go
type ValidationResult struct {
    Valid  bool              // Whether the workflow is valid
    Errors []ValidationError // List of validation errors
}
```

Each `ValidationError` contains:

```go
type ValidationError struct {
    Field   string // The field that failed validation
    Message string // Human-readable error message
    NodeID  string // Optional: The node ID related to the error
}
```

## Best Practices

### 1. Always Validate Before Execution

```go
// Recommended pattern
result, err := workflow.ValidatePayload(payload)
if err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}

if !result.Valid {
    return fmt.Errorf("workflow validation failed: %v", result.Errors)
}

engine, err := workflow.NewEngine(payload)
if err != nil {
    return fmt.Errorf("failed to create engine: %w", err)
}

result, err := engine.Execute()
// ... handle execution
```

### 2. Provide User-Friendly Error Messages

```go
if !result.Valid {
    fmt.Println("Workflow validation errors:")
    for i, err := range result.Errors {
        if err.NodeID != "" {
            fmt.Printf("  %d. Node '%s': %s\n", i+1, err.NodeID, err.Message)
        } else {
            fmt.Printf("  %d. %s\n", i+1, err.Message)
        }
    }
}
```

### 3. Validate in Frontend Before Sending to Backend

Use the same validation logic in your frontend to catch errors before workflow execution:

```typescript
// Example: Validate in frontend
function validateWorkflow(nodes, edges) {
  const errors = [];
  
  // Check for empty workflow
  if (nodes.length === 0) {
    errors.push({ field: 'nodes', message: 'Workflow must contain at least one node' });
  }
  
  // Check for duplicate IDs
  const nodeIds = new Set();
  for (const node of nodes) {
    if (nodeIds.has(node.id)) {
      errors.push({ field: 'id', nodeId: node.id, message: `Duplicate node ID: ${node.id}` });
    }
    nodeIds.add(node.id);
  }
  
  // Check edges reference valid nodes
  for (const edge of edges) {
    if (!nodeIds.has(edge.source)) {
      errors.push({ field: 'edge.source', message: `Edge source '${edge.source}' does not exist` });
    }
    if (!nodeIds.has(edge.target)) {
      errors.push({ field: 'edge.target', message: `Edge target '${edge.target}' does not exist` });
    }
  }
  
  return errors;
}
```

## Error Examples

### Example 1: Cyclic Dependency

**Workflow:**
```json
{
  "nodes": [
    {"id": "1", "data": {"value": 10}},
    {"id": "2", "data": {"value": 20}},
    {"id": "3", "data": {"op": "add"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "1"}
  ]
}
```

**Error:**
```
validation error on field 'edges': workflow contains cycles (circular dependencies)
```

### Example 2: Invalid Node Data

**Workflow:**
```json
{
  "nodes": [
    {"id": "1", "data": {"op": "invalid_operation"}}
  ],
  "edges": []
}
```

**Error:**
```
validation error on node '1' field 'data.op': invalid operation 'invalid_operation', must be one of: add, subtract, multiply, divide, modulo, power
```

### Example 3: Missing Required Fields

**Workflow:**
```json
{
  "nodes": [
    {"id": "1", "data": {}}
  ],
  "edges": []
}
```

**Error (if node type is Number):**
```
validation error on node '1' field 'data.value': number node requires a 'value' field
```

## Performance Impact

Validation adds minimal overhead:
- **Time**: < 1ms for typical workflows (< 100 nodes)
- **Memory**: Negligible (creates temporary maps for validation)
- **CPU**: One-time cost before execution

The validation runs once before execution and does not impact workflow execution performance.

## Integration with Existing Code

Validation is backward compatible and optional:

### Without Validation (Works as Before)
```go
engine, err := workflow.NewEngine(payload)
if err != nil {
    log.Fatal(err)
}
result, err := engine.Execute()
// Errors caught during execution
```

### With Validation (Recommended)
```go
// Validate first
validationResult, err := workflow.ValidatePayload(payload)
if err != nil {
    log.Fatal(err)
}
if !validationResult.Valid {
    log.Fatalf("Invalid workflow: %v", validationResult.Errors)
}

// Then execute
engine, err := workflow.NewEngine(payload)
if err != nil {
    log.Fatal(err)
}
result, err := engine.Execute()
// Better error handling with early validation
```

## Testing

The validation layer includes 24 comprehensive tests covering:
- Empty workflows
- Duplicate node IDs
- Invalid edges
- Cyclic dependencies
- Node-specific validation
- Valid workflows

Run validation tests:
```bash
cd backend
go test -v -run TestValidation
```

## Future Enhancements

Potential future improvements:
1. **Warnings**: Non-fatal issues (e.g., unused nodes)
2. **Performance validation**: Check for potential performance issues
3. **Security validation**: Additional security checks beyond existing SSRF protection
4. **Custom validators**: Allow users to add custom validation rules
5. **Validation levels**: Strict vs lenient validation modes

## Related Documentation

- [ARCHITECTURE.md](../ARCHITECTURE.md) - Overall architecture
- [backend/README.md](README.md) - Backend workflow engine documentation
- [docs/NODES.md](../docs/NODES.md) - Complete node type reference
- [QUICK_WINS_SUMMARY.md](../QUICK_WINS_SUMMARY.md) - Previous quick wins

---

**Added**: October 30, 2025  
**Version**: 1.0  
**Status**: ✅ Complete
