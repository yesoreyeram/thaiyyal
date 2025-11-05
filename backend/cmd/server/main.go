// Command server starts the Thaiyyal workflow engine HTTP API server.
//
// Usage:
//
//	server [flags]
//
// Flags:
//
//	-addr string
//	    Server address (default ":8080")
//	-read-timeout duration
//	    HTTP read timeout (default 30s)
//	-write-timeout duration
//	    HTTP write timeout (default 30s)
//	-max-execution-time duration
//	    Maximum workflow execution time (default 1m)
//	-max-node-executions int
//	    Maximum node executions per workflow (default 10000)
//
// Example:
//
//	# Start server on default port
//	server
//
//	# Start server on custom port with strict limits
//	server -addr :9090 -max-execution-time 30s -max-node-executions 1000
//
// The server exposes the following endpoints:
//
//	POST   /api/v1/workflow/execute        - Execute a workflow
//	POST   /api/v1/workflow/validate       - Validate a workflow
//	POST   /api/v1/workflow/save           - Save a workflow
//	GET    /api/v1/workflow/list           - List all saved workflows
//	GET    /api/v1/workflow/load/{id}      - Load a workflow by ID
//	DELETE /api/v1/workflow/delete/{id}    - Delete a workflow by ID
//	POST   /api/v1/workflow/execute/{id}   - Execute a workflow by ID
//	POST   /api/v1/httpclient/register     - Register an HTTP client
//	GET    /api/v1/httpclient/list         - List registered HTTP clients
//	GET    /health                         - Health check
//	GET    /health/live                    - Liveness probe
//	GET    /health/ready                   - Readiness probe
//	GET    /metrics                        - Prometheus metrics
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/server"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func main() {
	// Define flags
	addr := flag.String("addr", ":8080", "Server address")
	readTimeout := flag.Duration("read-timeout", 30*time.Second, "HTTP read timeout")
	writeTimeout := flag.Duration("write-timeout", 30*time.Second, "HTTP write timeout")
	maxExecutionTime := flag.Duration("max-execution-time", 1*time.Minute, "Maximum workflow execution time")
	maxNodeExecutions := flag.Int("max-node-executions", 10000, "Maximum node executions per workflow")
	maxHTTPCalls := flag.Int("max-http-calls", 100, "Maximum HTTP calls per execution")
	maxLoopIterations := flag.Int("max-loop-iterations", 10000, "Maximum loop iterations")

	flag.Parse()

	// Create server config
	serverConfig := server.Config{
		Address:            *addr,
		ReadTimeout:        *readTimeout,
		WriteTimeout:       *writeTimeout,
		ShutdownTimeout:    10 * time.Second,
		MaxRequestBodySize: 10 * 1024 * 1024, // 10MB
		EnableCORS:         true,
	}

	// Create engine config
	engineConfig := types.DefaultConfig()
	engineConfig.AllowHTTP = true
	engineConfig.MaxExecutionTime = *maxExecutionTime
	engineConfig.MaxNodeExecutions = *maxNodeExecutions
	engineConfig.MaxHTTPCallsPerExec = *maxHTTPCalls
	engineConfig.MaxIterations = *maxLoopIterations

	// Create server
	srv, err := server.New(serverConfig, engineConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("Starting Thaiyyal Workflow Engine Server on %s\n", *addr)
		fmt.Printf("Health check:     http://localhost%s/health\n", *addr)
		fmt.Printf("Liveness probe:   http://localhost%s/health/live\n", *addr)
		fmt.Printf("Readiness probe:  http://localhost%s/health/ready\n", *addr)
		fmt.Printf("Metrics:          http://localhost%s/metrics\n", *addr)
		fmt.Printf("API endpoint:     http://localhost%s/api/v1/workflow/execute\n", *addr)
		fmt.Println("\nPress Ctrl+C to shutdown")

		if err := srv.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-errChan:
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal: %v\n", sig)
		fmt.Println("Shutting down gracefully...")

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), serverConfig.ShutdownTimeout)
		defer cancel()

		// Shutdown server
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Shutdown error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Server stopped")
	}
}
