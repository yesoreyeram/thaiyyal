package security

import (
	"net"
	"testing"
)

// TestSSRFProtection_ValidateURL_AllowedURLs tests valid URLs
func TestSSRFProtection_ValidateURL_AllowedURLs(t *testing.T) {
	p := NewSSRFProtection()
	
	validURLs := []string{
		"https://example.com",
		"http://example.com",
		"https://api.example.com/data",
		"https://example.com:8080/path",
	}
	
	for _, urlStr := range validURLs {
		err := p.ValidateURL(urlStr)
		if err != nil {
			t.Errorf("URL should be valid: %s, error: %v", urlStr, err)
		}
	}
}

// TestSSRFProtection_ValidateURL_BlockedSchemes tests blocked URL schemes
func TestSSRFProtection_ValidateURL_BlockedSchemes(t *testing.T) {
	p := NewSSRFProtection()
	
	blockedURLs := []string{
		"ftp://example.com",
		"file:///etc/passwd",
		"gopher://example.com",
		"dict://example.com",
	}
	
	for _, urlStr := range blockedURLs {
		err := p.ValidateURL(urlStr)
		if err == nil {
			t.Errorf("URL should be blocked (scheme): %s", urlStr)
		}
	}
}

// TestSSRFProtection_ValidateURL_BlockedLocalhost tests localhost blocking
func TestSSRFProtection_ValidateURL_BlockedLocalhost(t *testing.T) {
	p := NewSSRFProtection()
	
	localhostURLs := []string{
		"http://localhost",
		"http://127.0.0.1",
		"http://127.0.0.1:8080",
		"http://[::1]",
		"http://0.0.0.0",
	}
	
	for _, urlStr := range localhostURLs {
		err := p.ValidateURL(urlStr)
		if err == nil {
			t.Errorf("URL should be blocked (localhost): %s", urlStr)
		}
	}
}

// TestSSRFProtection_ValidateURL_BlockedPrivateIPs tests private IP blocking
func TestSSRFProtection_ValidateURL_BlockedPrivateIPs(t *testing.T) {
	p := NewSSRFProtection()
	
	privateIPURLs := []string{
		"http://10.0.0.1",
		"http://10.255.255.255",
		"http://172.16.0.1",
		"http://172.31.255.255",
		"http://192.168.0.1",
		"http://192.168.255.255",
	}
	
	for _, urlStr := range privateIPURLs {
		err := p.ValidateURL(urlStr)
		if err == nil {
			t.Errorf("URL should be blocked (private IP): %s", urlStr)
		}
	}
}

// TestSSRFProtection_ValidateURL_BlockedLinkLocal tests link-local blocking
func TestSSRFProtection_ValidateURL_BlockedLinkLocal(t *testing.T) {
	p := NewSSRFProtection()
	
	linkLocalURLs := []string{
		"http://169.254.0.1",
		"http://169.254.255.255",
	}
	
	for _, urlStr := range linkLocalURLs {
		err := p.ValidateURL(urlStr)
		if err == nil {
			t.Errorf("URL should be blocked (link-local): %s", urlStr)
		}
	}
}

// TestSSRFProtection_ValidateURL_BlockedCloudMetadata tests cloud metadata blocking
func TestSSRFProtection_ValidateURL_BlockedCloudMetadata(t *testing.T) {
	p := NewSSRFProtection()
	
	metadataURLs := []string{
		"http://169.254.169.254",
		"http://169.254.169.254/latest/meta-data",
	}
	
	for _, urlStr := range metadataURLs {
		err := p.ValidateURL(urlStr)
		if err == nil {
			t.Errorf("URL should be blocked (cloud metadata): %s", urlStr)
		}
	}
}

// TestSSRFProtection_CustomConfig tests custom configuration
func TestSSRFProtection_CustomConfig(t *testing.T) {
	// Allow localhost but block everything else
	config := SSRFConfig{
		AllowedSchemes:     []string{"http", "https"},
		BlockPrivateIPs:    true,
		BlockLocalhost:     false, // Allow localhost
		BlockLinkLocal:     true,
		BlockCloudMetadata: true,
		AllowedDomains:     []string{},
		BlockedDomains:     []string{},
	}
	
	p := NewSSRFProtectionWithConfig(config)
	
	// Localhost should be allowed
	err := p.ValidateURL("http://localhost")
	if err != nil {
		t.Errorf("localhost should be allowed with custom config: %v", err)
	}
	
	// Private IPs should still be blocked
	err = p.ValidateURL("http://192.168.1.1")
	if err == nil {
		t.Error("private IPs should still be blocked")
	}
}

// TestSSRFProtection_DomainWhitelist tests domain whitelisting
func TestSSRFProtection_DomainWhitelist(t *testing.T) {
	config := SSRFConfig{
		AllowedSchemes:     []string{"http", "https"},
		BlockPrivateIPs:    true,
		BlockLocalhost:     true,
		BlockLinkLocal:     true,
		BlockCloudMetadata: true,
		AllowedDomains:     []string{"example.com", "api.example.com"},
		BlockedDomains:     []string{},
	}
	
	p := NewSSRFProtectionWithConfig(config)
	
	// Whitelisted domain should be allowed
	err := p.ValidateURL("https://example.com")
	if err != nil {
		t.Errorf("whitelisted domain should be allowed: %v", err)
	}
	
	// Non-whitelisted domain should be blocked
	err = p.ValidateURL("https://other.com")
	if err == nil {
		t.Error("non-whitelisted domain should be blocked")
	}
}

// TestSSRFProtection_DomainBlacklist tests domain blacklisting
func TestSSRFProtection_DomainBlacklist(t *testing.T) {
	config := SSRFConfig{
		AllowedSchemes:     []string{"http", "https"},
		BlockPrivateIPs:    false,
		BlockLocalhost:     false,
		BlockLinkLocal:     false,
		BlockCloudMetadata: false,
		AllowedDomains:     []string{},
		BlockedDomains:     []string{"evil.com", "malicious.com"},
	}
	
	p := NewSSRFProtectionWithConfig(config)
	
	// Blacklisted domain should be blocked
	err := p.ValidateURL("https://evil.com")
	if err == nil {
		t.Error("blacklisted domain should be blocked")
	}
	
	// Non-blacklisted domain should be allowed
	err = p.ValidateURL("https://example.com")
	if err != nil {
		t.Errorf("non-blacklisted domain should be allowed: %v", err)
	}
}

// TestSSRFProtection_InvalidURL tests invalid URLs
func TestSSRFProtection_InvalidURL(t *testing.T) {
	p := NewSSRFProtection()
	
	invalidURLs := []string{
		"",
		"not-a-url",
		"://missing-scheme",
		"http://",
	}
	
	for _, urlStr := range invalidURLs {
		err := p.ValidateURL(urlStr)
		if err == nil {
			t.Errorf("invalid URL should be rejected: %s", urlStr)
		}
	}
}

// TestIsLocalhost tests localhost detection
func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"127.0.0.1", true},
		{"127.0.0.2", true},
		{"::1", true},
		{"8.8.8.8", false},
		{"192.168.1.1", false},
	}
	
	for _, tt := range tests {
		ip := parseIP(t, tt.ip)
		result := isLocalhost(ip)
		if result != tt.expected {
			t.Errorf("isLocalhost(%s) = %v, want %v", tt.ip, result, tt.expected)
		}
	}
}

// TestIsPrivateIP tests private IP detection
func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"10.0.0.1", true},
		{"10.255.255.255", true},
		{"172.16.0.1", true},
		{"172.31.255.255", true},
		{"192.168.0.1", true},
		{"192.168.255.255", true},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"127.0.0.1", false}, // Loopback, not private
	}
	
	for _, tt := range tests {
		ip := parseIP(t, tt.ip)
		result := isPrivateIP(ip)
		if result != tt.expected {
			t.Errorf("isPrivateIP(%s) = %v, want %v", tt.ip, result, tt.expected)
		}
	}
}

// TestIsLinkLocal tests link-local detection
func TestIsLinkLocal(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"169.254.0.1", true},
		{"169.254.255.255", true},
		{"169.253.0.1", false},
		{"170.254.0.1", false},
		{"8.8.8.8", false},
	}
	
	for _, tt := range tests {
		ip := parseIP(t, tt.ip)
		result := isLinkLocal(ip)
		if result != tt.expected {
			t.Errorf("isLinkLocal(%s) = %v, want %v", tt.ip, result, tt.expected)
		}
	}
}

// TestIsCloudMetadata tests cloud metadata endpoint detection
func TestIsCloudMetadata(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"169.254.169.254", true},
		{"169.254.169.253", false},
		{"169.254.170.254", false},
		{"8.8.8.8", false},
	}
	
	for _, tt := range tests {
		ip := parseIP(t, tt.ip)
		result := isCloudMetadata(ip)
		if result != tt.expected {
			t.Errorf("isCloudMetadata(%s) = %v, want %v", tt.ip, result, tt.expected)
		}
	}
}

// Helper function to parse IP for tests
func parseIP(t *testing.T, ipStr string) net.IP {
	t.Helper()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		t.Fatalf("failed to parse IP: %s", ipStr)
	}
	return ip
}
