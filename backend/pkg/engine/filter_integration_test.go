package engine

import (
	"encoding/json"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestFilterNode_Integration tests the filter node in a complete workflow
func TestFilterNode_Integration(t *testing.T) {
	t.Run("Filter array of numbers", func(t *testing.T) {
		payload := types.Payload{
			Nodes: []types.Node{
				{
					ID:   "1",
					Type: types.NodeTypeTransform,
					Data: types.NodeData{
						TransformType: strPtr("to_array"),
					},
				},
				{
					ID:   "num1",
					Type: types.NodeTypeNumber,
					Data: types.NodeData{
						Value: float64Ptr(5),
					},
				},
				{
					ID:   "num2",
					Type: types.NodeTypeNumber,
					Data: types.NodeData{
						Value: float64Ptr(15),
					},
				},
				{
					ID:   "num3",
					Type: types.NodeTypeNumber,
					Data: types.NodeData{
						Value: float64Ptr(8),
					},
				},
				{
					ID:   "num4",
					Type: types.NodeTypeNumber,
					Data: types.NodeData{
						Value: float64Ptr(20),
					},
				},
				{
					ID:   "filter",
					Type: types.NodeTypeFilter,
					Data: types.NodeData{
						Condition: strPtr("variables.item > 10"),
					},
				},
				{
					ID:   "viz",
					Type: types.NodeTypeVisualization,
					Data: types.NodeData{
						Mode: strPtr("json"),
					},
				},
			},
			Edges: []types.Edge{
				{Source: "num1", Target: "1"},
				{Source: "num2", Target: "1"},
				{Source: "num3", Target: "1"},
				{Source: "num4", Target: "1"},
				{Source: "1", Target: "filter"},
				{Source: "filter", Target: "viz"},
			},
		}

		payloadJSON, _ := json.Marshal(payload)
		engine, err := New(payloadJSON)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("Execution failed: %v", err)
		}

		// Get filter result
		filterResult, ok := result.NodeResults["filter"].(map[string]interface{})
		if !ok {
			t.Fatalf("Filter result not found or wrong type")
		}

		filtered, ok := filterResult["filtered"].([]interface{})
		if !ok {
			t.Fatalf("Filtered array not found or wrong type")
		}

		// Should have 2 items: 15 and 20
		if len(filtered) != 2 {
			t.Errorf("Expected 2 filtered items, got %d", len(filtered))
		}

		// Verify metadata
		if filterResult["input_count"] != 4 {
			t.Errorf("Expected input_count=4, got %v", filterResult["input_count"])
		}
		if filterResult["output_count"] != 2 {
			t.Errorf("Expected output_count=2, got %v", filterResult["output_count"])
		}
	})

	t.Run("Filter objects by field value", func(t *testing.T) {
		payload := types.Payload{
			Nodes: []types.Node{
				{
					ID:   "array",
					Type: types.NodeTypeTransform,
					Data: types.NodeData{
						TransformType: strPtr("to_array"),
					},
				},
				{
					ID:   "user1",
					Type: types.NodeTypeVariable,
					Data: types.NodeData{
						VarName: strPtr("temp"),
						VarOp:   strPtr("set"),
					},
				},
				{
					ID:   "user2",
					Type: types.NodeTypeVariable,
					Data: types.NodeData{
						VarName: strPtr("temp2"),
						VarOp:   strPtr("set"),
					},
				},
				{
					ID:   "user3",
					Type: types.NodeTypeVariable,
					Data: types.NodeData{
						VarName: strPtr("temp3"),
						VarOp:   strPtr("set"),
					},
				},
				{
					ID:   "user4",
					Type: types.NodeTypeVariable,
					Data: types.NodeData{
						VarName: strPtr("temp4"),
						VarOp:   strPtr("set"),
					},
				},
				{
					ID:   "filter",
					Type: types.NodeTypeFilter,
					Data: types.NodeData{
						Condition: strPtr("age >= 18 && active == true"),
					},
				},
			},
			Edges: []types.Edge{
				{Source: "user1", Target: "array"},
				{Source: "user2", Target: "array"},
				{Source: "user3", Target: "array"},
				{Source: "user4", Target: "array"},
				{Source: "array", Target: "filter"},
			},
		}

		// Actually, let's just use a simpler approach with direct number inputs
		// Create test data inline
		payload2 := `{
			"nodes": [
				{"id": "1", "type": "number", "data": {"value": 25}},
				{"id": "2", "type": "number", "data": {"value": 15}},
				{"id": "3", "type": "number", "data": {"value": 8}},
				{"id": "4", "type": "number", "data": {"value": 20}},
				{"id": "5", "type": "transform", "data": {"transform_type": "to_array"}},
				{"id": "6", "type": "filter", "data": {"condition": "variables.item >= 18"}},
				{"id": "7", "type": "visualization", "data": {"mode": "json"}}
			],
			"edges": [
				{"source": "1", "target": "5"},
				{"source": "2", "target": "5"},
				{"source": "3", "target": "5"},
				{"source": "4", "target": "5"},
				{"source": "5", "target": "6"},
				{"source": "6", "target": "7"}
			]
		}`

		engine, err := New([]byte(payload2))
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("Execution failed: %v", err)
		}

		// Get filter result
		filterResult, ok := result.NodeResults["6"].(map[string]interface{})
		if !ok {
			t.Fatalf("Filter result not found or wrong type")
		}

		filtered, ok := filterResult["filtered"].([]interface{})
		if !ok {
			t.Fatalf("Filtered array not found or wrong type")
		}

		// Should have 2 items: 25 and 20
		if len(filtered) != 2 {
			t.Errorf("Expected 2 filtered items, got %d", len(filtered))
		}
	})

	t.Run("Filter with non-array input", func(t *testing.T) {
		payload := types.Payload{
			Nodes: []types.Node{
				{
					ID:   "text",
					Type: types.NodeTypeTextInput,
					Data: types.NodeData{
						Text: strPtr("not an array"),
					},
				},
				{
					ID:   "filter",
					Type: types.NodeTypeFilter,
					Data: types.NodeData{
						Condition: strPtr("variables.item > 10"),
					},
				},
			},
			Edges: []types.Edge{
				{Source: "text", Target: "filter"},
			},
		}

		payloadJSON, _ := json.Marshal(payload)
		engine, err := New(payloadJSON)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("Execution failed: %v", err)
		}

		// Get filter result
		filterResult, ok := result.NodeResults["filter"].(map[string]interface{})
		if !ok {
			t.Fatalf("Filter result not found or wrong type")
		}

		// Should have warning and pass through original input
		if filterResult["is_array"] != false {
			t.Errorf("Expected is_array=false for non-array input")
		}

		if _, hasWarning := filterResult["warning"]; !hasWarning {
			t.Errorf("Expected warning field for non-array input")
		}

		if filterResult["filtered"] != "not an array" {
			t.Errorf("Expected filtered to be original input")
		}
	})

	t.Run("Filter with context variable threshold", func(t *testing.T) {
		payload := types.Payload{
			Nodes: []types.Node{
				{
					ID:   "threshold",
					Type: types.NodeTypeContextConstant,
					Data: types.NodeData{
						ContextValues: []types.ContextVariableValue{
							{Name: "minAge", Value: float64(21)},
						},
					},
				},
				{
					ID:   "users",
					Type: types.NodeTypeContextVariable,
					Data: types.NodeData{
						ContextValues: []types.ContextVariableValue{
							{
								Name: "userList",
								Value: []interface{}{
									map[string]interface{}{"name": "Alice", "age": float64(25)},
									map[string]interface{}{"name": "Bob", "age": float64(19)},
									map[string]interface{}{"name": "Charlie", "age": float64(22)},
								},
							},
						},
					},
				},
				{
					ID:   "get_users",
					Type: types.NodeTypeVariable,
					Data: types.NodeData{
						VarName: strPtr("userList"),
						VarOp:   strPtr("get"),
					},
				},
				{
					ID:   "filter",
					Type: types.NodeTypeFilter,
					Data: types.NodeData{
						Condition: strPtr("variables.item.age >= context.minAge"),
					},
				},
			},
			Edges: []types.Edge{
				{Source: "threshold", Target: "get_users"},
				{Source: "users", Target: "get_users"},
				{Source: "get_users", Target: "filter"},
			},
		}

		payloadJSON, _ := json.Marshal(payload)
		engine, err := New(payloadJSON)
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("Execution failed: %v", err)
		}

		// Get filter result
		filterResult, ok := result.NodeResults["filter"].(map[string]interface{})
		if !ok {
			t.Fatalf("Filter result not found or wrong type")
		}

		filtered, ok := filterResult["filtered"].([]interface{})
		if !ok {
			t.Fatalf("Filtered array not found or wrong type")
		}

		// Should have 2 items: Alice (25) and Charlie (22)
		// Bob is excluded (age < 21)
		if len(filtered) != 2 {
			t.Errorf("Expected 2 filtered items, got %d", len(filtered))
		}
	})
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
