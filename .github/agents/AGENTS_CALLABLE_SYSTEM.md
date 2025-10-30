# Callable Agent System for Thaiyyal

**Version**: 1.0  
**Created**: October 30, 2025  
**Status**: Active

---

## Overview

This document describes how to "call" the agent specifications in `.github/agents/` as if they were callable tools. While these are specification documents (not system-level callable tools), this framework provides a structured way to invoke and apply them programmatically.

---

## How to Call an Agent

### Agent Invocation Syntax

To "call" an agent, use this structured approach:

```markdown
@agent:<agent-name> <task-description>

Context:
- <relevant context item 1>
- <relevant context item 2>

Files:
- <file path 1>
- <file path 2>

Expected Output:
- <what you expect the agent to produce>
```

### Available Agents

| Agent Name | Invocation | Purpose |
|------------|------------|---------|
| Security Code Review | `@agent:security-review` | Security vulnerability analysis |
| System Architecture | `@agent:architecture` | Architecture design and review |
| Observability | `@agent:observability` | Monitoring, logging, tracing |
| Multi-Tenancy | `@agent:multi-tenancy` | Multi-tenant architecture |
| Testing & QA | `@agent:testing` | Testing strategies and QA |
| Performance | `@agent:performance` | Performance optimization |
| Documentation | `@agent:documentation` | Technical documentation |
| DevOps & CI/CD | `@agent:devops` | CI/CD and infrastructure |

---

## Agent Call Examples

### Example 1: Security Review

```markdown
@agent:security-review Review HTTP node for vulnerabilities

Context:
- Backend Go application
- HTTP node in backend/nodes_http.go
- Currently accepts any URL without validation

Files:
- backend/nodes_http.go

Expected Output:
- List of security vulnerabilities
- Remediation recommendations with code examples
- OWASP Top 10 compliance check
```

### Example 2: Architecture Review

```markdown
@agent:architecture Review workflow.go organization

Context:
- Backend workflow engine
- Current file is 1,173 lines
- Monolithic structure

Files:
- backend/workflow.go

Expected Output:
- Code organization issues
- Recommended package structure
- Refactoring strategy
```

### Example 3: Multi-Agent Call (Sequential)

```markdown
@agent:architecture -> @agent:security-review -> @agent:performance

Task: Review new API endpoint implementation

Context:
- REST API for workflow execution
- PostgreSQL database backend
- Multi-tenant system

Files:
- backend/api/handlers/workflow.go
- backend/repository/workflow_repo.go

Expected Output:
- Architecture assessment
- Security vulnerabilities
- Performance bottlenecks
- Consolidated recommendations
```

---

## Implementation Guide

### For Human Reviewers

When you see an agent call like `@agent:security-review`:

1. **Open the specification**: `.github/agents/security-code-review.md`
2. **Read the relevant sections**: Focus on the specific task domain
3. **Apply the checklists**: Use the specification's guidelines
4. **Document findings**: Follow the specification's output format
5. **Provide recommendations**: Based on the spec's best practices

### For AI Assistants

When processing an agent call:

1. **Parse the invocation**: Extract agent name, task, context, files
2. **Load the specification**: Read `.github/agents/<agent-name>.md`
3. **Apply the guidelines**: Use the spec's checklists and patterns
4. **Generate output**: Follow the spec's recommended format
5. **Cross-reference**: Link to other specs when needed

### For Automated Systems

To integrate agent calls in automation:

```bash
#!/bin/bash
# Example: Call security agent on changed files

AGENT="security-review"
FILES=$(git diff --name-only main...HEAD | grep "\.go$")

# Generate agent call
cat > agent-call.md << EOF
@agent:${AGENT} Review security of changed files

Context:
- Pull request review
- Go backend files
- Production deployment target

Files:
${FILES}

Expected Output:
- Security vulnerabilities (if any)
- OWASP Top 10 compliance
- Remediation steps
EOF

# Process with agent specification
./scripts/process-agent-call.sh agent-call.md
```

---

## Agent Response Format

Each agent should respond in this standardized format:

```markdown
# Agent Response: <Agent Name>

**Task**: <Original task description>
**Date**: <ISO 8601 timestamp>
**Status**: ✅ Complete | ⚠️ Issues Found | ❌ Failed

## Summary
<2-3 sentence overview of findings>

## Findings

### Issue 1: <Issue Title>
**Severity**: Critical | High | Medium | Low
**Location**: <file:line>
**Description**: <What's wrong>
**Impact**: <What could happen>
**Recommendation**: <How to fix>

```code
<Example fix code>
```

### Issue 2: ...

## Recommendations
1. <Action item 1>
2. <Action item 2>

## References
- <Link to spec section>
- <Related documentation>

## Metrics
- Issues Found: X
- Critical: X
- High: X
- Medium: X
- Low: X
```

---

## Agent Chaining

### Sequential Chaining

Apply agents one after another:

```markdown
@agent:architecture
  ↓ (output becomes context)
@agent:security-review
  ↓ (output becomes context)
@agent:performance
```

**Implementation**:
```bash
# Step 1: Architecture review
architecture_output=$(invoke-agent architecture "$files")

# Step 2: Security review (with architecture context)
security_output=$(invoke-agent security-review "$files" --context "$architecture_output")

# Step 3: Performance review (with both contexts)
performance_output=$(invoke-agent performance "$files" \
  --context "$architecture_output" \
  --context "$security_output")
```

### Parallel Chaining

Apply multiple agents simultaneously:

```markdown
@agent:security-review ──┐
@agent:performance    ────┤──→ Merge Results
@agent:testing        ────┘
```

**Implementation**:
```bash
# Run in parallel
invoke-agent security-review "$files" &
PID1=$!

invoke-agent performance "$files" &
PID2=$!

invoke-agent testing "$files" &
PID3=$!

# Wait for all to complete
wait $PID1 $PID2 $PID3

# Merge results
merge-agent-outputs security-review.md performance.md testing.md > final-review.md
```

---

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Agent Review

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  agent-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Security Agent Review
        run: |
          ./scripts/invoke-agent.sh security-review \
            --files "$(git diff --name-only origin/main...HEAD)" \
            --context "Pull Request #${{ github.event.pull_request.number }}" \
            --output security-review.md
      
      - name: Architecture Agent Review
        run: |
          ./scripts/invoke-agent.sh architecture \
            --files "$(git diff --name-only origin/main...HEAD)" \
            --context "Pull Request #${{ github.event.pull_request.number }}" \
            --output architecture-review.md
      
      - name: Post Results as PR Comment
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const security = fs.readFileSync('security-review.md', 'utf8');
            const architecture = fs.readFileSync('architecture-review.md', 'utf8');
            
            await github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: `## 🤖 Agent Review\n\n${security}\n\n${architecture}`
            });
```

### Pre-commit Hook Example

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running agent reviews..."

# Get staged files
STAGED_FILES=$(git diff --cached --name-only)

# Security review
if echo "$STAGED_FILES" | grep -q "\.go$"; then
  echo "→ Security review of Go files..."
  ./scripts/invoke-agent.sh security-review \
    --files "$STAGED_FILES" \
    --quick-check
  
  if [ $? -ne 0 ]; then
    echo "❌ Security issues found. Fix before committing."
    exit 1
  fi
fi

# Architecture review for large changes
if [ $(echo "$STAGED_FILES" | wc -l) -gt 10 ]; then
  echo "→ Architecture review (large changeset)..."
  ./scripts/invoke-agent.sh architecture \
    --files "$STAGED_FILES" \
    --quick-check
fi

echo "✅ Agent reviews passed"
```

---

## Agent Invocation Script

Create `scripts/invoke-agent.sh`:

```bash
#!/bin/bash
# scripts/invoke-agent.sh
# Invokes an agent specification against files

set -e

AGENT_NAME=$1
shift

# Parse arguments
FILES=""
CONTEXT=""
OUTPUT=""
QUICK_CHECK=false

while [[ $# -gt 0 ]]; do
  case $1 in
    --files)
      FILES="$2"
      shift 2
      ;;
    --context)
      CONTEXT="$2"
      shift 2
      ;;
    --output)
      OUTPUT="$2"
      shift 2
      ;;
    --quick-check)
      QUICK_CHECK=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Map agent name to spec file
case $AGENT_NAME in
  security-review|security)
    SPEC_FILE=".github/agents/security-code-review.md"
    ;;
  architecture|arch)
    SPEC_FILE=".github/agents/system-architecture.md"
    ;;
  observability|obs)
    SPEC_FILE=".github/agents/observability.md"
    ;;
  multi-tenancy|mt)
    SPEC_FILE=".github/agents/multi-tenancy.md"
    ;;
  testing|qa)
    SPEC_FILE=".github/agents/testing-qa.md"
    ;;
  performance|perf)
    SPEC_FILE=".github/agents/performance.md"
    ;;
  documentation|docs)
    SPEC_FILE=".github/agents/documentation.md"
    ;;
  devops|cicd)
    SPEC_FILE=".github/agents/devops-cicd.md"
    ;;
  *)
    echo "Unknown agent: $AGENT_NAME"
    echo "Available agents: security-review, architecture, observability, multi-tenancy, testing, performance, documentation, devops"
    exit 1
    ;;
esac

# Check if spec file exists
if [ ! -f "$SPEC_FILE" ]; then
  echo "Agent specification not found: $SPEC_FILE"
  exit 1
fi

# Create agent call document
CALL_FILE=$(mktemp)
cat > "$CALL_FILE" << EOF
@agent:${AGENT_NAME}

Context:
${CONTEXT}

Files:
${FILES}

Quick Check: ${QUICK_CHECK}
EOF

echo "📋 Agent Call: $AGENT_NAME"
echo "📄 Specification: $SPEC_FILE"
echo "📁 Files: $(echo "$FILES" | wc -w) files"
echo ""

# TODO: Actual implementation would:
# 1. Parse the specification file
# 2. Extract relevant checklists and guidelines
# 3. Apply them to the specified files
# 4. Generate findings report

# For now, display what would be checked
echo "Would apply checklist from: $SPEC_FILE"
echo "To files: $FILES"

# If output specified, create placeholder
if [ -n "$OUTPUT" ]; then
  cat > "$OUTPUT" << EOF
# Agent Response: ${AGENT_NAME}

**Date**: $(date -Iseconds)
**Status**: ⚠️ Manual Review Required

## Summary
Agent specification loaded from $SPEC_FILE
Manual review required to apply guidelines to specified files.

## Next Steps
1. Review specification: $SPEC_FILE
2. Apply checklists to files: $FILES
3. Document findings in this file

## Files Reviewed
${FILES}
EOF
  echo "📝 Output written to: $OUTPUT"
fi

rm "$CALL_FILE"

exit 0
```

Make it executable:
```bash
chmod +x scripts/invoke-agent.sh
```

---

## Agent Call Tracking

### Track Agent Calls in Issues

When using agents in issue/PR reviews, track them:

```markdown
## Agent Reviews Applied

- [x] @agent:security-review - Completed - [Results](./reviews/security-001.md)
- [x] @agent:architecture - Completed - [Results](./reviews/architecture-001.md)
- [ ] @agent:performance - In Progress
- [ ] @agent:testing - Pending

### Findings Summary
- **Critical**: 2 issues (security)
- **High**: 5 issues (architecture, performance)
- **Medium**: 8 issues (various)
- **Low**: 12 issues (documentation, style)
```

### Agent Call History

Keep a log in `.github/agent-calls.log`:

```
2025-10-30T14:30:00Z | security-review | PR#123 | backend/nodes_http.go | 3 issues found
2025-10-30T14:35:00Z | architecture | PR#123 | backend/workflow.go | 7 issues found
2025-10-30T15:00:00Z | performance | PR#124 | backend/executor.go | 2 issues found
```

---

## Best Practices

### 1. Always Provide Context
Bad:
```markdown
@agent:security-review backend/api/
```

Good:
```markdown
@agent:security-review Review authentication implementation

Context:
- New JWT-based auth system
- Used in multi-tenant environment
- Production deployment planned

Files:
- backend/api/auth/jwt.go
- backend/api/middleware/auth.go
```

### 2. Specify Expected Output
```markdown
Expected Output:
- OWASP Top 10 compliance check
- JWT token security analysis
- Session management review
- Input validation assessment
```

### 3. Chain Agents Appropriately
For new features: Architecture → Security → Performance → Testing
For bug fixes: Testing → Security → Documentation
For refactoring: Architecture → Performance → Testing

### 4. Document Agent Decisions
When an agent recommends changes, document why:
```markdown
## Agent Recommendation Adopted

**Agent**: @agent:security-review
**Recommendation**: Add CSRF protection to API endpoints
**Rationale**: API is used from browser, vulnerable to CSRF
**Implementation**: Added CSRF middleware (commit: abc123)
**Verification**: Unit tests added, penetration test passed
```

---

## Troubleshooting

### Agent Call Not Working?

1. **Check specification exists**: Verify `.github/agents/<agent>.md` exists
2. **Verify syntax**: Use correct `@agent:<name>` format
3. **Provide context**: Always include context and files
4. **Check file paths**: Ensure files are relative to repo root

### No Output Generated?

1. **Manual review required**: Some agents need human expertise
2. **Increase verbosity**: Add `--verbose` flag
3. **Check logs**: Review agent execution logs
4. **Verify spec format**: Ensure spec file is well-formed

---

## Future Enhancements

### Planned Features

1. **AI-Powered Agent Execution**: Use LLM to automatically apply specs
2. **Agent Result Caching**: Cache results for unchanged files
3. **Custom Agent Creation**: Allow teams to create custom agent specs
4. **Agent Metrics Dashboard**: Track agent usage and findings
5. **Agent Learning**: Improve specs based on findings patterns

### Integration Roadmap

- **Month 1**: Basic invocation script and CI/CD integration
- **Month 2**: Agent chaining and result merging
- **Month 3**: AI-powered automatic application
- **Month 4**: Metrics dashboard and reporting
- **Month 5**: Custom agent framework
- **Month 6**: Full automation and learning

---

## Contributing

To add a new callable agent:

1. Create specification in `.github/agents/<agent-name>.md`
2. Add agent mapping to `invoke-agent.sh`
3. Add to "Available Agents" table above
4. Add example agent call
5. Test with sample files
6. Update documentation

---

**Version History**:
- v1.0 (2025-10-30): Initial callable agent system

**Maintained By**: Thaiyyal Development Team  
**Questions**: Open an issue with label `agent-system`
