# Expression Evaluation System - Implementation Tasks

**Project**: Advanced Expression Evaluation for Thaiyyal  
**Status**: Ready to Start  
**Timeline**: 6 weeks (30 working days)  
**Full Design**: [docs/EXPRESSION_SYSTEM_DESIGN.md](docs/EXPRESSION_SYSTEM_DESIGN.md)

---

## Phase 1: Core Infrastructure (10 days)

### Week 1-2: Lexer, Parser, Basic Evaluator

#### Task 1.1: Project Setup (Day 1)
- [ ] Create `backend/pkg/expr/` directory structure
- [ ] Set up package documentation
- [ ] Create test file structure
- [ ] Add package to module imports
- [ ] Create benchmark file structure

**Files**:
- `backend/pkg/expr/expr.go`
- `backend/pkg/expr/expr_test.go`
- `backend/pkg/expr/expr_bench_test.go`
- `backend/pkg/expr/examples_test.go`

---

#### Task 1.2: Token Types and Lexer (Days 2-3)
- [ ] Define `TokenType` enum (30+ types)
- [ ] Define `Token` struct (Type, Literal, Line, Column)
- [ ] Implement `Lexer` struct
- [ ] Implement `NewLexer(input string) *Lexer`
- [ ] Implement `NextToken() Token`
- [ ] Handle all operators: `&&`, `||`, `!`, `==`, `!=`, `<`, `>`, `<=`, `>=`
- [ ] Handle delimiters: `(`, `)`, `[`, `]`, `.`, `,`
- [ ] Handle literals: numbers, strings (single/double quotes), booleans
- [ ] Handle identifiers and keywords
- [ ] Implement line/column tracking for error reporting
- [ ] Add comprehensive unit tests (50+ test cases)

**Files**:
- `backend/pkg/expr/lexer.go`
- `backend/pkg/expr/token.go`
- `backend/pkg/expr/lexer_test.go`

**Test Cases**:
```go
TestLexer/simple_number
TestLexer/operators
TestLexer/identifiers
TestLexer/strings
TestLexer/delimiters
TestLexer/whitespace
TestLexer/line_tracking
TestLexer/error_cases
```

---

#### Task 1.3: AST Node Definitions (Day 4)
- [ ] Define `Node` interface (marker method)
- [ ] Implement `BinaryExpr` (left, operator, right)
- [ ] Implement `UnaryExpr` (operator, operand)
- [ ] Implement `LiteralExpr` (value interface{})
- [ ] Implement `IdentifierExpr` (name)
- [ ] Implement `FieldAccessExpr` (object, field)
- [ ] Implement `IndexAccessExpr` (object, index)
- [ ] Implement `FunctionCallExpr` (function name, args)
- [ ] Implement `ParenExpr` (expr)
- [ ] Add AST node documentation
- [ ] Add AST visualization helpers (for debugging)

**Files**:
- `backend/pkg/expr/ast.go`
- `backend/pkg/expr/ast_test.go`

---

#### Task 1.4: Recursive Descent Parser (Days 5-7)
- [ ] Define `Parser` struct (lexer, current, peek)
- [ ] Implement `NewParser(lexer *Lexer) *Parser`
- [ ] Implement `Parse() (Node, error)`
- [ ] Implement precedence climbing:
  - `parseExpression()` - entry point
  - `parseLogicalOr()` - `||` operator
  - `parseLogicalAnd()` - `&&` operator
  - `parseEquality()` - `==`, `!=`
  - `parseComparison()` - `<`, `>`, `<=`, `>=`
  - `parseAdditive()` - `+`, `-`
  - `parseMultiplicative()` - `*`, `/`, `%`
  - `parseUnary()` - `!`, `-`
  - `parsePostfix()` - `.`, `[]`, `()`
  - `parsePrimary()` - literals, identifiers, `(expr)`
- [ ] Implement error recovery
- [ ] Add detailed error messages with line/column
- [ ] Add comprehensive unit tests (100+ test cases)

**Files**:
- `backend/pkg/expr/parser.go`
- `backend/pkg/expr/parser_test.go`

**Test Cases**:
```go
TestParser/literals
TestParser/binary_operators
TestParser/unary_operators
TestParser/precedence
TestParser/parentheses
TestParser/field_access
TestParser/function_calls
TestParser/complex_expressions
TestParser/error_cases
```

---

#### Task 1.5: Basic Evaluator (Days 8-9)
- [ ] Define `EvaluationContext` interface
- [ ] Implement `Evaluator` struct
- [ ] Implement `NewEvaluator(ast Node) *Evaluator`
- [ ] Implement `Evaluate(ctx EvaluationContext) (interface{}, error)`
- [ ] Implement evaluation for each AST node type:
  - `BinaryExpr`: Arithmetic, comparison, logical operators
  - `UnaryExpr`: Negation, NOT
  - `LiteralExpr`: Return value directly
  - `IdentifierExpr`: Look up in context
  - `ParenExpr`: Evaluate inner expression
- [ ] Implement operator short-circuiting (`&&`, `||`)
- [ ] Add comprehensive unit tests (80+ test cases)

**Files**:
- `backend/pkg/expr/evaluator.go`
- `backend/pkg/expr/evaluator_test.go`

**Test Cases**:
```go
TestEvaluator/arithmetic
TestEvaluator/comparisons
TestEvaluator/logical_operators
TestEvaluator/short_circuit
TestEvaluator/type_coercion
TestEvaluator/error_cases
```

---

#### Task 1.6: Type System (Day 10)
- [ ] Define `Type` enum (Unknown, Null, Bool, Number, String, Array, Object, Time)
- [ ] Implement `detectType(value interface{}) Type`
- [ ] Implement `coerceTypes(left, right interface{}) (interface{}, interface{}, error)`
- [ ] Implement type coercion rules:
  - Number + String → String concatenation
  - Number == String → Parse string as number
  - Bool && Number → Convert number to bool
  - Time comparison → Unix timestamp
- [ ] Add comprehensive type coercion tests

**Files**:
- `backend/pkg/expr/types.go`
- `backend/pkg/expr/types_test.go`

---

#### Task 1.7: Error Handling (Day 10)
- [ ] Define `ErrorType` enum (Syntax, Type, UndefinedReference, etc.)
- [ ] Implement `ExprError` struct (Type, Message, Line, Column, Source)
- [ ] Implement `Error() string` with formatted output
- [ ] Add error wrapping helpers
- [ ] Add comprehensive error tests

**Files**:
- `backend/pkg/expr/errors.go`
- `backend/pkg/expr/errors_test.go`

---

#### Deliverables Phase 1:
- ✅ Working lexer (tokenization)
- ✅ Working parser (AST construction)
- ✅ Working evaluator (simple expressions)
- ✅ Type system with coercion
- ✅ Error handling framework
- ✅ 95%+ test coverage
- ✅ Benchmarks for performance baseline

---

## Phase 2: Reference Resolution (5 days)

### Week 3: Node, Variable, Context References

#### Task 2.1: Field Access Evaluation (Day 11)
- [ ] Implement `evaluateFieldAccess(obj Node, field string) (interface{}, error)`
- [ ] Handle struct field access
- [ ] Handle map key access
- [ ] Handle nested field access (object.field1.field2)
- [ ] Add nil checking and error handling
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/expr/evaluator.go`
- Update `backend/pkg/expr/evaluator_test.go`

---

#### Task 2.2: Index Access Evaluation (Day 12)
- [ ] Implement `evaluateIndexAccess(obj, index Node) (interface{}, error)`
- [ ] Handle array indexing
- [ ] Handle map key indexing
- [ ] Handle string character access
- [ ] Add bounds checking
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/expr/evaluator.go`
- Update `backend/pkg/expr/evaluator_test.go`

---

#### Task 2.3: Reference Parsing (Day 13)
- [ ] Implement `parseNodeReference()` - `node.nodeId.output.field`
- [ ] Implement `parseVariableReference()` - `variables.varName`
- [ ] Implement `parseContextReference()` - `context.constName`
- [ ] Add reference validation
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/expr/parser.go`
- Update `backend/pkg/expr/parser_test.go`

---

#### Task 2.4: Reference Evaluation (Day 14)
- [ ] Implement node result lookup in `EvaluationContext`
- [ ] Implement variable lookup in `EvaluationContext`
- [ ] Implement context lookup in `EvaluationContext`
- [ ] Add `GetNodeResult(nodeID string) (interface{}, bool)` to context interface
- [ ] Add `GetVariable(name string) (interface{}, error)` to context interface
- [ ] Add `GetContextVariable(name string) (interface{}, bool)` to context interface
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/expr/evaluator.go`
- Update `backend/pkg/expr/evaluator_test.go`

---

#### Task 2.5: Dependency Extraction (Day 15)
- [ ] Implement `DependencyExtractor` struct
- [ ] Implement AST visitor pattern
- [ ] Extract node IDs from `node.nodeId.field` references
- [ ] Implement `Extract(ast Node) []string`
- [ ] Add comprehensive tests (including nested references)

**Files**:
- `backend/pkg/expr/dependencies.go`
- `backend/pkg/expr/dependencies_test.go`

---

#### Task 2.6: Integration with Executor Context (Day 15)
- [ ] Verify `ExecutionContext` interface compatibility
- [ ] Create adapter if needed
- [ ] Add integration tests with real workflow data
- [ ] Test end-to-end reference resolution

**Files**:
- `backend/pkg/expr/integration_test.go`

---

#### Deliverables Phase 2:
- ✅ Field access working (`.field`)
- ✅ Index access working (`[i]`)
- ✅ Node references working (`node.http1.output.status`)
- ✅ Variable references working (`variables.counter`)
- ✅ Context references working (`context.apiKey`)
- ✅ Dependency extraction working
- ✅ Integration tests passing

---

## Phase 3: Functions and Advanced Features (5 days)

### Week 4: Built-in Functions

#### Task 3.1: Function Call Infrastructure (Day 16)
- [ ] Implement `evaluateFunctionCall(fn *FunctionCallExpr) (interface{}, error)`
- [ ] Create function registry
- [ ] Define `FunctionHandler` type
- [ ] Implement function validation (argument count, types)
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/expr/evaluator.go`
- `backend/pkg/expr/functions.go`
- `backend/pkg/expr/functions_test.go`

---

#### Task 3.2: String Functions (Day 17)
- [ ] Implement `len(array|string) -> number`
- [ ] Implement `contains(string, substring) -> bool`
- [ ] Implement `startsWith(string, prefix) -> bool`
- [ ] Implement `endsWith(string, suffix) -> bool`
- [ ] Implement `lower(string) -> string`
- [ ] Implement `upper(string) -> string`
- [ ] Implement `trim(string) -> string`
- [ ] Add comprehensive tests for each function

**Files**:
- Update `backend/pkg/expr/functions.go`
- Update `backend/pkg/expr/functions_test.go`

---

#### Task 3.3: Regex and Parsing Functions (Day 18)
- [ ] Implement `matches(string, pattern) -> bool`
- [ ] Add regex timeout protection (100ms)
- [ ] Add regex cache
- [ ] Implement `parseInt(string) -> number`
- [ ] Implement `parseFloat(string) -> number`
- [ ] Implement `toString(value) -> string`
- [ ] Add comprehensive tests
- [ ] Add security tests (ReDoS protection)

**Files**:
- Update `backend/pkg/expr/functions.go`
- Update `backend/pkg/expr/functions_test.go`

---

#### Task 3.4: Time Functions (Day 19)
- [ ] Implement `now() -> time.Time`
- [ ] Add time comparison support
- [ ] Add time arithmetic support
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/expr/functions.go`
- Update `backend/pkg/expr/functions_test.go`

---

#### Task 3.5: Array Functions (Day 20)
- [ ] Ensure `len()` works with arrays
- [ ] Add array type validation
- [ ] Add comprehensive array tests

**Files**:
- Update `backend/pkg/expr/functions.go`
- Update `backend/pkg/expr/functions_test.go`

---

#### Deliverables Phase 3:
- ✅ All built-in functions implemented
- ✅ Function documentation complete
- ✅ Regex timeout protection working
- ✅ Time operations working
- ✅ Array operations working
- ✅ Security tests passing

---

## Phase 4: Integration and Migration (5 days)

### Week 5: Node Integration and Backward Compatibility

#### Task 4.1: Expression Public API (Day 21)
- [ ] Implement `Expression` struct (source, ast, dependencies)
- [ ] Implement `Compile(source string) (*Expression, error)`
- [ ] Implement `Evaluate(ctx EvaluationContext) (interface{}, error)`
- [ ] Implement `Dependencies() []string`
- [ ] Implement expression caching (sync.Map)
- [ ] Add `CompileWithCache(source string) (*Expression, error)`
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/expr/expr.go`
- Update `backend/pkg/expr/expr_test.go`

---

#### Task 4.2: Modify Condition Executor (Day 22)
- [ ] Update `ConditionExecutor.Execute()` to detect template expressions
- [ ] Add expression compilation and evaluation
- [ ] Maintain legacy `evaluateCondition()` for simple conditions
- [ ] Add backward compatibility tests
- [ ] Add integration tests with real workflows

**Files**:
- Update `backend/pkg/executor/condition.go`
- Update `backend/pkg/executor/condition_test.go` (if exists)
- Create integration tests

---

#### Task 4.3: Modify Switch Executor (Day 23)
- [ ] Update `SwitchExecutor.Execute()` for expression support
- [ ] Update case matching logic
- [ ] Add backward compatibility
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/executor/switch.go`
- Update tests

---

#### Task 4.4: Modify WhileLoop Executor (Day 23)
- [ ] Update `WhileLoopExecutor.Execute()` for expression support
- [ ] Add backward compatibility
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/executor/whileloop.go`
- Update tests

---

#### Task 4.5: Engine Dependency Extraction (Day 24)
- [ ] Implement `extractExpressionDependencies(nodes []Node) ([]Edge, error)`
- [ ] Add implicit edge creation for referenced nodes
- [ ] Update graph building to include implicit edges
- [ ] Ensure cycle detection works with implicit edges
- [ ] Add comprehensive tests

**Files**:
- Update `backend/pkg/engine/engine.go`
- Update `backend/pkg/engine/engine_test.go`

---

#### Task 4.6: Integration Testing (Day 25)
- [ ] Create end-to-end workflow tests
- [ ] Test complex multi-node workflows with expressions
- [ ] Test dependency extraction accuracy
- [ ] Test cycle detection with expressions
- [ ] Test backward compatibility thoroughly
- [ ] Performance regression testing

**Files**:
- `backend/pkg/engine/expression_integration_test.go`

---

#### Deliverables Phase 4:
- ✅ Public API complete
- ✅ All node executors updated
- ✅ Dependency extraction working
- ✅ Backward compatibility 100%
- ✅ Integration tests passing
- ✅ No performance regression

---

## Phase 5: Documentation and Polish (5 days)

### Week 6: Production Readiness

#### Task 5.1: API Documentation (Day 26)
- [ ] Document all public functions with godoc
- [ ] Add package-level documentation
- [ ] Create usage examples in godoc
- [ ] Generate godoc HTML

**Files**:
- Update all `*.go` files with godoc comments
- Create `backend/pkg/expr/doc.go`

---

#### Task 5.2: User Documentation (Day 27)
- [ ] Create expression syntax guide
- [ ] Create migration guide (simple → template expressions)
- [ ] Create built-in function reference
- [ ] Create troubleshooting guide
- [ ] Add real-world examples

**Files**:
- `docs/EXPRESSION_SYNTAX_GUIDE.md`
- `docs/EXPRESSION_MIGRATION_GUIDE.md`
- `docs/EXPRESSION_FUNCTIONS.md`
- Update `README.md` with expression feature

---

#### Task 5.3: Performance Optimization (Day 28)
- [ ] Profile expression evaluation
- [ ] Optimize hot paths
- [ ] Implement token pooling if needed
- [ ] Implement AST node pooling if needed
- [ ] Ensure benchmarks meet targets
- [ ] Document performance characteristics

**Files**:
- Update `backend/pkg/expr/expr_bench_test.go`
- Create `docs/EXPRESSION_PERFORMANCE.md`

---

#### Task 5.4: Security Audit (Day 29)
- [ ] Review for injection vulnerabilities
- [ ] Test resource limit enforcement
- [ ] Test timeout protection
- [ ] Test ReDoS protection
- [ ] Verify type safety
- [ ] Document security measures

**Files**:
- Create `docs/EXPRESSION_SECURITY.md`
- Update security tests

---

#### Task 5.5: Final Testing and Bug Fixes (Day 30)
- [ ] Final end-to-end testing
- [ ] User acceptance testing
- [ ] Fix any remaining bugs
- [ ] Verify all success criteria
- [ ] Create release notes

**Files**:
- `RELEASE_NOTES_EXPRESSION_SYSTEM.md`

---

#### Deliverables Phase 5:
- ✅ Complete API documentation
- ✅ Complete user documentation
- ✅ Performance targets met
- ✅ Security audit passed
- ✅ All tests passing
- ✅ Production ready

---

## Success Criteria Checklist

### Functional Requirements
- [ ] All expression types working (literals, operators, references, functions)
- [ ] Node output references functional
- [ ] Variable references functional
- [ ] Context references functional
- [ ] Boolean logic (AND, OR, NOT) working
- [ ] String operations working
- [ ] Array operations working
- [ ] Time operations working
- [ ] Type coercion working correctly

### Performance Requirements
- [ ] Simple comparison: <250ns
- [ ] Node reference: <400ns
- [ ] Complex expression: <800ns
- [ ] Compilation (cached): <100ns
- [ ] Cache hit rate >90%
- [ ] No performance regression vs. legacy

### Security Requirements
- [ ] Zero code injection vulnerabilities
- [ ] Resource limits enforced (depth, timeout)
- [ ] ReDoS protection working
- [ ] Type safety verified
- [ ] Security audit passed

### Quality Requirements
- [ ] Test coverage ≥95%
- [ ] All benchmarks passing
- [ ] All integration tests passing
- [ ] Documentation complete
- [ ] Code review passed

### Compatibility Requirements
- [ ] Backward compatibility 100%
- [ ] All existing tests passing
- [ ] Migration guide complete
- [ ] Legacy syntax still works

---

## Risk Mitigation

### Risk 1: Performance Regression
**Mitigation**:
- Implement caching early
- Profile continuously
- Optimize hot paths
- Keep legacy path for simple conditions

### Risk 2: Security Vulnerabilities
**Mitigation**:
- Security review after each phase
- Implement limits early (timeout, depth, length)
- Whitelist functions only
- No dynamic code execution

### Risk 3: Scope Creep
**Mitigation**:
- Stick to designed feature set
- Defer "nice-to-have" features to v2
- Focus on core functionality first

### Risk 4: Integration Issues
**Mitigation**:
- Frequent integration testing
- Early integration with condition node
- Comprehensive compatibility tests

---

## Dependencies

### External Dependencies
- None (standard library only) ✅

### Internal Dependencies
- `backend/pkg/types` - Type definitions
- `backend/pkg/executor` - Node executors
- `backend/pkg/engine` - Workflow engine
- `backend/pkg/state` - State management

---

## Team Allocation

**Recommended**:
- 1 Senior Go Developer (lead)
- 1 Mid-level Go Developer (implementation)
- 1 QA Engineer (testing)
- 0.5 Technical Writer (documentation)

**Estimated Effort**: 6 person-weeks

---

## Progress Tracking

### Phase 1: Core Infrastructure
- [ ] Task 1.1: Project setup
- [ ] Task 1.2: Lexer
- [ ] Task 1.3: AST
- [ ] Task 1.4: Parser
- [ ] Task 1.5: Evaluator
- [ ] Task 1.6: Type system
- [ ] Task 1.7: Error handling

### Phase 2: Reference Resolution
- [ ] Task 2.1: Field access
- [ ] Task 2.2: Index access
- [ ] Task 2.3: Reference parsing
- [ ] Task 2.4: Reference evaluation
- [ ] Task 2.5: Dependency extraction
- [ ] Task 2.6: Integration

### Phase 3: Functions
- [ ] Task 3.1: Function infrastructure
- [ ] Task 3.2: String functions
- [ ] Task 3.3: Regex/parsing functions
- [ ] Task 3.4: Time functions
- [ ] Task 3.5: Array functions

### Phase 4: Integration
- [ ] Task 4.1: Public API
- [ ] Task 4.2: Condition executor
- [ ] Task 4.3: Switch executor
- [ ] Task 4.4: WhileLoop executor
- [ ] Task 4.5: Engine dependency extraction
- [ ] Task 4.6: Integration testing

### Phase 5: Polish
- [ ] Task 5.1: API documentation
- [ ] Task 5.2: User documentation
- [ ] Task 5.3: Performance optimization
- [ ] Task 5.4: Security audit
- [ ] Task 5.5: Final testing

---

## References

- **Design Document**: [docs/EXPRESSION_SYSTEM_DESIGN.md](docs/EXPRESSION_SYSTEM_DESIGN.md)
- **Executive Summary**: [EXPRESSION_SYSTEM_EXECUTIVE_SUMMARY.md](EXPRESSION_SYSTEM_EXECUTIVE_SUMMARY.md)
- **Current Code**: `backend/pkg/executor/helpers.go:evaluateCondition()`
- **Architecture**: `backend/pkg/` packages

---

**Last Updated**: November 1, 2025  
**Status**: Ready to Start - Awaiting Approval
