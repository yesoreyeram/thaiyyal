# Non-Functional Requirements: Governance

## Governance Requirements

### GOV-1: Code Review
- **REQ-GOV-1.1**: All changes SHALL be code reviewed
- **REQ-GOV-1.2**: Security changes SHALL have security-focused review
- **REQ-GOV-1.3**: Reviews SHALL check for code quality
- **REQ-GOV-1.4**: Reviews SHALL verify test coverage

### GOV-2: Documentation
- **REQ-GOV-2.1**: Public APIs SHALL have documentation
- **REQ-GOV-2.2**: Architecture changes SHALL update docs
- **REQ-GOV-2.3**: Breaking changes SHALL be documented
- **REQ-GOV-2.4**: Examples SHALL be provided for new features

### GOV-3: Version Control
- **REQ-GOV-3.1**: Commits SHALL follow conventional commits format
- **REQ-GOV-3.2**: Branches SHALL be short-lived
- **REQ-GOV-3.3**: No force push to main branch
- **REQ-GOV-3.4**: Release tags SHALL follow semantic versioning

### GOV-4: Security
- **REQ-GOV-4.1**: Security vulnerabilities SHALL be reported privately
- **REQ-GOV-4.2**: Dependencies SHALL be regularly updated
- **REQ-GOV-4.3**: Security scans SHALL run on PRs
- **REQ-GOV-4.4**: No secrets in repository

---

**Last Updated:** 2025-11-03
**Version:** 1.0
