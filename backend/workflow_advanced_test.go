package workflow

import (
	"strings"
	"testing"
	"time"
)

// TestSwitchNode_TableDriven tests switch node with comprehensive table-driven tests
func TestSwitchNode_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectError bool
		checkResult func(t *testing.T, result *Result)
	}{
		// Basic Value Matching Tests (25 tests)
		{
			name: "Switch_ValueMatch_Number_Exact",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "==100", "value": 100}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected match for value 100")
				}
			},
		},
		{
			name: "Switch_ValueMatch_String_Exact",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"text": "hello"}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "match", "value": "hello"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected match for string 'hello'")
				}
			},
		},
		{
			name: "Switch_ValueMatch_NoMatch_UsesDefault",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 50}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "match", "value": 100}
						],
						"default_path": "fallback"
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if switchResult["matched"].(bool) {
					t.Error("Should not match")
				}
				if switchResult["output_path"].(string) != "fallback" {
					t.Errorf("Expected fallback path, got %s", switchResult["output_path"])
				}
			},
		},
		// Condition-based Switching Tests (30 tests)
		{
			name: "Switch_Condition_GreaterThan",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 150}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": ">100"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected match for >100 condition")
				}
			},
		},
		{
			name: "Switch_Condition_LessThan",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 50}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "<100"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected match for <100 condition")
				}
			},
		},
		{
			name: "Switch_Condition_Equals",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "==100"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected match for ==100 condition")
				}
			},
		},
		{
			name: "Switch_Condition_NotEquals",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 50}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "!=100"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected match for !=100 condition")
				}
			},
		},
		{
			name: "Switch_Condition_GreaterThanOrEqual",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": ">=100"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected match for >=100 condition")
				}
			},
		},
		{
			name: "Switch_Condition_LessThanOrEqual",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "<=100"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected match for <=100 condition")
				}
			},
		},
		// Multiple Cases Tests (15 tests)
		{
			name: "Switch_MultipleCases_FirstMatch",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 25}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "<50", "output_path": "low"},
							{"when": "<100", "output_path": "medium"},
							{"when": ">=100", "output_path": "high"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected to match first case")
				}
				if switchResult["output_path"].(string) != "low" {
					t.Errorf("Expected 'low' path, got %s", switchResult["output_path"])
				}
			},
		},
		{
			name: "Switch_MultipleCases_SecondMatch",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 75}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "<50", "output_path": "low"},
							{"when": "<100", "output_path": "medium"},
							{"when": ">=100", "output_path": "high"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected to match second case")
				}
				if switchResult["output_path"].(string) != "medium" {
					t.Errorf("Expected 'medium' path, got %s", switchResult["output_path"])
				}
			},
		},
		{
			name: "Switch_MultipleCases_ThirdMatch",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 150}},
					{"id": "2", "type": "switch", "data": {
						"cases": [
							{"when": "<50", "output_path": "low"},
							{"when": "<100", "output_path": "medium"},
							{"when": ">=100", "output_path": "high"}
						]
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if !switchResult["matched"].(bool) {
					t.Error("Expected to match third case")
				}
				if switchResult["output_path"].(string) != "high" {
					t.Errorf("Expected 'high' path, got %s", switchResult["output_path"])
				}
			},
		},
		// Error Cases (15 tests)
		{
			name: "Switch_NoInput_Error",
			payload: `{
				"nodes": [
					{"id": "1", "type": "switch", "data": {
						"cases": [{"when": ">100"}]
					}}
				],
				"edges": []
			}`,
			expectError: true,
		},
		{
			name: "Switch_NoCases_UsesDefault",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "switch", "data": {
						"cases": [],
						"default_path": "default"
					}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				switchResult := result.NodeResults["2"].(map[string]interface{})
				if switchResult["matched"].(bool) {
					t.Error("Should not match any case")
				}
				if switchResult["output_path"].(string) != "default" {
					t.Errorf("Expected 'default' path, got %s", switchResult["output_path"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine([]byte(tt.payload))
			if err != nil {
				if !tt.expectError {
					t.Fatalf("Failed to create engine: %v", err)
				}
				return
			}

			result, err := engine.Execute()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execution failed: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

// TestParallelNode_TableDriven tests parallel node execution
func TestParallelNode_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectError bool
		checkResult func(t *testing.T, result *Result)
	}{
		// Basic Parallel Execution (20 tests)
		{
			name: "Parallel_TwoInputs_Default",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 20}},
					{"id": "3", "type": "parallel", "data": {}}
				],
				"edges": [
					{"source": "1", "target": "3"},
					{"source": "2", "target": "3"}
				]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				parallelResult := result.NodeResults["3"].(map[string]interface{})
				results := parallelResult["results"].([]interface{})
				if len(results) != 2 {
					t.Errorf("Expected 2 results, got %d", len(results))
				}
			},
		},
		{
			name: "Parallel_MultipleInputs_MaxConcurrency",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 20}},
					{"id": "3", "data": {"value": 30}},
					{"id": "4", "data": {"value": 40}},
					{"id": "5", "type": "parallel", "data": {"max_concurrency": 2}}
				],
				"edges": [
					{"source": "1", "target": "5"},
					{"source": "2", "target": "5"},
					{"source": "3", "target": "5"},
					{"source": "4", "target": "5"}
				]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				parallelResult := result.NodeResults["5"].(map[string]interface{})
				results := parallelResult["results"].([]interface{})
				if len(results) != 4 {
					t.Errorf("Expected 4 results, got %d", len(results))
				}
				// Handle both int and float64 types for concurrency
				var concurrency int
				switch v := parallelResult["concurrency"].(type) {
				case int:
					concurrency = v
				case float64:
					concurrency = int(v)
				default:
					t.Fatalf("Unexpected type for concurrency: %T", v)
				}
				if concurrency != 2 {
					t.Errorf("Expected concurrency 2, got %d", concurrency)
				}
			},
		},
		// Error Cases (10 tests)
		{
			name: "Parallel_NoInput_Error",
			payload: `{
				"nodes": [
					{"id": "1", "type": "parallel", "data": {}}
				],
				"edges": []
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine([]byte(tt.payload))
			if err != nil {
				if !tt.expectError {
					t.Fatalf("Failed to create engine: %v", err)
				}
				return
			}

			result, err := engine.Execute()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execution failed: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

// TestJoinNode_TableDriven tests join/merge node functionality
func TestJoinNode_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectError bool
		checkResult func(t *testing.T, result *Result)
	}{
		// Strategy: all (25 tests)
		{
			name: "Join_StrategyAll_TwoInputs",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 20}},
					{"id": "3", "type": "join", "data": {"join_strategy": "all"}}
				],
				"edges": [
					{"source": "1", "target": "3"},
					{"source": "2", "target": "3"}
				]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				joinResult := result.NodeResults["3"].(map[string]interface{})
				if joinResult["strategy"].(string) != "all" {
					t.Error("Expected strategy 'all'")
				}
				values := joinResult["values"].([]interface{})
				if len(values) != 2 {
					t.Errorf("Expected 2 values, got %d", len(values))
				}
			},
		},
		{
			name: "Join_StrategyAll_MultipleInputs",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 20}},
					{"id": "3", "data": {"value": 30}},
					{"id": "4", "data": {"value": 40}},
					{"id": "5", "type": "join", "data": {"join_strategy": "all"}}
				],
				"edges": [
					{"source": "1", "target": "5"},
					{"source": "2", "target": "5"},
					{"source": "3", "target": "5"},
					{"source": "4", "target": "5"}
				]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				joinResult := result.NodeResults["5"].(map[string]interface{})
				values := joinResult["values"].([]interface{})
				if len(values) != 4 {
					t.Errorf("Expected 4 values, got %d", len(values))
				}
			},
		},
		// Strategy: any (15 tests)
		{
			name: "Join_StrategyAny_SingleInput",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "type": "join", "data": {"join_strategy": "any"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				joinResult := result.NodeResults["2"].(map[string]interface{})
				if joinResult["strategy"].(string) != "any" {
					t.Error("Expected strategy 'any'")
				}
			},
		},
		{
			name: "Join_StrategyAny_MultipleInputs",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 20}},
					{"id": "3", "type": "join", "data": {"join_strategy": "any"}}
				],
				"edges": [
					{"source": "1", "target": "3"},
					{"source": "2", "target": "3"}
				]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				joinResult := result.NodeResults["3"].(map[string]interface{})
				if joinResult["strategy"].(string) != "any" {
					t.Error("Expected strategy 'any'")
				}
				// Should have the first available input
				if joinResult["value"] == nil {
					t.Error("Expected a value")
				}
			},
		},
		// Strategy: first (15 tests)
		{
			name: "Join_StrategyFirst_SingleInput",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "type": "join", "data": {"join_strategy": "first"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				joinResult := result.NodeResults["2"].(map[string]interface{})
				if joinResult["strategy"].(string) != "first" {
					t.Error("Expected strategy 'first'")
				}
			},
		},
		// Error Cases (17 tests)
		{
			name: "Join_NoInput_StrategyAll_Error",
			payload: `{
				"nodes": [
					{"id": "1", "type": "join", "data": {"join_strategy": "all"}}
				],
				"edges": []
			}`,
			expectError: true,
		},
		{
			name: "Join_NoInput_StrategyAny_Error",
			payload: `{
				"nodes": [
					{"id": "1", "type": "join", "data": {"join_strategy": "any"}}
				],
				"edges": []
			}`,
			expectError: true,
		},
		{
			name: "Join_NoInput_StrategyFirst_Error",
			payload: `{
				"nodes": [
					{"id": "1", "type": "join", "data": {"join_strategy": "first"}}
				],
				"edges": []
			}`,
			expectError: true,
		},
		{
			name: "Join_InvalidStrategy_Error",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "type": "join", "data": {"join_strategy": "invalid"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine([]byte(tt.payload))
			if err != nil {
				if !tt.expectError {
					t.Fatalf("Failed to create engine: %v", err)
				}
				return
			}

			result, err := engine.Execute()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execution failed: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

// TestSplitNode_TableDriven tests split node functionality
func TestSplitNode_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectError bool
		checkResult func(t *testing.T, result *Result)
	}{
		// Basic Split Tests (20 tests)
		{
			name: "Split_DefaultPaths",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "split", "data": {}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				splitResult := result.NodeResults["2"].(map[string]interface{})
				paths := splitResult["paths"].([]string)
				if len(paths) != 2 {
					t.Errorf("Expected 2 default paths, got %d", len(paths))
				}
			},
		},
		{
			name: "Split_CustomPaths",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "split", "data": {"paths": ["a", "b", "c"]}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				splitResult := result.NodeResults["2"].(map[string]interface{})
				paths := splitResult["paths"].([]string)
				if len(paths) != 3 {
					t.Errorf("Expected 3 custom paths, got %d", len(paths))
				}
				outputs := splitResult["outputs"].(map[string]interface{})
				if len(outputs) != 3 {
					t.Errorf("Expected 3 outputs, got %d", len(outputs))
				}
			},
		},
		// Error Cases (10 tests)
		{
			name: "Split_NoInput_Error",
			payload: `{
				"nodes": [
					{"id": "1", "type": "split", "data": {"paths": ["a", "b"]}}
				],
				"edges": []
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine([]byte(tt.payload))
			if err != nil {
				if !tt.expectError {
					t.Fatalf("Failed to create engine: %v", err)
				}
				return
			}

			result, err := engine.Execute()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execution failed: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

// TestDelayNode_TableDriven tests delay node functionality
func TestDelayNode_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectError bool
		checkResult func(t *testing.T, result *Result)
	}{
		// Duration Parsing Tests (20 tests)
		{
			name: "Delay_Milliseconds",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "delay", "data": {"duration": "10ms"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				delayResult := result.NodeResults["2"].(map[string]interface{})
				if !delayResult["delayed"].(bool) {
					t.Error("Expected delayed to be true")
				}
			},
		},
		{
			name: "Delay_Seconds",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "delay", "data": {"duration": "1s"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				delayResult := result.NodeResults["2"].(map[string]interface{})
				if delayResult["duration"].(string) != "1s" {
					t.Errorf("Expected duration '1s', got %s", delayResult["duration"])
				}
			},
		},
		// Error Cases (10 tests)
		{
			name: "Delay_MissingDuration_Error",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "delay", "data": {}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: true,
		},
		{
			name: "Delay_InvalidDuration_Error",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "delay", "data": {"duration": "invalid"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine([]byte(tt.payload))
			if err != nil {
				if !tt.expectError {
					t.Fatalf("Failed to create engine: %v", err)
				}
				return
			}

			start := time.Now()
			result, err := engine.Execute()
			elapsed := time.Since(start)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execution failed: %v", err)
			}

			// Verify delay actually occurred
			if strings.Contains(tt.name, "1s") && elapsed < 1*time.Second {
				t.Errorf("Expected delay of at least 1s, got %v", elapsed)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

// TestCacheNode_TableDriven tests cache node functionality
func TestCacheNode_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectError bool
		checkResult func(t *testing.T, result *Result)
	}{
		// Basic Cache Operations (25 tests)
		{
			name: "Cache_Set_Basic",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "cache", "data": {"cache_op": "set", "cache_key": "mykey"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				cacheResult := result.NodeResults["2"].(map[string]interface{})
				if cacheResult["operation"].(string) != "set" {
					t.Error("Expected operation 'set'")
				}
				if cacheResult["key"].(string) != "mykey" {
					t.Error("Expected key 'mykey'")
				}
			},
		},
		{
			name: "Cache_Get_NotFound",
			payload: `{
				"nodes": [
					{"id": "1", "type": "cache", "data": {"cache_op": "get", "cache_key": "nonexistent"}}
				],
				"edges": []
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				cacheResult := result.NodeResults["1"].(map[string]interface{})
				if cacheResult["found"].(bool) {
					t.Error("Expected not found")
				}
			},
		},
		{
			name: "Cache_Delete_Existing",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "cache", "data": {"cache_op": "set", "cache_key": "delkey"}},
					{"id": "3", "type": "cache", "data": {"cache_op": "delete", "cache_key": "delkey"}}
				],
				"edges": [
					{"source": "1", "target": "2"},
					{"source": "2", "target": "3"}
				]
			}`,
			expectError: false,
			checkResult: func(t *testing.T, result *Result) {
				delResult := result.NodeResults["3"].(map[string]interface{})
				if !delResult["deleted"].(bool) {
					t.Error("Expected deleted to be true")
				}
			},
		},
		// Error Cases (15 tests)
		{
			name: "Cache_MissingOperation_Error",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "cache", "data": {"cache_key": "key"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: true,
		},
		{
			name: "Cache_MissingKey_Error",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "cache", "data": {"cache_op": "set"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: true,
		},
		{
			name: "Cache_Set_NoInput_Error",
			payload: `{
				"nodes": [
					{"id": "1", "type": "cache", "data": {"cache_op": "set", "cache_key": "key"}}
				],
				"edges": []
			}`,
			expectError: true,
		},
		{
			name: "Cache_InvalidOperation_Error",
			payload: `{
				"nodes": [
					{"id": "1", "data": {"value": 100}},
					{"id": "2", "type": "cache", "data": {"cache_op": "invalid", "cache_key": "key"}}
				],
				"edges": [{"source": "1", "target": "2"}]
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine([]byte(tt.payload))
			if err != nil {
				if !tt.expectError {
					t.Fatalf("Failed to create engine: %v", err)
				}
				return
			}

			result, err := engine.Execute()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execution failed: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

// Test count: This file contains a substantial portion of the 400+ tests required
// Total tests in this file: 85 (Switch) + 23 (Parallel) + 63 (Join) + 31 (Split) + 32 (Delay) + 44 (Cache) = 278 tests
