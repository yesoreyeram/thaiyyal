package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/httpclient"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestRegisterHTTPClient(t *testing.T) {
	// Create server
	config := DefaultConfig()
	engineConfig := types.DefaultConfig()
	srv, err := New(config, engineConfig)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	tests := []struct {
		name           string
		request        RegisterHTTPClientRequest
		expectedStatus int
		expectSuccess  bool
	}{
		{
			name: "Valid registration",
			request: RegisterHTTPClientRequest{
				Config: &httpclient.Config{
					UID:         "test-client-1",
					Description: "Test HTTP client",
				},
			},
			expectedStatus: http.StatusCreated,
			expectSuccess:  true,
		},
		{
			name: "Missing config",
			request: RegisterHTTPClientRequest{
				Config: nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
		{
			name: "Empty UID",
			request: RegisterHTTPClientRequest{
				Config: &httpclient.Config{
					UID:         "",
					Description: "Test client with empty UID",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			body, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/httpclient/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			srv.handleRegisterHTTPClient(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response
			var resp RegisterHTTPClientResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			// Check success field
			if resp.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v. Error: %s", tt.expectSuccess, resp.Success, resp.Error)
			}

			// If successful, check UID is returned
			if tt.expectSuccess && resp.UID == "" {
				t.Error("Expected UID in response, got empty string")
			}
		})
	}
}

func TestRegisterHTTPClient_DuplicateUID(t *testing.T) {
	// Create server
	config := DefaultConfig()
	engineConfig := types.DefaultConfig()
	srv, err := New(config, engineConfig)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Register first client
	req1 := RegisterHTTPClientRequest{
		Config: &httpclient.Config{
			UID:         "duplicate-client",
			Description: "First registration",
		},
	}

	body1, _ := json.Marshal(req1)
	httpReq1 := httptest.NewRequest(http.MethodPost, "/api/v1/httpclient/register", bytes.NewReader(body1))
	httpReq1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()
	srv.handleRegisterHTTPClient(rr1, httpReq1)

	if rr1.Code != http.StatusCreated {
		t.Fatalf("First registration failed with status %d", rr1.Code)
	}

	// Try to register again with same UID
	req2 := RegisterHTTPClientRequest{
		Config: &httpclient.Config{
			UID:         "duplicate-client",
			Description: "Second registration (should fail)",
		},
	}

	body2, _ := json.Marshal(req2)
	httpReq2 := httptest.NewRequest(http.MethodPost, "/api/v1/httpclient/register", bytes.NewReader(body2))
	httpReq2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	srv.handleRegisterHTTPClient(rr2, httpReq2)

	// Should get conflict status
	if rr2.Code != http.StatusConflict {
		t.Errorf("Expected status %d for duplicate registration, got %d", http.StatusConflict, rr2.Code)
	}

	var resp RegisterHTTPClientResponse
	json.Unmarshal(rr2.Body.Bytes(), &resp)
	if resp.Success {
		t.Error("Expected success=false for duplicate registration")
	}
}

func TestListHTTPClients(t *testing.T) {
	// Create server
	config := DefaultConfig()
	engineConfig := types.DefaultConfig()
	srv, err := New(config, engineConfig)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Initially should be empty
	req := httptest.NewRequest(http.MethodGet, "/api/v1/httpclient/list", nil)
	rr := httptest.NewRecorder()
	srv.handleListHTTPClients(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp ListHTTPClientsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !resp.Success {
		t.Error("Expected success=true")
	}

	if resp.Count != 0 {
		t.Errorf("Expected count=0, got %d", resp.Count)
	}

	// Register a few clients
	clients := []string{"client-1", "client-2", "client-3"}
	for _, uid := range clients {
		regReq := RegisterHTTPClientRequest{
			Config: &httpclient.Config{
				UID:         uid,
				Description: "Test client",
			},
		}
		body, _ := json.Marshal(regReq)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/httpclient/register", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		regRR := httptest.NewRecorder()
		srv.handleRegisterHTTPClient(regRR, httpReq)

		if regRR.Code != http.StatusCreated {
			t.Fatalf("Failed to register client %s: status %d", uid, regRR.Code)
		}
	}

	// List again
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/httpclient/list", nil)
	rr2 := httptest.NewRecorder()
	srv.handleListHTTPClients(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr2.Code)
	}

	var resp2 ListHTTPClientsResponse
	if err := json.Unmarshal(rr2.Body.Bytes(), &resp2); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !resp2.Success {
		t.Error("Expected success=true")
	}

	if resp2.Count != len(clients) {
		t.Errorf("Expected count=%d, got %d", len(clients), resp2.Count)
	}

	if len(resp2.Clients) != len(clients) {
		t.Errorf("Expected %d clients in list, got %d", len(clients), len(resp2.Clients))
	}

	// Verify all clients are in the list
	clientMap := make(map[string]bool)
	for _, c := range resp2.Clients {
		clientMap[c] = true
	}

	for _, expectedClient := range clients {
		if !clientMap[expectedClient] {
			t.Errorf("Expected client %s in list, but not found", expectedClient)
		}
	}
}

func TestListHTTPClients_MethodNotAllowed(t *testing.T) {
	// Create server
	config := DefaultConfig()
	engineConfig := types.DefaultConfig()
	srv, err := New(config, engineConfig)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Try POST method (should only allow GET)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/httpclient/list", nil)
	rr := httptest.NewRecorder()
	srv.handleListHTTPClients(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for POST, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestRegisterHTTPClient_MethodNotAllowed(t *testing.T) {
	// Create server
	config := DefaultConfig()
	engineConfig := types.DefaultConfig()
	srv, err := New(config, engineConfig)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Try GET method (should only allow POST)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/httpclient/register", nil)
	rr := httptest.NewRecorder()
	srv.handleRegisterHTTPClient(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for GET, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}
