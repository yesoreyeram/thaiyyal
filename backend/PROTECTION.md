# Workflow Protection Measures

This document describes the comprehensive protection measures implemented in the Thaiyyal workflow engine to prevent resource exhaustion and security vulnerabilities.

## Overview

The workflow engine implements multiple layers of protection to ensure safe execution of workflows, even when processing untrusted or malicious workflows. These protections are configurable and have safe defaults suitable for production use.

## Protection Categories

### 1. Execution Time Limits

#### MaxExecutionTime
- **Default**: 5 minutes
- **Validation**: 1 minute
- **Development**: 30 minutes
- **Description**: Maximum time for entire workflow execution. Prevents workflows from running indefinitely.
- **Enforcement**: Enforced via context timeout in `Engine.Execute()`

#### MaxNodeExecutionTime
- **Default**: 30 seconds
- **Validation**: 10 seconds
- **Development**: 5 minutes
- **Description**: Maximum time for single node execution (currently advisory, not enforced per-node)

### 2. Iteration Limits

#### MaxIterations
- **Default**: 1000
- **Validation**: 100
- **Development**: 10000
- **Description**: Default maximum iterations for loop nodes (for_each, while_loop)
- **Enforcement**: Checked in loop executors

#### MaxNodeExecutions
- **Default**: 10000
- **Validation**: 1000
- **Development**: 100000
- **Description**: Maximum total node executions including loop iterations. Prevents excessive execution through loops.
- **Enforcement**: Tracked and enforced in `Engine.executeNode()` via `IncrementNodeExecution()`
- **Error**: "maximum node executions exceeded: X (limit: Y)"

### 3. HTTP Protection

#### HTTPTimeout
- **Default**: 30 seconds
- **Description**: Timeout for individual HTTP requests
- **Enforcement**: Applied in HTTP client configuration

#### MaxHTTPRedirects
- **Default**: 10
- **Description**: Maximum number of HTTP redirects to follow
- **Enforcement**: Applied in HTTP client transport

#### MaxResponseSize
- **Default**: 10MB
- **Validation**: 10MB
- **Development**: 10MB
- **Description**: Maximum size of HTTP response body
- **Enforcement**: Enforced in `HTTPExecutor.Execute()` using `io.LimitReader`
- **Error**: "response too large (exceeds X bytes limit)"

#### MaxHTTPCallsPerExec
- **Default**: 100
- **Validation**: 10
- **Development**: 1000
- **Description**: Maximum HTTP calls allowed per workflow execution. Prevents excessive HTTP requests.
- **Enforcement**: Tracked and enforced in `HTTPExecutor.Execute()` via `IncrementHTTPCall()`
- **Error**: "maximum HTTP calls per execution exceeded: X (limit: Y)"

#### BlockInternalIPs
- **Default**: true
- **Description**: Block requests to internal/private IP addresses (SSRF protection)
- **Enforcement**: Enforced in `HTTPExecutor.Execute()` via `isAllowedURL()`

#### AllowedURLPatterns
- **Default**: [] (empty = allow all external URLs)
- **Description**: Whitelist of allowed URL patterns for additional control

### 4. Data Size Limits

#### MaxInputSize
- **Default**: 10MB
- **Description**: Maximum size of input data (currently advisory)

#### MaxPayloadSize
- **Default**: 1MB
- **Description**: Maximum size of workflow payload JSON

#### MaxStringLength
- **Default**: 1MB
- **Validation**: 100KB
- **Development**: 10MB
- **Description**: Maximum length of string values in workflow state
- **Enforcement**: Enforced in `ValidateValue()` for variables and node results
- **Error**: "string too long: X bytes (limit: Y)"

#### MaxArrayLength
- **Default**: 10000
- **Validation**: 1000
- **Development**: 100000
- **Description**: Maximum number of elements in array values
- **Enforcement**: Enforced recursively in `ValidateValue()` for variables and node results
- **Error**: "array too large: X elements (limit: Y)"

### 5. Workflow Structure Limits

#### MaxNodes
- **Default**: 1000
- **Validation**: 100
- **Development**: 10000
- **Description**: Maximum number of nodes in workflow definition

#### MaxEdges
- **Default**: 10000
- **Validation**: 1000
- **Development**: 100000
- **Description**: Maximum number of edges in workflow definition

### 6. State Management Limits

#### MaxVariables
- **Default**: 1000
- **Validation**: 100
- **Development**: 10000
- **Description**: Maximum number of variables in workflow state
- **Enforcement**: Enforced in `Engine.SetVariable()` before creating new variables
- **Error**: "maximum variables exceeded: X (limit: Y)"

#### MaxContextDepth
- **Default**: 32
- **Validation**: 16
- **Development**: 64
- **Description**: Maximum nesting depth of objects/arrays. Prevents stack overflow from deeply nested structures.
- **Enforcement**: Enforced in `ValidateValue()` via `getValueDepth()`
- **Error**: "value too deeply nested: X levels (limit: Y)"

#### MaxCacheSize
- **Default**: 1000
- **Description**: Maximum number of cache entries (LRU eviction)

### 7. Retry Configuration

#### DefaultMaxAttempts
- **Default**: 3
- **Description**: Default maximum retry attempts for retry nodes

#### DefaultBackoff
- **Default**: 1 second
- **Description**: Default initial backoff delay for retry nodes

## Configuration Presets

### Default Configuration
```go
config := types.DefaultConfig()
```
Balanced settings suitable for production use with reasonable resource limits.

### Validation Limits
```go
config := types.ValidationLimits()
```
Strict limits suitable for validating untrusted workflows before execution.

### Development Configuration
```go
config := types.DevelopmentConfig()
```
Relaxed limits for development and testing environments.

## Runtime Counters

The engine maintains runtime counters to track resource usage during execution:

### Node Execution Counter
- **Access**: `engine.GetNodeExecutionCount()`
- **Description**: Total number of nodes executed, including loop iterations
- **Thread-safe**: Yes (protected by mutex)

### HTTP Call Counter
- **Access**: `engine.GetHTTPCallCount()`
- **Description**: Total number of HTTP requests made during execution
- **Thread-safe**: Yes (protected by mutex)

## Validation Functions

### ValidateValue
```go
err := types.ValidateValue(value, config)
```
Validates a value against resource limits:
- Checks string length (MaxStringLength)
- Checks array length (MaxArrayLength) recursively
- Checks nesting depth (MaxContextDepth)

### getValueDepth
Internal function that calculates nesting depth of maps, slices, and arrays recursively.

## Best Practices

### 1. Choose Appropriate Configuration
- **Production**: Use `DefaultConfig()` as starting point
- **User-provided workflows**: Use `ValidationLimits()` first
- **Testing**: Use `DevelopmentConfig()` for convenience

### 2. Monitor Counters
After execution, check counters to understand resource usage:
```go
result, err := engine.Execute()
fmt.Printf("Nodes executed: %d\n", engine.GetNodeExecutionCount())
fmt.Printf("HTTP calls made: %d\n", engine.GetHTTPCallCount())
```

### 3. Customize Limits
Adjust limits based on your use case:
```go
config := types.DefaultConfig()
config.MaxHTTPCallsPerExec = 50  // Reduce for security
config.MaxNodeExecutions = 5000   // Increase for complex workflows
```

### 4. Zero Values Mean Unlimited
Setting a limit to 0 disables that check:
```go
config.MaxHTTPCallsPerExec = 0  // Unlimited HTTP calls (not recommended)
config.MaxNodeExecutions = 0     // Unlimited node executions (not recommended)
```

### 5. Validate Early
For untrusted workflows, use ValidationLimits first to fail fast:
```go
// Quick validation with strict limits
validationEngine, _ := NewWithConfig(payload, types.ValidationLimits())
_, err := validationEngine.Execute()
if err != nil {
    return fmt.Errorf("workflow validation failed: %w", err)
}

// Execute with production limits
productionEngine, _ := NewWithConfig(payload, types.DefaultConfig())
result, err := productionEngine.Execute()
```

## Error Messages

All protection violations return descriptive errors:

| Protection | Error Pattern |
|------------|---------------|
| Node executions | `maximum node executions exceeded: X (limit: Y)` |
| HTTP calls | `maximum HTTP calls per execution exceeded: X (limit: Y)` |
| String length | `string too long: X bytes (limit: Y)` |
| Array length | `array too large: X elements (limit: Y)` |
| Nesting depth | `value too deeply nested: X levels (limit: Y)` |
| Variables | `maximum variables exceeded: X (limit: Y)` |
| HTTP response | `response too large (exceeds X bytes limit)` |
| Workflow timeout | `workflow execution timeout: exceeded Xm` |

## Security Considerations

### SSRF Protection
The HTTP executor includes SSRF (Server-Side Request Forgery) protection:
- Blocks internal IP addresses (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16)
- Blocks localhost and loopback addresses
- Blocks link-local addresses
- Blocks cloud metadata endpoints (AWS, GCP, Azure)

### DoS Prevention
Multiple layers prevent denial-of-service attacks:
1. **Workflow timeout**: Prevents infinite execution
2. **Node execution limit**: Prevents runaway loops
3. **HTTP call limit**: Prevents HTTP flooding
4. **Response size limit**: Prevents memory exhaustion
5. **String/array limits**: Prevents memory bombs
6. **Nesting depth limit**: Prevents stack overflow

### Resource Exhaustion
Protection against resource exhaustion:
- Memory: String, array, and response size limits
- CPU: Execution time and iteration limits
- Network: HTTP call limits and timeouts
- State: Variable count and nesting depth limits

## Testing

Comprehensive test coverage ensures protections work correctly:

### Protection Tests
- `protection_test.go`: Tests for node execution and HTTP call limits
- `validation_test.go`: Tests for data validation limits

### Test Coverage
- Node execution limits (under/at/over)
- HTTP call limits (under/at/over)
- String length limits
- Array length limits
- Nesting depth limits
- Variable count limits
- Multiple limits working together
- Configuration presets (default/validation/development)

Run tests:
```bash
cd backend
go test ./pkg/engine -v -run Protection
go test ./pkg/engine -v -run Validation
```

## Future Enhancements

Potential future improvements:

1. **Memory Monitoring**: Track actual memory usage during execution
2. **Rate Limiting**: Per-time-window limits for HTTP calls
3. **Per-Node Timeouts**: Enforce MaxNodeExecutionTime per node
4. **Regex Complexity**: Prevent ReDoS attacks in text operations
5. **Logging**: Structured logging of protection violations
6. **Metrics**: Export protection metrics for monitoring
7. **Custom Validators**: Allow user-defined validation rules

## References

- Implementation: `backend/pkg/types/types.go` (Config definition)
- Defaults: `backend/pkg/types/helpers.go` (Config presets)
- Enforcement: `backend/pkg/engine/engine.go` (Engine implementation)
- Validation: `backend/pkg/types/helpers.go` (ValidateValue)
- HTTP Protection: `backend/pkg/executor/http.go` (HTTPExecutor)
- Tests: `backend/pkg/engine/protection_test.go`, `backend/pkg/engine/validation_test.go`
