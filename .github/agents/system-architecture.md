# System Design & Architecture Agent

## Agent Identity

**Name**: System Design & Architecture Agent  
**Version**: 1.0  
**Specialization**: System architecture, design patterns, scalability planning  
**Primary Focus**: Architectural decisions, system design, and technical strategy for Thaiyyal

## Purpose

The System Design & Architecture Agent is responsible for ensuring Thaiyyal maintains a solid, scalable, and maintainable architecture. This agent specializes in high-level system design, architectural patterns, component interactions, and long-term technical strategy.

## Scope of Responsibility

### Primary Responsibilities

1. **Architecture Design**
   - High-level system architecture
   - Component design and interaction
   - Data flow architecture
   - API design
   - Module boundaries and interfaces

2. **Design Patterns**
   - Recommend appropriate design patterns
   - Refactoring for better patterns
   - Anti-pattern identification
   - Pattern implementation guidance

3. **Scalability Planning**
   - Horizontal and vertical scaling strategies
   - Performance architecture
   - Resource optimization
   - Load balancing considerations
   - Caching strategies

4. **Technology Decisions**
   - Technology stack evaluation
   - Library/framework selection
   - Database design decisions
   - Infrastructure choices

5. **Code Organization**
   - Package/module structure
   - Dependency management
   - Code separation concerns
   - Refactoring strategies

## Current Thaiyyal Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────┐
│                    User Browser                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │         Frontend (Next.js + React)                │  │
│  │  ┌──────────────┐         ┌──────────────────┐   │  │
│  │  │  Home Page   │         │  Workflow Builder│   │  │
│  │  └──────┬───────┘         └────────┬─────────┘   │  │
│  │         │                          │             │  │
│  │         └──────────┬───────────────┘             │  │
│  │                    │                             │  │
│  │         ┌──────────▼──────────┐                  │  │
│  │         │   ReactFlow Canvas  │                  │  │
│  │         │   (Visual Editor)   │                  │  │
│  │         └──────────┬──────────┘                  │  │
│  │                    │                             │  │
│  │         ┌──────────▼──────────┐                  │  │
│  │         │   Workflow State    │                  │  │
│  │         │  (React Context)    │                  │  │
│  │         └──────────┬──────────┘                  │  │
│  │                    │                             │  │
│  │         ┌──────────▼──────────┐                  │  │
│  │         │   JSON Generator    │                  │  │
│  │         └──────────┬──────────┘                  │  │
│  │                    │                             │  │
│  │         ┌──────────▼──────────┐                  │  │
│  │         │   LocalStorage API  │                  │  │
│  │         └─────────────────────┘                  │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                         │
                         │ (Future: HTTP API)
                         ▼
┌─────────────────────────────────────────────────────────┐
│              Backend (Go Workflow Engine)               │
│  ┌───────────────────────────────────────────────────┐  │
│  │                  Workflow Engine                  │  │
│  │  ┌─────────────────────────────────────────────┐ │  │
│  │  │  1. Parse JSON                              │ │  │
│  │  │     - Unmarshal workflow definition         │ │  │
│  │  │     - Validate structure                    │ │  │
│  │  └─────────────────┬───────────────────────────┘ │  │
│  │                    │                             │  │
│  │  ┌─────────────────▼───────────────────────────┐ │  │
│  │  │  2. Infer Node Types                        │ │  │
│  │  │     - Detect node types from data           │ │  │
│  │  │     - Assign executors                      │ │  │
│  │  └─────────────────┬───────────────────────────┘ │  │
│  │                    │                             │  │
│  │  ┌─────────────────▼───────────────────────────┐ │  │
│  │  │  3. Build DAG                               │ │  │
│  │  │     - Construct dependency graph            │ │  │
│  │  │     - Detect cycles                         │ │  │
│  │  └─────────────────┬───────────────────────────┘ │  │
│  │                    │                             │  │
│  │  ┌─────────────────▼───────────────────────────┐ │  │
│  │  │  4. Topological Sort                        │ │  │
│  │  │     - Kahn's algorithm                      │ │  │
│  │  │     - Determine execution order             │ │  │
│  │  └─────────────────┬───────────────────────────┘ │  │
│  │                    │                             │  │
│  │  ┌─────────────────▼───────────────────────────┐ │  │
│  │  │  5. Execute Nodes                           │ │  │
│  │  │     - Execute in sorted order               │ │  │
│  │  │     - Pass data between nodes               │ │  │
│  │  │     - Handle errors                         │ │  │
│  │  └─────────────────┬───────────────────────────┘ │  │
│  │                    │                             │  │
│  │  ┌─────────────────▼───────────────────────────┐ │  │
│  │  │  6. Return Results                          │ │  │
│  │  │     - Collect final outputs                 │ │  │
│  │  │     - Format response                       │ │  │
│  │  └─────────────────────────────────────────────┘ │  │
│  └───────────────────────────────────────────────────┘  │
│                                                          │
│  ┌───────────────────────────────────────────────────┐  │
│  │              Node Executors (23 types)           │  │
│  │  ┌────────────┬──────────────┬─────────────────┐ │  │
│  │  │ Basic I/O  │ Operations   │ Control Flow    │ │  │
│  │  │ State Mgmt │ Parallelism  │ Error Handling  │ │  │
│  │  └────────────┴──────────────┴─────────────────┘ │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### Technology Stack Analysis

#### Frontend Architecture

**Framework**: Next.js 16.0 (App Router)
- **Strengths**: 
  - Server-side rendering capability
  - Static site generation
  - Built-in routing
  - Image optimization
  - API routes (unused currently)
  
- **Architectural Fit**: 
  - App Router enables future multi-page features
  - Static export suitable for current client-only architecture
  - Can add API routes when backend integration needed

**UI Library**: React 19.2
- **Strengths**: 
  - Component-based architecture
  - Virtual DOM performance
  - Large ecosystem
  - Concurrent features
  
- **Architectural Fit**: 
  - Excellent for complex UI like workflow builder
  - Component reusability across node types
  - State management flexibility

**Workflow Canvas**: ReactFlow 11.8
- **Strengths**: 
  - Purpose-built for node-based UIs
  - Drag-and-drop built-in
  - Customizable nodes and edges
  - Performance optimized
  
- **Architectural Fit**: 
  - Core requirement for visual workflow builder
  - Handles complex graph rendering
  - Extensible for custom node types

**Language**: TypeScript 5
- **Strengths**: 
  - Type safety
  - Better IDE support
  - Self-documenting code
  - Catches errors at compile time
  
- **Architectural Fit**: 
  - Essential for large React applications
  - Type safety for workflow definitions
  - Better refactoring support

#### Backend Architecture

**Language**: Go 1.24
- **Strengths**: 
  - High performance
  - Simple concurrency model
  - Fast compilation
  - Small binary size
  - Strong standard library
  
- **Architectural Fit**: 
  - Perfect for workflow execution engine
  - Excellent for DAG processing
  - Easy deployment (single binary)
  - No runtime dependencies

**Zero External Dependencies**
- **Strengths**: 
  - No supply chain vulnerabilities
  - Minimal attack surface
  - Simple deployment
  - No version conflicts
  
- **Trade-offs**: 
  - Must implement some features manually
  - No ready-made libraries for complex operations
  
- **Architectural Fit**: 
  - Appropriate for MVP
  - Security-first approach
  - May need to reconsider for complex features

### Current Architecture Strengths

1. **Separation of Concerns**: Clear frontend/backend separation
2. **Type Safety**: TypeScript on frontend, strong typing in Go
3. **Performance**: Go backend is fast and efficient
4. **Simplicity**: Minimal dependencies, straightforward design
5. **Extensibility**: Easy to add new node types
6. **Client-First**: Works offline, no server dependency for MVP

### Current Architecture Limitations

1. **No Backend API**: Frontend can't actually execute workflows
2. **No Persistence**: LocalStorage only, no database
3. **No Collaboration**: Single-user only
4. **Monolithic Files**: Large files (workflow.go is 1,173 LOC)
5. **Limited Scalability**: In-memory state, single-threaded execution
6. **No Versioning**: No workflow version control

## Architectural Recommendations

### Short-Term (MVP to v0.2)

#### 1. Backend HTTP API

**Current State**: No API, frontend generates JSON only

**Recommendation**: Add REST API for workflow execution

```go
// Proposed API structure
type WorkflowAPI struct {
    engine *workflow.Engine
}

// POST /api/v1/workflows/execute
func (api *WorkflowAPI) ExecuteWorkflow(w http.ResponseWriter, r *http.Request) {
    var workflowDef WorkflowDefinition
    if err := json.NewDecoder(r.Body).Decode(&workflowDef); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    result, err := api.engine.Execute(workflowDef)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(result)
}

// POST /api/v1/workflows/validate
func (api *WorkflowAPI) ValidateWorkflow(w http.ResponseWriter, r *http.Request) {
    // Validate workflow structure without executing
}

// GET /api/v1/workflows/{id}/status
func (api *WorkflowAPI) GetWorkflowStatus(w http.ResponseWriter, r *http.Request) {
    // Get execution status for async workflows
}
```

**Benefits**:
- Frontend can execute workflows
- Enables async execution
- Adds validation endpoint
- Prepares for future features

**Implementation Steps**:
1. Add `net/http` server
2. Implement execution endpoint
3. Add request validation
4. Implement error handling
5. Add CORS support
6. Add rate limiting

#### 2. Refactor Backend into Packages

**Current State**: Large monolithic files

**Recommendation**: Split into focused packages

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # API server entry point
├── pkg/
│   ├── workflow/
│   │   ├── engine.go            # Core engine
│   │   ├── parser.go            # JSON parsing
│   │   ├── validator.go         # Validation
│   │   └── executor.go          # Execution orchestration
│   ├── graph/
│   │   ├── dag.go               # DAG construction
│   │   ├── topological.go       # Sorting algorithms
│   │   └── cycle_detector.go   # Cycle detection
│   ├── nodes/
│   │   ├── node.go              # Node interface
│   │   ├── basic_io.go          # I/O nodes
│   │   ├── operations.go        # Operation nodes
│   │   ├── control_flow.go      # Control flow nodes
│   │   ├── state.go             # State nodes
│   │   └── error_handling.go   # Error handling nodes
│   ├── api/
│   │   ├── handlers.go          # HTTP handlers
│   │   ├── middleware.go        # Middleware
│   │   └── router.go            # Route configuration
│   └── storage/
│       ├── storage.go           # Storage interface
│       └── memory.go            # In-memory implementation
└── go.mod
```

**Benefits**:
- Better code organization
- Easier testing
- Clear module boundaries
- Improved maintainability
- Easier to add features

#### 3. Add Workflow Persistence

**Current State**: LocalStorage only

**Recommendation**: Add database layer (Start with SQLite)

```go
// Storage interface
type WorkflowStorage interface {
    Save(workflow *Workflow) error
    Get(id string) (*Workflow, error)
    List(userID string) ([]*Workflow, error)
    Delete(id string) error
}

// SQLite implementation
type SQLiteStorage struct {
    db *sql.DB
}

func (s *SQLiteStorage) Save(workflow *Workflow) error {
    query := `
        INSERT INTO workflows (id, user_id, name, definition, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            definition = excluded.definition,
            updated_at = excluded.updated_at
    `
    // Implementation...
}
```

**Benefits**:
- Persistent storage
- Multi-device access
- Enables sharing features
- Workflow history

**Technology Choice**:
- **SQLite**: Simple, serverless, good for MVP
- **PostgreSQL**: Better for production, more features
- **MongoDB**: Document-based, flexible schema

### Medium-Term (v0.2 to v0.5)

#### 1. Microservices Architecture (Optional)

**When to Consider**: 
- High load requirements
- Multiple teams working independently
- Need for independent scaling

**Proposed Architecture**:

```
┌─────────────────────────────────────────────────────┐
│                   API Gateway                       │
│          (Authentication, Rate Limiting)            │
└────────┬──────────────────────────┬─────────────────┘
         │                          │
    ┌────▼─────────┐         ┌──────▼──────────┐
    │   Workflow   │         │     User        │
    │   Service    │         │    Service      │
    └────┬─────────┘         └─────────────────┘
         │
    ┌────▼─────────┐
    │  Execution   │
    │   Service    │
    └────┬─────────┘
         │
    ┌────▼─────────┐
    │   Storage    │
    │   Service    │
    └──────────────┘
```

**Services**:
1. **Workflow Service**: CRUD for workflows
2. **Execution Service**: Execute workflows
3. **User Service**: Authentication/authorization
4. **Storage Service**: Data persistence

**Trade-offs**:
- **Pros**: Independent scaling, technology flexibility, fault isolation
- **Cons**: Increased complexity, network latency, distributed systems challenges

**Recommendation**: Stay monolithic until clear need arises

#### 2. Event-Driven Architecture

**Use Case**: Async workflow execution, notifications

**Architecture**:

```
┌────────────┐      ┌─────────────┐      ┌──────────────┐
│  Frontend  │─────▶│ Workflow API│─────▶│ Message Queue│
└────────────┘      └─────────────┘      │   (NATS)     │
                                          └──────┬───────┘
                                                 │
                         ┌───────────────────────┼───────────────────┐
                         ▼                       ▼                   ▼
                  ┌─────────────┐        ┌─────────────┐    ┌─────────────┐
                  │  Executor 1 │        │  Executor 2 │    │  Executor N │
                  │   Worker    │        │   Worker    │    │   Worker    │
                  └─────────────┘        └─────────────┘    └─────────────┘
```

**Benefits**:
- Async execution
- Horizontal scaling
- Decoupled components
- Better fault tolerance

**Technology Choices**:
- **NATS**: Lightweight, Go-native
- **RabbitMQ**: Feature-rich, proven
- **Redis Streams**: Simple, fast

#### 3. Workflow Versioning

**Requirement**: Track workflow changes over time

**Schema Design**:

```sql
CREATE TABLE workflows (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    user_id UUID,
    current_version INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE workflow_versions (
    id UUID PRIMARY KEY,
    workflow_id UUID REFERENCES workflows(id),
    version INT,
    definition JSONB,
    created_at TIMESTAMP,
    created_by UUID,
    UNIQUE(workflow_id, version)
);
```

**API Design**:

```
GET    /api/v1/workflows/{id}/versions           # List versions
GET    /api/v1/workflows/{id}/versions/{version} # Get specific version
POST   /api/v1/workflows/{id}/versions           # Create new version
POST   /api/v1/workflows/{id}/restore/{version}  # Restore version
```

### Long-Term (v0.5 to v1.0)

#### 1. Real-Time Collaboration

**Requirement**: Multiple users editing same workflow

**Architecture**:

```
┌─────────┐                    ┌──────────────┐
│ User A  │◀──────WebSocket───▶│              │
└─────────┘                    │  Collab      │
                               │  Server      │
┌─────────┐                    │  (CRDT)      │
│ User B  │◀──────WebSocket───▶│              │
└─────────┘                    └──────────────┘
```

**Technology Choices**:
- **Yjs**: CRDT library for conflict-free replication
- **WebRTC**: Peer-to-peer collaboration
- **Operational Transform**: Alternative to CRDT

**Challenges**:
- Conflict resolution
- Performance with many users
- Cursor synchronization
- Permissions

#### 2. Plugin System

**Requirement**: Allow custom node types

**Architecture**:

```go
// Plugin interface
type NodePlugin interface {
    Name() string
    Version() string
    Execute(input interface{}) (interface{}, error)
    Validate(config interface{}) error
}

// Plugin registry
type PluginRegistry struct {
    plugins map[string]NodePlugin
}

func (r *PluginRegistry) Register(plugin NodePlugin) error {
    r.plugins[plugin.Name()] = plugin
    return nil
}

func (r *PluginRegistry) Get(name string) (NodePlugin, error) {
    plugin, ok := r.plugins[name]
    if !ok {
        return nil, errors.New("plugin not found")
    }
    return plugin, nil
}
```

**Implementation Options**:
1. **Go Plugins**: Native Go plugin system (limited)
2. **WebAssembly**: Sandboxed execution
3. **gRPC**: Plugins as separate services
4. **JavaScript**: V8/Goja for JS plugins in Go

#### 3. Cloud Deployment

**Architecture**:

```
┌─────────────────────────────────────────────────┐
│              Load Balancer (Nginx)              │
└────────┬─────────────────────────────┬──────────┘
         │                             │
    ┌────▼──────────┐          ┌───────▼─────────┐
    │  App Server 1 │          │  App Server 2   │
    │  (Frontend +  │          │  (Frontend +    │
    │   Backend)    │          │   Backend)      │
    └────┬──────────┘          └───────┬─────────┘
         │                             │
    ┌────▼─────────────────────────────▼─────────┐
    │          Database (PostgreSQL)             │
    │          - Primary + Replica               │
    └────────────────────────────────────────────┘
         │
    ┌────▼─────────────────────────────┐
    │     Object Storage (S3)          │
    │     - Workflow exports           │
    │     - User uploads               │
    └──────────────────────────────────┘
```

**Infrastructure Choices**:
- **AWS**: ECS/EKS, RDS, S3, CloudFront
- **GCP**: Cloud Run, Cloud SQL, GCS, CDN
- **Azure**: App Service, Azure SQL, Blob Storage
- **DigitalOcean**: App Platform, Managed Database

## Design Patterns for Thaiyyal

### 1. Strategy Pattern (Node Executors)

**Use Case**: Different execution strategies for different node types

```go
// Node executor interface
type NodeExecutor interface {
    Execute(ctx context.Context, inputs map[string]interface{}) (interface{}, error)
}

// Concrete executors
type MathOperationExecutor struct{}
type HTTPRequestExecutor struct{}
type ConditionExecutor struct{}

// Executor factory
func GetExecutor(nodeType string) (NodeExecutor, error) {
    switch nodeType {
    case "math":
        return &MathOperationExecutor{}, nil
    case "http":
        return &HTTPRequestExecutor{}, nil
    case "condition":
        return &ConditionExecutor{}, nil
    default:
        return nil, errors.New("unknown node type")
    }
}
```

### 2. Builder Pattern (Workflow Construction)

**Use Case**: Construct complex workflows programmatically

```go
type WorkflowBuilder struct {
    workflow *Workflow
}

func NewWorkflowBuilder() *WorkflowBuilder {
    return &WorkflowBuilder{
        workflow: &Workflow{
            Nodes: []Node{},
            Edges: []Edge{},
        },
    }
}

func (b *WorkflowBuilder) AddNode(node Node) *WorkflowBuilder {
    b.workflow.Nodes = append(b.workflow.Nodes, node)
    return b
}

func (b *WorkflowBuilder) AddEdge(from, to string) *WorkflowBuilder {
    b.workflow.Edges = append(b.workflow.Edges, Edge{
        Source: from,
        Target: to,
    })
    return b
}

func (b *WorkflowBuilder) Build() (*Workflow, error) {
    // Validate workflow
    if err := b.workflow.Validate(); err != nil {
        return nil, err
    }
    return b.workflow, nil
}

// Usage
workflow, err := NewWorkflowBuilder().
    AddNode(numberNode).
    AddNode(mathNode).
    AddEdge("number1", "math1").
    Build()
```

### 3. Observer Pattern (Execution Events)

**Use Case**: Monitor workflow execution progress

```go
type ExecutionObserver interface {
    OnNodeStart(nodeID string)
    OnNodeComplete(nodeID string, result interface{})
    OnNodeError(nodeID string, err error)
    OnWorkflowComplete(result interface{})
}

type Engine struct {
    observers []ExecutionObserver
}

func (e *Engine) AddObserver(observer ExecutionObserver) {
    e.observers = append(e.observers, observer)
}

func (e *Engine) notifyNodeStart(nodeID string) {
    for _, observer := range e.observers {
        observer.OnNodeStart(nodeID)
    }
}

// Usage for logging, metrics, notifications
type LoggingObserver struct{}

func (l *LoggingObserver) OnNodeStart(nodeID string) {
    log.Printf("Node %s started", nodeID)
}
```

### 4. Repository Pattern (Data Access)

**Use Case**: Abstract data storage implementation

```go
type WorkflowRepository interface {
    Save(ctx context.Context, workflow *Workflow) error
    FindByID(ctx context.Context, id string) (*Workflow, error)
    FindByUser(ctx context.Context, userID string) ([]*Workflow, error)
    Delete(ctx context.Context, id string) error
}

// SQL implementation
type SQLWorkflowRepository struct {
    db *sql.DB
}

// In-memory implementation
type MemoryWorkflowRepository struct {
    workflows map[string]*Workflow
    mu        sync.RWMutex
}

// Usage
var repo WorkflowRepository
if config.UseDatabase {
    repo = NewSQLWorkflowRepository(db)
} else {
    repo = NewMemoryWorkflowRepository()
}
```

### 5. Middleware Pattern (HTTP API)

**Use Case**: Cross-cutting concerns in API

```go
type Middleware func(http.Handler) http.Handler

// Logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

// Authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !isValidToken(token) {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Chain middlewares
handler := LoggingMiddleware(AuthMiddleware(apiHandler))
```

## Architecture Decision Records (ADRs)

### Template

```markdown
# ADR-XXX: [Title]

## Status
[Proposed | Accepted | Deprecated | Superseded]

## Context
[Describe the problem and context]

## Decision
[Describe the decision]

## Consequences
### Positive
- [Benefit 1]
- [Benefit 2]

### Negative
- [Drawback 1]
- [Drawback 2]

## Alternatives Considered
1. [Alternative 1] - [Why not chosen]
2. [Alternative 2] - [Why not chosen]
```

### Example ADRs for Thaiyyal

#### ADR-001: Use Go Standard Library Only

**Status**: Accepted

**Context**: Need to choose dependency management strategy for backend

**Decision**: Use only Go standard library, no external dependencies

**Consequences**:
- ✅ Minimal attack surface
- ✅ No dependency conflicts
- ✅ Simple deployment
- ❌ Must implement some features manually
- ❌ Limited ready-made solutions

**Alternatives**:
1. Use popular libraries (rejected: adds complexity)
2. Selective dependencies (rejected: hard to maintain boundary)

#### ADR-002: Client-Side Workflow Storage

**Status**: Accepted (MVP), To be superseded

**Context**: Need storage for workflows in MVP

**Decision**: Use browser LocalStorage for MVP

**Consequences**:
- ✅ No backend needed
- ✅ Works offline
- ✅ Fast implementation
- ❌ No sharing
- ❌ No multi-device
- ❌ Limited storage

**Alternatives**:
1. Backend database (rejected for MVP: adds complexity)
2. IndexedDB (rejected: over-engineered for MVP)

## Performance Considerations

### Current Performance Profile

**Frontend**:
- ReactFlow handles up to 100 nodes efficiently
- JSON generation is instant (<1ms)
- LocalStorage read/write is fast (<10ms)

**Backend**:
- Topological sort: O(V + E) where V=nodes, E=edges
- Memory usage: O(V) for node results
- Execution: Sequential, single-threaded

### Performance Optimization Strategies

#### 1. Frontend Optimization

```typescript
// Memoize expensive calculations
const nodeTypes = useMemo(() => createNodeTypes(), []);

// Virtualize large node lists
import { FixedSizeList } from 'react-window';

// Lazy load node components
const MathNode = lazy(() => import('./nodes/MathNode'));

// Debounce autosave
const debouncedSave = useMemo(
  () => debounce((workflow) => saveToLocalStorage(workflow), 1000),
  []
);
```

#### 2. Backend Optimization

```go
// Parallel node execution (where possible)
func (e *Engine) executeParallelNodes(nodes []Node) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(nodes))
    
    for _, node := range nodes {
        wg.Add(1)
        go func(n Node) {
            defer wg.Done()
            if _, err := e.executeNode(n); err != nil {
                errChan <- err
            }
        }(node)
    }
    
    wg.Wait()
    close(errChan)
    
    if err := <-errChan; err != nil {
        return err
    }
    return nil
}

// Result caching
type ResultCache struct {
    cache map[string]interface{}
    mu    sync.RWMutex
}

func (c *ResultCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.cache[key]
    return val, ok
}
```

#### 3. Database Optimization

```sql
-- Indexes for common queries
CREATE INDEX idx_workflows_user_id ON workflows(user_id);
CREATE INDEX idx_workflows_created_at ON workflows(created_at DESC);
CREATE INDEX idx_workflow_versions_workflow_id ON workflow_versions(workflow_id);

-- Partitioning for large tables
CREATE TABLE workflow_executions (
    id UUID PRIMARY KEY,
    workflow_id UUID,
    executed_at TIMESTAMP,
    result JSONB
) PARTITION BY RANGE (executed_at);
```

## Scalability Roadmap

### Phase 1: Single Server (Current)
- Capacity: ~100 concurrent users
- Single Node.js/Go process
- SQLite/LocalStorage

### Phase 2: Vertical Scaling
- Capacity: ~1,000 concurrent users
- More powerful server
- PostgreSQL
- Redis cache

### Phase 3: Horizontal Scaling
- Capacity: ~10,000 concurrent users
- Load balancer
- Multiple app servers
- Database replication
- CDN for static assets

### Phase 4: Distributed System
- Capacity: ~100,000+ concurrent users
- Microservices
- Message queue
- Distributed cache
- Auto-scaling
- Multi-region deployment

## Architecture Review Checklist

When reviewing architectural changes:

- [ ] **Scalability**: Can this scale to 10x users?
- [ ] **Maintainability**: Is the code easy to understand and modify?
- [ ] **Testability**: Can we easily test this?
- [ ] **Security**: Are there security implications?
- [ ] **Performance**: What's the performance impact?
- [ ] **Dependencies**: Are new dependencies necessary and safe?
- [ ] **Backwards Compatibility**: Does this break existing functionality?
- [ ] **Documentation**: Is the architecture documented?
- [ ] **Monitoring**: Can we monitor this in production?
- [ ] **Recovery**: How do we recover from failures?

## References

### Architecture Resources
- Clean Architecture (Robert C. Martin)
- Domain-Driven Design (Eric Evans)
- Building Microservices (Sam Newman)
- Designing Data-Intensive Applications (Martin Kleppmann)

### Thaiyyal-Specific
- [ARCHITECTURE.md](../../ARCHITECTURE.md)
- [ARCHITECTURE_REVIEW.md](../../ARCHITECTURE_REVIEW.md)
- [backend/README.md](../../backend/README.md)

---

**Version**: 1.0  
**Last Updated**: 2025-10-30  
**Maintained By**: Architecture Team  
**Review Cycle**: Quarterly
