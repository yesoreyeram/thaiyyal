package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestNewProvider(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "default config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "custom config",
			config: Config{
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				Environment:    "test",
				EnableTracing:  true,
				EnableMetrics:  true,
			},
			wantErr: false,
		},
		{
			name: "metrics only",
			config: Config{
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				Environment:    "test",
				EnableTracing:  false,
				EnableMetrics:  true,
			},
			wantErr: false,
		},
		{
			name: "tracing only",
			config: Config{
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				Environment:    "test",
				EnableTracing:  true,
				EnableMetrics:  false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewProvider(ctx, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if provider == nil {
					t.Error("NewProvider() returned nil provider")
					return
				}

				// Verify tracer
				if tt.config.EnableTracing && provider.Tracer() == nil {
					t.Error("Tracer() returned nil when tracing is enabled")
				}

				// Verify meter
				if tt.config.EnableMetrics && provider.Meter() == nil {
					t.Error("Meter() returned nil when metrics are enabled")
				}

				// Clean up
				if err := provider.Shutdown(ctx); err != nil {
					t.Errorf("Shutdown() error = %v", err)
				}
			}
		})
	}
}

func TestRecordWorkflowExecution(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()

	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}
	defer provider.Shutdown(ctx)

	tests := []struct {
		name          string
		workflowID    string
		duration      time.Duration
		success       bool
		nodesExecuted int
	}{
		{
			name:          "successful workflow",
			workflowID:    "wf-123",
			duration:      100 * time.Millisecond,
			success:       true,
			nodesExecuted: 5,
		},
		{
			name:          "failed workflow",
			workflowID:    "wf-456",
			duration:      50 * time.Millisecond,
			success:       false,
			nodesExecuted: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			provider.RecordWorkflowExecution(ctx, tt.workflowID, tt.duration, tt.success, tt.nodesExecuted)
		})
	}
}

func TestRecordNodeExecution(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()

	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}
	defer provider.Shutdown(ctx)

	tests := []struct {
		name     string
		nodeID   string
		nodeType types.NodeType
		duration time.Duration
		success  bool
	}{
		{
			name:     "successful number node",
			nodeID:   "node-1",
			nodeType: types.NodeTypeNumber,
			duration: 10 * time.Millisecond,
			success:  true,
		},
		{
			name:     "failed operation node",
			nodeID:   "node-2",
			nodeType: types.NodeTypeOperation,
			duration: 5 * time.Millisecond,
			success:  false,
		},
		{
			name:     "successful http node",
			nodeID:   "node-3",
			nodeType: types.NodeTypeHTTP,
			duration: 200 * time.Millisecond,
			success:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			provider.RecordNodeExecution(ctx, tt.nodeID, tt.nodeType, tt.duration, tt.success)
		})
	}
}

func TestRecordHTTPCall(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()

	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}
	defer provider.Shutdown(ctx)

	tests := []struct {
		name       string
		method     string
		url        string
		statusCode int
		duration   time.Duration
	}{
		{
			name:       "successful GET",
			method:     "GET",
			url:        "https://api.example.com/data",
			statusCode: 200,
			duration:   150 * time.Millisecond,
		},
		{
			name:       "failed POST",
			method:     "POST",
			url:        "https://api.example.com/submit",
			statusCode: 500,
			duration:   100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			provider.RecordHTTPCall(ctx, tt.method, tt.url, tt.statusCode, tt.duration)
		})
	}
}

func TestShutdown(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()

	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}

	// First shutdown should succeed
	if err := provider.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}

	// Second shutdown should handle already shut down state gracefully
	// Note: The underlying SDK may return an error when shutting down twice
	// This is expected behavior and we just verify it doesn't panic
	_ = provider.Shutdown(ctx)
}

func TestProviderWithNilMetrics(t *testing.T) {
	ctx := context.Background()

	// Create provider with metrics disabled
	config := Config{
		ServiceName:    "test",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		EnableTracing:  true,
		EnableMetrics:  false,
	}

	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}
	defer provider.Shutdown(ctx)

	// These should not panic even with nil metrics
	provider.RecordWorkflowExecution(ctx, "test", time.Second, true, 1)
	provider.RecordNodeExecution(ctx, "node1", types.NodeTypeNumber, time.Millisecond, true)
	provider.RecordHTTPCall(ctx, "GET", "http://example.com", 200, time.Second)
}
