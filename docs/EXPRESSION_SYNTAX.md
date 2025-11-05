# Expression Syntax Guide

## Overview

The Thaiyyal expression engine provides a powerful, JavaScript-like syntax for evaluating conditions, transforming data, and performing calculations in workflows. This guide covers all supported features and syntax.

## Table of Contents

1. [Basic Values](#basic-values)
2. [Field Access](#field-access)
3. [Array Operations](#array-operations)
4. [String Operations](#string-operations)
5. [Arithmetic Operations](#arithmetic-operations)
6. [Comparison Operators](#comparison-operators)
7. [Logical Operators](#logical-operators)
8. [Method Calls](#method-calls)
9. [Variable References](#variable-references)
10. [Examples](#examples)

## Basic Values

### Literals

```javascript
// Numbers
42
3.14
-10

// Strings  
"hello"
'world'

// Booleans
true
false

// Null
null
```

### Input Reference

Access the current input value using `input` or `item`:

```javascript
// In filter expressions
item > 10
item.age >= 18

// In map expressions
item * 2
item.price * 1.1
```

## Field Access

Access nested object fields using dot notation:

```javascript
// Simple field
item.name

// Nested fields
item.user.profile.verified
item.metadata.tags.primary

// Deep nesting
item.data.results.items.first.value
```

## Array Operations

### Array Length

Get the length of an array or string:

```javascript
// Array length
item.tags.length == 3
item.users.length > 0

// String length
item.name.length > 5
item.email.length >= 10
```

### Array Indexing

Access array elements by index (0-based):

```javascript
// Simple indexing
item.tags[0]
item.users[1]

// Nested indexing
item.users[0].name
item.data[2].value

// Multiple levels
item.results[0].items[1].price
```

### Array Methods

```javascript
// Check if array includes a value
item.tags.includes('admin')
item.roles.includes('moderator')

// Get first/last element
item.list.first()
item.history.last()

// Join array into string
item.tags.join(', ')
item.items.join(' | ')

// Reverse array
item.numbers.reverse()
```

## String Operations

### String Methods

```javascript
// Convert case
item.name.toUpperCase()
item.email.toLowerCase()

// Check string contents
item.email.includes('@example.com')
item.filename.startsWith('report')
item.filename.endsWith('.pdf')

// Trim whitespace
item.text.trim()

// Replace text
item.message.replace('old', 'new')

// Split string into array
item.csv.split(',')
item.text.split(' ')
```

### Method Chaining

Combine multiple operations:

```javascript
// Trim and convert to uppercase
item.name.trim().toUpperCase()

// Lowercase and check
item.email.toLowerCase().includes('@test.com')

// Multiple transformations
item.text.trim().toLowerCase().replace('hello', 'hi')
```

## Arithmetic Operations

### Basic Arithmetic

```javascript
// Addition and subtraction
item.price + 10
item.quantity - 5

// Multiplication and division
item.price * 1.1
item.total / item.count

// Modulo
item.value % 2

// Parentheses for grouping
(item.price + item.tax) * item.quantity
```

### Math Functions

```javascript
// Power
pow(item.value, 2)

// Square root
sqrt(item.area)

// Absolute value
abs(item.difference)

// Rounding
floor(item.value)
ceil(item.value)
round(item.value)

// Min/max
min(item.a, item.b, item.c)
max(item.x, item.y)
```

### With Array Properties

```javascript
// Use array length in calculations
item.items.length * 2
item.tags.length + item.categories.length

// Complex expressions
(item.scores[0] + item.scores[1]) / 2
```

## Comparison Operators

```javascript
// Equality
item.status == 'active'
item.count != 0

// Numeric comparison
item.age > 18
item.price >= 100
item.stock < 10
item.score <= 50

// String comparison
item.name == 'Alice'
item.role != 'guest'
```

## Logical Operators

### AND, OR, NOT

```javascript
// AND
item.age > 18 && item.verified == true
item.price > 100 && item.stock > 0

// OR
item.role == 'admin' || item.role == 'moderator'
item.status == 'active' || item.status == 'pending'

// NOT
!item.deleted
!(item.age < 18)

// Complex expressions
(item.age > 18 && item.verified) || item.role == 'admin'
```

## Method Calls

### String Methods

| Method | Description | Example |
|--------|-------------|---------|
| `toUpperCase()` | Convert to uppercase | `item.name.toUpperCase()` |
| `toLowerCase()` | Convert to lowercase | `item.email.toLowerCase()` |
| `trim()` | Remove whitespace | `item.text.trim()` |
| `includes(substring)` | Check if contains substring | `item.email.includes('@')` |
| `startsWith(prefix)` | Check if starts with | `item.filename.startsWith('data')` |
| `endsWith(suffix)` | Check if ends with | `item.filename.endsWith('.pdf')` |
| `replace(old, new)` | Replace text | `item.text.replace('a', 'b')` |
| `split(separator)` | Split into array | `item.csv.split(',')` |

### Array Methods

| Method | Description | Example |
|--------|-------------|---------|
| `includes(value)` | Check if array contains value | `item.tags.includes('important')` |
| `join(separator)` | Join array elements | `item.items.join(', ')` |
| `reverse()` | Reverse array order | `item.list.reverse()` |
| `first()` | Get first element | `item.items.first()` |
| `last()` | Get last element | `item.items.last()` |

## Variable References

Access workflow variables, node results, and context:

```javascript
// Variables
variables.count
variables.user.name
variables.items[0]
variables.config.maxRetries

// Node results
node.http1.response.status
node.filter1.filtered.length
node.map1.results[0]

// Context variables
context.apiKey
context.environment
context.maxItems
```

## Examples

### Filtering Examples

```javascript
// Filter users by age
item.age >= 18

// Filter active items
item.status == 'active' && !item.deleted

// Filter by tag
item.tags.includes('important')

// Filter by email domain
item.email.toLowerCase().includes('@company.com')

// Complex filter
item.score > 80 && item.verified && item.tags.length > 0
```

### Mapping Examples

```javascript
// Double values
item * 2

// Calculate total with tax
item.price * 1.1

// Extract field
item.user.name

// Transform string
item.name.toUpperCase()

// Complex transformation
{
  name: item.name.trim(),
  email: item.email.toLowerCase(),
  score: item.score * 1.5
}
```

### Conditional Examples

```javascript
// Check age requirement
item.age > 18

// Verify permissions
item.role == 'admin' || item.role == 'moderator'

// Check array conditions
item.tags.length > 0 && item.tags.includes('verified')

// String validation
item.email.includes('@') && item.email.endsWith('.com')

// Multiple conditions
item.active && item.verified && item.score >= 75
```

### Reduce Examples

```javascript
// Sum values
accumulator + item

// Calculate average
accumulator + item.score

// Concatenate strings
accumulator + item.name + ', '

// Complex accumulation
{
  total: accumulator.total + item.price,
  count: accumulator.count + 1,
  items: accumulator.items.concat(item)
}
```

## Best Practices

### 1. Use Explicit References

Prefer `item.field` over just `field` for clarity:

```javascript
// Good
item.age > 18

// Works but less clear
age > 18
```

### 2. Handle Edge Cases

Check for empty arrays and null values:

```javascript
// Check length before indexing
item.tags.length > 0 && item.tags[0] == 'important'

// Use method chaining safely
item.email.toLowerCase().includes('@')
```

### 3. Keep Expressions Simple

Break complex logic into multiple nodes:

```javascript
// Good - simple and clear
item.verified && item.age >= 18

// Too complex - consider splitting
item.verified && item.age >= 18 && item.tags.includes('premium') && item.score > 80 && !item.suspended
```

### 4. Use Type-Appropriate Methods

Match methods to data types:

```javascript
// String methods on strings
item.name.toUpperCase()

// Array methods on arrays
item.tags.includes('admin')

// Number operations on numbers
item.price * 1.1
```

## Common Patterns

### Email Validation

```javascript
item.email.includes('@') && item.email.includes('.')
```

### Role-Based Access

```javascript
item.role == 'admin' || item.role == 'moderator'
```

### Active Status Check

```javascript
item.status == 'active' && !item.deleted
```

### Price Range Filter

```javascript
item.price >= 10 && item.price <= 100
```

### Tag-Based Filtering

```javascript
item.tags.includes('featured') || item.tags.includes('trending')
```

### Name Normalization

```javascript
item.name.trim().toLowerCase()
```

### File Type Validation

```javascript
item.filename.endsWith('.pdf') || item.filename.endsWith('.doc')
```

## Error Handling

Expressions that fail to evaluate will:
- Return `false` for boolean conditions (filter)
- Return `nil` for value expressions (map)
- Log errors for debugging

Common errors:
- Field not found: accessing non-existent fields
- Type mismatch: calling string methods on numbers
- Index out of bounds: accessing invalid array indexes

## Performance Tips

1. **Simple expressions are faster** - avoid unnecessary complexity
2. **Method chaining is efficient** - operations are lazy where possible
3. **Array indexing is fast** - O(1) access time
4. **Use built-in methods** - optimized implementations

## Limitations

- No custom function definitions
- No loops within expressions (use ForEach node)
- No assignments (use Variable/Extract nodes)
- Limited to single expression per node

## See Also

- [Node Types Documentation](./NODE_TYPES.md)
- [Workflow Examples](../src/data/workflowExamples.ts)
- [Expression Engine Implementation](../backend/pkg/expression/)

## Version

This documentation applies to expression engine version 2.0+ (with string/array methods support).

---

For questions or feature requests, please open an issue on GitHub.
