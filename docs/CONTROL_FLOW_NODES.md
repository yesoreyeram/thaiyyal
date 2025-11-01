# Control Flow Nodes - Complete Reference

**Document Version:** 1.0  
**Last Updated:** 2025-11-01  
**Status:** Production Ready

---

## Table of Contents

1. [Overview](#overview)
2. [Implementation Status Summary](#implementation-status-summary)
3. [Available Control Flow Nodes](#available-control-flow-nodes)
4. [Condition Node](#1-condition-node)
5. [Filter Node](#2-filter-node)
6. [Map Node](#3-map-node)
7. [Reduce Node](#4-reduce-reduce)
8. [ForEach Node](#5-foreach-node)
9. [WhileLoop Node](#6-whileloop-node)
10. [Switch Node](#7-switch-node)
11. [Recommended New Node Types](#recommended-new-node-types)
12. [Expression Language Reference](#expression-language-reference)
13. [Design Patterns](#design-patterns)
14. [Troubleshooting](#troubleshooting)

---

## Overview

Control flow nodes in Thaiyyal enable conditional execution, iteration, and data transformation within workflows. These nodes follow the **Single Responsibility Principle** and are designed to be **composable** - you can combine them in flexible ways to create sophisticated data processing pipelines.

### Design Philosophy

✅ **Composability** - Mix and match primitives freely  
✅ **Single Responsibility** - Each node does ONE thing well  
✅ **Reusability** - Nodes work standalone or in combination  
✅ **Clarity** - Obvious what each node does  
✅ **Safety** - Built-in limits prevent resource exhaustion

---

## Implementation Status Summary

### ✅ Fully Implemented (Production Ready)

**1. Condition Node** - 100% Complete
- ✅ Backend executor implemented
- ✅ Frontend component with title editing
- ✅ Full expression support (variables, context, node refs)
- ✅ Comprehensive test coverage (11 test cases)
- ✅ Documentation complete
- ✅ Error handling robust

**2. Filter Node** - 100% Complete
- ✅ Backend executor implemented
- ✅ Frontend component (FilterNode.tsx)
- ✅ `item.field` syntax fully working
- ✅ Expression evaluation with variables/context
- ✅ Comprehensive test coverage (13 test suites)
- ✅ Non-array input handling
- ✅ Documentation complete

**3. Switch Node** - 100% Complete
- ✅ Backend executor implemented
- ✅ Value matching mode working
- ✅ Condition matching mode working
- ✅ Default path support
- ✅ Test coverage (7 test cases)
- ✅ Documentation complete

### ⚠️ Partially Implemented (Core Working, Enhancement Needed)

**4. Map Node** - 75% Complete
- ✅ Backend executor implemented
- ✅ Field extraction mode fully working
- ✅ Validation and error handling
- ✅ Test coverage for field extraction
- ⚠️ Expression transformation mode needs enhancement
- ⚠️ Variable resolution in arithmetic expressions (TODO)
- ❌ Frontend component not yet implemented
- ✅ Documentation complete

**5. Reduce Node** - 70% Complete
- ✅ Backend executor structure implemented
- ✅ Initial value configuration working
- ✅ Validation complete
- ✅ Test coverage for structure
- ⚠️ Expression evaluation with accumulator needs enhancement
- ⚠️ Ternary operator support needed
- ⚠️ Variable resolution in arithmetic (TODO)
- ❌ Frontend component not yet implemented
- ✅ Documentation complete

**6. ForEach Node** - 70% Complete
- ✅ Backend executor implemented
- ✅ Iteration context prepared (item/index/items)
- ✅ Max iteration safety limits
- ✅ Error handling (continue on error)
- ✅ Test coverage for iteration logic
- ⚠️ Child node execution requires engine integration (TODO)
- ⚠️ Variable injection handled by engine (not executor)
- ❌ Frontend component not yet implemented
- ✅ Documentation complete

**7. WhileLoop Node** - 60% Complete
- ✅ Backend executor structure implemented
- ✅ Condition evaluation working
- ✅ Max iteration safety
- ✅ Test coverage (4 test cases)
- ⚠️ Value update between iterations requires engine integration
- ⚠️ Child node execution not yet integrated
- ❌ Frontend component basic, needs enhancement
- ✅ Documentation complete

### 📊 Overall Status

| Component | Backend | Frontend | Tests | Docs | Status |
|-----------|---------|----------|-------|------|--------|
| Condition | ✅ 100% | ✅ 100% | ✅ 11 cases | ✅ Complete | 🟢 Production |
| Filter | ✅ 100% | ✅ 100% | ✅ 13 suites | ✅ Complete | 🟢 Production |
| Map | ✅ 85% | ❌ 0% | ✅ Field tests | ✅ Complete | 🟡 Beta |
| Reduce | ✅ 70% | ❌ 0% | ✅ Basic tests | ✅ Complete | 🟡 Beta |
| ForEach | ✅ 70% | ❌ 0% | ✅ 7 cases | ✅ Complete | 🟡 Beta |
| WhileLoop | ✅ 60% | 🟡 50% | ✅ 4 cases | ✅ Complete | 🟡 Beta |
| Switch | ✅ 100% | 🟡 50% | ✅ 7 cases | ✅ Complete | 🟢 Production |

### 🔧 Key Enhancement Needed

**Expression System Enhancement** (affects Map, Reduce):
- Variable resolution in arithmetic contexts (`item * 2`, `accumulator + item`)
- Ternary operator support (`a > b ? a : b`)
- Complex object operations in expressions
- Current workaround: Use dedicated nodes (Math, Transform) after Map/Reduce

**Workflow Engine Integration** (affects ForEach, WhileLoop):
- Child node execution within loops
- Variable injection to child nodes
- Loop state management
- Current workaround: Basic iteration works, full loop execution pending

---

## Available Control Flow Nodes

| Node Type | Purpose | Input Type | Output Type | Branch Count |
|-----------|---------|------------|-------------|--------------|
| **Condition** | Evaluate boolean expression and route data | Any | Object with metadata | 2 (true/false) |
| **Filter** | Filter array elements by condition | Array | Array (filtered) | 1 |
| **Map** | Transform array elements | Array | Array (transformed) | 1 |
| **Reduce** | Aggregate array to single value | Array | Single value | 1 |
| **ForEach** | Iterate over array with side effects | Array | Metadata | 1 |
| **WhileLoop** | Loop while condition is true | Any | Object with metadata | 1 |
| **Switch** | Multi-way branching | Any | Object with metadata | N (cases + default) |

---

## 1. Condition Node

**Implementation Status:** 🟢 **Production Ready** (100% Complete)

### Description

The Condition node evaluates a boolean expression and passes through the input value along with metadata about whether the condition was met. It's used for **conditional branching** in workflows.

### Current Implementation Status

✅ **Backend:** Fully implemented (`backend/pkg/executor/condition.go`)
- Expression evaluation with fallback
- Variable, context, and node reference support
- Comprehensive error handling
- Test coverage: 11 test cases in `condition_test.go`

✅ **Frontend:** Fully implemented (`src/components/nodes/ConditionNode.tsx`)
- React component with title editing
- Condition input field
- Context menu integration
- Node info documentation

✅ **Features:** All planned features complete
- Advanced expression support
- Boolean logic (AND, OR, NOT)
- Metadata output with path information
- Production-ready error handling

### Node Type

`condition`

### Responsibilities

- Evaluate a boolean expression against input data
- Pass through the original input value unchanged
- Add metadata indicating which branch (true/false) should be taken
- Support advanced expressions with variables, context, and node references

### Input Types

**Accepts:** Any type (number, string, object, array, boolean, null)

### Output Format

```json
{
  "value": <original_input>,
  "condition_met": true|false,
  "condition": "<expression_string>",
  "path": "true"|"false",
  "true_path": true|false,
  "false_path": true|false
}
```

### Node Properties

```typescript
{
  "type": "condition",
  "data": {
    "condition": string  // Required: Boolean expression to evaluate
  }
}
```

**Property Details:**

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `condition` | string | ✅ Yes | N/A | Boolean expression using expression language |

### Expression Syntax

The condition supports the full expression language:

**Simple comparisons:**
```javascript
value > 10
value == "active"
value >= 18 && value <= 65
```

**Variable references:**
```javascript
value > variables.threshold
age >= variables.minAge
```

**Context variables:**
```javascript
score > context.passingScore
status == context.expectedStatus
```

**Node references:**
```javascript
value > node.sensor1.value
currentTemp >= node.tempSensor.reading
```

**Complex boolean logic:**
```javascript
(age >= 18 && status == "active") || admin == true
temperature > 100 && pressure < 50
```

### Error Scenarios

| Scenario | Behavior | Error Message |
|----------|----------|---------------|
| Missing condition | ❌ Execution fails | "condition node missing condition" |
| No input | ❌ Execution fails | "condition node needs at least 1 input" |
| Expression evaluation error | ⚠️ Falls back to simple evaluation | Logs warning, continues execution |
| Invalid expression syntax | ⚠️ Returns false | Logs error, treats as condition not met |

### Branches

**Branch 1: True Path**
- Taken when `condition_met == true`
- Output field: `true_path: true`

**Branch 2: False Path**
- Taken when `condition_met == false`
- Output field: `false_path: true`

### Example Workflows

#### Example 1: Age Verification
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {"value": 25}
    },
    {
      "id": "2",
      "type": "condition",
      "data": {"condition": "value >= 18"}
    },
    {
      "id": "3",
      "type": "text_input",
      "data": {"value": "Adult content"}
    },
    {
      "id": "4",
      "type": "text_input",
      "data": {"value": "Age restricted"}
    }
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3", "sourceHandle": "true"},
    {"source": "2", "target": "4", "sourceHandle": "false"}
  ]
}
```
**Result:** Takes true path → "Adult content"

#### Example 2: Variable-Based Threshold
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "variable",
      "data": {"name": "threshold", "value": 100}
    },
    {
      "id": "2",
      "type": "number",
      "data": {"value": 150}
    },
    {
      "id": "3",
      "type": "condition",
      "data": {"condition": "value > variables.threshold"}
    }
  ]
}
```
**Result:** `condition_met: true` (150 > 100)

#### Example 3: Complex Business Logic
```json
{
  "data": {
    "condition": "(age >= 18 && age <= 65) && (status == \"active\" || premium == true)"
  }
}
```

#### Example 4: Node Reference Comparison
```json
{
  "data": {
    "condition": "value > node.sensor1.value && value < node.sensor2.value"
  }
}
```

#### Example 5: Temperature Range Check
```json
{
  "data": {
    "condition": "temperature >= 20 && temperature <= 25"
  }
}
```

#### Example 6: Status Validation
```json
{
  "data": {
    "condition": "status == \"approved\" && balance > 0"
  }
}
```

#### Example 7: Multi-Factor Authentication
```json
{
  "data": {
    "condition": "passwordCorrect == true && (twoFactorEnabled == false || twoFactorVerified == true)"
  }
}
```

#### Example 8: Inventory Check
```json
{
  "data": {
    "condition": "quantity > 0 && price <= variables.budget && inStock == true"
  }
}
```

### Limitations

⚠️ **Current Limitations:**
- Expression fallback to simple evaluation if advanced evaluation fails
- Complex object comparisons may not work as expected
- Date/time comparisons require parseDate() function

### TODOs

- [ ] Add support for ternary operator (a ? b : c)
- [ ] Improve error messages for invalid expressions
- [ ] Add expression validation before execution
- [ ] Support for regular expression matching

### Related Nodes

- **Switch Node** - For multi-way branching (more than 2 paths)
- **Filter Node** - For filtering arrays with similar expression logic
- **WhileLoop Node** - For condition-based iteration

---

## 2. Filter Node

**Implementation Status:** 🟢 **Production Ready** (100% Complete)

### Description

The Filter node filters JSON array elements based on an expression condition. Elements where the expression evaluates to `true` are included in the output. Uses intuitive `item.field` syntax where `item` references the current array element.

### Current Implementation Status

✅ **Backend:** Fully implemented (`backend/pkg/executor/filter.go`)
- `item.field` syntax working perfectly
- Expression evaluation with variables/context/node refs
- Non-array input graceful handling
- Comprehensive logging and metadata
- Test coverage: 13 test suites in `filter_test.go`

✅ **Frontend:** Fully implemented (`src/components/nodes/FilterNode.tsx`)
- Purple gradient styling
- Condition input field with `item.field` placeholder
- Title editing support
- Full context menu integration

✅ **Features:** All planned features complete
- Expression engine integration
- Error handling (continue on error)
- Statistics tracking (input/output/skipped/error counts)
- Production-ready

### Node Type

`filter`

### Responsibilities

- Filter array elements using boolean expression
- Provide access to current element via `item` variable
- Handle non-array inputs gracefully (pass-through with warning)
- Track filtering statistics (input/output/skipped/error counts)
- Support full expression language with variables and context

### Input Types

**Primary:** Array of objects or primitives  
**Fallback:** Any type (passed through with warning)

### Output Format

```json
{
  "filtered": [...],         // Filtered array
  "input_count": number,     // Original array length
  "output_count": number,    // Filtered array length
  "skipped_count": number,   // Elements filtered out
  "error_count": number,     // Evaluation errors
  "condition": "string",     // Expression used
  "is_array": true|false     // Input was array
}
```

**Non-array input output:**
```json
{
  "input": <original_value>,
  "filtered": <original_value>,
  "is_array": false,
  "warning": "input is not an array, passed through unchanged",
  "original_type": "string|number|object"
}
```

### Node Properties

```typescript
{
  "type": "filter",
  "data": {
    "condition": string  // Required: Boolean expression with item syntax
  }
}
```

**Property Details:**

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `condition` | string | ✅ Yes | N/A | Expression to filter elements (use `item` variable) |

### Expression Syntax

The filter node uses the `item` variable to reference the current array element:

**Simple field comparison:**
```javascript
item.age > 18
item.status == "active"
item.price <= 100
```

**Nested field access:**
```javascript
item.profile.verified == true
item.address.city == "New York"
item.settings.notifications.email == true
```

**Primitive arrays:**
```javascript
item > 10           // For array of numbers: [5, 15, 20]
item == "apple"     // For array of strings
```

**Complex conditions:**
```javascript
item.age >= 18 && item.active == true
item.score > 80 || item.bonus == true
item.price >= 10 && item.price <= 100
```

**With variables:**
```javascript
item.age >= variables.minAge
item.score > variables.threshold
item.price <= variables.budget
```

**With context:**
```javascript
item.department == context.currentDepartment
item.priority >= context.minPriority
```

### Error Scenarios

| Scenario | Behavior | Error Message |
|----------|----------|---------------|
| Missing condition | ❌ Execution fails | "filter node missing condition expression" |
| Empty condition | ❌ Execution fails | "filter node missing condition expression" |
| No input | ❌ Execution fails | "filter node needs at least 1 input" |
| Non-array input | ⚠️ Passes through with warning | Logs warning, returns original input |
| Expression error | ⚠️ Skips element, continues | Increments error_count, logs debug message |
| Missing field | ⚠️ Skips element | Element skipped if field doesn't exist |

### Branches

**Single Output Branch:**
- Contains filtered array
- Access via `filtered` field in output object

### Example Workflows

#### Example 1: Filter Adults
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {
        "value": [
          {"name": "Alice", "age": 25},
          {"name": "Bob", "age": 15},
          {"name": "Charlie", "age": 30}
        ]
      }
    },
    {
      "id": "2",
      "type": "filter",
      "data": {"condition": "item.age >= 18"}
    }
  ]
}
```
**Output:**
```json
{
  "filtered": [
    {"name": "Alice", "age": 25},
    {"name": "Charlie", "age": 30}
  ],
  "input_count": 3,
  "output_count": 2,
  "skipped_count": 1
}
```

#### Example 2: Filter Active Users
```json
{
  "data": {
    "condition": "item.active == true && item.verified == true"
  }
}
```
**Input:** `[{active:true, verified:true}, {active:true, verified:false}]`  
**Output:** `[{active:true, verified:true}]`

#### Example 3: Price Range Filter
```json
{
  "data": {
    "condition": "item.price >= 10 && item.price <= 100"
  }
}
```
**Input:** `[{price:5}, {price:50}, {price:150}]`  
**Output:** `[{price:50}]`

#### Example 4: Filter with Variable Threshold
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "variable",
      "data": {"name": "minAge", "value": 21}
    },
    {
      "id": "2",
      "type": "number",
      "data": {
        "value": [
          {"name": "Alice", "age": 25},
          {"name": "Bob", "age": 18}
        ]
      }
    },
    {
      "id": "3",
      "type": "filter",
      "data": {"condition": "item.age >= variables.minAge"}
    }
  ]
}
```
**Output:** `[{"name": "Alice", "age": 25}]`

#### Example 5: Nested Field Filter
```json
{
  "data": {
    "condition": "item.profile.verified == true && item.profile.premium == true"
  }
}
```

#### Example 6: Filter Primitive Array
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {"value": [1, 5, 10, 15, 20]}
    },
    {
      "id": "2",
      "type": "filter",
      "data": {"condition": "item > 10"}
    }
  ]
}
```
**Output:** `{"filtered": [15, 20], "output_count": 2}`

#### Example 7: Complex Business Logic
```json
{
  "data": {
    "condition": "(item.age >= 18 && item.age <= 65) && (item.employed == true || item.student == true)"
  }
}
```

#### Example 8: Multi-Criteria Product Filter
```json
{
  "data": {
    "condition": "item.inStock == true && item.price <= variables.budget && item.rating >= 4.0"
  }
}
```

### Limitations

⚠️ **Current Limitations:**
- `item` must be used to reference array elements
- Direct field access without `item.` may not work reliably
- Expression errors skip elements rather than failing the entire operation
- No support for filtering with child node results (use Map + Filter composition)

### TODOs

- [ ] Add support for filter expressions that return array indices
- [ ] Add option to collect error items separately
- [ ] Support for negative filtering (inverse condition)
- [ ] Add support for filtering by index (`item` + `index` variable)

### Related Nodes

- **Map Node** - Transform filtered results
- **Reduce Node** - Aggregate filtered results
- **Condition Node** - Similar expression syntax for branching
- **ForEach Node** - Iterate over filtered results

---

## 3. Map Node

**Implementation Status:** 🟡 **Beta** (75% Complete - Field Extraction Working)

### Description

The Map node transforms each element of an array into a new value. It supports two modes: **field extraction** (extract a specific property from objects) and **expression transformation** (apply an expression to each element).

### Current Implementation Status

✅ **Backend:** Core implemented (`backend/pkg/executor/map.go`)
- ✅ Field extraction mode fully working
- ✅ Error handling and validation complete
- ✅ Statistics tracking (successful/failed counts)
- ⚠️ Expression transformation mode needs enhancement
- ⚠️ Variable resolution in arithmetic expressions (TODO)
- Test coverage: Field extraction tests in `map_test.go`

❌ **Frontend:** Not yet implemented
- Need MapNode.tsx component
- Should match Filter node styling
- Configuration for field vs expression mode

⚠️ **Limitations:**
- Expression-based transformations (`item * 2`) require expression system enhancement
- Only top-level field extraction (no nested paths yet)
- Workaround: Use Field extraction + Math node for transformations

### Node Type

`map`

### Responsibilities

- Transform each array element to a new value
- Extract fields from objects
- Apply expressions to elements
- Provide access to `item`, `index`, `items` variables
- Handle transformation errors gracefully
- Track success/failure statistics

### Input Types

**Required:** Array of any type (objects, primitives, mixed)

### Output Format

```json
{
  "results": [...],          // Transformed array
  "input_count": number,     // Original array length
  "output_count": number,    // Result array length (same as input)
  "successful": number,      // Successful transformations
  "failed": number           // Failed transformations (result is null)
}
```

### Node Properties

```typescript
{
  "type": "map",
  "data": {
    "field": string,        // Field to extract (mode 1)
    "expression": string    // Expression to evaluate (mode 2)
  }
}
```

**Property Details:**

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `field` | string | Conditional | N/A | Field name to extract from each object |
| `expression` | string | Conditional | N/A | Expression to transform each element |

**⚠️ Constraint:** Exactly ONE of `field` or `expression` must be specified, not both.

### Expression Syntax

#### Field Extraction Mode
```json
{
  "field": "name"
}
```
**Input:** `[{name:"Alice", age:25}, {name:"Bob", age:30}]`  
**Output:** `["Alice", "Bob"]`

#### Expression Mode
```javascript
item * 2                    // Double each number
item.age * 1.1              // Increase age by 10%
item.price * 0.9            // Apply 10% discount
```

**Variables available in expressions:**
- `item` - Current array element
- `index` - Current index (0-based)
- `items` - Full input array

### Error Scenarios

| Scenario | Behavior | Error Message |
|----------|----------|---------------|
| No input | ❌ Execution fails | "map node needs at least 1 input" |
| Non-array input | ❌ Execution fails | "map node requires array input, got <type>" |
| No field or expression | ❌ Validation fails | "map node requires either 'expression' or 'field'" |
| Both field and expression | ❌ Validation fails | "map node cannot have both 'expression' and 'field'" |
| Field not in object | ⚠️ Sets null, continues | Increments failed count, logs debug |
| Expression error | ⚠️ Sets null, continues | Increments failed count, logs debug |
| Non-object for field extraction | ⚠️ Sets null, continues | "cannot extract field from non-object" |

### Branches

**Single Output Branch:**
- Contains `results` array with transformed values
- Failed transformations result in `null` at that index

### Example Workflows

#### Example 1: Extract Names
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {
        "value": [
          {"name": "Alice", "age": 25},
          {"name": "Bob", "age": 30}
        ]
      }
    },
    {
      "id": "2",
      "type": "map",
      "data": {"field": "name"}
    }
  ]
}
```
**Output:**
```json
{
  "results": ["Alice", "Bob"],
  "successful": 2,
  "failed": 0
}
```

#### Example 2: Extract Nested Field
```json
{
  "data": {"field": "profile.email"}
}
```
**Note:** Currently only supports top-level fields. Nested field support is TODO.

#### Example 3: Double Numbers (TODO - requires expression enhancement)
```json
{
  "data": {"expression": "item * 2"}
}
```
**Input:** `[1, 2, 3]`  
**Expected Output:** `[2, 4, 6]`  
**Status:** ⚠️ Requires expression enhancement

#### Example 4: Calculate Discount
```json
{
  "data": {"expression": "item.price * 0.9"}
}
```
**Status:** ⚠️ Requires expression enhancement for arithmetic

#### Example 5: Filter + Map Composition
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {
        "value": [
          {"name": "Alice", "age": 25, "active": true},
          {"name": "Bob", "age": 15, "active": true},
          {"name": "Charlie", "age": 30, "active": false}
        ]
      }
    },
    {
      "id": "2",
      "type": "filter",
      "data": {"condition": "item.age >= 18 && item.active == true"}
    },
    {
      "id": "3",
      "type": "map",
      "data": {"field": "name"}
    }
  ]
}
```
**Output:** `["Alice"]` (adult active users' names)

#### Example 6: Extract IDs
```json
{
  "data": {"field": "id"}
}
```
**Input:** `[{id:1, name:"A"}, {id:2, name:"B"}]`  
**Output:** `[1, 2]`

#### Example 7: Extract Email Addresses
```json
{
  "data": {"field": "email"}
}
```

#### Example 8: Multi-Step Transformation
```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": [/*users*/]}},
    {"id": "2", "type": "filter", "data": {"condition": "item.verified == true"}},
    {"id": "3", "type": "map", "data": {"field": "email"}},
    {"id": "4", "type": "visualization", "data": {}}
  ]
}
```
**Flow:** Filter verified users → Extract emails → Display

### Limitations

⚠️ **Current Limitations:**
- Expression-based transformations require expression system enhancement
- Only supports top-level field extraction (no nested paths like `user.profile.email`)
- Failed transformations result in `null`, not omitted
- No support for conditional mapping (use Filter + Map)
- Cannot map with child node execution

### TODOs

- [ ] Add support for nested field paths (`user.profile.email`)
- [ ] Implement expression-based transformations with variable resolution
- [ ] Add option to omit failed transformations instead of `null`
- [ ] Add support for mapping with custom functions
- [ ] Add support for parallel processing
- [ ] Add mapIndexed variant that includes index in transformation

### Related Nodes

- **Filter Node** - Pre-filter before mapping
- **Reduce Node** - Aggregate mapped results
- **Extract Node** - Single-value field extraction (not array-based)
- **Transform Node** - General transformation (not array-specific)

---

## 4. Reduce Node

**Implementation Status:** 🟡 **Beta** (70% Complete - Structure Complete)

### Description

The Reduce node aggregates an array to a single value by iteratively applying an expression that combines an accumulator with each element. Common use cases include summing, finding max/min, concatenating, or building complex aggregated objects.

### Current Implementation Status

✅ **Backend:** Structure implemented (`backend/pkg/executor/reduce.go`)
- ✅ Executor structure complete
- ✅ Initial value configuration working
- ✅ Iteration framework in place
- ✅ Validation complete
- ⚠️ Expression evaluation with accumulator needs enhancement
- ⚠️ Ternary operator support needed (`a > b ? a : b`)
- ⚠️ Variable resolution in arithmetic (TODO)
- Test coverage: Basic validation tests in `reduce_test.go`

❌ **Frontend:** Not yet implemented
- Need ReduceNode.tsx component
- Configuration for initial_value and expression
- Should show accumulator concept visually

⚠️ **Limitations:**
- Expression-based reductions (`accumulator + item`) require expression enhancement
- No ternary operator yet (needed for max/min)
- Workaround: Use Map to extract values + Math for aggregation

### Node Type

`reduce`

### Responsibilities

- Reduce array to single accumulated value
- Maintain accumulator state across iterations
- Provide access to `accumulator`, `item`, `index`, `items` variables
- Support custom initial value
- Handle reduction errors gracefully
- Track iteration statistics

### Input Types

**Required:** Array of any type

### Output Format

```json
{
  "result": <any>,           // Final accumulated value
  "initial_value": <any>,    // Starting accumulator value
  "final_value": <any>,      // Same as result
  "input_count": number,     // Array length
  "iterations": number,      // Total iterations attempted
  "successful": number,      // Successful iterations
  "failed": number           // Failed iterations
}
```

### Node Properties

```typescript
{
  "type": "reduce",
  "data": {
    "initial_value": any,     // Optional: Starting accumulator (default: 0)
    "expression": string      // Required: Expression to compute new accumulator
  }
}
```

**Property Details:**

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `initial_value` | any | ❌ No | 0 | Starting value for accumulator |
| `expression` | string | ✅ Yes | N/A | Expression to update accumulator |

### Expression Syntax

The reduce expression has access to special variables:

```javascript
accumulator + item              // Sum numbers
accumulator + item.age          // Sum ages
item > accumulator ? item : accumulator  // Find max (requires ternary support)
accumulator + item.price        // Sum prices
```

**Variables available:**
- `accumulator` - Current accumulated value
- `item` - Current array element
- `index` - Current index (0-based)
- `items` - Full input array
- All workflow variables and context

### Error Scenarios

| Scenario | Behavior | Error Message |
|----------|----------|---------------|
| No input | ❌ Execution fails | "reduce node needs at least 1 input" |
| Non-array input | ❌ Execution fails | "reduce node requires array input, got <type>" |
| No expression | ❌ Validation fails | "reduce node requires an 'expression'" |
| Empty expression | ❌ Validation fails | "reduce node requires an 'expression'" |
| Expression error | ⚠️ Skips iteration | Increments failed count, logs debug |

### Branches

**Single Output Branch:**
- Contains `result` field with final accumulated value

### Example Workflows

#### Example 1: Sum Numbers (TODO - requires expression enhancement)
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {"value": [1, 2, 3, 4, 5]}
    },
    {
      "id": "2",
      "type": "reduce",
      "data": {
        "initial_value": 0,
        "expression": "accumulator + item"
      }
    }
  ]
}
```
**Expected Output:** `{"result": 15, "successful": 5}`  
**Status:** ⚠️ Requires expression enhancement

#### Example 2: Sum Ages
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {
        "value": [
          {"name": "Alice", "age": 25},
          {"name": "Bob", "age": 30}
        ]
      }
    },
    {
      "id": "2",
      "type": "reduce",
      "data": {
        "initial_value": 0,
        "expression": "accumulator + item.age"
      }
    }
  ]
}
```
**Expected Output:** `{"result": 55}`  
**Status:** ⚠️ Requires expression enhancement

#### Example 3: Find Maximum (requires ternary operator)
```json
{
  "data": {
    "initial_value": 0,
    "expression": "item > accumulator ? item : accumulator"
  }
}
```
**Status:** ⚠️ Requires ternary operator support

#### Example 4: Concatenate Strings
```json
{
  "data": {
    "initial_value": "",
    "expression": "accumulator + item"
  }
}
```
**Input:** `["Hello", " ", "World"]`  
**Expected:** `"Hello World"`

#### Example 5: Calculate Average (multi-step)
```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": [10, 20, 30]}},
    {"id": "2", "type": "reduce", "data": {"initial_value": 0, "expression": "accumulator + item"}},
    {"id": "3", "type": "math", "data": {"operation": "divide", "value": 3}}
  ]
}
```
**Flow:** Sum → Divide by count → Average

#### Example 6: Sum Prices with Discount
```json
{
  "data": {
    "initial_value": 0,
    "expression": "accumulator + (item.price * 0.9)"
  }
}
```

#### Example 7: Count Active Users
```json
{
  "data": {
    "initial_value": 0,
    "expression": "item.active == true ? accumulator + 1 : accumulator"
  }
}
```

#### Example 8: Build Object (advanced)
```json
{
  "data": {
    "initial_value": {},
    "expression": "{...accumulator, [item.id]: item.value}"
  }
}
```
**Status:** ⚠️ Requires advanced object spread support

### Limitations

⚠️ **Current Limitations:**
- Expression-based reductions require expression system enhancement
- No ternary operator support yet (`a ? b : c`)
- Variable resolution in arithmetic expressions needs work
- Cannot use complex object operations
- No support for early termination
- No support for async operations in expressions

### TODOs

- [ ] Implement expression evaluation with variable resolution
- [ ] Add ternary operator support for conditional reductions
- [ ] Add support for object spread in expressions
- [ ] Add reduceRight variant (right-to-left)
- [ ] Add support for early termination (break condition)
- [ ] Add built-in aggregation functions (sum, avg, max, min, etc.)

### Related Nodes

- **Map Node** - Transform before reducing
- **Filter Node** - Filter before reducing
- **Accumulator Node** - Alternative for simple accumulation
- **Math Node** - Post-process reduced result

---

## 5. ForEach Node

**Implementation Status:** 🟡 **Beta** (70% Complete - Iteration Working)

### Description

The ForEach node is a **pure iterator** that executes for each element in an array. Unlike Map or Reduce, it doesn't collect results - it's designed for **side effects** like making API calls, sending notifications, or logging. It injects iteration variables (`item`, `index`, `items`) that child nodes can access.

### Current Implementation Status

✅ **Backend:** Core implemented (`backend/pkg/executor/foreach.go`)
- ✅ Iteration logic complete
- ✅ Context variables prepared (item/index/items)
- ✅ Max iteration safety limits
- ✅ Error handling (continue on error)
- ✅ Statistics tracking
- ⚠️ Child node execution requires workflow engine integration (TODO)
- ⚠️ Variable injection handled by engine, not executor
- Test coverage: 7 test cases in `foreach_simplified_test.go`

❌ **Frontend:** Not yet implemented
- Need ForEachNode.tsx component
- Configuration for max_iterations
- Visual indication of iteration concept

⚠️ **Limitations:**
- Child nodes are not yet executed (marked TODO in code)
- Variables are prepared but injection happens at engine level
- Currently returns metadata only
- Full functionality requires workflow engine integration

### Node Type

`for_each`

### Responsibilities

- Iterate over array elements
- Inject iteration context variables for each element
- Execute child nodes for each iteration (via workflow engine)
- Track iteration statistics
- Prevent resource exhaustion with max iteration limits
- Handle iteration errors gracefully (continue on error)

### Input Types

**Required:** Array of any type

### Output Format

```json
{
  "input_count": number,    // Original array length
  "iterations": number,     // Total iterations executed
  "successful": number,     // Successful iterations
  "failed": number          // Failed iterations
}
```

### Node Properties

```typescript
{
  "type": "for_each",
  "data": {
    "max_iterations": number  // Optional: Max iterations allowed (default: 1000)
  }
}
```

**Property Details:**

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `max_iterations` | number | ❌ No | 1000 | Maximum iterations to prevent resource exhaustion |

### Iteration Variables

For each iteration, the following variables are injected:

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `variables.item` | any | Current array element | `{"name": "Alice", "age": 25}` |
| `variables.index` | number | Current index (0-based) | `0`, `1`, `2` |
| `variables.items` | array | Full input array | `[{...}, {...}]` |

**Accessing in child nodes:**
```javascript
variables.item.name          // "Alice"
variables.item.age           // 25
variables.index              // 0
variables.items.length       // 3
```

### Error Scenarios

| Scenario | Behavior | Error Message |
|----------|----------|---------------|
| No input | ❌ Execution fails | "for_each node needs at least 1 input" |
| Non-array input | ❌ Execution fails | "for_each node requires array input, got <type>" |
| Exceeds max iterations | ❌ Execution fails | "for_each exceeds max iterations: X > Y" |
| Iteration error | ⚠️ Continues, increments failed | Logs debug, continues to next iteration |

### Branches

**Single Output Branch:**
- Contains execution metadata
- No data collection (use Map/Reduce for that)

### Example Workflows

#### Example 1: Send Email to Each User
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {
        "value": [
          {"email": "alice@example.com", "name": "Alice"},
          {"email": "bob@example.com", "name": "Bob"}
        ]
      }
    },
    {
      "id": "2",
      "type": "for_each",
      "data": {}
    },
    {
      "id": "3",
      "type": "http",
      "data": {
        "url": "https://api.example.com/send-email",
        "method": "POST",
        "body": {
          "to": "{{variables.item.email}}",
          "name": "{{variables.item.name}}"
        }
      }
    }
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3"}
  ]
}
```
**Result:** Sends email to Alice and Bob

#### Example 2: Log Each Item
```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": [1, 2, 3]}},
    {"id": "2", "type": "for_each", "data": {}},
    {"id": "3", "type": "visualization", "data": {"value": "Item {{variables.index}}: {{variables.item}}"}}
  ]
}
```

#### Example 3: Batch Processing with Rate Limiting
```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": [/*large array*/]}},
    {"id": "2", "type": "for_each", "data": {"max_iterations": 100}},
    {"id": "3", "type": "delay", "data": {"ms": 100}},
    {"id": "4", "type": "http", "data": {/*API call*/}}
  ]
}
```

#### Example 4: Create Resources
```json
{
  "nodes": [
    {
      "id": "1",
      "type": "number",
      "data": {
        "value": [
          {"name": "Project A", "owner": "Alice"},
          {"name": "Project B", "owner": "Bob"}
        ]
      }
    },
    {"id": "2", "type": "for_each", "data": {}},
    {
      "id": "3",
      "type": "http",
      "data": {
        "url": "https://api.example.com/projects",
        "method": "POST",
        "body": "{{variables.item}}"
      }
    }
  ]
}
```

#### Example 5: Conditional Side Effects
```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": [/*users*/]}},
    {"id": "2", "type": "for_each", "data": {}},
    {"id": "3", "type": "condition", "data": {"condition": "variables.item.verified == true"}},
    {"id": "4", "type": "http", "data": {/*send to verified users only*/}}
  ]
}
```

#### Example 6: Update Multiple Records
```json
{
  "data": {
    "max_iterations": 50
  }
}
```
**Use case:** Update 50 database records via HTTP PATCH

#### Example 7: Parallel Notifications
```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": [/*subscribers*/]}},
    {"id": "2", "type": "for_each", "data": {}},
    {"id": "3", "type": "http", "data": {"url": "https://api.example.com/notify", "body": "{{variables.item}}"}}
  ]
}
```

#### Example 8: Data Migration
```json
{
  "nodes": [
    {"id": "1", "type": "http", "data": {"url": "https://old-system.com/data"}},
    {"id": "2", "type": "extract", "data": {"field": "items"}},
    {"id": "3", "type": "for_each", "data": {"max_iterations": 1000}},
    {"id": "4", "type": "transform", "data": {/*transform data*/}},
    {"id": "5", "type": "http", "data": {"url": "https://new-system.com/import", "method": "POST"}}
  ]
}
```

### Limitations

⚠️ **Current Limitations:**
- Child node execution requires workflow engine integration (marked TODO)
- Variables are prepared but injection happens at engine level
- No support for parallel execution
- No break/continue support
- No support for nested loops (use multiple ForEach nodes)
- Cannot return collected results (use Map/Reduce instead)

### TODOs

- [ ] Implement child node execution in workflow engine
- [ ] Add support for parallel execution mode
- [ ] Add break condition support
- [ ] Add continue condition support
- [ ] Add support for batch processing (process N items at once)
- [ ] Add support for nested loop detection and warnings

### Related Nodes

- **Map Node** - Use when you need to collect transformed results
- **Reduce Node** - Use when you need to aggregate results
- **Filter Node** - Filter array before iteration
- **WhileLoop Node** - Condition-based iteration instead of array-based

---

## 6. WhileLoop Node

**Implementation Status:** 🟡 **Beta** (60% Complete - Basic Loop Working)

### Description

The WhileLoop node executes a loop while a boolean condition remains true. It's useful for iteration when you don't have a fixed-size array, such as polling, retry logic, or stateful iteration.

### Current Implementation Status

✅ **Backend:** Basic implementation (`backend/pkg/executor/whileloop.go`)
- ✅ Condition evaluation working
- ✅ Max iteration safety limits
- ✅ Simple loop structure
- ⚠️ Value update between iterations requires engine integration
- ⚠️ Child node execution not yet integrated
- ⚠️ Currently just counts iterations without meaningful processing
- Test coverage: 4 test cases in `whileloop_test.go`

🟡 **Frontend:** Basic component exists (`src/components/nodes/WhileLoopNode.tsx`)
- Needs enhancement for better UX
- Condition input field
- Max iterations configuration

⚠️ **Limitations:**
- Value doesn't update between iterations (static without engine integration)
- Child nodes not executed yet
- Limited to simple condition checking
- Full loop execution requires workflow engine integration

### Node Type

`while_loop`

### Responsibilities

- Evaluate condition before each iteration
- Loop while condition is true
- Prevent infinite loops with max iteration limit
- Execute child nodes on each iteration (via workflow engine)
- Track iteration count
- Pass through input value with metadata

### Input Types

**Accepts:** Any type (used to evaluate condition)

### Output Format

```json
{
  "final_value": <any>,     // Final value after loop completes
  "iterations": number,     // Number of iterations executed
  "condition": "string"     // Condition expression used
}
```

### Node Properties

```typescript
{
  "type": "while_loop",
  "data": {
    "condition": string,      // Required: Boolean expression
    "max_iterations": number  // Optional: Safety limit (default: 100)
  }
}
```

**Property Details:**

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `condition` | string | ✅ Yes | N/A | Boolean expression to evaluate each iteration |
| `max_iterations` | number | ❌ No | 100 | Maximum iterations to prevent infinite loops |

### Expression Syntax

Similar to Condition node:

```javascript
value < 100             // Loop while value is less than 100
count > 0               // Loop while count is positive
retries < maxRetries    // Retry logic
status != "complete"    // Poll until complete
```

### Error Scenarios

| Scenario | Behavior | Error Message |
|----------|----------|---------------|
| Missing condition | ❌ Execution fails | "while_loop node missing condition" |
| No input | ❌ Execution fails | "while_loop node needs at least 1 input" |
| Exceeds max iterations | ❌ Execution fails | "while_loop exceeded max iterations: X" |
| Condition evaluation error | ⚠️ Treats as false | Loop terminates |

### Branches

**Single Output Branch:**
- Contains final value and iteration metadata

### Example Workflows

#### Example 1: Countdown Loop
```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": 5}},
    {"id": "2", "type": "while_loop", "data": {"condition": "value > 0", "max_iterations": 10}},
    {"id": "3", "type": "math", "data": {"operation": "subtract", "value": 1}}
  ]
}
```
**Note:** Full implementation requires workflow engine to update value each iteration

#### Example 2: Polling API Until Complete
```json
{
  "nodes": [
    {"id": "1", "type": "http", "data": {"url": "https://api.example.com/job/123"}},
    {"id": "2", "type": "extract", "data": {"field": "status"}},
    {"id": "3", "type": "while_loop", "data": {"condition": "value != \"complete\"", "max_iterations": 50}},
    {"id": "4", "type": "delay", "data": {"ms": 1000}},
    {"id": "5", "type": "http", "data": {"url": "https://api.example.com/job/123"}}
  ]
}
```

#### Example 3: Retry Logic
```json
{
  "data": {
    "condition": "retries < 3 && success != true",
    "max_iterations": 3
  }
}
```

#### Example 4: Accumulate Until Threshold
```json
{
  "data": {
    "condition": "total < variables.targetAmount",
    "max_iterations": 100
  }
}
```

#### Example 5: Process Until Empty
```json
{
  "data": {
    "condition": "queue.length > 0",
    "max_iterations": 1000
  }
}
```

#### Example 6: Fibonacci Sequence
```json
{
  "data": {
    "condition": "n < 100",
    "max_iterations": 20
  }
}
```

#### Example 7: Binary Search
```json
{
  "data": {
    "condition": "left <= right",
    "max_iterations": 32
  }
}
```

#### Example 8: Drain Buffer
```json
{
  "data": {
    "condition": "buffer.hasData == true",
    "max_iterations": 500
  }
}
```

### Limitations

⚠️ **Current Limitations:**
- Child node execution requires workflow engine integration (marked TODO)
- Value doesn't update between iterations in current implementation
- No support for loop variables beyond input value
- No support for break/continue
- Condition is re-evaluated but value is static without engine integration

### TODOs

- [ ] Implement child node execution in workflow engine
- [ ] Add support for loop variables that persist across iterations
- [ ] Add break condition support
- [ ] Add support for loop state management
- [ ] Add iteration timeout support
- [ ] Improve infinite loop detection

### Related Nodes

- **ForEach Node** - Use for array-based iteration
- **Condition Node** - Single condition evaluation (no loop)
- **Retry Node** - Specialized retry logic

---

## 7. Switch Node

**Implementation Status:** 🟢 **Production Ready** (100% Complete)

### Description

The Switch node enables **multi-way branching** based on value matching or condition evaluation. Similar to switch/case statements in programming languages, it routes execution to different paths based on the input value.

### Current Implementation Status

✅ **Backend:** Fully implemented (`backend/pkg/executor/switch.go`)
- Value matching mode complete
- Condition matching mode complete
- Default path support
- Case evaluation logic
- Test coverage: 7 test cases in `switch_test.go`

🟡 **Frontend:** Basic component exists (`src/components/nodes/SwitchNode.tsx`)
- Core functionality working
- Could use UI enhancement for better case management
- Multiple output handles

✅ **Features:** All core features complete
- Both value and condition matching
- Custom output paths
- First-match-wins semantics
- Production-ready

### Node Type

`switch`

### Responsibilities

- Evaluate input against multiple cases
- Support value matching (equality) and condition matching (expressions)
- Route to first matching case
- Provide default path when no cases match
- Pass through input value with metadata about matched case

### Input Types

**Accepts:** Any type

### Output Format

```json
{
  "value": <input>,         // Original input value
  "matched": true|false,    // Whether any case matched
  "case": "string",         // Matched case expression (if matched)
  "output_path": "string"   // Path taken (case name or "default")
}
```

### Node Properties

```typescript
{
  "type": "switch",
  "data": {
    "cases": [               // Required: Array of case definitions
      {
        "when": string,      // Case expression or label
        "value": any,        // Optional: Value to match against
        "output_path": string  // Optional: Custom path name
      }
    ],
    "default_path": string   // Optional: Default path name (default: "default")
  }
}
```

**Property Details:**

| Property | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `cases` | array | ✅ Yes | N/A | List of case definitions (min: 1) |
| `cases[].when` | string | ✅ Yes | N/A | Case label or condition expression |
| `cases[].value` | any | ❌ No | N/A | Value to match (if set, uses value matching) |
| `cases[].output_path` | string | ❌ No | "matched" | Custom output path name |
| `default_path` | string | ❌ No | "default" | Path name when no cases match |

### Case Matching Modes

#### Mode 1: Value Matching
When `value` is specified in a case, uses equality comparison:

```json
{
  "cases": [
    {"when": "admin", "value": "admin"},
    {"when": "user", "value": "user"},
    {"when": "guest", "value": "guest"}
  ]
}
```

#### Mode 2: Condition Matching
When `value` is NOT specified, evaluates `when` as an expression:

```json
{
  "cases": [
    {"when": "value > 100"},
    {"when": "value >= 50 && value <= 100"},
    {"when": "value < 50"}
  ]
}
```

### Error Scenarios

| Scenario | Behavior | Error Message |
|----------|----------|---------------|
| No input | ❌ Execution fails | "switch node requires at least one input" |
| No cases | ❌ Validation fails | "switch node requires at least one case" |
| Empty cases array | ❌ Validation fails | "switch node requires at least one case" |
| Case evaluation error | ⚠️ Treats as no match | Continues to next case |
| All cases fail | ✅ Takes default path | `matched: false` |

### Branches

**Multiple Output Branches:**
- One branch per case (named by `output_path` or "matched")
- One default branch (named by `default_path` or "default")

### Example Workflows

#### Example 1: Role-Based Routing
```json
{
  "nodes": [
    {"id": "1", "type": "text_input", "data": {"value": "admin"}},
    {
      "id": "2",
      "type": "switch",
      "data": {
        "cases": [
          {"when": "admin", "value": "admin", "output_path": "admin_flow"},
          {"when": "user", "value": "user", "output_path": "user_flow"},
          {"when": "guest", "value": "guest", "output_path": "guest_flow"}
        ],
        "default_path": "unknown_role"
      }
    },
    {"id": "3", "type": "text_input", "data": {"value": "Admin Dashboard"}},
    {"id": "4", "type": "text_input", "data": {"value": "User Dashboard"}},
    {"id": "5", "type": "text_input", "data": {"value": "Guest Landing"}}
  ],
  "edges": [
    {"source": "1", "target": "2"},
    {"source": "2", "target": "3", "sourceHandle": "admin_flow"},
    {"source": "2", "target": "4", "sourceHandle": "user_flow"},
    {"source": "2", "target": "5", "sourceHandle": "guest_flow"}
  ]
}
```

#### Example 2: Number Range Classification
```json
{
  "data": {
    "cases": [
      {"when": "value >= 90", "output_path": "excellent"},
      {"when": "value >= 70", "output_path": "good"},
      {"when": "value >= 50", "output_path": "average"},
      {"when": "value < 50", "output_path": "poor"}
    ]
  }
}
```

#### Example 3: HTTP Status Code Handling
```json
{
  "data": {
    "cases": [
      {"when": "2xx", "value": 200, "output_path": "success"},
      {"when": "4xx", "value": 404, "output_path": "not_found"},
      {"when": "5xx", "value": 500, "output_path": "server_error"}
    ],
    "default_path": "unexpected"
  }
}
```

#### Example 4: Order Status Workflow
```json
{
  "data": {
    "cases": [
      {"when": "pending", "value": "pending", "output_path": "process_order"},
      {"when": "paid", "value": "paid", "output_path": "fulfill_order"},
      {"when": "shipped", "value": "shipped", "output_path": "track_shipment"},
      {"when": "delivered", "value": "delivered", "output_path": "request_review"}
    ]
  }
}
```

#### Example 5: Priority Routing
```json
{
  "data": {
    "cases": [
      {"when": "value == \"critical\"", "output_path": "immediate"},
      {"when": "value == \"high\"", "output_path": "urgent"},
      {"when": "value == \"medium\"", "output_path": "normal"},
      {"when": "value == \"low\"", "output_path": "backlog"}
    ]
  }
}
```

#### Example 6: Date-Based Routing
```json
{
  "data": {
    "cases": [
      {"when": "dayOfWeek(value) >= 1 && dayOfWeek(value) <= 5", "output_path": "weekday"},
      {"when": "dayOfWeek(value) == 0 || dayOfWeek(value) == 6", "output_path": "weekend"}
    ]
  }
}
```

#### Example 7: Temperature Classification
```json
{
  "data": {
    "cases": [
      {"when": "value > 30", "output_path": "hot"},
      {"when": "value >= 20 && value <= 30", "output_path": "warm"},
      {"when": "value >= 10 && value < 20", "output_path": "cool"},
      {"when": "value < 10", "output_path": "cold"}
    ]
  }
}
```

#### Example 8: Content Type Router
```json
{
  "data": {
    "cases": [
      {"when": "image", "value": "image/png", "output_path": "image_processor"},
      {"when": "video", "value": "video/mp4", "output_path": "video_processor"},
      {"when": "document", "value": "application/pdf", "output_path": "doc_processor"}
    ],
    "default_path": "unsupported_type"
  }
}
```

### Limitations

⚠️ **Current Limitations:**
- Cases are evaluated sequentially (first match wins)
- No fallthrough support (unlike C switch)
- Cannot have overlapping conditions (first match takes precedence)
- Value matching uses simple equality (no pattern matching)
- No range matching syntax (must use condition expressions)

### TODOs

- [ ] Add support for multiple values per case (OR logic)
- [ ] Add support for range matching syntax (`10..20`)
- [ ] Add support for pattern matching
- [ ] Add support for case priority/ordering
- [ ] Add warnings for unreachable cases
- [ ] Add support for computed case values

### Related Nodes

- **Condition Node** - Two-way branching (simpler)
- **Filter Node** - Array filtering (different purpose)

---

## Recommended New Node Types

To make Thaiyyal workflows more flexible and powerful, the following new node types are recommended based on common workflow patterns and gaps in current functionality:

### 🔴 High Priority (Essential for Common Workflows)

#### 1. **FlatMap Node** 
**Purpose:** Transform each element to an array and flatten the results

**Use Cases:**
- Expand nested arrays: `[{tags:["a","b"]}, {tags:["c"]}]` → `["a","b","c"]`
- Generate multiple items from one: `[1,2,3]` → `[1,1, 2,2, 3,3]`
- Combine Map + Flatten in one step

**Configuration:**
```json
{
  "type": "flat_map",
  "data": {
    "field": "tags",  // or
    "expression": "item.items"
  }
}
```

**Why Important:** Very common pattern in data processing (logs, tags, nested data)

---

#### 2. **Slice Node**
**Purpose:** Extract a portion of an array (pagination, windowing)

**Use Cases:**
- Pagination: Get first 10 items
- Batch processing: Process array in chunks
- Windowing: Get last N items
- Skip and take patterns

**Configuration:**
```json
{
  "type": "slice",
  "data": {
    "start": 0,      // Start index (default: 0)
    "end": 10,       // End index (optional)
    "length": 10     // Alternative to end (take N items)
  }
}
```

**Examples:**
- `slice(0, 10)` - First 10 items
- `slice(-5)` - Last 5 items
- `slice(10, 20)` - Items 10-20

**Why Important:** Essential for pagination, batch processing, limiting results

---

#### 3. **Sort Node**
**Purpose:** Sort array elements by field or expression

**Use Cases:**
- Sort users by age
- Sort products by price
- Sort events by timestamp
- Custom sort logic

**Configuration:**
```json
{
  "type": "sort",
  "data": {
    "field": "age",           // Sort by field
    "order": "asc" | "desc",  // Sort direction
    "expression": "item.price * item.quantity"  // Or custom expression
  }
}
```

**Why Important:** Very common requirement, currently requires external processing

---

#### 4. **GroupBy Node**
**Purpose:** Group array elements by a field value

**Use Cases:**
- Group users by department
- Group orders by status
- Group events by date
- Aggregate by category

**Configuration:**
```json
{
  "type": "group_by",
  "data": {
    "field": "department",  // Group by this field
    "aggregate": "count" | "sum" | "avg" | "max" | "min"  // Optional
  }
}
```

**Output:**
```json
{
  "groups": {
    "Engineering": [{...}, {...}],
    "Sales": [{...}],
    "Marketing": [{...}]
  },
  "counts": {
    "Engineering": 15,
    "Sales": 8,
    "Marketing": 5
  }
}
```

**Why Important:** Common analytics/reporting pattern

---

#### 5. **Find Node**
**Purpose:** Find first element matching condition

**Use Cases:**
- Find user by ID
- Find first item above threshold
- Search in array
- Existence checks

**Configuration:**
```json
{
  "type": "find",
  "data": {
    "condition": "item.id == variables.searchId",
    "return_index": false  // Return index instead of element
  }
}
```

**Why Important:** More efficient than Filter when you only need one result

---

### 🟡 Medium Priority (Enhance Workflow Capabilities)

#### 6. **Unique Node**
**Purpose:** Remove duplicate elements from array

**Use Cases:**
- Deduplicate list
- Get unique values
- Clean data

**Configuration:**
```json
{
  "type": "unique",
  "data": {
    "field": "id"  // Use specific field for uniqueness
  }
}
```

---

#### 7. **Zip Node**
**Purpose:** Combine multiple arrays element-wise

**Use Cases:**
- Combine parallel arrays
- Merge data from multiple sources
- Create pairs/tuples

**Configuration:**
```json
{
  "type": "zip",
  "data": {
    "arrays": ["node.array1", "node.array2"],
    "fill_missing": null  // Value for shorter arrays
  }
}
```

**Example:**
```javascript
[1,2,3] + [a,b,c] → [[1,a], [2,b], [3,c]]
```

---

#### 8. **Partition Node**
**Purpose:** Split array into two groups based on condition

**Use Cases:**
- Separate valid/invalid items
- Split pass/fail results
- Bifurcate data flow

**Configuration:**
```json
{
  "type": "partition",
  "data": {
    "condition": "item.age >= 18"
  }
}
```

**Output:**
```json
{
  "passed": [...],   // Items where condition is true
  "failed": [...]    // Items where condition is false
}
```

**Why Useful:** Alternative to multiple Filter nodes

---

#### 9. **Chunk Node**
**Purpose:** Split array into fixed-size chunks

**Use Cases:**
- Batch processing
- Paginate large arrays
- Rate limiting groups

**Configuration:**
```json
{
  "type": "chunk",
  "data": {
    "size": 10  // Chunk size
  }
}
```

**Example:**
```javascript
[1,2,3,4,5,6,7] → [[1,2,3], [4,5,6], [7]]
```

---

#### 10. **Reverse Node**
**Purpose:** Reverse array order

**Use Cases:**
- Reverse chronological order
- Reverse processing
- LIFO operations

**Configuration:**
```json
{
  "type": "reverse",
  "data": {}
}
```

---

### 🟢 Low Priority (Nice to Have)

#### 11. **Sample Node**
**Purpose:** Get random sample from array

**Use Cases:**
- Random selection
- Testing with subset
- Statistical sampling

**Configuration:**
```json
{
  "type": "sample",
  "data": {
    "count": 10,        // Number of items
    "method": "random" | "first" | "last"
  }
}
```

---

#### 12. **Range Node**
**Purpose:** Generate array of numbers

**Use Cases:**
- Create sequential IDs
- Generate test data
- Iteration counts

**Configuration:**
```json
{
  "type": "range",
  "data": {
    "start": 1,
    "end": 100,
    "step": 1
  }
}
```

**Output:** `[1, 2, 3, ..., 100]`

---

#### 13. **Join Node**
**Purpose:** Join two arrays like SQL JOIN

**Use Cases:**
- Combine related data
- Lookup enrichment
- Data merging

**Configuration:**
```json
{
  "type": "join",
  "data": {
    "left": "node.users",
    "right": "node.orders",
    "on": "user_id",
    "type": "inner" | "left" | "right" | "full"
  }
}
```

---

#### 14. **Compact Node**
**Purpose:** Remove null/undefined/falsy values

**Use Cases:**
- Clean data
- Remove failed transformations
- Filter out empties

**Configuration:**
```json
{
  "type": "compact",
  "data": {
    "remove": "null" | "undefined" | "falsy" | "empty"
  }
}
```

---

#### 15. **Transpose Node**
**Purpose:** Transpose array of arrays (rows ↔ columns)

**Use Cases:**
- Matrix operations
- Pivot data
- Reshape data

---

### 🎯 Priority Implementation Roadmap

**Phase 1: Essential Array Operations** (Most Impact)
1. Slice - Pagination essential
2. Sort - Very common need
3. Find - More efficient than filter for single results
4. FlatMap - Common nested data pattern

**Phase 2: Grouping & Aggregation**
5. GroupBy - Analytics workflows
6. Partition - Alternative to multiple filters
7. Unique - Data cleaning

**Phase 3: Advanced Operations**
8. Chunk - Batch processing
9. Zip - Combining data sources
10. Join - Complex data merging

**Phase 4: Utility Operations**
11. Reverse, Sample, Range, Compact, Transpose

---

### Implementation Notes

**Frontend Requirements:**
- All new nodes need React components
- Consistent styling with existing control flow nodes
- Input validation in UI
- Node info documentation

**Backend Requirements:**
- Executor implementation following existing patterns
- Comprehensive test coverage
- Error handling and validation
- Performance optimization for large arrays

**Expression Integration:**
- Nodes should support expression language where applicable
- Variable access (variables, context, node refs)
- Consistent error handling

---

### Alternative Approaches

**Instead of many specialized nodes, consider:**

1. **Enhanced Map/Filter/Reduce** with built-in functions
   - `map(sort)`, `filter(unique)`, etc.
   
2. **General Array Node** with operation modes
   - Risk: Violates Single Responsibility Principle
   - Benefit: Fewer node types

3. **Expression Functions**
   - Add array functions to expression language
   - `sort(items, "field")`, `slice(items, 0, 10)`
   - More flexible but less visual

**Recommendation:** Stick with specialized nodes for commonly used operations (Slice, Sort, Find, GroupBy), use expression functions for rare cases.

---

## Expression Language Reference

All control flow nodes use a unified expression language for evaluating conditions and transformations.

### Supported Operators

**Comparison:**
- `>` - Greater than
- `<` - Less than
- `>=` - Greater than or equal
- `<=` - Less than or equal
- `==` - Equal (value comparison)
- `!=` - Not equal

**Boolean Logic:**
- `&&` - Logical AND
- `||` - Logical OR
- `!` - Logical NOT

**Arithmetic:** (in expression evaluation)
- `+` - Addition
- `-` - Subtraction
- `*` - Multiplication
- `/` - Division

### Variable Access

**Workflow Variables:**
```javascript
variables.threshold
variables.minAge
variables.config.timeout
```

**Context Variables:**
```javascript
context.currentUser
context.environment
context.settings.maxRetries
```

**Node Results:**
```javascript
node.sensor1.value
node.calculation.result
node.http1.data.items
```

**Iteration Variables (in Filter/Map/Reduce/ForEach):**
```javascript
item                  // Current element
item.field           // Field access
item.nested.field    // Nested field access
index                // Current index
items                // Full array
accumulator          // Current accumulated value (Reduce only)
```

### Built-in Functions

**Math Functions:**
```javascript
sqrt(value)
pow(base, exponent)
abs(value)
floor(value)
ceil(value)
round(value)
```

**Date/Time Functions:**
```javascript
year(date)
month(date)
day(date)
hour(date)
minute(date)
second(date)
parseDate(string, format)
formatDate(date, format)
```

**String Functions:**
```javascript
length(string)
contains(string, substring)
startsWith(string, prefix)
endsWith(string, suffix)
```

### Expression Examples

```javascript
// Simple comparison
age > 18

// Complex boolean logic
(age >= 18 && age <= 65) && (active == true || premium == true)

// Variable reference
price <= variables.budget

// Context reference
role == context.requiredRole

// Node reference
currentTemp > node.sensor1.maxTemp

// Array iteration (Filter node)
item.age >= 18 && item.verified == true

// Nested field access
item.profile.settings.notifications == true

// Arithmetic
price * 0.9  // 10% discount

// Function calls
year(item.birthdate) <= 2005
```

---

## Design Patterns

### Pattern 1: Filter → Map → Reduce

**Goal:** Find total price of active items

```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": [/*items*/]}},
    {"id": "2", "type": "filter", "data": {"condition": "item.active == true"}},
    {"id": "3", "type": "map", "data": {"field": "price"}},
    {"id": "4", "type": "reduce", "data": {"initial_value": 0, "expression": "accumulator + item"}}
  ]
}
```

### Pattern 2: Condition-Based ForEach

**Goal:** Process only qualified items

```json
{
  "nodes": [
    {"id": "1", "type": "number", "data": {"value": [/*items*/]}},
    {"id": "2", "type": "filter", "data": {"condition": "item.qualified == true"}},
    {"id": "3", "type": "for_each", "data": {}},
    {"id": "4", "type": "http", "data": {/*process each*/}}
  ]
}
```

### Pattern 3: Switch → Different Transforms

**Goal:** Route data to different processors

```json
{
  "nodes": [
    {"id": "1", "type": "http", "data": {/*fetch data*/}},
    {"id": "2", "type": "switch", "data": {/*route by type*/}},
    {"id": "3", "type": "transform", "data": {/*process images*/}},
    {"id": "4", "type": "transform", "data": {/*process videos*/}}
  ]
}
```

### Pattern 4: While → Polling

**Goal:** Poll until complete

```json
{
  "nodes": [
    {"id": "1", "type": "http", "data": {/*start job*/}},
    {"id": "2", "type": "while_loop", "data": {"condition": "status != \"complete\""}},
    {"id": "3", "type": "delay", "data": {"ms": 1000}},
    {"id": "4", "type": "http", "data": {/*check status*/}}
  ]
}
```

---

## Troubleshooting

### Common Issues

**Issue:** Filter returns empty array

**Solutions:**
- Check that `item.field` syntax is used correctly
- Verify field names match the data
- Test expression in Condition node first
- Check for typos in field names

---

**Issue:** Map returns all `null` values

**Solutions:**
- Ensure field name exists in objects
- Check for case sensitivity (`Name` vs `name`)
- Verify input is array of objects
- Check expression syntax

---

**Issue:** ForEach doesn't seem to do anything

**Solutions:**
- ForEach is for side effects, not data collection
- Use Map/Reduce if you need to collect results
- Check that child nodes reference `variables.item`

---

**Issue:** Reduce expression doesn't work

**Solutions:**
- Expression-based reductions require enhancement (see TODOs)
- Use field extraction with Map first
- Check that expression syntax is correct
- Verify `accumulator` and `item` are accessible

---

**Issue:** WhileLoop runs forever

**Solutions:**
- Set appropriate `max_iterations`
- Ensure condition eventually becomes false
- Check that child nodes update the value
- Add safety limits

---

### Debug Tips

1. **Use Visualization nodes** to inspect intermediate results
2. **Start simple** - test with small arrays and simple conditions
3. **Check logs** - look for debug messages in console
4. **Test expressions** in Condition node before using in Filter/Map
5. **Use Variable nodes** to set test data
6. **Compose incrementally** - add one node at a time

---

## Summary

Thaiyyal's control flow nodes provide a **powerful, composable toolkit** for building sophisticated data processing workflows:

- ✅ **7 specialized nodes** each with single responsibility
- ✅ **Unified expression language** across all nodes
- ✅ **Composable design** - mix and match freely
- ✅ **Safety built-in** - iteration limits, error handling
- ✅ **Production ready** - comprehensive test coverage

**Next Steps:**
- Explore the [Expression Language](#expression-language-reference)
- Try the [Design Patterns](#design-patterns)
- Check [Example Workflows](#example-workflows) for each node
- Review [Limitations](#limitations) and [TODOs](#todos)

---

**Document Maintained By:** Thaiyyal Team  
**Last Review:** 2025-11-01  
**Version:** 1.0
