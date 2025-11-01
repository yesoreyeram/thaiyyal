# Workflow Snapshot/Restore Feature

## Overview

The workflow snapshot/restore mechanism enables workflows to be paused, serialized to JSON, and resumed later. This is a critical feature for:

- **Long-running workflows**: Pause and resume workflows that run for hours or days
- **Crash recovery**: Recover from failures without re-executing completed nodes
- **Debugging**: Inspect and restore workflow state at any point
- **Migration**: Move workflow execution between servers
- **Testing**: Create test fixtures from real execution states

## Features

- ✅ Complete state capture (node results, variables, cache, counters)
- ✅ JSON serialization for portability
- ✅ Version-aware snapshot format
- ✅ Protection counter preservation
- ✅ Cache TTL handling
- ✅ Context variable/constant restoration
- ✅ Comprehensive test coverage

## Usage

### Basic Snapshot Creation

```go
// Create and execute a workflow
engine, err := engine.New(payloadJSON)
if err != nil {
    log.Fatal(err)
}

result, err := engine.Execute()
if err != nil {
    log.Fatal(err)
}

// Save a snapshot
snapshot, err := engine.SaveSnapshot()
if err != nil {
    log.Fatal(err)
}

// Serialize to JSON for storage
data, err := engine.SerializeSnapshot(snapshot)
if err != nil {
    log.Fatal(err)
}

// Save to file
os.WriteFile("workflow-snapshot.json", data, 0644)
```

### Restoring from Snapshot

```go
// Load snapshot from file
data, err := os.ReadFile("workflow-snapshot.json")
if err != nil {
    log.Fatal(err)
}

// Deserialize snapshot
snapshot, err := engine.DeserializeSnapshot(data)
if err != nil {
    log.Fatal(err)
}

// Restore engine state
restored, err := engine.LoadSnapshot(snapshot, nil)
if err != nil {
    log.Fatal(err)
}

// Continue execution from where it left off
result, err := restored.Execute()
if err != nil {
    log.Fatal(err)
}
```

### Convenience Method

```go
// One-step restore and execute
result, err := engine.ExecuteFromSnapshot(snapshot, nil)
if err != nil {
    log.Fatal(err)
}
```

## Snapshot Structure

A snapshot contains:

```json
{
  "version": "1.0.0",
  "snapshot_time": "2025-11-01T17:00:00Z",
  "workflow_id": "my-workflow",
  "execution_id": "a1b2c3d4e5f6g7h8",
  "nodes": [ ...node definitions... ],
  "edges": [ ...edge definitions... ],
  "results": {
    "node1": 42,
    "node2": "hello"
  },
  "completed_nodes": ["node1", "node2"],
  "variables": {
    "myVar": 100
  },
  "accumulator": 150,
  "counter": 5.0,
  "cache": {
    "key1": {
      "value": "cached data",
      "expiration": "2025-11-01T18:00:00Z"
    }
  },
  "context_vars": {
    "env": "production"
  },
  "context_consts": {
    "version": "1.0"
  },
  "node_execution_count": 2,
  "http_call_count": 0,
  "config": { ...engine configuration... }
}
```

## Implementation Details

### State Captured

1. **Workflow Metadata**
   - Workflow ID
   - Execution ID
   - Snapshot timestamp
   - Version number

2. **Workflow Definition**
   - Node definitions
   - Edge definitions

3. **Execution Progress**
   - Node results map
   - List of completed nodes
   - Current execution level

4. **State Manager Data**
   - Variables (name → value map)
   - Accumulator value
   - Counter value
   - Cache entries (with TTL)
   - Context variables
   - Context constants

5. **Runtime Protection**
   - Node execution counter
   - HTTP call counter

6. **Configuration**
   - Engine configuration settings

### Cache TTL Handling

When a snapshot is created, cache entries include their expiration timestamps. During restoration:

- Expired cache entries are **not** restored
- Non-expired entries are restored with **remaining TTL**
- TTL is calculated as: `original_expiration - current_time`

Example:
```
Original cache entry: expires at 18:00:00
Snapshot created at:   17:30:00  (30 minutes remaining)
Snapshot restored at:  17:45:00  (15 minutes remaining)
Restored TTL:          15 minutes
```

### Thread Safety

The snapshot mechanism is thread-safe:
- Uses read locks when accessing engine state
- Creates defensive copies of all mutable data
- No data races during concurrent access

### Version Compatibility

The current snapshot version is `1.0.0`. Future versions will:
- Maintain backward compatibility for reading old snapshots
- Include migration logic for format changes
- Provide clear error messages for unsupported versions

## Testing

Comprehensive test coverage includes:

- ✅ Basic snapshot creation
- ✅ Snapshot with state manager data
- ✅ Snapshot restoration
- ✅ Serialization/deserialization
- ✅ Convenience function
- ✅ Protection counter preservation
- ✅ Invalid version handling
- ✅ Nil snapshot handling
- ✅ Complex workflow scenarios
- ✅ Empty workflow handling
- ✅ Timestamp verification

Run tests:
```bash
cd backend
go test -v ./pkg/engine -run Snapshot
```

## Performance Considerations

### Snapshot Creation

- **Time complexity**: O(n) where n = number of nodes + state entries
- **Space complexity**: O(m) where m = total size of results and state
- **Lock duration**: Minimal (read locks only, no blocking)

### Snapshot Restoration

- **Time complexity**: O(n) for state restoration
- **Space complexity**: O(m) for new engine instance
- **Validation overhead**: Minimal (version check only)

### JSON Serialization

- **Format**: Pretty-printed JSON for human readability
- **Size**: Typically 10-100KB for small workflows
- **Compression**: Not included (can be added externally)

## Limitations

1. **Execution Resume**: Current implementation re-executes all nodes (optimization planned)
2. **Incremental Snapshots**: Not supported (full snapshots only)
3. **Binary Format**: JSON only (binary format not implemented)
4. **Compression**: Not built-in (apply externally if needed)
5. **Encryption**: Not built-in (apply externally for sensitive data)

## Future Enhancements

- [ ] Incremental execution (skip completed nodes on restore)
- [ ] Binary snapshot format for efficiency
- [ ] Built-in compression
- [ ] Snapshot diff functionality
- [ ] Snapshot versioning and rollback
- [ ] Encrypted snapshots for sensitive workflows
- [ ] Automatic snapshot intervals
- [ ] Snapshot size limits
- [ ] Snapshot garbage collection

## Related Documentation

- [Engine Documentation](README.md)
- [State Manager](../state/manager.go)
- [Types](../types/types.go)

## Examples

See `snapshot_test.go` for comprehensive usage examples.

## License

MIT License - See LICENSE file for details.
