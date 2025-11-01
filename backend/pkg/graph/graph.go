// Package graph provides DAG (Directed Acyclic Graph) operations for workflow execution.
// This includes topological sorting, cycle detection, and graph traversal utilities.
package graph

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Graph represents a workflow graph with nodes and edges
type Graph struct {
	nodes []types.Node
	edges []types.Edge
}

// New creates a new Graph from nodes and edges
func New(nodes []types.Node, edges []types.Edge) *Graph {
	return &Graph{
		nodes: nodes,
		edges: edges,
	}
}

// TopologicalSort performs topological sorting on the workflow graph using Kahn's algorithm.
// This determines the correct execution order for nodes in a directed acyclic graph (DAG).
//
// Returns:
//   - []string: Ordered list of node IDs for sequential execution
//   - error: If the workflow contains cycles (circular dependencies)
//
// Algorithm:
//  1. Calculate in-degree (number of incoming edges) for each node
//  2. Start with nodes that have no dependencies (in-degree = 0)
//  3. Process nodes and reduce in-degree of their neighbors
//  4. If all nodes processed, we have a valid execution order
//  5. If nodes remain, there's a cycle in the graph
//
// Optimizations:
//  - Pre-allocated slices with exact capacity to minimize allocations
//  - Ring buffer for queue to avoid expensive slice operations
//  - Insertion sort for small orphan node sets (faster than generic sort for small n)
//  - Single pass edge processing to build both adjacency list and in-degree
func (g *Graph) TopologicalSort() ([]string, error) {
	numNodes := len(g.nodes)
	
	// Early return for empty graph
	if numNodes == 0 {
		return []string{}, nil
	}
	
	// Pre-allocate with exact capacity to avoid reallocation
	inDegree := make(map[string]int, numNodes)
	adjacency := make(map[string][]string, numNodes)
	
	// Initialize in-degree for all nodes to zero
	for i := range g.nodes {
		inDegree[g.nodes[i].ID] = 0
	}
	
	// Build the graph structure in a single pass
	for i := range g.edges {
		edge := &g.edges[i]
		adjacency[edge.Source] = append(adjacency[edge.Source], edge.Target)
		inDegree[edge.Target]++
	}
	
	// Find all nodes with no dependencies (in-degree = 0)
	// Pre-allocate with capacity to avoid growing
	orphanNodes := make([]string, 0, numNodes)
	for nodeID, degree := range inDegree {
		if degree == 0 {
			orphanNodes = append(orphanNodes, nodeID)
		}
	}
	
	// Sort orphan nodes by ID to ensure deterministic execution order
	// Use insertion sort for small arrays (typically faster than quicksort for n < 20)
	// This is important for context nodes that need to execute before other nodes
	insertionSort(orphanNodes)
	
	// Use a ring buffer for the queue to avoid expensive slice operations
	// Pre-allocate with capacity for all nodes
	queue := make([]string, numNodes)
	queueStart := 0
	queueEnd := len(orphanNodes)
	copy(queue, orphanNodes)
	
	// Pre-allocate result with exact capacity
	order := make([]string, 0, numNodes)
	
	// Process nodes in topological order
	for queueStart < queueEnd {
		// Dequeue using ring buffer (O(1) instead of O(n))
		current := queue[queueStart]
		queueStart++
		order = append(order, current)
		
		// Reduce in-degree for all neighbors
		neighbors := adjacency[current]
		for i := range neighbors {
			neighbor := neighbors[i]
			inDegree[neighbor]--
			// If neighbor has no more dependencies, add to queue
			if inDegree[neighbor] == 0 {
				queue[queueEnd] = neighbor
				queueEnd++
			}
		}
	}
	
	// Check if all nodes were processed
	// If not, there's a cycle in the graph
	if len(order) != numNodes {
		return nil, fmt.Errorf("workflow contains cycles (circular dependencies)")
	}
	
	return order, nil
}

// insertionSort sorts a slice of strings in place using insertion sort.
// This is faster than the standard library sort for small slices (n < ~20).
func insertionSort(arr []string) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

// GetNode retrieves a node by its ID
func (g *Graph) GetNode(nodeID string) *types.Node {
	for i := range g.nodes {
		if g.nodes[i].ID == nodeID {
			return &g.nodes[i]
		}
	}
	return nil
}

// GetNodeInputEdges returns all edges where the given node is the target
func (g *Graph) GetNodeInputEdges(nodeID string) []types.Edge {
	var edges []types.Edge
	for _, edge := range g.edges {
		if edge.Target == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetNodeOutputEdges returns all edges where the given node is the source
func (g *Graph) GetNodeOutputEdges(nodeID string) []types.Edge {
	var edges []types.Edge
	for _, edge := range g.edges {
		if edge.Source == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetTerminalNodes returns all nodes that have no outgoing edges
func (g *Graph) GetTerminalNodes() []string {
	// Build a set of all terminal nodes
	terminalNodes := make(map[string]bool)

	// Initially, all nodes are considered terminal
	for _, node := range g.nodes {
		terminalNodes[node.ID] = true
	}

	// Remove nodes that have outgoing edges
	for _, edge := range g.edges {
		terminalNodes[edge.Source] = false
	}

	// Convert to slice
	result := []string{}
	for nodeID, isTerminal := range terminalNodes {
		if isTerminal {
			result = append(result, nodeID)
		}
	}

	return result
}

// DetectCycles detects if the graph contains any cycles
func (g *Graph) DetectCycles() error {
	_, err := g.TopologicalSort()
	return err
}
