# Examples & Tutorials

## Basic Examples

### Example 1: Simple Addition

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

### Example 2: Text Processing

```json
{
  "nodes": [
    {"id": "1", "type": "text_input", "data": {"text": "hello world"}},
    {"id": "2", "type": "text_operation", "data": {"textOp": "uppercase"}},
    {"id": "3", "type": "visualization", "data": {"mode": "text"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"}
  ]
}
```

### Example 3: HTTP Request

```json
{
  "nodes": [
    {"id": "1", "type": "http", "data": {
      "url": "https://api.github.com/users/github"
    }},
    {"id": "2", "type": "extract", "data": {"field": "name"}},
    {"id": "3", "type": "visualization", "data": {"mode": "text"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"}
  ]
}
```

### Example 4: Array Processing

```json
{
  "nodes": [
    {"id": "1", "type": "range", "data": {"start": 1, "end": 100}},
    {"id": "2", "type": "filter", "data": {"expression": "x % 2 == 0"}},
    {"id": "3", "type": "map", "data": {"expression": "x * x"}},
    {"id": "4", "type": "reduce", "data": {"op": "sum"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"},
    {"source": "3", "target": "4"}
  ]
}
```

For more examples, see the `examples/` directory.

---

**Last Updated:** 2025-11-03
**Version:** 1.0
