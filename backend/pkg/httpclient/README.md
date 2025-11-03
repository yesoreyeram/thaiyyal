# HTTP Client Package

A standalone, zero-dependency HTTP client builder with middleware support for Go.

## Features

- **Standalone Package**: No dependencies on other project packages
- **Immutable UIDs**: Clients identified by unique, immutable UIDs
- **Middleware Pattern**: Composable roundtrippers for extensibility
- **Duplicate Keys Support**: Headers and query params use `[]KeyValue` structure
- **Multiple Auth Types**: None, Basic Auth, Bearer Token (OAuth2-ready)
- **SSRF Protection**: Built-in security against Server-Side Request Forgery
- **Thread-Safe Registry**: Concurrent-safe client management
- **Security First**: Configurable protection against private IPs, localhost, cloud metadata

## Installation

```bash
go get github.com/yesoreyeram/thaiyyal/backend/pkg/httpclient
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/yesoreyeram/thaiyyal/backend/pkg/httpclient"
)

func main() {
    // Create configuration
    config := &httpclient.Config{
        UID:      "github-api-client",
        AuthType: httpclient.AuthTypeBearer,
        Token:    "ghp_your_token_here",
        Timeout:  30 * time.Second,
        Headers: []httpclient.KeyValue{
            {Key: "Accept", Value: "application/vnd.github.v3+json"},
            {Key: "User-Agent", Value: "MyApp/1.0"},
        },
    }

    // Create HTTP client
    client, err := httpclient.New(context.Background(), config)
    if err != nil {
        log.Fatal(err)
    }

    // Use the client
    resp, err := client.Get("https://api.github.com/users/octocat")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Status:", resp.Status)
}
```

### Using Registry

```go
// Create registry
registry := httpclient.NewRegistry()

// Build and register clients
configs := []*httpclient.Config{
    {
        UID:      "github-api",
        AuthType: httpclient.AuthTypeBearer,
        Token:    "token1",
    },
    {
        UID:      "internal-api",
        AuthType: httpclient.AuthTypeBasic,
        Username: "user",
        Password: "pass",
    },
}

for _, cfg := range configs {
    client, err := httpclient.New(context.Background(), cfg)
    if err != nil {
        log.Fatal(err)
    }
    registry.Register(cfg.UID, client)
}

// Retrieve and use clients
client, err := registry.Get("github-api")
if err != nil {
    log.Fatal(err)
}
resp, _ := client.Get("https://api.github.com/repos/owner/repo")
```

## Configuration Options

### Authentication

```go
// No authentication
config := &httpclient.Config{
    UID:      "public-api",
    AuthType: httpclient.AuthTypeNone,
}

// Basic authentication
config := &httpclient.Config{
    UID:      "basic-auth-api",
    AuthType: httpclient.AuthTypeBasic,
    Username: "user",
    Password: "secret",
}

// Bearer token
config := &httpclient.Config{
    UID:      "bearer-api",
    AuthType: httpclient.AuthTypeBearer,
    Token:    "your-token-here",
}
```

### Headers and Query Parameters

```go
config := &httpclient.Config{
    UID: "custom-client",
    
    // Headers support duplicate keys
    Headers: []httpclient.KeyValue{
        {Key: "X-Custom", Value: "value1"},
        {Key: "X-Custom", Value: "value2"}, // Duplicate key
        {Key: "User-Agent", Value: "MyApp/1.0"},
    },
    
    // Query params support duplicate keys
    QueryParams: []httpclient.KeyValue{
        {Key: "api_key", Value: "secret"},
        {Key: "format", Value: "json"},
        {Key: "tag", Value: "v1"},
        {Key: "tag", Value: "v2"}, // Duplicate key
    },
}
```

### Network Configuration

```go
config := &httpclient.Config{
    UID:                 "optimized-client",
    Timeout:             60 * time.Second,
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    MaxConnsPerHost:     100,
    IdleConnTimeout:     90 * time.Second,
    TLSHandshakeTimeout: 10 * time.Second,
    DisableKeepAlives:   false,
}
```

### Security Configuration

```go
config := &httpclient.Config{
    UID: "secure-client",
    
    // SSRF Protection
    BlockPrivateIPs:    true,
    BlockLocalhost:     true,
    BlockLinkLocal:     true,
    BlockCloudMetadata: true,
    AllowedDomains:     []string{"api.example.com", "data.example.com"},
    
    // Response limits
    MaxResponseSize:    10 * 1024 * 1024, // 10MB
    
    // Redirect handling
    FollowRedirects: true,
    MaxRedirects:    5,
}
```

## Architecture

### Middleware Pattern

The package uses the middleware/roundtripper pattern for composability:

```
Request
    ↓
[Auth Middleware] → Adds authentication headers
    ↓
[Headers Middleware] → Adds default headers
    ↓
[Query Params Middleware] → Adds default query params
    ↓
[SSRF Protection Middleware] → Validates URL
    ↓
[Base Transport] → Makes actual HTTP request
    ↓
Response
```

### Components

- **Config**: Configuration structure with validation
- **Client Builder**: Creates HTTP clients from configuration
- **Middlewares**: Composable request/response interceptors
- **Registry**: Thread-safe client management
- **SSRF Protection**: Security validation for URLs and IPs

## API Reference

### Config

```go
type Config struct {
    UID                 string        // Unique immutable identifier (required)
    Description         string        // Human-readable description
    AuthType            AuthType      // "none", "basic", or "bearer"
    Username            string        // For basic auth
    Password            string        // For basic auth
    Token               string        // For bearer token
    Timeout             time.Duration // Request timeout (default: 30s)
    MaxIdleConns        int          // Max idle connections (default: 100)
    MaxIdleConnsPerHost int          // Max idle conns per host (default: 10)
    MaxConnsPerHost     int          // Max conns per host (default: 100)
    IdleConnTimeout     time.Duration // Idle conn timeout (default: 90s)
    TLSHandshakeTimeout time.Duration // TLS timeout (default: 10s)
    DisableKeepAlives   bool         // Disable keep-alives (default: false)
    MaxRedirects        int          // Max redirects (default: 10)
    MaxResponseSize     int64        // Max response size (default: 10MB)
    FollowRedirects     bool         // Follow redirects (default: true)
    BlockPrivateIPs     bool         // Block private IPs
    BlockLocalhost      bool         // Block localhost
    BlockLinkLocal      bool         // Block link-local addresses
    BlockCloudMetadata  bool         // Block cloud metadata endpoints
    AllowedDomains      []string     // Whitelist of allowed domains
    Headers             []KeyValue   // Default headers
    QueryParams         []KeyValue   // Default query parameters
    BaseURL             string       // Base URL for requests
}

func (c *Config) Validate() error
func (c *Config) ApplyDefaults()
func (c *Config) Clone() *Config
```

### Client Builder

```go
func New(ctx context.Context, config *Config) (*http.Client, error)
```

Creates a new HTTP client with the given configuration. The context parameter is for future extensibility.

### Registry

```go
type Registry struct { /* ... */ }

func NewRegistry() *Registry
func (r *Registry) Register(uid string, client *http.Client) error
func (r *Registry) Get(uid string) (*http.Client, error)
func (r *Registry) Has(uid string) bool
func (r *Registry) List() []string
func (r *Registry) Count() int
func (r *Registry) Clear()
func (r *Registry) Unregister(uid string) error
```

## Security

### SSRF Protection

The package includes built-in SSRF (Server-Side Request Forgery) protection:

- **Private IP blocking**: Blocks 10.x, 172.16.x, 192.168.x ranges
- **Localhost blocking**: Blocks loopback addresses
- **Link-local blocking**: Blocks 169.254.x.x addresses
- **Cloud metadata blocking**: Blocks AWS/Azure/GCP metadata endpoints
- **Domain whitelisting**: Restrict requests to specific domains

```go
config := &httpclient.Config{
    UID:                "secure-client",
    BlockPrivateIPs:    true,
    BlockLocalhost:     true,
    BlockCloudMetadata: true,
    AllowedDomains:     []string{"trusted.com"},
}
```

### Best Practices

1. **Use environment variables** for sensitive data:
   ```go
   Token: os.Getenv("API_TOKEN"),
   ```

2. **Enable SSRF protection** in production:
   ```go
   BlockPrivateIPs:    true,
   BlockLocalhost:     true,
   BlockCloudMetadata: true,
   ```

3. **Set response size limits**:
   ```go
   MaxResponseSize: 10 * 1024 * 1024, // 10MB
   ```

4. **Use domain whitelisting** when possible:
   ```go
   AllowedDomains: []string{"api.trusted.com"},
   ```

## Examples

See the [examples directory](../../examples/httpclient_standalone/) for complete working examples:

- Basic usage
- Multiple authentication types
- Registry management
- SSRF protection demo
- Duplicate headers/params

## Testing

Run tests:
```bash
go test ./pkg/httpclient/...
```

Run benchmarks:
```bash
go test -bench=. ./pkg/httpclient/...
```

## License

MIT License - See repository LICENSE file for details.
