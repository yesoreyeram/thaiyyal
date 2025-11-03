# Principles: Modular Design

## Overview

Thaiyyal uses modular design with clear package boundaries and minimal dependencies.

## Package Organization

```
backend/pkg/
├── types/          # Foundation - no dependencies
├── config/         # Configuration
├── logging/        # Structured logging
├── security/       # Security utilities
├── expression/     # Expression evaluation
├── graph/          # Graph algorithms
├── state/          # State management
├── observer/       # Observer pattern
├── httpclient/     # HTTP client
├── middleware/     # Execution middleware
├── executor/       # Node executors
└── engine/         # Orchestration
```

## Dependency Rules

1. **No Circular Dependencies**: Packages form a DAG
2. **Minimal Dependencies**: Each package depends only on what it needs
3. **Clear Interfaces**: Well-defined package APIs
4. **Standard Library Preferred**: Minimal external dependencies

## Module Benefits

- **Independent Testing**: Test packages in isolation
- **Easy Understanding**: Clear responsibilities
- **Maintainability**: Changes localized to packages
- **Reusability**: Packages can be used independently

---

**Last Updated:** 2025-11-03
**Version:** 1.0
