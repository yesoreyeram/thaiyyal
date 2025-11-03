package executor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// BenchmarkHTTPExecutor_Sequential benchmarks sequential HTTP requests with pooling
func BenchmarkHTTPExecutor_Sequential(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	executor := NewHTTPExecutor()
	config := types.Config{
		HTTPTimeout:      30 * time.Second,
		MaxHTTPRedirects: 10,
		MaxResponseSize:  10 * 1024 * 1024,
	}
	ctx := &mockExecutionContext{config: config}

	url := server.URL
	node := types.Node{
		ID:   "1",
		Type: types.NodeTypeHTTP,
		Data: types.NodeData{URL: &url},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := executor.Execute(ctx, node)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
	}
}

// BenchmarkHTTPExecutor_NoPooling benchmarks without connection pooling (creates new client each time)
func BenchmarkHTTPExecutor_NoPooling(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	config := types.Config{
		HTTPTimeout:      30 * time.Second,
		MaxHTTPRedirects: 10,
		MaxResponseSize:  10 * 1024 * 1024,
	}

	url := server.URL

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Create new client each time (simulating old behavior)
		client := &http.Client{
			Timeout: config.HTTPTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				IdleConnTimeout:     30 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
				DisableKeepAlives:   false,
			},
		}

		resp, err := client.Get(url)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkHTTPExecutor_Concurrent benchmarks concurrent requests
func BenchmarkHTTPExecutor_Concurrent(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	executor := NewHTTPExecutor()
	config := types.Config{
		HTTPTimeout:      30 * time.Second,
		MaxHTTPRedirects: 10,
		MaxResponseSize:  10 * 1024 * 1024,
	}
	ctx := &mockExecutionContext{config: config}

	url := server.URL
	node := types.Node{
		ID:   "1",
		Type: types.NodeTypeHTTP,
		Data: types.NodeData{URL: &url},
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := executor.Execute(ctx, node)
			if err != nil {
				b.Fatalf("Request failed: %v", err)
			}
		}
	})
}

// BenchmarkHTTPExecutor_MultipleHosts benchmarks requests to different hosts
func BenchmarkHTTPExecutor_MultipleHosts(b *testing.B) {
	// Create 3 test servers
	servers := make([]*httptest.Server, 3)
	urls := make([]string, 3)
	for i := 0; i < 3; i++ {
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))
		urls[i] = servers[i].URL
		defer servers[i].Close()
	}

	executor := NewHTTPExecutor()
	config := types.Config{
		HTTPTimeout:      30 * time.Second,
		MaxHTTPRedirects: 10,
		MaxResponseSize:  10 * 1024 * 1024,
	}
	ctx := &mockExecutionContext{config: config}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Round-robin across servers
		url := urls[i%3]
		node := types.Node{
			ID:   "1",
			Type: types.NodeTypeHTTP,
			Data: types.NodeData{URL: &url},
		}

		_, err := executor.Execute(ctx, node)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
	}
}
