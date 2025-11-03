package httpclient

import (
	"fmt"
	"net/http"
	"sync"
)

// Registry manages named HTTP clients by their UIDs
type Registry struct {
	clients map[string]*http.Client
	mu      sync.RWMutex
}

// NewRegistry creates a new HTTP client registry
func NewRegistry() *Registry {
	return &Registry{
		clients: make(map[string]*http.Client),
	}
}

// Register adds a client to the registry with the given UID
func (r *Registry) Register(uid string, client *http.Client) error {
	if uid == "" {
		return fmt.Errorf("client UID cannot be empty")
	}

	if client == nil {
		return fmt.Errorf("client cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.clients[uid]; exists {
		return fmt.Errorf("client with UID %q already exists", uid)
	}

	r.clients[uid] = client
	return nil
}

// Get retrieves a client by UID
func (r *Registry) Get(uid string) (*http.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, exists := r.clients[uid]
	if !exists {
		return nil, fmt.Errorf("client with UID %q not found", uid)
	}

	return client, nil
}

// Has checks if a client exists with the given UID
func (r *Registry) Has(uid string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.clients[uid]
	return exists
}

// List returns all registered client UIDs
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	uids := make([]string, 0, len(r.clients))
	for uid := range r.clients {
		uids = append(uids, uid)
	}
	return uids
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

	r.clients = make(map[string]*http.Client)
}

// Unregister removes a client from the registry
func (r *Registry) Unregister(uid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.clients[uid]; !exists {
		return fmt.Errorf("client with UID %q not found", uid)
	}

	delete(r.clients, uid)
	return nil
}
