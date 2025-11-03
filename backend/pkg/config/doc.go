// Package config provides configuration management for the Thaiyyal workflow engine.
//
// # Overview
//
// The config package centralizes all configuration for the workflow engine,
// providing a pluggable, replaceable configuration system with validation,
// defaults, and environment variable support.
//
// # Features
//
//   - Centralized configuration: Single source of truth
//   - Type-safe: Strongly typed configuration options
//   - Validation: Automatic validation of configuration values
//   - Defaults: Sensible default values
//   - Environment variables: Override from environment
//   - File-based: Load from JSON/YAML files
//   - Runtime updates: Hot-reload configuration
//   - Thread-safe: Concurrent access support
//
// # Configuration Structure
//
// The configuration is organized into logical sections:
//
//   - Execution limits: Timeouts and iteration limits
//   - HTTP settings: HTTP request configuration
//   - Security: Network access control and restrictions
//   - Cache settings: Cache TTL and size limits
//   - Resource limits: Memory and size constraints
//   - Retry settings: Retry and backoff configuration
//
// # Basic Usage
//
//import "github.com/yesoreyeram/thaiyyal/backend/pkg/config"
//
//// Create default configuration
//cfg := config.Default()
//
//// Use in engine
//engine := engine.New(engine.WithConfig(cfg))
//
// # Custom Configuration
//
//cfg := config.New(
//    config.WithMaxExecutionTime(10 * time.Minute),
//    config.WithHTTPTimeout(30 * time.Second),
//    config.WithMaxNodes(1000),
//)
//
// # Default Configuration
//
// The default configuration provides secure, production-ready defaults:
//
//MaxExecutionTime: 5 minutes
//MaxNodeExecutionTime: 30 seconds
//MaxIterations: 10000
//HTTPTimeout: 30 seconds
//MaxHTTPRedirects: 10
//MaxResponseSize: 10MB
//AllowHTTP: false (HTTPS only)
//BlockPrivateIPs: true
//BlockLocalhost: true
//BlockCloudMetadata: true
//DefaultCacheTTL: 1 hour
//MaxCacheSize: 1000
//MaxNodes: 1000
//MaxEdges: 5000
//DefaultMaxAttempts: 3
//DefaultBackoff: 1 second
//
// # Thread Safety
//
// Configuration objects are safe for concurrent read access.
package config
