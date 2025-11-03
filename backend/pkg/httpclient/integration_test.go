package httpclient_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/config"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/httpclient"
)

// TestNamedHTTPClient_Integration tests the complete flow of using named HTTP clients in workflows
func TestNamedHTTPClient_Integration(t *testing.T) {
	// Create test servers
	basicAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "testuser" || password != "testpass" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated with basic auth"))
	}))
	defer basicAuthServer.Close()

	bearerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer secret-token-123" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated with bearer token"))
	}))
	defer bearerServer.Close()

	customHeaderServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "my-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("missing api key"))
			return
		}
		if r.Header.Get("User-Agent") != "MyApp/1.0" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid user agent"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("custom headers validated"))
	}))
	defer customHeaderServer.Close()

	// Create engine config with named HTTP clients
	engineConfig := config.Testing()
	engineConfig.HTTPClients = []config.HTTPClientConfig{
		{
			Name:        "basic-auth-client",
			Description: "Client with basic authentication",
			AuthType:    "basic",
			Username:    "testuser",
			Password:    "testpass",
			Timeout:     30 * time.Second,
		},
		{
			Name:        "bearer-token-client",
			Description: "Client with bearer token",
			AuthType:    "bearer",
			Token:       "secret-token-123",
			Timeout:     30 * time.Second,
		},
		{
			Name:        "custom-headers-client",
			Description: "Client with custom headers",
			AuthType:    "none",
			Timeout:     30 * time.Second,
			DefaultHeaders: map[string]string{
				"X-API-Key":  "my-api-key",
				"User-Agent": "MyApp/1.0",
			},
		},
	}

	// Build HTTP client registry
	builder := httpclient.NewBuilder(*engineConfig)
	registry := httpclient.NewRegistry()

	for _, clientConfig := range engineConfig.HTTPClients {
		httpClientConfig := httpclient.FromConfigHTTPClient(clientConfig)
		client, err := builder.Build(httpClientConfig)
		if err != nil {
			t.Fatalf("Failed to build HTTP client %q: %v", clientConfig.Name, err)
		}
		if err := registry.Register(clientConfig.Name, client); err != nil {
			t.Fatalf("Failed to register HTTP client %q: %v", clientConfig.Name, err)
		}
	}

	t.Run("basic auth client", func(t *testing.T) {
		clientName := "basic-auth-client"
		payload := createHTTPWorkflow(basicAuthServer.URL, &clientName)

		eng, err := engine.NewWithConfig([]byte(payload), *engineConfig)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		// Set HTTP client registry
		eng.SetHTTPClientRegistry(registry)

		result, err := eng.Execute()
		if err != nil {
			t.Fatalf("Workflow execution failed: %v", err)
		}

		if result.FinalOutput != "authenticated with basic auth" {
			t.Errorf("Expected 'authenticated with basic auth', got %v", result.FinalOutput)
		}
	})

	t.Run("bearer token client", func(t *testing.T) {
		clientName := "bearer-token-client"
		payload := createHTTPWorkflow(bearerServer.URL, &clientName)

		eng, err := engine.NewWithConfig([]byte(payload), *engineConfig)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		eng.SetHTTPClientRegistry(registry)

		result, err := eng.Execute()
		if err != nil {
			t.Fatalf("Workflow execution failed: %v", err)
		}

		if result.FinalOutput != "authenticated with bearer token" {
			t.Errorf("Expected 'authenticated with bearer token', got %v", result.FinalOutput)
		}
	})

	t.Run("custom headers client", func(t *testing.T) {
		clientName := "custom-headers-client"
		payload := createHTTPWorkflow(customHeaderServer.URL, &clientName)

		eng, err := engine.NewWithConfig([]byte(payload), *engineConfig)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		eng.SetHTTPClientRegistry(registry)

		result, err := eng.Execute()
		if err != nil {
			t.Fatalf("Workflow execution failed: %v", err)
		}

		if result.FinalOutput != "custom headers validated" {
			t.Errorf("Expected 'custom headers validated', got %v", result.FinalOutput)
		}
	})

	t.Run("default client (no client name)", func(t *testing.T) {
		// Create a simple test server that doesn't require auth
		simpleServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("default client response"))
		}))
		defer simpleServer.Close()

		payload := createHTTPWorkflow(simpleServer.URL, nil)

		eng, err := engine.NewWithConfig([]byte(payload), *engineConfig)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		eng.SetHTTPClientRegistry(registry)

		result, err := eng.Execute()
		if err != nil {
			t.Fatalf("Workflow execution failed: %v", err)
		}

		if result.FinalOutput != "default client response" {
			t.Errorf("Expected 'default client response', got %v", result.FinalOutput)
		}
	})

	t.Run("non-existent client", func(t *testing.T) {
		clientName := "non-existent-client"
		payload := createHTTPWorkflow(basicAuthServer.URL, &clientName)

		eng, err := engine.NewWithConfig([]byte(payload), *engineConfig)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		eng.SetHTTPClientRegistry(registry)

		_, err = eng.Execute()
		if err == nil {
			t.Error("Expected error for non-existent client, got nil")
		}
	})

	t.Run("no registry configured", func(t *testing.T) {
		clientName := "basic-auth-client"
		payload := createHTTPWorkflow(basicAuthServer.URL, &clientName)

		eng, err := engine.NewWithConfig([]byte(payload), *engineConfig)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		// Don't set registry - should fail
		_, err = eng.Execute()
		if err == nil {
			t.Error("Expected error when registry not configured, got nil")
		}
	})
}

// createHTTPWorkflow creates a simple workflow with an HTTP node
func createHTTPWorkflow(url string, clientName *string) string {
	clientNameStr := ""
	if clientName != nil {
		clientNameStr = `,"client_name":"` + *clientName + `"`
	}

	return `{
		"nodes": [
			{
				"id": "http-1",
				"type": "http",
				"data": {
					"url": "` + url + `"` + clientNameStr + `
				}
			}
		],
		"edges": []
	}`
}

// TestHTTPClientConfig_FromConfig tests the conversion from config.HTTPClientConfig
func TestHTTPClientConfig_FromConfig(t *testing.T) {
	configClient := config.HTTPClientConfig{
		Name:                "test-client",
		Description:         "Test client",
		AuthType:            "basic",
		Username:            "user",
		Password:            "pass",
		Timeout:             60 * time.Second,
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 5,
		MaxConnsPerHost:     50,
		IdleConnTimeout:     120 * time.Second,
		TLSHandshakeTimeout: 15 * time.Second,
		DisableKeepAlives:   true,
		MaxRedirects:        5,
		MaxResponseSize:     5 * 1024 * 1024,
		FollowRedirects:     false,
		DefaultHeaders: map[string]string{
			"X-Custom": "value",
		},
		DefaultQueryParams: map[string]string{
			"api_key": "secret",
		},
		BaseURL: "https://api.example.com",
	}

	httpClient := httpclient.FromConfigHTTPClient(configClient)

	if httpClient.Name != configClient.Name {
		t.Errorf("Name = %v, want %v", httpClient.Name, configClient.Name)
	}
	if httpClient.Description != configClient.Description {
		t.Errorf("Description = %v, want %v", httpClient.Description, configClient.Description)
	}
	if string(httpClient.AuthType) != configClient.AuthType {
		t.Errorf("AuthType = %v, want %v", httpClient.AuthType, configClient.AuthType)
	}
	if httpClient.Username != configClient.Username {
		t.Errorf("Username = %v, want %v", httpClient.Username, configClient.Username)
	}
	if httpClient.Password != configClient.Password {
		t.Errorf("Password = %v, want %v", httpClient.Password, configClient.Password)
	}
	if httpClient.Timeout != configClient.Timeout {
		t.Errorf("Timeout = %v, want %v", httpClient.Timeout, configClient.Timeout)
	}
	if httpClient.MaxIdleConns != configClient.MaxIdleConns {
		t.Errorf("MaxIdleConns = %v, want %v", httpClient.MaxIdleConns, configClient.MaxIdleConns)
	}
	if httpClient.BaseURL != configClient.BaseURL {
		t.Errorf("BaseURL = %v, want %v", httpClient.BaseURL, configClient.BaseURL)
	}

	// Verify maps are copied correctly
	if httpClient.DefaultHeaders["X-Custom"] != "value" {
		t.Error("DefaultHeaders not copied correctly")
	}
	if httpClient.DefaultQueryParams["api_key"] != "secret" {
		t.Error("DefaultQueryParams not copied correctly")
	}
}
