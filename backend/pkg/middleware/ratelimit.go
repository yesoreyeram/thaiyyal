package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// RateLimiter defines the interface for rate limiting implementations
type RateLimiter interface {
	// Allow checks if a request is allowed based on rate limits
	// Returns true if allowed, false if rate limit exceeded
	Allow(key string) bool
	
	// Reset clears all rate limit state
	Reset()
}

// RateLimitMiddleware enforces rate limits to prevent DoS attacks.
// It uses the token bucket algorithm for smooth rate limiting.
type RateLimitMiddleware struct {
	globalLimiter     RateLimiter
	nodeTypeLimiters  map[types.NodeType]RateLimiter
	workflowLimiters  map[string]RateLimiter
	mu                sync.RWMutex
	
	// Configuration
	enableGlobal      bool
	enablePerNodeType bool
	enablePerWorkflow bool
	
	// Metrics
	rejectedCount     int64
	rejectedCountMu   sync.Mutex
}

// RateLimitConfig configures rate limiting behavior
type RateLimitConfig struct {
	// Global rate limit (requests per second across all nodes)
	GlobalRPS float64
	
	// Per-node-type rate limits
	NodeTypeRPS map[types.NodeType]float64
	
	// Per-workflow rate limits (requests per second per workflow)
	WorkflowRPS float64
	
	// Enable flags
	EnableGlobal      bool
	EnablePerNodeType bool
	EnablePerWorkflow bool
}

// DefaultRateLimitConfig returns default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		GlobalRPS:         100,  // 100 requests/sec globally
		WorkflowRPS:       10,   // 10 requests/sec per workflow
		EnableGlobal:      true,
		EnablePerNodeType: false,
		EnablePerWorkflow: false,
		NodeTypeRPS:       make(map[types.NodeType]float64),
	}
}

// NewRateLimitMiddleware creates a new rate limiting middleware with default config
func NewRateLimitMiddleware() *RateLimitMiddleware {
	return NewRateLimitMiddlewareWithConfig(DefaultRateLimitConfig())
}

// NewRateLimitMiddlewareWithConfig creates a new rate limiting middleware with custom config
func NewRateLimitMiddlewareWithConfig(config RateLimitConfig) *RateLimitMiddleware {
	m := &RateLimitMiddleware{
		nodeTypeLimiters:  make(map[types.NodeType]RateLimiter),
		workflowLimiters:  make(map[string]RateLimiter),
		enableGlobal:      config.EnableGlobal,
		enablePerNodeType: config.EnablePerNodeType,
		enablePerWorkflow: config.EnablePerWorkflow,
	}
	
	// Create global limiter
	if config.EnableGlobal && config.GlobalRPS > 0 {
		m.globalLimiter = NewTokenBucket(config.GlobalRPS, int64(config.GlobalRPS))
	}
	
	// Create per-node-type limiters
	if config.EnablePerNodeType {
		for nodeType, rps := range config.NodeTypeRPS {
			if rps > 0 {
				m.nodeTypeLimiters[nodeType] = NewTokenBucket(rps, int64(rps))
			}
		}
	}
	
	return m
}

// Process enforces rate limits before node execution
func (m *RateLimitMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	// Check global rate limit
	if m.enableGlobal && m.globalLimiter != nil {
		if !m.globalLimiter.Allow("global") {
			m.incrementRejected()
			return nil, fmt.Errorf("global rate limit exceeded")
		}
	}
	
	// Check per-node-type rate limit
	if m.enablePerNodeType {
		m.mu.RLock()
		limiter, exists := m.nodeTypeLimiters[node.Type]
		m.mu.RUnlock()
		
		if exists && !limiter.Allow(string(node.Type)) {
			m.incrementRejected()
			return nil, fmt.Errorf("rate limit exceeded for node type: %s", node.Type)
		}
	}
	
	// Check per-workflow rate limit
	if m.enablePerWorkflow {
		workflowID := getWorkflowID(ctx)
		if workflowID != "" {
			limiter := m.getOrCreateWorkflowLimiter(workflowID)
			if !limiter.Allow(workflowID) {
				m.incrementRejected()
				return nil, fmt.Errorf("rate limit exceeded for workflow: %s", workflowID)
			}
		}
	}
	
	// Rate limits passed, execute node
	return next(ctx, node)
}

// Name returns the middleware name
func (m *RateLimitMiddleware) Name() string {
	return "RateLimit"
}

// GetRejectedCount returns the number of rejected requests
func (m *RateLimitMiddleware) GetRejectedCount() int64 {
	m.rejectedCountMu.Lock()
	defer m.rejectedCountMu.Unlock()
	return m.rejectedCount
}

// incrementRejected increments the rejected request counter
func (m *RateLimitMiddleware) incrementRejected() {
	m.rejectedCountMu.Lock()
	m.rejectedCount++
	m.rejectedCountMu.Unlock()
}

// getOrCreateWorkflowLimiter gets or creates a rate limiter for a workflow
func (m *RateLimitMiddleware) getOrCreateWorkflowLimiter(workflowID string) RateLimiter {
	m.mu.RLock()
	limiter, exists := m.workflowLimiters[workflowID]
	m.mu.RUnlock()
	
	if exists {
		return limiter
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Double-check after acquiring write lock
	limiter, exists = m.workflowLimiters[workflowID]
	if exists {
		return limiter
	}
	
	// Create new limiter (default: 10 RPS per workflow)
	limiter = NewTokenBucket(10, 10)
	m.workflowLimiters[workflowID] = limiter
	return limiter
}

// getWorkflowID extracts workflow ID from context (placeholder implementation)
func getWorkflowID(ctx executor.ExecutionContext) string {
	// In a real implementation, this would extract the workflow ID from the context
	// For now, return empty string to disable per-workflow limiting
	return ""
}

// TokenBucket implements the token bucket algorithm for rate limiting
type TokenBucket struct {
	rate       float64   // tokens per second
	capacity   int64     // maximum tokens
	tokens     float64   // current tokens
	lastRefill time.Time // last refill time
	mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(rate float64, capacity int64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     float64(capacity),
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed based on available tokens
func (tb *TokenBucket) Allow(key string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = min(tb.tokens+elapsed*tb.rate, float64(tb.capacity))
	tb.lastRefill = now
	
	// Check if we have at least 1 token
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	
	return false
}

// Reset clears the token bucket state
func (tb *TokenBucket) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	tb.tokens = float64(tb.capacity)
	tb.lastRefill = time.Now()
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
