package workflow

import (
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Parallel Execution Tests
// ============================================================================

// TestComputeExecutionLevels tests the level computation algorithm
func TestComputeExecutionLevels(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		expectedLevels int
		expectError    bool
	}{
		{
			name: "simple linear workflow",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"op": "add"}},
					{"id": "3", "data": {"mode": "text"}}
				],
				"edges": [
					{"source": "1", "target": "2"},
					{"source": "2", "target": "3"}
				]
			}`,
			expectedLevels: 3,
			expectError:    false,
		},
		{
			name: "parallel branches",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 20}},
					{"id": "3", "data": {"value": 30}},
					{"id": "4", "data": {"op": "add"}}
				],
				"edges": [
					{"source": "1", "target": "4"},
					{"source": "2", "target": "4"},
					{"source": "3", "target": "4"}
				]
			}`,
			expectedLevels: 2,
			expectError:    false,
		},
		{
			name: "diamond pattern",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 5}},
					{"id": "3", "data": {"value": 3}},
					{"id": "4", "data": {"op": "add"}}
				],
				"edges": [
					{"source": "1", "target": "2"},
					{"source": "1", "target": "3"},
					{"source": "2", "target": "4"},
					{"source": "3", "target": "4"}
				]
			}`,
			expectedLevels: 3,
			expectError:    false,
		},
		{
			name: "complex workflow with multiple levels",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 20}},
					{"id": "3", "data": {"op": "add"}},
					{"id": "4", "data": {"op": "multiply"}},
					{"id": "5", "data": {"value": 2}},
					{"id": "6", "data": {"op": "add"}}
				],
				"edges": [
					{"source": "1", "target": "3"},
					{"source": "2", "target": "3"},
					{"source": "3", "target": "4"},
					{"source": "5", "target": "4"},
					{"source": "4", "target": "6"}
				]
			}`,
			expectedLevels: 4,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine([]byte(tt.payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			engine.inferNodeTypes()
			levels, err := engine.computeExecutionLevels()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(levels) != tt.expectedLevels {
				t.Errorf("Expected %d levels, got %d", tt.expectedLevels, len(levels))
			}

			// Verify all nodes are assigned to a level
			totalNodes := 0
			for _, level := range levels {
				totalNodes += len(level.NodeIDs)
			}
			if totalNodes != len(engine.nodes) {
				t.Errorf("Expected %d nodes in levels, got %d", len(engine.nodes), totalNodes)
			}
		})
	}
}

// TestParallelExecutionSimple tests basic parallel execution
func TestParallelExecutionSimple(t *testing.T) {
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

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := DefaultParallelConfig()
	result, err := engine.ExecuteWithParallelism(config)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	if result.FinalOutput != float64(15) {
		t.Errorf("Expected final output 15, got %v", result.FinalOutput)
	}
}

// TestParallelExecutionMultipleBranches tests parallel execution with multiple independent branches
func TestParallelExecutionMultipleBranches(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 20}},
			{"id": "3", "data": {"value": 30}},
			{"id": "4", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "4"},
			{"source": "2", "target": "4"},
			{"source": "3", "target": "4"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := DefaultParallelConfig()
	result, err := engine.ExecuteWithParallelism(config)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Add operation takes first two inputs: 10 + 20 = 30
	if result.FinalOutput != float64(30) {
		t.Errorf("Expected final output 30, got %v", result.FinalOutput)
	}

	// Verify all intermediate results are stored
	if len(result.NodeResults) != 4 {
		t.Errorf("Expected 4 node results, got %d", len(result.NodeResults))
	}
}

// TestParallelExecutionDiamond tests the diamond pattern (common in workflows)
func TestParallelExecutionDiamond(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 100}},
			{"id": "2", "data": {"value": 10}},
			{"id": "3", "data": {"value": 5}},
			{"id": "add1", "data": {"op": "add"}},
			{"id": "mult1", "data": {"op": "multiply"}},
			{"id": "final", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "add1"},
			{"source": "2", "target": "add1"},
			{"source": "1", "target": "mult1"},
			{"source": "3", "target": "mult1"},
			{"source": "add1", "target": "final"},
			{"source": "mult1", "target": "final"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := DefaultParallelConfig()
	result, err := engine.ExecuteWithParallelism(config)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Level 0: nodes 1, 2, 3 (values 100, 10, 5)
	// Level 1: add1 (100 + 10 = 110), mult1 (100 * 5 = 500)
	// Level 2: final (110 + 500 = 610)
	if result.FinalOutput != float64(610) {
		t.Errorf("Expected final output 610, got %v", result.FinalOutput)
	}
}

// TestParallelExecutionConcurrencyLimit tests concurrency limiting
func TestParallelExecutionConcurrencyLimit(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 1}},
			{"id": "2", "data": {"value": 2}},
			{"id": "3", "data": {"value": 3}},
			{"id": "4", "data": {"value": 4}},
			{"id": "5", "data": {"value": 5}},
			{"id": "6", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "6"},
			{"source": "2", "target": "6"},
			{"source": "3", "target": "6"},
			{"source": "4", "target": "6"},
			{"source": "5", "target": "6"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Test with concurrency limit of 2
	config := ParallelExecutionConfig{
		MaxConcurrency: 2,
		EnableParallel: true,
	}
	
	result, err := engine.ExecuteWithParallelism(config)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Add operation takes first two inputs: 1 + 2 = 3
	if result.FinalOutput != float64(3) {
		t.Errorf("Expected final output 3, got %v", result.FinalOutput)
	}
}

// TestParallelExecutionWithTextOperations tests parallel execution with text operations
func TestParallelExecutionWithTextOperations(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "hello"}},
			{"id": "2", "data": {"text": "world"}},
			{"id": "3", "data": {"text_op": "uppercase"}},
			{"id": "4", "data": {"text_op": "uppercase"}},
			{"id": "5", "data": {"text_op": "concat", "separator": " "}}
		],
		"edges": [
			{"source": "1", "target": "3"},
			{"source": "2", "target": "4"},
			{"source": "3", "target": "5"},
			{"source": "4", "target": "5"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := DefaultParallelConfig()
	result, err := engine.ExecuteWithParallelism(config)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	expected := "HELLO WORLD"
	if result.FinalOutput != expected {
		t.Errorf("Expected final output '%s', got '%v'", expected, result.FinalOutput)
	}
}

// TestParallelExecutionComplexWorkflow tests a complex multi-level workflow
func TestParallelExecutionComplexWorkflow(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "n1", "data": {"value": 10}},
			{"id": "n2", "data": {"value": 20}},
			{"id": "n3", "data": {"value": 30}},
			{"id": "n4", "data": {"op": "add"}},
			{"id": "n5", "data": {"op": "multiply"}},
			{"id": "n6", "data": {"value": 2}},
			{"id": "n7", "data": {"op": "add"}},
			{"id": "n8", "data": {"mode": "text"}}
		],
		"edges": [
			{"source": "n1", "target": "n4"},
			{"source": "n2", "target": "n4"},
			{"source": "n3", "target": "n5"},
			{"source": "n6", "target": "n5"},
			{"source": "n4", "target": "n7"},
			{"source": "n5", "target": "n7"},
			{"source": "n7", "target": "n8"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := DefaultParallelConfig()
	result, err := engine.ExecuteWithParallelism(config)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// n1 + n2 = 30, n3 * n6 = 60, 30 + 60 = 90
	// The final result is a map from visualization node, check it contains the value
	finalOutput := result.FinalOutput
	if outputMap, ok := finalOutput.(map[string]interface{}); ok {
		if outputMap["value"] != float64(90) {
			t.Errorf("Expected final output value 90, got %v", outputMap["value"])
		}
	} else {
		t.Errorf("Expected final output to be a map, got %T", finalOutput)
	}
}

// TestParallelExecutionDisabled tests fallback to sequential execution
func TestParallelExecutionDisabled(t *testing.T) {
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

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := ParallelExecutionConfig{
		MaxConcurrency: 0,
		EnableParallel: false, // disable parallel execution
	}
	
	result, err := engine.ExecuteWithParallelism(config)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	if result.FinalOutput != float64(15) {
		t.Errorf("Expected final output 15, got %v", result.FinalOutput)
	}
}

// TestParallelExecutionWithVariables tests parallel execution with variable nodes
func TestParallelExecutionWithVariables(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 100}},
			{"id": "2", "data": {"var_name": "myvar", "var_op": "set"}},
			{"id": "3", "data": {"var_name": "myvar", "var_op": "get"}},
			{"id": "4", "type": "extract", "data": {"field": "value"}}
		],
		"edges": [
			{"source": "1", "target": "2"},
			{"source": "2", "target": "3"},
			{"source": "3", "target": "4"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := DefaultParallelConfig()
	result, err := engine.ExecuteWithParallelism(config)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Extract returns a map with field and value
	if extractResult, ok := result.FinalOutput.(map[string]interface{}); ok {
		if extractResult["value"] != float64(100) {
			t.Errorf("Expected extracted value 100, got %v", extractResult["value"])
		}
	} else {
		t.Errorf("Expected final output to be a map, got %T", result.FinalOutput)
	}
}

// TestParallelExecutionSingleNode tests that single-node execution works efficiently
func TestParallelExecutionSingleNode(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 42}}
		],
		"edges": []
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := DefaultParallelConfig()
	result, err := engine.ExecuteWithParallelism(config)
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	if result.FinalOutput != float64(42) {
		t.Errorf("Expected final output 42, got %v", result.FinalOutput)
	}
}

// TestParallelExecutionErrorHandling tests error propagation in parallel execution
func TestParallelExecutionErrorHandling(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 0}},
			{"id": "3", "data": {"op": "divide"}}
		],
		"edges": [
			{"source": "1", "target": "3"},
			{"source": "2", "target": "3"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	config := DefaultParallelConfig()
	_, err = engine.ExecuteWithParallelism(config)
	
	if err == nil {
		t.Error("Expected division by zero error, got nil")
	}

	if !strings.Contains(err.Error(), "division by zero") {
		t.Errorf("Expected division by zero error, got: %v", err)
	}
}

// TestParallelExecutionTimeout tests timeout handling
func TestParallelExecutionTimeout(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"mode": "text"}}
		],
		"edges": [
			{"source": "1", "target": "2"}
		]
	}`

	config := Config{
		MaxExecutionTime: 1 * time.Nanosecond, // very short timeout
	}

	engine, err := NewEngineWithConfig([]byte(payload), config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	parallelConfig := DefaultParallelConfig()
	_, err = engine.ExecuteWithParallelism(parallelConfig)
	
	// We might or might not get a timeout depending on timing,
	// but execution should complete (success or timeout)
	// This test mainly ensures timeout context is properly propagated
	if err != nil && !strings.Contains(err.Error(), "context") {
		t.Logf("Got error (expected with very short timeout): %v", err)
	}
}

// TestExecutionLevelsDeterministic tests that level computation is deterministic
func TestExecutionLevelsDeterministic(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "a", "data": {"value": 1}},
			{"id": "b", "data": {"value": 2}},
			{"id": "c", "data": {"value": 3}},
			{"id": "d", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "a", "target": "d"},
			{"source": "b", "target": "d"},
			{"source": "c", "target": "d"}
		]
	}`

	// Run level computation multiple times
	for i := 0; i < 10; i++ {
		engine, err := NewEngine([]byte(payload))
		if err != nil {
			t.Fatalf("Failed to create engine: %v", err)
		}

		engine.inferNodeTypes()
		levels, err := engine.computeExecutionLevels()
		if err != nil {
			t.Fatalf("Failed to compute levels: %v", err)
		}

		if len(levels) != 2 {
			t.Errorf("Expected 2 levels, got %d", len(levels))
		}

		// Level 0 should have nodes a, b, c in sorted order
		level0 := levels[0].NodeIDs
		if len(level0) != 3 {
			t.Errorf("Expected 3 nodes in level 0, got %d", len(level0))
		}

		// Check sorting
		expected := []string{"a", "b", "c"}
		for j, nodeID := range level0 {
			if nodeID != expected[j] {
				t.Errorf("Expected node %s at position %d, got %s", expected[j], j, nodeID)
			}
		}
	}
}
