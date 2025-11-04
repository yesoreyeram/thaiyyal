package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Status represents the health status
type Status string

const (
	// StatusHealthy indicates the component is healthy
	StatusHealthy Status = "healthy"
	
	// StatusUnhealthy indicates the component is unhealthy
	StatusUnhealthy Status = "unhealthy"
	
	// StatusDegraded indicates the component is degraded but functional
	StatusDegraded Status = "degraded"
)

// CheckFunc is a function that performs a health check
type CheckFunc func(ctx context.Context) error

// Check represents a single health check
type Check struct {
	Name        string
	CheckFunc   CheckFunc
	Timeout     time.Duration
	Critical    bool // If true, failure makes the entire service unhealthy
	lastChecked time.Time
	lastStatus  Status
	lastError   error
	mu          sync.RWMutex
}

// Checker manages health checks for the service
type Checker struct {
	checks map[string]*Check
	mu     sync.RWMutex
	
	// Service metadata
	serviceName    string
	serviceVersion string
	startTime      time.Time
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status         Status                    `json:"status"`
	ServiceName    string                    `json:"service_name"`
	ServiceVersion string                    `json:"service_version"`
	Uptime         string                    `json:"uptime"`
	Timestamp      time.Time                 `json:"timestamp"`
	Checks         map[string]CheckResult    `json:"checks,omitempty"`
}

// CheckResult represents the result of a single health check
type CheckResult struct {
	Status      Status    `json:"status"`
	LastChecked time.Time `json:"last_checked"`
	Error       string    `json:"error,omitempty"`
}

// NewChecker creates a new health checker
func NewChecker(serviceName, serviceVersion string) *Checker {
	return &Checker{
		checks:         make(map[string]*Check),
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		startTime:      time.Now(),
	}
}

// RegisterCheck registers a new health check
func (c *Checker) RegisterCheck(name string, checkFunc CheckFunc, timeout time.Duration, critical bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.checks[name] = &Check{
		Name:      name,
		CheckFunc: checkFunc,
		Timeout:   timeout,
		Critical:  critical,
	}
}

// Check performs all registered health checks
func (c *Checker) Check(ctx context.Context) HealthResponse {
	c.mu.RLock()
	checks := make(map[string]*Check, len(c.checks))
	for name, check := range c.checks {
		checks[name] = check
	}
	c.mu.RUnlock()
	
	results := make(map[string]CheckResult)
	overallStatus := StatusHealthy
	
	// Run all checks
	for name, check := range checks {
		result := c.runCheck(ctx, check)
		results[name] = result
		
		// Determine overall status
		if check.Critical && result.Status == StatusUnhealthy {
			overallStatus = StatusUnhealthy
		} else if result.Status == StatusDegraded && overallStatus == StatusHealthy {
			overallStatus = StatusDegraded
		}
	}
	
	return HealthResponse{
		Status:         overallStatus,
		ServiceName:    c.serviceName,
		ServiceVersion: c.serviceVersion,
		Uptime:         time.Since(c.startTime).String(),
		Timestamp:      time.Now(),
		Checks:         results,
	}
}

// runCheck executes a single health check with timeout
func (c *Checker) runCheck(ctx context.Context, check *Check) CheckResult {
	check.mu.Lock()
	defer check.mu.Unlock()
	
	// Create context with timeout
	checkCtx, cancel := context.WithTimeout(ctx, check.Timeout)
	defer cancel()
	
	// Run the check
	errChan := make(chan error, 1)
	go func() {
		errChan <- check.CheckFunc(checkCtx)
	}()
	
	var err error
	select {
	case err = <-errChan:
	case <-checkCtx.Done():
		err = fmt.Errorf("health check timed out after %v", check.Timeout)
	}
	
	// Update check state
	check.lastChecked = time.Now()
	check.lastError = err
	
	if err != nil {
		check.lastStatus = StatusUnhealthy
	} else {
		check.lastStatus = StatusHealthy
	}
	
	result := CheckResult{
		Status:      check.lastStatus,
		LastChecked: check.lastChecked,
	}
	
	if err != nil {
		result.Error = err.Error()
	}
	
	return result
}

// Liveness returns a simple liveness check (always healthy if service is running)
func (c *Checker) Liveness(ctx context.Context) HealthResponse {
	return HealthResponse{
		Status:         StatusHealthy,
		ServiceName:    c.serviceName,
		ServiceVersion: c.serviceVersion,
		Uptime:         time.Since(c.startTime).String(),
		Timestamp:      time.Now(),
	}
}

// Readiness performs all checks and returns readiness status
func (c *Checker) Readiness(ctx context.Context) HealthResponse {
	return c.Check(ctx)
}

// HTTPHandler returns an HTTP handler for health checks
func (c *Checker) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := c.Check(r.Context())
		
		w.Header().Set("Content-Type", "application/json")
		
		// Set status code based on health
		statusCode := http.StatusOK
		if response.Status == StatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		} else if response.Status == StatusDegraded {
			statusCode = http.StatusOK // Still return 200 for degraded
		}
		
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}
}

// LivenessHandler returns an HTTP handler for liveness probes
func (c *Checker) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := c.Liveness(r.Context())
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// ReadinessHandler returns an HTTP handler for readiness probes
func (c *Checker) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := c.Readiness(r.Context())
		
		w.Header().Set("Content-Type", "application/json")
		
		// Set status code based on health
		statusCode := http.StatusOK
		if response.Status == StatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		}
		
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}
}
