# API Protection - Quick Reference Guide

> **⚠️ Important**: Default MaxHTTPCallsPerExec changed from unlimited to 100 calls per execution

## Quick Configuration

### Production-Ready Config

```go
config := config.Production()
// Defaults:
// - MaxHTTPCallsPerExec: 100 (NEW - was unlimited)
// - MaxExecutionTime: 5 minutes
// - AllowHTTP: false (HTTPS only)
// - AllowPrivateIPs: false (blocked)
// - AllowLocalhost: false (blocked)
```

### High-Volume Workflows

```go
config := config.Default()
config.MaxHTTPCallsPerExec = 500  // Explicitly increase if needed
```

### Development Config

```go
config := config.Development()
// Relaxed limits for dev/testing
// - MaxHTTPCallsPerExec: 100 (same as default, can be increased if needed)
// - AllowHTTP: true
// - AllowPrivateIPs: true
// - AllowLocalhost: true
```

## Protection Layers

| Layer | Default | Description | Config |
|-------|---------|-------------|--------|
| **HTTP Call Limit** | 100 calls | Per workflow execution | `MaxHTTPCallsPerExec` |
| **Execution Timeout** | 5 minutes | Total workflow time | `MaxExecutionTime` |
| **Node Executions** | 10,000 | Total node executions | `MaxNodeExecutions` |
| **Response Size** | 10 MB | Max HTTP response | `MaxResponseSize` |
| **SSRF Protection** | Enabled | Private IPs blocked | `AllowPrivateIPs`, etc. |
| **Connection Pool** | 100/host | Max connections | Built-in |

## Common Issues

### "maximum HTTP calls per execution exceeded"

**Cause:** Workflow made more than configured HTTP calls (default: 100)

**Solutions:**
1. **Increase limit** if legitimate:
   ```go
   config.MaxHTTPCallsPerExec = 500
   ```

2. **Optimize workflow**: Reduce redundant API calls

3. **Use caching**: Cache API responses

4. **Batch requests**: Combine multiple calls if API supports it

### "HTTP requests are not allowed (AllowHTTP=false)"

**Cause:** Trying to use HTTP instead of HTTPS in production

**Solutions:**
1. **Use HTTPS** (recommended):
   ```json
   {"url": "https://api.example.com"}
   ```

2. **Allow HTTP** in dev only:
   ```go
   config.AllowHTTP = true  // Development only!
   ```

### "domain X is not in allowed domains list"

**Cause:** Domain whitelisting is configured and domain not allowed

**Solutions:**
1. **Add domain** to whitelist:
   ```go
   config.AllowedDomains = []string{
       "api.github.com",
       "api.example.com",
   }
   ```

2. **Clear whitelist** to allow all:
   ```go
   config.AllowedDomains = nil  // Allow all external domains
   ```

### "private IP addresses are blocked"

**Cause:** Trying to call private/internal IPs

**Solutions:**
1. **Use public endpoint** (recommended)

2. **Allow private IPs** (dev/testing only):
   ```go
   config.AllowPrivateIPs = true  // BE CAREFUL!
   ```

## Monitoring Checklist

### Metrics to Watch

```promql
# HTTP calls per workflow
sum(rate(http_calls_total[5m])) by (workflow_id)

# HTTP call limit violations
sum(rate(http_calls_exceeded_total[5m]))

# Workflow execution duration
histogram_quantile(0.95, workflow_execution_duration)

# Workflow timeout rate
sum(rate(workflow_timeouts_total[5m]))
```

### Alerts to Configure

```yaml
# Critical: HTTP call limit frequently hit
- alert: HighHTTPCallLimitViolations
  expr: rate(http_calls_exceeded_total[5m]) > 1
  severity: warning
  
# Warning: Workflow approaching execution timeout
- alert: LongRunningWorkflows
  expr: workflow_execution_duration > 240s
  severity: info
```

## Security Checklist

### Production Deployment

- [ ] `MaxHTTPCallsPerExec` configured appropriately
- [ ] `AllowHTTP = false` (require HTTPS)
- [ ] `AllowPrivateIPs = false` (block private networks)
- [ ] `AllowLocalhost = false` (block localhost)
- [ ] `AllowCloudMetadata = false` (block cloud metadata)
- [ ] Domain whitelist configured (optional but recommended)
- [ ] Monitoring dashboards set up
- [ ] Alerts configured
- [ ] Incident response plan documented

### Development Environment

- [ ] `AllowHTTP = true` (optional for dev)
- [ ] `AllowPrivateIPs = true` (if testing locally)
- [ ] `AllowLocalhost = true` (if testing locally)
- [ ] Still use reasonable `MaxHTTPCallsPerExec`

## Migration Guide

### Updating from Unlimited HTTP Calls

**Before (unlimited):**
```go
config := config.Default()
// MaxHTTPCallsPerExec was 0 (unlimited)
```

**After (limited to 100):**
```go
config := config.Default()
// MaxHTTPCallsPerExec is now 100

// If workflow needs more:
config.MaxHTTPCallsPerExec = 500  // Explicitly set
```

**Testing Migration:**
1. Check workflow execution logs for "maximum HTTP calls" errors
2. Count actual HTTP calls in typical execution
3. Set limit to `actual_calls * 1.5` (50% buffer)
4. Monitor for violations

## Best Practices

### 1. Start with Strict Limits
```go
// Production
config := config.Production()

// Or custom strict limits
config := config.Default()
config.MaxHTTPCallsPerExec = 50
config.MaxExecutionTime = 30 * time.Second
```

### 2. Monitor and Adjust
- Start conservative
- Monitor violations
- Increase only when necessary
- Document reasons for increases

### 3. Use Different Configs per Use Case
```go
// User-facing workflows: strict
apiConfig := config.ValidationLimits()
apiConfig.MaxHTTPCallsPerExec = 10

// Batch processing: relaxed
batchConfig := config.Default()
batchConfig.MaxHTTPCallsPerExec = 1000
batchConfig.MaxExecutionTime = 30 * time.Minute
```

### 4. Implement Caching
```json
{
  "nodes": [
    {
      "id": "cache1",
      "type": "cache",
      "data": {
        "key": "api-response-{{input}}",
        "ttl": 3600
      }
    },
    {
      "id": "http1",
      "type": "http",
      "data": {
        "url": "https://api.example.com/data"
      }
    }
  ]
}
```

### 5. Use Retry Wisely
```json
{
  "type": "retry",
  "data": {
    "maxAttempts": 3,
    "backoff": "1s"
  }
}
```

## Quick Troubleshooting

| Error | Check | Fix |
|-------|-------|-----|
| HTTP call limit | Count API calls in workflow | Increase limit or optimize |
| Execution timeout | Check slow API responses | Increase timeout or optimize |
| SSRF blocked | Check URL targets private IP | Use public endpoint |
| Domain not allowed | Check whitelist | Add domain to AllowedDomains |
| Connection refused | Check API availability | Verify endpoint URL |

## Emergency Procedures

### API Under Attack
1. **Identify** attacking workflow/user
2. **Block** workflow execution
3. **Apply** stricter rate limits globally
4. **Monitor** for continued abuse
5. **Review** logs for patterns

### API Rate Limit Hit
1. **Check** which workflows are calling API
2. **Reduce** `MaxHTTPCallsPerExec` temporarily
3. **Enable** circuit breaker (Phase 1 feature)
4. **Communicate** with API provider
5. **Optimize** workflows to reduce calls

### Service Degradation
1. **Check** metrics dashboard
2. **Review** recent workflow changes
3. **Check** external API status
4. **Apply** temporary stricter limits
5. **Scale** horizontally if needed

## Related Documentation

- [Full Security Analysis](API_PROTECTION_SECURITY_ANALYSIS.md) - Comprehensive 36KB analysis
- [Security Summary](../SECURITY_SUMMARY.md) - Executive summary
- [Workload Protection](PRINCIPLES_WORKLOAD_PROTECTION.md) - Protection mechanisms
- [Security Best Practices](SECURITY_BEST_PRACTICES.md) - Implementation guide
- [Operations Guide](OPERATIONS_GUIDE.md) - Production operations

---

**Quick Reference Version:** 1.0  
**Last Updated:** 2025-11-05  
**For:** Operators, SREs, Platform Engineers
