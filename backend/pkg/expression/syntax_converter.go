package expression

import (
	"regexp"
	"strings"
)

// convertSyntax converts Thaiyyal expression syntax to expr-lang syntax
// This provides compatibility with existing expressions while using expr-lang internally
func convertSyntax(expression string) string {
	// Convert .length property to len() function
	// Match: somevar.field.length or item.array.length
	lengthRe := regexp.MustCompile(`(\w+(?:\.\w+|\[\d+\])*?)\.length\b`)
	expression = lengthRe.ReplaceAllString(expression, "len($1)")

	// Convert map(array, expr) syntax to map(array, {closure})
	// This is more complex as we need to handle nested expressions
	expression = convertMapSyntax(expression)

	return expression
}

// convertMapSyntax converts map() function calls from Thaiyyal syntax to expr-lang syntax
// Thaiyyal: map(users, item.age * 2)
// expr-lang: map(users, {#.age * 2})
func convertMapSyntax(expression string) string {
	// Find map() calls
	mapRe := regexp.MustCompile(`map\s*\(\s*([^,]+),\s*(.+?)\s*\)`)
	
	// Process each map() call
	for {
		matches := mapRe.FindStringSubmatch(expression)
		if matches == nil {
			break
		}
		
		fullMatch := matches[0]
		arrayExpr := strings.TrimSpace(matches[1])
		itemExpr := strings.TrimSpace(matches[2])
		
		// Convert item references to # in the closure
		// Replace 'item.' with '#.' and standalone 'item' with '#'
		closureExpr := itemExpr
		closureExpr = regexp.MustCompile(`\bitem\.`).ReplaceAllString(closureExpr, "#.")
		closureExpr = regexp.MustCompile(`\bitem\b`).ReplaceAllString(closureExpr, "#")
		
		// Reconstruct the map call with closure syntax
		newMapCall := "map(" + arrayExpr + ", {" + closureExpr + "})"
		
		// Replace in expression
		expression = strings.Replace(expression, fullMatch, newMapCall, 1)
	}
	
	return expression
}
