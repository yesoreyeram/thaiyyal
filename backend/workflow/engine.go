package workflow

import (
	"fmt"
	"sort"
)

// Engine is the workflow execution engine
type Engine struct {
	workflow   *Workflow
	nodeMap    map[string]*Node     // Map of node ID to node for O(1) lookup
	inputsMap  map[string][]string  // Map of node ID to input source IDs for O(1) lookup
}

// NewEngine creates a new workflow execution engine
func NewEngine(workflow *Workflow) *Engine {
	engine := &Engine{
		workflow:  workflow,
		nodeMap:   make(map[string]*Node),
		inputsMap: make(map[string][]string),
	}
	
	// Build node map for O(1) lookups
	for i := range workflow.Nodes {
		engine.nodeMap[workflow.Nodes[i].ID] = &workflow.Nodes[i]
	}
	
	// Build inputs map for O(1) lookups
	for _, edge := range workflow.Edges {
		engine.inputsMap[edge.Target] = append(engine.inputsMap[edge.Target], edge.Source)
	}
	
	return engine
}

// Execute runs the workflow and returns the execution result
func (e *Engine) Execute() (*ExecutionResult, error) {
	// Build dependency graph
	dependencies := e.buildDependencyGraph()

	// Topologically sort nodes
	sortedNodes, err := e.topologicalSort(dependencies)
	if err != nil {
		return nil, fmt.Errorf("failed to sort workflow nodes: %w", err)
	}

	// Execute nodes in order
	results := make(map[string]*NodeResult)
	var finalOutput interface{}

	for _, nodeID := range sortedNodes {
		node := e.getNodeByID(nodeID)
		if node == nil {
			return nil, fmt.Errorf("node not found: %s", nodeID)
		}

		result, err := e.executeNode(node, results)
		if err != nil {
			return &ExecutionResult{
				Results: results,
			}, fmt.Errorf("error executing node %s: %w", nodeID, err)
		}

		results[nodeID] = result

		// Check if this is a visualization node (final output)
		nodeType := e.inferNodeType(node)
		if nodeType == "vizNode" {
			finalOutput = result.Value
		}
	}

	return &ExecutionResult{
		Results: results,
		Output:  finalOutput,
	}, nil
}

// buildDependencyGraph creates a map of node dependencies
func (e *Engine) buildDependencyGraph() map[string][]string {
	dependencies := make(map[string][]string)

	// Initialize all nodes
	for _, node := range e.workflow.Nodes {
		dependencies[node.ID] = []string{}
	}

	// Add dependencies based on edges
	for _, edge := range e.workflow.Edges {
		dependencies[edge.Target] = append(dependencies[edge.Target], edge.Source)
	}

	return dependencies
}

// topologicalSort performs topological sorting on the workflow nodes
func (e *Engine) topologicalSort(dependencies map[string][]string) ([]string, error) {
	var sorted []string
	visited := make(map[string]bool)
	visiting := make(map[string]bool)

	var visit func(string) error
	visit = func(nodeID string) error {
		if visited[nodeID] {
			return nil
		}
		if visiting[nodeID] {
			return fmt.Errorf("circular dependency detected at node: %s", nodeID)
		}

		visiting[nodeID] = true

		// Visit dependencies first
		for _, dep := range dependencies[nodeID] {
			if err := visit(dep); err != nil {
				return err
			}
		}

		visiting[nodeID] = false
		visited[nodeID] = true
		sorted = append(sorted, nodeID)
		return nil
	}

	// Get sorted node IDs for deterministic ordering
	nodeIDs := make([]string, 0, len(dependencies))
	for nodeID := range dependencies {
		nodeIDs = append(nodeIDs, nodeID)
	}
	sort.Strings(nodeIDs)

	// Visit all nodes
	for _, nodeID := range nodeIDs {
		if err := visit(nodeID); err != nil {
			return nil, err
		}
	}

	return sorted, nil
}

// executeNode executes a single node
func (e *Engine) executeNode(node *Node, results map[string]*NodeResult) (*NodeResult, error) {
	nodeType := e.inferNodeType(node)

	switch nodeType {
	case "numberNode":
		return e.executeNumberNode(node)
	case "opNode":
		return e.executeOperationNode(node, results)
	case "vizNode":
		return e.executeVisualizationNode(node, results)
	default:
		return nil, fmt.Errorf("unknown node type for node %s", node.ID)
	}
}

// executeNumberNode executes a number input node
func (e *Engine) executeNumberNode(node *Node) (*NodeResult, error) {
	if node.Data.Value == nil {
		return &NodeResult{
			NodeID: node.ID,
			Value:  0.0,
		}, nil
	}

	return &NodeResult{
		NodeID: node.ID,
		Value:  *node.Data.Value,
	}, nil
}

// executeOperationNode executes an operation node
func (e *Engine) executeOperationNode(node *Node, results map[string]*NodeResult) (*NodeResult, error) {
	// Get input values from dependencies
	inputs := e.getNodeInputs(node.ID, results)

	if len(inputs) == 0 {
		return &NodeResult{
			NodeID: node.ID,
			Value:  0.0,
		}, nil
	}

	// Convert inputs to floats
	var values []float64
	for _, input := range inputs {
		if val, ok := input.(float64); ok {
			values = append(values, val)
		} else {
			return nil, fmt.Errorf("invalid input type for operation node %s", node.ID)
		}
	}

	// Determine operation
	op := "add"
	if node.Data.Op != nil {
		op = *node.Data.Op
	}

	// Execute operation
	var result float64
	switch op {
	case "add":
		result = 0
		for _, v := range values {
			result += v
		}
	case "subtract":
		if len(values) == 0 {
			result = 0
		} else {
			result = values[0]
			for i := 1; i < len(values); i++ {
				result -= values[i]
			}
		}
	case "multiply":
		result = 1
		for _, v := range values {
			result *= v
		}
	case "divide":
		if len(values) == 0 {
			result = 0
		} else {
			result = values[0]
			for i := 1; i < len(values); i++ {
				if values[i] == 0 {
					return nil, fmt.Errorf("division by zero in node %s", node.ID)
				}
				result /= values[i]
			}
		}
	default:
		return nil, fmt.Errorf("unknown operation: %s", op)
	}

	return &NodeResult{
		NodeID: node.ID,
		Value:  result,
	}, nil
}

// executeVisualizationNode executes a visualization node
func (e *Engine) executeVisualizationNode(node *Node, results map[string]*NodeResult) (*NodeResult, error) {
	// Get input values from dependencies
	inputs := e.getNodeInputs(node.ID, results)

	mode := "text"
	if node.Data.Mode != nil {
		mode = *node.Data.Mode
	}

	var output interface{}
	switch mode {
	case "text":
		// Simple text output
		if len(inputs) > 0 {
			output = fmt.Sprintf("Result: %v", inputs[0])
		} else {
			output = "Result: (no input)"
		}
	case "table":
		// Table format output
		output = map[string]interface{}{
			"type":   "table",
			"values": inputs,
		}
	default:
		output = inputs
	}

	return &NodeResult{
		NodeID: node.ID,
		Value:  output,
	}, nil
}

// getNodeInputs retrieves the input values for a node from its dependencies
func (e *Engine) getNodeInputs(nodeID string, results map[string]*NodeResult) []interface{} {
	var inputs []interface{}

	// Use the pre-built inputs map for O(1) lookup
	sourceIDs := e.inputsMap[nodeID]
	for _, sourceID := range sourceIDs {
		if result, ok := results[sourceID]; ok && result.Error == nil {
			inputs = append(inputs, result.Value)
		}
	}

	return inputs
}

// getNodeByID retrieves a node by its ID
func (e *Engine) getNodeByID(nodeID string) *Node {
	return e.nodeMap[nodeID]
}

// inferNodeType infers the node type from node data
func (e *Engine) inferNodeType(node *Node) string {
	// Use explicit type if available
	if node.Type != "" {
		return node.Type
	}

	// Infer from data
	if node.Data.Value != nil {
		return "numberNode"
	}
	if node.Data.Op != nil {
		return "opNode"
	}
	if node.Data.Mode != nil {
		return "vizNode"
	}

	return "unknown"
}
