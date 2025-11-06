package executor

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// FormatExecutor executes Format nodes
// Converts structured data to string formats (JSON, CSV, TSV)
type FormatExecutor struct{}

// Execute runs the Format node
// Formats structured input data into string output based on the specified output type.
func (e *FormatExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsFormatData(node.Data)
if err != nil {
return nil, err
}
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("format node requires input")
	}

	input := inputs[0]

	// Get output type, default to JSON
	outputType := "JSON"
	if data.OutputType != nil {
		outputType = strings.ToUpper(*data.OutputType)
	}

	// Get formatting options
	prettyPrint := false
	if data.PrettyPrint != nil {
		prettyPrint = *data.PrettyPrint
	}

	includeHeaders := true
	if data.IncludeHeaders != nil {
		includeHeaders = *data.IncludeHeaders
	}

	delimiter := ','
	if data.Delimiter != nil {
		if len(*data.Delimiter) > 0 {
			delimiter = rune((*data.Delimiter)[0])
		}
	}

	// Format based on type
	switch outputType {
	case "JSON":
		return formatJSON(input, prettyPrint)
	case "CSV":
		return formatCSV(input, delimiter, includeHeaders)
	case "TSV":
		return formatCSV(input, '\t', includeHeaders)
	default:
		return nil, fmt.Errorf("unsupported output type: %s (supported: JSON, CSV, TSV)", outputType)
	}
}

// NodeType returns the node type this executor handles
func (e *FormatExecutor) NodeType() types.NodeType {
	return types.NodeTypeFormat
}

// Validate checks if node configuration is valid
func (e *FormatExecutor) Validate(node types.Node) error {
data, err := types.AsFormatData(node.Data)
if err != nil {
return err
}
	if data.OutputType != nil {
		outputType := strings.ToUpper(*data.OutputType)
		validTypes := map[string]bool{
			"JSON": true,
			"CSV":  true,
			"TSV":  true,
		}
		if !validTypes[outputType] {
			return fmt.Errorf("invalid output_type: %s (must be one of: JSON, CSV, TSV)", outputType)
		}
	}
	return nil
}

// formatJSON converts data to JSON string
func formatJSON(data interface{}, prettyPrint bool) (interface{}, error) {
	var output []byte
	var err error

	if prettyPrint {
		output, err = json.MarshalIndent(data, "", "  ")
	} else {
		output, err = json.Marshal(data)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to format as JSON: %w", err)
	}

	return string(output), nil
}

// formatCSV converts array of objects to CSV string
func formatCSV(data interface{}, delimiter rune, includeHeaders bool) (interface{}, error) {
	// Convert to array if not already
	var records []map[string]interface{}

	switch v := data.(type) {
	case []interface{}:
		// Convert each element to map
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				records = append(records, m)
			} else {
				return nil, fmt.Errorf("CSV formatting requires array of objects, got array element of type %T", item)
			}
		}
	case []map[string]interface{}:
		records = v
	case map[string]interface{}:
		// Single object - treat as single-row CSV
		records = []map[string]interface{}{v}
	default:
		return nil, fmt.Errorf("CSV formatting requires array of objects or single object, got %T", data)
	}

	if len(records) == 0 {
		return "", nil
	}

	// Collect all unique headers from all records
	headersMap := make(map[string]bool)
	for _, record := range records {
		for key := range record {
			headersMap[key] = true
		}
	}

	// Convert headers map to ordered slice (deterministic order)
	headers := make([]string, 0, len(headersMap))
	for key := range headersMap {
		headers = append(headers, key)
	}

	// Sort headers for consistent output (simple bubble sort)
	for i := 0; i < len(headers)-1; i++ {
		for j := i + 1; j < len(headers); j++ {
			if headers[i] > headers[j] {
				headers[i], headers[j] = headers[j], headers[i]
			}
		}
	}

	// Create CSV writer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = delimiter

	// Write headers if requested
	if includeHeaders {
		if err := writer.Write(headers); err != nil {
			return nil, fmt.Errorf("failed to write CSV headers: %w", err)
		}
	}

	// Write data rows
	for _, record := range records {
		row := make([]string, len(headers))
		for i, header := range headers {
			if val, exists := record[header]; exists && val != nil {
				row[i] = formatValue(val)
			} else {
				row[i] = ""
			}
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.String(), nil
}

// formatValue converts a value to its string representation for CSV
func formatValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		// Format numbers without unnecessary decimals
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%v", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		// For complex types, use JSON encoding
		if data, err := json.Marshal(v); err == nil {
			return string(data)
		}
		return fmt.Sprintf("%v", v)
	}
}
