package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestSwitchExecutor_ValueMatching tests switch with value matching
func TestSwitchExecutor_ValueMatching(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		cases         []types.SwitchCase
		defaultPath   *string
		expectedMatch bool
		expectedPath  string
		description   string
	}{
		{
			name:  "Exact number match",
			input: float64(10),
			cases: []types.SwitchCase{
				{When: "==10", Value: float64(10), OutputPath: strPtr("path1")},
				{When: "==20", Value: float64(20), OutputPath: strPtr("path2")},
			},
			defaultPath:   strPtr("default"),
			expectedMatch: true,
			expectedPath:  "path1",
			description:   "Should match first case with value 10",
		},
		{
			name:  "String match",
			input: "hello",
			cases: []types.SwitchCase{
				{When: "==hello", Value: "hello", OutputPath: strPtr("greet")},
				{When: "==goodbye", Value: "goodbye", OutputPath: strPtr("farewell")},
			},
			defaultPath:   strPtr("default"),
			expectedMatch: true,
			expectedPath:  "greet",
			description:   "Should match string value",
		},
		{
			name:  "No match - use default",
			input: float64(30),
			cases: []types.SwitchCase{
				{When: "==10", Value: float64(10), OutputPath: strPtr("path1")},
				{When: "==20", Value: float64(20), OutputPath: strPtr("path2")},
			},
			defaultPath:   strPtr("fallback"),
			expectedMatch: false,
			expectedPath:  "fallback",
			description:   "Should use default path when no match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &SwitchExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeSwitch,
				Data: types.NodeData{
					Cases:       tt.cases,
					DefaultPath: tt.defaultPath,
				},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected result to be map, got %T", result)
			}

			matched, ok := resultMap["matched"].(bool)
			if !ok {
				t.Fatalf("Expected matched to be bool, got %T", resultMap["matched"])
			}

			if matched != tt.expectedMatch {
				t.Errorf("Expected matched=%v, got %v. Description: %s",
					tt.expectedMatch, matched, tt.description)
			}

			outputPath, ok := resultMap["output_path"].(string)
			if !ok {
				t.Fatalf("Expected output_path to be string, got %T", resultMap["output_path"])
			}

			if outputPath != tt.expectedPath {
				t.Errorf("Expected output_path='%s', got '%s'", tt.expectedPath, outputPath)
			}

			// Verify input value is preserved
			if resultMap["value"] != tt.input {
				t.Errorf("Expected value to be preserved: %v, got %v", tt.input, resultMap["value"])
			}
		})
	}
}

// TestSwitchExecutor_ConditionMatching tests switch with condition expressions
func TestSwitchExecutor_ConditionMatching(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		cases         []types.SwitchCase
		expectedMatch bool
		expectedCase  string
		description   string
	}{
		{
			name:  "Range match - less than",
			input: float64(5),
			cases: []types.SwitchCase{
				{When: "<10", OutputPath: strPtr("small")},
				{When: ">=10", OutputPath: strPtr("large")},
			},
			expectedMatch: true,
			expectedCase:  "<10",
			description:   "Should match <10 condition",
		},
		{
			name:  "Range match - greater than",
			input: float64(15),
			cases: []types.SwitchCase{
				{When: "<10", OutputPath: strPtr("small")},
				{When: ">=10", OutputPath: strPtr("large")},
			},
			expectedMatch: true,
			expectedCase:  ">=10",
			description:   "Should match >=10 condition",
		},
		{
			name:  "First matching case wins",
			input: float64(15),
			cases: []types.SwitchCase{
				{When: ">10", OutputPath: strPtr("first")},
				{When: ">5", OutputPath: strPtr("second")},
			},
			expectedMatch: true,
			expectedCase:  ">10",
			description:   "Should match first case when multiple match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &SwitchExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeSwitch,
				Data: types.NodeData{
					Cases: tt.cases,
				},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resultMap := result.(map[string]interface{})
			matched := resultMap["matched"].(bool)

			if matched != tt.expectedMatch {
				t.Errorf("Expected matched=%v, got %v. Description: %s",
					tt.expectedMatch, matched, tt.description)
			}

			if matched {
				caseStr, ok := resultMap["case"].(string)
				if !ok {
					t.Fatalf("Expected case to be string, got %T", resultMap["case"])
				}

				if caseStr != tt.expectedCase {
					t.Errorf("Expected case='%s', got '%s'", tt.expectedCase, caseStr)
				}
			}
		})
	}
}

// TestSwitchExecutor_DefaultPath tests default path handling
func TestSwitchExecutor_DefaultPath(t *testing.T) {
	exec := &SwitchExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {float64(100)},
		},
	}

	// No matching cases
	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeSwitch,
		Data: types.NodeData{
			Cases: []types.SwitchCase{
				{When: "<10", OutputPath: strPtr("small")},
			},
			DefaultPath: strPtr("other"),
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})

	if resultMap["matched"].(bool) {
		t.Error("Expected no match, got matched=true")
	}

	if resultMap["output_path"].(string) != "other" {
		t.Errorf("Expected output_path='other', got '%s'", resultMap["output_path"].(string))
	}
}

// TestSwitchExecutor_Validation tests node validation
func TestSwitchExecutor_Validation(t *testing.T) {
	tests := []struct {
		name        string
		node        types.Node
		expectError bool
	}{
		{
			name: "Valid node with cases",
			node: types.Node{
				Type: types.NodeTypeSwitch,
				Data: types.NodeData{
					Cases: []types.SwitchCase{
						{When: ">0"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "No cases",
			node: types.Node{
				Type: types.NodeTypeSwitch,
				Data: types.NodeData{
					Cases: []types.SwitchCase{},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &SwitchExecutor{}
			err := exec.Validate(tt.node)

			if tt.expectError && err == nil {
				t.Error("Expected validation error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
		})
	}
}

// TestSwitchExecutor_MissingInput tests error handling for missing input
func TestSwitchExecutor_MissingInput(t *testing.T) {
	exec := &SwitchExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{}, // No inputs
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeSwitch,
		Data: types.NodeData{
			Cases: []types.SwitchCase{
				{When: ">0"},
			},
		},
	}

	_, err := exec.Execute(ctx, node)
	if err == nil {
		t.Error("Expected error for missing input")
	}
}

// TestSwitchExecutor_NodeType tests NodeType method
func TestSwitchExecutor_NodeType(t *testing.T) {
	exec := &SwitchExecutor{}
	if exec.NodeType() != types.NodeTypeSwitch {
		t.Errorf("Expected NodeType to be %s, got %s", types.NodeTypeSwitch, exec.NodeType())
	}
}
