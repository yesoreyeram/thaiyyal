package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with no auth",
			config: &Config{
				UID:      "test-client",
				AuthType: AuthTypeNone,
			},
			wantErr: false,
		},
		{
			name: "valid config with basic auth",
			config: &Config{
				UID:      "test-client",
				AuthType: AuthTypeBasic,
				Username: "user",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "valid config with bearer token",
			config: &Config{
				UID:      "test-client",
				AuthType: AuthTypeBearer,
				Token:    "token123",
			},
			wantErr: false,
		},
		{
			name: "missing UID",
			config: &Config{
				AuthType: AuthTypeNone,
			},
			wantErr: true,
			errMsg:  "client UID is required",
		},
		{
			name: "invalid auth type",
			config: &Config{
				UID:      "test-client",
				AuthType: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid auth_type",
		},
		{
			name: "basic auth missing username",
			config: &Config{
				UID:      "test-client",
				AuthType: AuthTypeBasic,
				Password: "pass",
			},
			wantErr: true,
			errMsg:  "username is required",
		},
		{
			name: "bearer auth missing token",
			config: &Config{
				UID:      "test-client",
				AuthType: AuthTypeBearer,
			},
			wantErr: true,
			errMsg:  "token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error, got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					if len(err.Error()) < len(tt.errMsg) || err.Error()[:len(tt.errMsg)] != tt.errMsg {
						t.Errorf("Validate() error = %v, want error containing %v", err, tt.errMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestConfig_ApplyDefaults(t *testing.T) {
	config := &Config{
		UID: "test-client",
	}

	config.ApplyDefaults()

	if config.AuthType != AuthTypeNone {
		t.Errorf("AuthType = %v, want %v", config.AuthType, AuthTypeNone)
	}
	if config.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want %v", config.Timeout, 30*time.Second)
	}
	if config.MaxIdleConns != 100 {
		t.Errorf("MaxIdleConns = %v, want 100", config.MaxIdleConns)
	}
	if config.MaxRedirects != 10 {
		t.Errorf("MaxRedirects = %v, want 10", config.MaxRedirects)
	}
	if config.MaxResponseSize != 10*1024*1024 {
		t.Errorf("MaxResponseSize = %v, want 10MB", config.MaxResponseSize)
	}
}

func TestConfig_Clone(t *testing.T) {
	original := &Config{
		UID:      "test-client",
		AuthType: AuthTypeBasic,
		Username: "user",
		Password: "pass",
		Headers: []KeyValue{
			{Key: "X-Custom", Value: "value"},
		},
		QueryParams: []KeyValue{
			{Key: "api_key", Value: "secret"},
		},
		AllowedDomains: []string{"example.com"},
	}

	clone := original.Clone()

	// Verify clone is equal
	if clone.UID != original.UID {
		t.Errorf("Clone UID = %v, want %v", clone.UID, original.UID)
	}

	// Verify deep copy of slices
	clone.Headers[0].Value = "modified"
	if original.Headers[0].Value == "modified" {
		t.Error("Clone modified original Headers")
	}

	clone.QueryParams[0].Value = "modified"
	if original.QueryParams[0].Value == "modified" {
		t.Error("Clone modified original QueryParams")
	}

	clone.AllowedDomains[0] = "modified.com"
	if original.AllowedDomains[0] == "modified.com" {
		t.Error("Clone modified original AllowedDomains")
	}
}

func TestNew_BasicAuth(t *testing.T) {
	// Create a test server that checks for basic auth
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("BasicAuth not found in request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != "testuser" || password != "testpass" {
			t.Errorf("BasicAuth = %v:%v, want testuser:testpass", username, password)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	}))
	defer server.Close()

	// Build client with basic auth
	config := &Config{
		UID:      "test-client",
		AuthType: AuthTypeBasic,
		Username: "testuser",
		Password: "testpass",
	}

	client, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Make request
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestNew_BearerToken(t *testing.T) {
	expectedToken := "test-token-123"

	// Create a test server that checks for bearer token
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		expected := "Bearer " + expectedToken
		if auth != expected {
			t.Errorf("Authorization header = %v, want %v", auth, expected)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	}))
	defer server.Close()

	// Build client with bearer token
	config := &Config{
		UID:      "test-client",
		AuthType: AuthTypeBearer,
		Token:    expectedToken,
	}

	client, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Make request
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestNew_DefaultHeaders(t *testing.T) {
	// Create a test server that checks for custom headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("X-Custom-Header = %v, want custom-value", r.Header.Get("X-Custom-Header"))
		}
		if r.Header.Get("User-Agent") != "TestAgent/1.0" {
			t.Errorf("User-Agent = %v, want TestAgent/1.0", r.Header.Get("User-Agent"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Build client with default headers
	config := &Config{
		UID: "test-client",
		Headers: []KeyValue{
			{Key: "X-Custom-Header", Value: "custom-value"},
			{Key: "User-Agent", Value: "TestAgent/1.0"},
		},
	}

	client, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Make request
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestNew_DefaultQueryParams(t *testing.T) {
	// Create a test server that checks for query parameters
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("api_key") != "secret123" {
			t.Errorf("api_key = %v, want secret123", r.URL.Query().Get("api_key"))
		}
		if r.URL.Query().Get("format") != "json" {
			t.Errorf("format = %v, want json", r.URL.Query().Get("format"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Build client with default query params
	config := &Config{
		UID: "test-client",
		QueryParams: []KeyValue{
			{Key: "api_key", Value: "secret123"},
			{Key: "format", Value: "json"},
		},
	}

	client, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Make request
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestNew_DuplicateHeaders(t *testing.T) {
	// Create a test server that checks for duplicate headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.Header.Values("X-Multi")
		if len(values) != 2 {
			t.Errorf("Expected 2 X-Multi headers, got %d", len(values))
		}
		if values[0] != "value1" || values[1] != "value2" {
			t.Errorf("X-Multi values = %v, want [value1, value2]", values)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Build client with duplicate headers
	config := &Config{
		UID: "test-client",
		Headers: []KeyValue{
			{Key: "X-Multi", Value: "value1"},
			{Key: "X-Multi", Value: "value2"},
		},
	}

	client, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Make request
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}
