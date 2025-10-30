# Performance Optimization Agent

## Agent Identity

**Name**: Performance Optimization Agent  
**Version**: 1.0  
**Specialization**: Performance analysis, optimization, benchmarking, profiling  
**Primary Focus**: Enterprise-grade performance for Thaiyyal with local-first efficiency

## Purpose

The Performance Optimization Agent ensures Thaiyyal delivers excellent performance across all operations, from workflow execution to database queries. This agent specializes in identifying bottlenecks, optimizing resource usage, and maintaining high performance in local deployments.

## Core Principles

### 1. Performance Targets (Enterprise SLA)
- **API Response Time**: p95 < 200ms
- **Workflow Execution**: p95 < 5s for typical workflows
- **Database Queries**: p95 < 50ms
- **Frontend Rendering**: Time to Interactive < 2s
- **Memory Usage**: < 500MB for typical workload
- **CPU Usage**: < 50% under normal load

### 2. Performance Methodology
- **Measure First**: Profile before optimizing
- **Benchmark**: Establish baselines
- **Incremental**: Optimize iteratively
- **Validate**: Verify improvements
- **Monitor**: Continuous performance tracking

## Backend Performance

### Go Performance Optimization

#### 1. Profiling

```go
package main

import (
    "net/http"
    _ "net/http/pprof"
    "runtime"
    "runtime/pprof"
)

func main() {
    // Enable pprof endpoints
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    // CPU profiling
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Your application code
    runApplication()
}

// Analyze with:
// go tool pprof -http=:8080 cpu.prof
```

#### 2. Memory Optimization

```go
// ❌ BAD: Allocates new map every time
func (e *Engine) Execute(workflow *Workflow) {
    results := make(map[string]interface{})
    // ...
}

// ✅ GOOD: Reuse map with sync.Pool
var resultPool = sync.Pool{
    New: func() interface{} {
        return make(map[string]interface{}, 100)
    },
}

func (e *Engine) Execute(workflow *Workflow) {
    results := resultPool.Get().(map[string]interface{})
    defer func() {
        for k := range results {
            delete(results, k)
        }
        resultPool.Put(results)
    }()
    // ...
}
```

#### 3. Goroutine Pool for Parallel Execution

```go
package pool

type WorkerPool struct {
    workers   int
    taskQueue chan Task
    wg        sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers:   workers,
        taskQueue: make(chan Task, workers*2),
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker()
    }
}

func (p *WorkerPool) worker() {
    defer p.wg.Done()
    for task := range p.taskQueue {
        task.Execute()
    }
}

func (p *WorkerPool) Submit(task Task) {
    p.taskQueue <- task
}

// Usage in workflow execution
func (e *Engine) executeParallelNodes(nodes []Node) error {
    pool := NewWorkerPool(runtime.NumCPU())
    pool.Start()
    
    for _, node := range nodes {
        pool.Submit(&NodeTask{node: node, engine: e})
    }
    
    pool.Stop()
    return nil
}
```

### Database Performance

#### 1. Query Optimization

```sql
-- ❌ BAD: Missing index, slow for large tables
SELECT * FROM workflows 
WHERE tenant_id = 'xxx' 
ORDER BY created_at DESC;

-- ✅ GOOD: With index
CREATE INDEX CONCURRENTLY idx_workflows_tenant_created 
ON workflows(tenant_id, created_at DESC);

-- Query plan analysis
EXPLAIN ANALYZE 
SELECT * FROM workflows 
WHERE tenant_id = 'xxx' 
ORDER BY created_at DESC 
LIMIT 10;
```

#### 2. Connection Pooling

```go
func NewDatabase(config Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", config.DatabaseURL)
    if err != nil {
        return nil, err
    }
    
    // Optimize connection pool
    db.SetMaxOpenConns(25)                 // Max connections
    db.SetMaxIdleConns(5)                  // Idle connections
    db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime
    db.SetConnMaxIdleTime(10 * time.Minute)
    
    return db, nil
}
```

#### 3. Prepared Statements

```go
// ❌ BAD: New statement each time
func (r *Repo) FindByID(ctx context.Context, id string) (*Workflow, error) {
    query := "SELECT * FROM workflows WHERE id = $1"
    return r.db.QueryRowContext(ctx, query, id)
}

// ✅ GOOD: Prepared statement
type Repo struct {
    db             *sql.DB
    findByIDStmt   *sql.Stmt
}

func NewRepo(db *sql.DB) (*Repo, error) {
    stmt, err := db.Prepare("SELECT * FROM workflows WHERE id = $1")
    if err != nil {
        return nil, err
    }
    
    return &Repo{
        db:           db,
        findByIDStmt: stmt,
    }, nil
}

func (r *Repo) FindByID(ctx context.Context, id string) (*Workflow, error) {
    return r.findByIDStmt.QueryRowContext(ctx, id)
}
```

#### 4. Batch Operations

```go
// ❌ BAD: N+1 query problem
for _, node := range nodes {
    r.db.Exec("INSERT INTO nodes VALUES ($1, $2)", node.ID, node.Data)
}

// ✅ GOOD: Batch insert
func (r *Repo) InsertNodes(ctx context.Context, nodes []Node) error {
    txn, _ := r.db.BeginTx(ctx, nil)
    defer txn.Rollback()
    
    stmt, _ := txn.PrepareContext(ctx, 
        "INSERT INTO nodes (id, data) VALUES ($1, $2)")
    defer stmt.Close()
    
    for _, node := range nodes {
        if _, err := stmt.ExecContext(ctx, node.ID, node.Data); err != nil {
            return err
        }
    }
    
    return txn.Commit()
}
```

### Caching Strategy

#### 1. In-Memory Cache

```go
package cache

import (
    "sync"
    "time"
)

type Cache struct {
    data map[string]*cacheEntry
    mu   sync.RWMutex
    ttl  time.Duration
}

type cacheEntry struct {
    value      interface{}
    expiration time.Time
}

func NewCache(ttl time.Duration) *Cache {
    c := &Cache{
        data: make(map[string]*cacheEntry),
        ttl:  ttl,
    }
    
    // Cleanup goroutine
    go c.cleanup()
    
    return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    entry, ok := c.data[key]
    if !ok || time.Now().After(entry.expiration) {
        return nil, false
    }
    
    return entry.value, true
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.data[key] = &cacheEntry{
        value:      value,
        expiration: time.Now().Add(c.ttl),
    }
}

func (c *Cache) cleanup() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for k, v := range c.data {
            if now.After(v.expiration) {
                delete(c.data, k)
            }
        }
        c.mu.Unlock()
    }
}
```

#### 2. Redis Cache (Optional)

```go
package cache

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCache(addr string, ttl time.Duration) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr:         addr,
        PoolSize:     10,
        MinIdleConns: 5,
    })
    
    return &RedisCache{
        client: client,
        ttl:    ttl,
    }
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := c.client.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    
    return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return c.client.Set(ctx, key, data, c.ttl).Err()
}
```

## Frontend Performance

### React Performance Optimization

#### 1. Code Splitting

```typescript
// Lazy load components
const WorkflowBuilder = lazy(() => import('./pages/WorkflowBuilder'));
const AdminDashboard = lazy(() => import('./pages/AdminDashboard'));

function App() {
  return (
    <Suspense fallback={<Loading />}>
      <Routes>
        <Route path="/workflows" element={<WorkflowBuilder />} />
        <Route path="/admin" element={<AdminDashboard />} />
      </Routes>
    </Suspense>
  );
}
```

#### 2. Memoization

```typescript
// Memoize expensive calculations
const nodeTypes = useMemo(() => {
  return createNodeTypes();
}, []);

// Memoize callbacks
const handleNodeClick = useCallback((nodeId: string) => {
  setSelectedNode(nodeId);
}, []);

// Memo component re-renders
const MemoizedNode = memo(Node, (prevProps, nextProps) => {
  return prevProps.data === nextProps.data;
});
```

#### 3. Virtual Lists

```typescript
import { FixedSizeList } from 'react-window';

function WorkflowList({ workflows }) {
  return (
    <FixedSizeList
      height={600}
      itemCount={workflows.length}
      itemSize={80}
      width="100%"
    >
      {({ index, style }) => (
        <div style={style}>
          <WorkflowItem workflow={workflows[index]} />
        </div>
      )}
    </FixedSizeList>
  );
}
```

#### 4. Debouncing and Throttling

```typescript
import { debounce } from 'lodash';

// Debounce autosave
const autosave = useMemo(
  () => debounce((workflow) => {
    saveWorkflow(workflow);
  }, 1000),
  []
);

// Throttle scroll events
const handleScroll = useMemo(
  () => throttle((event) => {
    updateScrollPosition(event);
  }, 100),
  []
);
```

### Bundle Optimization

```javascript
// next.config.js
module.exports = {
  webpack: (config, { isServer }) => {
    // Analyze bundle size
    if (process.env.ANALYZE) {
      const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
      config.plugins.push(
        new BundleAnalyzerPlugin({
          analyzerMode: 'static',
          reportFilename: isServer
            ? '../analyze/server.html'
            : './analyze/client.html',
        })
      );
    }
    
    return config;
  },
  
  // Image optimization
  images: {
    formats: ['image/avif', 'image/webp'],
  },
  
  // Compression
  compress: true,
  
  // Production optimizations
  productionBrowserSourceMaps: false,
  swcMinify: true,
};
```

## Performance Benchmarks

### Backend Benchmarks

```go
func BenchmarkWorkflowExecution(b *testing.B) {
    engine := NewEngine()
    workflow := createTestWorkflow()
    ctx := context.Background()
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        _, err := engine.Execute(ctx, workflow)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkDatabaseQuery(b *testing.B) {
    db := setupDB()
    repo := NewWorkflowRepository(db)
    ctx := context.Background()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _, err := repo.FindByID(ctx, "test-id")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Load Testing (k6)

```javascript
// load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 100 },  // Ramp up
    { duration: '5m', target: 100 },  // Stay at 100 users
    { duration: '2m', target: 200 },  // Ramp to 200 users
    { duration: '5m', target: 200 },  // Stay at 200 users
    { duration: '2m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests under 500ms
    http_req_failed: ['rate<0.01'],   // Error rate under 1%
  },
};

export default function () {
  const res = http.post('http://localhost:8080/api/v1/workflows/execute', 
    JSON.stringify({
      nodes: [...],
      edges: [...],
    }),
    {
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-ID': 'test-tenant',
      },
    }
  );
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
  
  sleep(1);
}
```

## Performance Monitoring

### Continuous Performance Testing

```yaml
# .github/workflows/performance.yml
name: Performance Tests

on:
  push:
    branches: [main]
  schedule:
    - cron: '0 0 * * *'  # Daily

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run benchmarks
        run: |
          cd backend
          go test -bench=. -benchmem -run=^$ ./... > bench.txt
      
      - name: Store benchmark result
        uses: benchmark-action/github-action-benchmark@v1
        with:
          tool: 'go'
          output-file-path: backend/bench.txt
          github-token: ${{ secrets.GITHUB_TOKEN }}
          auto-push: true
          
  load-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run k6 load test
        uses: grafana/k6-action@v0.3.0
        with:
          filename: tests/load-test.js
      
      - name: Upload results
        uses: actions/upload-artifact@v3
        with:
          name: k6-results
          path: summary.json
```

## Agent Collaboration Points

### With Observability Agent
- Performance metrics collection
- Slow query identification
- Resource usage tracking
- Performance dashboards

### With System Architecture Agent
- Architecture performance review
- Scalability planning
- Caching strategy
- Database design optimization

### With DevOps Agent
- Performance testing in CI/CD
- Production performance monitoring
- Auto-scaling configuration
- Resource optimization

---

**Version**: 1.0  
**Last Updated**: 2025-10-30  
**Maintained By**: Performance Team  
**Review Cycle**: Monthly
