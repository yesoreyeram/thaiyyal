# Agent Quick Reference Guide

This is a quick reference for using the Thaiyyal AI Agent system. For complete documentation, see [AGENTS.md](AGENTS.md).

## When to Use Which Agent?

### 🔒 Security Concerns
**Use: [Security Code Review Agent](.github/agents/security-code-review.md)**
- "Is this code secure?"
- "How do I prevent SQL injection?"
- "What security tests should I add?"
- "How do I implement authentication?"

### 🏗️ Architecture Questions
**Use: [System Architecture Agent](.github/agents/system-architecture.md)**
- "How should I structure this feature?"
- "What design pattern should I use?"
- "How do I make this scalable?"
- "Should I refactor this code?"

### 📊 Monitoring & Logging
**Use: [Observability Agent](.github/agents/observability.md)**
- "How do I add metrics?"
- "What should I log?"
- "How do I trace this workflow?"
- "How do I monitor tenant usage?"

### 👥 Multi-Tenant Features
**Use: [Multi-Tenancy Specialist](.github/agents/multi-tenancy.md)**
- "How do I isolate tenant data?"
- "How do I implement quotas?"
- "How do I add tenant-specific features?"
- "How do I design the tenant database schema?"

### ✅ Testing Strategy
**Use: [Testing & QA Agent](.github/agents/testing-qa.md)**
- "What tests should I write?"
- "How do I test multi-tenant isolation?"
- "How do I improve test coverage?"
- "How do I write E2E tests?"

### ⚡ Performance Issues
**Use: [Performance Agent](.github/agents/performance.md)**
- "Why is this slow?"
- "How do I optimize this query?"
- "How do I reduce memory usage?"
- "How do I profile my code?"

### 📝 Documentation Needs
**Use: [Documentation Agent](.github/agents/documentation.md)**
- "How do I document this API?"
- "How do I write user guides?"
- "How do I create architecture diagrams?"
- "How do I generate API docs?"

### 🚀 Deployment & CI/CD
**Use: [DevOps Agent](.github/agents/devops-cicd.md)**
- "How do I deploy this locally?"
- "How do I set up CI/CD?"
- "How do I containerize this?"
- "How do I deploy to production?"

## Common Workflows

### Adding a New Feature (End-to-End)

1. **Design Phase** → System Architecture Agent
   - Design the feature architecture
   - Choose appropriate design patterns
   - Plan for scalability

2. **Security Review** → Security Code Review Agent
   - Review security implications
   - Design secure APIs
   - Plan authentication/authorization

3. **Multi-Tenant Design** → Multi-Tenancy Specialist
   - Ensure tenant isolation
   - Design quota management
   - Plan data segregation

4. **Implementation** → (You write the code)

5. **Performance Review** → Performance Agent
   - Optimize database queries
   - Add caching where needed
   - Profile critical paths

6. **Testing** → Testing & QA Agent
   - Write unit tests
   - Add integration tests
   - Create E2E tests

7. **Observability** → Observability Agent
   - Add metrics
   - Implement logging
   - Add tracing

8. **Documentation** → Documentation Agent
   - Document API
   - Update user guide
   - Create examples

9. **Deployment** → DevOps Agent
   - Update CI/CD pipeline
   - Deploy to staging
   - Deploy to production

### Code Review (Comprehensive)

Use these agents in parallel:
- **Security Agent**: Check for vulnerabilities
- **Performance Agent**: Review for bottlenecks
- **Testing Agent**: Verify test coverage
- **Documentation Agent**: Check documentation

### Fixing a Production Issue

1. **Observability Agent** → Analyze metrics and logs
2. **Performance Agent** → Profile and identify bottleneck
3. **Security Agent** → Check if it's a security issue
4. **System Architecture Agent** → Design fix
5. **Testing Agent** → Write regression tests
6. **DevOps Agent** → Deploy hotfix

## Quick Tips

### ✅ Do's
- Consult agents **before** implementing complex features
- Use multiple agents for comprehensive reviews
- Follow the agent collaboration patterns
- Validate agent recommendations against project requirements

### ❌ Don'ts
- Don't implement security-critical code without Security Agent review
- Don't deploy without DevOps Agent review
- Don't skip Testing Agent for new features
- Don't ignore Performance Agent recommendations

## Local Development Setup

All agents assume you can run Thaiyyal locally:

```bash
# Quick start (works offline)
git clone https://github.com/yesoreyeram/thaiyyal.git
cd thaiyyal
npm install
npm run dev
```

No cloud dependencies required!

## Agent Outputs

Each agent provides:
1. **Analysis** of the current situation
2. **Recommendations** with specific actions
3. **Code Examples** showing best practices
4. **Testing Strategies** to validate changes
5. **Documentation** updates needed

## Getting Help

1. **Start with AGENTS.md** - Overview and index
2. **Check specific agent file** - Detailed guidance
3. **Look for examples** - Each agent has code samples
4. **Follow collaboration patterns** - Multi-agent workflows

## Enterprise Standards

All agents enforce:
- ✅ 80%+ test coverage
- ✅ Security best practices
- ✅ Performance optimization
- ✅ Complete documentation
- ✅ Multi-tenant isolation
- ✅ Local-first architecture
- ✅ Production-ready code

---

**Quick Access:**
- [Main Documentation](AGENTS.md)
- [All Agent Files](.github/agents/)
- [Project README](README.md)
- [Architecture Docs](ARCHITECTURE.md)
