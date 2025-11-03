package executor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestHTTPExecutor_ConnectionPooling tests that connections are reused
func TestHTTPExecutor_ConnectionPooling(t *testing.T) {
	requestCount := 0
	mu := sync.Mutex{}

	// Create test server that tracks request count
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	executor := NewHTTPExecutor()
	config := types.Config{
		HTTPTimeout:       30 * time.Second,
		MaxHTTPRedirects:  10,
		MaxResponseSize:   10 * 1024 * 1024,
		AllowHTTP:         true, // Enable HTTP for this test
	}
	ctx := &mockExecutionContext{config: config}

	url := server.URL
	node := types.Node{
		ID:   "1",
		Type: types.NodeTypeHTTP,
		Data: types.NodeData{
			URL: &url,
		},
	}

	// Make multiple requests - they should reuse the same HTTP client
	for i := 0; i < 5; i++ {
		result, err := executor.Execute(ctx, node)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		if result != "OK" {
			t.Errorf("Expected 'OK', got %v", result)
		}
	}

	// Verify all requests were made
	mu.Lock()
	defer mu.Unlock()
	if requestCount != 5 {
		t.Errorf("Expected 5 requests, got %d", requestCount)
	}

	// Verify that the client is cached (same instance)
	if executor.client == nil {
		t.Error("Expected client to be cached")
	}
}

// TestHTTPExecutor_ConcurrentRequests tests thread-safe concurrent requests
func TestHTTPExecutor_ConcurrentRequests(t *testing.T) {
	requestCount := 0
	mu := sync.Mutex{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()
		// Simulate some processing time
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	executor := NewHTTPExecutor()
	config := types.Config{
		HTTPTimeout:       30 * time.Second,
		MaxHTTPRedirects:  10,
		MaxResponseSize:   10 * 1024 * 1024,
		AllowHTTP:         true, // Enable HTTP for this test
	}
	ctx := &mockExecutionContext{config: config}

	url := server.URL
	node := types.Node{
		ID:   "1",
		Type: types.NodeTypeHTTP,
		Data: types.NodeData{
			URL: &url,
		},
	}

	// Make concurrent requests
	concurrency := 20
	var wg sync.WaitGroup
	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result, err := executor.Execute(ctx, node)
			if err != nil {
				errors <- fmt.Errorf("request %d failed: %w", id, err)
				return
			}
			if result != "OK" {
				errors <- fmt.Errorf("request %d: expected 'OK', got %v", id, result)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Error(err)
	}

	// Verify all requests were made
	mu.Lock()
	defer mu.Unlock()
	if requestCount != concurrency {
		t.Errorf("Expected %d requests, got %d", concurrency, requestCount)
	}
}

// TestHTTPExecutor_ClientReuse tests that the same client is reused
func TestHTTPExecutor_ClientReuse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	executor := NewHTTPExecutor()
	config := types.Config{
		HTTPTimeout:       30 * time.Second,
		MaxHTTPRedirects:  10,
		MaxResponseSize:   10 * 1024 * 1024,
		AllowHTTP:         true, // Enable HTTP for this test
	}
	ctx := &mockExecutionContext{config: config}

	url := server.URL
	node := types.Node{
		ID:   "1",
		Type: types.NodeTypeHTTP,
		Data: types.NodeData{
			URL: &url,
		},
	}

	// First request - creates client
	_, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}

	client1 := executor.client
	if client1 == nil {
		t.Fatal("Client should be created after first request")
	}

	// Second request - should reuse client
	_, err = executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}

	client2 := executor.client
	if client1 != client2 {
		t.Error("Expected same client instance to be reused")
	}
}

// TestHTTPExecutor_MultipleHosts tests connection pooling with multiple hosts
func TestHTTPExecutor_MultipleHosts(t *testing.T) {
	// Create multiple test servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server 1"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server 2"))
	}))
	defer server2.Close()

	executor := NewHTTPExecutor()
	config := types.Config{
		HTTPTimeout:       30 * time.Second,
		MaxHTTPRedirects:  10,
		MaxResponseSize:   10 * 1024 * 1024,
		AllowHTTP:         true, // Enable HTTP for this test
	}
	ctx := &mockExecutionContext{config: config}

	// Request to server 1
	url1 := server1.URL
	node1 := types.Node{
		ID:   "1",
		Type: types.NodeTypeHTTP,
		Data: types.NodeData{URL: &url1},
	}

	result1, err := executor.Execute(ctx, node1)
	if err != nil {
		t.Fatalf("Request to server 1 failed: %v", err)
	}
	if result1 != "Server 1" {
		t.Errorf("Expected 'Server 1', got %v", result1)
	}

	// Request to server 2 - should use same client but different connection
	url2 := server2.URL
	node2 := types.Node{
		ID:   "2",
		Type: types.NodeTypeHTTP,
		Data: types.NodeData{URL: &url2},
	}

	result2, err := executor.Execute(ctx, node2)
	if err != nil {
		t.Fatalf("Request to server 2 failed: %v", err)
	}
	if result2 != "Server 2" {
		t.Errorf("Expected 'Server 2', got %v", result2)
	}

	// Both should use the same client
	if executor.client == nil {
		t.Error("Client should be created")
	}
}

// mockExecutionContext for testing
type mockExecutionContext struct {
	config types.Config
}

func (m *mockExecutionContext) GetNodeInputs(nodeID string) []interface{} {
	return nil
}

func (m *mockExecutionContext) GetNode(nodeID string) *types.Node {
	return nil
}

func (m *mockExecutionContext) GetVariable(name string) (interface{}, error) {
	return nil, nil
}

func (m *mockExecutionContext) SetVariable(name string, value interface{}) error {
	return nil
}

func (m *mockExecutionContext) GetAccumulator() interface{} {
	return nil
}

func (m *mockExecutionContext) SetAccumulator(value interface{}) {
}

func (m *mockExecutionContext) GetCounter() float64 {
	return 0
}

func (m *mockExecutionContext) SetCounter(value float64) {
}

func (m *mockExecutionContext) GetCache(key string) (interface{}, bool) {
	return nil, false
}

func (m *mockExecutionContext) SetCache(key string, value interface{}, ttl time.Duration) {
}

func (m *mockExecutionContext) GetWorkflowContext() map[string]interface{} {
	return nil
}

func (m *mockExecutionContext) GetContextVariable(name string) (interface{}, bool) {
	return nil, false
}

func (m *mockExecutionContext) SetContextVariable(name string, value interface{}) {
}

func (m *mockExecutionContext) GetContextConstant(name string) (interface{}, bool) {
	return nil, false
}

func (m *mockExecutionContext) SetContextConstant(name string, value interface{}) {
}

func (m *mockExecutionContext) InterpolateTemplate(template string) string {
	return template
}

func (m *mockExecutionContext) GetNodeResult(nodeID string) (interface{}, bool) {
	return nil, false
}

func (m *mockExecutionContext) SetNodeResult(nodeID string, result interface{}) {
}

func (m *mockExecutionContext) GetAllNodeResults() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *mockExecutionContext) GetVariables() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *mockExecutionContext) GetContextVariables() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *mockExecutionContext) GetConfig() types.Config {
	return m.config
}

func (m *mockExecutionContext) GetHTTPClientRegistry() interface{} {
	return nil
}

func (m *mockExecutionContext) IncrementNodeExecution() error {
	return nil
}

func (m *mockExecutionContext) IncrementHTTPCall() error {
	return nil
}

func (m *mockExecutionContext) GetNodeExecutionCount() int {
	return 0
}

func (m *mockExecutionContext) GetHTTPCallCount() int {
	return 0
}
