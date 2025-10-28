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

	// Infer node types if not explicitly set
	e.inferNodeTypes()

	// Build execution plan using topological sort
	executionOrder, err := e.buildExecutionOrder()
	if err != nil {
		return result, err
	}

	// Execute nodes in topological order
	if err := e.executeNodes(executionOrder, result); err != nil {
		return result, err
	}

	// Set final output
	e.setFinalOutput(result)

	return result, nil
}

// inferNodeTypes determines node types from their data fields if not explicitly set
func (e *Engine) inferNodeTypes() {
	for i := range e.payload.Nodes {
		node := &e.payload.Nodes[i]
		if node.Type != "" {
			continue // Type already set
		}

		// Infer type from data fields
		if node.Data.Value != nil {
			node.Type = NodeTypeNumber
		} else if node.Data.Op != nil {
			node.Type = NodeTypeOperation
		} else if node.Data.Mode != nil {
			node.Type = NodeTypeVisualization
		}
	}
}

// buildExecutionOrder creates a topological sort of nodes using Kahn's algorithm
func (e *Engine) buildExecutionOrder() ([]string, error) {
	// Build node map
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
		return nil, fmt.Errorf("workflow contains cycles or disconnected nodes")
	}

	return executionOrder, nil
}

// executeNodes executes all nodes in the given order
func (e *Engine) executeNodes(executionOrder []string, result *Result) error {
	nodeMap := make(map[string]Node)
	for _, node := range e.payload.Nodes {
		nodeMap[node.ID] = node
	}

	for _, nodeID := range executionOrder {
		node := nodeMap[nodeID]
		value, err := e.executeNode(node, result.NodeResults)
		if err != nil {
			errMsg := fmt.Sprintf("error executing node %s: %v", nodeID, err)
			result.Errors = append(result.Errors, errMsg)
			return fmt.Errorf("error executing node %s: %w", nodeID, err)
		}
		result.NodeResults[nodeID] = value
	}

	return nil
}

// executeNode executes a single node based on its type
func (e *Engine) executeNode(node Node, nodeResults map[string]interface{}) (interface{}, error) {
	switch node.Type {
	case NodeTypeNumber:
		return e.executeNumberNode(node)
	case NodeTypeOperation:
		return e.executeOperationNode(node, nodeResults)
	case NodeTypeVisualization:
		return e.executeVisualizationNode(node, nodeResults)
	default:
		return nil, fmt.Errorf("unknown node type '%s' for node %s", node.Type, node.ID)
	}
}

// executeNumberNode executes a number input node
func (e *Engine) executeNumberNode(node Node) (interface{}, error) {
	if node.Data.Value == nil {
		return nil, fmt.Errorf("number node %s missing value", node.ID)
	}
	return *node.Data.Value, nil
}

// executeOperationNode executes an arithmetic operation node
func (e *Engine) executeOperationNode(node Node, nodeResults map[string]interface{}) (interface{}, error) {
	if node.Data.Op == nil {
		return nil, fmt.Errorf("operation node %s missing operation", node.ID)
	}

	// Get inputs from incoming edges
	inputs, err := e.getNodeInputs(node.ID, nodeResults)
	if err != nil {
		return nil, err
	}

	if len(inputs) < 2 {
		return nil, fmt.Errorf("operation node %s requires at least 2 inputs, got %d", node.ID, len(inputs))
	}

	// Extract numeric values
	left, ok1 := inputs[0].(float64)
	right, ok2 := inputs[1].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("operation node %s inputs must be numbers", node.ID)
	}

	// Perform the operation
	return e.performOperation(*node.Data.Op, left, right, node.ID)
}

// performOperation executes the specified arithmetic operation
func (e *Engine) performOperation(op string, left, right float64, nodeID string) (float64, error) {
	switch OperationType(op) {
	case OperationAdd:
		return left + right, nil
	case OperationSubtract:
		return left - right, nil
	case OperationMultiply:
		return left * right, nil
	case OperationDivide:
		if right == 0 {
			return 0, fmt.Errorf("division by zero in node %s", nodeID)
		}
		return left / right, nil
	default:
		return 0, fmt.Errorf("unknown operation: %s", op)
	}
}

// executeVisualizationNode executes a visualization node
func (e *Engine) executeVisualizationNode(node Node, nodeResults map[string]interface{}) (interface{}, error) {
	if node.Data.Mode == nil {
		return nil, fmt.Errorf("visualization node %s missing mode", node.ID)
	}

	// Get inputs from incoming edges
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

// setFinalOutput determines and sets the final output of the workflow
func (e *Engine) setFinalOutput(result *Result) {
	// Find nodes with no outgoing edges (terminal nodes)
	hasOutgoing := make(map[string]bool)
	for _, edge := range e.payload.Edges {
		hasOutgoing[edge.Source] = true
	}

	finalNodes := []string{}
	for _, node := range e.payload.Nodes {
		if !hasOutgoing[node.ID] {
			finalNodes = append(finalNodes, node.ID)
		}
	}

	if len(finalNodes) > 0 {
		// Use the last final node as the output
		result.FinalOutput = result.NodeResults[finalNodes[len(finalNodes)-1]]
	}
}
