// Package logging provides structured logging capabilities for the workflow engine.
//
// # Overview
//
// The logging package implements a structured logging system with support for
// multiple output formats, log levels, contextual information, and integration
// with the workflow execution lifecycle.
//
// # Features
//
//   - Structured logging: JSON and text formats
//   - Log levels: DEBUG, INFO, WARN, ERROR
//   - Context propagation: Execution ID, workflow ID, node ID
//   - Conditional logging: Enable/disable per package or level
//   - Performance: Minimal overhead for disabled log levels
//   - Thread-safe: Safe for concurrent use
//   - Flexible output: Write to any io.Writer
//
// # Log Levels
//
// The package supports standard log levels:
//
//   - DEBUG: Detailed diagnostic information
//   - INFO: General informational messages
//   - WARN: Warning messages for potential issues
//   - ERROR: Error messages for failures
//
// # Basic Usage
//
//	import "github.com/yesoreyeram/thaiyyal/backend/pkg/logging"
//
//	// Create logger
//	logger := logging.New(logging.Config{
//	    Level:  logging.LevelInfo,
//	    Format: logging.FormatJSON,
//	    Output: os.Stdout,
//	})
//
//	// Log messages
//	logger.Info("Workflow started", map[string]interface{}{
//	    "workflow_id": "wf-123",
//	    "node_count": 42,
//	})
//
//	logger.Error("Execution failed", map[string]interface{}{
//	    "error": err.Error(),
//	    "node_id": "node-5",
//	})
//
// # Context Integration
//
// The logger integrates with Go contexts for automatic field extraction:
//
//	// Logger extracts execution_id and workflow_id from context
//	logger.WithContext(ctx).Info("Node executing", map[string]interface{}{
//	    "node_type": "http",
//	})
//
// # Structured Fields
//
// All log entries support structured fields:
//
//	logger.Info("HTTP request completed", map[string]interface{}{
//	    "method": "GET",
//	    "url": "https://api.example.com",
//	    "status": 200,
//	    "duration_ms": 145,
//	})
//
// # Output Formats
//
// JSON Format (production):
//
//	{
//	  "timestamp": "2024-01-15T10:30:00Z",
//	  "level": "INFO",
//	  "message": "Workflow started",
//	  "workflow_id": "wf-123",
//	  "execution_id": "exec-456"
//	}
//
// Text Format (development):
//
//	2024-01-15T10:30:00Z INFO Workflow started workflow_id=wf-123 execution_id=exec-456
//
// # Configuration
//
// Logger configuration options:
//
//	config := logging.Config{
//	    Level:      logging.LevelDebug,    // Minimum level to log
//	    Format:     logging.FormatJSON,    // Output format
//	    Output:     os.Stdout,             // Where to write logs
//	    AddSource:  true,                  // Include file:line
//	    TimeFormat: time.RFC3339Nano,      // Timestamp format
//	}
//
// # Performance Considerations
//
//   - Zero allocation for disabled log levels
//   - Lazy field evaluation
//   - Buffered output for high throughput
//   - Minimal lock contention
//
// # Common Logging Patterns
//
// Workflow execution:
//
//	logger.Info("Workflow execution started", map[string]interface{}{
//	    "workflow_id": workflow.ID,
//	    "node_count": len(workflow.Nodes),
//	})
//
// Node execution:
//
//	logger.Debug("Node executing", map[string]interface{}{
//	    "node_id": node.ID,
//	    "node_type": node.Type,
//	    "inputs": inputs,
//	})
//
// Error logging:
//
//	logger.Error("Node execution failed", map[string]interface{}{
//	    "node_id": node.ID,
//	    "error": err.Error(),
//	    "retry_count": retries,
//	})
//
// Performance metrics:
//
//	logger.Info("Workflow completed", map[string]interface{}{
//	    "workflow_id": workflow.ID,
//	    "duration_ms": elapsed.Milliseconds(),
//	    "nodes_executed": count,
//	})
//
// # Integration with Observability
//
// The logging package integrates with the observer package to automatically
// log workflow execution events:
//
//	observer := logging.NewLoggingObserver(logger)
//	engine.RegisterObserver(observer)
//
// # Best Practices
//
//   - Use structured fields instead of string formatting
//   - Include execution context (workflow_id, node_id, etc.)
//   - Log at appropriate levels (avoid debug in production)
//   - Add timing information for performance analysis
//   - Include error context (not just error message)
//   - Use consistent field names across the codebase
//
// # Thread Safety
//
// All logger operations are thread-safe and can be used concurrently
// from multiple goroutines without additional synchronization.
//
// # Testing
//
// For testing, use a logger with a buffer:
//
//	buf := &bytes.Buffer{}
//	logger := logging.New(logging.Config{
//	    Output: buf,
//	    Format: logging.FormatJSON,
//	})
//
//	// Execute code
//	// Verify log output
//	assert.Contains(t, buf.String(), "expected message")
package logging
