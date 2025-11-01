// Package security provides security utilities for the workflow engine.
// This includes SSRF protection, input validation, and other security features.
package security

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// SSRFProtection provides protection against Server-Side Request Forgery attacks
type SSRFProtection struct {
	allowedSchemes     map[string]bool
	blockPrivateIPs    bool
	blockLocalhost     bool
	blockLinkLocal     bool
	blockCloudMetadata bool
	allowedDomains     map[string]bool
	blockedDomains     map[string]bool
}

// SSRFConfig configures SSRF protection
type SSRFConfig struct {
	// AllowedSchemes lists allowed URL schemes (default: http, https)
	AllowedSchemes []string
	
	// BlockPrivateIPs blocks private IP ranges (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16)
	BlockPrivateIPs bool
	
	// BlockLocalhost blocks localhost and loopback addresses
	BlockLocalhost bool
	
	// BlockLinkLocal blocks link-local addresses (169.254.0.0/16)
	BlockLinkLocal bool
	
	// BlockCloudMetadata blocks cloud metadata endpoints (169.254.169.254, fd00:ec2::254)
	BlockCloudMetadata bool
	
	// AllowedDomains is a whitelist of allowed domains (empty = all allowed)
	AllowedDomains []string
	
	// BlockedDomains is a blacklist of blocked domains
	BlockedDomains []string
}

// DefaultSSRFConfig returns default SSRF protection configuration
func DefaultSSRFConfig() SSRFConfig {
	return SSRFConfig{
		AllowedSchemes:     []string{"http", "https"},
		BlockPrivateIPs:    true,
		BlockLocalhost:     true,
		BlockLinkLocal:     true,
		BlockCloudMetadata: true,
		AllowedDomains:     []string{},
		BlockedDomains:     []string{},
	}
}

// NewSSRFProtection creates a new SSRF protection instance with default config
func NewSSRFProtection() *SSRFProtection {
	return NewSSRFProtectionWithConfig(DefaultSSRFConfig())
}

// NewSSRFProtectionWithConfig creates a new SSRF protection instance with custom config
func NewSSRFProtectionWithConfig(config SSRFConfig) *SSRFProtection {
	p := &SSRFProtection{
		allowedSchemes:     make(map[string]bool),
		blockPrivateIPs:    config.BlockPrivateIPs,
		blockLocalhost:     config.BlockLocalhost,
		blockLinkLocal:     config.BlockLinkLocal,
		blockCloudMetadata: config.BlockCloudMetadata,
		allowedDomains:     make(map[string]bool),
		blockedDomains:     make(map[string]bool),
	}
	
	// Set allowed schemes
	if len(config.AllowedSchemes) == 0 {
		p.allowedSchemes["http"] = true
		p.allowedSchemes["https"] = true
	} else {
		for _, scheme := range config.AllowedSchemes {
			p.allowedSchemes[strings.ToLower(scheme)] = true
		}
	}
	
	// Set allowed domains
	for _, domain := range config.AllowedDomains {
		p.allowedDomains[strings.ToLower(domain)] = true
	}
	
	// Set blocked domains
	for _, domain := range config.BlockedDomains {
		p.blockedDomains[strings.ToLower(domain)] = true
	}
	
	return p
}

// ValidateURL validates a URL for SSRF protection
func (p *SSRFProtection) ValidateURL(urlStr string) error {
	// Parse URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	
	// Check scheme
	if !p.allowedSchemes[strings.ToLower(parsedURL.Scheme)] {
		return fmt.Errorf("URL scheme not allowed: %s (allowed: %v)", parsedURL.Scheme, p.getAllowedSchemesList())
	}
	
	// Get hostname
	hostname := parsedURL.Hostname()
	if hostname == "" {
		return fmt.Errorf("URL missing hostname")
	}
	
	// Check domain blocklist
	if p.blockedDomains[strings.ToLower(hostname)] {
		return fmt.Errorf("domain is blocked: %s", hostname)
	}
	
	// Check domain whitelist (if configured)
	if len(p.allowedDomains) > 0 && !p.allowedDomains[strings.ToLower(hostname)] {
		return fmt.Errorf("domain not in allowlist: %s", hostname)
	}
	
	// Try to parse as IP first
	ip := net.ParseIP(hostname)
	if ip != nil {
		// Hostname is an IP address
		if err := p.validateIP(ip); err != nil {
			return fmt.Errorf("IP validation failed for %s: %w", hostname, err)
		}
		return nil
	}
	
	// Check hostname-based validation first
	if err := p.validateHostname(hostname); err != nil {
		return err
	}
	
	// Try to resolve hostname to IP
	ips, err := net.LookupIP(hostname)
	if err != nil {
		// If we can't resolve, hostname validation already passed
		return nil
	}
	
	// Check each resolved IP
	for _, ip := range ips {
		if err := p.validateIP(ip); err != nil {
			return fmt.Errorf("IP validation failed for %s (%s): %w", hostname, ip, err)
		}
	}
	
	return nil
}

// validateIP validates an IP address
func (p *SSRFProtection) validateIP(ip net.IP) error {
	// Check localhost
	if p.blockLocalhost && isLocalhost(ip) {
		return fmt.Errorf("localhost addresses are blocked")
	}
	
	// Check private IPs
	if p.blockPrivateIPs && isPrivateIP(ip) {
		return fmt.Errorf("private IP addresses are blocked")
	}
	
	// Check link-local
	if p.blockLinkLocal && isLinkLocal(ip) {
		return fmt.Errorf("link-local addresses are blocked")
	}
	
	// Check cloud metadata
	if p.blockCloudMetadata && isCloudMetadata(ip) {
		return fmt.Errorf("cloud metadata endpoints are blocked")
	}
	
	return nil
}

// validateHostname validates a hostname (when IP resolution fails)
func (p *SSRFProtection) validateHostname(hostname string) error {
	hostname = strings.ToLower(hostname)
	
	// Check for localhost variations
	if p.blockLocalhost {
		localhostNames := []string{"localhost", "127.0.0.1", "::1", "0.0.0.0"}
		for _, localName := range localhostNames {
			if hostname == localName {
				return fmt.Errorf("localhost addresses are blocked")
			}
		}
	}
	
	// Check for cloud metadata hostnames
	if p.blockCloudMetadata {
		metadataHosts := []string{
			"169.254.169.254",
			"metadata.google.internal",
			"metadata.azure.com",
		}
		for _, metadataHost := range metadataHosts {
			if hostname == metadataHost {
				return fmt.Errorf("cloud metadata endpoints are blocked")
			}
		}
	}
	
	return nil
}

// getAllowedSchemesList returns list of allowed schemes for error messages
func (p *SSRFProtection) getAllowedSchemesList() []string {
	schemes := make([]string, 0, len(p.allowedSchemes))
	for scheme := range p.allowedSchemes {
		schemes = append(schemes, scheme)
	}
	return schemes
}

// isLocalhost checks if IP is localhost/loopback
func isLocalhost(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}
	
	// Also check for 0.0.0.0 (all interfaces - often treated as localhost)
	if ipv4 := ip.To4(); ipv4 != nil {
		if ipv4[0] == 0 && ipv4[1] == 0 && ipv4[2] == 0 && ipv4[3] == 0 {
			return true
		}
	}
	
	return false
}

// isPrivateIP checks if IP is in private ranges
func isPrivateIP(ip net.IP) bool {
	// Convert to 4-byte representation if IPv4
	if ipv4 := ip.To4(); ipv4 != nil {
		// 10.0.0.0/8
		if ipv4[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31 {
			return true
		}
		// 192.168.0.0/16
		if ipv4[0] == 192 && ipv4[1] == 168 {
			return true
		}
		return false
	}
	
	// Check for IPv6 private ranges (ULA: fc00::/7)
	if len(ip) == 16 && (ip[0]&0xfe) == 0xfc {
		return true
	}
	
	return false
}

// isLinkLocal checks if IP is link-local
func isLinkLocal(ip net.IP) bool {
	// Convert to 4-byte representation if IPv4
	if ipv4 := ip.To4(); ipv4 != nil {
		// IPv4 link-local: 169.254.0.0/16
		if ipv4[0] == 169 && ipv4[1] == 254 {
			return true
		}
		return false
	}
	
	// IPv6 link-local: fe80::/10
	if len(ip) == 16 && ip[0] == 0xfe && (ip[1]&0xc0) == 0x80 {
		return true
	}
	
	return ip.IsLinkLocalUnicast()
}

// isCloudMetadata checks if IP is a cloud metadata endpoint
func isCloudMetadata(ip net.IP) bool {
	// Convert to 4-byte representation if IPv4
	if ipv4 := ip.To4(); ipv4 != nil {
		// AWS, GCP, Azure metadata: 169.254.169.254
		if ipv4[0] == 169 && ipv4[1] == 254 && ipv4[2] == 169 && ipv4[3] == 254 {
			return true
		}
		return false
	}
	
	// AWS IMDSv2 (IPv6): fd00:ec2::254
	if len(ip) == 16 {
		if ip[0] == 0xfd && ip[1] == 0x00 && ip[2] == 0x0e && ip[3] == 0xc2 {
			// Check if it ends with ::254
			isZeros := true
			for i := 4; i < 14; i++ {
				if ip[i] != 0 {
					isZeros = false
					break
				}
			}
			if isZeros && ip[14] == 0x02 && ip[15] == 0x54 {
				return true
			}
		}
	}
	
	return false
}
