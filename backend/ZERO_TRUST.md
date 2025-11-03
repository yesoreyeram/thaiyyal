# Zero Trust Security Model

## Overview

Thaiyyal workflow engine implements a **zero trust / zero permission security model** where all privileged operations are **DENIED by default** and require **explicit opt-in** through configuration.

This document describes the security model, how to use it, and how to configure workflows securely.

## Zero Trust Principles

### 1. Deny All by Default

The workflow engine operates on a **deny-all** foundation:

- ❌ **No network access** - HTTP requests are disabled by default
- ❌ **No environment variables** - Workflow engine has no access to environment variables
- ❌ **No file system access** - Workflow engine has no file system capabilities
- ❌ **No external commands** - Cannot execute system commands
- ❌ **No code execution** - Cannot execute arbitrary code (only predefined node types)

### 2. Explicit Opt-In Required

To enable any privileged operation, you must **explicitly configure it**:

```go
config := types.DefaultConfig()

// MUST explicitly enable HTTP
config.AllowHTTP = true

// MUST explicitly whitelist domains (recommended)
config.AllowedDomains = []string{
    "api.trusted.com",
    "api.example.org",
}

engine, err := engine.NewWithConfig(payload, config)
```

### 3. Least Privilege

Enable only the minimum capabilities needed:

```go
// ✅ GOOD: Minimal, specific permissions
config.AllowHTTP = true
config.AllowedDomains = []string{"api.myservice.com"}
config.MaxHTTPCallsPerExec = 10

// ❌ BAD: Too permissive
config.AllowHTTP = true
config.AllowedDomains = []string{}  // Empty = allow all domains
config.MaxHTTPCallsPerExec = 0      // 0 = unlimited
```

### 4. Defense in Depth

Multiple layers of security protection:

- **Network isolation** (AllowHTTP flag)
- **Domain whitelisting** (AllowedDomains)
- **IP blocking** (BlockPrivateIPs, BlockLocalhost, etc.)
- **Rate limiting** (MaxHTTPCallsPerExec)
- **Size limits** (MaxResponseSize)
- **Timeout limits** (HTTPTimeout)
- **SSRF protection** (blocks metadata endpoints, private IPs)

## Configuration Presets

Thaiyyal provides three configuration presets for different use cases:

### DefaultConfig() - Zero Trust Production

**Use for**: Production deployments, untrusted workflows, sandboxed execution

```go
config := types.DefaultConfig()
// HTTP: DISABLED by default
// Localhost: BLOCKED
// Private IPs: BLOCKED
// Resource limits: Reasonable for production
```

**Key Settings**:
- `AllowHTTP = false` - No network access
- `BlockLocalhost = true` - Block localhost
- `BlockPrivateIPs = true` - Block private IP ranges
- `BlockCloudMetadata = true` - Block cloud metadata endpoints
- `MaxHTTPCallsPerExec = 100` - Limited calls (when HTTP enabled)
- `MaxExecutionTime = 5 minutes` - Reasonable timeout

### ZeroTrustConfig() - Maximum Security

**Use for**: High-security environments, untrusted code, public sandboxes

```go
config := types.ZeroTrustConfig()
// HTTP: DISABLED
// All protections: MAXIMUM
// Resource limits: MINIMAL
```

**Key Settings**:
- `AllowHTTP = false` - No network access
- `MaxExecutionTime = 30 seconds` - Short timeout
- `MaxNodes = 50` - Minimal workflow size
- `MaxHTTPCallsPerExec = 0` - No HTTP calls allowed
- `MaxIterations = 50` - Minimal iterations
- All security blocks enabled

### DevelopmentConfig() - Relaxed for Development

**Use for**: Local development, testing, debugging

```go
config := types.DevelopmentConfig()
// HTTP: ENABLED
// Localhost: ALLOWED
// Resource limits: RELAXED
```

**Key Settings**:
- `AllowHTTP = true` - HTTP enabled
- `BlockLocalhost = false` - Allow localhost
- `BlockPrivateIPs = false` - Allow private IPs
- `MaxHTTPCallsPerExec = 1000` - Higher limits
- `MaxExecutionTime = 30 minutes` - Long timeout

## Secure Configuration Examples

### Example 1: No Network Access (Default)

```go
// Most secure - no network access at all
config := types.DefaultConfig()
engine, err := engine.NewWithConfig(payload, config)

// Workflows can:
// ✅ Process data (math, text operations)
// ✅ Use control flow (loops, conditions)
// ✅ Store state (variables, cache)
// ❌ Make HTTP requests (disabled)
```

### Example 2: Whitelisted API Access

```go
// Allow HTTP only to specific trusted APIs
config := types.DefaultConfig()
config.AllowHTTP = true
config.AllowedDomains = []string{
    "api.github.com",
    "api.stripe.com",
}
config.MaxHTTPCallsPerExec = 20

engine, err := engine.NewWithConfig(payload, config)

// Workflows can:
// ✅ Call api.github.com
// ✅ Call api.stripe.com
// ❌ Call any other domain (blocked)
// ❌ Call localhost (blocked)
// ❌ Call private IPs (blocked)
```

### Example 3: Development/Testing

```go
// Relaxed for local development
config := types.DevelopmentConfig()

// Or customize from default:
config := types.DefaultConfig()
config.AllowHTTP = true
config.BlockLocalhost = false  // Allow localhost for testing
config.AllowedDomains = []string{}  // Allow all domains (development only!)

engine, err := engine.NewWithConfig(payload, config)
```

### Example 4: Maximum Security Sandbox

```go
// Ultra-restrictive for untrusted workflows
config := types.ZeroTrustConfig()

engine, err := engine.NewWithConfig(payload, config)

// Workflows can:
// ✅ Basic data processing only
// ❌ No network access
// ❌ Minimal resource usage
// ❌ Short execution time
```

## Security Features

### 1. Network Access Control

#### AllowHTTP Flag

The master switch for all network access:

```go
config.AllowHTTP = false  // DISABLED by default
```

When `false`:
- All HTTP nodes fail with error
- No network requests are made
- Workflows are completely isolated

When `true`:
- HTTP nodes can execute (subject to other restrictions)
- Domain whitelist applies (if configured)
- IP blocking applies
- Rate limits apply

#### Domain Whitelisting

Control exactly which domains workflows can access:

```go
// Empty = allow all domains (when AllowHTTP is true)
config.AllowedDomains = []string{}

// Whitelist specific domains (RECOMMENDED)
config.AllowedDomains = []string{
    "api.trusted.com",
    "api.example.org",
}
```

When configured:
- Only whitelisted domains are allowed
- All other domains are blocked
- Subdomains must be explicitly listed

### 2. SSRF Protection

Multiple layers protect against Server-Side Request Forgery:

#### Block Private IPs

```go
config.BlockPrivateIPs = true  // ENABLED by default
```

Blocks:
- `10.0.0.0/8` (Class A private)
- `172.16.0.0/12` (Class B private)
- `192.168.0.0/16` (Class C private)
- IPv6 ULA (`fc00::/7`)

#### Block Localhost

```go
config.BlockLocalhost = true  // ENABLED by default
```

Blocks:
- `127.0.0.1`, `::1` (loopback)
- `0.0.0.0` (all interfaces)
- `localhost` (hostname)

#### Block Link-Local

```go
config.BlockLinkLocal = true  // ENABLED by default
```

Blocks:
- `169.254.0.0/16` (IPv4 link-local)
- `fe80::/10` (IPv6 link-local)

#### Block Cloud Metadata

```go
config.BlockCloudMetadata = true  // ENABLED by default
```

Blocks:
- `169.254.169.254` (AWS, GCP, Azure)
- `fd00:ec2::254` (AWS IMDSv2 IPv6)
- `metadata.google.internal`
- `metadata.azure.com`

### 3. Resource Limits

Prevent resource exhaustion and DoS attacks:

```go
config.MaxHTTPCallsPerExec = 100      // Max HTTP calls per workflow
config.MaxResponseSize = 10 * 1024 * 1024  // 10MB response limit
config.HTTPTimeout = 30 * time.Second      // 30s per request
config.MaxHTTPRedirects = 10               // Limit redirect chains
```

See [PROTECTION.md](PROTECTION.md) for complete resource limits.

### 4. No Environment Variable Access

The workflow engine has **zero access** to environment variables:

- ✅ Workflows cannot read `os.Getenv()`
- ✅ No `$ENV_VAR` substitution
- ✅ Environment is completely isolated

All configuration must be passed explicitly through:
- Workflow JSON payload
- Engine configuration
- Context variables

### 5. No File System Access

The workflow engine has **zero access** to the file system:

- ✅ Cannot read files
- ✅ Cannot write files
- ✅ Cannot list directories
- ✅ Cannot execute binaries

All data must flow through:
- Workflow node inputs/outputs
- HTTP responses (if enabled)
- In-memory state

## Migration Guide

### Breaking Changes from Previous Version

**Network Access**: HTTP is now **DISABLED by default**

**Before** (old behavior):
```go
config := types.DefaultConfig()
engine, _ := engine.NewWithConfig(payload, config)
// HTTP was allowed
```

**After** (new behavior):
```go
config := types.DefaultConfig()
// HTTP is now DISABLED by default
// Must explicitly enable:
config.AllowHTTP = true
config.AllowedDomains = []string{"api.example.com"}
engine, _ := engine.NewWithConfig(payload, config)
```

**Localhost Blocking**: Localhost is now **BLOCKED by default**

**Before**:
```go
config := types.DefaultConfig()
// Localhost was allowed
```

**After**:
```go
config := types.DefaultConfig()
// Localhost is now blocked
// For development/testing only:
config.BlockLocalhost = false
```

### Migration Steps

1. **Audit your workflows**: Identify which workflows use HTTP nodes

2. **Enable HTTP explicitly**: Update your code to enable HTTP:
   ```go
   config := types.DefaultConfig()
   config.AllowHTTP = true
   ```

3. **Whitelist domains** (recommended): Specify allowed domains:
   ```go
   config.AllowedDomains = []string{
       "api.service1.com",
       "api.service2.com",
   }
   ```

4. **Development environments**: Use DevelopmentConfig:
   ```go
   config := types.DevelopmentConfig()  // Relaxed for development
   ```

5. **Test thoroughly**: Verify workflows work with new security model

## Best Practices

### 1. Use Least Privilege

✅ **DO**: Enable only what you need
```go
config := types.DefaultConfig()
config.AllowHTTP = true
config.AllowedDomains = []string{"api.myservice.com"}
config.MaxHTTPCallsPerExec = 5
```

❌ **DON'T**: Use permissive settings in production
```go
config.AllowHTTP = true
config.AllowedDomains = []string{}  // Allows all domains!
config.MaxHTTPCallsPerExec = 0      // Unlimited!
```

### 2. Separate Development and Production

```go
// In your application
var config types.Config
if os.Getenv("ENV") == "production" {
    config = types.DefaultConfig()  // Strict
    config.AllowHTTP = true
    config.AllowedDomains = loadProductionDomains()
} else {
    config = types.DevelopmentConfig()  // Relaxed
}
```

### 3. Validate Untrusted Workflows First

```go
// Step 1: Validate with strict limits
validationEngine, _ := engine.NewWithConfig(payload, types.ZeroTrustConfig())
_, err := validationEngine.Execute()
if err != nil {
    return fmt.Errorf("workflow validation failed: %w", err)
}

// Step 2: Execute with production config
productionConfig := types.DefaultConfig()
productionConfig.AllowHTTP = true
productionConfig.AllowedDomains = allowedDomains
productionEngine, _ := engine.NewWithConfig(payload, productionConfig)
result, err := productionEngine.Execute()
```

### 4. Monitor and Log

```go
result, err := engine.Execute()

// Log security-relevant metrics
log.Printf("Workflow executed: nodes=%d, http_calls=%d, duration=%v",
    engine.GetNodeExecutionCount(),
    engine.GetHTTPCallCount(),
    time.Since(startTime),
)
```

### 5. Defense in Depth

Use multiple security layers:

```go
config := types.DefaultConfig()

// Layer 1: Enable HTTP only if needed
config.AllowHTTP = true

// Layer 2: Whitelist specific domains
config.AllowedDomains = []string{"api.trusted.com"}

// Layer 3: Block dangerous IPs
config.BlockPrivateIPs = true
config.BlockLocalhost = true
config.BlockCloudMetadata = true

// Layer 4: Rate limiting
config.MaxHTTPCallsPerExec = 20

// Layer 5: Size limits
config.MaxResponseSize = 5 * 1024 * 1024  // 5MB

// Layer 6: Time limits
config.HTTPTimeout = 10 * time.Second
config.MaxExecutionTime = 2 * time.Minute
```

## Security Checklist

Before deploying to production:

- [ ] HTTP is disabled OR explicitly enabled with justification
- [ ] If HTTP enabled, AllowedDomains is configured (not empty)
- [ ] BlockPrivateIPs = true
- [ ] BlockLocalhost = true
- [ ] BlockCloudMetadata = true
- [ ] MaxHTTPCallsPerExec is set to reasonable limit
- [ ] MaxResponseSize is set to prevent memory exhaustion
- [ ] HTTPTimeout is set to prevent hanging requests
- [ ] MaxExecutionTime is set to prevent runaway workflows
- [ ] MaxNodeExecutions is set to prevent infinite loops
- [ ] Workflow validation is performed before execution
- [ ] Security-relevant events are logged
- [ ] Error messages don't leak sensitive information

## Threat Model

### What Zero Trust Protects Against

✅ **SSRF Attacks**: Cannot access internal services, cloud metadata
✅ **DoS Attacks**: Resource limits prevent exhaustion
✅ **Data Exfiltration**: No file system or environment access
✅ **Arbitrary Code Execution**: Only predefined node types allowed
✅ **Network Scanning**: Cannot probe internal networks
✅ **Credential Theft**: Cannot access cloud metadata endpoints

### What Zero Trust Does NOT Protect Against

❌ **Logic Bugs**: Workflow logic errors are not prevented
❌ **Data Validation**: Input data should be validated by application
❌ **Authentication**: Workflow engine doesn't handle auth (application responsibility)
❌ **Authorization**: Node-level permissions not enforced (future enhancement)

## FAQ

**Q: Why is HTTP disabled by default?**

A: Zero trust principle - deny all by default. Network access is a privileged operation that should require explicit opt-in.

**Q: How do I enable HTTP for development?**

A: Use `DevelopmentConfig()` or manually enable:
```go
config := types.DefaultConfig()
config.AllowHTTP = true
config.BlockLocalhost = false  // For local testing
```

**Q: Can I allow all domains?**

A: Yes, but not recommended for production:
```go
config.AllowHTTP = true
config.AllowedDomains = []string{}  // Empty = allow all (use with caution!)
```

**Q: What if I need to call localhost in production?**

A: This is a security risk. If absolutely necessary:
```go
config.BlockLocalhost = false  // NOT RECOMMENDED
```
Better: Use proper service discovery and internal DNS.

**Q: How do I test workflows that make HTTP calls?**

A: Three options:
1. Use `DevelopmentConfig()` for local testing
2. Set up a mock HTTP server and whitelist it
3. Use production config with test API endpoints whitelisted

**Q: Will old workflows break?**

A: Yes, if they use HTTP nodes. You must explicitly enable HTTP and optionally whitelist domains.

**Q: Is there a way to disable all security?**

A: Not recommended, but you can:
```go
config := types.DevelopmentConfig()
config.AllowHTTP = true
config.AllowedDomains = []string{}
config.BlockPrivateIPs = false
config.BlockLocalhost = false
// etc. (NOT RECOMMENDED FOR PRODUCTION)
```

## References

- [PROTECTION.md](PROTECTION.md) - Complete resource protection documentation
- [ARCHITECTURE.md](../ARCHITECTURE.md) - System architecture
- [README.md](../README.md) - Getting started guide

## Security Contact

For security issues, please open a GitHub issue with the `security` label.
