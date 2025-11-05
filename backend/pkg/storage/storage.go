package storage

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Workflow represents a stored workflow with metadata
type Workflow struct {
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

// Store defines the interface for workflow storage operations
type Store interface {
	// Save creates or updates a workflow
	Save(name, description string, data json.RawMessage) (string, error)
	
	// Update updates an existing workflow
	Update(id, name, description string, data json.RawMessage) error
	
	// Load retrieves a workflow by ID
	Load(id string) (*Workflow, error)
	
	// Delete removes a workflow by ID
	Delete(id string) error
	
	// List returns all workflow summaries
	List() []WorkflowSummary
	
	// Exists checks if a workflow exists
	Exists(id string) bool
}

// InMemoryStore implements Store using in-memory storage
type InMemoryStore struct {
	workflows map[string]*Workflow
	mu        sync.RWMutex
}

// NewInMemoryStore creates a new in-memory workflow store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		workflows: make(map[string]*Workflow),
	}
}

// Save creates a new workflow and returns its ID
func (s *InMemoryStore) Save(name, description string, data json.RawMessage) (string, error) {
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
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	id := uuid.New().String()
	now := time.Now()
	
	workflow := &Workflow{
		ID:          id,
		Name:        name,
		Description: description,
		Data:        data,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	s.workflows[id] = workflow
	
	return id, nil
}

// Update updates an existing workflow
func (s *InMemoryStore) Update(id, name, description string, data json.RawMessage) error {
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
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	workflow, exists := s.workflows[id]
	if !exists {
		return fmt.Errorf("workflow with ID %s not found", id)
	}
	
	workflow.Name = name
	workflow.Description = description
	workflow.Data = data
	workflow.UpdatedAt = time.Now()
	
	return nil
}

// Load retrieves a workflow by ID
func (s *InMemoryStore) Load(id string) (*Workflow, error) {
	if id == "" {
		return nil, fmt.Errorf("workflow ID is required")
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	workflow, exists := s.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow with ID %s not found", id)
	}
	
	// Return a copy to prevent external modifications
	workflowCopy := &Workflow{
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

// Delete removes a workflow by ID
func (s *InMemoryStore) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("workflow ID is required")
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.workflows[id]; !exists {
		return fmt.Errorf("workflow with ID %s not found", id)
	}
	
	delete(s.workflows, id)
	
	return nil
}

// List returns all workflow summaries
func (s *InMemoryStore) List() []WorkflowSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	summaries := make([]WorkflowSummary, 0, len(s.workflows))
	
	for _, workflow := range s.workflows {
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

// Exists checks if a workflow exists
func (s *InMemoryStore) Exists(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	_, exists := s.workflows[id]
	return exists
}
