#!/bin/bash
# scripts/invoke-agent.sh
# Invokes an agent specification against files

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Parse arguments
AGENT_NAME=$1
shift

FILES=""
CONTEXT=""
OUTPUT=""
QUICK_CHECK=false
VERBOSE=false

while [[ $# -gt 0 ]]; do
  case $1 in
    --files) FILES="$2"; shift 2 ;;
    --context) CONTEXT="$2"; shift 2 ;;
    --output) OUTPUT="$2"; shift 2 ;;
    --quick-check) QUICK_CHECK=true; shift ;;
    --verbose) VERBOSE=true; shift ;;
    *) echo -e "${RED}Unknown option: $1${NC}"; exit 1 ;;
  esac
done

# Map agent to spec file
case $AGENT_NAME in
  security-review|security) SPEC_FILE=".github/agents/security-code-review.md"; AGENT_DISPLAY="Security Code Review" ;;
  architecture|arch) SPEC_FILE=".github/agents/system-architecture.md"; AGENT_DISPLAY="System Architecture" ;;
  observability|obs) SPEC_FILE=".github/agents/observability.md"; AGENT_DISPLAY="Observability" ;;
  multi-tenancy|mt) SPEC_FILE=".github/agents/multi-tenancy.md"; AGENT_DISPLAY="Multi-Tenancy" ;;
  testing|qa) SPEC_FILE=".github/agents/testing-qa.md"; AGENT_DISPLAY="Testing & QA" ;;
  performance|perf) SPEC_FILE=".github/agents/performance.md"; AGENT_DISPLAY="Performance" ;;
  documentation|docs) SPEC_FILE=".github/agents/documentation.md"; AGENT_DISPLAY="Documentation" ;;
  devops|cicd) SPEC_FILE=".github/agents/devops-cicd.md"; AGENT_DISPLAY="DevOps & CI/CD" ;;
  ui-ux|ux) SPEC_FILE=".github/agents/ui-ux-architect.md"; AGENT_DISPLAY="UI/UX Architect" ;;
  product-manager|pm) SPEC_FILE=".github/agents/product-manager.md"; AGENT_DISPLAY="Product Manager" ;;
  marketing|market) SPEC_FILE=".github/agents/marketing.md"; AGENT_DISPLAY="Marketing" ;;
  *)
    echo -e "${RED}Unknown agent: $AGENT_NAME${NC}"
    echo "Available: security-review, architecture, observability, multi-tenancy, testing, performance, documentation, devops, ui-ux, product-manager, marketing"
    exit 1
    ;;
esac

[ ! -f "$SPEC_FILE" ] && echo -e "${RED}Spec not found: $SPEC_FILE${NC}" && exit 1

# Print header
echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}  🤖 Agent: ${GREEN}${AGENT_DISPLAY}${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}📄 Spec:${NC} $SPEC_FILE"
[ -n "$FILES" ] && echo -e "${YELLOW}📁 Files:${NC} $(echo "$FILES" | wc -w) file(s)"
[ -n "$CONTEXT" ] && echo -e "${YELLOW}📝 Context:${NC} $CONTEXT"
echo ""

# Generate output if specified
if [ -n "$OUTPUT" ]; then
  cat > "$OUTPUT" << EOF
# Agent Review: ${AGENT_DISPLAY}

**Date**: $(date -Iseconds)
**Specification**: $SPEC_FILE

## Files
${FILES:-All relevant files}

## Context
${CONTEXT:-No context provided}

## Next Steps
1. Review specification: $SPEC_FILE
2. Apply guidelines to files
3. Document findings here
EOF
  echo -e "${GREEN}✅ Output created: $OUTPUT${NC}"
fi

# Log
mkdir -p .github/agent-logs
echo "$(date -Iseconds) | $AGENT_NAME | ${CONTEXT:-no-context} | Review initiated" >> .github/agent-logs/agent-calls.log

echo -e "${GREEN}✅ Agent invocation complete${NC}"
exit 0
