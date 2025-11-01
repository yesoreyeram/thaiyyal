# Structured Logging with Context Propagation

**Status**: ✅ Complete  
**Task ID**: OBS-001  
**Priority**: P1  
**Implementation Date**: November 1, 2025

## Overview

Enterprise-grade structured logging has been implemented for the Thaiyyal workflow engine using Go's built-in `slog` package (introduced in Go 1.21), providing high-performance structured logging with zero external dependencies.

## Features

### 🎯 Core Capabilities

- **Structured JSON Logging**: All logs are emitted in JSON format for easy parsing and analysis
- **Context Propagation**: Workflow and execution context automatically flows through all log messages
- **Correlation IDs**: Each workflow execution gets a unique execution ID for end-to-end tracing
- **Built-in Package**: Uses Go's standard library `slog` for high-performance, memory-efficient logging
- **Configurable Log Levels**: Support for debug, info, warn, error levels
- **Contextual Fields**: Automatic inclusion of workflow_id, execution_id, node_id, and node_type

### 🔑 Key Fields

Every log message includes:
- `level`: Log level (DEBUG, INFO, WARN, ERROR)
- `time`: ISO 8601 timestamp
- `msg`: Human-readable message
- `workflow_id`: Workflow identifier (if available)
- `execution_id`: Unique execution identifier
- `node_id`: Node being executed (for node-level logs)
- `node_type`: Type of node being executed
- `duration_ms`: Execution duration in milliseconds (for completion logs)
- `error`: Error details (for error logs)

## Architecture

### Package Structure

```
backend/pkg/logging/
├── logger.go         # Main logger implementation
└── logger_test.go    # Comprehensive tests (22 test cases)
```

### Integration Points

The structured logger is integrated at key points in the workflow engine:

1. **Engine Initialization**: Logger created with workflow_id and execution_id
2. **Workflow Execution Start**: Log when workflow execution begins
3. **Workflow Execution End**: Log success/failure with duration and node count
4. **Node Execution**: Log each node execution with timing
5. **Validation Warnings**: Log when resource validation fails but execution continues
6. **Error Handling**: Log all errors with full context

## Usage

### Basic Usage

The logging system is automatically initialized when creating a workflow engine:

```go
import "github.com/yesoreyeram/thaiyyal/backend/pkg/engine"

// Create engine - logging is automatically set up
eng, err := engine.New(payloadJSON)
if err != nil {
    log.Fatalf("Failed to create engine: %v", err)
}

// Execute workflow - all logs will include context
result, err := eng.Execute()
```

### Custom Logger Configuration

For advanced use cases, you can create a custom logger:

```go
import "github.com/yesoreyeram/thaiyyal/backend/pkg/logging"

// Create custom logger
logger := logging.New(logging.Config{
    Level:            "debug",    // Set log level
    Output:           os.Stderr,  // Custom output
    Pretty:           true,       // Human-readable console output
    IncludeTimestamp: true,
    IncludeCaller:    false,      // Expensive, use only for debugging
})

// Add workflow context
logger = logger.
    WithWorkflowID("my-workflow").
    WithExecutionID("exec-123")

// Log messages
logger.Info("workflow started")
logger.WithField("custom_field", "value").Info("custom message")
```

### Log Levels

Configure minimum log level via `Config.Level`:

- `debug`: Detailed information for diagnosing problems
- `info`: General informational messages (default)
- `warn`: Warning messages for potentially harmful situations
- `error`: Error events that might still allow the application to continue
- `fatal`: Severe error events that will cause the application to abort
- `panic`: Very severe errors, followed by a panic

### Context Propagation

The logger supports method chaining for building rich context:

```go
logger := logging.New(logging.DefaultConfig()).
    WithWorkflowID("workflow-123").
    WithExecutionID("exec-456").
    WithNodeID("node-789").
    WithNodeType(types.NodeTypeHTTP).
    WithField("custom", "value").
    WithFields(map[string]interface{}{
        "field1": "value1",
        "field2": 42,
    })

logger.Info("executing node")
```

Output:
```json
{
  "level": "info",
  "workflow_id": "workflow-123",
  "execution_id": "exec-456",
  "node_id": "node-789",
  "node_type": "http",
  "custom": "value",
  "field1": "value1",
  "field2": 42,
  "time": "2025-11-01T14:54:29Z",
  "message": "executing node"
}
```

## Example Log Output

### Successful Workflow Execution

```json
{"time":"2025-11-01T15:50:44.325269983Z","level":"INFO","msg":"workflow execution started","workflow_id":"demo-workflow","execution_id":"ba6efab7886d6f1f"}
{"time":"2025-11-01T15:50:44.325359291Z","level":"INFO","msg":"node execution completed successfully","workflow_id":"demo-workflow","execution_id":"ba6efab7886d6f1f","node_id":"1","node_type":"number","duration_ms":0}
{"time":"2025-11-01T15:50:44.325373548Z","level":"INFO","msg":"node execution completed successfully","workflow_id":"demo-workflow","execution_id":"ba6efab7886d6f1f","node_id":"2","node_type":"number","duration_ms":0}
{"time":"2025-11-01T15:50:44.325385801Z","level":"INFO","msg":"node execution completed successfully","workflow_id":"demo-workflow","execution_id":"ba6efab7886d6f1f","node_id":"3","node_type":"operation","duration_ms":0}
{"time":"2025-11-01T15:50:44.325414124Z","level":"INFO","msg":"workflow execution completed successfully","workflow_id":"demo-workflow","execution_id":"ba6efab7886d6f1f","duration_ms":0,"nodes_executed":4}
```

### Failed Workflow Execution

```json
{"level":"info","workflow_id":"demo-workflow","execution_id":"7dae176e25650eab","time":"2025-11-01T14:54:11Z","message":"workflow execution started"}
{"level":"info","workflow_id":"demo-workflow","execution_id":"7dae176e25650eab","node_id":"1","node_type":"number","duration_ms":0,"time":"2025-11-01T14:54:11Z","message":"node execution completed successfully"}
{"level":"error","workflow_id":"demo-workflow","execution_id":"7dae176e25650eab","node_id":"3","node_type":"operation","error":"operation node missing op","time":"2025-11-01T14:54:11Z","message":"node execution failed"}
{"level":"error","workflow_id":"demo-workflow","execution_id":"7dae176e25650eab","error":"error executing node 3: operation node missing op","time":"2025-11-01T14:54:11Z","message":"workflow execution failed"}
```

## Performance

Go's `slog` package is designed for high performance:
- **Efficient**: Optimized for production use with minimal allocations
- **Fast**: Built into the Go standard library and well-optimized
- **Small Memory Footprint**: Minimal impact on application memory
- **No External Dependencies**: No need to manage third-party library versions

## Integration with Observability Stack

The structured logging implementation is designed to integrate with enterprise observability tools:

### Log Aggregation
- **Elasticsearch**: JSON logs can be directly ingested
- **Splunk**: Native JSON support
- **DataDog**: Structured log parsing
- **CloudWatch Logs Insights**: JSON query support

### Log Analysis
All logs include consistent fields for filtering:
```sql
-- Find all errors in a specific workflow
SELECT * FROM logs 
WHERE workflow_id = 'demo-workflow' 
  AND level = 'error'

-- Calculate average node execution time
SELECT node_type, AVG(duration_ms) 
FROM logs 
WHERE message = 'node execution completed successfully'
GROUP BY node_type
```

### Correlation with Traces
The `execution_id` field enables correlation between:
- Logs (this implementation)
- Distributed traces (future: OBS-002)
- Metrics (future: OBS-003)

## Testing

Comprehensive test suite with 22 test cases covering:
- Logger creation with various configurations
- All log levels (debug, info, warn, error)
- Context propagation (workflow_id, execution_id, node_id, etc.)
- Field chaining and custom fields
- Error logging
- JSON output validation
- Level filtering

Run tests:
```bash
cd backend
go test -v ./pkg/logging/
```

## Future Enhancements

### Planned (Future Tasks)
- **OBS-002**: Distributed tracing integration (OpenTelemetry)
- **OBS-003**: Metrics collection (Prometheus)
- **OBS-004**: Real-time monitoring dashboard
- **SEC-007**: Audit logging framework

### Possible Improvements
- Log sampling for high-volume scenarios
- Dynamic log level adjustment
- Log rotation and archival
- Sensitive data redaction
- Custom log formatters

## Migration Guide

### From No Logging (Previous State)
The logging system is automatically enabled. No changes required to existing code.

### Custom Integration
If you need to integrate the logger into custom node executors:

```go
import "github.com/yesoreyeram/thaiyyal/backend/pkg/logging"

type MyCustomExecutor struct {
    logger *logging.Logger
}

func (e *MyCustomExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    // Get logger from context or create new one
    logger := logging.FromContext(ctx.GetWorkflowContext()).
        WithNodeID(node.ID).
        WithNodeType(node.Type)
    
    logger.Info("custom node execution started")
    
    // Your execution logic...
    
    logger.WithField("result_size", len(result)).Info("custom node completed")
    return result, nil
}
```

## Configuration Reference

### Config Struct

```go
type Config struct {
    // Level is the minimum log level (debug, info, warn, error)
    Level string
    
    // Output is where logs are written (default: os.Stdout)
    Output io.Writer
    
    // Pretty enables human-readable console output (default: false for JSON)
    Pretty bool
    
    // IncludeTimestamp includes timestamps in logs (default: true)
    IncludeTimestamp bool
    
    // IncludeCaller includes file:line in logs (default: false, expensive)
    IncludeCaller bool
}
```

### Default Configuration

```go
logging.DefaultConfig() // Returns:
// {
//     Level: "info",
//     Output: os.Stdout,
//     Pretty: false,
//     IncludeTimestamp: true,
//     IncludeCaller: false,
// }
```

## Dependencies

**None** - Uses only Go's standard library `log/slog` package (Go 1.21+)

## Related Documentation

- [OBS-002: Distributed Tracing](OBS-002-distributed-tracing.md) (planned)
- [OBS-003: Metrics Collection](OBS-003-metrics.md) (planned)
- [SEC-007: Audit Logging](SEC-007-audit-logging.md) (planned)
- [Workflow Engine Documentation](../backend/pkg/engine/README.md)

## Support

For questions or issues related to logging:
1. Check this documentation
2. Review example code in `backend/examples/logging_demo/`
3. Review test cases in `backend/pkg/logging/logger_test.go`
4. Open an issue on GitHub

---

**Last Updated**: November 1, 2025  
**Task Status**: ✅ Complete  
**Test Coverage**: 22/22 tests passing
