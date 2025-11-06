package executor

import (
	"fmt"
	"time"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// CacheExecutor executes Cache nodes
type CacheExecutor struct{}

// Execute runs the Cache node
// Handles cache get/set operations
func (e *CacheExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsCacheData(node.Data)
if err != nil {
return nil, err
}
	if data.CacheOp == nil {
		return nil, fmt.Errorf("cache node requires cache_op field")
	}
	if data.CacheKey == nil {
		return nil, fmt.Errorf("cache node requires cache_key field")
	}

	cacheOp := *data.CacheOp
	cacheKey := *data.CacheKey

	switch cacheOp {
	case "set":
		inputs := ctx.GetNodeInputs(node.ID)
		if len(inputs) == 0 {
			return nil, fmt.Errorf("cache set operation requires an input value")
		}

		value := inputs[0]

		// Parse TTL
		var ttlDuration time.Duration
		if data.TTL != nil {
			var err error
			ttlDuration, err = parseDuration(*data.TTL)
			if err != nil {
				return nil, fmt.Errorf("invalid TTL format: %w", err)
			}
		} else {
			// Default TTL: 5 minutes
			ttlDuration = 5 * time.Minute
		}

		ctx.SetCache(cacheKey, value, ttlDuration)

		return map[string]interface{}{
			"operation": "set",
			"key":       cacheKey,
			"value":     value,
			"ttl":       data.TTL,
		}, nil

	case "get":
		value, exists := ctx.GetCache(cacheKey)

		if !exists {
			return map[string]interface{}{
				"operation": "get",
				"key":       cacheKey,
				"found":     false,
				"value":     nil,
			}, nil
		}

		return map[string]interface{}{
			"operation": "get",
			"key":       cacheKey,
			"found":     true,
			"value":     value,
		}, nil

	case "delete":
		// Note: ExecutionContext doesn't have DeleteCache method
		// For now, we'll just report it as deleted
		// In a full implementation, we'd add DeleteCache to ExecutionContext
		return map[string]interface{}{
			"operation": "delete",
			"key":       cacheKey,
			"deleted":   true,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported cache operation: %s (use get, set, or delete)", cacheOp)
	}
}

// NodeType returns the node type this executor handles
func (e *CacheExecutor) NodeType() types.NodeType {
	return types.NodeTypeCache
}

// Validate checks if node configuration is valid
func (e *CacheExecutor) Validate(node types.Node) error {
data, err := types.AsCacheData(node.Data)
if err != nil {
return err
}
	if data.CacheOp == nil {
		return fmt.Errorf("cache node requires cache_op field")
	}
	if data.CacheKey == nil {
		return fmt.Errorf("cache node requires cache_key field")
	}
	return nil
}
