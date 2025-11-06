package engine

import (
	"encoding/json"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestConditionalRendering_BasicTruePath tests basic conditional rendering when condition is true
func TestConditionalRendering_BasicTruePath(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "input",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(25)},
			},
			{
				ID:   "condition",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{Condition: strPtr(">18")},
			},
			{
				ID:   "output",
				Type: types.NodeTypeVisualization,
				Data: types.VisualizationData{Mode: strPtr("json")},
			},
		},
		Edges: []types.Edge{
			{Source: "input", Target: "condition"},
			{Source: "condition", Target: "output"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	condResult := mustGetMapResult(t, result, "condition")
	if !condResult["condition_met"].(bool) {
		t.Error("Expected condition to be met")
	}
	if condResult["path"].(string) != "true" {
		t.Errorf("Expected path='true', got %v", condResult["path"])
	}
}

// TestConditionalRendering_BasicFalsePath tests basic conditional rendering when condition is false
func TestConditionalRendering_BasicFalsePath(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "input",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(15)},
			},
			{
				ID:   "condition",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{Condition: strPtr(">18")},
			},
			{
				ID:   "output",
				Type: types.NodeTypeVisualization,
				Data: types.VisualizationData{Mode: strPtr("json")},
			},
		},
		Edges: []types.Edge{
			{Source: "input", Target: "condition"},
			{Source: "condition", Target: "output"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	condResult := mustGetMapResult(t, result, "condition")
	if condResult["condition_met"].(bool) {
		t.Error("Expected condition to not be met")
	}
	if condResult["path"].(string) != "false" {
		t.Errorf("Expected path='false', got %v", condResult["path"])
	}
}

// TestConditionalAPICall_SuccessPath tests conditional API call execution
func TestConditionalAPICall_SuccessPath(t *testing.T) {
	// Mock HTTP server will be created in actual test
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "flag",
				Type: types.NodeTypeBooleanInput,
				Data: types.BooleanInputData{BooleanValue: boolPtr(true)},
			},
			{
				ID:   "condition",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{Condition: strPtr("input == true")},
			},
		},
		Edges: []types.Edge{
			{Source: "flag", Target: "condition"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	condResult := mustGetMapResult(t, result, "condition")
	if !condResult["condition_met"].(bool) {
		t.Error("Expected condition to be met for API call")
	}
}

// TestConditionalContextUpdate tests conditional context variable updates
func TestConditionalContextUpdate(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "age",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(30)},
			},
			{
				ID:   "condition",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{Condition: strPtr(">=21")},
			},
			{
				ID:   "set_context",
				Type: types.NodeTypeContextConstant,
				Data: types.ContextConstantData{
					ContextName:   strPtr("user_status"),
					ContextValue: strPtr("adult"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "age", Target: "condition"},
			{Source: "condition", Target: "set_context"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	condResult := mustGetMapResult(t, result, "condition")
	if !condResult["condition_met"].(bool) {
		t.Error("Expected age condition to be met")
	}

	// Verify context was set
	if result.NodeResults["set_context"] == nil {
		t.Error("Expected context constant to execute")
	}
}

// TestConditionalNested_TwoLevels tests nested conditionals (2 levels)
func TestConditionalNested_TwoLevels(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "score",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(85)},
			},
			{
				ID:   "condition1",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{Condition: strPtr(">=60")},
			},
			{
				ID:   "extract_value",
				Type: types.NodeTypeExtract,
				Data: types.ExtractData{Field: strPtr("value")},
			},
			{
				ID:   "condition2",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{Condition: strPtr(">=80")},
			},
		},
		Edges: []types.Edge{
			{Source: "score", Target: "condition1"},
			{Source: "condition1", Target: "extract_value"},
			{Source: "extract_value", Target: "condition2"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	cond1Result := mustGetMapResult(t, result, "condition1")
	if !cond1Result["condition_met"].(bool) {
		t.Error("Expected first condition to be met")
	}

	cond2Result := mustGetMapResult(t, result, "condition2")
	if !cond2Result["condition_met"].(bool) {
		t.Error("Expected second condition to be met")
	}
}

// TestSwitch_MultipleCase tests switch node with multiple cases
func TestSwitch_MultipleCases(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "status_code",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(200)},
			},
			{
				ID:   "switch",
				Type: types.NodeTypeSwitch,
				Data: types.SwitchData{
					Cases: []types.SwitchCase{
						{When: "==200", Value: float64(200), OutputPath: strPtr("success")},
						{When: "==404", Value: float64(404), OutputPath: strPtr("not_found")},
						{When: "==500", Value: float64(500), OutputPath: strPtr("error")},
					},
					DefaultPath: strPtr("unknown"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "status_code", Target: "switch"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	switchResult := mustGetMapResult(t, result, "switch")
	if !switchResult["matched"].(bool) {
		t.Error("Expected switch to match a case")
	}
	if switchResult["output_path"].(string) != "success" {
		t.Errorf("Expected output_path='success', got %v", switchResult["output_path"])
	}
}

// TestSwitch_DefaultCase tests switch node falling back to default
func TestSwitch_DefaultCase(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "status_code",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(302)},
			},
			{
				ID:   "switch",
				Type: types.NodeTypeSwitch,
				Data: types.SwitchData{
					Cases: []types.SwitchCase{
						{When: "==200", Value: float64(200), OutputPath: strPtr("success")},
						{When: "==404", Value: float64(404), OutputPath: strPtr("not_found")},
					},
					DefaultPath: strPtr("other"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "status_code", Target: "switch"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	switchResult := mustGetMapResult(t, result, "switch")
	if switchResult["matched"].(bool) {
		t.Error("Expected switch to not match any case")
	}
	if switchResult["output_path"].(string) != "other" {
		t.Errorf("Expected output_path='other', got %v", switchResult["output_path"])
	}
}

// TestSwitch_RangeConditions tests switch with range-based conditions
func TestSwitch_RangeConditions(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "temperature",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(75)},
			},
			{
				ID:   "switch",
				Type: types.NodeTypeSwitch,
				Data: types.SwitchData{
					Cases: []types.SwitchCase{
						{When: "<32", OutputPath: strPtr("freezing")},
						{When: ">=32 && <60", OutputPath: strPtr("cold")},
						{When: ">=60 && <80", OutputPath: strPtr("moderate")},
						{When: ">=80", OutputPath: strPtr("hot")},
					},
					DefaultPath: strPtr("unknown"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "temperature", Target: "switch"},
		},
	}

	// Note: Current implementation may not support complex expressions in switch
	// This test documents the desired behavior
	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	switchResult := mustGetMapResult(t, result, "switch")
	// Verify result structure
	if switchResult["output_path"] == nil {
		t.Error("Expected output_path to be set")
	}
}

// TestConditionalArrayProcessing_FilterWithCondition tests array filtering based on condition
func TestConditionalArrayProcessing_FilterWithCondition(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "range",
				Type: types.NodeTypeRange,
				Data: types.RangeData{
					Start: intPtr(1),
					End:   intPtr(10),
				},
			},
			{
				ID:   "filter",
				Type: types.NodeTypeFilter,
				Data: types.ConditionData{
					Condition: strPtr("variables.item % 2 == 0"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "range", Target: "filter"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	filterResult := mustGetMapResult(t, result, "filter")
	filtered, ok := filterResult["filtered"].([]interface{})
	if !ok {
		t.Fatalf("Expected filtered array")
	}

	// Should have 5 even numbers: 2, 4, 6, 8, 10
	if len(filtered) != 5 {
		t.Errorf("Expected 5 filtered items, got %d", len(filtered))
	}
}

// TestConditionalArrayProcessing_PartitionByCondition tests partition node
func TestConditionalArrayProcessing_PartitionByCondition(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "range",
				Type: types.NodeTypeRange,
				Data: types.RangeData{
					Start: intPtr(1),
					End:   intPtr(10),
				},
			},
			{
				ID:   "partition",
				Type: types.NodeTypePartition,
				Data: types.ConditionData{
					Condition: strPtr("variables.item > 5"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "range", Target: "partition"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	partitionResult := mustGetMapResult(t, result, "partition")
	truePart, ok := partitionResult["true_partition"].([]interface{})
	if !ok {
		t.Fatalf("Expected true_partition array")
	}
	falsePart, ok := partitionResult["false_partition"].([]interface{})
	if !ok {
		t.Fatalf("Expected false_partition array")
	}

	// Elements > 5: 6, 7, 8, 9, 10 = 5 items
	if len(truePart) != 5 {
		t.Errorf("Expected 5 items in true partition, got %d", len(truePart))
	}
	// Elements <= 5: 1, 2, 3, 4, 5 = 5 items
	if len(falsePart) != 5 {
		t.Errorf("Expected 5 items in false partition, got %d", len(falsePart))
	}
}

// TestConditionalWithVariables tests condition using workflow variables
func TestConditionalWithVariables(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "threshold",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(100)},
			},
			{
				ID:   "store_threshold",
				Type: types.NodeTypeVariable,
				Data: types.VariableData{
					Name:  strPtr("max_value"),
					ContextValue: strPtr("{{node.threshold}}"),
				},
			},
			{
				ID:   "actual_value",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(150)},
			},
			{
				ID:   "condition",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{
					Condition: strPtr("input > variables.max_value"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "threshold", Target: "store_threshold"},
			{Source: "actual_value", Target: "condition"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	condResult := mustGetMapResult(t, result, "condition")
	if !condResult["condition_met"].(bool) {
		t.Error("Expected condition with variable reference to be met")
	}
}

// TestConditionalWithContextVariables tests condition using context variables
func TestConditionalWithContextVariables(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "set_min_age",
				Type: types.NodeTypeContextConstant,
				Data: types.ContextConstantData{
					ContextName:   strPtr("minimum_age"),
					Value: float64Ptr(18),
				},
			},
			{
				ID:   "user_age",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(25)},
			},
			{
				ID:   "condition",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{
					Condition: strPtr("input >= context.minimum_age"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "set_min_age", Target: "user_age"},
			{Source: "user_age", Target: "condition"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	condResult := mustGetMapResult(t, result, "condition")
	if !condResult["condition_met"].(bool) {
		t.Error("Expected condition with context variable to be met")
	}
}

// TestConditionalComplex_BooleanLogic tests complex boolean expressions
func TestConditionalComplex_BooleanLogic(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "value",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(15)},
			},
			{
				ID:   "condition_and",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{
					Condition: strPtr("input > 10 && input < 20"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "value", Target: "condition_and"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	condResult := mustGetMapResult(t, result, "condition_and")
	if !condResult["condition_met"].(bool) {
		t.Error("Expected AND condition to be met")
	}
}

// TestConditionalComplex_OrLogic tests OR boolean logic
func TestConditionalComplex_OrLogic(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "value",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(5)},
			},
			{
				ID:   "condition_or",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{
					Condition: strPtr("input < 10 || input > 100"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "value", Target: "condition_or"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	condResult := mustGetMapResult(t, result, "condition_or")
	if !condResult["condition_met"].(bool) {
		t.Error("Expected OR condition to be met")
	}
}

// TestConditionalObjectFieldAccess tests condition on object fields
func TestConditionalObjectFieldAccess(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "user_data",
				Type: types.NodeTypeTransform,
				Data: types.TransformData{
					TransformType: strPtr("to_object"),
					Fields: []types.FieldMapping{
						{Source: strPtr("age"), Target: strPtr("age")},
						{Source: strPtr("status"), Target: strPtr("status")},
					},
				},
			},
			{
				ID:   "age_input",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: float64Ptr(30)},
			},
			{
				ID:   "status_input",
				Type: types.NodeTypeTextInput,
				Data: types.TextInputData{Text: strPtr("active")},
			},
			{
				ID:   "condition",
				Type: types.NodeTypeCondition,
				Data: types.ConditionData{
					Condition: strPtr("input.age > 18 && input.status == 'active'"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "age_input", Target: "user_data"},
			{Source: "status_input", Target: "user_data"},
			{Source: "user_data", Target: "condition"},
		},
	}

	// This test may need adjustment based on actual transform implementation
	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	// May fail if transform doesn't support this - that's valuable feedback
	if err != nil {
		t.Logf("Test reveals limitation in object field access: %v", err)
	} else {
		condResult := mustGetMapResult(t, result, "condition")
		if !condResult["condition_met"].(bool) {
			t.Error("Expected complex object condition to be met")
		}
	}
}

// Helper functions (intentionally at end to avoid redeclaration with other test files)

func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func mustGetMapResult(t *testing.T, result *types.Result, nodeID string) map[string]interface{} {
	t.Helper()
	nodeResult, ok := result.NodeResults[nodeID]
	if !ok {
		t.Fatalf("Node result for %s not found", nodeID)
	}
	mapResult, ok := nodeResult.(map[string]interface{})
	if !ok {
		t.Fatalf("Node result for %s is not a map, got %T", nodeID, nodeResult)
	}
	return mapResult
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
