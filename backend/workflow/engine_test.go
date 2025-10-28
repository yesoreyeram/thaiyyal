package workflow

import (
	"encoding/json"
	"testing"
)

func TestNewEngine(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "valid payload",
			payload: `{"nodes":[{"id":"1","data":{"value":10}}],"edges":[]}`,
			wantErr: false,
		},
		{
			name:    "invalid json",
			payload: `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEngine([]byte(tt.payload))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEngine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEngine_Execute_SimpleAddition(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "data": {"op": "add"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(result.Errors) > 0 {
		t.Errorf("Execute() errors = %v", result.Errors)
	}

	// Check node 3 result (addition)
	node3Result, ok := result.NodeResults["3"].(float64)
	if !ok {
		t.Fatalf("node 3 result is not a float64: %v", result.NodeResults["3"])
	}

	expected := 15.0
	if node3Result != expected {
		t.Errorf("node 3 result = %v, want %v", node3Result, expected)
	}

	// Check final output
	finalOutput, ok := result.FinalOutput.(float64)
	if !ok {
		t.Fatalf("final output is not a float64: %v", result.FinalOutput)
	}

	if finalOutput != expected {
		t.Errorf("final output = %v, want %v", finalOutput, expected)
	}
}

func TestEngine_Execute_CompleteWorkflow(t *testing.T) {
	// Test the complete workflow from the frontend:
	// Two numbers -> operation -> visualization
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10, "label": "Node 1"}},
			{"id": "2", "data": {"value": 5, "label": "Node 2"}},
			{"id": "3", "data": {"op": "add", "label": "Node 3 (op)"}},
			{"id": "4", "data": {"mode": "text", "label": "Node 4 (viz)"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-4", "source": "3", "target": "4"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(result.Errors) > 0 {
		t.Errorf("Execute() errors = %v", result.Errors)
	}

	// Check visualization node result
	vizResult, ok := result.FinalOutput.(map[string]interface{})
	if !ok {
		t.Fatalf("visualization result is not a map: %v", result.FinalOutput)
	}

	if vizResult["mode"] != "text" {
		t.Errorf("visualization mode = %v, want text", vizResult["mode"])
	}

	vizValue, ok := vizResult["value"].(float64)
	if !ok {
		t.Fatalf("visualization value is not a float64: %v", vizResult["value"])
	}

	if vizValue != 15.0 {
		t.Errorf("visualization value = %v, want 15", vizValue)
	}
}

func TestEngine_Execute_AllOperations(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		left     float64
		right    float64
		expected float64
		wantErr  bool
	}{
		{"addition", "add", 10, 5, 15, false},
		{"subtraction", "subtract", 10, 5, 5, false},
		{"multiplication", "multiply", 10, 5, 50, false},
		{"division", "divide", 10, 5, 2, false},
		{"division by zero", "divide", 10, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := Payload{
				Nodes: []Node{
					{ID: "1", Data: NodeData{Value: &tt.left}},
					{ID: "2", Data: NodeData{Value: &tt.right}},
					{ID: "3", Data: NodeData{Op: &tt.op}},
				},
				Edges: []Edge{
					{ID: "e1-3", Source: "1", Target: "3"},
					{ID: "e2-3", Source: "2", Target: "3"},
				},
			}

			payloadJSON, _ := json.Marshal(payload)
			engine, err := NewEngine(payloadJSON)
			if err != nil {
				t.Fatalf("NewEngine() error = %v", err)
			}

			result, err := engine.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				node3Result, ok := result.NodeResults["3"].(float64)
				if !ok {
					t.Fatalf("node 3 result is not a float64: %v", result.NodeResults["3"])
				}

				if node3Result != tt.expected {
					t.Errorf("result = %v, want %v", node3Result, tt.expected)
				}
			}
		})
	}
}

func TestEngine_Execute_MultipleOperations(t *testing.T) {
	// Test: (10 + 5) * 2 = 30
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "data": {"op": "add"}},
			{"id": "4", "data": {"value": 2}},
			{"id": "5", "data": {"op": "multiply"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-5", "source": "3", "target": "5"},
			{"id": "e4-5", "source": "4", "target": "5"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	node5Result, ok := result.NodeResults["5"].(float64)
	if !ok {
		t.Fatalf("node 5 result is not a float64: %v", result.NodeResults["5"])
	}

	expected := 30.0
	if node5Result != expected {
		t.Errorf("result = %v, want %v", node5Result, expected)
	}
}

func TestEngine_Execute_CyclicGraph(t *testing.T) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"op": "add"}},
			{"id": "3", "data": {"op": "multiply"}}
		],
		"edges": [
			{"id": "e1-2", "source": "1", "target": "2"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-2", "source": "3", "target": "2"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	_, err = engine.Execute()
	if err == nil {
		t.Error("Execute() expected error for cyclic graph, got nil")
	}
}

func TestEngine_Execute_MissingInputs(t *testing.T) {
	// Operation node with only one input
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"op": "add"}}
		],
		"edges": [
			{"id": "e1-2", "source": "1", "target": "2"}
		]
	}`

	engine, err := NewEngine([]byte(payload))
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	_, err = engine.Execute()
	if err == nil {
		t.Error("Execute() expected error for operation with insufficient inputs, got nil")
	}
}

func TestEngine_Execute_InvalidOperation(t *testing.T) {
	invalidOp := "invalid"
	payload := Payload{
		Nodes: []Node{
			{ID: "1", Data: NodeData{Value: floatPtr(10)}},
			{ID: "2", Data: NodeData{Value: floatPtr(5)}},
			{ID: "3", Data: NodeData{Op: &invalidOp}},
		},
		Edges: []Edge{
			{ID: "e1-3", Source: "1", Target: "3"},
			{ID: "e2-3", Source: "2", Target: "3"},
		},
	}

	payloadJSON, _ := json.Marshal(payload)
	engine, err := NewEngine(payloadJSON)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	_, err = engine.Execute()
	if err == nil {
		t.Error("Execute() expected error for invalid operation, got nil")
	}
}

func TestEngine_Execute_VisualizationModes(t *testing.T) {
	modes := []string{"text", "table"}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			payload := Payload{
				Nodes: []Node{
					{ID: "1", Data: NodeData{Value: floatPtr(42)}},
					{ID: "2", Data: NodeData{Mode: &mode}},
				},
				Edges: []Edge{
					{ID: "e1-2", Source: "1", Target: "2"},
				},
			}

			payloadJSON, _ := json.Marshal(payload)
			engine, err := NewEngine(payloadJSON)
			if err != nil {
				t.Fatalf("NewEngine() error = %v", err)
			}

			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			vizResult, ok := result.FinalOutput.(map[string]interface{})
			if !ok {
				t.Fatalf("visualization result is not a map: %v", result.FinalOutput)
			}

			if vizResult["mode"] != mode {
				t.Errorf("visualization mode = %v, want %v", vizResult["mode"], mode)
			}

			vizValue, ok := vizResult["value"].(float64)
			if !ok {
				t.Fatalf("visualization value is not a float64: %v", vizResult["value"])
			}

			if vizValue != 42.0 {
				t.Errorf("visualization value = %v, want 42", vizValue)
			}
		})
	}
}

// Helper function to create float64 pointers
func floatPtr(f float64) *float64 {
	return &f
}
