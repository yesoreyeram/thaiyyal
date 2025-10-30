package workflow

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ============================================================================
// HTTP Security Utilities
// ============================================================================
// This file contains security utilities for HTTP operations, including SSRF
// protection and URL validation.
// ============================================================================

// isInternalIP checks if a hostname resolves to an internal/private IP address.
// This prevents SSRF attacks against internal infrastructure.
//
// Blocked IP ranges:
//   - Loopback: 127.0.0.0/8, ::1
//   - Private: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, fc00::/7
//   - Link-local: 169.254.0.0/16, fe80::/10
//   - Multicast: 224.0.0.0/4, ff00::/8
func isInternalIP(hostname string) bool {
	// Resolve hostname to IP addresses
	ips, err := net.LookupIP(hostname)
	if err != nil {
		// If we can't resolve, block it to be safe
		return true
	}

	// Check each resolved IP address
	for _, ip := range ips {
		if isPrivateOrSpecialIP(ip) {
			return true
		}
	}

	return false
}

// isPrivateOrSpecialIP checks if an IP is private, loopback, link-local, or multicast
func isPrivateOrSpecialIP(ip net.IP) bool {
	// Check for loopback
	if ip.IsLoopback() {
		return true
	}

	// Check for link-local
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	// Check for multicast
	if ip.IsMulticast() {
		return true
	}

	// Check for private IPv4 ranges
	if ip4 := ip.To4(); ip4 != nil {
		// 10.0.0.0/8
		if ip4[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return true
		}
		// 192.168.0.0/16
		if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
		// 169.254.0.0/16 (AWS metadata, link-local)
		if ip4[0] == 169 && ip4[1] == 254 {
			return true
		}
	}

	// Check for private IPv6 ranges
	// fc00::/7 (Unique Local Addresses)
	if len(ip) == 16 && (ip[0]&0xfe) == 0xfc {
		return true
	}

	return false
}

// isAllowedURL validates a URL against security policies.
// It checks:
//   - URL is well-formed
//   - Scheme is http or https only
//   - Hostname doesn't resolve to internal IPs (if BlockInternalIPs is true)
//   - URL matches allowed patterns (if AllowedURLPatterns is non-empty)
func isAllowedURL(rawURL string, config Config) error {
	// Parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Check scheme (only allow http/https)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme '%s': only http and https are allowed", parsedURL.Scheme)
	}

	// Check for empty hostname
	if parsedURL.Hostname() == "" {
		return fmt.Errorf("URL must have a hostname")
	}

	// Block internal IPs if configured
	if config.BlockInternalIPs {
		if isInternalIP(parsedURL.Hostname()) {
			return fmt.Errorf("access to internal/private IP addresses is blocked for security reasons")
		}
	}

	// Check URL patterns whitelist (if configured)
	if len(config.AllowedURLPatterns) > 0 {
		allowed := false
		for _, pattern := range config.AllowedURLPatterns {
			// Simple pattern matching: check if hostname contains or matches the pattern
			// For production use, consider using regexp for more sophisticated matching
			if strings.Contains(parsedURL.Hostname(), pattern) || parsedURL.Hostname() == pattern {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("URL hostname '%s' is not in the allowed list", parsedURL.Hostname())
		}
	}

	return nil
}
