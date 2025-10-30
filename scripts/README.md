# Agent Invocation Scripts

This directory contains scripts to make the agent specifications in `.github/agents/` callable.

## Quick Start

### Invoke an Agent

```bash
./scripts/invoke-agent.sh <agent-name> [options]
```

### Example: Security Review

```bash
./scripts/invoke-agent.sh security-review \
  --files "backend/nodes_http.go backend/api/*.go" \
  --context "Review for OWASP Top 10 compliance" \
  --output security-review-results.md
```

### Example: Architecture Review

```bash
./scripts/invoke-agent.sh architecture \
  --files "backend/workflow.go" \
  --context "Large file needs refactoring" \
  --output architecture-review.md \
  --verbose
```

## Available Agents

| Short Name | Full Name | Specification |
|------------|-----------|---------------|
| `security-review` | Security Code Review | `.github/agents/security-code-review.md` |
| `architecture` | System Architecture | `.github/agents/system-architecture.md` |
| `observability` | Observability | `.github/agents/observability.md` |
| `multi-tenancy` | Multi-Tenancy Specialist | `.github/agents/multi-tenancy.md` |
| `testing` | Testing & QA | `.github/agents/testing-qa.md` |
| `performance` | Performance Optimization | `.github/agents/performance.md` |
| `documentation` | Documentation | `.github/agents/documentation.md` |
| `devops` | DevOps & CI/CD | `.github/agents/devops-cicd.md` |
| `ui-ux` | UI/UX Architect | `.github/agents/ui-ux-architect.md` |
| `product-manager` | Product Manager | `.github/agents/product-manager.md` |
| `marketing` | Marketing | `.github/agents/marketing.md` |

## Options

- `--files <file-list>` - Space-separated list of files to review
- `--context <description>` - Context or description of the review
- `--output <file>` - Output file for review template
- `--quick-check` - Quick check mode (future: fail on issues)
- `--verbose` - Show detailed output

## Integration Examples

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

STAGED=$(git diff --cached --name-only | grep "\.go$")
[ -z "$STAGED" ] && exit 0

./scripts/invoke-agent.sh security-review \
  --files "$STAGED" \
  --context "Pre-commit security check" \
  --quick-check
```

### CI/CD (GitHub Actions)

```yaml
- name: Security Review
  run: |
    ./scripts/invoke-agent.sh security-review \
      --files "$(git diff --name-only origin/main...HEAD)" \
      --context "PR #${{ github.event.pull_request.number }}" \
      --output security-review.md
    
    # Upload as artifact
    - uses: actions/upload-artifact@v3
      with:
        name: agent-reviews
        path: security-review.md
```

### Review Automation Script

```bash
#!/bin/bash
# scripts/review-pr.sh - Comprehensive PR review

FILES=$(git diff --name-only main...HEAD)

# Run multiple agents
./scripts/invoke-agent.sh security-review --files "$FILES" --output reviews/security.md
./scripts/invoke-agent.sh architecture --files "$FILES" --output reviews/architecture.md
./scripts/invoke-agent.sh performance --files "$FILES" --output reviews/performance.md
./scripts/invoke-agent.sh testing --files "$FILES" --output reviews/testing.md

echo "✅ All agent reviews complete. Check reviews/ directory."
```

## Agent Call Logging

All agent invocations are logged to `.github/agent-logs/agent-calls.log`:

```
2025-10-30T14:31:16+00:00 | security-review | Review HTTP node security | Review initiated
2025-10-30T14:35:22+00:00 | architecture | Refactoring review | Review initiated
```

## For More Information

See [AGENTS_CALLABLE_SYSTEM.md](../.github/agents/AGENTS_CALLABLE_SYSTEM.md) for complete documentation on the callable agent system.

## Contributing

To add a new agent invocation feature:

1. Update `invoke-agent.sh` with new functionality
2. Add examples to this README
3. Update AGENTS_CALLABLE_SYSTEM.md
4. Test with various scenarios
5. Submit PR
