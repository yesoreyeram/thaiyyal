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

## Workflow Examples Library

The Thaiyyal Workflow Builder includes a comprehensive library of 150+ pre-built workflow examples that you can use as templates or learning resources. These examples cover a wide range of use cases and patterns.

### Accessing the Examples Library

1. Open the Workflow Builder page
2. Click the **"Open"** button in the top navigation bar
3. The Workflow Examples Library modal will appear

![Workflow Examples Modal](https://github.com/user-attachments/assets/02bf2008-efef-4477-8447-64bb4db3bdc9)

### Features

#### Search and Filter
- **Search Bar**: Search examples by title, description, or tags
- **Tag Filtering**: Filter examples by clicking popular tags (api, http, json, data, branching, variables, error-handling, performance)
- **Multiple Filters**: Combine multiple tag filters to narrow down results

![Filtered Examples](https://github.com/user-attachments/assets/b9ba63af-a10c-499e-b6af-80d19fc52bd1)

#### Browse Examples
- Browse through 150+ categorized workflow examples
- Each card displays:
  - **Title**: Clear, descriptive name (max 200 characters)
  - **Description**: Detailed explanation of what the workflow does (max 2000 characters)
  - **Tags**: 3-4 relevant tags for easy filtering

#### Load Examples
- Click any example card to instantly load it into the canvas
- The workflow title automatically updates to match the example name
- Start customizing immediately or use as-is

![Loaded Example](https://github.com/user-attachments/assets/56e3a202-1e87-4c8d-b1c8-973ba3a85396)

### Example Categories

The library includes examples across multiple categories:

- **API & HTTP** (20 examples): API calls, authentication, retry logic, webhooks, pagination
- **Data Processing** (20 examples): JSON parsing, filtering, transformation, aggregation, validation
- **Control Flow** (20 examples): Branching, loops, error handling, parallel execution, state machines
- **Variables & State** (10 examples): Variable storage, counters, accumulators, caching
- **Text Processing** (10 examples): String operations, pattern matching, formatting
- **Visualization** (10 examples): Charts, tables, dashboards, real-time displays
- **Integration** (10 examples): Database, email, Slack, file systems, cloud storage
- **Advanced Patterns** (20 examples): MapReduce, ETL, event sourcing, CQRS, distributed patterns
- **Testing** (10 examples): Unit tests, integration tests, load testing, mocking
- **Monitoring** (10 examples): Health checks, metrics, alerts, logging, tracing
- **Scheduling** (10 examples): Cron jobs, delays, triggers, orchestration

### Tips for Using Examples

1. **Start Simple**: Begin with basic examples and gradually explore more complex patterns
2. **Customize**: Use examples as templates and modify them for your specific needs
3. **Combine**: Mix and match concepts from different examples
4. **Learn Patterns**: Study the examples to understand workflow design patterns and best practices

---

**Last Updated:** 2025-11-04
**Version:** 1.1
