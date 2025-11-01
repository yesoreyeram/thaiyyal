package middleware

import (
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/logging"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// LoggingMiddleware logs node execution start and completion.
// It records execution time and logs errors if execution fails.
type LoggingMiddleware struct {
	logger *logging.Logger
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger *logging.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Process logs node execution
func (m *LoggingMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	// Add node context to logger
	nodeLogger := m.logger.
		WithNodeID(node.ID).
		WithNodeType(node.Type)

	nodeLogger.Debug("node execution started")
	startTime := time.Now()

	// Execute the node
	result, err := next(ctx, node)

	duration := time.Since(startTime)

	if err != nil {
		nodeLogger.
			WithError(err).
			WithField("duration_ms", duration.Milliseconds()).
			Error("node execution failed")
	} else {
		nodeLogger.
			WithField("duration_ms", duration.Milliseconds()).
			Debug("node execution completed")
	}

	return result, err
}

// Name returns the middleware name
func (m *LoggingMiddleware) Name() string {
	return "Logging"
}
