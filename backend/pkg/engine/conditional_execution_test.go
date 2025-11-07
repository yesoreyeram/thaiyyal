package engine

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestConditionalExecution_TruePathOnly tests execution only on true path
func TestConditionalExecution_TruePathOnly(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{ID: "age", Type: types.NodeTypeNumber, Data: types.NumberData{Value: float64Ptr(25)}},
			{ID: "check", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr(">=18")}},
			{ID: "adult_action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Adult profile")}},
			{ID: "minor_action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Minor education")}},
		},
		Edges: []types.Edge{
			{Source: "age", Target: "check"},
			{Source: "check", Target: "adult_action", SourceHandle: strPtr("true")},
			{Source: "check", Target: "minor_action", SourceHandle: strPtr("false")},
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

	// Adult action should execute (true path)
	if _, ok := result.NodeResults["adult_action"]; !ok {
		t.Error("Expected adult_action to execute on true path")
	}

	// Minor action should NOT execute (false path)
	if _, ok := result.NodeResults["minor_action"]; ok {
		t.Error("Expected minor_action to NOT execute on false path")
	}
}

// TestConditionalExecution_FalsePathOnly tests execution only on false path
func TestConditionalExecution_FalsePathOnly(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{ID: "age", Type: types.NodeTypeNumber, Data: types.NumberData{Value: float64Ptr(15)}},
			{ID: "check", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr(">=18")}},
			{ID: "adult_action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Adult profile")}},
			{ID: "minor_action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Minor education")}},
		},
		Edges: []types.Edge{
			{Source: "age", Target: "check"},
			{Source: "check", Target: "adult_action", SourceHandle: strPtr("true")},
			{Source: "check", Target: "minor_action", SourceHandle: strPtr("false")},
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

	// Minor action should execute (false path)
	if _, ok := result.NodeResults["minor_action"]; !ok {
		t.Error("Expected minor_action to execute on false path")
	}

	// Adult action should NOT execute (true path)
	if _, ok := result.NodeResults["adult_action"]; ok {
		t.Error("Expected adult_action to NOT execute on true path")
	}
}

// TestConditionalExecution_SwitchRouting tests switch-based routing
func TestConditionalExecution_SwitchRouting(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{ID: "status_code", Type: types.NodeTypeNumber, Data: types.NumberData{Value: float64Ptr(200)}},
			{ID: "router", Type: types.NodeTypeSwitch, Data: types.SwitchData{
				Cases: []types.SwitchCase{
					{When: "input == 200", OutputPath: strPtr("success")},
					{When: "input == 404", OutputPath: strPtr("not_found")},
					{When: "input >= 500", OutputPath: strPtr("error")},
					{When: "default", OutputPath: strPtr("other"), IsDefault: true},
				},
			}},
			{ID: "success_handler", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Success")}},
			{ID: "error_handler", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Error")}},
			{ID: "not_found_handler", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Not Found")}},
		},
		Edges: []types.Edge{
			{Source: "status_code", Target: "router"},
			{Source: "router", Target: "success_handler", SourceHandle: strPtr("success")},
			{Source: "router", Target: "error_handler", SourceHandle: strPtr("error")},
			{Source: "router", Target: "not_found_handler", SourceHandle: strPtr("not_found")},
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

	// Success handler should execute
	if _, ok := result.NodeResults["success_handler"]; !ok {
		t.Error("Expected success_handler to execute")
	}

	// Other handlers should NOT execute
	if _, ok := result.NodeResults["error_handler"]; ok {
		t.Error("Expected error_handler to NOT execute")
	}
	if _, ok := result.NodeResults["not_found_handler"]; ok {
		t.Error("Expected not_found_handler to NOT execute")
	}
}

// TestConditionalExecution_NestedConditions tests nested conditional branches
func TestConditionalExecution_NestedConditions(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{ID: "age", Type: types.NodeTypeNumber, Data: types.NumberData{Value: float64Ptr(30)}},
			{ID: "age_check", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr(">=21")}},
			{ID: "extract1", Type: types.NodeTypeExtract, Data: types.ExtractData{Field: strPtr("value")}},
			{ID: "senior_check", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr(">=65")}},
			{ID: "adult_action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Adult")}},
			{ID: "senior_action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Senior")}},
			{ID: "minor_action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Minor")}},
		},
		Edges: []types.Edge{
			{Source: "age", Target: "age_check"},
			{Source: "age_check", Target: "extract1", SourceHandle: strPtr("true")},
			{Source: "age_check", Target: "minor_action", SourceHandle: strPtr("false")},
			{Source: "extract1", Target: "senior_check"},
			{Source: "senior_check", Target: "adult_action", SourceHandle: strPtr("false")},
			{Source: "senior_check", Target: "senior_action", SourceHandle: strPtr("true")},
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

	// Adult action should execute (age 30: >=21 but <65)
	if _, ok := result.NodeResults["adult_action"]; !ok {
		t.Error("Expected adult_action to execute")
	}

	// Senior and minor should NOT execute
	if _, ok := result.NodeResults["senior_action"]; ok {
		t.Error("Expected senior_action to NOT execute")
	}
	if _, ok := result.NodeResults["minor_action"]; ok {
		t.Error("Expected minor_action to NOT execute")
	}
}

// TestConditionalExecution_MultipleConditionalEdges tests node with multiple conditional incoming edges
func TestConditionalExecution_MultipleConditionalEdges(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{ID: "flag1", Type: types.NodeTypeBooleanInput, Data: types.BooleanInputData{BooleanValue: boolPtr(true)}},
			{ID: "flag2", Type: types.NodeTypeBooleanInput, Data: types.BooleanInputData{BooleanValue: boolPtr(false)}},
			{ID: "check1", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr("input == true")}},
			{ID: "check2", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr("input == true")}},
			{ID: "action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Executed")}},
		},
		Edges: []types.Edge{
			{Source: "flag1", Target: "check1"},
			{Source: "flag2", Target: "check2"},
			// Action executes if EITHER check1 true OR check2 true
			{Source: "check1", Target: "action", SourceHandle: strPtr("true")},
			{Source: "check2", Target: "action", SourceHandle: strPtr("true")},
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

	// Action should execute because check1 is true (even though check2 is false)
	if _, ok := result.NodeResults["action"]; !ok {
		t.Error("Expected action to execute when at least one condition is true")
	}
}

// TestConditionalExecution_UnconditionalEdgeTakesPrecedence tests unconditional edge behavior
func TestConditionalExecution_UnconditionalEdgeTakesPrecedence(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{ID: "value", Type: types.NodeTypeNumber, Data: types.NumberData{Value: float64Ptr(10)}},
			{ID: "check", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr(">100")}},
			{ID: "action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Always runs")}},
		},
		Edges: []types.Edge{
			{Source: "value", Target: "check"},
			// Mix of conditional and unconditional edges to same target
			{Source: "check", Target: "action", SourceHandle: strPtr("true")},
			{Source: "value", Target: "action"}, // Unconditional edge
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

	// Action should execute due to unconditional edge (even though conditional is false)
	if _, ok := result.NodeResults["action"]; !ok {
		t.Error("Expected action to execute via unconditional edge")
	}
}

// TestConditionalExecution_BackwardCompatibility tests legacy "condition" field
func TestConditionalExecution_BackwardCompatibility(t *testing.T) {
	payload := types.Payload{
		Nodes: []types.Node{
			{ID: "age", Type: types.NodeTypeNumber, Data: types.NumberData{Value: float64Ptr(25)}},
			{ID: "check", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr(">=18")}},
			{ID: "action", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("Adult")}},
		},
		Edges: []types.Edge{
			{Source: "age", Target: "check"},
			{Source: "check", Target: "action", Condition: strPtr("true")}, // Legacy field
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

	// Action should execute using legacy condition field
	if _, ok := result.NodeResults["action"]; !ok {
		t.Error("Expected backward compatibility with legacy condition field")
	}
}
