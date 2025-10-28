package workflow

// NodeData represents the data field in a workflow node
type NodeData struct {
	Value *float64 `json:"value,omitempty"` // For number nodes
	Op    *string  `json:"op,omitempty"`    // For operation nodes
	Mode  *string  `json:"mode,omitempty"`  // For visualization nodes
	Label *string  `json:"label,omitempty"` // Node label
}

// Node represents a workflow node
type Node struct {
	ID   string   `json:"id"`
	Data NodeData `json:"data"`
	Type string   `json:"type,omitempty"` // numberNode, opNode, vizNode
}

// Edge represents a connection between two nodes
type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// Workflow represents the complete workflow structure
type Workflow struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// NodeResult holds the execution result of a node
type NodeResult struct {
	NodeID string
	Value  interface{} // Can be number, string, or structured data
	Error  error
}

// ExecutionResult holds the complete execution result
type ExecutionResult struct {
	Results map[string]*NodeResult // Map of nodeID to result
	Output  interface{}            // Final output/visualization
}
