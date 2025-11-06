package httpclient

import (
	"fmt"
	"time"
)

// AuthType represents the type of authentication to use
type AuthType string

const (
	// AuthTypeNone represents no authentication (default)
	AuthTypeNone AuthType = "none"
	// AuthTypeBasic represents HTTP Basic Authentication
	AuthTypeBasic AuthType = "basic"
	// AuthTypeBearer represents Bearer Token Authentication
	AuthTypeBearer AuthType = "bearer"
	// AuthTypeAPIKey represents API Key Authentication
	AuthTypeAPIKey AuthType = "apikey"
)

// KeyValue represents a key-value pair for headers and query parameters.
// This structure allows duplicate keys, which is important for some HTTP scenarios.
type KeyValue struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// BasicAuthConfig contains HTTP Basic Authentication configuration
type BasicAuthConfig struct {
	Username string       `json:"username" yaml:"username"`
	Password SecureString `json:"password" yaml:"password"`
}

// TokenAuthConfig contains Bearer Token Authentication configuration
type TokenAuthConfig struct {
	Token SecureString `json:"token" yaml:"token"`
}

// APIKeyAuthConfig contains API Key Authentication configuration
type APIKeyAuthConfig struct {
	Key      string       `json:"key" yaml:"key"`           // Header or query parameter name
	Value    SecureString `json:"value" yaml:"value"`       // API key value
	Location string       `json:"location" yaml:"location"` // "header" or "query"
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	Type      AuthType          `json:"type,omitempty" yaml:"type,omitempty"`             // Authentication type (default: "none")
	BasicAuth *BasicAuthConfig  `json:"basic_auth,omitempty" yaml:"basic_auth,omitempty"` // Basic auth credentials
	Token     *TokenAuthConfig  `json:"token,omitempty" yaml:"token,omitempty"`           // Bearer token
	APIKey    *APIKeyAuthConfig `json:"api_key,omitempty" yaml:"api_key,omitempty"`       // API key configuration
}

// NetworkConfig contains network-level configuration
type NetworkConfig struct {
	Timeout             time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`                                 // Request timeout (default: 30s)
	MaxIdleConns        int           `json:"max_idle_conns,omitempty" yaml:"max_idle_conns,omitempty"`                   // Max idle connections (default: 100)
	MaxIdleConnsPerHost int           `json:"max_idle_conns_per_host,omitempty" yaml:"max_idle_conns_per_host,omitempty"` // Max idle conns per host (default: 10)
	MaxConnsPerHost     int           `json:"max_conns_per_host,omitempty" yaml:"max_conns_per_host,omitempty"`           // Max conns per host (default: 100)
	IdleConnTimeout     time.Duration `json:"idle_conn_timeout,omitempty" yaml:"idle_conn_timeout,omitempty"`             // Idle conn timeout (default: 90s)
	TLSHandshakeTimeout time.Duration `json:"tls_handshake_timeout,omitempty" yaml:"tls_handshake_timeout,omitempty"`     // TLS timeout (default: 10s)
	DisableKeepAlives   bool          `json:"disable_keep_alives,omitempty" yaml:"disable_keep_alives,omitempty"`         // Disable keep-alives (default: false)
}

// SecurityConfig contains security-related configuration including SSRF protection
// All protection is DENY BY DEFAULT (zero trust). Use Allow* fields to explicitly permit access.
type SecurityConfig struct {
	MaxRedirects    int   `json:"max_redirects,omitempty" yaml:"max_redirects,omitempty"`         // Max redirects (default: 10)
	MaxResponseSize int64 `json:"max_response_size,omitempty" yaml:"max_response_size,omitempty"` // Max response size in bytes (default: 10MB)
	FollowRedirects bool  `json:"follow_redirects,omitempty" yaml:"follow_redirects,omitempty"`   // Follow redirects (default: true)

	// SSRF Protection - DENY BY DEFAULT (zero trust security model)
	// Private IPs, localhost, link-local, and cloud metadata are BLOCKED by default
	// Set Allow* to true to explicitly permit access
	AllowPrivateIPs    bool     `json:"allow_private_ips,omitempty" yaml:"allow_private_ips,omitempty"`       // Allow private IP ranges (10.x, 172.16.x, 192.168.x) - default: false
	AllowLocalhost     bool     `json:"allow_localhost,omitempty" yaml:"allow_localhost,omitempty"`           // Allow localhost/loopback - default: false
	AllowLinkLocal     bool     `json:"allow_link_local,omitempty" yaml:"allow_link_local,omitempty"`         // Allow link-local addresses (169.254.x.x) - default: false
	AllowCloudMetadata bool     `json:"allow_cloud_metadata,omitempty" yaml:"allow_cloud_metadata,omitempty"` // Allow cloud metadata endpoints - default: false
	AllowedDomains     []string `json:"allowed_domains,omitempty" yaml:"allowed_domains,omitempty"`           // Whitelist of allowed domains (if set, only these domains are allowed)
}

// Config represents the configuration for an HTTP client.
// All fields are self-contained with no external dependencies.
type Config struct {
	// UID is the unique immutable identifier for this HTTP client
	UID string `json:"uid" yaml:"uid"`

	// Description provides human-readable documentation for this client
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Auth contains authentication configuration
	Auth AuthConfig `json:"auth,omitempty" yaml:"auth,omitempty"`

	// Network contains network-level configuration
	Network NetworkConfig `json:"network,omitempty" yaml:"network,omitempty"`

	// Security contains security-related configuration including SSRF protection
	Security SecurityConfig `json:"security,omitempty" yaml:"security,omitempty"`

	// Default headers to include in all requests (supports duplicate keys)
	Headers []KeyValue `json:"headers,omitempty" yaml:"headers,omitempty"`

	// Default query parameters to include in all requests (supports duplicate keys)
	QueryParams []KeyValue `json:"query_params,omitempty" yaml:"query_params,omitempty"`

	// BaseURL is the base URL for all requests (optional)
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty"`
}

// Validate checks if the client configuration is valid
func (c *Config) Validate() error {
	if c.UID == "" {
		return fmt.Errorf("client UID is required")
	}

	// Validate auth type
	if c.Auth.Type != "" && c.Auth.Type != AuthTypeNone && c.Auth.Type != AuthTypeBasic &&
		c.Auth.Type != AuthTypeBearer && c.Auth.Type != AuthTypeAPIKey {
		return fmt.Errorf("invalid auth_type: %s (must be one of: none, basic, bearer, apikey)", c.Auth.Type)
	}

	// Validate basic auth
	if c.Auth.Type == AuthTypeBasic {
		if c.Auth.BasicAuth == nil {
			return fmt.Errorf("basic_auth configuration is required for basic auth")
		}
		if c.Auth.BasicAuth.Username == "" {
			return fmt.Errorf("username is required for basic auth")
		}
		if c.Auth.BasicAuth.Password.IsEmpty() {
			return fmt.Errorf("password is required for basic auth")
		}
	}

	// Validate bearer token
	if c.Auth.Type == AuthTypeBearer {
		if c.Auth.Token == nil {
			return fmt.Errorf("token configuration is required for bearer auth")
		}
		if c.Auth.Token.Token.IsEmpty() {
			return fmt.Errorf("token is required for bearer auth")
		}
	}

	// Validate API key
	if c.Auth.Type == AuthTypeAPIKey {
		if c.Auth.APIKey == nil {
			return fmt.Errorf("api_key configuration is required for apikey auth")
		}
		if c.Auth.APIKey.Key == "" {
			return fmt.Errorf("api_key.key is required for apikey auth")
		}
		if c.Auth.APIKey.Value.IsEmpty() {
			return fmt.Errorf("api_key.value is required for apikey auth")
		}
		if c.Auth.APIKey.Location != "header" && c.Auth.APIKey.Location != "query" {
			return fmt.Errorf("api_key.location must be 'header' or 'query'")
		}
	}

	// Validate network settings
	if c.Network.Timeout < 0 {
		return fmt.Errorf("timeout cannot be negative")
	}
	if c.Network.MaxIdleConns < 0 {
		return fmt.Errorf("max_idle_conns cannot be negative")
	}
	if c.Network.MaxIdleConnsPerHost < 0 {
		return fmt.Errorf("max_idle_conns_per_host cannot be negative")
	}
	if c.Network.MaxConnsPerHost < 0 {
		return fmt.Errorf("max_conns_per_host cannot be negative")
	}
	if c.Network.IdleConnTimeout < 0 {
		return fmt.Errorf("idle_conn_timeout cannot be negative")
	}
	if c.Network.TLSHandshakeTimeout < 0 {
		return fmt.Errorf("tls_handshake_timeout cannot be negative")
	}

	// Validate security settings
	if c.Security.MaxRedirects < 0 {
		return fmt.Errorf("max_redirects cannot be negative")
	}
	if c.Security.MaxResponseSize < 0 {
		return fmt.Errorf("max_response_size cannot be negative")
	}

	return nil
}

// ApplyDefaults fills in default values for unset fields
func (c *Config) ApplyDefaults() {
	if c.Auth.Type == "" {
		c.Auth.Type = AuthTypeNone
	}

	if c.Network.Timeout == 0 {
		c.Network.Timeout = 30 * time.Second
	}

	if c.Network.MaxIdleConns == 0 {
		c.Network.MaxIdleConns = 100
	}

	if c.Network.MaxIdleConnsPerHost == 0 {
		c.Network.MaxIdleConnsPerHost = 10
	}

	if c.Network.MaxConnsPerHost == 0 {
		c.Network.MaxConnsPerHost = 100
	}

	if c.Network.IdleConnTimeout == 0 {
		c.Network.IdleConnTimeout = 90 * time.Second
	}

	if c.Network.TLSHandshakeTimeout == 0 {
		c.Network.TLSHandshakeTimeout = 10 * time.Second
	}

	if c.Security.MaxRedirects == 0 {
		c.Security.MaxRedirects = 10
	}

	if c.Security.MaxResponseSize == 0 {
		c.Security.MaxResponseSize = 10 * 1024 * 1024 // 10MB
	}

	// FollowRedirects defaults to true (handled in client creation)
	// Security is deny-by-default, so all Allow* fields default to false (no action needed)
}

// Clone creates a deep copy of the configuration
func (c *Config) Clone() *Config {
	clone := *c

	// Deep copy Auth config
	if c.Auth.BasicAuth != nil {
		basicAuth := *c.Auth.BasicAuth
		clone.Auth.BasicAuth = &basicAuth
	}
	if c.Auth.Token != nil {
		token := *c.Auth.Token
		clone.Auth.Token = &token
	}
	if c.Auth.APIKey != nil {
		apiKey := *c.Auth.APIKey
		clone.Auth.APIKey = &apiKey
	}

	// Deep copy slices
	if c.Security.AllowedDomains != nil {
		clone.Security.AllowedDomains = make([]string, len(c.Security.AllowedDomains))
		copy(clone.Security.AllowedDomains, c.Security.AllowedDomains)
	}

	if c.Headers != nil {
		clone.Headers = make([]KeyValue, len(c.Headers))
		copy(clone.Headers, c.Headers)
	}

	if c.QueryParams != nil {
		clone.QueryParams = make([]KeyValue, len(c.QueryParams))
		copy(clone.QueryParams, c.QueryParams)
	}

	return &clone
}
