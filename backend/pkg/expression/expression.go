// Package expression provides simple expression evaluation for workflow conditions.
// Supports node references, variables, context values, and boolean logic WITHOUT template delimiters.
package expression

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Context provides access to workflow state during expression evaluation
type Context struct {
	NodeResults map[string]interface{} // Results from executed nodes
	Variables   map[string]interface{} // Workflow variables
	ContextVars map[string]interface{} // Context variables/constants
}

// Evaluate evaluates an expression and returns a boolean result
// Supports:
//   - Simple comparisons: ">100", "==5", "!=0"
//   - Node references: "node.id.output > 100"
//   - Variable references: "variables.count > 10"
//   - Context references: "context.maxValue < 50"
//   - Boolean operators: "&&", "||", "!"
//   - String operations: "contains()", "=="
func Evaluate(expression string, input interface{}, ctx *Context) (bool, error) {
	if ctx == nil {
		ctx = &Context{
			NodeResults: make(map[string]interface{}),
			Variables:   make(map[string]interface{}),
			ContextVars: make(map[string]interface{}),
		}
	}

	// Trim whitespace
	expression = strings.TrimSpace(expression)

	// Handle boolean constants
	if expression == "true" {
		return true, nil
	}
	if expression == "false" {
		return false, nil
	}

	// Check for boolean operators (&&, ||)
	if result, ok := evaluateBooleanExpression(expression, input, ctx); ok {
		return result, nil
	}

	// Check for NOT operator (but not != which is a comparison operator)
	if strings.HasPrefix(expression, "!") && !strings.HasPrefix(expression, "!=") {
		result, err := Evaluate(strings.TrimSpace(expression[1:]), input, ctx)
		if err != nil {
			return false, err
		}
		return !result, nil
	}

	// Check for function calls (date/time, null handling, strings)
	if idx := strings.Index(expression, "("); idx > 0 && strings.HasSuffix(expression, ")") {
		funcName := strings.TrimSpace(expression[:idx])
		// Check if it's a known function
		if isFunctionCall(funcName) {
			return evaluateFunctionCall(expression, input, ctx)
		}
	}

	// Check for contains() function (backward compatibility)
	if strings.HasPrefix(expression, "contains(") && strings.HasSuffix(expression, ")") {
		return evaluateContains(expression, input, ctx)
	}

	// Check for comparison operators with references
	if result, ok := evaluateComparison(expression, input, ctx); ok {
		return result, nil
	}

	// Try to resolve as a direct boolean value reference
	if val, err := resolveValue(expression, input, ctx); err == nil {
		// If it's a boolean, return it directly
		if boolVal, ok := val.(bool); ok {
			return boolVal, nil
		}
	}

	// Fallback to simple numeric comparison (backward compatible)
	return evaluateSimpleCondition(expression, input), nil
}

// EvaluateExpression evaluates an expression and returns its value (not just boolean)
// This is used for transformations in Map and Reduce nodes.
// Supports:
//   - Arithmetic expressions: "item.age * 2", "accumulator + item.value"
//   - Ternary operator: "condition ? value1 : value2"
//   - String concatenation: "accumulator + item"
//   - Field access: "item.field", "item.nested.field"
//   - All value references (variables, node, context)
func EvaluateExpression(expression string, input interface{}, ctx *Context) (interface{}, error) {
	if ctx == nil {
		ctx = &Context{
			NodeResults: make(map[string]interface{}),
			Variables:   make(map[string]interface{}),
			ContextVars: make(map[string]interface{}),
		}
	}

	expression = strings.TrimSpace(expression)

	// Handle ternary operator: condition ? value1 : value2
	if idx := strings.Index(expression, "?"); idx > 0 {
		colonIdx := strings.Index(expression[idx:], ":")
		if colonIdx > 0 {
			colonIdx += idx
			condition := strings.TrimSpace(expression[:idx])
			trueVal := strings.TrimSpace(expression[idx+1 : colonIdx])
			falseVal := strings.TrimSpace(expression[colonIdx+1:])

			// Evaluate condition
			condResult, err := Evaluate(condition, input, ctx)
			if err != nil {
				return nil, fmt.Errorf("ternary condition evaluation failed: %w", err)
			}

			// Return appropriate value
			if condResult {
				return EvaluateExpression(trueVal, input, ctx)
			}
			return EvaluateExpression(falseVal, input, ctx)
		}
	}

	// Try arithmetic evaluation first (handles +, -, *, /, %, math functions)
	if containsArithmeticOp(expression) {
		result, err := EvaluateArithmetic(expression, ctx)
		if err == nil {
			return result, nil
		}
		// If arithmetic fails, continue to other evaluation methods
	}

	// Try to resolve as a value reference (variable, node, context, field access)
	if val, err := resolveValue(expression, input, ctx); err == nil {
		return val, nil
	}

	// Try as literal value
	if val, ok := parseLiteral(expression); ok {
		return val, nil
	}

	return nil, fmt.Errorf("could not evaluate expression: %s", expression)
}

// containsArithmeticOp checks if expression contains arithmetic operators
func containsArithmeticOp(expr string) bool {
	// Check for arithmetic operators
	arithmeticOps := []string{"+", "-", "*", "/", "%", "(", "pow", "sqrt", "abs", "floor", "ceil", "round", "min", "max"}
	for _, op := range arithmeticOps {
		if strings.Contains(expr, op) {
			return true
		}
	}
	return false
}

// parseLiteral attempts to parse a string as a literal value
func parseLiteral(s string) (interface{}, bool) {
	s = strings.TrimSpace(s)

	// String literal (quoted)
	if (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, `'`) && strings.HasSuffix(s, `'`)) {
		return s[1 : len(s)-1], true
	}

	// Boolean
	if s == "true" {
		return true, true
	}
	if s == "false" {
		return false, true
	}

	// Null
	if s == "null" {
		return nil, true
	}

	// Number
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return num, true
	}

	return nil, false
}

// ExtractDependencies extracts node IDs referenced in an expression
// This is used to build the dependency graph for topological sorting
func ExtractDependencies(expression string) []string {
	var dependencies []string
	seen := make(map[string]bool)

	// Find all node.id references using regex
	re := regexp.MustCompile(`node\.([a-zA-Z0-9_-]+)`)
	matches := re.FindAllStringSubmatch(expression, -1)

	for _, match := range matches {
		if len(match) > 1 {
			nodeID := match[1]
			if !seen[nodeID] {
				dependencies = append(dependencies, nodeID)
				seen[nodeID] = true
			}
		}
	}

	return dependencies
}

// evaluateBooleanExpression handles && and || operators
func evaluateBooleanExpression(expr string, input interface{}, ctx *Context) (bool, bool) {
	// Check for || (OR) - lower precedence
	if idx := findOperator(expr, "||"); idx != -1 {
		left := strings.TrimSpace(expr[:idx])
		right := strings.TrimSpace(expr[idx+2:])

		leftResult, err := Evaluate(left, input, ctx)
		if err != nil {
			return false, false
		}

		rightResult, err := Evaluate(right, input, ctx)
		if err != nil {
			return false, false
		}

		return leftResult || rightResult, true
	}

	// Check for && (AND) - higher precedence
	if idx := findOperator(expr, "&&"); idx != -1 {
		left := strings.TrimSpace(expr[:idx])
		right := strings.TrimSpace(expr[idx+2:])

		leftResult, err := Evaluate(left, input, ctx)
		if err != nil {
			return false, false
		}

		rightResult, err := Evaluate(right, input, ctx)
		if err != nil {
			return false, false
		}

		return leftResult && rightResult, true
	}

	return false, false
}

// findOperator finds the position of an operator, respecting parentheses
func findOperator(expr string, op string) int {
	depth := 0
	for i := 0; i <= len(expr)-len(op); i++ {
		if expr[i] == '(' {
			depth++
		} else if expr[i] == ')' {
			depth--
		} else if depth == 0 && strings.HasPrefix(expr[i:], op) {
			return i
		}
	}
	return -1
}

// evaluateComparison handles comparison operations
func evaluateComparison(expr string, input interface{}, ctx *Context) (bool, bool) {
	operators := []string{"==", "!=", "<=", ">=", "<", ">"}

	for _, op := range operators {
		idx := findOperator(expr, op)
		if idx == -1 {
			continue
		}

		left := strings.TrimSpace(expr[:idx])
		right := strings.TrimSpace(expr[idx+len(op):])

		// Resolve left operand
		leftVal, err := resolveValue(left, input, ctx)
		if err != nil {
			return false, false
		}

		// Resolve right operand
		rightVal, err := resolveValue(right, input, ctx)
		if err != nil {
			return false, false
		}

		// Perform comparison
		result := compareValues(leftVal, rightVal, op)
		return result, true
	}

	return false, false
}

// resolveValue resolves a value from a reference or literal
func resolveValue(ref string, input interface{}, ctx *Context) (interface{}, error) {
	ref = strings.TrimSpace(ref)

	// Check for string literals
	if (strings.HasPrefix(ref, "\"") && strings.HasSuffix(ref, "\"")) ||
		(strings.HasPrefix(ref, "'") && strings.HasSuffix(ref, "'")) {
		return ref[1 : len(ref)-1], nil
	}

	// Check for numeric literals
	if num, err := strconv.ParseFloat(ref, 64); err == nil {
		return num, nil
	}

	// Check for boolean literals
	if ref == "true" {
		return true, nil
	}
	if ref == "false" {
		return false, nil
	}

	// Check if it contains arithmetic operators or function calls - try arithmetic evaluation
	if containsArithmetic(ref) {
		if val, err := EvaluateArithmetic(ref, ctx); err == nil {
			return val, nil
		}
		// If arithmetic evaluation fails, continue with reference resolution
	}

	// Check for node reference: node.id.output or node.id.value
	if strings.HasPrefix(ref, "node.") {
		return resolveNodeReference(ref, ctx)
	}

	// Check for variable reference: variables.name or variables.name.field
	if strings.HasPrefix(ref, "variables.") {
		// Parse: variables.name.field or just variables.name
		parts := strings.Split(ref, ".")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid variable reference: %s", ref)
		}

		varName := parts[1]
		val, ok := ctx.Variables[varName]
		if !ok {
			return nil, fmt.Errorf("variable not found: %s", varName)
		}

		// If just variables.name, return the whole value
		if len(parts) == 2 {
			return val, nil
		}

		// Navigate to nested field (like node references)
		current := val
		for i := 2; i < len(parts); i++ {
			field := parts[i]
			if m, ok := current.(map[string]interface{}); ok {
				if fieldVal, exists := m[field]; exists {
					current = fieldVal
				} else {
					return nil, fmt.Errorf("field not found: %s in variables.%s", field, varName)
				}
			} else {
				return nil, fmt.Errorf("cannot access field %s on non-object", field)
			}
		}

		return current, nil
	}

	// Check for context reference: context.name
	if strings.HasPrefix(ref, "context.") {
		ctxName := ref[8:] // Remove "context." prefix
		if val, ok := ctx.ContextVars[ctxName]; ok {
			return val, nil
		}
		return nil, fmt.Errorf("context variable not found: %s", ctxName)
	}

	// Check for input reference
	if ref == "input" {
		return input, nil
	}

	// Check for item reference: item.field or just item
	// This is the preferred syntax for filter expressions (e.g., "item.age >= 18")
	if strings.HasPrefix(ref, "item.") || ref == "item" {
		if ref == "item" {
			return input, nil
		}
		// Navigate to nested field starting from input
		fieldPath := ref[5:] // Remove "item." prefix
		return resolveFieldPath(fieldPath, input)
	}

	// Check for input reference: input.field or just input
	// Many condition expressions use 'input' as the value placeholder (e.g., "input > 10" or "input.age > 18")
	if strings.HasPrefix(ref, "input.") || ref == "input" {
		if ref == "input" {
			return input, nil
		}
		// Navigate to nested field starting from input
		fieldPath := ref[6:] // Remove "input." prefix
		return resolveFieldPath(fieldPath, input)
	}

	// Check for direct field access on input object (e.g., "age", "name", "profile.verified")
	// This also works but "item.age" is more explicit and recommended
	if input != nil {
		// Try to resolve as a field path on the input object
		if val, err := resolveFieldPath(ref, input); err == nil {
			return val, nil
		}
	}

	return nil, fmt.Errorf("unknown reference: %s", ref)
}

// resolveFieldPath resolves a field path (e.g., "age" or "profile.verified") from an object
func resolveFieldPath(path string, obj interface{}) (interface{}, error) {
	parts := strings.Split(path, ".")
	current := obj

	for _, field := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[field]; exists {
				current = val
			} else {
				return nil, fmt.Errorf("field not found: %s", field)
			}
		} else {
			return nil, fmt.Errorf("cannot access field %s on non-object", field)
		}
	}

	return current, nil
}

// containsArithmetic checks if an expression contains arithmetic operators or functions
// containsArithmetic checks if an expression contains arithmetic operators or functions
func containsArithmetic(expr string) bool {
	// Check for math functions first
	mathFuncs := []string{"pow(", "sqrt(", "abs(", "floor(", "ceil(", "round(", "min(", "max("}
	for _, fn := range mathFuncs {
		if strings.Contains(expr, fn) {
			return true
		}
	}

	// Don't treat simple variable/node references as arithmetic
	// Only return true if we find actual arithmetic operators used as operators (not comparison)
	hasArithOp := false

	// Look for +, *, /, % which are always arithmetic
	if strings.ContainsAny(expr, "*/%+") {
		hasArithOp = true
	}

	return hasArithOp
}

// resolveNodeReference resolves node.id.field references
func resolveNodeReference(ref string, ctx *Context) (interface{}, error) {
	// Parse: node.id.field or just node.id
	parts := strings.Split(ref, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid node reference: %s", ref)
	}

	nodeID := parts[1]
	result, ok := ctx.NodeResults[nodeID]
	if !ok {
		return nil, fmt.Errorf("node result not found: %s", nodeID)
	}

	// If just node.id, return the whole result
	if len(parts) == 2 {
		return result, nil
	}

	// Navigate to nested field
	current := result
	for i := 2; i < len(parts); i++ {
		field := parts[i]
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[field]; exists {
				current = val
			} else {
				return nil, fmt.Errorf("field not found: %s in node.%s", field, nodeID)
			}
		} else {
			return nil, fmt.Errorf("cannot access field %s on non-object", field)
		}
	}

	return current, nil
}

// compareValues compares two values using the specified operator
func compareValues(left, right interface{}, op string) bool {
	switch op {
	case "==":
		return compareEquality(left, right)
	case "!=":
		return !compareEquality(left, right)
	case "<", "<=", ">", ">=":
		return compareNumeric(left, right, op)
	}
	return false
}

// compareEquality compares two values for equality
func compareEquality(left, right interface{}) bool {
	// Handle nil
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	// Try time.Time comparison
	leftTime, leftIsTime := left.(time.Time)
	rightTime, rightIsTime := right.(time.Time)
	if leftIsTime && rightIsTime {
		return leftTime.Equal(rightTime)
	}

	// Try numeric comparison
	leftNum, leftIsNum := toFloat64(left)
	rightNum, rightIsNum := toFloat64(right)
	if leftIsNum && rightIsNum {
		return leftNum == rightNum
	}

	// Try string comparison
	leftStr, leftIsStr := left.(string)
	rightStr, rightIsStr := right.(string)
	if leftIsStr && rightIsStr {
		return leftStr == rightStr
	}

	// Try boolean comparison
	leftBool, leftIsBool := left.(bool)
	rightBool, rightIsBool := right.(bool)
	if leftIsBool && rightIsBool {
		return leftBool == rightBool
	}

	return false
}

// compareNumeric compares two values numerically
func compareNumeric(left, right interface{}, op string) bool {
	// Handle time.Time comparisons
	leftTime, leftIsTime := left.(time.Time)
	rightTime, rightIsTime := right.(time.Time)
	if leftIsTime && rightIsTime {
		switch op {
		case "<":
			return leftTime.Before(rightTime)
		case "<=":
			return leftTime.Before(rightTime) || leftTime.Equal(rightTime)
		case ">":
			return leftTime.After(rightTime)
		case ">=":
			return leftTime.After(rightTime) || leftTime.Equal(rightTime)
		}
		return false
	}

	leftNum, leftOk := toFloat64(left)
	rightNum, rightOk := toFloat64(right)

	if !leftOk || !rightOk {
		return false
	}

	switch op {
	case "<":
		return leftNum < rightNum
	case "<=":
		return leftNum <= rightNum
	case ">":
		return leftNum > rightNum
	case ">=":
		return leftNum >= rightNum
	}

	return false
}

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
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

// evaluateContains evaluates contains(haystack, needle) function
func evaluateContains(expr string, input interface{}, ctx *Context) (bool, error) {
	// Extract arguments from contains(arg1, arg2)
	inner := expr[9 : len(expr)-1] // Remove "contains(" and ")"
	args := splitArguments(inner)

	if len(args) != 2 {
		return false, fmt.Errorf("contains() requires exactly 2 arguments")
	}

	haystack, err := resolveValue(args[0], input, ctx)
	if err != nil {
		return false, err
	}

	needle, err := resolveValue(args[1], input, ctx)
	if err != nil {
		return false, err
	}

	haystackStr := fmt.Sprintf("%v", haystack)
	needleStr := fmt.Sprintf("%v", needle)

	return strings.Contains(haystackStr, needleStr), nil
}

// splitArguments splits function arguments by comma (respecting quotes)
func splitArguments(s string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if ch == '"' || ch == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuotes = false
			}
			current.WriteByte(ch)
		} else if ch == ',' && !inQuotes {
			args = append(args, strings.TrimSpace(current.String()))
			current.Reset()
		} else {
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		args = append(args, strings.TrimSpace(current.String()))
	}

	return args
}

// isFunctionCall checks if a name is a known function
func isFunctionCall(name string) bool {
	dateFuncs := []string{"now", "parseDate", "toEpoch", "toEpochMillis", "fromEpoch", "fromEpochMillis",
		"dateDiff", "dateAdd", "year", "month", "day", "hour", "minute", "isNull", "coalesce"}
	for _, fn := range dateFuncs {
		if name == fn {
			return true
		}
	}
	return false
}

// evaluateFunctionCall evaluates a function call and returns a boolean result
func evaluateFunctionCall(expr string, input interface{}, ctx *Context) (bool, error) {
	// Extract function name and arguments
	idx := strings.Index(expr, "(")
	if idx == -1 {
		return false, fmt.Errorf("invalid function call: %s", expr)
	}

	funcName := strings.TrimSpace(expr[:idx])
	argsStr := expr[idx+1 : len(expr)-1] // Remove "funcName(" and ")"

	// Parse arguments
	argStrs := splitArguments(argsStr)
	var args []interface{}
	for _, argStr := range argStrs {
		val, err := resolveValue(argStr, input, ctx)
		if err != nil {
			return false, fmt.Errorf("error resolving argument '%s': %w", argStr, err)
		}
		args = append(args, val)
	}

	// Call the function
	result, err := callDateTimeFunction(funcName, args, ctx)
	if err != nil {
		return false, err
	}

	// Convert result to boolean
	if boolResult, ok := result.(bool); ok {
		return boolResult, nil
	}

	// Non-boolean results from functions like coalesce might need comparison
	return false, fmt.Errorf("function %s() did not return a boolean value", funcName)
}

// evaluateSimpleCondition evaluates simple numeric conditions (backward compatible)
func evaluateSimpleCondition(condition string, value interface{}) bool {
	// Handle boolean constants
	if condition == "true" {
		return true
	}
	if condition == "false" {
		return false
	}

	// Extract numeric value from input
	numVal, ok := value.(float64)
	if !ok {
		// Try to extract value from map (common in node results)
		if m, isMap := value.(map[string]interface{}); isMap {
			if v, exists := m["value"]; exists {
				numVal, ok = v.(float64)
			}
		}
		if !ok {
			return false
		}
	}

	// Parse condition using a simple state machine
	var threshold float64
	var operator string

	if len(condition) >= 2 {
		// Check two-character operators first
		twoChar := condition[0:2]
		switch twoChar {
		case ">=":
			operator = ">="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		case "<=":
			operator = "<="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		case "==":
			operator = "=="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		case "!=":
			operator = "!="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		default:
			// Single-character operators
			switch condition[0] {
			case '>':
				operator = ">"
				fmt.Sscanf(condition[1:], "%f", &threshold)
			case '<':
				operator = "<"
				fmt.Sscanf(condition[1:], "%f", &threshold)
			}
		}
	}

	// Evaluate comparison
	switch operator {
	case ">":
		return numVal > threshold
	case "<":
		return numVal < threshold
	case ">=":
		return numVal >= threshold
	case "<=":
		return numVal <= threshold
	case "==":
		return numVal == threshold
	case "!=":
		return numVal != threshold
	default:
		return false
	}
}

// ============================================================================
// Arithmetic Expression Evaluation
// ============================================================================

// EvaluateArithmetic evaluates an arithmetic expression and returns a numeric result
// Supports:
//   - Basic operations: +, -, *, /, %
//   - Parentheses for grouping: (a + b) * c
//   - Math functions: pow, sqrt, abs, floor, ceil, round, min, max
//   - Variable references: variables.name
//   - Node references: node.id.value
func EvaluateArithmetic(expression string, ctx *Context) (float64, error) {
	if ctx == nil {
		ctx = &Context{
			NodeResults: make(map[string]interface{}),
			Variables:   make(map[string]interface{}),
			ContextVars: make(map[string]interface{}),
		}
	}

	expression = strings.TrimSpace(expression)
	if expression == "" {
		return 0, fmt.Errorf("empty expression")
	}

	// Parse and evaluate the expression
	parser := &arithmeticParser{
		expression: expression,
		pos:        0,
		ctx:        ctx,
	}

	result, err := parser.parseExpression()
	if err != nil {
		return 0, err
	}

	// Make sure we consumed the entire expression
	parser.skipWhitespace()
	if parser.pos < len(parser.expression) {
		return 0, fmt.Errorf("unexpected characters at position %d: %s", parser.pos, parser.expression[parser.pos:])
	}

	return result, nil
}

// arithmeticParser is a recursive descent parser for arithmetic expressions
type arithmeticParser struct {
	expression string
	pos        int
	ctx        *Context
}

// parseExpression parses addition and subtraction (lowest precedence)
func (p *arithmeticParser) parseExpression() (float64, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expression) {
			break
		}

		op := p.peek()
		if op != '+' && op != '-' {
			break
		}

		p.pos++
		right, err := p.parseTerm()
		if err != nil {
			return 0, err
		}

		if op == '+' {
			left = left + right
		} else {
			left = left - right
		}
	}

	return left, nil
}

// parseTerm parses multiplication, division, and modulo (higher precedence)
func (p *arithmeticParser) parseTerm() (float64, error) {
	left, err := p.parseFactor()
	if err != nil {
		return 0, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.expression) {
			break
		}

		op := p.peek()
		if op != '*' && op != '/' && op != '%' {
			break
		}

		p.pos++
		right, err := p.parseFactor()
		if err != nil {
			return 0, err
		}

		switch op {
		case '*':
			left = left * right
		case '/':
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left = left / right
		case '%':
			if right == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			left = float64(int(left) % int(right))
		}
	}

	return left, nil
}

// parseFactor parses unary operators, numbers, variables, function calls, and parentheses
func (p *arithmeticParser) parseFactor() (float64, error) {
	p.skipWhitespace()

	if p.pos >= len(p.expression) {
		return 0, fmt.Errorf("unexpected end of expression")
	}

	// Handle unary operators
	if p.peek() == '+' {
		p.pos++
		return p.parseFactor()
	}
	if p.peek() == '-' {
		p.pos++
		val, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		return -val, nil
	}

	// Handle parentheses
	if p.peek() == '(' {
		p.pos++
		val, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		p.skipWhitespace()
		if p.pos >= len(p.expression) || p.peek() != ')' {
			return 0, fmt.Errorf("unmatched parentheses at position %d", p.pos)
		}
		p.pos++
		return val, nil
	}

	// Handle numbers
	if p.isDigit(p.peek()) {
		return p.parseNumber()
	}

	// Handle identifiers (variables, node references, function calls)
	if p.isLetter(p.peek()) {
		return p.parseIdentifier()
	}

	return 0, fmt.Errorf("unexpected character '%c' at position %d", p.peek(), p.pos)
}

// parseNumber parses a numeric literal
func (p *arithmeticParser) parseNumber() (float64, error) {
	start := p.pos
	hasDecimal := false

	for p.pos < len(p.expression) {
		ch := p.expression[p.pos]
		if ch == '.' {
			if hasDecimal {
				break
			}
			hasDecimal = true
			p.pos++
		} else if p.isDigit(ch) {
			p.pos++
		} else {
			break
		}
	}

	numStr := p.expression[start:p.pos]
	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number '%s' at position %d", numStr, start)
	}

	return val, nil
}

// parseIdentifier parses an identifier (variable, node reference, or function call)
func (p *arithmeticParser) parseIdentifier() (float64, error) {
	start := p.pos

	// Read identifier
	for p.pos < len(p.expression) && (p.isLetter(p.expression[p.pos]) || p.isDigit(p.expression[p.pos]) || p.expression[p.pos] == '_') {
		p.pos++
	}

	ident := p.expression[start:p.pos]

	// Check if it's a function call
	p.skipWhitespace()
	if p.pos < len(p.expression) && p.peek() == '(' {
		return p.parseFunction(ident)
	}

	// Check if it's a dotted path (variables.x or node.id.field)
	if p.pos < len(p.expression) && p.peek() == '.' {
		// Read the full path
		path := ident
		for p.pos < len(p.expression) && p.peek() == '.' {
			p.pos++ // skip '.'
			start := p.pos
			for p.pos < len(p.expression) && (p.isLetter(p.expression[p.pos]) || p.isDigit(p.expression[p.pos]) || p.expression[p.pos] == '_' || p.expression[p.pos] == '-') {
				p.pos++
			}
			if p.pos == start {
				return 0, fmt.Errorf("expected identifier after '.' at position %d", p.pos)
			}
			path += "." + p.expression[start:p.pos]
		}

		// Resolve the path
		val, err := resolveValue(path, nil, p.ctx)
		if err != nil {
			return 0, err
		}

		// Convert to float64
		num, ok := toFloat64(val)
		if !ok {
			return 0, fmt.Errorf("value '%v' at '%s' cannot be converted to number", val, path)
		}

		return num, nil
	}

	return 0, fmt.Errorf("unknown identifier '%s' at position %d", ident, start)
}

// parseFunction parses a function call
func (p *arithmeticParser) parseFunction(name string) (float64, error) {
	p.pos++ // skip '('

	var args []float64

	p.skipWhitespace()
	if p.peek() == ')' {
		p.pos++
		return p.callFunction(name, args)
	}

	for {
		arg, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		args = append(args, arg)

		p.skipWhitespace()
		if p.pos >= len(p.expression) {
			return 0, fmt.Errorf("unmatched parentheses in function call")
		}

		if p.peek() == ')' {
			p.pos++
			break
		}

		if p.peek() == ',' {
			p.pos++
			p.skipWhitespace()
			continue
		}

		return 0, fmt.Errorf("expected ',' or ')' at position %d", p.pos)
	}

	return p.callFunction(name, args)
}

// callFunction executes a math function
func (p *arithmeticParser) callFunction(name string, args []float64) (float64, error) {

	switch name {
	case "pow":
		if len(args) != 2 {
			return 0, fmt.Errorf("pow() requires exactly 2 arguments, got %d", len(args))
		}
		return math.Pow(args[0], args[1]), nil

	case "sqrt":
		if len(args) != 1 {
			return 0, fmt.Errorf("sqrt() requires exactly 1 argument, got %d", len(args))
		}
		if args[0] < 0 {
			return 0, fmt.Errorf("sqrt() of negative number")
		}
		return math.Sqrt(args[0]), nil

	case "abs":
		if len(args) != 1 {
			return 0, fmt.Errorf("abs() requires exactly 1 argument, got %d", len(args))
		}
		return math.Abs(args[0]), nil

	case "floor":
		if len(args) != 1 {
			return 0, fmt.Errorf("floor() requires exactly 1 argument, got %d", len(args))
		}
		return math.Floor(args[0]), nil

	case "ceil":
		if len(args) != 1 {
			return 0, fmt.Errorf("ceil() requires exactly 1 argument, got %d", len(args))
		}
		return math.Ceil(args[0]), nil

	case "round":
		if len(args) != 1 {
			return 0, fmt.Errorf("round() requires exactly 1 argument, got %d", len(args))
		}
		return math.Round(args[0]), nil

	case "min":
		if len(args) < 2 {
			return 0, fmt.Errorf("min() requires at least 2 arguments, got %d", len(args))
		}
		min := args[0]
		for _, arg := range args[1:] {
			if arg < min {
				min = arg
			}
		}
		return min, nil

	case "max":
		if len(args) < 2 {
			return 0, fmt.Errorf("max() requires at least 2 arguments, got %d", len(args))
		}
		max := args[0]
		for _, arg := range args[1:] {
			if arg > max {
				max = arg
			}
		}
		return max, nil

	default:
		return 0, fmt.Errorf("unknown function '%s'", name)
	}
}

// skipWhitespace skips whitespace characters
func (p *arithmeticParser) skipWhitespace() {
	for p.pos < len(p.expression) && (p.expression[p.pos] == ' ' || p.expression[p.pos] == '\t' || p.expression[p.pos] == '\n' || p.expression[p.pos] == '\r') {
		p.pos++
	}
}

// peek returns the current character without advancing
func (p *arithmeticParser) peek() byte {
	if p.pos >= len(p.expression) {
		return 0
	}
	return p.expression[p.pos]
}

// isDigit checks if a character is a digit
func (p *arithmeticParser) isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// isLetter checks if a character is a letter
func (p *arithmeticParser) isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// ============================================================================
// Date/Time Functions
// ============================================================================

// Date/time helper functions for the expression evaluator

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

// isNull checks if a value is null/nil
func isNull(value interface{}) bool {
	return value == nil
}

// callDateTimeFunction handles date/time specific functions
func callDateTimeFunction(name string, args []interface{}, _ *Context) (interface{}, error) {
	switch name {
	case "now":
		// Current timestamp
		if len(args) != 0 {
			return nil, fmt.Errorf("now() takes no arguments, got %d", len(args))
		}
		return time.Now(), nil

	case "parseDate":
		// Parse date string
		if len(args) != 1 {
			return nil, fmt.Errorf("parseDate() requires exactly 1 argument, got %d", len(args))
		}
		return parseDateTime(args[0])

	case "toEpoch":
		// Convert to Unix timestamp (seconds)
		if len(args) != 1 {
			return nil, fmt.Errorf("toEpoch() requires exactly 1 argument, got %d", len(args))
		}
		t, err := parseDateTime(args[0])
		if err != nil {
			return nil, err
		}
		return float64(t.Unix()), nil

	case "toEpochMillis":
		// Convert to Unix timestamp (milliseconds)
		if len(args) != 1 {
			return nil, fmt.Errorf("toEpochMillis() requires exactly 1 argument, got %d", len(args))
		}
		t, err := parseDateTime(args[0])
		if err != nil {
			return nil, err
		}
		return float64(t.UnixMilli()), nil

	case "fromEpoch":
		// Create time from Unix timestamp (seconds)
		if len(args) != 1 {
			return nil, fmt.Errorf("fromEpoch() requires exactly 1 argument, got %d", len(args))
		}
		val, ok := toFloat64(args[0])
		if !ok {
			return nil, fmt.Errorf("fromEpoch() requires numeric argument")
		}
		return time.Unix(int64(val), 0), nil

	case "fromEpochMillis":
		// Create time from Unix timestamp (milliseconds)
		if len(args) != 1 {
			return nil, fmt.Errorf("fromEpochMillis() requires exactly 1 argument, got %d", len(args))
		}
		val, ok := toFloat64(args[0])
		if !ok {
			return nil, fmt.Errorf("fromEpochMillis() requires numeric argument")
		}
		return time.UnixMilli(int64(val)), nil

	case "dateDiff":
		// Difference between two dates in seconds
		if len(args) != 2 {
			return nil, fmt.Errorf("dateDiff() requires exactly 2 arguments, got %d", len(args))
		}
		t1, err := parseDateTime(args[0])
		if err != nil {
			return nil, fmt.Errorf("dateDiff() first argument: %w", err)
		}
		t2, err := parseDateTime(args[1])
		if err != nil {
			return nil, fmt.Errorf("dateDiff() second argument: %w", err)
		}
		return float64(t1.Sub(t2).Seconds()), nil

	case "dateAdd":
		// Add seconds to a date
		if len(args) != 2 {
			return nil, fmt.Errorf("dateAdd() requires exactly 2 arguments, got %d", len(args))
		}
		t, err := parseDateTime(args[0])
		if err != nil {
			return nil, fmt.Errorf("dateAdd() first argument: %w", err)
		}
		seconds, ok := toFloat64(args[1])
		if !ok {
			return nil, fmt.Errorf("dateAdd() second argument must be numeric")
		}
		return t.Add(time.Duration(seconds) * time.Second), nil

	case "year":
		// Get year from date
		if len(args) != 1 {
			return nil, fmt.Errorf("year() requires exactly 1 argument, got %d", len(args))
		}
		t, err := parseDateTime(args[0])
		if err != nil {
			return nil, err
		}
		return float64(t.Year()), nil

	case "month":
		// Get month from date (1-12)
		if len(args) != 1 {
			return nil, fmt.Errorf("month() requires exactly 1 argument, got %d", len(args))
		}
		t, err := parseDateTime(args[0])
		if err != nil {
			return nil, err
		}
		return float64(t.Month()), nil

	case "day":
		// Get day from date
		if len(args) != 1 {
			return nil, fmt.Errorf("day() requires exactly 1 argument, got %d", len(args))
		}
		t, err := parseDateTime(args[0])
		if err != nil {
			return nil, err
		}
		return float64(t.Day()), nil

	case "hour":
		// Get hour from date
		if len(args) != 1 {
			return nil, fmt.Errorf("hour() requires exactly 1 argument, got %d", len(args))
		}
		t, err := parseDateTime(args[0])
		if err != nil {
			return nil, err
		}
		return float64(t.Hour()), nil

	case "minute":
		// Get minute from date
		if len(args) != 1 {
			return nil, fmt.Errorf("minute() requires exactly 1 argument, got %d", len(args))
		}
		t, err := parseDateTime(args[0])
		if err != nil {
			return nil, err
		}
		return float64(t.Minute()), nil

	case "isNull":
		// Check if value is null
		if len(args) != 1 {
			return nil, fmt.Errorf("isNull() requires exactly 1 argument, got %d", len(args))
		}
		return isNull(args[0]), nil

	case "coalesce":
		// Return first non-null value
		if len(args) < 1 {
			return nil, fmt.Errorf("coalesce() requires at least 1 argument")
		}
		for _, arg := range args {
			if !isNull(arg) {
				return arg, nil
			}
		}
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown date/time function: %s", name)
	}
}
