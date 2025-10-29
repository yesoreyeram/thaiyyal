# Thaiyyal Architecture

## Overview

Thaiyyal is a visual workflow builder MVP with a Next.js frontend and Go backend.

```
┌─────────────────────────────────────────┐
│         Frontend (Next.js/React)        │
│  ┌───────────┐         ┌──────────────┐ │
│  │  Canvas   │◄────────┤ Node Palette │ │
│  │ (ReactFlow)│         └──────────────┘ │
│  └─────┬─────┘                           │
│        │ generates                       │
│        ▼                                 │
│  ┌──────────────────┐                   │
│  │  JSON Payload    │                   │
│  └────────┬─────────┘                   │
└───────────┼──────────────────────────────┘
            │
            │ HTTP (future)
            ▼
┌───────────────────────────────────────────┐
│      Backend (Go Workflow Engine)         │
│  ┌──────────────────────────────────────┐ │
│  │  Engine                              │ │
│  │  ┌──────────────┐   ┌─────────────┐ │ │
│  │  │ Parse JSON   │──►│ Infer Types │ │ │
│  │  └──────────────┘   └──────┬──────┘ │ │
│  │                            │        │ │
│  │  ┌─────────────────────────▼─────┐ │ │
│  │  │  Topological Sort (DAG)       │ │ │
│  │  └──────────────┬─────────────────┘ │ │
│  │                 │                   │ │
│  │  ┌──────────────▼─────────────────┐ │ │
│  │  │  Execute Nodes in Order        │ │ │
│  │  │  ┌──────────────────────────┐  │ │ │
│  │  │  │ Node Executors (23 types)│  │ │ │
│  │  │  └──────────────────────────┘  │ │ │
│  │  └──────────────┬─────────────────┘ │ │
│  │                 │                   │ │
│  │  ┌──────────────▼─────────────────┐ │ │
│  │  │  Return Results                │ │ │
│  │  └──────────────────────────────────┘ │ │
│  └──────────────────────────────────────┘ │
└───────────────────────────────────────────┘
```

## Technology Stack

### Frontend
- **Framework**: Next.js 16.0.1 (App Router)
- **UI Library**: React 19.2.0
- **Workflow Canvas**: ReactFlow 11.8.0
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS 4
- **Build Tool**: Next.js built-in

### Backend
- **Language**: Go 1.24.7
- **Dependencies**: None (standard library only)
- **Package**: `github.com/yesoreyeram/thaiyyal/backend/workflow`

## Directory Structure

```
thaiyyal/
├── src/                          # Frontend source
│   ├── app/                      # Next.js pages
│   │   ├── page.tsx              # Main application
│   │   ├── tests/page.tsx        # Test scenarios
│   │   └── pagination-tests/     # Pagination tests
│   └── components/
│       └── nodes/                # React Flow node components
├── backend/                      # Go workflow engine
│   ├── workflow.go               # Main engine (1,173 LOC)
│   ├── workflow_advanced_nodes.go
│   ├── workflow_errorhandling_nodes.go
│   ├── workflow_*_test.go        # Comprehensive tests
│   └── examples/                 # Example usage
├── docs/                         # Documentation
└── screenshots/                  # Visual documentation
```

## Core Concepts

### Frontend

#### Workflow Canvas
- Uses ReactFlow for visual node editor
- Drag-and-drop interface for node composition
- Real-time JSON payload generation
- Node palette organized by category

#### Node Components
- Each node type has dedicated React component
- Nodes are configurable via properties panel
- Connections define data flow between nodes

### Backend

#### Workflow Execution Model

1. **Parse**: JSON payload → Go structs
2. **Infer**: Determine node types from data (if not explicit)
3. **Sort**: Topological sort for execution order (DAG)
4. **Execute**: Process nodes sequentially
5. **Return**: Collect results and final output

#### Node Types (23 types)

**Basic I/O (3)**
- Number: Numeric input
- TextInput: String input  
- Visualization: Format output

**Operations (3)**
- Operation: Arithmetic (add, subtract, multiply, divide)
- TextOperation: Text transforms (uppercase, lowercase, etc.)
- HTTP: HTTP GET requests

**Control Flow (3)**
- Condition: Conditional branching
- ForEach: Array iteration
- WhileLoop: Conditional looping

**State & Memory (5)**
- Variable: Store/retrieve values
- Extract: Extract object fields
- Transform: Data structure transformations
- Accumulator: Accumulate values over time
- Counter: Simple counter

**Advanced Control Flow (6)**
- Switch: Multi-way branching
- Parallel: Concurrent execution
- Join: Combine multiple inputs (all/any/first)
- Split: Fan-out to multiple paths
- Delay: Pause execution
- Cache: LRU cache with TTL

**Error Handling & Resilience (3)**
- Retry: Retry with backoff strategies
- TryCatch: Error handling with fallback
- Timeout: Enforce time limits

### Data Flow

```
Input Node → Operation Node → Output Node
     │              │              │
     └──── value ───┴──── result ──┘
```

Each edge connects output of source node to input of target node. Engine resolves dependencies via topological sort.

## Key Design Decisions

### 1. No External Dependencies (Backend)
- **Rationale**: Simplicity, security, easy deployment
- **Trade-off**: Implement some features manually
- **Impact**: Minimal attack surface, portable binary

### 2. Type Inference
- **Rationale**: Reduce frontend complexity
- **Implementation**: Infer from NodeData fields
- **Limitation**: Some nodes require explicit type

### 3. DAG Execution
- **Rationale**: Predictable, deterministic execution
- **Algorithm**: Kahn's algorithm for topological sort
- **Limitation**: No cycles (by design)

### 4. In-Memory State
- **Rationale**: Fast for MVP, simple implementation
- **Limitation**: No persistence, single-workflow scope
- **Future**: Add state persistence layer

### 5. Monolithic Files
- **Current**: Large workflow.go (1,173 LOC)
- **Rationale**: Simplicity for MVP
- **Future**: Split into focused modules (see ARCHITECTURE_REVIEW.md)

## State Management

The Engine maintains workflow state:

```go
type Engine struct {
    nodes       []Node
    edges       []Edge
    nodeResults map[string]interface{}
    // State
    variables   map[string]interface{}  // Variable node state
    accumulator interface{}             // Accumulator node state
    counter     float64                 // Counter node state
    cache       map[string]*CacheEntry  // Cache node state
}
```

State is scoped to single workflow execution (not persistent).

## Error Handling

Errors halt execution immediately:
- Parse errors: Invalid JSON
- Validation errors: Missing fields, invalid types
- Runtime errors: Division by zero, HTTP failures
- Graph errors: Cycles detected, missing dependencies

All errors include descriptive messages for debugging.

## Testing Strategy

- **Unit Tests**: Per-node executor testing
- **Integration Tests**: Full workflow execution
- **Table-Driven Tests**: Systematic scenario coverage
- **Coverage**: 142+ test cases, ~95% code coverage

## Performance Characteristics

- **Time Complexity**: O(V + E) for topological sort
- **Space Complexity**: O(V) for node results
- **Execution**: Sequential, single-threaded
- **Scalability**: Suitable for small-to-medium workflows

## Security Considerations

Current security measures:
- Input validation in node executors
- Type checking for operations
- Cycle detection prevents infinite loops

Areas for improvement:
- Add timeouts for all operations
- Implement URL whitelist for HTTP nodes
- Add rate limiting
- Input size limits
- Execution quotas

See ARCHITECTURE_REVIEW.md for detailed recommendations.

## Extension Points

To add a new node type:

1. Add constant to `NodeType` enum
2. Add fields to `NodeData` struct
3. Add type inference logic in `inferNodeTypes()`
4. Implement executor function `executeXXXNode()`
5. Add case to `executeNode()` switch
6. Add frontend component in `src/components/nodes/`
7. Write tests in `workflow_test.go`

## Future Architecture Evolution

See ARCHITECTURE_REVIEW.md for:
- Recommended package structure
- Design patterns
- Scalability improvements
- Performance optimizations
- Security enhancements

## Quick Start

### Frontend
```bash
npm install
npm run dev
# Open http://localhost:3000
```

### Backend
```bash
cd backend
go test -v          # Run tests
cd examples
go run main.go      # Run example
```

## Resources

- **Main README**: `/README.md` - User documentation
- **Backend README**: `/backend/README.md` - Backend details  
- **Architecture Review**: `/ARCHITECTURE_REVIEW.md` - Detailed analysis
- **Node Types**: `/docs/NODES.md` - Complete node reference
- **Frontend Tests**: `/FRONTEND_TESTS.md` - Test scenarios
- **Integration**: `/backend/INTEGRATION.md` - Frontend-backend integration
