// Package expression provides a powerful expression evaluation engine for workflows.
//
// # Overview
//
// The expression package implements a domain-specific language (DSL) for dynamic
// evaluation of expressions within workflow nodes. It supports data access, operations,
// functions, and complex data transformations.
//
// # Features
//
//   - Field access: Navigate object hierarchies (user.profile.name)
//   - Array indexing: Access array elements (items[0], items[-1] for last)
//   - Operators: Arithmetic, comparison, logical, string operations
//   - Functions: Rich set of built-in functions (len, upper, lower, etc.)
//   - Type coercion: Automatic type conversion where appropriate
//   - Null safety: Graceful handling of null/undefined values
//   - Context variables: Access workflow-level variables
//
// # Expression Syntax
//
// Basic field access:
//
//	user.name           // Access field
//	user.profile.email  // Nested field access
//	items[0]            // Array index
//	items[-1]           // Last element
//	data.users[5].name  // Combined access
//
// Operators:
//
//	x + y               // Addition
//	x - y               // Subtraction
//	x * y               // Multiplication
//	x / y               // Division
//	x % y               // Modulo
//	x == y              // Equality
//	x != y              // Inequality
//	x > y, x < y        // Comparison
//	x >= y, x <= y      // Comparison
//	x && y              // Logical AND
//	x || y              // Logical OR
//	!x                  // Logical NOT
//
// String operations:
//
//	"Hello" + " " + "World"  // Concatenation
//	name + " (" + age + ")"  // Mixed types
//
// # Built-in Functions
//
// String functions:
//
//	upper(text)         // Convert to uppercase
//	lower(text)         // Convert to lowercase
//	trim(text)          // Remove whitespace
//	split(text, sep)    // Split into array
//	join(array, sep)    // Join array elements
//	replace(text, old, new)  // Replace substring
//	substr(text, start, length)  // Extract substring
//	contains(text, substr)  // Check if contains
//	startsWith(text, prefix)  // Check prefix
//	endsWith(text, suffix)  // Check suffix
//
// Array functions:
//
//	len(array)          // Array length
//	first(array)        // First element
//	last(array)         // Last element
//	reverse(array)      // Reverse array
//	sort(array)         // Sort array
//	unique(array)       // Remove duplicates
//	flatten(array)      // Flatten nested arrays
//	concat(arr1, arr2)  // Concatenate arrays
//
// Math functions:
//
//	abs(x)              // Absolute value
//	ceil(x)             // Ceiling
//	floor(x)            // Floor
//	round(x)            // Round to nearest
//	min(x, y)           // Minimum
//	max(x, y)           // Maximum
//	sqrt(x)             // Square root
//	pow(x, y)           // Power
//
// Type functions:
//
//	typeof(value)       // Get type name
//	isNull(value)       // Check if null
//	isNumber(value)     // Check if number
//	isString(value)     // Check if string
//	isArray(value)      // Check if array
//	isObject(value)     // Check if object
//
// Date/Time functions:
//
//	now()               // Current timestamp
//	date(timestamp)     // Format date
//	parseDate(string)   // Parse date string
//
// # Usage Examples
//
// Simple evaluation:
//
//	evaluator := expression.NewEvaluator()
//	result, err := evaluator.Evaluate("user.age >= 18", map[string]interface{}{
//	    "user": map[string]interface{}{
//	        "age": 25,
//	    },
//	})
//	// result: true
//
// Complex expressions:
//
//	expr := "upper(user.name) + ' is ' + (user.age >= 18 ? 'adult' : 'minor')"
//	result, err := evaluator.Evaluate(expr, context)
//
// Array filtering:
//
//	expr := "item.price > 100 && item.category == 'electronics'"
//	result, err := evaluator.Evaluate(expr, map[string]interface{}{
//	    "item": product,
//	})
//
// # Context Variables
//
// Expressions can access workflow-level context:
//
//	result, err := evaluator.EvaluateWithContext(
//	    "config.maxRetries * 2",
//	    data,
//	    map[string]interface{}{
//	        "config": workflowConfig,
//	    },
//	)
//
// # Type System
//
// The expression evaluator supports these types:
//
//   - Number: int, int64, float64
//   - String: string
//   - Boolean: bool
//   - Array: []interface{}
//   - Object: map[string]interface{}
//   - Null: nil
//
// Type coercion rules:
//
//   - Numbers are coerced for arithmetic operations
//   - Strings concatenate with + operator
//   - Booleans for logical operations
//   - Automatic conversion for comparisons
//
// # Error Handling
//
// The evaluator provides detailed error messages:
//
//   - Syntax errors: Invalid expression syntax
//   - Type errors: Incompatible types for operation
//   - Reference errors: Undefined field or variable
//   - Function errors: Invalid function arguments
//
// # Performance
//
//   - Expression parsing is optimized for common patterns
//   - Results can be cached for repeated evaluations
//   - Field access is optimized with reflection caching
//   - Function calls are fast-pathed for built-ins
//
// # Security Considerations
//
//   - No code execution: Expressions cannot execute arbitrary code
//   - Safe evaluation: No access to system resources
//   - Memory limits: Protection against resource exhaustion
//   - No recursion: Prevents stack overflow
//
// # Extension Points
//
// Custom functions can be registered:
//
//	evaluator.RegisterFunction("myFunc", func(args ...interface{}) (interface{}, error) {
//	    // Implementation
//	})
//
// # Thread Safety
//
// The Evaluator is safe for concurrent use by multiple goroutines.
// Internal state is protected with appropriate synchronization.
package expression
