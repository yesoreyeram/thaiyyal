# Observability Workflow Examples

This directory contains example workflows demonstrating observability patterns with Thaiyyal.

## Examples

### 1. Metrics Collection Workflow

Collects metrics from multiple endpoints and aggregates them.

**Use Case:** Monitor multiple services and aggregate their health metrics.

**File:** `metrics-collection.json`

### 2. Log Aggregation Workflow

Fetches logs from multiple sources, filters, and transforms them.

**Use Case:** Centralized log collection and analysis.

**File:** `log-aggregation.json`

### 3. Alert Enrichment Workflow

Enriches alert data with context from multiple sources.

**Use Case:** Add context to alerts before sending to notification systems.

**File:** `alert-enrichment.json`

### 4. Multi-Source Health Check

Checks health of multiple services and aggregates status.

**Use Case:** Service mesh health monitoring.

**File:** `health-check.json`

### 5. Time-Series Data Processing

Processes time-series data with filtering and aggregation.

**Use Case:** Metrics processing and downsampling.

**File:** `timeseries-processing.json`

## Running Examples

### Using the API Server

```bash
# Start the server
./backend/bin/thaiyyal-server

# Execute a workflow
curl -X POST http://localhost:8080/api/v1/workflow/execute \
  -H "Content-Type: application/json" \
  -d @examples/observability/metrics-collection.json

# Validate a workflow
curl -X POST http://localhost:8080/api/v1/workflow/validate \
  -H "Content-Type: application/json" \
  -d @examples/observability/health-check.json
```

### Using the Go Library

```go
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
)

func main() {
	// Read workflow definition
	data, err := ioutil.ReadFile("examples/observability/metrics-collection.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create engine
	eng, err := engine.New(data)
	if err != nil {
		log.Fatal(err)
	}

	// Execute workflow
	result, err := eng.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %+v\n", result)
}
```

## Best Practices

1. **Error Handling**: Always use `try_catch` nodes for HTTP calls
2. **Timeouts**: Set appropriate timeouts for external calls
3. **Retries**: Use `retry` nodes for flaky endpoints
4. **Caching**: Cache expensive operations with `cache` nodes
5. **Validation**: Validate inputs before processing
6. **Monitoring**: Use the `/metrics` endpoint to monitor workflow execution

## Metrics to Monitor

When running these workflows in production, monitor:

```promql
# Workflow execution rate
rate(workflow_executions_total{workflow_id="metrics-collection"}[5m])

# Execution duration (p95)
histogram_quantile(0.95, rate(workflow_execution_duration_bucket[5m]))

# Error rate
rate(workflow_executions_failure_total[5m]) / rate(workflow_executions_total[5m])

# Node execution time by type
histogram_quantile(0.95, rate(node_execution_duration_bucket[5m])) by (node_type)
```

## Customizing Examples

Each example can be customized by:

1. Modifying node configurations
2. Adding new nodes for additional processing
3. Changing data transformations
4. Adjusting retry and timeout policies

## Security Considerations

When using these examples in production:

1. **Enable HTTP allowlist**: Configure `AllowedDomains` in engine config
2. **Use HTTPS**: Always use HTTPS for external calls
3. **Authentication**: Add authentication tokens to HTTP headers
4. **Rate Limiting**: Implement rate limiting on workflow execution
5. **Input Validation**: Validate all external data

See [Security Best Practices](../../docs/SECURITY_BEST_PRACTICES.md) for more information.
