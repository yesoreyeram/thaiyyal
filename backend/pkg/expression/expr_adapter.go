package expression

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// ExprEngine wraps expr-lang/expr for expression evaluation
type ExprEngine struct {
	programCache map[string]*vm.Program
}

// NewExprEngine creates a new expression engine using expr-lang/expr
func NewExprEngine() *ExprEngine {
	return &ExprEngine{
		programCache: make(map[string]*vm.Program),
	}
}

// EvaluateBoolean evaluates an expression and returns a boolean result
// This is the expr-lang/expr implementation of Evaluate()
func (e *ExprEngine) EvaluateBoolean(expression string, input interface{}, ctx *Context) (bool, error) {
	if ctx == nil {
		ctx = &Context{
			NodeResults: make(map[string]interface{}),
			Variables:   make(map[string]interface{}),
			ContextVars: make(map[string]interface{}),
		}
	}

	// Convert Thaiyyal syntax to expr-lang syntax
	expression = convertSyntax(expression)

	// Build environment with all context data
	env := e.buildEnvironment(input, ctx)

	// Try to get cached program (cache key includes converted expression)
	program, exists := e.programCache[expression]
	if !exists {
		// Compile the expression
		var err error
		program, err = expr.Compile(expression, expr.Env(env), expr.AsBool())
		if err != nil {
			return false, fmt.Errorf("expression compilation failed: %w", err)
		}
		e.programCache[expression] = program
	}

	// Execute the program
	output, err := expr.Run(program, env)
	if err != nil {
		return false, fmt.Errorf("expression execution failed: %w", err)
	}

	// Convert to boolean
	result, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf("expression did not return boolean, got %T", output)
	}

	return result, nil
}

// EvaluateValue evaluates an expression and returns its value
// This is the expr-lang/expr implementation of EvaluateExpression()
func (e *ExprEngine) EvaluateValue(expression string, input interface{}, ctx *Context) (interface{}, error) {
	if ctx == nil {
		ctx = &Context{
			NodeResults: make(map[string]interface{}),
			Variables:   make(map[string]interface{}),
			ContextVars: make(map[string]interface{}),
		}
	}

	// Convert Thaiyyal syntax to expr-lang syntax
	expression = convertSyntax(expression)

	// Build environment with all context data
	env := e.buildEnvironment(input, ctx)

	// Try to get cached program
	program, exists := e.programCache[expression]
	if !exists {
		// Compile the expression
		var err error
		program, err = expr.Compile(expression, expr.Env(env))
		if err != nil {
			return nil, fmt.Errorf("expression compilation failed: %w", err)
		}
		e.programCache[expression] = program
	}

	// Execute the program
	output, err := expr.Run(program, env)
	if err != nil {
		return nil, fmt.Errorf("expression execution failed: %w", err)
	}

	return output, nil
}

// buildEnvironment creates the execution environment with all variables and functions
func (e *ExprEngine) buildEnvironment(input interface{}, ctx *Context) map[string]interface{} {
	env := make(map[string]interface{})

	// Add custom functions
	e.addCustomFunctions(env)

	// Add node results
	if ctx.NodeResults != nil {
		env["node"] = ctx.NodeResults
	}

	// Add variables
	if ctx.Variables != nil {
		env["variables"] = ctx.Variables
		// Also add variables directly for backward compatibility
		for k, v := range ctx.Variables {
			if k != "node" && k != "variables" && k != "context" {
				env[k] = v
			}
		}
	}

	// Add context variables
	if ctx.ContextVars != nil {
		env["context"] = ctx.ContextVars
	}

	// Add input as both "item" and "input"
	if input != nil {
		env["item"] = input
		env["input"] = input
	}

	return env
}

// addCustomFunctions adds all custom functions to the environment
func (e *ExprEngine) addCustomFunctions(env map[string]interface{}) {
	// String functions
	env["contains"] = func(s, substr string) bool {
		return strings.Contains(s, substr)
	}
	env["startsWith"] = func(s, prefix string) bool {
		return strings.HasPrefix(s, prefix)
	}
	env["endsWith"] = func(s, suffix string) bool {
		return strings.HasSuffix(s, suffix)
	}
	env["upper"] = strings.ToUpper
	env["lower"] = strings.ToLower
	env["trim"] = strings.TrimSpace
	env["toUpperCase"] = strings.ToUpper
	env["toLowerCase"] = strings.ToLower
	env["split"] = strings.Split
	env["replace"] = strings.ReplaceAll
	env["join"] = func(arr []interface{}, sep string) string {
		strArr := make([]string, len(arr))
		for i, v := range arr {
			strArr[i] = fmt.Sprintf("%v", v)
		}
		return strings.Join(strArr, sep)
	}

	// Math functions
	env["pow"] = math.Pow
	env["sqrt"] = math.Sqrt
	// abs, floor, ceil, round are built-in in expr-lang

	// Array functions
	env["reverse"] = func(arr []interface{}) []interface{} {
		result := make([]interface{}, len(arr))
		for i, v := range arr {
			result[len(arr)-1-i] = v
		}
		return result
	}
	env["unique"] = func(arr []interface{}) []interface{} {
		seen := make(map[string]bool)
		result := make([]interface{}, 0)
		for _, item := range arr {
			key := fmt.Sprintf("%v", item)
			if !seen[key] {
				seen[key] = true
				result = append(result, item)
			}
		}
		return result
	}
	env["flatten"] = func(arr []interface{}) []interface{} {
		result := make([]interface{}, 0)
		var flattenRec func([]interface{})
		flattenRec = func(items []interface{}) {
			for _, item := range items {
				if subArr, ok := item.([]interface{}); ok {
					flattenRec(subArr)
				} else {
					result = append(result, item)
				}
			}
		}
		flattenRec(arr)
		return result
	}
	env["slice"] = func(arr []interface{}, start int, args ...int) []interface{} {
		end := len(arr)
		if len(args) > 0 {
			end = args[0]
		}
		if start < 0 {
			start = len(arr) + start
		}
		if end < 0 {
			end = len(arr) + end
		}
		if start < 0 {
			start = 0
		}
		if end > len(arr) {
			end = len(arr)
		}
		if start > end {
			return []interface{}{}
		}
		return arr[start:end]
	}
	env["first"] = func(arr []interface{}) interface{} {
		if len(arr) == 0 {
			return nil
		}
		return arr[0]
	}
	env["last"] = func(arr []interface{}) interface{} {
		if len(arr) == 0 {
			return nil
		}
		return arr[len(arr)-1]
	}

	// Aggregation functions - expr-lang has sum, min, max built-in
	// but we add avg for compatibility and make sum variadic
	env["avg"] = func(args ...interface{}) float64 {
		if len(args) == 0 {
			return 0
		}
		// Check if first arg is an array
		if arr, ok := args[0].([]interface{}); ok && len(args) == 1 {
			if len(arr) == 0 {
				return 0
			}
			sum := 0.0
			for _, v := range arr {
				if n, ok := toFloat64(v); ok {
					sum += n
				}
			}
			return sum / float64(len(arr))
		}
		// Multiple arguments
		sum := 0.0
		for _, v := range args {
			if n, ok := toFloat64(v); ok {
				sum += n
			}
		}
		return sum / float64(len(args))
	}
	
	// Override sum to support variadic args (expr-lang's sum only takes array)
	env["sum"] = func(args ...interface{}) float64 {
		if len(args) == 0 {
			return 0
		}
		// Check if first arg is an array
		if arr, ok := args[0].([]interface{}); ok && len(args) == 1 {
			sum := 0.0
			for _, v := range arr {
				if n, ok := toFloat64(v); ok {
					sum += n
				}
			}
			return sum
		}
		// Multiple arguments
		sum := 0.0
		for _, v := range args {
			if n, ok := toFloat64(v); ok {
				sum += n
			}
		}
		return sum
	}
	
	// Override min/max to support variadic args
	env["min"] = func(args ...interface{}) (float64, error) {
		if len(args) == 0 {
			return 0, fmt.Errorf("min() requires at least 1 argument")
		}
		// Check if first arg is an array
		if arr, ok := args[0].([]interface{}); ok && len(args) == 1 {
			if len(arr) == 0 {
				return 0, fmt.Errorf("min() on empty array")
			}
			minVal, ok := toFloat64(arr[0])
			if !ok {
				return 0, fmt.Errorf("min() requires numeric values")
			}
			for _, v := range arr[1:] {
				if n, ok := toFloat64(v); ok && n < minVal {
					minVal = n
				}
			}
			return minVal, nil
		}
		// Multiple arguments
		minVal, ok := toFloat64(args[0])
		if !ok {
			return 0, fmt.Errorf("min() requires numeric values")
		}
		for _, v := range args[1:] {
			if n, ok := toFloat64(v); ok && n < minVal {
				minVal = n
			}
		}
		return minVal, nil
	}
	
	env["max"] = func(args ...interface{}) (float64, error) {
		if len(args) == 0 {
			return 0, fmt.Errorf("max() requires at least 1 argument")
		}
		// Check if first arg is an array
		if arr, ok := args[0].([]interface{}); ok && len(args) == 1 {
			if len(arr) == 0 {
				return 0, fmt.Errorf("max() on empty array")
			}
			maxVal, ok := toFloat64(arr[0])
			if !ok {
				return 0, fmt.Errorf("max() requires numeric values")
			}
			for _, v := range arr[1:] {
				if n, ok := toFloat64(v); ok && n > maxVal {
					maxVal = n
				}
			}
			return maxVal, nil
		}
		// Multiple arguments
		maxVal, ok := toFloat64(args[0])
		if !ok {
			return 0, fmt.Errorf("max() requires numeric values")
		}
		for _, v := range args[1:] {
			if n, ok := toFloat64(v); ok && n > maxVal {
				maxVal = n
			}
		}
		return maxVal, nil
	}
	
	// Add zip function
	env["zip"] = func(args ...interface{}) []interface{} {
		if len(args) < 2 {
			return []interface{}{}
		}
		
		// Convert all args to arrays
		arrays := make([][]interface{}, 0, len(args))
		maxLen := 0
		for _, arg := range args {
			if arr, ok := arg.([]interface{}); ok {
				arrays = append(arrays, arr)
				if len(arr) > maxLen {
					maxLen = len(arr)
				}
			}
		}
		
		// Zip the arrays
		result := make([]interface{}, maxLen)
		for i := 0; i < maxLen; i++ {
			tuple := make([]interface{}, len(arrays))
			for j, arr := range arrays {
				if i < len(arr) {
					tuple[j] = arr[i]
				} else {
					tuple[j] = nil
				}
			}
			result[i] = tuple
		}
		return result
	}

	// Math functions that can work on arrays
	env["round"] = func(arg interface{}) interface{} {
		if arr, ok := arg.([]interface{}); ok {
			result := make([]interface{}, len(arr))
			for i, v := range arr {
				if n, ok := toFloat64(v); ok {
					result[i] = math.Round(n)
				}
			}
			return result
		}
		if n, ok := toFloat64(arg); ok {
			return math.Round(n)
		}
		return arg
	}
	
	env["floor"] = func(arg interface{}) interface{} {
		if arr, ok := arg.([]interface{}); ok {
			result := make([]interface{}, len(arr))
			for i, v := range arr {
				if n, ok := toFloat64(v); ok {
					result[i] = math.Floor(n)
				}
			}
			return result
		}
		if n, ok := toFloat64(arg); ok {
			return math.Floor(n)
		}
		return arg
	}
	
	env["ceil"] = func(arg interface{}) interface{} {
		if arr, ok := arg.([]interface{}); ok {
			result := make([]interface{}, len(arr))
			for i, v := range arr {
				if n, ok := toFloat64(v); ok {
					result[i] = math.Ceil(n)
				}
			}
			return result
		}
		if n, ok := toFloat64(arg); ok {
			return math.Ceil(n)
		}
		return arg
	}
	
	env["abs"] = func(arg interface{}) interface{} {
		if arr, ok := arg.([]interface{}); ok {
			result := make([]interface{}, len(arr))
			for i, v := range arr {
				if n, ok := toFloat64(v); ok {
					result[i] = math.Abs(n)
				}
			}
			return result
		}
		if n, ok := toFloat64(arg); ok {
			return math.Abs(n)
		}
		return arg
	}

	// Date/Time functions
	env["now"] = time.Now
	env["parseDate"] = parseDateTime
	env["toEpoch"] = func(val interface{}) (float64, error) {
		t, err := parseDateTime(val)
		if err != nil {
			return 0, err
		}
		return float64(t.Unix()), nil
	}
	env["toEpochMillis"] = func(val interface{}) (float64, error) {
		t, err := parseDateTime(val)
		if err != nil {
			return 0, err
		}
		return float64(t.UnixMilli()), nil
	}
	env["fromEpoch"] = func(seconds float64) time.Time {
		return time.Unix(int64(seconds), 0)
	}
	env["fromEpochMillis"] = func(millis float64) time.Time {
		return time.UnixMilli(int64(millis))
	}
	env["dateDiff"] = func(t1, t2 interface{}) (float64, error) {
		time1, err := parseDateTime(t1)
		if err != nil {
			return 0, err
		}
		time2, err := parseDateTime(t2)
		if err != nil {
			return 0, err
		}
		return time1.Sub(time2).Seconds(), nil
	}
	env["dateAdd"] = func(t interface{}, seconds float64) (time.Time, error) {
		time1, err := parseDateTime(t)
		if err != nil {
			return time.Time{}, err
		}
		return time1.Add(time.Duration(seconds) * time.Second), nil
	}
	env["year"] = func(t interface{}) (float64, error) {
		time1, err := parseDateTime(t)
		if err != nil {
			return 0, err
		}
		return float64(time1.Year()), nil
	}
	env["month"] = func(t interface{}) (float64, error) {
		time1, err := parseDateTime(t)
		if err != nil {
			return 0, err
		}
		return float64(time1.Month()), nil
	}
	env["day"] = func(t interface{}) (float64, error) {
		time1, err := parseDateTime(t)
		if err != nil {
			return 0, err
		}
		return float64(time1.Day()), nil
	}
	env["hour"] = func(t interface{}) (float64, error) {
		time1, err := parseDateTime(t)
		if err != nil {
			return 0, err
		}
		return float64(time1.Hour()), nil
	}
	env["minute"] = func(t interface{}) (float64, error) {
		time1, err := parseDateTime(t)
		if err != nil {
			return 0, err
		}
		return float64(time1.Minute()), nil
	}

	// Null handling
	env["isNull"] = func(v interface{}) bool {
		return v == nil
	}
	env["coalesce"] = func(args ...interface{}) interface{} {
		for _, arg := range args {
			if arg != nil {
				return arg
			}
		}
		return nil
	}
}

// Helper functions

// toFloat64 converts a value to float64
func toFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint64:
		return float64(v), true
	case uint32:
		return float64(v), true
	case string:
		// Try to parse as number
		var f float64
		if _, err := fmt.Sscanf(v, "%f", &f); err == nil {
			return f, true
		}
	}
	return 0, false
}

// parseDateTime parses various date/time formats into time.Time
func parseDateTime(value interface{}) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		// Try common formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			time.RFC822,
			time.RFC1123,
			"2006-01-02",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("unable to parse date/time: %s", v)
	case float64:
		// Assume Unix timestamp in seconds
		return time.Unix(int64(v), 0), nil
	case int64:
		// Unix timestamp in seconds
		return time.Unix(v, 0), nil
	case int:
		// Unix timestamp in seconds
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported date/time type: %T", value)
	}
}
