package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestGroupByExecutor_Basic(t *testing.T) {
	tests := []struct {
		name           string
		inputArray     []interface{}
		field          string
		aggregate      *string
		valueField     *string
		expectedGroups int
		description    string
	}{
		{
			name: "Group by status with count",
			inputArray: []interface{}{
				map[string]interface{}{"name": "Alice", "status": "active"},
				map[string]interface{}{"name": "Bob", "status": "inactive"},
				map[string]interface{}{"name": "Charlie", "status": "active"},
			},
			field:          "status",
			aggregate:      stringPtr("count"),
			expectedGroups: 2,
			description:    "Should group by status and count",
		},
		{
			name: "Group by category with sum",
			inputArray: []interface{}{
				map[string]interface{}{"category": "A", "value": 10.0},
				map[string]interface{}{"category": "B", "value": 20.0},
				map[string]interface{}{"category": "A", "value": 15.0},
			},
			field:          "category",
			aggregate:      stringPtr("sum"),
			valueField:     stringPtr("value"),
			expectedGroups: 2,
			description:    "Should group by category and sum values",
		},
		{
			name: "Group by type with values",
			inputArray: []interface{}{
				map[string]interface{}{"type": "X", "data": "d1"},
				map[string]interface{}{"type": "Y", "data": "d2"},
				map[string]interface{}{"type": "X", "data": "d3"},
			},
			field:          "type",
			aggregate:      stringPtr("values"),
			expectedGroups: 2,
			description:    "Should group and return all values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &GroupByExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
			}

			nodeData := types.ExtractData{
				Field: &tt.field,
			}
			if tt.aggregate != nil {
				nodeData.Aggregate = tt.aggregate
			}
			if tt.valueField != nil {
				nodeData.ValueField = tt.valueField
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeGroupBy,
				Data: nodeData,
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected result to be map, got %T", result)
			}

			groupCount, ok := resultMap["group_count"].(int)
			if !ok {
				t.Fatalf("Expected 'group_count' to be int, got %T", resultMap["group_count"])
			}

			if groupCount != tt.expectedGroups {
				t.Errorf("Expected %d groups, got %d", tt.expectedGroups, groupCount)
			}
		})
	}
}

func TestGroupByExecutor_Validate(t *testing.T) {
	exec := &GroupByExecutor{}

	tests := []struct {
		name        string
		field       *string
		aggregate   *string
		valueField  *string
		expectError bool
	}{
		{"Valid count aggregate", stringPtr("category"), stringPtr("count"), nil, false},
		{"Valid sum aggregate", stringPtr("category"), stringPtr("sum"), stringPtr("value"), false},
		{"Missing field", nil, stringPtr("count"), nil, true},
		{"Invalid aggregate", stringPtr("category"), stringPtr("invalid"), nil, true},
		{"Sum without valueField", stringPtr("category"), stringPtr("sum"), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeData := types.NodeData{}
			if tt.field != nil {
				nodeData.Field = tt.field
			}
			if tt.aggregate != nil {
				nodeData.Aggregate = tt.aggregate
			}
			if tt.valueField != nil {
				nodeData.ValueField = tt.valueField
			}

			node := types.Node{
				Type: types.NodeTypeGroupBy,
				Data: nodeData,
			}

			err := exec.Validate(node)
			if tt.expectError && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}
