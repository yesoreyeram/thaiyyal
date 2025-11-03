// HTTP Client Builder Example
//
// This example demonstrates how to use named HTTP clients in workflows.
// It shows different authentication types and configurations.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/config"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/httpclient"
)

func main() {
	// Create engine configuration with named HTTP clients
	cfg := config.Development() // Use development config for demo
	cfg.HTTPClients = []config.HTTPClientConfig{
		{
			Name:        "github-api",
			Description: "GitHub API client with bearer token",
			AuthType:    "bearer",
			Token:       "ghp_your_github_token_here", // Replace with actual token
			Timeout:     30 * time.Second,
			DefaultHeaders: map[string]string{
				"Accept":     "application/vnd.github.v3+json",
				"User-Agent": "Thaiyyal-Example/1.0",
			},
		},
		{
			Name:        "httpbin-basic",
			Description: "HTTPBin with basic auth for testing",
			AuthType:    "basic",
			Username:    "user",
			Password:    "passwd",
			Timeout:     15 * time.Second,
		},
		{
			Name:        "public-api",
			Description: "Public API without authentication",
			AuthType:    "none",
			Timeout:     10 * time.Second,
			DefaultHeaders: map[string]string{
				"User-Agent": "Thaiyyal-Example/1.0",
			},
		},
	}

	// Build HTTP client registry
	fmt.Println("Building HTTP client registry...")
	builder := httpclient.NewBuilder(*cfg)
	registry := httpclient.NewRegistry()

	for _, clientCfg := range cfg.HTTPClients {
		httpClientCfg := httpclient.FromConfigHTTPClient(clientCfg)
		client, err := builder.Build(httpClientCfg)
		if err != nil {
			log.Fatalf("Failed to build HTTP client %q: %v", clientCfg.Name, err)
		}

		if err := registry.Register(clientCfg.Name, client); err != nil {
			log.Fatalf("Failed to register HTTP client %q: %v", clientCfg.Name, err)
		}

		fmt.Printf("  ✓ Registered client: %s\n", clientCfg.Name)
	}

	// Example 1: Workflow using GitHub API client
	fmt.Println("\nExample 1: Using GitHub API client")
	githubWorkflow := `{
		"workflow_id": "github-example",
		"nodes": [
			{
				"id": "fetch-user",
				"type": "http",
				"data": {
					"url": "https://api.github.com/users/octocat",
					"client_name": "github-api"
				}
			}
		],
		"edges": []
	}`

	executeWorkflow(githubWorkflow, cfg, registry)

	// Example 2: Workflow using basic auth
	fmt.Println("\nExample 2: Using HTTPBin with basic auth")
	basicAuthWorkflow := `{
		"workflow_id": "basic-auth-example",
		"nodes": [
			{
				"id": "test-auth",
				"type": "http",
				"data": {
					"url": "https://httpbin.org/basic-auth/user/passwd",
					"client_name": "httpbin-basic"
				}
			}
		],
		"edges": []
	}`

	executeWorkflow(basicAuthWorkflow, cfg, registry)

	// Example 3: Workflow using public API (no auth)
	fmt.Println("\nExample 3: Using public API without authentication")
	publicWorkflow := `{
		"workflow_id": "public-api-example",
		"nodes": [
			{
				"id": "get-ip",
				"type": "http",
				"data": {
					"url": "https://httpbin.org/ip",
					"client_name": "public-api"
				}
			}
		],
		"edges": []
	}`

	executeWorkflow(publicWorkflow, cfg, registry)

	// Example 4: Workflow without named client (uses default)
	fmt.Println("\nExample 4: Using default HTTP client (no client_name)")
	defaultWorkflow := `{
		"workflow_id": "default-client-example",
		"nodes": [
			{
				"id": "get-uuid",
				"type": "http",
				"data": {
					"url": "https://httpbin.org/uuid"
				}
			}
		],
		"edges": []
	}`

	executeWorkflow(defaultWorkflow, cfg, registry)

	// List all registered clients
	fmt.Println("\nRegistered HTTP clients:")
	for _, name := range registry.List() {
		client, _ := registry.Get(name)
		cfg := client.GetConfig()
		fmt.Printf("  - %s: auth=%s, timeout=%v\n", name, cfg.AuthType, cfg.Timeout)
	}
}

func executeWorkflow(payload string, cfg *config.Config, registry *httpclient.Registry) {
	// Create engine
	eng, err := engine.NewWithConfig([]byte(payload), *cfg)
	if err != nil {
		log.Printf("  ✗ Failed to create engine: %v\n", err)
		return
	}

	// Set HTTP client registry
	eng.SetHTTPClientRegistry(registry)

	// Execute workflow
	result, err := eng.Execute()
	if err != nil {
		log.Printf("  ✗ Workflow execution failed: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Workflow ID: %s\n", result.WorkflowID)
	fmt.Printf("  ✓ Execution ID: %s\n", result.ExecutionID)
	fmt.Printf("  ✓ Result: %v\n", truncate(fmt.Sprintf("%v", result.FinalOutput), 100))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
