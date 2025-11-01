# Workflow Node Types Reference

This document provides a comprehensive reference of all node types in the Thaiyyal workflow engine, including implemented nodes and planned nodes for observability and analytics pipelines.

## Implementation Status Legend

- ✅ **Implemented** - Fully implemented and tested
- 🚧 **In Progress** - Partially implemented
- 📋 **Planned** - Documented but not yet implemented

---

## Currently Implemented Nodes

| Node Type | Description | Inputs | Outputs | Example | Status |
|-----------|-------------|--------|---------|---------|--------|
| **Number** | Numeric input node | None | `number` | `{"data": {"value": 42}}` | ✅ |
| **Operation** | Arithmetic operations (add, subtract, multiply, divide) | 2+ `number` | `number` | `{"data": {"op": "add"}}` | ✅ |
| **Visualization** | Output formatting for display | 1+ `any` | `object` with mode and value | `{"data": {"mode": "text"}}` | ✅ |
| **Text Input** | Text string input | None | `string` | `{"data": {"text": "Hello"}}` | ✅ |
| **Text Operation** | Text transformations | 1+ `string` | `string` | See operations below | ✅ |
| **HTTP** | HTTP GET request | None | `string` (response body) | `{"data": {"url": "https://api.example.com"}}` | ✅ |
| **Condition** | Conditional branching/validation | 1 `any` | `object` with value and condition_met | `{"type": "condition", "data": {"condition": ">100"}}` | ✅ |
| **For Each** | Iterate over array elements | 1 `array` | `object` with iteration metadata | `{"type": "for_each", "data": {"max_iterations": 1000}}` | ✅ |
| **While Loop** | Loop while condition is true | 1 `any` | `object` with final value and iterations | `{"type": "while_loop", "data": {"condition": "<10", "max_iterations": 100}}` | ✅ |
| **Filter** | Filter JSON array elements by condition | 1 `array` or `any` | `object` with filtered array and metadata | `{"type": "filter", "data": {"condition": "variables.item.age > 18"}}` | ✅ |
| **Variable** | Store/retrieve values across workflow | 1 `any` (for set) / None (for get) | `object` with var_name, operation, value | `{"data": {"var_name": "x", "var_op": "set"}}` | ✅ |
| **Extract** | Extract fields from objects | 1 `object` | `object` with extracted fields | `{"type": "extract", "data": {"field": "name"}}` | ✅ |
| **Transform** | Transform data structures | 1 `any` | Varies by transform type | `{"type": "transform", "data": {"transform_type": "to_array"}}` | ✅ |
| **Accumulator** | Accumulate values over time | 1 `any` | `object` with operation and accumulated value | `{"type": "accumulator", "data": {"accum_op": "sum"}}` | ✅ |
| **Counter** | Increment/decrement/reset counter | None | `object` with operation and counter value | `{"type": "counter", "data": {"counter_op": "increment"}}` | ✅ |
| **Switch** | Multi-way branching based on value/condition | 1 `any` | `object` with matched case info | `{"type": "switch", "data": {"cases": [{"when": ">100"}]}}` | ✅ |
| **Parallel** | Execute multiple branches concurrently | Multiple `any` | `object` with results array | `{"type": "parallel", "data": {"max_concurrency": 10}}` | ✅ |
| **Join** | Combine outputs from multiple nodes | Multiple `any` | `object` with combined values | `{"type": "join", "data": {"join_strategy": "all"}}` | ✅ |
| **Split** | Split single input to multiple paths | 1 `any` | `object` with path outputs | `{"type": "split", "data": {"paths": ["a", "b"]}}` | ✅ |
| **Delay** | Pause execution for specified duration | 1 `any` | Same as input with delay metadata | `{"type": "delay", "data": {"duration": "5s"}}` | ✅ |
| **Cache** | Get/set cached values with TTL | 1 `any` (for set) | `object` with cache operation result | `{"type": "cache", "data": {"cache_op": "get", "cache_key": "k1"}}` | ✅ |

### Text Operations (Sub-types)

| Operation | Description | Additional Fields | Example Input | Example Output |
|-----------|-------------|-------------------|---------------|----------------|
| `uppercase` | Convert to uppercase | None | "hello" | "HELLO" |
| `lowercase` | Convert to lowercase | None | "HELLO" | "hello" |
| `titlecase` | Capitalize each word | None | "hello world" | "Hello World" |
| `camelcase` | Convert to camelCase | None | "hello world" | "helloWorld" |
| `inversecase` | Swap character case | None | "HeLLo" | "hEllO" |
| `concat` | Concatenate multiple inputs | `separator` (optional) | ["Hello", "World"] | "HelloWorld" or "Hello World" |
| `repeat` | Repeat text n times | `repeat_n` (required) | "Ha" with n=3 | "HaHaHa" |

### Control Flow Node Details

#### Condition Node
- **Purpose**: Evaluate conditions and pass through input values
- **Condition Syntax**: `>N`, `<N`, `>=N`, `<=N`, `==N`, `!=N`, `true`, `false`
- **Output**: Returns `{"value": input, "condition_met": boolean, "condition": string}`
- **Use Cases**: Validate thresholds, filter data, conditional routing

#### For Each Node
- **Purpose**: Iterate over array elements
- **Configuration**: Optional `max_iterations` (default: 1000)
- **Output**: Returns `{"items": array, "count": int, "iterations": int}`
- **Use Cases**: Batch processing, parallel operations, array transformations

#### While Loop Node
- **Purpose**: Loop while a condition remains true
- **Configuration**: Required `condition`, optional `max_iterations` (default: 100)
- **Output**: Returns `{"final_value": any, "iterations": int, "condition": string}`
- **Use Cases**: Retry logic, iterative processing, threshold monitoring
- **Note**: Current implementation tracks iterations but doesn't modify values in loop

#### Filter Node
- **Purpose**: Filter JSON array elements based on expression conditions
- **Configuration**: Required `condition` (expression string)
- **Input Behavior**:
  - **Array input**: Filters elements where expression evaluates to true
  - **Non-array input**: Passes through unchanged with warning log
- **Expression Syntax**: ✨ **Use `item.field` to reference array elements**
  - **`item`** - refers to the current array element being evaluated
  - **`item.field`** - access a field of the current element
  - **`item.nested.field`** - access nested fields
  - Full expression language support:
    - Comparisons: `>`, `<`, `>=`, `<=`, `==`, `!=`
    - Boolean operators: `&&`, `||`, `!`
    - Variable references: `variables.name`
    - Context references: `context.name`
    - Node references: `node.id.value`
    - Functions: Math, date/time, string functions
- **Output**: Returns object with:
  - `filtered`: Filtered array (or original input if not array)
  - `input_count`: Number of input elements
  - `output_count`: Number of filtered elements
  - `skipped_count`: Number of elements that didn't match
  - `error_count`: Number of evaluation errors
  - `condition`: The filter expression used
  - `is_array`: Boolean indicating if input was an array
  - `warning`: Present if input was not an array
- **Use Cases**: 
  - Filter API response arrays by field values
  - Remove invalid/incomplete data from arrays
  - Select items matching business rules
  - Data quality filtering
- **Example Expressions** (RECOMMENDED SYNTAX):
  ```json
  // Filter numbers greater than 10
  {"condition": "item > 10"}
  
  // Filter objects by field value (MOST COMMON)
  {"condition": "item.age >= 18"}
  
  // Filter by string field
  {"condition": "item.status == \"active\""}
  
  // Complex condition with AND
  {"condition": "item.price < 50 && item.category == \"books\""}
  
  // Filter using variable threshold
  {"condition": "item.age >= variables.minAge"}
  
  // Filter using context variable
  {"condition": "item.score > context.passingScore"}
  
  // Filter with nested field access
  {"condition": "item.profile.verified == true"}
  
  // Complex business logic
  {"condition": "item.age >= 18 && item.active == true && item.verified == true"}
  ```
- **Alternative Syntaxes** (also supported):
  ```json
  // Explicit variable reference (verbose)
  {"condition": "variables.item.age >= 18"}
  
  // Direct field access (less explicit, but works)
  {"condition": "age >= 18"}
  ```
- **Complete Example**:
  ```json
  {
    "nodes": [
      {
        "id": "1",
        "type": "http",
        "data": {"url": "https://api.example.com/users"}
      },
      {
        "id": "2", 
        "type": "filter",
        "data": {"condition": "item.age >= 18 && item.active == true"}
      },
      {
        "id": "3",
        "type": "visualization",
        "data": {"mode": "json"}
      }
    ],
    "edges": [
      {"source": "1", "target": "2"},
      {"source": "2", "target": "3"}
    ]
  }
  ```

### State & Memory Node Details

#### Variable Node
- **Purpose**: Store and retrieve values across the workflow execution
- **Operations**:
  - `set`: Store a value with a given name
  - `get`: Retrieve a previously stored value
- **Configuration**: 
  - `var_name`: Name of the variable (required)
  - `var_op`: Operation to perform - "set" or "get" (required)
- **Output**: Returns `{"var_name": string, "operation": string, "value": any}`
- **Use Cases**: 
  - Store intermediate calculation results
  - Share data between disconnected parts of workflow
  - Cache values for reuse
- **Example**:
  ```json
  {
    "nodes": [
      {"id": "1", "data": {"value": 42}},
      {"id": "2", "data": {"var_name": "result", "var_op": "set"}},
      {"id": "3", "data": {"var_name": "result", "var_op": "get"}}
    ],
    "edges": [
      {"source": "1", "target": "2"},
      {"source": "2", "target": "3"}
    ]
  }
  ```

#### Extract Node
- **Purpose**: Extract specific fields from objects
- **Configuration**:
  - `field`: Extract single field (string)
  - `fields`: Extract multiple fields (array of strings)
- **Input**: Object/map
- **Output**: 
  - Single field: `{"field": string, "value": any}`
  - Multiple fields: Object with requested fields
- **Use Cases**:
  - Extract specific metrics from API responses
  - Filter object properties
  - Destructure complex data
- **Example**:
  ```json
  {
    "nodes": [
      {"id": "1", "type": "variable", "data": {"var_name": "user", "var_op": "get"}},
      {"id": "2", "type": "extract", "data": {"fields": ["name", "email"]}}
    ],
    "edges": [{"source": "1", "target": "2"}]
  }
  ```

#### Transform Node
- **Purpose**: Transform data structures between different formats
- **Transform Types**:
  - `to_array`: Convert inputs to array
  - `to_object`: Convert array of key-value pairs to object
  - `flatten`: Flatten nested arrays
  - `keys`: Extract keys from object
  - `values`: Extract values from object
- **Configuration**: `transform_type` (required)
- **Use Cases**:
  - Prepare data for different node types
  - Restructure API responses
  - Aggregate multiple inputs
- **Examples**:
  ```json
  // to_array: Collect multiple inputs
  {
    "nodes": [
      {"id": "1", "data": {"value": 10}},
      {"id": "2", "data": {"value": 20}},
      {"id": "3", "type": "transform", "data": {"transform_type": "to_array"}}
    ],
    "edges": [
      {"source": "1", "target": "3"},
      {"source": "2", "target": "3"}
    ]
  }
  
  // keys: Extract object keys
  {
    "nodes": [
      {"id": "1", "type": "variable", "data": {"var_name": "metrics", "var_op": "get"}},
      {"id": "2", "type": "transform", "data": {"transform_type": "keys"}}
    ],
    "edges": [{"source": "1", "target": "2"}]
  }
  ```

#### Accumulator Node
- **Purpose**: Accumulate values across multiple node executions
- **Operations**:
  - `sum`: Sum numeric values
  - `product`: Multiply numeric values
  - `concat`: Concatenate strings
  - `array`: Collect values into array
  - `count`: Count number of values
- **Configuration**:
  - `accum_op`: Operation to perform (required)
  - `initial_value`: Starting value (optional)
- **Output**: Returns `{"operation": string, "value": any}`
- **Use Cases**:
  - Calculate running totals
  - Build result arrays
  - Count processed items
  - Concatenate log messages
- **Example**:
  ```json
  {
    "nodes": [
      {"id": "1", "data": {"value": 10}},
      {"id": "2", "type": "accumulator", "data": {"accum_op": "sum"}},
      {"id": "3", "data": {"value": 20}},
      {"id": "4", "type": "accumulator", "data": {"accum_op": "sum"}}
    ],
    "edges": [
      {"source": "1", "target": "2"},
      {"source": "2", "target": "3"},
      {"source": "3", "target": "4"}
    ]
  }
  // Result: {"operation": "sum", "value": 30}
  ```

#### Counter Node
- **Purpose**: Simple counter with increment/decrement/reset operations
- **Operations**:
  - `increment`: Increase counter (default delta: 1)
  - `decrement`: Decrease counter (default delta: 1)
  - `reset`: Reset to initial value (default: 0)
  - `get`: Get current counter value
- **Configuration**:
  - `counter_op`: Operation to perform (required)
  - `delta`: Amount to increment/decrement (optional, default: 1)
  - `initial_value`: Value for reset operation (optional, default: 0)
- **Output**: Returns `{"operation": string, "value": number}`
- **Use Cases**:
  - Track number of iterations
  - Count events or occurrences
  - Generate sequence numbers
- **Example**:
  ```json
  {
    "nodes": [
      {"id": "1", "type": "counter", "data": {"counter_op": "increment"}},
      {"id": "2", "type": "counter", "data": {"counter_op": "increment"}},
      {"id": "3", "type": "counter", "data": {"counter_op": "get"}}
    ],
    "edges": [
      {"source": "1", "target": "2"},
      {"source": "2", "target": "3"}
    ]
  }
  // Result: {"operation": "get", "value": 2}
  ```

---

## Planned Nodes for Observability & Analytics

### Data Ingestion & Connectors

| Node Type | Description | Inputs | Outputs | Example Config | Priority | Status |
|-----------|-------------|--------|---------|----------------|----------|--------|
| **Prometheus Query** | Query Prometheus metrics | Time range params | Array of metric samples | `{"query": "up", "start": "1h"}` | High | 📋 |
| **Loki Query** | Query Loki logs | Time range, label filters | Array of log lines | `{"query": "{app=\"web\"}", "limit": 100}` | High | 📋 |
| **Database Query** | SQL query execution | Query params | Array of rows | `{"sql": "SELECT * FROM metrics", "db": "postgres"}` | Medium | 📋 |
| **Elasticsearch Query** | Query Elasticsearch | Query DSL | Array of documents | `{"index": "logs-*", "query": {...}}` | Medium | 📋 |
| **ClickHouse Query** | Query ClickHouse | SQL query | Array of rows | `{"sql": "SELECT * FROM events"}` | Medium | 📋 |
| **S3 Reader** | Read objects from S3 | Bucket/key | Object content | `{"bucket": "logs", "key": "file.json"}` | Low | 📋 |
| **Kafka Consumer** | Consume Kafka messages | Topic, group | Message batch | `{"topic": "metrics", "group": "processor"}` | Medium | 📋 |

### Data Transformation

| Node Type | Description | Inputs | Outputs | Example Config | Priority | Status |
|-----------|-------------|--------|---------|----------------|----------|--------|
| **JSON Parse** | Parse JSON string | `string` | `object` | `{"strict": true}` | High | 📋 |
| **JSON Path** | Extract fields via JSONPath | `object` | `any` | `{"path": "$.data.items[*].name"}` | High | 📋 |
| **JQ Transform** | Transform using jq syntax | `object` | `any` | `{"filter": ".items | map(.name)"}` | Medium | 📋 |
| **CSV Parse** | Parse CSV to array | `string` | `array` | `{"delimiter": ",", "headers": true}` | Medium | 📋 |
| **Template** | String interpolation | Multiple inputs | `string` | `{"template": "Alert: {{.metric}} = {{.value}}"}` | High | 📋 |
| **Regex Extract** | Extract using regex | `string` | `object` with captures | `{"pattern": "error: (.+)", "group": 1}` | Medium | 📋 |
| **Filter** | Filter array elements | `array` | `array` | `{"condition": "value > 100"}` | High | 📋 |
| **Map** | Transform each array element | `array` | `array` | `{"transform": "uppercase"}` | High | 📋 |
| **Reduce** | Aggregate array values | `array` | `any` | `{"operation": "sum"}` | Medium | 📋 |

### Analytics & Aggregation

| Node Type | Description | Inputs | Outputs | Example Config | Priority | Status |
|-----------|-------------|--------|---------|----------------|----------|--------|
| **Statistical** | Calculate statistics | Array of numbers | Object with stats | `{"operations": ["mean", "p95", "stddev"]}` | High | 📋 |
| **Time Bucket** | Group by time windows | Array with timestamps | Bucketed array | `{"window": "5m", "aggregation": "avg"}` | High | 📋 |
| **Moving Average** | Calculate moving average | Array of numbers | Array of numbers | `{"window": 10}` | Medium | 📋 |
| **Anomaly Detection** | Detect anomalies | Array of numbers | Array with scores | `{"method": "zscore", "threshold": 3}` | Medium | 📋 |
| **Correlation** | Calculate correlation | 2 number arrays | Number (correlation) | `{"method": "pearson"}` | Low | 📋 |
| **Group By** | Group and aggregate | Array of objects | Grouped object | `{"key": "host", "aggregation": "count"}` | High | 📋 |

### Control Flow & Looping

| Node Type | Description | Inputs | Outputs | Example Config | Priority | Status |
|-----------|-------------|--------|---------|----------------|----------|--------|
| **If/Condition** | Conditional branching | 1+ inputs | Same as input | `{"condition": "value > 100"}` | High | ✅ |
| **Switch** | Multi-way branching | 1 input | Same as input | `{"cases": [{"when": ">100", "output_path": "high"}]}` | High | ✅ |
| **For Each** | Iterate over array | Array | Array (processed) | `{"max_iterations": 1000}` | High | ✅ |
| **While Loop** | Loop with condition | 1+ inputs | Last iteration output | `{"condition": "count < 10", "max_iterations": 100}` | Medium | ✅ |
| **Parallel** | Execute in parallel | Array | Array (results) | `{"max_concurrency": 10}` | Medium | ✅ |
| **Join/Merge** | Combine multiple inputs | Multiple inputs | Combined output | `{"join_strategy": "all"}` | High | ✅ |
| **Split** | Split into multiple paths | 1 input | Multiple outputs | `{"paths": ["path1", "path2"]}` | Medium | ✅ |
| **Delay** | Wait/pause execution | 1+ inputs | Same as input | `{"duration": "5s"}` | Low | ✅ |
| **Cache** | Cache get/set operations | Value (for set) | Cached value or status | `{"cache_op": "get", "cache_key": "key1", "ttl": "5m"}` | Medium | ✅ |
| **Retry** | Retry failed operations with backoff | 1+ inputs | Result or error | `{"max_attempts": 3, "backoff": "exponential"}` | High | 📋 |
| **Try-Catch** | Handle errors gracefully | 1+ inputs | Value or fallback | `{"fallback_value": null}` | High | 📋 |
| **Timeout** | Enforce time limits | 1+ inputs | Result or timeout error | `{"timeout": "30s"}` | Medium | 📋 |
| **Throttle** | Rate limit execution | 1+ inputs | Same as input (delayed) | `{"rate": "100/s", "burst": 10}` | High | 📋 |
| **Batch** | Collect items into batches | Array or stream | Batched arrays | `{"batch_size": 100, "batch_timeout": "5s"}` | High | 📋 |
| **Filter** | Filter array elements | Array | Filtered array | `{"condition": "value > 100"}` | High | 📋 |
| **Route** | Content-based routing | 1 input | Routed to path | `{"routes": [{"condition": "critical", "path": "alert"}]}` | High | 📋 |
| **Map** | Transform array elements | Array | Transformed array | `{"operation": "uppercase", "parallel": true}` | High | 📋 |
| **Reduce** | Aggregate array to value | Array | Single value | `{"operation": "sum"}` | High | 📋 |
| **Window** | Time-based windowing | Array with timestamps | Windowed arrays | `{"window_type": "tumbling", "window_size": "5m"}` | High | 📋 |
| **Barrier** | Wait for N inputs | Multiple inputs | Combined when all arrive | `{"wait_for": "all", "timeout": "30s"}` | Medium | 📋 |
| **Debounce** | Emit only after stabilization | 1+ inputs | Stabilized value | `{"wait_time": "30s"}` | Medium | 📋 |
| **Gate** | Control flow with external signal | 1+ inputs | Same as input (if open) | `{"gate_id": "maintenance", "default_state": "open"}` | Medium | 📋 |
| **Until** | Loop until condition true | 1 input | Final value | `{"condition": "ready == true", "max_iterations": 100}` | Medium | 📋 |
| **Multiplex** | Fan-out to multiple outputs | 1 input | Multiple copies | `{"outputs": ["path1", "path2", "path3"]}` | Medium | 📋 |
| **Alert** | Trigger alerts | 1+ inputs | Alert status | `{"severity": "critical", "title": "High CPU"}` | High | 📋 |

### Output & Actions

| Node Type | Description | Inputs | Outputs | Example Config | Priority | Status |
|-----------|-------------|--------|---------|----------------|----------|--------|
| **HTTP POST** | Send HTTP POST request | `object` or `string` | Response object | `{"url": "https://webhook.site", "headers": {...}}` | High | 📋 |
| **Slack Alert** | Send Slack message | `string` or `object` | Status object | `{"webhook": "...", "channel": "#alerts"}` | High | 📋 |
| **Email** | Send email | `string` or `object` | Status object | `{"to": "ops@example.com", "subject": "Alert"}` | Medium | 📋 |
| **PagerDuty** | Create PagerDuty incident | `object` | Incident ID | `{"severity": "critical", "service": "web"}` | Medium | 📋 |
| **Write to Database** | Insert/update data | `object` or `array` | Row count | `{"table": "metrics", "operation": "insert"}` | Medium | 📋 |
| **S3 Writer** | Write to S3 | `string` or `binary` | Object URL | `{"bucket": "archives", "key": "data.json"}` | Low | 📋 |
| **Kafka Producer** | Publish to Kafka | `object` or `string` | Offset | `{"topic": "events", "partition": 0}` | Medium | 📋 |

### State & Memory

| Node Type | Description | Inputs | Outputs | Example Config | Priority | Status |
|-----------|-------------|--------|---------|----------------|----------|--------|
| **Variable** | Store/retrieve values | Value (for set) / None (for get) | Object with var_name and value | `{"var_name": "x", "var_op": "set"}` | High | ✅ |
| **Extract** | Extract fields from objects | Object | Object with extracted fields | `{"field": "name"}` or `{"fields": ["name", "age"]}` | High | ✅ |
| **Transform** | Transform data structures | Any | Varies by type | `{"transform_type": "to_array"}` | High | ✅ |
| **Accumulator** | Accumulate values | Values | Accumulated result | `{"accum_op": "sum"}` | High | ✅ |
| **Counter** | Increment/decrement counter | None | Current count | `{"counter_op": "increment", "delta": 1}` | High | ✅ |
| **Cache Get** | Retrieve from cache | Key | Cached value or null | `{"key": "last_value", "ttl": "5m"}` | Medium | 📋 |
| **Cache Set** | Store in cache | Value | Success status | `{"key": "last_value", "ttl": "5m"}` | Medium | 📋 |

---

## Looping Patterns for Observability Pipelines

### 1. For-Each Pattern (Parallel Processing)

**Use Case**: Process each host's metrics independently

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "prometheus_query",
      "data": {"query": "up", "groupBy": ["host"]}
    },
    {
      "id": "2",
      "type": "for_each",
      "data": {"max_concurrency": 5}
    },
    {
      "id": "3",
      "type": "statistical",
      "data": {"operations": ["mean", "max"]}
    },
    {
      "id": "4",
      "type": "join",
      "data": {"strategy": "all"}
    }
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3", "scope": "iteration"},
    {"source": "3", "target": "4", "scope": "iteration"}
  ]
}
```

**Benefits**:
- Process metrics from multiple hosts in parallel
- Independent failure handling per host
- Scalable for large infrastructure

### 2. Time Window Iteration Pattern

**Use Case**: Analyze metrics across multiple time windows

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "time_range_generator",
      "data": {"start": "now-24h", "window": "1h"}
    },
    {
      "id": "2",
      "type": "for_each"
    },
    {
      "id": "3",
      "type": "prometheus_query",
      "data": {"query": "rate(requests_total[5m])"}
    },
    {
      "id": "4",
      "type": "anomaly_detection",
      "data": {"method": "zscore", "threshold": 3}
    },
    {
      "id": "5",
      "type": "filter",
      "data": {"condition": "is_anomaly == true"}
    },
    {
      "id": "6",
      "type": "accumulator"
    }
  ]
}
```

**Benefits**:
- Analyze historical patterns
- Detect anomalies across time periods
- Build aggregate reports

### 3. Retry Loop with Backoff

**Use Case**: Resilient API calls with exponential backoff

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "http",
      "data": {"url": "https://api.unreliable.com/metrics"}
    },
    {
      "id": "2",
      "type": "if",
      "data": {"condition": "status_code >= 500"}
    },
    {
      "id": "3",
      "type": "delay",
      "data": {"duration_expr": "2^attempt * 1s"}
    },
    {
      "id": "4",
      "type": "while",
      "data": {"condition": "attempt < 5", "increment": "attempt"}
    }
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3", "condition": "true"},
    {"source": "3", "target": "4"},
    {"source": "4", "target": "1", "scope": "loop"}
  ]
}
```

**Benefits**:
- Automatic retry on transient failures
- Exponential backoff prevents overwhelming services
- Configurable max attempts

### 4. Consolidation/Aggregation Pattern

**Use Case**: Collect metrics from multiple sources and aggregate

```json
{
  "nodes": [
    {
      "id": "sources",
      "type": "parallel",
      "data": {
        "branches": [
          {"type": "prometheus_query", "query": "cpu_usage"},
          {"type": "prometheus_query", "query": "memory_usage"},
          {"type": "database_query", "sql": "SELECT avg(response_time)"}
        ]
      }
    },
    {
      "id": "normalize",
      "type": "for_each",
      "data": {"transform": "normalize_metrics"}
    },
    {
      "id": "join",
      "type": "join",
      "data": {"strategy": "all", "timeout": "10s"}
    },
    {
      "id": "aggregate",
      "type": "reduce",
      "data": {"operation": "merge"}
    },
    {
      "id": "alert_check",
      "type": "if",
      "data": {"condition": "any(values) > threshold"}
    },
    {
      "id": "alert",
      "type": "slack_alert"
    }
  ]
}
```

**Benefits**:
- Unified view of metrics from different sources
- Parallel data collection
- Consolidated alerting

### 5. Batch Processing with Chunking

**Use Case**: Process large log files in chunks

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "s3_reader",
      "data": {"bucket": "logs", "key": "access.log"}
    },
    {
      "id": "2",
      "type": "chunk",
      "data": {"size": 1000, "by": "lines"}
    },
    {
      "id": "3",
      "type": "for_each",
      "data": {"max_concurrency": 10}
    },
    {
      "id": "4",
      "type": "parse_log",
      "data": {"format": "nginx"}
    },
    {
      "id": "5",
      "type": "filter",
      "data": {"condition": "status >= 400"}
    },
    {
      "id": "6",
      "type": "accumulator",
      "data": {"operation": "concat"}
    },
    {
      "id": "7",
      "type": "group_by",
      "data": {"key": "status", "aggregation": "count"}
    }
  ]
}
```

**Benefits**:
- Memory-efficient processing of large files
- Parallel chunk processing
- Progressive results

---

## Types of Looping for Observability

### 1. **Iteration Loops** (For-Each)
- **Purpose**: Process each element in a collection
- **Use Cases**:
  - Process metrics from each host/pod/service
  - Analyze each log file in a directory
  - Query each time bucket independently
- **Key Features**:
  - Fixed number of iterations (array length)
  - Can be parallelized
  - Independent execution per item

### 2. **Conditional Loops** (While)
- **Purpose**: Loop until a condition is met
- **Use Cases**:
  - Retry failed API calls
  - Wait for system to reach steady state
  - Poll until data is available
- **Key Features**:
  - Dynamic iteration count
  - Must have max iteration limit
  - Condition evaluated each iteration

### 3. **Recursive Loops** (Self-referencing)
- **Purpose**: Traverse hierarchical structures
- **Use Cases**:
  - Process nested JSON structures
  - Follow service dependency chains
  - Navigate distributed traces
- **Key Features**:
  - Stack-based execution
  - Depth limits required
  - Memory considerations

### 4. **Time-based Loops** (Scheduled)
- **Purpose**: Execute at regular intervals
- **Use Cases**:
  - Periodic metric collection
  - Regular health checks
  - Scheduled reports
- **Key Features**:
  - Cron or interval-based
  - Runs indefinitely
  - External orchestration

### 5. **Fan-out/Fan-in** (Parallel + Join)
- **Purpose**: Split work, process in parallel, then combine
- **Use Cases**:
  - Query multiple data sources simultaneously
  - Parallel metric aggregation
  - Multi-region data collection
- **Key Features**:
  - Explicit parallelism
  - Join/merge semantics
  - Timeout handling

### 6. **Stream Processing** (Continuous)
- **Purpose**: Process infinite streams
- **Use Cases**:
  - Real-time log processing
  - Metric streaming pipelines
  - Event-driven workflows
- **Key Features**:
  - No defined end
  - Windowing semantics
  - Backpressure handling

---

## Proof of Concept: Looping Implementation

See the example implementation in `backend/examples/looping_poc.go` which demonstrates:

1. **Basic For-Each Loop**: Iterate over array and transform each element
2. **Conditional Loop with Retry**: Retry logic with max attempts
3. **Parallel Processing**: Fan-out work to multiple goroutines
4. **Aggregation/Consolidation**: Collect and merge results from multiple sources

### Running the POC

```bash
cd backend/examples
go run looping_poc.go
```

### POC Output Examples

```
=== For-Each Loop Example ===
Processing item: metric-1
Processing item: metric-2
Processing item: metric-3
Results: [METRIC-1, METRIC-2, METRIC-3]

=== Retry Loop Example ===
Attempt 1: Failed
Attempt 2: Failed  
Attempt 3: Success
Result: Data retrieved successfully

=== Parallel Processing Example ===
Querying source-1...
Querying source-2...
Querying source-3...
Combined results: 150 metrics

=== Consolidation Example ===
CPU: 45%, Memory: 78%, Disk: 23%
Alert: Memory usage high
```

---

## Implementation Roadmap

### Phase 1: Core Loop Support (High Priority)
- [ ] For-Each node with parallelism control
- [ ] If/Condition node
- [ ] Join/Merge node
- [ ] Basic error handling in loops

### Phase 2: Advanced Control Flow (Medium Priority)
- [ ] While loop with conditions
- [ ] Switch/case node
- [ ] Parallel execution node
- [ ] Timeout and cancellation

### Phase 3: Observability-Specific (High Priority)
- [ ] Prometheus query node
- [ ] Time bucketing node
- [ ] Statistical aggregation node
- [ ] Alert output nodes

### Phase 4: Advanced Features (Low Priority)
- [ ] Recursive loops with depth limits
- [ ] Stream processing support
- [ ] State management
- [ ] Circuit breaker patterns

---

## Design Considerations

### Loop Safety
- **Max Iterations**: All loops must have maximum iteration limits
- **Timeouts**: Loops should support timeouts to prevent infinite execution
- **Resource Limits**: Memory and CPU constraints for loop bodies
- **Circuit Breakers**: Fail fast when loops consistently fail

### Performance
- **Parallelism**: For-each loops should support concurrent execution
- **Lazy Evaluation**: Don't materialize entire arrays unnecessarily
- **Streaming**: Support streaming for large datasets
- **Backpressure**: Handle slow downstream consumers

### Debugging
- **Loop Metadata**: Track iteration count, elapsed time
- **Breakpoints**: Ability to pause on specific iterations
- **Logging**: Detailed logs for loop execution
- **Metrics**: Emit metrics about loop performance

---

## Contributing

To add a new node type:

1. Update this document with the node specification
2. Implement the node in `backend/workflow.go`
3. Add comprehensive tests in `backend/workflow_test.go`
4. Update examples in `backend/examples/`
5. Document usage patterns

For questions or suggestions, please open an issue.
