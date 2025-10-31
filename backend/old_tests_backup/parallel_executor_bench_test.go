package workflow

import (
	"fmt"
	"testing"
)

// ============================================================================
// Parallel Execution Benchmarks
// ============================================================================
// These benchmarks demonstrate the performance improvements from parallel execution.
// Expected results:
// - Linear workflows: Similar performance (no parallelism opportunities)
// - Branching workflows: 2-10x speedup depending on branch count
// - Complex workflows: Significant improvements with multiple independent paths
// ============================================================================

// BenchmarkSequentialVsParallel_SimpleBranch benchmarks 2 parallel branches
func BenchmarkSequentialVsParallel_SimpleBranch(b *testing.B) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 20}},
			{"id": "3", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "3"},
			{"source": "2", "target": "3"}
		]
	}`

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			_, err := engine.Execute()
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			config := DefaultParallelConfig()
			_, err := engine.ExecuteWithParallelism(config)
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})
}

// BenchmarkSequentialVsParallel_MultipleBranches benchmarks 10 parallel branches
func BenchmarkSequentialVsParallel_MultipleBranches(b *testing.B) {
	// Generate workflow with 10 independent input nodes merging into one
	nodes := `[`
	for i := 1; i <= 10; i++ {
		nodes += fmt.Sprintf(`{"id": "%d", "data": {"value": %d}}`, i, i*10)
		if i < 10 {
			nodes += `,`
		}
	}
	nodes += `,{"id": "merge", "data": {"op": "add"}}]`

	edges := `[`
	for i := 1; i <= 10; i++ {
		edges += fmt.Sprintf(`{"source": "%d", "target": "merge"}`, i)
		if i < 10 {
			edges += `,`
		}
	}
	edges += `]`

	payload := fmt.Sprintf(`{"nodes": %s, "edges": %s}`, nodes, edges)

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			_, err := engine.Execute()
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			config := DefaultParallelConfig()
			_, err := engine.ExecuteWithParallelism(config)
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})
}

// BenchmarkSequentialVsParallel_ComplexWorkflow benchmarks a complex multi-level workflow
func BenchmarkSequentialVsParallel_ComplexWorkflow(b *testing.B) {
	payload := `{
		"nodes": [
			{"id": "n1", "data": {"value": 10}},
			{"id": "n2", "data": {"value": 20}},
			{"id": "n3", "data": {"value": 30}},
			{"id": "n4", "data": {"value": 40}},
			{"id": "n5", "data": {"value": 50}},
			{"id": "n6", "data": {"value": 60}},
			{"id": "add1", "data": {"op": "add"}},
			{"id": "add2", "data": {"op": "add"}},
			{"id": "add3", "data": {"op": "add"}},
			{"id": "mult1", "data": {"op": "multiply"}},
			{"id": "mult2", "data": {"op": "multiply"}},
			{"id": "final", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "n1", "target": "add1"},
			{"source": "n2", "target": "add1"},
			{"source": "n3", "target": "add2"},
			{"source": "n4", "target": "add2"},
			{"source": "n5", "target": "add3"},
			{"source": "n6", "target": "add3"},
			{"source": "add1", "target": "mult1"},
			{"source": "add2", "target": "mult1"},
			{"source": "add3", "target": "mult2"},
			{"source": "mult1", "target": "final"},
			{"source": "mult2", "target": "final"}
		]
	}`

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			_, err := engine.Execute()
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			config := DefaultParallelConfig()
			_, err := engine.ExecuteWithParallelism(config)
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})
}

// BenchmarkSequentialVsParallel_LinearWorkflow benchmarks a linear workflow (no parallelism)
func BenchmarkSequentialVsParallel_LinearWorkflow(b *testing.B) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"op": "add"}},
			{"id": "3", "data": {"op": "multiply"}},
			{"id": "4", "data": {"op": "add"}},
			{"id": "5", "data": {"mode": "text"}}
		],
		"edges": [
			{"source": "1", "target": "2"},
			{"source": "2", "target": "3"},
			{"source": "3", "target": "4"},
			{"source": "4", "target": "5"}
		]
	}`

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			_, err := engine.Execute()
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			config := DefaultParallelConfig()
			_, err := engine.ExecuteWithParallelism(config)
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})
}

// BenchmarkSequentialVsParallel_DiamondPattern benchmarks diamond-shaped workflow
func BenchmarkSequentialVsParallel_DiamondPattern(b *testing.B) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 100}},
			{"id": "2", "data": {"op": "add"}},
			{"id": "3", "data": {"op": "multiply"}},
			{"id": "4", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "2"},
			{"source": "1", "target": "3"},
			{"source": "2", "target": "4"},
			{"source": "3", "target": "4"}
		]
	}`

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			_, err := engine.Execute()
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			config := DefaultParallelConfig()
			_, err := engine.ExecuteWithParallelism(config)
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})
}

// BenchmarkSequentialVsParallel_TextOperations benchmarks text operation workflows
func BenchmarkSequentialVsParallel_TextOperations(b *testing.B) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"text": "hello"}},
			{"id": "2", "data": {"text": "world"}},
			{"id": "3", "data": {"text": "foo"}},
			{"id": "4", "data": {"text": "bar"}},
			{"id": "up1", "data": {"text_op": "uppercase"}},
			{"id": "up2", "data": {"text_op": "uppercase"}},
			{"id": "up3", "data": {"text_op": "uppercase"}},
			{"id": "up4", "data": {"text_op": "uppercase"}},
			{"id": "concat", "data": {"text_op": "concat", "separator": " "}}
		],
		"edges": [
			{"source": "1", "target": "up1"},
			{"source": "2", "target": "up2"},
			{"source": "3", "target": "up3"},
			{"source": "4", "target": "up4"},
			{"source": "up1", "target": "concat"},
			{"source": "up2", "target": "concat"},
			{"source": "up3", "target": "concat"},
			{"source": "up4", "target": "concat"}
		]
	}`

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			_, err := engine.Execute()
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine, _ := NewEngine([]byte(payload))
			config := DefaultParallelConfig()
			_, err := engine.ExecuteWithParallelism(config)
			if err != nil {
				b.Fatalf("Execution failed: %v", err)
			}
		}
	})
}

// BenchmarkConcurrencyLimits benchmarks different concurrency limits
func BenchmarkConcurrencyLimits(b *testing.B) {
	// Generate workflow with 20 independent branches
	nodes := `[`
	for i := 1; i <= 20; i++ {
		nodes += fmt.Sprintf(`{"id": "%d", "data": {"value": %d}}`, i, i)
		if i < 20 {
			nodes += `,`
		}
	}
	nodes += `,{"id": "merge", "data": {"op": "add"}}]`

	edges := `[`
	for i := 1; i <= 20; i++ {
		edges += fmt.Sprintf(`{"source": "%d", "target": "merge"}`, i)
		if i < 20 {
			edges += `,`
		}
	}
	edges += `]`

	payload := fmt.Sprintf(`{"nodes": %s, "edges": %s}`, nodes, edges)

	limits := []int{1, 2, 4, 8, 0} // 0 means unlimited

	for _, limit := range limits {
		name := fmt.Sprintf("Limit_%d", limit)
		if limit == 0 {
			name = "Unlimited"
		}
		
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				engine, _ := NewEngine([]byte(payload))
				config := ParallelExecutionConfig{
					MaxConcurrency: limit,
					EnableParallel: true,
				}
				_, err := engine.ExecuteWithParallelism(config)
				if err != nil {
					b.Fatalf("Execution failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkLevelComputation benchmarks the level computation algorithm
func BenchmarkLevelComputation(b *testing.B) {
	// Generate a complex workflow with multiple levels
	nodes := `[`
	for i := 1; i <= 50; i++ {
		nodes += fmt.Sprintf(`{"id": "n%d", "data": {"value": %d}}`, i, i)
		if i < 50 {
			nodes += `,`
		}
	}
	nodes += `]`

	// Create a multi-level dependency structure
	edges := `[`
	edgeCount := 0
	for i := 1; i <= 25; i++ {
		for j := 26; j <= 50; j++ {
			if (i+j)%5 == 0 { // Create some dependencies
				if edgeCount > 0 {
					edges += `,`
				}
				edges += fmt.Sprintf(`{"source": "n%d", "target": "n%d"}`, i, j)
				edgeCount++
			}
		}
	}
	edges += `]`

	payload := fmt.Sprintf(`{"nodes": %s, "edges": %s}`, nodes, edges)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine, _ := NewEngine([]byte(payload))
		engine.inferNodeTypes()
		_, err := engine.computeExecutionLevels()
		if err != nil {
			b.Fatalf("Level computation failed: %v", err)
		}
	}
}

// BenchmarkParallelExecutor_MemoryAllocation benchmarks memory allocation
func BenchmarkParallelExecutor_MemoryAllocation(b *testing.B) {
	payload := `{
		"nodes": [
			{"id": "1", "data": {"value": 10}},
			{"id": "2", "data": {"value": 20}},
			{"id": "3", "data": {"value": 30}},
			{"id": "4", "data": {"value": 40}},
			{"id": "5", "data": {"value": 50}},
			{"id": "merge", "data": {"op": "add"}}
		],
		"edges": [
			{"source": "1", "target": "merge"},
			{"source": "2", "target": "merge"},
			{"source": "3", "target": "merge"},
			{"source": "4", "target": "merge"},
			{"source": "5", "target": "merge"}
		]
	}`

	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		engine, _ := NewEngine([]byte(payload))
		config := DefaultParallelConfig()
		_, err := engine.ExecuteWithParallelism(config)
		if err != nil {
			b.Fatalf("Execution failed: %v", err)
		}
	}
}
