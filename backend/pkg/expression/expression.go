// Package expression provides simple expression evaluation for workflow conditions.
// Supports node references, variables, context values, and boolean logic WITHOUT template delimiters.
package expression

import (
"fmt"
"math"
"regexp"
"strconv"
"strings"
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

// Check for contains() function
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

// Check for variable reference: variables.name
if strings.HasPrefix(ref, "variables.") {
varName := ref[10:] // Remove "variables." prefix
if val, ok := ctx.Variables[varName]; ok {
return val, nil
}
return nil, fmt.Errorf("variable not found: %s", varName)
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

return nil, fmt.Errorf("unknown reference: %s", ref)
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
