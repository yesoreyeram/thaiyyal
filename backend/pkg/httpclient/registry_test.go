package httpclient

import (
	"context"
	"net/http"
	"strconv"
	"testing"
)

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	config := &Config{
		UID:      "test-client",
		Auth: AuthConfig{Type: AuthTypeNone},
	}

	client, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test successful registration
	err = registry.Register("test-client", client)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// Test duplicate registration
	err = registry.Register("test-client", client)
	if err == nil {
		t.Error("Register() expected error for duplicate, got nil")
	}

	// Test empty UID
	err = registry.Register("", client)
	if err == nil {
		t.Error("Register() expected error for empty UID, got nil")
	}

	// Test nil client
	err = registry.Register("nil-client", nil)
	if err == nil {
		t.Error("Register() expected error for nil client, got nil")
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()

	config := &Config{
		UID:      "test-client",
		Auth: AuthConfig{Type: AuthTypeNone},
	}

	client, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Register client
	registry.Register("test-client", client)

	// Test successful get
	retrieved, err := registry.Get("test-client")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if retrieved == nil {
		t.Error("Get() returned nil client")
	}
	if retrieved != client {
		t.Error("Get() returned different client instance")
	}

	// Test get non-existent client
	_, err = registry.Get("non-existent")
	if err == nil {
		t.Error("Get() expected error for non-existent client, got nil")
	}
}

func TestRegistry_Has(t *testing.T) {
	registry := NewRegistry()

	config := &Config{
		UID:      "test-client",
		Auth: AuthConfig{Type: AuthTypeNone},
	}

	client, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test before registration
	if registry.Has("test-client") {
		t.Error("Has() returned true before registration")
	}

	// Register client
	registry.Register("test-client", client)

	// Test after registration
	if !registry.Has("test-client") {
		t.Error("Has() returned false after registration")
	}

	// Test non-existent client
	if registry.Has("non-existent") {
		t.Error("Has() returned true for non-existent client")
	}
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()

	// Empty registry
	if len(registry.List()) != 0 {
		t.Error("List() should return empty slice for empty registry")
	}

	// Add clients
	for i, uid := range []string{"client1", "client2", "client3"} {
		config := &Config{
			UID:      uid,
			Auth: AuthConfig{Type: AuthTypeNone},
		}
		client, err := New(context.Background(), config)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		registry.Register(uid, client)

		list := registry.List()
		if len(list) != i+1 {
			t.Errorf("List() length = %v, want %v", len(list), i+1)
		}
	}

	list := registry.List()
	if len(list) != 3 {
		t.Errorf("List() length = %v, want 3", len(list))
	}

	// Verify all UIDs are present
	uids := make(map[string]bool)
	for _, uid := range list {
		uids[uid] = true
	}

	for _, expected := range []string{"client1", "client2", "client3"} {
		if !uids[expected] {
			t.Errorf("List() missing expected UID %v", expected)
		}
	}
}

func TestRegistry_Count(t *testing.T) {
	registry := NewRegistry()

	// Empty registry
	if registry.Count() != 0 {
		t.Error("Count() should return 0 for empty registry")
	}

	// Add clients
	for i := 1; i <= 3; i++ {
		config := &Config{
			UID:      "client" + strconv.Itoa(i),
			Auth: AuthConfig{Type: AuthTypeNone},
		}
		client, err := New(context.Background(), config)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		registry.Register(config.UID, client)

		if registry.Count() != i {
			t.Errorf("Count() = %v, want %v", registry.Count(), i)
		}
	}
}

func TestRegistry_Clear(t *testing.T) {
	registry := NewRegistry()

	// Add clients
	for _, uid := range []string{"client1", "client2", "client3"} {
		config := &Config{
			UID:      uid,
			Auth: AuthConfig{Type: AuthTypeNone},
		}
		client, err := New(context.Background(), config)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		registry.Register(uid, client)
	}

	if registry.Count() != 3 {
		t.Errorf("Count() before clear = %v, want 3", registry.Count())
	}

	// Clear registry
	registry.Clear()

	if registry.Count() != 0 {
		t.Errorf("Count() after clear = %v, want 0", registry.Count())
	}

	// Verify clients are gone
	for _, uid := range []string{"client1", "client2", "client3"} {
		if registry.Has(uid) {
			t.Errorf("Has(%v) returned true after clear", uid)
		}
	}
}

func TestRegistry_Unregister(t *testing.T) {
	registry := NewRegistry()

	config := &Config{
		UID:      "test-client",
		Auth: AuthConfig{Type: AuthTypeNone},
	}

	client, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Register client
	registry.Register("test-client", client)

	// Verify it's there
	if !registry.Has("test-client") {
		t.Error("Client not found after registration")
	}

	// Unregister
	err = registry.Unregister("test-client")
	if err != nil {
		t.Errorf("Unregister() error = %v", err)
	}

	// Verify it's gone
	if registry.Has("test-client") {
		t.Error("Client still exists after unregister")
	}

	// Try to unregister again (should error)
	err = registry.Unregister("test-client")
	if err == nil {
		t.Error("Unregister() expected error for non-existent client, got nil")
	}
}

func TestRegistry_Concurrent(t *testing.T) {
	registry := NewRegistry()

	done := make(chan bool)

	// Concurrent registrations
	for i := 0; i < 10; i++ {
		go func(id int) {
			config := &Config{
				UID:      "client" + strconv.Itoa(id),
				Auth: AuthConfig{Type: AuthTypeNone},
			}
			client, err := New(context.Background(), config)
			if err != nil {
				t.Errorf("New() error = %v", err)
			}
			registry.Register(config.UID, client)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify count
	count := registry.Count()
	if count != 10 {
		t.Errorf("Count() after concurrent registrations = %v, want 10", count)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func(id int) {
			uid := "client" + strconv.Itoa(id)
			if !registry.Has(uid) {
				t.Errorf("Has(%v) returned false", uid)
			}
			_, err := registry.Get(uid)
			if err != nil {
				t.Errorf("Get(%v) error = %v", uid, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Benchmark registry operations
func BenchmarkRegistry_Register(b *testing.B) {
	registry := NewRegistry()
	clients := make([]*http.Client, b.N)

	for i := 0; i < b.N; i++ {
		config := &Config{
			UID:      "client" + strconv.Itoa(i),
			Auth: AuthConfig{Type: AuthTypeNone},
		}
		clients[i], _ = New(context.Background(), config)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.Register("client"+strconv.Itoa(i), clients[i])
	}
}

func BenchmarkRegistry_Get(b *testing.B) {
	registry := NewRegistry()

	// Pre-populate registry
	for i := 0; i < 100; i++ {
		config := &Config{
			UID:      "client" + strconv.Itoa(i),
			Auth: AuthConfig{Type: AuthTypeNone},
		}
		client, _ := New(context.Background(), config)
		registry.Register(config.UID, client)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.Get("client" + strconv.Itoa(i%100))
	}
}
