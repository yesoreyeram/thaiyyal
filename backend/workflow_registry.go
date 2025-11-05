package workflow

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// WorkflowMeta represents a stored workflow with metadata
type WorkflowMeta struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Data        json.RawMessage `json:"data"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// WorkflowSummary represents a lightweight workflow reference for listing
type WorkflowSummary struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WorkflowRegistry manages stored workflows by their IDs
type WorkflowRegistry struct {
	workflows map[string]*WorkflowMeta
	mu        sync.RWMutex
}

// NewWorkflowRegistry creates a new workflow registry
func NewWorkflowRegistry() *WorkflowRegistry {
	return &WorkflowRegistry{
		workflows: make(map[string]*WorkflowMeta),
	}
}

// Register adds a workflow to the registry and returns its ID
func (r *WorkflowRegistry) Register(name, description string, data json.RawMessage) (string, error) {
	if name == "" {
		return "", fmt.Errorf("workflow name is required")
	}

	if len(data) == 0 {
		return "", fmt.Errorf("workflow data is required")
	}

	// Validate that data is valid JSON
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return "", fmt.Errorf("invalid workflow data: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	now := time.Now()

	workflow := &WorkflowMeta{
		ID:          id,
		Name:        name,
		Description: description,
		Data:        data,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	r.workflows[id] = workflow

	return id, nil
}

// Update updates an existing workflow
func (r *WorkflowRegistry) Update(id, name, description string, data json.RawMessage) error {
	if id == "" {
		return fmt.Errorf("workflow ID is required")
	}

	if name == "" {
		return fmt.Errorf("workflow name is required")
	}

	if len(data) == 0 {
		return fmt.Errorf("workflow data is required")
	}

	// Validate that data is valid JSON
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("invalid workflow data: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	workflow, exists := r.workflows[id]
	if !exists {
		return fmt.Errorf("workflow with ID %s not found", id)
	}

	workflow.Name = name
	workflow.Description = description
	workflow.Data = data
	workflow.UpdatedAt = time.Now()

	return nil
}

// Get retrieves a workflow by ID
func (r *WorkflowRegistry) Get(id string) (*WorkflowMeta, error) {
	if id == "" {
		return nil, fmt.Errorf("workflow ID is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	workflow, exists := r.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow with ID %s not found", id)
	}

	// Return a copy to prevent external modifications
	workflowCopy := &WorkflowMeta{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		Data:        make(json.RawMessage, len(workflow.Data)),
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
	}
	copy(workflowCopy.Data, workflow.Data)

	return workflowCopy, nil
}

// Unregister removes a workflow by ID
func (r *WorkflowRegistry) Unregister(id string) error {
	if id == "" {
		return fmt.Errorf("workflow ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workflows[id]; !exists {
		return fmt.Errorf("workflow with ID %s not found", id)
	}

	delete(r.workflows, id)

	return nil
}

// List returns all workflow summaries
func (r *WorkflowRegistry) List() []WorkflowSummary {
	r.mu.RLock()
	defer r.mu.RUnlock()

	summaries := make([]WorkflowSummary, 0, len(r.workflows))

	for _, workflow := range r.workflows {
		summaries = append(summaries, WorkflowSummary{
			ID:          workflow.ID,
			Name:        workflow.Name,
			Description: workflow.Description,
			CreatedAt:   workflow.CreatedAt,
			UpdatedAt:   workflow.UpdatedAt,
		})
	}

	return summaries
}

// Has checks if a workflow exists
func (r *WorkflowRegistry) Has(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.workflows[id]
	return exists
}

// Count returns the number of registered workflows
func (r *WorkflowRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.workflows)
}

// Clear removes all workflows from the registry
func (r *WorkflowRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.workflows = make(map[string]*WorkflowMeta)
}
