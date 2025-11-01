package executor

import (
	"fmt"
	"sync"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Registry manages node executor registration and lookup.
// It provides thread-safe registration and execution of node executors.
type Registry struct {
	executors map[types.NodeType]NodeExecutor
	mu        sync.RWMutex
}

// NewRegistry creates a new executor registry
func NewRegistry() *Registry {
	return &Registry{
		executors: make(map[types.NodeType]NodeExecutor),
	}
}

// Register adds an executor to the registry.
// Returns error if an executor for this type already exists.
func (r *Registry) Register(exec NodeExecutor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodeType := exec.NodeType()
	if _, exists := r.executors[nodeType]; exists {
		return fmt.Errorf("executor already registered for type: %s", nodeType)
	}

	r.executors[nodeType] = exec
	return nil
}

// MustRegister registers an executor and panics on error.
// Useful for initialization where executor registration must succeed.
func (r *Registry) MustRegister(exec NodeExecutor) {
	if err := r.Register(exec); err != nil {
		panic(err)
	}
}

// Execute dispatches execution to the appropriate executor for the node type.
// Returns error if no executor is registered for the node type.
func (r *Registry) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	r.mu.RLock()
	exec, exists := r.executors[node.Type]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no executor registered for type: %s", node.Type)
	}

	return exec.Execute(ctx, node)
}

// Validate validates a node using its registered executor.
// Returns error if no executor is registered or if validation fails.
func (r *Registry) Validate(node types.Node) error {
	r.mu.RLock()
	exec, exists := r.executors[node.Type]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no executor registered for type: %s", node.Type)
	}

	return exec.Validate(node)
}

// GetExecutor returns the executor for a given node type.
// Returns nil if no executor is registered.
func (r *Registry) GetExecutor(nodeType types.NodeType) NodeExecutor {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.executors[nodeType]
}

// ListRegisteredTypes returns all registered node types
func (r *Registry) ListRegisteredTypes() []types.NodeType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]types.NodeType, 0, len(r.executors))
	for nodeType := range r.executors {
		types = append(types, nodeType)
	}
	return types
}
