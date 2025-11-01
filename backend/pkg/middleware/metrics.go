package middleware

import (
	"sync"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// MetricsCollector defines the interface for metrics collection
type MetricsCollector interface {
	RecordNodeExecution(nodeType types.NodeType, duration time.Duration, success bool)
	RecordNodeError(nodeType types.NodeType, errorType string)
}

// MetricsMiddleware collects execution metrics for nodes.
// It records execution time, success/failure rates, and error types.
type MetricsMiddleware struct {
	collector MetricsCollector
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(collector MetricsCollector) *MetricsMiddleware {
	return &MetricsMiddleware{
		collector: collector,
	}
}

// Process records metrics for node execution
func (m *MetricsMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	startTime := time.Now()

	// Execute the node
	result, err := next(ctx, node)

	duration := time.Since(startTime)
	success := err == nil

	// Record metrics
	if m.collector != nil {
		m.collector.RecordNodeExecution(node.Type, duration, success)
		if err != nil {
			m.collector.RecordNodeError(node.Type, err.Error())
		}
	}

	return result, err
}

// Name returns the middleware name
func (m *MetricsMiddleware) Name() string {
	return "Metrics"
}

// InMemoryMetricsCollector is a simple in-memory metrics collector for testing
type InMemoryMetricsCollector struct {
	mu               sync.RWMutex
	executionCount   map[types.NodeType]int64
	successCount     map[types.NodeType]int64
	failureCount     map[types.NodeType]int64
	totalDuration    map[types.NodeType]time.Duration
	errorCount       map[string]int64
}

// NewInMemoryMetricsCollector creates a new in-memory metrics collector
func NewInMemoryMetricsCollector() *InMemoryMetricsCollector {
	return &InMemoryMetricsCollector{
		executionCount: make(map[types.NodeType]int64),
		successCount:   make(map[types.NodeType]int64),
		failureCount:   make(map[types.NodeType]int64),
		totalDuration:  make(map[types.NodeType]time.Duration),
		errorCount:     make(map[string]int64),
	}
}

// RecordNodeExecution records a node execution
func (c *InMemoryMetricsCollector) RecordNodeExecution(nodeType types.NodeType, duration time.Duration, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.executionCount[nodeType]++
	c.totalDuration[nodeType] += duration

	if success {
		c.successCount[nodeType]++
	} else {
		c.failureCount[nodeType]++
	}
}

// RecordNodeError records a node error
func (c *InMemoryMetricsCollector) RecordNodeError(nodeType types.NodeType, errorType string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.errorCount[errorType]++
}

// GetExecutionCount returns the total execution count for a node type
func (c *InMemoryMetricsCollector) GetExecutionCount(nodeType types.NodeType) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.executionCount[nodeType]
}

// GetSuccessCount returns the success count for a node type
func (c *InMemoryMetricsCollector) GetSuccessCount(nodeType types.NodeType) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.successCount[nodeType]
}

// GetFailureCount returns the failure count for a node type
func (c *InMemoryMetricsCollector) GetFailureCount(nodeType types.NodeType) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.failureCount[nodeType]
}

// GetAverageDuration returns the average execution duration for a node type
func (c *InMemoryMetricsCollector) GetAverageDuration(nodeType types.NodeType) time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := c.executionCount[nodeType]
	if count == 0 {
		return 0
	}

	return c.totalDuration[nodeType] / time.Duration(count)
}

// GetErrorCount returns the count for a specific error type
func (c *InMemoryMetricsCollector) GetErrorCount(errorType string) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.errorCount[errorType]
}

// Reset clears all metrics
func (c *InMemoryMetricsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.executionCount = make(map[types.NodeType]int64)
	c.successCount = make(map[types.NodeType]int64)
	c.failureCount = make(map[types.NodeType]int64)
	c.totalDuration = make(map[types.NodeType]time.Duration)
	c.errorCount = make(map[string]int64)
}
