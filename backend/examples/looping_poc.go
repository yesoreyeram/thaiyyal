package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Looping Patterns Proof of Concept ===\n")

	// 1. For-Each Loop Example
	forEachExample()

	// 2. Retry Loop with Backoff
	retryLoopExample()

	// 3. Parallel Processing (Fan-out)
	parallelProcessingExample()

	// 4. Consolidation/Aggregation
	consolidationExample()

	// 5. Time Window Iteration
	timeWindowExample()

	// 6. Conditional Loop (While)
	conditionalLoopExample()
}

// 1. For-Each Loop: Process array items independently
func forEachExample() {
	fmt.Println("=== 1. For-Each Loop Example ===")

	metrics := []string{"cpu_usage", "memory_usage", "disk_usage", "network_io"}
	results := make([]string, len(metrics))

	// Simulate processing each metric
	for i, metric := range metrics {
		result := processMetric(metric)
		results[i] = result
		fmt.Printf("  Processed: %s -> %s\n", metric, result)
	}

	fmt.Printf("Results: %v\n\n", results)
}

func processMetric(metric string) string {
	// Simulate some processing
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("%s_transformed", metric)
}

// 2. Retry Loop with Exponential Backoff
func retryLoopExample() {
	fmt.Println("=== 2. Retry Loop with Backoff ===")

	maxAttempts := 5
	attempt := 0
	var result string
	var err error

	for attempt < maxAttempts {
		attempt++
		result, err = unreliableOperation(attempt)

		if err == nil {
			fmt.Printf("  Attempt %d: SUCCESS - %s\n", attempt, result)
			break
		}

		fmt.Printf("  Attempt %d: FAILED - %s\n", attempt, err)

		if attempt < maxAttempts {
			backoff := time.Duration(1<<uint(attempt-1)) * 100 * time.Millisecond
			fmt.Printf("  Waiting %v before retry...\n", backoff)
			time.Sleep(backoff)
		}
	}

	if err != nil {
		fmt.Printf("  All attempts failed\n")
	}
	fmt.Println()
}

func unreliableOperation(attempt int) (string, error) {
	// Succeeds on 3rd attempt
	if attempt < 3 {
		return "", fmt.Errorf("temporary failure")
	}
	return "Data retrieved successfully", nil
}

// 3. Parallel Processing (Fan-out/Fan-in)
func parallelProcessingExample() {
	fmt.Println("=== 3. Parallel Processing Example ===")

	sources := []string{"prometheus", "loki", "elasticsearch", "clickhouse"}
	results := make(chan MetricResult, len(sources))
	var wg sync.WaitGroup

	// Fan-out: query all sources in parallel
	for _, source := range sources {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			result := queryDataSource(s)
			results <- result
		}(source)
	}

	// Wait for all to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Fan-in: collect all results
	totalMetrics := 0
	fmt.Println("  Collecting results from parallel queries:")
	for result := range results {
		fmt.Printf("    %s: %d metrics (took %v)\n", result.Source, result.Count, result.Duration)
		totalMetrics += result.Count
	}

	fmt.Printf("  Total metrics collected: %d\n\n", totalMetrics)
}

type MetricResult struct {
	Source   string
	Count    int
	Duration time.Duration
}

func queryDataSource(source string) MetricResult {
	start := time.Now()
	// Simulate query with random delay
	time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)

	return MetricResult{
		Source:   source,
		Count:    rand.Intn(100) + 50,
		Duration: time.Since(start),
	}
}

// 4. Consolidation/Aggregation Pattern
func consolidationExample() {
	fmt.Println("=== 4. Consolidation/Aggregation Example ===")

	// Simulate collecting metrics from different sources
	cpuMetrics := collectCPUMetrics()
	memMetrics := collectMemoryMetrics()
	diskMetrics := collectDiskMetrics()

	// Consolidate into single report
	report := ConsolidatedReport{
		Timestamp: time.Now(),
		CPU:       cpuMetrics,
		Memory:    memMetrics,
		Disk:      diskMetrics,
	}

	// Analyze and alert
	fmt.Printf("  Consolidated Metrics at %s:\n", report.Timestamp.Format("15:04:05"))
	fmt.Printf("    CPU Usage: %.1f%%\n", report.CPU.Average)
	fmt.Printf("    Memory Usage: %.1f%%\n", report.Memory.Average)
	fmt.Printf("    Disk Usage: %.1f%%\n", report.Disk.Average)

	// Check thresholds
	alerts := []string{}
	if report.CPU.Average > 80 {
		alerts = append(alerts, "CPU usage high")
	}
	if report.Memory.Average > 85 {
		alerts = append(alerts, "Memory usage high")
	}
	if report.Disk.Average > 90 {
		alerts = append(alerts, "Disk usage critical")
	}

	if len(alerts) > 0 {
		fmt.Printf("  Alerts: %v\n", alerts)
	} else {
		fmt.Printf("  Status: All systems normal\n")
	}
	fmt.Println()
}

type ResourceMetrics struct {
	Average float64
	Max     float64
	Min     float64
}

type ConsolidatedReport struct {
	Timestamp time.Time
	CPU       ResourceMetrics
	Memory    ResourceMetrics
	Disk      ResourceMetrics
}

func collectCPUMetrics() ResourceMetrics {
	return ResourceMetrics{Average: 45.3, Max: 78.2, Min: 12.1}
}

func collectMemoryMetrics() ResourceMetrics {
	return ResourceMetrics{Average: 67.8, Max: 89.3, Min: 45.2}
}

func collectDiskMetrics() ResourceMetrics {
	return ResourceMetrics{Average: 54.2, Max: 72.1, Min: 32.8}
}

// 5. Time Window Iteration
func timeWindowExample() {
	fmt.Println("=== 5. Time Window Iteration Example ===")

	// Analyze metrics over 6 time windows (1 hour each)
	windows := 6
	now := time.Now()

	anomalies := []TimeWindowAnomaly{}

	for i := 0; i < windows; i++ {
		windowStart := now.Add(time.Duration(-windows+i) * time.Hour)
		windowEnd := windowStart.Add(1 * time.Hour)

		metrics := queryTimeWindow(windowStart, windowEnd)
		hasAnomaly := detectAnomaly(metrics)

		status := "normal"
		if hasAnomaly {
			status = "ANOMALY"
			anomalies = append(anomalies, TimeWindowAnomaly{
				Window: fmt.Sprintf("%s to %s", windowStart.Format("15:04"), windowEnd.Format("15:04")),
				Metric: metrics,
			})
		}

		fmt.Printf("  Window %d (%s - %s): avg=%.2f, status=%s\n",
			i+1, windowStart.Format("15:04"), windowEnd.Format("15:04"),
			metrics.Average, status)
	}

	if len(anomalies) > 0 {
		fmt.Printf("  Total anomalies detected: %d\n", len(anomalies))
	}
	fmt.Println()
}

type TimeWindowAnomaly struct {
	Window string
	Metric WindowMetrics
}

type WindowMetrics struct {
	Average float64
	StdDev  float64
}

func queryTimeWindow(start, end time.Time) WindowMetrics {
	// Simulate querying metrics for time window
	avg := 50.0 + rand.Float64()*30.0
	stdDev := 5.0 + rand.Float64()*5.0
	return WindowMetrics{Average: avg, StdDev: stdDev}
}

func detectAnomaly(metrics WindowMetrics) bool {
	// Simple anomaly detection: avg > 75 or stddev > 8
	return metrics.Average > 75 || metrics.StdDev > 8
}

// 6. Conditional Loop (While)
func conditionalLoopExample() {
	fmt.Println("=== 6. Conditional Loop (While) Example ===")

	maxIterations := 10
	iteration := 0
	threshold := 95.0
	currentValue := 0.0

	fmt.Printf("  Target: Reach value >= %.0f\n", threshold)

	for iteration < maxIterations && currentValue < threshold {
		iteration++
		// Simulate gradual increase
		increment := rand.Float64()*20.0 + 5.0
		currentValue += increment

		fmt.Printf("  Iteration %d: value=%.2f (increment=%.2f)\n", iteration, currentValue, increment)

		if currentValue >= threshold {
			fmt.Printf("  SUCCESS: Threshold reached in %d iterations\n", iteration)
			break
		}
	}

	if currentValue < threshold {
		fmt.Printf("  TIMEOUT: Max iterations reached, value=%.2f\n", currentValue)
	}
	fmt.Println()
}
