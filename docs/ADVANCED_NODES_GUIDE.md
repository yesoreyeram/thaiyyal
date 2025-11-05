# Advanced Workflow Nodes Guide

This guide covers advanced workflow nodes for rate limiting, request throttling, and data validation use cases.

## Overview

Phase 4 introduced three specialized node types to support advanced workflow patterns:

1. **RateLimiter**: Control request rates to comply with API limits
2. **Throttle**: Simple delay-based request spacing
3. **SchemaValidator**: Validate data against JSON schemas

These nodes enable sophisticated workflows for API integration, data quality, and performance optimization.

## Node Types

### 1. RateLimiter Node

**Purpose**: Control request rates with time-window tracking to prevent overwhelming APIs or hitting rate limits.

**Node Type**: `rate_limiter`

#### Configuration

```javascript
{
  type: "rate_limiter",
  data: {
    max_requests: 10,        // Maximum requests allowed
    per_duration: "1s",      // Time window (1s, 1m, 1h)
    strategy: "fixed_window" // Currently only fixed_window supported
  }
}
```

#### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| max_requests | number | Yes | - | Maximum number of requests allowed in the time window |
| per_duration | string | Yes | - | Time window duration (1s, 1m, 1h) |
| strategy | string | No | "fixed_window" | Rate limiting strategy (only fixed_window currently) |

#### Behavior

- Tracks request counts within time windows
- Automatically delays requests when approaching limit
- Thread-safe per-node bucket tracking
- Passes through input values with metadata
- Resets counter after time window expires

#### Use Cases

1. **API Rate Limit Compliance**
   - GitHub API: 5000 requests/hour
   - Stripe API: 100 requests/second
   - Twitter API: 300 requests/15min

2. **Resource Protection**
   - Prevent overwhelming downstream services
   - Protect database connection pools
   - Control batch operation rates

3. **Fair Resource Allocation**
   - Multi-tenant systems
   - Shared service limits
   - Queue management

#### Example: GitHub API Rate Limiting

```javascript
// Workflow: Fetch multiple GitHub repositories with rate limiting
{
  nodes: [
    {
      id: "repos",
      type: "range",
      data: { start: 1, end: 100 }
    },
    {
      id: "limit",
      type: "rate_limiter",
      data: {
        max_requests: 60,      // GitHub allows 60/hour for unauthenticated
        per_duration: "1h"
      }
    },
    {
      id: "fetch",
      type: "http",
      data: {
        url: "https://api.github.com/repositories/{{item}}",
        method: "GET"
      }
    }
  ]
}
```

### 2. Throttle Node

**Purpose**: Simple delay-based request throttling for consistent request spacing.

**Node Type**: `throttle`

#### Configuration

```javascript
{
  type: "throttle",
  data: {
    requests_per_second: 5  // Target request rate
  }
}
```

#### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| requests_per_second | number | Yes | - | Target number of requests per second |

#### Behavior

- Calculates minimum delay between requests
- Ensures consistent spacing (delay = 1000ms / requests_per_second)
- Lightweight alternative to full rate limiting
- Simpler than RateLimiter for basic throttling needs
- Passes through input values unchanged

#### Use Cases

1. **Simple Request Spacing**
   - Background job processing
   - Polling with fixed intervals
   - Bulk operations

2. **Resource-Friendly Operations**
   - Web scraping with respect
   - Bulk email sending
   - Database batch operations

3. **Consistent Load Distribution**
   - Avoid traffic spikes
   - Smooth resource usage
   - Predictable performance

#### Example: Throttled Web Scraping

```javascript
// Workflow: Scrape multiple pages with throttling
{
  nodes: [
    {
      id: "urls",
      type: "variable",
      data: {
        value: [
          "https://example.com/page1",
          "https://example.com/page2",
          "https://example.com/page3"
        ]
      }
    },
    {
      id: "throttle",
      type: "throttle",
      data: {
        requests_per_second: 2  // 2 requests per second = 500ms between requests
      }
    },
    {
      id: "fetch",
      type: "http",
      data: {
        url: "{{item}}",
        method: "GET"
      }
    }
  ]
}
```

### 3. SchemaValidator Node

**Purpose**: Validate data against JSON Schema draft-07 for type safety and data quality.

**Node Type**: `schema_validator`

#### Configuration

```javascript
{
  type: "schema_validator",
  data: {
    schema: {
      type: "object",
      properties: {
        name: {
          type: "string",
          minLength: 1,
          maxLength: 100
        },
        age: {
          type: "number",
          minimum: 0,
          maximum: 150
        },
        email: {
          type: "string",
          pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
        },
        role: {
          type: "string",
          enum: ["admin", "user", "guest"]
        }
      },
      required: ["name", "email"]
    },
    strict: false  // true = fail on errors, false = return errors as metadata
  }
}
```

#### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| schema | object | Yes | - | JSON Schema draft-07 definition |
| strict | boolean | No | false | Fail on validation errors (true) or return errors (false) |

#### JSON Schema Features Supported

**Type Validation**:
- `string`, `number`, `integer`, `boolean`, `object`, `array`, `null`

**String Constraints**:
- `minLength`, `maxLength`: String length limits
- `pattern`: Regular expression validation
- `enum`: Allowed values list
- `format`: Built-in formats (email, uri, date-time, etc.)

**Number Constraints**:
- `minimum`, `maximum`: Value range
- `exclusiveMinimum`, `exclusiveMaximum`: Exclusive bounds
- `multipleOf`: Must be multiple of value

**Array Constraints**:
- `minItems`, `maxItems`: Array length limits
- `items`: Item schema validation
- `uniqueItems`: Require unique elements

**Object Constraints**:
- `properties`: Property schemas
- `required`: Required property list
- `additionalProperties`: Allow extra properties
- `minProperties`, `maxProperties`: Property count limits

**Advanced Features**:
- Nested object validation
- Array of objects validation
- Conditional schemas (if/then/else)
- Schema composition (allOf, anyOf, oneOf, not)

#### Behavior

**Strict Mode (`strict: true`)**:
- Validation errors cause node execution to fail
- Returns error immediately
- Use for critical validation where invalid data cannot proceed

**Lenient Mode (`strict: false`)**:
- Validation errors returned as metadata
- Execution continues with original data
- Errors available in `validation_errors` field
- Use for data quality checks where errors are logged but don't block

#### Use Cases

1. **API Request Validation**
   - Validate incoming request bodies
   - Ensure required fields present
   - Type safety for API contracts

2. **Data Quality Checks**
   - Validate imported data
   - Check data completeness
   - Enforce business rules

3. **Form Validation**
   - User registration validation
   - Profile update validation
   - Configuration validation

4. **Contract Testing**
   - API response validation
   - Integration testing
   - Schema evolution testing

#### Example: User Registration Validation

```javascript
// Workflow: Validate user registration data
{
  nodes: [
    {
      id: "input",
      type: "variable",
      data: {
        value: {
          name: "Alice Johnson",
          email: "alice@example.com",
          age: 25,
          role: "user"
        }
      }
    },
    {
      id: "validate",
      type: "schema_validator",
      data: {
        schema: {
          type: "object",
          properties: {
            name: {
              type: "string",
              minLength: 2,
              maxLength: 100
            },
            email: {
              type: "string",
              pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
            },
            age: {
              type: "number",
              minimum: 18,
              maximum: 120
            },
            role: {
              type: "string",
              enum: ["admin", "user", "guest"]
            }
          },
          required: ["name", "email"]
        },
        strict: true  // Fail if validation errors
      }
    },
    {
      id: "save",
      type: "http",
      data: {
        url: "https://api.example.com/users",
        method: "POST",
        body: "{{item}}"
      }
    }
  ]
}
```

#### Validation Error Format

When `strict: false`, validation errors are returned in this format:

```javascript
{
  validation_errors: [
    {
      field: "email",
      type: "pattern",
      description: "Does not match pattern '^[a-zA-Z0-9._%+-]+@...'",
      value: "invalid-email"
    },
    {
      field: "age",
      type: "minimum",
      description: "Must be greater than or equal to 18",
      value: 15
    }
  ],
  valid: false,
  data: { /* original data */ }
}
```

## Integration Patterns

### Pattern 1: Rate-Limited API Batch Processing

Combine RateLimiter with batch operations for safe API processing:

```javascript
{
  nodes: [
    // Generate batch of items
    { id: "items", type: "range", data: { start: 1, end: 100 } },
    
    // Apply rate limit
    {
      id: "limit",
      type: "rate_limiter",
      data: { max_requests: 10, per_duration: "1s" }
    },
    
    // Process each item
    {
      id: "process",
      type: "http",
      data: {
        url: "https://api.example.com/items/{{item}}",
        method: "POST"
      }
    }
  ]
}
```

### Pattern 2: Validated Data Pipeline

Validate data at each stage of processing:

```javascript
{
  nodes: [
    // Parse CSV data
    { id: "parse", type: "parse", data: { format: "CSV" } },
    
    // Validate each record
    {
      id: "validate",
      type: "schema_validator",
      data: {
        schema: { /* user schema */ },
        strict: false  // Continue with errors for logging
      }
    },
    
    // Filter valid records
    {
      id: "filter",
      type: "filter",
      data: { expression: "item.valid == true" }
    },
    
    // Export to JSON
    {
      id: "export",
      type: "format",
      data: { output_type: "JSON", pretty_print: true }
    }
  ]
}
```

### Pattern 3: Multi-API Coordination with Throttling

Coordinate multiple APIs with different rate limits:

```javascript
{
  nodes: [
    // Get user IDs
    { id: "users", type: "http", data: { url: "/api/users" } },
    
    // Throttle GitHub API (60/hour = 1/min)
    {
      id: "throttle_github",
      type: "throttle",
      data: { requests_per_second: 0.0166 }  // 1 per minute
    },
    
    // Fetch GitHub data
    {
      id: "github",
      type: "http",
      data: { url: "https://api.github.com/users/{{item.github_username}}" }
    },
    
    // Rate limit internal API (100/sec)
    {
      id: "limit_internal",
      type: "rate_limiter",
      data: { max_requests: 100, per_duration: "1s" }
    },
    
    // Save to internal DB
    {
      id: "save",
      type: "http",
      data: {
        url: "/api/profiles",
        method: "POST",
        body: "{{item}}"
      }
    }
  ]
}
```

## Real-World Examples

### Example 1: GitHub Repository Scraper

Fetch multiple repositories while respecting rate limits:

```javascript
{
  nodes: [
    {
      id: "repo_ids",
      type: "range",
      data: { start: 1, end: 1000 }
    },
    {
      id: "rate_limit",
      type: "rate_limiter",
      data: {
        max_requests: 60,
        per_duration: "1h"  // GitHub: 60 requests/hour
      }
    },
    {
      id: "fetch",
      type: "http",
      data: {
        url: "https://api.github.com/repositories/{{item}}",
        method: "GET",
        headers: {
          "Accept": "application/vnd.github.v3+json"
        }
      }
    },
    {
      id: "format",
      type: "format",
      data: { output_type: "JSON", pretty_print: true }
    }
  ]
}
```

### Example 2: E-commerce Product Validation

Validate product data before import:

```javascript
{
  nodes: [
    {
      id: "parse",
      type: "parse",
      data: { format: "CSV" }
    },
    {
      id: "validate",
      type: "schema_validator",
      data: {
        schema: {
          type: "object",
          properties: {
            sku: { type: "string", pattern: "^[A-Z0-9-]+$" },
            name: { type: "string", minLength: 1, maxLength: 200 },
            price: { type: "number", minimum: 0.01 },
            quantity: { type: "integer", minimum: 0 },
            category: {
              type: "string",
              enum: ["electronics", "clothing", "books", "home"]
            }
          },
          required: ["sku", "name", "price", "quantity"]
        },
        strict: false
      }
    },
    {
      id: "split_valid_invalid",
      type: "partition",
      data: { expression: "item.valid == true" }
    }
  ]
}
```

### Example 3: Bulk Email Sender with Throttling

Send bulk emails with consistent spacing:

```javascript
{
  nodes: [
    {
      id: "recipients",
      type: "http",
      data: { url: "/api/subscribers" }
    },
    {
      id: "throttle",
      type: "throttle",
      data: { requests_per_second: 2 }  // 2 emails per second
    },
    {
      id: "send",
      type: "http",
      data: {
        url: "https://api.sendgrid.com/v3/mail/send",
        method: "POST",
        headers: {
          "Authorization": "Bearer {{env.SENDGRID_API_KEY}}"
        },
        body: {
          personalizations: [{ to: [{ email: "{{item.email}}" }] }],
          from: { email: "noreply@example.com" },
          subject: "Your Newsletter",
          content: [{ type: "text/html", value: "{{item.content}}" }]
        }
      }
    }
  ]
}
```

### Example 4: API Response Validation

Validate API responses match expected schema:

```javascript
{
  nodes: [
    {
      id: "fetch",
      type: "http",
      data: { url: "https://api.example.com/data" }
    },
    {
      id: "validate",
      type: "schema_validator",
      data: {
        schema: {
          type: "object",
          properties: {
            status: { type: "string", enum: ["success", "error"] },
            data: {
              type: "array",
              items: {
                type: "object",
                properties: {
                  id: { type: "integer" },
                  name: { type: "string" },
                  created_at: { type: "string", format: "date-time" }
                },
                required: ["id", "name"]
              }
            }
          },
          required: ["status", "data"]
        },
        strict: true  // Fail if API response doesn't match schema
      }
    }
  ]
}
```

### Example 5: Multi-Source Data Aggregation

Aggregate data from multiple rate-limited APIs:

```javascript
{
  nodes: [
    {
      id: "users",
      type: "variable",
      data: { value: ["user1", "user2", "user3"] }
    },
    
    // Fetch from API 1 (rate limited)
    {
      id: "limit_api1",
      type: "rate_limiter",
      data: { max_requests: 10, per_duration: "1s" }
    },
    {
      id: "api1",
      type: "http",
      data: { url: "https://api1.example.com/users/{{item}}" }
    },
    
    // Fetch from API 2 (throttled)
    {
      id: "throttle_api2",
      type: "throttle",
      data: { requests_per_second: 5 }
    },
    {
      id: "api2",
      type: "http",
      data: { url: "https://api2.example.com/profiles/{{item}}" }
    },
    
    // Merge and validate
    {
      id: "merge",
      type: "transform",
      data: {
        expression: "{ user: item.api1, profile: item.api2 }"
      }
    },
    {
      id: "validate",
      type: "schema_validator",
      data: {
        schema: {
          type: "object",
          properties: {
            user: { type: "object" },
            profile: { type: "object" }
          },
          required: ["user", "profile"]
        },
        strict: false
      }
    }
  ]
}
```

## Best Practices

### 1. Choose the Right Tool

- **Use RateLimiter** when you need precise rate limit compliance (e.g., API limits)
- **Use Throttle** for simple, consistent request spacing
- **Use SchemaValidator** for data quality and type safety

### 2. Set Appropriate Limits

```javascript
// Too aggressive - may miss rate limit
{ max_requests: 100, per_duration: "1s" }  // If API allows 60/sec

// Better - leave safety margin
{ max_requests: 50, per_duration: "1s" }  // 50/sec with buffer

// Best - match API limits exactly
{ max_requests: 60, per_duration: "1s" }  // Match documented limit
```

### 3. Handle Validation Errors

```javascript
// Lenient mode for data quality checks
{
  type: "schema_validator",
  data: {
    schema: { /* ... */ },
    strict: false  // Log errors, continue processing
  }
}

// Strict mode for critical validation
{
  type: "schema_validator",
  data: {
    schema: { /* ... */ },
    strict: true  // Fail fast on errors
  }
}
```

### 4. Schema Design

```javascript
// Good - specific, meaningful constraints
{
  type: "object",
  properties: {
    email: {
      type: "string",
      pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
      maxLength: 254  // RFC 5321 limit
    },
    age: {
      type: "number",
      minimum: 0,
      maximum: 150  // Realistic limit
    }
  }
}

// Avoid - overly permissive
{
  type: "object",
  properties: {
    email: { type: "string" },  // No format validation
    age: { type: "number" }      // No range limits
  }
}
```

### 5. Combine Nodes Effectively

```javascript
// Good pattern - validate then rate limit
[
  { type: "schema_validator", data: { strict: true } },
  { type: "rate_limiter", data: { max_requests: 10, per_duration: "1s" } },
  { type: "http", data: { url: "..." } }
]

// This ensures only valid data consumes rate limit quota
```

### 6. Monitor and Log

```javascript
// Use lenient mode to collect validation errors
{
  type: "schema_validator",
  data: {
    schema: { /* ... */ },
    strict: false
  }
}

// Then analyze errors to improve data quality
```

### 7. Test Rate Limits

```javascript
// Start conservative, tune based on monitoring
{
  type: "rate_limiter",
  data: {
    max_requests: 50,     // Start at 50% of limit
    per_duration: "1s"
  }
}

// Monitor for 429 errors, adjust if needed
// Increase to 75%, then 90% if stable
```

## Troubleshooting

### RateLimiter Issues

**Problem**: Still getting 429 (Too Many Requests) errors

**Solutions**:
- Reduce `max_requests` to leave safety margin
- Check if API has multiple rate limits (per second AND per hour)
- Verify time window matches API's reset period
- Consider bursting behavior (some APIs allow bursts)

**Problem**: Workflow running slower than expected

**Solutions**:
- Check if rate limit is too conservative
- Verify `per_duration` is set correctly (1s vs 1m vs 1h)
- Consider if you need rate limiting at all for this API
- Use Throttle instead if simple spacing is sufficient

### Throttle Issues

**Problem**: Requests still too fast

**Solutions**:
- Reduce `requests_per_second` value
- Check if calculation is correct (delay = 1000ms / rps)
- Verify throttle node is in the correct position in workflow

**Problem**: Workflow too slow

**Solutions**:
- Increase `requests_per_second` value
- Consider using RateLimiter for more precise control
- Check if throttling is necessary

### SchemaValidator Issues

**Problem**: Validation always fails

**Solutions**:
- Check schema syntax (JSON Schema draft-07)
- Verify data types match (string vs number)
- Check required fields are present
- Test schema with online validator first
- Use lenient mode (`strict: false`) to see detailed errors

**Problem**: Pattern validation fails unexpectedly

**Solutions**:
- Escape special regex characters (backslash in JSON)
- Test regex pattern separately
- Check for ReDoS vulnerability in complex patterns
- Consider using simpler validation with `minLength`/`maxLength`

**Problem**: Nested object validation not working

**Solutions**:
- Ensure nested schemas are properly defined
- Check `required` fields at each level
- Verify `additionalProperties` settings
- Use JSON Schema validator tools to debug

## Performance Considerations

### RateLimiter

- **Memory**: O(1) per node (single bucket)
- **CPU**: O(1) per request (simple counter)
- **Latency**: Adds delay only when limit approached
- **Concurrency**: Thread-safe with mutex

**Optimization Tips**:
- Use per-node instances (automatic)
- Don't over-configure (more nodes = more overhead)
- Monitor actual throughput vs limits

### Throttle

- **Memory**: O(1) per node (last request time)
- **CPU**: O(1) per request (time calculation)
- **Latency**: Consistent delay between requests
- **Concurrency**: Thread-safe

**Optimization Tips**:
- Lighter than RateLimiter for simple cases
- No bucket tracking overhead
- Predictable performance

### SchemaValidator

- **Memory**: O(schema size) per node
- **CPU**: O(data size × schema complexity)
- **Latency**: Varies with data and schema complexity
- **Concurrency**: Stateless, fully concurrent

**Optimization Tips**:
- Keep schemas simple and focused
- Avoid deeply nested schemas when possible
- Cache compiled schemas (done automatically)
- Use `strict: true` to fail fast
- Consider validating only critical fields

## Security Considerations

### RateLimiter

**Threat**: Resource exhaustion via excessive requests

**Mitigations**:
- Set reasonable `max_requests` limits
- Use multiple rate limiters for layered protection
- Monitor for abuse patterns
- Implement circuit breakers for failures

### SchemaValidator

**Threat**: ReDoS (Regular Expression Denial of Service)

**Mitigations**:
- Avoid complex regex patterns in schema
- Test patterns for ReDoS vulnerability
- Set schema complexity limits
- Use timeout mechanisms

**Threat**: Schema injection

**Mitigations**:
- Don't accept schemas from untrusted sources
- Validate schemas before use
- Use predefined schemas when possible
- Limit schema size

### General

**Best Practices**:
- Always validate input data
- Use strict mode for critical validations
- Monitor for unusual patterns
- Implement proper error handling
- Log validation failures for security auditing

## Dependencies

### External Libraries

**gojsonschema** (github.com/xeipuuv/gojsonschema v1.2.0):
- License: BSD-2-Clause
- Purpose: JSON Schema draft-07 validation
- Well-maintained, 2k+ stars on GitHub
- Good performance and comprehensive feature support

## Future Enhancements

### Planned Features

1. **Advanced Rate Limiting**:
   - Token bucket algorithm
   - Leaky bucket algorithm
   - Sliding window strategy
   - Distributed rate limiting

2. **Schema Evolution**:
   - Schema versioning
   - Migration support
   - Backward compatibility checks

3. **Enhanced Validation**:
   - Custom validation functions
   - Async validation
   - Conditional validation rules

4. **Performance**:
   - Schema compilation caching
   - Validation result caching
   - Parallel validation

## Additional Resources

- [JSON Schema Documentation](https://json-schema.org/)
- [Rate Limiting Algorithms](https://en.wikipedia.org/wiki/Rate_limiting)
- [API Design Best Practices](https://restfulapi.net/)
- [ReDoS Prevention](https://owasp.org/www-community/attacks/Regular_expression_Denial_of_Service_-_ReDoS)

---

**Version**: 1.0  
**Last Updated**: 2025-11-05  
**Related Guides**: EXPRESSION_SYNTAX.md, DATA_FORMAT_GUIDE.md, WORKFLOW_EXECUTION_GUIDE.md
