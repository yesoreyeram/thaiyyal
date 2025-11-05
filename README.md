# Thaiyyal - Visual Workflow Engine

[![Go Version](https://img.shields.io/badge/Go-1.24.7-blue.svg)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-16.0-black)](https://nextjs.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen.svg)](backend/pkg/)

A powerful, secure, and extensible visual workflow engine that combines a Next.js-based visual editor with a robust Go backend execution engine. Build complex data processing workflows through an intuitive drag-and-drop interface.

## 🌟 Key Features

### Visual Workflow Builder
- **Drag-and-Drop Interface**: Intuitive React Flow-based visual editor
- **Real-time Validation**: Instant feedback on workflow configuration
- **Rich Node Palette**: 40+ built-in node types across multiple categories
- **Workflow Persistence**: Save and load workflows with version control
- **Live Execution Results**: Resizable panel with real-time execution feedback
- **Interactive Canvas**: Pan, zoom, and organize workflows visually

### Powerful Execution Engine
- **Type-Safe Execution**: Strong typing with automatic type inference
- **Topological Sorting**: Automatic dependency resolution and execution ordering
- **Parallel Processing**: Execute independent nodes concurrently
- **State Management**: Variables, accumulators, counters, and caching
- **Real-time Feedback**: Live execution results with node-by-node breakdown
- **Cancellation Support**: Stop long-running workflows on demand

### Execution Results Display
- **Resizable Results Panel**: Drag-to-resize execution panel (100px-600px)
- **Loading Indicators**: Professional animated spinners and progress display
- **Detailed Results**: Execution summary, final output, and per-node results
- **Error Handling**: Comprehensive error messages with detailed context
- **Execution History**: View previous execution results
- **Cancel Execution**: Abort long-running workflows instantly

### Security & Protection
- **Zero-Trust Architecture**: All inputs validated, all outputs sanitized
- **SSRF Protection**: Comprehensive protection against Server-Side Request Forgery
- **Resource Limits**: Configurable limits for execution time, memory, and iterations
- **Input Sanitization**: Automatic validation and sanitization of all inputs
- **API Protection**: Per-execution HTTP call limits, rate limiting, and circuit breakers (planned)

### Extensibility
- **Custom Node Executors**: Plugin architecture for custom node types
- **Middleware System**: Composable middleware for cross-cutting concerns
- **Observer Pattern**: Hook into workflow execution events
- **HTTP Client Registry**: Named HTTP clients with custom configurations

## 📚 Documentation

### Getting Started
- [Quick Start Guide](#quick-start)
- [Installation](#installation)
- [Basic Concepts](#basic-concepts)
- [**Workflow Execution Guide**](docs/WORKFLOW_EXECUTION_GUIDE.md) - Step-by-step guide with screenshots

### Core Documentation
- [Architecture Overview](docs/ARCHITECTURE.md)
- [Design Patterns](docs/ARCHITECTURE_DESIGN_PATTERNS.md)
- [Server Implementation](docs/SERVER_IMPLEMENTATION.md)
- [API Examples](docs/API_EXAMPLES.md)
- [Developer Guide](DEV_GUIDE.md)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Agent System](AGENTS.md)

### Principles & Best Practices
- [Zero-Trust Security](docs/PRINCIPLES_ZERO_TRUST.md)
- [Workload Protection](docs/PRINCIPLES_WORKLOAD_PROTECTION.md)
- [No Runtime Errors](docs/PRINCIPLES_NO_RUNTIME_ERRORS.md)
- [Pluggable Architecture](docs/PRINCIPLES_PLUGGABLE_ARCHITECTURE.md)
- [Modular Design](docs/PRINCIPLES_MODULAR_DESIGN.md)
- [Single Responsibility](docs/PRINCIPLES_SINGLE_RESPONSIBILITY.md)

### Requirements
- [Functional Requirements](docs/REQUIREMENTS_FUNCTIONAL.md)
- [Non-Functional Requirements](docs/REQUIREMENTS_NON_FUNCTIONAL.md)
  - [Security](docs/REQUIREMENTS_NON_FUNCTIONAL_SECURITY.md)
  - [Observability](docs/REQUIREMENTS_NON_FUNCTIONAL_OBSERVABILITY.md)
  - [Code Quality](docs/REQUIREMENTS_NON_FUNCTIONAL_CODE_QUALITY.md)
  - [Logging](docs/REQUIREMENTS_NON_FUNCTIONAL_LOGGING.md)
  - [Testing](docs/REQUIREMENTS_NON_FUNCTIONAL_TESTING.md)
  - [Governance](docs/REQUIREMENTS_NON_FUNCTIONAL_GOVERNANCE.md)
  - [Deployment](docs/REQUIREMENTS_NON_FUNCTIONAL_DEPLOYMENT.md)

### Additional Resources
- [Node Types Reference](docs/NODE_TYPES.md)
- [API Reference](docs/API_REFERENCE.md)
- [Security Best Practices](docs/SECURITY_BEST_PRACTICES.md)
- [API Protection Analysis](docs/API_PROTECTION_SECURITY_ANALYSIS.md) - **NEW**
- [API Protection Quick Reference](docs/API_PROTECTION_QUICK_REFERENCE.md) - **NEW**
- [Performance Tuning](docs/PERFORMANCE_TUNING.md)
- [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
- [Examples & Tutorials](docs/EXAMPLES.md)

## 🚀 Quick Start

### Prerequisites

- **Go**: 1.24.7 or later
- **Node.js**: 20.x or later
- **npm**: Latest version

### Installation

```bash
# Clone the repository
git clone https://github.com/yesoreyeram/thaiyyal.git
cd thaiyyal

# Install frontend dependencies
npm install

# Build backend (optional - for development)
cd backend && go build ./...
```

### Running the Application

#### Development Mode with Hot Reloading (Recommended)

Run both frontend and backend with automatic hot reloading:

```bash
# Start both Next.js dev server and Go backend
./dev.sh

# Or using npm
npm run dev:full

# Access the application at http://localhost:8080
# Frontend changes reload instantly via Next.js Fast Refresh
```

**Alternative: Manual Setup**

Terminal 1 - Frontend only:
```bash
npm run dev
# Next.js dev server on http://localhost:3000
```

Terminal 2 - Backend with dev proxy:
```bash
cd backend/cmd/server
go run -tags dev . -addr :8080
# Backend on http://localhost:8080 (proxies to Next.js)
```

**Docker Development**

```bash
# Start development environment with hot reloading
docker-compose -f docker-compose.dev.yml up
```

See [DEV_HOT_RELOAD.md](DEV_HOT_RELOAD.md) for detailed information about hot module reloading.

#### Production Build

```bash
# Build the frontend
npm run build

# Start the production server
npm start
```

### Your First Workflow

1. **Open the Application**: Navigate to http://localhost:3000
2. **Create a New Workflow**: Click "New Workflow" button
3. **Add Nodes**: Drag nodes from the palette onto the canvas
4. **Connect Nodes**: Draw edges between nodes to define data flow
5. **Configure Nodes**: Click nodes to configure their properties
6. **Execute**: Click the "Run" button to execute your workflow
7. **View Results**: See results in the visualization panel

## 🎥 Demo Video

See Thaiyyal in action! This video demonstrates the complete workflow creation and execution process:

<!-- Video will be embedded here once recorded -->
<!--
![Thaiyyal Workflow Demo](docs/demo/thaiyyal-demo.mp4)

Or view on: [YouTube Link](https://youtube.com/watch?v=...)
-->

### What the Demo Shows

1. **Creating a Workflow** (0:00-0:30)
   - Opening the workflow builder
   - Adding nodes from the palette
   - Connecting nodes visually

2. **Configuring Nodes** (0:30-1:00)
   - Setting node properties
   - Configuring data sources
   - Adjusting parameters

3. **Running the Workflow** (1:00-1:30)
   - Clicking the Run button
   - Viewing the execution panel expand
   - Watching the loading indicator

4. **Viewing Results** (1:30-2:00)
   - Reviewing execution summary
   - Examining final output
   - Exploring node-by-node results

5. **Advanced Features** (2:00-2:30)
   - Resizing the results panel
   - Handling errors gracefully
   - Canceling long-running executions

📖 For detailed instructions, see the [Workflow Execution Guide](docs/WORKFLOW_EXECUTION_GUIDE.md)

## 💡 Basic Concepts

### Workflows

A workflow is a directed acyclic graph (DAG) consisting of:
- **Nodes**: Processing units that perform specific operations
- **Edges**: Connections that define data flow between nodes
- **Data**: Values that flow through the workflow

### Node Categories

```
📥 Basic I/O
   └─ Number, Text Input, Visualization

⚙️ Operations
   └─ Mathematical Operations, Text Operations, HTTP Requests

🔀 Control Flow
   └─ Condition, ForEach, While Loop, Filter, Map, Reduce

📊 Array Processing
   └─ Slice, Sort, Find, FlatMap, GroupBy, Unique, Chunk, Reverse, 
      Partition, Zip, Sample, Range, Compact, Transpose

💾 State & Memory
   └─ Variable, Extract, Transform, Accumulator, Counter, Parse

🔧 Advanced Control
   └─ Switch, Parallel, Join, Split, Delay, Cache

🛡️ Error Handling
   └─ Retry, Try/Catch, Timeout

🌐 Context
   └─ Context Variable, Context Constant
```

### Execution Model

```
┌─────────────────────────────────────────────────────────┐
│                    Workflow Execution                    │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  1. Parse JSON Payload                                  │
│     ↓                                                    │
│  2. Infer Node Types (if not explicit)                  │
│     ↓                                                    │
│  3. Validate Structure & Configuration                  │
│     ↓                                                    │
│  4. Topological Sort (determine execution order)        │
│     ↓                                                    │
│  5. Execute Nodes in Order                              │
│     ├─ Apply Middleware                                 │
│     ├─ Interpolate Templates                            │
│     ├─ Execute Node Logic                               │
│     ├─ Store Results                                    │
│     └─ Notify Observers                                 │
│     ↓                                                    │
│  6. Return Final Output                                 │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

## 🏗️ Architecture

Thaiyyal follows a **modular, layered architecture** with clear separation of concerns:

```
┌──────────────────────────────────────────────────────┐
│                    Frontend (Next.js)                 │
│  ┌────────────┬────────────┬─────────────────────┐  │
│  │  UI Layer  │ Components │   React Flow Editor  │  │
│  └────────────┴────────────┴─────────────────────┘  │
└──────────────────────────────────────────────────────┘
                         │
                         │ JSON Payload
                         ↓
┌──────────────────────────────────────────────────────┐
│              Backend (Go - Workflow Engine)          │
│  ┌────────────────────────────────────────────────┐ │
│  │              Engine Layer                       │ │
│  │  • Workflow Parsing & Validation                │ │
│  │  • Topological Sorting                          │ │
│  │  • Execution Orchestration                      │ │
│  └────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────┐ │
│  │           Executor Layer                        │ │
│  │  • Node Executors (40+ types)                   │ │
│  │  • Registry Pattern                             │ │
│  │  • Middleware Chain                             │ │
│  └────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────┐ │
│  │          Support Services                       │ │
│  │  • State Management   • Graph Analysis          │ │
│  │  • HTTP Client        • Security (SSRF)         │ │
│  │  • Logging            • Observer Pattern        │ │
│  └────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

### Package Structure

```
backend/
├── pkg/
│   ├── config/          # Configuration management
│   ├── engine/          # Workflow execution engine
│   ├── executor/        # Node executors and registry
│   ├── expression/      # Expression evaluation
│   ├── graph/           # Graph algorithms (topological sort)
│   ├── httpclient/      # HTTP client with middleware
│   ├── logging/         # Structured logging (slog)
│   ├── middleware/      # Execution middleware
│   ├── observer/        # Observer pattern implementation
│   ├── security/        # Security utilities (SSRF protection)
│   ├── state/           # State management
│   └── types/           # Shared type definitions
└── workflow.go          # Backward-compatible facade
```

## 🔐 Security Features

### SSRF Protection
- Blocks private IP ranges (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16)
- Blocks localhost and loopback addresses
- Blocks link-local addresses (169.254.0.0/16)
- Blocks cloud metadata endpoints (169.254.169.254)
- Configurable domain allow/block lists

### Input Validation
- Maximum string length limits
- Maximum array size limits
- Maximum object depth limits
- Type validation
- Sanitization of user inputs

### Resource Protection
- Maximum execution time
- Maximum node executions per workflow
- **Maximum HTTP calls per execution (default: 100)** - **NEW**
- Maximum loop iterations
- Maximum recursion depth
- Maximum variables count

## 🧩 Extensibility

### Custom Node Executors

```go
// Define a custom executor
type MyCustomExecutor struct{}

func (e *MyCustomExecutor) Type() types.NodeType {
    return "my_custom_node"
}

func (e *MyCustomExecutor) Execute(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
    // Your custom logic here
    return result, nil
}

func (e *MyCustomExecutor) Validate(node types.Node) error {
    // Validation logic
    return nil
}

// Register the executor
registry := engine.DefaultRegistry()
registry.MustRegister(&MyCustomExecutor{})
engine, err := engine.NewWithRegistry(payload, config, registry)
```

### Middleware

```go
// Add custom middleware
type MyMiddleware struct{}

func (m *MyMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next middleware.Handler) (interface{}, error) {
    // Pre-processing
    result, err := next(ctx, node)
    // Post-processing
    return result, err
}

func (m *MyMiddleware) Name() string {
    return "MyMiddleware"
}
```

### Observers

```go
// Monitor workflow execution
type MyObserver struct{}

func (o *MyObserver) OnEvent(ctx context.Context, event observer.Event) {
    // Handle execution events
}

engine.RegisterObserver(&MyObserver{})
```

## 📊 Example Workflows

### Simple Addition

```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": 10}},
    {"id": "2", "type": "number", "data": {"value": 5}},
    {"id": "3", "type": "operation", "data": {"op": "add"}}
  ],
  "edges": [
    {"source": "1", "target": "3"},
    {"source": "2", "target": "3"}
  ]
}
```

### HTTP Request with Error Handling

```json
{
  "nodes": [
    {"id": "1", "type": "http", "data": {"url": "https://api.example.com/data"}},
    {"id": "2", "type": "try_catch", "data": {"fallbackValue": {"error": "API unavailable"}}},
    {"id": "3", "type": "extract", "data": {"field": "results"}},
    {"id": "4", "type": "visualization", "data": {"mode": "table"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "4"}
  ]
}
```

### Array Processing Pipeline

```json
{
  "nodes": [
    {"id": "1", "type": "range", "data": {"start": 1, "end": 100}},
    {"id": "2", "type": "filter", "data": {"expression": "x % 2 == 0"}},
    {"id": "3", "type": "map", "data": {"expression": "x * x"}},
    {"id": "4", "type": "reduce", "data": {"op": "sum"}},
    {"id": "5", "type": "visualization", "data": {"mode": "text"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "4"},
    {"source": "4", "target": "5"}
  ]
}
```

## 🚀 Running the Server

### HTTP API Server

Start the standalone HTTP server for programmatic workflow execution:

```bash
# Build the server
cd backend/cmd/server
go build -o ../../bin/thaiyyal-server .

# Run with default settings
../../bin/thaiyyal-server

# Run with custom configuration
../../bin/thaiyyal-server \
  -addr :9090 \
  -max-execution-time 2m \
  -max-node-executions 5000
```

The server exposes:
- `POST /api/v1/workflow/execute` - Execute workflows
- `POST /api/v1/workflow/validate` - Validate workflows
- `GET /health` - Health check
- `GET /health/live` - Liveness probe (K8s)
- `GET /health/ready` - Readiness probe (K8s)
- `GET /metrics` - Prometheus metrics

### Docker Deployment

```bash
# Build Docker image
docker build -t thaiyyal/workflow-engine:latest .

# Run container
docker run -d -p 8080:8080 thaiyyal/workflow-engine:latest

# Or use Docker Compose (includes Prometheus & Grafana)
docker-compose up -d
```

Access:
- Workflow API: http://localhost:8080
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001 (admin/admin)

### Kubernetes Deployment

```bash
# Deploy to Kubernetes
kubectl apply -f deployments/kubernetes/deployment.yaml

# Verify deployment
kubectl get pods -n thaiyyal

# Port forward for local access
kubectl port-forward -n thaiyyal svc/thaiyyal-server 8080:8080
```

See [Operations Guide](docs/OPERATIONS_GUIDE.md) for detailed deployment and scaling information.

## 📊 Observability & Monitoring

### Metrics

The server exports Prometheus metrics including:

```promql
# Workflow execution metrics
workflow_executions_total          # Total workflow executions
workflow_execution_duration        # Execution duration histogram
workflow_executions_success_total  # Successful executions
workflow_executions_failure_total  # Failed executions

# Node execution metrics
node_executions_total              # Total node executions
node_execution_duration            # Node execution duration
node_executions_success_total      # Successful node executions
node_executions_failure_total      # Failed node executions

# HTTP call metrics
http_calls_total                   # Total HTTP calls
http_call_duration                 # HTTP call duration
```

### Distributed Tracing

Integrated with OpenTelemetry for distributed tracing:
- Workflow-level spans
- Node-level spans
- Automatic trace ID generation
- Context propagation

### Health Checks

Multiple health check endpoints for monitoring:
- **Liveness**: Server is running
- **Readiness**: Server can handle requests
- **Health**: Comprehensive checks with dependency status

## 📚 Documentation

### Getting Started
- [Quick Start Guide](#quick-start)
- [Installation](#installation)
- [Basic Concepts](#basic-concepts)

### Core Documentation
- [Architecture Overview](docs/ARCHITECTURE.md)
- [Design Patterns](docs/ARCHITECTURE_DESIGN_PATTERNS.md)
- [Developer Guide](DEV_GUIDE.md)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Agent System](AGENTS.md)

### Operations
- [Operations Guide](docs/OPERATIONS_GUIDE.md) - **NEW**
- [API Reference](docs/api/openapi.yaml) - **NEW**
- [Performance Tuning](docs/PERFORMANCE_TUNING.md)
- [Security Best Practices](docs/SECURITY_BEST_PRACTICES.md)

### Principles & Best Practices
- [Zero-Trust Security](docs/PRINCIPLES_ZERO_TRUST.md)
- [Workload Protection](docs/PRINCIPLES_WORKLOAD_PROTECTION.md)
- [No Runtime Errors](docs/PRINCIPLES_NO_RUNTIME_ERRORS.md)
- [Pluggable Architecture](docs/PRINCIPLES_PLUGGABLE_ARCHITECTURE.md)
- [Modular Design](docs/PRINCIPLES_MODULAR_DESIGN.md)
- [Single Responsibility](docs/PRINCIPLES_SINGLE_RESPONSIBILITY.md)

### Requirements
- [Functional Requirements](docs/REQUIREMENTS_FUNCTIONAL.md)
- [Non-Functional Requirements](docs/REQUIREMENTS_NON_FUNCTIONAL.md)
  - [Security](docs/REQUIREMENTS_NON_FUNCTIONAL_SECURITY.md)
  - [Observability](docs/REQUIREMENTS_NON_FUNCTIONAL_OBSERVABILITY.md)
  - [Code Quality](docs/REQUIREMENTS_NON_FUNCTIONAL_CODE_QUALITY.md)
  - [Logging](docs/REQUIREMENTS_NON_FUNCTIONAL_LOGGING.md)
  - [Testing](docs/REQUIREMENTS_NON_FUNCTIONAL_TESTING.md)
  - [Governance](docs/REQUIREMENTS_NON_FUNCTIONAL_GOVERNANCE.md)
  - [Deployment](docs/REQUIREMENTS_NON_FUNCTIONAL_DEPLOYMENT.md)

### Additional Resources
- [Node Types Reference](docs/NODE_TYPES.md)
- [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
- [Examples & Tutorials](docs/EXAMPLES.md)

## 🧪 Testing

```bash
# Run all tests
cd backend
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/engine/...

# Run with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

Test coverage: **95%+** across all packages

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on:
- Code of Conduct
- Development workflow
- Coding standards
- Submitting pull requests
- Reporting issues

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [React Flow](https://reactflow.dev/) - Visual graph editor
- [Next.js](https://nextjs.org/) - React framework
- [Tailwind CSS](https://tailwindcss.com/) - Styling
- Go community for excellent tooling and libraries

## 📧 Contact & Support

- **Issues**: [GitHub Issues](https://github.com/yesoreyeram/thaiyyal/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yesoreyeram/thaiyyal/discussions)
- **Security**: See [Security Policy](SECURITY.md) for reporting vulnerabilities

## 🗺️ Roadmap

### Completed ✅
- [x] OpenTelemetry integration
- [x] Prometheus metrics
- [x] HTTP API server
- [x] Health checks and probes
- [x] Docker deployment
- [x] Kubernetes deployment templates
- [x] Distributed tracing

### In Progress 🚧
- [ ] Workflow persistence layer
- [ ] Execution history API
- [ ] Workflow versioning

### Planned 📋
- [ ] GraphQL API support
- [ ] WebSocket-based real-time execution
- [ ] More built-in node types
- [ ] Visual debugging tools
- [ ] Workflow templates marketplace
- [ ] Multi-language expression support
- [ ] Circuit breaker pattern
- [ ] Rate limiting
- [ ] API authentication/authorization

---

**Built with ❤️ by the Thaiyyal team**
