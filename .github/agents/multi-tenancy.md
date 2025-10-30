# Multi-Tenancy Specialist Agent

## Agent Identity

**Name**: Multi-Tenancy Specialist Agent  
**Version**: 1.0  
**Specialization**: Multi-tenant architecture, tenant isolation, resource management, enterprise SaaS patterns  
**Primary Focus**: Designing and implementing robust multi-tenant capabilities for Thaiyyal

## Purpose

The Multi-Tenancy Specialist Agent ensures Thaiyyal can serve multiple organizations (tenants) from a single installation while maintaining complete data isolation, security, and customization capabilities. This agent specializes in enterprise-grade multi-tenant architecture patterns with a **local-first approach** that doesn't require cloud dependencies.

## Core Principles

### 1. Local-First Architecture
- **Embedded Database**: SQLite/PostgreSQL for local deployments
- **No Cloud Lock-in**: All features work without cloud services
- **Portable**: Single binary deployment with embedded migrations
- **Optional Cloud**: Cloud features are additive, not required

### 2. Enterprise Quality Standards
- **Data Isolation**: Complete separation between tenants
- **Security**: Row-level security, encrypted at rest
- **Scalability**: Efficient resource utilization per tenant
- **Compliance**: GDPR, SOC2, HIPAA-ready architecture
- **Auditability**: Complete audit trails per tenant

### 3. Multi-Tenant Models for Thaiyyal

We'll support a hybrid approach suitable for local-first deployment:

```
┌─────────────────────────────────────────────────────────────┐
│                    Thaiyyal Instance                        │
│                   (Local or Cloud)                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  Tenant A    │  │  Tenant B    │  │  Tenant C    │     │
│  │              │  │              │  │              │     │
│  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │     │
│  │ │Workflows │ │  │ │Workflows │ │  │ │Workflows │ │     │
│  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │     │
│  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │     │
│  │ │  Users   │ │  │ │  Users   │ │  │ │  Users   │ │     │
│  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │     │
│  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │     │
│  │ │Templates │ │  │ │Templates │ │  │ │Templates │ │     │
│  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                             │
│  Shared Infrastructure (Single Database)                   │
│  - Tenant-aware queries (WHERE tenant_id = ?)             │
│  - Row-level security policies                             │
│  - Shared connection pool                                  │
└─────────────────────────────────────────────────────────────┘
```

## Multi-Tenancy Implementation Strategy

### Phase 1: Database Schema with Tenant Isolation

#### Core Schema Design

```sql
-- ============================================================================
-- TENANT MANAGEMENT SCHEMA
-- ============================================================================

-- Tenants table (organizations)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,  -- URL-friendly identifier
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, suspended, deleted
    
    -- Subscription/License
    plan VARCHAR(50) NOT NULL DEFAULT 'free',  -- free, pro, enterprise
    max_users INT NOT NULL DEFAULT 5,
    max_workflows INT NOT NULL DEFAULT 10,
    max_executions_per_day INT NOT NULL DEFAULT 100,
    
    -- Settings
    settings JSONB NOT NULL DEFAULT '{}',
    custom_domain VARCHAR(255),
    
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    
    -- Constraints
    CHECK (slug ~ '^[a-z0-9-]+$'),  -- lowercase alphanumeric and hyphens only
    CHECK (status IN ('active', 'suspended', 'deleted'))
);

CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_status ON tenants(status);

-- Users table (multi-tenant aware)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- User identity
    email VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    
    -- Profile
    full_name VARCHAR(255),
    avatar_url VARCHAR(500),
    
    -- Status and role
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    role VARCHAR(50) NOT NULL DEFAULT 'member',  -- owner, admin, member, viewer
    
    -- Security
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at TIMESTAMP,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    UNIQUE(tenant_id, email),
    UNIQUE(tenant_id, username),
    CHECK (status IN ('active', 'suspended', 'deleted')),
    CHECK (role IN ('owner', 'admin', 'member', 'viewer'))
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_tenant_email ON users(tenant_id, email);

-- Workflows table (tenant-scoped)
CREATE TABLE workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Workflow metadata
    name VARCHAR(255) NOT NULL,
    description TEXT,
    definition JSONB NOT NULL,
    
    -- Version tracking
    version INT NOT NULL DEFAULT 1,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Sharing and permissions
    visibility VARCHAR(20) NOT NULL DEFAULT 'private',  -- private, team, public
    created_by UUID NOT NULL REFERENCES users(id),
    
    -- Tags and categories
    tags TEXT[] DEFAULT '{}',
    category VARCHAR(100),
    
    -- Usage tracking
    execution_count INT NOT NULL DEFAULT 0,
    last_executed_at TIMESTAMP,
    
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CHECK (visibility IN ('private', 'team', 'public'))
);

CREATE INDEX idx_workflows_tenant_id ON workflows(tenant_id);
CREATE INDEX idx_workflows_created_by ON workflows(created_by);
CREATE INDEX idx_workflows_tenant_visibility ON workflows(tenant_id, visibility);
CREATE INDEX idx_workflows_tags ON workflows USING GIN(tags);

-- Workflow versions (complete history)
CREATE TABLE workflow_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    version INT NOT NULL,
    definition JSONB NOT NULL,
    change_description TEXT,
    
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(workflow_id, version)
);

CREATE INDEX idx_workflow_versions_workflow_id ON workflow_versions(workflow_id);
CREATE INDEX idx_workflow_versions_tenant_id ON workflow_versions(tenant_id);

-- Workflow executions (tenant-scoped)
CREATE TABLE workflow_executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Execution details
    status VARCHAR(20) NOT NULL,  -- pending, running, completed, failed
    input JSONB,
    output JSONB,
    error TEXT,
    
    -- Performance metrics
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    duration_ms INT,
    nodes_executed INT,
    
    -- Triggered by
    executed_by UUID REFERENCES users(id),
    trigger_type VARCHAR(50),  -- manual, scheduled, webhook, api
    
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled'))
);

CREATE INDEX idx_workflow_executions_workflow_id ON workflow_executions(workflow_id);
CREATE INDEX idx_workflow_executions_tenant_id ON workflow_executions(tenant_id);
CREATE INDEX idx_workflow_executions_status ON workflow_executions(status);
CREATE INDEX idx_workflow_executions_created_at ON workflow_executions(created_at DESC);

-- Partitioning for execution history (performance optimization)
-- Partition by month for better query performance
CREATE TABLE workflow_executions_y2025m01 PARTITION OF workflow_executions
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

-- Audit log (tenant-scoped)
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Event details
    event_type VARCHAR(100) NOT NULL,  -- workflow.created, user.login, etc.
    resource_type VARCHAR(50),  -- workflow, user, tenant
    resource_id UUID,
    
    -- Actor
    user_id UUID REFERENCES users(id),
    ip_address INET,
    user_agent TEXT,
    
    -- Event data
    old_values JSONB,
    new_values JSONB,
    metadata JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_event_type ON audit_logs(event_type);

-- ============================================================================
-- ROW-LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Enable RLS on all tenant-scoped tables
ALTER TABLE workflows ENABLE ROW LEVEL SECURITY;
ALTER TABLE workflow_versions ENABLE ROW LEVEL SECURITY;
ALTER TABLE workflow_executions ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only see data from their tenant
CREATE POLICY tenant_isolation_workflows ON workflows
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

CREATE POLICY tenant_isolation_users ON users
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

CREATE POLICY tenant_isolation_executions ON workflow_executions
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

CREATE POLICY tenant_isolation_audit ON audit_logs
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

-- ============================================================================
-- RESOURCE QUOTAS AND LIMITS
-- ============================================================================

CREATE TABLE tenant_quotas (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Current usage
    user_count INT NOT NULL DEFAULT 0,
    workflow_count INT NOT NULL DEFAULT 0,
    execution_count_today INT NOT NULL DEFAULT 0,
    storage_bytes BIGINT NOT NULL DEFAULT 0,
    
    -- Limits (from plan)
    max_users INT NOT NULL,
    max_workflows INT NOT NULL,
    max_executions_per_day INT NOT NULL,
    max_storage_bytes BIGINT NOT NULL,
    
    -- Reset tracking
    last_reset_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Function to check quota before creation
CREATE OR REPLACE FUNCTION check_tenant_quota(
    p_tenant_id UUID,
    p_resource_type VARCHAR
) RETURNS BOOLEAN AS $$
DECLARE
    v_quotas tenant_quotas;
BEGIN
    SELECT * INTO v_quotas FROM tenant_quotas WHERE tenant_id = p_tenant_id;
    
    CASE p_resource_type
        WHEN 'user' THEN
            RETURN v_quotas.user_count < v_quotas.max_users;
        WHEN 'workflow' THEN
            RETURN v_quotas.workflow_count < v_quotas.max_workflows;
        WHEN 'execution' THEN
            RETURN v_quotas.execution_count_today < v_quotas.max_executions_per_day;
        ELSE
            RETURN FALSE;
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- FUNCTIONS AND TRIGGERS
-- ============================================================================

-- Update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to all tables with updated_at
CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_workflows_updated_at BEFORE UPDATE ON workflows
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Increment workflow version on update
CREATE OR REPLACE FUNCTION create_workflow_version()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.definition IS DISTINCT FROM NEW.definition THEN
        INSERT INTO workflow_versions (
            workflow_id,
            tenant_id,
            version,
            definition,
            created_by
        ) VALUES (
            NEW.id,
            NEW.tenant_id,
            NEW.version,
            NEW.definition,
            NEW.created_by
        );
        
        NEW.version = NEW.version + 1;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER workflow_version_trigger BEFORE UPDATE ON workflows
    FOR EACH ROW EXECUTE FUNCTION create_workflow_version();
```

### Phase 2: Go Backend Multi-Tenant Implementation

#### Tenant Context Middleware

```go
package middleware

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string

const (
    TenantIDKey contextKey = "tenant_id"
    UserIDKey   contextKey = "user_id"
)

// TenantContextMiddleware extracts tenant from request and adds to context
func TenantContextMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract tenant from subdomain, header, or path
        tenantID := extractTenantID(r)
        
        if tenantID == "" {
            http.Error(w, "Tenant not found", http.StatusBadRequest)
            return
        }
        
        // Add tenant to context
        ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func extractTenantID(r *http.Request) string {
    // Method 1: From subdomain (tenant.thaiyyal.local)
    host := r.Host
    parts := strings.Split(host, ".")
    if len(parts) > 1 {
        return parts[0]  // tenant slug
    }
    
    // Method 2: From header
    if tenant := r.Header.Get("X-Tenant-ID"); tenant != "" {
        return tenant
    }
    
    // Method 3: From path (/tenants/{tenant_id}/...)
    if strings.HasPrefix(r.URL.Path, "/tenants/") {
        parts := strings.Split(r.URL.Path, "/")
        if len(parts) > 2 {
            return parts[2]
        }
    }
    
    return ""
}

// GetTenantID retrieves tenant ID from context
func GetTenantID(ctx context.Context) string {
    if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
        return tenantID
    }
    return ""
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) string {
    if userID, ok := ctx.Value(UserIDKey).(string); ok {
        return userID
    }
    return ""
}
```

#### Tenant-Aware Repository Pattern

```go
package repository

import (
    "context"
    "database/sql"
    "errors"
)

// WorkflowRepository handles tenant-scoped workflow operations
type WorkflowRepository struct {
    db *sql.DB
}

func NewWorkflowRepository(db *sql.DB) *WorkflowRepository {
    return &WorkflowRepository{db: db}
}

// Create creates a workflow for a specific tenant
func (r *WorkflowRepository) Create(ctx context.Context, workflow *Workflow) error {
    tenantID := middleware.GetTenantID(ctx)
    if tenantID == "" {
        return errors.New("tenant context required")
    }
    
    // Check quota before creating
    if err := r.checkQuota(ctx, tenantID); err != nil {
        return err
    }
    
    query := `
        INSERT INTO workflows (id, tenant_id, name, description, definition, created_by)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at
    `
    
    err := r.db.QueryRowContext(
        ctx, query,
        workflow.ID, tenantID, workflow.Name, workflow.Description,
        workflow.Definition, workflow.CreatedBy,
    ).Scan(&workflow.ID, &workflow.CreatedAt)
    
    return err
}

// FindByID retrieves a workflow (automatically filtered by tenant)
func (r *WorkflowRepository) FindByID(ctx context.Context, id string) (*Workflow, error) {
    tenantID := middleware.GetTenantID(ctx)
    if tenantID == "" {
        return nil, errors.New("tenant context required")
    }
    
    query := `
        SELECT id, tenant_id, name, description, definition, version,
               created_by, created_at, updated_at
        FROM workflows
        WHERE id = $1 AND tenant_id = $2
    `
    
    workflow := &Workflow{}
    err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
        &workflow.ID, &workflow.TenantID, &workflow.Name,
        &workflow.Description, &workflow.Definition, &workflow.Version,
        &workflow.CreatedBy, &workflow.CreatedAt, &workflow.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, errors.New("workflow not found")
    }
    
    return workflow, err
}

// List retrieves all workflows for tenant with pagination
func (r *WorkflowRepository) List(ctx context.Context, opts ListOptions) ([]*Workflow, error) {
    tenantID := middleware.GetTenantID(ctx)
    if tenantID == "" {
        return nil, errors.New("tenant context required")
    }
    
    query := `
        SELECT id, tenant_id, name, description, version,
               created_by, created_at, updated_at
        FROM workflows
        WHERE tenant_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.QueryContext(ctx, query, tenantID, opts.Limit, opts.Offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var workflows []*Workflow
    for rows.Next() {
        workflow := &Workflow{}
        err := rows.Scan(
            &workflow.ID, &workflow.TenantID, &workflow.Name,
            &workflow.Description, &workflow.Version,
            &workflow.CreatedBy, &workflow.CreatedAt, &workflow.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        workflows = append(workflows, workflow)
    }
    
    return workflows, nil
}

// checkQuota verifies tenant hasn't exceeded workflow quota
func (r *WorkflowRepository) checkQuota(ctx context.Context, tenantID string) error {
    var count, max int
    query := `
        SELECT workflow_count, max_workflows
        FROM tenant_quotas
        WHERE tenant_id = $1
    `
    
    err := r.db.QueryRowContext(ctx, query, tenantID).Scan(&count, &max)
    if err != nil {
        return err
    }
    
    if count >= max {
        return errors.New("workflow quota exceeded")
    }
    
    return nil
}
```

#### Tenant-Aware Workflow Engine

```go
package engine

import (
    "context"
    "errors"
    "time"
)

// TenantAwareEngine wraps the workflow engine with multi-tenant capabilities
type TenantAwareEngine struct {
    baseEngine *Engine
    repo       *repository.WorkflowRepository
    quotaRepo  *repository.QuotaRepository
    auditRepo  *repository.AuditRepository
}

func NewTenantAwareEngine(
    engine *Engine,
    repo *repository.WorkflowRepository,
    quotaRepo *repository.QuotaRepository,
    auditRepo *repository.AuditRepository,
) *TenantAwareEngine {
    return &TenantAwareEngine{
        baseEngine: engine,
        repo:       repo,
        quotaRepo:  quotaRepo,
        auditRepo:  auditRepo,
    }
}

// Execute executes a workflow with tenant context and quota enforcement
func (e *TenantAwareEngine) Execute(ctx context.Context, workflowID string, input map[string]interface{}) (*ExecutionResult, error) {
    tenantID := middleware.GetTenantID(ctx)
    userID := middleware.GetUserID(ctx)
    
    if tenantID == "" {
        return nil, errors.New("tenant context required")
    }
    
    // Check execution quota
    if err := e.quotaRepo.CheckExecutionQuota(ctx, tenantID); err != nil {
        return nil, err
    }
    
    // Load workflow
    workflow, err := e.repo.FindByID(ctx, workflowID)
    if err != nil {
        return nil, err
    }
    
    // Create execution record
    execution := &WorkflowExecution{
        WorkflowID: workflowID,
        TenantID:   tenantID,
        ExecutedBy: userID,
        Status:     "running",
        Input:      input,
        StartedAt:  time.Now(),
    }
    
    // Execute workflow
    startTime := time.Now()
    result, execErr := e.baseEngine.Execute(workflow.Definition, input)
    duration := time.Since(startTime)
    
    // Update execution record
    execution.CompletedAt = time.Now()
    execution.DurationMs = int(duration.Milliseconds())
    
    if execErr != nil {
        execution.Status = "failed"
        execution.Error = execErr.Error()
    } else {
        execution.Status = "completed"
        execution.Output = result
    }
    
    // Save execution record
    if err := e.saveExecution(ctx, execution); err != nil {
        // Log error but don't fail the request
        logger.Error("Failed to save execution record", "error", err)
    }
    
    // Increment quota counter
    if err := e.quotaRepo.IncrementExecutionCount(ctx, tenantID); err != nil {
        logger.Error("Failed to increment quota", "error", err)
    }
    
    // Audit log
    e.auditRepo.Log(ctx, AuditEvent{
        TenantID:     tenantID,
        UserID:       userID,
        EventType:    "workflow.executed",
        ResourceType: "workflow",
        ResourceID:   workflowID,
        Metadata: map[string]interface{}{
            "duration_ms": duration.Milliseconds(),
            "status":      execution.Status,
        },
    })
    
    return &ExecutionResult{
        ExecutionID: execution.ID,
        Status:      execution.Status,
        Output:      result,
        DurationMs:  execution.DurationMs,
        Error:       execErr,
    }, execErr
}
```

### Phase 3: Local-First Configuration

#### Embedded SQLite Configuration

```go
package config

import (
    "database/sql"
    "embed"
    "errors"
    
    _ "github.com/mattn/go-sqlite3"  // Or modernc.org/sqlite for pure Go
)

//go:embed migrations/*.sql
var migrations embed.FS

type Config struct {
    DatabaseType string  // "sqlite" or "postgres"
    DatabaseURL  string
    DataDir      string  // Local data directory
}

// NewLocalConfig creates configuration for local-first deployment
func NewLocalConfig(dataDir string) *Config {
    return &Config{
        DatabaseType: "sqlite",
        DatabaseURL:  dataDir + "/thaiyyal.db",
        DataDir:      dataDir,
    }
}

// InitializeDatabase sets up the database with migrations
func (c *Config) InitializeDatabase() (*sql.DB, error) {
    var db *sql.DB
    var err error
    
    switch c.DatabaseType {
    case "sqlite":
        db, err = sql.Open("sqlite3", c.DatabaseURL+"?_foreign_keys=on")
    case "postgres":
        db, err = sql.Open("postgres", c.DatabaseURL)
    default:
        return nil, errors.New("unsupported database type")
    }
    
    if err != nil {
        return nil, err
    }
    
    // Run migrations
    if err := c.runMigrations(db); err != nil {
        return nil, err
    }
    
    return db, nil
}

// runMigrations applies all embedded SQL migrations
func (c *Config) runMigrations(db *sql.DB) error {
    // Create migrations table
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version INT PRIMARY KEY,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        return err
    }
    
    // Read and apply migrations
    files, err := migrations.ReadDir("migrations")
    if err != nil {
        return err
    }
    
    for _, file := range files {
        // Check if already applied
        var count int
        version := extractVersion(file.Name())
        err := db.QueryRow(
            "SELECT COUNT(*) FROM schema_migrations WHERE version = ?",
            version,
        ).Scan(&count)
        
        if err != nil || count > 0 {
            continue  // Skip if already applied
        }
        
        // Read migration file
        content, err := migrations.ReadFile("migrations/" + file.Name())
        if err != nil {
            return err
        }
        
        // Execute migration
        if _, err := db.Exec(string(content)); err != nil {
            return err
        }
        
        // Record migration
        _, err = db.Exec(
            "INSERT INTO schema_migrations (version) VALUES (?)",
            version,
        )
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

#### Local-First Deployment Configuration

```yaml
# config.yaml - Local deployment configuration

server:
  host: "localhost"
  port: 8080
  mode: "local"  # local, cloud, hybrid

database:
  type: "sqlite"  # sqlite, postgres
  path: "./data/thaiyyal.db"
  
  # Optional: PostgreSQL for local deployment
  # type: "postgres"
  # host: "localhost"
  # port: 5432
  # database: "thaiyyal"
  # username: "thaiyyal_user"
  # password: "secure_password"

storage:
  type: "local"  # local, s3
  path: "./data/storage"

multi_tenant:
  enabled: true
  mode: "shared_database"  # shared_database, separate_databases
  
  # Default tenant for single-tenant deployments
  default_tenant:
    slug: "default"
    name: "Default Organization"
    plan: "enterprise"

security:
  jwt_secret: "change-me-in-production"
  session_timeout: "24h"
  enable_cors: true
  allowed_origins:
    - "http://localhost:3000"
    - "http://localhost:8080"

limits:
  # Default plan limits
  free:
    max_users: 5
    max_workflows: 10
    max_executions_per_day: 100
    max_storage_mb: 100
  
  pro:
    max_users: 25
    max_workflows: 100
    max_executions_per_day: 1000
    max_storage_mb: 1000
  
  enterprise:
    max_users: -1  # unlimited
    max_workflows: -1
    max_executions_per_day: -1
    max_storage_mb: -1

features:
  enable_authentication: true
  enable_audit_logs: true
  enable_workflow_versioning: true
  enable_shared_workflows: true
  enable_api_access: true

observability:
  logging:
    level: "info"  # debug, info, warn, error
    format: "json"  # json, text
    output: "./data/logs/thaiyyal.log"
  
  metrics:
    enabled: true
    type: "prometheus"
    endpoint: "/metrics"
  
  tracing:
    enabled: false  # Enable for debugging
```

### Phase 4: Agent Collaboration for Multi-Tenancy

#### Coordination with Other Agents

**1. Security Agent Collaboration**
```markdown
Topic: Tenant Data Isolation Security Review
Security Agent Tasks:
- Review row-level security policies
- Verify tenant context middleware
- Check for SQL injection in tenant queries
- Validate JWT token tenant claims
- Review API authorization logic
```

**2. System Architecture Agent Collaboration**
```markdown
Topic: Multi-Tenant Architecture Review
Architecture Agent Tasks:
- Review database schema design
- Evaluate shared vs separate database approach
- Assess scalability of tenant isolation
- Design tenant migration strategies
- Plan for tenant-specific customizations
```

**3. Performance Agent Collaboration**
```markdown
Topic: Multi-Tenant Performance Optimization
Performance Agent Tasks:
- Index strategy for tenant-scoped queries
- Connection pooling per tenant
- Query performance with tenant_id filters
- Caching strategies for tenant data
- Database partitioning recommendations
```

**4. Observability Agent Collaboration**
```markdown
Topic: Tenant-Specific Monitoring
Observability Agent Tasks:
- Tenant-level metrics collection
- Per-tenant resource usage tracking
- Tenant-specific log aggregation
- Quota monitoring and alerting
- Tenant health dashboards
```

### Enterprise Quality Checklist

- [ ] **Data Isolation**: Complete tenant separation enforced at database level
- [ ] **Security**: Row-level security, encrypted connections, audit logs
- [ ] **Scalability**: Efficient multi-tenant queries with proper indexing
- [ ] **Local-First**: Works with SQLite without any cloud services
- [ ] **Quota Management**: Resource limits enforced per tenant
- [ ] **Audit Trail**: Complete audit log for compliance
- [ ] **Versioning**: Workflow version history per tenant
- [ ] **Authentication**: Tenant-aware JWT authentication
- [ ] **Authorization**: Role-based access control per tenant
- [ ] **Migration Path**: Easy migration from single to multi-tenant
- [ ] **Backup/Restore**: Tenant-level backup capabilities
- [ ] **Monitoring**: Per-tenant resource monitoring
- [ ] **Documentation**: Clear multi-tenant setup guide
- [ ] **Testing**: Multi-tenant isolation test suite

## Best Practices for Multi-Tenant Thaiyyal

### 1. Always Use Tenant Context
```go
// ❌ BAD: Direct query without tenant context
query := "SELECT * FROM workflows WHERE id = ?"

// ✅ GOOD: Tenant-scoped query
tenantID := middleware.GetTenantID(ctx)
query := "SELECT * FROM workflows WHERE id = ? AND tenant_id = ?"
```

### 2. Enforce Quotas Before Operations
```go
// ✅ Always check quota before resource creation
func (s *WorkflowService) Create(ctx context.Context, workflow *Workflow) error {
    if err := s.quotaRepo.Check(ctx, "workflow"); err != nil {
        return ErrQuotaExceeded
    }
    // Proceed with creation
}
```

### 3. Audit All Operations
```go
// ✅ Log all tenant operations for compliance
defer s.auditRepo.Log(ctx, AuditEvent{
    EventType: "workflow.created",
    ResourceID: workflow.ID,
    Metadata: map[string]interface{}{
        "workflow_name": workflow.Name,
    },
})
```

### 4. Test Tenant Isolation
```go
// Test that tenants can't access each other's data
func TestTenantIsolation(t *testing.T) {
    // Create workflow for tenant A
    ctxA := withTenant(context.Background(), "tenant-a")
    workflowA, _ := repo.Create(ctxA, workflow)
    
    // Try to access from tenant B
    ctxB := withTenant(context.Background(), "tenant-b")
    _, err := repo.FindByID(ctxB, workflowA.ID)
    
    assert.Error(t, err, "Should not access other tenant's workflow")
}
```

## Migration Strategy

### Single-Tenant to Multi-Tenant

```sql
-- Migration script
BEGIN;

-- 1. Add default tenant
INSERT INTO tenants (id, name, slug, plan)
VALUES ('00000000-0000-0000-0000-000000000001', 'Default Tenant', 'default', 'enterprise');

-- 2. Add tenant_id to existing tables
ALTER TABLE workflows ADD COLUMN tenant_id UUID;
ALTER TABLE users ADD COLUMN tenant_id UUID;

-- 3. Assign all existing data to default tenant
UPDATE workflows SET tenant_id = '00000000-0000-0000-0000-000000000001';
UPDATE users SET tenant_id = '00000000-0000-0000-0000-000000000001';

-- 4. Make tenant_id required
ALTER TABLE workflows ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE users ALTER COLUMN tenant_id SET NOT NULL;

-- 5. Add foreign keys
ALTER TABLE workflows ADD CONSTRAINT fk_workflows_tenant
    FOREIGN KEY (tenant_id) REFERENCES tenants(id);
ALTER TABLE users ADD CONSTRAINT fk_users_tenant
    FOREIGN KEY (tenant_id) REFERENCES tenants(id);

COMMIT;
```

---

**Version**: 1.0  
**Last Updated**: 2025-10-30  
**Maintained By**: Multi-Tenancy Team  
**Review Cycle**: Quarterly
