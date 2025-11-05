# HTTP Client Management API Examples

This document provides examples of how to use all the server APIs including workflow management, HTTP client management, and frontend serving.

## Overview

The Thaiyyal server provides APIs for:
1. **Workflow Management**
   - Executing workflows (`/api/v1/workflow/execute`)
   - Validating workflows (`/api/v1/workflow/validate`)
   - Saving workflows (`/api/v1/workflow/save`)
   - Listing workflows (`/api/v1/workflow/list`)
   - Loading workflows (`/api/v1/workflow/load/{id}`)
   - Deleting workflows (`/api/v1/workflow/delete/{id}`)
   - Executing by ID (`/api/v1/workflow/execute/{id}`)
2. **HTTP Client Management**
   - Registering HTTP clients (`/api/v1/httpclient/register`)
   - Listing registered HTTP clients (`/api/v1/httpclient/list`)
3. **Frontend Serving**
   - Static files served from root (`/`)

## Starting the Server

```bash
# Start with default settings
./server

# Start with custom settings
./server -addr :9090 -max-execution-time 30s -max-node-executions 1000
```

## Workflow Management

### Save a Workflow

Save a workflow with a name, description, and workflow data.

**Endpoint:** `POST /api/v1/workflow/save`

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/workflow/save \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Addition Workflow",
    "description": "Simple addition of two numbers",
    "data": {
      "nodes": [
        {"id": "1", "data": {"value": 10}},
        {"id": "2", "data": {"value": 5}},
        {"id": "3", "data": {"op": "add"}}
      ],
      "edges": [
        {"source": "1", "target": "3"},
        {"source": "2", "target": "3"}
      ]
    }
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3e4e4585-2b18-4db1-9968-ad2d8649c64c",
  "message": "Workflow saved successfully"
}
```

### List All Workflows

List all saved workflows with their metadata.

**Endpoint:** `GET /api/v1/workflow/list`

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/workflow/list
```

**Response:**
```json
{
  "success": true,
  "workflows": [
    {
      "id": "3e4e4585-2b18-4db1-9968-ad2d8649c64c",
      "name": "My Addition Workflow",
      "description": "Simple addition of two numbers",
      "created_at": "2025-11-05T02:08:10.964574825Z",
      "updated_at": "2025-11-05T02:08:10.964574825Z"
    }
  ],
  "count": 1
}
```

### Load a Workflow by ID

Load a complete workflow including its data.

**Endpoint:** `GET /api/v1/workflow/load/{id}`

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/workflow/load/3e4e4585-2b18-4db1-9968-ad2d8649c64c
```

**Response:**
```json
{
  "success": true,
  "workflow": {
    "id": "3e4e4585-2b18-4db1-9968-ad2d8649c64c",
    "name": "My Addition Workflow",
    "description": "Simple addition of two numbers",
    "data": {
      "nodes": [
        {"id": "1", "data": {"value": 10}},
        {"id": "2", "data": {"value": 5}},
        {"id": "3", "data": {"op": "add"}}
      ],
      "edges": [
        {"source": "1", "target": "3"},
        {"source": "2", "target": "3"}
      ]
    },
    "created_at": "2025-11-05T02:08:10.964574825Z",
    "updated_at": "2025-11-05T02:08:10.964574825Z"
  }
}
```

### Execute a Workflow by ID

Execute a previously saved workflow by its ID.

**Endpoint:** `POST /api/v1/workflow/execute/{id}`

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/workflow/execute/3e4e4585-2b18-4db1-9968-ad2d8649c64c
```

**Response:**
```json
{
  "success": true,
  "workflow_id": "3e4e4585-2b18-4db1-9968-ad2d8649c64c",
  "workflow_name": "My Addition Workflow",
  "results": {
    "execution_id": "76537766d2651745",
    "node_results": {
      "1": 10,
      "2": 5,
      "3": 15
    },
    "final_output": 15
  }
}
```

### Delete a Workflow

Delete a saved workflow by its ID.

**Endpoint:** `DELETE /api/v1/workflow/delete/{id}`

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/workflow/delete/3e4e4585-2b18-4db1-9968-ad2d8649c64c
```

**Response:**
```json
{
  "success": true,
  "message": "Workflow deleted successfully"
}
```

## Workflow Execution

### Execute a Workflow

### Register an HTTP Client

Register a new HTTP client with authentication and security settings.

**Endpoint:** `POST /api/v1/httpclient/register`

**Example: Basic Bearer Token Authentication**
```bash
curl -X POST http://localhost:8080/api/v1/httpclient/register \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "uid": "my-api-client",
      "description": "Client for accessing external API",
      "auth": {
        "type": "bearer",
        "token": {
          "token": "my-secret-token"
        }
      },
      "network": {
        "timeout": 30000000000
      },
      "security": {
        "max_redirects": 5,
        "follow_redirects": true
      }
    }
  }'
```

> **Note:** Timeout values are specified in nanoseconds. For example:
> - 30s = 30000000000 nanoseconds
> - 1m = 60000000000 nanoseconds
> - Use the formula: `seconds * 1000000000` to convert

**Response:**
```json
{
  "success": true,
  "message": "HTTP client registered successfully",
  "uid": "my-api-client"
}
```

**Example: Basic Authentication**
```bash
curl -X POST http://localhost:8080/api/v1/httpclient/register \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "uid": "basic-auth-client",
      "description": "Client with basic auth",
      "auth": {
        "type": "basic",
        "basic_auth": {
          "username": "myuser",
          "password": "mypassword"
        }
      }
    }
  }'
```

**Example: API Key Authentication**
```bash
curl -X POST http://localhost:8080/api/v1/httpclient/register \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "uid": "apikey-client",
      "description": "Client with API key",
      "auth": {
        "type": "apikey",
        "api_key": {
          "key": "X-API-Key",
          "value": "my-api-key-value",
          "location": "header"
        }
      }
    }
  }'
```

**Example: With Custom Headers and Base URL**
```bash
curl -X POST http://localhost:8080/api/v1/httpclient/register \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "uid": "github-api-client",
      "description": "Client for GitHub API",
      "base_url": "https://api.github.com",
      "headers": [
        {"key": "Accept", "value": "application/vnd.github.v3+json"},
        {"key": "User-Agent", "value": "Thaiyyal-Workflow-Engine"}
      ],
      "network": {
        "timeout": 60000000000
      }
    }
  }'
```

> **Note:** The timeout above is 60s (60000000000 nanoseconds)

### List Registered HTTP Clients

List all registered HTTP clients.

**Endpoint:** `GET /api/v1/httpclient/list`

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/httpclient/list
```

**Response:**
```json
{
  "success": true,
  "clients": [
    "my-api-client",
    "github-api-client",
    "basic-auth-client"
  ],
  "count": 3
}
```

## Workflow Execution

### Execute a Workflow

Execute a workflow with nodes and edges.

**Endpoint:** `POST /api/v1/workflow/execute`

**Example: Simple Addition**
```bash
curl -X POST http://localhost:8080/api/v1/workflow/execute \
  -H "Content-Type: application/json" \
  -d '{
    "nodes": [
      {"id": "1", "data": {"value": 10}},
      {"id": "2", "data": {"value": 5}},
      {"id": "3", "data": {"op": "add"}}
    ],
    "edges": [
      {"source": "1", "target": "3"},
      {"source": "2", "target": "3"}
    ]
  }'
```

**Response:**
```json
{
  "success": true,
  "execution_time": "414.175µs",
  "results": {
    "execution_id": "27820e2b232465fe",
    "node_results": {
      "1": 10,
      "2": 5,
      "3": 15
    },
    "final_output": 15
  }
}
```

### Validate a Workflow

Validate a workflow without executing it.

**Endpoint:** `POST /api/v1/workflow/validate`

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/workflow/validate \
  -H "Content-Type: application/json" \
  -d '{
    "nodes": [
      {"id": "1", "data": {"value": 10}},
      {"id": "2", "data": {"value": 5}},
      {"id": "3", "data": {"op": "add"}}
    ],
    "edges": [
      {"source": "1", "target": "3"},
      {"source": "2", "target": "3"}
    ]
  }'
```

**Response (valid):**
```json
{
  "valid": true
}
```

**Response (invalid):**
```json
{
  "valid": false,
  "error": "validation error message"
}
```

## Health Check Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### Liveness Probe
```bash
curl http://localhost:8080/health/live
```

### Readiness Probe
```bash
curl http://localhost:8080/health/ready
```

### Metrics (Prometheus)
```bash
curl http://localhost:8080/metrics
```

## Error Handling

### Duplicate Client Registration
```bash
# Try to register the same client twice
curl -X POST http://localhost:8080/api/v1/httpclient/register \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "uid": "duplicate-client",
      "description": "First registration"
    }
  }'

# Second attempt will fail
curl -X POST http://localhost:8080/api/v1/httpclient/register \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "uid": "duplicate-client",
      "description": "Second registration"
    }
  }'
```

**Error Response:**
```json
{
  "success": false,
  "error": "Failed to register HTTP client: client with UID \"duplicate-client\" already exists"
}
```

### Invalid Configuration
```bash
curl -X POST http://localhost:8080/api/v1/httpclient/register \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "uid": "",
      "description": "Invalid client with empty UID"
    }
  }'
```

**Error Response:**
```json
{
  "success": false,
  "error": "Failed to create HTTP client: invalid config: client UID is required"
}
```

## Security Configuration

HTTP clients support SSRF protection with a zero-trust security model.

**Example: SSRF Protection**
```bash
curl -X POST http://localhost:8080/api/v1/httpclient/register \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "uid": "secure-client",
      "description": "Client with SSRF protection",
      "security": {
        "allow_private_ips": false,
        "allow_localhost": false,
        "allow_link_local": false,
        "allow_cloud_metadata": false,
        "allowed_domains": ["api.example.com", "trusted.example.com"],
        "max_redirects": 3,
        "follow_redirects": true
      }
    }
  }'
```

## Frontend Serving

The server serves the frontend application from the root path (`/`). The frontend and API are served from the same origin, eliminating CORS issues.

**Example:**
```bash
# Access the frontend in your browser
open http://localhost:8080/

# Or use curl to test
curl http://localhost:8080/
```

The frontend is served using Go's embedded filesystem, so no separate build step is needed at runtime.

## Complete Example Workflow

```bash
#!/bin/bash

# 1. Start the server
./server -addr :8080 &
SERVER_PID=$!

# Wait for server to start
sleep 2

# 2. Save a workflow
SAVE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/workflow/save \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Production Workflow",
    "description": "Multiply two numbers",
    "data": {
      "nodes": [
        {"id": "1", "data": {"value": 100}},
        {"id": "2", "data": {"value": 50}},
        {"id": "3", "data": {"op": "multiply"}}
      ],
      "edges": [
        {"source": "1", "target": "3"},
        {"source": "2", "target": "3"}
      ]
    }
  }')

WORKFLOW_ID=$(echo $SAVE_RESPONSE | jq -r '.id')
echo "Workflow ID: $WORKFLOW_ID"

# 3. List all workflows
curl -s http://localhost:8080/api/v1/workflow/list

# 4. Execute the workflow by ID
curl -s -X POST http://localhost:8080/api/v1/workflow/execute/$WORKFLOW_ID

# 5. Register an HTTP client
curl -X POST http://localhost:8080/api/v1/httpclient/register \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "uid": "external-api-client",
      "description": "Client for external API",
      "base_url": "https://api.example.com",
      "auth": {
        "type": "bearer",
        "token": {
          "token": "secret-token-here"
        }
      }
    }
  }'

# 6. List all registered clients
curl -X GET http://localhost:8080/api/v1/httpclient/list

# 7. Check health
curl http://localhost:8080/health

# 8. Access frontend
echo "Open http://localhost:8080/ in your browser"

# Stop the server
kill $SERVER_PID
```

## Configuration Reference

### HTTP Client Configuration

```json
{
  "config": {
    "uid": "string (required)",
    "description": "string (optional)",
    "base_url": "string (optional)",
    "auth": {
      "type": "none|basic|bearer|apikey",
      "basic_auth": {
        "username": "string",
        "password": "string"
      },
      "token": {
        "token": "string"
      },
      "api_key": {
        "key": "string",
        "value": "string",
        "location": "header|query"
      }
    },
    "network": {
      "timeout": "duration in nanoseconds (default: 30s)",
      "max_idle_conns": "int (default: 100)",
      "max_idle_conns_per_host": "int (default: 10)",
      "max_conns_per_host": "int (default: 100)",
      "idle_conn_timeout": "duration in nanoseconds (default: 90s)",
      "tls_handshake_timeout": "duration in nanoseconds (default: 10s)",
      "disable_keep_alives": "bool (default: false)"
    },
    "security": {
      "max_redirects": "int (default: 10)",
      "max_response_size": "int64 (default: 10MB)",
      "follow_redirects": "bool (default: true)",
      "allow_private_ips": "bool (default: false)",
      "allow_localhost": "bool (default: false)",
      "allow_link_local": "bool (default: false)",
      "allow_cloud_metadata": "bool (default: false)",
      "allowed_domains": ["array of allowed domains"]
    },
    "headers": [
      {"key": "string", "value": "string"}
    ],
    "query_params": [
      {"key": "string", "value": "string"}
    ]
  }
}
```
