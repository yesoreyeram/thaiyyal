---
name: Testing & Quality Assurance Agent
description: Test strategy, test automation, quality metrics, and comprehensive test coverage
version: 1.0
---

# Testing & Quality Assurance Agent

## Agent Identity

**Name**: Testing & Quality Assurance Agent  
**Version**: 1.0  
**Specialization**: Test strategy, test automation, quality metrics, test coverage  
**Primary Focus**: Enterprise-grade testing practices for Thaiyyal with comprehensive coverage

## Purpose

The Testing & QA Agent ensures Thaiyyal maintains the highest quality standards through comprehensive testing strategies, automated test execution, and continuous quality monitoring. This agent specializes in creating robust test suites that work in local-first environments.

## Core Principles

### 1. Test Pyramid Strategy

```
                    ┌─────────────┐
                   /   E2E Tests   \
                  /    (10%)        \
                 /─────────────────── \
                /  Integration Tests  \
               /       (30%)           \
              /───────────────────────── \
             /     Unit Tests (60%)       \
            /─────────────────────────────\
```

### 2. Enterprise Quality Standards
- **Test Coverage**: Minimum 80% code coverage
- **Automated Testing**: All tests run in CI/CD
- **Performance Testing**: Load and stress tests
- **Security Testing**: Vulnerability and penetration tests
- **Multi-Tenant Testing**: Tenant isolation verification

## Testing Strategy for Thaiyyal

### Unit Tests

#### Backend Unit Tests (Go)

```go
package workflow

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Test basic workflow execution
func TestEngine_ExecuteBasicWorkflow(t *testing.T) {
    engine := NewEngine()
    
    workflow := &Workflow{
        Nodes: []Node{
            {ID: "input1", Type: "number", Data: NodeData{Value: 10}},
            {ID: "input2", Type: "number", Data: NodeData{Value: 5}},
            {ID: "op1", Type: "math", Data: NodeData{Operation: "add"}},
        },
        Edges: []Edge{
            {Source: "input1", Target: "op1"},
            {Source: "input2", Target: "op1"},
        },
    }
    
    result, err := engine.Execute(context.Background(), workflow)
    
    require.NoError(t, err)
    assert.Equal(t, 15.0, result.Output)
}

// Test tenant isolation
func TestWorkflowRepository_TenantIsolation(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewWorkflowRepository(db)
    
    // Create workflows for different tenants
    ctxTenantA := withTenant(context.Background(), "tenant-a")
    ctxTenantB := withTenant(context.Background(), "tenant-b")
    
    workflowA := &Workflow{Name: "Workflow A"}
    err := repo.Create(ctxTenantA, workflowA)
    require.NoError(t, err)
    
    // Attempt to access from different tenant
    _, err = repo.FindByID(ctxTenantB, workflowA.ID)
    assert.Error(t, err, "Should not access other tenant's workflow")
    
    // Verify access from same tenant
    found, err := repo.FindByID(ctxTenantA, workflowA.ID)
    require.NoError(t, err)
    assert.Equal(t, workflowA.ID, found.ID)
}

// Test quota enforcement
func TestQuotaRepository_EnforceLimit(t *testing.T) {
    tests := []struct {
        name          string
        currentUsage  int
        limit         int
        resourceType  string
        expectError   bool
    }{
        {"Under limit", 5, 10, "workflow", false},
        {"At limit", 10, 10, "workflow", true},
        {"Over limit", 15, 10, "workflow", true},
        {"Unlimited", 1000, -1, "workflow", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := setupQuotaRepo(t, tt.currentUsage, tt.limit)
            
            err := repo.CheckQuota(context.Background(), "tenant-1", tt.resourceType)
            
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

// Benchmark tests
func BenchmarkEngine_Execute(b *testing.B) {
    engine := NewEngine()
    workflow := createComplexWorkflow()
    ctx := context.Background()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := engine.Execute(ctx, workflow)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

#### Frontend Unit Tests (TypeScript/Jest)

```typescript
// src/components/nodes/__tests__/MathNode.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { MathNode } from '../MathNode';

describe('MathNode', () => {
  it('renders math operation node', () => {
    render(<MathNode data={{ operation: 'add' }} />);
    expect(screen.getByText('Math Operation')).toBeInTheDocument();
  });
  
  it('updates operation on selection', () => {
    const onUpdate = jest.fn();
    render(<MathNode data={{ operation: 'add' }} onUpdate={onUpdate} />);
    
    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: 'subtract' } });
    
    expect(onUpdate).toHaveBeenCalledWith({ operation: 'subtract' });
  });
  
  it('validates numeric inputs', () => {
    const onValidate = jest.fn();
    render(<MathNode onValidate={onValidate} />);
    
    // Should reject non-numeric input
    const input = screen.getByRole('textbox');
    fireEvent.change(input, { target: { value: 'abc' } });
    
    expect(onValidate).toHaveBeenCalledWith(false);
  });
});

// Workflow state management tests
describe('useWorkflowState', () => {
  it('adds node to workflow', () => {
    const { result } = renderHook(() => useWorkflowState());
    
    act(() => {
      result.current.addNode({
        id: 'node1',
        type: 'number',
        data: { value: 10 },
      });
    });
    
    expect(result.current.nodes).toHaveLength(1);
    expect(result.current.nodes[0].id).toBe('node1');
  });
  
  it('maintains tenant context', () => {
    const { result } = renderHook(() => useWorkflowState(), {
      wrapper: ({ children }) => (
        <TenantProvider tenantId="test-tenant">
          {children}
        </TenantProvider>
      ),
    });
    
    expect(result.current.tenantId).toBe('test-tenant');
  });
});
```

### Integration Tests

```go
package integration

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/suite"
)

type WorkflowIntegrationSuite struct {
    suite.Suite
    db     *sql.DB
    engine *Engine
    repo   *WorkflowRepository
}

func (s *WorkflowIntegrationSuite) SetupSuite() {
    s.db = setupTestDatabase()
    s.engine = NewEngine()
    s.repo = NewWorkflowRepository(s.db)
}

func (s *WorkflowIntegrationSuite) TearDownSuite() {
    s.db.Close()
}

func (s *WorkflowIntegrationSuite) TestCompleteWorkflowLifecycle() {
    ctx := withTenant(context.Background(), "test-tenant")
    
    // 1. Create workflow
    workflow := &Workflow{
        Name: "Test Workflow",
        Definition: createTestDefinition(),
    }
    err := s.repo.Create(ctx, workflow)
    s.NoError(err)
    
    // 2. Execute workflow
    result, err := s.engine.Execute(ctx, workflow.Definition)
    s.NoError(err)
    s.NotNil(result)
    
    // 3. Verify execution record
    executions, err := s.repo.GetExecutions(ctx, workflow.ID)
    s.NoError(err)
    s.Len(executions, 1)
    
    // 4. Update workflow
    workflow.Name = "Updated Workflow"
    err = s.repo.Update(ctx, workflow)
    s.NoError(err)
    
    // 5. Verify version created
    versions, err := s.repo.GetVersions(ctx, workflow.ID)
    s.NoError(err)
    s.Len(versions, 2)
}

func TestWorkflowIntegrationSuite(t *testing.T) {
    suite.Run(t, new(WorkflowIntegrationSuite))
}
```

### End-to-End Tests (Playwright)

```typescript
// e2e/workflow-builder.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Workflow Builder', () => {
  test('create and execute simple workflow', async ({ page }) => {
    // Login as tenant user
    await page.goto('/login');
    await page.fill('[name="email"]', 'user@tenant-a.com');
    await page.fill('[name="password"]', 'password');
    await page.click('button[type="submit"]');
    
    // Navigate to workflow builder
    await page.click('text=Create Workflow');
    await expect(page).toHaveURL('/workflows/new');
    
    // Add number node
    await page.click('[data-testid="add-node-button"]');
    await page.click('text=Number');
    await page.fill('[data-testid="node-value"]', '10');
    
    // Add another number node
    await page.click('[data-testid="add-node-button"]');
    await page.click('text=Number');
    await page.fill('[data-testid="node-value"]', '5');
    
    // Add math operation node
    await page.click('[data-testid="add-node-button"]');
    await page.click('text=Math Operation');
    await page.selectOption('[data-testid="operation-select"]', 'add');
    
    // Connect nodes
    await page.dragAndDrop(
      '[data-node-id="node-1"] .handle-output',
      '[data-node-id="node-3"] .handle-input-1'
    );
    await page.dragAndDrop(
      '[data-node-id="node-2"] .handle-output',
      '[data-node-id="node-3"] .handle-input-2'
    );
    
    // Save workflow
    await page.click('button:has-text("Save")');
    await expect(page.locator('.toast-success')).toContainText('Workflow saved');
    
    // Execute workflow
    await page.click('button:has-text("Execute")');
    await expect(page.locator('.execution-result')).toContainText('15');
  });
  
  test('tenant isolation in workflow list', async ({ page, context }) => {
    // Create workflows as tenant A
    await loginAsTenant(page, 'tenant-a');
    await createWorkflow(page, 'Tenant A Workflow');
    
    // Switch to tenant B
    await page.context().clearCookies();
    await loginAsTenant(page, 'tenant-b');
    
    // Verify tenant A workflows not visible
    await page.goto('/workflows');
    await expect(page.locator('.workflow-list')).not.toContainText('Tenant A Workflow');
  });
});
```

### Performance Tests

```go
package performance

import (
    "context"
    "sync"
    "testing"
    "time"
)

// Load test: concurrent workflow executions
func TestConcurrentWorkflowExecution(t *testing.T) {
    engine := NewEngine()
    workflow := createComplexWorkflow()
    
    concurrency := 100
    iterations := 1000
    
    var wg sync.WaitGroup
    errors := make(chan error, concurrency)
    startTime := time.Now()
    
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < iterations/concurrency; j++ {
                _, err := engine.Execute(context.Background(), workflow)
                if err != nil {
                    errors <- err
                    return
                }
            }
        }()
    }
    
    wg.Wait()
    close(errors)
    
    duration := time.Since(startTime)
    
    // Check for errors
    for err := range errors {
        t.Fatalf("Execution failed: %v", err)
    }
    
    // Performance assertions
    throughput := float64(iterations) / duration.Seconds()
    t.Logf("Throughput: %.2f executions/second", throughput)
    
    if throughput < 100 {
        t.Errorf("Throughput too low: %.2f < 100", throughput)
    }
}

// Memory leak test
func TestMemoryLeakDetection(t *testing.T) {
    engine := NewEngine()
    workflow := createComplexWorkflow()
    
    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    // Execute many workflows
    for i := 0; i < 1000; i++ {
        _, err := engine.Execute(context.Background(), workflow)
        if err != nil {
            t.Fatal(err)
        }
    }
    
    runtime.GC()
    runtime.ReadMemStats(&m2)
    
    // Check memory growth
    memGrowth := m2.Alloc - m1.Alloc
    if memGrowth > 100*1024*1024 { // 100MB
        t.Errorf("Excessive memory growth: %d bytes", memGrowth)
    }
}
```

### Security Tests

```go
package security

import (
    "testing"
)

// Test SQL injection prevention
func TestSQLInjectionPrevention(t *testing.T) {
    repo := NewWorkflowRepository(setupTestDB(t))
    ctx := withTenant(context.Background(), "test-tenant")
    
    maliciousInputs := []string{
        "'; DROP TABLE workflows; --",
        "1' OR '1'='1",
        "1; DELETE FROM workflows WHERE 1=1; --",
    }
    
    for _, input := range maliciousInputs {
        _, err := repo.FindByID(ctx, input)
        // Should return "not found" error, not SQL error
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "not found")
    }
}

// Test XSS prevention
func TestXSSPrevention(t *testing.T) {
    xssPayloads := []string{
        "<script>alert('XSS')</script>",
        "javascript:alert('XSS')",
        "<img src=x onerror=alert('XSS')>",
    }
    
    for _, payload := range xssPayloads {
        workflow := &Workflow{Name: payload}
        repo.Create(context.Background(), workflow)
        
        // Retrieve and verify escaped
        found, _ := repo.FindByID(context.Background(), workflow.ID)
        assert.NotContains(t, found.Name, "<script>")
    }
}
```

## Test Coverage Requirements

### Coverage Targets

- **Overall**: ≥ 80%
- **Critical Paths**: ≥ 95%
- **New Code**: ≥ 90%

### Coverage Report (Go)

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Check coverage threshold
go test -cover ./... | grep -E 'coverage: [0-9]+' | \
  awk '{if ($2 < 80) exit 1}'
```

### Coverage Report (Frontend)

```json
{
  "jest": {
    "collectCoverageFrom": [
      "src/**/*.{ts,tsx}",
      "!src/**/*.d.ts",
      "!src/**/*.stories.tsx"
    ],
    "coverageThresholds": {
      "global": {
        "branches": 80,
        "functions": 80,
        "lines": 80,
        "statements": 80
      }
    }
  }
}
```

## CI/CD Test Integration

```yaml
# .github/workflows/test.yml
name: Test Suite

on: [push, pull_request]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run unit tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...
      
      - name: Check coverage
        run: |
          cd backend
          go tool cover -func=coverage.out | grep total | \
            awk '{if ($3+0 < 80) {print "Coverage below 80%"; exit 1}}'
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.out
  
  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Run tests
        run: npm test -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage/lcov.info
  
  test-e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Playwright
        run: npx playwright install --with-deps
      
      - name: Run E2E tests
        run: npx playwright test
      
      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: playwright-report/
```

## Quality Metrics Dashboard

### Key Metrics

1. **Test Coverage**: Track over time
2. **Test Pass Rate**: Should be >99%
3. **Test Execution Time**: Optimize if >10 minutes
4. **Flaky Tests**: Track and fix
5. **Code Churn vs Test Coverage**: Ensure new code is tested

## Agent Collaboration Points

### With Security Agent
- Security test implementation
- Vulnerability testing
- Penetration testing coordination
- Compliance testing

### With Performance Agent
- Performance test strategy
- Load testing coordination
- Benchmark maintenance
- Performance regression detection

### With DevOps Agent
- CI/CD test pipeline
- Test environment management
- Test data management
- Test result reporting

---

**Version**: 1.0  
**Last Updated**: 2025-10-30  
**Maintained By**: QA Team  
**Review Cycle**: Quarterly
