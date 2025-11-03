package httpclient

import (
	"fmt"
	"sync"
)

// Registry manages named HTTP clients
type Registry struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

// NewRegistry creates a new HTTP client registry
func NewRegistry() *Registry {
	return &Registry{
		clients: make(map[string]*Client),
	}
}

// Register adds a client to the registry
func (r *Registry) Register(name string, client *Client) error {
	if name == "" {
		return fmt.Errorf("client name cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.clients[name]; exists {
		return fmt.Errorf("client with name %q already exists", name)
	}

	r.clients[name] = client
	return nil
}

// Get retrieves a client by name
func (r *Registry) Get(name string) (*Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, exists := r.clients[name]
	if !exists {
		return nil, fmt.Errorf("client %q not found", name)
	}

	return client, nil
}

// Has checks if a client exists
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.clients[name]
	return exists
}

// List returns all registered client names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.clients))
	for name := range r.clients {
		names = append(names, name)
	}
	return names
}

// Count returns the number of registered clients
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.clients)
}

// Clear removes all clients from the registry
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clients = make(map[string]*Client)
}
