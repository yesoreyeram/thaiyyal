package executor

import (
	"fmt"
	"sync"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// RateLimiterExecutor executes RateLimiter nodes
// Controls request rates to prevent overwhelming APIs
type RateLimiterExecutor struct {
	mu      sync.Mutex
	buckets map[string]*rateLimitBucket
}

// rateLimitBucket tracks requests within a time window
type rateLimitBucket struct {
	requests  []time.Time
	maxRequests int
	window    time.Duration
}

// NewRateLimiterExecutor creates a new RateLimiterExecutor
func NewRateLimiterExecutor() *RateLimiterExecutor {
	return &RateLimiterExecutor{
		buckets: make(map[string]*rateLimitBucket),
	}
}

// Execute runs the RateLimiter node
// Enforces rate limits and delays requests as needed
func (e *RateLimiterExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	var inputValue interface{}
	if len(inputs) > 0 {
		inputValue = inputs[0]
	}

	// Get configuration
	maxRequests := 10
	if node.Data.MaxRequests != nil {
		maxRequests = *node.Data.MaxRequests
	}

	perDuration := "1s"
	if node.Data.PerDuration != nil {
		perDuration = *node.Data.PerDuration
	}

	duration, err := parseDuration(perDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid per_duration format: %w", err)
	}

	strategy := "fixed_window"
	if node.Data.RateLimitStrategy != nil {
		strategy = *node.Data.RateLimitStrategy
	}

	// Only fixed_window strategy is implemented for now
	if strategy != "fixed_window" {
		return nil, fmt.Errorf("unsupported rate limit strategy: %s (only fixed_window supported)", strategy)
	}

	// Get or create bucket for this node
	e.mu.Lock()
	bucket, exists := e.buckets[node.ID]
	if !exists {
		bucket = &rateLimitBucket{
			requests:    make([]time.Time, 0),
			maxRequests: maxRequests,
			window:      duration,
		}
		e.buckets[node.ID] = bucket
	}
	e.mu.Unlock()

	// Check and enforce rate limit
	now := time.Now()
	
	e.mu.Lock()
	// Remove requests outside the current window
	cutoff := now.Add(-bucket.window)
	validRequests := make([]time.Time, 0)
	for _, t := range bucket.requests {
		if t.After(cutoff) {
			validRequests = append(validRequests, t)
		}
	}
	bucket.requests = validRequests

	// Check if we're at the limit
	if len(bucket.requests) >= bucket.maxRequests {
		e.mu.Unlock()
		// Calculate how long to wait
		oldestRequest := bucket.requests[0]
		waitTime := bucket.window - now.Sub(oldestRequest)
		if waitTime > 0 {
			time.Sleep(waitTime)
		}
		// After waiting, retry
		e.mu.Lock()
		// Clean up again after waiting
		cutoff = time.Now().Add(-bucket.window)
		validRequests = make([]time.Time, 0)
		for _, t := range bucket.requests {
			if t.After(cutoff) {
				validRequests = append(validRequests, t)
			}
		}
		bucket.requests = validRequests
	}

	// Record this request
	bucket.requests = append(bucket.requests, time.Now())
	currentCount := len(bucket.requests)
	e.mu.Unlock()

	return map[string]interface{}{
		"value":          inputValue,
		"rate_limited":   true,
		"requests_count": currentCount,
		"max_requests":   maxRequests,
		"window":         perDuration,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *RateLimiterExecutor) NodeType() types.NodeType {
	return types.NodeTypeRateLimiter
}

// Validate checks if node configuration is valid
func (e *RateLimiterExecutor) Validate(node types.Node) error {
	if node.Data.MaxRequests != nil && *node.Data.MaxRequests <= 0 {
		return fmt.Errorf("max_requests must be positive")
	}
	if node.Data.PerDuration != nil {
		if _, err := parseDuration(*node.Data.PerDuration); err != nil {
			return fmt.Errorf("invalid per_duration format: %w", err)
		}
	}
	if node.Data.RateLimitStrategy != nil {
		strategy := *node.Data.RateLimitStrategy
		if strategy != "fixed_window" && strategy != "sliding_window" && strategy != "token_bucket" {
			return fmt.Errorf("invalid strategy: %s (must be fixed_window, sliding_window, or token_bucket)", strategy)
		}
		if strategy != "fixed_window" {
			return fmt.Errorf("strategy %s not yet implemented (only fixed_window supported)", strategy)
		}
	}
	return nil
}
