// Package server provides HTTP API routes for the Thaiyyal workflow engine,
// including HTTP client management endpoints.
package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/httpclient"
)

// RegisterHTTPClientRequest represents the request body for registering an HTTP client
type RegisterHTTPClientRequest struct {
	Config *httpclient.Config `json:"config"`
}

// RegisterHTTPClientResponse represents the response for registering an HTTP client
type RegisterHTTPClientResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	UID     string `json:"uid,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ListHTTPClientsResponse represents the response for listing HTTP clients
type ListHTTPClientsResponse struct {
	Success bool     `json:"success"`
	Clients []string `json:"clients"`
	Count   int      `json:"count"`
}

// handleRegisterHTTPClient handles HTTP client registration requests
func (s *Server) handleRegisterHTTPClient(w http.ResponseWriter, r *http.Request) {
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
	var req RegisterHTTPClientRequest
	if err := json.Unmarshal(body, &req); err != nil {
		s.writeErrorResponse(w, "Failed to parse request", http.StatusBadRequest, err)
		return
	}

	// Validate config
	if req.Config == nil {
		s.writeJSONResponse(w, http.StatusBadRequest, RegisterHTTPClientResponse{
			Success: false,
			Error:   "config is required",
		})
		return
	}

	// Create HTTP client from config
	client, err := httpclient.New(context.Background(), req.Config)
	if err != nil {
		s.writeJSONResponse(w, http.StatusBadRequest, RegisterHTTPClientResponse{
			Success: false,
			Error:   "Failed to create HTTP client: " + err.Error(),
		})
		return
	}

	// Register the client
	if err := s.httpClientRegistry.Register(req.Config.UID, client); err != nil {
		s.writeJSONResponse(w, http.StatusConflict, RegisterHTTPClientResponse{
			Success: false,
			Error:   "Failed to register HTTP client: " + err.Error(),
		})
		return
	}

	s.logger.WithField("uid", req.Config.UID).Info("HTTP client registered")

	// Write successful response
	s.writeJSONResponse(w, http.StatusCreated, RegisterHTTPClientResponse{
		Success: true,
		Message: "HTTP client registered successfully",
		UID:     req.Config.UID,
	})
}

// handleListHTTPClients handles listing HTTP clients requests
func (s *Server) handleListHTTPClients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get list of registered clients
	clients := s.httpClientRegistry.List()

	// Write successful response
	s.writeJSONResponse(w, http.StatusOK, ListHTTPClientsResponse{
		Success: true,
		Clients: clients,
		Count:   len(clients),
	})
}
