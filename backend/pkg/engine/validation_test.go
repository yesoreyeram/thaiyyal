package engine

import (
	"strings"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestMaxStringLength tests that string length limits are enforced
func TestMaxStringLength(t *testing.T) {
	tests := []struct {
		name         string
		maxLength    int
		stringLength int
		shouldFail   bool
	}{
		{
			name:         "under limit",
			maxLength:    1000,
			stringLength: 500,
			shouldFail:   false,
		},
		{
			name:         "at limit",
			maxLength:    1000,
			stringLength: 1000,
			shouldFail:   false,
		},
		{
			name:         "exceed limit",
			maxLength:    100,
			stringLength: 200,
			shouldFail:   true,
		},
		{
			name:         "unlimited",
			maxLength:    0,
			stringLength: 1000000,
			shouldFail:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a long string
			longString := strings.Repeat("x", tt.stringLength)

			payload := `{
				"nodes": [
					{"id": "1", "type": "text_input", "data": {"text": "` + longString + `"}},
					{"id": "2", "type": "variable", "data": {"var_name": "test", "var_op": "set"}},
					{"id": "3", "type": "visualization", "data": {"mode": "text"}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"},
					{"id": "e2", "source": "2", "target": "3"}
				]
			}`

			config := types.DefaultConfig()
			config.MaxStringLength = tt.maxLength

			engine, err := NewWithConfig([]byte(payload), config)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			_, err = engine.Execute()

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if !strings.Contains(err.Error(), "string too long") && !strings.Contains(err.Error(), "validation failed") {
					t.Errorf("Expected string length error, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestMaxArrayLength tests that array length limits are enforced
func TestMaxArrayLength(t *testing.T) {
	tests := []struct {
		name        string
		maxLength   int
		arrayLength int
		shouldFail  bool
	}{
		{
			name:        "under limit",
			maxLength:   1000,
			arrayLength: 500,
			shouldFail:  false,
		},
		{
			name:        "exceed limit",
			maxLength:   10,
			arrayLength: 20,
			shouldFail:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.DefaultConfig()
			config.MaxArrayLength = tt.maxLength

			// Create array programmatically
			arr := make([]interface{}, tt.arrayLength)
			for i := 0; i < tt.arrayLength; i++ {
				arr[i] = i
			}

			engine, err := NewWithConfig([]byte(`{
				"nodes": [
					{"id": "1", "type": "number", "data": {"value": 1}},
					{"id": "2", "type": "variable", "data": {"var_name": "test", "var_op": "set"}}
				],
				"edges": [
					{"id": "e1", "source": "1", "target": "2"}
				]
			}`), config)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			// Try to set a variable with the array
			err = engine.SetVariable("testArray", arr)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if !strings.Contains(err.Error(), "array too large") {
					t.Errorf("Expected array length error, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestMaxContextDepth tests that nesting depth limits are enforced
func TestMaxContextDepth(t *testing.T) {
	tests := []struct {
		name       string
		maxDepth   int
		nestDepth  int
		shouldFail bool
	}{
		{
			name:       "under limit",
			maxDepth:   10,
			nestDepth:  5,
			shouldFail: false,
		},
		{
			name:       "exceed limit",
			maxDepth:   5,
			nestDepth:  10,
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.DefaultConfig()
			config.MaxContextDepth = tt.maxDepth

			engine, err := NewWithConfig([]byte(`{
				"nodes": [
					{"id": "1", "type": "number", "data": {"value": 1}}
				],
				"edges": []
			}`), config)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			// Create deeply nested object
			var nested interface{} = "value"
			for i := 0; i < tt.nestDepth; i++ {
				nested = map[string]interface{}{"level": nested}
			}

			// Try to set a variable with the nested object
			err = engine.SetVariable("nested", nested)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if !strings.Contains(err.Error(), "too deeply nested") {
					t.Errorf("Expected nesting depth error, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestMaxVariables tests that variable count limits are enforced
func TestMaxVariables(t *testing.T) {
	tests := []struct {
		name       string
		maxVars    int
		varCount   int
		shouldFail bool
	}{
		{
			name:       "under limit",
			maxVars:    10,
			varCount:   5,
			shouldFail: false,
		},
		{
			name:       "at limit",
			maxVars:    10,
			varCount:   10,
			shouldFail: false,
		},
		{
			name:       "exceed limit",
			maxVars:    5,
			varCount:   10,
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.DefaultConfig()
			config.MaxVariables = tt.maxVars

			engine, err := NewWithConfig([]byte(`{
				"nodes": [
					{"id": "1", "type": "number", "data": {"value": 1}}
				],
				"edges": []
			}`), config)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			// Try to set multiple variables
			var lastErr error
			for i := 0; i < tt.varCount; i++ {
				err := engine.SetVariable(string(rune('a'+i)), i)
				if err != nil {
					lastErr = err
					break
				}
			}

			if tt.shouldFail {
				if lastErr == nil {
					t.Errorf("Expected error but got none")
				}
				if !strings.Contains(lastErr.Error(), "maximum variables exceeded") {
					t.Errorf("Expected variable count error, got: %v", lastErr)
				}
			} else {
				if lastErr != nil {
					t.Errorf("Expected no error but got: %v", lastErr)
				}
			}
		})
	}
}

// TestResourceLimitsInDefaultConfigs verifies resource limits are set in all configs
func TestResourceLimitsInDefaultConfigs(t *testing.T) {
	configs := map[string]types.Config{
		"default":     types.DefaultConfig(),
		"validation":  types.ValidationLimits(),
		"development": types.DevelopmentConfig(),
	}

	for name, config := range configs {
		t.Run(name, func(t *testing.T) {
			if config.MaxStringLength <= 0 {
				t.Errorf("%s: MaxStringLength should be positive, got %d", name, config.MaxStringLength)
			}
			if config.MaxArrayLength <= 0 {
				t.Errorf("%s: MaxArrayLength should be positive, got %d", name, config.MaxArrayLength)
			}
			if config.MaxVariables <= 0 {
				t.Errorf("%s: MaxVariables should be positive, got %d", name, config.MaxVariables)
			}
			if config.MaxContextDepth <= 0 {
				t.Errorf("%s: MaxContextDepth should be positive, got %d", name, config.MaxContextDepth)
			}
		})
	}
}

// TestValidationStricterThanDefault verifies validation limits are stricter
func TestValidationStricterThanDefault(t *testing.T) {
	defaultConfig := types.DefaultConfig()
	validationConfig := types.ValidationLimits()

	if validationConfig.MaxStringLength >= defaultConfig.MaxStringLength {
		t.Errorf("Validation MaxStringLength should be stricter")
	}
	if validationConfig.MaxArrayLength >= defaultConfig.MaxArrayLength {
		t.Errorf("Validation MaxArrayLength should be stricter")
	}
	if validationConfig.MaxVariables >= defaultConfig.MaxVariables {
		t.Errorf("Validation MaxVariables should be stricter")
	}
	if validationConfig.MaxContextDepth >= defaultConfig.MaxContextDepth {
		t.Errorf("Validation MaxContextDepth should be stricter")
	}
}

// TestDevelopmentMorePermissiveThanDefault verifies development limits are more permissive
func TestDevelopmentMorePermissiveThanDefault(t *testing.T) {
	defaultConfig := types.DefaultConfig()
	devConfig := types.DevelopmentConfig()

	if devConfig.MaxStringLength <= defaultConfig.MaxStringLength {
		t.Errorf("Development MaxStringLength should be more permissive")
	}
	if devConfig.MaxArrayLength <= defaultConfig.MaxArrayLength {
		t.Errorf("Development MaxArrayLength should be more permissive")
	}
	if devConfig.MaxVariables <= defaultConfig.MaxVariables {
		t.Errorf("Development MaxVariables should be more permissive")
	}
	if devConfig.MaxContextDepth <= defaultConfig.MaxContextDepth {
		t.Errorf("Development MaxContextDepth should be more permissive")
	}
}
