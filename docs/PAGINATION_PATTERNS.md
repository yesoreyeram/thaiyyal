# HTTP Pagination Using Composable Nodes

This guide demonstrates how to implement HTTP pagination using Thaiyyal's composable workflow nodes. Instead of a monolithic pagination node, we use existing control flow, HTTP, and data processing nodes to create flexible, reusable pagination patterns.

## Philosophy: Composition Over Monoliths

The composable approach provides:
- ✅ **Flexibility**: Customize pagination logic for different APIs
- ✅ **Reusability**: Same nodes work for various use cases
- ✅ **Extensibility**: Add retry logic, rate limiting, caching
- ✅ **Clarity**: Each node has a single, clear responsibility

## Available Building Blocks

| Node | Purpose | Example |
|------|---------|---------|
| **Counter** | Track page/offset numbers | Current page: 1, 2, 3... |
| **Condition** | Check loop continuation | `page <= 5` |
| **HTTP** | Make API requests | `GET /items?page={N}` |
| **Variable** | Store state (URL, cursors) | Base URL, next cursor |
| **Extract** | Parse response data | Get `items` array |
| **Accumulator** | Collect all results | Merge all pages |
| **While Loop** | Repeat until condition false | Loop while has_more |
| **Transform** | Reshape data | Flatten nested arrays |

## Pagination Pattern 1: Page-Based (Common)

**Use Case**: Fetch 5 pages with 10 items each (50 items total)

### Workflow Pattern

```
┌──────────┐     ┌───────────┐     ┌──────┐     ┌──────────┐
│ Counter  │ ──▶ │ Condition │ ──▶ │ HTTP │ ──▶ │ Extract  │
│ (page=1) │     │ page<=5   │     │      │     │ .items   │
└──────────┘     └───────────┘     └──────┘     └──────────┘
                                                      │
                                                      ▼
                                              ┌──────────────┐
                                              │ Accumulator  │
                                              │ (all items)  │
                                              └──────────────┘
```

### JSON Workflow

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "counter",
      "data": {
        "counter_op": "increment",
        "delta": 1,
        "label": "Page Counter"
      }
    },
    {
      "id": "2",
      "type": "condition",
      "data": {
        "condition": "<=5",
        "label": "Check Page Limit"
      }
    },
    {
      "id": "3",
      "type": "http",
      "data": {
        "url": "https://api.example.com/items?page={page}",
        "label": "Fetch Page"
      }
    },
    {
      "id": "4",
      "type": "extract",
      "data": {
        "field": "items",
        "label": "Extract Items"
      }
    },
    {
      "id": "5",
      "type": "accumulator",
      "data": {
        "accum_op": "array",
        "label": "Collect All Items"
      }
    }
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "4"},
    {"source": "4", "target": "5"}
  ]
}
```

### How It Works

1. **Counter** starts at 1, increments each iteration
2. **Condition** checks if page ≤ 5, continues if true
3. **HTTP** fetches `https://api.example.com/items?page=1`
4. **Extract** gets the `items` array from response
5. **Accumulator** collects items from all pages
6. Repeat until page > 5

## Pagination Pattern 2: Offset-Based

**Use Case**: APIs that use `?offset=0&limit=10`

### Workflow Pattern

```
┌──────────┐     ┌──────────┐     ┌──────┐     ┌──────────┐
│ Counter  │ ──▶ │ Multiply │ ──▶ │ HTTP │ ──▶ │ Extract  │
│ (n=0,1,2)│     │ n * 10   │     │      │     │          │
└──────────┘     └──────────┘     └──────┘     └──────────┘
```

### JSON Workflow

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "counter",
      "data": {
        "counter_op": "increment",
        "delta": 1,
        "initial_value": 0,
        "label": "Offset Counter"
      }
    },
    {
      "id": "2",
      "type": "operation",
      "data": {
        "op": "multiply",
        "label": "Calculate Offset (n * 10)"
      }
    },
    {
      "id": "3",
      "type": "number",
      "data": {
        "value": 10,
        "label": "Page Size"
      }
    },
    {
      "id": "4",
      "type": "http",
      "data": {
        "url": "https://api.example.com/items?offset={offset}&limit=10",
        "label": "Fetch Page"
      }
    }
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "3", "target": "2"},
    {"source": "2", "target": "4"}
  ]
}
```

### How It Works

1. **Counter** generates 0, 1, 2, 3, 4...
2. **Multiply** by 10 → offsets: 0, 10, 20, 30, 40...
3. **HTTP** fetches with calculated offset
4. Continue until desired number of pages

## Pagination Pattern 3: Cursor-Based

**Use Case**: APIs that return `next_cursor` in responses

### Workflow Pattern

```
┌──────────┐     ┌──────┐     ┌─────────┐     ┌──────────┐
│ Variable │ ──▶ │ HTTP │ ──▶ │ Extract │ ──▶ │ Variable │
│ (cursor) │     │      │     │ .next   │     │ (update) │
└──────────┘     └──────┘     └─────────┘     └──────────┘
     ▲                                              │
     └──────────────────────────────────────────────┘
```

### JSON Workflow

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "variable",
      "data": {
        "var_name": "cursor",
        "var_op": "get",
        "label": "Get Cursor"
      }
    },
    {
      "id": "2",
      "type": "http",
      "data": {
        "url": "https://api.example.com/items?cursor={cursor}",
        "label": "Fetch Page"
      }
    },
    {
      "id": "3",
      "type": "extract",
      "data": {
        "field": "next_cursor",
        "label": "Extract Next Cursor"
      }
    },
    {
      "id": "4",
      "type": "variable",
      "data": {
        "var_name": "cursor",
        "var_op": "set",
        "label": "Update Cursor"
      }
    },
    {
      "id": "5",
      "type": "condition",
      "data": {
        "condition": "!=null",
        "label": "Has More Pages?"
      }
    }
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "4"},
    {"source": "3", "target": "5"}
  ]
}
```

### How It Works

1. **Variable (get)** retrieves current cursor (starts empty/null)
2. **HTTP** fetches with cursor parameter
3. **Extract** gets `next_cursor` from response
4. **Variable (set)** updates cursor for next iteration
5. **Condition** checks if cursor exists (has more pages)
6. Repeat until cursor is null

## Pagination Pattern 4: Until-Empty

**Use Case**: Fetch until API returns empty array

### Workflow Pattern

```
┌──────┐     ┌─────────┐     ┌───────────┐     ┌──────────┐
│ HTTP │ ──▶ │ Extract │ ──▶ │ Condition │ ──▶ │Continue? │
│      │     │ .items  │     │ not empty │     │          │
└──────┘     └─────────┘     └───────────┘     └──────────┘
```

### JSON Workflow

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "counter",
      "data": {
        "counter_op": "increment",
        "label": "Page Counter"
      }
    },
    {
      "id": "2",
      "type": "http",
      "data": {
        "url": "https://api.example.com/items?page={page}",
        "label": "Fetch Page"
      }
    },
    {
      "id": "3",
      "type": "extract",
      "data": {
        "field": "items",
        "label": "Extract Items"
      }
    },
    {
      "id": "4",
      "type": "transform",
      "data": {
        "transform_type": "keys",
        "label": "Check Array Length"
      }
    },
    {
      "id": "5",
      "type": "condition",
      "data": {
        "condition": ">0",
        "label": "Has Items?"
      }
    }
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "4"},
    {"source": "4", "target": "5"}
  ]
}
```

### How It Works

1. **Counter** increments page number
2. **HTTP** fetches current page
3. **Extract** gets items array
4. **Transform** checks array length
5. **Condition** stops if array is empty
6. Repeat while items exist

## Advanced Patterns

### Error Handling

Add error detection and retry logic:

```json
{
  "nodes": [
    {
      "id": "http",
      "type": "http",
      "data": {"url": "..."}
    },
    {
      "id": "error_check",
      "type": "condition",
      "data": {"condition": "==200"}
    },
    {
      "id": "retry_counter",
      "type": "counter",
      "data": {"counter_op": "increment"}
    },
    {
      "id": "retry_limit",
      "type": "condition",
      "data": {"condition": "<=3"}
    }
  ]
}
```

### Rate Limiting

Add delays between requests:

```json
{
  "nodes": [
    {
      "id": "http",
      "type": "http",
      "data": {"url": "..."}
    },
    {
      "id": "delay",
      "type": "delay",
      "data": {"duration": "1s"}
    },
    {
      "id": "next_page",
      "type": "counter",
      "data": {"counter_op": "increment"}
    }
  ],
  "edges": [
    {"source": "http", "target": "delay"},
    {"source": "delay", "target": "next_page"}
  ]
}
```

### Parallel Page Fetching

Fetch multiple pages concurrently:

```json
{
  "nodes": [
    {
      "id": "page_numbers",
      "type": "transform",
      "data": {
        "transform_type": "to_array",
        "label": "Pages [1,2,3,4,5]"
      }
    },
    {
      "id": "parallel",
      "type": "parallel",
      "data": {
        "max_concurrency": 3,
        "label": "Fetch 3 pages at once"
      }
    },
    {
      "id": "http",
      "type": "http",
      "data": {"url": "..."}
    }
  ]
}
```

## Real-World Example: GitHub API

Fetch 5 pages of repositories:

```json
{
  "nodes": [
    {
      "id": "page",
      "type": "counter",
      "data": {"counter_op": "increment", "delta": 1}
    },
    {
      "id": "check_limit",
      "type": "condition",
      "data": {"condition": "<=5"}
    },
    {
      "id": "github_api",
      "type": "http",
      "data": {
        "url": "https://api.github.com/users/nodejs/repos?page={page}&per_page=10"
      }
    },
    {
      "id": "extract_repos",
      "type": "extract",
      "data": {"field": "items"}
    },
    {
      "id": "collect",
      "type": "accumulator",
      "data": {"accum_op": "array"}
    },
    {
      "id": "visualize",
      "type": "visualization",
      "data": {"mode": "text"}
    }
  ],
  "edges": [
    {"source": "page", "target": "check_limit"},
    {"source": "check_limit", "target": "github_api"},
    {"source": "github_api", "target": "extract_repos"},
    {"source": "extract_repos", "target": "collect"},
    {"source": "collect", "target": "visualize"}
  ]
}
```

## Comparison: Composable vs Monolithic

### Composable Approach (Current)

**Pros:**
- ✅ Flexible for different pagination types
- ✅ Easy to add custom logic (retry, cache, etc.)
- ✅ Reusable nodes for other workflows
- ✅ Clear visual representation
- ✅ Easy to understand and debug

**Cons:**
- ⚠️ Requires more nodes for simple cases
- ⚠️ User needs to understand composition

### Monolithic Pagination Node (Alternative)

**Pros:**
- ✅ Simple for basic pagination
- ✅ Fewer nodes in workflow

**Cons:**
- ❌ Limited to built-in pagination types
- ❌ Hard to extend with custom logic
- ❌ Not reusable for other patterns
- ❌ Black box - harder to debug

## Best Practices

1. **Start Simple**: Begin with page-based pattern, add complexity as needed
2. **Test Incrementally**: Test with 1-2 pages before full pagination
3. **Handle Errors**: Always add condition nodes to check for errors
4. **Use Variables**: Store API base URLs and parameters in variables
5. **Add Delays**: Use delay nodes to respect rate limits
6. **Visualize Results**: Add visualization nodes to see intermediate results
7. **Document Workflows**: Add labels to nodes explaining their purpose

## Summary

HTTP pagination is achievable using Thaiyyal's existing composable nodes:

| Pattern | Key Nodes | Use Case |
|---------|-----------|----------|
| Page-based | Counter + Condition + HTTP | Most common REST APIs |
| Offset-based | Counter + Multiply + HTTP | Older APIs with offset/limit |
| Cursor-based | Variable + Extract + HTTP | Modern paginated APIs |
| Until-empty | HTTP + Extract + Condition | Unknown page count |

The composable approach provides flexibility and reusability while maintaining clarity in workflow design.
