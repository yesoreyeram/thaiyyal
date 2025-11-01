# Advanced Expression Evaluation System Design
# Thaiyyal Workflow Engine

**Document Version**: 1.0  
**Date**: November 1, 2025  
**Author**: Enterprise Architecture Team  
**Status**: Design Specification

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Current State Analysis](#current-state-analysis)
3. [Requirements](#requirements)
4. [Syntax Style Comparison](#syntax-style-comparison)
5. [Recommended Architecture](#recommended-architecture)
6. [Expression Grammar](#expression-grammar)
7. [Dependency Extraction](#dependency-extraction)
8. [Evaluation Strategy](#evaluation-strategy)
9. [Type System](#type-system)
10. [Implementation Phases](#implementation-phases)
11. [Testing Strategy](#testing-strategy)
12. [Security Considerations](#security-considerations)
13. [Migration Path](#migration-path)
14. [Example Expressions](#example-expressions)

---

## Executive Summary

This document presents a comprehensive design for implementing an advanced expression evaluation system for Thaiyyal's IF/Condition nodes. The system will enable:

- **Node output references**: Access other nodes' results (e.g., `{{ node.http1.output.status == 200 }}`)
- **Variable references**: Reference workflow state (e.g., `{{ variables.counter > 10 }}`)
- **Complex boolean logic**: AND, OR, NOT operators
- **String operations**: Equality, contains, regex matching
- **Type-safe evaluation**: Automatic type checking and conversion
- **Dependency tracking**: Auto-detect node dependencies for topological sorting

**Key Design Decisions:**
- ✅ **Zero external dependencies** (standard library only)
- ✅ **Template-style syntax** (`{{ }}`) for familiarity and clarity
- ✅ **Recursive descent parser** for maintainability
- ✅ **Compiled expressions** for performance
- ✅ **Backward compatible** with simple conditions (">100", "==5")

---

## Current State Analysis

### Technology Stack
- **Backend**: Go 1.24.7 (zero external dependencies)
- **Architecture**: DAG-based workflow engine with topological sorting
- **Pattern**: Strategy pattern with executor registry

### Current Limitations

**File**: `backend/pkg/executor/helpers.go:evaluateCondition()`

```go
// Current implementation only supports:
// - "true", "false" - boolean constants
// - ">N", "<N", ">=N", "<=N", "==N", "!=N" - numeric comparisons
func evaluateCondition(condition string, value interface{}) bool {
    // ... simple string parsing ...
}
```

**Problems:**
1. ❌ Cannot reference other nodes' outputs
2. ❌ Cannot reference workflow variables
3. ❌ Cannot reference context values
4. ❌ No complex boolean expressions (AND, OR, NOT)
5. ❌ No string comparisons or pattern matching
6. ❌ Limited to single value comparisons
7. ❌ No dependency extraction for graph building

### Existing Infrastructure

**Advantages** ✅:
- Clean package structure (`pkg/engine`, `pkg/executor`, `pkg/state`, `pkg/graph`)
- Executor registry pattern for extensibility
- State management with variables, accumulator, counter, cache
- Template interpolation exists (`{{ variable.name }}`, `{{ const.name }}`)
- Topological sorting with cycle detection
- Context variables and constants support

---

## Requirements

### Functional Requirements

1. **Expression Language**:
   - Node output references: `{{ node.nodeId.output.field }}`
   - Variable references: `{{ variables.varName }}`
   - Context references: `{{ context.constName }}`
   - Current timestamp: `{{ now() }}`
   - Boolean logic: `&&`, `||`, `!`
   - Comparisons: `>`, `<`, `>=`, `<=`, `==`, `!=`
   - String operations: `contains()`, `matches()`, `startsWith()`, `endsWith()`
   - Array operations: `len()`, `in`, indexing `[i]`
   - Object field access: dot notation `.field`

2. **Dependency Graph**:
   - Extract all node references from expressions
   - Add implicit edges for referenced nodes
   - Validate no circular dependencies

3. **Error Handling**:
   - Compile-time syntax errors
   - Runtime type errors with clear messages
   - Undefined reference detection
   - Circular dependency detection

4. **Performance**:
   - Compile expressions once, evaluate many times
   - Minimal overhead per evaluation
   - Cache compiled expressions

### Non-Functional Requirements

1. **Zero External Dependencies**: Use only Go standard library
2. **Backward Compatibility**: Support existing simple conditions
3. **Security**: Prevent code injection, limit recursion, enforce timeouts
4. **Maintainability**: Clean, well-documented code
5. **Extensibility**: Easy to add new functions and operators

---

## Syntax Style Comparison

### Option 1: Template-Based (Recommended ⭐)

**Syntax**: `{{ expression }}`

**Examples**:
```javascript
{{ node.http1.output.status == 200 }}
{{ variables.counter > 10 && variables.enabled }}
{{ context.apiKey != "" }}
{{ len(node.array1.output.items) > 0 }}
{{ node.text1.output contains "error" }}
```

**Pros**:
- ✅ Consistent with existing template interpolation (`{{ variable.name }}`)
- ✅ Visually distinct from surrounding text
- ✅ Familiar to users of Jinja2, Handlebars, Go templates
- ✅ Clear boundaries for parser
- ✅ Works well in UI (syntax highlighting)
- ✅ Easy to extend with functions

**Cons**:
- ❌ Slightly more verbose than bare expressions
- ❌ Requires escaping if literal `{{` needed (rare)

**Parsing Complexity**: Low (clear delimiters)

**User Friendliness**: ⭐⭐⭐⭐⭐ (Excellent - familiar pattern)

---

### Option 2: JSONPath-Like

**Syntax**: `$.path.to.value`

**Examples**:
```
$.nodes.http1.output.status == 200
$.variables.counter > 10 && $.variables.enabled
$.context.apiKey != ""
```

**Pros**:
- ✅ JSON users may find familiar
- ✅ Clear hierarchy representation
- ✅ Less verbose

**Cons**:
- ❌ Inconsistent with existing template system
- ❌ Confusing with JSON data structures in workflow
- ❌ Less intuitive for non-JSON users
- ❌ Harder to distinguish from literal strings

**Parsing Complexity**: Medium

**User Friendliness**: ⭐⭐⭐ (Good for JSON experts)

---

### Option 3: JavaScript-Like

**Syntax**: Bare expressions with implicit scope

**Examples**:
```javascript
nodes.http1.output.status == 200
variables.counter > 10 && variables.enabled
context.apiKey != ""
```

**Pros**:
- ✅ Concise
- ✅ Familiar to JavaScript developers
- ✅ Natural for complex expressions

**Cons**:
- ❌ No clear boundary (starts/ends anywhere)
- ❌ Ambiguous with literal values
- ❌ Inconsistent with existing template system
- ❌ Parser must handle entire condition as expression

**Parsing Complexity**: High (ambiguous boundaries)

**User Friendliness**: ⭐⭐⭐⭐ (Very good for programmers)

---

### Option 4: Prefix-Based

**Syntax**: `@type.path`

**Examples**:
```
@node.http1.output.status == 200
@var.counter > 10 && @var.enabled
@ctx.apiKey != ""
```

**Pros**:
- ✅ Clear type prefix
- ✅ Concise
- ✅ Unambiguous references

**Cons**:
- ❌ Inconsistent with existing `{{ }}` templates
- ❌ Less familiar syntax
- ❌ Multiple prefix types to remember

**Parsing Complexity**: Medium

**User Friendliness**: ⭐⭐⭐ (Good but unfamiliar)

---

### Option 5: Custom DSL

**Syntax**: Domain-specific language

**Examples**:
```
WHEN node.http1.output.status EQUALS 200
WHEN variables.counter GREATER_THAN 10 AND variables.enabled
```

**Pros**:
- ✅ Very readable
- ✅ Self-documenting
- ✅ Non-programmers friendly

**Cons**:
- ❌ Verbose
- ❌ Complex parser implementation
- ❌ Limited expressiveness
- ❌ Inconsistent with ecosystem

**Parsing Complexity**: High

**User Friendliness**: ⭐⭐⭐⭐ (Good for non-programmers)

---

### Recommendation: Template-Based (Option 1) ⭐

**Rationale**:
1. **Consistency**: Matches existing `{{ variable.name }}` interpolation
2. **Familiarity**: Widely used pattern (Jinja2, Handlebars, Django templates)
3. **Clarity**: Clear visual boundaries
4. **Extensibility**: Easy to add functions and complex expressions
5. **Parsing**: Simple delimiter-based parsing
6. **UI Integration**: Excellent syntax highlighting support

**Example Evolution**:
```javascript
// Current (simple)
">100"

// New (advanced) - both work!
">100"  // backward compatible
"{{ value > 100 }}"  // new syntax
"{{ node.sensor1.output.temperature > 100 }}"  // full power
```

---

## Recommended Architecture

### Package Structure

```
backend/pkg/
├── expr/                    # NEW: Expression evaluation package
│   ├── expr.go              # Public API and Expression struct
│   ├── parser.go            # Recursive descent parser
│   ├── lexer.go             # Tokenizer/lexer
│   ├── evaluator.go         # Expression evaluation
│   ├── compiler.go          # Expression compilation
│   ├── ast.go               # Abstract syntax tree nodes
│   ├── types.go             # Type system and coercion
│   ├── functions.go         # Built-in functions (len, contains, etc.)
│   ├── dependencies.go      # Dependency extraction
│   ├── errors.go            # Error types
│   ├── expr_test.go         # Unit tests
│   └── examples_test.go     # Example usage tests
│
├── executor/
│   ├── condition.go         # MODIFIED: Use expr package
│   ├── switch.go            # MODIFIED: Use expr package
│   └── whileloop.go         # MODIFIED: Use expr package
│
├── engine/
│   └── engine.go            # MODIFIED: Dependency extraction
│
└── graph/
    └── graph.go             # UNMODIFIED: Existing cycle detection works
```

### Key Components

#### 1. Lexer (Tokenizer)

**Responsibility**: Convert expression string into tokens

```go
package expr

type TokenType int

const (
    // Literals
    TokenNumber TokenType = iota
    TokenString
    TokenTrue
    TokenFalse
    TokenNull
    
    // Identifiers and keywords
    TokenIdentifier  // node, variables, context, etc.
    TokenDot         // .
    TokenComma       // ,
    
    // Operators
    TokenAnd         // &&
    TokenOr          // ||
    TokenNot         // !
    TokenEQ          // ==
    TokenNE          // !=
    TokenLT          // <
    TokenLE          // <=
    TokenGT          // >
    TokenGE          // >=
    TokenPlus        // +
    TokenMinus       // -
    TokenStar        // *
    TokenSlash       // /
    TokenPercent     // %
    
    // Delimiters
    TokenLParen      // (
    TokenRParen      // )
    TokenLBracket    // [
    TokenRBracket    // ]
    TokenLBrace      // {
    TokenRBrace      // }
    
    // Special
    TokenEOF
    TokenError
)

type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}

type Lexer struct {
    input   string
    pos     int
    line    int
    column  int
}

func NewLexer(input string) *Lexer
func (l *Lexer) NextToken() Token
```

#### 2. Parser (Recursive Descent)

**Responsibility**: Build Abstract Syntax Tree (AST) from tokens

```go
package expr

// AST Node types
type Node interface {
    node()  // marker method
}

type BinaryExpr struct {
    Left     Node
    Operator TokenType
    Right    Node
}

type UnaryExpr struct {
    Operator TokenType
    Operand  Node
}

type LiteralExpr struct {
    Value interface{}
}

type IdentifierExpr struct {
    Name string
}

type FieldAccessExpr struct {
    Object Node
    Field  string
}

type IndexAccessExpr struct {
    Object Node
    Index  Node
}

type FunctionCallExpr struct {
    Function string
    Args     []Node
}

// Parser
type Parser struct {
    lexer   *Lexer
    current Token
    peek    Token
}

func NewParser(lexer *Lexer) *Parser
func (p *Parser) Parse() (Node, error)
```

#### 3. Evaluator

**Responsibility**: Execute AST against runtime context

```go
package expr

type EvaluationContext interface {
    // Get node result
    GetNodeResult(nodeID string) (interface{}, bool)
    
    // Get variable
    GetVariable(name string) (interface{}, error)
    
    // Get context value
    GetContextVariable(name string) (interface{}, bool)
    GetContextConstant(name string) (interface{}, bool)
    
    // Current time
    Now() time.Time
}

type Evaluator struct {
    ast AST Node
}

func NewEvaluator(ast Node) *Evaluator
func (e *Evaluator) Evaluate(ctx EvaluationContext) (interface{}, error)
```

#### 4. Expression (Public API)

**Responsibility**: High-level API for expression evaluation

```go
package expr

// Expression represents a compiled expression
type Expression struct {
    source       string
    ast          Node
    dependencies []string  // Node IDs referenced
}

// Compile parses and validates an expression
func Compile(source string) (*Expression, error)

// Evaluate executes the expression against a context
func (e *Expression) Evaluate(ctx EvaluationContext) (interface{}, error)

// Dependencies returns all node IDs referenced in expression
func (e *Expression) Dependencies() []string

// String returns the original source
func (e *Expression) String() string
```

#### 5. Dependency Extractor

**Responsibility**: Extract node references from AST

```go
package expr

type DependencyExtractor struct {
    dependencies map[string]bool
}

func NewDependencyExtractor() *DependencyExtractor
func (d *DependencyExtractor) Extract(ast Node) []string
func (d *DependencyExtractor) visit(node Node)  // AST visitor pattern
```

### Integration Points

#### 1. Condition Node Executor

```go
// backend/pkg/executor/condition.go

func (e *ConditionExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    condition := *node.Data.Condition
    
    // Check if it's a template expression
    if strings.HasPrefix(condition, "{{") && strings.HasSuffix(condition, "}}") {
        // New expression evaluation
        exprStr := strings.TrimSpace(condition[2:len(condition)-2])
        compiled, err := expr.Compile(exprStr)
        if err != nil {
            return nil, fmt.Errorf("invalid expression: %w", err)
        }
        
        result, err := compiled.Evaluate(ctx)
        if err != nil {
            return nil, fmt.Errorf("evaluation error: %w", err)
        }
        
        conditionMet, ok := result.(bool)
        if !ok {
            return nil, fmt.Errorf("expression must evaluate to boolean, got %T", result)
        }
        
        // ... rest of execution
    } else {
        // Legacy simple condition evaluation
        conditionMet := evaluateCondition(condition, input)
        // ... rest of execution
    }
}
```

#### 2. Engine Dependency Extraction

```go
// backend/pkg/engine/engine.go

func (e *Engine) extractDependencies(nodes []types.Node) ([]types.Edge, error) {
    var implicitEdges []types.Edge
    
    for _, node := range nodes {
        switch node.Type {
        case types.NodeTypeCondition, types.NodeTypeSwitch, types.NodeTypeWhileLoop:
            condition := getNodeCondition(node)
            if isTemplateExpression(condition) {
                compiled, err := expr.Compile(extractExpression(condition))
                if err != nil {
                    return nil, fmt.Errorf("node %s: %w", node.ID, err)
                }
                
                // Add implicit edges for referenced nodes
                for _, depNodeID := range compiled.Dependencies() {
                    implicitEdges = append(implicitEdges, types.Edge{
                        ID:     fmt.Sprintf("implicit_%s_%s", depNodeID, node.ID),
                        Source: depNodeID,
                        Target: node.ID,
                    })
                }
            }
        }
    }
    
    return implicitEdges, nil
}
```

---

## Expression Grammar

### EBNF Grammar Specification

```ebnf
expression      = logical_or ;

logical_or      = logical_and { "||" logical_and } ;

logical_and     = equality { "&&" equality } ;

equality        = comparison { ( "==" | "!=" ) comparison } ;

comparison      = additive { ( ">" | ">=" | "<" | "<=" ) additive } ;

additive        = multiplicative { ( "+" | "-" ) multiplicative } ;

multiplicative  = unary { ( "*" | "/" | "%" ) unary } ;

unary           = "!" unary
                | "-" unary
                | postfix ;

postfix         = primary { field_access | index_access | function_call } ;

field_access    = "." IDENTIFIER ;

index_access    = "[" expression "]" ;

function_call   = "(" [ expression { "," expression } ] ")" ;

primary         = NUMBER
                | STRING
                | "true"
                | "false"
                | "null"
                | IDENTIFIER
                | "(" expression ")"
                | reference ;

reference       = "node" "." IDENTIFIER { field_access }
                | "variables" "." IDENTIFIER
                | "context" "." IDENTIFIER
                | "now" "(" ")" ;

IDENTIFIER      = LETTER { LETTER | DIGIT | "_" } ;
NUMBER          = DIGIT { DIGIT } [ "." DIGIT { DIGIT } ] ;
STRING          = '"' { ANY_CHAR except '"' } '"'
                | "'" { ANY_CHAR except "'" } "'" ;
LETTER          = "a" ... "z" | "A" ... "Z" ;
DIGIT           = "0" ... "9" ;
```

### Operator Precedence (Highest to Lowest)

1. **Primary**: Literals, identifiers, parentheses, references
2. **Postfix**: Field access `.`, indexing `[]`, function calls `()`
3. **Unary**: `!` (NOT), `-` (negation)
4. **Multiplicative**: `*`, `/`, `%`
5. **Additive**: `+`, `-`
6. **Comparison**: `<`, `<=`, `>`, `>=`
7. **Equality**: `==`, `!=`
8. **Logical AND**: `&&`
9. **Logical OR**: `||`

### Example Parse Trees

**Example 1**: `{{ node.http1.output.status == 200 }}`

```
BinaryExpr (==)
├── FieldAccessExpr
│   ├── FieldAccessExpr
│   │   ├── FieldAccessExpr
│   │   │   ├── IdentifierExpr("node")
│   │   │   └── "http1"
│   │   └── "output"
│   └── "status"
└── LiteralExpr(200)
```

**Example 2**: `{{ variables.count > 10 && variables.enabled }}`

```
BinaryExpr (&&)
├── BinaryExpr (>)
│   ├── FieldAccessExpr
│   │   ├── IdentifierExpr("variables")
│   │   └── "count"
│   └── LiteralExpr(10)
└── FieldAccessExpr
    ├── IdentifierExpr("variables")
    └── "enabled"
```

---

## Dependency Extraction

### Algorithm

**Goal**: Extract all node IDs referenced in an expression

**Approach**: AST Visitor Pattern

```go
package expr

type DependencyExtractor struct {
    dependencies map[string]bool
}

func (d *DependencyExtractor) Extract(ast Node) []string {
    d.dependencies = make(map[string]bool)
    d.visit(ast)
    
    // Convert to slice
    result := make([]string, 0, len(d.dependencies))
    for dep := range d.dependencies {
        result = append(result, dep)
    }
    return result
}

func (d *DependencyExtractor) visit(node Node) {
    switch n := node.(type) {
    case *BinaryExpr:
        d.visit(n.Left)
        d.visit(n.Right)
        
    case *UnaryExpr:
        d.visit(n.Operand)
        
    case *FieldAccessExpr:
        // Check if this is a node reference
        if id, ok := extractNodeReference(n); ok {
            d.dependencies[id] = true
        }
        d.visit(n.Object)
        
    case *IndexAccessExpr:
        d.visit(n.Object)
        d.visit(n.Index)
        
    case *FunctionCallExpr:
        for _, arg := range n.Args {
            d.visit(arg)
        }
        
    // Leaf nodes: LiteralExpr, IdentifierExpr
    }
}

func extractNodeReference(n *FieldAccessExpr) (string, bool) {
    // node.http1.output.status -> "http1"
    current := n
    var path []string
    
    for current != nil {
        if fieldAccess, ok := current.(*FieldAccessExpr); ok {
            path = append([]string{fieldAccess.Field}, path...)
            current = fieldAccess.Object.(*FieldAccessExpr)
        } else if ident, ok := current.(*IdentifierExpr); ok {
            if ident.Name == "node" && len(path) > 0 {
                return path[0], true  // Return node ID
            }
            return "", false
        } else {
            break
        }
    }
    
    return "", false
}
```

### Integration with Graph Builder

```go
// backend/pkg/engine/engine.go

func (e *Engine) buildGraph(payload types.Payload) error {
    // Extract implicit dependencies from expressions
    implicitEdges, err := e.extractExpressionDependencies(payload.Nodes)
    if err != nil {
        return fmt.Errorf("dependency extraction failed: %w", err)
    }
    
    // Combine explicit and implicit edges
    allEdges := append(payload.Edges, implicitEdges...)
    
    // Create graph with all edges
    e.graph = graph.New(payload.Nodes, allEdges)
    
    // Detect cycles (including implicit dependencies)
    if err := e.graph.DetectCycles(); err != nil {
        return fmt.Errorf("circular dependency: %w", err)
    }
    
    return nil
}
```

---

## Evaluation Strategy

### Compilation vs Interpretation

**Decision**: **Compile once, evaluate many times** ✅

**Rationale**:
1. Expressions are evaluated multiple times (loops, retries)
2. Parsing overhead is eliminated after first compile
3. Dependency extraction happens once
4. Syntax errors caught early

### Caching Strategy

```go
package expr

var (
    // Global expression cache
    expressionCache sync.Map  // map[string]*Expression
    maxCacheSize    = 1000
)

// CompileWithCache compiles or retrieves cached expression
func CompileWithCache(source string) (*Expression, error) {
    // Check cache first
    if cached, ok := expressionCache.Load(source); ok {
        return cached.(*Expression), nil
    }
    
    // Compile new expression
    expr, err := Compile(source)
    if err != nil {
        return nil, err
    }
    
    // Store in cache (with LRU eviction if needed)
    expressionCache.Store(source, expr)
    
    return expr, nil
}
```

### Performance Optimizations

1. **Token Pooling**: Reuse token objects
2. **AST Node Pooling**: Reuse AST nodes
3. **String Interning**: Common identifiers ("node", "variables", "context")
4. **Operator Short-Circuit**: `&&` and `||` short-circuit evaluation
5. **Type Caching**: Cache type conversions

### Benchmarks (Target)

```
BenchmarkSimpleCondition-8         5000000    250 ns/op    0 allocs/op
BenchmarkComplexExpression-8       2000000    650 ns/op    2 allocs/op
BenchmarkNodeReference-8           3000000    400 ns/op    1 allocs/op
BenchmarkCompilation-8              100000  12000 ns/op  150 allocs/op
```

---

## Type System

### Type Hierarchy

```go
package expr

type Type int

const (
    TypeUnknown Type = iota
    TypeNull
    TypeBool
    TypeNumber   // float64
    TypeString
    TypeArray
    TypeObject
    TypeTime
)
```

### Type Coercion Rules

```go
// Automatic coercion rules (safe conversions)
// Number + String -> String (concatenation)
// Number == String -> parse string as number
// Bool && Number -> convert number to bool (0 = false, non-zero = true)
// Time comparison -> convert to Unix timestamp

func coerceTypes(left, right interface{}) (interface{}, interface{}, error) {
    leftType := detectType(left)
    rightType := detectType(right)
    
    // Same types - no coercion needed
    if leftType == rightType {
        return left, right, nil
    }
    
    // Number + String
    if leftType == TypeNumber && rightType == TypeString {
        return left, parseNumber(right), nil
    }
    if leftType == TypeString && rightType == TypeNumber {
        return parseNumber(left), right, nil
    }
    
    // Bool + Number
    if leftType == TypeBool && rightType == TypeNumber {
        return left, numberToBool(right), nil
    }
    if leftType == TypeNumber && rightType == TypeBool {
        return numberToBool(left), right, nil
    }
    
    return nil, nil, fmt.Errorf("cannot coerce %v to %v", leftType, rightType)
}
```

### Type Checking

```go
// Type checking during compilation (static analysis)
type TypeChecker struct {
    errors []error
}

func (tc *TypeChecker) Check(ast Node) error {
    tc.visit(ast)
    if len(tc.errors) > 0 {
        return tc.errors[0]  // Return first error
    }
    return nil
}

func (tc *TypeChecker) visit(node Node) Type {
    switch n := node.(type) {
    case *BinaryExpr:
        leftType := tc.visit(n.Left)
        rightType := tc.visit(n.Right)
        
        switch n.Operator {
        case TokenAnd, TokenOr:
            if leftType != TypeBool || rightType != TypeBool {
                tc.errors = append(tc.errors, 
                    fmt.Errorf("logical operator requires boolean operands"))
            }
            return TypeBool
            
        case TokenEQ, TokenNE:
            return TypeBool
            
        case TokenLT, TokenLE, TokenGT, TokenGE:
            if !isComparable(leftType) || !isComparable(rightType) {
                tc.errors = append(tc.errors,
                    fmt.Errorf("comparison requires comparable types"))
            }
            return TypeBool
        }
    }
    // ... more cases
}
```

---

## Implementation Phases

### Phase 1: Core Infrastructure (Week 1-2)

**Goal**: Lexer, parser, and basic evaluator

**Tasks**:
- [ ] Create `pkg/expr` package structure
- [ ] Implement lexer with token types
- [ ] Implement recursive descent parser
- [ ] Define AST node types
- [ ] Implement basic evaluator (literals, binary ops)
- [ ] Add comprehensive unit tests
- [ ] Document grammar specification

**Deliverables**:
- Working lexer that tokenizes expressions
- Parser that builds AST from tokens
- Evaluator for simple expressions: `10 > 5`, `true && false`

**Success Criteria**:
- All unit tests passing
- 95%+ code coverage
- Benchmarks under 1µs for simple expressions

---

### Phase 2: Reference Resolution (Week 3)

**Goal**: Support node, variable, and context references

**Tasks**:
- [ ] Implement `FieldAccessExpr` and `IndexAccessExpr`
- [ ] Add reference parsing (`node.`, `variables.`, `context.`)
- [ ] Implement `EvaluationContext` interface
- [ ] Add dependency extraction from AST
- [ ] Integrate with existing executor context
- [ ] Add tests for all reference types

**Deliverables**:
- Working references: `{{ node.http1.output.status }}`
- Dependency extraction working
- Integration tests with executor context

**Success Criteria**:
- Can reference any node output
- Dependencies correctly extracted
- Type-safe field access

---

### Phase 3: Functions and Advanced Features (Week 4)

**Goal**: Built-in functions and complex operations

**Tasks**:
- [ ] Implement function call parsing and evaluation
- [ ] Add built-in functions:
  - `len(array)` - array/string length
  - `contains(str, substr)` - substring check
  - `matches(str, pattern)` - regex matching
  - `startsWith(str, prefix)` - prefix check
  - `endsWith(str, suffix)` - suffix check
  - `now()` - current timestamp
  - `parseInt(str)` - string to number
  - `toString(val)` - value to string
- [ ] Add array indexing support
- [ ] Add object field navigation
- [ ] Comprehensive function tests

**Deliverables**:
- All built-in functions working
- Function documentation
- Example expressions for each function

**Success Criteria**:
- All functions tested
- Clear error messages for invalid calls
- Performance benchmarks meet targets

---

### Phase 4: Integration and Migration (Week 5)

**Goal**: Integrate with condition nodes and migrate existing code

**Tasks**:
- [ ] Modify `ConditionExecutor` to use expression system
- [ ] Modify `SwitchExecutor` to use expressions
- [ ] Modify `WhileLoopExecutor` to use expressions
- [ ] Update engine to extract dependencies
- [ ] Add backward compatibility layer
- [ ] Update frontend to support new syntax
- [ ] Migrate existing tests
- [ ] Add migration documentation

**Deliverables**:
- All node types using new system
- Backward compatibility working
- Integration tests passing
- Migration guide for users

**Success Criteria**:
- Zero breaking changes for existing workflows
- All new features working
- Performance equal or better than before

---

### Phase 5: Documentation and Polish (Week 6)

**Goal**: Production-ready with full documentation

**Tasks**:
- [ ] Write comprehensive API documentation
- [ ] Create user guide with examples
- [ ] Add expression playground in UI
- [ ] Implement syntax validation in frontend
- [ ] Add syntax highlighting for expressions
- [ ] Performance optimization pass
- [ ] Security audit
- [ ] Final testing and bug fixes

**Deliverables**:
- Complete user documentation
- API reference
- Interactive examples
- Performance report
- Security audit report

**Success Criteria**:
- All documentation complete
- No known security issues
- Performance targets met
- User acceptance testing passed

---

## Testing Strategy

### Unit Tests

**Coverage Target**: ≥95%

```go
// expr_test.go
func TestLexer(t *testing.T) {
    tests := []struct {
        input    string
        expected []Token
    }{
        {"10 > 5", []Token{Number, GT, Number, EOF}},
        {"node.http1", []Token{Identifier, Dot, Identifier, EOF}},
    }
    // ... test implementation
}

func TestParser(t *testing.T) {
    tests := []struct {
        input       string
        expectedAST Node
        expectError bool
    }{
        // ... test cases
    }
}

func TestEvaluator(t *testing.T) {
    tests := []struct {
        expression string
        context    map[string]interface{}
        expected   interface{}
    }{
        {
            "{{ node.http1.output.status == 200 }}",
            map[string]interface{}{
                "node.http1.output.status": 200,
            },
            true,
        },
    }
}
```

### Integration Tests

```go
// integration_test.go
func TestConditionNodeWithExpression(t *testing.T) {
    // Create workflow with condition node using new expression
    workflow := `{
        "nodes": [
            {"id": "http1", "type": "http", "data": {"url": "https://api.example.com"}},
            {"id": "cond1", "type": "condition", "data": {
                "condition": "{{ node.http1.output.status == 200 }}"
            }}
        ],
        "edges": [
            {"source": "http1", "target": "cond1"}
        ]
    }`
    
    engine, _ := engine.New([]byte(workflow))
    result, err := engine.Execute(context.Background())
    
    // Assert result
}
```

### Benchmark Tests

```go
// expr_bench_test.go
func BenchmarkSimpleExpression(b *testing.B) {
    expr, _ := Compile("10 > 5")
    ctx := newMockContext()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        expr.Evaluate(ctx)
    }
}

func BenchmarkComplexExpression(b *testing.B) {
    expr, _ := Compile("{{ node.http1.output.status == 200 && len(node.array1.output) > 0 }}")
    ctx := newMockContext()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        expr.Evaluate(ctx)
    }
}
```

### Edge Case Tests

```go
func TestEdgeCases(t *testing.T) {
    cases := []struct {
        name        string
        expression  string
        expectError bool
    }{
        {"Empty expression", "{{}}", true},
        {"Unclosed brace", "{{ 10 > 5", true},
        {"Undefined node", "{{ node.unknown.value }}", true},
        {"Type mismatch", "{{ 'text' > 10 }}", true},
        {"Division by zero", "{{ 10 / 0 }}", true},
        {"Circular reference", "{{ node.A.value == node.B.value }}", false},  // Caught by graph
        {"Deep nesting", "{{ ((((10)))) }}", false},
        {"Unicode identifiers", "{{ variables.café > 0 }}", false},
    }
}
```

---

## Security Considerations

### 1. Injection Prevention

**Threat**: User-controlled expressions executing arbitrary code

**Mitigation**:
```go
// Only allow safe operations
// No: eval(), exec(), system calls, file I/O
// Yes: Arithmetic, comparison, field access, safe functions

// Whitelist of allowed functions
var allowedFunctions = map[string]bool{
    "len": true, "contains": true, "matches": true,
    "startsWith": true, "endsWith": true, "now": true,
}

func (e *Evaluator) evaluateFunctionCall(fn *FunctionCallExpr) (interface{}, error) {
    if !allowedFunctions[fn.Function] {
        return nil, fmt.Errorf("function not allowed: %s", fn.Function)
    }
    // ... safe execution
}
```

### 2. Resource Limits

**Threat**: Infinite loops, stack overflow, memory exhaustion

**Mitigation**:
```go
type EvaluatorConfig struct {
    MaxDepth      int           // Max recursion depth (default: 100)
    MaxIterations int           // Max loop iterations (default: 10000)
    Timeout       time.Duration // Max evaluation time (default: 1s)
}

func (e *Evaluator) Evaluate(ctx EvaluationContext) (interface{}, error) {
    // Set timeout
    evalCtx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
    defer cancel()
    
    return e.evaluateWithContext(evalCtx, e.ast, 0)
}

func (e *Evaluator) evaluateWithContext(ctx context.Context, node Node, depth int) (interface{}, error) {
    // Check timeout
    select {
    case <-ctx.Done():
        return nil, fmt.Errorf("evaluation timeout")
    default:
    }
    
    // Check depth
    if depth > e.config.MaxDepth {
        return nil, fmt.Errorf("maximum recursion depth exceeded")
    }
    
    // ... evaluation
}
```

### 3. Type Safety

**Threat**: Type confusion attacks

**Mitigation**:
```go
// Strict type checking during evaluation
func (e *Evaluator) evaluateBinaryOp(op TokenType, left, right interface{}) (interface{}, error) {
    // Validate types for operation
    switch op {
    case TokenAnd, TokenOr:
        leftBool, ok1 := left.(bool)
        rightBool, ok2 := right.(bool)
        if !ok1 || !ok2 {
            return nil, fmt.Errorf("logical operator requires boolean operands")
        }
        // ... safe operation
    }
}
```

### 4. Regular Expression Safety

**Threat**: ReDoS (Regular Expression Denial of Service)

**Mitigation**:
```go
import "regexp"

var (
    regexCache   sync.Map
    regexTimeout = 100 * time.Millisecond
)

func matches(str, pattern string) (bool, error) {
    // Compile with timeout
    re, err := regexp.Compile(pattern)
    if err != nil {
        return false, fmt.Errorf("invalid regex: %w", err)
    }
    
    // Execute with timeout
    done := make(chan bool, 1)
    go func() {
        done <- re.MatchString(str)
    }()
    
    select {
    case result := <-done:
        return result, nil
    case <-time.After(regexTimeout):
        return false, fmt.Errorf("regex timeout")
    }
}
```

### 5. Input Validation

**Threat**: Malformed or malicious input

**Mitigation**:
```go
func Compile(source string) (*Expression, error) {
    // Validate length
    if len(source) > 10000 {
        return nil, fmt.Errorf("expression too long (max 10000 chars)")
    }
    
    // Validate characters (prevent control characters)
    if !isValidExpression(source) {
        return nil, fmt.Errorf("expression contains invalid characters")
    }
    
    // ... parsing
}
```

---

## Migration Path

### Backward Compatibility

**Guarantee**: All existing simple conditions continue to work

```go
// backend/pkg/executor/condition.go

func (e *ConditionExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
    condition := *node.Data.Condition
    input := ctx.GetNodeInputs(node.ID)[0]
    
    var conditionMet bool
    
    // Check if it's a template expression
    if isTemplateExpression(condition) {
        // NEW: Expression evaluation
        conditionMet, err = evaluateTemplateExpression(condition, ctx)
        if err != nil {
            return nil, err
        }
    } else {
        // LEGACY: Simple condition evaluation (unchanged)
        conditionMet = evaluateCondition(condition, input)
    }
    
    return map[string]interface{}{
        "value":         input,
        "condition_met": conditionMet,
    }, nil
}

func isTemplateExpression(s string) bool {
    return strings.HasPrefix(s, "{{") && strings.HasSuffix(s, "}}")
}
```

### Migration Guide for Users

**Example Migration**:

```javascript
// OLD (Simple Condition)
{
  "type": "condition",
  "data": {
    "condition": ">100"
  }
}

// NEW (Same Behavior)
{
  "type": "condition",
  "data": {
    "condition": ">100"  // Still works!
  }
}

// NEW (Advanced Expression)
{
  "type": "condition",
  "data": {
    "condition": "{{ value > 100 }}"  // Explicit template
  }
}

// NEW (Node Reference)
{
  "type": "condition",
  "data": {
    "condition": "{{ node.sensor1.output.temperature > 100 }}"
  }
}
```

### Deprecation Strategy

**Timeline**:
- **Phase 1** (Now): Both syntaxes supported
- **Phase 2** (+6 months): Encourage migration via documentation
- **Phase 3** (+12 months): Deprecation warning for simple syntax
- **Phase 4** (+24 months): Simple syntax becomes legacy-only

**No Breaking Changes**: Simple syntax remains supported indefinitely for backward compatibility

---

## Example Expressions

### Basic Comparisons

```javascript
// Numeric comparisons
{{ value > 100 }}
{{ value >= 100 }}
{{ value < 50 }}
{{ value <= 50 }}
{{ value == 42 }}
{{ value != 0 }}

// String comparisons
{{ text == "hello" }}
{{ text != "error" }}
{{ text > "abc" }}  // Lexicographic comparison
```

### Node Output References

```javascript
// HTTP status check
{{ node.http1.output.status == 200 }}

// Nested field access
{{ node.http1.output.data.user.age > 18 }}

// Array length
{{ len(node.http1.output.items) > 0 }}

// Multiple node references
{{ node.sensor1.output.temperature > node.threshold.output.max }}
```

### Variable References

```javascript
// Simple variable check
{{ variables.counter > 10 }}
{{ variables.enabled == true }}
{{ variables.name != "" }}

// Combined with node reference
{{ variables.threshold < node.sensor1.output.value }}
```

### Context References

```javascript
// Context constants
{{ context.apiKey != "" }}
{{ context.environment == "production" }}

// Context variables (mutable)
{{ context.retryCount < 3 }}
```

### Boolean Logic

```javascript
// AND
{{ node.http1.output.status == 200 && len(node.http1.output.data) > 0 }}

// OR
{{ node.sensor1.output.temperature > 100 || node.sensor2.output.temperature > 100 }}

// NOT
{{ !node.http1.output.error }}

// Complex combinations
{{ (node.http1.output.status == 200 || node.http1.output.status == 201) && !node.http1.output.error }}
```

### String Operations

```javascript
// Contains
{{ contains(node.http1.output.message, "success") }}

// Starts with
{{ startsWith(node.http1.output.url, "https://") }}

// Ends with
{{ endsWith(node.file1.output.name, ".json") }}

// Regex match
{{ matches(node.text1.output, "^[A-Z]{3}-\\d{4}$") }}
```

### Array and Object Operations

```javascript
// Array length
{{ len(node.http1.output.items) > 0 }}

// Array indexing
{{ node.http1.output.items[0].status == "active" }}

// Object field check
{{ node.http1.output.data.user != null }}
{{ node.http1.output.data.user.verified == true }}

// Nested access
{{ node.http1.output.users[0].permissions.admin == true }}
```

### Time-Based Conditions

```javascript
// Current time comparison
{{ now() > node.schedule.output.startTime }}

// Time difference
{{ now() - node.cache1.output.timestamp > 3600 }}  // More than 1 hour
```

### Complex Real-World Examples

```javascript
// API health check
{{ node.http1.output.status == 200 && 
   node.http1.output.responseTime < 500 &&
   contains(node.http1.output.body, "healthy") }}

// Multi-sensor threshold
{{ (node.sensor1.output.temperature > variables.threshold ||
    node.sensor2.output.temperature > variables.threshold) &&
   variables.alertsEnabled == true }}

// Retry logic
{{ node.http1.output.status >= 500 && 
   context.retryCount < context.maxRetries }}

// Data validation
{{ len(node.form1.output.email) > 0 &&
   matches(node.form1.output.email, "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$") &&
   len(node.form1.output.password) >= 8 }}

// Cache hit decision
{{ node.cache1.output.hit == true &&
   now() - node.cache1.output.timestamp < variables.cacheTTL }}
```

---

## Appendix A: AST Node Definitions

```go
package expr

// Base node interface
type Node interface {
    node()  // marker method
}

// Expression nodes
type BinaryExpr struct {
    Left     Node
    Operator TokenType
    Right    Node
}

type UnaryExpr struct {
    Operator TokenType
    Operand  Node
}

type LiteralExpr struct {
    Value interface{}  // bool, float64, string, nil
}

type IdentifierExpr struct {
    Name string
}

type FieldAccessExpr struct {
    Object Node
    Field  string
}

type IndexAccessExpr struct {
    Object Node
    Index  Node
}

type FunctionCallExpr struct {
    Function string
    Args     []Node
}

type ParenExpr struct {
    Expr Node
}

// Reference nodes
type NodeReference struct {
    NodeID string
    Path   []string  // output.status.code
}

type VariableReference struct {
    Name string
}

type ContextReference struct {
    Name string
}
```

---

## Appendix B: Error Types

```go
package expr

type ErrorType int

const (
    SyntaxError ErrorType = iota
    TypeError
    UndefinedReference
    CircularDependency
    EvaluationError
    TimeoutError
    SecurityError
)

type ExprError struct {
    Type    ErrorType
    Message string
    Line    int
    Column  int
    Source  string
}

func (e *ExprError) Error() string {
    return fmt.Sprintf("[%d:%d] %s: %s", e.Line, e.Column, e.Type, e.Message)
}
```

---

## Appendix C: Performance Benchmarks

**Target Performance** (on modern hardware):

| Operation | Time | Allocations |
|-----------|------|-------------|
| Simple comparison (`10 > 5`) | 200ns | 0 |
| Field access (`node.http1.output.status`) | 400ns | 1 |
| Boolean logic (`a && b || c`) | 300ns | 0 |
| Function call (`len(array)`) | 500ns | 2 |
| Complex expression | 800ns | 3 |
| Compilation (first time) | 10µs | 150 |
| Compilation (cached) | 50ns | 0 |

---

## Appendix D: Built-in Functions Reference

| Function | Signature | Description | Example |
|----------|-----------|-------------|---------|
| `len()` | `len(array\|string) -> number` | Returns length | `len(items) > 0` |
| `contains()` | `contains(string, substring) -> bool` | Substring check | `contains(text, "error")` |
| `matches()` | `matches(string, pattern) -> bool` | Regex match | `matches(email, ".*@.*")` |
| `startsWith()` | `startsWith(string, prefix) -> bool` | Prefix check | `startsWith(url, "https")` |
| `endsWith()` | `endsWith(string, suffix) -> bool` | Suffix check | `endsWith(file, ".json")` |
| `now()` | `now() -> time` | Current timestamp | `now() > startTime` |
| `parseInt()` | `parseInt(string) -> number` | Parse integer | `parseInt("42") > 0` |
| `parseFloat()` | `parseFloat(string) -> number` | Parse float | `parseFloat("3.14")` |
| `toString()` | `toString(any) -> string` | Convert to string | `toString(123)` |
| `lower()` | `lower(string) -> string` | Lowercase | `lower(text) == "hello"` |
| `upper()` | `upper(string) -> string` | Uppercase | `upper(text) == "ERROR"` |
| `trim()` | `trim(string) -> string` | Remove whitespace | `trim(text) != ""` |

---

## Conclusion

This design specification provides a complete blueprint for implementing an enterprise-grade expression evaluation system in Thaiyyal. The system balances:

- **Power**: Complex expressions with node references, variables, and functions
- **Simplicity**: Clean syntax, easy to learn and use
- **Performance**: Compiled expressions, caching, minimal overhead
- **Security**: Safe evaluation, resource limits, no code injection
- **Compatibility**: Backward compatible with existing simple conditions
- **Maintainability**: Clean architecture, zero dependencies, well-tested

The phased implementation approach allows incremental delivery of value while maintaining quality and stability.

**Next Steps**:
1. Review and approve this design
2. Begin Phase 1 implementation (Core Infrastructure)
3. Set up project tracking for remaining phases
4. Schedule regular design reviews

**Questions or Feedback?** Please review and provide input on:
- Syntax preference (confirmed template-based?)
- Built-in functions to add/remove
- Security requirements
- Performance targets
- Implementation timeline

---

**Document History**:
- v1.0 (2025-11-01): Initial design specification
