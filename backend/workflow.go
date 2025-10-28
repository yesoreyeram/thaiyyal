package workflow

import (
	"encoding/json"
	"fmt"
)

// NodeType represents the type of a workflow node
type NodeType string

const (
	NodeTypeNumber        NodeType = "number"
	NodeTypeOperation     NodeType = "operation"
	NodeTypeVisualization NodeType = "visualization"
	NodeTypeTextInput     NodeType = "text_input"
	NodeTypeTextOperation NodeType = "text_operation"
)

// Payload represents the JSON payload from the frontend
type Payload struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Node represents a workflow node
type Node struct {
	ID   string   `json:"id"`
	Type NodeType `json:"type,omitempty"`
	Data NodeData `json:"data"`
}

// NodeData contains the node-specific configuration
type NodeData struct {
	Value  *float64 `json:"value,omitempty"`   // for number nodes
	Op     *string  `json:"op,omitempty"`      // for operation nodes
	Mode   *string  `json:"mode,omitempty"`    // for visualization nodes
	Label  *string  `json:"label,omitempty"`   // optional label
	Text   *string  `json:"text,omitempty"`    // for text input nodes
	TextOp *string  `json:"text_op,omitempty"` // for text operation nodes
}

// Edge represents a connection between nodes
type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// Result represents the execution result of the workflow
type Result struct {
	NodeResults map[string]interface{} `json:"node_results"`
	FinalOutput interface{}            `json:"final_output"`
	Errors      []string               `json:"errors,omitempty"`
}

// Engine is the workflow execution engine
type Engine struct {
	nodes       []Node
	edges       []Edge
	nodeResults map[string]interface{}
}

// NewEngine creates a new workflow engine from JSON payload
func NewEngine(payloadJSON []byte) (*Engine, error) {
	var payload Payload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	return &Engine{
		nodes:       payload.Nodes,
		edges:       payload.Edges,
		nodeResults: make(map[string]interface{}),
	}, nil
}

// Execute runs the workflow and returns the result
func (e *Engine) Execute() (*Result, error) {
	result := &Result{
		NodeResults: make(map[string]interface{}),
		Errors:      []string{},
	}

	// Infer node types if not set
	e.inferNodeTypes()

	// Get execution order using topological sort
	executionOrder, err := e.topologicalSort()
	if err != nil {
		return result, err
	}

	// Execute each node in order
	for _, nodeID := range executionOrder {
		node := e.getNode(nodeID)
		value, err := e.executeNode(node)
		if err != nil {
			errMsg := fmt.Sprintf("error executing node %s: %v", nodeID, err)
			result.Errors = append(result.Errors, errMsg)
			return result, fmt.Errorf("%s", errMsg)
		}
		e.nodeResults[nodeID] = value
	}

	// Copy results and set final output
	result.NodeResults = e.nodeResults
	result.FinalOutput = e.getFinalOutput()

	return result, nil
}

// inferNodeTypes determines node types from data if not explicitly set
func (e *Engine) inferNodeTypes() {
	for i := range e.nodes {
		if e.nodes[i].Type != "" {
			continue
		}
		// Infer type from data fields
		if e.nodes[i].Data.Value != nil {
			e.nodes[i].Type = NodeTypeNumber
		} else if e.nodes[i].Data.Op != nil {
			e.nodes[i].Type = NodeTypeOperation
		} else if e.nodes[i].Data.Mode != nil {
			e.nodes[i].Type = NodeTypeVisualization
		} else if e.nodes[i].Data.Text != nil {
			e.nodes[i].Type = NodeTypeTextInput
		} else if e.nodes[i].Data.TextOp != nil {
			e.nodes[i].Type = NodeTypeTextOperation
		}
	}
}

// topologicalSort returns execution order using Kahn's algorithm
func (e *Engine) topologicalSort() ([]string, error) {
	// Build in-degree map and adjacency list
	inDegree := make(map[string]int)
	adjacency := make(map[string][]string)

	// Initialize in-degree for all nodes
	for _, node := range e.nodes {
		inDegree[node.ID] = 0
	}

	// Build graph
	for _, edge := range e.edges {
		adjacency[edge.Source] = append(adjacency[edge.Source], edge.Target)
		inDegree[edge.Target]++
	}

	// Find nodes with no dependencies
	queue := []string{}
	for nodeID, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, nodeID)
		}
	}

	// Process nodes
	order := []string{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)

		// Reduce in-degree for neighbors
		for _, neighbor := range adjacency[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for cycles
	if len(order) != len(e.nodes) {
		return nil, fmt.Errorf("workflow contains cycles")
	}

	return order, nil
}

// executeNode executes a single node based on its type
func (e *Engine) executeNode(node Node) (interface{}, error) {
	switch node.Type {
	case NodeTypeNumber:
		return e.executeNumberNode(node)
	case NodeTypeOperation:
		return e.executeOperationNode(node)
	case NodeTypeVisualization:
		return e.executeVisualizationNode(node)
	case NodeTypeTextInput:
		return e.executeTextInputNode(node)
	case NodeTypeTextOperation:
		return e.executeTextOperationNode(node)
	default:
		return nil, fmt.Errorf("unknown node type: %s", node.Type)
	}
}

// executeNumberNode returns the number value
func (e *Engine) executeNumberNode(node Node) (interface{}, error) {
	if node.Data.Value == nil {
		return nil, fmt.Errorf("number node missing value")
	}
	return *node.Data.Value, nil
}

// executeOperationNode performs arithmetic operation
func (e *Engine) executeOperationNode(node Node) (interface{}, error) {
	if node.Data.Op == nil {
		return nil, fmt.Errorf("operation node missing op")
	}

	// Get inputs from predecessor nodes
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) < 2 {
		return nil, fmt.Errorf("operation needs 2 inputs, got %d", len(inputs))
	}

	// Convert to numbers
	left, ok1 := inputs[0].(float64)
	right, ok2 := inputs[1].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("operation inputs must be numbers")
	}

	// Perform operation
	switch *node.Data.Op {
	case "add":
		return left + right, nil
	case "subtract":
		return left - right, nil
	case "multiply":
		return left * right, nil
	case "divide":
		if right == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return left / right, nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", *node.Data.Op)
	}
}

// executeVisualizationNode formats output for display
func (e *Engine) executeVisualizationNode(node Node) (interface{}, error) {
	if node.Data.Mode == nil {
		return nil, fmt.Errorf("visualization node missing mode")
	}

	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("visualization needs at least 1 input")
	}

	return map[string]interface{}{
		"mode":  *node.Data.Mode,
		"value": inputs[0],
	}, nil
}

// executeTextInputNode returns the text value
func (e *Engine) executeTextInputNode(node Node) (interface{}, error) {
	if node.Data.Text == nil {
		return nil, fmt.Errorf("text input node missing text")
	}
	return *node.Data.Text, nil
}

// executeTextOperationNode performs text transformation
func (e *Engine) executeTextOperationNode(node Node) (interface{}, error) {
	if node.Data.TextOp == nil {
		return nil, fmt.Errorf("text operation node missing text_op")
	}

	// Get input from predecessor node
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("text operation needs at least 1 input")
	}

	// Validate that input is a string
	inputText, ok := inputs[0].(string)
	if !ok {
		return nil, fmt.Errorf("text operation input must be text/string")
	}

	// Perform text operation
	switch *node.Data.TextOp {
	case "uppercase":
		return toUpperCase(inputText), nil
	case "lowercase":
		return toLowerCase(inputText), nil
	case "titlecase":
		return toTitleCase(inputText), nil
	case "camelcase":
		return toCamelCase(inputText), nil
	case "inversecase":
		return toInverseCase(inputText), nil
	default:
		return nil, fmt.Errorf("unknown text operation: %s", *node.Data.TextOp)
	}
}

// Text transformation helper functions
func toUpperCase(s string) string {
	result := []rune{}
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			result = append(result, r-32)
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func toLowerCase(s string) string {
	result := []rune{}
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result = append(result, r+32)
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func toTitleCase(s string) string {
	result := []rune{}
	capitalizeNext := true
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' {
			result = append(result, r)
			capitalizeNext = true
		} else if capitalizeNext && r >= 'a' && r <= 'z' {
			result = append(result, r-32)
			capitalizeNext = false
		} else if !capitalizeNext && r >= 'A' && r <= 'Z' {
			result = append(result, r+32)
			capitalizeNext = false
		} else {
			result = append(result, r)
			capitalizeNext = false
		}
	}
	return string(result)
}

func toCamelCase(s string) string {
	result := []rune{}
	capitalizeNext := false
	firstChar := true

	for _, r := range s {
		if r == ' ' || r == '_' || r == '-' {
			capitalizeNext = true
			continue
		}

		if firstChar {
			if r >= 'A' && r <= 'Z' {
				result = append(result, r+32)
			} else {
				result = append(result, r)
			}
			firstChar = false
		} else if capitalizeNext {
			if r >= 'a' && r <= 'z' {
				result = append(result, r-32)
			} else {
				result = append(result, r)
			}
			capitalizeNext = false
		} else {
			if r >= 'A' && r <= 'Z' {
				result = append(result, r+32)
			} else {
				result = append(result, r)
			}
		}
	}
	return string(result)
}

func toInverseCase(s string) string {
	result := []rune{}
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			result = append(result, r-32)
		} else if r >= 'A' && r <= 'Z' {
			result = append(result, r+32)
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// getNodeInputs returns input values for a node from its predecessors
func (e *Engine) getNodeInputs(nodeID string) []interface{} {
	inputs := []interface{}{}
	for _, edge := range e.edges {
		if edge.Target == nodeID {
			if result, exists := e.nodeResults[edge.Source]; exists {
				inputs = append(inputs, result)
			}
		}
	}
	return inputs
}

// getNode returns a node by ID
func (e *Engine) getNode(nodeID string) Node {
	for _, node := range e.nodes {
		if node.ID == nodeID {
			return node
		}
	}
	return Node{}
}

// getFinalOutput returns the output from terminal nodes (nodes with no outgoing edges)
func (e *Engine) getFinalOutput() interface{} {
	// Find which nodes have outgoing edges
	hasOutgoing := make(map[string]bool)
	for _, edge := range e.edges {
		hasOutgoing[edge.Source] = true
	}

	// Find terminal nodes
	var finalNodeID string
	for _, node := range e.nodes {
		if !hasOutgoing[node.ID] {
			finalNodeID = node.ID
		}
	}

	// Return result of the last terminal node
	if finalNodeID != "" {
		return e.nodeResults[finalNodeID]
	}
	return nil
}
