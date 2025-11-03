package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/observer"
)

// TelemetryObserver implements observer.Observer and records telemetry data
// for workflow execution events.
type TelemetryObserver struct {
	provider *Provider
	
	// Track active spans for workflow and nodes
	workflowSpan trace.Span
	nodeSpans    map[string]trace.Span
	
	// Track execution times
	workflowStartTime time.Time
	nodeStartTimes    map[string]time.Time
}

// NewTelemetryObserver creates a new telemetry observer
func NewTelemetryObserver(provider *Provider) *TelemetryObserver {
	return &TelemetryObserver{
		provider:       provider,
		nodeSpans:      make(map[string]trace.Span),
		nodeStartTimes: make(map[string]time.Time),
	}
}

// OnEvent handles execution events and records telemetry data
func (o *TelemetryObserver) OnEvent(ctx context.Context, event observer.Event) {
	switch event.Type {
	case observer.EventWorkflowStart:
		o.handleWorkflowStart(ctx, event)
	case observer.EventWorkflowEnd:
		o.handleWorkflowEnd(ctx, event)
	case observer.EventNodeStart:
		o.handleNodeStart(ctx, event)
	case observer.EventNodeSuccess:
		o.handleNodeSuccess(ctx, event)
	case observer.EventNodeFailure:
		o.handleNodeFailure(ctx, event)
	}
}

func (o *TelemetryObserver) handleWorkflowStart(ctx context.Context, event observer.Event) {
	// Start workflow span
	// Note: spanCtx can be used for context propagation if needed in the future
	_, span := o.provider.Tracer().Start(ctx, "workflow.execute",
		trace.WithAttributes(
			attribute.String("workflow.id", event.WorkflowID),
			attribute.String("execution.id", event.ExecutionID),
		),
	)
	
	o.workflowSpan = span
	o.workflowStartTime = event.Timestamp
}

func (o *TelemetryObserver) handleWorkflowEnd(ctx context.Context, event observer.Event) {
	// Calculate duration
	duration := time.Since(o.workflowStartTime)
	
	// Get nodes executed count from metadata
	nodesExecuted := 0
	if val, ok := event.Metadata["nodes_executed"]; ok {
		if count, ok := val.(int); ok {
			nodesExecuted = count
		}
	}
	
	// Record metrics
	success := event.Status == observer.StatusSuccess
	o.provider.RecordWorkflowExecution(ctx, event.WorkflowID, duration, success, nodesExecuted)
	
	// End workflow span
	if o.workflowSpan != nil {
		if event.Error != nil {
			o.workflowSpan.RecordError(event.Error)
			o.workflowSpan.SetStatus(codes.Error, event.Error.Error())
		} else {
			o.workflowSpan.SetStatus(codes.Ok, "workflow completed successfully")
		}
		o.workflowSpan.End()
	}
}

func (o *TelemetryObserver) handleNodeStart(ctx context.Context, event observer.Event) {
	// Start node span as child of workflow span
	var spanCtx context.Context
	if o.workflowSpan != nil {
		spanCtx = trace.ContextWithSpan(ctx, o.workflowSpan)
	} else {
		spanCtx = ctx
	}
	
	_, span := o.provider.Tracer().Start(spanCtx, "node.execute",
		trace.WithAttributes(
			attribute.String("node.id", event.NodeID),
			attribute.String("node.type", string(event.NodeType)),
			attribute.String("execution.id", event.ExecutionID),
		),
	)
	
	o.nodeSpans[event.NodeID] = span
	o.nodeStartTimes[event.NodeID] = event.Timestamp
}

func (o *TelemetryObserver) handleNodeSuccess(ctx context.Context, event observer.Event) {
	o.handleNodeEnd(ctx, event, true)
}

func (o *TelemetryObserver) handleNodeFailure(ctx context.Context, event observer.Event) {
	o.handleNodeEnd(ctx, event, false)
}

func (o *TelemetryObserver) handleNodeEnd(ctx context.Context, event observer.Event, success bool) {
	// Calculate duration
	var duration time.Duration
	if startTime, ok := o.nodeStartTimes[event.NodeID]; ok {
		duration = time.Since(startTime)
		delete(o.nodeStartTimes, event.NodeID)
	}
	
	// Record metrics
	o.provider.RecordNodeExecution(ctx, event.NodeID, event.NodeType, duration, success)
	
	// End node span
	if span, ok := o.nodeSpans[event.NodeID]; ok {
		if event.Error != nil {
			span.RecordError(event.Error)
			span.SetStatus(codes.Error, event.Error.Error())
		} else {
			span.SetStatus(codes.Ok, "node completed successfully")
		}
		span.End()
		delete(o.nodeSpans, event.NodeID)
	}
}
