# Conditional Branching Examples

This directory contains comprehensive examples demonstrating conditional branching and **conditional execution** capabilities in Thaiyyal workflows.

## Overview

Conditional branching allows workflows to make decisions and route data based on conditions. Thaiyyal provides two primary nodes for conditional logic:

- **Condition Node**: Evaluates a boolean expression and provides "true" and "false" output paths
- **Switch Node**: Multi-way branching with pattern matching and custom output paths

## ✨ NEW: Conditional Execution (Path Termination)

**As of this update**, Thaiyyal now supports **conditional execution** where nodes are only executed if their incoming edge conditions are satisfied. This enables true branching workflows where only one path executes based on runtime conditions.

### How It Works

1. **Condition Nodes** have two output handles:
   - 🟢 **True Handle** (top-right): Connects to nodes that execute when condition is true
   - 🔴 **False Handle** (bottom-right): Connects to nodes that execute when condition is false

2. **Switch Nodes** have output paths based on cases:
   - Each case can have a custom `output_path` (e.g., "success", "error", "not_found")
   - Edges from switch nodes use `sourceHandle` to specify which path they represent

3. **The Engine** automatically:
   - Evaluates conditions during execution
   - Skips nodes whose incoming edge conditions aren't satisfied
   - Executes only the nodes in the active path

### Example: Age-Based Routing

```json
{
  "nodes": [
    {"id": "age", "type": "number", "data": {"value": 25}},
    {"id": "check", "type": "condition", "data": {"condition": ">=18"}},
    {"id": "adult_path", "type": "text_input", "data": {"text": "Adult"}},
    {"id": "minor_path", "type": "text_input", "data": {"text": "Minor"}}
  ],
  "edges": [
    {"source": "age", "target": "check"},
    {"source": "check", "target": "adult_path", "sourceHandle": "true"},
    {"source": "check", "target": "minor_path", "sourceHandle": "false"}
  ]
}
```

**Result**: If age is 25, only `adult_path` executes. The `minor_path` node is **skipped entirely**.

## Examples

### 01. Basic Age Check
**File**: `01-basic-age-check.json`

Simple conditional check to determine if a user is an adult (age >= 18).

**Demonstrates**:
- Basic greater-than-or-equal comparison
- Condition node usage
- Path metadata (true/false paths)

**Concepts**: Simple conditional, comparison operators

---

### 02. Grade Calculation
**File**: `02-grade-calculation.json`

Uses switch statement to assign letter grades (A-F) based on numeric test scores.

**Demonstrates**:
- Switch node with multiple cases
- Range-based pattern matching
- Default path handling

**Concepts**: Multi-way branching, grade thresholds, fallback logic

---

### 03. Nested Eligibility Check
**File**: `03-nested-eligibility.json`

Multi-level conditional checks for loan eligibility (age range validation).

**Demonstrates**:
- Nested conditional logic (2 levels)
- Sequential condition evaluation
- Value extraction between conditions

**Concepts**: Nested conditions, multi-stage validation

---

### 04. Data Validation
**File**: `04-data-validation.json`

Validates input data is within acceptable range (0-100).

**Demonstrates**:
- Complex AND logic (`&&`)
- Range validation
- Input sanitization pattern

**Concepts**: Data validation, boundary checking, boolean AND

---

### 05. A/B Testing
**File**: `05-ab-testing.json`

Routes users to different feature variants based on user ID (even/odd split).

**Demonstrates**:
- Modulo operator for distribution
- Feature flag pattern
- User segmentation

**Concepts**: A/B testing, feature flags, user routing

---

### 06. Multi-Tenant Processing
**File**: `06-multi-tenant.json`

Routes data processing based on tenant tier (premium vs standard).

**Demonstrates**:
- Tenant-based routing
- Tier identification
- Multi-tenant architecture pattern

**Concepts**: Multi-tenancy, tier-based processing

---

### 07. Complex Boolean Logic
**File**: `07-complex-boolean.json`

Temperature comfort zone check using complex AND/OR combinations.

**Demonstrates**:
- Mixed AND (`&&`) and OR (`||`) operators
- Grouped conditions with parentheses
- Multiple logical clauses

**Concepts**: Complex boolean expressions, operator precedence

---

### 08. Arithmetic Then Condition
**File**: `08-arithmetic-condition.json`

Calculates discounted price then validates against minimum price threshold.

**Demonstrates**:
- Chaining arithmetic operations with conditions
- Post-calculation validation
- Workflow composition

**Concepts**: Operation sequencing, calculated conditions

---

## Condition Node

### Syntax

The condition node accepts expressions in the `condition` field:

```json
{
  "type": "condition",
  "data": {
    "condition": "<expression>"
  }
}
```

### Supported Operators

#### Comparison Operators
- `==` - Equality
- `!=` - Inequality
- `>`  - Greater than
- `<`  - Less than
- `>=` - Greater than or equal
- `<=` - Less than or equal

#### Logical Operators
- `&&` - AND
- `||` - OR
- `!`  - NOT

#### Arithmetic Operators (in expressions)
- `+` - Addition
- `-` - Subtraction
- `*` - Multiplication
- `/` - Division
- `%` - Modulo

### Expression Context

Conditions can reference:
- `input` - The input value to the condition node
- `variables.<name>` - Workflow variables
- `context.<name>` - Context variables
- `node.<id>.<field>` - Other node results

### Output Structure

Condition nodes return:

```json
{
  "value": <original-input>,
  "condition_met": true|false,
  "condition": "<expression>",
  "path": "true"|"false",
  "true_path": true|false,
  "false_path": true|false
}
```

## Switch Node

### Syntax

The switch node evaluates cases in order and takes the first match:

```json
{
  "type": "switch",
  "data": {
    "cases": [
      {
        "when": "<condition>",
        "value": <optional-match-value>,
        "output_path": "<path-name>"
      }
    ],
    "default_path": "<default-path-name>"
  }
}
```

### Output Structure

Switch nodes return:

```json
{
  "value": <original-input>,
  "matched": true|false,
  "case": "<matched-condition>",
  "output_path": "<path-name>"
}
```

## Common Patterns

### Pattern 1: Simple Validation

```
Input → Condition (>=0) → Process if valid
```

Use when: Single validation check needed

### Pattern 2: Multi-Stage Validation

```
Input → Check1 → Extract → Check2 → Extract → Process
```

Use when: Multiple independent validations required

### Pattern 3: Route by Value

```
Input → Switch(cases) → Different processing paths
```

Use when: Need different logic for different input values

### Pattern 4: Grade/Score Ranges

```
Score → Switch(>=90, >=80, >=70, ...) → Grade Assignment
```

Use when: Categorizing continuous values into discrete buckets

### Pattern 5: Feature Flags

```
UserID → Condition(id % 2) → Feature A or B
```

Use when: A/B testing, gradual rollouts

### Pattern 6: Eligibility Checks

```
Input → AgeCheck → IncomeCheck → CreditCheck → Approve/Deny
```

Use when: Multiple criteria must all be satisfied

## Best Practices

### 1. Keep Conditions Simple
✅ Good: `input > 18`  
❌ Avoid: Deeply nested complex expressions

### 2. Use Descriptive Labels
```json
{
  "data": {
    "condition": ">=18",
    "label": "Is Adult?"  // Clear intent
  }
}
```

### 3. Handle Edge Cases
Always consider:
- Zero values
- Negative numbers
- Boundary conditions (==, not just > or <)
- Null/undefined inputs

### 4. Document Complex Logic
For complex conditions, add comments in the workflow description:
```json
{
  "description": "Checks: age 21-65 AND income >= $30K AND credit score > 650"
}
```

### 5. Test Both Paths
Always test workflows with inputs that:
- Satisfy the condition (true path)
- Don't satisfy the condition (false path)
- Boundary values (exactly on the threshold)

### 6. Use Switch for Multiple Ranges
When you have 3+ discrete ranges, prefer switch over multiple nested conditions:

✅ Good:
```json
{
  "type": "switch",
  "cases": [
    {"when": ">=90", "output_path": "a"},
    {"when": ">=80", "output_path": "b"},
    {"when": ">=70", "output_path": "c"}
  ]
}
```

❌ Avoid: Multiple nested if-else conditions

## Testing Scenarios

Each example can be tested by:

1. **Loading the JSON**: Import into Thaiyyal workflow editor
2. **Modify Input Values**: Change the numeric/text inputs to test different paths
3. **Execute**: Run the workflow
4. **Verify Output**: Check the visualization node for results

### Example Test Cases

For **01-basic-age-check.json**:
- Test with age = 17 → condition_met should be false
- Test with age = 18 → condition_met should be true
- Test with age = 25 → condition_met should be true

For **02-grade-calculation.json**:
- Test with score = 95 → output_path should be "grade_a"
- Test with score = 85 → output_path should be "grade_b"
- Test with score = 50 → output_path should be "grade_f"

For **05-ab-testing.json**:
- Test with user_id = 12344 (even) → condition_met should be true
- Test with user_id = 12345 (odd) → condition_met should be false

## Advanced Topics

### Chaining Conditions

Multiple conditions can be chained for complex decision trees:

```
Input → Condition1 → Extract → Condition2 → Extract → Condition3 → Result
```

Each extract node pulls out the `value` field to pass to the next condition.

### Dynamic Thresholds

Use context variables or workflow variables to make thresholds configurable:

```json
{
  "condition": "input >= context.threshold"
}
```

### Calculated Conditions

Perform arithmetic before conditionals:

```
Value1 → Operation → Condition → Result
Value2 →
```

## Troubleshooting

### Common Issues

**Issue**: Condition always evaluates to false
- Check that the input type matches (number vs string)
- Verify the comparison operator direction
- Check for floating-point precision issues

**Issue**: Switch doesn't match any case
- Verify case order (first match wins)
- Check if default_path is set
- Ensure value types match (number vs string)

**Issue**: Complex boolean expression doesn't work as expected
- Use parentheses to control operator precedence
- Test sub-expressions independently
- Break into multiple sequential conditions if too complex

## Security Considerations

1. **Input Validation**: Always validate external inputs before processing
2. **Range Checks**: Prevent out-of-bounds values with range conditions
3. **Type Checking**: Ensure input types match expected types
4. **Default Paths**: Always provide default/fallback paths in switch statements

## Performance Notes

- Conditions are evaluated quickly (microseconds)
- Switch cases are evaluated in order - put most common cases first
- Complex nested conditions may be harder to maintain than switch statements
- Extract nodes between conditions add minimal overhead

## Related Documentation

- [Node Types Reference](../../docs/NODE_TYPES.md)
- [Expression Evaluation](../../docs/EXPRESSIONS.md)
- [Workflow Testing Guide](../../docs/TESTING.md)
- [Best Practices](../../docs/BEST_PRACTICES.md)

## Example Workflow Execution

See [Workflow Execution Guide](../../docs/WORKFLOW_EXECUTION_GUIDE.md) for step-by-step instructions on:
- Loading example workflows
- Executing workflows
- Viewing results
- Debugging issues

---

**Total Examples**: 8  
**Coverage**: Basic, Nested, Switch, Validation, Feature Flags, Multi-tenant, Boolean Logic, Arithmetic + Conditions

For more examples and tutorials, see the main [Examples Directory](../README.md).

---

### 09. Age-Based API Routing ⭐ NEW
**File**: `09-age-based-api-routing.json`

**Real-world conditional execution example**: Routes users to different APIs based on age.
- If age >= 18: Fetch profile → Register for sports
- If age < 18: Register for education

**Demonstrates**:
- Conditional execution (path termination)
- True/false path routing with `sourceHandle`
- Nodes execute only in active path

**Use Case**: User registration systems, age-gated content

---

### 10. Multi-Step Registration Flow ⭐ NEW
**File**: `10-multi-step-registration.json`

**Complex conditional workflow** with multiple steps in each branch:
- Adult path: Parse user → Fetch profile → Extract interests → Register programs
- Minor path: Direct registration to education

**Demonstrates**:
- Multi-node conditional branches
- Sequential processing within a branch
- Path convergence (both paths lead to confirmation)

**Use Case**: SaaS onboarding, tiered service registration

---

### 11. HTTP Status Code Routing ⭐ NEW
**File**: `11-http-status-routing.json`

**Switch-based conditional execution** routing based on API response codes:
- 200 → Success handler
- 404 → Retry handler
- 500+ → Error handler

**Demonstrates**:
- Switch node with custom output paths
- Multiple conditional paths from single node
- Error handling patterns

**Use Case**: API orchestration, error recovery workflows

---

## Conditional Execution API

### Edge Schema

Edges now support conditional routing through the `sourceHandle` field:

```typescript
{
  "id": "edge1",
  "source": "condition_node",
  "target": "target_node",
  "sourceHandle": "true"  // or "false", "success", "error", etc.
}
```

### Condition Node Handles

Condition nodes provide two output handles:

- **`true`**: Top-right handle (green) - connects to nodes that execute when condition is met
- **`false`**: Bottom-right handle (red) - connects to nodes that execute when condition is not met

### Switch Node Handles

Switch nodes dynamically create handles based on defined cases:

- Each case's `output_path` becomes a handle name
- Example: `"success"`, `"error"`, `"not_found"`, `"grade_a"`, etc.

### Frontend Integration (React Flow)

The frontend automatically handles conditional edges using React Flow's handle system:

```tsx
// Condition node with multiple handles
<Handle 
  type="source" 
  position={Position.Right}
  id="true"
  style={{ top: "30%" }}
  className="w-2 h-2 bg-green-500"
/>
<Handle 
  type="source" 
  position={Position.Right}
  id="false"
  style={{ top: "70%" }}
  className="w-2 h-2 bg-red-500"
/>
```

When users connect edges in the UI, React Flow automatically includes `sourceHandle` and `targetHandle` in the edge data.

### Backend Processing

The workflow engine processes conditional edges:

1. **Topological Sort**: All nodes are sorted for execution order
2. **Condition Check**: Before executing each node, check incoming edges
3. **Path Evaluation**: 
   - If any incoming edge has `sourceHandle`, evaluate the source node's result
   - Match `sourceHandle` value against source node's output (`path`, `output_path`, etc.)
   - Skip node if no incoming edge condition is satisfied
4. **Execution**: Only execute nodes in the active path

## Visual Guide

### Creating Conditional Workflows

1. **Add Condition Node**
   - Drag "Condition" node from palette
   - Set condition expression (e.g., ">18")

2. **Connect True Path**
   - Click on the **green handle** (top-right)
   - Drag to target node
   - Edge automatically gets `sourceHandle: "true"`

3. **Connect False Path**
   - Click on the **red handle** (bottom-right)
   - Drag to alternate target node
   - Edge automatically gets `sourceHandle: "false"`

4. **Execute Workflow**
   - Only nodes in the active path execute
   - Inactive path nodes are skipped

### Creating Switch-Based Workflows

1. **Configure Switch Node**
   - Set cases with `output_path` values
   - Example: `[{when: "==200", output_path: "success"}]`

2. **Connect Paths**
   - Connect edges using appropriate handles
   - Each edge's `sourceHandle` matches a case's `output_path`

3. **Runtime Routing**
   - Engine evaluates which case matches
   - Only the matching path executes

## Migration Guide

### From Metadata-Only to Conditional Execution

**Before** (all nodes execute, check metadata manually):
```json
{
  "nodes": [
    {"id": "check", "type": "condition"},
    {"id": "action1", "type": "text_input"},
    {"id": "action2", "type": "text_input"}
  ],
  "edges": [
    {"source": "check", "target": "action1"},
    {"source": "check", "target": "action2"}
  ]
}
```
Result: Both `action1` and `action2` execute. You must manually check condition metadata.

**After** (only active path executes):
```json
{
  "nodes": [
    {"id": "check", "type": "condition"},
    {"id": "action1", "type": "text_input"},
    {"id": "action2", "type": "text_input"}
  ],
  "edges": [
    {"source": "check", "target": "action1", "sourceHandle": "true"},
    {"source": "check", "target": "action2", "sourceHandle": "false"}
  ]
}
```
Result: Only one action executes based on condition result!

### Backward Compatibility

- Edges without `sourceHandle` are **unconditional** (always execute)
- Mix conditional and unconditional edges as needed
- Legacy `condition` field still supported (deprecated in favor of `sourceHandle`)

## Best Practices

### 1. Use Conditional Execution for Mutually Exclusive Paths
✅ **Good**: Different API calls based on user type
```
Condition → [True: Premium API] or [False: Standard API]
```

### 2. Always Provide Both Paths
✅ **Good**: Handle both true and false cases
❌ **Avoid**: Leaving one path unhandled (leads to incomplete workflows)

### 3. Use Switch for 3+ Options
✅ **Good**: HTTP status routing with switch
❌ **Avoid**: Nested if-else for multiple options

### 4. Converge Paths When Needed
Multiple paths can converge to a single node:
```
[True Path] ──┐
              ├──→ [Final Step]
[False Path] ─┘
```

### 5. Label Your Paths Clearly
Use descriptive `output_path` names:
- ✅ "premium_user", "standard_user"
- ❌ "path1", "path2"

## Troubleshooting

**Q: My node doesn't execute even though it should**
- Check if incoming edge has correct `sourceHandle`
- Verify source node's output contains matching path field
- Ensure condition evaluates to expected result

**Q: Both paths execute (I expected only one)**
- Ensure edges have `sourceHandle` set
- Check for unconditional edges to the same node

**Q: No nodes execute after condition**
- Verify at least one edge's condition is satisfied
- Check condition expression syntax

## Performance Notes

- Conditional execution **improves performance** by skipping unnecessary nodes
- Use for expensive operations (API calls, heavy computations)
- Engine evaluates conditions in microseconds

## Related Documentation

- [Workflow Execution Guide](../../docs/WORKFLOW_EXECUTION_GUIDE.md)
- [Conditional Branching Testing Summary](../../docs/CONDITIONAL_BRANCHING_TESTING_SUMMARY.md)
- [Node Types Reference](../../docs/NODE_TYPES.md)
- [React Flow Documentation](https://reactflow.dev/learn/advanced-use/computing-flows)

---

**Total Examples**: 11 (8 original + 3 new conditional execution examples)  
**New Feature**: Conditional Execution (Path Termination) ✨  
**Coverage**: Basic branching, nested conditions, switch routing, real-world API flows

For more examples and tutorials, see the main [Examples Directory](../README.md).
