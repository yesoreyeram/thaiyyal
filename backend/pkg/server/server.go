package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/health"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/logging"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/telemetry"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// Config holds server configuration
type Config struct {
	// Address to listen on (e.g., ":8080")
	Address string
	
	// ReadTimeout for HTTP requests
	ReadTimeout time.Duration
	
	// WriteTimeout for HTTP responses
	WriteTimeout time.Duration
	
	// ShutdownTimeout for graceful shutdown
	ShutdownTimeout time.Duration
	
	// MaxRequestBodySize limits request body size
	MaxRequestBodySize int64
	
	// EnableCORS enables CORS headers
	EnableCORS bool
}

// DefaultConfig returns default server configuration
func DefaultConfig() Config {
	return Config{
		Address:            ":8080",
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		ShutdownTimeout:    10 * time.Second,
		MaxRequestBodySize: 10 * 1024 * 1024, // 10MB
		EnableCORS:         true,
	}
}

// Server is the HTTP API server
type Server struct {
	config           Config
	httpServer       *http.Server
	healthChecker    *health.Checker
	telemetryProvider *telemetry.Provider
	logger           *logging.Logger
	engineConfig     types.Config
}

// New creates a new server instance
func New(config Config, engineConfig types.Config) (*Server, error) {
	// Create logger
	logger := logging.New(logging.DefaultConfig())
	
	// Create telemetry provider
	telemetryConfig := telemetry.DefaultConfig()
	telemetryProvider, err := telemetry.NewProvider(context.Background(), telemetryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry provider: %w", err)
	}
	
	// Create health checker
	healthChecker := health.NewChecker("thaiyyal-workflow-engine", "0.1.0")
	
	// Register basic health checks
	healthChecker.RegisterCheck("engine", func(ctx context.Context) error {
		// Basic check - always healthy if server is running
		return nil
	}, 5*time.Second, true)
	
	server := &Server{
		config:            config,
		healthChecker:     healthChecker,
		telemetryProvider: telemetryProvider,
		logger:            logger,
		engineConfig:      engineConfig,
	}
	
	// Create HTTP server
	mux := http.NewServeMux()
	server.registerRoutes(mux)
	
	server.httpServer = &http.Server{
		Addr:         config.Address,
		Handler:      server.middlewareChain(mux),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	
	return server, nil
}

// registerRoutes registers all HTTP routes
func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Health endpoints
	mux.HandleFunc("/health", s.healthChecker.HTTPHandler())
	mux.HandleFunc("/health/live", s.healthChecker.LivenessHandler())
	mux.HandleFunc("/health/ready", s.healthChecker.ReadinessHandler())
	
	// Metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())
	
	// API endpoints
	mux.HandleFunc("/api/v1/workflow/execute", s.handleExecuteWorkflow)
	mux.HandleFunc("/api/v1/workflow/validate", s.handleValidateWorkflow)
}

// middlewareChain applies middleware to the handler
func (s *Server) middlewareChain(handler http.Handler) http.Handler {
	// Apply CORS if enabled
	if s.config.EnableCORS {
		handler = s.corsMiddleware(handler)
	}
	
	// Apply logging middleware
	handler = s.loggingMiddleware(handler)
	
	// Apply recovery middleware
	handler = s.recoveryMiddleware(handler)
	
	return handler
}

// handleExecuteWorkflow handles workflow execution requests
func (s *Server) handleExecuteWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, s.config.MaxRequestBodySize)
	
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeErrorResponse(w, "Failed to read request body", http.StatusBadRequest, err)
		return
	}
	
	// Execute workflow
	startTime := time.Now()
	eng, err := engine.NewWithConfig(body, s.engineConfig)
	if err != nil {
		s.writeErrorResponse(w, "Failed to create engine", http.StatusBadRequest, err)
		return
	}
	
	// Register telemetry observer
	telemetryObserver := telemetry.NewTelemetryObserver(s.telemetryProvider)
	eng.RegisterObserver(telemetryObserver)
	
	// Execute
	result, err := eng.Execute()
	duration := time.Since(startTime)
	
	// Record metrics
	success := err == nil
	nodesExecuted := 0
	if result != nil {
		nodesExecuted = len(result.NodeResults)
	}
	s.telemetryProvider.RecordWorkflowExecution(r.Context(), "", duration, success, nodesExecuted)
	
	if err != nil {
		s.writeErrorResponse(w, "Workflow execution failed", http.StatusInternalServerError, err)
		return
	}
	
	// Write successful response
	s.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success":        true,
		"results":        result,
		"execution_time": duration.String(),
	})
}

// handleValidateWorkflow handles workflow validation requests
func (s *Server) handleValidateWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, s.config.MaxRequestBodySize)
	
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeErrorResponse(w, "Failed to read request body", http.StatusBadRequest, err)
		return
	}
	
	// Try to create engine (validates the workflow)
	_, err = engine.NewWithConfig(body, s.engineConfig)
	if err != nil {
		s.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
			"valid": false,
			"error": err.Error(),
		})
		return
	}
	
	s.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"valid": true,
	})
}

// writeJSONResponse writes a JSON response
func (s *Server) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.WithError(err).Error("failed to encode response")
	}
}

// writeErrorResponse writes an error response
func (s *Server) writeErrorResponse(w http.ResponseWriter, message string, statusCode int, err error) {
	s.logger.WithError(err).WithField("status_code", statusCode).Error(message)
	
	s.writeJSONResponse(w, statusCode, map[string]interface{}{
		"success": false,
		"error":   message,
		"details": err.Error(),
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.WithField("address", s.config.Address).Info("starting server")
	
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	
	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}
	
	// Shutdown telemetry
	if err := s.telemetryProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown telemetry: %w", err)
	}
	
	s.logger.Info("server shutdown complete")
	return nil
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		
		// Create response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(rw, r)
		
		duration := time.Since(startTime)
		
		s.logger.WithFields(map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": rw.statusCode,
			"duration_ms": duration.Milliseconds(),
			"remote_addr": r.RemoteAddr,
		}).Info("http request")
	})
}

// recoveryMiddleware recovers from panics
func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.WithField("error", fmt.Sprintf("%v", err)).
					WithField("path", r.URL.Path).
					Error("panic recovered")
				
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
