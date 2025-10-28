package workflow

import (
	"encoding/json"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "valid simple workflow",
			json: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "2", "data": {"value": 5}},
					{"id": "3", "data": {"op": "add"}}
				],
				"edges": [
					{"id": "e1-3", "source": "1", "target": "3"},
					{"id": "e2-3", "source": "2", "target": "3"}
				]
			}`,
			wantErr: false,
		},
		{
			name:    "empty workflow",
			json:    `{"nodes": [], "edges": []}`,
			wantErr: true,
		},
		{
			name: "invalid edge reference",
			json: `{
				"nodes": [{"id": "1", "data": {"value": 10}}],
				"edges": [{"id": "e1-2", "source": "1", "target": "999"}]
			}`,
			wantErr: true,
		},
		{
			name: "duplicate node ID",
			json: `{
				"nodes": [
					{"id": "1", "data": {"value": 10}},
					{"id": "1", "data": {"value": 20}}
				],
				"edges": []
			}`,
			wantErr: true,
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEngine_Execute_NumberNodes(t *testing.T) {
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Data: NodeData{Value: floatPtr(42.5)}},
		},
		Edges: []Edge{},
	}

	engine := NewEngine(workflow)
	result, err := engine.Execute()

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Results["1"].Value != 42.5 {
		t.Errorf("Expected value 42.5, got %v", result.Results["1"].Value)
	}
}

func TestEngine_Execute_Addition(t *testing.T) {
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Type: "numberNode", Data: NodeData{Value: floatPtr(10)}},
			{ID: "2", Type: "numberNode", Data: NodeData{Value: floatPtr(5)}},
			{ID: "3", Type: "opNode", Data: NodeData{Op: stringPtr("add")}},
		},
		Edges: []Edge{
			{ID: "e1-3", Source: "1", Target: "3"},
			{ID: "e2-3", Source: "2", Target: "3"},
		},
	}

	engine := NewEngine(workflow)
	result, err := engine.Execute()

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Results["3"].Value != 15.0 {
		t.Errorf("Expected 15, got %v", result.Results["3"].Value)
	}
}

func TestEngine_Execute_Subtraction(t *testing.T) {
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Type: "numberNode", Data: NodeData{Value: floatPtr(10)}},
			{ID: "2", Type: "numberNode", Data: NodeData{Value: floatPtr(3)}},
			{ID: "3", Type: "opNode", Data: NodeData{Op: stringPtr("subtract")}},
		},
		Edges: []Edge{
			{ID: "e1-3", Source: "1", Target: "3"},
			{ID: "e2-3", Source: "2", Target: "3"},
		},
	}

	engine := NewEngine(workflow)
	result, err := engine.Execute()

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Results["3"].Value != 7.0 {
		t.Errorf("Expected 7, got %v", result.Results["3"].Value)
	}
}

func TestEngine_Execute_Multiplication(t *testing.T) {
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Type: "numberNode", Data: NodeData{Value: floatPtr(4)}},
			{ID: "2", Type: "numberNode", Data: NodeData{Value: floatPtr(5)}},
			{ID: "3", Type: "opNode", Data: NodeData{Op: stringPtr("multiply")}},
		},
		Edges: []Edge{
			{ID: "e1-3", Source: "1", Target: "3"},
			{ID: "e2-3", Source: "2", Target: "3"},
		},
	}

	engine := NewEngine(workflow)
	result, err := engine.Execute()

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Results["3"].Value != 20.0 {
		t.Errorf("Expected 20, got %v", result.Results["3"].Value)
	}
}

func TestEngine_Execute_Division(t *testing.T) {
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Type: "numberNode", Data: NodeData{Value: floatPtr(20)}},
			{ID: "2", Type: "numberNode", Data: NodeData{Value: floatPtr(4)}},
			{ID: "3", Type: "opNode", Data: NodeData{Op: stringPtr("divide")}},
		},
		Edges: []Edge{
			{ID: "e1-3", Source: "1", Target: "3"},
			{ID: "e2-3", Source: "2", Target: "3"},
		},
	}

	engine := NewEngine(workflow)
	result, err := engine.Execute()

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Results["3"].Value != 5.0 {
		t.Errorf("Expected 5, got %v", result.Results["3"].Value)
	}
}

func TestEngine_Execute_DivisionByZero(t *testing.T) {
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Type: "numberNode", Data: NodeData{Value: floatPtr(20)}},
			{ID: "2", Type: "numberNode", Data: NodeData{Value: floatPtr(0)}},
			{ID: "3", Type: "opNode", Data: NodeData{Op: stringPtr("divide")}},
		},
		Edges: []Edge{
			{ID: "e1-3", Source: "1", Target: "3"},
			{ID: "e2-3", Source: "2", Target: "3"},
		},
	}

	engine := NewEngine(workflow)
	_, err := engine.Execute()

	if err == nil {
		t.Error("Expected error for division by zero")
	}
}

func TestEngine_Execute_CompleteWorkflow(t *testing.T) {
	// This mimics the frontend's default workflow
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Type: "numberNode", Data: NodeData{Value: floatPtr(10), Label: stringPtr("Node 1")}},
			{ID: "2", Type: "numberNode", Data: NodeData{Value: floatPtr(5), Label: stringPtr("Node 2")}},
			{ID: "3", Type: "opNode", Data: NodeData{Op: stringPtr("add"), Label: stringPtr("Node 3 (op)")}},
			{ID: "4", Type: "vizNode", Data: NodeData{Mode: stringPtr("text"), Label: stringPtr("Node 4 (viz)")}},
		},
		Edges: []Edge{
			{ID: "e1-3", Source: "1", Target: "3"},
			{ID: "e2-3", Source: "2", Target: "3"},
			{ID: "e3-4", Source: "3", Target: "4"},
		},
	}

	engine := NewEngine(workflow)
	result, err := engine.Execute()

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Check intermediate results
	if result.Results["1"].Value != 10.0 {
		t.Errorf("Node 1: expected 10, got %v", result.Results["1"].Value)
	}
	if result.Results["2"].Value != 5.0 {
		t.Errorf("Node 2: expected 5, got %v", result.Results["2"].Value)
	}
	if result.Results["3"].Value != 15.0 {
		t.Errorf("Node 3: expected 15, got %v", result.Results["3"].Value)
	}

	// Check final output exists
	if result.Output == nil {
		t.Error("Expected output from visualization node")
	}
}

func TestEngine_Execute_VisualizationTable(t *testing.T) {
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Type: "numberNode", Data: NodeData{Value: floatPtr(100)}},
			{ID: "2", Type: "vizNode", Data: NodeData{Mode: stringPtr("table")}},
		},
		Edges: []Edge{
			{ID: "e1-2", Source: "1", Target: "2"},
		},
	}

	engine := NewEngine(workflow)
	result, err := engine.Execute()

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Check that table output is structured correctly
	if output, ok := result.Output.(map[string]interface{}); ok {
		if output["type"] != "table" {
			t.Errorf("Expected type 'table', got %v", output["type"])
		}
	} else {
		t.Error("Expected output to be a map for table mode")
	}
}

func TestEngine_Execute_CircularDependency(t *testing.T) {
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Type: "opNode", Data: NodeData{Op: stringPtr("add")}},
			{ID: "2", Type: "opNode", Data: NodeData{Op: stringPtr("add")}},
		},
		Edges: []Edge{
			{ID: "e1-2", Source: "1", Target: "2"},
			{ID: "e2-1", Source: "2", Target: "1"},
		},
	}

	engine := NewEngine(workflow)
	_, err := engine.Execute()

	if err == nil {
		t.Error("Expected error for circular dependency")
	}
}

func TestExecuteWorkflow_Integration(t *testing.T) {
	jsonData := []byte(`{
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
	}`)

	result, err := ExecuteWorkflow(jsonData)
	if err != nil {
		t.Fatalf("ExecuteWorkflow() error = %v", err)
	}

	if result.Results["3"].Value != 15.0 {
		t.Errorf("Expected operation result 15, got %v", result.Results["3"].Value)
	}

	if result.Output == nil {
		t.Error("Expected visualization output")
	}
}

func TestExecuteWorkflow_InvalidJSON(t *testing.T) {
	jsonData := []byte(`invalid json`)

	_, err := ExecuteWorkflow(jsonData)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestEngine_TopologicalSort(t *testing.T) {
	workflow := &Workflow{
		Nodes: []Node{
			{ID: "1", Type: "numberNode", Data: NodeData{Value: floatPtr(10)}},
			{ID: "2", Type: "numberNode", Data: NodeData{Value: floatPtr(5)}},
			{ID: "3", Type: "opNode", Data: NodeData{Op: stringPtr("add")}},
			{ID: "4", Type: "opNode", Data: NodeData{Op: stringPtr("multiply")}},
			{ID: "5", Type: "vizNode", Data: NodeData{Mode: stringPtr("text")}},
		},
		Edges: []Edge{
			{ID: "e1-3", Source: "1", Target: "3"},
			{ID: "e2-3", Source: "2", Target: "3"},
			{ID: "e3-4", Source: "3", Target: "4"},
			{ID: "e2-4", Source: "2", Target: "4"},
			{ID: "e4-5", Source: "4", Target: "5"},
		},
	}

	engine := NewEngine(workflow)
	deps := engine.buildDependencyGraph()
	sorted, err := engine.topologicalSort(deps)

	if err != nil {
		t.Fatalf("topologicalSort() error = %v", err)
	}

	// Verify that dependencies come before dependents
	positions := make(map[string]int)
	for i, nodeID := range sorted {
		positions[nodeID] = i
	}

	for _, edge := range workflow.Edges {
		if positions[edge.Source] >= positions[edge.Target] {
			t.Errorf("Source node %s should come before target node %s", edge.Source, edge.Target)
		}
	}
}

// Helper functions
func floatPtr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

// Benchmark tests
func BenchmarkExecuteWorkflow(b *testing.B) {
	jsonData := []byte(`{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 5}},
			{"id": "3", "data": {"op": "add"}},
			{"id": "4", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1-3", "source": "1", "target": "3"},
			{"id": "e2-3", "source": "2", "target": "3"},
			{"id": "e3-4", "source": "3", "target": "4"}
		]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ExecuteWorkflow(jsonData)
	}
}

func TestNodeData_JSONMarshaling(t *testing.T) {
	node := Node{
		ID: "1",
		Data: NodeData{
			Value: floatPtr(42.5),
			Label: stringPtr("Test Node"),
		},
		Type: "numberNode",
	}

	// Marshal
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Unmarshal
	var decoded Node
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if *decoded.Data.Value != 42.5 {
		t.Errorf("Expected value 42.5, got %v", *decoded.Data.Value)
	}
}
