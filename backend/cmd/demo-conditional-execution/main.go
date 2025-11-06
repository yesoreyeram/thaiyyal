package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/engine"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func main() {
	fmt.Println("=================================================")
	fmt.Println("Conditional Execution Demo")
	fmt.Println("=================================================\n")

	// Demo 1: Age-based routing (adult vs minor)
	demo1AgeBasedRouting()

	// Demo 2: Switch-based HTTP status routing
	demo2SwitchRouting()

	// Demo 3: Nested conditions
	demo3NestedConditions()
}

func demo1AgeBasedRouting() {
	fmt.Println("📋 DEMO 1: Age-Based API Routing")
	fmt.Println("----------------------------------")
	fmt.Println("Scenario: If age >= 18, call profile API -> sports API")
	fmt.Println("          If age < 18, call education API")
	fmt.Println()

	testAges := []float64{25, 15}

	for _, age := range testAges {
		fmt.Printf("Testing with age = %.0f:\n", age)

		payload := types.Payload{
			Nodes: []types.Node{
				{ID: "user_age", Type: types.NodeTypeNumber, Data: types.NumberData{Value: &age}},
				{ID: "age_check", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr(">=18")}},
				{ID: "profile_api", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ Fetched user profile")}},
				{ID: "sports_api", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ Registered for sports")}},
				{ID: "education_api", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ Registered for education")}},
			},
			Edges: []types.Edge{
				{Source: "user_age", Target: "age_check"},
				// True path (age >= 18): profile -> sports
				{Source: "age_check", Target: "profile_api", SourceHandle: strPtr("true")},
				{Source: "profile_api", Target: "sports_api"},
				// False path (age < 18): education
				{Source: "age_check", Target: "education_api", SourceHandle: strPtr("false")},
			},
		}

		eng, err := engine.New(mustMarshal(payload))
		if err != nil {
			fmt.Printf("❌ Error creating engine: %v\n", err)
			continue
		}

		result, err := eng.Execute()
		if err != nil {
			fmt.Printf("❌ Execution error: %v\n", err)
			continue
		}

		// Show which nodes executed
		fmt.Println("  Executed nodes:")
		for nodeID := range result.NodeResults {
			if nodeID == "user_age" || nodeID == "age_check" {
				continue // Skip setup nodes
			}
			if resultMap, ok := result.NodeResults[nodeID].(map[string]interface{}); ok {
				if text, ok := resultMap["text"].(string); ok {
					fmt.Printf("    - %s: %s\n", nodeID, text)
				}
			} else if text, ok := result.NodeResults[nodeID].(string); ok {
				fmt.Printf("    - %s: %s\n", nodeID, text)
			}
		}

		// Show which nodes DIDN'T execute
		allNodes := map[string]bool{
			"profile_api":    true,
			"sports_api":     true,
			"education_api":  true,
		}
		fmt.Println("  Skipped nodes:")
		for nodeID := range allNodes {
			if _, executed := result.NodeResults[nodeID]; !executed {
				fmt.Printf("    - %s (not in active path)\n", nodeID)
			}
		}
		fmt.Println()
	}

	fmt.Println()
}

func demo2SwitchRouting() {
	fmt.Println("📋 DEMO 2: HTTP Status Code Routing with Switch")
	fmt.Println("------------------------------------------------")
	fmt.Println("Scenario: Route to different handlers based on HTTP status code")
	fmt.Println()

	statusCodes := []float64{200, 404, 500}

	for _, code := range statusCodes {
		fmt.Printf("Testing with status_code = %.0f:\n", code)

		payload := types.Payload{
			Nodes: []types.Node{
				{ID: "status_code", Type: types.NodeTypeNumber, Data: types.NumberData{Value: &code}},
				{ID: "router", Type: types.NodeTypeSwitch, Data: types.SwitchData{
					Cases: []types.SwitchCase{
						{When: "==200", Value: 200.0, OutputPath: strPtr("success")},
						{When: "==404", Value: 404.0, OutputPath: strPtr("not_found")},
						{When: ">=500", OutputPath: strPtr("error")},
					},
					DefaultPath: strPtr("other"),
				}},
				{ID: "success_handler", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ Processed successful response")}},
				{ID: "error_handler", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ Logged server error")}},
				{ID: "not_found_handler", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ Handled not found")}},
				{ID: "other_handler", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ Other status code")}},
			},
			Edges: []types.Edge{
				{Source: "status_code", Target: "router"},
				{Source: "router", Target: "success_handler", SourceHandle: strPtr("success")},
				{Source: "router", Target: "error_handler", SourceHandle: strPtr("error")},
				{Source: "router", Target: "not_found_handler", SourceHandle: strPtr("not_found")},
				{Source: "router", Target: "other_handler", SourceHandle: strPtr("other")},
			},
		}

		eng, err := engine.New(mustMarshal(payload))
		if err != nil {
			fmt.Printf("❌ Error creating engine: %v\n", err)
			continue
		}

		result, err := eng.Execute()
		if err != nil {
			fmt.Printf("❌ Execution error: %v\n", err)
			continue
		}

		// Show which handler executed
		handlers := []string{"success_handler", "error_handler", "not_found_handler", "other_handler"}
		fmt.Println("  Executed handlers:")
		for _, handler := range handlers {
			if resultData, ok := result.NodeResults[handler]; ok {
				if resultMap, ok := resultData.(map[string]interface{}); ok {
					if text, ok := resultMap["text"].(string); ok {
						fmt.Printf("    ✅ %s: %s\n", handler, text)
					}
				}
			}
		}

		fmt.Println("  Skipped handlers:")
		for _, handler := range handlers {
			if _, executed := result.NodeResults[handler]; !executed {
				fmt.Printf("    ⏭️  %s (not in active path)\n", handler)
			}
		}
		fmt.Println()
	}

	fmt.Println()
}

func demo3NestedConditions() {
	fmt.Println("📋 DEMO 3: Nested Conditional Logic")
	fmt.Println("------------------------------------")
	fmt.Println("Scenario: Age >= 18 AND country == 'US' -> special_offer")
	fmt.Println("          Age >= 18 AND country != 'US' -> standard_offer")
	fmt.Println("          Age < 18 -> parental_consent")
	fmt.Println()

	testCases := []struct {
		age     float64
		country string
	}{
		{25, "US"},
		{25, "UK"},
		{15, "US"},
	}

	for _, tc := range testCases {
		fmt.Printf("Testing with age = %.0f, country = %s:\n", tc.age, tc.country)

		payload := types.Payload{
			Nodes: []types.Node{
				{ID: "user_age", Type: types.NodeTypeNumber, Data: types.NumberData{Value: &tc.age}},
				{ID: "user_country", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: &tc.country}},
				{ID: "age_check", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr(">=18")}},
				{ID: "country_check", Type: types.NodeTypeCondition, Data: types.ConditionData{Condition: strPtr("=='US'")}},
				{ID: "special_offer", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ US special offer applied!")}},
				{ID: "standard_offer", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ Standard offer applied!")}},
				{ID: "parental_consent", Type: types.NodeTypeTextInput, Data: types.TextInputData{Text: strPtr("✓ Parental consent required")}},
			},
			Edges: []types.Edge{
				{Source: "user_age", Target: "age_check"},
				{Source: "user_country", Target: "country_check"},
				// Adult path: check country
				{Source: "age_check", Target: "country_check", SourceHandle: strPtr("true")},
				{Source: "country_check", Target: "special_offer", SourceHandle: strPtr("true")},
				{Source: "country_check", Target: "standard_offer", SourceHandle: strPtr("false")},
				// Minor path: parental consent
				{Source: "age_check", Target: "parental_consent", SourceHandle: strPtr("false")},
			},
		}

		eng, err := engine.New(mustMarshal(payload))
		if err != nil {
			fmt.Printf("❌ Error creating engine: %v\n", err)
			continue
		}

		result, err := eng.Execute()
		if err != nil {
			fmt.Printf("❌ Execution error: %v\n", err)
			continue
		}

		// Show which action executed
		actions := []string{"special_offer", "standard_offer", "parental_consent"}
		fmt.Println("  Result:")
		for _, action := range actions {
			if resultData, ok := result.NodeResults[action]; ok {
				if resultMap, ok := resultData.(map[string]interface{}); ok {
					if text, ok := resultMap["text"].(string); ok {
						fmt.Printf("    ✅ %s\n", text)
					}
				}
			}
		}
		fmt.Println()
	}
}

func strPtr(s string) *string {
	return &s
}

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	return b
}
