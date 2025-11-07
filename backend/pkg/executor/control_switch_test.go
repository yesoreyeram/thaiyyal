package executor

import (
"testing"

"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestSwitchExecutor_ExpressionMatching tests switch with expression-based matching
func TestSwitchExecutor_ExpressionMatching(t *testing.T) {
tests := []struct {
name          string
input         interface{}
cases         []types.SwitchCase
expectedMatch bool
expectedPath  string
description   string
}{
{
name:  "Exact number match with expression",
input: float64(10),
cases: []types.SwitchCase{
{When: "input == 10", OutputPath: strPtr("path1")},
{When: "input == 20", OutputPath: strPtr("path2")},
{When: "default", OutputPath: strPtr("default"), IsDefault: true},
},
expectedMatch: true,
expectedPath:  "path1",
description:   "Should match first case with expression input == 10",
},
{
name:  "Greater than expression",
input: float64(15),
cases: []types.SwitchCase{
{When: "input > 10", OutputPath: strPtr("large")},
{When: "input <= 10", OutputPath: strPtr("small")},
{When: "default", OutputPath: strPtr("default"), IsDefault: true},
},
expectedMatch: true,
expectedPath:  "large",
description:   "Should match input > 10",
},
{
name:  "Range expression",
input: float64(85),
cases: []types.SwitchCase{
{When: "input >= 90", OutputPath: strPtr("A")},
{When: "input >= 80", OutputPath: strPtr("B")},
{When: "input >= 70", OutputPath: strPtr("C")},
{When: "default", OutputPath: strPtr("F"), IsDefault: true},
},
expectedMatch: true,
expectedPath:  "B",
description:   "Should match input >= 80 (grade B)",
},
{
name:  "No match - use default",
input: float64(30),
cases: []types.SwitchCase{
{When: "input == 10", OutputPath: strPtr("path1")},
{When: "input == 20", OutputPath: strPtr("path2")},
{When: "default", OutputPath: strPtr("fallback"), IsDefault: true},
},
expectedMatch: false,
expectedPath:  "fallback",
description:   "Should use default when no match",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
exec := &SwitchExecutor{}
ctx := &MockExecutionContext{
inputs: map[string][]interface{}{
"test-node": {tt.input},
},
}

node := types.Node{
ID:   "test-node",
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: tt.cases,
},
}

result, err := exec.Execute(ctx, node)
if err != nil {
t.Fatalf("Unexpected error: %v", err)
}

resultMap, ok := result.(map[string]interface{})
if !ok {
t.Fatalf("Expected result to be map, got %T", result)
}

matched, ok := resultMap["matched"].(bool)
if !ok {
t.Fatalf("Expected matched to be bool, got %T", resultMap["matched"])
}

if matched != tt.expectedMatch {
t.Errorf("Expected matched=%v, got %v. Description: %s",
tt.expectedMatch, matched, tt.description)
}

outputPath, ok := resultMap["output_path"].(string)
if !ok {
t.Fatalf("Expected output_path to be string, got %T", resultMap["output_path"])
}

if outputPath != tt.expectedPath {
t.Errorf("Expected output_path='%s', got '%s'", tt.expectedPath, outputPath)
}

// Verify input value is preserved
if resultMap["value"] != tt.input {
t.Errorf("Expected value to be preserved: %v, got %v", tt.input, resultMap["value"])
}
})
}
}

// TestSwitchExecutor_FirstMatchWins tests that first matching case is selected
func TestSwitchExecutor_FirstMatchWins(t *testing.T) {
exec := &SwitchExecutor{}
ctx := &MockExecutionContext{
inputs: map[string][]interface{}{
"test-node": {float64(15)},
},
}

node := types.Node{
ID:   "test-node",
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: []types.SwitchCase{
{When: "input > 10", OutputPath: strPtr("first")},
{When: "input > 5", OutputPath: strPtr("second")}, // Also matches but should not be used
{When: "default", OutputPath: strPtr("default"), IsDefault: true},
},
},
}

result, err := exec.Execute(ctx, node)
if err != nil {
t.Fatalf("Unexpected error: %v", err)
}

resultMap := result.(map[string]interface{})
outputPath := resultMap["output_path"].(string)

if outputPath != "first" {
t.Errorf("Expected first matching case to win, got path '%s'", outputPath)
}

if resultMap["case"].(string) != "input > 10" {
t.Errorf("Expected case='input > 10', got '%s'", resultMap["case"].(string))
}
}

// TestSwitchExecutor_DefaultCase tests default case handling
func TestSwitchExecutor_DefaultCase(t *testing.T) {
exec := &SwitchExecutor{}
ctx := &MockExecutionContext{
inputs: map[string][]interface{}{
"test-node": {float64(100)},
},
}

node := types.Node{
ID:   "test-node",
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: []types.SwitchCase{
{When: "input < 10", OutputPath: strPtr("small")},
{When: "input < 50", OutputPath: strPtr("medium")},
{When: "default", OutputPath: strPtr("other"), IsDefault: true},
},
},
}

result, err := exec.Execute(ctx, node)
if err != nil {
t.Fatalf("Unexpected error: %v", err)
}

resultMap := result.(map[string]interface{})

if resultMap["matched"].(bool) {
t.Error("Expected no match (using default), got matched=true")
}

if resultMap["output_path"].(string) != "other" {
t.Errorf("Expected output_path='other', got '%s'", resultMap["output_path"].(string))
}

if resultMap["case"].(string) != "default" {
t.Errorf("Expected case='default', got '%s'", resultMap["case"].(string))
}
}

// TestSwitchExecutor_Validation tests node validation
func TestSwitchExecutor_Validation(t *testing.T) {
tests := []struct {
name        string
node        types.Node
expectError bool
errorMsg    string
}{
{
name: "Valid node with default case",
node: types.Node{
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: []types.SwitchCase{
{When: "input > 0", OutputPath: strPtr("positive")},
{When: "default", OutputPath: strPtr("default"), IsDefault: true},
},
},
},
expectError: false,
},
{
name: "No cases",
node: types.Node{
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: []types.SwitchCase{},
},
},
expectError: true,
errorMsg:    "cases",
},
{
name: "No default case",
node: types.Node{
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: []types.SwitchCase{
{When: "input > 0", OutputPath: strPtr("positive")},
},
},
},
expectError: true,
errorMsg:    "must have exactly one default case (found 0)",
},
{
name: "Multiple default cases",
node: types.Node{
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: []types.SwitchCase{
{When: "default", OutputPath: strPtr("default1"), IsDefault: true},
{When: "default", OutputPath: strPtr("default2"), IsDefault: true},
},
},
},
expectError: true,
errorMsg:    "default case must be the last case",
},
{
name: "Default case not last",
node: types.Node{
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: []types.SwitchCase{
{When: "default", OutputPath: strPtr("default"), IsDefault: true},
{When: "input > 0", OutputPath: strPtr("positive")},
},
},
},
expectError: true,
errorMsg:    "default case must be the last case",
},
{
name: "Non-default case without output_path",
node: types.Node{
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: []types.SwitchCase{
{When: "input > 0"}, // Missing OutputPath
{When: "default", OutputPath: strPtr("default"), IsDefault: true},
},
},
},
expectError: true,
errorMsg:    "must have output_path",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
exec := &SwitchExecutor{}
err := exec.Validate(tt.node)

if tt.expectError && err == nil {
t.Error("Expected validation error, got nil")
}
if !tt.expectError && err != nil {
t.Errorf("Expected no validation error, got: %v", err)
}
if tt.expectError && err != nil && tt.errorMsg != "" {
if !contains(err.Error(), tt.errorMsg) {
t.Errorf("Expected error to contain '%s', got: %v", tt.errorMsg, err)
}
}
})
}
}

// TestSwitchExecutor_MissingInput tests error handling for missing input
func TestSwitchExecutor_MissingInput(t *testing.T) {
exec := &SwitchExecutor{}
ctx := &MockExecutionContext{
inputs: map[string][]interface{}{}, // No inputs
}

node := types.Node{
ID:   "test-node",
Type: types.NodeTypeSwitch,
Data: types.SwitchData{
Cases: []types.SwitchCase{
{When: "input > 0", OutputPath: strPtr("positive")},
{When: "default", OutputPath: strPtr("default"), IsDefault: true},
},
},
}

_, err := exec.Execute(ctx, node)
if err == nil {
t.Error("Expected error for missing input")
}
}

// TestSwitchExecutor_NodeType tests NodeType method
func TestSwitchExecutor_NodeType(t *testing.T) {
exec := &SwitchExecutor{}
if exec.NodeType() != types.NodeTypeSwitch {
t.Errorf("Expected NodeType to be %s, got %s", types.NodeTypeSwitch, exec.NodeType())
}
}

// Helper functions

func contains(s, substr string) bool {
return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
for i := 0; i <= len(s)-len(substr); i++ {
if s[i:i+len(substr)] == substr {
return true
}
}
return false
}
