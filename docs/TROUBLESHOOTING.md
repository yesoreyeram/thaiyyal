# Troubleshooting Guide

## Common Issues

### Workflow Timeout

**Symptom:** Workflow times out after 30 seconds

**Solution:**
```go
// Increase timeout
config := types.DefaultConfig()
config.MaxExecutionTime = 5 * time.Minute
engine, _ := workflow.NewWithConfig(payload, config)
```

### Node Execution Limit

**Symptom:** Error "maximum node executions exceeded"

**Solution:**
```go
// Increase limit or optimize workflow
config.MaxNodeExecutions = 100000
```

### SSRF Blocked

**Symptom:** HTTP request blocked

**Solution:**
- Check if URL uses private IP
- Check if URL is localhost
- Use allowlist for trusted domains

### Memory Issues

**Symptom:** Out of memory errors

**Solution:**
- Reduce array sizes
- Limit object depth
- Process data in chunks
- Use pagination

## Debugging

### Enable Debug Logging

```go
cfg := logging.DefaultConfig()
cfg.Level = "debug"
cfg.Pretty = true
logger := logging.New(cfg)
```

### Use Observers

```go
engine.RegisterObserver(workflow.NewConsoleObserver())
```

### Check Execution Metrics

```go
metrics := engine.GetMetrics()
fmt.Printf("Nodes: %d, HTTP calls: %d\n", 
    metrics.NodeExecutions, metrics.HTTPCalls)
```

---

**Last Updated:** 2025-11-03
**Version:** 1.0
