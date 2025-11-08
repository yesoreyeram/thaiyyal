package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/security"
)

// PlaygroundExecuteRequest represents a request to execute an HTTP call from the playground
type PlaygroundExecuteRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Timeout int               `json:"timeout,omitempty"` // timeout in seconds
}

// PlaygroundExecuteResponse represents the response from executing an HTTP call
type PlaygroundExecuteResponse struct {
	Status     int                  `json:"status"`
	StatusText string               `json:"statusText"`
	Headers    map[string]string    `json:"headers"`
	Data       interface{}          `json:"data"`
	Timing     PlaygroundTimingInfo `json:"timing"`
	Error      string               `json:"error,omitempty"`
}

// PlaygroundTimingInfo contains timing information about the request
type PlaygroundTimingInfo struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Duration  int64  `json:"duration"` // duration in milliseconds
}

// handlePlaygroundExecute handles HTTP execution requests from the playground
func (s *Server) handlePlaygroundExecute(w http.ResponseWriter, r *http.Request) {
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

	// Parse request
	var req PlaygroundExecuteRequest
	if err := json.Unmarshal(body, &req); err != nil {
		s.writeErrorResponse(w, "Failed to parse request", http.StatusBadRequest, err)
		return
	}

	// Validate request
	if err := validatePlaygroundRequest(&req); err != nil {
		s.writeJSONResponse(w, http.StatusBadRequest, PlaygroundExecuteResponse{
			Error: err.Error(),
		})
		return
	}

	// Execute HTTP request
	response := s.executePlaygroundRequest(r.Context(), &req)

	// Write response
	s.writeJSONResponse(w, http.StatusOK, response)
}

// validatePlaygroundRequest validates the playground execute request
func validatePlaygroundRequest(req *PlaygroundExecuteRequest) error {
	// Validate HTTP method
	validMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"PATCH":  true,
		"DELETE": true,
	}
	if !validMethods[req.Method] {
		return fmt.Errorf("invalid HTTP method: %s", req.Method)
	}

	// Validate URL format
	if req.URL == "" {
		return fmt.Errorf("URL is required")
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s (only http and https are allowed)", parsedURL.Scheme)
	}

	// Validate timeout (1-60 seconds)
	if req.Timeout < 0 || req.Timeout > 60 {
		return fmt.Errorf("timeout must be between 0 and 60 seconds")
	}

	// SSRF Protection - validate URL
	ssrfConfig := security.SSRFConfig{
		AllowedSchemes:     []string{"http", "https"},
		AllowPrivateIPs:    false, // Block private IPs by default
		AllowLocalhost:     true,  // Allow localhost for testing
		AllowLinkLocal:     false,
		AllowCloudMetadata: false, // Block cloud metadata endpoints
		AllowedDomains:     []string{},
		BlockedDomains:     []string{},
	}

	protection := security.NewSSRFProtectionWithConfig(ssrfConfig)
	if err := protection.ValidateURL(req.URL); err != nil {
		return fmt.Errorf("URL security validation failed: %w", err)
	}

	return nil
}

// executePlaygroundRequest executes the HTTP request and returns the response
func (s *Server) executePlaygroundRequest(ctx context.Context, req *PlaygroundExecuteRequest) PlaygroundExecuteResponse {
	startTime := time.Now()

	// Set default timeout if not specified
	timeout := 30 * time.Second
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create HTTP request
	var bodyReader io.Reader
	if req.Body != "" && (req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH") {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyReader)
	if err != nil {
		return PlaygroundExecuteResponse{
			Error: fmt.Sprintf("Failed to create HTTP request: %v", err),
		}
	}

	// Add headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Auto-add Content-Type if not present and body is provided
	if bodyReader != nil && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Execute request
	resp, err := client.Do(httpReq)
	if err != nil {
		// Check if it's a timeout error
		if ctx.Err() == context.DeadlineExceeded {
			return PlaygroundExecuteResponse{
				Status:     408,
				StatusText: "Request Timeout",
				Headers:    make(map[string]string),
				Data:       nil,
				Error:      fmt.Sprintf("Request timed out after %v", timeout),
				Timing: PlaygroundTimingInfo{
					StartTime: startTime.Format(time.RFC3339),
					EndTime:   time.Now().Format(time.RFC3339),
					Duration:  time.Since(startTime).Milliseconds(),
				},
			}
		}

		return PlaygroundExecuteResponse{
			Status:     503,
			StatusText: "Service Unavailable",
			Headers:    make(map[string]string),
			Data:       nil,
			Error:      fmt.Sprintf("HTTP request failed: %v", err),
			Timing: PlaygroundTimingInfo{
				StartTime: startTime.Format(time.RFC3339),
				EndTime:   time.Now().Format(time.RFC3339),
				Duration:  time.Since(startTime).Milliseconds(),
			},
		}
	}
	defer resp.Body.Close()

	// Read response body with size limit (10MB)
	limitedReader := io.LimitReader(resp.Body, 10*1024*1024)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return PlaygroundExecuteResponse{
			Status:     resp.StatusCode,
			StatusText: resp.Status,
			Headers:    extractHeaders(resp.Header),
			Data:       nil,
			Error:      fmt.Sprintf("Failed to read response body: %v", err),
			Timing: PlaygroundTimingInfo{
				StartTime: startTime.Format(time.RFC3339),
				EndTime:   time.Now().Format(time.RFC3339),
				Duration:  time.Since(startTime).Milliseconds(),
			},
		}
	}

	// Parse response body
	var responseData interface{}
	contentType := resp.Header.Get("Content-Type")

	// Try to parse as JSON if content-type is JSON or body looks like JSON
	if isJSONContentType(contentType) || looksLikeJSON(bodyBytes) {
		var jsonData interface{}
		dec := json.NewDecoder(bytes.NewReader(bodyBytes))
		dec.UseNumber() // Preserve numbers as json.Number
		if err := dec.Decode(&jsonData); err == nil {
			responseData = jsonData
		} else {
			// If JSON parsing fails, return as string
			responseData = string(bodyBytes)
		}
	} else {
		// Return as string for non-JSON content
		responseData = string(bodyBytes)
	}

	endTime := time.Now()

	return PlaygroundExecuteResponse{
		Status:     resp.StatusCode,
		StatusText: resp.Status,
		Headers:    extractHeaders(resp.Header),
		Data:       responseData,
		Timing: PlaygroundTimingInfo{
			StartTime: startTime.Format(time.RFC3339),
			EndTime:   endTime.Format(time.RFC3339),
			Duration:  endTime.Sub(startTime).Milliseconds(),
		},
	}
}

// extractHeaders extracts headers from http.Header to a simple map
func extractHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

// isJSONContentType returns true if the provided Content-Type header value denotes JSON
func isJSONContentType(ct string) bool {
	if ct == "" {
		return false
	}
	ct = strings.ToLower(ct)
	// Strip parameters such as charset
	if idx := strings.Index(ct, ";"); idx >= 0 {
		ct = ct[:idx]
	}
	ct = strings.TrimSpace(ct)
	if ct == "application/json" {
		return true
	}
	// Accept vendor-specific types like application/problem+json, application/vnd.api+json
	return strings.HasPrefix(ct, "application/") && strings.HasSuffix(ct, "+json")
}

// looksLikeJSON performs a lightweight check on the raw body to see if it appears to be JSON
func looksLikeJSON(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	s := bytes.TrimSpace(b)
	if len(s) == 0 {
		return false
	}
	first := s[0]
	return first == '{' || first == '['
}
