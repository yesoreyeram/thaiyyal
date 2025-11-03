// Package health provides health check and readiness probe functionality.
// It enables monitoring of service health with support for:
//   - Liveness probes to detect if the service is running
//   - Readiness probes to detect if the service can handle requests
//   - Custom health checks for dependencies
//   - HTTP handlers for health endpoints
package health
