// Package executor provides the Strategy Pattern implementation for node execution.
// This replaces the large switch statement with a registry of executor strategies.
package executor

import (
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ExecutionContext provides access to workflow state and operations.
// This interface breaks the circular dependency between executor and engine.
//
// Executors receive this context and can access workflow state without
// directly depending on the engine implementation.
type ExecutionContext interface {
	// Input retrieval
	GetNodeInputs(nodeID string) []interface{}
	GetNode(nodeID string) *types.Node

	// State management
	GetVariable(name string) (interface{}, error)
	SetVariable(name string, value interface{}) error
	GetAccumulator() interface{}
	SetAccumulator(value interface{})
	GetCounter() float64
	SetCounter(value float64)
	GetCache(key string) (interface{}, bool)
	SetCache(key string, value interface{}, ttl time.Duration)

	// Context operations
	GetWorkflowContext() map[string]interface{}
	GetContextVariable(name string) (interface{}, bool)
	SetContextVariable(name string, value interface{})
	GetContextConstant(name string) (interface{}, bool)
	SetContextConstant(name string, value interface{})
	InterpolateTemplate(template string) string

	// Result management
	GetNodeResult(nodeID string) (interface{}, bool)
	SetNodeResult(nodeID string, result interface{})
	GetAllNodeResults() map[string]interface{}
	GetVariables() map[string]interface{}
	GetContextVariables() map[string]interface{}
	
	// Configuration
	GetConfig() types.Config
}

// NodeExecutor defines the interface for node execution strategies.
// Each node type has its own executor implementation.
type NodeExecutor interface {
	// Execute runs the node with given context
	Execute(ctx ExecutionContext, node types.Node) (interface{}, error)

	// NodeType returns the type this executor handles
	NodeType() types.NodeType

	// Validate checks if node configuration is valid
	Validate(node types.Node) error
}
