// Standalone HTTP Client Example
//
// This example demonstrates the new standalone httpclient package with:
// - Immutable UID-based client identification
// - Middleware pattern for composability
// - KeyValue structure for duplicate headers/params
// - Zero dependencies on other packages

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/httpclient"
)

func main() {
	fmt.Println("=== Standalone HTTP Client Examples ===\n")

	// Example 1: Basic Usage with Bearer Token
	example1_BearerAuth()

	// Example 2: Basic Authentication
	example2_BasicAuth()

	// Example 3: Custom Headers (including duplicates)
	example3_CustomHeaders()

	// Example 4: Query Parameters (including duplicates)
	example4_QueryParams()

	// Example 5: SSRF Protection
	example5_SSRFProtection()

	// Example 6: Using Registry
	example6_Registry()
}

func example1_BearerAuth() {
	fmt.Println("Example 1: Bearer Token Authentication")
	fmt.Println("---------------------------------------")

	// Create a test server that checks for bearer token
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "Bearer secret-token-123" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "authenticated", "user": "test"}`))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized"}`))
		}
	}))
	defer server.Close()

	// Create HTTP client with bearer token
	config := &httpclient.Config{
		UID:      "bearer-api-client",
		AuthType: httpclient.AuthTypeBearer,
		Token:    "secret-token-123",
		Timeout:  10 * time.Second,
	}

	client, err := httpclient.New(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Make request
	resp, err := client.Get(server.URL + "/api/user")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✓ UID: %s\n", config.UID)
	fmt.Printf("✓ Status: %s\n", resp.Status)
	fmt.Printf("✓ Response: %s\n\n", string(body))
}

func example2_BasicAuth() {
	fmt.Println("Example 2: Basic Authentication")
	fmt.Println("--------------------------------")

	// Create a test server that checks for basic auth
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok && username == "admin" && password == "secret" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "authenticated", "user": "admin"}`))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized"}`))
		}
	}))
	defer server.Close()

	// Create HTTP client with basic auth
	config := &httpclient.Config{
		UID:      "basic-auth-client",
		AuthType: httpclient.AuthTypeBasic,
		Username: "admin",
		Password: "secret",
		Timeout:  10 * time.Second,
	}

	client, err := httpclient.New(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Make request
	resp, err := client.Get(server.URL + "/api/protected")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✓ UID: %s\n", config.UID)
	fmt.Printf("✓ Status: %s\n", resp.Status)
	fmt.Printf("✓ Response: %s\n\n", string(body))
}

func example3_CustomHeaders() {
	fmt.Println("Example 3: Custom Headers (with duplicates)")
	fmt.Println("--------------------------------------------")

	// Create a test server that checks for custom headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Values("Accept")
		custom := r.Header.Values("X-Custom-Header")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"accept_headers": %d,
			"custom_headers": %d,
			"user_agent": "%s"
		}`, len(accept), len(custom), r.Header.Get("User-Agent"))))
	}))
	defer server.Close()

	// Create HTTP client with custom headers (including duplicates)
	config := &httpclient.Config{
		UID: "custom-headers-client",
		Headers: []httpclient.KeyValue{
			{Key: "Accept", Value: "application/json"},
			{Key: "Accept", Value: "application/xml"},  // Duplicate key
			{Key: "X-Custom-Header", Value: "value1"},
			{Key: "X-Custom-Header", Value: "value2"}, // Duplicate key
			{Key: "User-Agent", Value: "MyApp/1.0"},
		},
		Timeout: 10 * time.Second,
	}

	client, err := httpclient.New(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Make request
	resp, err := client.Get(server.URL)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✓ UID: %s\n", config.UID)
	fmt.Printf("✓ Headers configured: %d\n", len(config.Headers))
	fmt.Printf("✓ Response: %s\n\n", string(body))
}

func example4_QueryParams() {
	fmt.Println("Example 4: Query Parameters (with duplicates)")
	fmt.Println("----------------------------------------------")

	// Create a test server that checks for query params
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags := r.URL.Query()["tag"]
		apiKey := r.URL.Query().Get("api_key")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"api_key": "%s",
			"tags_count": %d,
			"tags": %v
		}`, apiKey, len(tags), tags)))
	}))
	defer server.Close()

	// Create HTTP client with query parameters (including duplicates)
	config := &httpclient.Config{
		UID: "query-params-client",
		QueryParams: []httpclient.KeyValue{
			{Key: "api_key", Value: "secret123"},
			{Key: "format", Value: "json"},
			{Key: "tag", Value: "important"},
			{Key: "tag", Value: "urgent"}, // Duplicate key
		},
		Timeout: 10 * time.Second,
	}

	client, err := httpclient.New(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Make request
	resp, err := client.Get(server.URL + "/api/search")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✓ UID: %s\n", config.UID)
	fmt.Printf("✓ Query params configured: %d\n", len(config.QueryParams))
	fmt.Printf("✓ Response: %s\n\n", string(body))
}

func example5_SSRFProtection() {
	fmt.Println("Example 5: SSRF Protection")
	fmt.Println("--------------------------")

	// Create HTTP client with SSRF protection
	config := &httpclient.Config{
		UID:                "secure-client",
		BlockPrivateIPs:    true,
		BlockLocalhost:     true,
		BlockCloudMetadata: true,
		AllowedDomains:     []string{"api.github.com"},
		Timeout:            10 * time.Second,
	}

	client, err := httpclient.New(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Try to access localhost (should be blocked)
	fmt.Println("Attempting to access localhost (should be blocked)...")
	_, err = client.Get("http://localhost:8080")
	if err != nil {
		fmt.Printf("✓ Blocked as expected: %v\n", err)
	} else {
		fmt.Println("✗ Should have been blocked!")
	}

	// Try to access allowed domain
	fmt.Println("\nAttempting to access allowed domain...")
	resp, err := client.Get("https://api.github.com")
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		fmt.Printf("✓ Allowed: Status %s\n", resp.Status)
		resp.Body.Close()
	}

	// Try to access non-allowed domain
	fmt.Println("\nAttempting to access non-allowed domain (should be blocked)...")
	_, err = client.Get("https://example.com")
	if err != nil {
		fmt.Printf("✓ Blocked as expected: %v\n\n", err)
	} else {
		fmt.Println("✗ Should have been blocked!\n")
	}
}

func example6_Registry() {
	fmt.Println("Example 6: Using Registry for Multiple Clients")
	fmt.Println("-----------------------------------------------")

	// Create registry
	registry := httpclient.NewRegistry()

	// Create multiple clients
	clients := []struct {
		uid      string
		authType httpclient.AuthType
		token    string
	}{
		{"github-api", httpclient.AuthTypeBearer, "ghp_token1"},
		{"internal-api", httpclient.AuthTypeBearer, "internal_token"},
		{"public-api", httpclient.AuthTypeNone, ""},
	}

	for _, c := range clients {
		config := &httpclient.Config{
			UID:      c.uid,
			AuthType: c.authType,
			Token:    c.token,
			Timeout:  10 * time.Second,
		}

		client, err := httpclient.New(context.Background(), config)
		if err != nil {
			log.Fatalf("Failed to create client %s: %v", c.uid, err)
		}

		if err := registry.Register(c.uid, client); err != nil {
			log.Fatalf("Failed to register client %s: %v", c.uid, err)
		}

		fmt.Printf("✓ Registered: %s (auth: %s)\n", c.uid, c.authType)
	}

	// List all registered clients
	fmt.Printf("\n✓ Total clients in registry: %d\n", registry.Count())
	fmt.Printf("✓ Client UIDs: %v\n\n", registry.List())

	// Retrieve and use a specific client
	client, err := registry.Get("github-api")
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}

	fmt.Printf("✓ Retrieved client: github-api\n")
	fmt.Printf("✓ Client type: %T\n", client)
}
