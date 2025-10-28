package workflow

import (
	"encoding/json"
	"testing"
)

// TestExplicitNodeTypes tests workflows with explicit node type fields
func TestExplicitNodeTypes(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 10}},
			{"id": "2", "type": "number", "data": {"value": 5}},
			{"id": "3", "type": "operation", "data": {"op": "add"}},
			{"id": "4", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-4", "source": "3", "target": "4"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	// Verify execution
	if result.NodeResults["3"] != 15.0 {
		t.Errorf("Expected node 3 result to be 15, got %v", result.NodeResults["3"])
	}
}

// TestNodeTypeInference tests that node types are correctly inferred when not explicit
func TestNodeTypeInference(t *testing.T) {
	tests := []struct {
		name         string
		payload      string
		nodeID       string
		expectedType NodeType
	}{
		{
			name: "infer number type",
			payload: `{
				"nodes": [{"id": "1", "data": {"value": 42}}],
				"edges": []
			}`,
			nodeID:       "1",
			expectedType: NodeTypeNumber,
		},
		{
			name: "infer operation type",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 5}},
					{"id": "3", "data": {"op": "add"}}
				],
				"edges": [
					{"id": "e1-3", "source": "1", "target": "3"},
					{"id": "e2-3", "source": "2", "target": "3"}
				]
			}`,
			nodeID:       "3",
			expectedType: NodeTypeOperation,
		},
		{
			name: "infer visualization type",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 42}},
					{"id": "2", "data": {"mode": "text"}}
				],
				"edges": [{"id": "e1-2", "source": "1", "target": "2"}]
			}`,
			nodeID:       "2",
			expectedType: NodeTypeVisualization,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine([]byte(tt.payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			// Check that type was inferred
			var node *Node
			for i := range engine.payload.Nodes {
				if engine.payload.Nodes[i].ID == tt.nodeID {
					node = &engine.payload.Nodes[i]
					break
				}
			}

			if node == nil {
				t.Fatalf("Node %s not found", tt.nodeID)
			}

			// Execute to trigger type inference
			_, err = engine.Execute()
			if err != nil {
				t.Fatalf("Failed to execute: %v", err)
			}

			if node.Type != tt.expectedType {
				t.Errorf("Expected node type %s, got %s", tt.expectedType, node.Type)
			}
		})
	}
}

// TestMixedExplicitAndInferredTypes tests workflows with both explicit and inferred types
func TestMixedExplicitAndInferredTypes(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "type": "operation", "data": {"op": "multiply"}},
			{"id": "4", "data": {"mode": "table"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-4", "source": "3", "target": "4"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	// Verify execution
	if result.NodeResults["3"] != 50.0 {
		t.Errorf("Expected node 3 result to be 50, got %v", result.NodeResults["3"])
	}

	// Verify visualization
	vizResult, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected visualization result to be a map")
	}

	if vizResult["mode"] != "table" {
		t.Errorf("Expected mode 'table', got %v", vizResult["mode"])
	}
}

// TestNodeTypeConstants verifies the node type constants are properly defined
func TestNodeTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		nodeType NodeType
		expected string
	}{
		{"number type", NodeTypeNumber, "number"},
		{"operation type", NodeTypeOperation, "operation"},
		{"visualization type", NodeTypeVisualization, "visualization"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.nodeType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.nodeType))
			}
		})
	}
}

// TestOperationTypeConstants verifies the operation type constants
func TestOperationTypeConstants(t *testing.T) {
	tests := []struct {
		name      string
		operation OperationType
		expected  string
	}{
		{"add", OperationAdd, "add"},
		{"subtract", OperationSubtract, "subtract"},
		{"multiply", OperationMultiply, "multiply"},
		{"divide", OperationDivide, "divide"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.operation) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.operation))
			}
		})
	}
}

// TestVisualizationModeConstants verifies the visualization mode constants
func TestVisualizationModeConstants(t *testing.T) {
	tests := []struct {
		name     string
		mode     VisualizationMode
		expected string
	}{
		{"text mode", VisualizationModeText, "text"},
		{"table mode", VisualizationModeTable, "table"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.mode) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.mode))
			}
		})
	}
}

// TestTypedNodeDataStructures tests the typed node data structures
func TestTypedNodeDataStructures(t *testing.T) {
	t.Run("NumberNodeData", func(t *testing.T) {
		data := NumberNodeData{
			Value: 42.5,
			Label: "Test Number",
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var unmarshaled NumberNodeData
		if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if unmarshaled.Value != 42.5 {
			t.Errorf("Expected value 42.5, got %v", unmarshaled.Value)
		}
		if unmarshaled.Label != "Test Number" {
			t.Errorf("Expected label 'Test Number', got %s", unmarshaled.Label)
		}
	})

	t.Run("OperationNodeData", func(t *testing.T) {
		data := OperationNodeData{
			Operation: OperationAdd,
			Label:     "Test Operation",
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var unmarshaled OperationNodeData
		if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if unmarshaled.Operation != OperationAdd {
			t.Errorf("Expected operation 'add', got %s", unmarshaled.Operation)
		}
	})

	t.Run("VisualizationNodeData", func(t *testing.T) {
		data := VisualizationNodeData{
			Mode:  VisualizationModeText,
			Label: "Test Viz",
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var unmarshaled VisualizationNodeData
		if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if unmarshaled.Mode != VisualizationModeText {
			t.Errorf("Expected mode 'text', got %s", unmarshaled.Mode)
		}
	})
}

// TestInvalidNodeType tests error handling for invalid node types
func TestInvalidNodeType(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "type": "invalid", "data": {"value": 10}}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	_, err = engine.Execute()
	if err == nil {
		t.Error("Expected error for invalid node type, got nil")
	}
}
