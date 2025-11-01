#!/bin/bash
# Script to create critical GitHub issues for Thaiyyal project
# Usage: GH_TOKEN=your_token ./create-github-issues.sh

set -e

REPO="yesoreyeram/thaiyyal"

echo "Creating critical GitHub issues for Thaiyyal..."
echo "Repository: $REPO"
echo ""

# Check if gh CLI is available
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed."
    echo "Install from: https://cli.github.com/"
    exit 1
fi

# Verify authentication
if ! gh auth status &> /dev/null; then
    echo "Error: Not authenticated with GitHub."
    echo "Run: gh auth login"
    exit 1
fi

echo "✓ GitHub CLI authenticated"
echo ""

# Issue 1: Distributed Workflow Execution Engine
echo "Creating Issue 1: Distributed Workflow Execution Engine..."
gh issue create \
  --repo "$REPO" \
  --title "[EPIC] Distributed Workflow Execution Engine" \
  --label "epic,enhancement,priority:high,complexity:very-high,area:backend,area:infrastructure" \
  --body - <<'EOF'
## Overview

Implement a distributed workflow execution engine that enables horizontal scaling, fault tolerance, and resource isolation across multiple worker nodes.

## Problem Statement

**Current State**: Workflows execute in a single process with in-memory state, limiting scalability and fault tolerance.

**Desired State**: Distributed execution across multiple workers with:
- Horizontal scaling for high-load scenarios
- Fault tolerance for long-running workflows
- Resource isolation between workflows
- Automatic failover and recovery

## Business Value

- **Scalability**: Handle 10x more concurrent workflow executions
- **Reliability**: 99.9% uptime with automatic recovery
- **Performance**: Reduce execution time through parallel processing
- **Cost**: Better resource utilization and cost efficiency

## Acceptance Criteria

### Phase 1: Foundation (Sprint 1-2)
- [ ] Design distributed architecture document
- [ ] Select task queue technology (Redis/RabbitMQ/NATS/Kafka)
- [ ] Select state store technology (Redis/etcd/PostgreSQL)
- [ ] Create ADR (Architecture Decision Record)
- [ ] Define communication protocols
- [ ] Design failure scenarios and recovery strategies

### Phase 2: Core Implementation (Sprint 3-5)
- [ ] Implement coordinator service
- [ ] Implement worker service  
- [ ] Implement task queue integration
- [ ] Implement distributed state management

### Phase 3: Fault Tolerance (Sprint 6-7)
- [ ] Implement failure detection
- [ ] Implement recovery mechanisms
- [ ] Implement circuit breakers
- [ ] Add backpressure mechanisms

### Phase 4: Observability (Sprint 8)
- [ ] Add comprehensive logging with distributed tracing
- [ ] Add metrics collection
- [ ] Create monitoring dashboards

### Phase 5: Testing & Optimization (Sprint 9-10)
- [ ] Write unit tests (>80% coverage)
- [ ] Write integration tests
- [ ] Chaos engineering tests
- [ ] Load testing
- [ ] Complete documentation

## Timeline

**Estimated Effort**: 20-30 person-days  
**Duration**: 10 weeks
EOF

echo "✓ Issue 1 created"
sleep 2

# Continue with remaining issues...
echo "Script completed. Run the full version to create all 12 issues."
