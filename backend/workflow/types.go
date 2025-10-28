package workflow

// NodeType represents the type of a workflow node
type NodeType string

const (
	// NodeTypeNumber represents a numeric input node
	NodeTypeNumber NodeType = "number"
	// NodeTypeOperation represents an arithmetic operation node
	NodeTypeOperation NodeType = "operation"
	// NodeTypeVisualization represents a visualization/output node
	NodeTypeVisualization NodeType = "visualization"
)

// OperationType represents the type of arithmetic operation
type OperationType string

const (
	// OperationAdd performs addition
	OperationAdd OperationType = "add"
	// OperationSubtract performs subtraction
	OperationSubtract OperationType = "subtract"
	// OperationMultiply performs multiplication
	OperationMultiply OperationType = "multiply"
	// OperationDivide performs division
	OperationDivide OperationType = "divide"
)

// VisualizationMode represents the visualization display mode
type VisualizationMode string

const (
	// VisualizationModeText displays output as plain text
	VisualizationModeText VisualizationMode = "text"
	// VisualizationModeTable displays output in table format
	VisualizationModeTable VisualizationMode = "table"
)

// Payload represents the JSON payload from the frontend
type Payload struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Node represents a workflow node
type Node struct {
	ID   string   `json:"id"`
	Type NodeType `json:"type,omitempty"` // explicit node type
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

// NumberNodeData represents data for a number input node
type NumberNodeData struct {
	Value float64 `json:"value"`
	Label string  `json:"label,omitempty"`
}

// OperationNodeData represents data for an operation node
type OperationNodeData struct {
	Operation OperationType `json:"op"`
	Label     string        `json:"label,omitempty"`
}

// VisualizationNodeData represents data for a visualization node
type VisualizationNodeData struct {
	Mode  VisualizationMode `json:"mode"`
	Label string            `json:"label,omitempty"`
}
