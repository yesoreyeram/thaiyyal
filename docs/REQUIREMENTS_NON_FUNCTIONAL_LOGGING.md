# Non-Functional Requirements: Logging

## Logging Requirements

### LOG-1: Format
- **REQ-LOG-1.1**: Logs SHALL use JSON format
- **REQ-LOG-1.2**: Logs SHALL include timestamp (ISO 8601)
- **REQ-LOG-1.3**: Logs SHALL include log level
- **REQ-LOG-1.4**: Logs SHALL include execution context

### LOG-2: Content
- **REQ-LOG-2.1**: Workflow execution start/end SHALL be logged
- **REQ-LOG-2.2**: Node execution SHALL be logged (with node ID, type)
- **REQ-LOG-2.3**: Errors SHALL be logged with stack traces
- **REQ-LOG-2.4**: Resource limit violations SHALL be logged

### LOG-3: Security
- **REQ-LOG-3.1**: Sensitive data SHALL be sanitized before logging
- **REQ-LOG-3.2**: URLs SHALL have credentials removed
- **REQ-LOG-3.3**: Error messages SHALL not leak internal details

### LOG-4: Performance
- **REQ-LOG-4.1**: Logging SHALL not significantly impact performance
- **REQ-LOG-4.2**: High-frequency logs SHALL use appropriate levels
- **REQ-LOG-4.3**: Structured logging SHALL use slog package

## Log Levels

- **DEBUG**: Detailed troubleshooting information
- **INFO**: General workflow execution information
- **WARN**: Warning conditions (approaching limits)
- **ERROR**: Error conditions (execution failures)

## Example Logs

```json
{"time":"2025-11-03T15:52:04Z","level":"INFO","msg":"workflow execution started","workflow_id":"wf-123","execution_id":"exec-456"}
{"time":"2025-11-03T15:52:04Z","level":"INFO","msg":"node execution completed","workflow_id":"wf-123","execution_id":"exec-456","node_id":"node-1","node_type":"number","duration_ms":0}
{"time":"2025-11-03T15:52:05Z","level":"ERROR","msg":"node execution failed","workflow_id":"wf-123","execution_id":"exec-456","node_id":"node-2","error":"division by zero"}
```

---

**Last Updated:** 2025-11-03
**Version:** 1.0
