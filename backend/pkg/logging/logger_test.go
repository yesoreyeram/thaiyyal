package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name:   "default config",
			config: DefaultConfig(),
		},
		{
			name: "debug level",
			config: Config{
				Level:  "debug",
				Output: &bytes.Buffer{},
				Pretty: false,
			},
		},
		{
			name: "pretty output",
			config: Config{
				Level:  "info",
				Output: &bytes.Buffer{},
				Pretty: true,
			},
		},
		{
			name: "with caller",
			config: Config{
				Level:         "info",
				Output:        &bytes.Buffer{},
				Pretty:        false,
				IncludeCaller: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.config)
			if logger == nil {
				t.Error("Expected logger to be created, got nil")
			}
		})
	}
}

func TestLogger_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("Expected log to contain level INFO, got: %s", output)
	}
}

func TestLogger_Debug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "debug",
		Output: buf,
		Pretty: false,
	})

	logger.Debug("debug message")

	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected log to contain 'debug message', got: %s", output)
	}
	if !strings.Contains(output, `"level":"DEBUG"`) {
		t.Errorf("Expected log to contain level DEBUG, got: %s", output)
	}
}

func TestLogger_DebugNotLogged(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info", // Debug should not be logged
		Output: buf,
		Pretty: false,
	})

	logger.Debug("debug message")

	output := buf.String()
	if output != "" {
		t.Errorf("Expected no log output for debug when level is info, got: %s", output)
	}
}

func TestLogger_Warn(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "warn",
		Output: buf,
		Pretty: false,
	})

	logger.Warn("warning message")

	output := buf.String()
	if !strings.Contains(output, "warning message") {
		t.Errorf("Expected log to contain 'warning message', got: %s", output)
	}
	if !strings.Contains(output, `"level":"WARN"`) {
		t.Errorf("Expected log to contain level WARN, got: %s", output)
	}
}

func TestLogger_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "error",
		Output: buf,
		Pretty: false,
	})

	logger.Error("error message")

	output := buf.String()
	if !strings.Contains(output, "error message") {
		t.Errorf("Expected log to contain 'error message', got: %s", output)
	}
	if !strings.Contains(output, `"level":"ERROR"`) {
		t.Errorf("Expected log to contain level ERROR, got: %s", output)
	}
}

func TestLogger_WithWorkflowID(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger = logger.WithWorkflowID("workflow-123")
	logger.Info("test")

	output := buf.String()
	if !strings.Contains(output, `"workflow_id":"workflow-123"`) {
		t.Errorf("Expected log to contain workflow_id, got: %s", output)
	}
}

func TestLogger_WithExecutionID(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger = logger.WithExecutionID("exec-456")
	logger.Info("test")

	output := buf.String()
	if !strings.Contains(output, `"execution_id":"exec-456"`) {
		t.Errorf("Expected log to contain execution_id, got: %s", output)
	}
}

func TestLogger_WithNodeID(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger = logger.WithNodeID("node-789")
	logger.Info("test")

	output := buf.String()
	if !strings.Contains(output, `"node_id":"node-789"`) {
		t.Errorf("Expected log to contain node_id, got: %s", output)
	}
}

func TestLogger_WithNodeType(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger = logger.WithNodeType(types.NodeTypeHTTP)
	logger.Info("test")

	output := buf.String()
	if !strings.Contains(output, `"node_type":"http"`) {
		t.Errorf("Expected log to contain node_type, got: %s", output)
	}
}

func TestLogger_WithField(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger = logger.WithField("custom_field", "custom_value")
	logger.Info("test")

	output := buf.String()
	if !strings.Contains(output, `"custom_field":"custom_value"`) {
		t.Errorf("Expected log to contain custom_field, got: %s", output)
	}
}

func TestLogger_WithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger = logger.WithFields(map[string]interface{}{
		"field1": "value1",
		"field2": 42,
	})
	logger.Info("test")

	output := buf.String()
	if !strings.Contains(output, `"field1":"value1"`) {
		t.Errorf("Expected log to contain field1, got: %s", output)
	}
	if !strings.Contains(output, `"field2":42`) {
		t.Errorf("Expected log to contain field2, got: %s", output)
	}
}

func TestLogger_WithError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "error",
		Output: buf,
		Pretty: false,
	})

	err := &testError{"test error"}
	logger = logger.WithError(err)
	logger.Error("error occurred")

	output := buf.String()
	if !strings.Contains(output, "test error") {
		t.Errorf("Expected log to contain error message, got: %s", output)
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestLogger_ChainedContext(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger = logger.
		WithWorkflowID("wf-123").
		WithExecutionID("exec-456").
		WithNodeID("node-789").
		WithNodeType(types.NodeTypeHTTP)

	logger.Info("test")

	output := buf.String()

	// Parse JSON to verify all fields
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse log JSON: %v", err)
	}

	expectedFields := map[string]string{
		"workflow_id":  "wf-123",
		"execution_id": "exec-456",
		"node_id":      "node-789",
		"node_type":    "http",
		"level":        "INFO",
		"msg":          "test",
	}

	for key, expectedValue := range expectedFields {
		if value, ok := logEntry[key]; !ok {
			t.Errorf("Expected field %s in log, got: %v", key, logEntry)
		} else if value != expectedValue {
			t.Errorf("Expected %s=%s, got %s=%v", key, expectedValue, key, value)
		}
	}
}

func TestLogger_WithContext(t *testing.T) {
	logger := New(DefaultConfig())
	ctx := context.Background()

	// Add logger to context
	ctx = logger.WithContext(ctx)

	// Retrieve logger from context
	retrieved := FromContext(ctx)
	if retrieved == nil {
		t.Error("Expected logger from context, got nil")
	}
}

func TestLogger_FromContextDefault(t *testing.T) {
	ctx := context.Background()

	// Should return default logger when not in context
	logger := FromContext(ctx)
	if logger == nil {
		t.Error("Expected default logger, got nil")
	}
}

func TestLogger_Infof(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger.Infof("formatted message: %s %d", "test", 42)

	output := buf.String()
	if !strings.Contains(output, "formatted message: test 42") {
		t.Errorf("Expected formatted message, got: %s", output)
	}
}

func TestLogger_Debugf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "debug",
		Output: buf,
		Pretty: false,
	})

	logger.Debugf("debug: %d", 123)

	output := buf.String()
	if !strings.Contains(output, "debug: 123") {
		t.Errorf("Expected formatted debug message, got: %s", output)
	}
}

func TestLogger_Warnf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "warn",
		Output: buf,
		Pretty: false,
	})

	logger.Warnf("warning: %s", "test")

	output := buf.String()
	if !strings.Contains(output, "warning: test") {
		t.Errorf("Expected formatted warning message, got: %s", output)
	}
}

func TestLogger_Errorf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "error",
		Output: buf,
		Pretty: false,
	})

	logger.Errorf("error: %d", 500)

	output := buf.String()
	if !strings.Contains(output, "error: 500") {
		t.Errorf("Expected formatted error message, got: %s", output)
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"debug", "DEBUG"},
		{"info", "INFO"},
		{"warn", "WARN"},
		{"warning", "WARN"},
		{"error", "ERROR"},
		{"invalid", "INFO"}, // Should default to info
		{"", "INFO"},        // Should default to info
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			level := parseLevel(tt.input)
			if level.String() != tt.expected {
				t.Errorf("parseLevel(%s) = %s, want %s", tt.input, level.String(), tt.expected)
			}
		})
	}
}

func TestLogger_JSONOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  "info",
		Output: buf,
		Pretty: false,
	})

	logger.Info("test message")

	// Verify output is valid JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Errorf("Log output is not valid JSON: %v", err)
	}
}
