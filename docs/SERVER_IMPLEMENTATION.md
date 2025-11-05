# Server Implementation Summary

This document summarizes the server package implementation for the Thaiyyal workflow engine.

## Overview

The server package provides a complete HTTP API server that:
1. Executes and validates workflows
2. Manages workflow storage (save, load, list, delete, execute by ID)
3. Manages HTTP client configurations
4. Serves the frontend application via embedded files
5. Provides health checks and metrics

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Server (:8080)                        │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Frontend (/)                                            │
│  ├── index.html                                          │
│  ├── workflow.html                                       │
│  └── _next/* (static assets)                            │
│                                                          │
│  API Routes                                              │
│  ├── /api/v1/workflow/*                                 │
│  │   ├── execute (POST) - Execute workflow              │
│  │   ├── validate (POST) - Validate workflow            │
│  │   ├── save (POST) - Save workflow                    │
│  │   ├── list (GET) - List workflows                    │
│  │   ├── load/{id} (GET) - Load workflow                │
│  │   ├── delete/{id} (DELETE) - Delete workflow         │
│  │   └── execute/{id} (POST) - Execute by ID            │
│  │                                                       │
│  └── /api/v1/httpclient/*                               │
│      ├── register (POST) - Register HTTP client         │
│      └── list (GET) - List HTTP clients                 │
│                                                          │
│  System Routes                                           │
│  ├── /health - Health check                             │
│  ├── /health/live - Liveness probe                      │
│  ├── /health/ready - Readiness probe                    │
│  └── /metrics - Prometheus metrics                      │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

## Key Features

### 1. Single-Port Architecture
- Frontend and API served from the same port
- Eliminates CORS issues
- Simplifies deployment and configuration

### 2. Embedded Frontend
- Frontend files embedded using Go's `embed` package
- No separate build step needed at runtime
- Versioned together with backend code

### 3. Workflow Storage
- In-memory storage for development/testing
- Thread-safe operations
- UUID-based workflow IDs
- Full CRUD operations

### 4. HTTP Client Management
- Configurable HTTP clients with authentication
- Support for Basic, Bearer, and API Key auth
- SSRF protection with zero-trust model
- Reusable across workflow executions

### 5. Middleware Stack
- CORS handling
- Request logging
- Panic recovery
- Request size limiting

## Package Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Server entry point
├── workflow_registry.go         # Workflow storage registry
├── workflow_registry_test.go    # Workflow registry tests
└── pkg/
    ├── server/
    │   ├── server.go            # Server core
    │   ├── embed.go             # Frontend file embedding
    │   ├── static/              # Frontend build files (gitignored)
    │   │   └── README.md        # Placeholder
    │   ├── routes_httpclient.go # HTTP client routes
    │   ├── routes_workflow.go   # Workflow routes
    │   ├── routes_static.go     # Static file serving
    │   └── routes_httpclient_test.go
    ├── engine/                  # Workflow engine
    ├── executor/                # Node executors
    ├── httpclient/              # HTTP client registry
    └── ...
```

## Building

### Local Development Build

Use the provided build script:

```bash
# Build frontend and backend together
./build.sh

# Then run the server
./server
```

The build script:
1. Builds the Next.js frontend
2. Copies static files to `backend/pkg/server/static/`
3. Builds the Go backend with embedded files

### Docker Build

```bash
# Build with Docker (multi-stage build)
docker build -t thaiyyal .

# Or use docker-compose
docker-compose up --build
```

The Dockerfile uses a 3-stage multi-stage build:
1. **Frontend stage**: Builds Next.js application
2. **Backend stage**: Builds Go binary and embeds frontend files  
3. **Runtime stage**: Minimal Alpine image with just the binary


## Configuration

### Server Configuration
```go
type Config struct {
    Address            string        // Server address (e.g., ":8080")
    ReadTimeout        time.Duration // HTTP read timeout
    WriteTimeout       time.Duration // HTTP write timeout
    ShutdownTimeout    time.Duration // Graceful shutdown timeout
    MaxRequestBodySize int64         // Max request body size
    EnableCORS         bool          // Enable CORS headers
}
```

### Engine Configuration
```go
type Config struct {
    MaxExecutionTime    time.Duration // Max workflow execution time
    MaxNodeExecutions   int           // Max node executions per workflow
    MaxHTTPCallsPerExec int           // Max HTTP calls per execution
    MaxIterations       int           // Max loop iterations
}
```

## Usage

### Starting the Server

```bash
# Default settings
./server

# Custom configuration
./server -addr :9090 \
         -read-timeout 30s \
         -write-timeout 30s \
         -max-execution-time 1m \
         -max-node-executions 10000
```

### Accessing the Frontend

Open http://localhost:8080/ in your browser. The frontend will load and can make API calls to the same origin without CORS issues.

### API Examples

See [API_EXAMPLES.md](./API_EXAMPLES.md) for comprehensive API usage examples.

## Security Considerations

### SSRF Protection
- HTTP clients use zero-trust security model
- Private IPs, localhost, link-local addresses blocked by default
- Cloud metadata endpoints blocked by default
- Explicit allow-list required for restricted resources

### Input Validation
- All endpoints validate input
- Request body size limits enforced
- JSON validation on workflow data
- ID validation for resource access

### Error Handling
- Descriptive error messages without sensitive data leakage
- Proper HTTP status codes
- Request logging for debugging
- Panic recovery middleware

## Testing

### Unit Tests
```bash
# Test workflow registry
go test ./backend/workflow_registry_test.go ./backend/workflow_registry.go -v

# Test server package
go test ./backend/pkg/server -v

# Test all packages
go test ./backend/... -short
```

### Integration Tests
Manual integration testing performed with curl scripts. See test output in PR for verification.

### Security Scanning
CodeQL security scanning passed with 0 alerts for both Go and JavaScript code.

## Deployment

### Docker
```dockerfile
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./backend/cmd/server

FROM debian:bookworm-slim
COPY --from=builder /app/server /server
EXPOSE 8080
CMD ["/server"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  thaiyyal:
    build: .
    ports:
      - "8080:8080"
    environment:
      - MAX_EXECUTION_TIME=1m
      - MAX_NODE_EXECUTIONS=10000
```

## Future Enhancements

### Workflow Storage
- Add persistent storage backend (PostgreSQL, Redis, etc.)
- Implement workflow versioning
- Add workflow sharing/permissions
- Support workflow templates

**Note**: Workflow storage now follows the same registry pattern as HTTP clients, located in `backend/workflow_registry.go`

### HTTP Client Management
- Add client update/delete endpoints
- Implement client usage tracking
- Add client health checks
- Support client pooling

### Authentication & Authorization
- Add user authentication
- Implement role-based access control
- Support API keys for programmatic access
- Add OAuth2/OIDC integration

### Monitoring
- Enhanced metrics collection
- Distributed tracing
- Request/response logging
- Performance profiling

## Troubleshooting

### Server won't start
- Check if port is already in use: `lsof -i :8080`
- Verify configuration values are valid
- Check logs for error messages

### Frontend not loading
- Verify frontend files are embedded: `ls backend/pkg/frontend/static`
- Check browser console for errors
- Ensure correct base path in requests

### API requests failing
- Check request format matches documentation
- Verify Content-Type header is set
- Check server logs for detailed error messages
- Ensure request body size is within limits

### CORS errors
- Should not occur since frontend and API share same origin
- If seeing CORS errors, verify accessing via same domain/port
- Check browser is not caching old responses

## Contributing

When adding new endpoints:
1. Create route handler in appropriate routes_*.go file
2. Add request/response types
3. Register route in server.RegisterRoutes()
4. Add comprehensive unit tests
5. Update API documentation
6. Update this summary document

## License

See repository LICENSE file.

## Support

For issues or questions:
- Open an issue on GitHub
- Check documentation in docs/ directory
- Review API examples in API_EXAMPLES.md
