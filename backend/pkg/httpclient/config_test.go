package httpclient

import (
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
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeNone,
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with basic auth",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeBasic,
					BasicAuth: &BasicAuthConfig{
						Username: "user",
						Password: NewSecureString("pass"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with bearer token",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeBearer,
					Token: &TokenAuthConfig{
						Token: NewSecureString("token123"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with API key in header",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeAPIKey,
					APIKey: &APIKeyAuthConfig{
						Key:      "X-API-Key",
						Value:    NewSecureString("secret-key"),
						Location: "header",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with API key in query",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeAPIKey,
					APIKey: &APIKeyAuthConfig{
						Key:      "api_key",
						Value:    NewSecureString("secret-key"),
						Location: "query",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing UID",
			config: &Config{
				Auth: AuthConfig{
					Type: AuthTypeNone,
				},
			},
			wantErr: true,
			errMsg:  "client UID is required",
		},
		{
			name: "invalid auth type",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "invalid auth_type",
		},
		{
			name: "basic auth missing username",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeBasic,
					BasicAuth: &BasicAuthConfig{
						Password: NewSecureString("pass"),
					},
				},
			},
			wantErr: true,
			errMsg:  "username is required",
		},
		{
			name: "basic auth missing password",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeBasic,
					BasicAuth: &BasicAuthConfig{
						Username: "user",
					},
				},
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "basic auth missing config",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeBasic,
				},
			},
			wantErr: true,
			errMsg:  "basic_auth configuration is required",
		},
		{
			name: "bearer auth missing token",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeBearer,
				},
			},
			wantErr: true,
			errMsg:  "token configuration is required",
		},
		{
			name: "api key missing key",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeAPIKey,
					APIKey: &APIKeyAuthConfig{
						Value:    NewSecureString("secret"),
						Location: "header",
					},
				},
			},
			wantErr: true,
			errMsg:  "api_key.key is required",
		},
		{
			name: "api key invalid location",
			config: &Config{
				UID: "test-client",
				Auth: AuthConfig{
					Type: AuthTypeAPIKey,
					APIKey: &APIKeyAuthConfig{
						Key:      "api_key",
						Value:    NewSecureString("secret"),
						Location: "invalid",
					},
				},
			},
			wantErr: true,
			errMsg:  "api_key.location must be",
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

	if config.Auth.Type != AuthTypeNone {
		t.Errorf("Auth.Type = %v, want %v", config.Auth.Type, AuthTypeNone)
	}
	if config.Network.Timeout != 30*time.Second {
		t.Errorf("Network.Timeout = %v, want %v", config.Network.Timeout, 30*time.Second)
	}
	if config.Network.MaxIdleConns != 100 {
		t.Errorf("Network.MaxIdleConns = %v, want 100", config.Network.MaxIdleConns)
	}
	if config.Security.MaxRedirects != 10 {
		t.Errorf("Security.MaxRedirects = %v, want 10", config.Security.MaxRedirects)
	}
	if config.Security.MaxResponseSize != 10*1024*1024 {
		t.Errorf("Security.MaxResponseSize = %v, want 10MB", config.Security.MaxResponseSize)
	}
}

func TestConfig_Clone(t *testing.T) {
	original := &Config{
		UID: "test-client",
		Auth: AuthConfig{
			Type: AuthTypeBasic,
			BasicAuth: &BasicAuthConfig{
				Username: "user",
				Password: NewSecureString("pass"),
			},
		},
		Headers: []KeyValue{
			{Key: "X-Custom", Value: "value"},
		},
		QueryParams: []KeyValue{
			{Key: "api_key", Value: "secret"},
		},
		Security: SecurityConfig{
			AllowedDomains: []string{"example.com"},
		},
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

	clone.Security.AllowedDomains[0] = "modified.com"
	if original.Security.AllowedDomains[0] == "modified.com" {
		t.Error("Clone modified original AllowedDomains")
	}

	// Verify deep copy of auth config
	if clone.Auth.BasicAuth != nil {
		clone.Auth.BasicAuth.Username = "modified"
		if original.Auth.BasicAuth.Username == "modified" {
			t.Error("Clone modified original BasicAuth")
		}
	}
}
