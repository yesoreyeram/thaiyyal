// Package engine provides workflow snapshot and restore functionality.
// This enables workflows to be paused, serialized, and resumed later.
package engine

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/graph"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/logging"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/observer"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/state"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Snapshot represents a complete snapshot of a workflow execution state.
// It contains all information needed to restore and resume execution from a specific point.
type Snapshot struct {
	// Metadata
	Version      string    `json:"version"`       // Snapshot format version
	SnapshotTime time.Time `json:"snapshot_time"` // When snapshot was created
	WorkflowID   string    `json:"workflow_id"`   // Workflow definition ID
	ExecutionID  string    `json:"execution_id"`  // Unique execution ID

	// Workflow Definition
	Nodes []types.Node `json:"nodes"` // Node definitions
	Edges []types.Edge `json:"edges"` // Edge definitions

	// Execution State
	Results        map[string]interface{} `json:"results"`         // Node execution results
	CompletedNodes []string               `json:"completed_nodes"` // IDs of completed nodes
	CurrentLevel   int                    `json:"current_level"`   // Current execution level in DAG

	// State Manager Data
	Variables     map[string]interface{}       `json:"variables"`      // Workflow variables
	Accumulator   interface{}                  `json:"accumulator"`    // Accumulator value
	Counter       float64                      `json:"counter"`        // Counter value
	Cache         map[string]*types.CacheEntry `json:"cache"`          // Cached entries
	ContextVars   map[string]interface{}       `json:"context_vars"`   // Context variables
	ContextConsts map[string]interface{}       `json:"context_consts"` // Context constants

	// Runtime Protection Counters
	NodeExecutionCount int `json:"node_execution_count"` // Number of nodes executed
	HTTPCallCount      int `json:"http_call_count"`      // Number of HTTP calls made

	// Configuration
	Config types.Config `json:"config"` // Engine configuration
}

// snapshotVersion is the current snapshot format version
const snapshotVersion = "1.0.0"

// SaveSnapshot creates a snapshot of the current execution state.
// This captures all state needed to resume execution from this point.
//
// The snapshot includes:
//   - Workflow metadata (IDs, timestamps)
//   - Node execution results
//   - State manager data (variables, cache, counters, accumulators)
//   - Execution progress (completed nodes, current level)
//   - Runtime protection counters
//   - Engine configuration
//
// Returns:
//   - *Snapshot: Complete execution state snapshot
//   - error: If snapshot creation fails
func (e *Engine) SaveSnapshot() (*Snapshot, error) {
	e.resultsMu.RLock()
	e.countersMu.RLock()
	defer e.resultsMu.RUnlock()
	defer e.countersMu.RUnlock()

	// Copy results to avoid race conditions
	results := make(map[string]interface{}, len(e.results))
	for k, v := range e.results {
		results[k] = v
	}

	// Get state manager data
	variables := e.state.ListVariables()
	accumulator := e.state.GetAccumulator()
	counter := e.state.GetCounter()
	cache := e.state.GetAllCache()
	contextVars := e.state.GetContextVariables()
	contextConsts := e.state.GetContextConstants()

	// Determine completed nodes from results
	completedNodes := make([]string, 0, len(results))
	for nodeID := range results {
		completedNodes = append(completedNodes, nodeID)
	}

	snapshot := &Snapshot{
		Version:            snapshotVersion,
		SnapshotTime:       time.Now(),
		WorkflowID:         e.workflowID,
		ExecutionID:        e.executionID,
		Nodes:              e.nodes,
		Edges:              e.edges,
		Results:            results,
		CompletedNodes:     completedNodes,
		CurrentLevel:       0, // Will be computed during restore if needed
		Variables:          variables,
		Accumulator:        accumulator,
		Counter:            counter,
		Cache:              cache,
		ContextVars:        contextVars,
		ContextConsts:      contextConsts,
		NodeExecutionCount: e.nodeExecutionCount,
		HTTPCallCount:      e.httpCallCount,
		Config:             e.config,
	}

	return snapshot, nil
}

// LoadSnapshot restores a workflow execution from a snapshot.
// This creates a new Engine instance with state restored from the snapshot.
//
// The restored engine will:
//   - Have the same execution ID as the original
//   - Have all node results from before the snapshot
//   - Have all state manager data restored
//   - Have runtime counters restored
//   - Be ready to resume execution from where it left off
//
// Parameters:
//   - snapshot: Previously saved snapshot
//   - registry: Executor registry (can be nil to use default)
//
// Returns:
//   - *Engine: Restored engine ready for execution
//   - error: If snapshot is invalid or restoration fails
func LoadSnapshot(snapshot *Snapshot, registry *executor.Registry) (*Engine, error) {
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot cannot be nil")
	}

	// Validate snapshot version
	if snapshot.Version != snapshotVersion {
		return nil, fmt.Errorf("unsupported snapshot version: %s (expected %s)", snapshot.Version, snapshotVersion)
	}

	// Use default registry if none provided
	if registry == nil {
		registry = DefaultRegistry()
	}

	// Create structured logger with workflow and execution context
	structuredLogger := logging.New(logging.DefaultConfig()).
		WithWorkflowID(snapshot.WorkflowID).
		WithExecutionID(snapshot.ExecutionID)

	// Create new engine with restored state
	engine := &Engine{
		state:              state.New(),
		registry:           registry,
		config:             snapshot.Config,
		results:            make(map[string]interface{}),
		executionID:        snapshot.ExecutionID, // Keep original execution ID
		workflowID:         snapshot.WorkflowID,
		nodes:              snapshot.Nodes,
		edges:              snapshot.Edges,
		nodeExecutionCount: snapshot.NodeExecutionCount,
		httpCallCount:      snapshot.HTTPCallCount,
		observerMgr:        observer.NewManager(),
		logger:             &observer.NoOpLogger{},
		structuredLogger:   structuredLogger,
	}

	// Restore results
	for nodeID, result := range snapshot.Results {
		engine.results[nodeID] = result
	}

	// Restore state manager data
	for name, value := range snapshot.Variables {
		engine.state.SetVariable(name, value)
	}

	if snapshot.Accumulator != nil {
		engine.state.SetAccumulator(snapshot.Accumulator)
	}

	engine.state.SetCounter(snapshot.Counter)

	// Restore cache entries (with updated TTL)
	now := time.Now()
	for key, entry := range snapshot.Cache {
		// Only restore non-expired entries
		if now.Before(entry.Expiration) {
			ttl := entry.Expiration.Sub(now)
			engine.state.SetCache(key, entry.Value, ttl)
		}
	}

	// Restore context variables and constants
	for name, value := range snapshot.ContextVars {
		engine.state.SetContextVariable(name, value)
	}

	for name, value := range snapshot.ContextConsts {
		engine.state.SetContextConstant(name, value)
	}

	// Create graph for topological sorting
	engine.graph = graph.New(snapshot.Nodes, snapshot.Edges)

	return engine, nil
}

// SerializeSnapshot converts a snapshot to JSON bytes.
// The JSON format is human-readable and suitable for storage or transmission.
//
// Parameters:
//   - snapshot: Snapshot to serialize
//
// Returns:
//   - []byte: JSON representation
//   - error: If serialization fails
func SerializeSnapshot(snapshot *Snapshot) ([]byte, error) {
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot cannot be nil")
	}
	return json.MarshalIndent(snapshot, "", "  ")
}

// DeserializeSnapshot converts JSON bytes to a Snapshot.
// The JSON must be in the format produced by SerializeSnapshot.
//
// Parameters:
//   - data: JSON bytes to deserialize
//
// Returns:
//   - *Snapshot: Deserialized snapshot
//   - error: If deserialization fails or JSON is invalid
func DeserializeSnapshot(data []byte) (*Snapshot, error) {
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to deserialize snapshot: %w", err)
	}

	// Validate deserialized snapshot
	if snapshot.Version == "" {
		return nil, fmt.Errorf("invalid snapshot: missing version")
	}

	if snapshot.ExecutionID == "" {
		return nil, fmt.Errorf("invalid snapshot: missing execution_id")
	}

	return &snapshot, nil
}

// ExecuteFromSnapshot is a convenience method that loads a snapshot and executes the workflow.
// This combines LoadSnapshot and Execute in a single call.
//
// Note: The current implementation will re-execute all nodes in the workflow.
// Future versions may support incremental execution to skip already-completed nodes.
//
// Parameters:
//   - snapshot: Previously saved snapshot
//   - registry: Executor registry (can be nil to use default)
//
// Returns:
//   - *types.Result: Workflow execution results (including pre-snapshot results)
//   - error: If restoration or execution fails
func ExecuteFromSnapshot(snapshot *Snapshot, registry *executor.Registry) (*types.Result, error) {
	engine, err := LoadSnapshot(snapshot, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to load snapshot: %w", err)
	}

	return engine.Execute()
}
