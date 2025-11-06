# Visual Workflow Guide: Age-Based API Routing

This document provides visual representations of the conditional branching workflow requested.

## Workflow Diagram

### Complete Workflow Structure

```
┌──────────────┐
│  User Age    │
│  (Number:25) │
│              │
└──────┬───────┘
       │
       │ Input value = 25
       ▼
┌──────────────┐
│  Age Check   │
│(Condition:   │
│   >= 18)     │
└──┬─────────┬─┘
   │         │
   │         └─────────────────────┐
   │                               │
   │ TRUE path                     │ FALSE path
   │ (Green handle)                │ (Red handle)
   │                               │
   ▼                               ▼
┌──────────────┐            ┌─────────────┐
│ Profile API  │            │ Education   │
│ (HTTP)       │            │ API (HTTP)  │
│              │            │             │
└──────┬───────┘            └─────────────┘
       │                          │
       │                          │
       ▼                          ▼
┌──────────────┐            ┌─────────────┐
│  Sports API  │            │   Render    │
│  (HTTP)      │            │   Result    │
│              │            └─────────────┘
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Render     │
│   Result     │
└──────────────┘
```

### Execution Flow for Age = 25 (Adult)

```
✅ EXECUTED PATH:

User Age (25) 
    ↓
Age Check (>=18) → TRUE ✓
    ↓
Profile API 
    ↓
Sports API
    ↓
Render Result

❌ SKIPPED:
Education API (not connected via true path)
```

**Nodes Executed:** 5 out of 6 total nodes
**Performance:** 16.7% reduction in node execution

### Execution Flow for Age = 15 (Minor)

```
✅ EXECUTED PATH:

User Age (15)
    ↓
Age Check (>=18) → FALSE ✗
    ↓
Education API
    ↓
Render Result

❌ SKIPPED (Transitive):
Profile API (not connected via false path)
Sports API (parent Profile API skipped)
```

**Nodes Executed:** 3 out of 6 total nodes  
**Performance:** 50% reduction in node execution

---

## Node Visual Reference

### 1. Number Node (User Age)

```
┌─────────────────────┐
│ 📝 User Age         │ ← Node title (editable)
├─────────────────────┤
│ [    25      ]  🟢 │ ← Value input | Source handle (right)
└─────────────────────┘
   ^
   │
   Target handle (left, blue)
```

**Configuration:**
- **Value:** 25 (or any age to test)
- **Type:** number
- **Output:** The number value (25)

### 2. Condition Node (Age Check)

```
┌─────────────────────┐
│ 🚦 Age Check        │ ← Node title
├─────────────────────┤
│ 🔵 [  >= 18  ] 🟢  │ ← Input | Condition | TRUE handle (green, top-right)
│                 🔴  │ ← FALSE handle (red, bottom-right)
└─────────────────────┘
   ^
   │
   Target handle (left, blue)
```

**Configuration:**
- **Condition:** `>=18` (checks if input >= 18)
- **Type:** condition
- **Output:** 
  - Metadata with `path: "true"` if condition true
  - Metadata with `path: "false"` if condition false

**Handles:**
- **Green (top-right):** Connects to nodes that run when TRUE
- **Red (bottom-right):** Connects to nodes that run when FALSE

### 3. HTTP Node (Profile API, Sports API, Education API)

```
┌─────────────────────┐
│ 🌐 Profile API      │ ← Node title
├─────────────────────┤
│ 🔵 [URL input ] 🟢 │ ← Target | URL field | Source
└─────────────────────┘
```

**Configuration:**
- **URL:** API endpoint (e.g., `https://api.example.com/profile`)
- **Method:** GET/POST/etc.
- **Type:** http
- **Output:** API response data

### 4. Renderer Node (Display Results)

```
┌─────────────────────┐
│ 📺 Adult Result     │ ← Node title
├─────────────────────┤
│ 🔵 [Display Area]   │ ← Shows the input data
└─────────────────────┘
   ^
   │
   Target handle (left, blue)
```

**Purpose:** Displays the workflow output for debugging/verification

---

## Step-by-Step Visual Connection Guide

### Step 1-3: Add and Configure Age Input

```
Action: Drag "Number" node from palette

Before:                    After:
[Empty Canvas]      →      ┌──────────────┐
                           │ User Age     │
                           │ Value: 25    │
                           └──────────────┘
```

### Step 4-5: Add and Configure Condition

```
Action: Drag "Condition" node, configure condition

Before:                    After:
┌──────────────┐          ┌──────────────┐     ┌──────────────┐
│ User Age     │    →     │ User Age     │     │ Age Check    │
│ Value: 25    │          │ Value: 25    │     │ Cond: >=18   │
└──────────────┘          └──────────────┘     └──────────────┘
```

### Step 6: Connect Age to Condition

```
Action: Drag from Age's RIGHT handle to Condition's LEFT handle

Before:                            After:
┌──────────────┐  ┌──────────────┐       ┌──────────────┐      ┌──────────────┐
│ User Age  🟢 │  │🔵 Age Check  │  →    │ User Age  🟢─────→🔵 Age Check    │
│ Value: 25    │  │   Cond: >=18 │       │ Value: 25    │      │   Cond: >=18 │
└──────────────┘  └──────────────┘       └──────────────┘      └──────────────┘
```

### Step 7-8: Add Profile API and Connect TRUE Path

```
Action: Drag from Age Check's GREEN handle to Profile API's LEFT handle

CRITICAL: Use the GREEN handle (top-right) for TRUE path!

Before:                                    After:
┌──────────────┐      ┌──────────────┐          ┌──────────────┐      ┌──────────────┐
│ User Age  🟢─────→🔵 Age Check  🟢│          │ User Age  🟢─────→🔵 Age Check  🟢──┐
│ Value: 25    │      │   Cond: >=18 🔴│   →      │ Value: 25    │      │   Cond: >=18 🔴│ │
└──────────────┘      └──────────────┘          └──────────────┘      └──────────────┘ │
                                                                                         │
                      ┌──────────────┐                              ┌──────────────┐    │
                      │🔵Profile API │                              │🔵Profile API │←───┘
                      │              │                              │              │
                      └──────────────┘                              └──────────────┘
```

**Note:** The green handle ensures Profile API only runs when age >= 18.

### Step 9: Add Sports API and Chain

```
Action: Connect Profile API → Sports API (regular connection)

Before:                                          After:
┌──────────────┐      ┌──────────────┐                 ┌──────────────┐      ┌──────────────┐
│ User Age  🟢─────→🔵 Age Check  🟢──┐                │ User Age  🟢─────→🔵 Age Check  🟢──┐
│ Value: 25    │      │   Cond: >=18 🔴│ │               │ Value: 25    │      │   Cond: >=18 🔴│ │
└──────────────┘      └──────────────┘ │               └──────────────┘      └──────────────┘ │
                                        │                                                       │
                      ┌──────────────┐ │  ┌──────────┐       ┌──────────────┐                │
                      │🔵Profile API │←┘  │Sports API│       │🔵Profile API │←───────────────┘
                      │           🟢 │    │          │  →    │           🟢───────→🔵Sports API│
                      └──────────────┘    └──────────┘       └──────────────┘        │          │
                                                                                      └──────────┘
```

### Step 10: Add Education API and Connect FALSE Path

```
Action: Drag from Age Check's RED handle to Education API's LEFT handle

CRITICAL: Use the RED handle (bottom-right) for FALSE path!

Final Workflow:

                      ┌──────────────┐      ┌──────────────┐
                      │ User Age  🟢─────→🔵 Age Check  🟢──────┐
                      │ Value: 25    │      │   Cond: >=18 🔴─┐ │
                      └──────────────┘      └──────────────┘ │ │
                                                              │ │
                                         TRUE path (green) ──┘ └── FALSE path (red)
                                                              │   │
                                    ┌──────────────┐          │   │    ┌──────────────┐
                                    │🔵Profile API │←─────────┘   └───→│🔵Education   │
                                    │           🟢─────┐                │    API       │
                                    └──────────────┘   │                └──────────────┘
                                                       │
                                              ┌────────▼────┐
                                              │🔵Sports API │
                                              │             │
                                              └─────────────┘
```

---

## Handle Connection Reference

### Condition Node Handles

```
                    ┌─────────────────────┐
                    │ 🚦 Condition Node   │
                    ├─────────────────────┤
     Input ────────►│ 🔵 [condition] 🟢  │◄──── TRUE output handle
     (Blue, Left)   │                 🔴  │◄──── FALSE output handle
                    └─────────────────────┘
                           30% from top
                           70% from top
```

**Connection Rules:**
1. **Input (Left, Blue):** Connect from any source node that provides a value
2. **TRUE Output (Right-Top, Green ~30%):** Connect to nodes that should run when condition is TRUE
3. **FALSE Output (Right-Bottom, Red ~70%):** Connect to nodes that should run when condition is FALSE

### Regular Node Handles

```
     ┌─────────────────────┐
     │ 📄 Regular Node     │
     ├─────────────────────┤
────►│ 🔵      [data]  🟢  │◄──── Output handle
     └─────────────────────┘      (Green, Right)
          │
          └──── Input handle
               (Blue, Left)
```

---

## Color Coding Guide

### Handle Colors

- **🔵 Blue (Target/Input):** Where connections come IN to a node
- **🟢 Green (Source/Output):** Where connections go OUT from a node
- **🟢 Green (Condition TRUE):** Specific to condition's true path
- **🔴 Red (Condition FALSE):** Specific to condition's false path

### Node Colors (in UI)

- **Number nodes:** Blue tint
- **HTTP nodes:** Purple tint  
- **Condition nodes:** Green/Emerald background
- **Renderer nodes:** Pink tint
- **Switch nodes:** Default styling

### Execution States

- **✅ Green highlight:** Node executed successfully
- **❌ Red/grayed out:** Node was skipped (not executed)
- **⚠️ Yellow:** Node encountered an error
- **🔵 Blue outline:** Node is selected

---

## Verification Checklist

After building your workflow, verify:

- [ ] Age input node is configured with a test value
- [ ] Condition node has the correct expression (`>=18`)
- [ ] Profile API is connected to the GREEN (true) handle
- [ ] Education API is connected to the RED (false) handle
- [ ] Sports API is connected downstream from Profile API
- [ ] No loose/dangling connections
- [ ] All nodes are properly labeled

**Test Execution:**

- [ ] Test with age >= 18: Profile API + Sports API execute, Education API skipped
- [ ] Test with age < 18: Education API executes, Profile + Sports APIs skipped
- [ ] Execution count matches expected (4 vs 3 nodes)
- [ ] Results appear in the correct renderer nodes

---

## Common Visual Patterns

### Pattern 1: Simple Binary Branch

```
    Input
      ↓
  Condition
   /     \
True    False
  ↓       ↓
Path A  Path B
```

### Pattern 2: Chained TRUE Path, Simple FALSE

```
    Input
      ↓
  Condition
   /     \
True    False
  ↓       ↓
Node A  Node B
  ↓
Node C
  ↓
Node D
```

### Pattern 3: Nested Conditions

```
    Input
      ↓
 Condition 1
   /     \
True    False
  ↓       ↓
Cond 2  Path C
 /   \
T     F
↓     ↓
A     B
```

### Pattern 4: Multiple Conditions Converging

```
Input A  Input B
   ↓        ↓
Cond A   Cond B
   ↓        ↓
   └───┬────┘
       ↓
   Join Node
       ↓
    Output
```

---

## Keyboard Shortcuts (Common in React Flow)

- **Ctrl + Mouse Wheel:** Zoom in/out
- **Space + Drag:** Pan the canvas
- **Delete/Backspace:** Delete selected node/edge
- **Ctrl + C / Ctrl + V:** Copy/paste nodes
- **Ctrl + Z / Ctrl + Y:** Undo/redo
- **Ctrl + A:** Select all nodes
- **Click + Drag:** Select multiple nodes (box select)

---

## Next Steps

1. **Build the workflow** following the visual guide above
2. **Test both paths** (age >= 18 and age < 18)
3. **Verify execution** in the execution panel
4. **Experiment** with different condition expressions
5. **Try nested conditions** for more complex logic
6. **Explore switch nodes** for multi-way branching

---

## Additional Visual Resources

See also:
- `examples/conditional-branching/09-age-based-api-routing.json` - The actual JSON for this workflow
- `docs/UI_GUIDE_CONDITIONAL_BRANCHING.md` - Detailed step-by-step text guide
- `docs/CONDITIONAL_EXECUTION_DEMO.md` - Execution logs and proof
- `docs/CONDITIONAL_EXECUTION_IMPLEMENTATION.md` - Technical implementation details
