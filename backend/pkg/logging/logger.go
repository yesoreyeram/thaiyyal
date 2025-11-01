// Package logging provides structured logging with context propagation for the workflow engine.
// It uses zerolog for high-performance, zero-allocation structured logging.
package logging

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// contextKey is used for context keys to avoid collisions
type contextKey string

const (
	// ContextKeyLogger is the context key for the logger instance
	ContextKeyLogger contextKey = "logger"
)

// Logger wraps zerolog.Logger with workflow-specific functionality
type Logger struct {
	logger zerolog.Logger
}

// Config holds logging configuration
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

// DefaultConfig returns default logging configuration
func DefaultConfig() Config {
	return Config{
		Level:            "info",
		Output:           os.Stdout,
		Pretty:           false,
		IncludeTimestamp: true,
		IncludeCaller:    false,
	}
}

// New creates a new logger with the given configuration
func New(cfg Config) *Logger {
	// Set global log level
	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	// Configure output
	output := cfg.Output
	if output == nil {
		output = os.Stdout
	}

	// Use console writer for pretty output, JSON otherwise
	var writer io.Writer
	if cfg.Pretty {
		writer = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}
	} else {
		writer = output
	}

	// Create logger context
	loggerCtx := zerolog.New(writer)

	// Add timestamp if configured
	if cfg.IncludeTimestamp {
		loggerCtx = loggerCtx.With().Timestamp().Logger()
	}

	// Add caller info if configured (expensive, only for debugging)
	if cfg.IncludeCaller {
		loggerCtx = loggerCtx.With().Caller().Logger()
	}

	return &Logger{
		logger: loggerCtx,
	}
}

// parseLevel converts string level to zerolog.Level
func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
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
		logger: l.logger.With().Str("workflow_id", workflowID).Logger(),
	}
}

// WithExecutionID adds execution_id to the logger context
func (l *Logger) WithExecutionID(executionID string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("execution_id", executionID).Logger(),
	}
}

// WithNodeID adds node_id to the logger context
func (l *Logger) WithNodeID(nodeID string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("node_id", nodeID).Logger(),
	}
}

// WithNodeType adds node_type to the logger context
func (l *Logger) WithNodeType(nodeType types.NodeType) *Logger {
	return &Logger{
		logger: l.logger.With().Str("node_type", string(nodeType)).Logger(),
	}
}

// WithField adds a custom field to the logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: l.logger.With().Interface(key, value).Logger(),
	}
}

// WithFields adds multiple custom fields to the logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}

// WithError adds error to the logger context
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		logger: l.logger.With().Err(err).Logger(),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// Event returns a zerolog.Event for advanced logging
func (l *Logger) Event(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

// GetZerologLogger returns the underlying zerolog.Logger for advanced use cases
func (l *Logger) GetZerologLogger() zerolog.Logger {
	return l.logger
}
