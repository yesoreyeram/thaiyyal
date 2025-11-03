# Architecture Overview

This document provides a comprehensive overview of the Thaiyyal workflow engine architecture.

## Table of Contents

- [System Overview](#system-overview)
- [High-Level Architecture](#high-level-architecture)
- [Component Architecture](#component-architecture)
- [Data Flow](#data-flow)
- [Security Architecture](#security-architecture)
- [Deployment Architecture](#deployment-architecture)

## System Overview

Thaiyyal is a visual workflow engine that combines a modern web-based editor with a robust backend execution engine. The system is designed with security, extensibility, and performance as core principles.

### Key Characteristics

- **Hybrid Architecture**: Next.js frontend + Go backend
- **Visual-First**: Drag-and-drop workflow creation
- **Type-Safe**: Strong typing throughout the stack
- **Secure by Default**: Zero-trust security model
- **Extensible**: Plugin architecture for custom nodes
- **Observable**: Built-in logging, metrics, and tracing

## High-Level Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                         User Interface                          │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │              Next.js Frontend (Port 3000)                 │ │
│  │                                                            │ │
│  │  ┌────────────┐  ┌────────────┐  ┌──────────────────┐   │ │
│  │  │ React Flow │  │ Components │  │  State Management│   │ │
│  │  │   Editor   │  │  (UI/UX)   │  │    (React)       │   │ │
│  │  └────────────┘  └────────────┘  └──────────────────┘   │ │
│  │                                                            │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │          Workflow Definition (JSON)                  │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └──────────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/JSON
                              ↓
┌────────────────────────────────────────────────────────────────┐
│                      Backend Services                           │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │                 Go Workflow Engine                        │ │
│  │                                                            │ │
│  │  ┌────────────────────────────────────────────────────┐  │ │
│  │  │              Engine Layer                           │  │ │
│  │  │  • JSON Parsing                                     │  │ │
│  │  │  • Type Inference                                   │  │ │
│  │  │  • Validation                                       │  │ │
│  │  │  • Topological Sort                                 │  │ │
│  │  │  • Execution Orchestration                          │  │ │
│  │  └────────────────────────────────────────────────────┘  │ │
│  │                                                            │ │
│  │  ┌────────────────────────────────────────────────────┐  │ │
│  │  │              Executor Layer                         │  │ │
│  │  │  • 40+ Node Executors                               │  │ │
│  │  │  • Registry Pattern                                 │  │ │
│  │  │  • Middleware Chain                                 │  │ │
│  │  │  • Custom Executors                                 │  │ │
│  │  └────────────────────────────────────────────────────┘  │ │
│  │                                                            │ │
│  │  ┌────────────────────────────────────────────────────┐  │ │
│  │  │          Supporting Services                        │  │ │
│  │  │  ┌─────────────┐  ┌──────────────┐  ┌───────────┐ │  │ │
│  │  │  │    State    │  │    Graph     │  │  Security │ │  │ │
│  │  │  │  Manager    │  │  Algorithms  │  │   (SSRF)  │ │  │ │
│  │  │  └─────────────┘  └──────────────┘  └───────────┘ │  │ │
│  │  │  ┌─────────────┐  ┌──────────────┐  ┌───────────┐ │  │ │
│  │  │  │   Logging   │  │  Observer    │  │   HTTP    │ │  │ │
│  │  │  │   (slog)    │  │   Pattern    │  │  Client   │ │  │ │
│  │  │  └─────────────┘  └──────────────┘  └───────────┘ │  │ │
│  │  └────────────────────────────────────────────────────┘  │ │
│  └──────────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────────┘
```

## Component Architecture

### Frontend Components

```
src/
├── app/                    # Next.js App Router
│   ├── layout.tsx         # Root layout
│   ├── page.tsx           # Home page
│   └── workflow/          # Workflow editor
│       ├── page.tsx       # Editor page
│       └── layout.tsx     # Editor layout
│
├── components/            # React Components
│   ├── AppNavBar.tsx     # Main navigation
│   ├── WorkflowNavBar.tsx # Editor navigation
│   ├── NodePalette.tsx    # Node selection panel
│   ├── HeroComponent.tsx  # Landing page hero
│   ├── Toast.tsx          # Notifications
│   ├── JSONPayloadModal.tsx
│   ├── OpenWorkflowModal.tsx
│   ├── WorkflowStatusBar.tsx
│   └── nodes/             # Custom node components
│       ├── NumberNode.tsx
│       ├── OperationNode.tsx
│       └── ...
│
└── types/                 # TypeScript types
    └── workflow.ts        # Workflow definitions
```

### Backend Packages

```
backend/pkg/
├── config/          # Configuration management
│   ├── config.go   # Config structures
│   ├── errors.go   # Config errors
│   └── doc.go      # Package documentation
│
├── engine/          # Workflow execution engine
│   ├── engine.go   # Engine implementation
│   ├── errors.go   # Engine-specific errors
│   ├── snapshot.go # State snapshot for debugging
│   └── *_test.go   # Tests (95%+ coverage)
│
├── executor/        # Node executors
│   ├── executor.go      # Base executor interface
│   ├── registry.go      # Executor registry
│   ├── input_*.go       # Input nodes
│   ├── operation.go     # Operation nodes
│   ├── control_*.go     # Control flow nodes
│   ├── state_*.go       # State management nodes
│   └── *_test.go        # Comprehensive tests
│
├── graph/           # Graph algorithms
│   ├── graph.go    # Graph representation
│   ├── errors.go   # Graph errors (cycles, etc.)
│   └── *_test.go   # Algorithm tests
│
├── state/           # State management
│   ├── manager.go  # State manager
│   ├── errors.go   # State errors
│   └── doc.go      # Documentation
│
├── security/        # Security utilities
│   ├── ssrf.go     # SSRF protection
│   └── *_test.go   # Security tests
│
├── logging/         # Structured logging
│   ├── logger.go   # Logger implementation
│   └── *_test.go   # Logging tests
│
├── middleware/      # Execution middleware
│   ├── middleware.go     # Base middleware
│   ├── validation.go     # Validation middleware
│   ├── logging.go        # Logging middleware
│   ├── ratelimit.go      # Rate limiting
│   ├── timeout.go        # Timeout protection
│   └── *_test.go         # Middleware tests
│
├── observer/        # Observer pattern
│   ├── observer.go # Observer interface
│   ├── defaults.go # Default observers
│   └── *_test.go   # Observer tests
│
├── httpclient/      # HTTP client
│   ├── client.go   # HTTP client
│   ├── middleware.go # HTTP middleware
│   ├── ssrf.go     # SSRF integration
│   ├── registry.go # Named clients
│   └── *_test.go   # HTTP tests
│
├── expression/      # Expression evaluation
│   ├── expression.go # Expression parser
│   └── *_test.go     # Expression tests
│
└── types/           # Shared types
    ├── types.go    # Core type definitions
    ├── helpers.go  # Type utilities
    └── doc.go      # Type documentation
```

## Data Flow

### Workflow Creation Flow

```
┌──────────┐
│  User    │
└────┬─────┘
     │
     │ 1. Drag & Drop Nodes
     ↓
┌──────────────────┐
│  React Flow      │
│  Visual Editor   │
└────┬─────────────┘
     │
     │ 2. Create/Update Nodes
     ↓
┌──────────────────┐
│  Workflow State  │
│  (React State)   │
└────┬─────────────┘
     │
     │ 3. Serialize to JSON
     ↓
┌──────────────────┐
│  JSON Payload    │
│  {nodes, edges}  │
└────┬─────────────┘
     │
     │ 4. Store/Execute
     ↓
┌──────────────────┐
│  LocalStorage/   │
│  Backend API     │
└──────────────────┘
```

### Workflow Execution Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    Execution Pipeline                        │
└─────────────────────────────────────────────────────────────┘

1. JSON Payload
   ↓
2. Parse & Validate
   ├─ Parse JSON
   ├─ Validate structure
   └─ Check for cycles
   ↓
3. Type Inference
   ├─ Infer missing types
   └─ Set default values
   ↓
4. Topological Sort
   ├─ Build dependency graph
   ├─ Detect cycles
   └─ Determine execution order
   ↓
5. Create Execution Context
   ├─ Initialize state manager
   ├─ Set up observers
   ├─ Configure timeouts
   └─ Apply resource limits
   ↓
6. Execute Nodes (in order)
   ├─ For each node:
   │  ├─ Apply middleware (pre)
   │  ├─ Get inputs from predecessors
   │  ├─ Interpolate templates
   │  ├─ Execute node logic
   │  ├─ Apply middleware (post)
   │  ├─ Store result
   │  └─ Notify observers
   ↓
7. Return Result
   ├─ Node results map
   ├─ Final output
   ├─ Execution metadata
   └─ Any errors
```

### Node Execution Detail

```
┌─────────────────────────────────────────────────────┐
│              Node Execution Flow                     │
└─────────────────────────────────────────────────────┘

1. Get Node from Execution Queue
   ↓
2. Check Protection Limits
   ├─ Node execution count
   ├─ HTTP call count
   └─ Timeout
   ↓
3. Apply Pre-Execution Middleware
   ├─ Validation
   ├─ Logging
   ├─ Rate limiting
   └─ Input validation
   ↓
4. Get Inputs
   ├─ Retrieve from previous nodes
   └─ Validate input types
   ↓
5. Interpolate Templates
   ├─ Replace {{ variable.name }}
   └─ Replace {{ const.name }}
   ↓
6. Execute Node Logic
   ├─ Dispatch to executor
   ├─ Process inputs
   └─ Generate output
   ↓
7. Apply Post-Execution Middleware
   ├─ Output validation
   ├─ Size limits
   └─ Logging
   ↓
8. Store Result
   ├─ Node results map
   └─ State updates
   ↓
9. Notify Observers
   ├─ Success event
   └─ Or failure event
```

## Security Architecture

### Defense in Depth

```
┌───────────────────────────────────────────────────────────┐
│                   Security Layers                          │
├───────────────────────────────────────────────────────────┤
│                                                            │
│  Layer 1: Input Validation                                │
│  ├─ JSON schema validation                                │
│  ├─ Type checking                                         │
│  ├─ Size limits                                           │
│  └─ Sanitization                                          │
│                                                            │
│  Layer 2: Resource Protection                             │
│  ├─ Execution timeout                                     │
│  ├─ Node execution limits                                 │
│  ├─ Memory limits                                         │
│  ├─ Loop iteration limits                                 │
│  └─ Variable count limits                                 │
│                                                            │
│  Layer 3: Network Security                                │
│  ├─ SSRF protection                                       │
│  ├─ URL validation                                        │
│  ├─ Private IP blocking                                   │
│  ├─ Localhost blocking                                    │
│  └─ Cloud metadata blocking                               │
│                                                            │
│  Layer 4: Execution Isolation                             │
│  ├─ Context isolation                                     │
│  ├─ State sandboxing                                      │
│  └─ Error containment                                     │
│                                                            │
│  Layer 5: Observability                                   │
│  ├─ Structured logging                                    │
│  ├─ Execution tracing                                     │
│  ├─ Metrics collection                                    │
│  └─ Security event monitoring                             │
│                                                            │
└───────────────────────────────────────────────────────────┘
```

### SSRF Protection

```
HTTP Request Flow with SSRF Protection:

┌─────────────┐
│ HTTP Node   │
└──────┬──────┘
       │
       ↓
┌──────────────────┐
│ Validate URL     │
├──────────────────┤
│ • Check scheme   │
│ • Parse hostname │
│ • Check blocklist│
└──────┬───────────┘
       │
       ↓
┌──────────────────┐
│ Resolve IP       │
├──────────────────┤
│ • DNS lookup     │
│ • Get IPs        │
└──────┬───────────┘
       │
       ↓
┌──────────────────┐
│ Validate IPs     │
├──────────────────┤
│ • Block private  │
│ • Block localhost│
│ • Block metadata │
│ • Block link-lcl │
└──────┬───────────┘
       │
       ↓ (allowed)
┌──────────────────┐
│ Make Request     │
├──────────────────┤
│ • Apply timeout  │
│ • Size limits    │
│ • Follow redirects (with revalidation)
└──────┬───────────┘
       │
       ↓
┌──────────────────┐
│ Return Response  │
└──────────────────┘
```

## Deployment Architecture

### Development Environment

```
┌─────────────────────────────────────────┐
│         Developer Workstation            │
├─────────────────────────────────────────┤
│                                          │
│  ┌────────────────────────────────────┐ │
│  │  Next.js Dev Server (Port 3000)    │ │
│  │  • Hot reload                       │ │
│  │  • Fast refresh                     │ │
│  └────────────────────────────────────┘ │
│                                          │
│  ┌────────────────────────────────────┐ │
│  │  Go Development                     │ │
│  │  • Direct execution                 │ │
│  │  • Test runner                      │ │
│  │  • Debugger (delve)                 │ │
│  └────────────────────────────────────┘ │
│                                          │
└─────────────────────────────────────────┘
```

### Production Deployment Options

#### Option 1: Vercel Deployment

```
┌──────────────────────────────────────────────┐
│              Vercel Platform                  │
├──────────────────────────────────────────────┤
│                                               │
│  ┌─────────────────────────────────────────┐ │
│  │  Next.js App (Edge Functions)           │ │
│  │  • SSR/SSG                               │ │
│  │  • API Routes                            │ │
│  │  • CDN Distribution                      │ │
│  └─────────────────────────────────────────┘ │
│                                               │
└──────────────────────────────────────────────┘
         │
         │ API Calls (if backend separated)
         ↓
┌──────────────────────────────────────────────┐
│         Separate Go Backend Service           │
│  (Cloud Run / ECS / Kubernetes)              │
└──────────────────────────────────────────────┘
```

#### Option 2: Container Deployment

```
┌────────────────────────────────────────────┐
│          Kubernetes Cluster                 │
├────────────────────────────────────────────┤
│                                             │
│  ┌───────────────────────────────────────┐ │
│  │  Frontend Pods                        │ │
│  │  ├─ Next.js Container                 │ │
│  │  └─ Nginx (optional)                  │ │
│  └───────────────────────────────────────┘ │
│                                             │
│  ┌───────────────────────────────────────┐ │
│  │  Backend Pods                         │ │
│  │  ├─ Go Workflow Engine                │ │
│  │  └─ Auto-scaling                      │ │
│  └───────────────────────────────────────┘ │
│                                             │
│  ┌───────────────────────────────────────┐ │
│  │  Supporting Services                  │ │
│  │  ├─ Redis (caching)                   │ │
│  │  ├─ PostgreSQL (persistence)          │ │
│  │  └─ Prometheus (metrics)              │ │
│  └───────────────────────────────────────┘ │
│                                             │
└────────────────────────────────────────────┘
```

## Package Dependencies

### Backend Dependency Graph

```
types (foundation)
  ↓
  ├─→ config
  ├─→ logging
  ├─→ security
  ├─→ expression
  └─→ graph
       ↓
       ├─→ state
       ├─→ observer
       ├─→ httpclient
       └─→ middleware
            ↓
            └─→ executor
                 ↓
                 └─→ engine
```

**Dependency Rules:**
1. `types` package has no dependencies (foundation)
2. Packages only depend on layers below them
3. No circular dependencies allowed
4. Minimal external dependencies

### External Dependencies

**Go Dependencies:**
- Standard library only (no external dependencies for core functionality)
- Test dependencies: `testify`, `require` (testing only)

**Frontend Dependencies:**
- `next`: React framework
- `react`, `react-dom`: UI library
- `reactflow`: Visual workflow editor
- `tailwindcss`: Styling

## Scalability Considerations

### Horizontal Scaling

```
       ┌─────────────────┐
       │  Load Balancer  │
       └────────┬────────┘
                │
    ┌───────────┼───────────┐
    │           │           │
    ↓           ↓           ↓
┌────────┐  ┌────────┐  ┌────────┐
│Engine 1│  │Engine 2│  │Engine N│
└────────┘  └────────┘  └────────┘
    │           │           │
    └───────────┼───────────┘
                ↓
       ┌─────────────────┐
       │  Shared Storage │
       └─────────────────┘
```

**Scaling Strategies:**
- Stateless engine instances
- Shared cache layer (Redis)
- Workflow persistence (Database)
- Result streaming for large outputs

### Performance Optimization

1. **Caching**
   - Workflow definition caching
   - Compiled expression caching
   - HTTP response caching

2. **Parallelization**
   - Independent node parallel execution
   - Worker pool for concurrent workflows

3. **Resource Management**
   - Connection pooling
   - Memory pooling
   - Goroutine limits

## Related Documentation

- [Design Patterns](ARCHITECTURE_DESIGN_PATTERNS.md)
- [Security Architecture](PRINCIPLES_ZERO_TRUST.md)
- [Performance Tuning](PERFORMANCE_TUNING.md)
- [Deployment Guide](REQUIREMENTS_NON_FUNCTIONAL_DEPLOYMENT.md)

---

**Last Updated:** 2025-11-03
**Version:** 1.0
