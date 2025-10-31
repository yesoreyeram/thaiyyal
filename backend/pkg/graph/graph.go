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
func (g *Graph) TopologicalSort() ([]string, error) {
	// Build in-degree map and adjacency list
	inDegree := make(map[string]int)
	adjacency := make(map[string][]string)

	// Initialize in-degree for all nodes to zero
	for _, node := range g.nodes {
		inDegree[node.ID] = 0
	}

	// Build the graph structure
	for _, edge := range g.edges {
		adjacency[edge.Source] = append(adjacency[edge.Source], edge.Target)
		inDegree[edge.Target]++
	}

	// Find all nodes with no dependencies (in-degree = 0)
	// These are the starting points for execution
	// Process in deterministic order (sorted by node ID) for reproducibility
	queue := []string{}
	orphanNodes := []string{}
	for nodeID, degree := range inDegree {
		if degree == 0 {
			orphanNodes = append(orphanNodes, nodeID)
		}
	}

	// Sort orphan nodes by ID to ensure deterministic execution order
	// This is important for context nodes that need to execute before other nodes
	for i := 0; i < len(orphanNodes); i++ {
		for j := i + 1; j < len(orphanNodes); j++ {
			if orphanNodes[i] > orphanNodes[j] {
				orphanNodes[i], orphanNodes[j] = orphanNodes[j], orphanNodes[i]
			}
		}
	}
	queue = append(queue, orphanNodes...)

	// Process nodes in topological order
	order := []string{}
	for len(queue) > 0 {
		// Dequeue the first node
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)

		// Reduce in-degree for all neighbors
		for _, neighbor := range adjacency[current] {
			inDegree[neighbor]--
			// If neighbor has no more dependencies, add to queue
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check if all nodes were processed
	// If not, there's a cycle in the graph
	if len(order) != len(g.nodes) {
		return nil, fmt.Errorf("workflow contains cycles (circular dependencies)")
	}

	return order, nil
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
