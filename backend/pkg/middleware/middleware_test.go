package middleware

import (
	"errors"
	"fmt"
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/executor"
	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// mockMiddleware records execution order for testing
type mockMiddleware struct {
	name       string
	order      *[]string
	shouldFail bool
}

func (m *mockMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	*m.order = append(*m.order, m.name+":pre")
	
	if m.shouldFail {
		return nil, errors.New(m.name + " failed")
	}
	
	result, err := next(ctx, node)
	
	*m.order = append(*m.order, m.name+":post")
	return result, err
}

func (m *mockMiddleware) Name() string {
	return m.name
}

// TestChain_SingleMiddleware tests chain with one middleware
func TestChain_SingleMiddleware(t *testing.T) {
	order := []string{}
	
	chain := NewChain()
	chain.Use(&mockMiddleware{name: "M1", order: &order})
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		order = append(order, "handler")
		return "result", nil
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	result, err := chain.Execute(nil, node, handler)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != "result" {
		t.Errorf("expected 'result', got %v", result)
	}
	
	expected := []string{"M1:pre", "handler", "M1:post"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d executions, got %d", len(expected), len(order))
	}
	
	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("execution %d: expected %s, got %s", i, exp, order[i])
		}
	}
}

// TestChain_MultipleMiddleware tests chain with multiple middleware
func TestChain_MultipleMiddleware(t *testing.T) {
	order := []string{}
	
	chain := NewChain()
	chain.Use(&mockMiddleware{name: "M1", order: &order})
	chain.Use(&mockMiddleware{name: "M2", order: &order})
	chain.Use(&mockMiddleware{name: "M3", order: &order})
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		order = append(order, "handler")
		return "result", nil
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	result, err := chain.Execute(nil, node, handler)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != "result" {
		t.Errorf("expected 'result', got %v", result)
	}
	
	// Middleware execute in order: M1(pre) -> M2(pre) -> M3(pre) -> handler -> M3(post) -> M2(post) -> M1(post)
	expected := []string{
		"M1:pre", "M2:pre", "M3:pre", "handler", "M3:post", "M2:post", "M1:post",
	}
	
	if len(order) != len(expected) {
		t.Fatalf("expected %d executions, got %d: %v", len(expected), len(order), order)
	}
	
	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("execution %d: expected %s, got %s", i, exp, order[i])
		}
	}
}

// TestChain_EmptyChain tests chain with no middleware
func TestChain_EmptyChain(t *testing.T) {
	order := []string{}
	
	chain := NewChain()
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		order = append(order, "handler")
		return "result", nil
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	result, err := chain.Execute(nil, node, handler)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != "result" {
		t.Errorf("expected 'result', got %v", result)
	}
	
	expected := []string{"handler"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d executions, got %d", len(expected), len(order))
	}
	
	if order[0] != expected[0] {
		t.Errorf("expected %s, got %s", expected[0], order[0])
	}
}

// TestChain_ErrorPropagation tests error propagation through the chain
func TestChain_ErrorPropagation(t *testing.T) {
	order := []string{}
	
	chain := NewChain()
	chain.Use(&mockMiddleware{name: "M1", order: &order})
	chain.Use(&mockMiddleware{name: "M2", order: &order, shouldFail: true})
	chain.Use(&mockMiddleware{name: "M3", order: &order})
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		order = append(order, "handler")
		return "result", nil
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	result, err := chain.Execute(nil, node, handler)
	
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	
	if err.Error() != "M2 failed" {
		t.Errorf("expected 'M2 failed', got %v", err)
	}
	
	if result != nil {
		t.Errorf("expected nil result on error, got %v", result)
	}
	
	// M2 should fail before calling M3 or handler, but M1:post should still execute
	expected := []string{"M1:pre", "M2:pre", "M1:post"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d executions, got %d: %v", len(expected), len(order), order)
	}
	
	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("execution %d: expected %s, got %s", i, exp, order[i])
		}
	}
}

// TestChain_HandlerError tests error from handler
func TestChain_HandlerError(t *testing.T) {
	order := []string{}
	
	chain := NewChain()
	chain.Use(&mockMiddleware{name: "M1", order: &order})
	chain.Use(&mockMiddleware{name: "M2", order: &order})
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		order = append(order, "handler")
		return nil, errors.New("handler failed")
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	_, err := chain.Execute(nil, node, handler)
	
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	
	if err.Error() != "handler failed" {
		t.Errorf("expected 'handler failed', got %v", err)
	}
	
	// Middleware should still execute post processing even on handler error
	expected := []string{"M1:pre", "M2:pre", "handler", "M2:post", "M1:post"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d executions, got %d: %v", len(expected), len(order), order)
	}
}

// TestChain_Len tests the Len method
func TestChain_Len(t *testing.T) {
	chain := NewChain()
	
	if chain.Len() != 0 {
		t.Errorf("expected length 0, got %d", chain.Len())
	}
	
	chain.Use(&mockMiddleware{name: "M1", order: &[]string{}})
	if chain.Len() != 1 {
		t.Errorf("expected length 1, got %d", chain.Len())
	}
	
	chain.Use(&mockMiddleware{name: "M2", order: &[]string{}})
	chain.Use(&mockMiddleware{name: "M3", order: &[]string{}})
	if chain.Len() != 3 {
		t.Errorf("expected length 3, got %d", chain.Len())
	}
}

// TestChain_Middlewares tests the Middlewares method
func TestChain_Middlewares(t *testing.T) {
	chain := NewChain()
	
	m1 := &mockMiddleware{name: "M1", order: &[]string{}}
	m2 := &mockMiddleware{name: "M2", order: &[]string{}}
	
	chain.Use(m1).Use(m2)
	
	middlewares := chain.Middlewares()
	if len(middlewares) != 2 {
		t.Fatalf("expected 2 middleware, got %d", len(middlewares))
	}
	
	if middlewares[0].Name() != "M1" {
		t.Errorf("expected M1, got %s", middlewares[0].Name())
	}
	
	if middlewares[1].Name() != "M2" {
		t.Errorf("expected M2, got %s", middlewares[1].Name())
	}
}

// shortCircuitMiddleware demonstrates middleware that short-circuits execution
type shortCircuitMiddleware struct {
	returnValue interface{}
}

func (m *shortCircuitMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	// Short-circuit: return cached value without calling next
	return m.returnValue, nil
}

func (m *shortCircuitMiddleware) Name() string {
	return "ShortCircuit"
}

// TestChain_ShortCircuit tests middleware that doesn't call next
func TestChain_ShortCircuit(t *testing.T) {
	order := []string{}
	
	chain := NewChain()
	chain.Use(&mockMiddleware{name: "M1", order: &order})
	chain.Use(&shortCircuitMiddleware{returnValue: "cached"})
	chain.Use(&mockMiddleware{name: "M3", order: &order})
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		order = append(order, "handler")
		return "fresh", nil
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	result, err := chain.Execute(nil, node, handler)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if result != "cached" {
		t.Errorf("expected 'cached', got %v", result)
	}
	
	// Only M1:pre should execute, then short-circuit returns
	expected := []string{"M1:pre", "M1:post"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d executions, got %d: %v", len(expected), len(order), order)
	}
}

// modifyingMiddleware modifies the result
type modifyingMiddleware struct {
	prefix string
}

func (m *modifyingMiddleware) Process(ctx executor.ExecutionContext, node types.Node, next Handler) (interface{}, error) {
	result, err := next(ctx, node)
	if err != nil {
		return result, err
	}
	
	// Modify result
	if str, ok := result.(string); ok {
		return m.prefix + str, nil
	}
	return result, nil
}

func (m *modifyingMiddleware) Name() string {
	return "Modifying"
}

// TestChain_ResultModification tests middleware that modifies results
func TestChain_ResultModification(t *testing.T) {
	chain := NewChain()
	chain.Use(&modifyingMiddleware{prefix: "A:"})
	chain.Use(&modifyingMiddleware{prefix: "B:"})
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "result", nil
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	result, err := chain.Execute(nil, node, handler)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Middleware execute in order, so post-processing is reverse:
	// A(pre) -> B(pre) -> handler("result") -> B(post, "result" -> "B:result") -> A(post, "B:result" -> "A:B:result")
	expected := "A:B:result"
	if result != expected {
		t.Errorf("expected %s, got %v", expected, result)
	}
}

// BenchmarkChain_NoMiddleware benchmarks execution without middleware
func BenchmarkChain_NoMiddleware(b *testing.B) {
	chain := NewChain()
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "result", nil
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = chain.Execute(nil, node, handler)
	}
}

// BenchmarkChain_SingleMiddleware benchmarks with one middleware
func BenchmarkChain_SingleMiddleware(b *testing.B) {
	order := []string{}
	chain := NewChain()
	chain.Use(&mockMiddleware{name: "M1", order: &order})
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "result", nil
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = chain.Execute(nil, node, handler)
	}
}

// BenchmarkChain_FiveMiddleware benchmarks with five middleware
func BenchmarkChain_FiveMiddleware(b *testing.B) {
	order := []string{}
	chain := NewChain()
	for i := 0; i < 5; i++ {
		chain.Use(&mockMiddleware{name: fmt.Sprintf("M%d", i), order: &order})
	}
	
	handler := func(ctx executor.ExecutionContext, node types.Node) (interface{}, error) {
		return "result", nil
	}
	
	node := types.Node{ID: "test", Type: types.NodeTypeNumber}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = chain.Execute(nil, node, handler)
	}
}
