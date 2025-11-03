package graph

import (
	"sort"
	"strings"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestTopologicalSort_Simple tests basic topological sorting
func TestTopologicalSort_Simple(t *testing.T) {
	tests := []struct {
		name       string
		nodes      []types.Node
		edges      []types.Edge
		wantOrder  []string
		wantErr    bool
		checkOrder bool // if false, just check success/failure
	}{
		{
			name: "linear chain",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeOperation},
				{ID: "3", Type: types.NodeTypeOperation},
			},
			edges: []types.Edge{
				{Source: "1", Target: "2"},
				{Source: "2", Target: "3"},
			},
			wantOrder: []string{"1", "2", "3"},
		},
		{
			name: "diamond shape",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeOperation},
				{ID: "3", Type: types.NodeTypeOperation},
				{ID: "4", Type: types.NodeTypeOperation},
			},
			edges: []types.Edge{
				{Source: "1", Target: "2"},
				{Source: "1", Target: "3"},
				{Source: "2", Target: "4"},
				{Source: "3", Target: "4"},
			},
			// Multiple valid orders exist, just verify 1 before 2,3 and 2,3 before 4
			checkOrder: false,
		},
		{
			name: "single node",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
			},
			edges:     []types.Edge{},
			wantOrder: []string{"1"},
		},
		{
			name: "multiple roots",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeNumber},
				{ID: "3", Type: types.NodeTypeOperation},
			},
			edges: []types.Edge{
				{Source: "1", Target: "3"},
				{Source: "2", Target: "3"},
			},
			// 1 and 2 can be in any order, but must come before 3
			checkOrder: false,
		},
		{
			name:      "empty graph",
			nodes:     []types.Node{},
			edges:     []types.Edge{},
			wantOrder: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(tt.nodes, tt.edges)
			got, err := g.TopologicalSort()

			if (err != nil) != tt.wantErr {
				t.Errorf("TopologicalSort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if tt.checkOrder {
				if !equalSlices(got, tt.wantOrder) {
					t.Errorf("TopologicalSort() = %v, want %v", got, tt.wantOrder)
				}
			} else {
				// Verify it's a valid topological order
				if !isValidTopologicalOrder(got, tt.edges) {
					t.Errorf("TopologicalSort() returned invalid order: %v", got)
				}
			}
		})
	}
}

// TestTopologicalSort_Cycles tests cycle detection
func TestTopologicalSort_Cycles(t *testing.T) {
	tests := []struct {
		name  string
		nodes []types.Node
		edges []types.Edge
	}{
		{
			name: "simple cycle",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeOperation},
			},
			edges: []types.Edge{
				{Source: "1", Target: "2"},
				{Source: "2", Target: "1"},
			},
		},
		{
			name: "self loop",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
			},
			edges: []types.Edge{
				{Source: "1", Target: "1"},
			},
		},
		{
			name: "three node cycle",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeOperation},
				{ID: "3", Type: types.NodeTypeOperation},
			},
			edges: []types.Edge{
				{Source: "1", Target: "2"},
				{Source: "2", Target: "3"},
				{Source: "3", Target: "1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(tt.nodes, tt.edges)
			_, err := g.TopologicalSort()

			if err == nil {
				t.Error("TopologicalSort() expected error for cyclic graph, got nil")
			}
		})
	}
}

// TestTopologicalSort_Large tests performance with larger graphs
func TestTopologicalSort_Large(t *testing.T) {
	tests := []struct {
		name     string
		numNodes int
	}{
		{name: "100 nodes linear", numNodes: 100},
		{name: "1000 nodes linear", numNodes: 1000},
		{name: "100 nodes wide", numNodes: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nodes []types.Node
			var edges []types.Edge

			// Generate different graph types based on test name
			if strings.Contains(tt.name, "linear") {
				nodes, edges = generateLinearChain(tt.numNodes)
			} else if strings.Contains(tt.name, "wide") {
				nodes, edges = generateWideGraph(tt.numNodes)
			}

			g := New(nodes, edges)

			order, err := g.TopologicalSort()
			if err != nil {
				t.Errorf("TopologicalSort() unexpected error: %v", err)
				return
			}

			if len(order) != len(nodes) {
				t.Errorf("TopologicalSort() returned %d nodes, want %d", len(order), len(nodes))
			}

			if !isValidTopologicalOrder(order, edges) {
				t.Error("TopologicalSort() returned invalid order")
			}
		})
	}
}

// TestDetectCycles tests the cycle detection method
func TestDetectCycles(t *testing.T) {
	tests := []struct {
		name    string
		nodes   []types.Node
		edges   []types.Edge
		wantErr bool
	}{
		{
			name: "no cycle",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeOperation},
			},
			edges: []types.Edge{
				{Source: "1", Target: "2"},
			},
			wantErr: false,
		},
		{
			name: "cycle exists",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeOperation},
			},
			edges: []types.Edge{
				{Source: "1", Target: "2"},
				{Source: "2", Target: "1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(tt.nodes, tt.edges)
			err := g.DetectCycles()

			if (err != nil) != tt.wantErr {
				t.Errorf("DetectCycles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetNode tests node retrieval
func TestGetNode(t *testing.T) {
	nodes := []types.Node{
		{ID: "1", Type: types.NodeTypeNumber},
		{ID: "2", Type: types.NodeTypeOperation},
	}
	g := New(nodes, nil)

	tests := []struct {
		name   string
		nodeID string
		want   *types.Node
	}{
		{name: "existing node", nodeID: "1", want: &nodes[0]},
		{name: "another existing node", nodeID: "2", want: &nodes[1]},
		{name: "non-existing node", nodeID: "3", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.GetNode(tt.nodeID)
			if got == nil && tt.want == nil {
				return
			}
			if got == nil || tt.want == nil {
				t.Errorf("GetNode() = %v, want %v", got, tt.want)
				return
			}
			if got.ID != tt.want.ID || got.Type != tt.want.Type {
				t.Errorf("GetNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetNodeInputEdges tests retrieving input edges
func TestGetNodeInputEdges(t *testing.T) {
	edges := []types.Edge{
		{Source: "1", Target: "2"},
		{Source: "3", Target: "2"},
		{Source: "2", Target: "4"},
	}
	g := New(nil, edges)

	tests := []struct {
		name      string
		nodeID    string
		wantCount int
	}{
		{name: "node with 2 inputs", nodeID: "2", wantCount: 2},
		{name: "node with 1 input", nodeID: "4", wantCount: 1},
		{name: "node with no inputs", nodeID: "1", wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.GetNodeInputEdges(tt.nodeID)
			if len(got) != tt.wantCount {
				t.Errorf("GetNodeInputEdges() returned %d edges, want %d", len(got), tt.wantCount)
			}
		})
	}
}

// TestGetNodeOutputEdges tests retrieving output edges
func TestGetNodeOutputEdges(t *testing.T) {
	edges := []types.Edge{
		{Source: "1", Target: "2"},
		{Source: "1", Target: "3"},
		{Source: "2", Target: "4"},
	}
	g := New(nil, edges)

	tests := []struct {
		name      string
		nodeID    string
		wantCount int
	}{
		{name: "node with 2 outputs", nodeID: "1", wantCount: 2},
		{name: "node with 1 output", nodeID: "2", wantCount: 1},
		{name: "node with no outputs", nodeID: "4", wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.GetNodeOutputEdges(tt.nodeID)
			if len(got) != tt.wantCount {
				t.Errorf("GetNodeOutputEdges() returned %d edges, want %d", len(got), tt.wantCount)
			}
		})
	}
}

// TestGetTerminalNodes tests finding terminal nodes
func TestGetTerminalNodes(t *testing.T) {
	tests := []struct {
		name  string
		nodes []types.Node
		edges []types.Edge
		want  []string
	}{
		{
			name: "single terminal",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeOperation},
			},
			edges: []types.Edge{
				{Source: "1", Target: "2"},
			},
			want: []string{"2"},
		},
		{
			name: "multiple terminals",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeOperation},
				{ID: "3", Type: types.NodeTypeOperation},
			},
			edges: []types.Edge{
				{Source: "1", Target: "2"},
				{Source: "1", Target: "3"},
			},
			want: []string{"2", "3"},
		},
		{
			name: "all nodes terminal",
			nodes: []types.Node{
				{ID: "1", Type: types.NodeTypeNumber},
				{ID: "2", Type: types.NodeTypeNumber},
			},
			edges: []types.Edge{},
			want:  []string{"1", "2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(tt.nodes, tt.edges)
			got := g.GetTerminalNodes()

			sort.Strings(got)
			sort.Strings(tt.want)

			if !equalSlices(got, tt.want) {
				t.Errorf("GetTerminalNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isValidTopologicalOrder(order []string, edges []types.Edge) bool {
	// Create position map
	pos := make(map[string]int)
	for i, nodeID := range order {
		pos[nodeID] = i
	}

	// Check all edges respect the order
	for _, edge := range edges {
		sourcePos, sourceExists := pos[edge.Source]
		targetPos, targetExists := pos[edge.Target]

		if !sourceExists || !targetExists {
			return false
		}

		// Source must come before target
		if sourcePos >= targetPos {
			return false
		}
	}

	return true
}
