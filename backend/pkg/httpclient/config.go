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
)

// KeyValue represents a key-value pair for headers and query parameters.
// This structure allows duplicate keys, which is important for some HTTP scenarios.
type KeyValue struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// Config represents the configuration for an HTTP client.
// All fields are self-contained with no external dependencies.
type Config struct {
	// UID is the unique immutable identifier for this HTTP client
	UID string `json:"uid" yaml:"uid"`

	// Description provides human-readable documentation for this client
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Authentication configuration
	AuthType AuthType `json:"auth_type,omitempty" yaml:"auth_type,omitempty"` // Default: "none"
	Username string   `json:"username,omitempty" yaml:"username,omitempty"`   // For basic auth
	Password string   `json:"password,omitempty" yaml:"password,omitempty"`   // For basic auth
	Token    string   `json:"token,omitempty" yaml:"token,omitempty"`         // For bearer token

	// Network configuration
	Timeout             time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`                                 // Request timeout (default: 30s)
	MaxIdleConns        int           `json:"max_idle_conns,omitempty" yaml:"max_idle_conns,omitempty"`                   // Max idle connections (default: 100)
	MaxIdleConnsPerHost int           `json:"max_idle_conns_per_host,omitempty" yaml:"max_idle_conns_per_host,omitempty"` // Max idle conns per host (default: 10)
	MaxConnsPerHost     int           `json:"max_conns_per_host,omitempty" yaml:"max_conns_per_host,omitempty"`           // Max conns per host (default: 100)
	IdleConnTimeout     time.Duration `json:"idle_conn_timeout,omitempty" yaml:"idle_conn_timeout,omitempty"`             // Idle conn timeout (default: 90s)
	TLSHandshakeTimeout time.Duration `json:"tls_handshake_timeout,omitempty" yaml:"tls_handshake_timeout,omitempty"`     // TLS timeout (default: 10s)
	DisableKeepAlives   bool          `json:"disable_keep_alives,omitempty" yaml:"disable_keep_alives,omitempty"`         // Disable keep-alives (default: false)

	// Security configuration
	MaxRedirects    int   `json:"max_redirects,omitempty" yaml:"max_redirects,omitempty"`         // Max redirects (default: 10)
	MaxResponseSize int64 `json:"max_response_size,omitempty" yaml:"max_response_size,omitempty"` // Max response size in bytes (default: 10MB)
	FollowRedirects bool  `json:"follow_redirects,omitempty" yaml:"follow_redirects,omitempty"`   // Follow redirects (default: true)

	// SSRF Protection
	BlockPrivateIPs    bool     `json:"block_private_ips,omitempty" yaml:"block_private_ips,omitempty"`       // Block private IP ranges
	BlockLocalhost     bool     `json:"block_localhost,omitempty" yaml:"block_localhost,omitempty"`           // Block localhost
	BlockLinkLocal     bool     `json:"block_link_local,omitempty" yaml:"block_link_local,omitempty"`         // Block link-local addresses
	BlockCloudMetadata bool     `json:"block_cloud_metadata,omitempty" yaml:"block_cloud_metadata,omitempty"` // Block cloud metadata endpoints
	AllowedDomains     []string `json:"allowed_domains,omitempty" yaml:"allowed_domains,omitempty"`           // Whitelist of allowed domains

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
	if c.AuthType != "" && c.AuthType != AuthTypeNone && c.AuthType != AuthTypeBasic && c.AuthType != AuthTypeBearer {
		return fmt.Errorf("invalid auth_type: %s (must be one of: none, basic, bearer)", c.AuthType)
	}

	// Validate basic auth
	if c.AuthType == AuthTypeBasic {
		if c.Username == "" {
			return fmt.Errorf("username is required for basic auth")
		}
		if c.Password == "" {
			return fmt.Errorf("password is required for basic auth")
		}
	}

	// Validate bearer token
	if c.AuthType == AuthTypeBearer {
		if c.Token == "" {
			return fmt.Errorf("token is required for bearer auth")
		}
	}

	// Validate network settings
	if c.Timeout < 0 {
		return fmt.Errorf("timeout cannot be negative")
	}
	if c.MaxIdleConns < 0 {
		return fmt.Errorf("max_idle_conns cannot be negative")
	}
	if c.MaxIdleConnsPerHost < 0 {
		return fmt.Errorf("max_idle_conns_per_host cannot be negative")
	}
	if c.MaxConnsPerHost < 0 {
		return fmt.Errorf("max_conns_per_host cannot be negative")
	}
	if c.IdleConnTimeout < 0 {
		return fmt.Errorf("idle_conn_timeout cannot be negative")
	}
	if c.TLSHandshakeTimeout < 0 {
		return fmt.Errorf("tls_handshake_timeout cannot be negative")
	}

	// Validate security settings
	if c.MaxRedirects < 0 {
		return fmt.Errorf("max_redirects cannot be negative")
	}
	if c.MaxResponseSize < 0 {
		return fmt.Errorf("max_response_size cannot be negative")
	}

	return nil
}

// ApplyDefaults fills in default values for unset fields
func (c *Config) ApplyDefaults() {
	if c.AuthType == "" {
		c.AuthType = AuthTypeNone
	}

	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 100
	}

	if c.MaxIdleConnsPerHost == 0 {
		c.MaxIdleConnsPerHost = 10
	}

	if c.MaxConnsPerHost == 0 {
		c.MaxConnsPerHost = 100
	}

	if c.IdleConnTimeout == 0 {
		c.IdleConnTimeout = 90 * time.Second
	}

	if c.TLSHandshakeTimeout == 0 {
		c.TLSHandshakeTimeout = 10 * time.Second
	}

	if c.MaxRedirects == 0 {
		c.MaxRedirects = 10
	}

	if c.MaxResponseSize == 0 {
		c.MaxResponseSize = 10 * 1024 * 1024 // 10MB
	}

	// FollowRedirects defaults to true (handled in client creation)
}

// Clone creates a deep copy of the configuration
func (c *Config) Clone() *Config {
	clone := *c

	// Deep copy slices
	if c.AllowedDomains != nil {
		clone.AllowedDomains = make([]string, len(c.AllowedDomains))
		copy(clone.AllowedDomains, c.AllowedDomains)
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
