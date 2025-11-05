# Workflow Execution Guide

This guide explains how to create, run, and view results of workflows in Thaiyyal.

## Overview

Thaiyyal provides a visual workflow builder with real-time execution capabilities. You can create workflows by dragging and connecting nodes, then execute them to see immediate results.

## Step-by-Step Workflow Execution

### Step 1: Access the Workflow Builder

1. Open your browser and navigate to the Thaiyyal application
2. Click on "Create New Workflow" or navigate to `/workflow`
3. You'll see the visual workflow canvas

![Workflow Builder](screenshots/01-workflow-builder.png)

### Step 2: Add Nodes to Your Workflow

1. Click the "+ Add Node" button in the bottom-left corner
2. Browse through the categorized node palette
3. Click on a node type to add it to the canvas
4. Drag nodes to position them on the canvas

**Available Node Categories:**
- **I/O Nodes**: TextInput, TextOperation, HTTP
- **Data Operations**: Transform, Parse, Extract, Map, Filter
- **Control Flow**: Condition, Switch, ForEach, WhileLoop
- **Parallel & Join**: Parallel, Join, Split
- **State & Memory**: Variable, Accumulator, Counter, Cache
- **Resilience**: Retry, TryCatch, Timeout, Delay
- **Array Operations**: Slice, Sort, Find, FlatMap, GroupBy, Unique, Chunk, Reverse, Partition, Zip, Sample, Range, Transpose
- **Visualization**: BarChart

![Node Palette](screenshots/02-node-palette.png)

### Step 3: Connect Nodes

1. Click and drag from a node's output handle (right side)
2. Drop on another node's input handle (left side)
3. Edges will automatically connect and route between nodes
4. Multiple inputs and outputs are supported

![Connecting Nodes](screenshots/03-connecting-nodes.png)

### Step 4: Configure Node Properties

1. Click on a node to select it
2. Edit properties in the node's configuration panel
3. Each node type has specific configuration options
4. Changes are automatically saved

**Example Node Configurations:**
- **Number Node**: Set a numeric value
- **Add Operation**: Automatically adds incoming values
- **HTTP Node**: Configure URL, method, headers
- **Condition Node**: Set condition expression

![Node Configuration](screenshots/04-node-configuration.png)

### Step 5: Run the Workflow

1. Click the **"▶ Run"** button in the top navigation bar
2. The execution panel automatically opens at the bottom
3. A loading indicator shows the workflow is executing
4. You can cancel execution at any time by clicking "Cancel"

![Run Button](screenshots/05-run-button.png)

### Step 6: View Execution Results

The execution panel displays comprehensive results:

#### Execution Panel Features

**1. Resizable Panel**
- Drag the top edge to resize the panel height
- Panel can be adjusted from 100px to 600px
- Smooth drag interaction with visual feedback

![Resizable Panel](screenshots/06-resizable-panel.png)

**2. Loading State**
- Animated spinner during execution
- "Running workflow..." status message
- Cancel button to abort execution
- Real-time progress indication

![Loading State](screenshots/07-loading-state.png)

**3. Execution Summary**
- **Execution ID**: Unique identifier for this run
- **Duration**: Total execution time
- **Status**: Success or failure indicator
- **Workflow Info**: ID and name if saved

![Execution Summary](screenshots/08-execution-summary.png)

**4. Final Output**
- The final result of the workflow
- Formatted JSON display
- Syntax highlighting
- Scrollable for large outputs

![Final Output](screenshots/09-final-output.png)

**5. Node Results**
- Individual results for each node
- Node ID and output value
- Collapsible sections
- Formatted JSON for each node

![Node Results](screenshots/10-node-results.png)

**6. Error Handling**
- Detailed error messages
- Error location indication
- Validation feedback
- Network error handling

![Error Display](screenshots/11-error-display.png)

### Step 7: Manage Execution Panel

**Close the Panel:**
- Click the "✕" button in the top-right corner
- Panel closes and status bar reappears

**Resize the Panel:**
- Hover over the top edge until cursor changes
- Click and drag up or down to resize
- Release to set new height

**Cancel Execution:**
- Click "Cancel" button during execution
- Workflow execution stops immediately
- Panel shows cancellation message

![Panel Management](screenshots/12-panel-management.png)

## Example Workflows

### Example 1: Simple Addition

**Nodes:**
1. Number node (value: 10)
2. Number node (value: 5)
3. Add operation node

**Connections:**
- Number 1 → Add
- Number 2 → Add

**Expected Result:**
```json
{
  "execution_id": "abc123",
  "final_output": 15,
  "node_results": {
    "1": 10,
    "2": 5,
    "3": 15
  }
}
```

![Addition Example](screenshots/13-addition-example.png)

### Example 2: HTTP Request with Filtering

**Nodes:**
1. HTTP node (GET request to API)
2. Filter node (filter by criteria)
3. Transform node (format output)

**Expected Result:**
- Fetched data from API
- Filtered records
- Transformed output

![HTTP Filter Example](screenshots/14-http-filter-example.png)

### Example 3: Loop with Accumulation

**Nodes:**
1. Range node (generates array)
2. ForEach node (iterates)
3. Accumulator node (sums values)

**Expected Result:**
- Array generated
- Each value processed
- Final sum calculated

![Loop Accumulation Example](screenshots/15-loop-example.png)

## Advanced Features

### Canceling Long-Running Workflows

For workflows that take a long time:

1. Monitor the loading indicator
2. Click "Cancel" if needed
3. Execution stops immediately
4. Panel shows cancellation status

**Use Cases for Cancellation:**
- Infinite loops
- Long HTTP requests
- Large dataset processing
- Debugging workflow logic

### Viewing Detailed Node Results

For complex workflows:

1. Scroll through the node results section
2. Each node shows its individual output
3. Expand JSON for detailed inspection
4. Identify bottlenecks or errors

### Saving and Loading Workflows

**Save Workflow:**
1. Click "💾 Save" in the navigation bar
2. Workflow persists to the server
3. Get a unique workflow ID

**Load Workflow:**
1. Click "📂 Open" in the navigation bar
2. Select from saved workflows
3. Workflow loads on canvas

**Execute Saved Workflow:**
- Use the workflow ID
- Execute via API: `POST /api/v1/workflow/execute/{id}`

## API Integration

### Execute Workflow Endpoint

```bash
POST /api/v1/workflow/execute
Content-Type: application/json

{
  "nodes": [
    {"id": "1", "data": {"value": 10}},
    {"id": "2", "data": {"value": 5}},
    {"id": "3", "data": {"op": "add"}}
  ],
  "edges": [
    {"source": "1", "target": "3"},
    {"source": "2", "target": "3"}
  ]
}
```

**Response:**
```json
{
  "success": true,
  "execution_time": "245.3µs",
  "results": {
    "execution_id": "abc123def456",
    "node_results": {
      "1": 10,
      "2": 5,
      "3": 15
    },
    "final_output": 15
  }
}
```

### Execute by ID Endpoint

```bash
POST /api/v1/workflow/execute/{workflow-id}
```

**Response:**
```json
{
  "success": true,
  "workflow_id": "uuid-here",
  "workflow_name": "My Saved Workflow",
  "execution_time": "1.2ms",
  "results": {
    "execution_id": "abc123",
    "node_results": {...},
    "final_output": {...}
  }
}
```

## Keyboard Shortcuts

- **Delete Node**: Select node + Delete/Backspace
- **Deselect**: Escape
- **Zoom In**: Ctrl/Cmd + Plus
- **Zoom Out**: Ctrl/Cmd + Minus
- **Fit View**: Ctrl/Cmd + 0

## Troubleshooting

### Execution Fails

**Problem**: Workflow execution returns an error

**Solutions:**
1. Check node configurations
2. Verify all required connections
3. Review error message in execution panel
4. Check network connectivity
5. Validate input data formats

### Panel Not Opening

**Problem**: Run button clicked but panel doesn't appear

**Solutions:**
1. Check browser console for errors
2. Refresh the page
3. Clear browser cache
4. Check if workflow has nodes

### Slow Execution

**Problem**: Workflow takes too long

**Solutions:**
1. Optimize node count
2. Reduce loop iterations
3. Check HTTP timeout settings
4. Use Cancel button if needed
5. Review workflow logic

## Best Practices

### Workflow Design

1. **Keep it Simple**: Start with small workflows
2. **Test Incrementally**: Add and test one node at a time
3. **Use Descriptive Names**: Name workflows clearly
4. **Document Complex Logic**: Add comments/descriptions
5. **Handle Errors**: Use TryCatch nodes

### Execution Management

1. **Monitor Progress**: Watch the execution panel
2. **Check Results**: Review node-by-node outputs
3. **Save Successful Workflows**: Persist working configurations
4. **Cancel When Needed**: Don't wait for timeout
5. **Review Errors**: Read error messages carefully

### Performance Tips

1. **Limit Iterations**: Set max_iterations on loops
2. **Use Parallel Nodes**: Process arrays concurrently
3. **Cache Results**: Use Cache nodes for repeated data
4. **Optimize Connections**: Minimize unnecessary edges
5. **Test with Small Data**: Validate logic before scaling

## Summary

The Thaiyyal workflow execution system provides:

✅ **Visual Workflow Builder** - Drag-and-drop node creation
✅ **Real-time Execution** - Immediate results display
✅ **Resizable Results Panel** - Flexible viewing experience
✅ **Loading Indicators** - Clear execution state
✅ **Cancellation Support** - Stop long-running workflows
✅ **Detailed Results** - Node-by-node breakdown
✅ **Error Handling** - Comprehensive error messages
✅ **Professional UI** - Enterprise-grade design

For more information, see the [API Documentation](../API_EXAMPLES.md) and [Server Implementation](../SERVER_IMPLEMENTATION.md).
