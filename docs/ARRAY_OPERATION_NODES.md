# Array Operation Nodes

This document provides comprehensive documentation for all array operation nodes in Thaiyyal. These nodes enable powerful data transformation and manipulation capabilities for array-based workflows.

## Overview

The Array Operations category includes 15 specialized nodes for working with arrays and collections. These nodes are inspired by functional programming concepts and JavaScript's array methods, providing a familiar and powerful toolkit for data processing.

## Available Nodes

### 1. Map Node

**Purpose**: Transform each element in an array using an expression.

**Equivalent to**: JavaScript's `Array.map()`

**Configuration**:
- **Expression**: Transformation expression applied to each element
- Default: `item * 2`

**Input**: Array of any type

**Output**: Transformed array with same length

**Examples**:
```javascript
// Input: [1, 2, 3, 4, 5]
// Expression: item * 2
// Output: [2, 4, 6, 8, 10]

// Input: [{name: "Alice", age: 25}, {name: "Bob", age: 30}]
// Expression: item.name
// Output: ["Alice", "Bob"]

// Input: [10, 20, 30]
// Expression: item / 10
// Output: [1, 2, 3]
```

**Use Cases**:
- Transforming data structures
- Extracting specific fields from objects
- Applying mathematical operations to all elements
- Data normalization

---

### 2. Reduce Node

**Purpose**: Reduce an array to a single value using an accumulator expression.

**Equivalent to**: JavaScript's `Array.reduce()`

**Configuration**:
- **Expression**: Accumulator expression (e.g., `acc + item`)
- **Initial Value**: Starting value for the accumulator
- Default Expression: `acc + item`
- Default Initial Value: `0`

**Input**: Array of any type

**Output**: Single accumulated value

**Examples**:
```javascript
// Sum of numbers
// Input: [1, 2, 3, 4, 5]
// Expression: acc + item, Initial: 0
// Output: 15

// Concatenate strings
// Input: ["Hello", " ", "World"]
// Expression: acc + item, Initial: ""
// Output: "Hello World"

// Find maximum
// Input: [10, 5, 20, 15]
// Expression: item > acc ? item : acc, Initial: 0
// Output: 20
```

**Use Cases**:
- Summing numerical values
- Concatenating strings
- Finding min/max values
- Building aggregated objects

---

### 3. Slice Node

**Purpose**: Extract a portion of an array by start and end indices.

**Equivalent to**: JavaScript's `Array.slice()`

**Configuration**:
- **Start Index**: Beginning index (inclusive)
- **End Index**: Ending index (exclusive), -1 for end of array
- Default Start: `0`
- Default End: `-1`

**Input**: Array of any type

**Output**: Sliced array

**Examples**:
```javascript
// Get first 3 elements
// Input: [1, 2, 3, 4, 5]
// Start: 0, End: 3
// Output: [1, 2, 3]

// Get last 2 elements
// Input: [1, 2, 3, 4, 5]
// Start: -2, End: -1
// Output: [4, 5]

// Get middle elements
// Input: [1, 2, 3, 4, 5]
// Start: 1, End: 4
// Output: [2, 3, 4]
```

**Use Cases**:
- Pagination
- Getting first/last N elements
- Removing elements from start/end
- Data windowing

---

### 4. Sort Node

**Purpose**: Sort array elements by a specified field.

**Equivalent to**: JavaScript's `Array.sort()` with comparator

**Configuration**:
- **Field**: Field to sort by (leave empty for primitive arrays)
- **Order**: Ascending (`asc`) or Descending (`desc`)
- Default Field: `""`
- Default Order: `asc`

**Input**: Array of objects or primitives

**Output**: Sorted array

**Examples**:
```javascript
// Sort numbers ascending
// Input: [5, 2, 8, 1, 9]
// Field: "", Order: asc
// Output: [1, 2, 5, 8, 9]

// Sort objects by age
// Input: [{name: "Alice", age: 30}, {name: "Bob", age: 25}]
// Field: "age", Order: asc
// Output: [{name: "Bob", age: 25}, {name: "Alice", age: 30}]

// Sort descending
// Input: [5, 2, 8, 1, 9]
// Field: "", Order: desc
// Output: [9, 8, 5, 2, 1]
```

**Use Cases**:
- Ordering results
- Finding top N items
- Ranking data
- Alphabetical sorting

---

### 5. Find Node

**Purpose**: Find the first element matching a given expression.

**Equivalent to**: JavaScript's `Array.find()`

**Configuration**:
- **Expression**: Boolean expression to match
- Default: `item.id == 1`

**Input**: Array of any type

**Output**: First matching element or null

**Examples**:
```javascript
// Find by ID
// Input: [{id: 1, name: "Alice"}, {id: 2, name: "Bob"}]
// Expression: item.id == 2
// Output: {id: 2, name: "Bob"}

// Find by condition
// Input: [10, 20, 30, 40]
// Expression: item > 25
// Output: 30

// Not found
// Input: [1, 2, 3]
// Expression: item > 10
// Output: null
```

**Use Cases**:
- Searching for specific items
- Lookup operations
- Validation checks
- First match selection

---

### 6. FlatMap Node

**Purpose**: Map each element to an array and flatten the result by one level.

**Equivalent to**: JavaScript's `Array.flatMap()`

**Configuration**:
- **Expression**: Expression that returns an array for each element
- Default: `item.values`

**Input**: Array of any type

**Output**: Flattened array

**Examples**:
```javascript
// Flatten nested arrays
// Input: [[1, 2], [3, 4], [5, 6]]
// Expression: item
// Output: [1, 2, 3, 4, 5, 6]

// Extract and flatten properties
// Input: [{values: [1, 2]}, {values: [3, 4]}]
// Expression: item.values
// Output: [1, 2, 3, 4]

// Duplicate and flatten
// Input: [1, 2, 3]
// Expression: [item, item * 2]
// Output: [1, 2, 2, 4, 3, 6]
```

**Use Cases**:
- Flattening nested structures
- Combining multiple arrays
- Expanding data structures
- Data denormalization

---

### 7. Group By Node

**Purpose**: Group array elements by a specified field.

**Equivalent to**: Lodash's `_.groupBy()`

**Configuration**:
- **Key Field**: Field to group by
- Default: `category`

**Input**: Array of objects

**Output**: Object with grouped arrays

**Examples**:
```javascript
// Group by category
// Input: [
//   {name: "Apple", category: "fruit"},
//   {name: "Carrot", category: "vegetable"},
//   {name: "Banana", category: "fruit"}
// ]
// Key Field: category
// Output: {
//   fruit: [{name: "Apple", ...}, {name: "Banana", ...}],
//   vegetable: [{name: "Carrot", ...}]
// }

// Group by status
// Input: [{id: 1, status: "active"}, {id: 2, status: "inactive"}]
// Key Field: status
// Output: {active: [{id: 1, ...}], inactive: [{id: 2, ...}]}
```

**Use Cases**:
- Data aggregation
- Creating lookup tables
- Categorizing data
- Building hierarchies

---

### 8. Unique Node

**Purpose**: Remove duplicate elements from an array.

**Equivalent to**: Lodash's `_.uniq()` or `_.uniqBy()`

**Configuration**:
- **By Field**: Optional field for uniqueness comparison
- Default: `""` (compare entire element)

**Input**: Array of any type

**Output**: Array with unique elements

**Examples**:
```javascript
// Remove duplicate numbers
// Input: [1, 2, 2, 3, 3, 3, 4]
// By Field: ""
// Output: [1, 2, 3, 4]

// Unique by ID
// Input: [{id: 1, name: "A"}, {id: 2, name: "B"}, {id: 1, name: "C"}]
// By Field: "id"
// Output: [{id: 1, name: "A"}, {id: 2, name: "B"}]

// Unique strings
// Input: ["apple", "banana", "apple", "cherry"]
// By Field: ""
// Output: ["apple", "banana", "cherry"]
```

**Use Cases**:
- Deduplication
- Set operations
- Removing redundant data
- Data cleansing

---

### 9. Chunk Node

**Purpose**: Split an array into smaller arrays of specified size.

**Equivalent to**: Lodash's `_.chunk()`

**Configuration**:
- **Size**: Number of elements per chunk
- Default: `3`

**Input**: Array of any type

**Output**: Array of arrays (chunks)

**Examples**:
```javascript
// Chunk into groups of 3
// Input: [1, 2, 3, 4, 5, 6, 7, 8]
// Size: 3
// Output: [[1, 2, 3], [4, 5, 6], [7, 8]]

// Pair elements
// Input: [1, 2, 3, 4]
// Size: 2
// Output: [[1, 2], [3, 4]]

// Single element chunks
// Input: [1, 2, 3]
// Size: 1
// Output: [[1], [2], [3]]
```

**Use Cases**:
- Batch processing
- Pagination
- Creating rows/columns
- Data partitioning

---

### 10. Reverse Node

**Purpose**: Reverse the order of elements in an array.

**Equivalent to**: JavaScript's `Array.reverse()`

**Configuration**: None

**Input**: Array of any type

**Output**: Reversed array

**Examples**:
```javascript
// Reverse numbers
// Input: [1, 2, 3, 4, 5]
// Output: [5, 4, 3, 2, 1]

// Reverse strings
// Input: ["first", "second", "third"]
// Output: ["third", "second", "first"]

// Reverse objects
// Input: [{id: 1}, {id: 2}, {id: 3}]
// Output: [{id: 3}, {id: 2}, {id: 1}]
```

**Use Cases**:
- Changing sort order
- Stack operations (LIFO)
- Timeline reversal
- Data presentation

---

### 11. Partition Node

**Purpose**: Split an array into two groups based on a condition.

**Equivalent to**: Lodash's `_.partition()`

**Configuration**:
- **Expression**: Boolean expression for partitioning
- Default: `item > 0`

**Input**: Array of any type

**Output**: Two arrays - matching and non-matching elements

**Examples**:
```javascript
// Partition positive/negative
// Input: [1, -2, 3, -4, 5]
// Expression: item > 0
// Output: {
//   true: [1, 3, 5],
//   false: [-2, -4]
// }

// Partition by age
// Input: [{name: "Alice", age: 25}, {name: "Bob", age: 17}]
// Expression: item.age >= 18
// Output: {
//   true: [{name: "Alice", age: 25}],
//   false: [{name: "Bob", age: 17}]
// }
```

**Use Cases**:
- Filtering with both results
- Categorizing data
- Validation with rejected items
- Data quality checks

---

### 12. Zip Node

**Purpose**: Combine two arrays element-wise into pairs.

**Equivalent to**: Lodash's `_.zip()`

**Configuration**: None (accepts two input arrays)

**Input**: Two arrays

**Output**: Array of paired elements

**Examples**:
```javascript
// Zip two arrays
// Array1: [1, 2, 3]
// Array2: ['a', 'b', 'c']
// Output: [[1, 'a'], [2, 'b'], [3, 'c']]

// Zip different lengths
// Array1: [1, 2, 3]
// Array2: ['a', 'b']
// Output: [[1, 'a'], [2, 'b'], [3, undefined]]

// Combine coordinates
// Array1: [10, 20, 30]
// Array2: [100, 200, 300]
// Output: [[10, 100], [20, 200], [30, 300]]
```

**Use Cases**:
- Combining related data
- Creating key-value pairs
- Merging parallel arrays
- Coordinate transformation

---

### 13. Sample Node

**Purpose**: Randomly sample a specified number of elements from an array.

**Equivalent to**: Lodash's `_.sampleSize()`

**Configuration**:
- **Count**: Number of elements to sample
- Default: `1`

**Input**: Array of any type

**Output**: Array with randomly selected elements

**Examples**:
```javascript
// Sample 3 random elements
// Input: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
// Count: 3
// Output: [7, 2, 9] (random)

// Sample single element
// Input: ['apple', 'banana', 'cherry']
// Count: 1
// Output: ['banana'] (random)

// Sample all (shuffle)
// Input: [1, 2, 3]
// Count: 3
// Output: [2, 3, 1] (random order)
```

**Use Cases**:
- Random selection
- A/B testing
- Data sampling
- Shuffling

---

### 14. Range Node

**Purpose**: Generate an array of numbers from start to end with optional step.

**Equivalent to**: Lodash's `_.range()`

**Configuration**:
- **Start**: Starting number
- **End**: Ending number (exclusive)
- **Step**: Increment value
- Default Start: `0`
- Default End: `10`
- Default Step: `1`

**Input**: None (generates output)

**Output**: Array of numbers

**Examples**:
```javascript
// Generate 0 to 9
// Start: 0, End: 10, Step: 1
// Output: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

// Generate even numbers
// Start: 0, End: 10, Step: 2
// Output: [0, 2, 4, 6, 8]

// Generate countdown
// Start: 10, End: 0, Step: -1
// Output: [10, 9, 8, 7, 6, 5, 4, 3, 2, 1]
```

**Use Cases**:
- Generating sequences
- Loop iterations
- Index arrays
- Test data generation

---

### 15. Transpose Node

**Purpose**: Transpose a 2D array (matrix), swapping rows and columns.

**Equivalent to**: Matrix transpose operation

**Configuration**: None

**Input**: 2D array (array of arrays)

**Output**: Transposed 2D array

**Examples**:
```javascript
// Transpose matrix
// Input: [
//   [1, 2, 3],
//   [4, 5, 6]
// ]
// Output: [
//   [1, 4],
//   [2, 5],
//   [3, 6]
// ]

// Transpose single row
// Input: [[1, 2, 3]]
// Output: [[1], [2], [3]]

// Transpose square matrix
// Input: [
//   [1, 2],
//   [3, 4]
// ]
// Output: [
//   [1, 3],
//   [2, 4]
// ]
```

**Use Cases**:
- Matrix operations
- Data pivoting
- Row/column conversion
- Spreadsheet operations

---

## Node Styling

Each array operation node has a unique color gradient for easy visual identification:

| Node | Color Gradient |
|------|----------------|
| Map | Cyan (cyan-600 to cyan-700) |
| Reduce | Teal (teal-600 to teal-700) |
| Slice | Emerald (emerald-600 to emerald-700) |
| Sort | Lime (lime-600 to lime-700) |
| Find | Sky (sky-600 to sky-700) |
| FlatMap | Indigo (indigo-600 to indigo-700) |
| Group By | Violet (violet-600 to violet-700) |
| Unique | Fuchsia (fuchsia-600 to fuchsia-700) |
| Chunk | Pink (pink-600 to pink-700) |
| Reverse | Rose (rose-600 to rose-700) |
| Partition | Orange (orange-600 to orange-700) |
| Zip | Yellow (yellow-600 to yellow-700) |
| Sample | Blue (blue-600 to blue-700) |
| Range | Green (green-600 to green-700) |
| Transpose | Red (red-600 to red-700) |

## Common Workflows

### Example 1: Data Processing Pipeline

```
Range (0-100) → Map (square) → Filter (>50) → Sort (desc) → Slice (0,10)
```

This workflow:
1. Generates numbers 0-99
2. Squares each number
3. Filters values greater than 50
4. Sorts in descending order
5. Takes top 10 results

### Example 2: Data Aggregation

```
HTTP (fetch users) → Group By (department) → Map (count per group) → Reduce (sum)
```

This workflow:
1. Fetches user data from API
2. Groups users by department
3. Counts users in each department
4. Sums total users

### Example 3: Data Transformation

```
Input Array → FlatMap (expand) → Unique (deduplicate) → Sort (asc) → Chunk (size 10)
```

This workflow:
1. Starts with nested data
2. Flattens to single level
3. Removes duplicates
4. Sorts alphabetically
5. Splits into pages of 10

## Best Practices

1. **Use Map for Transformations**: When you need to transform every element in an array
2. **Use Filter Before Map**: Filter out unwanted data before expensive transformations
3. **Use Reduce for Aggregations**: When you need a single value from an array
4. **Use Find for Single Items**: More efficient than filter when you only need the first match
5. **Chain Operations**: Combine multiple array nodes for complex data pipelines
6. **Consider Performance**: Large arrays may benefit from chunking or limiting operations

## Performance Considerations

- **Map/Filter/Reduce**: O(n) time complexity
- **Sort**: O(n log n) time complexity
- **Find**: O(n) worst case, O(1) best case
- **Group By**: O(n) time complexity
- **Unique**: O(n) time complexity with Set
- **Chunk/Slice/Reverse**: O(n) time complexity
- **Transpose**: O(n*m) for n×m matrix

## See Also

- [Control Flow Nodes](CONTROL_FLOW_NODES.md)
- [Node Types](NODE_TYPES.md)
- [Examples](EXAMPLES.md)
