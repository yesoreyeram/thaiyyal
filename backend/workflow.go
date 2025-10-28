package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// NodeType represents the type of a workflow node
type NodeType string

const (
	NodeTypeNumber        NodeType = "number"
	NodeTypeOperation     NodeType = "operation"
	NodeTypeVisualization NodeType = "visualization"
	NodeTypeTextInput     NodeType = "text_input"
	NodeTypeTextOperation NodeType = "text_operation"
	NodeTypeHTTP          NodeType = "http"
	NodeTypeCondition     NodeType = "condition"
	NodeTypeForEach       NodeType = "for_each"
	NodeTypeWhileLoop     NodeType = "while_loop"
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
	Value         *float64 `json:"value,omitempty"`          // for number nodes
	Op            *string  `json:"op,omitempty"`             // for operation nodes
	Mode          *string  `json:"mode,omitempty"`           // for visualization nodes
	Label         *string  `json:"label,omitempty"`          // optional label
	Text          *string  `json:"text,omitempty"`           // for text input nodes
	TextOp        *string  `json:"text_op,omitempty"`        // for text operation nodes
	URL           *string  `json:"url,omitempty"`            // for HTTP nodes
	Separator     *string  `json:"separator,omitempty"`      // for concat text operation
	RepeatN       *int     `json:"repeat_n,omitempty"`       // for repeat text operation
	Condition     *string  `json:"condition,omitempty"`      // for condition nodes
	TruePath      *string  `json:"true_path,omitempty"`      // for condition nodes (output port name)
	FalsePath     *string  `json:"false_path,omitempty"`     // for condition nodes (output port name)
	MaxIterations *int     `json:"max_iterations,omitempty"` // for for_each and while_loop nodes
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
		} else if e.nodes[i].Data.URL != nil {
			e.nodes[i].Type = NodeTypeHTTP
		} else if e.nodes[i].Data.Condition != nil {
			e.nodes[i].Type = NodeTypeCondition
		}
		// Note: for_each and while_loop require explicit type as they have MaxIterations
		// which could be confused with other fields
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
	case NodeTypeHTTP:
		return e.executeHTTPNode(node)
	case NodeTypeCondition:
		return e.executeConditionNode(node)
	case NodeTypeForEach:
		return e.executeForEachNode(node)
	case NodeTypeWhileLoop:
		return e.executeWhileLoopNode(node)
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

	// Get input from predecessor node(s)
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("text operation needs at least 1 input")
	}

	// Handle concat operation (can accept multiple inputs)
	if *node.Data.TextOp == "concat" {
		// Validate all inputs are strings
		textInputs := []string{}
		for i, input := range inputs {
			text, ok := input.(string)
			if !ok {
				return nil, fmt.Errorf("concat operation input %d must be text/string", i)
			}
			textInputs = append(textInputs, text)
		}

		// Get separator (default to empty string)
		separator := ""
		if node.Data.Separator != nil {
			separator = *node.Data.Separator
		}

		// Concatenate all inputs
		result := ""
		for i, text := range textInputs {
			if i > 0 {
				result += separator
			}
			result += text
		}
		return result, nil
	}

	// Handle repeat operation
	if *node.Data.TextOp == "repeat" {
		// Validate single input is a string
		inputText, ok := inputs[0].(string)
		if !ok {
			return nil, fmt.Errorf("repeat operation input must be text/string")
		}

		// Get repeat count (required)
		if node.Data.RepeatN == nil {
			return nil, fmt.Errorf("repeat operation requires repeat_n field")
		}

		repeatCount := *node.Data.RepeatN
		if repeatCount < 0 {
			return nil, fmt.Errorf("repeat_n must be non-negative, got %d", repeatCount)
		}

		// Repeat the text
		result := ""
		for i := 0; i < repeatCount; i++ {
			result += inputText
		}
		return result, nil
	}

	// For other operations, validate single input is a string
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

// executeHTTPNode performs HTTP request and returns response body
func (e *Engine) executeHTTPNode(node Node) (interface{}, error) {
	if node.Data.URL == nil {
		return nil, fmt.Errorf("HTTP node missing url")
	}

	// Make HTTP GET request
	resp, err := http.Get(*node.Data.URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request returned error status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
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

// executeConditionNode evaluates a condition and passes through the input based on the result
func (e *Engine) executeConditionNode(node Node) (interface{}, error) {
	if node.Data.Condition == nil {
		return nil, fmt.Errorf("condition node missing condition")
	}

	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("condition node needs at least 1 input")
	}

	input := inputs[0]
	conditionMet := e.evaluateCondition(*node.Data.Condition, input)

	// Return the input value along with metadata about which path was taken
	return map[string]interface{}{
		"value":         input,
		"condition_met": conditionMet,
		"condition":     *node.Data.Condition,
	}, nil
}

// executeForEachNode iterates over an array input and returns processed results
func (e *Engine) executeForEachNode(node Node) (interface{}, error) {
	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("for_each node needs at least 1 input")
	}

	// Check if input is an array (slice)
	inputArray, ok := inputs[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("for_each node requires array input, got %T", inputs[0])
	}

	// Set default max iterations
	maxIter := 1000
	if node.Data.MaxIterations != nil && *node.Data.MaxIterations > 0 {
		maxIter = *node.Data.MaxIterations
	}

	// Limit iterations to prevent infinite loops
	iterCount := len(inputArray)
	if iterCount > maxIter {
		return nil, fmt.Errorf("for_each exceeds max iterations: %d > %d", iterCount, maxIter)
	}

	// For now, just return the array (in a real implementation, this would
	// iterate and execute child nodes for each element)
	return map[string]interface{}{
		"items":      inputArray,
		"count":      len(inputArray),
		"iterations": iterCount,
	}, nil
}

// executeWhileLoopNode executes a loop while a condition is true
func (e *Engine) executeWhileLoopNode(node Node) (interface{}, error) {
	if node.Data.Condition == nil {
		return nil, fmt.Errorf("while_loop node missing condition")
	}

	inputs := e.getNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("while_loop node needs at least 1 input")
	}

	// Set default max iterations
	maxIter := 100
	if node.Data.MaxIterations != nil && *node.Data.MaxIterations > 0 {
		maxIter = *node.Data.MaxIterations
	}

	currentValue := inputs[0]
	iterationCount := 0

	// Loop while condition is met
	for e.evaluateCondition(*node.Data.Condition, currentValue) && iterationCount < maxIter {
		iterationCount++
		// In a real implementation, this would execute child nodes
		// For now, we'll just track iterations
	}

	if iterationCount >= maxIter {
		return nil, fmt.Errorf("while_loop exceeded max iterations: %d", maxIter)
	}

	return map[string]interface{}{
		"final_value": currentValue,
		"iterations":  iterationCount,
		"condition":   *node.Data.Condition,
	}, nil
}

// evaluateCondition evaluates a condition string against an input value
func (e *Engine) evaluateCondition(condition string, value interface{}) bool {
	// Simple condition evaluation
	// Supports: ">N", "<N", ">=N", "<=N", "==N", "!=N", "true", "false"
	
	if condition == "true" {
		return true
	}
	if condition == "false" {
		return false
	}

	// Check for numeric comparisons
	numVal, ok := value.(float64)
	if !ok {
		// Try to get numeric value from maps
		if m, isMap := value.(map[string]interface{}); isMap {
			if v, exists := m["value"]; exists {
				numVal, ok = v.(float64)
			}
		}
		if !ok {
			return false
		}
	}

	// Parse condition
	var threshold float64
	var operator string
	
	if len(condition) >= 2 {
		if condition[0] == '>' && condition[1] == '=' {
			operator = ">="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		} else if condition[0] == '<' && condition[1] == '=' {
			operator = "<="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		} else if condition[0] == '=' && condition[1] == '=' {
			operator = "=="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		} else if condition[0] == '!' && condition[1] == '=' {
			operator = "!="
			fmt.Sscanf(condition[2:], "%f", &threshold)
		} else if condition[0] == '>' {
			operator = ">"
			fmt.Sscanf(condition[1:], "%f", &threshold)
		} else if condition[0] == '<' {
			operator = "<"
			fmt.Sscanf(condition[1:], "%f", &threshold)
		}
	}

	switch operator {
	case ">":
		return numVal > threshold
	case "<":
		return numVal < threshold
	case ">=":
		return numVal >= threshold
	case "<=":
		return numVal <= threshold
	case "==":
		return numVal == threshold
	case "!=":
		return numVal != threshold
	default:
		return false
	}
}
