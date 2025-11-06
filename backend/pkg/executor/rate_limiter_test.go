package executor

import (
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestRateLimiterExecutor_Execute(t *testing.T) {
	tests := []struct {
		name            string
		maxRequests     int
		perDuration     string
		numRequests     int
		expectedMinTime time.Duration
	}{
		{
			name:            "basic rate limiting - 2 requests per second",
			maxRequests:     2,
			perDuration:     "1s",
			numRequests:     4,
			expectedMinTime: 1 * time.Second, // Should take at least 1 second for 4 requests (2 per second)
		},
		{
			name:            "fast rate - 10 requests per 100ms",
			maxRequests:     10,
			perDuration:     "100ms",
			numRequests:     15,
			expectedMinTime: 100 * time.Millisecond, // Should take at least 100ms for 15 requests (10 per 100ms)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewRateLimiterExecutor()
			ctx := &MockExecutionContext{
				inputs: make(map[string][]interface{}),
			}

			node := types.Node{
				ID:   "rate1",
				Type: types.NodeTypeRateLimiter,
				Data: types.RateLimiterData{
					MaxRequests: &tt.maxRequests,
					PerDuration: &tt.perDuration,
				},
			}

			start := time.Now()
			for i := 0; i < tt.numRequests; i++ {
				_, err := executor.Execute(ctx, node)
				if err != nil {
					t.Fatalf("request %d failed: %v", i, err)
				}
			}
			elapsed := time.Since(start)

			if elapsed < tt.expectedMinTime {
				t.Errorf("rate limiting not working: expected at least %v, got %v", tt.expectedMinTime, elapsed)
			}
		})
	}
}

func TestRateLimiterExecutor_Validate(t *testing.T) {
	executor := NewRateLimiterExecutor()

	tests := []struct {
		name    string
		node    types.Node
		wantErr bool
	}{
		{
			name: "valid configuration",
			node: types.Node{
				ID:   "rate1",
				Type: types.NodeTypeRateLimiter,
				Data: types.RateLimiterData{
					MaxRequests: intPtr(10),
					PerDuration: strPtr("1s"),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid max_requests (zero)",
			node: types.Node{
				ID:   "rate1",
				Type: types.NodeTypeRateLimiter,
				Data: types.RateLimiterData{
					MaxRequests: intPtr(0),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid max_requests (negative)",
			node: types.Node{
				ID:   "rate1",
				Type: types.NodeTypeRateLimiter,
				Data: types.RateLimiterData{
					MaxRequests: intPtr(-5),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid duration format",
			node: types.Node{
				ID:   "rate1",
				Type: types.NodeTypeRateLimiter,
				Data: types.RateLimiterData{
					PerDuration: strPtr("invalid"),
				},
			},
			wantErr: true,
		},
		{
			name: "unsupported strategy",
			node: types.Node{
				ID:   "rate1",
				Type: types.NodeTypeRateLimiter,
				Data: types.RateLimiterData{
					RateLimitStrategy: strPtr("sliding_window"),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.Validate(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRateLimiterExecutor_PassThrough(t *testing.T) {
	executor := NewRateLimiterExecutor()
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"rate1": {"test value"},
		},
	}

	maxReq := 10
	node := types.Node{
		ID:   "rate1",
		Type: types.NodeTypeRateLimiter,
		Data: types.RateLimiterData{
			MaxRequests: &maxReq,
			PerDuration: strPtr("1s"),
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("result is not a map")
	}

	if resultMap["value"] != "test value" {
		t.Errorf("value not passed through correctly: got %v", resultMap["value"])
	}

	if resultMap["rate_limited"] != true {
		t.Errorf("rate_limited flag not set")
	}
}

func TestThrottleExecutor_Execute(t *testing.T) {
	tests := []struct {
		name              string
		requestsPerSecond float64
		numRequests       int
		expectedMinTime   time.Duration
	}{
		{
			name:              "throttle to 10 rps",
			requestsPerSecond: 10,
			numRequests:       5,
			expectedMinTime:   400 * time.Millisecond, // 5 requests at 10/s = 500ms, but first is immediate
		},
		{
			name:              "throttle to 2 rps",
			requestsPerSecond: 2,
			numRequests:       3,
			expectedMinTime:   1 * time.Second, // 3 requests at 2/s = 1.5s, but first is immediate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewThrottleExecutor()
			ctx := &MockExecutionContext{
				inputs: make(map[string][]interface{}),
			}

			node := types.Node{
				ID:   "throttle1",
				Type: types.NodeTypeThrottle,
				Data: types.ThrottleData{
					RequestsPerSecond: &tt.requestsPerSecond,
				},
			}

			start := time.Now()
			for i := 0; i < tt.numRequests; i++ {
				_, err := executor.Execute(ctx, node)
				if err != nil {
					t.Fatalf("request %d failed: %v", i, err)
				}
			}
			elapsed := time.Since(start)

			if elapsed < tt.expectedMinTime {
				t.Errorf("throttling not working: expected at least %v, got %v", tt.expectedMinTime, elapsed)
			}
		})
	}
}

func TestThrottleExecutor_Validate(t *testing.T) {
	executor := NewThrottleExecutor()

	tests := []struct {
		name    string
		node    types.Node
		wantErr bool
	}{
		{
			name: "valid configuration",
			node: types.Node{
				ID:   "throttle1",
				Type: types.NodeTypeThrottle,
				Data: types.ThrottleData{
					RequestsPerSecond: float64Ptr(5.0),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid requests_per_second (zero)",
			node: types.Node{
				ID:   "throttle1",
				Type: types.NodeTypeThrottle,
				Data: types.ThrottleData{
					RequestsPerSecond: float64Ptr(0.0),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid requests_per_second (negative)",
			node: types.Node{
				ID:   "throttle1",
				Type: types.NodeTypeThrottle,
				Data: types.ThrottleData{
					RequestsPerSecond: float64Ptr(-2.5),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.Validate(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestThrottleExecutor_PassThrough(t *testing.T) {
	executor := NewThrottleExecutor()
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"throttle1": {map[string]interface{}{
				"data": "test",
			}},
		},
	}

	rps := 10.0
	node := types.Node{
		ID:   "throttle1",
		Type: types.NodeTypeThrottle,
		Data: types.ThrottleData{
			RequestsPerSecond: &rps,
		},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("result is not a map")
	}

	valueMap, ok := resultMap["value"].(map[string]interface{})
	if !ok {
		t.Fatalf("value is not a map")
	}

	if valueMap["data"] != "test" {
		t.Errorf("value not passed through correctly")
	}

	if resultMap["throttled"] != true {
		t.Errorf("throttled flag not set")
	}
}

// Helper functions
func float64Ptr(f float64) *float64 {
	return &f
}
