export interface WorkflowExample {
  id: string;
  title: string;
  description: string;
  tags: string[];
  nodes: unknown[];
  edges: unknown[];
}

export const workflowExamples: WorkflowExample[] =
[
  {
    "id": "example-1",
    "title": "Simple API Call",
    "description": "Make a basic HTTP GET request to fetch data from an API endpoint and display the results.",
    "tags": [
      "api",
      "http",
      "get",
      "basic"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-2",
    "title": "API Call with Retry Logic",
    "description": "Make an HTTP request with automatic retry mechanism in case of failures. Implements exponential backoff strategy.",
    "tags": [
      "api",
      "http",
      "retry",
      "error-handling"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "retryNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "max_attempts": 3,
          "backoff_strategy": "exponential",
          "initial_delay": "1s",
          "label": "Retry Handler"
        }
      },
      {
        "id": "2",
        "type": "httpNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 500,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-3",
    "title": "API Call with Timeout",
    "description": "Execute an HTTP request with a timeout protection to prevent hanging requests. Falls back gracefully on timeout.",
    "tags": [
      "api",
      "http",
      "timeout",
      "error-handling"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "timeoutNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "timeout": "30s",
          "timeout_action": "error",
          "label": "Timeout Protection"
        }
      },
      {
        "id": "2",
        "type": "httpNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 500,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-4",
    "title": "POST Request with JSON Payload",
    "description": "Send data to an API endpoint using POST method with JSON body. Handles response and error codes appropriately.",
    "tags": [
      "api",
      "http",
      "post",
      "json"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "{\"name\": \"John\"}",
          "label": "JSON Payload"
        }
      },
      {
        "id": "2",
        "type": "parseNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "input_type": "JSON",
          "label": "Parse JSON"
        }
      },
      {
        "id": "3",
        "type": "httpNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/users",
          "method": "POST",
          "label": "POST Request"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Response"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-5",
    "title": "API Authentication Flow",
    "description": "Implement OAuth 2.0 authentication flow with token refresh. Securely manages access tokens and handles expiration.",
    "tags": [
      "api",
      "http",
      "auth",
      "oauth"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://auth.example.com/token",
          "method": "POST",
          "label": "Get Token"
        }
      },
      {
        "id": "2",
        "type": "extractNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "field": "access_token",
          "label": "Extract Token"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 500,
          "y": 100
        },
        "data": {
          "var_name": "auth_token",
          "var_op": "set",
          "label": "Store Token"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 700,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Success"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-6",
    "title": "Parallel API Calls",
    "description": "Execute multiple API requests in parallel for improved performance. Aggregates results from multiple endpoints.",
    "tags": [
      "api",
      "http",
      "parallel",
      "performance"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "parallelNode",
        "position": {
          "x": 100,
          "y": 150
        },
        "data": {
          "max_concurrency": 3,
          "label": "Parallel Execute"
        }
      },
      {
        "id": "2",
        "type": "httpNode",
        "position": {
          "x": 300,
          "y": 50
        },
        "data": {
          "url": "https://api.example.com/endpoint1",
          "method": "GET",
          "label": "API 1"
        }
      },
      {
        "id": "3",
        "type": "httpNode",
        "position": {
          "x": 300,
          "y": 150
        },
        "data": {
          "url": "https://api.example.com/endpoint2",
          "method": "GET",
          "label": "API 2"
        }
      },
      {
        "id": "4",
        "type": "httpNode",
        "position": {
          "x": 300,
          "y": 250
        },
        "data": {
          "url": "https://api.example.com/endpoint3",
          "method": "GET",
          "label": "API 3"
        }
      },
      {
        "id": "5",
        "type": "joinNode",
        "position": {
          "x": 500,
          "y": 150
        },
        "data": {
          "join_strategy": "all",
          "label": "Join Results"
        }
      },
      {
        "id": "6",
        "type": "vizNode",
        "position": {
          "x": 700,
          "y": 150
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e1-3",
        "source": "1",
        "target": "3"
      },
      {
        "id": "e1-4",
        "source": "1",
        "target": "4"
      },
      {
        "id": "e2-5",
        "source": "2",
        "target": "5"
      },
      {
        "id": "e3-5",
        "source": "3",
        "target": "5"
      },
      {
        "id": "e4-5",
        "source": "4",
        "target": "5"
      },
      {
        "id": "e5-6",
        "source": "5",
        "target": "6"
      }
    ]
  },
  {
    "id": "example-7",
    "title": "API Rate Limiting Handler",
    "description": "Implement rate limiting logic to comply with API quotas. Handles 429 responses and implements exponential backoff.",
    "tags": [
      "api",
      "http",
      "rate-limiting",
      "error-handling"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "retryNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "max_attempts": 5,
          "backoff_strategy": "exponential",
          "initial_delay": "2s",
          "label": "Retry Logic"
        }
      },
      {
        "id": "2",
        "type": "delayNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "duration": "1s",
          "label": "Delay"
        }
      },
      {
        "id": "3",
        "type": "httpNode",
        "position": {
          "x": 500,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 700,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-8",
    "title": "REST API CRUD Operations",
    "description": "Complete CRUD workflow demonstrating Create, Read, Update, Delete operations on a REST API resource.",
    "tags": [
      "api",
      "http",
      "crud",
      "rest"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-9",
    "title": "GraphQL Query Execution",
    "description": "Execute GraphQL queries with variables and fragments. Handles complex nested data structures.",
    "tags": [
      "api",
      "http",
      "graphql",
      "query"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-10",
    "title": "API Response Caching",
    "description": "Cache API responses to reduce redundant calls and improve performance. Implements TTL-based invalidation.",
    "tags": [
      "api",
      "http",
      "cache",
      "performance"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "cacheNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "cache_op": "get",
          "cache_key": "api_data",
          "label": "Check Cache"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "condition": "== null",
          "label": "Cache Miss?"
        }
      },
      {
        "id": "3",
        "type": "httpNode",
        "position": {
          "x": 500,
          "y": 50
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "Fetch from API"
        }
      },
      {
        "id": "4",
        "type": "cacheNode",
        "position": {
          "x": 700,
          "y": 50
        },
        "data": {
          "cache_op": "set",
          "cache_key": "api_data",
          "label": "Store in Cache"
        }
      },
      {
        "id": "5",
        "type": "vizNode",
        "position": {
          "x": 900,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      },
      {
        "id": "e4-5",
        "source": "4",
        "target": "5"
      },
      {
        "id": "e2-5",
        "source": "2",
        "target": "5"
      }
    ]
  },
  {
    "id": "example-11",
    "title": "Webhook Handler",
    "description": "Process incoming webhook payloads with validation and signature verification. Implements retry logic for failed webhooks.",
    "tags": [
      "api",
      "http",
      "webhook",
      "validation"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-12",
    "title": "API Pagination Handler",
    "description": "Automatically handle paginated API responses. Fetches all pages and aggregates results.",
    "tags": [
      "api",
      "http",
      "pagination",
      "iteration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-13",
    "title": "Multi-Step API Workflow",
    "description": "Chain multiple API calls where each step depends on the previous response. Handles data transformation between steps.",
    "tags": [
      "api",
      "http",
      "chaining",
      "transform"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-14",
    "title": "API Error Recovery",
    "description": "Comprehensive error handling for API calls with fallback strategies and graceful degradation.",
    "tags": [
      "api",
      "http",
      "error-handling",
      "fallback"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "tryCatchNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "continue_on_error": true,
          "label": "Try-Catch"
        }
      },
      {
        "id": "2",
        "type": "httpNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 500,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-15",
    "title": "File Upload via API",
    "description": "Upload files to an API endpoint with multipart/form-data encoding. Includes progress tracking.",
    "tags": [
      "api",
      "http",
      "upload",
      "file"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-16",
    "title": "API Response Validation",
    "description": "Validate API responses against JSON schema. Ensures data integrity and type safety.",
    "tags": [
      "api",
      "http",
      "validation",
      "json"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-17",
    "title": "Bulk API Operations",
    "description": "Perform bulk operations on API resources with batch processing. Implements chunking and rate limiting.",
    "tags": [
      "api",
      "http",
      "bulk",
      "batch"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-18",
    "title": "API Health Check",
    "description": "Monitor API availability and response times. Implements health check pattern with alerting.",
    "tags": [
      "api",
      "http",
      "health-check",
      "monitoring"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-19",
    "title": "API Version Handling",
    "description": "Handle multiple API versions gracefully with version detection and routing logic.",
    "tags": [
      "api",
      "http",
      "versioning",
      "routing"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-20",
    "title": "API Request Queue",
    "description": "Queue API requests for sequential processing. Useful for rate-limited APIs.",
    "tags": [
      "api",
      "http",
      "queue",
      "sequential"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "httpNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "url": "https://api.example.com/data",
          "method": "GET",
          "label": "API Request"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-21",
    "title": "JSON Data Parsing",
    "description": "Parse JSON data from various sources and extract specific fields. Handles nested structures and arrays.",
    "tags": [
      "json",
      "parse",
      "extract",
      "data"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-22",
    "title": "Array Filtering",
    "description": "Filter array elements based on conditions. Supports complex boolean expressions and nested properties.",
    "tags": [
      "array",
      "filter",
      "condition",
      "data"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-23",
    "title": "Data Transformation Pipeline",
    "description": "Transform data through multiple stages including mapping, filtering, and aggregation operations.",
    "tags": [
      "transform",
      "pipeline",
      "data",
      "map"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-24",
    "title": "Data Aggregation",
    "description": "Aggregate data using various operations like sum, average, count, min, and max values.",
    "tags": [
      "aggregate",
      "reduce",
      "data",
      "statistics"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-25",
    "title": "CSV to JSON Converter",
    "description": "Parse CSV data and convert it to JSON format with type inference and validation.",
    "tags": [
      "csv",
      "json",
      "convert",
      "parse"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-26",
    "title": "Data Deduplication",
    "description": "Remove duplicate entries from datasets based on specified fields. Maintains data integrity.",
    "tags": [
      "unique",
      "dedupe",
      "data",
      "filter"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-27",
    "title": "Data Sorting and Ranking",
    "description": "Sort data by multiple fields with custom ordering. Calculate rankings and percentiles.",
    "tags": [
      "sort",
      "rank",
      "data",
      "order"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-28",
    "title": "Data Validation Pipeline",
    "description": "Validate data against rules and schemas. Filters out invalid entries and reports errors.",
    "tags": [
      "validation",
      "filter",
      "data",
      "schema"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-29",
    "title": "Nested Data Flattening",
    "description": "Flatten nested JSON structures into a flat representation for easier processing and analysis.",
    "tags": [
      "flatten",
      "json",
      "transform",
      "data"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-30",
    "title": "Data Enrichment",
    "description": "Enrich data by joining multiple data sources and adding computed fields.",
    "tags": [
      "join",
      "enrich",
      "data",
      "transform"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-31",
    "title": "Time Series Data Processing",
    "description": "Process time series data with windowing, aggregation, and trend analysis.",
    "tags": [
      "time-series",
      "aggregate",
      "data",
      "analysis"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-32",
    "title": "Data Sampling",
    "description": "Sample data using various strategies: random, stratified, or systematic sampling.",
    "tags": [
      "sample",
      "data",
      "statistics",
      "random"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-33",
    "title": "Data Chunking",
    "description": "Split large datasets into smaller chunks for batch processing. Useful for memory management.",
    "tags": [
      "chunk",
      "batch",
      "data",
      "split"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-34",
    "title": "Data Merging",
    "description": "Merge multiple data sources with configurable join strategies (inner, outer, left, right).",
    "tags": [
      "merge",
      "join",
      "data",
      "combine"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-35",
    "title": "Statistical Analysis",
    "description": "Calculate statistical measures including mean, median, mode, standard deviation, and variance.",
    "tags": [
      "statistics",
      "analysis",
      "data",
      "math"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-36",
    "title": "Data Pivot Table",
    "description": "Create pivot tables from raw data with grouping and aggregation capabilities.",
    "tags": [
      "pivot",
      "group",
      "aggregate",
      "data"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-37",
    "title": "Data Cleansing",
    "description": "Clean data by removing nulls, trimming whitespace, and normalizing formats.",
    "tags": [
      "clean",
      "normalize",
      "data",
      "transform"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-38",
    "title": "Complex Data Extraction",
    "description": "Extract data from deeply nested structures using path expressions and queries.",
    "tags": [
      "extract",
      "query",
      "data",
      "nested"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-39",
    "title": "Data Type Conversion",
    "description": "Convert data between different types with validation and error handling.",
    "tags": [
      "convert",
      "type",
      "data",
      "transform"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-40",
    "title": "Batch Data Processing",
    "description": "Process large datasets in batches with progress tracking and error recovery.",
    "tags": [
      "batch",
      "process",
      "data",
      "bulk"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data Input"
        }
      },
      {
        "id": "2",
        "type": "mapNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "expression": "item * 2",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-41",
    "title": "Conditional Branching",
    "description": "Route data flow based on conditions. Implements if-then-else logic with multiple branches.",
    "tags": [
      "condition",
      "branching",
      "control-flow",
      "logic"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-42",
    "title": "Switch Case Logic",
    "description": "Implement switch-case pattern with multiple paths based on input value matching.",
    "tags": [
      "switch",
      "case",
      "branching",
      "routing"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-43",
    "title": "For Each Loop",
    "description": "Iterate over arrays and process each element individually. Supports break and continue logic.",
    "tags": [
      "loop",
      "for-each",
      "iteration",
      "array"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-44",
    "title": "While Loop with Condition",
    "description": "Execute operations repeatedly while a condition is true. Includes max iteration safety.",
    "tags": [
      "loop",
      "while",
      "iteration",
      "condition"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-45",
    "title": "Try-Catch Error Handling",
    "description": "Implement error handling with try-catch pattern. Provides fallback values and error recovery.",
    "tags": [
      "error-handling",
      "try-catch",
      "fallback",
      "recovery"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-46",
    "title": "Parallel Execution",
    "description": "Execute multiple operations in parallel for improved performance. Aggregates results.",
    "tags": [
      "parallel",
      "concurrent",
      "performance",
      "async"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-47",
    "title": "Sequential Processing",
    "description": "Process operations sequentially with data passing between steps. Maintains execution order.",
    "tags": [
      "sequential",
      "chain",
      "pipeline",
      "order"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-48",
    "title": "Conditional Retry",
    "description": "Retry operations based on specific error conditions with configurable strategies.",
    "tags": [
      "retry",
      "condition",
      "error-handling",
      "resilience"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-49",
    "title": "Circuit Breaker Pattern",
    "description": "Implement circuit breaker for fault tolerance. Prevents cascading failures.",
    "tags": [
      "circuit-breaker",
      "fault-tolerance",
      "resilience",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-50",
    "title": "Timeout with Fallback",
    "description": "Execute operations with timeout protection and fallback to default values.",
    "tags": [
      "timeout",
      "fallback",
      "error-handling",
      "resilience"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-51",
    "title": "Rate Limiting Logic",
    "description": "Control execution rate to prevent system overload. Implements token bucket algorithm.",
    "tags": [
      "rate-limit",
      "throttle",
      "control-flow",
      "performance"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-52",
    "title": "Batch Processing with Delay",
    "description": "Process items in batches with configurable delays between batches.",
    "tags": [
      "batch",
      "delay",
      "processing",
      "throttle"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-53",
    "title": "Recursive Processing",
    "description": "Implement recursive data processing with depth limiting and cycle detection.",
    "tags": [
      "recursive",
      "iteration",
      "deep",
      "processing"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-54",
    "title": "Dynamic Routing",
    "description": "Route data dynamically based on runtime conditions and data content.",
    "tags": [
      "routing",
      "dynamic",
      "branching",
      "condition"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-55",
    "title": "Event-Driven Workflow",
    "description": "React to events and trigger appropriate processing flows based on event types.",
    "tags": [
      "event",
      "trigger",
      "reactive",
      "async"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-56",
    "title": "State Machine",
    "description": "Implement state machine pattern for complex state management and transitions.",
    "tags": [
      "state",
      "machine",
      "transition",
      "control-flow"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-57",
    "title": "Saga Pattern",
    "description": "Implement saga pattern for distributed transactions with compensation logic.",
    "tags": [
      "saga",
      "transaction",
      "distributed",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-58",
    "title": "Fork-Join Pattern",
    "description": "Fork execution into multiple paths and join results at the end.",
    "tags": [
      "fork",
      "join",
      "parallel",
      "merge"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-59",
    "title": "Pipeline with Checkpoints",
    "description": "Execute pipeline with checkpoint saves for recovery and resumption.",
    "tags": [
      "pipeline",
      "checkpoint",
      "recovery",
      "resilience"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-60",
    "title": "Conditional Aggregation",
    "description": "Aggregate data based on dynamic conditions and grouping criteria.",
    "tags": [
      "aggregate",
      "condition",
      "group",
      "data"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 150
        },
        "data": {
          "value": 10,
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "conditionNode",
        "position": {
          "x": 200,
          "y": 150
        },
        "data": {
          "condition": "> 5",
          "label": "Check Condition"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "True Path"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 200
        },
        "data": {
          "mode": "text",
          "label": "False Path"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e2-4",
        "source": "2",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-61",
    "title": "Variable Storage",
    "description": "Store and retrieve variables across workflow execution. Supports global and local scopes.",
    "tags": [
      "variables",
      "state",
      "storage",
      "memory"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-62",
    "title": "Counter Implementation",
    "description": "Implement counter with increment, decrement, and reset operations.",
    "tags": [
      "counter",
      "increment",
      "variables",
      "state"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-63",
    "title": "Accumulator Pattern",
    "description": "Accumulate values using various operations like sum, product, and concatenation.",
    "tags": [
      "accumulator",
      "aggregate",
      "variables",
      "state"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-64",
    "title": "Session State Management",
    "description": "Manage session state across multiple workflow executions with persistence.",
    "tags": [
      "session",
      "state",
      "variables",
      "persistence"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-65",
    "title": "Cache Implementation",
    "description": "Implement caching mechanism with TTL and eviction policies.",
    "tags": [
      "cache",
      "variables",
      "storage",
      "performance"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-66",
    "title": "Variable Interpolation",
    "description": "Interpolate variables into strings and expressions dynamically.",
    "tags": [
      "variables",
      "template",
      "interpolation",
      "string"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-67",
    "title": "Shared State Pattern",
    "description": "Share state between parallel execution branches with proper synchronization.",
    "tags": [
      "shared",
      "state",
      "variables",
      "parallel"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-68",
    "title": "Variable Validation",
    "description": "Validate variable values against constraints and schemas before use.",
    "tags": [
      "variables",
      "validation",
      "schema",
      "check"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-69",
    "title": "Dynamic Configuration",
    "description": "Load and use configuration values dynamically from variables.",
    "tags": [
      "config",
      "variables",
      "dynamic",
      "settings"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-70",
    "title": "Context Propagation",
    "description": "Propagate context variables through nested workflow executions.",
    "tags": [
      "context",
      "variables",
      "propagate",
      "scope"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "numberNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "value": 42,
          "label": "Value"
        }
      },
      {
        "id": "2",
        "type": "variableNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "set",
          "label": "Store"
        }
      },
      {
        "id": "3",
        "type": "variableNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "var_name": "myVar",
          "var_op": "get",
          "label": "Retrieve"
        }
      },
      {
        "id": "4",
        "type": "vizNode",
        "position": {
          "x": 650,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Display"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      },
      {
        "id": "e3-4",
        "source": "3",
        "target": "4"
      }
    ]
  },
  {
    "id": "example-71",
    "title": "Text Transformation",
    "description": "Transform text using operations like uppercase, lowercase, trim, and replace.",
    "tags": [
      "text",
      "transform",
      "string",
      "manipulation"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-72",
    "title": "String Concatenation",
    "description": "Concatenate multiple strings with custom delimiters and formatting.",
    "tags": [
      "text",
      "concat",
      "string",
      "join"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-73",
    "title": "Text Pattern Matching",
    "description": "Match text against patterns using regular expressions and extract matches.",
    "tags": [
      "text",
      "regex",
      "pattern",
      "match"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-74",
    "title": "Text Splitting",
    "description": "Split text into parts using delimiters and patterns. Handles edge cases.",
    "tags": [
      "text",
      "split",
      "string",
      "parse"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-75",
    "title": "Text Validation",
    "description": "Validate text against patterns, length constraints, and custom rules.",
    "tags": [
      "text",
      "validation",
      "pattern",
      "check"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-76",
    "title": "Text Formatting",
    "description": "Format text using templates with variable substitution and custom formatters.",
    "tags": [
      "text",
      "format",
      "template",
      "string"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-77",
    "title": "Text Search and Replace",
    "description": "Search for patterns in text and replace with new values. Supports regex.",
    "tags": [
      "text",
      "search",
      "replace",
      "regex"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-78",
    "title": "Text Extraction",
    "description": "Extract specific portions of text using patterns and positional indexing.",
    "tags": [
      "text",
      "extract",
      "substring",
      "parse"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-79",
    "title": "Text Normalization",
    "description": "Normalize text by removing special characters, standardizing case, and trimming.",
    "tags": [
      "text",
      "normalize",
      "clean",
      "transform"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-80",
    "title": "Markdown Processing",
    "description": "Process markdown text and convert to HTML or extract structure.",
    "tags": [
      "text",
      "markdown",
      "parse",
      "convert"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "text": "hello world",
          "label": "Input Text"
        }
      },
      {
        "id": "2",
        "type": "textOpNode",
        "position": {
          "x": 250,
          "y": 100
        },
        "data": {
          "text_op": "uppercase",
          "label": "Transform"
        }
      },
      {
        "id": "3",
        "type": "vizNode",
        "position": {
          "x": 450,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Result"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      },
      {
        "id": "e2-3",
        "source": "2",
        "target": "3"
      }
    ]
  },
  {
    "id": "example-81",
    "title": "Bar Chart Visualization",
    "description": "Create bar charts from data with customizable colors, orientation, and styling.",
    "tags": [
      "visualization",
      "chart",
      "bar",
      "display"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-82",
    "title": "Data Table Display",
    "description": "Display data in tabular format with sorting, filtering, and pagination.",
    "tags": [
      "visualization",
      "table",
      "display",
      "data"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-83",
    "title": "JSON Viewer",
    "description": "Display JSON data in a formatted, syntax-highlighted viewer with expand/collapse.",
    "tags": [
      "visualization",
      "json",
      "display",
      "viewer"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-84",
    "title": "Dashboard Layout",
    "description": "Create dashboard with multiple visualizations arranged in a grid layout.",
    "tags": [
      "visualization",
      "dashboard",
      "layout",
      "display"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-85",
    "title": "Real-time Data Display",
    "description": "Display data that updates in real-time with live refresh capabilities.",
    "tags": [
      "visualization",
      "real-time",
      "live",
      "display"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-86",
    "title": "Summary Statistics",
    "description": "Display summary statistics including counts, averages, and distributions.",
    "tags": [
      "visualization",
      "statistics",
      "summary",
      "display"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-87",
    "title": "Comparison View",
    "description": "Compare multiple datasets side-by-side with highlighting of differences.",
    "tags": [
      "visualization",
      "compare",
      "diff",
      "display"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-88",
    "title": "Timeline Visualization",
    "description": "Visualize time-based data on an interactive timeline with zoom capabilities.",
    "tags": [
      "visualization",
      "timeline",
      "time",
      "display"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-89",
    "title": "Hierarchical Data View",
    "description": "Display hierarchical data in tree or nested structure format.",
    "tags": [
      "visualization",
      "tree",
      "hierarchy",
      "display"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-90",
    "title": "Progress Indicator",
    "description": "Show progress of long-running operations with percentage and status.",
    "tags": [
      "visualization",
      "progress",
      "status",
      "display"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "rangeNode",
        "position": {
          "x": 50,
          "y": 100
        },
        "data": {
          "start": 1,
          "end": 10,
          "step": 1,
          "label": "Data"
        }
      },
      {
        "id": "2",
        "type": "barChartNode",
        "position": {
          "x": 300,
          "y": 100
        },
        "data": {
          "orientation": "vertical",
          "bar_color": "#3b82f6",
          "show_values": true,
          "label": "Chart"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-91",
    "title": "Database Query",
    "description": "Execute database queries and process results. Supports parameterized queries.",
    "tags": [
      "database",
      "query",
      "sql",
      "integration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-92",
    "title": "Email Notification",
    "description": "Send email notifications with templates, attachments, and HTML formatting.",
    "tags": [
      "email",
      "notification",
      "smtp",
      "integration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-93",
    "title": "Slack Integration",
    "description": "Send messages to Slack channels with rich formatting and attachments.",
    "tags": [
      "slack",
      "chat",
      "notification",
      "integration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-94",
    "title": "File System Operations",
    "description": "Read, write, and manipulate files on the file system with error handling.",
    "tags": [
      "file",
      "io",
      "storage",
      "integration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-95",
    "title": "Cloud Storage Upload",
    "description": "Upload files to cloud storage services like S3, Azure Blob, or GCS.",
    "tags": [
      "cloud",
      "storage",
      "upload",
      "integration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-96",
    "title": "Message Queue Publisher",
    "description": "Publish messages to message queues like RabbitMQ, Kafka, or SQS.",
    "tags": [
      "queue",
      "message",
      "publish",
      "integration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-97",
    "title": "SMS Notification",
    "description": "Send SMS notifications using Twilio or similar services.",
    "tags": [
      "sms",
      "notification",
      "twilio",
      "integration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-98",
    "title": "Payment Processing",
    "description": "Process payments using Stripe, PayPal, or similar payment gateways.",
    "tags": [
      "payment",
      "stripe",
      "transaction",
      "integration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-99",
    "title": "Analytics Tracking",
    "description": "Track events and metrics to analytics platforms like Google Analytics.",
    "tags": [
      "analytics",
      "tracking",
      "metrics",
      "integration"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-100",
    "title": "External API Integration",
    "description": "Integrate with third-party APIs with authentication and error handling.",
    "tags": [
      "api",
      "integration",
      "external",
      "third-party"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-101",
    "title": "MapReduce Pattern",
    "description": "Implement MapReduce pattern for distributed data processing at scale.",
    "tags": [
      "mapreduce",
      "distributed",
      "pattern",
      "data"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-102",
    "title": "Stream Processing",
    "description": "Process data streams in real-time with windowing and aggregation.",
    "tags": [
      "stream",
      "real-time",
      "processing",
      "async"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-103",
    "title": "ETL Pipeline",
    "description": "Extract, Transform, Load pipeline for data warehousing and analytics.",
    "tags": [
      "etl",
      "pipeline",
      "data",
      "warehouse"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-104",
    "title": "Event Sourcing",
    "description": "Implement event sourcing pattern for audit trails and state reconstruction.",
    "tags": [
      "event",
      "sourcing",
      "audit",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-105",
    "title": "CQRS Pattern",
    "description": "Separate read and write operations using Command Query Responsibility Segregation.",
    "tags": [
      "cqrs",
      "pattern",
      "architecture",
      "separation"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-106",
    "title": "Bulkhead Pattern",
    "description": "Isolate resources to prevent cascading failures using bulkhead pattern.",
    "tags": [
      "bulkhead",
      "isolation",
      "resilience",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-107",
    "title": "Cache-Aside Pattern",
    "description": "Implement cache-aside pattern for optimized data access and caching.",
    "tags": [
      "cache",
      "pattern",
      "performance",
      "data"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-108",
    "title": "Strangler Fig Pattern",
    "description": "Gradually migrate from legacy systems using strangler fig pattern.",
    "tags": [
      "migration",
      "legacy",
      "pattern",
      "gradual"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-109",
    "title": "Scatter-Gather Pattern",
    "description": "Send requests to multiple services and aggregate responses.",
    "tags": [
      "scatter",
      "gather",
      "parallel",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-110",
    "title": "Publish-Subscribe Pattern",
    "description": "Implement pub-sub pattern for event-driven communication.",
    "tags": [
      "pubsub",
      "event",
      "pattern",
      "messaging"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-111",
    "title": "Request-Response Pattern",
    "description": "Implement synchronous request-response communication pattern.",
    "tags": [
      "request",
      "response",
      "sync",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-112",
    "title": "Polling Consumer",
    "description": "Poll for messages or events at regular intervals with backoff.",
    "tags": [
      "polling",
      "consumer",
      "pattern",
      "async"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-113",
    "title": "Competing Consumers",
    "description": "Multiple consumers compete for messages from a shared queue.",
    "tags": [
      "consumer",
      "queue",
      "concurrent",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-114",
    "title": "Priority Queue",
    "description": "Process messages based on priority levels with queue management.",
    "tags": [
      "queue",
      "priority",
      "pattern",
      "ordering"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-115",
    "title": "Dead Letter Queue",
    "description": "Handle failed messages using dead letter queue pattern.",
    "tags": [
      "dlq",
      "error",
      "queue",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-116",
    "title": "Idempotent Consumer",
    "description": "Ensure operations can be safely retried without side effects.",
    "tags": [
      "idempotent",
      "retry",
      "pattern",
      "safety"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-117",
    "title": "Compensation Transaction",
    "description": "Implement compensating transactions for rollback in distributed systems.",
    "tags": [
      "compensation",
      "transaction",
      "rollback",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-118",
    "title": "Two-Phase Commit",
    "description": "Implement two-phase commit protocol for distributed transactions.",
    "tags": [
      "2pc",
      "transaction",
      "distributed",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-119",
    "title": "Leader Election",
    "description": "Implement leader election for distributed coordination.",
    "tags": [
      "leader",
      "election",
      "distributed",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-120",
    "title": "Sharding Strategy",
    "description": "Distribute data across multiple nodes using sharding.",
    "tags": [
      "sharding",
      "distributed",
      "partition",
      "pattern"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-121",
    "title": "Unit Test Workflow",
    "description": "Create unit tests for workflow components with assertions and mocking.",
    "tags": [
      "testing",
      "unit",
      "validation",
      "qa"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-122",
    "title": "Integration Test",
    "description": "Test integration between multiple components and external services.",
    "tags": [
      "testing",
      "integration",
      "e2e",
      "qa"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-123",
    "title": "Load Testing",
    "description": "Perform load testing to measure performance under heavy traffic.",
    "tags": [
      "testing",
      "load",
      "performance",
      "stress"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-124",
    "title": "Data Validation Test",
    "description": "Validate data quality and integrity with comprehensive checks.",
    "tags": [
      "testing",
      "validation",
      "data",
      "quality"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-125",
    "title": "Mock Service Responses",
    "description": "Mock external service responses for testing purposes.",
    "tags": [
      "testing",
      "mock",
      "stub",
      "qa"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-126",
    "title": "Chaos Engineering",
    "description": "Test system resilience by introducing controlled failures.",
    "tags": [
      "testing",
      "chaos",
      "resilience",
      "failure"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-127",
    "title": "Smoke Test",
    "description": "Quick smoke tests to verify basic functionality after deployment.",
    "tags": [
      "testing",
      "smoke",
      "deployment",
      "qa"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-128",
    "title": "Regression Test Suite",
    "description": "Comprehensive regression tests to catch breaking changes.",
    "tags": [
      "testing",
      "regression",
      "qa",
      "validation"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-129",
    "title": "Performance Benchmark",
    "description": "Benchmark performance metrics for optimization analysis.",
    "tags": [
      "testing",
      "benchmark",
      "performance",
      "metrics"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-130",
    "title": "Security Testing",
    "description": "Test for common security vulnerabilities and exploits.",
    "tags": [
      "testing",
      "security",
      "vulnerability",
      "qa"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-131",
    "title": "Health Check Monitor",
    "description": "Monitor system health with periodic checks and alerting.",
    "tags": [
      "monitoring",
      "health",
      "alert",
      "observability"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-132",
    "title": "Performance Metrics",
    "description": "Collect and analyze performance metrics like latency and throughput.",
    "tags": [
      "monitoring",
      "metrics",
      "performance",
      "observability"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-133",
    "title": "Error Tracking",
    "description": "Track and analyze errors with stack traces and context.",
    "tags": [
      "monitoring",
      "error",
      "logging",
      "observability"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-134",
    "title": "Resource Utilization",
    "description": "Monitor CPU, memory, and network utilization metrics.",
    "tags": [
      "monitoring",
      "resources",
      "utilization",
      "observability"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-135",
    "title": "SLA Monitoring",
    "description": "Monitor service level agreements and generate compliance reports.",
    "tags": [
      "monitoring",
      "sla",
      "compliance",
      "observability"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-136",
    "title": "Custom Metrics Collection",
    "description": "Collect custom business metrics for analytics and reporting.",
    "tags": [
      "monitoring",
      "metrics",
      "custom",
      "analytics"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-137",
    "title": "Alert Rules Engine",
    "description": "Define and execute alert rules based on metric thresholds.",
    "tags": [
      "monitoring",
      "alert",
      "rules",
      "notification"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-138",
    "title": "Log Aggregation",
    "description": "Aggregate logs from multiple sources for centralized analysis.",
    "tags": [
      "monitoring",
      "logging",
      "aggregate",
      "observability"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-139",
    "title": "Trace Collection",
    "description": "Collect distributed traces for request flow analysis.",
    "tags": [
      "monitoring",
      "trace",
      "distributed",
      "observability"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-140",
    "title": "Dashboard Metrics",
    "description": "Display key metrics on a dashboard for real-time monitoring.",
    "tags": [
      "monitoring",
      "dashboard",
      "metrics",
      "visualization"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-141",
    "title": "Cron Job Workflow",
    "description": "Execute workflows on a schedule using cron expressions.",
    "tags": [
      "scheduling",
      "cron",
      "automation",
      "periodic"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-142",
    "title": "Delayed Execution",
    "description": "Delay execution by a specified duration before processing.",
    "tags": [
      "scheduling",
      "delay",
      "timeout",
      "timing"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-143",
    "title": "Time-Based Trigger",
    "description": "Trigger workflows at specific times or dates.",
    "tags": [
      "scheduling",
      "trigger",
      "time",
      "automation"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-144",
    "title": "Recurring Tasks",
    "description": "Execute tasks repeatedly at configured intervals.",
    "tags": [
      "scheduling",
      "recurring",
      "periodic",
      "automation"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-145",
    "title": "Batch Job Scheduling",
    "description": "Schedule batch jobs for off-peak hours with dependencies.",
    "tags": [
      "scheduling",
      "batch",
      "job",
      "automation"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-146",
    "title": "Event-Driven Scheduling",
    "description": "Schedule workflows based on external events and triggers.",
    "tags": [
      "scheduling",
      "event",
      "trigger",
      "reactive"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-147",
    "title": "Dynamic Scheduling",
    "description": "Dynamically calculate and schedule execution based on conditions.",
    "tags": [
      "scheduling",
      "dynamic",
      "conditional",
      "automation"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-148",
    "title": "Priority-Based Scheduling",
    "description": "Schedule tasks based on priority levels and resource availability.",
    "tags": [
      "scheduling",
      "priority",
      "queue",
      "ordering"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-149",
    "title": "Deadline-Aware Scheduling",
    "description": "Schedule tasks with deadlines and SLA constraints.",
    "tags": [
      "scheduling",
      "deadline",
      "sla",
      "constraint"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  },
  {
    "id": "example-150",
    "title": "Workflow Orchestration",
    "description": "Orchestrate complex workflows with dependencies and timing.",
    "tags": [
      "scheduling",
      "orchestration",
      "workflow",
      "automation"
    ],
    "nodes": [
      {
        "id": "1",
        "type": "textInputNode",
        "position": {
          "x": 100,
          "y": 100
        },
        "data": {
          "text": "Input",
          "label": "Input"
        }
      },
      {
        "id": "2",
        "type": "vizNode",
        "position": {
          "x": 400,
          "y": 100
        },
        "data": {
          "mode": "text",
          "label": "Output"
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2"
      }
    ]
  }
];

