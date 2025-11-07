# Switch Node Examples

This directory contains comprehensive examples demonstrating the **Switch Node** in Thaiyyal workflows. The switch node enables multi-way branching based on input values or conditions, making it ideal for routing, categorization, and decision-making scenarios.

## Overview

The **Switch Node** is a powerful control flow primitive that:

- Routes execution to multiple different paths (vs binary true/false in Condition Node)
- Supports both **exact value matching** and **condition-based matching**
- Enables **conditional execution** where only the matched path executes
- Provides custom output handles for each case
- Includes a default fallback path for unmatched values

## Examples

### 01. HTTP Status Code Routing
**File**: `01-http-status-routing.json`

Routes HTTP responses to different handlers based on status code using **value matching**.

**Use Case**: API error handling, retry logic, success/failure processing

**Key Features**:
- Exact value matching (200, 201, 404, 500)
- Multiple output paths for different handlers
- Default path for unknown status codes
- Demonstrates real-world HTTP response handling

**Pattern**: Discrete value routing

**Try it**:
- Change `status_code` value to 200, 404, or 500 to see different handlers execute
- Try an unknown status like 418 to trigger the default handler

---

### 02. User Role-Based Routing
**File**: `02-user-role-routing.json`

Routes users to appropriate interfaces based on their role using **string value matching**.

**Use Case**: Authentication, authorization, multi-tenant systems, access control

**Key Features**:
- String value matching ("admin", "moderator", "user")
- Role-based interface selection
- Guest fallback for unauthenticated users
- Demonstrates RBAC (Role-Based Access Control)

**Pattern**: String-based categorization

**Try it**:
- Change `user_role` to "admin", "moderator", "user", or "guest"
- See how each role routes to a different interface

---

### 03. Priority Queue Routing
**File**: `03-priority-queue-routing.json`

Routes tasks to different processing queues based on priority level using **range conditions**.

**Use Case**: Task scheduling, SLA management, workload distribution, queue systems

**Key Features**:
- Range-based matching (>=9, >=7, >=4, >=1)
- Priority levels (Critical, High, Medium, Low)
- Different SLA guarantees per queue
- Invalid priority detection

**Pattern**: Range-based routing with ordered cases

**Try it**:
- Set `task_priority` to values between 1-10
- Try invalid values like 0 or 15 to see default handling
- Note: Cases are evaluated in order, first match wins

---

### 04. Content-Type Routing
**File**: `04-content-type-routing.json`

Routes data to appropriate parsers based on content type using **exact value matching**.

**Use Case**: Data ingestion, API routing, file processing, content negotiation

**Key Features**:
- Content-type header matching
- Multiple parser types (JSON, XML, CSV, Text)
- Raw/binary fallback handler
- Demonstrates data pipeline routing

**Pattern**: String value routing with fallback

**Try it**:
- Change `content_type` to different MIME types
- See appropriate parser selected for each type
- Try unknown types to see raw handler activation

---

### 05. Age Category Routing
**File**: `05-age-category-routing.json`

Categorizes and routes based on age groups using **range conditions**.

**Use Case**: Demographics, age restrictions, pricing tiers, content filtering

**Key Features**:
- Age range conditions (<18, <65, >=65)
- Three distinct categories (Minor, Adult, Senior)
- Different processing per age group
- Age validation

**Pattern**: Boundary-based categorization

**Try it**:
- Set `age_input` to different values: 10, 35, 70
- Observe category boundaries at 18 and 65
- Try edge cases: exactly 18, exactly 65

---

## Switch Node Configuration

### Basic Structure

```json
{
  "type": "switch",
  "data": {
    "label": "Descriptive Label",
    "cases": [
      {
        "when": "Case description",
        "value": <optional-exact-value>,
        "output_path": "output_handle_name"
      }
    ],
    "default_path": "default_handle_name"
  }
}
```

### Matching Strategies

#### 1. Exact Value Matching

Use when you have specific, discrete values to match:

```json
{
  "when": "HTTP 200 OK",
  "value": 200,
  "output_path": "success"
}
```

**Supports**:
- Numbers: `200`, `404`, `500`
- Strings: `"admin"`, `"JSON"`, `"active"`
- Booleans: `true`, `false`

#### 2. Condition Matching

Use when you need ranges or comparisons:

```json
{
  "when": "High Priority",
  "output_path": "high_queue"
}
```

**Note**: When `value` is omitted, the `when` field is evaluated as a condition.

**Supported operators**:
- `>` Greater than
- `<` Less than
- `>=` Greater than or equal
- `<=` Less than or equal
- `==` Equality
- `!=` Inequality

### Important: Case Order Matters!

Cases are evaluated **in order**, and the **first match wins**:

❌ **Bad** (later cases never match):
```json
{
  "cases": [
    {"when": ">0", "output_path": "positive"},
    {"when": ">100", "output_path": "large"}  // Never reached!
  ]
}
```

✅ **Good** (most specific first):
```json
{
  "cases": [
    {"when": ">100", "output_path": "large"},
    {"when": ">0", "output_path": "positive"}
  ]
}
```

## Conditional Execution

Switch nodes support **conditional execution** (path termination), meaning only the nodes connected to the matched output path will execute.

### How It Works

1. Switch node evaluates input against all cases
2. First matching case determines the `output_path`
3. Engine only executes nodes connected via that `sourceHandle`
4. All other paths are **skipped entirely**

### Example

```json
{
  "edges": [
    {"source": "switch", "sourceHandle": "success", "target": "success_handler"},
    {"source": "switch", "sourceHandle": "error", "target": "error_handler"}
  ]
}
```

If switch matches "success", only `success_handler` executes. The `error_handler` is **not executed**.

## Common Patterns

### Pattern 1: Multi-Way API Routing

**Scenario**: Route API responses to different handlers

**Implementation**:
- Input: HTTP status code or API response type
- Switch: Match status codes to handlers
- Outputs: Success, retry, error, fallback paths

**Example**: `01-http-status-routing.json`

### Pattern 2: Role-Based Access Control (RBAC)

**Scenario**: Grant access based on user roles

**Implementation**:
- Input: User role string
- Switch: Match role to interface
- Outputs: Admin panel, user dashboard, guest view

**Example**: `02-user-role-routing.json`

### Pattern 3: Priority/SLA Routing

**Scenario**: Route tasks by priority level

**Implementation**:
- Input: Priority number (1-10)
- Switch: Match ranges to queues
- Outputs: Critical, high, medium, low queues

**Example**: `03-priority-queue-routing.json`

### Pattern 4: Content-Type Negotiation

**Scenario**: Select parser based on content type

**Implementation**:
- Input: Content-Type header string
- Switch: Match MIME type to parser
- Outputs: JSON, XML, CSV, text, raw parsers

**Example**: `04-content-type-routing.json`

### Pattern 5: Demographic Segmentation

**Scenario**: Categorize by age, income, or other ranges

**Implementation**:
- Input: Age/value number
- Switch: Match ranges to categories
- Outputs: Minor, adult, senior categories

**Example**: `05-age-category-routing.json`

## Best Practices

### 1. Order Cases from Specific to General

✅ Most specific conditions first
```json
[
  {"when": ">1000", "output_path": "huge"},
  {"when": ">100", "output_path": "large"},
  {"when": ">0", "output_path": "positive"}
]
```

### 2. Always Provide a Default Path

✅ Handle unexpected/invalid inputs
```json
{
  "cases": [...],
  "default_path": "error_handler"
}
```

### 3. Use Descriptive When Labels

✅ Self-documenting:
```json
{"when": "Premium Customer (Tier 1)", "value": "premium", "output_path": "vip_service"}
```

❌ Unclear:
```json
{"when": "T1", "value": "premium", "output_path": "vip"}
```

### 4. Use Value Matching for Discrete Values

✅ Exact matching when possible:
```json
{"when": "HTTP 404", "value": 404, "output_path": "not_found"}
```

### 5. Use Meaningful Output Path Names

✅ Descriptive:
```json
{"output_path": "high_priority_queue"}
```

❌ Generic:
```json
{"output_path": "path2"}
```

### 6. Test All Paths

Test scenarios:
- ✅ First case matches
- ✅ Last case matches
- ✅ Middle cases match
- ✅ Default path (no match)
- ✅ Boundary values

## Common Pitfalls

### Pitfall 1: Incorrect Case Order

❌ **Problem**: `>50` matches before `>100` is ever checked
```json
[
  {"when": ">50", "output_path": "high"},
  {"when": ">100", "output_path": "very_high"}  // Unreachable!
]
```

### Pitfall 2: Missing Default Path

❌ **Problem**: Unhandled inputs cause silent failures

✅ **Solution**: Always include `default_path`

### Pitfall 3: Type Mismatches

❌ **Problem**: String "200" doesn't match number 200
```json
{"value": "200"}  // String
// Input: 200 (number) → No match!
```

✅ **Solution**: Ensure types align

### Pitfall 4: Complex Conditions

❌ **Not Supported**: AND/OR logic in When field
```json
{"when": ">=50 && <=100"}  // Not supported!
```

✅ **Solution**: Use multiple cases or Condition node before switch

## Comparison: Switch vs Condition Node

| Feature | Switch Node | Condition Node |
|---------|-------------|----------------|
| **Output Paths** | Multiple (unlimited) | 2 (true/false) |
| **Use Case** | Multi-way routing | Binary decisions |
| **Complexity** | More complex | Simpler |
| **When to Use** | 3+ possible outcomes | Yes/no decisions |
| **Matching** | Value or condition | Expression only |

### When to Use Switch

✅ Multiple discrete values (status codes, roles, types)  
✅ 3+ possible outcomes  
✅ Range-based categorization  
✅ Cleaner than nested conditions  

### When to Use Condition

✅ Simple yes/no decision  
✅ Binary branching  
✅ Boolean logic  
✅ Single condition check  

## Testing Your Workflows

### Using the API

```bash
# Load workflow
curl -X POST http://localhost:8080/api/workflow/execute \
  -H "Content-Type: application/json" \
  -d @examples/switch-node/01-http-status-routing.json
```

### Testing Different Values

Modify the input value in the JSON file:

```json
{
  "id": "status_code",
  "type": "number",
  "data": {
    "value": 404  // Change this to test different cases
  }
}
```

### Verifying Conditional Execution

Check the response - only nodes in the matched path should appear in `node_results`:

```json
{
  "execution_id": "...",
  "node_results": {
    "status_code": 404,
    "status_router": {"matched": true, "output_path": "not_found"},
    "not_found_handler": "⚠️ Not found - retry..."
    // Other handlers NOT present - they were skipped!
  }
}
```

## Real-World Scenarios

### API Gateway Routing

```
[HTTP Request] → [Extract Status] → [Switch] ─┬─[200]──→ [Success Processing]
                                                ├─[404]──→ [Retry with Cache]
                                                ├─[429]──→ [Rate Limit Handler]
                                                └─[500]──→ [Error Alerting]
```

### Multi-Tenant Systems

```
[User Login] → [Get Tenant] → [Switch] ─┬─[enterprise]──→ [Premium Features]
                                         ├─[pro]─────────→ [Pro Features]
                                         └─[free]────────→ [Basic Features]
```

### Data Pipeline

```
[Incoming Data] → [Detect Type] → [Switch] ─┬─[json]──→ [JSON Validator]
                                             ├─[xml]───→ [XML Parser]
                                             ├─[csv]───→ [CSV Importer]
                                             └─[other]─→ [Raw Storage]
```

## Additional Resources

- **Main Documentation**: [`docs/SWITCH_NODE_ANALYSIS.md`](../../docs/SWITCH_NODE_ANALYSIS.md)
- **Integration Tests**: [`backend/pkg/engine/switch_node_scenarios_test.go`](../../backend/pkg/engine/switch_node_scenarios_test.go)
- **Unit Tests**: [`backend/pkg/executor/control_switch_test.go`](../../backend/pkg/executor/control_switch_test.go)
- **Conditional Execution Docs**: [`docs/CONDITIONAL_EXECUTION_IMPLEMENTATION.md`](../../docs/CONDITIONAL_EXECUTION_IMPLEMENTATION.md)
- **Conditional Branching Examples**: [`examples/conditional-branching/`](../conditional-branching/)

## Troubleshooting

### Issue: No cases matching

**Symptom**: Default path always executing

**Solutions**:
1. Check input value type matches case value type
2. Verify case conditions are correct
3. Test with simple exact value match first

### Issue: Wrong case matching

**Symptom**: Unexpected path taken

**Solutions**:
1. Check case order (first match wins)
2. Review condition operators (>=, >, <, etc.)
3. Add logging to see what value is being matched

### Issue: Multiple handlers executing

**Symptom**: Both matched and unmatched paths executing

**Solutions**:
1. Verify edges have `sourceHandle` set
2. Check edge connections in workflow JSON
3. Ensure conditional execution is properly configured

## Contributing

To add more examples:

1. Create a new JSON file following the naming pattern: `##-descriptive-name.json`
2. Include descriptive `name` and `description` fields
3. Use meaningful node IDs and labels
4. Add appropriate `sourceHandle` values to edges
5. Update this README with the new example
6. Test the workflow to ensure it works correctly

## Support

For issues or questions:
- Open an issue on GitHub
- Check existing documentation
- Review integration tests for usage examples

---

**Last Updated**: 2025-11-07  
**Examples**: 5 workflows  
**Node Type**: `switch`  
**Status**: Production Ready
