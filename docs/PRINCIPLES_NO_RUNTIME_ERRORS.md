# Principles: No Runtime Errors (Input Validation and Sanitization)

This document outlines the principles and practices for preventing runtime errors through comprehensive input validation and output sanitization.

## Overview

The "No Runtime Errors" principle ensures:
- **All inputs are validated before processing**
- **Type safety is enforced**
- **Errors are caught early**
- **Graceful error handling**
- **Predictable behavior**

## Input Validation Strategy

### 1. Validate Early, Validate Often

**Validation Layers:**

```
User Input
  ↓
Layer 1: JSON Schema Validation
  ↓
Layer 2: Type Validation  
  ↓
Layer 3: Business Logic Validation
  ↓
Layer 4: Security Validation
  ↓
Safe to Process
```

### 2. Type-Safe Validation

**Go Type System:**

```go
// Strong typing prevents many errors
type NodeData struct {
    Value          *float64                 `json:"value"`
    Text           *string                  `json:"text"`
    Op             *string                  `json:"op"`
    URL            *string                  `json:"url"`
    Condition      *string                  `json:"condition"`
    // All fields are pointers to distinguish between zero and absent
}

// Validation function
func (e *OperationExecutor) Validate(node types.Node) error {
    if node.Data.Op == nil {
        return fmt.Errorf("operation is required")
    }
    
    validOps := []string{"add", "subtract", "multiply", "divide"}
    if !contains(validOps, *node.Data.Op) {
        return fmt.Errorf("invalid operation: %s", *node.Data.Op)
    }
    
    return nil
}
```

### 3. Comprehensive Validation Functions

**String Validation:**

```go
func ValidateString(s string, config types.Config) error {
    // Check length
    if config.MaxStringLength > 0 && len(s) > config.MaxStringLength {
        return fmt.Errorf("string too long: %d bytes (max %d)", 
            len(s), config.MaxStringLength)
    }
    
    // Check for null bytes (security)
    if strings.Contains(s, "\x00") {
        return fmt.Errorf("string contains null bytes")
    }
    
    // Valid UTF-8
    if !utf8.ValidString(s) {
        return fmt.Errorf("string is not valid UTF-8")
    }
    
    return nil
}
```

**Number Validation:**

```go
func ValidateNumber(n float64) error {
    // Check for NaN
    if math.IsNaN(n) {
        return fmt.Errorf("number is NaN")
    }
    
    // Check for infinity
    if math.IsInf(n, 0) {
        return fmt.Errorf("number is infinite")
    }
    
    return nil
}
```

**Array Validation:**

```go
func ValidateArray(arr []interface{}, config types.Config) error {
    // Check size
    if config.MaxArraySize > 0 && len(arr) > config.MaxArraySize {
        return fmt.Errorf("array too large: %d (max %d)", 
            len(arr), config.MaxArraySize)
    }
    
    // Validate elements
    for i, elem := range arr {
        if err := ValidateValue(elem, config); err != nil {
            return fmt.Errorf("invalid array element at index %d: %w", i, err)
        }
    }
    
    return nil
}
```

**Object Validation:**

```go
func ValidateObject(obj map[string]interface{}, config types.Config) error {
    // Check key count
    if config.MaxObjectKeys > 0 && len(obj) > config.MaxObjectKeys {
        return fmt.Errorf("too many keys: %d (max %d)", 
            len(obj), config.MaxObjectKeys)
    }
    
    // Check depth
    if err := validateObjectDepth(obj, config.MaxObjectDepth, 0); err != nil {
        return err
    }
    
    // Validate values
    for key, value := range obj {
        if err := ValidateValue(value, config); err != nil {
            return fmt.Errorf("invalid value for key %s: %w", key, err)
        }
    }
    
    return nil
}

func validateObjectDepth(obj map[string]interface{}, maxDepth, current int) error {
    if current > maxDepth {
        return fmt.Errorf("object depth exceeded: %d (max %d)", current, maxDepth)
    }
    
    for _, value := range obj {
        if nested, ok := value.(map[string]interface{}); ok {
            if err := validateObjectDepth(nested, maxDepth, current+1); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

## Error Handling Patterns

### 1. Explicit Error Returns

```go
// Good: Explicit error handling
func (e *Executor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    result, err := processInput(node.Data)
    if err != nil {
        return nil, fmt.Errorf("failed to process input: %w", err)
    }
    return result, nil
}

// Bad: Panic on error
func (e *Executor) Execute(ctx ExecutionContext, node types.Node) interface{} {
    result := processInput(node.Data)  // Might panic
    return result
}
```

### 2. Wrapped Errors with Context

```go
// Wrap errors to preserve context
if err := validateNode(node); err != nil {
    return nil, fmt.Errorf("node %s validation failed: %w", node.ID, err)
}

// Error chain provides full context:
// "node abc123 validation failed: operation validation failed: invalid operator: unknown"
```

### 3. Error Types

```go
// Define error types for different categories
var (
    ErrValidation      = errors.New("validation error")
    ErrExecution       = errors.New("execution error")
    ErrResourceLimit   = errors.New("resource limit exceeded")
    ErrTimeout         = errors.New("timeout")
)

// Check error type
if errors.Is(err, ErrValidation) {
    // Handle validation error
}
```

## Sanitization

### 1. Output Sanitization

```go
// Sanitize URLs before logging
func sanitizeURL(rawURL string) string {
    parsed, err := url.Parse(rawURL)
    if err != nil {
        return "[invalid-url]"
    }
    
    // Remove credentials
    parsed.User = nil
    
    // Remove query parameters (might contain tokens)
    parsed.RawQuery = ""
    
    return parsed.String()
}

// Usage in logs
logger.WithField("url", sanitizeURL(requestURL)).Info("Making HTTP request")
```

### 2. Error Message Sanitization

```go
// Don't expose internal details in errors
func sanitizeError(err error) error {
    msg := err.Error()
    
    // Remove file paths
    msg = removeFilePaths(msg)
    
    // Remove IP addresses
    msg = removeIPAddresses(msg)
    
    // Remove internal identifiers
    msg = removeInternalIDs(msg)
    
    return errors.New(msg)
}
```

### 3. Log Sanitization

```go
// Sanitize before logging
type SanitizedLogger struct {
    logger *logging.Logger
}

func (l *SanitizedLogger) Info(msg string, fields map[string]interface{}) {
    sanitized := make(map[string]interface{})
    for k, v := range fields {
        sanitized[k] = sanitizeValue(v)
    }
    l.logger.WithFields(sanitized).Info(msg)
}

func sanitizeValue(v interface{}) interface{} {
    switch val := v.(type) {
    case string:
        return sanitizeString(val)
    case map[string]interface{}:
        return sanitizeMap(val)
    default:
        return v
    }
}
```

## Defensive Programming

### 1. Null Safety

```go
// Always check for nil pointers
if node.Data.Op == nil {
    return nil, errors.New("operation is required")
}

// Safe dereferencing with default
op := "add"  // default
if node.Data.Op != nil {
    op = *node.Data.Op
}
```

### 2. Bounds Checking

```go
// Always check array bounds
func getArrayElement(arr []interface{}, index int) (interface{}, error) {
    if index < 0 || index >= len(arr) {
        return nil, fmt.Errorf("index out of bounds: %d (length: %d)", 
            index, len(arr))
    }
    return arr[index], nil
}
```

### 3. Type Assertions with Checks

```go
// Safe type assertion
if strVal, ok := value.(string); ok {
    // Use strVal
} else {
    return nil, fmt.Errorf("expected string, got %T", value)
}

// Alternative: Type switch
switch v := value.(type) {
case string:
    // Handle string
case float64:
    // Handle number
default:
    return nil, fmt.Errorf("unexpected type: %T", v)
}
```

### 4. Division by Zero Protection

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    
    result := a / b
    
    // Check result
    if math.IsInf(result, 0) || math.IsNaN(result) {
        return 0, errors.New("invalid division result")
    }
    
    return result, nil
}
```

## Testing for Edge Cases

### 1. Null/Empty Input Tests

```go
func TestNullInputs(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        wantErr bool
    }{
        {"nil", nil, true},
        {"empty string", "", true},
        {"empty array", []interface{}{}, true},
        {"empty object", map[string]interface{}{}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

### 2. Boundary Value Tests

```go
func TestBoundaryValues(t *testing.T) {
    tests := []struct {
        name  string
        value float64
        valid bool
    }{
        {"zero", 0, true},
        {"max int", math.MaxInt64, true},
        {"min int", math.MinInt64, true},
        {"max float", math.MaxFloat64, true},
        {"infinity", math.Inf(1), false},
        {"negative infinity", math.Inf(-1), false},
        {"NaN", math.NaN(), false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateNumber(tt.value)
            if (err == nil) != tt.valid {
                t.Errorf("validation mismatch")
            }
        })
    }
}
```

### 3. Malicious Input Tests

```go
func TestMaliciousInputs(t *testing.T) {
    tests := []struct {
        name  string
        input string
    }{
        {"SQL injection", "'; DROP TABLE users; --"},
        {"XSS", "<script>alert('xss')</script>"},
        {"Path traversal", "../../etc/passwd"},
        {"Null bytes", "hello\x00world"},
        {"Very long string", strings.Repeat("A", 10000000)},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateString(tt.input, DefaultConfig())
            assert.Error(t, err, "should reject malicious input")
        })
    }
}
```

## Best Practices

### 1. Validate at System Boundaries

```go
// Validate when data enters the system
func NewEngine(payloadJSON []byte) (*Engine, error) {
    // Parse
    var payload types.Payload
    if err := json.Unmarshal(payloadJSON, &payload); err != nil {
        return nil, fmt.Errorf("invalid JSON: %w", err)
    }
    
    // Validate
    if err := validatePayload(payload); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Now safe to use
    return createEngine(payload), nil
}
```

### 2. Use Type Safety

```go
// Use strong types instead of interface{}
type NodeID string
type NodeType string

func getNode(id NodeID) (*Node, error) {
    // Compiler ensures type safety
}

// Instead of:
func getNode(id interface{}) (*Node, error) {
    // Runtime type checking required
}
```

### 3. Fail Fast

```go
// Detect errors early
func Execute() error {
    // Validate immediately
    if err := validate(); err != nil {
        return err  // Fail fast
    }
    
    // Continue with valid data
    return process()
}

// Instead of collecting errors and checking later
```

### 4. Document Assumptions

```go
// Document preconditions and postconditions
// Execute executes a node.
//
// Preconditions:
//   - node must be valid (call Validate first)
//   - node must have required inputs available
//
// Postconditions:
//   - returns non-nil result on success
//   - returns error with context on failure
//
func Execute(node Node) (interface{}, error) {
    // Implementation
}
```

## Related Documentation

- [Zero-Trust Security](PRINCIPLES_ZERO_TRUST.md)
- [Workload Protection](PRINCIPLES_WORKLOAD_PROTECTION.md)
- [Code Quality](REQUIREMENTS_NON_FUNCTIONAL_CODE_QUALITY.md)
- [Testing](REQUIREMENTS_NON_FUNCTIONAL_TESTING.md)

---

**Last Updated:** 2025-11-03
**Version:** 1.0
