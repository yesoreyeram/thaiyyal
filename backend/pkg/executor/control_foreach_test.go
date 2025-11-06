package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestForEachExecutor_BasicIteration(t *testing.T) {
	executor := &ForEachExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"foreach1": {[]interface{}{"item1", "item2", "item3"}},
		},
	}

	node := types.Node{
		ID:   "foreach1",
		Type: types.NodeTypeForEach,
		Data: types.ForEachData{},
	}

	result, err := executor.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	if resultMap["iterations"] != 3 {
		t.Errorf("Expected iterations=3, got %v", resultMap["iterations"])
	}
}

func TestForEachExecutor_NonArrayInput(t *testing.T) {
	executor := &ForEachExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"foreach1": {"not an array"},
		},
	}

	node := types.Node{
		ID: "foreach1",
	}

	_, err := executor.Execute(ctx, node)
	if err == nil {
		t.Fatal("Expected error for non-array input")
	}
}
