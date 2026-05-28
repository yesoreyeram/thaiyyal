# Thaiyyal - Product Requirements Document (PRD)

**Version:** 1.0
**Date:** 2026-05-28
**Status:** Active
**Owner:** Product Team

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Problem Definition](#problem-definition)
3. [Goals & Objectives](#goals--objectives)
4. [User Personas](#user-personas)
5. [User Stories](#user-stories)
6. [Feature Requirements](#feature-requirements)
7. [User Experience](#user-experience)
8. [Success Criteria](#success-criteria)
9. [Out of Scope](#out-of-scope)
10. [Dependencies & Risks](#dependencies--risks)

---

## Executive Summary

### What is Thaiyyal?

Thaiyyal is a **visual workflow automation platform** that enables users to create, execute, and manage complex data processing workflows without writing code. By providing an intuitive drag-and-drop interface powered by a secure and performant execution engine, Thaiyyal democratizes workflow automation for users of all technical skill levels.

### Why Now?

- **Market Need**: Growing demand for no-code/low-code automation tools
- **Technical Readiness**: Modern web technologies (React, Next.js) enable rich visual experiences
- **Security Concerns**: Existing tools lack proper security controls (SSRF, resource limits)
- **Open Source Gap**: Few open-source alternatives with enterprise-grade features

### Key Value Propositions

1. **No Coding Required**: Visual interface accessible to non-technical users
2. **Security First**: Enterprise-grade security built-in, not bolted on
3. **Fully Open Source**: No vendor lock-in, complete transparency
4. **Production Ready**: Deploy anywhere (local, cloud, on-premise)
5. **Extensible**: Plugin architecture for custom integrations

---

## Problem Definition

### Current Challenges

#### For Non-Technical Users
- **Problem**: Cannot automate workflows without learning programming
- **Impact**: Manual, error-prone processes; inefficiency
- **Evidence**: 70% of data analysts spend time on manual data wrangling

#### For Developers
- **Problem**: Building workflow automation requires extensive boilerplate
- **Impact**: Slow iteration, high maintenance burden
- **Evidence**: Average time to build custom workflow engine: 3-6 months

#### For Organizations
- **Problem**: Existing tools are expensive and create vendor lock-in
- **Impact**: High TCO, limited flexibility, compliance concerns
- **Evidence**: Enterprise workflow platforms cost $50-500K+ annually

### Why Existing Solutions Fall Short

| Solution Type | Pros | Cons | Gap Thaiyyal Fills |
|--------------|------|------|-------------------|
| **Zapier, n8n** | Easy to use, many integrations | Expensive, limited control, vendor lock-in | Open source, self-hosted, extensible |
| **Apache Airflow** | Powerful, scalable | Requires coding, complex setup | Visual interface, easier onboarding |
| **Node-RED** | Visual, IoT-focused | Weak typing, security gaps, performance issues | Strong typing, security, performance |
| **Custom Scripts** | Full control | High maintenance, no standardization | Standardized, visual, maintainable |

---

## Goals & Objectives

### Primary Goals

1. **Enable Visual Automation**: Allow users to create workflows visually without coding
2. **Ensure Security**: Provide enterprise-grade security by default
3. **Deliver Performance**: Support complex workflows with 10,000+ nodes
4. **Foster Extensibility**: Enable community to build custom integrations

### Success Criteria (6 Months)

- [ ] 1,000 monthly active users
- [ ] 95%+ user satisfaction rating
- [ ] <5% workflow error rate
- [ ] 1,000+ GitHub stars
- [ ] 10+ community contributors

### Success Criteria (12 Months)

- [ ] 10,000 monthly active users
- [ ] 100,000+ workflow executions/month
- [ ] 5,000+ GitHub stars
- [ ] 50+ community contributors
- [ ] Enterprise deployment case studies

---

## User Personas

### Persona 1: Sarah - Data Analyst

**Demographics:**
- Age: 28
- Role: Data Analyst at mid-size company
- Experience: 3 years in analytics, basic SQL, no programming

**Goals:**
- Process data from multiple sources efficiently
- Create automated reports
- Reduce manual work

**Pain Points:**
- Excel limitations for complex transformations
- Waiting on engineering team for simple automation
- Difficulty debugging data pipelines

**How Thaiyyal Helps:**
- Visual interface requires no coding
- Built-in data transformation nodes
- Clear error messages and debugging

**User Journey:**
1. Needs to combine data from 3 APIs
2. Opens Thaiyyal, creates new workflow
3. Adds HTTP nodes for each API
4. Adds transformation nodes (extract, map, filter)
5. Adds visualization node to see results
6. Saves workflow for daily reuse

### Persona 2: Mike - DevOps Engineer

**Demographics:**
- Age: 35
- Role: Senior DevOps Engineer
- Experience: 10+ years, expert in Linux, Docker, K8s

**Goals:**
- Automate operational workflows
- Monitor system health
- Reduce on-call burden

**Pain Points:**
- Maintaining complex shell scripts
- No visibility into automation execution
- Scripts break in unexpected ways

**How Thaiyyal Helps:**
- Visual workflows easier to understand and maintain
- Built-in observability (logs, metrics, traces)
- Type safety reduces runtime errors

**User Journey:**
1. Needs to automate service health checks
2. Creates workflow with HTTP health check nodes
3. Adds conditional logic for alerting
4. Deploys to Kubernetes
5. Monitors via Prometheus/Grafana
6. Easily modifies workflow when requirements change

### Persona 3: Alex - Backend Developer

**Demographics:**
- Age: 30
- Role: Backend Developer
- Experience: 5 years, proficient in Python/JavaScript/Go

**Goals:**
- Quickly prototype business logic
- Build reliable integrations
- Minimize boilerplate code

**Pain Points:**
- Slow iteration on complex logic
- Testing and debugging integration code
- Maintaining custom workflow engines

**How Thaiyyal Helps:**
- Visual prototyping faster than coding
- Extensibility via custom executors
- Built-in testing and validation

**User Journey:**
1. Needs to prototype complex approval workflow
2. Uses Thaiyyal to visualize logic
3. Tests with different input data
4. Presents to stakeholders for feedback
5. Iterates quickly based on feedback
6. Eventually converts to production code or keeps in Thaiyyal

---

## User Stories

### Epic 1: Workflow Creation

#### Story 1.1: Create First Workflow
**As a** new user
**I want to** create my first workflow easily
**So that** I can understand how the platform works

**Acceptance Criteria:**
- [ ] Landing page has "Create Workflow" button
- [ ] Click opens visual editor
- [ ] Node palette is visible with categories
- [ ] Example workflows are available as templates
- [ ] Help tooltip guides first-time users

**Priority:** P0 (Must Have)

#### Story 1.2: Add Nodes to Workflow
**As a** workflow creator
**I want to** drag and drop nodes onto the canvas
**So that** I can build my workflow visually

**Acceptance Criteria:**
- [ ] Can drag node from palette to canvas
- [ ] Node appears on canvas at drop location
- [ ] Node has visual representation (icon, label)
- [ ] Can select node to view properties
- [ ] Can delete node

**Priority:** P0 (Must Have)

#### Story 1.3: Connect Nodes
**As a** workflow creator
**I want to** connect nodes with edges
**So that** I can define data flow

**Acceptance Criteria:**
- [ ] Can drag from node output to node input
- [ ] Edge is drawn between nodes
- [ ] Cannot create cycles (validation error shown)
- [ ] Can delete edge
- [ ] Visual feedback when hovering over valid targets

**Priority:** P0 (Must Have)

#### Story 1.4: Configure Node Properties
**As a** workflow creator
**I want to** configure node settings
**So that** I can customize behavior

**Acceptance Criteria:**
- [ ] Click node opens properties panel
- [ ] Panel shows relevant fields for node type
- [ ] Changes are saved automatically
- [ ] Invalid values show error messages
- [ ] Help text explains each field

**Priority:** P0 (Must Have)

### Epic 2: Workflow Execution

#### Story 2.1: Execute Workflow
**As a** workflow user
**I want to** run my workflow
**So that** I can see the results

**Acceptance Criteria:**
- [ ] "Run" button is visible
- [ ] Click executes workflow
- [ ] Loading indicator shows execution in progress
- [ ] Results panel opens automatically
- [ ] Can view results while workflow runs

**Priority:** P0 (Must Have)

#### Story 2.2: View Execution Results
**As a** workflow user
**I want to** see detailed execution results
**So that** I can verify the output

**Acceptance Criteria:**
- [ ] Results panel shows execution summary
- [ ] Can see final workflow output
- [ ] Can expand to see per-node results
- [ ] Execution time is displayed
- [ ] Success/failure status is clear

**Priority:** P0 (Must Have)

#### Story 2.3: Handle Execution Errors
**As a** workflow user
**I want to** understand why my workflow failed
**So that** I can fix the problem

**Acceptance Criteria:**
- [ ] Error message is displayed clearly
- [ ] Error identifies which node failed
- [ ] Error message explains what went wrong
- [ ] Suggests how to fix (if possible)
- [ ] Can view full error details

**Priority:** P0 (Must Have)

#### Story 2.4: Cancel Running Workflow
**As a** workflow user
**I want to** stop a long-running workflow
**So that** I can make changes and retry

**Acceptance Criteria:**
- [ ] "Cancel" button available during execution
- [ ] Click immediately stops execution
- [ ] Results panel shows "Cancelled" status
- [ ] Partial results are available
- [ ] Can edit and re-run immediately

**Priority:** P1 (Should Have)

### Epic 3: Workflow Management

#### Story 3.1: Save Workflow
**As a** workflow creator
**I want to** save my workflow
**So that** I can use it later

**Acceptance Criteria:**
- [ ] "Save" button available
- [ ] Can name the workflow
- [ ] Workflow saved to local storage
- [ ] Confirmation message shown
- [ ] Can overwrite existing workflow

**Priority:** P0 (Must Have)

#### Story 3.2: Load Workflow
**As a** workflow user
**I want to** load a previously saved workflow
**So that** I can continue working on it

**Acceptance Criteria:**
- [ ] "Open" button shows list of saved workflows
- [ ] Can search/filter workflows
- [ ] Click loads workflow to canvas
- [ ] Unsaved changes prompt confirmation
- [ ] Can delete saved workflows

**Priority:** P0 (Must Have)

#### Story 3.3: Export/Import Workflow
**As a** workflow creator
**I want to** export my workflow as JSON
**So that** I can share it or back it up

**Acceptance Criteria:**
- [ ] "Export" button downloads JSON file
- [ ] "Import" button accepts JSON file
- [ ] Validates imported workflows
- [ ] Shows error for invalid JSON
- [ ] Preserves all workflow settings

**Priority:** P1 (Should Have)

### Epic 4: Node Types & Operations

#### Story 4.1: Basic Math Operations
**As a** workflow user
**I want to** perform basic math operations
**So that** I can process numeric data

**Acceptance Criteria:**
- [ ] Number input node provides constants
- [ ] Operation node supports: +, -, *, /, %
- [ ] Results are accurate
- [ ] Handles edge cases (division by zero)
- [ ] Error messages are clear

**Priority:** P0 (Must Have)

#### Story 4.2: HTTP Requests
**As a** workflow user
**I want to** fetch data from APIs
**So that** I can process external data

**Acceptance Criteria:**
- [ ] HTTP node supports GET, POST, PUT, DELETE
- [ ] Can set headers and query parameters
- [ ] Automatically parses JSON responses
- [ ] Handles errors gracefully
- [ ] SSRF protection blocks dangerous URLs

**Priority:** P0 (Must Have)

#### Story 4.3: Array Transformations
**As a** workflow user
**I want to** transform arrays of data
**So that** I can filter, map, and aggregate

**Acceptance Criteria:**
- [ ] Filter node removes unwanted items
- [ ] Map node transforms each item
- [ ] Reduce node aggregates to single value
- [ ] Expression syntax is intuitive
- [ ] Examples provided in help text

**Priority:** P0 (Must Have)

#### Story 4.4: Conditional Logic
**As a** workflow user
**I want to** execute different paths based on conditions
**So that** I can handle different scenarios

**Acceptance Criteria:**
- [ ] Condition node evaluates boolean expression
- [ ] Executes true or false branch
- [ ] Switch node supports multiple cases
- [ ] Expression syntax documented
- [ ] Visual indication of which path executed

**Priority:** P0 (Must Have)

### Epic 5: Advanced Features

#### Story 5.1: Variables & State
**As a** workflow user
**I want to** store and reuse values
**So that** I can build complex workflows

**Acceptance Criteria:**
- [ ] Can define workflow-level variables
- [ ] Variable node gets variable value
- [ ] Accumulator node updates state
- [ ] Counter node increments/decrements
- [ ] Variables scoped properly

**Priority:** P1 (Should Have)

#### Story 5.2: Error Handling
**As a** workflow user
**I want to** handle errors gracefully
**So that** one failure doesn't stop the entire workflow

**Acceptance Criteria:**
- [ ] Try/catch node catches errors
- [ ] Can provide fallback value
- [ ] Retry node retries on failure
- [ ] Exponential backoff implemented
- [ ] Timeout node enforces limits

**Priority:** P1 (Should Have)

#### Story 5.3: Loops & Iteration
**As a** workflow user
**I want to** iterate over collections
**So that** I can process multiple items

**Acceptance Criteria:**
- [ ] ForEach node iterates array
- [ ] While node loops with condition
- [ ] Loop limits prevent infinite loops
- [ ] Can access loop index
- [ ] Can break early (via condition)

**Priority:** P1 (Should Have)

---

## Feature Requirements

### Must Have (P0) - Version 0.1

- [x] Visual workflow editor (React Flow)
- [x] 40+ built-in node types
- [x] Workflow execution engine
- [x] Real-time execution results
- [x] Local storage persistence
- [x] SSRF protection
- [x] Resource limits (timeout, node count, HTTP calls)
- [x] Error handling and display
- [x] Workflow validation
- [x] Import/export workflows

### Should Have (P1) - Version 1.0

- [ ] Workflow persistence (database)
- [ ] Execution history
- [ ] Workflow versioning
- [ ] User authentication
- [ ] Workflow templates library
- [ ] Advanced error recovery
- [ ] Performance optimizations
- [ ] Comprehensive documentation

### Nice to Have (P2) - Version 1.1+

- [ ] Real-time collaboration
- [ ] Workflow scheduling (cron)
- [ ] Visual debugger
- [ ] Breakpoints and step-through
- [ ] GraphQL API
- [ ] WebSocket execution
- [ ] Mobile-responsive UI
- [ ] Workflow analytics

### Won't Have (Out of Scope)

- Native mobile apps (web-first)
- Built-in authentication providers (use existing solutions)
- Database as a service (users provide their own)
- AI-powered workflow generation (future consideration)

---

## User Experience

### Design Principles

1. **Clarity Over Cleverness**: Prefer obvious over clever UI patterns
2. **Progressive Disclosure**: Show complexity only when needed
3. **Immediate Feedback**: Provide real-time validation and feedback
4. **Forgiving**: Easy to undo mistakes
5. **Accessible**: Follow WCAG 2.1 AA guidelines

### Key Workflows

#### New User Onboarding (5 minutes)
1. Land on homepage → Clear value proposition
2. Click "Try Demo" → Opens pre-built example workflow
3. Click "Run" → See results immediately
4. Click "Edit" → Modify a node value
5. Click "Run" again → See different results
6. Click "Create New" → Start blank workflow

**Success**: User understands core concepts and feels confident

#### Creating a Workflow (10 minutes)
1. Open visual editor
2. Browse node palette (categorized)
3. Drag nodes to canvas
4. Connect nodes with edges
5. Configure node properties
6. Click "Run" to test
7. Fix any errors
8. Save workflow

**Success**: User creates a working workflow independently

#### Debugging a Failed Workflow (5 minutes)
1. Execute workflow → Fails
2. View error in results panel
3. Error identifies failing node
4. Click error → Opens properties for that node
5. See what's wrong (e.g., invalid URL)
6. Fix the problem
7. Re-run → Success

**Success**: User fixes error without external help

### UI Components

#### Node Palette
- **Location**: Left sidebar
- **Features**:
  - Categorized node types (collapsible)
  - Search/filter nodes
  - Node icons and labels
  - Hover for description
  - Drag to canvas

#### Canvas
- **Location**: Main area
- **Features**:
  - Infinite scrollable canvas
  - Pan and zoom
  - Node selection (single/multi)
  - Edge creation (drag from output to input)
  - Right-click context menu
  - Minimap (optional)

#### Properties Panel
- **Location**: Right sidebar
- **Features**:
  - Node-specific fields
  - Type-appropriate inputs (text, number, dropdown, etc.)
  - Validation errors inline
  - Help text for each field
  - Expand/collapse sections

#### Results Panel
- **Location**: Bottom panel (resizable)
- **Features**:
  - Execution summary (time, status)
  - Final output (formatted)
  - Per-node results (expandable)
  - Error details with stack trace
  - Execution logs (optional)

#### Navigation Bar
- **Location**: Top
- **Features**:
  - New/Open/Save workflow
  - Run/Cancel execution
  - Import/Export
  - Settings
  - Help/Documentation

---

## Success Criteria

### User Metrics

| Metric | Target (6m) | Target (12m) | Measurement |
|--------|------------|--------------|-------------|
| Monthly Active Users | 1,000 | 10,000 | Analytics |
| Workflow Creation Rate | 5,000/mo | 50,000/mo | Database |
| Execution Success Rate | >95% | >97% | Metrics |
| User Retention (7-day) | >40% | >50% | Analytics |
| NPS Score | >50 | >60 | Survey |

### Feature Adoption

| Feature | Adoption Target | Measurement |
|---------|----------------|-------------|
| HTTP Nodes | >70% of users | Usage analytics |
| Array Operations | >60% of users | Usage analytics |
| Error Handling | >40% of users | Usage analytics |
| Variables | >30% of users | Usage analytics |
| Loops | >25% of users | Usage analytics |

### Performance Benchmarks

| Metric | Target | Measurement |
|--------|--------|-------------|
| Workflow Execution | <5s (p95) | Metrics |
| API Response Time | <200ms (p95) | Metrics |
| UI Responsiveness | <100ms input lag | Synthetic monitoring |
| Uptime | >99.9% | Uptime monitoring |

### Quality Indicators

| Metric | Target | Measurement |
|--------|--------|-------------|
| Bug Reports | <10/week | GitHub issues |
| Critical Bugs | 0 | Issue tracker |
| Support Requests | <50/week | Support system |
| Documentation Gaps | <5 reported/month | Feedback |

---

## Out of Scope

The following are explicitly out of scope for version 1.0:

### Not Building

1. **Mobile Native Apps**: Web-first approach, mobile web responsive only
2. **Enterprise SSO Integration**: Users should use reverse proxy (e.g., OAuth2 Proxy)
3. **Built-in Database**: Users provide their own persistence layer
4. **Email/SMS Notifications**: Use webhook nodes to integrate with external services
5. **Built-in Scheduling**: Use external scheduler (cron, Kubernetes CronJob)
6. **Visual Workflow Versioning**: Rely on Git for version control
7. **Real-time Collaboration**: Single-user workflows initially

### Relying on Existing Solutions

1. **Authentication**: Kubernetes RBAC, OAuth2 Proxy, etc.
2. **TLS/HTTPS**: Reverse proxy (Nginx, Traefik)
3. **Monitoring**: Prometheus + Grafana
4. **Logging**: Elasticsearch + Kibana or similar
5. **Secret Management**: Kubernetes Secrets, HashiCorp Vault

---

## Dependencies & Risks

### Technical Dependencies

| Dependency | Version | Risk | Mitigation |
|-----------|---------|------|------------|
| Next.js | 16.0+ | Breaking changes in future versions | Pin version, test upgrades |
| React Flow | 11.8+ | Library maintenance | Have migration plan |
| Go | 1.24.7+ | Language changes | Use stable features only |
| Docker | 20.10+ | Container runtime changes | Multi-runtime support |

### External Dependencies

| Dependency | Risk | Mitigation |
|-----------|------|------------|
| GitHub (hosting) | Outage affects development | Mirror repositories |
| npm registry | Outage affects builds | Private registry mirror |
| Docker Hub | Rate limits | Use GitHub Container Registry |

### Project Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| **Security Vulnerability** | High | Medium | Regular audits, bug bounty program |
| **Performance Issues** | Medium | Low | Load testing, benchmarks |
| **Lack of Adoption** | High | Medium | Marketing, community building |
| **Competition** | Medium | High | Focus on differentiators (security, open source) |
| **Resource Constraints** | Medium | Medium | Prioritize ruthlessly, seek contributors |
| **Technical Debt** | Medium | Medium | Regular refactoring, high test coverage |

### Mitigation Strategies

#### Security
- Regular security audits
- Automated vulnerability scanning
- Bug bounty program
- Responsible disclosure policy

#### Performance
- Continuous benchmarking
- Load testing before releases
- Performance budgets
- Profiling and optimization

#### Adoption
- Clear value proposition
- Comprehensive documentation
- Active community engagement
- Success stories and case studies

#### Competition
- Focus on unique differentiators
- Fast iteration on feedback
- Community-driven development
- Open source transparency

---

## Appendices

### A. User Research

**Survey Results (n=100 potential users):**
- 85% struggle with workflow automation
- 70% find existing tools too expensive
- 60% concerned about vendor lock-in
- 55% need better security controls
- 45% want visual debugging tools

**Interview Insights:**
- "I don't want to learn another programming language just to automate simple tasks"
- "We can't justify $50K/year for a workflow tool"
- "Security is our top concern - we need SSRF protection"
- "We need to self-host for compliance reasons"

### B. Competitive Analysis

| Competitor | Strength | Weakness | Our Advantage |
|-----------|----------|----------|---------------|
| **Zapier** | Easy to use, many integrations | Expensive, vendor lock-in | Open source, self-hosted |
| **n8n** | Open source, self-hostable | Limited security, UI complexity | Better security, simpler UI |
| **Airflow** | Powerful, scalable | Requires coding, complex setup | Visual interface, easier |
| **Node-RED** | Visual, IoT-focused | Weak typing, security gaps | Strong typing, enterprise security |

### C. Glossary

- **Workflow**: A sequence of automated tasks
- **Node**: A single task or operation in a workflow
- **Edge**: Connection between nodes defining data flow
- **Executor**: Code that implements a node type
- **SSRF**: Server-Side Request Forgery attack
- **MAU**: Monthly Active Users
- **NPS**: Net Promoter Score

---

**Document Status**: Active
**Review Cycle**: Monthly
**Next Review**: 2026-06-28

**Feedback**: Please submit feedback via GitHub Issues or Discussions

*This PRD is a living document and will evolve based on user feedback and market conditions.*
