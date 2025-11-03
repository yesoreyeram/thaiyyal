package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestMaxNodeExecutions tests that the maximum node execution limit is enforced
func TestMaxNodeExecutions(t *testing.T) {
	tests := []struct {
		name          string
		maxExecutions int
		nodeCount     int
		shouldFail    bool
		expectedCount int
	}{
		{
			name:          "under limit",
			maxExecutions: 10,
			nodeCount:     5,
			shouldFail:    false,
			expectedCount: 5,
		},
		{
			name:          "at limit",
			maxExecutions: 5,
			nodeCount:     5,
			shouldFail:    false,
			expectedCount: 5,
		},
		{
			name:          "exceed limit",
			maxExecutions: 3,
			nodeCount:     5,
			shouldFail:    true,
			expectedCount: 4, // Should fail after 4th execution
		},
		{
			name:          "unlimited executions",
			maxExecutions: 0, // 0 means unlimited
			nodeCount:     10,
			shouldFail:    false,
			expectedCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple linear workflow
			payload := createLinearWorkflow(tt.nodeCount)

			// Create engine with custom config
			config := types.DefaultConfig()
			config.MaxNodeExecutions = tt.maxExecutions

			engine, err := NewWithConfig([]byte(payload), config)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			// Execute workflow
			result, err := engine.Execute()

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if !strings.Contains(err.Error(), "maximum node executions exceeded") {
					t.Errorf("Expected 'maximum node executions exceeded' error, got: %v", err)
				}
				// Check that we stopped at the right count
				if engine.GetNodeExecutionCount() != tt.expectedCount {
					t.Errorf("Expected %d node executions, got %d", tt.expectedCount, engine.GetNodeExecutionCount())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if result == nil {
					t.Errorf("Expected result but got nil")
				}
				// Verify execution count
				if engine.GetNodeExecutionCount() != tt.expectedCount {
					t.Errorf("Expected %d node executions, got %d", tt.expectedCount, engine.GetNodeExecutionCount())
				}
			}
		})
	}
}

// TestMaxHTTPCallsPerExecution tests that the maximum HTTP calls limit is enforced
func TestMaxHTTPCallsPerExecution(t *testing.T) {
	// Create a test HTTP server
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Response %d", requestCount)))
	}))
	defer server.Close()

	tests := []struct {
		name          string
		maxHTTPCalls  int
		httpNodeCount int
		shouldFail    bool
	}{
		{
			name:          "under limit",
			maxHTTPCalls:  10,
			httpNodeCount: 5,
			shouldFail:    false,
		},
		{
			name:          "at limit",
			maxHTTPCalls:  5,
			httpNodeCount: 5,
			shouldFail:    false,
		},
		{
			name:          "exceed limit",
			maxHTTPCalls:  3,
			httpNodeCount: 5,
			shouldFail:    true,
		},
		{
			name:          "unlimited",
			maxHTTPCalls:  0, // 0 means unlimited
			httpNodeCount: 10,
			shouldFail:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCount = 0 // Reset for each test

			// Create workflow with multiple HTTP nodes
			payload := createHTTPWorkflow(server.URL, tt.httpNodeCount)

			// Create engine with custom config
			config := types.DefaultConfig()
			config.AllowHTTP = true       // Enable HTTP for this test
			config.BlockLocalhost = false // Allow localhost for test server
			config.MaxHTTPCallsPerExec = tt.maxHTTPCalls

			engine, err := NewWithConfig([]byte(payload), config)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			// Execute workflow
			result, err := engine.Execute()

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if !strings.Contains(err.Error(), "maximum HTTP calls per execution exceeded") {
					t.Errorf("Expected 'maximum HTTP calls per execution exceeded' error, got: %v", err)
				}
				// Verify we stopped after the limit
				expectedCalls := tt.maxHTTPCalls + 1
				if engine.GetHTTPCallCount() != expectedCalls {
					t.Errorf("Expected %d HTTP calls, got %d", expectedCalls, engine.GetHTTPCallCount())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if result == nil {
					t.Errorf("Expected result but got nil")
				}
				// Verify all HTTP calls were made
				if engine.GetHTTPCallCount() != tt.httpNodeCount {
					t.Errorf("Expected %d HTTP calls, got %d", tt.httpNodeCount, engine.GetHTTPCallCount())
				}
			}
		})
	}
}

// TestHTTPCallsInWhileLoop tests that multiple HTTP calls in a single workflow are counted correctly
func TestHTTPCallsInChain(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	tests := []struct {
		name          string
		maxHTTPCalls  int
		httpCallCount int
		shouldFail    bool
		expectedCalls int
	}{
		{
			name:          "HTTP chain under limit",
			maxHTTPCalls:  10,
			httpCallCount: 3,
			shouldFail:    false,
			expectedCalls: 3,
		},
		{
			name:          "HTTP chain exceed limit",
			maxHTTPCalls:  5,
			httpCallCount: 10,
			shouldFail:    true,
			expectedCalls: 6, // Should fail after 6th call
		},
		{
			name:          "HTTP chain unlimited",
			maxHTTPCalls:  0,
			httpCallCount: 20,
			shouldFail:    false,
			expectedCalls: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create workflow with chained HTTP nodes
			payload := createHTTPWorkflow(server.URL, tt.httpCallCount)

			// Create engine with custom config - ENABLE HTTP explicitly
			config := types.DefaultConfig()
			config.AllowHTTP = true       // Enable HTTP for this test
			config.BlockLocalhost = false // Allow localhost for test server
			config.MaxHTTPCallsPerExec = tt.maxHTTPCalls

			engine, err := NewWithConfig([]byte(payload), config)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			// Execute workflow
			_, err = engine.Execute()

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if !strings.Contains(err.Error(), "maximum HTTP calls per execution exceeded") {
					t.Errorf("Expected 'maximum HTTP calls per execution exceeded' error, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}

			// Verify HTTP call count
			if engine.GetHTTPCallCount() != tt.expectedCalls {
				t.Errorf("Expected %d HTTP calls, got %d", tt.expectedCalls, engine.GetHTTPCallCount())
			}
		})
	}
}

// TestMultipleProtectionLimits tests that multiple limits work together
func TestMultipleProtectionLimits(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Create a workflow with multiple HTTP nodes
	payload := createHTTPWorkflow(server.URL, 10)

	config := types.DefaultConfig()
	config.AllowHTTP = true       // Enable HTTP for this test
	config.BlockLocalhost = false // Allow localhost for test server
	config.MaxNodeExecutions = 50
	config.MaxHTTPCallsPerExec = 5

	engine, err := NewWithConfig([]byte(payload), config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Execute workflow - should fail on HTTP limit first (after 5 HTTP nodes)
	_, err = engine.Execute()
	if err == nil {
		t.Errorf("Expected error but got none")
	}

	// Verify it was HTTP limit that was exceeded
	if !strings.Contains(err.Error(), "maximum HTTP calls per execution exceeded") {
		t.Errorf("Expected HTTP limit error, got: %v", err)
	}
}

// TestDefaultProtectionLimits verifies default protection limits are reasonable
func TestDefaultProtectionLimits(t *testing.T) {
	config := types.DefaultConfig()

	if config.MaxNodeExecutions <= 0 {
		t.Errorf("Default MaxNodeExecutions should be positive, got %d", config.MaxNodeExecutions)
	}
	if config.MaxHTTPCallsPerExec <= 0 {
		t.Errorf("Default MaxHTTPCallsPerExec should be positive, got %d", config.MaxHTTPCallsPerExec)
	}

	// Verify reasonable defaults
	if config.MaxNodeExecutions < 1000 {
		t.Errorf("Default MaxNodeExecutions seems too low: %d", config.MaxNodeExecutions)
	}
	if config.MaxHTTPCallsPerExec < 10 {
		t.Errorf("Default MaxHTTPCallsPerExec seems too low: %d", config.MaxHTTPCallsPerExec)
	}
}

// TestValidationLimitsProtection verifies validation limits are stricter
func TestValidationLimitsProtection(t *testing.T) {
	config := types.ValidationLimits()
	defaultConfig := types.DefaultConfig()

	if config.MaxNodeExecutions >= defaultConfig.MaxNodeExecutions {
		t.Errorf("Validation MaxNodeExecutions should be stricter than default")
	}
	if config.MaxHTTPCallsPerExec >= defaultConfig.MaxHTTPCallsPerExec {
		t.Errorf("Validation MaxHTTPCallsPerExec should be stricter than default")
	}
}

// TestDevelopmentConfigProtection verifies development config is more permissive
func TestDevelopmentConfigProtection(t *testing.T) {
	config := types.DevelopmentConfig()
	defaultConfig := types.DefaultConfig()

	if config.MaxNodeExecutions <= defaultConfig.MaxNodeExecutions {
		t.Errorf("Development MaxNodeExecutions should be more permissive than default")
	}
	if config.MaxHTTPCallsPerExec <= defaultConfig.MaxHTTPCallsPerExec {
		t.Errorf("Development MaxHTTPCallsPerExec should be more permissive than default")
	}
}

// TestWorkflowTimeoutStillWorks verifies existing timeout protection still works
func TestWorkflowTimeoutStillWorks(t *testing.T) {
	// Create a workflow with a delay node that will timeout
	payload := `{
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 1}},
			{"id": "2", "type": "delay", "data": {"duration": "10s"}},
			{"id": "3", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"}
		]
	}`

	config := types.DefaultConfig()
	config.MaxExecutionTime = 100 * time.Millisecond // Very short timeout

	engine, err := NewWithConfig([]byte(payload), config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	_, err = engine.Execute()
	if err == nil {
		t.Errorf("Expected timeout error but got none")
	}
	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

// Helper functions to create test workflows

func createLinearWorkflow(nodeCount int) string {
	nodes := []map[string]interface{}{}
	edges := []map[string]interface{}{}

	// Create a linear chain of nodes
	for i := 0; i < nodeCount; i++ {
		nodes = append(nodes, map[string]interface{}{
			"id":   fmt.Sprintf("node%d", i),
			"type": "number",
			"data": map[string]interface{}{"value": float64(i + 1)},
		})

		if i > 0 {
			edges = append(edges, map[string]interface{}{
				"id":     fmt.Sprintf("e%d", i),
				"source": fmt.Sprintf("node%d", i-1),
				"target": fmt.Sprintf("node%d", i),
			})
		}
	}

	payload := map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	}

	jsonData, _ := json.Marshal(payload)
	return string(jsonData)
}

func createHTTPWorkflow(url string, httpNodeCount int) string {
	nodes := []map[string]interface{}{}
	edges := []map[string]interface{}{}

	// Create a chain of HTTP nodes
	for i := 0; i < httpNodeCount; i++ {
		nodes = append(nodes, map[string]interface{}{
			"id":   fmt.Sprintf("http%d", i),
			"type": "http",
			"data": map[string]interface{}{"url": url},
		})

		if i > 0 {
			edges = append(edges, map[string]interface{}{
				"id":     fmt.Sprintf("e%d", i),
				"source": fmt.Sprintf("http%d", i-1),
				"target": fmt.Sprintf("http%d", i),
			})
		}
	}

	payload := map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	}

	jsonData, _ := json.Marshal(payload)
	return string(jsonData)
}
