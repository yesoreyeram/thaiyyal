# Principles: Zero-Trust Security

This document outlines the zero-trust security principles implemented in Thaiyyal.

## Overview

Thaiyyal implements a **zero-trust security model** where:
- **Nothing is trusted by default**
- **Everything is verified**
- **All inputs are validated**
- **All outputs are sanitized**
- **Least privilege is enforced**

## Core Principles

### 1. Never Trust, Always Verify

**Principle:** Assume all inputs are potentially malicious until proven safe.

**Implementation:**

```go
// Every input is validated
func (e *Engine) Execute() (*types.Result, error) {
    // 1. Validate JSON structure
    var payload types.Payload
    if err := json.Unmarshal(payloadJSON, &payload); err != nil {
        return nil, fmt.Errorf("invalid JSON: %w", err)
    }
    
    // 2. Validate workflow structure
    if err := validateWorkflow(payload); err != nil {
        return nil, err
    }
    
    // 3. Validate node configurations
    for _, node := range payload.Nodes {
        if err := e.registry.Validate(node); err != nil {
            return nil, err
        }
    }
}
```

**Validation Layers:**
1. **JSON Schema Validation**: Structure correctness
2. **Type Validation**: Data type correctness
3. **Business Logic Validation**: Semantic correctness
4. **Security Validation**: Safety checks

### 2. Defense in Depth

**Principle:** Multiple layers of security controls.

**Security Layers:**

```
┌─────────────────────────────────────────┐
│  Layer 1: Input Validation              │
│  • JSON parsing with size limits        │
│  • Type checking                         │
│  • Format validation                     │
└─────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│  Layer 2: Resource Limits               │
│  • Execution timeout                     │
│  • Node execution limits                 │
│  • Memory limits                         │
│  • Loop iteration limits                 │
└─────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│  Layer 3: Network Security              │
│  • SSRF protection                       │
│  • URL validation                        │
│  • IP blocking (private, localhost)     │
│  • Domain filtering                      │
└─────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│  Layer 4: Execution Isolation           │
│  • Context isolation                     │
│  • State sandboxing                      │
│  • Error containment                     │
└─────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│  Layer 5: Monitoring & Logging          │
│  • All actions logged                    │
│  • Security events tracked               │
│  • Anomaly detection                     │
└─────────────────────────────────────────┘
```

### 3. Least Privilege

**Principle:** Grant minimum necessary permissions.

**Implementation:**

```go
// Executors only have access to what they need
type ExecutionContext interface {
    // Node-specific access
    GetNodeInputs(nodeID string) []interface{}
    GetNode(nodeID string) *types.Node
    
    // Read-only access to results
    GetNodeResult(nodeID string) (interface{}, bool)
    GetAllNodeResults() map[string]interface{}
    
    // Controlled write access
    SetNodeResult(nodeID string, result interface{})
    
    // Limited state access
    GetVariable(name string) (interface{}, error)
    SetVariable(name string, value interface{}) error
    
    // No direct access to engine internals
    // No file system access
    // No network access (except through HTTP node)
}
```

**Access Control:**
- Executors can't access arbitrary files
- Network access requires explicit HTTP node
- State modifications are tracked and limited
- No direct memory access

### 4. Fail Securely

**Principle:** Failures should not compromise security.

**Implementation:**

```go
// Secure error handling
func (e *HTTPExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    url := *node.Data.URL
    
    // Validate URL before making request
    if err := e.ssrf.ValidateURL(url); err != nil {
        // Fail securely: Don't reveal internal details
        return nil, fmt.Errorf("URL validation failed")
        // NOT: return nil, fmt.Errorf("private IP detected: %s", url)
    }
    
    // Make request with timeout
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    resp, err := e.client.Get(url)
    if err != nil {
        // Fail securely: Limit error information
        return nil, fmt.Errorf("HTTP request failed")
        // NOT: return nil, err // Might reveal network topology
    }
    
    return resp, nil
}
```

**Secure Failure Modes:**
- Generic error messages externally
- Detailed errors logged internally
- No sensitive information in errors
- Graceful degradation

### 5. Complete Mediation

**Principle:** Every access is checked every time.

**Implementation:**

```go
// Every HTTP call is checked
func (h *HTTPExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    // Check limits before execution
    if err := ctx.IncrementHTTPCall(); err != nil {
        return nil, err
    }
    
    // Validate URL every time (no caching of validation)
    if err := h.ssrf.ValidateURL(*node.Data.URL); err != nil {
        return nil, err
    }
    
    // Make request
    return h.makeRequest(ctx, node)
}
```

**No Bypass:**
- Validation not cached
- Checks not skipped for "trusted" inputs
- Every execution validated independently

### 6. Separation of Duties

**Principle:** No single component has complete control.

**Separation:**

```
┌─────────────────────┐
│   Engine            │
│   • Orchestration   │
│   • No execution    │
└─────────────────────┘
         ↓
┌─────────────────────┐
│   Registry          │
│   • Executor lookup │
│   • No execution    │
└─────────────────────┘
         ↓
┌─────────────────────┐
│   Executor          │
│   • Node execution  │
│   • No orchestration│
└─────────────────────┘
         ↓
┌─────────────────────┐
│   Middleware        │
│   • Validation      │
│   • No business logic│
└─────────────────────┘
```

## Security Controls

### Input Validation

**All inputs validated:**

```go
// String validation
func ValidateString(s string, maxLength int) error {
    if len(s) > maxLength {
        return fmt.Errorf("string too long: %d (max %d)", len(s), maxLength)
    }
    return nil
}

// Array validation
func ValidateArray(arr []interface{}, maxSize int) error {
    if len(arr) > maxSize {
        return fmt.Errorf("array too large: %d (max %d)", len(arr), maxSize)
    }
    return nil
}

// Object validation
func ValidateObject(obj map[string]interface{}, maxDepth, maxKeys int) error {
    if len(obj) > maxKeys {
        return fmt.Errorf("too many keys: %d (max %d)", len(obj), maxKeys)
    }
    return validateDepth(obj, maxDepth, 0)
}
```

### SSRF Protection

**Comprehensive SSRF prevention:**

```go
type SSRFProtection struct {
    // Blocked by default
    blockPrivateIPs    bool  // 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
    blockLocalhost     bool  // 127.0.0.1, ::1
    blockLinkLocal     bool  // 169.254.0.0/16
    blockCloudMetadata bool  // 169.254.169.254, fd00:ec2::254
    
    // Optional filters
    allowedDomains     map[string]bool
    blockedDomains     map[string]bool
    allowedSchemes     map[string]bool
}

func (p *SSRFProtection) ValidateURL(urlStr string) error {
    // 1. Parse URL
    parsedURL, err := url.Parse(urlStr)
    
    // 2. Check scheme
    if !p.allowedSchemes[parsedURL.Scheme] {
        return fmt.Errorf("scheme not allowed: %s", parsedURL.Scheme)
    }
    
    // 3. Check domain blocklist
    // 4. Check domain allowlist
    // 5. Resolve to IP
    // 6. Check IP blocklist
    
    return nil
}
```

**Blocked by Default:**
- Private IP ranges (RFC 1918)
- Localhost and loopback
- Link-local addresses
- Cloud metadata endpoints
- file:// protocol
- Unknown schemes

### Resource Limits

**Protection against resource exhaustion:**

```go
type Config struct {
    // Execution limits
    MaxExecutionTime    time.Duration  // Default: 30s
    MaxNodeExecutions   int            // Default: 10000
    MaxHTTPCallsPerExec int            // Default: 100
    
    // Data limits
    MaxStringLength     int            // Default: 1MB
    MaxArraySize        int            // Default: 10000
    MaxObjectDepth      int            // Default: 10
    MaxObjectKeys       int            // Default: 1000
    
    // State limits
    MaxVariables        int            // Default: 1000
    MaxLoopIterations   int            // Default: 10000
    MaxRecursionDepth   int            // Default: 100
}
```

**Enforcement:**

```go
// Check before execution
func (e *Engine) IncrementNodeExecution() error {
    e.countersMu.Lock()
    defer e.countersMu.Unlock()
    
    e.nodeExecutionCount++
    if e.config.MaxNodeExecutions > 0 && 
       e.nodeExecutionCount > e.config.MaxNodeExecutions {
        return fmt.Errorf("node execution limit exceeded")
    }
    return nil
}
```

### Data Sanitization

**Output sanitization:**

```go
// Sanitize sensitive data from logs
func sanitizeURL(url string) string {
    parsed, err := url.Parse(url)
    if err != nil {
        return "[invalid URL]"
    }
    
    // Remove query parameters (might contain tokens)
    parsed.RawQuery = ""
    
    // Remove user info (credentials)
    parsed.User = nil
    
    return parsed.String()
}

// Sanitize errors
func sanitizeError(err error) error {
    // Remove paths, IPs, internal details
    msg := err.Error()
    msg = removePaths(msg)
    msg = removeIPs(msg)
    return fmt.Errorf("%s", msg)
}
```

## Security Testing

### Security Test Cases

```go
func TestSSRFProtection(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        wantErr bool
    }{
        {"private IP", "http://10.0.0.1", true},
        {"localhost", "http://localhost", true},
        {"loopback", "http://127.0.0.1", true},
        {"link local", "http://169.254.169.254", true},
        {"cloud metadata", "http://metadata.google.internal", true},
        {"valid URL", "https://api.example.com", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ssrf.ValidateURL(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

### Penetration Testing Scenarios

1. **SSRF Attacks**
   - Private IP access attempts
   - Cloud metadata access
   - DNS rebinding
   - URL parser bypasses

2. **Resource Exhaustion**
   - Infinite loops
   - Memory bombs (large strings/arrays)
   - Fork bombs (excessive parallel nodes)
   - Recursive workflows

3. **Injection Attacks**
   - Expression injection
   - Template injection
   - Path traversal

4. **Information Disclosure**
   - Error message leakage
   - Timing attacks
   - Side-channel leaks

## Security Monitoring

### What to Monitor

```go
// Log security-relevant events
logger.WithFields(map[string]interface{}{
    "event_type":    "ssrf_blocked",
    "url":           sanitizeURL(url),
    "ip":            remoteIP,
    "execution_id":  executionID,
}).Warn("SSRF attempt blocked")

logger.WithFields(map[string]interface{}{
    "event_type":    "limit_exceeded",
    "limit_type":    "node_executions",
    "count":         count,
    "execution_id":  executionID,
}).Error("Execution limit exceeded")
```

### Security Metrics

- SSRF blocks per hour
- Rate limit violations
- Resource limit violations
- Failed validations
- Execution timeouts
- Error rates by type

### Alerting Rules

```yaml
# Example alert rules
- alert: HighSSRFBlockRate
  expr: rate(ssrf_blocks[5m]) > 10
  severity: high
  description: High rate of SSRF blocks detected

- alert: ResourceLimitViolations
  expr: rate(limit_violations[5m]) > 5
  severity: medium
  description: Frequent resource limit violations
```

## Security Checklist

### Development Checklist

- [ ] All inputs validated
- [ ] Resource limits enforced
- [ ] SSRF protection enabled
- [ ] Errors don't leak sensitive info
- [ ] Security tests written
- [ ] Code review completed
- [ ] Dependency audit performed

### Deployment Checklist

- [ ] Security configuration reviewed
- [ ] Monitoring enabled
- [ ] Alerts configured
- [ ] Logs collection enabled
- [ ] Incident response plan documented
- [ ] Security scan performed

## Related Documentation

- [Workload Protection](PRINCIPLES_WORKLOAD_PROTECTION.md)
- [No Runtime Errors](PRINCIPLES_NO_RUNTIME_ERRORS.md)
- [Security Requirements](REQUIREMENTS_NON_FUNCTIONAL_SECURITY.md)
- [Security Best Practices](SECURITY_BEST_PRACTICES.md)

---

**Last Updated:** 2025-11-03
**Version:** 1.0
