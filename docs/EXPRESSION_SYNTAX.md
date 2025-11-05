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
9. [Value Functions](#value-functions)
10. [Variable References](#variable-references)
11. [Examples](#examples)

## Basic Values

### Literals

```javascript
// Numbers
42;
3.14 - 10;

// Strings
("hello");
("world");

// Booleans
true;
false;

// Null
null;
```

### Input Reference

Access the current input value using `input` or `item`:

```javascript
// In filter expressions
item > 10;
item.age >= 18;

// In map expressions
item * 2;
item.price * 1.1;
```

## Field Access

Access nested object fields using dot notation:

```javascript
// Simple field
item.name;

// Nested fields
item.user.profile.verified;
item.metadata.tags.primary;

// Deep nesting
item.data.results.items.first.value;
```

## Array Operations

### Array Length

Get the length of an array or string:

```javascript
// Array length
item.tags.length == 3;
item.users.length > 0;

// String length
item.name.length > 5;
item.email.length >= 10;
```

### Array Indexing

Access array elements by index (0-based):

```javascript
// Simple indexing
item.tags[0];
item.users[1];

// Nested indexing
item.users[0].name;
item.data[2].value;

// Multiple levels
item.results[0].items[1].price;
```

### Array Methods

```javascript
// Check if array includes a value
item.tags.includes("admin");
item.roles.includes("moderator");

// Get first/last element
item.list.first();
item.history.last();

// Join array into string
item.tags.join(", ");
item.items.join(" | ");

// Reverse array
item.numbers.reverse();
```

### Value Functions (map, avg)

- map(arrayExpr, itemExpr): Project each element of an array using an expression evaluated with item bound to the element. The first argument must resolve to an array. The second argument is an expression using item (e.g., item.age, item.price \* 1.1).
- avg(arrayExpr) or avg(v1, v2, ...): Compute the arithmetic mean of numeric values. With a single argument that resolves to an array, all elements must be numeric. With multiple arguments, each is evaluated and must be numeric.

Examples:

```javascript
// Map users to their ages
map(item.users, item.age); // -> [30, 25, 40]

// Average of mapped ages
avg(map(item.users, item.age)); // -> 31.666...

// Average of constants/expressions
avg(10, 20, 30); // -> 20
```

Notes:

- If avg receives an empty array, it returns an error.
- Non-numeric elements in avg() will cause an error.

## String Operations

### String Methods

```javascript
// Convert case
item.name.toUpperCase();
item.email.toLowerCase();

// Check string contents
item.email.includes("@example.com");
item.filename.startsWith("report");
item.filename.endsWith(".pdf");

// Trim whitespace
item.text.trim();

// Replace text
item.message.replace("old", "new");

// Split string into array
item.csv.split(",");
item.text.split(" ");
```

### Method Chaining

Combine multiple operations:

```javascript
// Trim and convert to uppercase
item.name.trim().toUpperCase();

// Lowercase and check
item.email.toLowerCase().includes("@test.com");

// Multiple transformations
item.text.trim().toLowerCase().replace("hello", "hi");
```

## Arithmetic Operations

### Basic Arithmetic

```javascript
// Addition and subtraction
item.price + 10;
item.quantity - 5;

// Multiplication and division
item.price * 1.1;
item.total / item.count;

// Modulo
(item.value %
  2(
    // Parentheses for grouping
    item.price + item.tax
  )) *
  item.quantity;
```

### Math Functions

```javascript
// Power
pow(item.value, 2);

// Square root
sqrt(item.area);

// Absolute value
abs(item.difference);

// Rounding
floor(item.value);
ceil(item.value);
round(item.value);

// Min/max
min(item.a, item.b, item.c);
max(item.x, item.y);
```

### With Array Properties

```javascript
// Use array length in calculations
item.items.length * 2;
item.tags.length +
  item.categories.length(
    // Complex expressions
    item.scores[0] + item.scores[1]
  ) /
    2;
```

## Comparison Operators

```javascript
// Equality
item.status == "active";
item.count != 0;

// Numeric comparison
item.age > 18;
item.price >= 100;
item.stock < 10;
item.score <= 50;

// String comparison
item.name == "Alice";
item.role != "guest";
```

## Logical Operators

### AND, OR, NOT

```javascript
// AND
item.age > 18 && item.verified == true;
item.price > 100 && item.stock > 0;

// OR
item.role == "admin" || item.role == "moderator";
item.status == "active" || item.status == "pending";

// NOT
!item.deleted;
!(item.age < 18)(
  // Complex expressions
  item.age > 18 && item.verified
) || item.role == "admin";
```

## Method Calls

### String Method Reference

| Method                | Description                 | Example                            |
| --------------------- | --------------------------- | ---------------------------------- |
| `toUpperCase()`       | Convert to uppercase        | `item.name.toUpperCase()`          |
| `toLowerCase()`       | Convert to lowercase        | `item.email.toLowerCase()`         |
| `trim()`              | Remove whitespace           | `item.text.trim()`                 |
| `includes(substring)` | Check if contains substring | `item.email.includes('@')`         |
| `startsWith(prefix)`  | Check if starts with        | `item.filename.startsWith('data')` |
| `endsWith(suffix)`    | Check if ends with          | `item.filename.endsWith('.pdf')`   |
| `replace(old, new)`   | Replace text                | `item.text.replace('a', 'b')`      |
| `split(separator)`    | Split into array            | `item.csv.split(',')`              |

### Array Method Reference

| Method            | Description                   | Example                           |
| ----------------- | ----------------------------- | --------------------------------- |
| `includes(value)` | Check if array contains value | `item.tags.includes('important')` |
| `join(separator)` | Join array elements           | `item.items.join(', ')`           |
| `reverse()`       | Reverse array order           | `item.list.reverse()`             |
| `first()`         | Get first element             | `item.items.first()`              |
| `last()`          | Get last element              | `item.items.last()`               |

## Value Functions

Value functions are functions that return computed values and can be used in expressions. They are particularly useful for array processing, aggregation, and mathematical transformations.

### Array Transformation Functions

#### `map(array, expression)`

Projects each element of an array through an expression, returning a new array.

```javascript
// Extract ages from user objects
map(input, item.age); // [25, 30, 35]

// Calculate prices with tax
map(input, item.price * 1.1); // [11.0, 22.0, 33.0]

// Extract nested field
map(input.users, item.profile.name); // ["Alice", "Bob", "Charlie"]

// Transform strings
map(input.names, item.toUpperCase()); // ["ALICE", "BOB", "CHARLIE"]
```

**Notes:**

- First argument must be an array
- Second argument is an expression evaluated for each element
- Use `item` to reference the current element
- Returns a new array with transformed values

### Aggregate Functions

#### `avg(array)` or `avg(value1, value2, ...)`

Calculates the average (mean) of numeric values.

```javascript
// Average of array
avg([1, 2, 3, 4, 5]); // 3.0

// Average of multiple arguments
avg(10, 20, 30); // 20.0

// Combined with map
avg(map(input.users, item.age)); // Average age

// Empty array returns 0
avg([]); // 0
```

#### `sum(array)` or `sum(value1, value2, ...)`

Calculates the sum of numeric values.

```javascript
// Sum of array
sum([1, 2, 3, 4, 5]); // 15

// Sum of multiple arguments
sum(10, 20, 30); // 60

// Combined with map
sum(map(input.items, item.price)); // Total price

// Empty array returns 0
sum([]); // 0
```

#### `min(array)` or `min(value1, value2, ...)`

Finds the minimum value.

```javascript
// Min of array
min([5, 2, 8, 1, 9]); // 1

// Min of multiple arguments
min(10, 20, 5, 30); // 5

// Combined with map
min(map(input.products, item.price)); // Cheapest price
```

#### `max(array)` or `max(value1, value2, ...)`

Finds the maximum value.

```javascript
// Max of array
max([5, 2, 8, 1, 9]); // 9

// Max of multiple arguments
max(10, 20, 5, 30); // 30

// Combined with map
max(map(input.students, item.score)); // Highest score
```

### Mathematical Functions

These functions accept either a single numeric value or an array of numbers. When given an array, they apply the function to each element.

#### `round(value)` or `round(array)`

Rounds to the nearest integer.

```javascript
// Single value
round(3.7); // 4
round(3.2); // 3

// Array of values
round([1.2, 2.7, 3.5]); // [1, 3, 4]

// Combined with map
round(map(input.prices, item * 1.15)); // Rounded prices with markup
```

#### `floor(value)` or `floor(array)`

Rounds down to the nearest integer.

```javascript
// Single value
floor(3.9); // 3

// Array of values
floor([1.9, 2.5, 3.1]); // [1, 2, 3]
```

#### `ceil(value)` or `ceil(array)`

Rounds up to the nearest integer.

```javascript
// Single value
ceil(3.1); // 4

// Array of values
ceil([1.1, 2.5, 3.9]); // [2, 3, 4]
```

#### `abs(value)` or `abs(array)`

Returns the absolute value.

```javascript
// Single value
abs(-5); // 5
abs(3); // 3

// Array of values
abs([-1, -2, 3, -4]); // [1, 2, 3, 4]
```

### Array Manipulation Functions

#### `sort(array)`

Sorts an array in ascending order. Works with numbers, strings, and mixed types.

```javascript
// Sort numbers
sort([3, 1, 4, 1, 5]); // [1, 1, 3, 4, 5]

// Sort strings
sort(["banana", "apple", "cherry"]); // ["apple", "banana", "cherry"]

// Sort mapped values
sort(map(input.users, item.age)); // Sorted ages
```

**Note:** Numbers sort before strings. Mixed-type arrays are sorted with numbers first.

#### `slice(array, start)` or `slice(array, start, end)`

Extracts a portion of an array.

```javascript
// Slice from index 2 to end
slice([1, 2, 3, 4, 5], 2); // [3, 4, 5]

// Slice from index 1 to 3 (exclusive)
slice([1, 2, 3, 4, 5], 1, 3); // [2, 3]

// Negative indices (from end)
slice([1, 2, 3, 4, 5], -2); // [4, 5]
slice([1, 2, 3, 4, 5], 0, -1); // [1, 2, 3, 4]

// Get first N elements
slice(input.items, 0, 10); // First 10 items
```

**Note:** Negative indices count from the end of the array.

#### `unique(array)`

Removes duplicate values from an array.

```javascript
// Remove duplicates
unique([1, 2, 2, 3, 1, 4]); // [1, 2, 3, 4]

// Unique strings
unique(["a", "b", "a", "c"]); // ["a", "b", "c"]

// Unique mapped values
unique(map(input.users, item.role)); // Unique roles
```

#### `reverse(array)`

Reverses the order of an array.

```javascript
// Reverse array
reverse([1, 2, 3, 4, 5]); // [5, 4, 3, 2, 1]

// Reverse sorted array (descending)
reverse(sort(input.scores)); // Highest to lowest
```

#### `flatten(array)`

Flattens nested arrays recursively.

```javascript
// Flatten nested arrays
flatten([
  [1, 2],
  [3, 4],
]); // [1, 2, 3, 4]

// Deep nesting
flatten([
  [1, [2, 3]],
  [4, [5, 6]],
]); // [1, 2, 3, 4, 5, 6]

// Mixed depths
flatten([1, [2, 3], 4, [5]]); // [1, 2, 3, 4, 5]
```

#### `zip(array1, array2, ...)`

Combines multiple arrays into an array of tuples.

```javascript
// Zip two arrays
zip([1, 2, 3], ["a", "b", "c"]); // [[1, "a"], [2, "b"], [3, "c"]]

// Zip three arrays
zip([1, 2], ["a", "b"], [true, false]); // [[1, "a", true], [2, "b", false]]

// Combine user data
zip(map(input.users, item.name), map(input.users, item.email));
```

**Note:** Zip stops at the length of the shortest array.

#### `sample(array, n)`

Returns a sample of n elements from an array.

```javascript
// Sample 3 elements
sample([1, 2, 3, 4, 5], 3); // [1, 2, 3] (simplified implementation)

// Get random subset
sample(input.products, 5); // 5 products
```

**Note:** Current implementation uses a simplified deterministic sampling approach.

### Complex Compositions

Value functions can be combined for powerful transformations:

```javascript
// Average age of users
avg(map(input.users, item.age));

// Total price rounded
round(sum(map(input.items, item.price)));

// Get top 5 scores
slice(reverse(sort(map(input.students, item.score))), 0, 5);

// Unique sorted roles
sort(unique(map(input.users, item.role)));

// Average of rounded prices with tax
avg(round(map(input.products, item.price * 1.15)));

// Min score, max score
min(map(input.tests, item.score));
max(map(input.tests, item.score));
```

## Variable References

Access workflow variables, node results, and context:

```javascript
// Variables
variables.count;
variables.user.name;
variables.items[0];
variables.config.maxRetries;

// Node results
node.http1.response.status;
node.filter1.filtered.length;
node.map1.results[0];

// Context variables
context.apiKey;
context.environment;
context.maxItems;
```

## Examples

### Filtering Examples

```javascript
// Filter users by age
item.age >= 18;

// Filter active items
item.status == "active" && !item.deleted;

// Filter by tag
item.tags.includes("important");

// Filter by email domain
item.email.toLowerCase().includes("@company.com");

// Complex filter
item.score > 80 && item.verified && item.tags.length > 0;
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
item.age > 18;

// Verify permissions
item.role == "admin" || item.role == "moderator";

// Check array conditions
item.tags.length > 0 && item.tags.includes("verified");

// String validation
item.email.includes("@") && item.email.endsWith(".com");

// Multiple conditions
item.active && item.verified && item.score >= 75;
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
item.age > 18;

// Works but less clear
age > 18;
```

### 2. Handle Edge Cases

Check for empty arrays and null values:

```javascript
// Check length before indexing
item.tags.length > 0 && item.tags[0] == "important";

// Use method chaining safely
item.email.toLowerCase().includes("@");
```

### 3. Keep Expressions Simple

Break complex logic into multiple nodes:

```javascript
// Good - simple and clear
item.verified && item.age >= 18;

// Too complex - consider splitting
item.verified &&
  item.age >= 18 &&
  item.tags.includes("premium") &&
  item.score > 80 &&
  !item.suspended;
```

### 4. Use Type-Appropriate Methods

Match methods to data types:

```javascript
// String methods on strings
item.name.toUpperCase();

// Array methods on arrays
item.tags.includes("admin");

// Number operations on numbers
item.price * 1.1;
```

## Common Patterns

### Email Validation

```javascript
item.email.includes("@") && item.email.includes(".");
```

### Role-Based Access

```javascript
item.role == "admin" || item.role == "moderator";
```

### Active Status Check

```javascript
item.status == "active" && !item.deleted;
```

### Price Range Filter

```javascript
item.price >= 10 && item.price <= 100;
```

### Tag-Based Filtering

```javascript
item.tags.includes("featured") || item.tags.includes("trending");
```

### Name Normalization

```javascript
item.name.trim().toLowerCase();
```

### File Type Validation

```javascript
item.filename.endsWith(".pdf") || item.filename.endsWith(".doc");
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
