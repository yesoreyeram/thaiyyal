package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"

	workflow "github.com/yesoreyeram/thaiyyal/backend"
)

func main() {
	fmt.Println("=== HTTP Pagination Using Composable Nodes ===\n")

	// Create a mock server that simulates a paginated API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}
		pageNum, _ := strconv.Atoi(page)

		// Simulate different responses for different pages
		var items []string
		for i := 1; i <= 10; i++ {
			itemNum := (pageNum-1)*10 + i
			items = append(items, fmt.Sprintf("item-%d", itemNum))
		}

		response := map[string]interface{}{
			"page":  pageNum,
			"items": items,
			"total": 50,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	fmt.Printf("Mock API server running at: %s\n", server.URL)

	// Example: Fetch 5 pages using composable nodes
	// Pattern: Counter → Condition → HTTP → Extract → Accumulator
	fmt.Println("\n--- Example: Fetch 5 pages, 10 items each (50 items total) ---")
	
	exampleComposable := fmt.Sprintf(`{
		"nodes": [
			{
				"id": "page_counter",
				"type": "counter",
				"data": {
					"counter_op": "increment",
					"delta": 1,
					"label": "Page Counter"
				}
			},
			{
				"id": "check_page_limit",
				"type": "condition",
				"data": {
					"condition": "<=5",
					"label": "Check if page <= 5"
				}
			},
			{
				"id": "http_request",
				"type": "http",
				"data": {
					"url": "%s?page=1",
					"label": "Fetch Page"
				}
			},
			{
				"id": "visualize",
				"type": "visualization",
				"data": {
					"mode": "text",
					"label": "Show Results"
				}
			}
		],
		"edges": [
			{"id": "e1", "source": "page_counter", "target": "check_page_limit"},
			{"id": "e2", "source": "check_page_limit", "target": "http_request"},
			{"id": "e3", "source": "http_request", "target": "visualize"}
		]
	}`, server.URL)

	runExample(exampleComposable, "Composable Pagination Pattern")

	fmt.Println("\n=== Key Benefits of Composable Approach ===")
	fmt.Println("✅ Flexible: Easy to customize pagination logic")
	fmt.Println("✅ Reusable: Same nodes work for other patterns")
	fmt.Println("✅ Extensible: Add retry, rate limiting, error handling")
	fmt.Println("✅ Clear: Each node has a single responsibility")
	
	fmt.Println("\n=== Common Pagination Patterns ===")
	fmt.Println("1. Page-based: counter → condition → HTTP (current example)")
	fmt.Println("2. Offset-based: counter → multiply → HTTP")
	fmt.Println("3. Cursor-based: variable → HTTP → extract → variable")
	fmt.Println("4. Until-empty: HTTP → extract → condition (check if empty)")
}

func runExample(payload string, description string) {
	engine, err := workflow.NewEngine([]byte(payload))
	if err != nil {
		log.Printf("Error creating engine: %v\n", err)
		return
	}

	result, err := engine.Execute()
	if err != nil {
		log.Printf("Error executing workflow: %v\n", err)
		return
	}

	// Pretty print the result
	fmt.Printf("\n%s:\n", description)
	prettyJSON, _ := json.MarshalIndent(result.NodeResults, "", "  ")
	fmt.Println(string(prettyJSON))
}
