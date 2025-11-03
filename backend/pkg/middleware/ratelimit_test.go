package middleware

import (
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestTokenBucket_Allow tests basic token bucket functionality
func TestTokenBucket_Allow(t *testing.T) {
	tb := NewTokenBucket(10, 10) // 10 tokens/sec, capacity 10

	// Should allow first 10 requests immediately
	for i := 0; i < 10; i++ {
		if !tb.Allow("test") {
			t.Errorf("request %d should be allowed", i)
		}
	}

	// 11th request should be denied (bucket empty)
	if tb.Allow("test") {
		t.Error("request 11 should be denied (bucket empty)")
	}
}

// TestTokenBucket_Refill tests token refill over time
func TestTokenBucket_Refill(t *testing.T) {
	tb := NewTokenBucket(10, 10) // 10 tokens/sec

	// Drain the bucket
	for i := 0; i < 10; i++ {
		tb.Allow("test")
	}

	// Should be denied immediately
	if tb.Allow("test") {
		t.Error("should be denied immediately after draining")
	}

	// Wait for 0.2 seconds (should refill ~2 tokens)
	time.Sleep(200 * time.Millisecond)

	// Should allow 2 more requests
	if !tb.Allow("test") {
		t.Error("should allow request after refill (1)")
	}
	if !tb.Allow("test") {
		t.Error("should allow request after refill (2)")
	}

	// Should deny 3rd request
	if tb.Allow("test") {
		t.Error("should deny 3rd request after partial refill")
	}
}

// TestTokenBucket_Reset tests bucket reset
func TestTokenBucket_Reset(t *testing.T) {
	tb := NewTokenBucket(10, 10)

	// Drain the bucket
	for i := 0; i < 10; i++ {
		tb.Allow("test")
	}

	// Should be denied
	if tb.Allow("test") {
		t.Error("should be denied after draining")
	}

	// Reset
	tb.Reset()

	// Should allow requests again
	if !tb.Allow("test") {
		t.Error("should allow request after reset")
	}
}

// TestRateLimitMiddleware_GlobalLimit tests global rate limiting
func TestRateLimitMiddleware_GlobalLimit(t *testing.T) {
	config := RateLimitConfig{
		GlobalRPS:    5,
		EnableGlobal: true,
	}

	m := NewRateLimitMiddlewareWithConfig(config)

	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	executionCount := 0

	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		executionCount++
		return "ok", nil
	}

	// Should allow first 5 requests
	for i := 0; i < 5; i++ {
		result, err := m.Process(nil, node, handler)
		if err != nil {
			t.Errorf("request %d should be allowed: %v", i, err)
		}
		if result != "ok" {
			t.Errorf("expected 'ok', got %v", result)
		}
	}

	if executionCount != 5 {
		t.Errorf("expected 5 executions, got %d", executionCount)
	}

	// 6th request should be denied
	_, err := m.Process(nil, node, handler)
	if err == nil {
		t.Error("request 6 should be denied (global limit)")
	}

	if m.GetRejectedCount() != 1 {
		t.Errorf("expected 1 rejected request, got %d", m.GetRejectedCount())
	}

	// Handler should not have been called
	if executionCount != 5 {
		t.Errorf("handler should not be called when rate limited, got %d executions", executionCount)
	}
}

// TestRateLimitMiddleware_NodeTypeLimit tests per-node-type rate limiting
func TestRateLimitMiddleware_NodeTypeLimit(t *testing.T) {
	config := RateLimitConfig{
		EnablePerNodeType: true,
		NodeTypeRPS: map[types.NodeType]float64{
			types.NodeTypeHTTP: 3,
		},
	}

	m := NewRateLimitMiddlewareWithConfig(config)

	httpNode := types.Node{ID: "http1", Type: types.NodeTypeHTTP}
	numberNode := types.Node{ID: "num1", Type: types.NodeTypeNumber}

	executionCount := 0
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		executionCount++
		return "ok", nil
	}

	// Should allow 3 HTTP requests
	for i := 0; i < 3; i++ {
		_, err := m.Process(nil, httpNode, handler)
		if err != nil {
			t.Errorf("HTTP request %d should be allowed: %v", i, err)
		}
	}

	// 4th HTTP request should be denied
	_, err := m.Process(nil, httpNode, handler)
	if err == nil {
		t.Error("HTTP request 4 should be denied (node type limit)")
	}

	// Number node should still be allowed (no limit set)
	_, err = m.Process(nil, numberNode, handler)
	if err != nil {
		t.Errorf("Number node should be allowed: %v", err)
	}

	if executionCount != 4 {
		t.Errorf("expected 4 successful executions, got %d", executionCount)
	}
}

// TestRateLimitMiddleware_DisabledLimits tests middleware with all limits disabled
func TestRateLimitMiddleware_DisabledLimits(t *testing.T) {
	config := RateLimitConfig{
		EnableGlobal:      false,
		EnablePerNodeType: false,
		EnablePerWorkflow: false,
	}

	m := NewRateLimitMiddlewareWithConfig(config)

	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	executionCount := 0

	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		executionCount++
		return "ok", nil
	}

	// Should allow unlimited requests
	for i := 0; i < 100; i++ {
		_, err := m.Process(nil, node, handler)
		if err != nil {
			t.Errorf("request %d should be allowed (no limits): %v", i, err)
		}
	}

	if executionCount != 100 {
		t.Errorf("expected 100 executions, got %d", executionCount)
	}

	if m.GetRejectedCount() != 0 {
		t.Errorf("expected 0 rejected requests, got %d", m.GetRejectedCount())
	}
}

// TestRateLimitMiddleware_DefaultConfig tests default configuration
func TestRateLimitMiddleware_DefaultConfig(t *testing.T) {
	m := NewRateLimitMiddleware()

	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "ok", nil
	}

	// Default config should allow up to 100 requests
	for i := 0; i < 100; i++ {
		_, err := m.Process(nil, node, handler)
		if err != nil {
			t.Errorf("request %d should be allowed with default config: %v", i, err)
		}
	}

	// 101st should be denied
	_, err := m.Process(nil, node, handler)
	if err == nil {
		t.Error("request 101 should be denied (default global limit)")
	}
}

// TestRateLimitMiddleware_ConcurrentAccess tests thread safety
func TestRateLimitMiddleware_ConcurrentAccess(t *testing.T) {
	config := RateLimitConfig{
		GlobalRPS:    50,
		EnableGlobal: true,
	}

	m := NewRateLimitMiddlewareWithConfig(config)

	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "ok", nil
	}

	// Run concurrent requests
	concurrency := 100
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer func() { done <- true }()
			m.Process(nil, node, handler)
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrency; i++ {
		<-done
	}

	// Should have rejected some requests (we sent 100, limit is 50)
	rejectedCount := m.GetRejectedCount()
	if rejectedCount < 40 {
		t.Errorf("expected significant rejections with concurrent access, got %d", rejectedCount)
	}
}

// TestRateLimitMiddleware_Name tests the Name method
func TestRateLimitMiddleware_Name(t *testing.T) {
	m := NewRateLimitMiddleware()

	if m.Name() != "RateLimit" {
		t.Errorf("expected 'RateLimit', got %s", m.Name())
	}
}

// BenchmarkRateLimitMiddleware_GlobalLimit benchmarks global rate limiting
func BenchmarkRateLimitMiddleware_GlobalLimit(b *testing.B) {
	config := RateLimitConfig{
		GlobalRPS:    1000000, // High limit to avoid rate limiting during benchmark
		EnableGlobal: true,
	}

	m := NewRateLimitMiddlewareWithConfig(config)

	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "ok", nil
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m.Process(nil, node, handler)
	}
}

// BenchmarkTokenBucket_Allow benchmarks token bucket algorithm
func BenchmarkTokenBucket_Allow(b *testing.B) {
	tb := NewTokenBucket(1000000, 1000000) // High limit to avoid rate limiting

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tb.Allow("test")
	}
}
