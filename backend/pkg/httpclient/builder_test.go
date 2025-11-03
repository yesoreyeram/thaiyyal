package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/config"
)

func TestClientConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *ClientConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with no auth",
			config: &ClientConfig{
				Name:     "test-client",
				AuthType: AuthTypeNone,
			},
			wantErr: false,
		},
		{
			name: "valid config with basic auth",
			config: &ClientConfig{
				Name:     "test-client",
				AuthType: AuthTypeBasic,
				Username: "user",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "valid config with bearer token",
			config: &ClientConfig{
				Name:     "test-client",
				AuthType: AuthTypeBearer,
				Token:    "token123",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: &ClientConfig{
				AuthType: AuthTypeNone,
			},
			wantErr: true,
			errMsg:  "client name is required",
		},
		{
			name: "invalid auth type",
			config: &ClientConfig{
				Name:     "test-client",
				AuthType: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid auth_type",
		},
		{
			name: "basic auth missing username",
			config: &ClientConfig{
				Name:     "test-client",
				AuthType: AuthTypeBasic,
				Password: "pass",
			},
			wantErr: true,
			errMsg:  "username is required",
		},
		{
			name: "basic auth missing password",
			config: &ClientConfig{
				Name:     "test-client",
				AuthType: AuthTypeBasic,
				Username: "user",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "bearer auth missing token",
			config: &ClientConfig{
				Name:     "test-client",
				AuthType: AuthTypeBearer,
			},
			wantErr: true,
			errMsg:  "token is required",
		},
		{
			name: "negative timeout",
			config: &ClientConfig{
				Name:     "test-client",
				Timeout:  -1 * time.Second,
			},
			wantErr: true,
			errMsg:  "timeout cannot be negative",
		},
		{
			name: "negative max redirects",
			config: &ClientConfig{
				Name:         "test-client",
				MaxRedirects: -1,
			},
			wantErr: true,
			errMsg:  "max_redirects cannot be negative",
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

func TestClientConfig_ApplyDefaults(t *testing.T) {
	config := &ClientConfig{
		Name: "test-client",
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

func TestClientConfig_Clone(t *testing.T) {
	original := &ClientConfig{
		Name:     "test-client",
		AuthType: AuthTypeBasic,
		Username: "user",
		Password: "pass",
		DefaultHeaders: map[string]string{
			"X-Custom": "value",
		},
		DefaultQueryParams: map[string]string{
			"api_key": "secret",
		},
	}

	clone := original.Clone()

	// Verify clone is equal
	if clone.Name != original.Name {
		t.Errorf("Clone Name = %v, want %v", clone.Name, original.Name)
	}
	if clone.Username != original.Username {
		t.Errorf("Clone Username = %v, want %v", clone.Username, original.Username)
	}

	// Verify deep copy of maps
	clone.DefaultHeaders["X-Custom"] = "modified"
	if original.DefaultHeaders["X-Custom"] == "modified" {
		t.Error("Clone modified original DefaultHeaders")
	}

	clone.DefaultQueryParams["api_key"] = "modified"
	if original.DefaultQueryParams["api_key"] == "modified" {
		t.Error("Clone modified original DefaultQueryParams")
	}
}

func TestBuilder_Build(t *testing.T) {
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)

	tests := []struct {
		name    string
		config  *ClientConfig
		wantErr bool
	}{
		{
			name: "build client with no auth",
			config: &ClientConfig{
				Name:     "test-client",
				AuthType: AuthTypeNone,
			},
			wantErr: false,
		},
		{
			name: "build client with basic auth",
			config: &ClientConfig{
				Name:     "basic-client",
				AuthType: AuthTypeBasic,
				Username: "user",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "build client with bearer token",
			config: &ClientConfig{
				Name:     "bearer-client",
				AuthType: AuthTypeBearer,
				Token:    "token123",
			},
			wantErr: false,
		},
		{
			name: "build with custom timeout",
			config: &ClientConfig{
				Name:    "timeout-client",
				Timeout: 60 * time.Second,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := builder.Build(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Build() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Build() unexpected error = %v", err)
				}
				if client == nil {
					t.Error("Build() returned nil client")
				}
				if client.GetConfig().Name != tt.config.Name {
					t.Errorf("Client config name = %v, want %v", client.GetConfig().Name, tt.config.Name)
				}
			}
		})
	}
}

func TestAuthTransport_BasicAuth(t *testing.T) {
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
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)
	clientConfig := &ClientConfig{
		Name:     "test-client",
		AuthType: AuthTypeBasic,
		Username: "testuser",
		Password: "testpass",
	}

	client, err := builder.Build(clientConfig)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
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

func TestAuthTransport_BearerToken(t *testing.T) {
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
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)
	clientConfig := &ClientConfig{
		Name:     "test-client",
		AuthType: AuthTypeBearer,
		Token:    expectedToken,
	}

	client, err := builder.Build(clientConfig)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
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

func TestAuthTransport_DefaultHeaders(t *testing.T) {
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
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)
	clientConfig := &ClientConfig{
		Name: "test-client",
		DefaultHeaders: map[string]string{
			"X-Custom-Header": "custom-value",
			"User-Agent":      "TestAgent/1.0",
		},
	}

	client, err := builder.Build(clientConfig)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
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

func TestAuthTransport_DefaultQueryParams(t *testing.T) {
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
	engineConfig := config.Testing()
	builder := NewBuilder(*engineConfig)
	clientConfig := &ClientConfig{
		Name: "test-client",
		DefaultQueryParams: map[string]string{
			"api_key": "secret123",
			"format":  "json",
		},
	}

	client, err := builder.Build(clientConfig)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
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

func TestBuilder_Redirects(t *testing.T) {
	t.Run("follow redirects", func(t *testing.T) {
		redirectCount := 0
		var server *httptest.Server
		
		// Create a test server that redirects
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if redirectCount < 3 {
				redirectCount++
				http.Redirect(w, r, server.URL, http.StatusFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("final"))
		}))
		defer server.Close()

		engineConfig := config.Testing()
		builder := NewBuilder(*engineConfig)
		clientConfig := &ClientConfig{
			Name:            "test-client",
			FollowRedirects: true,
			MaxRedirects:    10,
		}

		client, err := builder.Build(clientConfig)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
		}
	})

	t.Run("don't follow redirects", func(t *testing.T) {
		redirectCount := 0
		var server *httptest.Server
		
		// Create a test server that redirects
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if redirectCount < 3 {
				redirectCount++
				http.Redirect(w, r, server.URL, http.StatusFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("final"))
		}))
		defer server.Close()

		engineConfig := config.Testing()
		builder := NewBuilder(*engineConfig)
		clientConfig := &ClientConfig{
			Name:            "test-client",
			FollowRedirects: false,
		}

		client, err := builder.Build(clientConfig)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusFound {
			t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusFound)
		}
	})
}
