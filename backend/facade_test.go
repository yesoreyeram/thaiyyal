package workflow_test

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend"
)

// TestFacadeBackwardCompatibility verifies that the facade maintains 100% backward compatibility
func TestFacadeBackwardCompatibility(t *testing.T) {
	// Test 1: Verify all type re-exports work
	t.Run("Type re-exports", func(t *testing.T) {
		var _ workflow.NodeType
		var _ workflow.Node
		var _ workflow.NodeData
		var _ workflow.Edge
		var _ workflow.Payload
		var _ workflow.Result
		var _ workflow.Config
		var _ workflow.SwitchCase
		var _ workflow.ContextVariableValue
		var _ workflow.CacheEntry
		var _ *workflow.Engine
	})

	// Test 2: Verify all constants are available
	t.Run("Constant re-exports", func(t *testing.T) {
		_ = workflow.ContextKeyExecutionID
		_ = workflow.ContextKeyWorkflowID

		// Basic I/O
		_ = workflow.NodeTypeNumber
		_ = workflow.NodeTypeTextInput
		_ = workflow.NodeTypeVisualization

		// Operations
		_ = workflow.NodeTypeOperation
		_ = workflow.NodeTypeTextOperation
		_ = workflow.NodeTypeHTTP

		// Control Flow
		_ = workflow.NodeTypeCondition
		_ = workflow.NodeTypeForEach
		_ = workflow.NodeTypeWhileLoop

		// State & Memory
		_ = workflow.NodeTypeVariable
		_ = workflow.NodeTypeExtract
		_ = workflow.NodeTypeTransform
		_ = workflow.NodeTypeAccumulator
		_ = workflow.NodeTypeCounter

		// Advanced Control
		_ = workflow.NodeTypeSwitch
		_ = workflow.NodeTypeParallel
		_ = workflow.NodeTypeJoin
		_ = workflow.NodeTypeSplit
		_ = workflow.NodeTypeDelay
		_ = workflow.NodeTypeCache

		// Error Handling
		_ = workflow.NodeTypeRetry
		_ = workflow.NodeTypeTryCatch
		_ = workflow.NodeTypeTimeout

		// Context
		_ = workflow.NodeTypeContextVariable
		_ = workflow.NodeTypeContextConstant
	})

	// Test 3: Verify all functions are available
	t.Run("Function re-exports", func(t *testing.T) {
		payload := `{"nodes": [], "edges": []}`

		// Engine constructors
		engine1, err := workflow.NewEngine([]byte(payload))
		if err != nil {
			t.Errorf("NewEngine failed: %v", err)
		}
		if engine1 == nil {
			t.Error("NewEngine returned nil")
		}

		config := workflow.DefaultConfig()
		engine2, err := workflow.NewEngineWithConfig([]byte(payload), config)
		if err != nil {
			t.Errorf("NewEngineWithConfig failed: %v", err)
		}
		if engine2 == nil {
			t.Error("NewEngineWithConfig returned nil")
		}

		// Configuration functions
		_ = workflow.DefaultConfig()
		_ = workflow.ValidationLimits()
		_ = workflow.DevelopmentConfig()
	})

	// Test 4: Verify Engine.Execute() works
	t.Run("Engine execution", func(t *testing.T) {
		payload := `{
			"nodes": [
				{"id": "1", "data": {"value": 10}},
				{"id": "2", "data": {"value": 5}},
				{"id": "3", "data": {"op": "add"}}
			],
			"edges": [
				{"source": "1", "target": "3"},
				{"source": "2", "target": "3"}
			]
		}`

		engine, err := workflow.NewEngine([]byte(payload))
		if err != nil {
			t.Fatalf("NewEngine failed: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if result == nil {
			t.Fatal("Execute returned nil result")
		}

		if result.FinalOutput == nil {
			t.Fatal("Execute returned nil FinalOutput")
		}

		// Verify the calculation worked
		output, ok := result.FinalOutput.(float64)
		if !ok {
			t.Fatalf("Expected float64 output, got %T", result.FinalOutput)
		}

		if output != 15 {
			t.Errorf("Expected output 15, got %v", output)
		}
	})
}
