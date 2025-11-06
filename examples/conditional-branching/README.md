# Conditional Branching Examples

This directory contains comprehensive examples demonstrating conditional branching capabilities in Thaiyyal workflows.

## Overview

Conditional branching allows workflows to make decisions and route data based on conditions. Thaiyyal provides two primary nodes for conditional logic:

- **Condition Node**: Evaluates a boolean expression and returns metadata about the result
- **Switch Node**: Multi-way branching with pattern matching

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
