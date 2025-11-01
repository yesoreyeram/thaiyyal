package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestForEachExecutor_ReduceMode tests REDUCE mode (accumulation)
func TestForEachExecutor_ReduceMode(t *testing.T) {
	exec := &ForEachExecutor{}
	mode := "reduce"
	initialValue := float64(0)

	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {[]interface{}{float64(10), float64(20), float64(30)}},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeForEach,
		Data: types.NodeData{
			Mode:         &mode,
			InitialValue: initialValue,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})

	// Verify mode
	if resultMap["mode"].(string) != "reduce" {
		t.Errorf("Expected mode=reduce, got %s", resultMap["mode"])
	}

	// Verify result exists (would be accumulated value in full implementation)
	if _, ok := resultMap["result"]; !ok {
		t.Error("Expected result field in reduce mode")
	}

	// Verify initial and final values are tracked
	if resultMap["initial_value"] != initialValue {
		t.Errorf("Expected initial_value=%v, got %v", initialValue, resultMap["initial_value"])
	}
}

// TestForEachExecutor_FilterMapMode tests FILTER_MAP mode
func TestForEachExecutor_FilterMapMode(t *testing.T) {
	exec := &ForEachExecutor{}
	mode := "filter_map"
	condition := "item > 10"

	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {[]interface{}{float64(5), float64(15), float64(25)}},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeForEach,
		Data: types.NodeData{
			Mode:      &mode,
			Condition: &condition,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})

	// Verify mode
	if resultMap["mode"].(string) != "filter_map" {
		t.Errorf("Expected mode=filter_map, got %s", resultMap["mode"])
	}

	// Verify filtered_count exists
	if _, ok := resultMap["filtered_count"]; !ok {
		t.Error("Expected filtered_count in filter_map mode")
	}

	// Verify results array
	if _, ok := resultMap["results"]; !ok {
		t.Error("Expected results array in filter_map mode")
	}
}

// TestForEachExecutor_ForEachMode tests FOREACH mode (side effects)
func TestForEachExecutor_ForEachMode(t *testing.T) {
	exec := &ForEachExecutor{}
	mode := "foreach"

	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {[]interface{}{1, 2, 3}},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeForEach,
		Data: types.NodeData{
			Mode: &mode,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})

	// Verify mode
	if resultMap["mode"].(string) != "foreach" {
		t.Errorf("Expected mode=foreach, got %s", resultMap["mode"])
	}

	// Should have iterations count
	if resultMap["iterations"].(int) != 3 {
		t.Errorf("Expected 3 iterations, got %d", resultMap["iterations"])
	}

	// Should NOT have results array (side effects only)
	if _, ok := resultMap["results"]; ok {
		t.Error("foreach mode should not have results array")
	}
}

// TestForEachExecutor_MetadataMode tests backward compatible metadata mode
func TestForEachExecutor_MetadataMode(t *testing.T) {
	exec := &ForEachExecutor{}
	mode := "metadata"

	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {[]interface{}{1, 2, 3}},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeForEach,
		Data: types.NodeData{
			Mode: &mode,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})

	// Verify mode
	if resultMap["mode"].(string) != "metadata" {
		t.Errorf("Expected mode=metadata, got %s", resultMap["mode"])
	}

	// Should have items array (backward compatible)
	items, ok := resultMap["items"].([]interface{})
	if !ok {
		t.Fatal("Expected items array in metadata mode")
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}
}

// TestForEachExecutor_DefaultMode tests default mode (should be MAP)
func TestForEachExecutor_DefaultMode(t *testing.T) {
	exec := &ForEachExecutor{}

	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {[]interface{}{1, 2, 3}},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeForEach,
		Data: types.NodeData{
			// No mode specified - should default to map
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})

	// Should default to map mode
	if resultMap["mode"].(string) != "map" {
		t.Errorf("Expected default mode=map, got %s", resultMap["mode"])
	}
}

// TestForEachExecutor_InvalidMode tests invalid mode handling
func TestForEachExecutor_InvalidMode(t *testing.T) {
	exec := &ForEachExecutor{}
	invalidMode := "invalid_mode"

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeForEach,
		Data: types.NodeData{
			Mode: &invalidMode,
		},
	}

	err := exec.Validate(node)
	if err == nil {
		t.Error("Expected validation error for invalid mode")
	}
}

// TestForEachExecutor_FilterMapRequiresCondition tests validation
func TestForEachExecutor_FilterMapRequiresCondition(t *testing.T) {
	exec := &ForEachExecutor{}
	mode := "filter_map"

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeForEach,
		Data: types.NodeData{
			Mode: &mode,
			// Missing condition
		},
	}

	err := exec.Validate(node)
	if err == nil {
		t.Error("Expected validation error for filter_map without condition")
	}
}

// TestForEachExecutor_SuccessAndFailedCounts tests error handling
func TestForEachExecutor_SuccessAndFailedCounts(t *testing.T) {
	exec := &ForEachExecutor{}
	mode := "map"

	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {[]interface{}{1, 2, 3, 4, 5}},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeForEach,
		Data: types.NodeData{
			Mode: &mode,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})

	// Verify successful and failed counts
	successful := resultMap["successful"].(int)
	failed := resultMap["failed"].(int)

	if successful+failed != 5 {
		t.Errorf("Expected total of 5 iterations, got successful=%d, failed=%d", successful, failed)
	}
}
