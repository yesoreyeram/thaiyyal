# Phase 1: Expression Engine Enhancement - Implementation Plan

## Overview

**Goal**: Enhance the expression engine to fully support map, filter, and reduce operations with expressions, unblocking 40+ workflow examples.

**Duration**: 2 weeks  
**Impact**: +15% workflow coverage (from 60% to 75%)

## Current State Analysis

### What Works ✅
- ✅ Arithmetic operations (+, -, *, /, %)
- ✅ Comparison operators (==, !=, <, >, <=, >=)
- ✅ Logical operators (&&, ||, !)
- ✅ Field access (item.field, item.nested.field)
- ✅ Variable references (variables.name)
- ✅ Node references (node.id.value)
- ✅ Math functions (pow, sqrt, abs, floor, ceil, round, min, max)
- ✅ Integration with Map, Filter, Reduce executors

### What Needs Enhancement 🔧

1. **Array/Object Methods**
   - .length property
   - .includes() method
   - .startsWith() / .endsWith() for strings
   - .toUpperCase() / .toLowerCase() for strings
   - Array indexing: items[0], item.tags[0]

2. **Enhanced Expression Support in Executors**
   - Better error messages when expressions fail
   - Support for more complex nested expressions
   - Type checking and validation

3. **Testing & Documentation**
   - Comprehensive test coverage for all expression types
   - Update documentation with examples
   - Integration tests with workflow examples

## Implementation Tasks

### Task 1: Add Array/Object Methods (3 days)

#### 1.1 Array Length Property
**File**: `backend/pkg/expression/expression.go`

```go
// In resolveFieldPath, add special handling for .length
if field == "length" {
    if arr, ok := current.([]interface{}); ok {
        return float64(len(arr)), nil
    }
    if str, ok := current.(string); ok {
        return float64(len(str)), nil
    }
}
```

**Test Cases**:
- `items.length == 10`
- `item.name.length > 5`
- `item.tags.length >= 3`

#### 1.2 Array Indexing
**File**: `backend/pkg/expression/expression.go`

```go
// In resolveFieldPath, add support for array indexing: field[index]
// Example: "tags[0]", "users[1].name"
func parseArrayAccess(path string, obj interface{}) (interface{}, error) {
    // Parse: field[index] or field[index].nested
    // Use regex to extract field name and index
}
```

**Test Cases**:
- `items[0] == "first"`
- `item.tags[0].name == "important"`
- `users[item.index].active == true`

#### 1.3 String Methods
**File**: `backend/pkg/expression/expression.go`

Add function calls for string operations:
- `item.name.toUpperCase() == "ADMIN"`
- `item.email.toLowerCase().includes("@example.com")`
- `item.status.startsWith("active")`
- `item.filename.endsWith(".pdf")`

**Implementation**:
```go
// Add to evaluateFunctionCall or create new evaluateMethodCall
func evaluateMethodCall(obj interface{}, method string, args []interface{}) (interface{}, error) {
    switch method {
    case "toUpperCase":
        if str, ok := obj.(string); ok {
            return strings.ToUpper(str), nil
        }
    case "toLowerCase":
        if str, ok := obj.(string); ok {
            return strings.ToLower(str), nil
        }
    case "includes":
        if str, ok := obj.(string); ok && len(args) == 1 {
            needle := fmt.Sprintf("%v", args[0])
            return strings.Contains(str, needle), nil
        }
    // ... more methods
    }
}
```

#### 1.4 Array Methods
**File**: `backend/pkg/expression/expression.go`

- `items.includes(item)` - check if array contains value
- `items.some(expr)` - check if any item matches condition
- `items.every(expr)` - check if all items match condition

### Task 2: Enhanced Error Messages (2 days)

#### 2.1 Expression Error Context
**File**: `backend/pkg/expression/errors.go`

```go
type ExpressionError struct {
    Expression string
    Position   int
    Message    string
    Context    string
}

func (e *ExpressionError) Error() string {
    return fmt.Sprintf("expression error at position %d: %s\n%s\n%s^",
        e.Position, e.Message, e.Expression, strings.Repeat(" ", e.Position))
}
```

#### 2.2 Better Type Mismatch Messages
When operations fail due to type mismatches, provide clear error messages:
- "Cannot apply operator '+' to string and number"
- "Field 'age' not found in object"
- "Array index 5 out of bounds (length: 3)"

### Task 3: Comprehensive Testing (3 days)

#### 3.1 Expression Engine Tests
**File**: `backend/pkg/expression/expression_test.go`

Add test suites for:
- Array methods and indexing
- String methods
- Complex nested expressions
- Error cases with proper error messages
- Performance tests for complex expressions

#### 3.2 Integration Tests
**File**: `backend/pkg/executor/workflow_examples_test.go`

Un-skip and implement:
- TestWorkflowExample22_ArrayFiltering
- TestWorkflowExample23_DataTransformationPipeline
- Additional tests for examples 24-33

#### 3.3 End-to-End Workflow Tests
Create comprehensive tests that:
- Load actual workflow examples
- Execute them with mock data
- Verify expected outputs
- Test error handling

### Task 4: Update Map/Filter/Reduce Executors (2 days)

#### 4.1 Enhanced Map Executor
**File**: `backend/pkg/executor/control_map.go`

Improvements:
- Better error reporting with expression context
- Support for index-based expressions
- Metrics on transformation success/failure

#### 4.2 Enhanced Filter Executor
**File**: `backend/pkg/executor/control_filter.go`

Improvements:
- Support for complex filter expressions
- Better error reporting
- Performance optimization for large arrays

#### 4.3 Enhanced Reduce Executor
**File**: `backend/pkg/executor/control_reduce.go`

Improvements:
- Support for accumulator in expressions
- Better initial value handling
- Enhanced error messages

### Task 5: Documentation (2 days)

#### 5.1 Expression Syntax Guide
**File**: `docs/EXPRESSION_SYNTAX.md`

Document:
- All supported operators
- Field access patterns
- Array/object methods
- Function calls
- Examples for each feature

#### 5.2 Update Node Documentation
**File**: `docs/NODE_TYPES.md`

Update sections for:
- Map node with expression examples
- Filter node with condition examples
- Reduce node with accumulator examples

#### 5.3 Tutorial Examples
**File**: `docs/EXPRESSION_EXAMPLES.md`

Create practical examples:
- Filtering users by age and status
- Transforming product prices with calculations
- Aggregating sales data
- Complex multi-condition filters

## Implementation Schedule

### Week 1: Core Functionality
- **Days 1-3**: Implement array/object methods (Task 1)
  - Day 1: Array length and indexing
  - Day 2: String methods
  - Day 3: Array methods (includes, some, every)
  
- **Days 4-5**: Enhanced error messages (Task 2)
  - Day 4: Expression error context
  - Day 5: Type mismatch messages

### Week 2: Testing & Documentation
- **Days 1-3**: Comprehensive testing (Task 3)
  - Day 1: Expression engine tests
  - Day 2: Integration tests
  - Day 3: End-to-end workflow tests
  
- **Days 4-5**: Executor updates (Task 4)
  - Day 4: Map and Filter executors
  - Day 5: Reduce executor
  
- **Weekend**: Documentation (Task 5)

## Success Criteria

### Functional Requirements ✅
- [ ] All array/object methods implemented and tested
- [ ] Array indexing works correctly
- [ ] String methods functional
- [ ] Error messages are clear and actionable
- [ ] All skipped workflow example tests pass
- [ ] 40+ workflow examples now fully functional

### Quality Requirements ✅
- [ ] Test coverage > 90% for expression package
- [ ] All integration tests pass
- [ ] Performance: expressions evaluate in < 1ms average
- [ ] Documentation complete and accurate
- [ ] No security vulnerabilities (CodeQL clean)

### Workflow Coverage ✅
- [ ] Examples 21-33 (data processing) all pass
- [ ] Examples with map/filter/reduce all pass
- [ ] Coverage increases from 60% to 75%

## Testing Strategy

### Unit Tests
```go
// Test array length
func TestExpression_ArrayLength(t *testing.T) {
    ctx := &Context{
        Variables: map[string]interface{}{
            "items": []interface{}{1, 2, 3, 4, 5},
        },
    }
    
    result, err := Evaluate("variables.items.length == 5", nil, ctx)
    assert.NoError(t, err)
    assert.True(t, result)
}

// Test array indexing
func TestExpression_ArrayIndexing(t *testing.T) {
    ctx := &Context{
        Variables: map[string]interface{}{
            "users": []interface{}{
                map[string]interface{}{"name": "Alice"},
                map[string]interface{}{"name": "Bob"},
            },
        },
    }
    
    result, err := Evaluate("variables.users[0].name == 'Alice'", nil, ctx)
    assert.NoError(t, err)
    assert.True(t, result)
}
```

### Integration Tests
```go
func TestWorkflow_ComplexDataTransformation(t *testing.T) {
    // Test complete workflow with:
    // 1. Range to generate data
    // 2. Map to transform with expressions
    // 3. Filter with complex conditions
    // 4. Reduce to aggregate
}
```

## Risk Mitigation

### Risks & Mitigation

1. **Risk**: Breaking existing expressions
   - **Mitigation**: Comprehensive backwards compatibility tests
   - **Action**: Run all existing tests before changes

2. **Risk**: Performance degradation
   - **Mitigation**: Benchmark tests for common operations
   - **Action**: Performance testing required for all changes

3. **Risk**: Security vulnerabilities in expression evaluation
   - **Mitigation**: CodeQL scanning, input validation
   - **Action**: Security review for all expression parsing code

4. **Risk**: Incomplete error handling
   - **Mitigation**: Extensive error case testing
   - **Action**: Test all error paths

## Deliverables

1. ✅ Enhanced expression engine with all features
2. ✅ Updated Map, Filter, Reduce executors
3. ✅ Comprehensive test suite (>90% coverage)
4. ✅ Complete documentation
5. ✅ 40+ workflow examples now passing
6. ✅ Performance benchmarks
7. ✅ Security validation (CodeQL clean)

## Next Steps After Phase 1

Once Phase 1 is complete, proceed to:
- **Phase 2**: HTTP Enhancements (multipart, GraphQL, OAuth)
- **Phase 3**: Data Format Support (CSV, XML, YAML)
- **Phase 4**: New Node Types (RateLimiter, SchemaValidator, etc.)
- **Phase 5**: Integration Tests for all 150 examples

---

**Status**: Planning Complete - Ready for Implementation  
**Assigned**: @copilot  
**Timeline**: 2 weeks from approval  
**Dependencies**: None (self-contained)
