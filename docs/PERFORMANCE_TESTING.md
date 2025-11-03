# Performance Testing Guide

This guide covers performance testing, profiling, and optimization for the Thaiyyal Workflow Engine.

## Table of Contents

- [Performance Profiling](#performance-profiling)
- [Benchmarking](#benchmarking)
- [Load Testing](#load-testing)
- [Optimization Guidelines](#optimization-guidelines)
- [Monitoring Performance](#monitoring-performance)

## Performance Profiling

### CPU Profiling with pprof

The server includes pprof endpoints for runtime profiling:

```bash
# Start the server
./backend/bin/thaiyyal-server

# Capture 30-second CPU profile
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof

# Analyze the profile
go tool pprof cpu.prof

# Interactive commands in pprof:
# - top: Show top functions by CPU usage
# - list <function>: Show source code for function
# - web: Generate graph visualization (requires graphviz)
```

### Memory Profiling

```bash
# Capture heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof

# Analyze memory allocation
go tool pprof heap.prof

# Commands:
# - top: Show top allocators
# - list <function>: Show allocation sources
# - web: Visualize allocation graph
```

### Goroutine Analysis

```bash
# Check goroutine count and stack traces
curl http://localhost:8080/debug/pprof/goroutine?debug=1

# Get goroutine profile
curl http://localhost:8080/debug/pprof/goroutine > goroutine.prof
go tool pprof goroutine.prof
```

### Available pprof Endpoints

- `/debug/pprof/` - Profile index
- `/debug/pprof/profile` - CPU profile
- `/debug/pprof/heap` - Heap profile
- `/debug/pprof/goroutine` - Goroutine stack traces
- `/debug/pprof/block` - Blocking profile
- `/debug/pprof/mutex` - Mutex contention profile
- `/debug/pprof/allocs` - All past memory allocations
- `/debug/pprof/threadcreate` - Thread creation profile

## Benchmarking

### Go Benchmarks

Run the built-in benchmarks:

```bash
cd backend

# Run all benchmarks
go test -bench=. ./...

# Run specific package benchmarks
go test -bench=. ./pkg/engine/

# With memory allocation stats
go test -bench=. -benchmem ./pkg/engine/

# Save benchmark results
go test -bench=. -benchmem ./pkg/engine/ > bench.txt

# Compare benchmarks
go test -bench=. -benchmem ./pkg/engine/ > new.txt
benchstat old.txt new.txt
```

### Create Custom Benchmarks

```go
func BenchmarkWorkflowExecution(b *testing.B) {
    payload := []byte(`{"nodes":[...],"edges":[...]}`)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        eng, _ := engine.New(payload)
        _, _ = eng.Execute()
    }
}
```

## Load Testing

### Using Apache Bench (ab)

```bash
# Simple load test
ab -n 1000 -c 10 \
  -p workflow.json \
  -T application/json \
  http://localhost:8080/api/v1/workflow/execute

# Results show:
# - Requests per second
# - Time per request
# - Transfer rate
# - Percentile distribution
```

### Using wrk

```bash
# Install wrk
# brew install wrk (macOS)
# apt install wrk (Ubuntu)

# Create Lua script for POST requests
cat > post.lua << 'EOF'
wrk.method = "POST"
wrk.body   = '{"nodes":[...],"edges":[...]}'
wrk.headers["Content-Type"] = "application/json"
EOF

# Run load test
wrk -t4 -c100 -d30s \
  -s post.lua \
  http://localhost:8080/api/v1/workflow/execute

# Results:
# - Latency distribution
# - Requests/sec
# - Transfer/sec
```

### Using k6

```javascript
// load-test.js
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 },  // Ramp up
    { duration: '1m', target: 50 },   // Steady state
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% requests < 500ms
    http_req_failed: ['rate<0.01'],   // Error rate < 1%
  },
};

const workflow = JSON.stringify({
  nodes: [...],
  edges: [...],
});

export default function () {
  const res = http.post(
    'http://localhost:8080/api/v1/workflow/execute',
    workflow,
    { headers: { 'Content-Type': 'application/json' } }
  );
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
}
```

Run the test:

```bash
k6 run load-test.js
```

## Optimization Guidelines

### Workflow Optimization

1. **Minimize HTTP Calls**
   - Batch requests when possible
   - Use caching for repeated calls
   - Consider parallel execution

2. **Reduce Node Count**
   - Combine simple operations
   - Use built-in aggregation functions
   - Eliminate unnecessary transformations

3. **Optimize Data Flow**
   - Minimize data passed between nodes
   - Filter early, reduce late
   - Use efficient data structures

### Server Configuration

```bash
# Increase execution limits for batch processing
./thaiyyal-server \
  -max-execution-time 10m \
  -max-node-executions 50000

# Use strict limits for API endpoints
./thaiyyal-server \
  -max-execution-time 30s \
  -max-node-executions 1000
```

### Kubernetes Resource Tuning

```yaml
resources:
  requests:
    cpu: 200m      # Baseline
    memory: 256Mi
  limits:
    cpu: 1000m     # Burst capacity
    memory: 1Gi
```

## Monitoring Performance

### Key Metrics to Watch

```promql
# Request rate
rate(workflow_executions_total[5m])

# Execution time (p50, p95, p99)
histogram_quantile(0.50, rate(workflow_execution_duration_bucket[5m]))
histogram_quantile(0.95, rate(workflow_execution_duration_bucket[5m]))
histogram_quantile(0.99, rate(workflow_execution_duration_bucket[5m]))

# Error rate
rate(workflow_executions_failure_total[5m]) / rate(workflow_executions_total[5m])

# Node execution time by type
histogram_quantile(0.95, rate(node_execution_duration_bucket[5m])) by (node_type)

# HTTP call latency
histogram_quantile(0.95, rate(http_call_duration_bucket[5m]))
```

### Setting Alerts

```yaml
# Prometheus alert rules
groups:
  - name: thaiyyal
    rules:
      - alert: HighErrorRate
        expr: rate(workflow_executions_failure_total[5m]) > 0.05
        for: 5m
        annotations:
          summary: High workflow error rate
          
      - alert: SlowExecution
        expr: histogram_quantile(0.95, rate(workflow_execution_duration_bucket[5m])) > 5000
        for: 5m
        annotations:
          summary: Slow workflow execution (p95 > 5s)
          
      - alert: HighCPU
        expr: rate(process_cpu_seconds_total[5m]) > 0.8
        for: 5m
        annotations:
          summary: High CPU usage
```

### Performance Checklist

- [ ] Profile CPU usage under load
- [ ] Check memory allocation patterns
- [ ] Monitor goroutine count
- [ ] Verify no memory leaks
- [ ] Test with production-like workloads
- [ ] Set up performance alerts
- [ ] Document performance baselines
- [ ] Create capacity planning guidelines

## Example Optimization Workflow

### Before Optimization

```json
{
  "nodes": [
    {"id": "fetch1", "type": "http", "data": {"url": "..."}},
    {"id": "extract1", "type": "extract", "data": {"field": "data"}},
    {"id": "fetch2", "type": "http", "data": {"url": "..."}},
    {"id": "extract2", "type": "extract", "data": {"field": "data"}},
    {"id": "merge", "type": "join"},
    {"id": "filter", "type": "filter", "data": {"expression": "..."}}
  ]
}
```

**Issues:**
- Sequential HTTP calls
- Extracting before merging
- No caching

### After Optimization

```json
{
  "nodes": [
    {"id": "parallel", "type": "parallel"},
    {"id": "fetch1", "type": "http", "data": {"url": "...", "cache": true}},
    {"id": "fetch2", "type": "http", "data": {"url": "...", "cache": true}},
    {"id": "merge", "type": "join"},
    {"id": "filter_extract", "type": "transform", "data": {"expression": "..."}}
  ]
}
```

**Improvements:**
- Parallel HTTP calls (2x faster)
- Combined filter + extract
- Caching enabled

## Resources

- [Go Profiling Guide](https://go.dev/blog/pprof)
- [Performance Tuning](PERFORMANCE_TUNING.md)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/)

---

**Last Updated:** 2025-11-03  
**Version:** 0.1.0
