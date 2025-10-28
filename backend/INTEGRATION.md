# Frontend to Backend Integration Guide

This document explains how the frontend-generated JSON payload integrates with the backend workflow engine.

## Overview

The Thaiyyal frontend (Next.js/React) generates a JSON workflow payload that can be directly consumed by the backend Go workflow engine for execution.

## Frontend Payload Generation

The frontend canvas generates payloads in this format (from `src/app/page.tsx`):

```javascript
const payload = {
  nodes: nodes.map((n) => ({ id: n.id, data: n.data })),
  edges: edges.map((e) => ({
    id: e.id,
    source: e.source,
    target: e.target,
  })),
};
```

### Example Frontend Payload

```json
{
  "nodes": [
    {"id": "1", "data": {"value": 10, "label": "Node 1"}},
    {"id": "2", "data": {"value": 5, "label": "Node 2"}},
    {"id": "3", "data": {"op": "add", "label": "Node 3 (op)"}},
    {"id": "4", "data": {"mode": "text", "label": "Node 4 (viz)"}}
  ],
  "edges": [
    {"id": "e1-3", "source": "1", "target": "3"},
    {"id": "e2-3", "source": "2", "target": "3"},
    {"id": "e3-4", "source": "3", "target": "4"}
  ]
}
```

## Backend Processing

The backend workflow engine can consume this payload directly:

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    
    "github.com/yesoreyeram/thaiyyal/backend/workflow"
)

func workflowHandler(w http.ResponseWriter, r *http.Request) {
    // Read the payload from the request
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request", http.StatusBadRequest)
        return
    }
    
    // Create the workflow engine
    engine, err := workflow.NewEngine(body)
    if err != nil {
        http.Error(w, fmt.Sprintf("Invalid workflow: %v", err), http.StatusBadRequest)
        return
    }
    
    // Execute the workflow
    result, err := engine.Execute()
    if err != nil {
        http.Error(w, fmt.Sprintf("Execution failed: %v", err), http.StatusInternalServerError)
        return
    }
    
    // Return the result as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func main() {
    http.HandleFunc("/api/execute", workflowHandler)
    fmt.Println("Server listening on :8080")
    http.ListenAndServe(":8080", nil)
}
```

## Node Type Mapping

| Frontend Node Type | Backend Processing | Example Data |
|-------------------|-------------------|--------------|
| `numberNode` | Number input | `{"value": 10}` |
| `opNode` | Arithmetic operation | `{"op": "add"}` |
| `vizNode` | Visualization formatter | `{"mode": "text"}` |

## Workflow Execution Flow

1. **Frontend**: User creates workflow visually
2. **Frontend**: Generates JSON payload on "Show payload" click
3. **Backend** (future): Receive payload via API endpoint
4. **Backend**: Parse JSON using `workflow.NewEngine()`
5. **Backend**: Execute workflow using topological sort
6. **Backend**: Return execution results
7. **Frontend** (future): Display results in visualization nodes

## Example Request/Response

### Request (Frontend → Backend)
```http
POST /api/execute
Content-Type: application/json

{
  "nodes": [
    {"id": "1", "data": {"value": 10}},
    {"id": "2", "data": {"value": 5}},
    {"id": "3", "data": {"op": "multiply"}}
  ],
  "edges": [
    {"id": "e1-3", "source": "1", "target": "3"},
    {"id": "e2-3", "source": "2", "target": "3"}
  ]
}
```

### Response (Backend → Frontend)
```json
{
  "node_results": {
    "1": 10,
    "2": 5,
    "3": 50
  },
  "final_output": 50,
  "errors": []
}
```

## Testing the Integration

### 1. Test with Frontend Payload

Copy the JSON payload from the frontend "Show payload" view and save it as `test-payload.json`:

```bash
cat > test-payload.json << 'EOF'
{
  "nodes": [
    {"id": "1", "data": {"value": 10, "label": "Node 1"}},
    {"id": "2", "data": {"value": 5, "label": "Node 2"}},
    {"id": "3", "data": {"op": "add", "label": "Node 3 (op)"}},
    {"id": "4", "data": {"mode": "text", "label": "Node 4 (viz)"}}
  ],
  "edges": [
    {"id": "e1-3", "source": "1", "target": "3"},
    {"id": "e2-3", "source": "2", "target": "3"},
    {"id": "e3-4", "source": "3", "target": "4"}
  ]
}
EOF
```

### 2. Test with Go Code

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/yesoreyeram/thaiyyal/backend/workflow"
)

func main() {
    // Read the frontend payload
    payload, err := os.ReadFile("test-payload.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Execute the workflow
    engine, err := workflow.NewEngine(payload)
    if err != nil {
        log.Fatal(err)
    }
    
    result, err := engine.Execute()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %+v\n", result)
}
```

## Future API Implementation

The next step would be to create a Next.js API route:

```typescript
// src/app/api/execute/route.ts
import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  const payload = await request.json();
  
  // Call Go backend (could be via HTTP or as a compiled WebAssembly module)
  const response = await fetch('http://localhost:8080/api/execute', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  
  const result = await response.json();
  return NextResponse.json(result);
}
```

## Notes

- The current implementation is a **library** - no HTTP wiring yet
- The backend can parse and execute any frontend-generated payload
- All frontend node types are supported: number, operation, visualization
- The payload format is identical between frontend generation and backend consumption
- Future work: Add API endpoints for remote execution
