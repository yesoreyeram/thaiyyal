package workflow

// Payload represents the JSON payload from the frontend
type Payload struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Node represents a workflow node
type Node struct {
	ID   string   `json:"id"`
	Data NodeData `json:"data"`
}

// NodeData contains the node-specific configuration
type NodeData struct {
	Value *float64 `json:"value,omitempty"` // for number nodes
	Op    *string  `json:"op,omitempty"`    // for operation nodes (add, subtract, multiply, divide)
	Mode  *string  `json:"mode,omitempty"`  // for visualization nodes (text, table)
	Label *string  `json:"label,omitempty"` // optional label
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
