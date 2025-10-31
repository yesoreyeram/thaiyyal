package workflow

import "fmt"

// ============================================================================
// Graph Utilities
// ============================================================================
// This file contains helper functions for workflow graph operations:
// - Topological sorting (Kahn's algorithm for DAG execution order)
// - Node lookups and relationships
// - Final output determination
// ============================================================================

// topologicalSort performs topological sorting on the workflow graph using Kahn's algorithm.
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
func (e *Engine) topologicalSort() ([]string, error) {
	// Build in-degree map and adjacency list
	inDegree := make(map[string]int)
	adjacency := make(map[string][]string)

	// Initialize in-degree for all nodes to zero
	for _, node := range e.nodes {
		inDegree[node.ID] = 0
	}

	// Build the graph structure
	for _, edge := range e.edges {
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
	if len(order) != len(e.nodes) {
		return nil, fmt.Errorf("workflow contains cycles (circular dependencies)")
	}

	return order, nil
}

// getNode retrieves a node by its ID.
// This is a linear search through the nodes slice.
//
// Returns:
//   - Node: The node with the matching ID, or empty Node{} if not found
func (e *Engine) getNode(nodeID string) Node {
	for _, node := range e.nodes {
		if node.ID == nodeID {
			return node
		}
	}
	return Node{}
}

// getNodeInputs retrieves all input values for a node from its predecessor nodes.
// It looks at all edges where the target is the specified node and returns
// the results of those source nodes.
//
// This function is thread-safe for use in parallel execution.
//
// Returns:
//   - []interface{}: Slice of input values from predecessor nodes, in edge order
func (e *Engine) getNodeInputs(nodeID string) []interface{} {
	inputs := []interface{}{}
	e.resultsMutex.RLock()
	defer e.resultsMutex.RUnlock()
	
	for _, edge := range e.edges {
		if edge.Target == nodeID {
			if result, ok := e.nodeResults[edge.Source]; ok {
				inputs = append(inputs, result)
			}
		}
	}
	return inputs
}

// getFinalOutput determines the final output of the workflow.
// The final output is the result of a terminal node (node with no outgoing edges).
// Context nodes (context_variable, context_constant) are excluded from being final output
// unless they are the ONLY nodes in the workflow.
//
// If multiple terminal nodes exist, returns the first non-context one found.
// If no terminal nodes exist (all nodes have outgoing edges), returns nil.
//
// This function is thread-safe for use in parallel execution.
//
// Returns:
//   - interface{}: The result value from a terminal node, or nil if none found
func (e *Engine) getFinalOutput() interface{} {
	// Build a set of all terminal nodes (nodes with no outgoing edges)
	terminalNodes := make(map[string]bool)
	
	// Initially, all nodes are considered terminal
	for _, node := range e.nodes {
		terminalNodes[node.ID] = true
	}
	
	// Remove nodes that have outgoing edges
	for _, edge := range e.edges {
		terminalNodes[edge.Source] = false
	}

	e.resultsMutex.RLock()
	defer e.resultsMutex.RUnlock()
	
	// First pass: Try to find a non-context terminal node
	for nodeID, isTerminal := range terminalNodes {
		if isTerminal {
			node := e.getNode(nodeID)
			if node.Type != NodeTypeContextVariable && node.Type != NodeTypeContextConstant {
				if result, ok := e.nodeResults[nodeID]; ok {
					return result
				}
			}
		}
	}
	
	// Second pass: If no non-context terminal found, return any terminal (including context)
	// This handles the case where workflow contains only context nodes
	for nodeID, isTerminal := range terminalNodes {
		if isTerminal {
			if result, ok := e.nodeResults[nodeID]; ok {
				return result
			}
		}
	}

	// No terminal node found
	return nil
}
