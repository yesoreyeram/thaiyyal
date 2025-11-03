package executor

import (
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// MockExecutionContext is a test implementation of ExecutionContext
type MockExecutionContext struct {
	inputs      map[string][]interface{}
	variables   map[string]interface{}
	nodeResults map[string]interface{}
	contextVars map[string]interface{}
}

func (m *MockExecutionContext) GetNodeInputs(nodeID string) []interface{} {
	if inputs, ok := m.inputs[nodeID]; ok {
		return inputs
	}
	return nil
}

func (m *MockExecutionContext) GetNode(nodeID string) *types.Node {
	return nil
}

func (m *MockExecutionContext) GetVariable(name string) (interface{}, error) {
	if val, ok := m.variables[name]; ok {
		return val, nil
	}
	return nil, nil
}

func (m *MockExecutionContext) SetVariable(name string, value interface{}) error {
	m.variables[name] = value
	return nil
}

func (m *MockExecutionContext) GetAccumulator() interface{} {
	return nil
}

func (m *MockExecutionContext) SetAccumulator(value interface{}) {
}

func (m *MockExecutionContext) GetCounter() float64 {
	return 0
}

func (m *MockExecutionContext) SetCounter(value float64) {
}

func (m *MockExecutionContext) GetCache(key string) (interface{}, bool) {
	return nil, false
}

func (m *MockExecutionContext) SetCache(key string, value interface{}, ttl time.Duration) {
}

func (m *MockExecutionContext) GetWorkflowContext() map[string]interface{} {
	return m.contextVars
}

func (m *MockExecutionContext) GetContextVariable(name string) (interface{}, bool) {
	if m.contextVars == nil {
		return nil, false
	}
	val, ok := m.contextVars[name]
	return val, ok
}

func (m *MockExecutionContext) SetContextVariable(name string, value interface{}) {
	if m.contextVars == nil {
		m.contextVars = make(map[string]interface{})
	}
	m.contextVars[name] = value
}

func (m *MockExecutionContext) GetContextConstant(name string) (interface{}, bool) {
	return m.GetContextVariable(name)
}

func (m *MockExecutionContext) SetContextConstant(name string, value interface{}) {
	m.SetContextVariable(name, value)
}

func (m *MockExecutionContext) InterpolateTemplate(template string) string {
	return template
}

func (m *MockExecutionContext) GetNodeResult(nodeID string) (interface{}, bool) {
	if m.nodeResults == nil {
		return nil, false
	}
	val, ok := m.nodeResults[nodeID]
	return val, ok
}

func (m *MockExecutionContext) SetNodeResult(nodeID string, result interface{}) {
	if m.nodeResults == nil {
		m.nodeResults = make(map[string]interface{})
	}
	m.nodeResults[nodeID] = result
}

func (m *MockExecutionContext) GetAllNodeResults() map[string]interface{} {
	if m.nodeResults == nil {
		return make(map[string]interface{})
	}
	return m.nodeResults
}

func (m *MockExecutionContext) GetVariables() map[string]interface{} {
	if m.variables == nil {
		return make(map[string]interface{})
	}
	return m.variables
}

func (m *MockExecutionContext) GetContextVariables() map[string]interface{} {
	if m.contextVars == nil {
		return make(map[string]interface{})
	}
	return m.contextVars
}

func (m *MockExecutionContext) GetConfig() types.Config {
	return types.DefaultConfig()
}

func (m *MockExecutionContext) GetHTTPClientRegistry() interface{} {
	return nil
}

func (m *MockExecutionContext) IncrementNodeExecution() error {
	return nil
}

func (m *MockExecutionContext) IncrementHTTPCall() error {
	return nil
}

func (m *MockExecutionContext) GetNodeExecutionCount() int {
	return 0
}

func (m *MockExecutionContext) GetHTTPCallCount() int {
	return 0
}

// TestFilterExecutor_Basic tests basic array filtering functionality
func TestFilterExecutor_Basic(t *testing.T) {
	tests := []struct {
		name           string
		condition      string
		inputArray     []interface{}
		expectedCount  int
		expectedError  bool
		description    string
	}{
		{
			name:      "Filter numbers greater than 10",
			condition: "variables.item > 10",
			inputArray: []interface{}{
				float64(5),
				float64(15),
				float64(8),
				float64(20),
				float64(12),
			},
			expectedCount: 3, // 15, 20, 12
			description:   "Should filter items where value > 10",
		},
		{
			name:      "Filter numbers equal to 5",
			condition: "variables.item == 5",
			inputArray: []interface{}{
				float64(5),
				float64(10),
				float64(5),
				float64(3),
			},
			expectedCount: 2, // two 5's
			description:   "Should filter items where value == 5",
		},
		{
			name:      "Filter all (condition always true)",
			condition: "true",
			inputArray: []interface{}{
				float64(1),
				float64(2),
				float64(3),
			},
			expectedCount: 3, // all items
			description:   "Should include all items when condition is always true",
		},
		{
			name:      "Filter none (condition always false)",
			condition: "false",
			inputArray: []interface{}{
				float64(1),
				float64(2),
				float64(3),
			},
			expectedCount: 0, // no items
			description:   "Should exclude all items when condition is always false",
		},
		{
			name:          "Empty array",
			condition:     "variables.item > 10",
			inputArray:    []interface{}{},
			expectedCount: 0,
			description:   "Should handle empty arrays gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FilterExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
				variables: make(map[string]interface{}),
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeFilter,
				Data: types.NodeData{
					Condition: &tt.condition,
				},
			}

			result, err := exec.Execute(ctx, node)
			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if err != nil {
				return
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected result to be map, got %T", result)
			}

			filtered, ok := resultMap["filtered"].([]interface{})
			if !ok {
				t.Fatalf("Expected filtered to be array, got %T", resultMap["filtered"])
			}

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d filtered items, got %d. Description: %s", 
					tt.expectedCount, len(filtered), tt.description)
			}

			// Verify metadata
			if resultMap["is_array"] != true {
				t.Errorf("Expected is_array to be true")
			}
			if resultMap["input_count"] != len(tt.inputArray) {
				t.Errorf("Expected input_count to be %d, got %v", 
					len(tt.inputArray), resultMap["input_count"])
			}
			if resultMap["output_count"] != tt.expectedCount {
				t.Errorf("Expected output_count to be %d, got %v", 
					tt.expectedCount, resultMap["output_count"])
			}
		})
	}
}

// TestFilterExecutor_ObjectFields tests filtering objects by field values
func TestFilterExecutor_ObjectFields(t *testing.T) {
	tests := []struct {
		name          string
		condition     string
		inputArray    []interface{}
		expectedCount int
		description   string
	}{
		{
			name:      "Filter objects by field value",
			condition: "variables.item.age > 18",
			inputArray: []interface{}{
				map[string]interface{}{"name": "Alice", "age": float64(25)},
				map[string]interface{}{"name": "Bob", "age": float64(17)},
				map[string]interface{}{"name": "Charlie", "age": float64(30)},
				map[string]interface{}{"name": "Dave", "age": float64(16)},
			},
			expectedCount: 2, // Alice and Charlie
			description:   "Should filter objects where age > 18",
		},
		{
			name:      "Filter by string field equality",
			condition: "variables.item.status == \"active\"",
			inputArray: []interface{}{
				map[string]interface{}{"id": float64(1), "status": "active"},
				map[string]interface{}{"id": float64(2), "status": "inactive"},
				map[string]interface{}{"id": float64(3), "status": "active"},
			},
			expectedCount: 2, // id 1 and 3
			description:   "Should filter objects where status == 'active'",
		},
		{
			name:      "Filter by nested field",
			condition: "variables.item.profile.verified == true",
			inputArray: []interface{}{
				map[string]interface{}{
					"id": float64(1),
					"profile": map[string]interface{}{
						"verified": true,
					},
				},
				map[string]interface{}{
					"id": float64(2),
					"profile": map[string]interface{}{
						"verified": false,
					},
				},
				map[string]interface{}{
					"id": float64(3),
					"profile": map[string]interface{}{
						"verified": true,
					},
				},
			},
			expectedCount: 2, // id 1 and 3
			description:   "Should filter objects by nested field values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FilterExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
				variables: make(map[string]interface{}),
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeFilter,
				Data: types.NodeData{
					Condition: &tt.condition,
				},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resultMap := result.(map[string]interface{})
			filtered := resultMap["filtered"].([]interface{})

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d filtered items, got %d. Description: %s",
					tt.expectedCount, len(filtered), tt.description)
			}
		})
	}
}

// TestFilterExecutor_ComplexConditions tests complex boolean conditions
func TestFilterExecutor_ComplexConditions(t *testing.T) {
	tests := []struct {
		name          string
		condition     string
		inputArray    []interface{}
		expectedCount int
		description   string
	}{
		{
			name:      "AND condition",
			condition: "variables.item.price < 50 && variables.item.category == \"books\"",
			inputArray: []interface{}{
				map[string]interface{}{"price": float64(30), "category": "books"},
				map[string]interface{}{"price": float64(60), "category": "books"},
				map[string]interface{}{"price": float64(25), "category": "electronics"},
				map[string]interface{}{"price": float64(45), "category": "books"},
			},
			expectedCount: 2, // price < 50 AND category == "books"
			description:   "Should filter with AND condition",
		},
		{
			name:      "OR condition",
			condition: "variables.item.priority == \"high\" || variables.item.urgent == true",
			inputArray: []interface{}{
				map[string]interface{}{"priority": "high", "urgent": false},
				map[string]interface{}{"priority": "low", "urgent": true},
				map[string]interface{}{"priority": "medium", "urgent": false},
				map[string]interface{}{"priority": "high", "urgent": true},
			},
			expectedCount: 3, // high priority OR urgent
			description:   "Should filter with OR condition",
		},
		{
			name:      "Range condition",
			condition: "variables.item >= 10 && variables.item <= 20",
			inputArray: []interface{}{
				float64(5),
				float64(10),
				float64(15),
				float64(20),
				float64(25),
			},
			expectedCount: 3, // 10, 15, 20
			description:   "Should filter values in range [10, 20]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FilterExecutor{}
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.inputArray},
				},
				variables: make(map[string]interface{}),
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeFilter,
				Data: types.NodeData{
					Condition: &tt.condition,
				},
			}

			result, err := exec.Execute(ctx, node)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resultMap := result.(map[string]interface{})
			filtered := resultMap["filtered"].([]interface{})

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d filtered items, got %d. Description: %s",
					tt.expectedCount, len(filtered), tt.description)
			}
		})
	}
}

// TestFilterExecutor_NonArrayInput tests behavior with non-array inputs
func TestFilterExecutor_NonArrayInput(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		checkResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:  "String input",
			input: "hello world",
			checkResult: func(t *testing.T, result map[string]interface{}) {
				if result["filtered"] != "hello world" {
					t.Errorf("Expected filtered to be 'hello world', got %v", result["filtered"])
				}
			},
		},
		{
			name:  "Number input",
			input: float64(42),
			checkResult: func(t *testing.T, result map[string]interface{}) {
				if result["filtered"] != float64(42) {
					t.Errorf("Expected filtered to be 42, got %v", result["filtered"])
				}
			},
		},
		{
			name:  "Object input",
			input: map[string]interface{}{"key": "value"},
			checkResult: func(t *testing.T, result map[string]interface{}) {
				filtered, ok := result["filtered"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected filtered to be a map, got %T", result["filtered"])
					return
				}
				if filtered["key"] != "value" {
					t.Errorf("Expected filtered['key'] to be 'value', got %v", filtered["key"])
				}
			},
		},
		{
			name:  "Boolean input",
			input: true,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				if result["filtered"] != true {
					t.Errorf("Expected filtered to be true, got %v", result["filtered"])
				}
			},
		},
		{
			name:  "Null input",
			input: nil,
			checkResult: func(t *testing.T, result map[string]interface{}) {
				if result["filtered"] != nil {
					t.Errorf("Expected filtered to be nil, got %v", result["filtered"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FilterExecutor{}
			condition := "true"
			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"test-node": {tt.input},
				},
				variables: make(map[string]interface{}),
			}

			node := types.Node{
				ID:   "test-node",
				Type: types.NodeTypeFilter,
				Data: types.NodeData{
					Condition: &condition,
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

			// Check the filtered value
			tt.checkResult(t, resultMap)

			if resultMap["is_array"] != false {
				t.Errorf("Expected is_array to be false")
			}

			if _, hasWarning := resultMap["warning"]; !hasWarning {
				t.Errorf("Expected warning field to be present")
			}
		})
	}
}

// TestFilterExecutor_Validation tests input validation
func TestFilterExecutor_Validation(t *testing.T) {
	tests := []struct {
		name        string
		node        types.Node
		expectError bool
		description string
	}{
		{
			name: "Valid node",
			node: types.Node{
				Type: types.NodeTypeFilter,
				Data: types.NodeData{
					Condition: strPtr("variables.item > 10"),
				},
			},
			expectError: false,
			description: "Should validate valid node",
		},
		{
			name: "Missing condition",
			node: types.Node{
				Type: types.NodeTypeFilter,
				Data: types.NodeData{},
			},
			expectError: true,
			description: "Should reject node without condition",
		},
		{
			name: "Empty condition",
			node: types.Node{
				Type: types.NodeTypeFilter,
				Data: types.NodeData{
					Condition: strPtr(""),
				},
			},
			expectError: true,
			description: "Should reject node with empty condition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &FilterExecutor{}
			err := exec.Validate(tt.node)

			if tt.expectError && err == nil {
				t.Errorf("Expected validation error but got none. Description: %s", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected validation error: %v. Description: %s", err, tt.description)
			}
		})
	}
}

// TestFilterExecutor_MissingInput tests error handling for missing input
func TestFilterExecutor_MissingInput(t *testing.T) {
	exec := &FilterExecutor{}
	condition := "variables.item > 10"
	ctx := &MockExecutionContext{
		inputs:    map[string][]interface{}{},
		variables: make(map[string]interface{}),
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeFilter,
		Data: types.NodeData{
			Condition: &condition,
		},
	}

	_, err := exec.Execute(ctx, node)
	if err == nil {
		t.Errorf("Expected error for missing input but got none")
	}
}

// TestFilterExecutor_NodeType tests that NodeType returns correct value
func TestFilterExecutor_NodeType(t *testing.T) {
	exec := &FilterExecutor{}
	if exec.NodeType() != types.NodeTypeFilter {
		t.Errorf("Expected NodeType to be %s, got %s", types.NodeTypeFilter, exec.NodeType())
	}
}

// TestFilterExecutor_WithContextVariables tests filtering with context variables
func TestFilterExecutor_WithContextVariables(t *testing.T) {
	exec := &FilterExecutor{}
	condition := "variables.item > context.threshold"
	
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {[]interface{}{
				float64(5),
				float64(15),
				float64(8),
				float64(20),
			}},
		},
		variables: make(map[string]interface{}),
		contextVars: map[string]interface{}{
			"threshold": float64(10),
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeFilter,
		Data: types.NodeData{
			Condition: &condition,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	filtered := resultMap["filtered"].([]interface{})

	// Should filter items > 10 (threshold from context)
	if len(filtered) != 2 { // 15, 20
		t.Errorf("Expected 2 filtered items, got %d", len(filtered))
	}
}

// TestFilterExecutor_WithNodeReferences tests filtering with references to other nodes
func TestFilterExecutor_WithNodeReferences(t *testing.T) {
	exec := &FilterExecutor{}
	condition := "variables.item.value > node.threshold.value"
	
	ctx := &MockExecutionContext{
		inputs: map[string][]interface{}{
			"test-node": {[]interface{}{
				map[string]interface{}{"value": float64(5)},
				map[string]interface{}{"value": float64(15)},
				map[string]interface{}{"value": float64(25)},
			}},
		},
		variables: make(map[string]interface{}),
		nodeResults: map[string]interface{}{
			"threshold": map[string]interface{}{"value": float64(10)},
		},
	}

	node := types.Node{
		ID:   "test-node",
		Type: types.NodeTypeFilter,
		Data: types.NodeData{
			Condition: &condition,
		},
	}

	result, err := exec.Execute(ctx, node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resultMap := result.(map[string]interface{})
	filtered := resultMap["filtered"].([]interface{})

	// Should filter items where value > 10
	if len(filtered) != 2 { // 15, 25
		t.Errorf("Expected 2 filtered items, got %d", len(filtered))
	}
}

// Helper function for string pointers
func strPtr(s string) *string {
	return &s
}

// TestFilterExecutor_DirectFieldAccess tests the new direct field access syntax
func TestFilterExecutor_DirectFieldAccess(t *testing.T) {
tests := []struct {
name          string
condition     string
inputArray    []interface{}
expectedCount int
description   string
}{
{
name:      "Direct field access - simple",
condition: "age > 18",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "age": float64(25)},
map[string]interface{}{"name": "Bob", "age": float64(17)},
map[string]interface{}{"name": "Charlie", "age": float64(30)},
},
expectedCount: 2, // Alice and Charlie
description:   "Should filter using direct field access: age > 18",
},
{
name:      "Direct field access with variable comparison",
condition: "age >= variables.minAge",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "age": float64(25)},
map[string]interface{}{"name": "Bob", "age": float64(19)},
map[string]interface{}{"name": "Charlie", "age": float64(22)},
},
expectedCount: 2, // Alice and Charlie (if minAge=21)
description:   "Should filter using direct field and variable: age >= variables.minAge",
},
{
name:      "Direct nested field access",
condition: "profile.verified == true",
inputArray: []interface{}{
map[string]interface{}{
"name": "Alice",
"profile": map[string]interface{}{
"verified": true,
},
},
map[string]interface{}{
"name": "Bob",
"profile": map[string]interface{}{
"verified": false,
},
},
map[string]interface{}{
"name": "Charlie",
"profile": map[string]interface{}{
"verified": true,
},
},
},
expectedCount: 2, // Alice and Charlie
description:   "Should filter using nested field access: profile.verified == true",
},
{
name:      "Direct field with AND operator",
condition: "age >= 18 && status == \"active\"",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "age": float64(25), "status": "active"},
map[string]interface{}{"name": "Bob", "age": float64(19), "status": "inactive"},
map[string]interface{}{"name": "Charlie", "age": float64(22), "status": "active"},
map[string]interface{}{"name": "Dave", "age": float64(17), "status": "active"},
},
expectedCount: 2, // Alice and Charlie
description:   "Should filter using direct fields with AND: age >= 18 && status == active",
},
{
name:      "Direct field with context variable",
condition: "score > context.threshold",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "score": float64(85)},
map[string]interface{}{"name": "Bob", "score": float64(65)},
map[string]interface{}{"name": "Charlie", "score": float64(95)},
},
expectedCount: 2, // Alice and Charlie (if threshold=70)
description:   "Should filter using direct field and context: score > context.threshold",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
exec := &FilterExecutor{}
ctx := &MockExecutionContext{
inputs: map[string][]interface{}{
"test-node": {tt.inputArray},
},
variables: map[string]interface{}{
"minAge": float64(21),
},
contextVars: map[string]interface{}{
"threshold": float64(70),
},
}

node := types.Node{
ID:   "test-node",
Type: types.NodeTypeFilter,
Data: types.NodeData{
Condition: &tt.condition,
},
}

result, err := exec.Execute(ctx, node)
if err != nil {
t.Fatalf("Unexpected error: %v", err)
}

resultMap := result.(map[string]interface{})
filtered := resultMap["filtered"].([]interface{})

if len(filtered) != tt.expectedCount {
t.Errorf("Expected %d filtered items, got %d. Description: %s",
tt.expectedCount, len(filtered), tt.description)
}
})
}
}

// TestFilterExecutor_BackwardCompatibility ensures old syntax still works
func TestFilterExecutor_BackwardCompatibility(t *testing.T) {
tests := []struct {
name          string
condition     string
inputArray    []interface{}
expectedCount int
}{
{
name:      "Old syntax: variables.item.age",
condition: "variables.item.age > 18",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "age": float64(25)},
map[string]interface{}{"name": "Bob", "age": float64(17)},
},
expectedCount: 1, // Alice
},
{
name:      "Mixed: new and old syntax",
condition: "age > 18 && variables.item.status == \"active\"",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "age": float64(25), "status": "active"},
map[string]interface{}{"name": "Bob", "age": float64(19), "status": "inactive"},
},
expectedCount: 1, // Alice
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
exec := &FilterExecutor{}
ctx := &MockExecutionContext{
inputs: map[string][]interface{}{
"test-node": {tt.inputArray},
},
variables: make(map[string]interface{}),
}

node := types.Node{
ID:   "test-node",
Type: types.NodeTypeFilter,
Data: types.NodeData{
Condition: &tt.condition,
},
}

result, err := exec.Execute(ctx, node)
if err != nil {
t.Fatalf("Unexpected error: %v", err)
}

resultMap := result.(map[string]interface{})
filtered := resultMap["filtered"].([]interface{})

if len(filtered) != tt.expectedCount {
t.Errorf("Expected %d filtered items, got %d",
tt.expectedCount, len(filtered))
}
})
}
}

// TestFilterExecutor_ItemSyntax tests the recommended item.field syntax
func TestFilterExecutor_ItemSyntax(t *testing.T) {
tests := []struct {
name          string
condition     string
inputArray    []interface{}
expectedCount int
description   string
}{
{
name:      "item.field syntax - simple",
condition: "item.age > 18",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "age": float64(25)},
map[string]interface{}{"name": "Bob", "age": float64(17)},
map[string]interface{}{"name": "Charlie", "age": float64(30)},
},
expectedCount: 2, // Alice and Charlie
description:   "Should filter using item.age > 18 (RECOMMENDED SYNTAX)",
},
{
name:      "item.field with variable comparison",
condition: "item.age >= variables.minAge",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "age": float64(25)},
map[string]interface{}{"name": "Bob", "age": float64(19)},
map[string]interface{}{"name": "Charlie", "age": float64(22)},
},
expectedCount: 2, // Alice and Charlie (>= 21)
description:   "Should filter using item.age >= variables.minAge",
},
{
name:      "item.nested.field syntax",
condition: "item.profile.verified == true",
inputArray: []interface{}{
map[string]interface{}{
"name": "Alice",
"profile": map[string]interface{}{"verified": true},
},
map[string]interface{}{
"name": "Bob",
"profile": map[string]interface{}{"verified": false},
},
map[string]interface{}{
"name": "Charlie",
"profile": map[string]interface{}{"verified": true},
},
},
expectedCount: 2, // Alice and Charlie
description:   "Should filter using nested item.profile.verified",
},
{
name:      "item syntax with complex AND condition",
condition: "item.age >= 18 && item.status == \"active\"",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "age": float64(25), "status": "active"},
map[string]interface{}{"name": "Bob", "age": float64(19), "status": "inactive"},
map[string]interface{}{"name": "Charlie", "age": float64(22), "status": "active"},
map[string]interface{}{"name": "Dave", "age": float64(17), "status": "active"},
},
expectedCount: 2, // Alice and Charlie
description:   "Should filter with item.age >= 18 AND item.status == active",
},
{
name:      "item syntax with context variable",
condition: "item.score > context.passingScore",
inputArray: []interface{}{
map[string]interface{}{"name": "Alice", "score": float64(85)},
map[string]interface{}{"name": "Bob", "score": float64(65)},
map[string]interface{}{"name": "Charlie", "score": float64(95)},
},
expectedCount: 2, // Alice and Charlie (> 70)
description:   "Should filter using item.score > context.passingScore",
},
{
name:      "item reference for primitive in array",
condition: "item > 10",
inputArray: []interface{}{
float64(5),
float64(15),
float64(8),
float64(20),
},
expectedCount: 2, // 15 and 20
description:   "Should filter primitive values using item > 10",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
exec := &FilterExecutor{}
ctx := &MockExecutionContext{
inputs: map[string][]interface{}{
"test-node": {tt.inputArray},
},
variables: map[string]interface{}{
"minAge": float64(21),
},
contextVars: map[string]interface{}{
"passingScore": float64(70),
},
}

node := types.Node{
ID:   "test-node",
Type: types.NodeTypeFilter,
Data: types.NodeData{
Condition: &tt.condition,
},
}

result, err := exec.Execute(ctx, node)
if err != nil {
t.Fatalf("Unexpected error: %v", err)
}

resultMap := result.(map[string]interface{})
filtered := resultMap["filtered"].([]interface{})

if len(filtered) != tt.expectedCount {
t.Errorf("Expected %d filtered items, got %d. Description: %s",
tt.expectedCount, len(filtered), tt.description)
}
})
}
}
