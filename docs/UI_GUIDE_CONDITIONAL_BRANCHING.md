# UI Guide: Building Conditional Branching Workflows

This guide provides detailed, step-by-step instructions for creating conditional branching workflows using the Thaiyyal visual workflow editor.

## Table of Contents

1. [Quick Overview](#quick-overview)
2. [Example: Age-Based API Routing](#example-age-based-api-routing)
3. [Understanding Conditional Nodes](#understanding-conditional-nodes)
4. [Tips and Tricks](#tips-and-tricks)
5. [Troubleshooting](#troubleshooting)

---

## Quick Overview

**What You'll Learn:**
- How to use the Condition node with multiple output paths
- How to connect specific handles (true/false paths)
- How to build workflows where only the active branch executes

**Prerequisites:**
- Basic familiarity with the Thaiyyal UI
- Understanding of workflow concepts (nodes, edges, execution flow)

---

## Example: Age-Based API Routing

### Scenario

**Goal:** Create a workflow that:
- If age >= 18: Call profile API → Sports registration API
- If age < 18: Call education registration API
- Only execute the nodes in the active path (skip inactive nodes)

### Step-by-Step Instructions

#### Step 1: Open the Node Palette

1. **Locate the left sidebar** - Look for the panel on the left side of the screen
2. **Find the toggle button** - If the palette is closed, look for a "hamburger menu" icon (≡) or "Show Nodes" button in the top-left corner
3. **Click to open** - The Node Palette should slide out, showing categories like "Input/Output", "Operations", "Control Flow", etc.

> **Visual cue:** The palette has a dark background with white text saying "Nodes" at the top

#### Step 2: Add the Age Input Node

1. **Find the "Input/Output" category** in the Node Palette
2. **Click to expand** if it's collapsed (look for a down arrow ▼)
3. **Locate "Number" node** in the list - it's usually the first item
4. **Add the node** - You have 2 options:
   - **Option A (Click):** Click on "Number" once - it will be added to the canvas center
   - **Option B (Drag):** Click and hold on "Number", drag to your desired position, then release
5. **Position the node** on the left side of the canvas (this will be your starting point)

> **Tip:** Use drag method for precise placement. The node shows a dragging cursor (⋮⋮) icon.

#### Step 3: Configure the Age Value

1. **Click on the Number node** you just added (if it's not already selected)
2. **Find the number input field** inside the node - it should show "0" by default
3. **Click in the input field** and type your test age value (e.g., `25` for adult or `15` for minor)
4. **Rename the node (optional):**
   - Click on the node title "Node X" at the top
   - Type a descriptive name like "User Age" or "Age Input"
   - Press Enter to save

> **Note:** The title is editable when you click on it - it will show a text cursor.

#### Step 4: Add the Condition Node

1. **Scroll down in the Node Palette** to find the **"Control Flow"** category
2. **Click to expand** the category if needed
3. **Locate "Condition" node** - it should have an amber/yellow color indicator
4. **Add the node** by clicking or dragging it onto the canvas
5. **Position it to the right** of your Number node (leave some space for the connection)

> **Visual cue:** The Condition node has a green background and shows ">0" as the default condition

#### Step 5: Configure the Age Check Condition

1. **Click on the Condition node** to focus it
2. **Find the condition input field** - it shows ">0" by default
3. **Click in the field** and replace the text with: `>=18`
   - This means "check if the input value is greater than or equal to 18"
4. **Rename the node:**
   - Click the title "Condition" at the top
   - Type "Age Check" or "Adult Check"
   - Press Enter

> **Important:** The condition is evaluated against the input value from the connected node. So `>=18` checks if `age >= 18`.

#### Step 6: Connect Age to Condition

1. **Hover over the Number node** - you'll see small circular handles appear
2. **Locate the RIGHT handle** (source handle) - it's on the right edge of the node, colored green
3. **Click and hold** on the right handle of the Number node
4. **Drag toward the Condition node** - you'll see a connection line following your cursor
5. **Hover over the Condition node** - its LEFT handle (target handle) will highlight in blue
6. **Release the mouse** over the blue target handle
7. **Verify the connection** - you should see a solid line connecting Number → Condition

> **Visual feedback:** The connection line appears when you drag, and snaps into place when you release over a valid target handle.

#### Step 7: Add Profile API Node (Adult Path)

1. **Go back to the Node Palette**
2. **Expand "Input/Output" category**
3. **Locate "HTTP" node**
4. **Add it to the canvas** - position it to the right and slightly above the Condition node
5. **Rename it:**
   - Click the title "HTTP"
   - Type "Profile API"
   - Press Enter
6. **Configure the API (optional):**
   - You can click on the node to see configuration options
   - For testing, you can leave the URL field empty or use a placeholder like `https://api.example.com/profile`

> **Note:** For this demonstration, the actual HTTP call configuration is secondary - focus on the conditional routing structure.

#### Step 8: Connect Condition TRUE Path to Profile API

**This is the critical step for conditional execution!**

1. **Hover over the Condition node** - you'll see **TWO handles on the right side**:
   - **Top-right handle (green)** - labeled "True path" when you hover
   - **Bottom-right handle (red)** - labeled "False path" when you hover
2. **Click and hold the GREEN (true) handle** at the top-right
3. **Drag toward the Profile API node**
4. **Release over the Profile API's LEFT (target) handle**
5. **Verify:** The connection line should originate from the GREEN handle

> **Critical:** Make sure you connect from the GREEN handle, not the red one! This ensures the Profile API only executes when age >= 18.

#### Step 9: Add Sports API Node

1. **Add another HTTP node** from the Node Palette
2. **Position it to the right** of the Profile API node (creating a chain)
3. **Rename it to "Sports API"**
4. **Connect Profile API → Sports API:**
   - Drag from the RIGHT handle of Profile API
   - Release on the LEFT handle of Sports API

> **Result:** Now you have a chain: Age → Condition → (true path) → Profile API → Sports API

#### Step 10: Add Education API Node (Minor Path)

1. **Add another HTTP node** from the Node Palette
2. **Position it to the right** and **below** the Condition node (parallel to the Profile API chain)
3. **Rename it to "Education API"**
4. **Connect Condition FALSE path to Education API:**
   - **Important:** Click and hold the **RED (false) handle** at the bottom-right of the Condition node
   - Drag to the Education API node
   - Release on its LEFT target handle
5. **Verify:** The connection should come from the RED handle

> **Result:** You now have TWO branches:
> - True path (green): Condition → Profile API → Sports API  
> - False path (red): Condition → Education API

#### Step 11: Add Output Nodes (Optional)

To see the results, you can add Renderer nodes:

1. **Add a "Renderer" node** from the Input/Output category
2. **Position it to the right** of Sports API
3. **Connect:** Sports API → Renderer (name it "Adult Result")
4. **Add another Renderer node**
5. **Position it to the right** of Education API
6. **Connect:** Education API → Renderer (name it "Minor Result")

> **Purpose:** Renderer nodes display the output, helping you verify which path executed.

#### Step 12: Test the Workflow

**Test Case 1: Adult (age >= 18)**

1. **Set the age to 25** (or any value >= 18) in the Number node
2. **Click the "Execute" or "Run" button** (usually in the top toolbar)
3. **Observe the results:**
   - ✅ Profile API should execute
   - ✅ Sports API should execute
   - ✅ "Adult Result" renderer shows output
   - ❌ Education API should be **skipped** (grayed out or not executed)
   - ❌ "Minor Result" renderer shows nothing

**Test Case 2: Minor (age < 18)**

1. **Change the age to 15** (or any value < 18) in the Number node
2. **Click "Execute" again**
3. **Observe the results:**
   - ✅ Education API should execute
   - ✅ "Minor Result" renderer shows output
   - ❌ Profile API should be **skipped**
   - ❌ Sports API should be **skipped** (transitive skip!)
   - ❌ "Adult Result" renderer shows nothing

> **Expected behavior:** Only ONE path executes based on the condition. The inactive path is completely skipped, including downstream nodes (Sports API when Education API runs).

#### Step 13: View Execution Results

1. **Look for the Execution Panel** - usually on the right side or bottom of the screen
2. **Expand it** if it's minimized
3. **Check the execution log:**
   - You should see which nodes executed
   - Skipped nodes may be marked or omitted from the log
4. **Verify node count:**
   - Adult path (age=25): 4 nodes execute (Age, Condition, Profile API, Sports API)
   - Minor path (age=15): 3 nodes execute (Age, Condition, Education API)

---

## Understanding Conditional Nodes

### Condition Node

**Visual Features:**
- Green/emerald background color
- Two output handles on the right:
  - 🟢 **Green handle (top)** - True path
  - 🔴 **Red handle (bottom)** - False path
- One input handle on the left (blue)
- Condition input field in the middle

**How It Works:**
1. Receives input from connected nodes
2. Evaluates the condition expression (e.g., `>=18`, `>0`, `==true`)
3. Outputs metadata indicating which path was taken ("true" or "false")
4. Only nodes connected via the satisfied handle execute

**Condition Syntax:**
```
>=18          // Greater than or equal to 18
<100          // Less than 100
==true        // Equal to true
!=null        // Not equal to null
>0 && <100    // Between 0 and 100 (AND logic)
<0 || >100    // Less than 0 OR greater than 100 (OR logic)
```

### Switch Node

**Visual Features:**
- Default background color (varies)
- One or more output handles (depending on cases)
- One input handle on the left
- Configuration for cases and default path

**How It Works:**
1. Receives input value
2. Evaluates each case condition in order
3. Takes the first matching case path
4. If no case matches, takes the default path
5. Only nodes on the matched path execute

**Use Cases:**
- HTTP status code routing (200, 404, 500, etc.)
- User role-based routing (admin, user, guest)
- Multi-way branching (A/B/C testing)

---

## Tips and Tricks

### Connecting Handles

1. **Always connect FROM the correct source handle:**
   - For Condition: Use green (true) or red (false)
   - For regular nodes: Use the single green handle
2. **Hover to see handle labels** - they show "True path", "False path", etc.
3. **Connection line color** doesn't indicate the condition - check the source handle position
4. **Zoom in** if you have trouble clicking small handles (Ctrl + Mouse Wheel)

### Organizing Your Workflow

1. **Arrange branches vertically:**
   - True path: Top branch
   - False path: Bottom branch
2. **Use consistent spacing** - leave room between parallel paths
3. **Align nodes** - select multiple nodes and use alignment tools if available
4. **Add labels** - rename nodes to describe their purpose (e.g., "Adult Check", "Sports API")

### Testing

1. **Test both paths** - change the input value to test true and false cases
2. **Check execution counts** - verify only the expected nodes execute
3. **Use Renderer nodes** - add them at branch endpoints to see outputs
4. **Check the execution panel** - review which nodes ran and their outputs

### Debugging

1. **Connection issues:**
   - If a path doesn't execute, check that you connected to the correct handle (green vs red)
   - Verify the condition syntax (e.g., `>=18` not `= 18`)
2. **Both paths executing:**
   - Old workflows without sourceHandle will execute all nodes
   - Ensure you connected from the colored handles (green/red), not just clicking between nodes
3. **Neither path executing:**
   - Check that the Condition node receives input
   - Verify the condition expression is valid

---

## Troubleshooting

### Problem: Both branches execute regardless of condition

**Solution:**
- Check that your edges are connected from the **specific handles** (green or red)
- Edges connected by clicking between nodes (auto-connect) may not have sourceHandle set
- Reconnect by explicitly dragging from the green/red handles

### Problem: Condition node has no colored handles

**Solution:**
- Ensure you're using the latest version of the UI
- The Condition node should automatically show green and red handles
- Try refreshing the page or clearing your browser cache

### Problem: Can't see which handle to connect from

**Solution:**
- **Zoom in** using Ctrl + Mouse Wheel (or trackpad pinch)
- **Hover over the handle** - a tooltip should appear showing "True path" or "False path"
- The green handle is at ~30% from top, red handle at ~70% from top on the right edge

### Problem: Don't know which path was taken

**Solution:**
- Add **Renderer nodes** at the end of each branch
- Check the **Execution Panel** for execution logs
- Look at the **node highlighting** during execution (if supported)

### Problem: Downstream nodes in inactive path still execute

**Solution:**
- This was a bug in earlier versions, fixed in commit a393f78
- Ensure you're using the latest version of the backend
- Inactive paths should be completely skipped, including transitive dependencies

---

## Advanced Patterns

### Nested Conditions

You can chain multiple Condition nodes:

```
Age Check (>=18)
├─ True → Country Check (=="US")
│         ├─ True → Special Offer
│         └─ False → Standard Offer
└─ False → Parental Consent
```

**Steps:**
1. Create first Condition (Age Check)
2. Connect TRUE path to second Condition (Country Check)
3. Connect Country Check TRUE → Special Offer
4. Connect Country Check FALSE → Standard Offer
5. Connect Age Check FALSE → Parental Consent

### Switch with Multiple Outputs

For multi-way branching, use a Switch node:

```
HTTP Status Check
├─ 200 → Success Handler
├─ 404 → Not Found Handler
├─ 500 → Error Handler
└─ default → Other Handler
```

**Note:** Switch node configuration is done through the options panel (right-click or click the options button).

### Combining with Loops

You can use conditions inside loops for filtering:

```
For Each User
├─ Age Check (>=18)
│  ├─ True → Add to Adults List
│  └─ False → Add to Minors List
└─ Continue Loop
```

---

## Summary

**You've learned:**
- ✅ How to open the Node Palette and find nodes
- ✅ How to add and configure Number and Condition nodes
- ✅ How to connect from specific handles (green for true, red for false)
- ✅ How to build the age-based API routing workflow
- ✅ How to test both execution paths
- ✅ How to verify only the active path executes

**Key Concepts:**
- 🟢 **Green handle** = True path (condition satisfied)
- 🔴 **Red handle** = False path (condition not satisfied)
- 🎯 **Conditional execution** = Only active branch runs
- 🔄 **Transitive skipping** = Downstream nodes skip if parent skips

**Next Steps:**
- Try building nested conditions (multi-level checks)
- Experiment with Switch nodes for multi-way branching
- Explore combining conditionals with loops and parallel execution
- Check out the example workflows in `examples/conditional-branching/`

---

## Additional Resources

- **Example Workflows:** `examples/conditional-branching/` directory
- **Technical Documentation:** `docs/CONDITIONAL_EXECUTION_IMPLEMENTATION.md`
- **Demo Applications:** `backend/cmd/demo-conditional-execution/main.go`
- **Testing Summary:** `docs/CONDITIONAL_BRANCHING_TESTING_SUMMARY.md`
- **Feature Summary:** `docs/CONDITIONAL_EXECUTION_SUMMARY.md`

---

## Known Issues & Fixes

### Issue: Both Conditional Branches Execute (FIXED)

**Problem:** Prior to commit `[hash]`, when users created conditional workflows following this guide, both conditional branches would execute regardless of the condition result.

**Root Cause:** The frontend was not serializing the `sourceHandle` and `targetHandle` fields when sending the workflow payload to the backend for execution. The backend requires these fields to determine which conditional path should be taken.

**Fix Applied:** Updated `src/app/workflow/page.tsx` line 803-809 to include `sourceHandle` and `targetHandle` in the edge serialization:

```typescript
edges: edges.map((e) => ({
  id: e.id,
  source: e.source,
  target: e.target,
  sourceHandle: e.sourceHandle,  // ← ADDED
  targetHandle: e.targetHandle,  // ← ADDED
})),
```

**Verification:** After this fix:
- ✅ Only the true path executes when condition is met
- ✅ Only the false path executes when condition is not met  
- ✅ Transitive dependencies are properly skipped
- ✅ All 66 conditional execution tests pass

**Testing:** To verify the fix works in your workflow:
1. Create a simple condition workflow (e.g., age >= 18)
2. Add different nodes on true and false paths
3. Execute with age = 25: Only true path nodes should show results
4. Execute with age = 15: Only false path nodes should show results
5. Check execution panel - skipped nodes should not appear in results

---

**Questions or Issues?**

If you encounter problems not covered in this guide, please:
1. Check the Troubleshooting section above
2. Review the example workflows
3. Verify you have the latest version with the sourceHandle fix
4. Open an issue on GitHub with:
   - Steps to reproduce
   - Expected vs actual behavior
   - Screenshot of your workflow
   - Workflow JSON payload (from "View JSON" button)
