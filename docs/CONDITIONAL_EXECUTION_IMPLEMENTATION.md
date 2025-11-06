# Conditional Execution Implementation Summary

## Overview

This document summarizes the implementation of **conditional execution** (path termination) in Thaiyyal, enabling workflows to execute only the nodes in the active branch based on runtime conditions.

## User Request

> "With the current implementation of the workflow engine execution, how do we terminate the workflow path based on condition? Say for example if the user age is above 18, i want to call the user profile api to list their interests and based on that i will call another api for example sports registration api. If the user age is less than 18, i will call some education registration related api."

## Implementation

### Backend Changes

#### 1. Edge Type Enhancement (`backend/pkg/types/types.go`)

```go
type Edge struct {
    ID           string  `json:"id"`
    Source       string  `json:"source"`
    Target       string  `json:"target"`
    SourceHandle *string `json:"sourceHandle,omitempty"` // NEW: Output port (e.g., "true", "false", "success")
    TargetHandle *string `json:"targetHandle,omitempty"` // NEW: Input port
    Condition    *string `json:"condition,omitempty"`    // DEPRECATED: Backward compatibility
}
```

**Key Features**:
- `sourceHandle`: Specifies which output path from the source node (e.g., "true", "false", "success", "error")
- Backward compatible with legacy `condition` field
- React Flow automatically populates these fields when connecting to/from specific handles

#### 2. Engine Execution Logic (`backend/pkg/engine/engine.go`)

**New Methods**:

1. **`shouldExecuteNode(nodeID string) bool`**
   - Checks all incoming edges to the node
   - Returns `true` if node should execute:
     - No incoming edges (orphan/start node)
     - Has unconditional incoming edges
     - At least one conditional edge's condition is satisfied
   - Returns `false` if all edges are conditional and none are satisfied

2. **`getIncomingEdges(nodeID string) []Edge`**
   - Returns all edges targeting a specific node

3. **`isConditionSatisfied(sourceResult interface{}, condition string) bool`**
   - Evaluates if an edge condition is satisfied
   - Supports:
     - `"true"` / `"false"` for condition nodes (checks `path`, `true_path`, `condition_met` fields)
     - Custom paths for switch nodes (checks `output_path` field)

**Modified Execution Loop**:
```go
for _, nodeID := range executionOrder {
    // NEW: Check if node should execute based on conditional edges
    if !e.shouldExecuteNode(nodeID) {
        e.structuredLogger.WithNodeID(nodeID).Debug("node skipped due to conditional edge")
        continue  // Skip this node!
    }
    
    // Execute node (existing code)
    node := e.getNode(nodeID)
    value, err := e.executeNode(ctx, node)
    // ...
}
```

### Frontend Compatibility

**React Flow Integration** (Already Working):

The frontend uses React Flow which automatically handles multiple output handles:

```tsx
// Condition node already has true/false handles
<Handle type="source" id="true" position={Position.Right} style={{top: "30%"}} />
<Handle type="source" id="false" position={Position.Right} style={{top: "70%"}} />
```

When users connect edges in the UI, React Flow's `onConnect` callback receives:
```typescript
{
  source: "node1",
  target: "node2",
  sourceHandle: "true",  // Automatically included!
  targetHandle: null
}
```

**No frontend changes needed** - React Flow's built-in handle system works perfectly!

### Test Suite

**7 Comprehensive Tests** (`backend/pkg/engine/conditional_execution_test.go`):

1. ✅ `TestConditionalExecution_TruePathOnly` - Adult path only
2. ✅ `TestConditionalExecution_FalsePathOnly` - Minor path only
3. ✅ `TestConditionalExecution_SwitchRouting` - HTTP status routing
4. ✅ `TestConditionalExecution_NestedConditions` - Multi-level branches
5. ✅ `TestConditionalExecution_MultipleConditionalEdges` - OR logic (any condition satisfied)
6. ✅ `TestConditionalExecution_UnconditionalEdgeTakesPrecedence` - Mixed edges
7. ✅ `TestConditionalExecution_BackwardCompatibility` - Legacy condition field

**All tests passing**: 7/7 ✅

### Example Workflows

**3 New Real-World Examples**:

#### 1. Age-Based API Routing (`09-age-based-api-routing.json`)
Demonstrates the **exact user scenario**:
- If age >= 18: Fetch profile API → Sports registration
- If age < 18: Education registration
- **Result**: Only one path executes

#### 2. Multi-Step Registration (`10-multi-step-registration.json`)
Complex workflow with multiple steps per branch:
- Adult path: Parse → Fetch profile → Extract interests → Register programs
- Minor path: Direct education registration
- Both paths converge to confirmation step

#### 3. HTTP Status Routing (`11-http-status-routing.json`)
Switch-based conditional execution:
- 200 → Success handler
- 404 → Retry handler
- 500+ → Error handler

### Documentation

**Comprehensive Documentation Added** to `examples/conditional-branching/README.md`:

1. **Conditional Execution Overview**
   - How it works
   - Engine behavior
   - Path skipping logic

2. **Visual Guide**
   - Creating conditional workflows in UI
   - Connecting true/false paths
   - Switch-based workflows

3. **API Documentation**
   - Edge schema
   - Condition node handles
   - Switch node handles
   - Backend processing flow

4. **Migration Guide**
   - From metadata-only to conditional execution
   - Backward compatibility notes

5. **Best Practices**
   - When to use conditional execution
   - Path convergence patterns
   - Naming conventions

6. **Troubleshooting**
   - Common issues and solutions
   - Performance notes

## How It Works: Step-by-Step

### Example: Age-Based User Registration

**Workflow**:
```
[Age Input: 25] → [Condition: >=18] ─┬─[True]──→ [Fetch Profile] → [Sports Registration]
                                      │
                                      └─[False]─→ [Education Registration]
```

**Execution Flow**:

1. **Parse Workflow**: Engine loads nodes and edges
2. **Topological Sort**: Determines execution order: `[age, condition, fetch_profile, register_sports, register_education]`
3. **Execute Age Node**: Result = `25`
4. **Execute Condition Node**: Evaluates `25 >= 18` → Result = `{path: "true", condition_met: true, value: 25}`
5. **Check Fetch Profile**:
   - Incoming edge has `sourceHandle: "true"`
   - Source (condition) result has `path: "true"` ✅
   - **Execute fetch_profile**
6. **Check Register Sports**:
   - Incoming edge is unconditional (no sourceHandle)
   - **Execute register_sports**
7. **Check Education Registration**:
   - Incoming edge has `sourceHandle: "false"`
   - Source (condition) result has `path: "true"` ❌
   - **SKIP education registration** ✅

**Result**: Only adult path executes!

## Key Features

### ✅ Conditional Path Execution
- Nodes execute only if incoming edge conditions are satisfied
- Inactive paths are completely skipped
- Improves performance by avoiding unnecessary work

### ✅ Multiple Output Handles
- Condition nodes: `"true"` and `"false"` handles
- Switch nodes: Custom handles per case (`"success"`, `"error"`, etc.)
- React Flow automatically manages handle connections

### ✅ Flexible Routing Logic
- **OR Logic**: Node executes if ANY incoming edge condition is satisfied
- **Unconditional Precedence**: Unconditional edges always allow execution
- **Mixed Edges**: Combine conditional and unconditional as needed

### ✅ Backward Compatible
- Workflows without `sourceHandle` work as before (all nodes execute)
- Legacy `condition` field still supported
- No breaking changes to existing workflows

### ✅ Comprehensive Testing
- 7 backend tests covering all scenarios
- Real-world example workflows
- Edge cases handled (multiple edges, nested conditions, etc.)

## Usage Examples

### Frontend (React Flow UI)

Users can create conditional workflows visually:

1. **Drag Condition node** from palette
2. **Set condition**: Type `">=18"` in condition field
3. **Connect true path**: Click green handle → drag to "adult" node
4. **Connect false path**: Click red handle → drag to "minor" node
5. **Run workflow**: Only active path executes!

### Backend (JSON Payload)

```json
{
  "nodes": [
    {"id": "age", "type": "number", "data": {"value": 25}},
    {"id": "check", "type": "condition", "data": {"condition": ">=18"}},
    {"id": "adult", "type": "text_input", "data": {"text": "Adult profile"}},
    {"id": "minor", "type": "text_input", "data": {"text": "Minor education"}}
  ],
  "edges": [
    {"source": "age", "target": "check"},
    {"source": "check", "target": "adult", "sourceHandle": "true"},
    {"source": "check", "target": "minor", "sourceHandle": "false"}
  ]
}
```

**Execute**:
```bash
curl -X POST http://localhost:8080/api/workflow/execute \
  -H "Content-Type: application/json" \
  -d @workflow.json
```

**Result**:
```json
{
  "execution_id": "abc123",
  "node_results": {
    "age": 25,
    "check": {"path": "true", "condition_met": true, "value": 25},
    "adult": "Adult profile"
    // "minor" node NOT in results - it was skipped!
  }
}
```

## Performance Impact

- **Positive**: Skipping nodes improves performance (no unnecessary API calls, computations)
- **Minimal Overhead**: Condition checking adds microseconds per node
- **Recommended**: Use for expensive operations (HTTP calls, heavy processing)

## Future Enhancements

Potential improvements identified:

1. **Visual Path Highlighting**: Show active path in UI during/after execution
2. **Path Analytics**: Track which paths execute most frequently
3. **Conditional Loops**: While loops with exit conditions
4. **Parallel Conditional Branches**: Execute multiple paths in parallel
5. **Path Variables**: Pass different data to different paths

## Migration Notes

### For Existing Workflows

✅ **No changes required** - existing workflows continue to work
- Edges without `sourceHandle` are unconditional (all nodes execute as before)
- Add `sourceHandle` to enable conditional execution

### For New Workflows

✅ **Use conditional execution** for:
- Mutually exclusive operations (either A or B, not both)
- Error handling (try → success/failure paths)
- Feature flags (enabled → new feature, disabled → old feature)
- Multi-tenant routing (premium → advanced features, standard → basic features)

## Testing

### Run Tests

```bash
# Run all conditional execution tests
cd backend
go test -v ./pkg/engine -run TestConditionalExecution

# Expected output:
# PASS: TestConditionalExecution_TruePathOnly
# PASS: TestConditionalExecution_FalsePathOnly
# PASS: TestConditionalExecution_SwitchRouting
# PASS: TestConditionalExecution_NestedConditions
# PASS: TestConditionalExecution_MultipleConditionalEdges
# PASS: TestConditionalExecution_UnconditionalEdgeTakesPrecedence
# PASS: TestConditionalExecution_BackwardCompatibility
# ok   github.com/yesoreyeram/thaiyyal/backend/pkg/engine
```

### Try Examples

```bash
# Load example workflow
curl http://localhost:8080/api/workflow/load \
  -d @examples/conditional-branching/09-age-based-api-routing.json

# Change age to 15 and see different path execute
# Edit JSON: "value": 15
# Both paths won't execute - only the active one!
```

## Files Changed

### Backend
- ✅ `backend/pkg/types/types.go` - Edge type enhancement
- ✅ `backend/pkg/engine/engine.go` - Conditional execution logic (+120 lines)
- ✅ `backend/pkg/engine/conditional_execution_test.go` - 7 new tests (+270 lines)

### Frontend
- ✅ No changes needed - React Flow handles already support this!

### Examples
- ✅ `examples/conditional-branching/09-age-based-api-routing.json` - User's scenario
- ✅ `examples/conditional-branching/10-multi-step-registration.json` - Complex flow
- ✅ `examples/conditional-branching/11-http-status-routing.json` - Switch routing

### Documentation
- ✅ `examples/conditional-branching/README.md` - (+150 lines)
- ✅ `docs/CONDITIONAL_BRANCHING_TESTING_SUMMARY.md` - Updated
- ✅ `docs/CONDITIONAL_EXECUTION_IMPLEMENTATION.md` - This document

## Summary

**✅ IMPLEMENTATION COMPLETE**

The workflow engine now supports **full conditional execution** as requested:

- ✅ Nodes execute only in the active path
- ✅ Inactive paths are skipped entirely
- ✅ Condition nodes provide true/false routing
- ✅ Switch nodes provide multi-way routing
- ✅ Frontend UI already supports this (React Flow handles)
- ✅ Backward compatible with existing workflows
- ✅ Comprehensive test coverage (7/7 tests passing)
- ✅ Real-world examples demonstrating the feature
- ✅ Extensive documentation with visual guides

**User's scenario is now fully supported**:
- Age >= 18: Profile API → Sports registration ✅
- Age < 18: Education registration ✅
- Only one path executes based on runtime condition ✅

---

**Implementation Date**: 2025-11-06  
**Tests Passing**: 7/7 ✅  
**Examples Created**: 3 new workflows  
**Documentation**: 150+ lines added  
**Breaking Changes**: None (backward compatible)
