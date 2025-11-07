package engine

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TestSwitchNode_HTTPStatusRouting tests real-world HTTP status code routing
func TestSwitchNode_HTTPStatusRouting(t *testing.T) {
	tests := []struct {
		name             string
		statusCode       float64
		expectedMatched  bool
		expectedPath     string
		expectedCase     string
		description      string
	}{
		{
			name:            "Success status 200",
			statusCode:      200,
			expectedMatched: true,
			expectedPath:    "success",
			expectedCase:    "200 OK",
			description:     "200 status should route to success handler",
		},
		{
			name:            "Created status 201",
			statusCode:      201,
			expectedMatched: true,
			expectedPath:    "created",
			expectedCase:    "201 Created",
			description:     "201 status should route to created handler",
		},
		{
			name:            "Not found status 404",
			statusCode:      404,
			expectedMatched: true,
			expectedPath:    "not_found",
			expectedCase:    "404 Not Found",
			description:     "404 status should route to retry logic",
		},
		{
			name:            "Server error 500",
			statusCode:      500,
			expectedMatched: true,
			expectedPath:    "error",
			expectedCase:    "500 Server Error",
			description:     "500 status should route to error handler",
		},
		{
			name:            "Unknown status 418",
			statusCode:      418,
			expectedMatched: false,
			expectedPath:    "unknown",
			description:     "Unknown status should use default path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{
				Nodes: []types.Node{
					{
						ID:   "status-input",
						Type: types.NodeTypeNumber,
						Data: types.NumberData{
							Value: &tt.statusCode,
						},
					},
					{
						ID:   "status-switch",
						Type: types.NodeTypeSwitch,
						Data: types.SwitchData{
							Cases: []types.SwitchCase{
								{When: "200 OK", Value: float64(200), OutputPath: strPtr("success")},
								{When: "201 Created", Value: float64(201), OutputPath: strPtr("created")},
								{When: "404 Not Found", Value: float64(404), OutputPath: strPtr("not_found")},
								{When: "500 Server Error", Value: float64(500), OutputPath: strPtr("error")},
							},
							DefaultPath: strPtr("unknown"),
						},
					},
				},
				Edges: []types.Edge{
					{Source: "status-input", Target: "status-switch"},
				},
			}

			engine, err := New(mustMarshal(payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Failed to execute workflow: %v", err)
			}

			switchResult := mustGetMapResult(t, result, "status-switch")
			matched := switchResult["matched"].(bool)
			outputPath := switchResult["output_path"].(string)

			if matched != tt.expectedMatched {
				t.Errorf("Expected matched=%v, got %v. Description: %s",
					tt.expectedMatched, matched, tt.description)
			}

			if outputPath != tt.expectedPath {
				t.Errorf("Expected output_path='%s', got '%s'", tt.expectedPath, outputPath)
			}

			if tt.expectedMatched {
				if caseField, ok := switchResult["case"].(string); ok {
					if caseField != tt.expectedCase {
						t.Errorf("Expected case='%s', got '%s'", tt.expectedCase, caseField)
					}
				} else {
					t.Error("Expected case field to be present and string")
				}
			}

			// Verify value is preserved
			if switchResult["value"].(float64) != tt.statusCode {
				t.Errorf("Expected value to be preserved: %v, got %v", tt.statusCode, switchResult["value"])
			}
		})
	}
}

// TestSwitchNode_GradeAssignment tests grade calculation with ranges
func TestSwitchNode_GradeAssignment(t *testing.T) {
	tests := []struct {
		name         string
		score        float64
		expectedGrade string
		description  string
	}{
		{
			name:          "Perfect score",
			score:         100,
			expectedGrade: "A",
			description:   "100 should get grade A",
		},
		{
			name:          "High A",
			score:         95,
			expectedGrade: "A",
			description:   "95 should get grade A",
		},
		{
			name:          "Boundary A",
			score:         90,
			expectedGrade: "A",
			description:   "90 (boundary) should get grade A",
		},
		{
			name:          "High B",
			score:         89,
			expectedGrade: "B",
			description:   "89 should get grade B",
		},
		{
			name:          "Mid B",
			score:         85,
			expectedGrade: "B",
			description:   "85 should get grade B",
		},
		{
			name:          "Boundary B",
			score:         80,
			expectedGrade: "B",
			description:   "80 (boundary) should get grade B",
		},
		{
			name:          "High C",
			score:         79,
			expectedGrade: "C",
			description:   "79 should get grade C",
		},
		{
			name:          "Boundary C",
			score:         70,
			expectedGrade: "C",
			description:   "70 (boundary) should get grade C",
		},
		{
			name:          "High D",
			score:         69,
			expectedGrade: "D",
			description:   "69 should get grade D",
		},
		{
			name:          "Boundary D",
			score:         60,
			expectedGrade: "D",
			description:   "60 (boundary) should get grade D",
		},
		{
			name:          "Failing grade",
			score:         59,
			expectedGrade: "F",
			description:   "59 should get grade F (default)",
		},
		{
			name:          "Very low score",
			score:         25,
			expectedGrade: "F",
			description:   "25 should get grade F (default)",
		},
		{
			name:          "Zero score",
			score:         0,
			expectedGrade: "F",
			description:   "0 should get grade F (default)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{
				Nodes: []types.Node{
					{
						ID:   "score-input",
						Type: types.NodeTypeNumber,
						Data: types.NumberData{
							Value: &tt.score,
						},
					},
					{
						ID:   "grade-switch",
						Type: types.NodeTypeSwitch,
						Data: types.SwitchData{
							Cases: []types.SwitchCase{
								{When: ">=90", OutputPath: strPtr("A")},
								{When: ">=80", OutputPath: strPtr("B")},
								{When: ">=70", OutputPath: strPtr("C")},
								{When: ">=60", OutputPath: strPtr("D")},
							},
							DefaultPath: strPtr("F"),
						},
					},
				},
				Edges: []types.Edge{
					{Source: "score-input", Target: "grade-switch"},
				},
			}

			engine, err := New(mustMarshal(payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Failed to execute workflow: %v", err)
			}

			switchResult := mustGetMapResult(t, result, "grade-switch")
			outputPath := switchResult["output_path"].(string)

			if outputPath != tt.expectedGrade {
				t.Errorf("Expected grade '%s', got '%s'. Description: %s",
					tt.expectedGrade, outputPath, tt.description)
			}
		})
	}
}

// TestSwitchNode_UserRoleRouting tests string-based role routing
func TestSwitchNode_UserRoleRouting(t *testing.T) {
	tests := []struct {
		name         string
		role         string
		expectedPath string
		description  string
	}{
		{
			name:         "Admin role",
			role:         "admin",
			expectedPath: "admin_panel",
			description:  "Admin should access admin panel",
		},
		{
			name:         "Moderator role",
			role:         "moderator",
			expectedPath: "mod_tools",
			description:  "Moderator should access mod tools",
		},
		{
			name:         "User role",
			role:         "user",
			expectedPath: "user_dashboard",
			description:  "User should access user dashboard",
		},
		{
			name:         "Guest user",
			role:         "guest",
			expectedPath: "login",
			description:  "Guest should be redirected to login",
		},
		{
			name:         "Unknown role",
			role:         "unknown",
			expectedPath: "login",
			description:  "Unknown role should use default (login)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{
				Nodes: []types.Node{
					{
						ID:   "role-input",
						Type: types.NodeTypeTextInput,
						Data: types.TextInputData{
							Text: &tt.role,
						},
					},
					{
						ID:   "role-switch",
						Type: types.NodeTypeSwitch,
						Data: types.SwitchData{
							Cases: []types.SwitchCase{
								{When: "admin", Value: "admin", OutputPath: strPtr("admin_panel")},
								{When: "moderator", Value: "moderator", OutputPath: strPtr("mod_tools")},
								{When: "user", Value: "user", OutputPath: strPtr("user_dashboard")},
							},
							DefaultPath: strPtr("login"),
						},
					},
				},
				Edges: []types.Edge{
					{Source: "role-input", Target: "role-switch"},
				},
			}

			engine, err := New(mustMarshal(payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Failed to execute workflow: %v", err)
			}

			switchResult := mustGetMapResult(t, result, "role-switch")
			outputPath := switchResult["output_path"].(string)

			if outputPath != tt.expectedPath {
				t.Errorf("Expected path '%s', got '%s'. Description: %s",
					tt.expectedPath, outputPath, tt.description)
			}
		})
	}
}

// TestSwitchNode_ContentTypeRouting tests content-type based routing
func TestSwitchNode_ContentTypeRouting(t *testing.T) {
	tests := []struct {
		name          string
		contentType   string
		expectedPath  string
		description   string
	}{
		{
			name:         "JSON content",
			contentType:  "application/json",
			expectedPath: "json_parser",
			description:  "JSON should route to JSON parser",
		},
		{
			name:         "XML content",
			contentType:  "application/xml",
			expectedPath: "xml_parser",
			description:  "XML should route to XML parser",
		},
		{
			name:         "CSV content",
			contentType:  "text/csv",
			expectedPath: "csv_parser",
			description:  "CSV should route to CSV parser",
		},
		{
			name:         "Plain text",
			contentType:  "text/plain",
			expectedPath: "text_parser",
			description:  "Plain text should route to text parser",
		},
		{
			name:         "Unknown type",
			contentType:  "application/octet-stream",
			expectedPath: "raw_handler",
			description:  "Unknown type should use default (raw handler)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{
				Nodes: []types.Node{
					{
						ID:   "content-type-input",
						Type: types.NodeTypeTextInput,
						Data: types.TextInputData{
							Text: &tt.contentType,
						},
					},
					{
						ID:   "content-switch",
						Type: types.NodeTypeSwitch,
						Data: types.SwitchData{
							Cases: []types.SwitchCase{
								{When: "JSON", Value: "application/json", OutputPath: strPtr("json_parser")},
								{When: "XML", Value: "application/xml", OutputPath: strPtr("xml_parser")},
								{When: "CSV", Value: "text/csv", OutputPath: strPtr("csv_parser")},
								{When: "Text", Value: "text/plain", OutputPath: strPtr("text_parser")},
							},
							DefaultPath: strPtr("raw_handler"),
						},
					},
				},
				Edges: []types.Edge{
					{Source: "content-type-input", Target: "content-switch"},
				},
			}

			engine, err := New(mustMarshal(payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Failed to execute workflow: %v", err)
			}

			switchResult := mustGetMapResult(t, result, "content-switch")
			outputPath := switchResult["output_path"].(string)

			if outputPath != tt.expectedPath {
				t.Errorf("Expected path '%s', got '%s'. Description: %s",
					tt.expectedPath, outputPath, tt.description)
			}
		})
	}
}

// TestSwitchNode_PriorityQueueRouting tests priority-based routing
func TestSwitchNode_PriorityQueueRouting(t *testing.T) {
	tests := []struct {
		name         string
		priority     float64
		expectedPath string
		description  string
	}{
		{
			name:         "Critical priority",
			priority:     10,
			expectedPath: "critical",
			description:  "Priority 10 should route to critical queue",
		},
		{
			name:         "Boundary critical",
			priority:     9,
			expectedPath: "critical",
			description:  "Priority 9 should route to critical queue",
		},
		{
			name:         "High priority",
			priority:     8,
			expectedPath: "high",
			description:  "Priority 8 should route to high queue",
		},
		{
			name:         "Boundary high",
			priority:     7,
			expectedPath: "high",
			description:  "Priority 7 should route to high queue",
		},
		{
			name:         "Medium priority",
			priority:     5,
			expectedPath: "medium",
			description:  "Priority 5 should route to medium queue",
		},
		{
			name:         "Boundary medium",
			priority:     4,
			expectedPath: "medium",
			description:  "Priority 4 should route to medium queue",
		},
		{
			name:         "Low priority",
			priority:     2,
			expectedPath: "low",
			description:  "Priority 2 should route to low queue",
		},
		{
			name:         "Boundary low",
			priority:     1,
			expectedPath: "low",
			description:  "Priority 1 should route to low queue",
		},
		{
			name:         "Invalid priority",
			priority:     0,
			expectedPath: "invalid",
			description:  "Priority 0 should use default (invalid)",
		},
		{
			name:         "Negative priority",
			priority:     -5,
			expectedPath: "invalid",
			description:  "Negative priority should use default (invalid)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{
				Nodes: []types.Node{
					{
						ID:   "priority-input",
						Type: types.NodeTypeNumber,
						Data: types.NumberData{
							Value: &tt.priority,
						},
					},
					{
						ID:   "priority-switch",
						Type: types.NodeTypeSwitch,
						Data: types.SwitchData{
							Cases: []types.SwitchCase{
								{When: ">=9", OutputPath: strPtr("critical")},
								{When: ">=7", OutputPath: strPtr("high")},
								{When: ">=4", OutputPath: strPtr("medium")},
								{When: ">=1", OutputPath: strPtr("low")},
							},
							DefaultPath: strPtr("invalid"),
						},
					},
				},
				Edges: []types.Edge{
					{Source: "priority-input", Target: "priority-switch"},
				},
			}

			engine, err := New(mustMarshal(payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Failed to execute workflow: %v", err)
			}

			switchResult := mustGetMapResult(t, result, "priority-switch")
			outputPath := switchResult["output_path"].(string)

			if outputPath != tt.expectedPath {
				t.Errorf("Expected path '%s', got '%s'. Description: %s",
					tt.expectedPath, outputPath, tt.description)
			}
		})
	}
}

// TestSwitchNode_BooleanRouting tests boolean value routing
func TestSwitchNode_BooleanRouting(t *testing.T) {
	tests := []struct {
		name         string
		value        bool
		expectedPath string
		description  string
	}{
		{
			name:         "True value",
			value:        true,
			expectedPath: "enabled",
			description:  "True should route to enabled path",
		},
		{
			name:         "False value",
			value:        false,
			expectedPath: "disabled",
			description:  "False should route to disabled path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{
				Nodes: []types.Node{
					{
						ID:   "bool-input",
						Type: types.NodeTypeBooleanInput,
						Data: types.BooleanInputData{
							BooleanValue: &tt.value,
						},
					},
					{
						ID:   "bool-switch",
						Type: types.NodeTypeSwitch,
						Data: types.SwitchData{
							Cases: []types.SwitchCase{
								{When: "true", Value: true, OutputPath: strPtr("enabled")},
								{When: "false", Value: false, OutputPath: strPtr("disabled")},
							},
							DefaultPath: strPtr("error"),
						},
					},
				},
				Edges: []types.Edge{
					{Source: "bool-input", Target: "bool-switch"},
				},
			}

			engine, err := New(mustMarshal(payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Failed to execute workflow: %v", err)
			}

			switchResult := mustGetMapResult(t, result, "bool-switch")
			outputPath := switchResult["output_path"].(string)

			if outputPath != tt.expectedPath {
				t.Errorf("Expected path '%s', got '%s'. Description: %s",
					tt.expectedPath, outputPath, tt.description)
			}
		})
	}
}

// TestSwitchNode_FirstMatchWins tests that first matching case is selected
func TestSwitchNode_FirstMatchWins(t *testing.T) {
	value := float64(50)
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "input",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{
					Value: &value,
				},
			},
			{
				ID:   "switch",
				Type: types.NodeTypeSwitch,
				Data: types.SwitchData{
					Cases: []types.SwitchCase{
						{When: ">10", OutputPath: strPtr("first")},   // Matches
						{When: ">40", OutputPath: strPtr("second")},  // Also matches but shouldn't be used
						{When: ">50", OutputPath: strPtr("third")},   // Doesn't match
					},
				},
			},
		},
		Edges: []types.Edge{
			{Source: "input", Target: "switch"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	switchResult := mustGetMapResult(t, result, "switch")
	outputPath := switchResult["output_path"].(string)

	if outputPath != "first" {
		t.Errorf("Expected first matching case to win, got path '%s'", outputPath)
	}

	// Verify the case field contains the first match
	if caseField, ok := switchResult["case"].(string); ok {
		if caseField != ">10" {
			t.Errorf("Expected case='>10', got '%s'", caseField)
		}
	}
}

// TestSwitchNode_WithMultipleInputs tests switch with value from previous node
func TestSwitchNode_WithMultipleInputs(t *testing.T) {
	value := float64(50)

	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "input",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: &value},
			},
			{
				ID:   "result-switch",
				Type: types.NodeTypeSwitch,
				Data: types.SwitchData{
					Cases: []types.SwitchCase{
						{When: "<25", OutputPath: strPtr("low")},
						{When: ">=25", OutputPath: strPtr("medium")},
						{When: ">=75", OutputPath: strPtr("high")},
					},
					DefaultPath: strPtr("unknown"),
				},
			},
		},
		Edges: []types.Edge{
			{Source: "input", Target: "result-switch"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	switchResult := mustGetMapResult(t, result, "result-switch")
	outputPath := switchResult["output_path"].(string)

	// Value 50 should match >=25 (first match wins)
	if outputPath != "medium" {
		t.Errorf("Expected output_path='medium', got '%s'", outputPath)
	}
}

// TestSwitchNode_NestedInWorkflow tests switch in multi-stage workflow
func TestSwitchNode_NestedInWorkflow(t *testing.T) {
	age := float64(25)

	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "age-input",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: &age},
			},
			{
				ID:   "age-category",
				Type: types.NodeTypeSwitch,
				Data: types.SwitchData{
					Cases: []types.SwitchCase{
						{When: "<18", OutputPath: strPtr("minor")},
						{When: "<65", OutputPath: strPtr("adult")},
						{When: ">=65", OutputPath: strPtr("senior")},
					},
				},
			},
			{
				ID:   "message",
				Type: types.NodeTypeTextInput,
				Data: types.TextInputData{Text: strPtr("Processing complete")},
			},
		},
		Edges: []types.Edge{
			{Source: "age-input", Target: "age-category"},
			{Source: "age-category", Target: "message"},
		},
	}

	engine, err := New(mustMarshal(payload))
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	switchResult := mustGetMapResult(t, result, "age-category")
	outputPath := switchResult["output_path"].(string)

	if outputPath != "adult" {
		t.Errorf("Expected path 'adult' for age 25, got '%s'", outputPath)
	}

	// Verify message node executed (workflow continued)
	if _, exists := result.NodeResults["message"]; !exists {
		t.Error("Expected message node to execute after switch")
	}
}

// TestSwitchNode_EmptyCases tests validation catches empty cases
func TestSwitchNode_EmptyCases(t *testing.T) {
	value := float64(10)
	payload := types.Payload{
		Nodes: []types.Node{
			{
				ID:   "input",
				Type: types.NodeTypeNumber,
				Data: types.NumberData{Value: &value},
			},
			{
				ID:   "switch",
				Type: types.NodeTypeSwitch,
				Data: types.SwitchData{
					Cases: []types.SwitchCase{}, // Empty!
				},
			},
		},
		Edges: []types.Edge{
			{Source: "input", Target: "switch"},
		},
	}

	_, err := New(mustMarshal(payload))
	if err == nil {
		t.Error("Expected validation error for empty cases, got nil")
	}
}

// TestSwitchNode_PreservesValueType tests that value types are preserved
func TestSwitchNode_PreservesValueType(t *testing.T) {
	tests := []struct {
		name     string
		nodeType types.NodeType
		nodeData types.NodeDataInterface
		value    interface{}
	}{
		{
			name:     "Number type",
			nodeType: types.NodeTypeNumber,
			nodeData: types.NumberData{Value: floatPtr(42.5)},
			value:    float64(42.5),
		},
		{
			name:     "String type",
			nodeType: types.NodeTypeTextInput,
			nodeData: types.TextInputData{Text: strPtr("hello")},
			value:    "hello",
		},
		{
			name:     "Boolean type",
			nodeType: types.NodeTypeBooleanInput,
			nodeData: types.BooleanInputData{BooleanValue: boolPtr(true)},
			value:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{
				Nodes: []types.Node{
					{
						ID:   "input",
						Type: tt.nodeType,
						Data: tt.nodeData,
					},
					{
						ID:   "switch",
						Type: types.NodeTypeSwitch,
						Data: types.SwitchData{
							Cases: []types.SwitchCase{
								{When: "any", Value: tt.value, OutputPath: strPtr("matched")},
							},
							DefaultPath: strPtr("default"),
						},
					},
				},
				Edges: []types.Edge{
					{Source: "input", Target: "switch"},
				},
			}

			engine, err := New(mustMarshal(payload))
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			result, err := engine.Execute()
			if err != nil {
				t.Fatalf("Failed to execute workflow: %v", err)
			}

			switchResult := mustGetMapResult(t, result, "switch")
			
			// Verify value is preserved and has correct type
			resultValue := switchResult["value"]
			if resultValue != tt.value {
				t.Errorf("Expected value %v (%T), got %v (%T)", 
					tt.value, tt.value, resultValue, resultValue)
			}
		})
	}
}

// Helper functions - using package-level helpers from other test files

func floatPtr(f float64) *float64 {
	return &f
}
