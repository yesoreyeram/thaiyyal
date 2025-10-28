package workflow

import (
	"encoding/json"
	"fmt"
)

// Engine is the workflow execution engine
type Engine struct {
	payload Payload
}

// NewEngine creates a new workflow engine
func NewEngine(payloadJSON []byte) (*Engine, error) {
	var payload Payload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}
	return &Engine{payload: payload}, nil
}

// Execute runs the workflow and returns the result
func (e *Engine) Execute() (*Result, error) {
	result := &Result{
		NodeResults: make(map[string]interface{}),
		Errors:      []string{},
	}

	// Build adjacency map and in-degree count for topological sort
	nodeMap := make(map[string]Node)
	for _, node := range e.payload.Nodes {
		nodeMap[node.ID] = node
	}

	// Build graph structure
	inDegree := make(map[string]int)
	adjacency := make(map[string][]string) // source -> targets

	for _, node := range e.payload.Nodes {
		if _, exists := inDegree[node.ID]; !exists {
			inDegree[node.ID] = 0
		}
	}

	for _, edge := range e.payload.Edges {
		adjacency[edge.Source] = append(adjacency[edge.Source], edge.Target)
		inDegree[edge.Target]++
	}

	// Topological sort using Kahn's algorithm
	queue := []string{}
	for nodeID, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, nodeID)
		}
	}

	executionOrder := []string{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		executionOrder = append(executionOrder, current)

		for _, neighbor := range adjacency[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for cycles
	if len(executionOrder) != len(e.payload.Nodes) {
		return result, fmt.Errorf("workflow contains cycles or disconnected nodes")
	}

	// Execute nodes in topological order
	for _, nodeID := range executionOrder {
		node := nodeMap[nodeID]
		value, err := e.executeNode(node, result.NodeResults, adjacency)
		if err != nil {
			errMsg := fmt.Sprintf("error executing node %s: %v", nodeID, err)
			result.Errors = append(result.Errors, errMsg)
			return result, fmt.Errorf("error executing node %s: %w", nodeID, err)
		}
		result.NodeResults[nodeID] = value
	}

	// Find the final output (nodes with no outgoing edges)
	finalNodes := []string{}
	for _, node := range e.payload.Nodes {
		if len(adjacency[node.ID]) == 0 {
			finalNodes = append(finalNodes, node.ID)
		}
	}

	if len(finalNodes) > 0 {
		// Use the last final node as the output
		result.FinalOutput = result.NodeResults[finalNodes[len(finalNodes)-1]]
	}

	return result, nil
}

// executeNode executes a single node based on its type
func (e *Engine) executeNode(node Node, nodeResults map[string]interface{}, adjacency map[string][]string) (interface{}, error) {
	// Determine node type by checking which data fields are set
	if node.Data.Value != nil {
		// Number node - return the value directly
		return *node.Data.Value, nil
	}

	if node.Data.Op != nil {
		// Operation node - get inputs from incoming edges
		inputs, err := e.getNodeInputs(node.ID, nodeResults)
		if err != nil {
			return nil, err
		}

		if len(inputs) < 2 {
			return nil, fmt.Errorf("operation node %s requires at least 2 inputs, got %d", node.ID, len(inputs))
		}

		// Perform operation on the first two inputs
		left, ok1 := inputs[0].(float64)
		right, ok2 := inputs[1].(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("operation node %s inputs must be numbers", node.ID)
		}

		switch *node.Data.Op {
		case "add":
			return left + right, nil
		case "subtract":
			return left - right, nil
		case "multiply":
			return left * right, nil
		case "divide":
			if right == 0 {
				return nil, fmt.Errorf("division by zero in node %s", node.ID)
			}
			return left / right, nil
		default:
			return nil, fmt.Errorf("unknown operation: %s", *node.Data.Op)
		}
	}

	if node.Data.Mode != nil {
		// Visualization node - return formatted output
		inputs, err := e.getNodeInputs(node.ID, nodeResults)
		if err != nil {
			return nil, err
		}

		if len(inputs) == 0 {
			return nil, fmt.Errorf("visualization node %s requires at least 1 input", node.ID)
		}

		// Return the input value with visualization metadata
		return map[string]interface{}{
			"mode":  *node.Data.Mode,
			"value": inputs[0],
		}, nil
	}

	return nil, fmt.Errorf("unknown node type for node %s", node.ID)
}

// getNodeInputs retrieves the input values for a node from its predecessors
func (e *Engine) getNodeInputs(nodeID string, nodeResults map[string]interface{}) ([]interface{}, error) {
	inputs := []interface{}{}

	// Find all edges that target this node
	for _, edge := range e.payload.Edges {
		if edge.Target == nodeID {
			sourceResult, exists := nodeResults[edge.Source]
			if !exists {
				return nil, fmt.Errorf("source node %s has not been executed yet", edge.Source)
			}
			inputs = append(inputs, sourceResult)
		}
	}

	return inputs, nil
}
