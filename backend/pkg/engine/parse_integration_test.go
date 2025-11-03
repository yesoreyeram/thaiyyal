package engine

import (
	"encoding/json"
	"testing"
)

// TestParseNode_Integration tests the parse node in an integrated workflow
func TestParseNode_Integration(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		checkNodeID    string // Which node to check the result from (defaults to "2" - the parse node)
		expectedOutput interface{}
		expectError    bool
	}{
		{
			name: "Parse JSON object",
			payload: `{
				"nodes": [
					{"id": "1", "type": "text_input", "data": {"text": "{\"name\": \"Alice\", \"age\": 30}"}},
					{"id": "2", "type": "parse", "data": {"input_type": "JSON"}},
					{"id": "3", "type": "visualization", "data": {"mode": "text"}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"},
					{"id": "e2", "source": "2", "target": "3"}
				]
			}`,
			expectedOutput: map[string]interface{}{
				"name": "Alice",
				"age":  float64(30),
			},
			expectError: false,
		},
		{
			name: "Parse CSV",
			payload: `{
				"nodes": [
					{"id": "1", "type": "text_input", "data": {"text": "name,age,city\nJohn,25,NYC\nJane,30,LA"}},
					{"id": "2", "type": "parse", "data": {"input_type": "CSV"}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"}
				]
			}`,
			expectedOutput: []map[string]interface{}{
				{"name": "John", "age": float64(25), "city": "NYC"},
				{"name": "Jane", "age": float64(30), "city": "LA"},
			},
			expectError: false,
		},
		{
			name: "Parse with AUTO detection - JSON",
			payload: `{
				"nodes": [
					{"id": "1", "type": "text_input", "data": {"text": "[1, 2, 3, 4, 5]"}},
					{"id": "2", "type": "parse", "data": {"input_type": "AUTO"}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"}
				]
			}`,
			expectedOutput: []interface{}{float64(1), float64(2), float64(3), float64(4), float64(5)},
			expectError:    false,
		},
		{
			name: "Parse with default AUTO detection",
			payload: `{
				"nodes": [
					{"id": "1", "type": "text_input", "data": {"text": "true"}},
					{"id": "2", "type": "parse", "data": {}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"}
				]
			}`,
			expectedOutput: true,
			expectError:    false,
		},
		{
			name: "Parse TSV",
			payload: `{
				"nodes": [
					{"id": "1", "type": "text_input", "data": {"text": "product\tprice\niPhone\t999.99\niPad\t799.99"}},
					{"id": "2", "type": "parse", "data": {"input_type": "TSV"}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"}
				]
			}`,
			expectedOutput: []map[string]interface{}{
				{"product": "iPhone", "price": float64(999.99)},
				{"product": "iPad", "price": float64(799.99)},
			},
			expectError: false,
		},
		{
			name: "Parse YAML",
			payload: `{
				"nodes": [
					{"id": "1", "type": "text_input", "data": {"text": "name: Bob\nage: 42\nactive: true"}},
					{"id": "2", "type": "parse", "data": {"input_type": "YAML"}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"}
				]
			}`,
			expectedOutput: map[string]interface{}{
				"name":   "Bob",
				"age":    float64(42),
				"active": true,
			},
			expectError: false,
		},
		{
			name: "Parse number string to number",
			payload: `{
				"nodes": [
					{"id": "1", "type": "text_input", "data": {"text": "42.5"}},
					{"id": "2", "type": "parse", "data": {"input_type": "JSON"}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"}
				]
			}`,
			expectedOutput: float64(42.5),
			expectError:    false,
		},
		{
			name: "Chain Parse with Extract",
			payload: `{
				"nodes": [
					{"id": "1", "type": "text_input", "data": {"text": "{\"user\": {\"name\": \"Charlie\", \"email\": \"charlie@example.com\"}}"}},
					{"id": "2", "type": "parse", "data": {"input_type": "JSON"}},
					{"id": "3", "type": "extract", "data": {"field": "user"}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"},
					{"id": "e2", "source": "2", "target": "3"}
				]
			}`,
			checkNodeID: "3", // Check extract node result
			expectedOutput: map[string]interface{}{
				"field": "user",
				"value": map[string]interface{}{
					"name":  "Charlie",
					"email": "charlie@example.com",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := New([]byte(tt.payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			result, err := engine.Execute()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Find the parse node result
			checkNodeID := tt.checkNodeID
			if checkNodeID == "" {
				checkNodeID = "2" // Default to parse node
			}

			var parseNodeID string
			for nodeID := range result.NodeResults {
				if nodeID == checkNodeID {
					parseNodeID = nodeID
					break
				}
			}

			if parseNodeID == "" {
				t.Fatal("Parse node result not found")
			}

			actualResult := result.NodeResults[parseNodeID]

			// Compare results
			if !compareResults(actualResult, tt.expectedOutput) {
				actualJSON, _ := json.MarshalIndent(actualResult, "", "  ")
				expectedJSON, _ := json.MarshalIndent(tt.expectedOutput, "", "  ")
				t.Errorf("Parse result mismatch.\nActual:\n%s\n\nExpected:\n%s", actualJSON, expectedJSON)
			}
		})
	}
}

// compareResults compares two values for equality
func compareResults(a, b interface{}) bool {
	aJSON, err1 := json.Marshal(a)
	bJSON, err2 := json.Marshal(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return string(aJSON) == string(bJSON)
}
