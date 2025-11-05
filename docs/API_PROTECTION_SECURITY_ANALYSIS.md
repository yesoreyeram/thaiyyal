# Third-Party API Protection - Security Analysis and Implementation Plan

## Executive Summary

This document provides a comprehensive security analysis of how the Thaiyyal workflow engine protects third-party APIs from abuse, DDOS attacks, and misconfigurations. It evaluates current protections, identifies gaps, and proposes enhancements following zero-trust security principles.

**Current Security Posture: GOOD ✅ with identified areas for enhancement**

## Table of Contents

1. [Threat Model](#threat-model)
2. [Current Protection Mechanisms](#current-protection-mechanisms)
3. [Security Gap Analysis](#security-gap-analysis)
4. [Proposed Enhancements](#proposed-enhancements)
5. [Implementation Plan](#implementation-plan)
6. [Operational Procedures](#operational-procedures)
7. [Compliance and Governance](#compliance-and-governance)

---

## 1. Threat Model

### 1.1 Attack Vectors

#### A. Intentional Attacks

**Workflow-Based DDOS:**
- **Threat:** Malicious user creates workflow with loop calling API repeatedly
- **Impact:** API rate limits exceeded, service degradation, financial costs
- **Likelihood:** HIGH
- **Current Mitigation:** MaxHTTPCallsPerExec, MaxLoopIterations

**Resource Exhaustion:**
- **Threat:** Workflows designed to consume maximum resources
- **Impact:** Platform instability, denial of service to other users
- **Likelihood:** MEDIUM
- **Current Mitigation:** Multiple resource limits (CPU, memory, time)

**Credential Harvesting:**
- **Threat:** Workflow attempts to access internal endpoints or metadata
- **Impact:** Credential theft, security breach
- **Likelihood:** MEDIUM
- **Current Mitigation:** SSRF protection (private IP blocking, metadata blocking)

**Data Exfiltration:**
- **Threat:** Using HTTP nodes to send sensitive data to attacker-controlled endpoints
- **Impact:** Data breach, compliance violations
- **Likelihood:** MEDIUM
- **Current Mitigation:** Domain whitelisting, logging

#### B. Accidental Misconfigurations

**Infinite Loop with API Calls:**
- **Threat:** Developer error creates infinite loop
- **Impact:** API quota exhaustion, unexpected costs
- **Likelihood:** HIGH
- **Current Mitigation:** MaxLoopIterations, MaxHTTPCallsPerExec

**Retry Storm:**
- **Threat:** Aggressive retry configuration during API outage
- **Impact:** Amplified load on failing API, cascade failures
- **Likelihood:** HIGH
- **Current Mitigation:** Retry middleware with backoff

**Large Data Transfer:**
- **Threat:** Workflow processes/sends very large payloads
- **Impact:** Memory exhaustion, bandwidth costs
- **Likelihood:** MEDIUM
- **Current Mitigation:** MaxResponseSize, MaxArrayLength

**Fork Bomb Pattern:**
- **Threat:** Parallel node creating exponential executions
- **Impact:** Resource exhaustion, platform crash
- **Likelihood:** LOW
- **Current Mitigation:** MaxNodeExecutions

### 1.2 Impact Classification

| Impact Level | Description | Examples |
|--------------|-------------|----------|
| **Critical** | Platform-wide outage, data breach | Complete DDOS, credential theft |
| **High** | Single API service degradation | API rate limit exceeded, quota exhausted |
| **Medium** | Performance degradation | Slow responses, increased latency |
| **Low** | Minor resource waste | Slightly elevated API call counts |

---

## 2. Current Protection Mechanisms

### 2.1 Execution-Level Protections ✅

#### A. Global Execution Timeout
```go
MaxExecutionTime: 5 * time.Minute  // Default
```

**Purpose:** Prevent runaway workflows
**Effectiveness:** HIGH
**Coverage:** All workflows
**Gaps:** None - works as intended

#### B. Node Execution Counter
```go
MaxNodeExecutions: 10000  // Default: 10k nodes per execution
```

**Purpose:** Prevent infinite loops and excessive iterations
**Effectiveness:** HIGH
**Coverage:** All node types, including loops
**Gaps:** 
- ⚠️ Counter is per-execution, not per-API-endpoint
- ⚠️ Large limit (10k) may still cause issues with fast APIs

#### C. HTTP Call Counter
```go
MaxHTTPCallsPerExec: 0  // Default: unlimited
```

**Purpose:** Limit total HTTP calls per workflow execution
**Effectiveness:** MEDIUM (when configured)
**Coverage:** All HTTP nodes
**Gaps:**
- ⚠️ Default is unlimited (must be explicitly configured)
- ⚠️ No per-endpoint limits (can hit single API many times)
- ⚠️ No time-window based limiting (all calls allowed if under limit)
- ⚠️ No distinction between different APIs

**Recommendation:** Change default to reasonable limit (e.g., 100)

### 2.2 Network-Level Protections ✅

#### A. SSRF Protection
```go
// Zero-trust defaults
AllowHTTP:          false  // Require HTTPS
AllowPrivateIPs:    false  // Block 10.x, 172.16.x, 192.168.x
AllowLocalhost:     false  // Block 127.0.0.1, ::1
AllowLinkLocal:     false  // Block 169.254.x.x
AllowCloudMetadata: false  // Block 169.254.169.254
```

**Purpose:** Prevent SSRF attacks
**Effectiveness:** EXCELLENT
**Coverage:** All HTTP nodes
**Gaps:** None - comprehensive protection

#### B. Domain Whitelisting
```go
AllowedDomains: []string{"api.example.com", "api2.example.com"}
```

**Purpose:** Restrict HTTP calls to approved domains
**Effectiveness:** HIGH (when configured)
**Coverage:** All HTTP nodes
**Gaps:**
- ⚠️ Optional (empty = allow all external domains)
- ⚠️ No default whitelist for production

**Recommendation:** Encourage domain whitelisting in production

#### C. Response Size Limits
```go
MaxResponseSize: 10 * 1024 * 1024  // 10MB default
```

**Purpose:** Prevent memory exhaustion from large responses
**Effectiveness:** HIGH
**Coverage:** All HTTP nodes
**Gaps:** None - works well

### 2.3 Rate Limiting ⚠️

#### A. Middleware Rate Limiter
```go
// Exists but not enabled by default
type RateLimitMiddleware struct {
    globalLimiter    RateLimiter
    nodeTypeLimiters map[types.NodeType]RateLimiter
}
```

**Purpose:** Token bucket rate limiting
**Effectiveness:** HIGH (when enabled)
**Coverage:** Can limit by node type
**Gaps:**
- ⚠️ Not enabled by default
- ⚠️ No per-API-endpoint limits
- ⚠️ No distributed rate limiting (single instance only)
- ⚠️ No time-window based quotas (hourly, daily)

**Status:** EXISTS BUT NOT ACTIVELY USED

### 2.4 Retry Protection ✅

#### A. Retry Middleware
```go
DefaultMaxAttempts: 3
DefaultBackoff:     1 * time.Second
```

**Purpose:** Exponential backoff for failed requests
**Effectiveness:** MEDIUM
**Coverage:** Retry nodes
**Gaps:**
- ⚠️ No jitter in backoff (can cause thundering herd)
- ⚠️ Fixed backoff, not adaptive
- ⚠️ No circuit breaker pattern

### 2.5 Connection Pooling ✅

#### A. HTTP Client Pool
```go
MaxIdleConns:        100
MaxIdleConnsPerHost: 10
MaxConnsPerHost:     100
```

**Purpose:** Reuse connections, prevent connection exhaustion
**Effectiveness:** EXCELLENT
**Coverage:** All HTTP nodes
**Gaps:** None - well implemented

### 2.6 Observability 🔄

#### A. Metrics (via OpenTelemetry)
```go
metricHTTPCalls     = "http.calls.total"
metricHTTPDuration  = "http.call.duration"
```

**Purpose:** Monitor API usage patterns
**Effectiveness:** GOOD
**Coverage:** HTTP calls tracked
**Gaps:**
- ⚠️ No per-endpoint metrics (all HTTP calls aggregated)
- ⚠️ No error rate tracking by endpoint
- ⚠️ No rate limit violation metrics
- ⚠️ No anomaly detection

---

## 3. Security Gap Analysis

### 3.1 Critical Gaps (Must Fix)

| Gap ID | Description | Risk | Priority |
|--------|-------------|------|----------|
| **GAP-1** | No per-API-endpoint rate limiting | HIGH | P0 |
| **GAP-2** | No circuit breaker for failing APIs | HIGH | P0 |
| **GAP-3** | No audit logging for API calls | MEDIUM | P0 |
| **GAP-4** | MaxHTTPCallsPerExec defaults to unlimited | HIGH | P0 |

### 3.2 High Priority Gaps (Should Fix)

| Gap ID | Description | Risk | Priority |
|--------|-------------|------|----------|
| **GAP-5** | No time-window based quotas (hourly/daily) | MEDIUM | P1 |
| **GAP-6** | No retry jitter (thundering herd risk) | MEDIUM | P1 |
| **GAP-7** | No anomaly detection for API usage | MEDIUM | P1 |
| **GAP-8** | No per-workflow API budgets | MEDIUM | P1 |

### 3.3 Medium Priority Gaps (Nice to Have)

| Gap ID | Description | Risk | Priority |
|--------|-------------|------|----------|
| **GAP-9** | No distributed rate limiting | LOW | P2 |
| **GAP-10** | No adaptive retry backoff | LOW | P2 |
| **GAP-11** | No API cost tracking | LOW | P2 |
| **GAP-12** | No request prioritization | LOW | P2 |

---

## 4. Proposed Enhancements

### 4.1 Per-Endpoint Rate Limiting (GAP-1) 🔧

**Design:**

```go
// New configuration structure
type APIEndpointLimits struct {
    // Per-endpoint limits
    EndpointLimits map[string]EndpointLimit
    
    // Default limits for unlisted endpoints
    DefaultLimit EndpointLimit
}

type EndpointLimit struct {
    // Requests per time window
    RequestsPerSecond  int
    RequestsPerMinute  int
    RequestsPerHour    int
    RequestsPerDay     int
    
    // Burst allowance
    BurstSize int
    
    // Circuit breaker thresholds
    ErrorThreshold     float64  // e.g., 0.5 = 50% error rate
    FailureWindow      time.Duration
    RecoveryTimeout    time.Duration
}

// Example configuration
config := APIEndpointLimits{
    EndpointLimits: map[string]EndpointLimit{
        "api.github.com": {
            RequestsPerSecond: 10,
            RequestsPerHour:   5000,
            BurstSize:         20,
            ErrorThreshold:    0.3,
            FailureWindow:     1 * time.Minute,
            RecoveryTimeout:   5 * time.Minute,
        },
        "api.stripe.com": {
            RequestsPerSecond: 5,
            RequestsPerMinute: 100,
            BurstSize:         10,
        },
    },
    DefaultLimit: {
        RequestsPerSecond: 5,
        RequestsPerMinute: 100,
        RequestsPerHour:   1000,
        BurstSize:         10,
    },
}
```

**Implementation:**

```go
// New middleware: APIRateLimitMiddleware
type APIRateLimitMiddleware struct {
    limiters map[string]*EndpointRateLimiter
    config   APIEndpointLimits
    mu       sync.RWMutex
}

type EndpointRateLimiter struct {
    perSecond  *TokenBucket
    perMinute  *TokenBucket
    perHour    *SlidingWindow
    perDay     *SlidingWindow
}

func (m *APIRateLimitMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
    // Extract endpoint from URL
    endpoint := extractEndpoint(node.Data.URL)
    
    // Get or create limiter for this endpoint
    limiter := m.getOrCreateLimiter(endpoint)
    
    // Check all time windows
    if !limiter.Allow() {
        return nil, fmt.Errorf(
            "rate limit exceeded for %s: please retry later",
            endpoint,
        )
    }
    
    // Execute request
    return next(ctx, node)
}
```

**Benefits:**
- ✅ Prevents overwhelming single API
- ✅ Different limits for different APIs
- ✅ Time-window based quotas
- ✅ Burst allowance for spiky traffic

### 4.2 Circuit Breaker Pattern (GAP-2) 🔧

**Design:**

```go
// Circuit breaker states
type CircuitState int

const (
    StateClosed   CircuitState = iota  // Normal operation
    StateOpen                           // Circuit broken, reject requests
    StateHalfOpen                       // Testing if service recovered
)

type CircuitBreaker struct {
    state           CircuitState
    failureCount    int
    successCount    int
    lastFailureTime time.Time
    
    // Configuration
    failureThreshold int           // Failures before opening
    successThreshold int           // Successes to close from half-open
    openTimeout      time.Duration // Time to wait before trying half-open
    
    mu sync.RWMutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.RLock()
    state := cb.state
    cb.mu.RUnlock()
    
    switch state {
    case StateOpen:
        // Check if we should try half-open
        if time.Since(cb.lastFailureTime) > cb.openTimeout {
            cb.setState(StateHalfOpen)
            return cb.tryHalfOpen(fn)
        }
        return ErrCircuitOpen
        
    case StateHalfOpen:
        return cb.tryHalfOpen(fn)
        
    case StateClosed:
        return cb.tryCall(fn)
    }
    
    return nil
}

func (cb *CircuitBreaker) tryCall(fn func() error) error {
    err := fn()
    
    if err != nil {
        cb.recordFailure()
        return err
    }
    
    cb.recordSuccess()
    return nil
}
```

**Integration with HTTP Executor:**

```go
type HTTPExecutorWithCircuitBreaker struct {
    *HTTPExecutor
    circuitBreakers map[string]*CircuitBreaker
    mu              sync.RWMutex
}

func (e *HTTPExecutorWithCircuitBreaker) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    endpoint := extractEndpoint(*node.Data.URL)
    cb := e.getOrCreateCircuitBreaker(endpoint)
    
    var result interface{}
    err := cb.Call(func() error {
        var execErr error
        result, execErr = e.HTTPExecutor.Execute(ctx, node)
        return execErr
    })
    
    if err == ErrCircuitOpen {
        return nil, fmt.Errorf(
            "circuit breaker open for %s: service experiencing issues, please retry later",
            endpoint,
        )
    }
    
    return result, err
}
```

**Benefits:**
- ✅ Prevents cascade failures
- ✅ Fast-fail when API is down
- ✅ Automatic recovery testing
- ✅ Protects both client and server

### 4.3 Comprehensive Audit Logging (GAP-3) 📝

**Design:**

```go
// Audit log entry
type APICallAuditLog struct {
    // Identifiers
    ExecutionID string
    WorkflowID  string
    NodeID      string
    
    // Request details
    Timestamp   time.Time
    Method      string
    URL         string
    Endpoint    string
    Headers     map[string]string  // Sanitized
    
    // Response details
    StatusCode  int
    Duration    time.Duration
    ResponseSize int64
    
    // Security context
    UserID      string
    IPAddress   string
    
    // Compliance
    DataClass   string  // PII, financial, etc.
    
    // Result
    Success     bool
    ErrorType   string
    ErrorMsg    string  // Sanitized
}

// Audit logger interface
type APIAuditLogger interface {
    LogAPICall(ctx context.Context, log APICallAuditLog)
    QueryLogs(filter AuditFilter) ([]APICallAuditLog, error)
    DetectAnomalies() ([]Anomaly, error)
}

// Implementation
type StructuredAuditLogger struct {
    logger *logging.Logger
    store  AuditStore  // Could be database, S3, etc.
}

func (l *StructuredAuditLogger) LogAPICall(ctx context.Context, log APICallAuditLog) {
    // Sanitize sensitive data
    log.Headers = sanitizeHeaders(log.Headers)
    log.ErrorMsg = sanitizeError(log.ErrorMsg)
    
    // Log to structured logger
    l.logger.WithContext(ctx).WithFields(map[string]interface{}{
        "execution_id": log.ExecutionID,
        "workflow_id":  log.WorkflowID,
        "endpoint":     log.Endpoint,
        "status_code":  log.StatusCode,
        "duration_ms":  log.Duration.Milliseconds(),
        "success":      log.Success,
    }).Info("API call executed")
    
    // Store for audit trail
    l.store.Save(ctx, log)
}
```

**Anomaly Detection:**

```go
type AnomalyDetector struct {
    // Baseline metrics per endpoint
    baselines map[string]*EndpointBaseline
    
    // Thresholds
    volumeThreshold   float64  // e.g., 3x normal
    errorThreshold    float64  // e.g., 2x normal error rate
    latencyThreshold  float64  // e.g., 2x normal latency
}

type EndpointBaseline struct {
    avgRequestsPerMinute float64
    avgErrorRate         float64
    avgLatency           time.Duration
    
    // Rolling statistics
    recentRequests  *RollingWindow
    recentErrors    *RollingWindow
    recentLatencies *RollingWindow
}

func (d *AnomalyDetector) DetectAnomalies() []Anomaly {
    var anomalies []Anomaly
    
    for endpoint, baseline := range d.baselines {
        // Check volume spike
        currentRate := baseline.recentRequests.Rate()
        if currentRate > baseline.avgRequestsPerMinute * d.volumeThreshold {
            anomalies = append(anomalies, Anomaly{
                Type:     AnomalyTypeVolume,
                Endpoint: endpoint,
                Severity: SeverityHigh,
                Message:  fmt.Sprintf("Request volume spike: %.0f req/min (baseline: %.0f)", currentRate, baseline.avgRequestsPerMinute),
            })
        }
        
        // Check error rate spike
        currentErrorRate := baseline.recentErrors.Rate()
        if currentErrorRate > baseline.avgErrorRate * d.errorThreshold {
            anomalies = append(anomalies, Anomaly{
                Type:     AnomalyTypeErrors,
                Endpoint: endpoint,
                Severity: SeverityHigh,
                Message:  fmt.Sprintf("Error rate spike: %.2f%% (baseline: %.2f%%)", currentErrorRate*100, baseline.avgErrorRate*100),
            })
        }
        
        // Check latency spike
        currentLatency := baseline.recentLatencies.Avg()
        if currentLatency > float64(baseline.avgLatency) * d.latencyThreshold {
            anomalies = append(anomalies, Anomaly{
                Type:     AnomalyTypeLatency,
                Endpoint: endpoint,
                Severity: SeverityMedium,
                Message:  fmt.Sprintf("Latency spike: %.0fms (baseline: %.0fms)", currentLatency, float64(baseline.avgLatency.Milliseconds())),
            })
        }
    }
    
    return anomalies
}
```

**Benefits:**
- ✅ Complete audit trail for compliance
- ✅ Anomaly detection for security
- ✅ Troubleshooting and debugging
- ✅ Cost tracking and attribution

### 4.4 API Call Budgeting (GAP-8) 💰

**Design:**

```go
// Budget configuration
type APIBudgetConfig struct {
    // Per-workflow budgets
    WorkflowBudgets map[string]WorkflowBudget
    
    // Per-user budgets
    UserBudgets map[string]UserBudget
    
    // Global budget
    GlobalBudget GlobalBudget
}

type WorkflowBudget struct {
    // Total calls allowed
    MaxCallsPerExecution int
    MaxCallsPerHour      int
    MaxCallsPerDay       int
    
    // Per-endpoint budgets within workflow
    EndpointBudgets map[string]int
    
    // Cost limits (if API charges per call)
    MaxCostPerExecution float64
    MaxCostPerDay       float64
}

// Budget tracker
type BudgetTracker struct {
    used      map[string]*BudgetUsage
    config    APIBudgetConfig
    costTable map[string]float64  // Endpoint -> cost per call
    mu        sync.RWMutex
}

type BudgetUsage struct {
    callCount int
    cost      float64
    resetTime time.Time
}

func (bt *BudgetTracker) CheckBudget(workflowID, endpoint string) error {
    bt.mu.RLock()
    defer bt.mu.RUnlock()
    
    budget := bt.config.WorkflowBudgets[workflowID]
    usage := bt.used[workflowID]
    
    // Check total calls
    if usage.callCount >= budget.MaxCallsPerExecution {
        return fmt.Errorf(
            "workflow budget exceeded: %d/%d calls used",
            usage.callCount,
            budget.MaxCallsPerExecution,
        )
    }
    
    // Check endpoint-specific budget
    if endpointBudget, ok := budget.EndpointBudgets[endpoint]; ok {
        endpointUsage := bt.getEndpointUsage(workflowID, endpoint)
        if endpointUsage >= endpointBudget {
            return fmt.Errorf(
                "endpoint budget exceeded for %s: %d/%d calls used",
                endpoint,
                endpointUsage,
                endpointBudget,
            )
        }
    }
    
    // Check cost budget
    callCost := bt.costTable[endpoint]
    if usage.cost + callCost > budget.MaxCostPerExecution {
        return fmt.Errorf(
            "cost budget exceeded: $%.2f/$%.2f used",
            usage.cost,
            budget.MaxCostPerExecution,
        )
    }
    
    return nil
}
```

**Benefits:**
- ✅ Prevent cost overruns
- ✅ Per-workflow resource allocation
- ✅ Multi-tenant isolation
- ✅ Predictable API usage

### 4.5 Enhanced Retry with Jitter (GAP-6) 🔄

**Design:**

```go
type AdaptiveRetryStrategy struct {
    baseBackoff      time.Duration
    maxBackoff       time.Duration
    multiplier       float64
    jitterFactor     float64
    
    // Adaptive behavior
    successiveFailures int
    lastErrorType      string
}

func (s *AdaptiveRetryStrategy) NextBackoff() time.Duration {
    // Exponential backoff
    backoff := time.Duration(float64(s.baseBackoff) * math.Pow(s.multiplier, float64(s.successiveFailures)))
    
    // Cap at max
    if backoff > s.maxBackoff {
        backoff = s.maxBackoff
    }
    
    // Add jitter to prevent thundering herd
    jitter := time.Duration(rand.Float64() * float64(backoff) * s.jitterFactor)
    backoff = backoff + jitter
    
    return backoff
}

func (s *AdaptiveRetryStrategy) ShouldRetry(err error, attempt int) bool {
    // Don't retry certain errors
    if isNonRetryable(err) {
        return false
    }
    
    // Exponentially reduce retry probability on successive failures
    if s.successiveFailures > 5 {
        // Probabilistic retry
        probability := 1.0 / float64(s.successiveFailures)
        return rand.Float64() < probability
    }
    
    return true
}

// Example configuration
retry := AdaptiveRetryStrategy{
    baseBackoff:   100 * time.Millisecond,
    maxBackoff:    30 * time.Second,
    multiplier:    2.0,
    jitterFactor:  0.3,  // 30% jitter
}
```

**Benefits:**
- ✅ Prevents thundering herd
- ✅ Better distribution of retry attempts
- ✅ Adaptive to API behavior
- ✅ Probabilistic backoff for persistent failures

### 4.6 Enhanced Observability Metrics (GAP-7) 📊

**New Metrics:**

```go
// Per-endpoint metrics
const (
    metricAPICallsByEndpoint      = "api.calls.by_endpoint"
    metricAPIErrorsByEndpoint     = "api.errors.by_endpoint"
    metricAPILatencyByEndpoint    = "api.latency.by_endpoint"
    metricAPIRateLimitHits        = "api.rate_limit_hits"
    metricAPICircuitBreakerState  = "api.circuit_breaker_state"
    metricAPIBudgetUsage          = "api.budget_usage"
    metricAPIAnomalies            = "api.anomalies_detected"
)

// Enhanced telemetry
func (p *Provider) RecordAPICall(
    ctx context.Context,
    endpoint string,
    method string,
    statusCode int,
    duration time.Duration,
    success bool,
) {
    attrs := []attribute.KeyValue{
        attribute.String("api.endpoint", endpoint),
        attribute.String("http.method", method),
        attribute.Int("http.status_code", statusCode),
        attribute.Bool("success", success),
    }
    
    // Record call count
    p.apiCalls.Add(ctx, 1, metric.WithAttributes(attrs...))
    
    // Record latency
    p.apiLatency.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
    
    // Record errors
    if !success {
        p.apiErrors.Add(ctx, 1, metric.WithAttributes(attrs...))
    }
}

// Circuit breaker state metric
func (p *Provider) RecordCircuitBreakerState(ctx context.Context, endpoint string, state CircuitState) {
    attrs := []attribute.KeyValue{
        attribute.String("api.endpoint", endpoint),
        attribute.String("state", state.String()),
    }
    p.circuitBreakerState.Record(ctx, float64(state), metric.WithAttributes(attrs...))
}
```

**Grafana Dashboard Queries:**

```promql
# API calls per endpoint
sum(rate(api_calls_by_endpoint[5m])) by (api_endpoint)

# Error rate by endpoint
sum(rate(api_errors_by_endpoint[5m])) by (api_endpoint) 
  / 
sum(rate(api_calls_by_endpoint[5m])) by (api_endpoint)

# P95 latency by endpoint
histogram_quantile(0.95, api_latency_by_endpoint)

# Rate limit violations
sum(rate(api_rate_limit_hits[5m])) by (api_endpoint)

# Circuit breaker open count
count(api_circuit_breaker_state == 1) by (api_endpoint)
```

**Benefits:**
- ✅ Per-endpoint visibility
- ✅ Early warning of issues
- ✅ Capacity planning
- ✅ Cost attribution

---

## 5. Implementation Plan

### Phase 1: Critical Fixes (Week 1-2) 🔴

**Priority: P0**

**Tasks:**

1. **Change MaxHTTPCallsPerExec Default** (GAP-4)
   - File: `backend/pkg/config/config.go`
   - Change: Set default to `100` instead of `0` (unlimited)
   - Test: Update tests expecting unlimited
   - Documentation: Update README and config docs

2. **Add Per-Endpoint Rate Limiting** (GAP-1)
   - New file: `backend/pkg/middleware/api_ratelimit.go`
   - Implementation: `APIRateLimitMiddleware`
   - Config: Add `APIEndpointLimits` to main config
   - Tests: Comprehensive rate limiting tests
   - Documentation: Add usage examples

3. **Implement Circuit Breaker** (GAP-2)
   - New file: `backend/pkg/middleware/circuitbreaker.go`
   - Implementation: `CircuitBreaker` pattern
   - Integration: Wrap HTTP executor
   - Tests: Circuit breaker state transitions
   - Metrics: Add circuit breaker metrics

4. **Add API Audit Logging** (GAP-3)
   - New file: `backend/pkg/audit/api_audit.go`
   - Implementation: `APIAuditLogger`
   - Integration: Hook into HTTP executor
   - Storage: JSON file + stdout (pluggable)
   - Tests: Audit log generation and querying

**Success Criteria:**
- ✅ All tests pass
- ✅ No regressions in existing functionality
- ✅ Documentation updated
- ✅ Metrics exported correctly

### Phase 2: High Priority Enhancements (Week 3-4) 🟡

**Priority: P1**

**Tasks:**

1. **Add Retry Jitter** (GAP-6)
   - File: `backend/pkg/middleware/retry.go`
   - Enhancement: Add jitter to backoff calculation
   - Config: Add jitter factor configuration
   - Tests: Verify jitter distribution

2. **Implement Anomaly Detection** (GAP-7)
   - New file: `backend/pkg/audit/anomaly.go`
   - Implementation: `AnomalyDetector`
   - Integration: Run periodic detection
   - Alerts: Log anomalies, emit metrics
   - Tests: Anomaly detection scenarios

3. **Add API Budgeting** (GAP-8)
   - New file: `backend/pkg/budget/budget.go`
   - Implementation: `BudgetTracker`
   - Config: Budget configuration structure
   - Integration: Check budget before API calls
   - Tests: Budget enforcement tests

4. **Time-Window Quotas** (GAP-5)
   - Enhancement: Add sliding window counters
   - Implementation: Hourly, daily quotas
   - Storage: In-memory with TTL
   - Tests: Quota reset and enforcement

**Success Criteria:**
- ✅ Anomaly detection working
- ✅ Budget tracking functional
- ✅ Time-window quotas enforced
- ✅ Jitter prevents thundering herd

### Phase 3: Medium Priority Features (Week 5-6) 🟢

**Priority: P2**

**Tasks:**

1. **Distributed Rate Limiting** (GAP-9)
   - Implementation: Redis-backed rate limiter
   - Config: Redis connection settings
   - Fallback: Local rate limiting if Redis unavailable
   - Tests: Multi-instance coordination

2. **Adaptive Retry** (GAP-10)
   - Enhancement: Adaptive backoff based on API behavior
   - Implementation: Success/failure tracking
   - Config: Adaptive strategy parameters
   - Tests: Adaptive behavior scenarios

3. **API Cost Tracking** (GAP-11)
   - Implementation: Cost accumulation per endpoint
   - Integration: Budget tracker
   - Reporting: Cost reports and dashboards
   - Tests: Cost calculation accuracy

4. **Request Prioritization** (GAP-12)
   - Implementation: Priority queue for API calls
   - Config: Priority levels and scheduling
   - Integration: Workflow-level priorities
   - Tests: Priority enforcement

**Success Criteria:**
- ✅ Distributed deployments supported
- ✅ Cost tracking accurate
- ✅ Priority queue working

### Phase 4: Documentation & Operations (Week 7-8) 📚

**Tasks:**

1. **Security Runbook**
   - Incident response procedures
   - Common attack patterns and mitigations
   - Escalation paths
   - Recovery procedures

2. **Operations Guide**
   - Configuration best practices
   - Monitoring and alerting setup
   - Capacity planning
   - Troubleshooting guide

3. **Compliance Documentation**
   - Audit logging procedures
   - Data retention policies
   - Access control guidelines
   - Compliance checklist

4. **Training Materials**
   - Developer guide for secure workflows
   - Admin guide for platform configuration
   - Security awareness training
   - Best practices handbook

**Success Criteria:**
- ✅ Complete runbooks
- ✅ Clear operational procedures
- ✅ Training materials ready

---

## 6. Operational Procedures

### 6.1 Monitoring & Alerting

#### Critical Alerts

```yaml
# Prometheus alert rules
groups:
  - name: api_protection_critical
    interval: 30s
    rules:
      # High API error rate
      - alert: HighAPIErrorRate
        expr: |
          sum(rate(api_errors_by_endpoint[5m])) by (api_endpoint)
          / sum(rate(api_calls_by_endpoint[5m])) by (api_endpoint)
          > 0.5
        for: 2m
        labels:
          severity: critical
          component: api_protection
        annotations:
          summary: "High error rate for {{ $labels.api_endpoint }}"
          description: "{{ $labels.api_endpoint }} error rate is {{ $value | humanizePercentage }}"
          runbook: "docs/runbooks/high-api-error-rate.md"
          
      # Circuit breaker open
      - alert: CircuitBreakerOpen
        expr: api_circuit_breaker_state{state="open"} == 1
        for: 1m
        labels:
          severity: critical
          component: api_protection
        annotations:
          summary: "Circuit breaker open for {{ $labels.api_endpoint }}"
          description: "API {{ $labels.api_endpoint }} circuit breaker is open"
          runbook: "docs/runbooks/circuit-breaker-open.md"
          
      # Rate limit violations
      - alert: HighRateLimitViolations
        expr: sum(rate(api_rate_limit_hits[5m])) by (api_endpoint) > 10
        for: 5m
        labels:
          severity: warning
          component: api_protection
        annotations:
          summary: "High rate limit violations for {{ $labels.api_endpoint }}"
          description: "{{ $value }} violations/sec for {{ $labels.api_endpoint }}"
          
      # Anomaly detected
      - alert: APIAnomalyDetected
        expr: api_anomalies_detected > 0
        labels:
          severity: warning
          component: api_protection
        annotations:
          summary: "API anomaly detected"
          description: "{{ $value }} anomalies detected in API usage patterns"
```

#### Warning Alerts

```yaml
  - name: api_protection_warnings
    interval: 1m
    rules:
      # Budget approaching limit
      - alert: APIBudgetApproachingLimit
        expr: |
          api_budget_usage / api_budget_limit > 0.8
        labels:
          severity: warning
        annotations:
          summary: "API budget at {{ $value | humanizePercentage }} for {{ $labels.workflow_id }}"
          
      # High API latency
      - alert: HighAPILatency
        expr: |
          histogram_quantile(0.95, api_latency_by_endpoint) > 5000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High P95 latency for {{ $labels.api_endpoint }}"
          description: "P95 latency is {{ $value }}ms"
```

### 6.2 Incident Response

#### High Error Rate Runbook

**Symptoms:**
- API error rate > 50%
- Users reporting workflow failures
- Circuit breakers opening

**Investigation:**
1. Check API status page
2. Review error logs for patterns
3. Check if issue is endpoint-specific
4. Verify network connectivity
5. Check rate limiting status

**Mitigation:**
1. If API is down:
   - Circuit breaker will auto-open
   - Monitor for recovery
   - Communicate status to users

2. If rate limiting:
   - Reduce rate limits temporarily
   - Identify abusive workflows
   - Apply stricter budgets

3. If internal issue:
   - Check application logs
   - Verify configuration
   - Rollback recent changes if needed

**Recovery:**
- Circuit breaker will auto-recover
- Monitor metrics post-recovery
- Document incident for post-mortem

#### Circuit Breaker Open Runbook

**Symptoms:**
- Circuit breaker state = "open"
- Fast-fail errors in workflows
- No requests reaching API

**Investigation:**
1. Check API health independently
2. Review recent error patterns
3. Check if manual intervention needed
4. Verify circuit breaker configuration

**Actions:**
1. If API healthy:
   - Wait for half-open state
   - Monitor recovery attempts
   - Manually reset if necessary

2. If API down:
   - Communicate to stakeholders
   - Estimate recovery time
   - Prepare fallback strategies

### 6.3 Security Incident Response

#### Suspected DDOS Attack

**Detection:**
- Abnormal spike in API calls
- Rate limit violations from single workflow/user
- Anomaly alerts firing
- API quotas exhausted

**Immediate Actions:**
1. Identify attacking workflow/user
2. Block execution immediately
3. Apply emergency rate limits
4. Review audit logs

**Investigation:**
1. Check workflow definition
2. Review execution history
3. Identify pattern (loop, parallel, etc.)
4. Determine if malicious or accidental

**Mitigation:**
1. Apply per-workflow rate limits
2. Reduce global rate limits temporarily
3. Enable stricter circuit breakers
4. Block malicious users/workflows

**Recovery:**
1. Restore normal rate limits gradually
2. Monitor for recurrence
3. Update security controls
4. Document lessons learned

---

## 7. Compliance and Governance

### 7.1 Audit Requirements

#### What to Log

**Required Fields:**
- ✅ Execution ID, Workflow ID, User ID
- ✅ Timestamp (UTC, ISO 8601)
- ✅ API endpoint and method
- ✅ Request/response metadata (size, duration)
- ✅ Success/failure status
- ✅ Error details (sanitized)
- ✅ Rate limit status
- ✅ Circuit breaker state

**Optional Fields:**
- IP address (for security)
- Geographic location
- Cost attribution
- Data classification tags

#### Retention Policies

| Log Type | Retention | Storage | Access |
|----------|-----------|---------|--------|
| API Call Logs | 90 days | Hot storage | All admins |
| Security Events | 1 year | Warm storage | Security team |
| Compliance Logs | 7 years | Cold storage | Compliance team |
| Debug Logs | 7 days | Hot storage | Developers |

### 7.2 Access Controls

#### RBAC Model

```yaml
roles:
  - name: developer
    permissions:
      - view_workflow
      - execute_workflow
      - view_logs
      
  - name: admin
    permissions:
      - all_developer_permissions
      - configure_rate_limits
      - view_audit_logs
      - manage_budgets
      
  - name: security
    permissions:
      - view_all_logs
      - block_workflows
      - configure_security
      - review_anomalies
      
  - name: compliance
    permissions:
      - view_audit_logs
      - export_compliance_reports
      - configure_retention
```

### 7.3 Security Controls Checklist

**Pre-Production:**
- [ ] All rate limits configured
- [ ] Circuit breakers enabled
- [ ] Audit logging verified
- [ ] Budget limits set
- [ ] Monitoring dashboards created
- [ ] Alerts configured and tested
- [ ] Runbooks documented
- [ ] Incident response team trained

**Production:**
- [ ] Daily anomaly review
- [ ] Weekly security log review
- [ ] Monthly budget review
- [ ] Quarterly security audit
- [ ] Annual penetration testing
- [ ] Continuous compliance monitoring

---

## 8. Conclusion

### 8.1 Current State Summary

**Strengths:**
- ✅ Strong SSRF protection
- ✅ Good connection pooling
- ✅ Basic rate limiting exists
- ✅ Retry with backoff
- ✅ Comprehensive resource limits
- ✅ Good observability foundation

**Weaknesses:**
- ⚠️ No per-endpoint rate limiting
- ⚠️ No circuit breaker pattern
- ⚠️ Limited audit logging
- ⚠️ Default MaxHTTPCallsPerExec = unlimited
- ⚠️ No anomaly detection

### 8.2 Proposed State

**After Implementation:**
- ✅ Per-endpoint rate limiting with time windows
- ✅ Circuit breaker pattern for resilience
- ✅ Comprehensive audit logging with anomaly detection
- ✅ API budgeting and cost tracking
- ✅ Enhanced retry with jitter
- ✅ Rich metrics and dashboards
- ✅ Complete incident response procedures
- ✅ Compliance-ready audit trail

### 8.3 Risk Reduction

| Risk | Before | After | Reduction |
|------|--------|-------|-----------|
| DDOS Attack | HIGH | LOW | 80% |
| Retry Storm | HIGH | LOW | 85% |
| Cost Overrun | MEDIUM | LOW | 70% |
| Cascade Failure | MEDIUM | LOW | 75% |
| Security Breach | LOW | VERY LOW | 50% |

### 8.4 Next Steps

1. **Review and Approve** this security analysis
2. **Prioritize** implementation phases
3. **Assign** development resources
4. **Begin Phase 1** critical fixes
5. **Establish** monitoring and alerting
6. **Train** operations team
7. **Document** all procedures
8. **Continuous improvement** based on metrics

---

## Appendices

### A. Configuration Examples

See: `docs/examples/api-protection-config.yaml`

### B. Metrics Reference

See: `docs/API_PROTECTION_METRICS.md`

### C. Incident Response Procedures

See: `docs/runbooks/`

### D. Compliance Checklist

See: `docs/COMPLIANCE_CHECKLIST.md`

---

**Document Version:** 1.0  
**Last Updated:** 2025-11-05  
**Author:** Security Architecture Team  
**Reviewers:** Engineering, Operations, Security, Compliance  
**Status:** DRAFT - Pending Review
