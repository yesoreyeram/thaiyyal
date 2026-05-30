# Thaiyyal - Master Product & Technical Specification

**Version:** 1.0
**Last Updated:** 2026-05-28
**Status:** Active

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Vision & Goals](#2-vision--goals)
3. [Problem Statement](#3-problem-statement)
4. [Target Users & Use Cases](#4-target-users--use-cases)
5. [Functional Requirements](#5-functional-requirements)
6. [Non-Functional Requirements](#6-non-functional-requirements)
7. [Technical Architecture](#7-technical-architecture)
8. [API Specification](#8-api-specification)
9. [Security Architecture](#9-security-architecture)
10. [Testing Strategy](#10-testing-strategy)
11. [Deployment & Operations](#11-deployment--operations)
12. [Development Roadmap](#12-development-roadmap)
13. [Success Metrics](#13-success-metrics)

---

## 1. Executive Summary

### 1.1 Product Overview

Thaiyyal is a **visual workflow engine** that empowers users to create, execute, and manage complex data processing workflows through an intuitive drag-and-drop interface. The system combines a modern Next.js-based visual editor with a robust Go backend execution engine, delivering a secure, extensible, and performant solution for workflow automation.

### 1.2 Key Differentiators

- **Visual-First Design**: Intuitive drag-and-drop interface eliminates the need for coding
- **Security by Default**: Zero-trust architecture with comprehensive SSRF protection and resource limits
- **High Performance**: Go-powered execution engine supporting 10,000+ nodes per workflow
- **Extensibility**: Plugin architecture allowing custom node types and middleware
- **Enterprise-Ready**: Built-in observability, health checks, and deployment templates

### 1.3 Current Status

- **Version**: 0.1.0 (Active Development)
- **Test Coverage**: 95%+ across all backend packages
- **Supported Node Types**: 40+ built-in node types across 8 categories
- **Deployment**: Production-ready with Docker and Kubernetes support

### 1.4 Technology Stack

- **Frontend**: Next.js 16.0, React 19.2, React Flow 11.8, Tailwind CSS 4.0
- **Backend**: Go 1.24.7, Standard Library (zero external dependencies for core)
- **Infrastructure**: Docker, Kubernetes, Prometheus, OpenTelemetry

---

## 2. Vision & Goals

### 2.1 Vision Statement

To become the leading open-source visual workflow engine that democratizes workflow automation by making it accessible, secure, and powerful for users of all technical skill levels.

### 2.2 Mission

Empower individuals and organizations to automate complex data processing tasks through visual workflows without requiring programming expertise, while maintaining enterprise-grade security and performance.

### 2.3 Strategic Goals

#### Short-Term Goals (3-6 months)
- Achieve 1.0 release with stable API
- Complete workflow persistence layer
- Implement execution history API
- Add workflow versioning support
- Build comprehensive documentation and tutorials

#### Medium-Term Goals (6-12 months)
- Reach 1,000 active users
- Support 100+ node types (including community contributions)
- Implement GraphQL API
- Add WebSocket-based real-time execution
- Create workflow templates marketplace

#### Long-Term Goals (12+ months)
- Establish active community of 10,000+ users
- Support multi-language expression evaluation
- Implement visual debugging tools
- Enable enterprise features (RBAC, audit logs, compliance)
- Build integrations with major platforms (AWS, GCP, Azure, etc.)

### 2.4 Core Principles

1. **Security First**: All features must meet zero-trust security standards
2. **No Runtime Errors**: Comprehensive validation prevents runtime failures
3. **Modularity**: Every component should be independently testable and replaceable
4. **Observability**: All system behavior must be observable and measurable
5. **Developer Experience**: API design prioritizes clarity and ease of use

---

## 3. Problem Statement

### 3.1 The Challenge

Organizations and individuals face several challenges when automating data processing workflows:

1. **Complexity Barrier**: Traditional workflow tools require programming knowledge
2. **Security Risks**: Many solutions lack proper security controls (SSRF, injection, resource exhaustion)
3. **Vendor Lock-in**: Proprietary platforms limit flexibility and portability
4. **Debugging Difficulty**: Limited visibility into workflow execution
5. **Integration Overhead**: Connecting different data sources and services is time-consuming

### 3.2 Current Solutions & Limitations

| Solution Type | Examples | Limitations |
|--------------|----------|------------|
| **Low-Code Platforms** | Zapier, n8n | - Expensive<br>- Limited customization<br>- Vendor lock-in |
| **Enterprise Tools** | Apache Airflow, Prefect | - Complex setup<br>- Requires coding<br>- Steep learning curve |
| **Visual Tools** | Node-RED | - Limited security<br>- Performance issues<br>- Weak typing |
| **Custom Solutions** | In-house scripts | - High maintenance<br>- No standardization<br>- Security gaps |

### 3.3 Thaiyyal's Solution

Thaiyyal addresses these challenges by providing:

1. **Visual Interface**: Drag-and-drop workflow creation for all skill levels
2. **Zero-Trust Security**: Comprehensive protection against common vulnerabilities
3. **Open Source**: No vendor lock-in, full transparency and customizability
4. **Observability**: Built-in logging, metrics, and tracing
5. **Extensibility**: Plugin architecture for custom integrations
6. **Type Safety**: Strong typing prevents common errors
7. **Performance**: Optimized Go engine handles complex workflows efficiently

---

## 4. Target Users & Use Cases

### 4.1 Primary User Personas

#### 4.1.1 Data Analyst (Sarah)
- **Background**: Works with data but limited programming experience
- **Goals**: Process and transform data from multiple sources
- **Pain Points**: Excel limitations, manual data manipulation
- **How Thaiyyal Helps**: Visual interface for creating ETL pipelines without coding

#### 4.1.2 DevOps Engineer (Mike)
- **Background**: Manages infrastructure and automation
- **Goals**: Automate operational workflows and monitoring
- **Pain Points**: Maintaining complex scripts, lack of visibility
- **How Thaiyyal Helps**: Visual workflows with built-in observability and security

#### 4.1.3 Backend Developer (Alex)
- **Background**: Builds APIs and backend services
- **Goals**: Quickly prototype data processing logic
- **Pain Points**: Boilerplate code, testing complexity
- **How Thaiyyal Helps**: Extensible platform with custom executors and type safety

#### 4.1.4 Integration Specialist (Jordan)
- **Background**: Connects different systems and services
- **Goals**: Build reliable integrations between platforms
- **Pain Points**: Error handling, rate limiting, monitoring
- **How Thaiyyal Helps**: Built-in HTTP clients, retry logic, and observability

### 4.2 Use Cases

#### 4.2.1 Data Transformation Pipeline
**Scenario**: Sarah needs to fetch data from multiple APIs, transform it, and generate reports.

**Workflow**:
1. HTTP nodes fetch data from 3 different APIs
2. Extract nodes pull relevant fields
3. Transform nodes normalize data formats
4. Join node combines datasets
5. Filter/map nodes clean and enrich data
6. Visualization node displays results

**Benefits**: No coding required, visual debugging, reusable workflows

#### 4.2.2 Monitoring & Alerting
**Scenario**: Mike wants to monitor service health and send alerts.

**Workflow**:
1. HTTP nodes check service endpoints
2. Conditional nodes evaluate health status
3. Accumulator tracks failures over time
4. Try/catch handles API failures gracefully
5. HTTP node sends alerts to Slack/PagerDuty

**Benefits**: Self-documenting, easy to modify, observable execution

#### 4.2.3 API Development Prototype
**Scenario**: Alex needs to prototype complex business logic quickly.

**Workflow**:
1. Input nodes accept request parameters
2. Variable nodes manage state
3. Operation nodes implement business rules
4. Conditional nodes handle edge cases
5. Transform node formats response

**Benefits**: Rapid iteration, visual testing, easy conversion to code

#### 4.2.4 ETL Automation
**Scenario**: Jordan integrates CRM data with analytics platform.

**Workflow**:
1. HTTP node fetches from Salesforce API
2. Parse node handles JSON/XML
3. Map/filter nodes transform data
4. Chunk node batches for API limits
5. ForEach with retry posts to analytics platform
6. Accumulator tracks success/failure rates

**Benefits**: Resilient execution, built-in rate limiting, comprehensive logging

---

## 5. Functional Requirements

### 5.1 Workflow Management

#### FR-1.1: Workflow Creation
- **REQ-1.1.1**: System SHALL allow creating workflows via visual drag-and-drop editor
- **REQ-1.1.2**: System SHALL support JSON-based workflow definition for programmatic creation
- **REQ-1.1.3**: System SHALL validate workflow structure (DAG, no cycles)
- **REQ-1.1.4**: System SHALL provide real-time validation feedback in the editor
- **REQ-1.1.5**: System SHALL support workflow templates for common patterns

#### FR-1.2: Workflow Execution
- **REQ-1.2.1**: System SHALL execute workflows in topological order based on node dependencies
- **REQ-1.2.2**: System SHALL execute independent nodes in parallel where possible
- **REQ-1.2.3**: System SHALL enforce configurable timeout limits (default: 30s, max: 10m)
- **REQ-1.2.4**: System SHALL support workflow cancellation during execution
- **REQ-1.2.5**: System SHALL provide real-time execution feedback

#### FR-1.3: Workflow Persistence
- **REQ-1.3.1**: System SHALL save workflows to browser local storage
- **REQ-1.3.2**: System SHALL support workflow import/export as JSON
- **REQ-1.3.3**: System SHALL support workflow versioning (planned)
- **REQ-1.3.4**: System SHALL preserve workflow metadata (name, description, timestamps)

### 5.2 Node Types

The system SHALL support the following node categories:

#### FR-2.1: Basic I/O Nodes
- **Number Input**: Provide numeric constants
- **Text Input**: Provide string constants
- **Visualization**: Display results (text, table, JSON, chart modes)

#### FR-2.2: Operation Nodes
- **Mathematical Operations**: add, subtract, multiply, divide, modulo, power
- **Text Operations**: concat, uppercase, lowercase, trim, replace, split
- **HTTP Requests**: GET, POST, PUT, DELETE with SSRF protection

#### FR-2.3: Control Flow Nodes
- **Condition**: If/else branching based on expressions
- **Switch**: Multi-way branching
- **ForEach**: Iterate over arrays
- **While**: Loop with condition
- **Parallel**: Execute multiple branches concurrently
- **Join**: Wait for multiple inputs

#### FR-2.4: Array Processing Nodes
- **Filter**: Filter elements by expression
- **Map**: Transform elements
- **Reduce**: Aggregate to single value
- **Sort**: Sort by field or expression
- **Slice**: Extract sub-array
- **Find**: Locate element
- **FlatMap**: Map and flatten
- **GroupBy**: Group by field
- **Unique**: Remove duplicates
- **Chunk**: Split into batches
- **Reverse**: Reverse order
- **Partition**: Split by predicate
- **Zip**: Combine arrays
- **Sample**: Random sampling
- **Range**: Generate number sequence
- **Compact**: Remove null/undefined
- **Transpose**: Matrix transpose

#### FR-2.5: State Management Nodes
- **Variable**: Get/set workflow variables
- **Extract**: Extract fields from objects
- **Transform**: Transform object structure
- **Accumulator**: Accumulate values across iterations
- **Counter**: Increment/decrement counters
- **Parse**: Parse JSON, CSV, XML
- **Cache**: Cache values with TTL

#### FR-2.6: Error Handling Nodes
- **Try/Catch**: Handle errors with fallback
- **Retry**: Retry with exponential backoff
- **Timeout**: Enforce operation timeout

#### FR-2.7: Context Nodes
- **Context Variable**: Access workflow variables
- **Context Constant**: Access workflow constants

### 5.3 Data Processing

#### FR-3.1: Type System
- **REQ-3.1.1**: System SHALL support primitive types: string, number, boolean, null
- **REQ-3.1.2**: System SHALL support complex types: array, object
- **REQ-3.1.3**: System SHALL perform automatic type inference where possible
- **REQ-3.1.4**: System SHALL validate type compatibility between connected nodes

#### FR-3.2: Expression Engine
- **REQ-3.2.1**: System SHALL support JavaScript-like expression syntax
- **REQ-3.2.2**: System SHALL support arithmetic operators: +, -, *, /, %, **
- **REQ-3.2.3**: System SHALL support comparison operators: ==, !=, <, >, <=, >=
- **REQ-3.2.4**: System SHALL support logical operators: &&, ||, !
- **REQ-3.2.5**: System SHALL support field access: object.field, array[index]
- **REQ-3.2.6**: System SHALL support function calls: built-in functions

#### FR-3.3: Template Interpolation
- **REQ-3.3.1**: System SHALL support variable interpolation: {{ variable.name }}
- **REQ-3.3.2**: System SHALL support constant interpolation: {{ const.name }}
- **REQ-3.3.3**: System SHALL validate interpolated values at execution time

### 5.4 Integration

#### FR-4.1: HTTP Integration
- **REQ-4.1.1**: System SHALL support HTTP methods: GET, POST, PUT, DELETE, PATCH
- **REQ-4.1.2**: System SHALL support custom headers
- **REQ-4.1.3**: System SHALL support query parameters
- **REQ-4.1.4**: System SHALL support request body (JSON, form data)
- **REQ-4.1.5**: System SHALL parse response automatically (JSON, text)
- **REQ-4.1.6**: System SHALL follow redirects with re-validation
- **REQ-4.1.7**: System SHALL enforce HTTP call limits per workflow (default: 100)

#### FR-4.2: Extensibility
- **REQ-4.2.1**: System SHALL support custom node executors via plugin interface
- **REQ-4.2.2**: System SHALL support custom middleware for cross-cutting concerns
- **REQ-4.2.3**: System SHALL support custom observers for monitoring
- **REQ-4.2.4**: System SHALL support named HTTP clients with custom configurations

### 5.5 User Interface

#### FR-5.1: Visual Editor
- **REQ-5.1.1**: System SHALL provide drag-and-drop node placement
- **REQ-5.1.2**: System SHALL support edge creation between nodes
- **REQ-5.1.3**: System SHALL support node configuration via properties panel
- **REQ-5.1.4**: System SHALL support canvas pan and zoom
- **REQ-5.1.5**: System SHALL auto-layout nodes (optional)
- **REQ-5.1.6**: System SHALL highlight validation errors visually

#### FR-5.2: Execution Results
- **REQ-5.2.1**: System SHALL display execution summary (duration, node count, status)
- **REQ-5.2.2**: System SHALL display final workflow output
- **REQ-5.2.3**: System SHALL display per-node results
- **REQ-5.2.4**: System SHALL display errors with context
- **REQ-5.2.5**: System SHALL support resizable results panel (100px-600px)

---

## 6. Non-Functional Requirements

### 6.1 Performance

#### NFR-1.1: Execution Performance
- **REQ-6.1.1**: Workflow execution SHALL complete within configured timeout (default: 30s)
- **REQ-6.1.2**: System SHALL support workflows with 10,000+ nodes
- **REQ-6.1.3**: Node execution overhead SHALL be <1ms per node
- **REQ-6.1.4**: Memory usage SHALL scale linearly with workflow size
- **REQ-6.1.5**: System SHALL support 100+ concurrent workflow executions

#### NFR-1.2: UI Performance
- **REQ-6.1.6**: Editor SHALL render workflows with 1,000+ nodes smoothly (60fps)
- **REQ-6.1.7**: Node drag operations SHALL have <16ms latency
- **REQ-6.1.8**: Workflow validation SHALL complete within 100ms
- **REQ-6.1.9**: Initial page load SHALL complete within 2s

### 6.2 Security

#### NFR-2.1: Input Validation
- **REQ-6.2.1**: System SHALL validate all user inputs before processing
- **REQ-6.2.2**: System SHALL enforce maximum string length (default: 1MB)
- **REQ-6.2.3**: System SHALL enforce maximum array size (default: 10,000 elements)
- **REQ-6.2.4**: System SHALL enforce maximum object depth (default: 32 levels)
- **REQ-6.2.5**: System SHALL sanitize all outputs for XSS prevention

#### NFR-2.2: SSRF Protection
- **REQ-6.2.6**: System SHALL block private IP ranges (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16)
- **REQ-6.2.7**: System SHALL block localhost and loopback (127.0.0.0/8, ::1)
- **REQ-6.2.8**: System SHALL block link-local addresses (169.254.0.0/16)
- **REQ-6.2.9**: System SHALL block cloud metadata endpoints (169.254.169.254)
- **REQ-6.2.10**: System SHALL re-validate URLs after redirects
- **REQ-6.2.11**: System SHALL support domain allow/block lists

#### NFR-2.3: Resource Protection
- **REQ-6.2.12**: System SHALL enforce maximum execution time per workflow
- **REQ-6.2.13**: System SHALL enforce maximum node executions per workflow (default: 10,000)
- **REQ-6.2.14**: System SHALL enforce maximum HTTP calls per workflow (default: 100)
- **REQ-6.2.15**: System SHALL enforce maximum loop iterations (default: 10,000)
- **REQ-6.2.16**: System SHALL enforce maximum recursion depth (default: 100)
- **REQ-6.2.17**: System SHALL enforce maximum variables count (default: 1,000)

### 6.3 Reliability

#### NFR-3.1: Availability
- **REQ-6.3.1**: System SHALL target 99.9% uptime
- **REQ-6.3.2**: System SHALL gracefully degrade on partial failures
- **REQ-6.3.3**: System SHALL provide health check endpoints
- **REQ-6.3.4**: System SHALL support liveness and readiness probes

#### NFR-3.2: Error Handling
- **REQ-6.3.5**: System SHALL never crash on invalid input
- **REQ-6.3.6**: System SHALL provide actionable error messages
- **REQ-6.3.7**: System SHALL log all errors with context
- **REQ-6.3.8**: System SHALL support error recovery via try/catch nodes

### 6.4 Observability

#### NFR-4.1: Logging
- **REQ-6.4.1**: System SHALL use structured logging (JSON format)
- **REQ-6.4.2**: System SHALL log all workflow executions
- **REQ-6.4.3**: System SHALL log security events (SSRF blocks, limit violations)
- **REQ-6.4.4**: System SHALL support configurable log levels

#### NFR-4.2: Metrics
- **REQ-6.4.5**: System SHALL expose Prometheus metrics
- **REQ-6.4.6**: System SHALL track workflow execution counts (success/failure)
- **REQ-6.4.7**: System SHALL track execution duration (histogram)
- **REQ-6.4.8**: System SHALL track node execution counts per type
- **REQ-6.4.9**: System SHALL track HTTP call counts and duration

#### NFR-4.3: Tracing
- **REQ-6.4.10**: System SHALL support OpenTelemetry distributed tracing
- **REQ-6.4.11**: System SHALL create spans for workflow execution
- **REQ-6.4.12**: System SHALL create spans for node execution
- **REQ-6.4.13**: System SHALL propagate trace context

### 6.5 Maintainability

#### NFR-5.1: Code Quality
- **REQ-6.5.1**: Test coverage SHALL be >80% (target: >95%)
- **REQ-6.5.2**: Code SHALL pass all linter checks
- **REQ-6.5.3**: Code SHALL follow language-specific style guides
- **REQ-6.5.4**: All public APIs SHALL have documentation
- **REQ-6.5.5**: Complex logic SHALL have explanatory comments

#### NFR-5.2: Documentation
- **REQ-6.5.6**: All features SHALL have user documentation
- **REQ-6.5.7**: All APIs SHALL have reference documentation
- **REQ-6.5.8**: Architecture decisions SHALL be documented
- **REQ-6.5.9**: Security considerations SHALL be documented

### 6.6 Scalability

#### NFR-6.1: Horizontal Scaling
- **REQ-6.6.1**: Backend SHALL be stateless for horizontal scaling
- **REQ-6.6.2**: System SHALL support shared state via external storage
- **REQ-6.6.3**: System SHALL support load balancing across instances

#### NFR-6.2: Resource Efficiency
- **REQ-6.6.4**: System SHALL use connection pooling for HTTP clients
- **REQ-6.6.5**: System SHALL support result streaming for large outputs
- **REQ-6.6.6**: System SHALL implement caching where appropriate

### 6.7 Usability

#### NFR-7.1: User Experience
- **REQ-6.7.1**: New users SHALL be able to create basic workflow within 5 minutes
- **REQ-6.7.2**: Error messages SHALL be clear and actionable
- **REQ-6.7.3**: UI SHALL be responsive on desktop (1024px+)
- **REQ-6.7.4**: UI SHALL provide keyboard shortcuts for common actions

---

## 7. Technical Architecture

### 7.1 System Architecture

Thaiyyal uses a **hybrid architecture** combining:
- **Frontend**: Next.js-based single-page application
- **Backend**: Go-based workflow execution engine
- **Communication**: HTTP/JSON API

See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture documentation.

### 7.2 Component Overview

```
┌─────────────────────────────────────────────────────┐
│                 Frontend (Next.js)                   │
│  ┌────────────────────────────────────────────────┐ │
│  │  UI Components (React + Tailwind)              │ │
│  │  React Flow Editor                             │ │
│  │  State Management                              │ │
│  └────────────────────────────────────────────────┘ │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP/JSON
                       ↓
┌─────────────────────────────────────────────────────┐
│              Backend (Go Engine)                     │
│  ┌────────────────────────────────────────────────┐ │
│  │  HTTP Server (cmd/server)                      │ │
│  │  ├─ /api/v1/workflow/execute                   │ │
│  │  ├─ /api/v1/workflow/validate                  │ │
│  │  └─ /health, /metrics                          │ │
│  └────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────┐ │
│  │  Engine (pkg/engine)                           │ │
│  │  ├─ Workflow parsing & validation              │ │
│  │  ├─ Topological sorting                        │ │
│  │  └─ Execution orchestration                    │ │
│  └────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────┐ │
│  │  Executors (pkg/executor)                      │ │
│  │  ├─ Registry (40+ node types)                  │ │
│  │  └─ Custom executor support                    │ │
│  └────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────┐ │
│  │  Support Services                              │ │
│  │  ├─ State (pkg/state)                          │ │
│  │  ├─ Graph (pkg/graph)                          │ │
│  │  ├─ Security (pkg/security)                    │ │
│  │  ├─ HTTP Client (pkg/httpclient)               │ │
│  │  ├─ Middleware (pkg/middleware)                │ │
│  │  ├─ Observer (pkg/observer)                    │ │
│  │  └─ Telemetry (pkg/telemetry)                  │ │
│  └────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

### 7.3 Backend Package Structure

```
backend/pkg/
├── types/           # Core type definitions (zero dependencies)
├── config/          # Configuration management
├── logging/         # Structured logging (slog)
├── security/        # SSRF protection
├── expression/      # Expression evaluation
├── graph/           # Topological sort algorithms
├── state/           # State management
├── observer/        # Observer pattern
├── httpclient/      # HTTP client with middleware
├── middleware/      # Execution middleware
├── executor/        # Node executors & registry
├── engine/          # Workflow execution engine
├── server/          # HTTP server
├── health/          # Health checks
└── telemetry/       # OpenTelemetry integration
```

### 7.4 Data Flow

#### Workflow Creation
1. User interacts with React Flow editor
2. Node/edge changes update React state
3. State serialized to JSON on save
4. JSON stored in browser localStorage

#### Workflow Execution
1. User clicks "Run" button
2. Frontend sends JSON payload to backend API
3. Backend parses and validates JSON
4. Engine performs topological sort
5. Engine executes nodes in order
6. Results streamed back to frontend
7. Frontend displays results in panel

### 7.5 Technology Decisions

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Frontend Framework | Next.js 16 | SSR/SSG support, optimal performance, React Server Components |
| UI Library | React 19 | Modern hooks, concurrent features, large ecosystem |
| Visual Editor | React Flow 11 | Mature, performant, customizable graph editor |
| Styling | Tailwind CSS 4 | Utility-first, minimal bundle size, responsive |
| Backend Language | Go 1.24 | Performance, concurrency, static typing, simple deployment |
| HTTP Server | net/http (stdlib) | Zero dependencies, production-ready, well-tested |
| Logging | slog (stdlib) | Structured logging, performance, standard library |
| Testing | Go testing + testify | Standard tools, comprehensive mocking |
| Metrics | Prometheus | Industry standard, rich ecosystem |
| Tracing | OpenTelemetry | Vendor-neutral, comprehensive instrumentation |
| Containerization | Docker | Standard deployment, multi-stage builds |
| Orchestration | Kubernetes | Scalability, health checks, auto-scaling |

---

## 8. API Specification

### 8.1 HTTP API Endpoints

#### POST /api/v1/workflow/execute

Execute a workflow.

**Request:**
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {"value": 10}
    }
  ],
  "edges": [],
  "context": {
    "variables": {},
    "constants": {}
  }
}
```

**Response (Success):**
```json
{
  "success": true,
  "output": 10,
  "results": {
    "1": 10
  },
  "metadata": {
    "duration_ms": 5,
    "nodes_executed": 1
  }
}
```

**Response (Error):**
```json
{
  "success": false,
  "error": "validation error: cycle detected",
  "details": {
    "type": "validation_error",
    "cycle": ["1", "2", "1"]
  }
}
```

#### POST /api/v1/workflow/validate

Validate a workflow without executing.

**Request:** Same as execute

**Response:**
```json
{
  "valid": true,
  "execution_order": ["1", "2", "3"]
}
```

#### GET /health

Comprehensive health check.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2026-05-28T22:55:00Z",
  "version": "0.1.0"
}
```

#### GET /health/live

Liveness probe (Kubernetes).

**Response:** 200 OK or 503 Service Unavailable

#### GET /health/ready

Readiness probe (Kubernetes).

**Response:** 200 OK or 503 Service Unavailable

#### GET /metrics

Prometheus metrics endpoint.

**Response:** Prometheus text format

### 8.2 Workflow JSON Schema

See [DATA_FORMAT_GUIDE.md](docs/DATA_FORMAT_GUIDE.md) for complete schema.

**Core Structure:**
```typescript
interface Workflow {
  nodes: Node[];
  edges: Edge[];
  context?: Context;
}

interface Node {
  id: string;
  type: string;
  data?: Record<string, any>;
}

interface Edge {
  source: string;
  target: string;
  sourceHandle?: string;
}

interface Context {
  variables?: Record<string, any>;
  constants?: Record<string, any>;
}
```

### 8.3 Error Codes

| Code | Description |
|------|-------------|
| `validation_error` | Invalid workflow structure |
| `cycle_detected` | Workflow contains cycle |
| `execution_timeout` | Execution exceeded timeout |
| `node_limit_exceeded` | Too many node executions |
| `http_limit_exceeded` | Too many HTTP calls |
| `type_error` | Type mismatch |
| `runtime_error` | Node execution failed |
| `ssrf_blocked` | SSRF protection triggered |

---

## 9. Security Architecture

### 9.1 Security Model

Thaiyyal implements a **zero-trust security model**:
- All inputs are untrusted and validated
- All outputs are sanitized
- All external requests are protected
- All resources are limited

### 9.2 Threat Model

#### Threats Addressed

1. **Server-Side Request Forgery (SSRF)**
   - Mitigation: IP validation, domain filtering, redirect validation

2. **Denial of Service (DoS)**
   - Mitigation: Execution timeouts, node limits, loop limits, HTTP call limits

3. **Code Injection**
   - Mitigation: Expression sandboxing, input validation, no eval()

4. **Data Exfiltration**
   - Mitigation: Output size limits, logging, monitoring

5. **Resource Exhaustion**
   - Mitigation: Memory limits, goroutine limits, connection pooling

### 9.3 SSRF Protection

**Implementation Layers:**

1. **URL Validation**
   - Scheme must be http/https
   - Parse and validate hostname

2. **DNS Resolution**
   - Resolve hostname to IP addresses
   - Validate all IPs before connection

3. **IP Validation**
   - Block private ranges (RFC 1918)
   - Block loopback (127.0.0.0/8, ::1)
   - Block link-local (169.254.0.0/16)
   - Block cloud metadata (169.254.169.254)

4. **Redirect Validation**
   - Re-validate URLs after redirects
   - Prevent bypass via redirects

5. **Allow/Block Lists**
   - Support domain allowlist
   - Support domain blocklist

**Code Location:** `backend/pkg/security/ssrf.go`

### 9.4 Resource Limits

| Resource | Default Limit | Configurable | Rationale |
|----------|--------------|--------------|-----------|
| Execution Time | 30s | Yes (max 10m) | Prevent runaway workflows |
| Node Executions | 10,000 | Yes | Prevent infinite loops |
| HTTP Calls | 100 | Yes | Prevent DoS of external services |
| Loop Iterations | 10,000 | Yes | Prevent infinite loops |
| Recursion Depth | 100 | Yes | Prevent stack overflow |
| Variable Count | 1,000 | Yes | Prevent memory exhaustion |
| String Length | 1MB | Yes | Prevent memory exhaustion |
| Array Size | 10,000 | Yes | Prevent memory exhaustion |
| Object Depth | 32 | Yes | Prevent stack overflow |

### 9.5 Input Validation

All inputs undergo validation:
1. **Structure Validation**: JSON schema compliance
2. **Type Validation**: Correct types for all fields
3. **Size Validation**: Respect size limits
4. **Content Validation**: No malicious patterns

### 9.6 Secure Defaults

- SSRF protection enabled by default
- All resource limits enabled
- Structured logging with security events
- No external dependencies in core (minimize attack surface)

---

## 10. Testing Strategy

### 10.1 Testing Pyramid

```
        ┌──────────────┐
        │     E2E      │  10% - Browser automation
        ├──────────────┤
        │ Integration  │  20% - API tests
        ├──────────────┤
        │     Unit     │  70% - Function/package tests
        └──────────────┘
```

### 10.2 Backend Testing

#### Unit Tests
- **Coverage Target**: >95%
- **Tool**: Go testing framework + testify
- **Scope**: All packages in `backend/pkg/`
- **Run**: `cd backend && go test ./...`

**Example:**
```go
func TestNumberExecutor(t *testing.T) {
    tests := []struct {
        name    string
        value   float64
        want    interface{}
        wantErr bool
    }{
        {"positive", 10.5, 10.5, false},
        {"negative", -5.0, -5.0, false},
        {"zero", 0, 0.0, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

#### Integration Tests
- **Scope**: Engine + executors + middleware
- **Approach**: Test complete workflow execution
- **Location**: `backend/pkg/engine/*_test.go`

#### Benchmark Tests
- **Scope**: Performance-critical paths
- **Run**: `go test -bench=. ./...`

### 10.3 Frontend Testing

#### E2E Tests
- **Tool**: Playwright
- **Scope**: Critical user flows
- **Run**: `npm run test:e2e`

**Test Cases:**
1. Create workflow
2. Add nodes
3. Connect nodes
4. Execute workflow
5. View results

#### Component Tests
- **Tool**: React Testing Library (planned)
- **Scope**: Individual React components

### 10.4 Security Testing

#### SSRF Tests
- Comprehensive test suite in `backend/pkg/security/ssrf_test.go`
- Tests for all blocked IP ranges
- Tests for redirect validation
- Tests for allowlist/blocklist

#### Validation Tests
- Input validation edge cases
- Resource limit enforcement
- Error handling

### 10.5 Performance Testing

#### Load Testing
- **Tool**: k6 or custom Go benchmarks
- **Metrics**:
  - Requests per second
  - Latency percentiles (p50, p95, p99)
  - Memory usage
  - CPU usage

#### Scalability Testing
- Test with increasing workflow sizes (100, 1k, 10k nodes)
- Test concurrent executions (10, 100, 1000)

### 10.6 Regression Testing

- All bug fixes must include regression tests
- Test suite runs on every commit (CI/CD)
- No PR merges without passing tests

### 10.7 Test Data Management

- Use table-driven tests for multiple cases
- Use test fixtures for complex workflows
- Generate test data programmatically
- Separate test data from production data

---

## 11. Deployment & Operations

### 11.1 Deployment Options

#### Option 1: Local Development
```bash
npm run dev                    # Start Next.js dev server
cd backend && go run ./cmd/server  # Start Go server (optional)
```

#### Option 2: Docker
```bash
docker build -t thaiyyal:latest .
docker run -p 8080:8080 thaiyyal:latest
```

#### Option 3: Docker Compose (with monitoring)
```bash
docker-compose up -d
# Includes: Thaiyyal, Prometheus, Grafana
```

#### Option 4: Kubernetes
```bash
kubectl apply -f deployments/kubernetes/
```

### 11.2 Configuration

Configuration via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | HTTP server port |
| `MAX_EXECUTION_TIME` | 30s | Workflow timeout |
| `MAX_NODE_EXECUTIONS` | 10000 | Node limit |
| `MAX_HTTP_CALLS` | 100 | HTTP call limit |
| `LOG_LEVEL` | info | Logging level |
| `ENABLE_TELEMETRY` | false | OpenTelemetry |
| `TELEMETRY_ENDPOINT` | - | OTLP endpoint |

### 11.3 Monitoring

#### Health Checks
- `/health` - Overall health
- `/health/live` - Liveness (process running)
- `/health/ready` - Readiness (can serve traffic)

#### Metrics (Prometheus)
```promql
# Workflow metrics
workflow_executions_total
workflow_execution_duration_seconds
workflow_executions_success_total
workflow_executions_failure_total

# Node metrics
node_executions_total{type="http"}
node_execution_duration_seconds{type="http"}

# HTTP metrics
http_calls_total
http_call_duration_seconds
```

#### Tracing (OpenTelemetry)
- Workflow-level spans
- Node-level spans
- HTTP call spans
- Custom spans in middleware

#### Logging
- Structured JSON logging
- Log levels: debug, info, warn, error
- Correlation IDs for request tracing

### 11.4 Scaling

#### Horizontal Scaling
1. Deploy multiple backend instances
2. Use load balancer (Nginx, ALB, K8s Service)
3. Share state via Redis/PostgreSQL (planned)

#### Vertical Scaling
- Increase memory for larger workflows
- Increase CPU for concurrent executions

#### Auto-Scaling (Kubernetes)
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: thaiyyal-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: thaiyyal-server
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### 11.5 Backup & Recovery

#### Workflow Backup
- Export workflows as JSON files
- Version control recommended (Git)

#### State Backup (Future)
- Database backups (when persistence added)
- Point-in-time recovery

### 11.6 Security Operations

#### Secrets Management
- Use environment variables
- Use Kubernetes secrets
- Use vault solutions (HashiCorp Vault)

#### Network Security
- Use HTTPS in production
- Configure firewall rules
- Use private networks for backend

#### Compliance
- GDPR: No PII stored by default
- SOC2: Audit logs available
- HIPAA: Encryption at rest/transit (TLS)

---

## 12. Development Roadmap

### 12.1 Version 1.0 (Target: Q3 2026)

**Goals**: Production-ready stable release

**Features:**
- [ ] Stable API (no breaking changes)
- [ ] Workflow persistence (PostgreSQL backend)
- [ ] Execution history API
- [ ] Workflow versioning
- [ ] Advanced error handling
- [ ] Comprehensive documentation
- [ ] Production deployment guides

**Non-Functional:**
- [ ] 99.9% uptime in production
- [ ] <100ms p99 latency for validation
- [ ] Support 1000 concurrent executions
- [ ] Complete security audit

### 12.2 Version 1.1 (Target: Q4 2026)

**Goals**: Enhanced integration capabilities

**Features:**
- [ ] GraphQL API support
- [ ] WebSocket real-time execution
- [ ] Workflow templates marketplace
- [ ] More node types (50+ total):
  - Database nodes (PostgreSQL, MySQL)
  - Message queue nodes (RabbitMQ, Kafka)
  - File operations (CSV, Excel)
- [ ] Workflow scheduling (cron)

**Non-Functional:**
- [ ] Support 10,000 concurrent executions
- [ ] Multi-region deployment support

### 12.3 Version 2.0 (Target: Q2 2027)

**Goals**: Enterprise features

**Features:**
- [ ] Authentication & authorization (OAuth, SAML)
- [ ] Role-based access control (RBAC)
- [ ] Team collaboration features
- [ ] Audit logging
- [ ] Compliance certifications (SOC2, ISO 27001)
- [ ] Visual debugger
- [ ] Breakpoint support
- [ ] Step-through execution

**Non-Functional:**
- [ ] 99.99% uptime SLA
- [ ] Multi-tenancy support
- [ ] Enterprise support plans

### 12.4 Future Enhancements

**Advanced Features:**
- [ ] Multi-language expression support (Python, JavaScript)
- [ ] Custom function definitions
- [ ] Workflow composition (sub-workflows)
- [ ] Conditional workflow routing
- [ ] Machine learning integration
- [ ] AI-assisted workflow generation

**Platform Integrations:**
- [ ] AWS services (S3, Lambda, DynamoDB)
- [ ] GCP services (Cloud Storage, BigQuery)
- [ ] Azure services (Blob Storage, Functions)
- [ ] SaaS integrations (Salesforce, HubSpot, Stripe)

**Developer Experience:**
- [ ] VS Code extension
- [ ] CLI tool
- [ ] SDKs (Python, JavaScript, Java)
- [ ] Workflow testing framework
- [ ] CI/CD integration

---

## 13. Success Metrics

### 13.1 Product Metrics

#### Adoption Metrics
- **Active Users**: Monthly active users creating/executing workflows
  - Target (6 months): 1,000 MAU
  - Target (12 months): 10,000 MAU

- **Workflow Creation**: Number of workflows created per month
  - Target (6 months): 5,000 workflows
  - Target (12 months): 50,000 workflows

- **Workflow Executions**: Total executions per month
  - Target (6 months): 100,000 executions
  - Target (12 months): 1,000,000 executions

#### Engagement Metrics
- **Retention Rate**: % of users returning after 7 days
  - Target: >40%

- **Session Duration**: Average time per session
  - Target: >10 minutes

- **Workflows per User**: Average workflows created per user
  - Target: >5

#### Quality Metrics
- **Error Rate**: % of workflow executions that fail
  - Target: <5%

- **Execution Success Rate**: % of workflows completing successfully
  - Target: >95%

### 13.2 Technical Metrics

#### Performance
- **Execution Time**: p95 workflow execution time
  - Target: <5s for typical workflows

- **API Latency**: p95 API response time
  - Target: <200ms

- **Throughput**: Workflows executed per second
  - Target: >100 workflows/sec

#### Reliability
- **Uptime**: % availability
  - Target: >99.9%

- **Mean Time to Recovery**: Average time to recover from incidents
  - Target: <30 minutes

- **Error Budget**: Allowed downtime per month
  - Target: <43 minutes (99.9% uptime)

#### Quality
- **Test Coverage**: % of code covered by tests
  - Target: >95%

- **Bug Density**: Bugs per 1000 lines of code
  - Target: <0.5

- **Security Vulnerabilities**: Critical/High severity issues
  - Target: 0

### 13.3 Business Metrics

#### Growth
- **GitHub Stars**: Community interest
  - Target (6 months): 1,000 stars
  - Target (12 months): 5,000 stars

- **Contributors**: Active open-source contributors
  - Target (6 months): 10 contributors
  - Target (12 months): 50 contributors

- **Docker Pulls**: Container downloads
  - Target (6 months): 10,000 pulls
  - Target (12 months): 100,000 pulls

#### Community
- **Documentation Views**: Monthly docs page views
  - Target: >10,000/month

- **Forum Activity**: Questions/discussions per month
  - Target: >100/month

- **Issue Resolution Time**: Average time to close issues
  - Target: <7 days

### 13.4 User Satisfaction

- **Net Promoter Score (NPS)**: Willingness to recommend
  - Target: >50

- **Customer Satisfaction (CSAT)**: Overall satisfaction rating
  - Target: >4.5/5

- **Feature Request Satisfaction**: % of requested features implemented
  - Target: >30%

---

## Appendices

### A. Glossary

- **Node**: A processing unit in a workflow that performs a specific operation
- **Edge**: A connection between nodes defining data flow
- **Workflow**: A directed acyclic graph (DAG) of nodes and edges
- **Executor**: Code that implements the logic for a node type
- **Middleware**: Cross-cutting logic applied to all node executions
- **Observer**: Component that monitors workflow execution events
- **SSRF**: Server-Side Request Forgery attack
- **DAG**: Directed Acyclic Graph
- **Topological Sort**: Algorithm to determine execution order

### B. References

- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - Detailed architecture
- [REQUIREMENTS_FUNCTIONAL.md](docs/REQUIREMENTS_FUNCTIONAL.md) - Functional requirements
- [REQUIREMENTS_NON_FUNCTIONAL.md](docs/REQUIREMENTS_NON_FUNCTIONAL.md) - Non-functional requirements
- [SECURITY_SUMMARY.md](SECURITY_SUMMARY.md) - Security overview
- [DEV_GUIDE.md](DEV_GUIDE.md) - Developer guide
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines

### C. Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-05-28 | Thaiyyal Team | Initial comprehensive specification |

---

**Document Owner**: Thaiyyal Core Team
**Review Cycle**: Quarterly
**Next Review**: 2026-08-28

*This specification serves as the authoritative reference for what Thaiyyal is, why it exists, and how it should evolve.*
