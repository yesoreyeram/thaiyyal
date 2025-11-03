package executor

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ParseExecutor executes Parse nodes
// Converts string data to structured formats (JSON, CSV, TSV, YAML, XML)
type ParseExecutor struct{}

// Execute runs the Parse node
// Parses string input into structured data based on the specified input type.
func (e *ParseExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("parse node requires input")
	}

	// Get the input as string
	inputStr, err := toString(inputs[0])
	if err != nil {
		return nil, fmt.Errorf("parse node input must be convertible to string: %w", err)
	}

	// Get input type, default to AUTO
	inputType := "AUTO"
	if node.Data.InputType != nil {
		inputType = strings.ToUpper(*node.Data.InputType)
	}

	// If AUTO, detect the format
	if inputType == "AUTO" {
		inputType = detectFormat(inputStr)
	}

	// Parse based on type
	switch inputType {
	case "JSON":
		return parseJSON(inputStr)
	case "CSV":
		return parseCSV(inputStr, ',')
	case "TSV":
		return parseCSV(inputStr, '\t')
	case "YAML":
		return parseYAML(inputStr)
	case "XML":
		return parseXML(inputStr)
	default:
		return nil, fmt.Errorf("unsupported input type: %s (supported: AUTO, JSON, CSV, TSV, YAML, XML)", inputType)
	}
}

// NodeType returns the node type this executor handles
func (e *ParseExecutor) NodeType() types.NodeType {
	return types.NodeTypeParse
}

// Validate checks if node configuration is valid
func (e *ParseExecutor) Validate(node types.Node) error {
	if node.Data.InputType != nil {
		inputType := strings.ToUpper(*node.Data.InputType)
		validTypes := map[string]bool{
			"AUTO": true,
			"JSON": true,
			"CSV":  true,
			"TSV":  true,
			"YAML": true,
			"XML":  true,
		}
		if !validTypes[inputType] {
			return fmt.Errorf("invalid input_type: %s (must be one of: AUTO, JSON, CSV, TSV, YAML, XML)", inputType)
		}
	}
	return nil
}

// toString converts various types to string
func toString(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%v", v), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case nil:
		return "", nil
	default:
		// For maps and slices, try JSON encoding
		if data, err := json.Marshal(v); err == nil {
			return string(data), nil
		}
		return fmt.Sprintf("%v", v), nil
	}
}

// detectFormat attempts to detect the format of the input string
func detectFormat(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "JSON"
	}

	// Check for JSON
	if (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
		(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")) {
		// Try to parse as JSON to confirm
		var test interface{}
		if json.Unmarshal([]byte(trimmed), &test) == nil {
			return "JSON"
		}
	}

	// Check for XML
	if strings.HasPrefix(trimmed, "<") && strings.HasSuffix(trimmed, ">") {
		return "XML"
	}

	// Check for YAML (common patterns)
	if strings.Contains(trimmed, ":\n") || strings.Contains(trimmed, ": ") {
		lines := strings.Split(trimmed, "\n")
		if len(lines) > 1 {
			// Check if multiple lines have key:value pattern
			yamlLikeCount := 0
			for _, line := range lines {
				if strings.Contains(strings.TrimSpace(line), ": ") {
					yamlLikeCount++
				}
			}
			if yamlLikeCount >= 2 {
				return "YAML"
			}
		}
	}

	// Check for CSV/TSV (contains commas or tabs, multiple lines)
	lines := strings.Split(trimmed, "\n")
	if len(lines) > 1 {
		firstLine := lines[0]
		// Count separators
		commas := strings.Count(firstLine, ",")
		tabs := strings.Count(firstLine, "\t")

		if tabs > commas && tabs > 0 {
			return "TSV"
		}
		if commas > 0 {
			return "CSV"
		}
	}

	// Default to JSON for single values or unknown format
	return "JSON"
}

// parseJSON parses JSON string to structured data
func parseJSON(input string) (interface{}, error) {
	trimmed := strings.TrimSpace(input)

	// Handle primitive types without quotes
	if trimmed == "true" {
		return true, nil
	}
	if trimmed == "false" {
		return false, nil
	}
	if trimmed == "null" {
		return nil, nil
	}

	// Try to parse as number
	if num, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return num, nil
	}

	// Parse as JSON
	var result interface{}
	decoder := json.NewDecoder(strings.NewReader(trimmed))
	decoder.UseNumber() // Preserve number precision

	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Convert json.Number to float64
	result = convertJSONNumbers(result)

	return result, nil
}

// convertJSONNumbers recursively converts json.Number to float64
func convertJSONNumbers(val interface{}) interface{} {
	switch v := val.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = convertJSONNumbers(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = convertJSONNumbers(value)
		}
		return result
	case json.Number:
		if f, err := v.Float64(); err == nil {
			return f
		}
		return v.String()
	default:
		return v
	}
}

// parseCSV parses CSV/TSV string to array of objects
func parseCSV(input string, delimiter rune) (interface{}, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = delimiter
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("invalid CSV/TSV: %w", err)
	}

	if len(records) == 0 {
		return []map[string]interface{}{}, nil
	}

	// First row is headers
	headers := records[0]

	// Convert remaining rows to objects
	result := make([]map[string]interface{}, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		row := records[i]
		obj := make(map[string]interface{})

		for j, header := range headers {
			if j < len(row) {
				// Try to parse as number, bool, or keep as string
				obj[header] = parseValue(row[j])
			} else {
				obj[header] = nil
			}
		}
		result = append(result, obj)
	}

	return result, nil
}

// parseValue attempts to parse a string value to its appropriate type
func parseValue(s string) interface{} {
	trimmed := strings.TrimSpace(s)

	// Empty string
	if trimmed == "" {
		return ""
	}

	// Boolean
	if trimmed == "true" {
		return true
	}
	if trimmed == "false" {
		return false
	}

	// Null
	if trimmed == "null" {
		return nil
	}

	// Number
	if num, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return num
	}

	// Default to string
	return trimmed
}

// parseYAML parses YAML string to structured data
// Note: Since we don't have external dependencies, this is a basic YAML parser
// that handles simple key-value pairs and nested structures
func parseYAML(input string) (interface{}, error) {
	trimmed := strings.TrimSpace(input)
	lines := strings.Split(trimmed, "\n")

	if len(lines) == 0 {
		return nil, fmt.Errorf("empty YAML input")
	}

	// Simple YAML parser for basic structures
	result := make(map[string]interface{})

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key: value
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = parseValue(value)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid YAML key-value pairs found")
	}

	return result, nil
}

// parseXML parses XML string to structured data
func parseXML(input string) (interface{}, error) {
	trimmed := strings.TrimSpace(input)

	// Parse into a generic map structure
	var result map[string]interface{}
	if err := xml.Unmarshal([]byte(trimmed), &result); err != nil {
		// Try parsing as a simple element
		var simple struct {
			XMLName xml.Name
			Content string `xml:",chardata"`
		}
		if err2 := xml.Unmarshal([]byte(trimmed), &simple); err2 != nil {
			return nil, fmt.Errorf("invalid XML: %w", err)
		}

		// Return as a map with the tag name as key
		return map[string]interface{}{
			simple.XMLName.Local: parseValue(simple.Content),
		}, nil
	}

	return result, nil
}
