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
	fmt.Println("=== HTTP Pagination Example ===\n")

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

	// Example 1: Pagination with max_pages
	fmt.Println("\n--- Example 1: Fetch 5 pages (50 items, 10 per page) ---")
	example1 := fmt.Sprintf(`{
		"nodes": [
			{
				"id": "1",
				"type": "http_pagination",
				"data": {
					"base_url": "%s",
					"max_pages": 5
				}
			},
			{
				"id": "2",
				"type": "visualization",
				"data": {
					"mode": "text"
				}
			}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`, server.URL)

	runExample(example1, "Pagination Result")

	// Example 2: Pagination with total_items and page_size
	fmt.Println("\n--- Example 2: Using total_items and page_size ---")
	example2 := fmt.Sprintf(`{
		"nodes": [
			{
				"id": "1",
				"type": "http_pagination",
				"data": {
					"base_url": "%s",
					"total_items": 50,
					"page_size": 10
				}
			}
		],
		"edges": []
	}`, server.URL)

	runExample(example2, "Pagination with total_items")

	// Example 3: Pagination with URL placeholder
	fmt.Println("\n--- Example 3: Using URL placeholder {page} ---")
	example3 := fmt.Sprintf(`{
		"nodes": [
			{
				"id": "1",
				"type": "http_pagination",
				"data": {
					"base_url": "%s/api/items?p={page}",
					"max_pages": 3
				}
			}
		],
		"edges": []
	}`, server.URL)

	runExample(example3, "Pagination with placeholder")

	// Example 4: Custom start page
	fmt.Println("\n--- Example 4: Start from page 3 ---")
	example4 := fmt.Sprintf(`{
		"nodes": [
			{
				"id": "1",
				"type": "http_pagination",
				"data": {
					"base_url": "%s",
					"start_page": 3,
					"max_pages": 2
				}
			}
		],
		"edges": []
	}`, server.URL)

	runExample(example4, "Pagination starting at page 3")

	// Example 5: Error handling with break_on_error
	fmt.Println("\n--- Example 5: Error handling (simulated) ---")
	fmt.Println("This example shows how pagination handles errors.")
	fmt.Println("With break_on_error=true (default), it stops on first error.")
	fmt.Println("With break_on_error=false, it continues and collects partial results.")

	fmt.Println("\n=== All Examples Completed ===")
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
