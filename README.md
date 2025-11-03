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

### Powerful Execution Engine
- **Type-Safe Execution**: Strong typing with automatic type inference
- **Topological Sorting**: Automatic dependency resolution and execution ordering
- **Parallel Processing**: Execute independent nodes concurrently
- **State Management**: Variables, accumulators, counters, and caching

### Security & Protection
- **Zero-Trust Architecture**: All inputs validated, all outputs sanitized
- **SSRF Protection**: Comprehensive protection against Server-Side Request Forgery
- **Resource Limits**: Configurable limits for execution time, memory, and iterations
- **Input Sanitization**: Automatic validation and sanitization of all inputs

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

### Core Documentation
- [Architecture Overview](docs/ARCHITECTURE.md)
- [Design Patterns](docs/ARCHITECTURE_DESIGN_PATTERNS.md)
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

#### Development Mode

```bash
# Start the Next.js development server
npm run dev

# Open your browser to http://localhost:3000
```

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
- Maximum HTTP calls per execution
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

- [ ] GraphQL API support
- [ ] WebSocket-based real-time execution
- [ ] Cloud deployment templates
- [ ] More built-in node types
- [ ] Visual debugging tools
- [ ] Workflow templates marketplace
- [ ] Multi-language expression support
- [ ] Advanced analytics and monitoring

---

**Built with ❤️ by the Thaiyyal team**
