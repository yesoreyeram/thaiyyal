package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestForEachExecutor_MapMode tests MAP mode (transform each element)
func TestForEachExecutor_MapMode(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		mode          string
		expectedCount int
		expectError   bool
		description   string
	}{
		{
			name:          "Small array",
			input:         []interface{}{1, 2, 3},
			mode:          "map",
			expectedCount: 3,
			expectError:   false,
			description:   "Should map over 3 elements",
		},
		{
			name:          "Empty array",
			input:         []interface{}{},
			mode:          "map",
			expectedCount: 0,
			expectError:   false,
			description:   "Should handle empty array",
		},
		{
			name:          "Array of objects",
			input:         []interface{}{map[string]interface{}{"id": 1}, map[string]interface{}{"id": 2}},
			mode:          "map",
			expectedCount: 2,
			expectError:   false,
			description:   "Should map over objects",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ForEachExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeForEach,
				Data: types.NodeData{
					Mode: &tt.mode,
				},
			}

			result, err := exec.Execute(ctx, node)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected result to be map, got %T", result)
			}

			// Verify mode
			if resultMap["mode"].(string) != "map" {
				t.Errorf("Expected mode=map, got %s", resultMap["mode"])
			}

			// Verify results array
			results, ok := resultMap["results"].([]interface{})
			if !ok {
				t.Fatalf("Expected results to be array, got %T", resultMap["results"])
			}

			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
			}
		})
	}
}

// TestForEachExecutor_MaxIterations tests iteration limits
func TestForEachExecutor_MaxIterations(t *testing.T) {
	tests := []struct {
		name          string
		arraySize     int
		maxIterations *int
		expectError   bool
		description   string
	}{
		{
			name:          "Within default limit",
			arraySize:     500,
			maxIterations: nil,
			expectError:   false,
			description:   "500 elements should be within default 1000 limit",
		},
		{
			name:          "Exceeds default limit",
			arraySize:     1500,
			maxIterations: nil,
			expectError:   true,
			description:   "1500 elements should exceed default 1000 limit",
		},
		{
			name:          "Within custom limit",
			arraySize:     50,
			maxIterations: intPtr(100),
			expectError:   false,
			description:   "50 elements should be within custom 100 limit",
		},
		{
			name:          "Exceeds custom limit",
			arraySize:     150,
			maxIterations: intPtr(100),
			expectError:   true,
			description:   "150 elements should exceed custom 100 limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ForEachExecutor{}
			
			// Create array of specified size
			input := make([]interface{}, tt.arraySize)
			for i := range input {
				input[i] = i
			}

			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {input},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeForEach,
				Data: types.NodeData{
					MaxIterations: tt.maxIterations,
				},
			}

			_, err := exec.Execute(ctx, node)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s, got nil", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for %s, got: %v", tt.description, err)
			}
		})
	}
}

// TestForEachExecutor_NonArrayInput tests error handling for non-array inputs
func TestForEachExecutor_NonArrayInput(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"String input", "not an array"},
		{"Number input", float64(42)},
		{"Object input", map[string]interface{}{"key": "value"}},
		{"Boolean input", true},
		{"Null input", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ForEachExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeForEach,
				Data: types.NodeData{},
			}

			_, err := exec.Execute(ctx, node)
			if err == nil {
				t.Errorf("Expected error for non-array input %T, got nil", tt.input)
			}
		})
	}
}

// TestForEachExecutor_MissingInput tests error handling for missing input
func TestForEachExecutor_MissingInput(t *testing.T) {
	exec := &ForEachExecutor{}
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{}, // No inputs
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeForEach,
		Data: types.NodeData{},
	}

	_, err := exec.Execute(ctx, node)
	if err == nil {
		t.Error("Expected error for missing input")
	}
}

// TestForEachExecutor_Validation tests node validation
func TestForEachExecutor_Validation(t *testing.T) {
	exec := &ForEachExecutor{}
	
	// ForEach has no required fields, so validation should always pass
	node := types.Node{
		Type: types.NodeTypeForEach,
		Data: types.NodeData{},
	}

	err := exec.Validate(node)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

// TestForEachExecutor_NodeType tests NodeType method
func TestForEachExecutor_NodeType(t *testing.T) {
	exec := &ForEachExecutor{}
	if exec.NodeType() != types.NodeTypeForEach {
		t.Errorf("Expected NodeType to be %s, got %s", types.NodeTypeForEach, exec.NodeType())
	}
}

// Helper function for int pointers
func intPtr(i int) *int {
	return &i
}
