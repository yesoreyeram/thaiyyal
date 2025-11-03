package telemetry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

const (
	// Service name for telemetry
	serviceName = "thaiyyal-workflow-engine"
	
	// Metric names
	metricWorkflowExecutions     = "workflow.executions.total"
	metricWorkflowDuration       = "workflow.execution.duration"
	metricWorkflowSuccess        = "workflow.executions.success.total"
	metricWorkflowFailure        = "workflow.executions.failure.total"
	metricNodeExecutions         = "node.executions.total"
	metricNodeDuration           = "node.execution.duration"
	metricNodeSuccess            = "node.executions.success.total"
	metricNodeFailure            = "node.executions.failure.total"
	metricHTTPCalls              = "http.calls.total"
	metricHTTPDuration           = "http.call.duration"
)

// Provider manages OpenTelemetry setup and provides access to tracers and meters.
type Provider struct {
	meterProvider  *sdkmetric.MeterProvider
	tracerProvider trace.TracerProvider
	meter          metric.Meter
	tracer         trace.Tracer
	
	// Metrics instruments
	workflowExecutions metric.Int64Counter
	workflowDuration   metric.Float64Histogram
	workflowSuccess    metric.Int64Counter
	workflowFailure    metric.Int64Counter
	nodeExecutions     metric.Int64Counter
	nodeDuration       metric.Float64Histogram
	nodeSuccess        metric.Int64Counter
	nodeFailure        metric.Int64Counter
	httpCalls          metric.Int64Counter
	httpDuration       metric.Float64Histogram
	
	mu sync.RWMutex
}

// Config holds telemetry configuration
type Config struct {
	// ServiceName is the name of the service for telemetry
	ServiceName string
	
	// ServiceVersion is the version of the service
	ServiceVersion string
	
	// Environment (e.g., "production", "staging", "development")
	Environment string
	
	// EnableTracing enables distributed tracing
	EnableTracing bool
	
	// EnableMetrics enables metrics collection
	EnableMetrics bool
}

// DefaultConfig returns default telemetry configuration
func DefaultConfig() Config {
	return Config{
		ServiceName:    serviceName,
		ServiceVersion: "0.1.0",
		Environment:    "development",
		EnableTracing:  true,
		EnableMetrics:  true,
	}
}

// NewProvider creates a new telemetry provider with Prometheus metrics exporter.
// It initializes OpenTelemetry with the given configuration and returns a provider
// that can be used to create tracers and record metrics.
func NewProvider(ctx context.Context, config Config) (*Provider, error) {
	provider := &Provider{}
	
	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			attribute.String("environment", config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	
	// Initialize metrics if enabled
	if config.EnableMetrics {
		if err := provider.initMetrics(res); err != nil {
			return nil, fmt.Errorf("failed to initialize metrics: %w", err)
		}
	}
	
	// Initialize tracing if enabled
	if config.EnableTracing {
		provider.initTracing()
	}
	
	return provider, nil
}

// initMetrics initializes the metrics provider with Prometheus exporter
func (p *Provider) initMetrics(res *resource.Resource) error {
	// Create Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}
	
	// Create meter provider with the exporter
	p.meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)
	
	// Set as global meter provider
	otel.SetMeterProvider(p.meterProvider)
	
	// Create meter
	p.meter = p.meterProvider.Meter(serviceName)
	
	// Create metric instruments
	if err := p.createMetricInstruments(); err != nil {
		return fmt.Errorf("failed to create metric instruments: %w", err)
	}
	
	return nil
}

// initTracing initializes the tracing provider
func (p *Provider) initTracing() {
	// For now, use the global tracer provider
	// In production, this should be configured with appropriate exporters (OTLP, Jaeger, etc.)
	p.tracerProvider = otel.GetTracerProvider()
	p.tracer = p.tracerProvider.Tracer(serviceName)
}

// createMetricInstruments creates all metric instruments
func (p *Provider) createMetricInstruments() error {
	var err error
	
	// Workflow metrics
	p.workflowExecutions, err = p.meter.Int64Counter(
		metricWorkflowExecutions,
		metric.WithDescription("Total number of workflow executions"),
	)
	if err != nil {
		return err
	}
	
	p.workflowDuration, err = p.meter.Float64Histogram(
		metricWorkflowDuration,
		metric.WithDescription("Workflow execution duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}
	
	p.workflowSuccess, err = p.meter.Int64Counter(
		metricWorkflowSuccess,
		metric.WithDescription("Total number of successful workflow executions"),
	)
	if err != nil {
		return err
	}
	
	p.workflowFailure, err = p.meter.Int64Counter(
		metricWorkflowFailure,
		metric.WithDescription("Total number of failed workflow executions"),
	)
	if err != nil {
		return err
	}
	
	// Node metrics
	p.nodeExecutions, err = p.meter.Int64Counter(
		metricNodeExecutions,
		metric.WithDescription("Total number of node executions"),
	)
	if err != nil {
		return err
	}
	
	p.nodeDuration, err = p.meter.Float64Histogram(
		metricNodeDuration,
		metric.WithDescription("Node execution duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}
	
	p.nodeSuccess, err = p.meter.Int64Counter(
		metricNodeSuccess,
		metric.WithDescription("Total number of successful node executions"),
	)
	if err != nil {
		return err
	}
	
	p.nodeFailure, err = p.meter.Int64Counter(
		metricNodeFailure,
		metric.WithDescription("Total number of failed node executions"),
	)
	if err != nil {
		return err
	}
	
	// HTTP metrics
	p.httpCalls, err = p.meter.Int64Counter(
		metricHTTPCalls,
		metric.WithDescription("Total number of HTTP calls"),
	)
	if err != nil {
		return err
	}
	
	p.httpDuration, err = p.meter.Float64Histogram(
		metricHTTPDuration,
		metric.WithDescription("HTTP call duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}
	
	return nil
}

// Tracer returns the tracer for creating spans
func (p *Provider) Tracer() trace.Tracer {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.tracer
}

// Meter returns the meter for recording metrics
func (p *Provider) Meter() metric.Meter {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.meter
}

// RecordWorkflowExecution records metrics for a workflow execution
func (p *Provider) RecordWorkflowExecution(ctx context.Context, workflowID string, duration time.Duration, success bool, nodesExecuted int) {
	if p.meter == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("workflow.id", workflowID),
		attribute.Int("nodes.executed", nodesExecuted),
	}
	
	// Record execution count
	p.workflowExecutions.Add(ctx, 1, metric.WithAttributes(attrs...))
	
	// Record duration
	p.workflowDuration.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
	
	// Record success/failure
	if success {
		p.workflowSuccess.Add(ctx, 1, metric.WithAttributes(attrs...))
	} else {
		p.workflowFailure.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// RecordNodeExecution records metrics for a node execution
func (p *Provider) RecordNodeExecution(ctx context.Context, nodeID string, nodeType types.NodeType, duration time.Duration, success bool) {
	if p.meter == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("node.id", nodeID),
		attribute.String("node.type", string(nodeType)),
	}
	
	// Record execution count
	p.nodeExecutions.Add(ctx, 1, metric.WithAttributes(attrs...))
	
	// Record duration
	p.nodeDuration.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
	
	// Record success/failure
	if success {
		p.nodeSuccess.Add(ctx, 1, metric.WithAttributes(attrs...))
	} else {
		p.nodeFailure.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// RecordHTTPCall records metrics for an HTTP call
func (p *Provider) RecordHTTPCall(ctx context.Context, method, url string, statusCode int, duration time.Duration) {
	if p.meter == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.url", url),
		attribute.Int("http.status_code", statusCode),
	}
	
	// Record HTTP call count
	p.httpCalls.Add(ctx, 1, metric.WithAttributes(attrs...))
	
	// Record duration
	p.httpDuration.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
}

// Shutdown gracefully shuts down the telemetry provider
func (p *Provider) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.meterProvider != nil {
		if err := p.meterProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown meter provider: %w", err)
		}
	}
	
	return nil
}
