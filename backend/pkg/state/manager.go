// Package state provides state management for workflow execution.
// This includes variables, accumulators, counters, and caching.
package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Manager manages workflow execution state
type Manager struct {
	variables map[string]interface{}
	accumulator interface{}
	counter     float64
	cache       map[string]*types.CacheEntry
	
	// Context for template interpolation
	contextVariables map[string]interface{}
	contextConstants map[string]interface{}
	
	mu sync.RWMutex
}

// New creates a new state manager
func New() *Manager {
	return &Manager{
		variables:        make(map[string]interface{}),
		cache:            make(map[string]*types.CacheEntry),
		contextVariables: make(map[string]interface{}),
		contextConstants: make(map[string]interface{}),
	}
}

// ============================================================================
// Variable Operations
// ============================================================================

// GetVariable retrieves a variable value
func (m *Manager) GetVariable(name string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.variables[name]
	if !ok {
		return nil, fmt.Errorf("variable not found: %s", name)
	}
	return val, nil
}

// SetVariable sets a variable value
func (m *Manager) SetVariable(name string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.variables[name] = value
	return nil
}

// DeleteVariable removes a variable
func (m *Manager) DeleteVariable(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.variables, name)
	return nil
}

// ListVariables returns all variables
func (m *Manager) ListVariables() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[string]interface{}, len(m.variables))
	for k, v := range m.variables {
		result[k] = v
	}
	return result
}

// GetAllVariables is an alias for ListVariables
func (m *Manager) GetAllVariables() map[string]interface{} {
	return m.ListVariables()
}

// ============================================================================
// Accumulator Operations
// ============================================================================

// GetAccumulator returns the current accumulator value
func (m *Manager) GetAccumulator() interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.accumulator
}

// SetAccumulator sets the accumulator value
func (m *Manager) SetAccumulator(value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.accumulator = value
}

// ResetAccumulator resets the accumulator to nil
func (m *Manager) ResetAccumulator() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.accumulator = nil
}

// ============================================================================
// Counter Operations
// ============================================================================

// GetCounter returns the current counter value
func (m *Manager) GetCounter() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.counter
}

// SetCounter sets the counter value
func (m *Manager) SetCounter(value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter = value
}

// IncrementCounter increments the counter by delta
func (m *Manager) IncrementCounter(delta float64) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter += delta
	return m.counter
}

// ResetCounter resets the counter to a specific value
func (m *Manager) ResetCounter(value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter = value
}

// ============================================================================
// Cache Operations
// ============================================================================

// GetCache retrieves a cached value
// Returns the value, whether it was found, and any error
func (m *Manager) GetCache(key string) (interface{}, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.cache[key]
	if !exists {
		return nil, false, nil
	}

	// Check if entry has expired
	if time.Now().After(entry.Expiration) {
		return nil, false, nil
	}

	return entry.Value, true, nil
}

// SetCache sets a cached value with TTL
func (m *Manager) SetCache(key string, value interface{}, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache[key] = &types.CacheEntry{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
	return nil
}

// DeleteCache removes a cached entry
func (m *Manager) DeleteCache(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.cache, key)
	return nil
}

// ClearCache removes all cached entries
func (m *Manager) ClearCache() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache = make(map[string]*types.CacheEntry)
	return nil
}

// CleanExpiredCache removes expired cache entries
func (m *Manager) CleanExpiredCache() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for key, entry := range m.cache {
		if now.After(entry.Expiration) {
			delete(m.cache, key)
		}
	}
}

// ============================================================================
// Context Operations
// ============================================================================

// GetContextVariable retrieves a context variable
func (m *Manager) GetContextVariable(name string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.contextVariables[name]
	return val, ok
}

// SetContextVariable sets a context variable
func (m *Manager) SetContextVariable(name string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.contextVariables[name] = value
}

// GetContextConstant retrieves a context constant
func (m *Manager) GetContextConstant(name string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.contextConstants[name]
	return val, ok
}

// SetContextConstant sets a context constant
func (m *Manager) SetContextConstant(name string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.contextConstants[name] = value
}

// GetAllContext returns all context variables and constants
func (m *Manager) GetAllContext() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{})
	
	// Add constants first
	for k, v := range m.contextConstants {
		result[k] = v
	}
	
	// Add variables (which can override constants)
	for k, v := range m.contextVariables {
		result[k] = v
	}
	
	return result
}

// GetAllCache returns all cache entries
// Used for snapshot/restore functionality
func (m *Manager) GetAllCache() map[string]*types.CacheEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a shallow copy of cache entries to avoid race conditions
	// Note: This is a shallow copy. If Value contains complex types with pointers,
	// they will share references with the original. For true deep copy, use
	// serialization/deserialization or a deep copy library.
	result := make(map[string]*types.CacheEntry, len(m.cache))
	for k, v := range m.cache {
		// Shallow copy the cache entry
		result[k] = &types.CacheEntry{
			Value:      v.Value,
			Expiration: v.Expiration,
		}
	}
	return result
}

// GetContextVariables returns all context variables (without constants)
func (m *Manager) GetContextVariables() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{}, len(m.contextVariables))
	for k, v := range m.contextVariables {
		result[k] = v
	}
	return result
}

// GetContextConstants returns all context constants
func (m *Manager) GetContextConstants() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{}, len(m.contextConstants))
	for k, v := range m.contextConstants {
		result[k] = v
	}
	return result
}
