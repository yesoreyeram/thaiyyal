package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/yesoreyeram/thaiyyal/backend"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
)

// SaveWorkflowRequest represents the request to save a workflow
type SaveWorkflowRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Data        json.RawMessage `json:"data"`
}

// SaveWorkflowResponse represents the response from saving a workflow
type SaveWorkflowResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// LoadWorkflowResponse represents the response from loading a workflow
type LoadWorkflowResponse struct {
	Success  bool                `json:"success"`
	Workflow *workflow.WorkflowMeta `json:"workflow,omitempty"`
	Error    string              `json:"error,omitempty"`
}

// ListWorkflowsResponse represents the response from listing workflows
type ListWorkflowsResponse struct {
	Success   bool                      `json:"success"`
	Workflows []workflow.WorkflowSummary `json:"workflows"`
	Count     int                       `json:"count"`
}

// DeleteWorkflowResponse represents the response from deleting a workflow
type DeleteWorkflowResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// handleSaveWorkflow handles saving a workflow
func (s *Server) handleSaveWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, s.config.MaxRequestBodySize)

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeErrorResponse(w, "Failed to read request body", http.StatusBadRequest, err)
		return
	}

	// Parse request
	var req SaveWorkflowRequest
	if err := json.Unmarshal(body, &req); err != nil {
		s.writeErrorResponse(w, "Failed to parse request", http.StatusBadRequest, err)
		return
	}

	// Save workflow
	id, err := s.workflowRegistry.Register(req.Name, req.Description, req.Data)
	if err != nil {
		s.writeJSONResponse(w, http.StatusBadRequest, SaveWorkflowResponse{
			Success: false,
			Error:   "Failed to save workflow: " + err.Error(),
		})
		return
	}

	s.logger.WithField("id", id).WithField("name", req.Name).Info("Workflow saved")

	// Write successful response
	s.writeJSONResponse(w, http.StatusCreated, SaveWorkflowResponse{
		Success: true,
		ID:      id,
		Message: "Workflow saved successfully",
	})
}

// handleLoadWorkflow handles loading a workflow by ID
func (s *Server) handleLoadWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	// Path format: /api/v1/workflow/load/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/workflow/load/")
	id := strings.TrimSpace(path)

	if id == "" {
		s.writeJSONResponse(w, http.StatusBadRequest, LoadWorkflowResponse{
			Success: false,
			Error:   "Workflow ID is required",
		})
		return
	}

	// Load workflow
	workflow, err := s.workflowRegistry.Get(id)
	if err != nil {
		s.writeJSONResponse(w, http.StatusNotFound, LoadWorkflowResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Write successful response
	s.writeJSONResponse(w, http.StatusOK, LoadWorkflowResponse{
		Success:  true,
		Workflow: workflow,
	})
}

// handleListWorkflows handles listing all workflows
func (s *Server) handleListWorkflows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get list of workflows
	workflows := s.workflowRegistry.List()

	// Write successful response
	s.writeJSONResponse(w, http.StatusOK, ListWorkflowsResponse{
		Success:   true,
		Workflows: workflows,
		Count:     len(workflows),
	})
}

// handleDeleteWorkflow handles deleting a workflow by ID
func (s *Server) handleDeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	// Path format: /api/v1/workflow/delete/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/workflow/delete/")
	id := strings.TrimSpace(path)

	if id == "" {
		s.writeJSONResponse(w, http.StatusBadRequest, DeleteWorkflowResponse{
			Success: false,
			Error:   "Workflow ID is required",
		})
		return
	}

	// Delete workflow
	err := s.workflowRegistry.Unregister(id)
	if err != nil {
		s.writeJSONResponse(w, http.StatusNotFound, DeleteWorkflowResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.logger.WithField("id", id).Info("Workflow deleted")

	// Write successful response
	s.writeJSONResponse(w, http.StatusOK, DeleteWorkflowResponse{
		Success: true,
		Message: "Workflow deleted successfully",
	})
}

// handleExecuteWorkflowByID handles executing a workflow by ID
func (s *Server) handleExecuteWorkflowByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	// Path format: /api/v1/workflow/execute/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/workflow/execute/")
	id := strings.TrimSpace(path)

	if id == "" {
		s.writeErrorResponse(w, "Workflow ID is required", http.StatusBadRequest, nil)
		return
	}

	// Load workflow
	workflow, err := s.workflowRegistry.Get(id)
	if err != nil {
		s.writeErrorResponse(w, "Failed to load workflow", http.StatusNotFound, err)
		return
	}

	// Execute workflow using the loaded data
	eng, err := engine.NewWithConfig(workflow.Data, s.engineConfig)
	if err != nil {
		s.writeErrorResponse(w, "Failed to create engine", http.StatusBadRequest, err)
		return
	}

	result, err := eng.Execute()
	if err != nil {
		s.writeErrorResponse(w, "Workflow execution failed", http.StatusInternalServerError, err)
		return
	}

	s.logger.WithField("id", id).WithField("name", workflow.Name).Info("Workflow executed by ID")

	// Write successful response
	s.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success":       true,
		"workflow_id":   id,
		"workflow_name": workflow.Name,
		"results":       result,
	})
}
