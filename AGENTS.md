# AI Agent Specifications for Thaiyyal

## Overview

This document provides comprehensive specifications for AI agents that can assist with various aspects of the Thaiyyal project. Each agent is specialized for specific tasks within our visual workflow builder platform, ensuring **enterprise-grade quality** with **local-first architecture** and **seamless agent collaboration**.

> **📝 Note on Agent Types**
>
> The agents described in this document are **specification documents** that serve as guidelines for:
> - Human reviewers conducting code reviews
> - AI assistants performing architectural assessments
> - Development teams implementing enterprise features
>
> **✨ NEW: Callable Agent System**
>
> We've created a **callable agent invocation system** that allows you to programmatically apply these specifications:
>
> ```bash
> # Invoke an agent from command line
> ./scripts/invoke-agent.sh security-review \
>   --files "backend/nodes_http.go" \
>   --context "Review for OWASP compliance" \
>   --output security-review.md
> ```
>
> **Documentation:**
> - **Quick Start**: [scripts/README.md](scripts/README.md)
> - **Complete Guide**: [.github/agents/AGENTS_CALLABLE_SYSTEM.md](.github/agents/AGENTS_CALLABLE_SYSTEM.md)
> - **Invocation Script**: [scripts/invoke-agent.sh](scripts/invoke-agent.sh)
>
> **How to Use These Specifications:**
> 1. **Manual**: Review the specification and apply guidelines manually
> 2. **Scripted**: Use `./scripts/invoke-agent.sh` to invoke agents programmatically
> 3. **CI/CD**: Integrate agent calls into your pipeline (see documentation)
> 4. **Pre-commit**: Add agent checks to git hooks

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

### Using Agent Specifications in Reviews

These specification documents are designed to **guide comprehensive reviews**:
- **Shared Context**: All specifications share common understanding of Thaiyyal architecture
- **Coordinated Reviews**: Use multiple agent specs for comprehensive coverage (e.g., security + architecture)
- **Sequential Workflows**: Follow specs in logical order (e.g., architecture → security → performance)
- **Parallel Analysis**: Different team members can review different aspects simultaneously
- **Cross-Domain Expertise**: Specs reference each other for holistic guidance

**Example Review Workflow:**
1. Start with **System Architecture** spec to assess overall design
2. Apply **Security Code Review** spec to identify vulnerabilities
3. Use **Performance** spec to find optimization opportunities
4. Reference **Testing & QA** spec to ensure adequate test coverage
5. Apply **Observability** spec to verify monitoring capabilities

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
├── devops-cicd.md                 # DevOps and CI/CD pipelines
├── ui-ux-architect.md             # ✨ UI/UX design and user experience
├── product-manager.md             # ✨ Project tracking and planning
├── marketing.md                   # ✨ README maintenance and product messaging
└── AGENTS_CALLABLE_SYSTEM.md      # Callable agent invocation system
```

**Callable Agent Scripts:**

```
scripts/
├── invoke-agent.sh                # Main agent invocation script
└── README.md                      # Usage examples and integration guide
```

**Quick Invocation Examples:**

```bash
# Security review
./scripts/invoke-agent.sh security-review --files "backend/*.go"

# Architecture assessment
./scripts/invoke-agent.sh architecture --files "backend/workflow.go"

# UI/UX review
./scripts/invoke-agent.sh ui-ux --files "src/components/*.tsx"

# Product planning
./scripts/invoke-agent.sh product-manager --context "Sprint planning"

# README update check
./scripts/invoke-agent.sh marketing --files "README.md"

# Full review (multiple agents)
for agent in security-review architecture performance testing ui-ux; do
  ./scripts/invoke-agent.sh $agent --files "$(git diff --name-only main)" --output "review-$agent.md"
done
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
| UI/UX design review | UI/UX Architect Agent | [ui-ux-architect.md](.github/agents/ui-ux-architect.md) |
| Accessibility compliance | UI/UX Architect Agent | [ui-ux-architect.md](.github/agents/ui-ux-architect.md) |
| Product roadmap planning | Product Manager Agent | [product-manager.md](.github/agents/product-manager.md) |
| Sprint planning & tracking | Product Manager Agent | [product-manager.md](.github/agents/product-manager.md) |
| README.md maintenance | Marketing Agent | [marketing.md](.github/agents/marketing.md) |
| Product messaging & content | Marketing Agent | [marketing.md](.github/agents/marketing.md) |

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

### 9. UI/UX Architect Agent
**Focus**: User interface design, user experience optimization, accessibility

**Key Responsibilities**:
- Visual workflow builder interface design
- User experience optimization
- Accessibility compliance (WCAG 2.1 AA)
- Design system creation and maintenance
- Responsive design
- Usability testing
- Interactive prototypes

**[Full Specification →](.github/agents/ui-ux-architect.md)**

---

### 10. Product Manager Agent
**Focus**: Project tracking, planning, roadmap management

**Key Responsibilities**:
- Product strategy and vision
- Roadmap management
- Backlog prioritization
- Sprint planning and execution
- Stakeholder communication
- Metrics and analytics
- User feedback management

**[Full Specification →](.github/agents/product-manager.md)**

---

### 11. Marketing Agent
**Focus**: README maintenance, product messaging, content strategy

**Key Responsibilities**:
- README.md maintenance and accuracy
- Product messaging and positioning
- Content strategy and creation
- Documentation synchronization
- Community engagement
- Release announcements
- Feature highlights

**[Full Specification →](.github/agents/marketing.md)**

---

## Usage Guidelines

### How to Use Agent Specifications

These specifications serve as **comprehensive guides** for reviewing and improving Thaiyyal:

1. **Identify the Domain**: Determine which specification matches your work area (security, architecture, etc.)
2. **Review the Specification**: Read the detailed guidelines in the agent's dedicated file
3. **Apply the Standards**: Use the spec's checklists, patterns, and best practices
4. **Document Findings**: Record issues and recommendations based on the specification
5. **Validate Changes**: Ensure implementations follow the spec's guidelines
6. **Cross-Reference**: When work spans domains, consult multiple specifications

### Best Practices for Using Specifications

- **Domain Focus**: Apply one specification at a time for thorough coverage
- **Context Awareness**: Consider Thaiyyal's architecture when applying guidelines
- **Standards Validation**: Verify that implementations meet the specification's requirements
- **Incremental Application**: Apply recommendations incrementally, testing at each step
- **Documentation**: Record which specifications were applied and findings discovered
- **Multi-Domain Coordination**: For complex work, apply multiple specifications sequentially

### Specification Application Patterns

### Pattern 1: Sequential Application (Enterprise Feature Development)
When building a new enterprise feature, apply specifications in this order:
```
1. System Architecture Spec
   ↓ (Review architectural design)
2. Multi-Tenancy Spec
   ↓ (Verify tenant isolation)
3. Security Code Review Spec
   ↓ (Check for vulnerabilities)
4. Performance Spec
   ↓ (Optimize implementation)
5. Testing & QA Spec
   ↓ (Ensure test coverage)
6. Documentation Spec
   ↓ (Document the feature)
7. DevOps & CI/CD Spec
   ↓ (Automate deployment)
Final Implementation Review
```

### Pattern 2: Parallel Application (Code Review)
For comprehensive code reviews, multiple reviewers can apply different specs simultaneously:
```
Reviewer 1: Security Spec ──┐
                            │
Reviewer 2: Performance Spec├──→ Consolidated Findings
                            │      ↓
Reviewer 3: Testing Spec ───┤   Action Items
                            │      ↓
Reviewer 4: Docs Spec ──────┘   Implementation
```

### Pattern 3: Iterative Application (Quality Enhancement)
For improving existing code, apply specifications iteratively:
```
Current Implementation
    ↓
Apply Testing Spec (identify test gaps)
    ↓
Apply Performance Spec (optimize bottlenecks)
    ↓
Apply Security Spec (harden security)
    ↓
Apply Observability Spec (add monitoring)
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
