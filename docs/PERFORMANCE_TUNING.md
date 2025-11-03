# Performance Tuning

## Configuration Tuning

### Execution Limits

```go
// Relaxed limits for batch processing
config := types.DefaultConfig()
config.MaxExecutionTime = 5 * time.Minute
config.MaxNodeExecutions = 100000

// Strict limits for API endpoints
config := types.ValidationLimits()
config.MaxExecutionTime = 10 * time.Second
config.MaxNodeExecutions = 1000
```

### Caching

```go
// Use cache node for expensive operations
{
  "type": "cache",
  "data": {
    "cacheKey": "api_result",
    "cacheOp": "get",
    "ttl": 300
  }
}
```

## Optimization Strategies

1. **Minimize HTTP calls**: Batch requests when possible
2. **Use caching**: Cache expensive computations
3. **Optimize loops**: Reduce iteration count
4. **Parallel execution**: Independent nodes run concurrently
5. **Data size**: Keep payloads small

## Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./pkg/engine/

# Memory profiling
go test -memprofile=mem.prof -bench=. ./pkg/engine/

# Analyze profile
go tool pprof cpu.prof
```

## Monitoring

Monitor these metrics:
- Execution duration
- Node execution count
- Memory usage
- HTTP call count

---

**Last Updated:** 2025-11-03
**Version:** 1.0
