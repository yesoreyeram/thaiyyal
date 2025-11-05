# Phase 4: New Node Types Implementation Plan

## Overview

**Goal**: Add specialized workflow nodes to support advanced patterns and unblock 20+ workflows

**Duration**: 3 weeks (now accelerated to 1 week based on Phase 1 & 3 velocity)

**Impact**: Unlocks 20+ workflows, increases coverage from 80% to 88%

## Status

- [ ] Part 1: RateLimiter & Throttle Nodes
- [ ] Part 2: SchemaValidator Node
- [ ] Part 3: Paginator Node (Auto-pagination)
- [ ] Part 4: Enhanced Cache & Variable Scoping
- [ ] Part 5: Documentation & Tests

## Part 1: RateLimiter & Throttle Nodes (2 days)

### RateLimiter Node

**Purpose**: Control request rates to prevent overwhelming APIs or hitting rate limits

**Node Type**: `rate_limiter`

**Configuration**:
```javascript
{
  type: "rate_limiter",
  data: {
    max_requests: 10,        // Maximum requests
    per_duration: "1s",      // Time window (1s, 1m, 1h)
    strategy: "fixed_window" // fixed_window, sliding_window, token_bucket
  }
}
```

**Behavior**:
- Tracks request counts within time windows
- Delays requests when limit exceeded
- Returns 429-like errors if hard limit hit
- Supports multiple strategies

**Use Cases**:
- API rate limit compliance
- Throttling batch operations
- Preventing resource exhaustion

**Affected Workflows**: Examples with tags: rate-limit, throttle (~5 workflows)

### Throttle Node (Alternative simpler version)

**Purpose**: Simple delay-based request throttling

**Node Type**: `throttle`

**Configuration**:
```javascript
{
  type: "throttle",
  data: {
    requests_per_second: 5  // Simple rate
  }
}
```

**Implementation**:
- Simple delay calculation: 1000ms / requests_per_second
- Lightweight alternative to full rate limiter

## Part 2: SchemaValidator Node (2 days)

### Purpose
Validate data against JSON schemas for type safety and data quality

**Node Type**: `schema_validator`

**Configuration**:
```javascript
{
  type: "schema_validator",
  data: {
    schema: {
      type: "object",
      properties: {
        name: { type: "string" },
        age: { type: "number", minimum: 0 },
        email: { type: "string", format: "email" }
      },
      required: ["name", "email"]
    },
    strict: true  // Fail on validation errors vs warn
  }
}
```

**Features**:
- JSON Schema draft-07 support
- Comprehensive type validation
- Format validation (email, uri, date-time, etc.)
- Custom error messages
- Strict vs lenient modes

**Use Cases**:
- API request validation
- Data quality checks
- Input sanitization
- Contract testing

**Affected Workflows**: Examples with tags: validation, schema (~6 workflows)

**Implementation Notes**:
- Use existing JSON schema library (github.com/xeipuuv/gojsonschema)
- Comprehensive error reporting
- Support for custom formats
- Performance optimizations for repeated validation

## Part 3: Paginator Node (2 days)

### Purpose
Automatically handle API pagination to fetch all results

**Node Type**: `paginator`

**Configuration**:
```javascript
{
  type: "paginator",
  data: {
    strategy: "offset_limit",  // offset_limit, page_number, cursor, link_header
    
    // For offset_limit strategy
    offset_param: "offset",
    limit_param: "limit",
    page_size: 100,
    max_pages: 10,  // Safety limit
    
    // For page_number strategy
    page_param: "page",
    per_page_param: "per_page",
    
    // For cursor strategy
    cursor_param: "cursor",
    next_cursor_path: "response.next_cursor",
    
    // For link_header strategy
    link_header: "Link",  // HTTP header name
    
    // Common settings
    total_count_path: "response.total",  // Optional
    results_path: "response.data"  // Where results are in response
  }
}
```

**Strategies**:

1. **Offset/Limit**: `?offset=0&limit=100`, `?offset=100&limit=100`, etc.
2. **Page Number**: `?page=1&per_page=100`, `?page=2&per_page=100`, etc.
3. **Cursor-based**: `?cursor=abc`, `?cursor=def` (from response)
4. **Link Header**: Parse RFC 5988 Link headers (`<url>; rel="next"`)

**Behavior**:
- Automatically makes subsequent requests
- Aggregates results into single array
- Stops when no more pages (empty results, no next cursor, etc.)
- Respects max_pages safety limit
- Returns combined results

**Use Cases**:
- Fetching all users from an API
- Downloading complete datasets
- API scraping/archival
- Data migration

**Affected Workflows**: Examples with tags: pagination (~5 workflows)

**Implementation Notes**:
- Integrates with HTTP node
- Configurable delay between pages
- Error handling for partial results
- Memory-efficient streaming option for large datasets

## Part 4: Enhanced Cache & Variable Scoping (2 days)

### Enhanced Cache Node

**Enhancements to existing cache node**:

1. **Eviction Policies**:
   - LRU (Least Recently Used)
   - LFU (Least Frequently Used)
   - TTL (Time To Live)
   - Size-based eviction

2. **Configuration**:
```javascript
{
  type: "cache",
  data: {
    key: "user_{{id}}",
    ttl: "5m",              // Time to live
    max_size: 100,          // Max cache entries
    eviction: "lru",        // lru, lfu, ttl
    storage: "memory"       // memory, redis (future)
  }
}
```

3. **Features**:
   - Cache invalidation
   - Cache warming
   - Cache statistics
   - Conditional caching

### Variable Scoping Enhancement

**Enhancements to existing variable node**:

1. **Scope Levels**:
   - `global`: Shared across all workflows
   - `workflow`: Shared within workflow execution
   - `local`: Local to current node execution

2. **Configuration**:
```javascript
{
  type: "variable",
  data: {
    name: "user_data",
    scope: "workflow",  // global, workflow, local
    value: "{{ expression }}"
  }
}
```

3. **Features**:
   - Scope isolation
   - Variable namespacing
   - Immutability options
   - Type enforcement

**Affected Workflows**: Examples with tags: cache, state, variable (~4 workflows)

## Part 5: Documentation & Tests (1 day)

### Documentation

1. **Node Reference Guide** (`docs/ADVANCED_NODES.md`):
   - Complete reference for all new nodes
   - Configuration options
   - Usage examples
   - Best practices
   - Performance considerations

2. **Workflow Examples**:
   - Example 28: Rate-limited API calls
   - Example 29: Data validation with schemas
   - Example 30: Auto-pagination workflow
   - Example 31: Enhanced caching patterns

### Tests

1. **Unit Tests**:
   - RateLimiter: 20+ tests
   - SchemaValidator: 25+ tests
   - Paginator: 30+ tests
   - Enhanced Cache: 15+ tests
   - Variable Scoping: 10+ tests

2. **Integration Tests**:
   - End-to-end workflow tests
   - Multi-node interaction tests
   - Error handling tests

3. **Performance Tests**:
   - Rate limiter accuracy
   - Cache performance
   - Pagination efficiency

## Success Criteria

### Functionality
- [ ] All 4 new node types implemented
- [ ] All enhancement features working
- [ ] 100+ new tests passing
- [ ] 4 new workflow examples working

### Quality
- [ ] All tests passing
- [ ] No security vulnerabilities
- [ ] Comprehensive documentation
- [ ] Performance validated

### Impact
- [ ] 20+ workflows unblocked
- [ ] Coverage increased to 88%
- [ ] Developer experience improved

## Dependencies

### External Libraries

1. **JSON Schema**: `github.com/xeipuuv/gojsonschema` (BSD-2-Clause)
   - Well-maintained, 2k+ stars
   - JSON Schema draft-07 support
   - Good performance

2. **Rate Limiting**: Implement custom (simple fixed-window)
   - Alternative: `golang.org/x/time/rate` (official Go library)

3. **Caching**: Enhance existing in-memory implementation
   - Future: Optional Redis support

### Security Considerations

1. **RateLimiter**:
   - Prevent resource exhaustion
   - Protect against abuse
   - Ensure fairness

2. **SchemaValidator**:
   - Prevent ReDoS in regex patterns
   - Limit schema complexity
   - Validate schema itself

3. **Paginator**:
   - Limit max pages to prevent infinite loops
   - Validate URLs to prevent SSRF
   - Handle large datasets safely

4. **Cache**:
   - Prevent cache poisoning
   - Validate cache keys
   - Size limits to prevent memory exhaustion

## Timeline

| Part | Days | Deliverable |
|------|------|-------------|
| 1 | 2 | RateLimiter & Throttle nodes |
| 2 | 2 | SchemaValidator node |
| 3 | 2 | Paginator node |
| 4 | 2 | Enhanced Cache & Variables |
| 5 | 1 | Documentation & Tests |
| **Total** | **9 days** | **Phase 4 Complete** |

*Note*: Based on Phase 1 & 3 velocity (50-700% faster), actual completion likely 3-5 days.

## Risk Mitigation

1. **Complexity**: Start with simple implementations, enhance iteratively
2. **Dependencies**: Minimize external dependencies, prefer stdlib
3. **Performance**: Add benchmarks early, optimize as needed
4. **Security**: Security review for each node type
5. **Testing**: Write tests first (TDD approach)

## Future Enhancements (Post-Phase 4)

1. **Advanced Rate Limiting**:
   - Token bucket algorithm
   - Leaky bucket algorithm
   - Distributed rate limiting

2. **Schema Evolution**:
   - Schema versioning
   - Migration support
   - Backward compatibility

3. **Advanced Pagination**:
   - Parallel page fetching
   - Streaming results
   - Resume capability

4. **Distributed Caching**:
   - Redis backend
   - Memcached support
   - Cache synchronization

## Appendix: Workflow Examples Mapping

### Rate Limiter (5 workflows)
- Example with tag: rate-limit
- API compliance workflows
- Bulk operation throttling

### Schema Validator (6 workflows)
- Example with tag: validation
- Data quality workflows
- API contract testing

### Paginator (5 workflows)
- Example with tag: pagination
- Data collection workflows
- API scraping

### Enhanced Cache (4 workflows)
- Example with tag: cache
- Performance optimization
- Data persistence

**Total Impact**: 20 workflows unblocked by Phase 4
