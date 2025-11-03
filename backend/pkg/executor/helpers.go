package executor

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ============================================================================
// Type Conversion Helpers
// ============================================================================

// convertTypedValue converts a value to the specified type.
// Supported types:
//   - "string": Convert to string
//   - "number": Convert to float64
//   - "boolean": Convert to bool
//   - "time_string": Parse as RFC3339 time string, return as string
//   - "epoch_second": Parse as Unix epoch seconds, return as time.Time
//   - "epoch_ms": Parse as Unix epoch milliseconds, return as time.Time
//   - "null": Return nil
func convertTypedValue(value interface{}, valueType string) (interface{}, error) {
	switch valueType {
	case "string":
		return fmt.Sprintf("%v", value), nil

	case "number":
		// Try to convert to float64
		switch v := value.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %q to number: %w", v, err)
			}
			return f, nil
		default:
			return nil, fmt.Errorf("cannot convert type %T to number", value)
		}

	case "boolean":
		// Try to convert to bool
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %q to boolean: %w", v, err)
			}
			return b, nil
		case float64:
			return v != 0, nil
		default:
			return nil, fmt.Errorf("cannot convert type %T to boolean", value)
		}

	case "time_string":
		// Parse as RFC3339 time string
		var timeStr string
		switch v := value.(type) {
		case string:
			timeStr = v
		default:
			timeStr = fmt.Sprintf("%v", value)
		}

		// Validate it's a valid time string
		_, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid time string %q: %w", timeStr, err)
		}
		return timeStr, nil

	case "epoch_second":
		// Convert to Unix epoch seconds, return as time.Time
		var seconds int64
		switch v := value.(type) {
		case float64:
			seconds = int64(v)
		case int:
			seconds = int64(v)
		case int64:
			seconds = v
		case string:
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %q to epoch seconds: %w", v, err)
			}
			seconds = i
		default:
			return nil, fmt.Errorf("cannot convert type %T to epoch seconds", value)
		}
		return time.Unix(seconds, 0), nil

	case "epoch_ms":
		// Convert to Unix epoch milliseconds, return as time.Time
		var ms int64
		switch v := value.(type) {
		case float64:
			ms = int64(v)
		case int:
			ms = int64(v)
		case int64:
			ms = v
		case string:
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %q to epoch milliseconds: %w", v, err)
			}
			ms = i
		default:
			return nil, fmt.Errorf("cannot convert type %T to epoch milliseconds", value)
		}
		return time.Unix(ms/1000, (ms%1000)*1000000), nil

	case "null":
		return nil, nil

	default:
		return nil, fmt.Errorf("unsupported type %q", valueType)
	}
}

// ============================================================================
// Condition Evaluation Helpers
// ============================================================================

// evaluateCondition evaluates a condition string against an input value.
//
// Supported condition formats:
//   - "true" - Always true
//   - "false" - Always false
//   - ">N" - Greater than N
//   - "<N" - Less than N
//   - ">=N" - Greater than or equal to N
//   - "<=N" - Less than or equal to N
//   - "==N" - Equal to N
//   - "!=N" - Not equal to N
//
// The value can be a direct number or a map containing a "value" field.
//
// Returns:
//   - bool: true if condition is met, false otherwise
func evaluateCondition(condition string, value interface{}) bool {
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

	// Evaluate comparison using strategy pattern
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
// Text Transformation Helpers
// ============================================================================

// toTitleCase converts text to Title Case (first letter of each word capitalized).
func toTitleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

// toCamelCase converts text to camelCase.
// Example: "hello world" → "helloWorld"
func toCamelCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		word := words[i]
		if len(word) > 0 {
			// Capitalize first letter, lowercase rest
			result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return result
}

// toInverseCase inverts the case of each character.
// Example: "Hello" → "hELLO"
func toInverseCase(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			runes[i] = unicode.ToLower(r)
		} else if unicode.IsLower(r) {
			runes[i] = unicode.ToUpper(r)
		}
	}
	return string(runes)
}

// ============================================================================
// Duration Parsing Helpers
// ============================================================================

// parseDuration parses duration strings with support for ms, s, m, h
func parseDuration(durationStr string) (time.Duration, error) {
	// Support formats like "5s", "10m", "1h", "100ms"
	if duration, err := time.ParseDuration(durationStr); err == nil {
		return duration, nil
	}

	// Also support integer milliseconds
	if ms, err := strconv.Atoi(durationStr); err == nil {
		return time.Duration(ms) * time.Millisecond, nil
	}

	return 0, fmt.Errorf("invalid duration format: %s (use formats like '5s', '10m', '1h')", durationStr)
}

// ============================================================================
// Value Comparison Helpers
// ============================================================================

// compareValues compares values for switch cases
func compareValues(a, b interface{}) bool {
	// Simple equality check
	switch aVal := a.(type) {
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal == bVal
		}
	case string:
		if bVal, ok := b.(string); ok {
			return aVal == bVal
		}
	case bool:
		if bVal, ok := b.(bool); ok {
			return aVal == bVal
		}
	}
	return false
}

// ============================================================================
// Test Helpers
// ============================================================================

// intPtr returns a pointer to an int value (for tests)
func intPtr(i int) *int {
	return &i
}

// stringPtr returns a pointer to a string value (for tests)
func stringPtr(s string) *string {
	return &s
}

// // float64Ptr returns a pointer to a float64 value (for tests)
// func float64Ptr(f float64) *float64 {
// 	return &f
// }

// boolPtr returns a pointer to a bool value (for tests)
func boolPtr(b bool) *bool {
	return &b
}
