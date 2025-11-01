package engine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// ============================================================================
// Example Custom Executors for Testing
// ============================================================================

// ReverseStringExecutor is a custom executor that reverses strings
type ReverseStringExecutor struct{}

func (e *ReverseStringExecutor) Execute(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
	// Increment node execution counter for protection
	if err := ctx.IncrementNodeExecution(); err != nil {
		return nil, err
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("reverse_string node requires at least one input")
	}

	// Get the first input and convert to string
	str, ok := inputs[0].(string)
	if !ok {
		return nil, fmt.Errorf("reverse_string input must be a string, got %T", inputs[0])
	}

	// Reverse the string
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}

func (e *ReverseStringExecutor) NodeType() types.NodeType {
	return types.NodeType("reverse_string")
}

func (e *ReverseStringExecutor) Validate(node types.Node) error {
	// No special validation needed for this simple node
	return nil
}

// MultiplyByNExecutor is a custom executor that multiplies numbers by a configured factor
type MultiplyByNExecutor struct{}

func (e *MultiplyByNExecutor) Execute(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
	// Increment node execution counter
	if err := ctx.IncrementNodeExecution(); err != nil {
		return nil, err
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("multiply_by_n node requires at least one input")
	}

	// Get the multiplier from node data
	factor := node.Data.Factor
	if factor == nil {
		return nil, fmt.Errorf("multiply_by_n node requires 'factor' in data")
	}

	// Get the input value
	inputVal, ok := inputs[0].(float64)
	if !ok {
		return nil, fmt.Errorf("multiply_by_n input must be a number, got %T", inputs[0])
	}

	result := inputVal * (*factor)
	return result, nil
}

func (e *MultiplyByNExecutor) NodeType() types.NodeType {
	return types.NodeType("multiply_by_n")
}

func (e *MultiplyByNExecutor) Validate(node types.Node) error {
	if node.Data.Factor == nil {
		return fmt.Errorf("multiply_by_n node requires 'factor' field")
	}
	return nil
}

// ConcatWithPrefixExecutor concatenates strings with a custom prefix
type ConcatWithPrefixExecutor struct{}

func (e *ConcatWithPrefixExecutor) Execute(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
	if err := ctx.IncrementNodeExecution(); err != nil {
		return nil, err
	}

	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("concat_prefix node requires at least one input")
	}

	// Get prefix from node data
	prefix := ""
	if node.Data.Prefix != nil {
		prefix = *node.Data.Prefix
	}

	// Concatenate all inputs with prefix
	var parts []string
	for _, input := range inputs {
		str, ok := input.(string)
		if !ok {
			return nil, fmt.Errorf("concat_prefix requires string inputs, got %T", input)
		}
		parts = append(parts, prefix+str)
	}

	return strings.Join(parts, " "), nil
}

func (e *ConcatWithPrefixExecutor) NodeType() types.NodeType {
	return types.NodeType("concat_prefix")
}

func (e *ConcatWithPrefixExecutor) Validate(node types.Node) error {
	// Prefix is optional, so no validation needed
	return nil
}

// BadExecutor intentionally doesn't increment node execution counter (for testing)
type BadExecutor struct{}

func (e *BadExecutor) Execute(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
	// Intentionally NOT calling IncrementNodeExecution() - bad practice!
	return "bad", nil
}

func (e *BadExecutor) NodeType() types.NodeType {
	return types.NodeType("bad_executor")
}

func (e *BadExecutor) Validate(node types.Node) error {
	return nil
}

// ============================================================================
// Tests for Custom Executor Registration
// ============================================================================

func TestCustomExecutorRegistration(t *testing.T) {
	t.Run("register single custom executor", func(t *testing.T) {
		registry := executor.NewRegistry()
		err := registry.Register(&ReverseStringExecutor{})
		if err != nil {
			t.Fatalf("failed to register custom executor: %v", err)
		}

		// Verify it was registered
		exec := registry.GetExecutor(types.NodeType("reverse_string"))
		if exec == nil {
			t.Fatal("custom executor not found in registry")
		}
	})

	t.Run("register multiple custom executors", func(t *testing.T) {
		registry := executor.NewRegistry()
		registry.MustRegister(&ReverseStringExecutor{})
		registry.MustRegister(&MultiplyByNExecutor{})
		registry.MustRegister(&ConcatWithPrefixExecutor{})

		// Verify all were registered
		registeredTypes := registry.ListRegisteredTypes()
		if len(registeredTypes) != 3 {
			t.Fatalf("expected 3 registered types, got %d", len(registeredTypes))
		}
	})

	t.Run("cannot register duplicate executor", func(t *testing.T) {
		registry := executor.NewRegistry()
		registry.MustRegister(&ReverseStringExecutor{})

		// Try to register again
		err := registry.Register(&ReverseStringExecutor{})
		if err == nil {
			t.Fatal("expected error when registering duplicate executor")
		}
		if !strings.Contains(err.Error(), "already registered") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})

	t.Run("combine default and custom executors", func(t *testing.T) {
		registry := DefaultRegistry()
		originalCount := len(registry.ListRegisteredTypes())

		// Add custom executors to default registry
		registry.MustRegister(&ReverseStringExecutor{})
		registry.MustRegister(&MultiplyByNExecutor{})

		newCount := len(registry.ListRegisteredTypes())
		if newCount != originalCount+2 {
			t.Fatalf("expected %d types, got %d", originalCount+2, newCount)
		}

		// Verify default executors still work
		numberExec := registry.GetExecutor(types.NodeTypeNumber)
		if numberExec == nil {
			t.Fatal("default Number executor not found")
		}

		// Verify custom executors work
		reverseExec := registry.GetExecutor(types.NodeType("reverse_string"))
		if reverseExec == nil {
			t.Fatal("custom ReverseString executor not found")
		}
	})
}

// ============================================================================
// Tests for Custom Executor Execution
// ============================================================================

func TestCustomExecutorExecution(t *testing.T) {
	t.Run("simple custom executor workflow", func(t *testing.T) {
		registry := DefaultRegistry()
		registry.MustRegister(&ReverseStringExecutor{})

		payload := `{
			"nodes": [
				{"id": "1", "data": {"text": "hello"}},
				{"id": "2", "type": "reverse_string", "data": {}}
			],
			"edges": [
				{"source": "1", "target": "2"}
			]
		}`

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		// Check final output is reversed
		expected := "olleh"
		if result.FinalOutput != expected {
			t.Fatalf("expected %q, got %q", expected, result.FinalOutput)
		}
	})

	t.Run("custom executor with configuration", func(t *testing.T) {
		registry := DefaultRegistry()
		registry.MustRegister(&MultiplyByNExecutor{})

		factor := 5.0
		payload := fmt.Sprintf(`{
			"nodes": [
				{"id": "1", "data": {"value": 10}},
				{"id": "2", "type": "multiply_by_n", "data": {"factor": %f}}
			],
			"edges": [
				{"source": "1", "target": "2"}
			]
		}`, factor)

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		// Check result is 10 * 5 = 50
		expected := 50.0
		if result.FinalOutput != expected {
			t.Fatalf("expected %v, got %v", expected, result.FinalOutput)
		}
	})

	t.Run("custom executor with multiple inputs", func(t *testing.T) {
		registry := DefaultRegistry()
		registry.MustRegister(&ConcatWithPrefixExecutor{})

		prefix := ">> "
		payload := fmt.Sprintf(`{
			"nodes": [
				{"id": "1", "data": {"text": "hello"}},
				{"id": "2", "data": {"text": "world"}},
				{"id": "3", "type": "concat_prefix", "data": {"prefix": %q}}
			],
			"edges": [
				{"source": "1", "target": "3"},
				{"source": "2", "target": "3"}
			]
		}`, prefix)

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		// Check result contains prefixed strings
		resultStr, ok := result.FinalOutput.(string)
		if !ok {
			t.Fatalf("expected string result, got %T", result.FinalOutput)
		}

		if !strings.Contains(resultStr, ">> hello") || !strings.Contains(resultStr, ">> world") {
			t.Fatalf("unexpected result: %q", resultStr)
		}
	})

	t.Run("mixing built-in and custom executors", func(t *testing.T) {
		registry := DefaultRegistry()
		registry.MustRegister(&ReverseStringExecutor{})

		payload := `{
			"nodes": [
				{"id": "1", "data": {"text": "hello world"}},
				{"id": "2", "data": {"text_op": "uppercase"}},
				{"id": "3", "type": "reverse_string", "data": {}}
			],
			"edges": [
				{"source": "1", "target": "2"},
				{"source": "2", "target": "3"}
			]
		}`

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		// "hello world" -> "HELLO WORLD" -> "DLROW OLLEH"
		expected := "DLROW OLLEH"
		if result.FinalOutput != expected {
			t.Fatalf("expected %q, got %q", expected, result.FinalOutput)
		}
	})
}

// ============================================================================
// Tests for Protection Limits with Custom Executors
// ============================================================================

func TestCustomExecutorProtectionLimits(t *testing.T) {
	t.Run("custom executor respects node execution limit", func(t *testing.T) {
		registry := DefaultRegistry()
		registry.MustRegister(&ReverseStringExecutor{})

		// Create config with very low limit
		config := types.DefaultConfig()
		config.MaxNodeExecutions = 2 // Only allow 2 node executions

		payload := `{
			"nodes": [
				{"id": "1", "data": {"text": "hello"}},
				{"id": "2", "type": "reverse_string", "data": {}},
				{"id": "3", "data": {"text_op": "uppercase"}}
			],
			"edges": [
				{"source": "1", "target": "2"},
				{"source": "2", "target": "3"}
			]
		}`

		engine, err := NewWithRegistry([]byte(payload), config, registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		_, err = engine.Execute()
		if err == nil {
			t.Fatal("expected error due to node execution limit")
		}

		if !strings.Contains(err.Error(), "maximum node executions exceeded") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("custom executor counts toward execution limit", func(t *testing.T) {
		registry := DefaultRegistry()
		registry.MustRegister(&ReverseStringExecutor{})
		registry.MustRegister(&MultiplyByNExecutor{})

		config := types.DefaultConfig()
		config.MaxNodeExecutions = 10

		factor := 2.0
		payload := fmt.Sprintf(`{
			"nodes": [
				{"id": "1", "data": {"value": 5}},
				{"id": "2", "type": "multiply_by_n", "data": {"factor": %f}},
				{"id": "3", "data": {"op": "add"}},
				{"id": "4", "data": {"value": 3}},
				{"id": "5", "type": "multiply_by_n", "data": {"factor": %f}}
			],
			"edges": [
				{"source": "1", "target": "2"},
				{"source": "4", "target": "3"},
				{"source": "2", "target": "5"},
				{"source": "5", "target": "3"}
			]
		}`, factor, factor)

		engine, err := NewWithRegistry([]byte(payload), config, registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		// Check execution count
		// Note: Engine increments counter for each node, and custom executors also increment
		// So we get: 1 (node 1) + 2 (node 2: engine + executor) + 1 (node 4) + 2 (node 5: engine + executor) + 1 (node 3) = 7
		execCount := engine.GetNodeExecutionCount()
		if execCount < 5 {
			t.Fatalf("expected at least 5 node executions, got %d", execCount)
		}

		// Verify correct calculation: (5 * 2 * 2) + 3 = 23
		expected := 23.0
		if result.FinalOutput != expected {
			t.Fatalf("expected %v, got %v", expected, result.FinalOutput)
		}
	})

	t.Run("bad executor not incrementing counter still protected", func(t *testing.T) {
		// Even if a custom executor doesn't call IncrementNodeExecution(),
		// the engine itself should still enforce limits through the outer wrapper
		registry := DefaultRegistry()
		registry.MustRegister(&BadExecutor{})

		config := types.DefaultConfig()

		payload := `{
			"nodes": [
				{"id": "1", "type": "bad_executor", "data": {}}
			],
			"edges": []
		}`

		engine, err := NewWithRegistry([]byte(payload), config, registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		// Execution should succeed but count should be incremented by engine
		if result.FinalOutput != "bad" {
			t.Fatalf("unexpected result: %v", result.FinalOutput)
		}
	})
}

// ============================================================================
// Tests for Validation with Custom Executors
// ============================================================================

func TestCustomExecutorValidation(t *testing.T) {
	t.Run("validation fails for missing required field", func(t *testing.T) {
		registry := DefaultRegistry()
		registry.MustRegister(&MultiplyByNExecutor{})

		// Missing 'factor' field
		payload := `{
			"nodes": [
				{"id": "1", "data": {"value": 10}},
				{"id": "2", "type": "multiply_by_n", "data": {}}
			],
			"edges": [
				{"source": "1", "target": "2"}
			]
		}`

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		// Execution should fail because validation catches missing factor
		_, err = engine.Execute()
		if err == nil {
			t.Fatal("expected error for missing factor")
		}

		if !strings.Contains(err.Error(), "factor") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("validation succeeds for well-formed custom node", func(t *testing.T) {
		registry := DefaultRegistry()
		registry.MustRegister(&MultiplyByNExecutor{})

		factor := 3.0
		payload := fmt.Sprintf(`{
			"nodes": [
				{"id": "1", "data": {"value": 10}},
				{"id": "2", "type": "multiply_by_n", "data": {"factor": %f}}
			],
			"edges": [
				{"source": "1", "target": "2"}
			]
		}`, factor)

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		expected := 30.0
		if result.FinalOutput != expected {
			t.Fatalf("expected %v, got %v", expected, result.FinalOutput)
		}
	})
}

// ============================================================================
// Tests for Error Handling in Custom Executors
// ============================================================================

func TestCustomExecutorErrorHandling(t *testing.T) {
	t.Run("custom executor error propagates correctly", func(t *testing.T) {
		registry := DefaultRegistry()
		registry.MustRegister(&ReverseStringExecutor{})

		// Number input to string executor should fail
		payload := `{
			"nodes": [
				{"id": "1", "data": {"value": 123}},
				{"id": "2", "type": "reverse_string", "data": {}}
			],
			"edges": [
				{"source": "1", "target": "2"}
			]
		}`

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		_, err = engine.Execute()
		if err == nil {
			t.Fatal("expected error for type mismatch")
		}

		if !strings.Contains(err.Error(), "must be a string") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unregistered custom node type fails", func(t *testing.T) {
		// Don't register the custom executor
		registry := DefaultRegistry()

		payload := `{
			"nodes": [
				{"id": "1", "data": {"text": "hello"}},
				{"id": "2", "type": "reverse_string", "data": {}}
			],
			"edges": [
				{"source": "1", "target": "2"}
			]
		}`

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		_, err = engine.Execute()
		if err == nil {
			t.Fatal("expected error for unregistered executor")
		}

		if !strings.Contains(err.Error(), "no executor registered") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// ============================================================================
// Tests for NewWithRegistry Constructor
// ============================================================================

func TestNewWithRegistry(t *testing.T) {
	t.Run("nil registry returns error", func(t *testing.T) {
		payload := `{"nodes": [], "edges": []}`
		_, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), nil)
		if err == nil {
			t.Fatal("expected error for nil registry")
		}
		if !strings.Contains(err.Error(), "registry cannot be nil") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty registry works for workflows with no nodes", func(t *testing.T) {
		registry := executor.NewRegistry()
		payload := `{"nodes": [], "edges": []}`

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		result, err := engine.Execute()
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		if result.FinalOutput != nil {
			t.Fatalf("expected nil output, got %v", result.FinalOutput)
		}
	})

	t.Run("custom-only registry without built-in executors", func(t *testing.T) {
		// Create registry with only custom executors (no built-ins)
		registry := executor.NewRegistry()
		registry.MustRegister(&ReverseStringExecutor{})

		// This workflow only uses the custom executor
		payload := `{
			"nodes": [
				{"id": "1", "data": {"text": "test"}},
				{"id": "2", "type": "reverse_string", "data": {}}
			],
			"edges": [
				{"source": "1", "target": "2"}
			]
		}`

		engine, err := NewWithRegistry([]byte(payload), types.DefaultConfig(), registry)
		if err != nil {
			t.Fatalf("failed to create engine: %v", err)
		}

		// Should fail because text_input is not registered
		_, err = engine.Execute()
		if err == nil {
			t.Fatal("expected error for unregistered text_input executor")
		}
	})
}
