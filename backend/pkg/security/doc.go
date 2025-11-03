// Package security provides security controls and validation for workflow execution.
//
// # Overview
//
// The security package implements security measures to protect workflow execution
// from malicious inputs, resource exhaustion, and unauthorized access. It provides
// input validation, output sanitization, resource limits, and security policies.
//
// # Features
//
//   - Input validation: Sanitize and validate node inputs
//   - Output sanitization: Clean sensitive data from outputs
//   - Resource limits: Prevent resource exhaustion attacks
//   - URL validation: Whitelist/blacklist for HTTP nodes
//   - Expression safety: Prevent code injection in expressions
//   - Timeout enforcement: Limit execution time
//   - Rate limiting: Throttle execution requests
//   - Access control: Permission-based execution
//
// # Security Controls
//
// Input Validation:
//   - Type checking
//   - Size limits
//   - Format validation
//   - Injection prevention
//
// Resource Limits:
//   - Max execution time
//   - Max memory usage
//   - Max array size
//   - Max recursion depth
//
// Network Security:
//   - URL whitelisting
//   - URL blacklisting
//   - Protocol restrictions
//   - Request size limits
//
// Expression Security:
//   - No code execution
//   - Safe evaluation only
//   - Resource limits
//   - Injection prevention
//
// # Basic Usage
//
//	import "github.com/yesoreyeram/thaiyyal/backend/pkg/security"
//
//	// Create security validator
//	validator := security.NewValidator(security.Config{
//	    MaxExecutionTime: 5 * time.Minute,
//	    MaxArraySize: 10000,
//	    AllowedURLs: []string{"https://api.example.com/*"},
//	})
//
//	// Validate workflow before execution
//	if err := validator.ValidateWorkflow(workflow); err != nil {
//	    return fmt.Errorf("security validation failed: %w", err)
//	}
//
// # Input Validation
//
//	// Validate node inputs
//	if err := validator.ValidateInputs(node, inputs); err != nil {
//	    return fmt.Errorf("invalid inputs: %w", err)
//	}
//
// Validation rules:
//
//   - Strings: Max length, forbidden characters
//   - Numbers: Range validation
//   - Arrays: Max size, element validation
//   - Objects: Max depth, key validation
//   - URLs: Whitelist/blacklist checking
//
// # URL Validation
//
//	config := security.Config{
//	    AllowedURLs: []string{
//	        "https://api.example.com/*",
//	        "https://internal.company.com/api/*",
//	    },
//	    BlockedURLs: []string{
//	        "http://*",  // Block non-HTTPS
//	        "*/admin/*", // Block admin endpoints
//	    },
//	}
//
//	validator := security.NewValidator(config)
//
//	// Validate URL before HTTP request
//	if err := validator.ValidateURL(url); err != nil {
//	    return fmt.Errorf("URL not allowed: %w", err)
//	}
//
// # Resource Limits
//
//	config := security.Config{
//	    MaxExecutionTime: 5 * time.Minute,
//	    MaxArraySize: 10000,
//	    MaxObjectDepth: 10,
//	    MaxStringLength: 1000000,
//	    MaxMemoryMB: 512,
//	}
//
// Enforced limits:
//
//   - Execution time: Workflows timeout after max duration
//   - Array size: Reject arrays exceeding max size
//   - Object depth: Prevent deep nesting attacks
//   - String length: Limit string size
//   - Memory: Monitor and limit memory usage
//
// # Output Sanitization
//
//	// Remove sensitive data from outputs
//	sanitized := validator.SanitizeOutput(output, security.SanitizeConfig{
//	    RemoveFields: []string{"password", "secret", "token"},
//	    MaskFields: []string{"email", "phone"},
//	})
//
// # Rate Limiting
//
//	// Create rate limiter
//	limiter := security.NewRateLimiter(security.RateLimitConfig{
//	    RequestsPerSecond: 100,
//	    BurstSize: 200,
//	})
//
//	// Check before execution
//	if !limiter.Allow(userID) {
//	    return errors.New("rate limit exceeded")
//	}
//
// # Access Control
//
//	// Define permissions
//	permissions := security.Permissions{
//	    CanExecute: true,
//	    AllowedNodeTypes: []types.NodeType{
//	        types.NodeTypeNumber,
//	        types.NodeTypeOperation,
//	    },
//	    DeniedNodeTypes: []types.NodeType{
//	        types.NodeTypeHTTP, // Restrict HTTP access
//	    },
//	}
//
//	// Check permissions
//	if err := validator.CheckPermissions(user, workflow, permissions); err != nil {
//	    return fmt.Errorf("permission denied: %w", err)
//	}
//
// # Security Policies
//
// Policies can be defined and enforced:
//
//	policy := security.Policy{
//	    Name: "Production Policy",
//	    Rules: []security.Rule{
//	        {Type: "max_execution_time", Value: "5m"},
//	        {Type: "require_https", Value: "true"},
//	        {Type: "block_external_urls", Value: "true"},
//	    },
//	}
//
//	enforcer := security.NewPolicyEnforcer(policy)
//	if err := enforcer.Enforce(workflow); err != nil {
//	    return fmt.Errorf("policy violation: %w", err)
//	}
//
// # Threat Protection
//
// Protection against common threats:
//
//   - Injection attacks: Expression validation, input sanitization
//   - DoS attacks: Resource limits, rate limiting, timeout enforcement
//   - SSRF attacks: URL validation, network restrictions
//   - Data exfiltration: Output sanitization, audit logging
//   - Privilege escalation: Permission checks, access control
//
// # Audit Logging
//
// Security events are logged for audit:
//
//   - Validation failures
//   - Permission denials
//   - Rate limit violations
//   - Suspicious activities
//   - Policy violations
//
// # Integration
//
// Security controls integrate with the engine:
//
//	engine := engine.New(
//	    engine.WithValidator(validator),
//	    engine.WithRateLimiter(limiter),
//	    engine.WithAccessControl(permissions),
//	)
//
// # Performance Impact
//
// Security checks have minimal overhead:
//
//   - Validation: O(n) where n is input size
//   - URL checking: O(1) with efficient pattern matching
//   - Rate limiting: O(1) with token bucket algorithm
//   - Permission checks: O(1) with caching
//
// # Best Practices
//
//   - Always validate inputs before execution
//   - Use allowlists instead of blocklists for URLs
//   - Set reasonable resource limits
//   - Enable audit logging in production
//   - Regularly review security policies
//   - Test security controls thoroughly
//   - Keep security configuration external
//   - Follow principle of least privilege
//
// # Compliance
//
// The security package helps meet compliance requirements:
//
//   - OWASP Top 10 protection
//   - Input validation requirements
//   - Audit logging requirements
//   - Access control requirements
//   - Data protection requirements
//
// # Thread Safety
//
// Security validators and limiters are thread-safe and can be used
// concurrently from multiple goroutines.
package security
