// Package state provides state management for workflow execution.
//
// # Overview
//
// The state package implements thread-safe state management for workflow execution,
// enabling nodes to share data, maintain execution context, and persist values
// across workflow runs. It provides a key-value store with scope management,
// TTL support, and transaction semantics.
//
// # Features
//
//   - Key-value storage: Store and retrieve arbitrary values
//   - Scope isolation: Workflow, node, and global scopes
//   - Thread-safe: Concurrent access from multiple goroutines
//   - TTL support: Automatic expiration of values
//   - Transactional: Atomic multi-key operations
//   - Type preservation: Store any Go type
//   - Persistence: Optional persistence to disk
//
// # State Scopes
//
// The state manager supports three scopes:
//
// Global Scope:
//   - Shared across all workflows
//   - Useful for configuration, caches
//   - Persists across workflow runs
//
// Workflow Scope:
//   - Isolated per workflow execution
//   - Cleaned up after workflow completes
//   - Useful for workflow-level variables
//
// Node Scope:
//   - Isolated per node execution
//   - Cleaned up after node completes
//   - Useful for node-specific state
//
// # Basic Usage
//
//	import "github.com/yesoreyeram/thaiyyal/backend/pkg/state"
//
//	// Create state manager
//	sm := state.NewManager()
//
//	// Store value
//	sm.Set("counter", 0, state.ScopeGlobal)
//
//	// Retrieve value
//	value, exists := sm.Get("counter", state.ScopeGlobal)
//	if exists {
//	    count := value.(int)
//	}
//
//	// Delete value
//	sm.Delete("counter", state.ScopeGlobal)
//
// # Workflow-Scoped State
//
//	// Create workflow state
//	workflowID := "workflow-123"
//	sm.Set("workflow:"+workflowID+":counter", 0, state.ScopeWorkflow)
//
//	// Access in node
//	value, _ := sm.Get("workflow:"+workflowID+":counter", state.ScopeWorkflow)
//
//	// Cleanup after workflow
//	sm.ClearWorkflow(workflowID)
//
// # TTL Support
//
//	// Store value with 1-hour TTL
//	sm.SetWithTTL("cache:user:123", userData, state.ScopeGlobal, 1*time.Hour)
//
//	// Value automatically expires after TTL
//	time.Sleep(1 * time.Hour)
//	_, exists := sm.Get("cache:user:123", state.ScopeGlobal) // exists == false
//
// # Transactions
//
//	// Begin transaction
//	tx := sm.BeginTransaction()
//
//	// Perform multiple operations
//	tx.Set("key1", value1, state.ScopeGlobal)
//	tx.Set("key2", value2, state.ScopeGlobal)
//	tx.Delete("key3", state.ScopeGlobal)
//
//	// Commit atomically
//	if err := tx.Commit(); err != nil {
//	    tx.Rollback()
//	}
//
// # Atomic Operations
//
//	// Increment counter atomically
//	newValue := sm.Increment("counter", 1, state.ScopeGlobal)
//
//	// Compare-and-swap
//	success := sm.CompareAndSwap("key", oldValue, newValue, state.ScopeGlobal)
//
// # State Nodes Integration
//
// Variable Node:
//
//	// Set variable
//	sm.Set("myVariable", value, state.ScopeWorkflow)
//
//	// Get variable
//	value, _ := sm.Get("myVariable", state.ScopeWorkflow)
//
// Accumulator Node:
//
//	// Accumulate value
//	total := sm.Accumulate("accumulator", value, state.ScopeWorkflow)
//
// Counter Node:
//
//	// Increment counter
//	count := sm.Increment("counter", 1, state.ScopeWorkflow)
//
// Cache Node:
//
//	// Check cache
//	if cached, exists := sm.Get("cache:"+key, state.ScopeGlobal); exists {
//	    return cached
//	}
//
//	// Compute and cache
//	result := compute()
//	sm.SetWithTTL("cache:"+key, result, state.ScopeGlobal, 1*time.Hour)
//
// # Persistence
//
//	// Create persistent state manager
//	sm := state.NewManager(state.Config{
//	    Persistent: true,
//	    StoragePath: "/var/lib/thaiyyal/state",
//	})
//
//	// State survives process restarts
//	sm.Set("config", configData, state.ScopeGlobal)
//
// # Cleanup
//
//	// Clear workflow state
//	sm.ClearWorkflow(workflowID)
//
//	// Clear node state
//	sm.ClearNode(nodeID)
//
//	// Clear expired entries
//	sm.CleanupExpired()
//
//	// Clear all state
//	sm.Clear()
//
// # Memory Management
//
//	// Set max memory limit
//	sm.SetMaxMemory(100 * 1024 * 1024) // 100MB
//
//	// Evict old entries when limit reached
//	sm.EnableLRUEviction()
//
// # Monitoring
//
//	// Get state statistics
//	stats := sm.Stats()
//	fmt.Printf("Keys: %d, Memory: %d bytes\n", stats.KeyCount, stats.MemoryBytes)
//
// # Use Cases
//
//   - Workflow variables: Share data between nodes
//   - Accumulators: Aggregate values over iterations
//   - Counters: Track execution counts
//   - Caching: Cache expensive computations
//   - Configuration: Store workflow configuration
//   - Session state: Maintain user session data
//
// # Performance Characteristics
//
//   - Get: O(1) average case with map lookup
//   - Set: O(1) average case with map insertion
//   - Delete: O(1) average case with map deletion
//   - Clear: O(n) where n is number of keys
//   - Memory: Proportional to stored data size
//
// # Thread Safety
//
// All state manager operations are thread-safe using RWMutex:
//
//   - Concurrent reads are allowed
//   - Writes are serialized
//   - Transactions use optimistic locking
//
// # Error Handling
//
// State operations return errors for:
//
//   - Key not found
//   - Type mismatches
//   - Transaction conflicts
//   - Memory limits exceeded
//   - Persistence failures
//
// # Best Practices
//
//   - Use appropriate scope for isolation
//   - Set TTL for cache entries
//   - Clean up state after workflow completion
//   - Use transactions for multi-key updates
//   - Monitor memory usage
//   - Use atomic operations for counters
//   - Avoid storing large objects
//
// # Testing
//
// For testing, use an in-memory state manager:
//
//	sm := state.NewManager(state.Config{
//	    Persistent: false,
//	})
//
// # Integration
//
// State manager integrates with the engine:
//
//	engine := engine.New(
//	    engine.WithStateManager(sm),
//	)
package state
