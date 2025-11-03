// Package graph provides graph algorithms and utilities for workflow execution.
//
// # Overview
//
// The graph package implements essential graph algorithms used by the workflow
// engine for dependency resolution, cycle detection, and execution ordering.
// It provides a foundation for DAG-based workflow execution.
//
// # Key Algorithms
//
// Topological Sort:
//   - Implements Kahn's algorithm for topological ordering
//   - Ensures nodes execute in dependency order
//   - Detects cycles in the workflow graph
//   - Provides stable, deterministic ordering
//
// Cycle Detection:
//   - Identifies circular dependencies
//   - Reports all nodes involved in cycles
//   - Prevents infinite execution loops
//   - Fast detection using DFS-based approach
//
// Dependency Analysis:
//   - Computes transitive dependencies
//   - Identifies independent subgraphs
//   - Finds parallelizable nodes
//   - Analyzes execution paths
//
// # Graph Representation
//
// Workflows are represented as directed graphs where:
//
//   - Nodes represent workflow operations
//   - Edges represent data dependencies
//   - Direction indicates data flow (source → target)
//   - Multiple edges can connect the same nodes
//
// # Topological Sort Usage
//
//	import "github.com/yesoreyeram/thaiyyal/backend/pkg/graph"
//
//	// Build adjacency list
//	edges := []graph.Edge{
//	    {Source: "A", Target: "B"},
//	    {Source: "B", Target: "C"},
//	}
//
//	// Perform topological sort
//	sorted, err := graph.TopologicalSort(nodes, edges)
//	if err != nil {
//	    // Handle cycle or invalid graph
//	}
//
//	// Execute in sorted order
//	for _, nodeID := range sorted {
//	    execute(nodeID)
//	}
//
// # Cycle Detection
//
//	hasCycle, cyclePath := graph.DetectCycle(nodes, edges)
//	if hasCycle {
//	    fmt.Printf("Cycle detected: %v\n", cyclePath)
//	}
//
// # Parallel Execution Analysis
//
//	// Find nodes that can execute in parallel
//	levels := graph.ComputeExecutionLevels(nodes, edges)
//	for level, nodeIDs := range levels {
//	    // All nodes in same level can execute concurrently
//	    executeParallel(nodeIDs)
//	}
//
// # Performance Characteristics
//
//   - Topological sort: O(V + E) where V=nodes, E=edges
//   - Cycle detection: O(V + E) using DFS
//   - Parallel levels: O(V + E) single pass
//   - Memory efficient: Uses sparse graph representation
//
// # Error Handling
//
// Graph operations return errors for:
//
//   - Cycles: Circular dependencies detected
//   - Invalid edges: References to non-existent nodes
//   - Empty graph: No nodes to process
//   - Disconnected nodes: Orphaned nodes without connections
//
// # Algorithm Details
//
// Kahn's Algorithm (Topological Sort):
//  1. Calculate in-degree for all nodes
//  2. Add zero in-degree nodes to queue
//  3. Process queue: remove node, decrement neighbor in-degrees
//  4. Add newly zero in-degree nodes to queue
//  5. If processed count != node count, cycle exists
//
// DFS Cycle Detection:
//  1. Track visiting state for each node
//  2. Recursive DFS from each unvisited node
//  3. If visiting node encountered, cycle found
//  4. Otherwise, mark as visited
//
// # Use Cases
//
//   - Workflow execution ordering
//   - Dependency resolution
//   - Build systems
//   - Task scheduling
//   - Data pipeline orchestration
//
// # Design Principles
//
//   - Efficiency: Optimized for large graphs
//   - Clarity: Clear, well-documented algorithms
//   - Correctness: Extensively tested edge cases
//   - Flexibility: Generic graph operations
//
// # Thread Safety
//
// Graph algorithms are stateless and safe for concurrent use.
// Multiple goroutines can analyze the same graph simultaneously.
package graph
