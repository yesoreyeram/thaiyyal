package workflow

import (
	"encoding/json"
	"fmt"
)

// Parser handles parsing of workflow JSON payloads
type Parser struct{}

// NewParser creates a new workflow parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a JSON workflow payload into a Workflow struct
func (p *Parser) Parse(jsonData []byte) (*Workflow, error) {
	var workflow Workflow
	if err := json.Unmarshal(jsonData, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse workflow JSON: %w", err)
	}

	// Validate the workflow
	if err := p.validate(&workflow); err != nil {
		return nil, fmt.Errorf("workflow validation failed: %w", err)
	}

	return &workflow, nil
}

// validate performs basic validation on the workflow
func (p *Parser) validate(workflow *Workflow) error {
	if len(workflow.Nodes) == 0 {
		return fmt.Errorf("workflow must contain at least one node")
	}

	// Create a map of node IDs for validation
	nodeMap := make(map[string]bool)
	for _, node := range workflow.Nodes {
		if node.ID == "" {
			return fmt.Errorf("node must have an ID")
		}
		if nodeMap[node.ID] {
			return fmt.Errorf("duplicate node ID: %s", node.ID)
		}
		nodeMap[node.ID] = true
	}

	// Validate edges reference existing nodes
	for _, edge := range workflow.Edges {
		if !nodeMap[edge.Source] {
			return fmt.Errorf("edge references non-existent source node: %s", edge.Source)
		}
		if !nodeMap[edge.Target] {
			return fmt.Errorf("edge references non-existent target node: %s", edge.Target)
		}
	}

	return nil
}
