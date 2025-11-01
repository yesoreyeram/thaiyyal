# Expression Evaluation System - Executive Summary

**Date**: November 1, 2025  
**Status**: Design Complete - Ready for Implementation  
**Full Design**: [docs/EXPRESSION_SYSTEM_DESIGN.md](docs/EXPRESSION_SYSTEM_DESIGN.md)

---

## Overview

Comprehensive design document for implementing an advanced expression evaluation system for Thaiyyal's IF/Condition nodes. This will enable complex conditions with node references, variable access, and boolean logic while maintaining zero external dependencies.

---

## What's Being Built

### Current Limitation
```javascript
// Today: Only simple numeric comparisons
"condition": ">100"
"condition": "==5"
```

### Future Capability
```javascript
// Tomorrow: Full expression power
"condition": "{{ node.http1.output.status == 200 }}"
"condition": "{{ variables.counter > 10 && variables.enabled }}"
"condition": "{{ contains(node.text1.output, 'error') || node.retry.output.count < 3 }}"
```

---

## Key Features

✅ **Node Output References**: `{{ node.http1.output.status == 200 }}`  
✅ **Variable Access**: `{{ variables.counter > 10 }}`  
✅ **Context Values**: `{{ context.apiKey != "" }}`  
✅ **Boolean Logic**: `&&`, `||`, `!` operators  
✅ **String Functions**: `contains()`, `matches()`, `startsWith()`, `endsWith()`  
✅ **Array Operations**: `len()`, indexing `[i]`, field access `.field`  
✅ **Time Functions**: `now()` for current timestamp  
✅ **Type Safety**: Automatic type checking and coercion  
✅ **Dependency Tracking**: Auto-extract node dependencies for graph sorting

---

## Design Decisions

### 1. Syntax: Template-Based (Recommended ⭐)

**Selected**: `{{ expression }}` syntax

**Why?**
- Consistent with existing `{{ variable.name }}` interpolation
- Familiar to users of Jinja2, Handlebars, Go templates
- Clear visual boundaries for parser
- Excellent UI integration (syntax highlighting)

**Alternatives Considered**:
- JSONPath-like (`$.nodes.http1.output.status`)
- JavaScript-like (bare `nodes.http1.output.status`)
- Prefix-based (`@node.http1.output.status`)
- Custom DSL (`WHEN node.http1 EQUALS 200`)

### 2. Architecture: Zero Dependencies

**Packages**:
```
pkg/expr/          # NEW: Expression evaluation
  ├── lexer.go     # Tokenization
  ├── parser.go    # Recursive descent parser
  ├── evaluator.go # AST evaluation
  ├── compiler.go  # Expression compilation
  ├── ast.go       # Abstract syntax tree
  └── functions.go # Built-in functions
```

**Why?**
- Consistent with Thaiyyal's zero-dependency philosophy
- Full control over features and security
- No licensing concerns
- Minimal attack surface

### 3. Strategy: Compile Once, Evaluate Many

**Approach**:
- Parse expression into AST once
- Cache compiled expressions
- Evaluate AST multiple times (loops, retries)

**Performance Targets**:
- Simple comparison: 200ns
- Node reference: 400ns
- Complex expression: 800ns
- Compilation (cached): 50ns

### 4. Backward Compatibility: 100%

**Migration**:
```javascript
// OLD - Still works!
"condition": ">100"

// NEW - Enhanced syntax
"condition": "{{ value > 100 }}"

// ADVANCED - Full power
"condition": "{{ node.sensor1.output.temp > 100 }}"
```

**No breaking changes** - existing simple conditions continue to work indefinitely.

---

## Expression Grammar (EBNF)

```ebnf
expression      = logical_or ;
logical_or      = logical_and { "||" logical_and } ;
logical_and     = equality { "&&" equality } ;
equality        = comparison { ( "==" | "!=" ) comparison } ;
comparison      = additive { ( ">" | ">=" | "<" | "<=" ) additive } ;
additive        = multiplicative { ( "+" | "-" ) multiplicative } ;
multiplicative  = unary { ( "*" | "/" | "%" ) unary } ;
unary           = "!" unary | "-" unary | postfix ;
postfix         = primary { "." ID | "[" expr "]" | "(" args ")" } ;
primary         = NUMBER | STRING | BOOL | "(" expr ")" | reference ;
reference       = "node" "." ID | "variables" "." ID | "context" "." ID ;
```

---

## Built-in Functions

| Function | Description | Example |
|----------|-------------|---------|
| `len(x)` | Length of array/string | `len(items) > 0` |
| `contains(s, sub)` | Substring check | `contains(text, "error")` |
| `matches(s, regex)` | Regex matching | `matches(email, ".*@.*")` |
| `startsWith(s, prefix)` | Prefix check | `startsWith(url, "https")` |
| `endsWith(s, suffix)` | Suffix check | `endsWith(file, ".json")` |
| `now()` | Current timestamp | `now() > startTime` |
| `parseInt(s)` | String to integer | `parseInt("42")` |
| `toString(v)` | Value to string | `toString(123)` |
| `lower(s)` | Lowercase | `lower(text)` |
| `upper(s)` | Uppercase | `upper(text)` |

---

## Implementation Phases

### Phase 1: Core Infrastructure (Week 1-2)
- Lexer, parser, basic evaluator
- **Deliverable**: Simple expressions work (`10 > 5`, `true && false`)

### Phase 2: Reference Resolution (Week 3)
- Node, variable, context references
- **Deliverable**: `{{ node.http1.output.status == 200 }}` works

### Phase 3: Functions (Week 4)
- Built-in functions implementation
- **Deliverable**: All functions working (`len()`, `contains()`, etc.)

### Phase 4: Integration (Week 5)
- Update condition/switch/while nodes
- Dependency extraction in engine
- **Deliverable**: Full integration, backward compatibility

### Phase 5: Polish (Week 6)
- Documentation, UI improvements
- **Deliverable**: Production ready

**Total Timeline**: 6 weeks

---

## Security Measures

### 1. Injection Prevention
- ✅ No `eval()`, `exec()`, or system calls
- ✅ Whitelist of safe functions only
- ✅ No arbitrary code execution

### 2. Resource Limits
- ✅ Max recursion depth (default: 100)
- ✅ Evaluation timeout (default: 1s)
- ✅ Expression length limit (10,000 chars)

### 3. Type Safety
- ✅ Strict type checking during evaluation
- ✅ Safe type coercion rules
- ✅ Clear error messages

### 4. ReDoS Protection
- ✅ Regex timeout (100ms)
- ✅ Pattern validation

---

## Example Expressions

### HTTP Status Check
```javascript
{{ node.http1.output.status == 200 }}
```

### Multi-Condition Logic
```javascript
{{ (node.http1.output.status == 200 || node.http1.output.status == 201) && 
   !node.http1.output.error }}
```

### String Validation
```javascript
{{ len(node.form1.output.email) > 0 &&
   matches(node.form1.output.email, "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$") }}
```

### Sensor Threshold with Alerts
```javascript
{{ (node.sensor1.output.temperature > variables.threshold ||
    node.sensor2.output.temperature > variables.threshold) &&
   variables.alertsEnabled == true }}
```

### API Health Check
```javascript
{{ node.http1.output.status == 200 && 
   node.http1.output.responseTime < 500 &&
   contains(node.http1.output.body, "healthy") }}
```

---

## Architecture Integration

### Dependency Extraction
```go
// Automatic dependency graph building
expression := "{{ node.http1.output.status == 200 }}"
compiled := expr.Compile(expression)
dependencies := compiled.Dependencies()  // ["http1"]

// Add implicit edges
edge := Edge{
    Source: "http1",
    Target: "condition1",
}
```

### Condition Node Usage
```go
// In ConditionExecutor.Execute()
if isTemplateExpression(condition) {
    compiled := expr.Compile(extractExpression(condition))
    result := compiled.Evaluate(ctx)
    conditionMet := result.(bool)
} else {
    // Legacy: evaluateCondition(">100", value)
}
```

---

## Testing Strategy

### Coverage Targets
- Unit tests: ≥95% coverage
- Integration tests: All node types
- Benchmark tests: Performance validation
- Edge case tests: Error handling

### Test Categories
1. **Lexer**: Tokenization accuracy
2. **Parser**: AST correctness
3. **Evaluator**: Expression results
4. **Integration**: End-to-end workflows
5. **Performance**: Benchmark compliance
6. **Security**: Injection attempts, resource exhaustion

---

## Migration Path

### For Users
```javascript
// Step 1: No action required
// Existing workflows continue working

// Step 2: Try new syntax (optional)
"condition": "{{ value > 100 }}"

// Step 3: Adopt advanced features (when needed)
"condition": "{{ node.http1.output.status == 200 }}"
```

### For Developers
1. **No breaking changes**: Legacy `evaluateCondition()` preserved
2. **Gradual adoption**: New features opt-in
3. **Clear documentation**: Migration examples provided
4. **Testing**: Comprehensive backward compatibility tests

---

## Success Criteria

✅ **Functional**:
- All expression types working
- Node references functional
- Built-in functions operational
- Type coercion working

✅ **Performance**:
- Benchmarks meet targets (<1µs per evaluation)
- Cache hit rate >90%
- No performance regression

✅ **Security**:
- Zero code injection vulnerabilities
- Resource limits enforced
- Security audit passed

✅ **Quality**:
- Test coverage ≥95%
- Documentation complete
- Backward compatibility 100%

✅ **Usability**:
- Clear error messages
- UI syntax highlighting
- Interactive examples
- User documentation

---

## Next Steps

1. ✅ **Design Review** - Complete (this document)
2. ⏳ **Approval** - Awaiting stakeholder sign-off
3. ⏳ **Phase 1 Start** - Core infrastructure implementation
4. ⏳ **Sprint Planning** - Break down into tasks
5. ⏳ **Resource Allocation** - Assign developers

---

## Questions for Stakeholders

1. **Syntax Approval**: Confirm template-based syntax `{{ }}` is preferred?
2. **Built-in Functions**: Any additional functions needed?
3. **Security Requirements**: Any additional security concerns?
4. **Timeline**: Is 6-week timeline acceptable?
5. **Resources**: How many developers can be allocated?

---

## References

- **Full Design**: [docs/EXPRESSION_SYSTEM_DESIGN.md](docs/EXPRESSION_SYSTEM_DESIGN.md) (1773 lines)
- **Current Codebase**: `backend/pkg/executor/helpers.go:evaluateCondition()`
- **Architecture**: `backend/pkg/` (engine, executor, state, graph packages)
- **Agent Specs**: `.github/agents/system-architecture.md`, `security-code-review.md`

---

## Contact

For questions or feedback on this design:
- Review full document: `docs/EXPRESSION_SYSTEM_DESIGN.md`
- GitHub Issues: Tag with `enhancement`, `expression-system`
- Architecture Team: For technical questions

---

**Status**: ✅ Design Complete - Ready for Implementation Review
