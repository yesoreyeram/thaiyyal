# AI Agent Specifications for Thaiyyal

## Overview

This document provides comprehensive specifications for AI agents that can assist with various aspects of the Thaiyyal project. Each agent is specialized for specific tasks within our visual workflow builder platform, ensuring **enterprise-grade quality** with **local-first architecture** and **seamless agent collaboration**.

### Enterprise Quality Standards

All agents are designed to deliver:
- **Production-Ready Code**: Battle-tested patterns and best practices
- **Comprehensive Testing**: Unit, integration, and E2E test coverage ≥80%
- **Security First**: OWASP compliance, vulnerability scanning, secure defaults
- **Performance**: Optimized for low latency and high throughput
- **Observability**: Full metrics, logging, and tracing
- **Multi-Tenancy**: Complete tenant isolation and resource management
- **Documentation**: Extensive technical and user documentation

### Local-First Architecture

Thaiyyal is designed to run **locally without any cloud dependencies**:
- **Embedded Database**: SQLite for local deployments (PostgreSQL optional)
- **Single Binary**: Backend compiles to standalone executable
- **No Cloud Lock-in**: All features work offline
- **Optional Cloud**: Cloud services are additive, not required
- **Docker Support**: Easy containerized local deployment
- **Fast Setup**: Running locally in under 5 minutes

### Agent Collaboration Framework

Agents are designed to **work together efficiently**:
- **Shared Context**: Common understanding of Thaiyyal architecture
- **Coordinated Reviews**: Multi-agent code reviews for comprehensive coverage
- **Sequential Workflows**: Agents can hand-off tasks in logical order
- **Parallel Analysis**: Multiple agents can analyze different aspects simultaneously
- **Cross-Domain Expertise**: Agents consult each other for specialized knowledge

### Thaiyyal Technology Stack

- **Frontend**: Next.js 16.0 + React 19.2 + ReactFlow + TypeScript
- **Backend**: Go 1.24.7 workflow execution engine
- **Architecture**: Client-side workflow builder with DAG-based execution
- **Node Types**: 23 different node types for data processing, control flow, and error handling
- **Database**: SQLite (local) or PostgreSQL (production)
- **Deployment**: Local, Docker, Kubernetes, Cloud (AWS/GCP/Azure)

## Agent Directory Structure

All specialized agent specifications are located in `.github/agents/`:

```
.github/agents/
├── security-code-review.md       # Security analysis and code review
├── system-architecture.md         # Architecture design and review
├── observability.md               # Monitoring, logging, and tracing
├── multi-tenancy.md               # Multi-tenant architecture specialist
├── testing-qa.md                  # Testing strategies and quality assurance
├── performance.md                 # Performance optimization
├── documentation.md               # Technical documentation
└── devops-cicd.md                 # DevOps and CI/CD pipelines
```

## Quick Reference Guide

### When to Use Each Agent

| Task | Agent | File |
|------|-------|------|
| Security vulnerability assessment | Security Code Review Agent | [security-code-review.md](.github/agents/security-code-review.md) |
| Code review for security issues | Security Code Review Agent | [security-code-review.md](.github/agents/security-code-review.md) |
| Architecture design decisions | System Architecture Agent | [system-architecture.md](.github/agents/system-architecture.md) |
| System refactoring planning | System Architecture Agent | [system-architecture.md](.github/agents/system-architecture.md) |
| Adding monitoring/logging | Observability Agent | [observability.md](.github/agents/observability.md) |
| Performance metrics setup | Observability Agent | [observability.md](.github/agents/observability.md) |
| Multi-tenant feature design | Multi-Tenancy Specialist | [multi-tenancy.md](.github/agents/multi-tenancy.md) |
| Tenant isolation implementation | Multi-Tenancy Specialist | [multi-tenancy.md](.github/agents/multi-tenancy.md) |
| Test strategy development | Testing & QA Agent | [testing-qa.md](.github/agents/testing-qa.md) |
| Test coverage improvement | Testing & QA Agent | [testing-qa.md](.github/agents/testing-qa.md) |
| Performance optimization | Performance Agent | [performance.md](.github/agents/performance.md) |
| Bottleneck identification | Performance Agent | [performance.md](.github/agents/performance.md) |
| Documentation writing | Documentation Agent | [documentation.md](.github/agents/documentation.md) |
| API documentation | Documentation Agent | [documentation.md](.github/agents/documentation.md) |
| CI/CD pipeline setup | DevOps Agent | [devops-cicd.md](.github/agents/devops-cicd.md) |
| Deployment automation | DevOps Agent | [devops-cicd.md](.github/agents/devops-cicd.md) |

## Agent Specializations

### 1. Security Code Review Agent
**Focus**: Security analysis, vulnerability detection, secure coding practices

**Key Responsibilities**:
- Security vulnerability scanning and remediation
- Code review with security focus
- Authentication and authorization implementation
- Input validation and sanitization
- Secure data handling practices
- Dependency security auditing
- Security best practices enforcement

**[Full Specification →](.github/agents/security-code-review.md)**

---

### 2. System Design & Architecture Agent
**Focus**: System architecture, design patterns, scalability

**Key Responsibilities**:
- High-level architecture design
- Design pattern recommendations
- Component interaction design
- Scalability planning
- Technology stack decisions
- Refactoring strategies
- Architecture documentation

**[Full Specification →](.github/agents/system-architecture.md)**

---

### 3. Observability Agent
**Focus**: Monitoring, logging, tracing, metrics

**Key Responsibilities**:
- Logging strategy design
- Metrics collection implementation
- Distributed tracing setup
- Performance monitoring
- Error tracking and alerting
- Dashboard creation
- Observability best practices

**[Full Specification →](.github/agents/observability.md)**

---

### 4. Multi-Tenancy Specialist Agent
**Focus**: Multi-tenant architecture, tenant isolation, resource management

**Key Responsibilities**:
- Multi-tenant architecture design
- Tenant isolation strategies
- Resource quota management
- Data segregation implementation
- Tenant-specific customization
- Billing and metering
- Multi-tenant security

**[Full Specification →](.github/agents/multi-tenancy.md)**

---

### 5. Testing & Quality Assurance Agent
**Focus**: Test strategies, test automation, quality metrics

**Key Responsibilities**:
- Test strategy development
- Unit test implementation
- Integration test design
- E2E test automation
- Test coverage analysis
- Quality metrics tracking
- Testing best practices

**[Full Specification →](.github/agents/testing-qa.md)**

---

### 6. Performance Optimization Agent
**Focus**: Performance analysis, optimization, benchmarking

**Key Responsibilities**:
- Performance profiling
- Bottleneck identification
- Code optimization
- Resource usage optimization
- Caching strategies
- Query optimization
- Performance benchmarking

**[Full Specification →](.github/agents/performance.md)**

---

### 7. Documentation Agent
**Focus**: Technical documentation, API docs, user guides

**Key Responsibilities**:
- Technical documentation writing
- API documentation
- User guide creation
- Code documentation
- Architecture diagrams
- Tutorial development
- Documentation maintenance

**[Full Specification →](.github/agents/documentation.md)**

---

### 8. DevOps & CI/CD Agent
**Focus**: Deployment automation, CI/CD pipelines, infrastructure

**Key Responsibilities**:
- CI/CD pipeline design
- Deployment automation
- Infrastructure as Code
- Container orchestration
- Monitoring infrastructure
- Release management
- DevOps best practices

**[Full Specification →](.github/agents/devops-cicd.md)**

---

## Usage Guidelines

### How to Work with Agents

1. **Identify the Task**: Determine which agent specialization best matches your task
2. **Review Agent Specification**: Read the detailed specification in the agent's dedicated file
3. **Provide Context**: Give the agent relevant context about Thaiyyal's architecture
4. **Set Clear Objectives**: Define specific, measurable goals for the agent
5. **Review Agent Output**: Carefully review and validate agent recommendations
6. **Iterate as Needed**: Work iteratively with the agent to refine solutions

### Best Practices

- **Single Responsibility**: Use one agent per task for focused expertise
- **Context First**: Always provide project context before requesting agent assistance
- **Validate Recommendations**: Review agent suggestions against project requirements
- **Incremental Changes**: Implement changes incrementally, testing at each step
- **Documentation**: Document agent recommendations and implementation decisions
- **Cross-Agent Coordination**: When tasks span multiple domains, consult relevant agents sequentially

### Agent Interaction Patterns

### Pattern 1: Sequential Consultation (Enterprise Feature Development)
```
System Architecture Agent
    ↓ (Architectural design)
Multi-Tenancy Specialist
    ↓ (Tenant isolation design)
Security Code Review Agent
    ↓ (Security review)
Performance Optimization Agent
    ↓ (Performance optimization)
Testing & QA Agent
    ↓ (Test strategy and implementation)
Documentation Agent
    ↓ (Documentation)
DevOps & CI/CD Agent
    ↓ (Deployment automation)
Final Implementation
```

### Pattern 2: Parallel Consultation (Code Review)
```
Security Code Review Agent ──┐
                             │
Performance Agent        ────┤──→ Comprehensive Review
                             │      ↓
Testing & QA Agent       ────┤   Synthesis
                             │      ↓
Documentation Agent      ────┘   Implementation
```

### Pattern 3: Iterative Refinement (Quality Enhancement)
```
Initial Implementation
    ↓
Testing Agent (finds issues)
    ↓
Performance Agent (optimizes)
    ↓
Security Agent (hardens)
    ↓
Observability Agent (instruments)
    ↓
Final Production-Ready Code
```

### Pattern 4: Multi-Tenant Feature Addition
```
Requirements
    ↓
┌───────────────────────────────────────┐
│ Multi-Tenancy Specialist Agent        │
│ - Design tenant isolation             │
│ - Define data schema                  │
└───────────────┬───────────────────────┘
                ↓
┌───────────────────────────────────────┐
│ Security Code Review Agent            │
│ - Verify row-level security           │
│ - Check authorization logic           │
└───────────────┬───────────────────────┘
                ↓
┌───────────────────────────────────────┐
│ Performance Optimization Agent        │
│ - Optimize tenant-scoped queries      │
│ - Design efficient indexes            │
└───────────────┬───────────────────────┘
                ↓
┌───────────────────────────────────────┐
│ Testing & QA Agent                    │
│ - Test tenant isolation               │
│ - Verify quota enforcement            │
└───────────────┬───────────────────────┘
                ↓
Production-Ready Multi-Tenant Feature
```

### Example Collaboration: Adding User Authentication

**Step 1: Architecture Design**
- **System Architecture Agent**: Designs auth flow, token management
- **Multi-Tenancy Specialist**: Ensures tenant-aware authentication
- **Output**: Architecture decision record (ADR)

**Step 2: Security Review**
- **Security Code Review Agent**: Reviews auth implementation
- **Validates**: Password hashing, JWT tokens, session management
- **Output**: Security checklist, vulnerability report

**Step 3: Implementation Review**
- **Performance Agent**: Reviews auth query performance
- **Testing Agent**: Creates auth test suite
- **Output**: Optimized code with comprehensive tests

**Step 4: Deployment**
- **Observability Agent**: Adds auth metrics and logging
- **DevOps Agent**: Updates CI/CD for auth changes
- **Documentation Agent**: Documents auth API
- **Output**: Production-ready, documented, monitored feature

## Project-Specific Context for Agents

### Technology Stack

**Frontend**:
- Next.js 16.0.1 (App Router)
- React 19.2.0
- ReactFlow 11.8.0 (visual workflow canvas)
- TypeScript 5
- Tailwind CSS 4
- Browser LocalStorage for persistence

**Backend**:
- Go 1.24.7
- Standard library only (zero external dependencies)
- DAG-based workflow execution
- Topological sorting (Kahn's algorithm)
- In-memory state management

### Current Architecture

```
Frontend (Next.js) ──generates──> JSON Workflow
                                      │
                                      ▼
                                 Backend (Go)
                                      │
                        ┌─────────────┼─────────────┐
                        ▼             ▼             ▼
                    Parse JSON    Infer Types   Validate
                        │             │             │
                        └─────────────┼─────────────┘
                                      ▼
                              Topological Sort
                                      ▼
                              Execute Nodes
                                      ▼
                              Return Results
```

### Node Types (23 Types)

1. **Basic I/O**: Number, TextInput, Visualization
2. **Operations**: Math, Text, Transform, Extract
3. **HTTP**: HTTP Request
4. **Control Flow**: Condition, ForEach, WhileLoop, Switch
5. **Parallel**: Parallel, Join, Split
6. **State**: Variable, Accumulator, Counter, Cache
7. **Error Handling**: Retry, TryCatch, Timeout
8. **Utility**: Delay
9. **Context**: ContextVariable, ContextConstant

### Key Design Principles

1. **Minimal Dependencies**: Backend uses only Go standard library
2. **Type Inference**: Automatic node type detection from data
3. **DAG Execution**: Deterministic execution order via topological sort
4. **Client-Side First**: Browser-based workflow building and storage
5. **Progressive Enhancement**: Start simple, add complexity as needed

### Current Limitations

1. No backend HTTP API (frontend generates JSON only)
2. No persistence layer (LocalStorage only)
3. Single-workflow execution scope
4. No real-time collaboration
5. No workflow versioning
6. Limited error recovery mechanisms

### Security Considerations

**Current Measures**:
- Input validation in node executors
- Type checking for operations
- Cycle detection (prevents infinite loops)
- Client-side only (no server-side data exposure)

**Areas for Improvement** (See Security Agent):
- Operation timeouts
- HTTP URL whitelisting
- Rate limiting
- Input size limits
- Execution quotas
- Authentication/authorization for API

## Agent Output Expectations

### All Agents Should Provide

1. **Analysis**: Clear assessment of current state
2. **Recommendations**: Specific, actionable suggestions
3. **Rationale**: Explanation of why recommendations are made
4. **Implementation Guidance**: Step-by-step implementation approach
5. **Trade-offs**: Discussion of pros/cons of proposed solutions
6. **Testing Strategy**: How to validate the changes
7. **Documentation**: Updates needed to project documentation
8. **Local-First Compatibility**: Ensure solutions work without cloud dependencies
9. **Multi-Tenant Awareness**: Consider tenant isolation in all designs
10. **Enterprise Quality**: Production-ready code with proper error handling

### Code Contributions Should Include

- Clean, well-commented code following project conventions
- Tests covering new functionality (unit, integration, E2E)
- Updated documentation (code comments, README, API docs)
- Performance considerations and benchmarks
- Security implications addressed
- Backward compatibility notes
- Multi-tenant isolation verification
- Local deployment instructions
- Observability instrumentation (metrics, logs, traces)

## Local-First Deployment Guide

### Quick Local Setup (Under 5 Minutes)

```bash
# 1. Clone repository
git clone https://github.com/yesoreyeram/thaiyyal.git
cd thaiyyal

# 2. Install dependencies
npm install

# 3. Run development server
npm run dev

# 4. Access application
open http://localhost:3000
```

### Local Production Deployment

**Option 1: Single Binary (Recommended)**
```bash
# Build frontend
npm run build

# Build backend with embedded frontend
cd backend
go build -o thaiyyal ./cmd/server

# Run with SQLite (no database setup needed)
./thaiyyal --data-dir=./data
```

**Option 2: Docker Compose**
```bash
# Start all services (app + postgres + monitoring)
docker-compose up -d

# Access application
open http://localhost:8080

# View logs
docker-compose logs -f thaiyyal
```

**Option 3: Systemd Service (Linux)**
```bash
# Install binary
sudo cp thaiyyal /usr/local/bin/

# Create service
sudo cat > /etc/systemd/system/thaiyyal.service << EOF
[Unit]
Description=Thaiyyal Workflow Builder
After=network.target

[Service]
Type=simple
User=thaiyyal
WorkingDirectory=/opt/thaiyyal
ExecStart=/usr/local/bin/thaiyyal --config=/etc/thaiyyal/config.yaml
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

# Start service
sudo systemctl enable thaiyyal
sudo systemctl start thaiyyal
```

### Local Multi-Tenant Setup

```yaml
# config.yaml - Local multi-tenant configuration
server:
  host: "localhost"
  port: 8080

database:
  type: "sqlite"  # No PostgreSQL needed locally
  path: "./data/thaiyyal.db"

multi_tenant:
  enabled: true
  mode: "shared_database"
  
  # Pre-configured tenants for local development
  default_tenants:
    - slug: "acme"
      name: "ACME Corporation"
      plan: "enterprise"
    - slug: "demo"
      name: "Demo Organization"
      plan: "pro"

# All features work locally without cloud
features:
  enable_authentication: true
  enable_audit_logs: true
  enable_workflow_versioning: true
  enable_api_access: true

# Local observability stack
observability:
  logging:
    output: "./data/logs/thaiyyal.log"
  metrics:
    enabled: true
    endpoint: "/metrics"
```

### Cloud-Optional Features

**Works Locally (No Cloud Needed)**:
- ✅ Workflow creation and execution
- ✅ Multi-tenant isolation
- ✅ User authentication and authorization
- ✅ Workflow versioning
- ✅ Audit logging
- ✅ Metrics and monitoring (Prometheus + Grafana)
- ✅ API access
- ✅ Role-based access control

**Cloud-Enhanced (Optional)**:
- ☁️ Cloud storage for workflow exports (S3/GCS)
- ☁️ Cloud-based alerting (PagerDuty/Slack)
- ☁️ External identity providers (Auth0/Okta)
- ☁️ Cloud logging services (DataDog/LogDNA)
- ☁️ CDN for static assets (CloudFront/CloudFlare)

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-10-30 | Initial agent specification framework |

## Contributing to Agent Specifications

To improve or extend agent specifications:

1. Review existing agent files
2. Propose changes via pull request
3. Ensure consistency with project architecture
4. Update this index file if adding new agents
5. Test specifications with real use cases

## Related Documentation

- [README.md](README.md) - Project overview and quick start
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture details
- [ARCHITECTURE_REVIEW.md](ARCHITECTURE_REVIEW.md) - Architectural analysis
- [backend/README.md](backend/README.md) - Backend workflow engine
- [docs/NODES.md](docs/NODES.md) - Complete node type reference

## Support and Questions

For questions about using these agent specifications:
- Open an issue on GitHub
- Review existing agent documentation
- Consult the main project documentation

---

**Last Updated**: 2025-10-30  
**Maintainer**: Thaiyyal Team  
**License**: MIT
