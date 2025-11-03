package httpclient

import (
	"encoding/json"
	"fmt"
)

// SecureString represents a sensitive string value that is masked in logs and string representations.
// This type should be used for passwords, tokens, API keys, and other sensitive credentials.
type SecureString struct {
	value string
}

// NewSecureString creates a new SecureString from a plain string value.
func NewSecureString(value string) SecureString {
	return SecureString{value: value}
}

// String returns a masked representation of the secure string.
// This ensures sensitive values are not exposed in logs or debug output.
func (s SecureString) String() string {
	if s.value == "" {
		return ""
	}
	return "***REDACTED***"
}

// Value returns the actual string value.
// Use this method only when the actual value is needed (e.g., for authentication).
func (s SecureString) Value() string {
	return s.value
}

// IsEmpty returns true if the secure string is empty.
func (s SecureString) IsEmpty() bool {
	return s.value == ""
}

// MarshalJSON implements json.Marshaler interface.
// The value is masked when marshaled to JSON.
func (s SecureString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (s *SecureString) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	s.value = value
	return nil
}

// MarshalYAML implements yaml.Marshaler interface.
// The value is masked when marshaled to YAML.
func (s SecureString) MarshalYAML() (interface{}, error) {
	return s.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (s *SecureString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	s.value = value
	return nil
}

// GoString returns a Go-syntax representation with masked value.
// This is used by fmt.Printf with %#v format.
func (s SecureString) GoString() string {
	if s.value == "" {
		return `SecureString{}`
	}
	return fmt.Sprintf(`SecureString{value:"%s"}`, s.String())
}
