# IF Condition Branching Node - Complete Workflow Examples

## Overview

This document provides complete, real-world workflow examples using IF condition branching nodes with JSON payloads and visual representations.

---

## Example 1: HTTP Status Code Checking with Error Handling

### Description
A workflow that makes an HTTP request and branches based on the response status code. Success path (status 200) processes the data, while failure paths handle different error scenarios.

### Workflow Diagram
```
[Number Input] → [HTTP Request] → [Condition: Status Check]
                                         ├─ TRUE (200) → [Transform Data] → [Visualization]
                                         └─ FALSE (!200) → [Error Handler] → [Visualization]
```

### JSON Payload
```json
{
  "nodes": [
    {
      "id": "input1",
      "type": "number",
      "data": {
        "label": "User ID",
        "value": 123
      },
      "position": { "x": 100, "y": 100 }
    },
    {
      "id": "http1",
      "type": "http",
      "data": {
        "label": "Fetch User Data",
        "url": "https://api.example.com/users/{{input1}}",
        "method": "GET"
      },
      "position": { "x": 300, "y": 100 }
    },
    {
      "id": "condition1",
      "type": "condition",
      "data": {
        "label": "Check Success",
        "condition": "node.http1.output.status == 200"
      },
      "position": { "x": 500, "y": 100 }
    },
    {
      "id": "transform1",
      "type": "transform",
      "data": {
        "label": "Extract User Data",
        "expression": "node.http1.output.data"
      },
      "position": { "x": 700, "y": 50 }
    },
    {
      "id": "error1",
      "type": "transform",
      "data": {
        "label": "Format Error",
        "expression": "'Error: Status ' + node.http1.output.status"
      },
      "position": { "x": 700, "y": 150 }
    },
    {
      "id": "viz1",
      "type": "visualization",
      "data": {
        "label": "Success Result"
      },
      "position": { "x": 900, "y": 50 }
    },
    {
      "id": "viz2",
      "type": "visualization",
      "data": {
        "label": "Error Result"
      },
      "position": { "x": 900, "y": 150 }
    }
  ],
  "edges": [
    { "source": "input1", "target": "http1" },
    { "source": "http1", "target": "condition1" },
    { "source": "condition1", "target": "transform1", "sourceHandle": "true" },
    { "source": "condition1", "target": "error1", "sourceHandle": "false" },
    { "source": "transform1", "target": "viz1" },
    { "source": "error1", "target": "viz2" }
  ]
}
```

### Key Features
- **Condition Expression**: `node.http1.output.status == 200`
- **True Path** (green handle): Processes successful response
- **False Path** (red handle): Handles errors
- **Use Case**: API integration with error handling

---

## Example 2: Multi-Condition Workflow with Variables

### Description
A workflow that processes temperature and humidity sensor data, using multiple condition nodes and workflow variables to determine environmental status.

### Workflow Diagram
```
[Temp Sensor] ─┐
                ├→ [Condition 1: Temp Check]
[Humidity] ────┤      ├─ TRUE → [Set Alarm Variable]
[Variable]      │      └─ FALSE → [Condition 2: Humidity Check]
                │                      ├─ TRUE → [Warning]
                │                      └─ FALSE → [Normal Status]
```

### JSON Payload
```json
{
  "nodes": [
    {
      "id": "temp1",
      "type": "number",
      "data": {
        "label": "Temperature Sensor",
        "value": 28.5
      },
      "position": { "x": 100, "y": 100 }
    },
    {
      "id": "humidity1",
      "type": "number",
      "data": {
        "label": "Humidity Sensor",
        "value": 75
      },
      "position": { "x": 100, "y": 200 }
    },
    {
      "id": "condition1",
      "type": "condition",
      "data": {
        "label": "High Temperature Alert",
        "condition": "node.temp1.value > 30 || (node.temp1.value > 25 && node.humidity1.value > 70)"
      },
      "position": { "x": 350, "y": 150 }
    },
    {
      "id": "alarm",
      "type": "variable",
      "data": {
        "label": "Set Alarm",
        "name": "alarmTriggered",
        "value": true
      },
      "position": { "x": 600, "y": 100 }
    },
    {
      "id": "condition2",
      "type": "condition",
      "data": {
        "label": "Humidity Warning",
        "condition": "node.humidity1.value > 80"
      },
      "position": { "x": 600, "y": 200 }
    },
    {
      "id": "critical",
      "type": "visualization",
      "data": {
        "label": "CRITICAL ALERT",
        "color": "red"
      },
      "position": { "x": 850, "y": 100 }
    },
    {
      "id": "warning",
      "type": "visualization",
      "data": {
        "label": "Warning Status",
        "color": "orange"
      },
      "position": { "x": 850, "y": 200 }
    },
    {
      "id": "normal",
      "type": "visualization",
      "data": {
        "label": "Normal Status",
        "color": "green"
      },
      "position": { "x": 850, "y": 300 }
    }
  ],
  "edges": [
    { "source": "temp1", "target": "condition1" },
    { "source": "humidity1", "target": "condition1" },
    { "source": "condition1", "target": "alarm", "sourceHandle": "true" },
    { "source": "condition1", "target": "condition2", "sourceHandle": "false" },
    { "source": "alarm", "target": "critical" },
    { "source": "condition2", "target": "warning", "sourceHandle": "true" },
    { "source": "condition2", "target": "normal", "sourceHandle": "false" }
  ]
}
```

### Key Features
- **Complex Boolean Logic**: `node.temp1.value > 30 || (node.temp1.value > 25 && node.humidity1.value > 70)`
- **Cascading Conditions**: Multiple condition nodes in sequence
- **Variable Integration**: Sets alarm variable when conditions are met
- **Use Case**: Environmental monitoring systems, IoT sensor processing

---

## Example 3: Data Validation Pipeline with Null Handling

### Description
A workflow that validates incoming data, checking for null values and data quality before processing.

### Workflow Diagram
```
[API Response] → [Condition: Null Check]
                      ├─ FALSE (has data) → [Condition: Data Valid]
                      │                          ├─ TRUE → [Process Data]
                      │                          └─ FALSE → [Invalid Data Handler]
                      └─ TRUE (null) → [Missing Data Handler]
```

### JSON Payload
```json
{
  "nodes": [
    {
      "id": "api1",
      "type": "http",
      "data": {
        "label": "Fetch Customer Data",
        "url": "https://api.example.com/customer/123",
        "method": "GET"
      },
      "position": { "x": 100, "y": 150 }
    },
    {
      "id": "condition1",
      "type": "condition",
      "data": {
        "label": "Check Data Exists",
        "condition": "!isNull(node.api1.output.data) && !isNull(node.api1.output.data.customerId)"
      },
      "position": { "x": 350, "y": 150 }
    },
    {
      "id": "condition2",
      "type": "condition",
      "data": {
        "label": "Validate Customer ID",
        "condition": "node.api1.output.data.customerId > 0 && node.api1.output.data.status == 'active'"
      },
      "position": { "x": 600, "y": 100 }
    },
    {
      "id": "process1",
      "type": "transform",
      "data": {
        "label": "Process Valid Customer",
        "expression": "node.api1.output.data"
      },
      "position": { "x": 850, "y": 50 }
    },
    {
      "id": "invalid1",
      "type": "transform",
      "data": {
        "label": "Handle Invalid Data",
        "expression": "'Invalid customer: ' + node.api1.output.data.status"
      },
      "position": { "x": 850, "y": 150 }
    },
    {
      "id": "missing1",
      "type": "transform",
      "data": {
        "label": "Handle Missing Data",
        "expression": "'Customer data not found'"
      },
      "position": { "x": 600, "y": 250 }
    },
    {
      "id": "viz1",
      "type": "visualization",
      "data": {
        "label": "Valid Customer"
      },
      "position": { "x": 1050, "y": 50 }
    },
    {
      "id": "viz2",
      "type": "visualization",
      "data": {
        "label": "Invalid Customer"
      },
      "position": { "x": 1050, "y": 150 }
    },
    {
      "id": "viz3",
      "type": "visualization",
      "data": {
        "label": "Missing Data"
      },
      "position": { "x": 850, "y": 250 }
    }
  ],
  "edges": [
    { "source": "api1", "target": "condition1" },
    { "source": "condition1", "target": "condition2", "sourceHandle": "true" },
    { "source": "condition1", "target": "missing1", "sourceHandle": "false" },
    { "source": "condition2", "target": "process1", "sourceHandle": "true" },
    { "source": "condition2", "target": "invalid1", "sourceHandle": "false" },
    { "source": "process1", "target": "viz1" },
    { "source": "invalid1", "target": "viz2" },
    { "source": "missing1", "target": "viz3" }
  ]
}
```

### Key Features
- **Null Checking**: `!isNull(node.api1.output.data) && !isNull(node.api1.output.data.customerId)`
- **Data Validation**: Multiple validation layers
- **Error Handling**: Separate paths for different error types
- **Use Case**: API data validation, data quality pipelines

---

## Example 4: Time-Based Workflow Execution

### Description
A workflow that executes different actions based on the current time and date, useful for scheduled tasks and time-sensitive operations.

### Workflow Diagram
```
[Timer Input] → [Condition: Business Hours]
                     ├─ TRUE (9AM-5PM) → [Condition: Weekday]
                     │                        ├─ TRUE → [Process Immediately]
                     │                        └─ FALSE → [Queue for Monday]
                     └─ FALSE (After Hours) → [Queue for Next Day]
```

### JSON Payload
```json
{
  "nodes": [
    {
      "id": "timer1",
      "type": "contextVariable",
      "data": {
        "label": "Current Time",
        "variable": "now()"
      },
      "position": { "x": 100, "y": 150 }
    },
    {
      "id": "condition1",
      "type": "condition",
      "data": {
        "label": "Business Hours Check",
        "condition": "hour(now()) >= 9 && hour(now()) < 17"
      },
      "position": { "x": 300, "y": 150 }
    },
    {
      "id": "condition2",
      "type": "condition",
      "data": {
        "label": "Weekday Check",
        "condition": "day(now()) >= 1 && day(now()) <= 5"
      },
      "position": { "x": 550, "y": 100 }
    },
    {
      "id": "process1",
      "type": "transform",
      "data": {
        "label": "Process Now",
        "expression": "'Processing during business hours: ' + now()"
      },
      "position": { "x": 800, "y": 50 }
    },
    {
      "id": "queue1",
      "type": "transform",
      "data": {
        "label": "Queue for Monday",
        "expression": "'Queued for next Monday'"
      },
      "position": { "x": 800, "y": 150 }
    },
    {
      "id": "queue2",
      "type": "transform",
      "data": {
        "label": "Queue for Tomorrow",
        "expression": "'Queued for next business day'"
      },
      "position": { "x": 550, "y": 250 }
    },
    {
      "id": "viz1",
      "type": "visualization",
      "data": {
        "label": "Immediate Processing"
      },
      "position": { "x": 1000, "y": 50 }
    },
    {
      "id": "viz2",
      "type": "visualization",
      "data": {
        "label": "Queued Tasks"
      },
      "position": { "x": 1000, "y": 200 }
    }
  ],
  "edges": [
    { "source": "timer1", "target": "condition1" },
    { "source": "condition1", "target": "condition2", "sourceHandle": "true" },
    { "source": "condition1", "target": "queue2", "sourceHandle": "false" },
    { "source": "condition2", "target": "process1", "sourceHandle": "true" },
    { "source": "condition2", "target": "queue1", "sourceHandle": "false" },
    { "source": "process1", "target": "viz1" },
    { "source": "queue1", "target": "viz2" },
    { "source": "queue2", "target": "viz2" }
  ]
}
```

### Key Features
- **Date/Time Functions**: `hour(now())`, `day(now())`
- **Time-Based Routing**: Different paths for business hours vs after hours
- **Scheduler Logic**: Weekday vs weekend handling
- **Use Case**: Scheduled task execution, business hours automation

---

## Example 5: Counter and Loop Control with Variables

### Description
A workflow that uses a counter variable with condition nodes to control loop execution and process data in batches.

### Workflow Diagram
```
[Start] → [Counter Variable] → [ForEach Loop]
                                    └→ [Condition: Counter Check]
                                           ├─ TRUE (< 100) → [Increment Counter] → [Process Item]
                                           └─ FALSE (>= 100) → [Stop Processing]
```

### JSON Payload
```json
{
  "nodes": [
    {
      "id": "start1",
      "type": "number",
      "data": {
        "label": "Start Signal",
        "value": 1
      },
      "position": { "x": 100, "y": 150 }
    },
    {
      "id": "counter1",
      "type": "counter",
      "data": {
        "label": "Item Counter",
        "initialValue": 0
      },
      "position": { "x": 250, "y": 150 }
    },
    {
      "id": "foreach1",
      "type": "foreach",
      "data": {
        "label": "Process Items",
        "items": "[1,2,3,4,5,6,7,8,9,10]"
      },
      "position": { "x": 400, "y": 150 }
    },
    {
      "id": "condition1",
      "type": "condition",
      "data": {
        "label": "Check Limit",
        "condition": "variables.counter < 100 && !isNull(node.foreach1.current)"
      },
      "position": { "x": 600, "y": 150 }
    },
    {
      "id": "increment1",
      "type": "accumulator",
      "data": {
        "label": "Increment Counter",
        "operation": "add",
        "value": 1
      },
      "position": { "x": 850, "y": 100 }
    },
    {
      "id": "process1",
      "type": "transform",
      "data": {
        "label": "Process Item",
        "expression": "'Processing item #' + variables.counter"
      },
      "position": { "x": 1050, "y": 100 }
    },
    {
      "id": "stop1",
      "type": "visualization",
      "data": {
        "label": "Limit Reached",
        "message": "Processing stopped at limit"
      },
      "position": { "x": 850, "y": 200 }
    },
    {
      "id": "viz1",
      "type": "visualization",
      "data": {
        "label": "Processed Items"
      },
      "position": { "x": 1250, "y": 100 }
    }
  ],
  "edges": [
    { "source": "start1", "target": "counter1" },
    { "source": "counter1", "target": "foreach1" },
    { "source": "foreach1", "target": "condition1" },
    { "source": "condition1", "target": "increment1", "sourceHandle": "true" },
    { "source": "condition1", "target": "stop1", "sourceHandle": "false" },
    { "source": "increment1", "target": "process1" },
    { "source": "process1", "target": "viz1" }
  ]
}
```

### Key Features
- **Variable Reference**: `variables.counter < 100`
- **Null Checking with Logic**: `!isNull(node.foreach1.current)`
- **Loop Control**: Condition determines whether to continue processing
- **Use Case**: Batch processing, rate limiting, controlled iteration

---

## Example 6: String Pattern Matching and Text Processing

### Description
A workflow that processes log entries, checking for error patterns and routing to different handlers based on log severity.

### Workflow Diagram
```
[Log Input] → [Condition: Contains Error]
                  ├─ TRUE → [Condition: Critical Error]
                  │             ├─ TRUE → [Alert Handler]
                  │             └─ FALSE → [Error Logger]
                  └─ FALSE → [Info Logger]
```

### JSON Payload
```json
{
  "nodes": [
    {
      "id": "log1",
      "type": "textInput",
      "data": {
        "label": "Log Entry",
        "value": "ERROR: Database connection failed at 2024-01-15T10:30:00Z"
      },
      "position": { "x": 100, "y": 150 }
    },
    {
      "id": "condition1",
      "type": "condition",
      "data": {
        "label": "Check for Errors",
        "condition": "contains(node.log1.value, 'ERROR') || contains(node.log1.value, 'CRITICAL')"
      },
      "position": { "x": 350, "y": 150 }
    },
    {
      "id": "condition2",
      "type": "condition",
      "data": {
        "label": "Critical Error Check",
        "condition": "contains(node.log1.value, 'CRITICAL') || contains(node.log1.value, 'FATAL') || contains(node.log1.value, 'Database connection failed')"
      },
      "position": { "x": 600, "y": 100 }
    },
    {
      "id": "alert1",
      "type": "transform",
      "data": {
        "label": "Send Alert",
        "expression": "'ALERT: ' + node.log1.value"
      },
      "position": { "x": 850, "y": 50 }
    },
    {
      "id": "error1",
      "type": "transform",
      "data": {
        "label": "Log Error",
        "expression": "'Error logged: ' + node.log1.value"
      },
      "position": { "x": 850, "y": 150 }
    },
    {
      "id": "info1",
      "type": "transform",
      "data": {
        "label": "Log Info",
        "expression": "'Info: ' + node.log1.value"
      },
      "position": { "x": 600, "y": 250 }
    },
    {
      "id": "viz1",
      "type": "visualization",
      "data": {
        "label": "Critical Alerts",
        "color": "red"
      },
      "position": { "x": 1050, "y": 50 }
    },
    {
      "id": "viz2",
      "type": "visualization",
      "data": {
        "label": "Error Log",
        "color": "orange"
      },
      "position": { "x": 1050, "y": 150 }
    },
    {
      "id": "viz3",
      "type": "visualization",
      "data": {
        "label": "Info Log",
        "color": "blue"
      },
      "position": { "x": 850, "y": 250 }
    }
  ],
  "edges": [
    { "source": "log1", "target": "condition1" },
    { "source": "condition1", "target": "condition2", "sourceHandle": "true" },
    { "source": "condition1", "target": "info1", "sourceHandle": "false" },
    { "source": "condition2", "target": "alert1", "sourceHandle": "true" },
    { "source": "condition2", "target": "error1", "sourceHandle": "false" },
    { "source": "alert1", "target": "viz1" },
    { "source": "error1", "target": "viz2" },
    { "source": "info1", "target": "viz3" }
  ]
}
```

### Key Features
- **String Matching**: `contains(node.log1.value, 'ERROR')`
- **Multiple Patterns**: Checking for various error keywords
- **Cascading Severity**: Different handlers for different severity levels
- **Use Case**: Log processing, error monitoring, alerting systems

---

## Example 7: Arithmetic Calculations with Conditional Routing

### Description
A workflow that performs mathematical calculations and routes based on the results, useful for financial calculations, scoring systems, or data analysis.

### Workflow Diagram
```
[Price Input] → [Calculate Total] → [Condition: Discount Eligible]
[Quantity]                               ├─ TRUE (>1000) → [Apply 20% Discount]
[Tax Rate]                               └─ FALSE → [Apply 10% Discount]
                                                         ↓
                                              [Condition: Free Shipping]
                                                  ├─ TRUE → [Final Total]
                                                  └─ FALSE → [Add Shipping]
```

### JSON Payload
```json
{
  "nodes": [
    {
      "id": "price1",
      "type": "number",
      "data": {
        "label": "Unit Price",
        "value": 25.50
      },
      "position": { "x": 100, "y": 100 }
    },
    {
      "id": "qty1",
      "type": "number",
      "data": {
        "label": "Quantity",
        "value": 50
      },
      "position": { "x": 100, "y": 180 }
    },
    {
      "id": "tax1",
      "type": "number",
      "data": {
        "label": "Tax Rate %",
        "value": 8.5
      },
      "position": { "x": 100, "y": 260 }
    },
    {
      "id": "calc1",
      "type": "transform",
      "data": {
        "label": "Calculate Subtotal",
        "expression": "node.price1.value * node.qty1.value"
      },
      "position": { "x": 300, "y": 150 }
    },
    {
      "id": "condition1",
      "type": "condition",
      "data": {
        "label": "Large Order Discount",
        "condition": "(node.price1.value * node.qty1.value) > 1000"
      },
      "position": { "x": 500, "y": 150 }
    },
    {
      "id": "discount1",
      "type": "transform",
      "data": {
        "label": "20% Discount",
        "expression": "node.calc1.value * 0.80"
      },
      "position": { "x": 700, "y": 100 }
    },
    {
      "id": "discount2",
      "type": "transform",
      "data": {
        "label": "10% Discount",
        "expression": "node.calc1.value * 0.90"
      },
      "position": { "x": 700, "y": 200 }
    },
    {
      "id": "condition2",
      "type": "condition",
      "data": {
        "label": "Free Shipping Eligible",
        "condition": "node.calc1.value > 500"
      },
      "position": { "x": 900, "y": 150 }
    },
    {
      "id": "final1",
      "type": "transform",
      "data": {
        "label": "Final Total (Free Ship)",
        "expression": "node.discount1.value * (1 + node.tax1.value / 100)"
      },
      "position": { "x": 1100, "y": 100 }
    },
    {
      "id": "final2",
      "type": "transform",
      "data": {
        "label": "Final Total + Shipping",
        "expression": "(node.discount1.value * (1 + node.tax1.value / 100)) + 15"
      },
      "position": { "x": 1100, "y": 200 }
    },
    {
      "id": "viz1",
      "type": "visualization",
      "data": {
        "label": "Order Total"
      },
      "position": { "x": 1300, "y": 150 }
    }
  ],
  "edges": [
    { "source": "price1", "target": "calc1" },
    { "source": "qty1", "target": "calc1" },
    { "source": "tax1", "target": "calc1" },
    { "source": "calc1", "target": "condition1" },
    { "source": "condition1", "target": "discount1", "sourceHandle": "true" },
    { "source": "condition1", "target": "discount2", "sourceHandle": "false" },
    { "source": "discount1", "target": "condition2" },
    { "source": "discount2", "target": "condition2" },
    { "source": "condition2", "target": "final1", "sourceHandle": "true" },
    { "source": "condition2", "target": "final2", "sourceHandle": "false" },
    { "source": "final1", "target": "viz1" },
    { "source": "final2", "target": "viz1" }
  ]
}
```

### Key Features
- **Arithmetic in Conditions**: `(node.price1.value * node.qty1.value) > 1000`
- **Nested Calculations**: Multiple calculation nodes with conditional branching
- **Business Logic**: Discount tiers and shipping rules
- **Use Case**: E-commerce pricing, invoice calculation, financial workflows

---

## Visual Elements Guide

### Condition Node Appearance

```
┌─────────────────────────┐
│  Condition Node         │
│  "Check Success"        │
│                         │
│  Condition:             │
│  node.http1.status==200 │
│                         │
│  ●─ TRUE (green)        │  ← Top right handle
│  ●─ FALSE (red)         │  ← Bottom right handle
└─────────────────────────┘
```

### Edge Connection Pattern

- **Green Edge** (TRUE path): Connects from top-right handle
- **Red Edge** (FALSE path): Connects from bottom-right handle
- **Label**: Edge can optionally show "true" or "false" label

---

## Best Practices from Examples

### 1. Expression Complexity
- Start simple, add complexity as needed
- Use parentheses for clarity: `(a && b) || (c && d)`
- Break complex logic into multiple condition nodes

### 2. Error Handling
- Always have FALSE path connected
- Handle null values explicitly with `isNull()`
- Validate data before processing

### 3. Performance
- Place most common conditions first
- Avoid deeply nested conditions when possible
- Use variables to cache computed values

### 4. Maintainability
- Use descriptive node labels
- Document complex conditions
- Group related conditions visually

### 5. Testing
- Test with null values
- Test boundary conditions
- Test with invalid data

---

## Common Condition Patterns

### Pattern 1: Null-Safe Check
```javascript
!isNull(node.data.value) && node.data.value > 100
```

### Pattern 2: Range Check
```javascript
node.value >= 10 && node.value <= 100
```

### Pattern 3: Multiple Alternatives
```javascript
node.status == 'active' || node.status == 'pending' || node.status == 'review'
```

### Pattern 4: Combined Conditions
```javascript
(node.age.value >= 18 && node.verified.value) || node.admin.value
```

### Pattern 5: String Pattern
```javascript
contains(node.email.value, '@') && contains(node.email.value, '.')
```

---

## Troubleshooting

### Issue: Condition always evaluates to FALSE
**Solution**: Check for null values, use `isNull()` to verify data exists

### Issue: Edge doesn't connect to correct path
**Solution**: Ensure you're connecting to the correct handle (green=true, red=false)

### Issue: Complex expression errors
**Solution**: Break into multiple simpler conditions or verify syntax

### Issue: Variable not found
**Solution**: Ensure variable is set before the condition node in workflow execution order

---

## Additional Resources

- **Expression Syntax Guide**: See `CONDITION_NODE_GUIDE.md`
- **Date/Time Functions**: See `DATETIME_NULL_HANDLING.md`
- **Math Functions**: See main documentation
- **API Reference**: Check `EXPRESSION_SYSTEM_DESIGN.md`

---

**Last Updated**: 2025-11-01
**Version**: 1.0.0
