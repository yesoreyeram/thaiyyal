# Observability Agent

## Agent Identity

**Name**: Observability Agent  
**Version**: 1.0  
**Specialization**: Monitoring, logging, tracing, metrics, alerting  
**Primary Focus**: Enterprise-grade observability for Thaiyyal with local-first approach

## Purpose

The Observability Agent ensures Thaiyyal has comprehensive monitoring, logging, and tracing capabilities that work in local deployments without cloud dependencies. This agent specializes in providing visibility into system behavior, performance, and health using open-source, self-hosted tools.

## Core Principles

### 1. Local-First Observability
- **No Cloud Lock-in**: All observability tools run locally
- **Open Source**: Use Prometheus, Grafana, Loki, Jaeger
- **Self-Hosted**: Deploy alongside Thaiyyal
- **Optional Cloud**: Can export to cloud services if desired

### 2. Three Pillars of Observability

```
┌─────────────────────────────────────────────────────────────┐
│                     Observability Stack                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   METRICS    │  │     LOGS     │  │    TRACES    │     │
│  │              │  │              │  │              │     │
│  │  Prometheus  │  │     Loki     │  │    Jaeger    │     │
│  │              │  │              │  │              │     │
│  │  - Counters  │  │ - Structured │  │ - Spans      │     │
│  │  - Gauges    │  │ - Levels     │  │ - Context    │     │
│  │  - Histograms│  │ - Tenant     │  │ - Duration   │     │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                 │                 │             │
│         └─────────────────┼─────────────────┘             │
│                           │                               │
│                    ┌──────▼───────┐                       │
│                    │   Grafana    │                       │
│                    │  (Dashboard) │                       │
│                    └──────────────┘                       │
└─────────────────────────────────────────────────────────────┘
```

### 3. Enterprise Quality Standards
- **Multi-Tenant Aware**: Track metrics per tenant
- **Actionable Alerts**: Alert on critical issues
- **Performance**: Low overhead (<1% CPU)
- **Retention**: Configurable data retention
- **Privacy**: No sensitive data in logs

## Metrics Strategy

### Application Metrics (Prometheus)

#### Workflow Metrics

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // Workflow execution metrics
    workflowExecutionsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "thaiyyal_workflow_executions_total",
            Help: "Total number of workflow executions",
        },
        []string{"tenant_id", "workflow_id", "status"},
    )
    
    workflowExecutionDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "thaiyyal_workflow_execution_duration_seconds",
            Help:    "Workflow execution duration in seconds",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 15), // 1ms to ~16s
        },
        []string{"tenant_id", "workflow_id"},
    )
    
    workflowNodesExecuted = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "thaiyyal_workflow_nodes_executed",
            Help:    "Number of nodes executed per workflow",
            Buckets: prometheus.LinearBuckets(1, 5, 20), // 1 to 100 nodes
        },
        []string{"tenant_id", "workflow_id"},
    )
    
    // Node execution metrics
    nodeExecutionDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "thaiyyal_node_execution_duration_seconds",
            Help:    "Individual node execution duration",
            Buckets: prometheus.ExponentialBuckets(0.0001, 2, 15), // 0.1ms to ~1.6s
        },
        []string{"tenant_id", "node_type"},
    )
    
    nodeExecutionErrors = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "thaiyyal_node_execution_errors_total",
            Help: "Total number of node execution errors",
        },
        []string{"tenant_id", "node_type", "error_type"},
    )
    
    // HTTP metrics
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "thaiyyal_http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"tenant_id", "method", "path", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "thaiyyal_http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"tenant_id", "method", "path"},
    )
    
    // Database metrics
    dbQueriesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "thaiyyal_db_queries_total",
            Help: "Total database queries",
        },
        []string{"tenant_id", "operation", "table"},
    )
    
    dbQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "thaiyyal_db_query_duration_seconds",
            Help:    "Database query duration",
            Buckets: prometheus.ExponentialBuckets(0.0001, 2, 15),
        },
        []string{"tenant_id", "operation", "table"},
    )
    
    // Tenant quotas
    tenantQuotaUsage = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "thaiyyal_tenant_quota_usage",
            Help: "Current quota usage for tenant resources",
        },
        []string{"tenant_id", "resource_type"},
    )
    
    tenantQuotaLimit = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "thaiyyal_tenant_quota_limit",
            Help: "Quota limit for tenant resources",
        },
        []string{"tenant_id", "resource_type"},
    )
    
    // System metrics
    goroutinesCount = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "thaiyyal_goroutines_count",
            Help: "Current number of goroutines",
        },
    )
    
    memoryUsageBytes = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "thaiyyal_memory_usage_bytes",
            Help: "Current memory usage in bytes",
        },
    )
)

// RecordWorkflowExecution records workflow execution metrics
func RecordWorkflowExecution(tenantID, workflowID, status string, duration float64, nodesExecuted int) {
    workflowExecutionsTotal.WithLabelValues(tenantID, workflowID, status).Inc()
    workflowExecutionDuration.WithLabelValues(tenantID, workflowID).Observe(duration)
    workflowNodesExecuted.WithLabelValues(tenantID, workflowID).Observe(float64(nodesExecuted))
}

// RecordNodeExecution records individual node execution
func RecordNodeExecution(tenantID, nodeType string, duration float64, err error) {
    nodeExecutionDuration.WithLabelValues(tenantID, nodeType).Observe(duration)
    
    if err != nil {
        errorType := getErrorType(err)
        nodeExecutionErrors.WithLabelValues(tenantID, nodeType, errorType).Inc()
    }
}

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(tenantID, method, path string, status int, duration float64) {
    httpRequestsTotal.WithLabelValues(tenantID, method, path, fmt.Sprintf("%d", status)).Inc()
    httpRequestDuration.WithLabelValues(tenantID, method, path).Observe(duration)
}

// UpdateTenantQuota updates tenant quota metrics
func UpdateTenantQuota(tenantID string, resourceType string, usage, limit int) {
    tenantQuotaUsage.WithLabelValues(tenantID, resourceType).Set(float64(usage))
    tenantQuotaLimit.WithLabelValues(tenantID, resourceType).Set(float64(limit))
}

// StartMetricsCollector starts background metrics collection
func StartMetricsCollector(ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            collectSystemMetrics()
        case <-ctx.Done():
            return
        }
    }
}

func collectSystemMetrics() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    goroutinesCount.Set(float64(runtime.NumGoroutine()))
    memoryUsageBytes.Set(float64(m.Alloc))
}
```

#### Metrics HTTP Handler

```go
package api

import (
    "net/http"
    
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupMetricsEndpoint configures Prometheus metrics endpoint
func SetupMetricsEndpoint(mux *http.ServeMux) {
    // Metrics endpoint (no authentication for monitoring)
    mux.Handle("/metrics", promhttp.Handler())
}
```

### Structured Logging

#### Logger Implementation

```go
package logger

import (
    "context"
    "encoding/json"
    "io"
    "os"
    "time"
)

type Level string

const (
    LevelDebug Level = "debug"
    LevelInfo  Level = "info"
    LevelWarn  Level = "warn"
    LevelError Level = "error"
)

type Logger struct {
    level  Level
    output io.Writer
}

type LogEntry struct {
    Timestamp string                 `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    TenantID  string                 `json:"tenant_id,omitempty"`
    UserID    string                 `json:"user_id,omitempty"`
    Fields    map[string]interface{} `json:"fields,omitempty"`
    Error     string                 `json:"error,omitempty"`
}

func New(level Level) *Logger {
    return &Logger{
        level:  level,
        output: os.Stdout,
    }
}

func (l *Logger) log(ctx context.Context, level Level, message string, fields map[string]interface{}, err error) {
    // Check if level is enabled
    if !l.isLevelEnabled(level) {
        return
    }
    
    entry := LogEntry{
        Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
        Level:     string(level),
        Message:   message,
        Fields:    fields,
    }
    
    // Add tenant context if available
    if tenantID := middleware.GetTenantID(ctx); tenantID != "" {
        entry.TenantID = tenantID
    }
    
    if userID := middleware.GetUserID(ctx); userID != "" {
        entry.UserID = userID
    }
    
    if err != nil {
        entry.Error = err.Error()
    }
    
    // Write JSON log
    encoder := json.NewEncoder(l.output)
    encoder.Encode(entry)
}

func (l *Logger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
    l.log(ctx, LevelDebug, message, fields, nil)
}

func (l *Logger) Info(ctx context.Context, message string, fields map[string]interface{}) {
    l.log(ctx, LevelInfo, message, fields, nil)
}

func (l *Logger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
    l.log(ctx, LevelWarn, message, fields, nil)
}

func (l *Logger) Error(ctx context.Context, message string, err error, fields map[string]interface{}) {
    l.log(ctx, LevelError, message, fields, err)
}

func (l *Logger) isLevelEnabled(level Level) bool {
    levels := map[Level]int{
        LevelDebug: 0,
        LevelInfo:  1,
        LevelWarn:  2,
        LevelError: 3,
    }
    return levels[level] >= levels[l.level]
}

// Example usage
func (e *Engine) Execute(ctx context.Context, workflow *Workflow) (*Result, error) {
    log.Info(ctx, "Starting workflow execution", map[string]interface{}{
        "workflow_id": workflow.ID,
        "node_count":  len(workflow.Nodes),
    })
    
    startTime := time.Now()
    result, err := e.executeWorkflow(ctx, workflow)
    duration := time.Since(startTime)
    
    if err != nil {
        log.Error(ctx, "Workflow execution failed", err, map[string]interface{}{
            "workflow_id": workflow.ID,
            "duration_ms": duration.Milliseconds(),
        })
        return nil, err
    }
    
    log.Info(ctx, "Workflow execution completed", map[string]interface{}{
        "workflow_id": workflow.ID,
        "duration_ms": duration.Milliseconds(),
        "nodes_executed": result.NodesExecuted,
    })
    
    return result, nil
}
```

### Distributed Tracing

#### OpenTelemetry Integration

```go
package tracing

import (
    "context"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
    "go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// InitTracer initializes OpenTelemetry with Jaeger
func InitTracer(serviceName, jaegerEndpoint string) error {
    // Create Jaeger exporter
    exporter, err := jaeger.New(
        jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)),
    )
    if err != nil {
        return err
    }
    
    // Create resource
    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            semconv.ServiceNameKey.String(serviceName),
        ),
    )
    if err != nil {
        return err
    }
    
    // Create trace provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
    )
    
    otel.SetTracerProvider(tp)
    tracer = tp.Tracer("thaiyyal")
    
    return nil
}

// TraceWorkflowExecution creates a span for workflow execution
func TraceWorkflowExecution(ctx context.Context, workflowID string) (context.Context, trace.Span) {
    ctx, span := tracer.Start(ctx, "workflow.execute",
        trace.WithAttributes(
            attribute.String("workflow.id", workflowID),
            attribute.String("tenant.id", middleware.GetTenantID(ctx)),
        ),
    )
    return ctx, span
}

// TraceNodeExecution creates a span for node execution
func TraceNodeExecution(ctx context.Context, nodeID, nodeType string) (context.Context, trace.Span) {
    ctx, span := tracer.Start(ctx, "node.execute",
        trace.WithAttributes(
            attribute.String("node.id", nodeID),
            attribute.String("node.type", nodeType),
        ),
    )
    return ctx, span
}

// Example usage
func (e *Engine) Execute(ctx context.Context, workflow *Workflow) (*Result, error) {
    ctx, span := tracing.TraceWorkflowExecution(ctx, workflow.ID)
    defer span.End()
    
    for _, node := range workflow.Nodes {
        ctx, nodeSpan := tracing.TraceNodeExecution(ctx, node.ID, node.Type)
        
        result, err := e.executeNode(ctx, node)
        if err != nil {
            nodeSpan.RecordError(err)
            nodeSpan.SetStatus(codes.Error, err.Error())
        }
        
        nodeSpan.End()
    }
    
    return result, nil
}
```

## Docker Compose Observability Stack

```yaml
# docker-compose.observability.yml
version: '3.8'

services:
  # Prometheus - Metrics collection
  prometheus:
    image: prom/prometheus:latest
    container_name: thaiyyal-prometheus
    volumes:
      - ./observability/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=30d'
    ports:
      - "9090:9090"
    networks:
      - thaiyyal-network
    restart: unless-stopped

  # Grafana - Visualization
  grafana:
    image: grafana/grafana:latest
    container_name: thaiyyal-grafana
    volumes:
      - ./observability/grafana/provisioning:/etc/grafana/provisioning
      - ./observability/grafana/dashboards:/var/lib/grafana/dashboards
      - grafana-data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    ports:
      - "3001:3000"
    networks:
      - thaiyyal-network
    depends_on:
      - prometheus
      - loki
    restart: unless-stopped

  # Loki - Log aggregation
  loki:
    image: grafana/loki:latest
    container_name: thaiyyal-loki
    volumes:
      - ./observability/loki/loki-config.yml:/etc/loki/local-config.yaml
      - loki-data:/loki
    ports:
      - "3100:3100"
    command: -config.file=/etc/loki/local-config.yaml
    networks:
      - thaiyyal-network
    restart: unless-stopped

  # Promtail - Log shipper
  promtail:
    image: grafana/promtail:latest
    container_name: thaiyyal-promtail
    volumes:
      - ./observability/promtail/promtail-config.yml:/etc/promtail/config.yml
      - /var/log:/var/log
      - ./data/logs:/app/logs
    command: -config.file=/etc/promtail/config.yml
    networks:
      - thaiyyal-network
    depends_on:
      - loki
    restart: unless-stopped

  # Jaeger - Distributed tracing
  jaeger:
    image: jaegertracing/all-in-one:latest
    container_name: thaiyyal-jaeger
    environment:
      - COLLECTOR_ZIPKIN_HOST_PORT=:9411
    ports:
      - "5775:5775/udp"
      - "6831:6831/udp"
      - "6832:6832/udp"
      - "5778:5778"
      - "16686:16686"  # Jaeger UI
      - "14268:14268"
      - "14250:14250"
      - "9411:9411"
    networks:
      - thaiyyal-network
    restart: unless-stopped

volumes:
  prometheus-data:
  grafana-data:
  loki-data:

networks:
  thaiyyal-network:
    driver: bridge
```

## Prometheus Configuration

```yaml
# observability/prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'thaiyyal-local'

scrape_configs:
  # Thaiyyal application metrics
  - job_name: 'thaiyyal'
    static_configs:
      - targets: ['host.docker.internal:8080']
    metrics_path: '/metrics'
    
  # Prometheus self-monitoring
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

# Alerting rules
rule_files:
  - 'alerts.yml'

# Alertmanager configuration (optional)
alerting:
  alertmanagers:
    - static_configs:
        - targets: []
```

## Alerting Rules

```yaml
# observability/prometheus/alerts.yml
groups:
  - name: thaiyyal_alerts
    interval: 30s
    rules:
      # High error rate
      - alert: HighWorkflowErrorRate
        expr: |
          sum(rate(thaiyyal_workflow_executions_total{status="failed"}[5m])) by (tenant_id)
          /
          sum(rate(thaiyyal_workflow_executions_total[5m])) by (tenant_id)
          > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High workflow error rate for tenant {{ $labels.tenant_id }}"
          description: "Error rate is {{ $value | humanizePercentage }}"
      
      # Quota exceeded
      - alert: TenantQuotaNearLimit
        expr: |
          thaiyyal_tenant_quota_usage / thaiyyal_tenant_quota_limit > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Tenant {{ $labels.tenant_id }} near quota limit"
          description: "Resource {{ $labels.resource_type }} at {{ $value | humanizePercentage }}"
      
      # Slow workflow execution
      - alert: SlowWorkflowExecution
        expr: |
          histogram_quantile(0.95,
            rate(thaiyyal_workflow_execution_duration_seconds_bucket[5m])
          ) > 30
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Slow workflow execution detected"
          description: "95th percentile execution time is {{ $value }}s"
      
      # High memory usage
      - alert: HighMemoryUsage
        expr: thaiyyal_memory_usage_bytes > 1e9  # 1GB
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High memory usage detected"
          description: "Memory usage is {{ $value | humanize }}B"
```

## Grafana Dashboards

### Workflow Execution Dashboard (JSON)

```json
{
  "dashboard": {
    "title": "Thaiyyal Workflow Execution",
    "panels": [
      {
        "title": "Workflow Executions (Rate)",
        "targets": [
          {
            "expr": "sum(rate(thaiyyal_workflow_executions_total[5m])) by (status)"
          }
        ]
      },
      {
        "title": "Execution Duration (p95)",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(thaiyyal_workflow_execution_duration_seconds_bucket[5m]))"
          }
        ]
      },
      {
        "title": "Per-Tenant Execution Count",
        "targets": [
          {
            "expr": "sum(rate(thaiyyal_workflow_executions_total[5m])) by (tenant_id)"
          }
        ]
      },
      {
        "title": "Node Execution Errors",
        "targets": [
          {
            "expr": "sum(rate(thaiyyal_node_execution_errors_total[5m])) by (node_type)"
          }
        ]
      }
    ]
  }
}
```

## Health Checks

```go
package health

import (
    "context"
    "database/sql"
    "encoding/json"
    "net/http"
    "time"
)

type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp string            `json:"timestamp"`
    Checks    map[string]Check  `json:"checks"`
}

type Check struct {
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
}

type HealthChecker struct {
    db *sql.DB
}

func (h *HealthChecker) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    status := HealthStatus{
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Checks:    make(map[string]Check),
    }
    
    // Check database
    if err := h.checkDatabase(ctx); err != nil {
        status.Checks["database"] = Check{
            Status:  "unhealthy",
            Message: err.Error(),
        }
    } else {
        status.Checks["database"] = Check{Status: "healthy"}
    }
    
    // Determine overall status
    allHealthy := true
    for _, check := range status.Checks {
        if check.Status != "healthy" {
            allHealthy = false
            break
        }
    }
    
    if allHealthy {
        status.Status = "healthy"
        w.WriteHeader(http.StatusOK)
    } else {
        status.Status = "unhealthy"
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

func (h *HealthChecker) checkDatabase(ctx context.Context) error {
    return h.db.PingContext(ctx)
}
```

## Agent Collaboration Points

### With Multi-Tenancy Agent
- Per-tenant metrics and dashboards
- Tenant quota monitoring
- Tenant-specific log aggregation
- Tenant resource usage tracking

### With Security Agent
- Audit log monitoring
- Failed authentication attempts tracking
- Security event alerting
- Anomaly detection

### With Performance Agent
- Performance metrics correlation
- Bottleneck identification
- Resource usage optimization
- Query performance tracking

### With DevOps Agent
- CI/CD pipeline metrics
- Deployment tracking
- Infrastructure monitoring
- Alert routing configuration

---

**Version**: 1.0  
**Last Updated**: 2025-10-30  
**Maintained By**: Observability Team  
**Review Cycle**: Quarterly
