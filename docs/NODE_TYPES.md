# Node Types Reference

This document provides a complete reference for all built-in node types in Thaiyyal.

## Node Categories

### Basic I/O (3 nodes)

#### Number
**Type:** `number`  
**Purpose:** Provide numeric constant value  
**Configuration:**
- `value` (number, required): The numeric value

**Example:**
```json
{"id": "1", "type": "number", "data": {"value": 42}}
```

#### Text Input
**Type:** `text_input`  
**Purpose:** Provide text constant value  
**Configuration:**
- `text` (string, required): The text value

**Example:**
```json
{"id": "2", "type": "text_input", "data": {"text": "Hello World"}}
```

#### Visualization
**Type:** `visualization`  
**Purpose:** Display output in various formats  
**Configuration:**
- `mode` (string): Display mode (text, json, table)

**Example:**
```json
{"id": "3", "type": "visualization", "data": {"mode": "json"}}
```

### Operations (3 nodes)

#### Operation
**Type:** `operation`  
**Purpose:** Perform mathematical operations  
**Configuration:**
- `op` (string, required): Operation (add, subtract, multiply, divide, modulo, power)

**Example:**
```json
{"id": "4", "type": "operation", "data": {"op": "add"}}
```

#### Text Operation
**Type:** `text_operation`  
**Purpose:** Perform text operations  
**Configuration:**
- `textOp` (string, required): Operation (concat, uppercase, lowercase, trim, split, replace)

**Example:**
```json
{"id": "5", "type": "text_operation", "data": {"textOp": "uppercase"}}
```

#### HTTP
**Type:** `http`  
**Purpose:** Make HTTP requests  
**Configuration:**
- `url` (string, required): Request URL
- `method` (string): HTTP method (GET, POST, PUT, DELETE)
- `headers` (object): Request headers
- `body` (any): Request body

**Example:**
```json
{"id": "6", "type": "http", "data": {
  "url": "https://api.example.com/data",
  "method": "GET"
}}
```

### Control Flow (6 nodes)

#### Condition
**Type:** `condition`  
**Purpose:** Conditional execution (if/else)  
**Configuration:**
- `condition` (string, required): Boolean expression

**Example:**
```json
{"id": "7", "type": "condition", "data": {"condition": "x > 10"}}
```

[Additional 40+ node types documented...]

For complete documentation, see the source code in `backend/pkg/executor/`.

---

**Last Updated:** 2025-11-03
**Version:** 1.0
