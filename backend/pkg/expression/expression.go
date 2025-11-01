// Package expression provides simple expression evaluation for workflow conditions.
// Supports node references, variables, context values, and boolean logic WITHOUT template delimiters.
package expression

import (
"fmt"
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

// Check for NOT operator
if strings.HasPrefix(expression, "!") {
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
