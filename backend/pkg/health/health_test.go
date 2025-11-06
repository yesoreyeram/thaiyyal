package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewChecker(t *testing.T) {
	checker := NewChecker("test-service", "1.0.0")

	if checker == nil {
		t.Fatal("NewChecker() returned nil")
	}

	if checker.serviceName != "test-service" {
		t.Errorf("serviceName = %v, want %v", checker.serviceName, "test-service")
	}

	if checker.serviceVersion != "1.0.0" {
		t.Errorf("serviceVersion = %v, want %v", checker.serviceVersion, "1.0.0")
	}
}

func TestRegisterCheck(t *testing.T) {
	checker := NewChecker("test", "1.0")

	checkFunc := func(ctx context.Context) error {
		return nil
	}

	checker.RegisterCheck("test-check", checkFunc, 5*time.Second, true)

	checker.mu.RLock()
	defer checker.mu.RUnlock()

	if len(checker.checks) != 1 {
		t.Errorf("len(checks) = %v, want 1", len(checker.checks))
	}

	check, ok := checker.checks["test-check"]
	if !ok {
		t.Fatal("check not found")
	}

	if check.Name != "test-check" {
		t.Errorf("check.Name = %v, want %v", check.Name, "test-check")
	}

	if check.Critical != true {
		t.Errorf("check.Critical = %v, want true", check.Critical)
	}
}

func TestCheckHealthy(t *testing.T) {
	checker := NewChecker("test", "1.0")

	// Register healthy check
	checker.RegisterCheck("always-healthy", func(ctx context.Context) error {
		return nil
	}, 5*time.Second, true)

	ctx := context.Background()
	response := checker.Check(ctx)

	if response.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", response.Status, StatusHealthy)
	}

	if response.ServiceName != "test" {
		t.Errorf("ServiceName = %v, want %v", response.ServiceName, "test")
	}

	if len(response.Checks) != 1 {
		t.Errorf("len(Checks) = %v, want 1", len(response.Checks))
	}

	checkResult, ok := response.Checks["always-healthy"]
	if !ok {
		t.Fatal("check result not found")
	}

	if checkResult.Status != StatusHealthy {
		t.Errorf("check Status = %v, want %v", checkResult.Status, StatusHealthy)
	}
}

func TestCheckUnhealthy(t *testing.T) {
	checker := NewChecker("test", "1.0")

	// Register unhealthy critical check
	checker.RegisterCheck("always-fails", func(ctx context.Context) error {
		return errors.New("check failed")
	}, 5*time.Second, true)

	ctx := context.Background()
	response := checker.Check(ctx)

	if response.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", response.Status, StatusUnhealthy)
	}

	checkResult := response.Checks["always-fails"]
	if checkResult.Status != StatusUnhealthy {
		t.Errorf("check Status = %v, want %v", checkResult.Status, StatusUnhealthy)
	}

	if checkResult.Error == "" {
		t.Error("expected error message")
	}
}

func TestCheckTimeout(t *testing.T) {
	checker := NewChecker("test", "1.0")

	// Register check that times out
	checker.RegisterCheck("timeout", func(ctx context.Context) error {
		time.Sleep(2 * time.Second)
		return nil
	}, 100*time.Millisecond, true)

	ctx := context.Background()
	response := checker.Check(ctx)

	if response.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", response.Status, StatusUnhealthy)
	}

	checkResult := response.Checks["timeout"]
	if checkResult.Status != StatusUnhealthy {
		t.Errorf("check Status = %v, want %v", checkResult.Status, StatusUnhealthy)
	}
}

func TestCheckNonCritical(t *testing.T) {
	checker := NewChecker("test", "1.0")

	// Register non-critical failing check
	checker.RegisterCheck("non-critical", func(ctx context.Context) error {
		return errors.New("non-critical failure")
	}, 5*time.Second, false)

	// Register healthy critical check
	checker.RegisterCheck("critical", func(ctx context.Context) error {
		return nil
	}, 5*time.Second, true)

	ctx := context.Background()
	response := checker.Check(ctx)

	// Overall should still be healthy since failing check is not critical
	if response.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", response.Status, StatusHealthy)
	}
}

func TestLiveness(t *testing.T) {
	checker := NewChecker("test", "1.0")

	ctx := context.Background()
	response := checker.Liveness(ctx)

	if response.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", response.Status, StatusHealthy)
	}

	if len(response.Checks) != 0 {
		t.Errorf("len(Checks) = %v, want 0", len(response.Checks))
	}
}

func TestReadiness(t *testing.T) {
	checker := NewChecker("test", "1.0")

	checker.RegisterCheck("ready", func(ctx context.Context) error {
		return nil
	}, 5*time.Second, true)

	ctx := context.Background()
	response := checker.Readiness(ctx)

	if response.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", response.Status, StatusHealthy)
	}

	if len(response.Checks) != 1 {
		t.Errorf("len(Checks) = %v, want 1", len(response.Checks))
	}
}

func TestHTTPHandler(t *testing.T) {
	checker := NewChecker("test", "1.0")
	checker.RegisterCheck("test", func(ctx context.Context) error {
		return nil
	}, 5*time.Second, true)

	handler := checker.HTTPHandler()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}

	var response HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusHealthy {
		t.Errorf("response.Status = %v, want %v", response.Status, StatusHealthy)
	}
}

func TestHTTPHandlerUnhealthy(t *testing.T) {
	checker := NewChecker("test", "1.0")
	checker.RegisterCheck("test", func(ctx context.Context) error {
		return errors.New("failure")
	}, 5*time.Second, true)

	handler := checker.HTTPHandler()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusServiceUnavailable)
	}
}

func TestLivenessHandler(t *testing.T) {
	checker := NewChecker("test", "1.0")

	handler := checker.LivenessHandler()

	req := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestReadinessHandler(t *testing.T) {
	checker := NewChecker("test", "1.0")
	checker.RegisterCheck("test", func(ctx context.Context) error {
		return nil
	}, 5*time.Second, true)

	handler := checker.ReadinessHandler()

	req := httptest.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestReadinessHandlerUnhealthy(t *testing.T) {
	checker := NewChecker("test", "1.0")
	checker.RegisterCheck("test", func(ctx context.Context) error {
		return errors.New("not ready")
	}, 5*time.Second, true)

	handler := checker.ReadinessHandler()

	req := httptest.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusServiceUnavailable)
	}
}
