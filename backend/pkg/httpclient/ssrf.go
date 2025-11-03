package httpclient

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// validateURL validates a URL for SSRF protection based on configuration
func validateURL(urlStr string, config *Config) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Validate scheme
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("unsupported scheme: %s (only http and https are allowed)", scheme)
	}

	// Extract hostname
	hostname := parsedURL.Hostname()
	if hostname == "" {
		return fmt.Errorf("missing hostname in URL")
	}

	// Check domain whitelist
	if len(config.AllowedDomains) > 0 {
		allowed := false
		for _, domain := range config.AllowedDomains {
			if hostname == domain || strings.HasSuffix(hostname, "."+domain) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("domain %s is not in allowed domains list", hostname)
		}
	}

	// Resolve IP address
	ips, err := net.LookupIP(hostname)
	if err != nil {
		// If we can't resolve, we might be dealing with a domain
		// that will be resolved later. Allow it but apply domain checks.
		return nil
	}

	// Check each resolved IP
	for _, ip := range ips {
		if err := validateIP(ip, config); err != nil {
			return fmt.Errorf("IP validation failed for %s: %w", ip.String(), err)
		}
	}

	return nil
}

// validateIP validates an IP address for SSRF protection
func validateIP(ip net.IP, config *Config) error {
	// Block localhost
	if config.BlockLocalhost {
		if ip.IsLoopback() {
			return fmt.Errorf("localhost/loopback addresses are blocked")
		}
	}

	// Block link-local addresses (169.254.0.0/16)
	if config.BlockLinkLocal {
		if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("link-local addresses are blocked")
		}
	}

	// Block private IP ranges
	if config.BlockPrivateIPs {
		if isPrivateIP(ip) {
			return fmt.Errorf("private IP addresses are blocked")
		}
	}

	// Block cloud metadata endpoints
	if config.BlockCloudMetadata {
		if isCloudMetadataIP(ip) {
			return fmt.Errorf("cloud metadata endpoints are blocked")
		}
	}

	return nil
}

// isPrivateIP checks if an IP is in a private range
func isPrivateIP(ip net.IP) bool {
	// Private IPv4 ranges:
	// - 10.0.0.0/8
	// - 172.16.0.0/12
	// - 192.168.0.0/16
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7", // IPv6 private range
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// isCloudMetadataIP checks if an IP is a cloud metadata endpoint
func isCloudMetadataIP(ip net.IP) bool {
	// Common cloud metadata IPs:
	// - 169.254.169.254 (AWS, Azure, GCP)
	// - 100.100.100.200 (Alibaba Cloud)
	metadataIPs := []string{
		"169.254.169.254",
		"100.100.100.200",
		"fd00:ec2::254", // AWS IPv6 metadata
	}

	ipStr := ip.String()
	for _, metadataIP := range metadataIPs {
		if ipStr == metadataIP {
			return true
		}
	}

	return false
}
