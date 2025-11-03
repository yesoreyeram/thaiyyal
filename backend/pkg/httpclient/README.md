# HTTP Client Builder Package

The `httpclient` package provides a configurable HTTP client builder system for the Thaiyyal workflow engine. It allows SDK consumers to define named HTTP clients with custom authentication, headers, timeouts, and other options via configuration.

## Features

- 🔐 **Multiple Authentication Types**: None, Basic Auth, Bearer Token (extensible for OAuth2)
- 🏷️ **Named Clients**: Reference pre-configured clients by name in HTTP nodes
- ⚙️ **Highly Configurable**: Timeouts, connection pooling, headers, query params
- 🛡️ **Security First**: Integrated SSRF protection, validates URLs against engine security settings
- 🔄 **Thread-Safe**: Concurrent access to client registry is fully synchronized
- ⬅️ **Backward Compatible**: Works seamlessly with existing HTTP nodes

## Quick Start

### 1. Define HTTP Clients in Configuration

```yaml
http_clients:
  - name: "github-api"
    description: "GitHub API client"
    auth_type: "bearer"
    token: "${GITHUB_TOKEN}"
    timeout: "30s"
    default_headers:
      Accept: "application/vnd.github.v3+json"
      User-Agent: "MyApp/1.0"
  
  - name: "internal-api"
    description: "Internal service"
    auth_type: "basic"
    username: "${API_USERNAME}"
    password: "${API_PASSWORD}"
    timeout: "15s"
```

### 2. Build and Register Clients

```go
package main

import (
    "time"
    
    "github.com/yesoreyeram/thaiyyal/backend/pkg/config"
    "github.com/yesoreyeram/thaiyyal/backend/pkg/httpclient"
    "github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
)

func main() {
    // Load configuration
    cfg := config.Default()
    cfg.HTTPClients = []config.HTTPClientConfig{
        {
            Name:     "github-api",
            AuthType: "bearer",
            Token:    "ghp_your_token_here",
            Timeout:  30 * time.Second,
            DefaultHeaders: map[string]string{
                "Accept": "application/vnd.github.v3+json",
            },
        },
    }
    
    // Build HTTP client registry
    builder := httpclient.NewBuilder(*cfg)
    registry := httpclient.NewRegistry()
    
    for _, clientConfig := range cfg.HTTPClients {
        httpClient := httpclient.FromConfigHTTPClient(clientConfig)
        client, err := builder.Build(httpClient)
        if err != nil {
            panic(err)
        }
        registry.Register(clientConfig.Name, client)
    }
    
    // Create engine and set registry
    payload := []byte(`{
        "nodes": [{
            "id": "1",
            "type": "http",
            "data": {
                "url": "https://api.github.com/users/octocat",
                "client_name": "github-api"
            }
        }],
        "edges": []
    }`)
    
    engine, _ := engine.NewWithConfig(payload, *cfg)
    engine.SetHTTPClientRegistry(registry)
    
    result, _ := engine.Execute()
    println(result.FinalOutput)
}
```

### 3. Use Named Clients in Workflows

```json
{
  "nodes": [
    {
      "id": "fetch-data",
      "type": "http",
      "data": {
        "url": "https://api.example.com/data",
        "client_name": "github-api"
      }
    }
  ],
  "edges": []
}
```

## Authentication Types

### None (Default)

No authentication headers are added:

```yaml
- name: "public-api"
  auth_type: "none"
  timeout: "10s"
```

### Basic Authentication

Adds HTTP Basic Auth headers:

```yaml
- name: "protected-api"
  auth_type: "basic"
  username: "user"
  password: "secret"
```

### Bearer Token

Adds `Authorization: Bearer <token>` header:

```yaml
- name: "api-with-token"
  auth_type: "bearer"
  token: "your-secret-token"
```

## Configuration Options

### Network Settings

```yaml
- name: "custom-network"
  timeout: "60s"                  # Request timeout
  max_idle_conns: 100             # Max idle connections
  max_idle_conns_per_host: 10     # Max idle connections per host
  max_conns_per_host: 100         # Max connections per host
  idle_conn_timeout: "90s"        # Idle connection timeout
  tls_handshake_timeout: "10s"    # TLS handshake timeout
  disable_keep_alives: false      # Enable keep-alives
```

### Security Settings

```yaml
- name: "secure-api"
  max_redirects: 5                # Maximum redirects to follow
  max_response_size: 10485760     # Max response size (10MB)
  follow_redirects: true          # Follow HTTP redirects
```

### Default Headers

```yaml
- name: "custom-headers"
  default_headers:
    Content-Type: "application/json"
    User-Agent: "MyApp/1.0"
    X-Custom-Header: "custom-value"
```

### Default Query Parameters

```yaml
- name: "api-with-params"
  default_query_params:
    api_key: "secret-key"
    format: "json"
    version: "v2"
```

## API Reference

### Builder

```go
type Builder struct {
    // Creates configured HTTP clients
}

func NewBuilder(engineConfig types.Config) *Builder

func (b *Builder) Build(config *ClientConfig) (*Client, error)
```

### Registry

```go
type Registry struct {
    // Manages named HTTP clients
}

func NewRegistry() *Registry

func (r *Registry) Register(name string, client *Client) error
func (r *Registry) Get(name string) (*Client, error)
func (r *Registry) GetHTTPClient(name string) (*http.Client, int64, error)
func (r *Registry) Has(name string) bool
func (r *Registry) List() []string
func (r *Registry) Count() int
func (r *Registry) Clear()
```

### ClientConfig

```go
type ClientConfig struct {
    Name                string
    Description         string
    AuthType            AuthType  // "none", "basic", "bearer"
    Username            string    // For basic auth
    Password            string    // For basic auth
    Token               string    // For bearer token
    Timeout             time.Duration
    MaxIdleConns        int
    MaxIdleConnsPerHost int
    MaxConnsPerHost     int
    IdleConnTimeout     time.Duration
    TLSHandshakeTimeout time.Duration
    DisableKeepAlives   bool
    MaxRedirects        int
    MaxResponseSize     int64
    FollowRedirects     bool
    DefaultHeaders      map[string]string
    DefaultQueryParams  map[string]string
    BaseURL             string
}

func (c *ClientConfig) Validate() error
func (c *ClientConfig) ApplyDefaults()
func (c *ClientConfig) Clone() *ClientConfig
```

## Security

### SSRF Protection

All HTTP clients inherit SSRF protection from the engine configuration:

- URL validation against private IPs
- Blocking of localhost and link-local addresses
- Cloud metadata endpoint protection
- Domain whitelisting support
- Redirect validation

### Credential Management

**Best Practices:**

- ✅ Use environment variables for sensitive data
- ✅ Load credentials from secure vaults
- ❌ Never hardcode credentials in configuration files
- ❌ Never commit credentials to version control

Example using environment variables:

```yaml
http_clients:
  - name: "secure-api"
    auth_type: "bearer"
    token: "${API_TOKEN}"  # References $API_TOKEN environment variable
```

### Response Size Limits

Clients enforce maximum response sizes to prevent memory exhaustion:

```yaml
- name: "limited-api"
  max_response_size: 5242880  # 5MB limit
```

## Examples

### GitHub API Integration

```go
cfg := config.Default()
cfg.HTTPClients = []config.HTTPClientConfig{
    {
        Name:     "github",
        AuthType: "bearer",
        Token:    os.Getenv("GITHUB_TOKEN"),
        Timeout:  30 * time.Second,
        DefaultHeaders: map[string]string{
            "Accept": "application/vnd.github.v3+json",
        },
    },
}

builder := httpclient.NewBuilder(*cfg)
registry := httpclient.NewRegistry()

client, _ := builder.Build(httpclient.FromConfigHTTPClient(cfg.HTTPClients[0]))
registry.Register("github", client)

engine.SetHTTPClientRegistry(registry)
```

### Multiple APIs with Different Auth

```go
cfg.HTTPClients = []config.HTTPClientConfig{
    {
        Name:     "api-basic",
        AuthType: "basic",
        Username: "user",
        Password: "pass",
    },
    {
        Name:     "api-bearer",
        AuthType: "bearer",
        Token:    "token123",
    },
    {
        Name:     "api-key",
        AuthType: "none",
        DefaultHeaders: map[string]string{
            "X-API-Key": "secret-key",
        },
    },
}
```

## Testing

### Unit Tests

```bash
go test ./pkg/httpclient -v
```

### Integration Tests

```bash
go test ./pkg/httpclient -v -run TestNamedHTTPClient_Integration
```

## Future Enhancements

- [ ] OAuth2 authentication support
- [ ] Client certificate authentication
- [ ] Request/response interceptors
- [ ] Retry policies per client
- [ ] Circuit breaker pattern
- [ ] Metrics and observability hooks
- [ ] Dynamic client registration
- [ ] Client health checks

## Contributing

When adding new features:

1. Maintain backward compatibility
2. Add comprehensive tests
3. Update documentation
4. Follow existing code patterns
5. Ensure security best practices

## License

MIT License - See repository LICENSE file for details.
