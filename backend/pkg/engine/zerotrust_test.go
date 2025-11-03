package engine

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestZeroTrustConfig_HTTPDisabledByDefault tests that HTTP is disabled in zero trust mode
func TestZeroTrustConfig_HTTPDisabledByDefault(t *testing.T) {
	// Create a simple workflow with HTTP node
	payload := `{
		"nodes": [
			{"id": "http1", "type": "http", "data": {"url": "https://api.github.com"}}
		],
		"edges": []
	}`

	// Use zero trust config
	config := types.ZeroTrustConfig()

	engine, err := NewWithConfig([]byte(payload), config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Execute should fail because HTTP is disabled
	_, err = engine.Execute()
	if err == nil {
		t.Fatal("Expected error when HTTP is disabled, got nil")
	}

	// Check error message
	if !strings.Contains(err.Error(), "HTTP requests are not allowed") {
		t.Errorf("Expected 'HTTP requests are not allowed' error, got: %v", err)
	}
}

// TestZeroTrustConfig_HTTPEnabledWithWhitelist tests that HTTP works when explicitly enabled with whitelist
func TestZeroTrustConfig_HTTPEnabledWithWhitelist(t *testing.T) {
	// Create a simple workflow with HTTP node
	payload := `{
		"nodes": [
			{"id": "http1", "type": "http", "data": {"url": "https://httpbin.org/status/200"}}
		],
		"edges": []
	}`

	// Start with zero trust config and enable HTTP with whitelist
	config := types.ZeroTrustConfig()
	config.AllowHTTP = true
	config.AllowedDomains = []string{"httpbin.org"}
	config.MaxHTTPCallsPerExec = 1

	engine, err := NewWithConfig([]byte(payload), config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Execute should succeed because domain is whitelisted
	// Note: This will fail if no internet connection, but that's okay for this test
	result, err := engine.Execute()

	// We don't care about network errors, just that it got past the AllowHTTP check
	if err != nil && strings.Contains(err.Error(), "HTTP requests are not allowed") {
		t.Errorf("HTTP should be allowed when AllowHTTP=true, got: %v", err)
	}

	// If successful, verify result structure
	if err == nil && result != nil {
		if result.NodeResults == nil {
			t.Error("Expected NodeResults to be populated")
		}
	}
}

// TestZeroTrustConfig_DomainWhitelistBlocking tests that non-whitelisted domains are blocked
func TestZeroTrustConfig_DomainWhitelistBlocking(t *testing.T) {
	// Create a workflow with HTTP node to non-whitelisted domain
	payload := `{
		"nodes": [
			{"id": "http1", "type": "http", "data": {"url": "https://api.github.com"}}
		],
		"edges": []
	}`

	// Enable HTTP but whitelist only a different domain
	config := types.ZeroTrustConfig()
	config.AllowHTTP = true
	config.AllowedDomains = []string{"example.com"}
	config.MaxHTTPCallsPerExec = 1

	engine, err := NewWithConfig([]byte(payload), config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Execute should fail because domain is not whitelisted
	_, err = engine.Execute()
	if err == nil {
		t.Fatal("Expected error when accessing non-whitelisted domain, got nil")
	}

	// Check error message
	if !strings.Contains(err.Error(), "domain not in allowlist") {
		t.Errorf("Expected 'domain not in allowlist' error, got: %v", err)
	}
}

// TestZeroTrustConfig_PrivateIPsBlocked tests that private IPs are blocked
func TestZeroTrustConfig_PrivateIPsBlocked(t *testing.T) {
	testCases := []struct {
		name string
		url  string
	}{
		{"Private 10.x.x.x", "http://10.0.0.1"},
		{"Private 172.16.x.x", "http://172.16.0.1"},
		{"Private 192.168.x.x", "http://192.168.1.1"},
		{"Localhost IP", "http://127.0.0.1"},
		{"Localhost name", "http://localhost"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload := map[string]interface{}{
				"nodes": []map[string]interface{}{
					{
						"id":   "http1",
						"type": "http",
						"data": map[string]interface{}{"url": tc.url},
					},
				},
				"edges": []map[string]interface{}{},
			}

			payloadJSON, _ := json.Marshal(payload)

			// Use zero trust config with HTTP enabled (but private IPs blocked)
			config := types.ZeroTrustConfig()
			config.AllowHTTP = true
			config.AllowedDomains = []string{} // Allow all domains (but IPs still blocked)
			config.MaxHTTPCallsPerExec = 1

			engine, err := NewWithConfig(payloadJSON, config)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			// Execute should fail because private IPs are blocked
			_, err = engine.Execute()
			if err == nil {
				t.Errorf("Expected error when accessing %s, got nil", tc.url)
				return
			}

			// Check error message contains blocking reason
			errMsg := err.Error()
			validErrors := []string{
				"private IP addresses are blocked",
				"localhost addresses are blocked",
				"URL validation failed",
			}

			hasValidError := false
			for _, validErr := range validErrors {
				if strings.Contains(errMsg, validErr) {
					hasValidError = true
					break
				}
			}

			if !hasValidError {
				t.Errorf("Expected blocking error for %s, got: %v", tc.url, err)
			}
		})
	}
}

// TestZeroTrustConfig_CloudMetadataBlocked tests that cloud metadata endpoints are blocked
func TestZeroTrustConfig_CloudMetadataBlocked(t *testing.T) {
	testCases := []struct {
		name string
		url  string
	}{
		{"AWS/GCP/Azure metadata IP", "http://169.254.169.254"},
		{"AWS metadata IPv6", "http://[fd00:ec2::254]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload := map[string]interface{}{
				"nodes": []map[string]interface{}{
					{
						"id":   "http1",
						"type": "http",
						"data": map[string]interface{}{"url": tc.url},
					},
				},
				"edges": []map[string]interface{}{},
			}

			payloadJSON, _ := json.Marshal(payload)

			// Use zero trust config with HTTP enabled
			config := types.ZeroTrustConfig()
			config.AllowHTTP = true
			config.AllowedDomains = []string{}
			config.MaxHTTPCallsPerExec = 1

			engine, err := NewWithConfig(payloadJSON, config)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			// Execute should fail because cloud metadata is blocked
			_, err = engine.Execute()
			if err == nil {
				t.Errorf("Expected error when accessing %s, got nil", tc.url)
				return
			}

			// Check error message
			if !strings.Contains(err.Error(), "cloud metadata") &&
				!strings.Contains(err.Error(), "link-local") &&
				!strings.Contains(err.Error(), "private IP") {
				t.Errorf("Expected cloud metadata/link-local/private IP blocking error, got: %v", err)
			}
		})
	}
}

// TestZeroTrustConfig_ResourceLimits tests that zero trust has minimal resource limits
func TestZeroTrustConfig_ResourceLimits(t *testing.T) {
	config := types.ZeroTrustConfig()

	// Verify execution time limits are restrictive
	if config.MaxExecutionTime != 30*time.Second {
		t.Errorf("Expected MaxExecutionTime=30s, got %v", config.MaxExecutionTime)
	}

	// Verify iteration limits are restrictive
	if config.MaxIterations != 50 {
		t.Errorf("Expected MaxIterations=50, got %d", config.MaxIterations)
	}

	// Verify node limits are restrictive
	if config.MaxNodes != 50 {
		t.Errorf("Expected MaxNodes=50, got %d", config.MaxNodes)
	}

	// Verify string limits are restrictive
	if config.MaxStringLength != 50*1024 {
		t.Errorf("Expected MaxStringLength=50KB, got %d", config.MaxStringLength)
	}

	// Verify HTTP is disabled by default
	if config.AllowHTTP {
		t.Error("Expected AllowHTTP=false in zero trust config")
	}

	// Verify all security protections are enabled (DENY BY DEFAULT)
	// In zero trust mode, all Allow* should be false (meaning blocked)
	if config.AllowPrivateIPs {
		t.Error("Expected AllowPrivateIPs=false in zero trust config (blocked by default)")
	}
	if config.AllowLocalhost {
		t.Error("Expected AllowLocalhost=false in zero trust config (blocked by default)")
	}
	if config.AllowLinkLocal {
		t.Error("Expected AllowLinkLocal=false in zero trust config (blocked by default)")
	}
	if config.AllowCloudMetadata {
		t.Error("Expected AllowCloudMetadata=false in zero trust config (blocked by default)")
	}
}

// TestZeroTrustConfig_NoRetries tests that zero trust disables retries
func TestZeroTrustConfig_NoRetries(t *testing.T) {
	config := types.ZeroTrustConfig()

	if config.DefaultMaxAttempts != 1 {
		t.Errorf("Expected DefaultMaxAttempts=1 (no retries), got %d", config.DefaultMaxAttempts)
	}

	if config.DefaultBackoff != 0 {
		t.Errorf("Expected DefaultBackoff=0, got %v", config.DefaultBackoff)
	}
}

// TestZeroTrustConfig_ExceedsIterationLimit tests that iteration limits are enforced
func TestZeroTrustConfig_ExceedsIterationLimit(t *testing.T) {
	// Create a workflow with loop exceeding zero trust limit
	payload := `{
		"nodes": [
			{"id": "range", "type": "range", "data": {"start": 1, "end": 100}},
			{"id": "foreach", "type": "for_each", "data": {"max_iterations": 100}}
		],
		"edges": [
			{"source": "range", "target": "foreach"}
		]
	}`

	config := types.ZeroTrustConfig()

	engine, err := NewWithConfig([]byte(payload), config)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Execute should fail due to iteration limit (50 < 100)
	_, err = engine.Execute()
	if err == nil {
		t.Fatal("Expected error when exceeding iteration limit, got nil")
	}

	// Check error indicates limit was exceeded
	// Note: The actual error depends on implementation details
	// At minimum, verify execution didn't succeed
}

// TestDefaultConfig_ZeroTrustByDefault tests that default config uses zero-trust model
func TestDefaultConfig_ZeroTrustByDefault(t *testing.T) {
	config := types.DefaultConfig()

	// Verify HTTP is DISABLED by default (zero trust)
	if config.AllowHTTP {
		t.Error("Expected AllowHTTP=false in default config (zero trust by default)")
	}

	// Verify localhost is BLOCKED by default (AllowLocalhost=false means BLOCKED)
	if config.AllowLocalhost {
		t.Error("Expected AllowLocalhost=false in default config (zero trust by default - blocked)")
	}

	// Verify private IPs are BLOCKED by default (AllowPrivateIPs=false means BLOCKED)
	if config.AllowPrivateIPs {
		t.Error("Expected AllowPrivateIPs=false in default config (zero trust by default - blocked)")
	}

	// Verify cloud metadata is blocked
	if config.AllowCloudMetadata {
		t.Error("Expected AllowCloudMetadata=false in default config (blocked)")
	}

	// Verify link-local is blocked
	if config.AllowLinkLocal {
		t.Error("Expected AllowLinkLocal=false in default config (blocked)")
	}
}

// TestValidationLimits_HTTPEnabled tests that validation mode has HTTP enabled but restricted
func TestValidationLimits_HTTPEnabled(t *testing.T) {
	config := types.ValidationLimits()

	// Verify HTTP is allowed in validation mode (for testing workflows)
	if !config.AllowHTTP {
		t.Error("Expected AllowHTTP=true in validation config")
	}

	// But verify localhost is blocked in validation mode (AllowLocalhost=false means BLOCKED)
	if config.AllowLocalhost {
		t.Error("Expected AllowLocalhost=false in validation config for security (blocked)")
	}

	// Verify limits are restrictive
	if config.MaxHTTPCallsPerExec > 10 {
		t.Errorf("Expected restrictive MaxHTTPCallsPerExec in validation, got %d", config.MaxHTTPCallsPerExec)
	}
}

// TestDevelopmentConfig_Permissive tests that development mode is permissive
func TestDevelopmentConfig_Permissive(t *testing.T) {
	config := types.DevelopmentConfig()

	// Verify HTTP is allowed
	if !config.AllowHTTP {
		t.Error("Expected AllowHTTP=true in development config")
	}

	// Verify localhost is allowed (AllowLocalhost=true means ALLOWED)
	if !config.AllowLocalhost {
		t.Error("Expected AllowLocalhost=true in development config (allowed)")
	}

	// Verify private IPs are allowed (AllowPrivateIPs=true means ALLOWED)
	if !config.AllowPrivateIPs {
		t.Error("Expected AllowPrivateIPs=true in development config (allowed)")
	}

	// Verify limits are relaxed
	if config.MaxHTTPCallsPerExec < 1000 {
		t.Errorf("Expected relaxed MaxHTTPCallsPerExec in development, got %d", config.MaxHTTPCallsPerExec)
	}
}
