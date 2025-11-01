package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// SizeLimitMiddleware enforces size limits to prevent memory exhaustion
type SizeLimitMiddleware struct {
	maxInputSize      int64 // Maximum size of input data per node (bytes)
	maxResultSize     int64 // Maximum size of result data per node (bytes)
	maxStringLength   int   // Maximum length of string values
	maxArrayLength    int   // Maximum length of arrays
	maxWorkflowSize   int64 // Maximum total workflow size (all nodes + edges)
	maxNodeCount      int   // Maximum number of nodes
	maxEdgeCount      int   // Maximum number of edges
	enforceInputSize  bool  // Whether to enforce input size limits
	enforceResultSize bool  // Whether to enforce result size limits
}

// SizeLimitConfig configures size limit enforcement
type SizeLimitConfig struct {
	// Per-node limits
	MaxInputSize    int64 // Maximum input size per node (default: 10MB)
	MaxResultSize   int64 // Maximum result size per node (default: 50MB)
	MaxStringLength int   // Maximum string length (default: 1MB)
	MaxArrayLength  int   // Maximum array length (default: 10000)
	
	// Workflow limits
	MaxWorkflowSize int64 // Maximum total workflow size (default: 100MB)
	MaxNodeCount    int   // Maximum nodes in workflow (default: 1000)
	MaxEdgeCount    int   // Maximum edges in workflow (default: 5000)
	
	// Control flags
	EnforceInputSize  bool // Enforce input size limits (default: true)
	EnforceResultSize bool // Enforce result size limits (default: true)
}

// DefaultSizeLimitConfig returns default size limit configuration
func DefaultSizeLimitConfig() SizeLimitConfig {
	return SizeLimitConfig{
		MaxInputSize:      10 * 1024 * 1024,  // 10 MB
		MaxResultSize:     50 * 1024 * 1024,  // 50 MB
		MaxStringLength:   1 * 1024 * 1024,   // 1 MB
		MaxArrayLength:    10000,             // 10k elements
		MaxWorkflowSize:   100 * 1024 * 1024, // 100 MB
		MaxNodeCount:      1000,              // 1000 nodes
		MaxEdgeCount:      5000,              // 5000 edges
		EnforceInputSize:  true,
		EnforceResultSize: true,
	}
}

// NewSizeLimitMiddleware creates a new size limit middleware with default config
func NewSizeLimitMiddleware() *SizeLimitMiddleware {
	return NewSizeLimitMiddlewareWithConfig(DefaultSizeLimitConfig())
}

// NewSizeLimitMiddlewareWithConfig creates a new size limit middleware with custom config
func NewSizeLimitMiddlewareWithConfig(config SizeLimitConfig) *SizeLimitMiddleware {
	return &SizeLimitMiddleware{
		maxInputSize:      config.MaxInputSize,
		maxResultSize:     config.MaxResultSize,
		maxStringLength:   config.MaxStringLength,
		maxArrayLength:    config.MaxArrayLength,
		maxWorkflowSize:   config.MaxWorkflowSize,
		maxNodeCount:      config.MaxNodeCount,
		maxEdgeCount:      config.MaxEdgeCount,
		enforceInputSize:  config.EnforceInputSize,
		enforceResultSize: config.EnforceResultSize,
	}
}

// Process enforces size limits on inputs and results
func (m *SizeLimitMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	// Check input size if enabled
	if m.enforceInputSize {
		inputs := ctx.GetNodeInputs(node.ID)
		if err := m.validateInputSize(inputs); err != nil {
			return nil, fmt.Errorf("input size limit exceeded: %w", err)
		}
	}
	
	// Execute node
	result, err := next(ctx, node)
	if err != nil {
		return result, err
	}
	
	// Check result size if enabled
	if m.enforceResultSize && result != nil {
		if err := m.validateResultSize(result); err != nil {
			return nil, fmt.Errorf("result size limit exceeded: %w", err)
		}
	}
	
	return result, nil
}

// Name returns the middleware name
func (m *SizeLimitMiddleware) Name() string {
	return "SizeLimit"
}

// validateInputSize validates the size of inputs
func (m *SizeLimitMiddleware) validateInputSize(inputs []interface{}) error {
	for i, input := range inputs {
		size, err := estimateSize(input)
		if err != nil {
			return fmt.Errorf("failed to estimate size of input %d: %w", i, err)
		}
		
		if size > m.maxInputSize {
			return fmt.Errorf("input %d size %d bytes exceeds limit %d bytes", i, size, m.maxInputSize)
		}
		
		// Check specific type limits
		if err := m.validateValue(input); err != nil {
			return fmt.Errorf("input %d validation failed: %w", i, err)
		}
	}
	
	return nil
}

// validateResultSize validates the size of result
func (m *SizeLimitMiddleware) validateResultSize(result interface{}) error {
	size, err := estimateSize(result)
	if err != nil {
		return fmt.Errorf("failed to estimate result size: %w", err)
	}
	
	if size > m.maxResultSize {
		return fmt.Errorf("result size %d bytes exceeds limit %d bytes", size, m.maxResultSize)
	}
	
	// Check specific type limits
	return m.validateValue(result)
}

// validateValue validates type-specific limits
func (m *SizeLimitMiddleware) validateValue(value interface{}) error {
	switch v := value.(type) {
	case string:
		if m.maxStringLength > 0 && len(v) > m.maxStringLength {
			return fmt.Errorf("string length %d exceeds limit %d", len(v), m.maxStringLength)
		}
	case []interface{}:
		if m.maxArrayLength > 0 && len(v) > m.maxArrayLength {
			return fmt.Errorf("array length %d exceeds limit %d", len(v), m.maxArrayLength)
		}
		// Recursively validate array elements
		for i, elem := range v {
			if err := m.validateValue(elem); err != nil {
				return fmt.Errorf("array element %d: %w", i, err)
			}
		}
	case map[string]interface{}:
		// Validate map values recursively
		for key, val := range v {
			if err := m.validateValue(val); err != nil {
				return fmt.Errorf("map key %s: %w", key, err)
			}
		}
	}
	
	return nil
}

// estimateSize estimates the size of a value in bytes
func estimateSize(value interface{}) (int64, error) {
	// Use JSON marshaling as a rough estimate of size
	// This is not exact but provides a reasonable approximation
	data, err := json.Marshal(value)
	if err != nil {
		return 0, err
	}
	return int64(len(data)), nil
}

// ValidateWorkflowSize validates workflow size limits
// This should be called before workflow execution
func ValidateWorkflowSize(nodes []types.Node, edges []types.Edge, config SizeLimitConfig) error {
	// Check node count
	if config.MaxNodeCount > 0 && len(nodes) > config.MaxNodeCount {
		return fmt.Errorf("workflow has %d nodes, exceeds limit of %d", len(nodes), config.MaxNodeCount)
	}
	
	// Check edge count
	if config.MaxEdgeCount > 0 && len(edges) > config.MaxEdgeCount {
		return fmt.Errorf("workflow has %d edges, exceeds limit of %d", len(edges), config.MaxEdgeCount)
	}
	
	// Check total workflow size
	if config.MaxWorkflowSize > 0 {
		// Estimate total size by marshaling
		type workflow struct {
			Nodes []types.Node `json:"nodes"`
			Edges []types.Edge `json:"edges"`
		}
		
		wf := workflow{Nodes: nodes, Edges: edges}
		data, err := json.Marshal(wf)
		if err != nil {
			return fmt.Errorf("failed to marshal workflow for size check: %w", err)
		}
		
		size := int64(len(data))
		if size > config.MaxWorkflowSize {
			return fmt.Errorf("workflow size %d bytes exceeds limit %d bytes", size, config.MaxWorkflowSize)
		}
	}
	
	return nil
}
