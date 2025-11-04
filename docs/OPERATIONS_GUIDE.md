# Thaiyyal Workflow Engine - Operations Guide

## Table of Contents

- [Overview](#overview)
- [Deployment](#deployment)
- [Configuration](#configuration)
- [Monitoring](#monitoring)
- [Scaling](#scaling)
- [Troubleshooting](#troubleshooting)
- [Maintenance](#maintenance)

## Overview

This guide provides operational procedures for running the Thaiyyal Workflow Engine in production environments.

### Service Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Production Stack                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│  │  Load        │    │  Thaiyyal    │    │  Prometheus  │ │
│  │  Balancer    │───▶│  Server      │───▶│  (Metrics)   │ │
│  └──────────────┘    │  (3+ replicas)│    └──────────────┘ │
│                      └──────────────┘           │          │
│                                                 │          │
│                                      ┌──────────▼────────┐ │
│                                      │  Grafana          │ │
│                                      │  (Dashboards)     │ │
│                                      └───────────────────┘ │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Deployment

### Docker Deployment

#### Quick Start

```bash
# Build the image
docker build -t thaiyyal/workflow-engine:latest .

# Run the container
docker run -d \
  --name thaiyyal-server \
  -p 8080:8080 \
  -e MAX_EXECUTION_TIME=5m \
  -e MAX_NODE_EXECUTIONS=10000 \
  thaiyyal/workflow-engine:latest

# Check health
curl http://localhost:8080/health
```

#### Docker Compose

```bash
# Start all services (server, Prometheus, Grafana)
docker-compose up -d

# View logs
docker-compose logs -f thaiyyal-server

# Stop all services
docker-compose down
```

### Kubernetes Deployment

#### Prerequisites

- Kubernetes cluster (1.24+)
- kubectl configured
- Prometheus Operator (optional, for ServiceMonitor)

#### Deploy

```bash
# Create namespace and deploy
kubectl apply -f deployments/kubernetes/deployment.yaml

# Verify deployment
kubectl get pods -n thaiyyal
kubectl get svc -n thaiyyal

# Check logs
kubectl logs -n thaiyyal -l app=thaiyyal-server

# Port forward for local access
kubectl port-forward -n thaiyyal svc/thaiyyal-server 8080:8080
```

#### Scaling

```bash
# Manual scaling
kubectl scale deployment -n thaiyyal thaiyyal-server --replicas=5

# HPA automatically scales based on CPU/memory
kubectl get hpa -n thaiyyal
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_ADDR` | `:8080` | Server listen address |
| `MAX_EXECUTION_TIME` | `5m` | Maximum workflow execution time |
| `MAX_NODE_EXECUTIONS` | `10000` | Maximum node executions per workflow |
| `MAX_HTTP_CALLS` | `100` | Maximum HTTP calls per execution |
| `MAX_LOOP_ITERATIONS` | `10000` | Maximum loop iterations |

### Command-Line Flags

```bash
./thaiyyal-server \
  -addr :8080 \
  -read-timeout 30s \
  -write-timeout 30s \
  -max-execution-time 5m \
  -max-node-executions 10000 \
  -max-http-calls 100 \
  -max-loop-iterations 10000
```

### Configuration Profiles

#### Development

```yaml
MAX_EXECUTION_TIME: 10m
MAX_NODE_EXECUTIONS: 50000
MAX_HTTP_CALLS: 500
```

#### Production

```yaml
MAX_EXECUTION_TIME: 5m
MAX_NODE_EXECUTIONS: 10000
MAX_HTTP_CALLS: 100
```

#### Strict (API Gateway)

```yaml
MAX_EXECUTION_TIME: 30s
MAX_NODE_EXECUTIONS: 1000
MAX_HTTP_CALLS: 10
```

### Security Considerations

#### Production Deployment

**⚠️ IMPORTANT: Profiling Endpoints**

The server exposes pprof endpoints at `/debug/pprof/*` for performance profiling. These endpoints provide sensitive runtime information and should be restricted in production:

**Option 1: Use a separate admin port (Recommended)**
```go
// Run pprof on separate admin port with restricted access
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

**Option 2: Add authentication middleware**
- Require API keys or JWT tokens
- Implement IP allowlisting
- Use mutual TLS (mTLS)

**Option 3: Disable in production builds**
```bash
# Build without pprof
go build -tags nopprof ./cmd/server
```

**Network-level protection:**
- Use network policies to restrict access
- Place behind firewall or VPN
- Only expose to monitoring/admin networks

## Monitoring

### Health Endpoints

#### Liveness Probe

```bash
# Check if server is running
curl http://localhost:8080/health/live

# Expected response:
{
  "status": "healthy",
  "service_name": "thaiyyal-workflow-engine",
  "service_version": "0.1.0",
  "uptime": "2h15m30s",
  "timestamp": "2025-11-03T22:00:00Z"
}
```

#### Readiness Probe

```bash
# Check if server can handle requests
curl http://localhost:8080/health/ready

# HTTP 200 = ready, 503 = not ready
```

#### Full Health Check

```bash
# Run all health checks
curl http://localhost:8080/health
```

### Metrics

#### Prometheus Metrics

Access metrics at `http://localhost:8080/metrics`

**Key Metrics:**

```promql
# Workflow execution rate
rate(workflow_executions_total[5m])

# Workflow execution duration (p99)
histogram_quantile(0.99, rate(workflow_execution_duration_bucket[5m]))

# Node execution success rate
rate(node_executions_success_total[5m]) / rate(node_executions_total[5m])

# HTTP call duration
histogram_quantile(0.95, rate(http_call_duration_bucket[5m]))
```

#### Grafana Dashboards

Access Grafana at `http://localhost:3001` (default credentials: admin/admin)

**Pre-configured Dashboards:**
- Workflow Execution Overview
- Node Performance Analysis
- HTTP Call Statistics
- System Health

### Logging

Logs are output in structured JSON format to stdout:

```json
{
  "time": "2025-11-03T22:00:00Z",
  "level": "INFO",
  "msg": "workflow execution completed successfully",
  "workflow_id": "wf-123",
  "execution_id": "exec-abc",
  "duration_ms": 150,
  "nodes_executed": 5
}
```

**Log Levels:**
- `DEBUG`: Detailed execution information
- `INFO`: Normal operation events
- `WARN`: Warning conditions
- `ERROR`: Error conditions

## Scaling

### Horizontal Scaling

The server is stateless and can be scaled horizontally:

```bash
# Docker Compose
docker-compose up -d --scale thaiyyal-server=5

# Kubernetes
kubectl scale deployment -n thaiyyal thaiyyal-server --replicas=10
```

### Auto-scaling

Kubernetes HPA automatically scales based on:
- CPU utilization (target: 70%)
- Memory utilization (target: 80%)

```bash
# Monitor auto-scaling
kubectl get hpa -n thaiyyal -w
```

### Performance Tuning

#### Resource Limits

Adjust based on workflow complexity:

```yaml
resources:
  requests:
    cpu: 100m      # Minimum CPU
    memory: 128Mi  # Minimum memory
  limits:
    cpu: 500m      # Maximum CPU
    memory: 512Mi  # Maximum memory
```

#### Execution Limits

For high-throughput scenarios:

```bash
-max-execution-time 30s \
-max-node-executions 5000 \
-max-http-calls 50
```

## Troubleshooting

### Common Issues

#### 1. Workflow Execution Timeout

**Symptoms:**
- 500 error with "execution timeout" message
- Metrics show high execution duration

**Solutions:**
- Increase `MAX_EXECUTION_TIME`
- Optimize workflow (reduce node count)
- Check for slow HTTP calls

#### 2. High Memory Usage

**Symptoms:**
- Pods being OOMKilled
- High memory metrics

**Solutions:**
- Reduce `MAX_ARRAY_LENGTH` in workflow config
- Limit data size in HTTP responses
- Increase memory limits

#### 3. Service Unhealthy

**Symptoms:**
- `/health/ready` returns 503
- Pods marked as not ready

**Solutions:**
- Check health check logs
- Verify dependencies are accessible
- Review recent deployments

### Debug Mode

Enable verbose logging:

```bash
# Set log level to debug
export LOG_LEVEL=debug
./thaiyyal-server
```

### Metrics Analysis

```bash
# Check workflow execution rate
curl http://localhost:8080/metrics | grep workflow_executions_total

# Check error rate
curl http://localhost:8080/metrics | grep workflow_executions_failure_total
```

## Maintenance

### Graceful Shutdown

The server supports graceful shutdown with in-flight request completion:

```bash
# Send SIGTERM
kill -TERM <pid>

# Server will:
# 1. Stop accepting new requests
# 2. Complete in-flight requests (up to shutdown timeout)
# 3. Shutdown telemetry
# 4. Exit
```

### Backup and Recovery

**State:**
- Server is stateless
- No persistent data to back up

**Configuration:**
- Back up Kubernetes ConfigMaps
- Version control deployment manifests

### Updates and Rollouts

#### Rolling Update (Kubernetes)

```bash
# Update image
kubectl set image deployment/thaiyyal-server \
  -n thaiyyal \
  thaiyyal-server=thaiyyal/workflow-engine:v0.2.0

# Monitor rollout
kubectl rollout status deployment/thaiyyal-server -n thaiyyal

# Rollback if needed
kubectl rollout undo deployment/thaiyyal-server -n thaiyyal
```

#### Blue-Green Deployment

```bash
# Deploy new version alongside old
kubectl apply -f deployments/kubernetes/deployment-v2.yaml

# Test new version
kubectl port-forward svc/thaiyyal-server-v2 8081:8080

# Switch traffic
kubectl patch svc thaiyyal-server -p '{"spec":{"selector":{"version":"v2"}}}'
```

### Monitoring Best Practices

1. **Set up alerts** for:
   - High error rate (>5%)
   - Slow execution (p95 > 5s)
   - Service unavailability
   - Resource exhaustion

2. **Regular health checks**:
   - Monitor `/health` endpoint
   - Track uptime metrics
   - Review error logs

3. **Capacity planning**:
   - Monitor request rate trends
   - Track resource usage
   - Plan for peak loads

## Support

For issues and questions:
- GitHub Issues: https://github.com/yesoreyeram/thaiyyal/issues
- Documentation: https://github.com/yesoreyeram/thaiyyal/tree/main/docs

---

**Last Updated:** 2025-11-03  
**Version:** 0.1.0
