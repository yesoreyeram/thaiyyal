// Package observer provides the Observer pattern implementation for workflow execution monitoring.
// This allows library consumers to track and monitor workflow execution behavior.
package observer

import (
	"context"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// EventType represents the type of execution event
type EventType string

const (
	// Workflow-level events
	EventWorkflowStart EventType = "workflow_start"
	EventWorkflowEnd   EventType = "workflow_end"

	// Node-level events
	EventNodeStart   EventType = "node_start"
	EventNodeEnd     EventType = "node_end"
	EventNodeSuccess EventType = "node_success"
	EventNodeFailure EventType = "node_failure"
)

// ExecutionStatus represents the status of a node or workflow execution
type ExecutionStatus string

const (
	StatusStarted   ExecutionStatus = "started"
	StatusSuccess   ExecutionStatus = "success"
	StatusFailure   ExecutionStatus = "failure"
	StatusCompleted ExecutionStatus = "completed"
)

// Event represents an execution event with all relevant metadata
type Event struct {
	// Event identification
	Type      EventType       `json:"type"`
	Status    ExecutionStatus `json:"status"`
	Timestamp time.Time       `json:"timestamp"`

	// Execution context
	ExecutionID string `json:"execution_id"`
	WorkflowID  string `json:"workflow_id,omitempty"`

	// Node-specific data (empty for workflow-level events)
	NodeID   string         `json:"node_id,omitempty"`
	NodeType types.NodeType `json:"node_type,omitempty"`

	// Timing information
	StartTime   time.Time     `json:"start_time,omitempty"`
	ElapsedTime time.Duration `json:"elapsed_time,omitempty"`

	// Execution results
	Result interface{} `json:"result,omitempty"`
	Error  error       `json:"error,omitempty"`

	// Additional metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Observer defines the interface for workflow execution observers.
// Observers receive notifications about various stages of workflow execution.
type Observer interface {
	// OnEvent is called when an execution event occurs.
	// The context can be used for cancellation and passing request-scoped values.
	OnEvent(ctx context.Context, event Event)
}

// Logger defines the interface for custom logging.
// This allows library consumers to integrate with their own logging systems.
type Logger interface {
	// Debug logs debug-level messages
	Debug(msg string, fields map[string]interface{})

	// Info logs info-level messages
	Info(msg string, fields map[string]interface{})

	// Warn logs warning-level messages
	Warn(msg string, fields map[string]interface{})

	// Error logs error-level messages
	Error(msg string, fields map[string]interface{})
}
