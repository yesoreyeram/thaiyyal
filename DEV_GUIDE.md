# Developer Guide

This guide provides comprehensive instructions for setting up, developing, testing, and debugging the Thaiyyal workflow engine.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Building the Project](#building-the-project)
- [Running Tests](#running-tests)
- [Debugging](#debugging)
- [Development Workflows](#development-workflows)
- [Performance Profiling](#performance-profiling)
- [Troubleshooting](#troubleshooting)

## Development Environment Setup

### Prerequisites

#### Required Software

| Software | Version | Purpose |
|----------|---------|---------|
| Go | 1.24.7+ | Backend development |
| Node.js | 20.x+ | Frontend development |
| npm | Latest | Package management |
| Git | 2.x+ | Version control |

#### Optional Tools

| Tool | Purpose |
|------|---------|
| golangci-lint | Go code linting |
| delve | Go debugger |
| VS Code | Recommended IDE |
| Docker | Container development |

### Installation Steps

#### 1. Install Go

```bash
# macOS (using Homebrew)
brew install go@1.24

# Linux (using official installer)
wget https://go.dev/dl/go1.24.7.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.7.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

#### 2. Install Node.js

```bash
# macOS (using Homebrew)
brew install node@20

# Linux (using nvm)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 20
nvm use 20

# Verify installation
node --version
npm --version
```

#### 3. Install Optional Tools

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Add Go bin to PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

#### 4. Clone and Setup Repository

```bash
# Clone the repository
git clone https://github.com/yesoreyeram/thaiyyal.git
cd thaiyyal

# Install frontend dependencies
npm install

# Install backend dependencies
cd backend
go mod download
go mod tidy

# Verify setup
cd ..
npm run lint
cd backend && go test ./...
```

### IDE Configuration

#### VS Code (Recommended)

Install recommended extensions:

```json
{
  "recommendations": [
    "golang.go",
    "dbaeumer.vscode-eslint",
    "esbenp.prettier-vscode",
    "bradlc.vscode-tailwindcss"
  ]
}
```

Workspace settings (`.vscode/settings.json`):

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.testFlags": ["-v"],
  "go.coverOnSave": true,
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[typescriptreact]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

## Project Structure

```
thaiyyal/
├── backend/                    # Go backend
│   ├── pkg/                   # Reusable packages
│   │   ├── config/           # Configuration management
│   │   ├── engine/           # Workflow execution engine
│   │   ├── executor/         # Node executors (40+ types)
│   │   ├── expression/       # Expression evaluation
│   │   ├── graph/            # Graph algorithms
│   │   ├── httpclient/       # HTTP client with middleware
│   │   ├── logging/          # Structured logging
│   │   ├── middleware/       # Execution middleware
│   │   ├── observer/         # Observer pattern
│   │   ├── security/         # Security utilities
│   │   ├── state/            # State management
│   │   └── types/            # Shared types
│   └── workflow.go           # Backward-compatible facade
├── src/                       # Next.js frontend
│   ├── app/                  # Next.js app directory
│   │   ├── layout.tsx        # Root layout
│   │   ├── page.tsx          # Home page
│   │   └── workflow/         # Workflow editor pages
│   ├── components/           # React components
│   │   ├── nodes/            # Custom node components
│   │   ├── AppNavBar.tsx     # Application navbar
│   │   ├── NodePalette.tsx   # Node palette
│   │   └── ...
│   └── types/                # TypeScript types
├── public/                    # Static assets
├── docs/                      # Documentation
├── .vscode/                   # VS Code settings
├── go.mod                     # Go dependencies
├── package.json               # npm dependencies
├── next.config.ts             # Next.js config
├── tsconfig.json              # TypeScript config
└── README.md                  # Main documentation
```

### Package Responsibilities

| Package | Responsibility |
|---------|---------------|
| `engine` | Workflow parsing, validation, orchestration |
| `executor` | Node execution logic, registry pattern |
| `graph` | Topological sort, cycle detection |
| `state` | Variables, accumulators, cache |
| `security` | SSRF protection, validation |
| `logging` | Structured logging with slog |
| `middleware` | Cross-cutting concerns |
| `observer` | Event notification system |
| `httpclient` | HTTP client with middleware |
| `types` | Shared type definitions |

## Building the Project

### Frontend Build

```bash
# Development build (with hot reload)
npm run dev

# Production build
npm run build

# Start production server
npm start

# Run linter
npm run lint

# Fix linting issues
npm run lint -- --fix
```

### Backend Build

```bash
cd backend

# Build all packages
go build ./...

# Build with race detector
go build -race ./...

# Build with optimizations disabled (for debugging)
go build -gcflags="all=-N -l" ./...

# Cross-compilation
GOOS=linux GOARCH=amd64 go build ./...
```

### Build Artifacts

Build artifacts are gitignored:

```
/backend/examples/examples
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
/.next/
/out/
/build/
```

## Running Tests

### Backend Tests

#### Run All Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific package
go test ./pkg/engine/...

# Run specific test
go test -run TestSimpleAddition ./pkg/engine/...
```

#### Test Options

```bash
# Run tests in parallel
go test -parallel 4 ./...

# Run with race detector
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run with short flag (skip long tests)
go test -short ./...

# Set timeout
go test -timeout 30s ./...
```

#### Table-Driven Tests Example

```go
func TestOperation(t *testing.T) {
    tests := []struct {
        name     string
        op       string
        inputs   []float64
        expected float64
        wantErr  bool
    }{
        {"add", "add", []float64{5, 3}, 8, false},
        {"subtract", "subtract", []float64{10, 3}, 7, false},
        {"multiply", "multiply", []float64{4, 5}, 20, false},
        {"divide", "divide", []float64{10, 2}, 5, false},
        {"divide by zero", "divide", []float64{10, 0}, 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Frontend Tests

```bash
# Run linter
npm run lint

# Build to check for errors
npm run build
```

### Continuous Testing

```bash
# Watch mode for Go tests
find . -name "*.go" | entr -c go test ./...

# Watch mode for frontend
npm run dev
```

## Debugging

### Backend Debugging

#### Using Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug a test
cd backend
dlv test ./pkg/engine -- -test.run TestSimpleAddition

# Set breakpoints
(dlv) break engine.go:100
(dlv) continue
(dlv) print variable
(dlv) next
(dlv) step
```

#### VS Code Debug Configuration

`.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Current Test",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${fileDirname}",
      "args": [
        "-test.run",
        "${input:testName}"
      ]
    },
    {
      "name": "Debug Package",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}/backend/pkg/engine"
    }
  ],
  "inputs": [
    {
      "id": "testName",
      "type": "promptString",
      "description": "Test name pattern"
    }
  ]
}
```

#### Logging for Debugging

```go
import "github.com/yesoreyeram/thaiyyal/backend/pkg/logging"

// Create logger with debug level
cfg := logging.DefaultConfig()
cfg.Level = "debug"
cfg.Pretty = true
logger := logging.New(cfg)

// Use logger
logger.WithField("node_id", "123").Debug("executing node")
logger.WithError(err).Error("execution failed")
```

### Frontend Debugging

#### Browser DevTools

1. Open Chrome DevTools (F12)
2. Use React DevTools extension
3. Use Console for logging
4. Use Network tab for API calls
5. Use Performance tab for profiling

#### Debug Configuration

```typescript
// Enable debug logging
if (process.env.NODE_ENV === 'development') {
  console.log('Debug: Workflow state', workflowState);
}
```

## Development Workflows

### Adding a New Node Type

#### 1. Define Node Type

```go
// backend/pkg/types/types.go
const NodeTypeMyCustom NodeType = "my_custom"
```

#### 2. Create Executor

```go
// backend/pkg/executor/my_custom.go
package executor

import (
    "fmt"
    "github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

type MyCustomExecutor struct{}

func (e *MyCustomExecutor) Type() types.NodeType {
    return types.NodeTypeMyCustom
}

func (e *MyCustomExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    // Implementation
    return result, nil
}

func (e *MyCustomExecutor) Validate(node types.Node) error {
    // Validation
    return nil
}
```

#### 3. Add Tests

```go
// backend/pkg/executor/my_custom_test.go
func TestMyCustomExecutor(t *testing.T) {
    // Test implementation
}
```

#### 4. Register Executor

```go
// backend/pkg/engine/engine.go (DefaultRegistry function)
reg.MustRegister(&executor.MyCustomExecutor{})
```

#### 5. Update Documentation

- Add to [Node Types Reference](docs/NODE_TYPES.md)
- Update README.md node count
- Add usage examples

### Adding Middleware

```go
// backend/pkg/middleware/my_middleware.go
package middleware

type MyMiddleware struct {
    config Config
}

func NewMyMiddleware(config Config) *MyMiddleware {
    return &MyMiddleware{config: config}
}

func (m *MyMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
    // Pre-processing
    
    // Execute next middleware/executor
    result, err := next(ctx, node)
    
    // Post-processing
    
    return result, err
}

func (m *MyMiddleware) Name() string {
    return "MyMiddleware"
}
```

### Adding Security Checks

```go
// backend/pkg/security/my_check.go
package security

func ValidateInput(input string) error {
    // Validation logic
    if len(input) > MaxLength {
        return fmt.Errorf("input too long: %d (max %d)", len(input), MaxLength)
    }
    return nil
}
```

## Performance Profiling

### CPU Profiling

```bash
# Run tests with CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./pkg/engine/...

# Analyze profile
go tool pprof cpu.prof
(pprof) top10
(pprof) list functionName
(pprof) web  # Opens graph in browser
```

### Memory Profiling

```bash
# Run tests with memory profiling
go test -memprofile=mem.prof -bench=. ./pkg/engine/...

# Analyze profile
go tool pprof mem.prof
(pprof) top10
(pprof) list functionName
```

### Benchmarking

```go
// Write benchmarks
func BenchmarkWorkflowExecution(b *testing.B) {
    payload := []byte(`{...}`)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        engine, _ := workflow.NewEngine(payload)
        engine.Execute()
    }
}
```

```bash
# Run benchmarks
go test -bench=. -benchmem ./pkg/engine/...

# Compare benchmarks
go test -bench=. -benchmem ./pkg/engine/... > old.txt
# Make changes
go test -bench=. -benchmem ./pkg/engine/... > new.txt
benchstat old.txt new.txt
```

### Trace Analysis

```bash
# Generate trace
go test -trace=trace.out ./pkg/engine/...

# Analyze trace
go tool trace trace.out
```

## Troubleshooting

### Common Issues

#### Go Module Issues

```bash
# Problem: Dependency resolution fails
# Solution: Clean module cache
go clean -modcache
go mod download

# Problem: Inconsistent dependencies
# Solution: Tidy modules
go mod tidy
```

#### Frontend Build Issues

```bash
# Problem: npm install fails
# Solution: Clear cache
rm -rf node_modules package-lock.json
npm cache clean --force
npm install

# Problem: Next.js build errors
# Solution: Clean build cache
rm -rf .next
npm run build
```

#### Test Failures

```bash
# Problem: Tests fail intermittently
# Solution: Run with race detector
go test -race ./...

# Problem: Tests timeout
# Solution: Increase timeout
go test -timeout 60s ./...
```

### Debug Environment Variables

```bash
# Enable Go module debug output
export GOMODULE_DEBUG=1

# Enable Go build verbose output
go build -v ./...

# Enable Go test verbose output
go test -v ./...
```

### Getting Help

1. Check [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
2. Search [GitHub Issues](https://github.com/yesoreyeram/thaiyyal/issues)
3. Ask in [GitHub Discussions](https://github.com/yesoreyeram/thaiyyal/discussions)
4. Review [Architecture Documentation](docs/ARCHITECTURE.md)

## Best Practices

### Code Quality

1. **Run tests before committing**
   ```bash
   cd backend && go test ./...
   cd .. && npm run lint
   ```

2. **Format code**
   ```bash
   cd backend && go fmt ./...
   npm run lint -- --fix
   ```

3. **Check coverage**
   ```bash
   cd backend && go test -cover ./...
   ```

4. **Run linter**
   ```bash
   cd backend && golangci-lint run
   ```

### Git Workflow

```bash
# Create feature branch
git checkout -b feature/my-feature

# Make changes and commit frequently
git add .
git commit -m "feat(executor): add new node type"

# Keep branch updated
git fetch upstream
git rebase upstream/main

# Push changes
git push origin feature/my-feature
```

### Documentation

1. Update documentation with code changes
2. Include examples in documentation
3. Add comments for complex logic
4. Keep README.md updated

## Additional Resources

- [Architecture Documentation](docs/ARCHITECTURE.md)
- [API Reference](docs/API_REFERENCE.md)
- [Security Best Practices](docs/SECURITY_BEST_PRACTICES.md)
- [Performance Tuning](docs/PERFORMANCE_TUNING.md)
- [Contributing Guidelines](CONTRIBUTING.md)

---

Need help? Open an issue or start a discussion on GitHub!
