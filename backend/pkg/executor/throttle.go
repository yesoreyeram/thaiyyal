package executor

import (
	"fmt"
	"sync"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ThrottleExecutor executes Throttle nodes
// Simple delay-based request throttling
type ThrottleExecutor struct {
	mu         sync.Mutex
	lastRequest map[string]time.Time
}

// NewThrottleExecutor creates a new ThrottleExecutor
func NewThrottleExecutor() *ThrottleExecutor {
	return &ThrottleExecutor{
		lastRequest: make(map[string]time.Time),
	}
}

// Execute runs the Throttle node
// Adds delay between requests based on requests_per_second
func (e *ThrottleExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsThrottleData(node.Data)
if err != nil {
return nil, err
}
	inputs := ctx.GetNodeInputs(node.ID)
	var inputValue interface{}
	if len(inputs) > 0 {
		inputValue = inputs[0]
	}

	// Get requests per second (default: 5)
	requestsPerSecond := 5.0
	if data.RequestsPerSecond != nil {
		requestsPerSecond = *data.RequestsPerSecond
	}

	if requestsPerSecond <= 0 {
		return nil, fmt.Errorf("requests_per_second must be positive, got %f", requestsPerSecond)
	}

	// Calculate minimum delay between requests (in milliseconds)
	minDelay := time.Duration(float64(time.Second) / requestsPerSecond)

	// Get last request time for this node
	e.mu.Lock()
	lastReq, exists := e.lastRequest[node.ID]
	e.mu.Unlock()

	var actualDelay time.Duration
	if exists {
		// Calculate how long since last request
		elapsed := time.Since(lastReq)
		if elapsed < minDelay {
			// Need to delay
			waitTime := minDelay - elapsed
			time.Sleep(waitTime)
			actualDelay = waitTime
		}
	}

	// Update last request time
	e.mu.Lock()
	e.lastRequest[node.ID] = time.Now()
	e.mu.Unlock()

	return map[string]interface{}{
		"value":               inputValue,
		"throttled":           true,
		"requests_per_second": requestsPerSecond,
		"delay_ms":            actualDelay.Milliseconds(),
	}, nil
}

// NodeType returns the node type this executor handles
func (e *ThrottleExecutor) NodeType() types.NodeType {
	return types.NodeTypeThrottle
}

// Validate checks if node configuration is valid
func (e *ThrottleExecutor) Validate(node types.Node) error {
data, err := types.AsThrottleData(node.Data)
if err != nil {
return err
}
	if data.RequestsPerSecond != nil && *data.RequestsPerSecond <= 0 {
		return fmt.Errorf("requests_per_second must be positive, got %f", *data.RequestsPerSecond)
	}
	return nil
}
