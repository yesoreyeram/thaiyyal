package httpclient

import (
	"strconv"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/config"
)

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)

	clientConfig := &ClientConfig{
		Name:     "test-client",
		AuthType: AuthTypeNone,
	}

	client, err := builder.Build(clientConfig)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
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

	// Test empty name
	err = registry.Register("", client)
	if err == nil {
		t.Error("Register() expected error for empty name, got nil")
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)

	clientConfig := &ClientConfig{
		Name:     "test-client",
		AuthType: AuthTypeNone,
	}

	client, err := builder.Build(clientConfig)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
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
	if retrieved.GetConfig().Name != "test-client" {
		t.Errorf("Get() client name = %v, want test-client", retrieved.GetConfig().Name)
	}

	// Test get non-existent client
	_, err = registry.Get("non-existent")
	if err == nil {
		t.Error("Get() expected error for non-existent client, got nil")
	}
}

func TestRegistry_Has(t *testing.T) {
	registry := NewRegistry()
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)

	clientConfig := &ClientConfig{
		Name:     "test-client",
		AuthType: AuthTypeNone,
	}

	client, err := builder.Build(clientConfig)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
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
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)

	// Empty registry
	if len(registry.List()) != 0 {
		t.Error("List() should return empty slice for empty registry")
	}

	// Add clients
	for i, name := range []string{"client1", "client2", "client3"} {
		clientConfig := &ClientConfig{
			Name:     name,
			AuthType: AuthTypeNone,
		}
		client, err := builder.Build(clientConfig)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		registry.Register(name, client)

		list := registry.List()
		if len(list) != i+1 {
			t.Errorf("List() length = %v, want %v", len(list), i+1)
		}
	}

	list := registry.List()
	if len(list) != 3 {
		t.Errorf("List() length = %v, want 3", len(list))
	}

	// Verify all names are present
	names := make(map[string]bool)
	for _, name := range list {
		names[name] = true
	}

	for _, expected := range []string{"client1", "client2", "client3"} {
		if !names[expected] {
			t.Errorf("List() missing expected name %v", expected)
		}
	}
}

func TestRegistry_Count(t *testing.T) {
	registry := NewRegistry()
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)

	// Empty registry
	if registry.Count() != 0 {
		t.Error("Count() should return 0 for empty registry")
	}

	// Add clients
	for i := 1; i <= 3; i++ {
		clientConfig := &ClientConfig{
			Name:     "client" + strconv.Itoa(i),
			AuthType: AuthTypeNone,
		}
		client, err := builder.Build(clientConfig)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		registry.Register(clientConfig.Name, client)

		if registry.Count() != i {
			t.Errorf("Count() = %v, want %v", registry.Count(), i)
		}
	}
}

func TestRegistry_Clear(t *testing.T) {
	registry := NewRegistry()
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)

	// Add clients
	for _, name := range []string{"client1", "client2", "client3"} {
		clientConfig := &ClientConfig{
			Name:     name,
			AuthType: AuthTypeNone,
		}
		client, err := builder.Build(clientConfig)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		registry.Register(name, client)
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
	for _, name := range []string{"client1", "client2", "client3"} {
		if registry.Has(name) {
			t.Errorf("Has(%v) returned true after clear", name)
		}
	}
}

func TestRegistry_Concurrent(t *testing.T) {
	registry := NewRegistry()
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)

	done := make(chan bool)

	// Concurrent registrations
	for i := 0; i < 10; i++ {
		go func(id int) {
			clientConfig := &ClientConfig{
				Name:     "client" + strconv.Itoa(id),
				AuthType: AuthTypeNone,
			}
			client, err := builder.Build(clientConfig)
			if err != nil {
				t.Errorf("Build() error = %v", err)
			}
			registry.Register(clientConfig.Name, client)
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
			name := "client" + strconv.Itoa(id)
			if !registry.Has(name) {
				t.Errorf("Has(%v) returned false", name)
			}
			_, err := registry.Get(name)
			if err != nil {
				t.Errorf("Get(%v) error = %v", name, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
