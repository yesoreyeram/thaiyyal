# Principles: Workload Protection

This document outlines the workload protection mechanisms implemented in Thaiyyal to prevent resource exhaustion and ensure system stability.

## Overview

Workload protection ensures that:
- **Individual workflows cannot monopolize resources**
- **System remains responsive under load**
- **Resource exhaustion attacks are prevented**
- **Fair resource allocation across workflows**

## Protection Mechanisms

### 1. Execution Timeout

**Purpose:** Prevent infinite or long-running workflows.

**Implementation:**

```go
// Create context with timeout
ctx, cancel := context.WithTimeout(
    context.Background(), 
    config.MaxExecutionTime,  // Default: 30 seconds
)
defer cancel()

// Execute with timeout
done := make(chan error, 1)
go func() {
    done <- e.executeWorkflow(ctx)
}()

select {
case err := <-done:
    return result, err
case <-ctx.Done():
    return nil, fmt.Errorf("workflow timeout after %v", config.MaxExecutionTime)
}
```

**Configuration:**

```go
// Different timeout strategies
config := types.Config{
    MaxExecutionTime: 30 * time.Second,  // Default
}

// Development: Longer timeout
devConfig := types.DevelopmentConfig()  // 5 minutes

// Production: Strict timeout
prodConfig := types.ValidationLimits()  // 10 seconds
```

**Best Practices:**
- Set appropriate timeouts for use case
- Use shorter timeouts for user-facing workflows
- Use longer timeouts for batch processing
- Always clean up resources on timeout

### 2. Node Execution Limits

**Purpose:** Prevent workflows with excessive iterations or recursion.

**Implementation:**

```go
type Engine struct {
    nodeExecutionCount int
    countersMu         sync.RWMutex
}

func (e *Engine) IncrementNodeExecution() error {
    e.countersMu.Lock()
    defer e.countersMu.Unlock()
    
    e.nodeExecutionCount++
    if e.config.MaxNodeExecutions > 0 && 
       e.nodeExecutionCount > e.config.MaxNodeExecutions {
        return fmt.Errorf(
            "max node executions exceeded: %d (limit: %d)",
            e.nodeExecutionCount,
            e.config.MaxNodeExecutions,
        )
    }
    return nil
}
```

**When Counted:**
- Each node execution
- Each loop iteration
- Each recursive call
- Each retry attempt

**Configuration:**

```go
config := types.Config{
    MaxNodeExecutions: 10000,  // Default: 10k nodes
}

// For loop-heavy workflows
config.MaxNodeExecutions = 100000

// For simple workflows
config.MaxNodeExecutions = 1000
```

### 3. HTTP Call Limits

**Purpose:** Prevent excessive external API calls.

**Implementation:**

```go
func (e *Engine) IncrementHTTPCall() error {
    e.countersMu.Lock()
    defer e.countersMu.Unlock()
    
    e.httpCallCount++
    if e.config.MaxHTTPCallsPerExec > 0 && 
       e.httpCallCount > e.config.MaxHTTPCallsPerExec {
        return fmt.Errorf(
            "max HTTP calls exceeded: %d (limit: %d)",
            e.httpCallCount,
            e.config.MaxHTTPCallsPerExec,
        )
    }
    return nil
}
```

**Applied to:**
- HTTP node executions
- Retry attempts
- Redirects (counted separately)

**Configuration:**

```go
config := types.Config{
    MaxHTTPCallsPerExec: 100,  // Default: 100 calls
}
```

### 4. Loop Iteration Limits

**Purpose:** Prevent infinite loops in ForEach, WhileLoop nodes.

**Implementation:**

```go
func (e *WhileLoopExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    iterations := 0
    maxIterations := ctx.GetConfig().MaxLoopIterations
    
    for evaluateCondition(ctx, node) {
        iterations++
        if maxIterations > 0 && iterations > maxIterations {
            return nil, fmt.Errorf(
                "loop iteration limit exceeded: %d",
                maxIterations,
            )
        }
        
        // Execute loop body
        if err := executeLoopBody(ctx, node); err != nil {
            return nil, err
        }
    }
    
    return iterations, nil
}
```

**Configuration:**

```go
config := types.Config{
    MaxLoopIterations: 10000,  // Default: 10k iterations
}
```

### 5. Data Size Limits

**Purpose:** Prevent memory exhaustion from large data structures.

**String Limits:**

```go
func ValidateString(s string, config types.Config) error {
    if config.MaxStringLength > 0 && len(s) > config.MaxStringLength {
        return fmt.Errorf(
            "string too long: %d bytes (max %d)",
            len(s),
            config.MaxStringLength,
        )
    }
    return nil
}
```

**Array Limits:**

```go
func ValidateArray(arr []interface{}, config types.Config) error {
    if config.MaxArraySize > 0 && len(arr) > config.MaxArraySize {
        return fmt.Errorf(
            "array too large: %d elements (max %d)",
            len(arr),
            config.MaxArraySize,
        )
    }
    return nil
}
```

**Object Limits:**

```go
func ValidateObject(obj map[string]interface{}, config types.Config) error {
    // Check key count
    if config.MaxObjectKeys > 0 && len(obj) > config.MaxObjectKeys {
        return fmt.Errorf("too many object keys: %d (max %d)", 
            len(obj), config.MaxObjectKeys)
    }
    
    // Check depth
    if config.MaxObjectDepth > 0 {
        if err := validateDepth(obj, config.MaxObjectDepth, 0); err != nil {
            return err
        }
    }
    
    return nil
}

func validateDepth(obj map[string]interface{}, maxDepth, currentDepth int) error {
    if currentDepth > maxDepth {
        return fmt.Errorf("object depth exceeded: %d (max %d)", 
            currentDepth, maxDepth)
    }
    
    for _, value := range obj {
        if nested, ok := value.(map[string]interface{}); ok {
            if err := validateDepth(nested, maxDepth, currentDepth+1); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

**Configuration:**

```go
config := types.Config{
    MaxStringLength: 1 * 1024 * 1024,  // 1 MB
    MaxArraySize:    10000,             // 10k elements
    MaxObjectDepth:  10,                // 10 levels deep
    MaxObjectKeys:   1000,              // 1k keys
}
```

### 6. Variable Count Limits

**Purpose:** Prevent memory exhaustion from excessive variables.

**Implementation:**

```go
func (e *Engine) SetVariable(name string, value interface{}) error {
    // Validate value
    if err := types.ValidateValue(value, e.config); err != nil {
        return err
    }
    
    // Check variable count
    if e.config.MaxVariables > 0 {
        vars := e.state.GetAllVariables()
        if len(vars) >= e.config.MaxVariables {
            return fmt.Errorf(
                "max variables exceeded: %d (limit: %d)",
                len(vars),
                e.config.MaxVariables,
            )
        }
    }
    
    return e.state.SetVariable(name, value)
}
```

**Configuration:**

```go
config := types.Config{
    MaxVariables: 1000,  // Default: 1000 variables
}
```

### 7. Recursion Depth Limits

**Purpose:** Prevent stack overflow from deep recursion.

**Implementation:**

```go
type recursionTracker struct {
    depth    int
    maxDepth int
}

func (t *recursionTracker) Enter() error {
    t.depth++
    if t.maxDepth > 0 && t.depth > t.maxDepth {
        return fmt.Errorf(
            "recursion depth exceeded: %d (max %d)",
            t.depth,
            t.maxDepth,
        )
    }
    return nil
}

func (t *recursionTracker) Exit() {
    t.depth--
}
```

**Configuration:**

```go
config := types.Config{
    MaxRecursionDepth: 100,  // Default: 100 levels
}
```

## Resource Limit Profiles

### Default Profile

**Use Case:** General purpose workflows

```go
func DefaultConfig() Config {
    return Config{
        MaxExecutionTime:    30 * time.Second,
        MaxNodeExecutions:   10000,
        MaxHTTPCallsPerExec: 100,
        MaxStringLength:     1 * 1024 * 1024,  // 1 MB
        MaxArraySize:        10000,
        MaxObjectDepth:      10,
        MaxObjectKeys:       1000,
        MaxVariables:        1000,
        MaxLoopIterations:   10000,
        MaxRecursionDepth:   100,
    }
}
```

### Strict Validation Profile

**Use Case:** User-facing APIs, untrusted workflows

```go
func ValidationLimits() Config {
    return Config{
        MaxExecutionTime:    10 * time.Second,  // Stricter
        MaxNodeExecutions:   1000,              // Stricter
        MaxHTTPCallsPerExec: 10,                // Stricter
        MaxStringLength:     100 * 1024,        // 100 KB
        MaxArraySize:        1000,
        MaxObjectDepth:      5,
        MaxObjectKeys:       100,
        MaxVariables:        100,
        MaxLoopIterations:   1000,
        MaxRecursionDepth:   10,
    }
}
```

### Development Profile

**Use Case:** Development, testing, debugging

```go
func DevelopmentConfig() Config {
    return Config{
        MaxExecutionTime:    5 * time.Minute,   // Relaxed
        MaxNodeExecutions:   100000,            // Relaxed
        MaxHTTPCallsPerExec: 1000,              // Relaxed
        MaxStringLength:     10 * 1024 * 1024,  // 10 MB
        MaxArraySize:        100000,
        MaxObjectDepth:      20,
        MaxObjectKeys:       10000,
        MaxVariables:        10000,
        MaxLoopIterations:   100000,
        MaxRecursionDepth:   200,
    }
}
```

## Monitoring and Metrics

### Resource Usage Tracking

```go
type ExecutionMetrics struct {
    NodeExecutions    int
    HTTPCalls         int
    LoopIterations    int
    VariableCount     int
    MaxStringLength   int
    MaxArraySize      int
    MaxObjectDepth    int
    ExecutionDuration time.Duration
}

func (e *Engine) GetMetrics() ExecutionMetrics {
    return ExecutionMetrics{
        NodeExecutions:    e.GetNodeExecutionCount(),
        HTTPCalls:         e.GetHTTPCallCount(),
        VariableCount:     len(e.GetVariables()),
        ExecutionDuration: time.Since(e.startTime),
    }
}
```

### Metrics to Monitor

1. **Execution Metrics:**
   - Execution duration distribution
   - Node execution count distribution
   - Timeout rate
   - Success/failure rate

2. **Resource Metrics:**
   - Peak memory usage
   - HTTP call count
   - Variable count
   - Loop iteration count

3. **Limit Violations:**
   - Timeout violations
   - Node execution limit violations
   - HTTP call limit violations
   - Data size limit violations

### Alerting

```yaml
# Example Prometheus alerts
groups:
  - name: workload_protection
    rules:
      - alert: HighTimeoutRate
        expr: rate(workflow_timeouts[5m]) > 0.1
        severity: warning
        description: High rate of workflow timeouts
        
      - alert: ResourceLimitViolations
        expr: rate(limit_violations[5m]) > 5
        severity: high
        description: Frequent resource limit violations
        
      - alert: LongRunningWorkflows
        expr: workflow_duration_seconds > 30
        severity: info
        description: Workflow taking longer than expected
```

## Best Practices

### 1. Choose Appropriate Limits

```go
// API endpoints: Use strict limits
apiConfig := types.ValidationLimits()

// Batch processing: Use relaxed limits
batchConfig := types.DefaultConfig()
batchConfig.MaxExecutionTime = 10 * time.Minute

// Development: Use development config
devConfig := types.DevelopmentConfig()
```

### 2. Handle Limit Violations Gracefully

```go
result, err := engine.Execute()
if err != nil {
    // Check for specific limit violations
    if strings.Contains(err.Error(), "timeout") {
        // Handle timeout
        log.Warn("Workflow timeout, consider optimizing")
        return handleTimeout()
    }
    if strings.Contains(err.Error(), "limit exceeded") {
        // Handle limit violation
        log.Error("Resource limit exceeded")
        return handleLimitViolation()
    }
    return err
}
```

### 3. Monitor and Tune Limits

```go
// Collect metrics
metrics := engine.GetMetrics()

// Log for analysis
logger.WithFields(map[string]interface{}{
    "node_executions":    metrics.NodeExecutions,
    "http_calls":         metrics.HTTPCalls,
    "execution_duration": metrics.ExecutionDuration,
}).Info("Workflow completed")

// Tune limits based on metrics
if metrics.NodeExecutions > 0.8 * config.MaxNodeExecutions {
    log.Warn("Workflow approaching node execution limit")
}
```

### 4. Test with Limits

```go
func TestResourceLimits(t *testing.T) {
    // Test with strict limits
    config := types.ValidationLimits()
    config.MaxLoopIterations = 10
    
    engine, _ := engine.NewWithConfig(payload, config)
    _, err := engine.Execute()
    
    // Expect limit violation
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "loop iteration limit")
}
```

## Related Documentation

- [Zero-Trust Security](PRINCIPLES_ZERO_TRUST.md)
- [No Runtime Errors](PRINCIPLES_NO_RUNTIME_ERRORS.md)
- [Security Requirements](REQUIREMENTS_NON_FUNCTIONAL_SECURITY.md)
- [Performance Tuning](PERFORMANCE_TUNING.md)

---

**Last Updated:** 2025-11-03
**Version:** 1.0
