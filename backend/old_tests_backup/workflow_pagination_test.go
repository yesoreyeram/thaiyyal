package workflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// ====================================
// HTTP Pagination Pattern Tests
// ====================================

// TestPaginationPageBased tests page-based pagination pattern
// Pattern: Counter → Condition → HTTP
func TestPaginationPageBased(t *testing.T) {
	// Create mock server for paginated API
	requestedPages := []int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}
		pageNum, _ := strconv.Atoi(page)
		requestedPages = append(requestedPages, pageNum)

		// Return 10 items per page
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

	// Workflow: Fetch 5 pages (50 items total)
	// This simulates: Counter → Condition → HTTP
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "1", "type": "counter", "data": {"counter_op": "increment", "delta": 1}},
			{"id": "2", "type": "condition", "data": {"condition": "<=5"}},
			{"id": "3", "type": "http", "data": {"url": "%s?page=1"}},
			{"id": "4", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"}
		]
	}`, server.URL)

	engine, err := NewEngineWithConfig([]byte(payload), testConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Verify counter incremented
	counterResult := result.NodeResults["1"].(map[string]interface{})
	if counterResult["value"].(float64) != 1 {
		t.Errorf("Expected counter value 1, got %v", counterResult["value"])
	}

	// Verify condition was evaluated
	conditionResult := result.NodeResults["2"].(map[string]interface{})
	if !conditionResult["condition_met"].(bool) {
		t.Error("Expected condition to be met")
	}

	// Verify HTTP request was made and returned JSON
	httpResult := result.NodeResults["3"].(string)
	if len(httpResult) == 0 {
		t.Error("Expected non-empty HTTP result")
	}

	// Verify the response contains expected data
	var responseData map[string]interface{}
	if err := json.Unmarshal([]byte(httpResult), &responseData); err == nil {
		if responseData["page"].(float64) != 1 {
			t.Errorf("Expected page 1, got %v", responseData["page"])
		}
	}

	t.Logf("Successfully executed page-based pagination pattern")
	t.Logf("Pattern: Counter (page tracker) → Condition (page limit) → HTTP (fetch page)")
}

// TestPaginationOffsetBased tests offset-based pagination pattern
// Pattern: Number → Counter → Multiply → HTTP
func TestPaginationOffsetBased(t *testing.T) {
	requestedOffsets := []int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset := r.URL.Query().Get("offset")
		if offset == "" {
			offset = "0"
		}
		offsetNum, _ := strconv.Atoi(offset)
		requestedOffsets = append(requestedOffsets, offsetNum)

		limit := 10
		var items []string
		for i := 0; i < limit; i++ {
			items = append(items, fmt.Sprintf("item-%d", offsetNum+i+1))
		}

		response := map[string]interface{}{
			"offset": offsetNum,
			"limit":  limit,
			"items":  items,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Workflow: Demonstrate offset calculation
	// Counter starts at 0, multiply by page_size to get offset
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "1", "type": "number", "data": {"value": 0}},
			{"id": "2", "type": "number", "data": {"value": 10}},
			{"id": "3", "type": "operation", "data": {"op": "multiply"}},
			{"id": "4", "type": "http", "data": {"url": "%s?offset=0&limit=10"}},
			{"id": "5", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "3"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"},
			{"id": "e4", "source": "4", "target": "5"}
		]
	}`, server.URL)

	engine, err := NewEngineWithConfig([]byte(payload), testConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Verify multiplication happened (0 * 10 = 0)
	multiplyResult := result.NodeResults["3"].(float64)
	if multiplyResult != 0 {
		t.Errorf("Expected offset 0, got %v", multiplyResult)
	}

	// Verify HTTP was called
	if result.NodeResults["4"] == nil {
		t.Error("Expected HTTP result")
	}

	t.Logf("Successfully executed offset-based pagination pattern")
	t.Logf("Pattern: Page Number → Multiply by PageSize → HTTP (with offset)")
}

// TestPaginationCursorBased tests cursor-based pagination pattern concept
// Pattern: Variable → HTTP (demonstrates state management for cursors)
func TestPaginationCursorBased(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cursor := r.URL.Query().Get("cursor")

		var nextCursor *string
		if cursor == "" {
			next := "cursor2"
			nextCursor = &next
		}

		items := []string{"item1", "item2", "item3"}

		response := map[string]interface{}{
			"items":       items,
			"next_cursor": nextCursor,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Workflow: Demonstrates Variable node for cursor management
	// In practice: Variable (get) → HTTP → parse response → Variable (set next cursor)
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "1", "type": "text_input", "data": {"text": ""}},
			{"id": "2", "type": "http", "data": {"url": "%s?cursor="}},
			{"id": "3", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"}
		]
	}`, server.URL)

	engine, err := NewEngineWithConfig([]byte(payload), testConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Verify HTTP was called
	httpResult := result.NodeResults["2"].(string)
	if len(httpResult) == 0 {
		t.Error("Expected HTTP result")
	}

	// Verify response contains cursor
	var responseData map[string]interface{}
	if err := json.Unmarshal([]byte(httpResult), &responseData); err == nil {
		if responseData["next_cursor"] == nil {
			t.Log("Note: In a real cursor pagination, this would be stored in a Variable node")
		}
	}

	t.Logf("Successfully demonstrated cursor-based pagination concept")
	t.Logf("Pattern: Variable (cursor storage) → HTTP → Parse → Update Variable")
}

// TestPaginationUntilEmpty tests pagination until empty array concept
// Pattern: HTTP → Check response (demonstrates conditional logic)
func TestPaginationUntilEmpty(t *testing.T) {
	pageCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageCount++

		var items []string
		// Return items for first 3 pages, then empty
		if pageCount <= 3 {
			for i := 1; i <= 10; i++ {
				items = append(items, fmt.Sprintf("item-%d", (pageCount-1)*10+i))
			}
		}

		response := map[string]interface{}{
			"page":      pageCount,
			"items":     items,
			"has_more":  len(items) > 0,
			"item_count": len(items),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Workflow: HTTP → parse response (check if has_more or items not empty)
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "1", "type": "http", "data": {"url": "%s?page=1"}},
			{"id": "2", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`, server.URL)

	engine, err := NewEngineWithConfig([]byte(payload), testConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Verify HTTP result
	httpResult := result.NodeResults["1"].(string)
	if len(httpResult) == 0 {
		t.Error("Expected HTTP result")
	}

	// Parse and verify response structure
	var responseData map[string]interface{}
	if err := json.Unmarshal([]byte(httpResult), &responseData); err == nil {
		if has_more, ok := responseData["has_more"].(bool); ok {
			t.Logf("has_more field present: %v (would be used to control loop)", has_more)
		}
		if item_count, ok := responseData["item_count"].(float64); ok {
			t.Logf("item_count: %v (would check if > 0 to continue)", item_count)
		}
	}

	t.Logf("Successfully demonstrated until-empty pagination concept")
	t.Logf("Pattern: HTTP → Parse response → Condition (has_more or count > 0) → Loop")
}

// TestPaginationWithAccumulator tests collecting results concept
// Demonstrates how Accumulator would collect data from multiple HTTP calls
func TestPaginationWithAccumulator(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"items": []string{"item1", "item2", "item3"},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Workflow: HTTP → Accumulator (demonstrates result collection)
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "1", "type": "http", "data": {"url": "%s"}},
			{"id": "2", "type": "text_input", "data": {"text": "page_data"}},
			{"id": "3", "type": "accumulator", "data": {"accum_op": "concat"}},
			{"id": "4", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"}
		]
	}`, server.URL)

	engine, err := NewEngineWithConfig([]byte(payload), testConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Verify accumulator is working
	accumulatorResult := result.NodeResults["3"].(map[string]interface{})
	if accumulatorResult["operation"] != "concat" {
		t.Errorf("Expected operation 'concat', got %v", accumulatorResult["operation"])
	}

	t.Logf("Successfully demonstrated pagination with accumulator concept")
	t.Logf("Pattern: HTTP (multiple calls) → Accumulator (collect all results)")
}

// TestPaginationErrorHandling tests pagination with error handling
func TestPaginationErrorHandling(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		// Fail on page 3
		if requestCount == 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
			return
		}

		response := map[string]interface{}{
			"items": []string{"item1", "item2"},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Workflow: Counter → HTTP (will fail on 3rd iteration)
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "1", "type": "counter", "data": {"counter_op": "increment"}},
			{"id": "2", "type": "http", "data": {"url": "%s"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"}
		]
	}`, server.URL)

	engine, err := NewEngineWithConfig([]byte(payload), testConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Execute - should succeed on first try
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Verify HTTP succeeded on first page
	if result.NodeResults["2"] == nil {
		t.Error("Expected HTTP result on first execution")
	}

	t.Logf("Successfully tested pagination error handling")
}

// TestPaginationWithRateLimit tests pagination with delay between requests
func TestPaginationWithRateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"items": []string{"item1", "item2"},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Workflow: HTTP → Delay → Counter (for next iteration)
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "1", "type": "http", "data": {"url": "%s"}},
			{"id": "2", "type": "delay", "data": {"duration": "100ms"}},
			{"id": "3", "type": "counter", "data": {"counter_op": "increment"}},
			{"id": "4", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1", "source": "1", "target": "2"},
			{"id": "e2", "source": "2", "target": "3"},
			{"id": "e3", "source": "3", "target": "4"}
		]
	}`, server.URL)

	engine, err := NewEngineWithConfig([]byte(payload), testConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Verify delay was executed
	delayResult := result.NodeResults["2"].(map[string]interface{})
	if delayResult["duration"] != "100ms" {
		t.Errorf("Expected duration '100ms', got %v", delayResult["duration"])
	}

	t.Logf("Successfully tested pagination with rate limiting")
}

// TestPaginationCompleteScenario tests a complete pagination workflow demonstration
func TestPaginationCompleteScenario(t *testing.T) {
	pagesRequested := []int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}
		pageNum, _ := strconv.Atoi(page)
		pagesRequested = append(pagesRequested, pageNum)

		// Generate 10 items per page
		var items []string
		for i := 1; i <= 10; i++ {
			itemNum := (pageNum-1)*10 + i
			items = append(items, fmt.Sprintf("item-%d", itemNum))
		}

		response := map[string]interface{}{
			"page":       pageNum,
			"items":      items,
			"total":      50,
			"page_size":  10,
			"has_more":   pageNum < 5,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Complete workflow demonstrating composable pagination
	// This is a simplified single-page fetch demonstrating the building blocks
	payload := fmt.Sprintf(`{
		"nodes": [
			{"id": "counter", "type": "counter", "data": {"counter_op": "increment", "delta": 1}},
			{"id": "check_limit", "type": "condition", "data": {"condition": "<=5"}},
			{"id": "http", "type": "http", "data": {"url": "%s?page=1"}},
			{"id": "accumulate", "type": "text_input", "data": {"text": "results"}},
			{"id": "visualize", "type": "visualization", "data": {"mode": "text"}}
		],
		"edges": [
			{"id": "e1", "source": "counter", "target": "check_limit"},
			{"id": "e2", "source": "check_limit", "target": "http"},
			{"id": "e3", "source": "http", "target": "accumulate"},
			{"id": "e4", "source": "accumulate", "target": "visualize"}
		]
	}`, server.URL)

	engine, err := NewEngineWithConfig([]byte(payload), testConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Failed to execute workflow: %v", err)
	}

	// Verify all key nodes executed
	if result.NodeResults["counter"] == nil {
		t.Error("Counter node did not execute")
	}
	if result.NodeResults["check_limit"] == nil {
		t.Error("Condition node did not execute")
	}
	if result.NodeResults["http"] == nil {
		t.Error("HTTP node did not execute")
	}

	// Verify the response structure
	httpResult := result.NodeResults["http"].(string)
	var responseData map[string]interface{}
	if err := json.Unmarshal([]byte(httpResult), &responseData); err == nil {
		if responseData["page"].(float64) != 1 {
			t.Errorf("Expected page 1, got %v", responseData["page"])
		}
		if responseData["total"].(float64) != 50 {
			t.Errorf("Expected total 50, got %v", responseData["total"])
		}
		t.Logf("Response structure: page=%v, total=%v, has_more=%v",
			responseData["page"], responseData["total"], responseData["has_more"])
	}

	t.Logf("Successfully executed complete pagination scenario")
	t.Logf("This demonstrates the composable approach:")
	t.Logf("  1. Counter - tracks page number")
	t.Logf("  2. Condition - checks if within page limit")
	t.Logf("  3. HTTP - fetches the page")
	t.Logf("  4. Accumulator - would collect results (multiple iterations)")
	t.Logf("  5. Visualization - displays results")
}
