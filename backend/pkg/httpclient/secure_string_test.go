package httpclient

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestSecureString_String(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "non-empty value",
			value:    "my-secret-password",
			expected: "***REDACTED***",
		},
		{
			name:     "empty value",
			value:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSecureString(tt.value)
			if got := s.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSecureString_Value(t *testing.T) {
	value := "my-secret-password"
	s := NewSecureString(value)

	if got := s.Value(); got != value {
		t.Errorf("Value() = %v, want %v", got, value)
	}
}

func TestSecureString_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "empty",
			value:    "",
			expected: true,
		},
		{
			name:     "non-empty",
			value:    "password",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSecureString(tt.value)
			if got := s.IsEmpty(); got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSecureString_MarshalJSON(t *testing.T) {
	s := NewSecureString("my-secret")
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	// Check that the marshaled value is redacted
	if !strings.Contains(string(data), "REDACTED") {
		t.Errorf("MarshalJSON() = %s, expected to contain REDACTED", string(data))
	}

	// Check that the actual value is NOT in the marshaled JSON
	if strings.Contains(string(data), "my-secret") {
		t.Errorf("MarshalJSON() = %s, should not contain actual secret", string(data))
	}
}

func TestSecureString_UnmarshalJSON(t *testing.T) {
	jsonData := `"my-secret-password"`
	var s SecureString

	err := json.Unmarshal([]byte(jsonData), &s)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if s.Value() != "my-secret-password" {
		t.Errorf("UnmarshalJSON() value = %v, want my-secret-password", s.Value())
	}
}

func TestSecureString_MarshalUnmarshalJSON(t *testing.T) {
	original := "my-secret-token"
	s1 := NewSecureString(original)

	// Marshal to JSON
	data, err := json.Marshal(s1)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Verify the JSON is redacted
	if strings.Contains(string(data), original) {
		t.Errorf("Marshaled JSON contains original secret: %s", string(data))
	}
}

func TestSecureString_InStruct(t *testing.T) {
	type Config struct {
		Username string       `json:"username"`
		Password SecureString `json:"password"`
	}

	cfg := Config{
		Username: "admin",
		Password: NewSecureString("super-secret"),
	}

	// Marshal to JSON
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	jsonStr := string(data)

	// Verify username is in JSON
	if !strings.Contains(jsonStr, "admin") {
		t.Errorf("JSON should contain username: %s", jsonStr)
	}

	// Verify password is redacted
	if strings.Contains(jsonStr, "super-secret") {
		t.Errorf("JSON should not contain password: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, "REDACTED") {
		t.Errorf("JSON should contain REDACTED: %s", jsonStr)
	}
}

func TestSecureString_GoString(t *testing.T) {
	s := NewSecureString("my-secret")
	goStr := fmt.Sprintf("%#v", s)

	// Check that GoString returns a masked representation
	if strings.Contains(goStr, "my-secret") {
		t.Errorf("GoString() should not contain actual secret: %s", goStr)
	}
	if !strings.Contains(goStr, "SecureString") {
		t.Errorf("GoString() should contain SecureString: %s", goStr)
	}
}

func TestSecureString_EmptyGoString(t *testing.T) {
	s := NewSecureString("")
	goStr := fmt.Sprintf("%#v", s)

	expected := "SecureString{}"
	if goStr != expected {
		t.Errorf("GoString() = %v, want %v", goStr, expected)
	}
}

func TestSecureString_PrintFormatting(t *testing.T) {
	s := NewSecureString("my-secret")

	// Test %s format
	if str := fmt.Sprintf("%s", s); str != "***REDACTED***" {
		t.Errorf("%%s format = %v, want ***REDACTED***", str)
	}

	// Test %v format
	if str := fmt.Sprintf("%v", s); str != "***REDACTED***" {
		t.Errorf("%%v format = %v, want ***REDACTED***", str)
	}

	// Verify actual value is not exposed
	if str := fmt.Sprintf("%s", s); strings.Contains(str, "my-secret") {
		t.Errorf("Format should not expose actual secret")
	}
}
