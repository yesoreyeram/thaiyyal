// Package httpclient provides a configurable HTTP client builder for the workflow engine.
//
// This package allows SDK consumers to define named HTTP clients with custom authentication,
// headers, timeouts, and other options via configuration. HTTP nodes in workflows can then
// reference these named clients to avoid repeating authentication and configuration.
//
// # Features
//
//   - Multiple named HTTP clients with independent configurations
//   - Authentication: None (default), Basic Auth, Bearer Token
//   - Configurable timeouts, connection pooling, and network settings
//   - Default headers and query parameters
//   - SSRF protection integrated with engine security settings
//   - Thread-safe client registry
//
// # Authentication Types
//
// The package supports three authentication types:
//
//   - None: No authentication (default)
//   - Basic: HTTP Basic Authentication with username and password
//   - Bearer: Bearer Token authentication
//
// Future extensions may include OAuth2 and other authentication mechanisms.
//
// # Example Usage
//
//	// Create client configuration
//	config := &httpclient.ClientConfig{
//	    Name:     "api-client",
//	    AuthType: httpclient.AuthTypeBearer,
//	    Token:    "your-api-token",
//	    Timeout:  60 * time.Second,
//	    DefaultHeaders: map[string]string{
//	        "Content-Type": "application/json",
//	    },
//	}
//
//	// Build the client
//	builder := httpclient.NewBuilder(engineConfig)
//	client, err := builder.Build(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Register in registry
//	registry := httpclient.NewRegistry()
//	registry.Register("api-client", client)
//
//	// Use in HTTP nodes
//	// In workflow JSON: {"type": "http", "data": {"url": "...", "client_name": "api-client"}}
//
// # Configuration File Format
//
// HTTP clients can be configured via YAML or JSON:
//
//	http_clients:
//	  - name: "github-api"
//	    description: "GitHub API client"
//	    auth_type: "bearer"
//	    token: "${GITHUB_TOKEN}"
//	    base_url: "https://api.github.com"
//	    timeout: "30s"
//	    default_headers:
//	      Accept: "application/vnd.github.v3+json"
//	      User-Agent: "Thaiyyal Workflow"
//
//	  - name: "internal-api"
//	    description: "Internal service API"
//	    auth_type: "basic"
//	    username: "service-account"
//	    password: "${API_PASSWORD}"
//	    timeout: "10s"
//	    max_redirects: 5
//
// # Security Considerations
//
//   - All clients inherit SSRF protection from the engine configuration
//   - Credentials should be passed via environment variables, not hardcoded
//   - Maximum response sizes are enforced to prevent memory exhaustion
//   - Redirect validation prevents redirect-based SSRF attacks
//   - Connection pooling limits prevent resource exhaustion
package httpclient
