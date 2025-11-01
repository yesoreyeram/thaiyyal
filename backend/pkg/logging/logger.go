// Package logging provides structured logging with context propagation for the workflow engine.
// It uses Go's built-in slog package for high-performance structured logging.
package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// contextKey is used for context keys to avoid collisions
type contextKey string

const (
	// ContextKeyLogger is the context key for the logger instance
	ContextKeyLogger contextKey = "logger"
)

// Logger wraps slog.Logger with workflow-specific functionality
type Logger struct {
	logger *slog.Logger
}

// Config holds logging configuration
type Config struct {
	// Level is the minimum log level (debug, info, warn, error)
	Level string
	// Output is where logs are written (default: os.Stdout)
	Output io.Writer
	// Pretty enables human-readable text output (default: false for JSON)
	Pretty bool
	// IncludeCaller includes source location in logs (default: false)
	IncludeCaller bool
}

// DefaultConfig returns default logging configuration
func DefaultConfig() Config {
	return Config{
		Level:         "info",
		Output:        os.Stdout,
		Pretty:        false,
		IncludeCaller: false,
	}
}

// New creates a new logger with the given configuration
func New(cfg Config) *Logger {
	// Parse log level
	level := parseLevel(cfg.Level)

	// Configure output
	output := cfg.Output
	if output == nil {
		output = os.Stdout
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.IncludeCaller,
	}

	// Create appropriate handler
	var handler slog.Handler
	if cfg.Pretty {
		handler = slog.NewTextHandler(output, opts)
	} else {
		handler = slog.NewJSONHandler(output, opts)
	}

	return &Logger{
		logger: slog.New(handler),
	}
}

// parseLevel converts string level to slog.Level
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithContext adds the logger to a context
func (l *Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, l)
}

// FromContext retrieves the logger from context, or returns default logger
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(ContextKeyLogger).(*Logger); ok {
		return logger
	}
	// Return default logger if not in context
	return New(DefaultConfig())
}

// WithWorkflowID adds workflow_id to the logger context
func (l *Logger) WithWorkflowID(workflowID string) *Logger {
	return &Logger{
		logger: l.logger.With(slog.String("workflow_id", workflowID)),
	}
}

// WithExecutionID adds execution_id to the logger context
func (l *Logger) WithExecutionID(executionID string) *Logger {
	return &Logger{
		logger: l.logger.With(slog.String("execution_id", executionID)),
	}
}

// WithNodeID adds node_id to the logger context
func (l *Logger) WithNodeID(nodeID string) *Logger {
	return &Logger{
		logger: l.logger.With(slog.String("node_id", nodeID)),
	}
}

// WithNodeType adds node_type to the logger context
func (l *Logger) WithNodeType(nodeType types.NodeType) *Logger {
	return &Logger{
		logger: l.logger.With(slog.String("node_type", string(nodeType))),
	}
}

// WithField adds a custom field to the logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: l.logger.With(slog.Any(key, value)),
	}
}

// WithFields adds multiple custom fields to the logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, slog.Any(k, v))
	}
	return &Logger{
		logger: l.logger.With(args...),
	}
}

// WithError adds error to the logger context
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		logger: l.logger.With(slog.Any("error", err)),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.logger.Error(msg)
	os.Exit(1)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

// GetSlogLogger returns the underlying slog.Logger for advanced use cases
func (l *Logger) GetSlogLogger() *slog.Logger {
	return l.logger
}
