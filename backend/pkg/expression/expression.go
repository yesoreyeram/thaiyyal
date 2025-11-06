// Package expression provides simple expression evaluation for workflow conditions.
// Supports node references, variables, context values, and boolean logic WITHOUT template delimiters.
// Now powered by expr-lang/expr for better performance and maintainability.
package expression

import (
	"regexp"
	"sync"
)

// Context provides access to workflow state during expression evaluation
type Context struct {
	NodeResults map[string]interface{} // Results from executed nodes
	Variables   map[string]interface{} // Workflow variables
	ContextVars map[string]interface{} // Context variables/constants
}

var (
	// Global engine instance for reuse and caching
	globalEngine *ExprEngine
	engineOnce   sync.Once
)

// getEngine returns the singleton expression engine
func getEngine() *ExprEngine {
	engineOnce.Do(func() {
		globalEngine = NewExprEngine()
	})
	return globalEngine
}

// Evaluate evaluates an expression and returns a boolean result
// Now powered by expr-lang/expr for better performance.
// Supports:
//   - Simple comparisons: ">100", "==5", "!=0", "value > 100"
//   - Node references: "node.id.output > 100"
//   - Variable references: "variables.count > 10"
//   - Context references: "context.maxValue < 50"
//   - Boolean operators: "&&", "||", "!"
//   - String operations: "contains(str, substr)", "startsWith()", etc.
func Evaluate(expression string, input interface{}, ctx *Context) (bool, error) {
	if ctx == nil {
		ctx = &Context{
			NodeResults: make(map[string]interface{}),
			Variables:   make(map[string]interface{}),
			ContextVars: make(map[string]interface{}),
		}
	}

	// If input is provided, ensure it's available as both 'item' and 'input'
	if input != nil {
		_, hasItem := ctx.Variables["item"]
		_, hasInput := ctx.Variables["input"]
		if !hasItem || !hasInput {
			// Create a copy of the context with item and input added
			newCtx := &Context{
				NodeResults: ctx.NodeResults,
				Variables:   make(map[string]interface{}),
				ContextVars: ctx.ContextVars,
			}
			// Copy existing variables
			for k, v := range ctx.Variables {
				newCtx.Variables[k] = v
			}
			// Add item and input
			if !hasItem {
				newCtx.Variables["item"] = input
			}
			if !hasInput {
				newCtx.Variables["input"] = input
			}
			ctx = newCtx
		}
	}

	// Use expr-lang/expr engine
	engine := getEngine()
	return engine.EvaluateBoolean(expression, input, ctx)
}

// EvaluateExpression evaluates an expression and returns its value (not just boolean)
// Now powered by expr-lang/expr for better performance.
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

	// If input is provided, ensure it's available as both 'item' and 'input'
	if input != nil {
		_, hasItem := ctx.Variables["item"]
		_, hasInput := ctx.Variables["input"]
		if !hasItem || !hasInput {
			// Create a shallow copy of the context and variables map
			newCtx := &Context{
				NodeResults: ctx.NodeResults,
				Variables:   make(map[string]interface{}),
				ContextVars: ctx.ContextVars,
			}
			for k, v := range ctx.Variables {
				newCtx.Variables[k] = v
			}
			if !hasItem {
				newCtx.Variables["item"] = input
			}
			if !hasInput {
				newCtx.Variables["input"] = input
			}
			ctx = newCtx
		}
	}

	// Use expr-lang/expr engine
	engine := getEngine()
	return engine.EvaluateValue(expression, input, ctx)
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

// Deprecated: EvaluateArithmetic is maintained for backward compatibility but now uses expr-lang internally
// Use EvaluateExpression instead.
func EvaluateArithmetic(expression string, ctx *Context) (float64, error) {
	engine := getEngine()
	result, err := engine.EvaluateValue(expression, nil, ctx)
	if err != nil {
		return 0, err
	}
	
	// Try to convert result to float64
	if num, ok := toFloat64(result); ok {
		return num, nil
	}
	
	return 0, nil
}
