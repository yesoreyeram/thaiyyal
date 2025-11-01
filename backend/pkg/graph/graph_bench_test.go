package graph

import (
	"fmt"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Benchmark topological sort with different graph sizes and structures

// BenchmarkTopologicalSort_Linear benchmarks linear chains
func BenchmarkTopologicalSort_Linear(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d_nodes", size), func(b *testing.B) {
			nodes, edges := generateLinearChain(size)
			g := New(nodes, edges)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				_, err := g.TopologicalSort()
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkTopologicalSort_Wide benchmarks wide graphs (many parallel branches)
func BenchmarkTopologicalSort_Wide(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d_nodes", size), func(b *testing.B) {
			nodes, edges := generateWideGraph(size)
			g := New(nodes, edges)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				_, err := g.TopologicalSort()
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkTopologicalSort_Dense benchmarks dense graphs
func BenchmarkTopologicalSort_Dense(b *testing.B) {
	sizes := []int{10, 50, 100, 500}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d_nodes", size), func(b *testing.B) {
			nodes, edges := generateDenseDAG(size)
			g := New(nodes, edges)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				_, err := g.TopologicalSort()
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkTopologicalSort_Tree benchmarks tree structures
func BenchmarkTopologicalSort_Tree(b *testing.B) {
	sizes := []int{15, 31, 63, 127, 255, 511, 1023} // Binary tree sizes: 2^n - 1
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d_nodes", size), func(b *testing.B) {
			nodes, edges := generateBinaryTree(size)
			g := New(nodes, edges)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				_, err := g.TopologicalSort()
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkTopologicalSort_Diamond benchmarks diamond-shaped graphs
func BenchmarkTopologicalSort_Diamond(b *testing.B) {
	sizes := []int{10, 50, 100, 500}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d_layers", size), func(b *testing.B) {
			nodes, edges := generateDiamondGraph(size)
			g := New(nodes, edges)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				_, err := g.TopologicalSort()
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkTopologicalSort_RealWorld benchmarks realistic workflow patterns
func BenchmarkTopologicalSort_RealWorld(b *testing.B) {
	scenarios := []struct {
		name  string
		nodes []types.Node
		edges []types.Edge
	}{
		{
			name:  "simple_pipeline",
			nodes: generatePipelineNodes(20, 5), // 20 stages, 5 parallel per stage
			edges: generatePipelineEdges(20, 5),
		},
		{
			name:  "complex_pipeline",
			nodes: generatePipelineNodes(50, 10), // 50 stages, 10 parallel per stage
			edges: generatePipelineEdges(50, 10),
		},
		{
			name:  "fan_out_fan_in",
			nodes: generateFanOutFanInNodes(100),
			edges: generateFanOutFanInEdges(100),
		},
	}
	
	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			g := New(scenario.nodes, scenario.edges)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				_, err := g.TopologicalSort()
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkNew tests graph creation performance
func BenchmarkNew(b *testing.B) {
	nodes, edges := generateLinearChain(1000)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_ = New(nodes, edges)
	}
}

// Helper functions to generate test graphs

func generateLinearChain(size int) ([]types.Node, []types.Edge) {
	nodes := make([]types.Node, size)
	edges := make([]types.Edge, size-1)
	
	for i := 0; i < size; i++ {
		nodes[i] = types.Node{
			ID:   fmt.Sprintf("node-%d", i),
			Type: types.NodeTypeOperation,
		}
	}
	
	for i := 0; i < size-1; i++ {
		edges[i] = types.Edge{
			Source: nodes[i].ID,
			Target: nodes[i+1].ID,
		}
	}
	
	return nodes, edges
}

func generateWideGraph(size int) ([]types.Node, []types.Edge) {
	// Create a graph with one root, many parallel branches, and one sink
	nodes := make([]types.Node, size+2) // +2 for root and sink
	edges := make([]types.Edge, 0, size*2)
	
	nodes[0] = types.Node{ID: "root", Type: types.NodeTypeOperation}
	nodes[size+1] = types.Node{ID: "sink", Type: types.NodeTypeOperation}
	
	for i := 0; i < size; i++ {
		nodes[i+1] = types.Node{
			ID:   fmt.Sprintf("node-%d", i),
			Type: types.NodeTypeOperation,
		}
		edges = append(edges, types.Edge{Source: "root", Target: nodes[i+1].ID})
		edges = append(edges, types.Edge{Source: nodes[i+1].ID, Target: "sink"})
	}
	
	return nodes, edges
}

func generateDenseDAG(size int) ([]types.Node, []types.Edge) {
	nodes := make([]types.Node, size)
	edges := make([]types.Edge, 0)
	
	for i := 0; i < size; i++ {
		nodes[i] = types.Node{
			ID:   fmt.Sprintf("node-%d", i),
			Type: types.NodeTypeOperation,
		}
	}
	
	// Add edges from each node to several later nodes
	for i := 0; i < size; i++ {
		// Connect to next 3 nodes (or fewer if near the end)
		for j := 1; j <= 3 && i+j < size; j++ {
			edges = append(edges, types.Edge{
				Source: nodes[i].ID,
				Target: nodes[i+j].ID,
			})
		}
	}
	
	return nodes, edges
}

func generateBinaryTree(size int) ([]types.Node, []types.Edge) {
	nodes := make([]types.Node, size)
	edges := make([]types.Edge, 0, size-1)
	
	for i := 0; i < size; i++ {
		nodes[i] = types.Node{
			ID:   fmt.Sprintf("node-%d", i),
			Type: types.NodeTypeOperation,
		}
	}
	
	// Binary tree: node i has children at 2i+1 and 2i+2
	for i := 0; i < size; i++ {
		left := 2*i + 1
		right := 2*i + 2
		
		if left < size {
			edges = append(edges, types.Edge{
				Source: nodes[i].ID,
				Target: nodes[left].ID,
			})
		}
		if right < size {
			edges = append(edges, types.Edge{
				Source: nodes[i].ID,
				Target: nodes[right].ID,
			})
		}
	}
	
	return nodes, edges
}

func generateDiamondGraph(layers int) ([]types.Node, []types.Edge) {
	// Each layer has 2 nodes, creating diamond patterns
	numNodes := layers * 2
	nodes := make([]types.Node, numNodes)
	edges := make([]types.Edge, 0)
	
	for i := 0; i < numNodes; i++ {
		nodes[i] = types.Node{
			ID:   fmt.Sprintf("node-%d", i),
			Type: types.NodeTypeOperation,
		}
	}
	
	for layer := 0; layer < layers-1; layer++ {
		// Connect both nodes in current layer to both nodes in next layer
		curr1 := layer * 2
		curr2 := layer*2 + 1
		next1 := (layer + 1) * 2
		next2 := (layer+1)*2 + 1
		
		edges = append(edges, 
			types.Edge{Source: nodes[curr1].ID, Target: nodes[next1].ID},
			types.Edge{Source: nodes[curr1].ID, Target: nodes[next2].ID},
			types.Edge{Source: nodes[curr2].ID, Target: nodes[next1].ID},
			types.Edge{Source: nodes[curr2].ID, Target: nodes[next2].ID},
		)
	}
	
	return nodes, edges
}

func generatePipelineNodes(stages, parallelPerStage int) []types.Node {
	nodes := make([]types.Node, stages*parallelPerStage)
	
	for i := 0; i < stages; i++ {
		for j := 0; j < parallelPerStage; j++ {
			idx := i*parallelPerStage + j
			nodes[idx] = types.Node{
				ID:   fmt.Sprintf("stage-%d-node-%d", i, j),
				Type: types.NodeTypeOperation,
			}
		}
	}
	
	return nodes
}

func generatePipelineEdges(stages, parallelPerStage int) []types.Edge {
	edges := make([]types.Edge, 0)
	
	for i := 0; i < stages-1; i++ {
		// Connect each node in current stage to all nodes in next stage
		for j := 0; j < parallelPerStage; j++ {
			for k := 0; k < parallelPerStage; k++ {
				sourceIdx := i*parallelPerStage + j
				targetIdx := (i+1)*parallelPerStage + k
				edges = append(edges, types.Edge{
					Source: fmt.Sprintf("stage-%d-node-%d", i, j),
					Target: fmt.Sprintf("stage-%d-node-%d", i+1, k),
				})
				_, _ = sourceIdx, targetIdx // avoid unused variable
			}
		}
	}
	
	return edges
}

func generateFanOutFanInNodes(branchCount int) []types.Node {
	nodes := make([]types.Node, branchCount+2) // +2 for root and sink
	
	nodes[0] = types.Node{ID: "root", Type: types.NodeTypeOperation}
	nodes[branchCount+1] = types.Node{ID: "sink", Type: types.NodeTypeOperation}
	
	for i := 0; i < branchCount; i++ {
		nodes[i+1] = types.Node{
			ID:   fmt.Sprintf("branch-%d", i),
			Type: types.NodeTypeOperation,
		}
	}
	
	return nodes
}

func generateFanOutFanInEdges(branchCount int) []types.Edge {
	edges := make([]types.Edge, 0, branchCount*2)
	
	for i := 0; i < branchCount; i++ {
		// Fan out from root
		edges = append(edges, types.Edge{
			Source: "root",
			Target: fmt.Sprintf("branch-%d", i),
		})
		// Fan in to sink
		edges = append(edges, types.Edge{
			Source: fmt.Sprintf("branch-%d", i),
			Target: "sink",
		})
	}
	
	return edges
}
