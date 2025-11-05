# Data Format Support Guide

## Overview

Thaiyyal provides comprehensive support for parsing and formatting data in multiple formats, enabling seamless data conversion workflows. This guide covers all supported formats and their usage patterns.

## Table of Contents

1. [Supported Formats](#supported-formats)
2. [Parse Node](#parse-node)
3. [Format Node](#format-node)
4. [Format Conversion Workflows](#format-conversion-workflows)
5. [Best Practices](#best-practices)
6. [Examples](#examples)

## Supported Formats

### Input Formats (Parse Node)

| Format | Description | Auto-Detection | Type Inference |
|--------|-------------|----------------|----------------|
| JSON | JavaScript Object Notation | ✅ Yes | ✅ Yes |
| CSV | Comma-Separated Values | ✅ Yes | ✅ Yes |
| TSV | Tab-Separated Values | ✅ Yes | ✅ Yes |
| YAML | YAML Ain't Markup Language | ✅ Yes | ✅ Yes |
| XML | Extensible Markup Language | ✅ Yes | Limited |

### Output Formats (Format Node)

| Format | Description | Options |
|--------|-------------|---------|
| JSON | JavaScript Object Notation | Pretty-print |
| CSV | Comma-Separated Values | Headers, Delimiter |
| TSV | Tab-Separated Values | Headers |

## Parse Node

Converts string data to structured formats.

### Configuration

```javascript
{
  type: "parse",
  data: {
    input_type: "AUTO" | "JSON" | "CSV" | "TSV" | "YAML" | "XML"
  }
}
```

### Auto-Detection

When `input_type` is `"AUTO"` (default), the parser automatically detects the format:

- **JSON**: Begins with `{` or `[`
- **XML**: Begins with `<`
- **YAML**: Contains `key: value` patterns
- **CSV**: Contains commas with multiple lines
- **TSV**: Contains tabs with multiple lines

### Type Inference

The parser automatically infers data types:

```csv
name,age,active
Alice,30,true
Bob,25,false
```

Becomes:

```json
[
  {"name": "Alice", "age": 30, "active": true},
  {"name": "Bob", "age": 25, "active": false}
]
```

**Supported Types**:
- Numbers: `"123"` → `123`
- Booleans: `"true"` → `true`, `"false"` → `false`
- Null: `"null"` → `null`
- Strings: Everything else

### Format-Specific Parsing

#### JSON Parsing

```javascript
// Input
'{"name": "Alice", "age": 30}'

// Output
{name: "Alice", age: 30}
```

**Features**:
- Preserves number precision
- Handles nested objects and arrays
- Supports all JSON primitives

#### CSV Parsing

```javascript
// Input
"name,age\nAlice,30\nBob,25"

// Output
[
  {name: "Alice", age: 30},
  {name: "Bob", age: 25}
]
```

**Features**:
- First row used as headers
- Automatic type inference
- Trims leading/trailing whitespace
- Handles quoted values

#### TSV Parsing

Same as CSV but uses tab `\t` as delimiter.

#### YAML Parsing

```javascript
// Input
"name: Alice\nage: 30\nactive: true"

// Output
{name: "Alice", age: 30, active: true}
```

**Features**:
- Simple key-value pairs
- Type inference (numbers, booleans)
- Comment support (lines starting with #)

#### XML Parsing

```javascript
// Input
"<user><name>Alice</name><age>30</age></user>"

// Output
{user: {name: "Alice", age: 30}}
```

**Note**: Basic XML support for simple structures.

### Error Handling

If parsing fails:
- Returns descriptive error message
- Indicates the format and issue
- Example: `"invalid JSON: unexpected token at position 10"`

## Format Node

Converts structured data to string formats.

### Configuration

```javascript
{
  type: "format",
  data: {
    output_type: "JSON" | "CSV" | "TSV",  // Default: "JSON"
    pretty_print: boolean,                 // JSON only, default: false
    include_headers: boolean,              // CSV/TSV only, default: true
    delimiter: string                      // CSV only, default: ","
  }
}
```

### JSON Formatting

#### Compact Mode

```javascript
{
  output_type: "JSON",
  pretty_print: false  // default
}
```

**Input**:
```javascript
{name: "Alice", age: 30}
```

**Output**:
```json
{"name":"Alice","age":30}
```

#### Pretty-Print Mode

```javascript
{
  output_type: "JSON",
  pretty_print: true
}
```

**Output**:
```json
{
  "name": "Alice",
  "age": 30
}
```

### CSV Formatting

#### With Headers (Default)

```javascript
{
  output_type: "CSV",
  include_headers: true  // default
}
```

**Input**:
```javascript
[
  {name: "Alice", age: 30},
  {name: "Bob", age: 25}
]
```

**Output**:
```csv
age,name
30,Alice
25,Bob
```

**Note**: Headers are sorted alphabetically for consistency.

#### Without Headers

```javascript
{
  output_type: "CSV",
  include_headers: false
}
```

**Output**:
```csv
30,Alice
25,Bob
```

#### Custom Delimiter

```javascript
{
  output_type: "CSV",
  delimiter: "|"
}
```

**Output**:
```
age|name
30|Alice
25|Bob
```

**Common Delimiters**:
- Pipe: `|`
- Semicolon: `;`
- Tab: `\t` (use TSV output type instead)

### TSV Formatting

Same as CSV but uses tab `\t` as delimiter.

```javascript
{
  output_type: "TSV",
  include_headers: true
}
```

**Output**:
```tsv
age	name
30	Alice
25	Bob
```

### Type Handling

The formatter intelligently handles different data types:

| Type | CSV/TSV Output | JSON Output |
|------|----------------|-------------|
| String | `"hello"` | `"hello"` |
| Number (int) | `42` | `42` |
| Number (float) | `3.14` | `3.14` |
| Boolean | `true`/`false` | `true`/`false` |
| Null | `` (empty) | `null` |
| Array | JSON string | Array |
| Object | JSON string | Object |

### Error Handling

Common errors:
- **CSV/TSV with non-object array**: Use array of objects
- **Invalid output type**: Use JSON, CSV, or TSV
- **No input**: Provide data to format

## Format Conversion Workflows

### CSV to JSON

```
CSV String → Parse(CSV) → Format(JSON) → JSON String
```

**Use Case**: Convert CSV exports to JSON for API consumption.

```javascript
// Workflow
[
  {type: "parse", data: {input_type: "CSV"}},
  {type: "format", data: {output_type: "JSON", pretty_print: true}}
]
```

### JSON to CSV

```
JSON Array → Format(CSV) → CSV String
```

**Use Case**: Export JSON data to CSV for spreadsheets.

```javascript
// Workflow
[
  {type: "format", data: {output_type: "CSV", include_headers: true}}
]
```

### CSV Transform CSV

```
CSV → Parse → Filter/Map/Transform → Format → CSV
```

**Use Case**: Process and filter CSV data.

```javascript
// Workflow
[
  {type: "parse", data: {input_type: "CSV"}},
  {type: "filter", data: {condition: "item.age > 25"}},
  {type: "format", data: {output_type: "CSV"}}
]
```

### Multi-Format Pipeline

```
CSV → JSON → Process → CSV
```

**Example**:
1. Parse CSV to structured data
2. Format to JSON for intermediate storage
3. Process/transform data
4. Format back to CSV for export

### Round-Trip Conversion

Parse and Format nodes are complementary:

```javascript
// Original data
data = [{name: "Alice", age: 30}]

// Format to CSV
csv = Format(data, "CSV")
// "age,name\n30,Alice\n"

// Parse back
parsed = Parse(csv, "CSV")
// [{name: "Alice", age: 30}]
```

**Note**: Some type information may be lost in CSV (all values become strings or inferred types).

## Best Practices

### 1. Use Auto-Detection for Unknown Formats

```javascript
{
  type: "parse",
  data: {input_type: "AUTO"}
}
```

Let the parser detect the format automatically.

### 2. Validate Data Before Formatting

```javascript
// Check if data is an array of objects before CSV formatting
if (Array.isArray(data) && data.every(item => typeof item === 'object')) {
  // Safe to format as CSV
}
```

### 3. Handle Large Datasets Carefully

For large CSV files:
- Consider chunking data
- Use streaming if available
- Monitor memory usage

### 4. Preserve Headers in CSV

```javascript
{
  output_type: "CSV",
  include_headers: true  // Makes CSV self-documenting
}
```

### 5. Use Pretty-Print for Readability

```javascript
{
  output_type: "JSON",
  pretty_print: true  // For human-readable output
}
```

### 6. Choose Appropriate Delimiters

- **Comma** (`,`): Standard CSV
- **Tab** (`\t`): TSV, better for data with commas
- **Pipe** (`|`): When data contains commas and tabs
- **Semicolon** (`;`): European CSV standard

### 7. Test Round-Trip Conversions

Always test that `Parse(Format(data))` preserves your data structure.

## Examples

### Example 1: CSV to JSON API Response

**Scenario**: Convert CSV report to JSON for web API.

```javascript
// Input: CSV report
const csvReport = `
product,sales,revenue
Widget,100,1999.99
Gadget,75,2249.25
`;

// Workflow
[
  {
    id: "parse_csv",
    type: "parse",
    data: {input_type: "CSV"}
  },
  {
    id: "format_json",
    type: "format",
    data: {
      output_type: "JSON",
      pretty_print: true
    }
  }
]

// Output
{
  "sales": [
    {
      "product": "Widget",
      "sales": 100,
      "revenue": 1999.99
    },
    {
      "product": "Gadget",
      "sales": 75,
      "revenue": 2249.25
    }
  ]
}
```

### Example 2: Filter and Export to CSV

**Scenario**: Filter high-value sales and export to CSV.

```javascript
// Workflow
[
  {
    id: "parse_data",
    type: "parse",
    data: {input_type: "JSON"}
  },
  {
    id: "filter_high_value",
    type: "filter",
    data: {condition: "item.revenue > 2000"}
  },
  {
    id: "export_csv",
    type: "format",
    data: {
      output_type: "CSV",
      include_headers: true
    }
  }
]

// Output: CSV with only high-value sales
product,revenue,sales
Gadget,2249.25,75
```

### Example 3: Data Transformation Pipeline

**Scenario**: Transform CSV, aggregate, and export.

```javascript
// Workflow
[
  // 1. Parse input CSV
  {type: "parse", data: {input_type: "CSV"}},
  
  // 2. Filter active items
  {type: "filter", data: {condition: "item.status == 'active'"}},
  
  // 3. Transform data
  {type: "map", data: {expression: "item.price * 1.1"}},  // Add 10% markup
  
  // 4. Export to CSV
  {type: "format", data: {output_type: "CSV"}}
]
```

### Example 4: Multi-Format Export

**Scenario**: Export same data in multiple formats.

```javascript
// Workflow
[
  {id: "source", type: "variable", data: {var_op: "get"}},
  
  // Branch 1: JSON export
  {
    id: "json_export",
    type: "format",
    data: {output_type: "JSON", pretty_print: true}
  },
  
  // Branch 2: CSV export
  {
    id: "csv_export",
    type: "format",
    data: {output_type: "CSV", include_headers: true}
  },
  
  // Branch 3: TSV export
  {
    id: "tsv_export",
    type: "format",
    data: {output_type: "TSV"}
  }
]
```

### Example 5: Custom Delimiter CSV

**Scenario**: Export to pipe-delimited format.

```javascript
{
  type: "format",
  data: {
    output_type: "CSV",
    delimiter: "|",
    include_headers: true
  }
}

// Output
product|sales|revenue
Widget|100|1999.99
Gadget|75|2249.25
```

## Troubleshooting

### Issue: CSV Parsing Fails

**Symptom**: Error "invalid CSV"

**Solutions**:
- Check for malformed CSV (mismatched columns)
- Ensure consistent delimiter
- Verify headers are present
- Check for special characters in data

### Issue: Type Information Lost

**Symptom**: Numbers become strings after CSV round-trip

**Explanation**: CSV is text-based. The parser infers types but may not always match original types.

**Solution**: Use JSON for type preservation, or document expected types.

### Issue: Headers Not in Expected Order

**Symptom**: CSV headers sorted alphabetically

**Explanation**: Format node sorts headers for consistency.

**Solution**: This is by design for deterministic output. Parse back to get data in object form with correct keys.

### Issue: Empty Values in CSV

**Symptom**: Null values appear as empty strings

**Explanation**: CSV doesn't have a native null representation.

**Solution**: This is expected behavior. Empty string in CSV represents null.

### Issue: Large Arrays in CSV Cells

**Symptom**: Array values become JSON strings in CSV

**Explanation**: CSV cells are scalar values. Arrays/objects are serialized as JSON.

**Solution**: Flatten arrays before CSV export, or use JSON format.

## Performance Considerations

### Memory Usage

- **Small datasets** (< 1MB): No special handling needed
- **Medium datasets** (1-10MB): Monitor memory, consider chunking
- **Large datasets** (> 10MB): Use streaming or batch processing

### Format Choice

- **JSON**: Fast parsing, preserves types, larger size
- **CSV**: Slower parsing, type inference overhead, smaller size
- **TSV**: Similar to CSV, better for data with commas

### Optimization Tips

1. **Avoid Unnecessary Formatting**: Only format when needed for output
2. **Cache Formatted Results**: If formatting same data multiple times
3. **Use Compact JSON**: Unless readability is required
4. **Batch Process Large Files**: Split into chunks

## Security Considerations

### Input Validation

- Always validate input before parsing
- Set size limits for input strings
- Sanitize output for web display

### Malformed Data

- Parser handles malformed input gracefully
- Returns errors instead of crashing
- No code injection risks

### Data Sanitization

When exporting to CSV for Excel:
- Be aware of CSV injection (formulas starting with `=`, `+`, `-`, `@`)
- Consider escaping or quoting formula-like values

## See Also

- [Expression Syntax Guide](./EXPRESSION_SYNTAX.md)
- [Node Types Documentation](./NODE_TYPES.md)
- [Workflow Examples](../src/data/workflowExamples.ts)

## Version

This documentation applies to data format features in version 2.0+ (Phase 3 implementation).

---

For questions or feature requests, please open an issue on GitHub.
