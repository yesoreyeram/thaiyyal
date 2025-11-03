package engine

import (
	"encoding/json"
	"fmt"
	"testing"
)

// BenchmarkEngine_SimpleWorkflow benchmarks a simple workflow with basic operations
func BenchmarkEngine_SimpleWorkflow(b *testing.B) {
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "type": "number", "data": map[string]interface{}{"value": 10.0}},
			map[string]interface{}{"id": "2", "type": "number", "data": map[string]interface{}{"value": 20.0}},
			map[string]interface{}{"id": "3", "type": "operation", "data": map[string]interface{}{"op": "add"}},
		},
		"edges": []interface{}{
			map[string]interface{}{"source": "1", "target": "3"},
			map[string]interface{}{"source": "2", "target": "3"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		engine, err := New(jsonData)
		if err != nil {
			b.Fatalf("New failed: %v", err)
		}
		_, err = engine.Execute()
		if err != nil {
			b.Fatalf("Execute failed: %v", err)
		}
	}
}

// BenchmarkEngine_MultipleOperations benchmarks chained math operations
func BenchmarkEngine_MultipleOperations(b *testing.B) {
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "type": "number", "data": map[string]interface{}{"value": 10.0}},
			map[string]interface{}{"id": "2", "type": "number", "data": map[string]interface{}{"value": 5.0}},
			map[string]interface{}{"id": "3", "type": "operation", "data": map[string]interface{}{"op": "add"}},
			map[string]interface{}{"id": "4", "type": "number", "data": map[string]interface{}{"value": 3.0}},
			map[string]interface{}{"id": "5", "type": "operation", "data": map[string]interface{}{"op": "multiply"}},
			map[string]interface{}{"id": "6", "type": "number", "data": map[string]interface{}{"value": 2.0}},
			map[string]interface{}{"id": "7", "type": "operation", "data": map[string]interface{}{"op": "divide"}},
		},
		"edges": []interface{}{
			map[string]interface{}{"source": "1", "target": "3"},
			map[string]interface{}{"source": "2", "target": "3"},
			map[string]interface{}{"source": "3", "target": "5"},
			map[string]interface{}{"source": "4", "target": "5"},
			map[string]interface{}{"source": "5", "target": "7"},
			map[string]interface{}{"source": "6", "target": "7"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		engine, _ := New(jsonData)
		_, err := engine.Execute()
		if err != nil {
			b.Fatalf("Execute failed: %v", err)
		}
	}
}

// BenchmarkEngine_TextOperations benchmarks text processing
func BenchmarkEngine_TextOperations(b *testing.B) {
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "type": "text_input", "data": map[string]interface{}{"text": "hello world"}},
			map[string]interface{}{"id": "2", "type": "text_operation", "data": map[string]interface{}{"text_op": "uppercase"}},
			map[string]interface{}{"id": "3", "type": "text_operation", "data": map[string]interface{}{"text_op": "lowercase"}},
		},
		"edges": []interface{}{
			map[string]interface{}{"source": "1", "target": "2"},
			map[string]interface{}{"source": "2", "target": "3"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		engine, _ := New(jsonData)
		_, err := engine.Execute()
		if err != nil {
			b.Fatalf("Execute failed: %v", err)
		}
	}
}

// BenchmarkEngine_StateOperations benchmarks state management
func BenchmarkEngine_StateOperations(b *testing.B) {
	b.Run("variable", func(b *testing.B) {
		payload := map[string]interface{}{
			"nodes": []interface{}{
				map[string]interface{}{"id": "1", "type": "number", "data": map[string]interface{}{"value": 42.0}},
				map[string]interface{}{"id": "2", "type": "variable", "data": map[string]interface{}{"var_op": "set", "var_name": "x"}},
				map[string]interface{}{"id": "3", "type": "variable", "data": map[string]interface{}{"var_op": "get", "var_name": "x"}},
			},
			"edges": []interface{}{
				map[string]interface{}{"source": "1", "target": "2"},
				map[string]interface{}{"source": "2", "target": "3"},
			},
		}
		jsonData, _ := json.Marshal(payload)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			engine, _ := New(jsonData)
			_, err := engine.Execute()
			if err != nil {
				b.Fatalf("Execute failed: %v", err)
			}
		}
	})

	b.Run("counter", func(b *testing.B) {
		payload := map[string]interface{}{
			"nodes": []interface{}{
				map[string]interface{}{"id": "1", "type": "counter", "data": map[string]interface{}{"counter_op": "reset", "counter_value": 0.0}},
				map[string]interface{}{"id": "2", "type": "counter", "data": map[string]interface{}{"counter_op": "increment", "counter_delta": 1.0}},
				map[string]interface{}{"id": "3", "type": "counter", "data": map[string]interface{}{"counter_op": "get"}},
			},
			"edges": []interface{}{
				map[string]interface{}{"source": "1", "target": "2"},
				map[string]interface{}{"source": "2", "target": "3"},
			},
		}
		jsonData, _ := json.Marshal(payload)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			engine, _ := New(jsonData)
			_, err := engine.Execute()
			if err != nil {
				b.Fatalf("Execute failed: %v", err)
			}
		}
	})
}

// BenchmarkEngine_ComplexWorkflow benchmarks a realistic multi-stage workflow
func BenchmarkEngine_ComplexWorkflow(b *testing.B) {
	payload := map[string]interface{}{
		"nodes": []interface{}{
			// Stage 1: Input processing
			map[string]interface{}{"id": "input1", "type": "number", "data": map[string]interface{}{"value": 100.0}},
			map[string]interface{}{"id": "input2", "type": "number", "data": map[string]interface{}{"value": 50.0}},
			map[string]interface{}{"id": "add1", "type": "operation", "data": map[string]interface{}{"op": "add"}},

			// Stage 2: Arithmetic operations
			map[string]interface{}{"id": "num3", "type": "number", "data": map[string]interface{}{"value": 2.0}},
			map[string]interface{}{"id": "mul1", "type": "operation", "data": map[string]interface{}{"op": "multiply"}},

			// Stage 3: State management
			map[string]interface{}{"id": "var1", "type": "variable", "data": map[string]interface{}{"var_op": "set", "var_name": "result"}},
			map[string]interface{}{"id": "var2", "type": "variable", "data": map[string]interface{}{"var_op": "get", "var_name": "result"}},

			// Stage 4: Final output
			map[string]interface{}{"id": "viz", "type": "visualization", "data": map[string]interface{}{"mode": "text"}},
		},
		"edges": []interface{}{
			map[string]interface{}{"source": "input1", "target": "add1"},
			map[string]interface{}{"source": "input2", "target": "add1"},
			map[string]interface{}{"source": "add1", "target": "mul1"},
			map[string]interface{}{"source": "num3", "target": "mul1"},
			map[string]interface{}{"source": "mul1", "target": "var1"},
			map[string]interface{}{"source": "var1", "target": "var2"},
			map[string]interface{}{"source": "var2", "target": "viz"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		engine, _ := New(jsonData)
		_, err := engine.Execute()
		if err != nil {
			b.Fatalf("Execute failed: %v", err)
		}
	}
}

// BenchmarkEngine_LargeWorkflow benchmarks workflow with many nodes
func BenchmarkEngine_LargeWorkflow(b *testing.B) {
	sizes := []int{20, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d_nodes", size), func(b *testing.B) {
			nodes := make([]interface{}, size)
			edges := make([]interface{}, 0)

			// Create a workflow with alternating number and operation nodes
			for i := 0; i < size; i += 3 {
				// Number node
				nodes[i] = map[string]interface{}{
					"id":   fmt.Sprintf("num%d", i),
					"type": "number",
					"data": map[string]interface{}{"value": float64(i + 1)},
				}

				if i+1 < size {
					// Another number node
					nodes[i+1] = map[string]interface{}{
						"id":   fmt.Sprintf("num%d", i+1),
						"type": "number",
						"data": map[string]interface{}{"value": float64(i + 2)},
					}
				}

				if i+2 < size {
					// Operation node
					nodes[i+2] = map[string]interface{}{
						"id":   fmt.Sprintf("op%d", i),
						"type": "operation",
						"data": map[string]interface{}{"op": "add"},
					}

					edges = append(edges,
						map[string]interface{}{"source": fmt.Sprintf("num%d", i), "target": fmt.Sprintf("op%d", i)},
						map[string]interface{}{"source": fmt.Sprintf("num%d", i+1), "target": fmt.Sprintf("op%d", i)},
					)
				}
			}

			payload := map[string]interface{}{
				"nodes": nodes,
				"edges": edges,
			}
			jsonData, _ := json.Marshal(payload)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				engine, _ := New(jsonData)
				_, err := engine.Execute()
				if err != nil {
					b.Fatalf("Execute failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkEngine_EngineCreation benchmarks just the engine creation (parsing)
func BenchmarkEngine_EngineCreation(b *testing.B) {
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "type": "number", "data": map[string]interface{}{"value": 10.0}},
			map[string]interface{}{"id": "2", "type": "number", "data": map[string]interface{}{"value": 20.0}},
			map[string]interface{}{"id": "3", "type": "operation", "data": map[string]interface{}{"op": "add"}},
		},
		"edges": []interface{}{
			map[string]interface{}{"source": "1", "target": "3"},
			map[string]interface{}{"source": "2", "target": "3"},
		},
	}
	jsonData, _ := json.Marshal(payload)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := New(jsonData)
		if err != nil {
			b.Fatalf("New failed: %v", err)
		}
	}
}

// BenchmarkEngine_Execution benchmarks just the execution (after parsing)
func BenchmarkEngine_Execution(b *testing.B) {
	payload := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{"id": "1", "type": "number", "data": map[string]interface{}{"value": 10.0}},
			map[string]interface{}{"id": "2", "type": "number", "data": map[string]interface{}{"value": 20.0}},
			map[string]interface{}{"id": "3", "type": "operation", "data": map[string]interface{}{"op": "add"}},
		},
		"edges": []interface{}{
			map[string]interface{}{"source": "1", "target": "3"},
			map[string]interface{}{"source": "2", "target": "3"},
		},
	}
	jsonData, _ := json.Marshal(payload)
	engine, _ := New(jsonData)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := engine.Execute()
		if err != nil {
			b.Fatalf("Execute failed: %v", err)
		}
	}
}
