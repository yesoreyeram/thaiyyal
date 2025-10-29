# HTTP Pagination Node Documentation

## Overview

The HTTP Pagination node enables efficient retrieval of paginated data from REST APIs. It automatically handles multiple page requests, consolidates results, and provides robust error handling.

## Features

- **Flexible Configuration**: Support for both max pages and total items/page size
- **URL Placeholder Support**: Use `{page}` or custom placeholders in URLs
- **Custom Page Parameters**: Configure query parameter names
- **Error Handling**: Break on first error or continue collecting partial results
- **Configurable Start Page**: Begin pagination from any page number
- **Result Consolidation**: Automatically consolidates all page responses

## Configuration

### Basic Configuration

#### Base URL (Required)
The base URL for the API endpoint. Can include:
- Query parameters: `https://api.example.com/items`
- URL placeholders: `https://api.example.com/items/{page}`
- Custom placeholders: `https://api.example.com/page-{page}/items`

#### Pagination Mode

**Option 1: Max Pages**
- **max_pages**: Number of pages to fetch (e.g., 5)
- Use when you know exactly how many pages to retrieve

**Option 2: Total Items + Page Size**
- **total_items**: Total number of items to retrieve (e.g., 50)
- **page_size**: Number of items per page (e.g., 10)
- The node automatically calculates pages: `⌈total_items / page_size⌉`

### Advanced Options

#### Start Page (Optional)
- Default: `1`
- The page number to start pagination from
- Example: Start from page 5 to skip first 4 pages

#### Page Parameter (Optional)
- Default: `"page"`
- The query parameter name for the page number
- Example: Use `"p"` for URLs like `?p=1`

#### Break on Error (Optional)
- Default: `true`
- **true**: Stop pagination on first error and return error
- **false**: Continue fetching remaining pages and collect partial results

## Output Format

The node returns a comprehensive result object:

```json
{
  "success": true,
  "pages_fetched": 5,
  "total_pages": 5,
  "results": [
    "{\"items\": [\"item1\", \"item2\"]}",
    "{\"items\": [\"item3\", \"item4\"]}",
    "{\"items\": [\"item5\", \"item6\"]}",
    "{\"items\": [\"item7\", \"item8\"]}",
    "{\"items\": [\"item9\", \"item10\"]}"
  ],
  "errors": [],
  "error_count": 0
}
```

### Output Fields

- **success**: Boolean indicating if all pages were fetched without errors
- **pages_fetched**: Number of successfully fetched pages
- **total_pages**: Total number of pages attempted
- **results**: Array of response bodies (as strings)
- **errors**: Array of error messages for failed pages
- **error_count**: Number of pages that failed

## Usage Examples

### Example 1: Simple Pagination with Max Pages

Fetch 5 pages from an API:

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "http_pagination",
      "data": {
        "base_url": "https://api.example.com/items",
        "max_pages": 5
      }
    }
  ],
  "edges": []
}
```

This will make requests to:
- `https://api.example.com/items?page=1`
- `https://api.example.com/items?page=2`
- `https://api.example.com/items?page=3`
- `https://api.example.com/items?page=4`
- `https://api.example.com/items?page=5`

### Example 2: Using Total Items and Page Size

Fetch 50 items with 10 items per page:

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "http_pagination",
      "data": {
        "base_url": "https://api.example.com/users",
        "total_items": 50,
        "page_size": 10
      }
    }
  ],
  "edges": []
}
```

This automatically calculates 5 pages (50 ÷ 10 = 5).

### Example 3: URL with Placeholder

Use `{page}` placeholder in the URL path:

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "http_pagination",
      "data": {
        "base_url": "https://api.example.com/items/{page}",
        "max_pages": 3
      }
    }
  ],
  "edges": []
}
```

This will make requests to:
- `https://api.example.com/items/1`
- `https://api.example.com/items/2`
- `https://api.example.com/items/3`

### Example 4: Custom Start Page

Start pagination from page 3:

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "http_pagination",
      "data": {
        "base_url": "https://api.example.com/data",
        "start_page": 3,
        "max_pages": 2
      }
    }
  ],
  "edges": []
}
```

This will fetch pages 3 and 4.

### Example 5: Custom Page Parameter

Use a custom query parameter name:

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "http_pagination",
      "data": {
        "base_url": "https://api.example.com/search",
        "page_param": "p",
        "max_pages": 3
      }
    }
  ],
  "edges": []
}
```

This will make requests to:
- `https://api.example.com/search?p=1`
- `https://api.example.com/search?p=2`
- `https://api.example.com/search?p=3`

### Example 6: Continue on Error

Continue fetching even if some pages fail:

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "http_pagination",
      "data": {
        "base_url": "https://api.example.com/items",
        "max_pages": 5,
        "break_on_error": false
      }
    }
  ],
  "edges": []
}
```

If page 2 fails (e.g., 404 or 500 error), the node will:
- Continue fetching pages 3, 4, and 5
- Return partial results with error information
- Set `success: false` in the output
- Include error details in the `errors` array

### Example 7: Complete Workflow with Visualization

Fetch paginated data and visualize:

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "http_pagination",
      "data": {
        "base_url": "https://jsonplaceholder.typicode.com/posts",
        "max_pages": 5,
        "break_on_error": true
      }
    },
    {
      "id": "2",
      "type": "visualization",
      "data": {
        "mode": "text"
      }
    }
  ],
  "edges": [
    {
      "id": "e1-2",
      "source": "1",
      "target": "2"
    }
  ]
}
```

## Error Handling

### Break on Error (Default: true)

When `break_on_error` is `true`:
- The node stops at the first error
- Returns an error via the workflow engine
- Provides partial results in the output object
- Useful when you need complete data or nothing

Example error output:
```json
{
  "success": false,
  "error": "pagination failed on page 3: status 500",
  "pages_fetched": 2,
  "total_pages": 5,
  "results": ["page1_data", "page2_data"],
  "errors": ["page 3: HTTP request returned error status: 500"]
}
```

### Continue on Error

When `break_on_error` is `false`:
- The node continues fetching remaining pages
- Collects all successful results
- Stores errors in the `errors` array
- Returns normally (no workflow error)
- Useful when partial data is acceptable

Example output with errors:
```json
{
  "success": false,
  "pages_fetched": 4,
  "total_pages": 5,
  "results": ["page1", "page2", "page4", "page5"],
  "errors": ["page 3: HTTP request returned error status: 404"],
  "error_count": 1
}
```

## Common Error Cases

### Missing Base URL
```
Error: http_pagination node missing base_url
```

### Missing Pagination Configuration
```
Error: http_pagination node requires either max_pages or both total_items and page_size
```

### HTTP Errors
```
Error: pagination failed on page 2: HTTP request returned error status: 404
Error: pagination failed on page 3: HTTP request failed: connection timeout
```

## Best Practices

1. **Start Small**: Test with 1-2 pages before fetching large datasets
2. **Use break_on_error**: Keep default `true` for critical data
3. **Monitor Errors**: Always check the `errors` array in results
4. **Rate Limiting**: Consider adding Delay nodes between pagination requests if needed
5. **URL Placeholders**: Use `{page}` for RESTful URLs without query parameters
6. **Validate URLs**: Ensure base_url is correctly formatted before execution

## Integration with Other Nodes

### Extract Results
Use **Extract** node to get specific fields from paginated results:
```
HTTP Pagination → Extract (field: "results") → Transform
```

### Process Each Page
Use **For Each** node to process individual page results:
```
HTTP Pagination → Extract → For Each → Custom Processing
```

### Accumulate Data
Use **Accumulator** node to combine results:
```
HTTP Pagination → Extract → Transform → Accumulator (array)
```

## Comparison with Loop-based Pagination

### Using HTTP Pagination Node (Recommended)
✅ Single node
✅ Built-in error handling
✅ Automatic result consolidation
✅ Cleaner workflow visualization

### Using While Loop + HTTP Node
❌ Multiple nodes required
❌ Manual error handling
❌ Manual result accumulation
❌ More complex workflow

## Performance Considerations

- **Sequential Requests**: Pages are fetched sequentially (not in parallel)
- **Memory Usage**: All results are kept in memory
- **Large Datasets**: For very large datasets (1000+ pages), consider processing in batches
- **Timeouts**: HTTP requests use default Go HTTP client timeout

## Screenshots

### HTTP Pagination Node in Canvas
![HTTP Pagination Node](https://github.com/user-attachments/assets/84b7cb21-7baf-4a16-93a4-951d06320493)

### Advanced Options Expanded
![HTTP Pagination Advanced Options](https://github.com/user-attachments/assets/b4a18b9a-ef96-4dea-aae5-65e490c940c6)

### Generated JSON Payload
![HTTP Pagination with Payload](https://github.com/user-attachments/assets/4f19c963-a938-4398-ba56-110acbbef5f0)

## Testing

The HTTP Pagination node includes comprehensive test coverage:
- ✅ Pagination with max_pages
- ✅ Pagination with total_items and page_size
- ✅ URL placeholder support
- ✅ Custom start page
- ✅ Break on error behavior
- ✅ Continue on error behavior
- ✅ Custom page parameter names
- ✅ Missing required fields validation

See `backend/workflow_test.go` for test details.

## Related Nodes

- **HTTP**: Single HTTP request
- **For Each**: Iterate over results
- **While Loop**: Conditional iteration
- **Accumulator**: Aggregate results
- **Extract**: Extract fields from responses
- **Transform**: Transform data structures

## Support

For issues or questions:
1. Check the error message in the output
2. Verify URL formatting and parameters
3. Test with a small number of pages first
4. Review the comprehensive test suite for examples
