package types

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements custom JSON unmarshaling for Node with type-safe decoding
func (n *Node) UnmarshalJSON(data []byte) error {
	// First, parse to get ID and Type
	type NodeTemp struct {
		ID   string          `json:"id"`
		Type NodeType        `json:"type,omitempty"`
		Data json.RawMessage `json:"data"`
	}

	var temp NodeTemp
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal node: %w", err)
	}

	n.ID = temp.ID
	n.Type = temp.Type

	// If no data field, return early
	if len(temp.Data) == 0 || string(temp.Data) == "null" {
		return nil
	}

	// Decode data based on node type
	nodeData, err := unmarshalNodeData(temp.Type, temp.Data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data for node %s (type %s): %w", n.ID, n.Type, err)
	}

	n.Data = nodeData
	return nil
}

// unmarshalNodeData decodes the JSON data into the appropriate NodeData type
func unmarshalNodeData(nodeType NodeType, data json.RawMessage) (NodeDataInterface, error) {
	switch nodeType {
	// Basic Input Nodes
	case NodeTypeNumber:
		var d NumberData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeTextInput:
		var d TextInputData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeBooleanInput:
		var d BooleanInputData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeDateInput:
		var d DateInputData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeDateTimeInput:
		var d DateTimeInputData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	// Operation Nodes
	case NodeTypeOperation:
		var d OperationData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeTextOperation:
		var d TextOperationData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeHTTP:
		var d HTTPData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeExpression:
		var d ExpressionData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	// Control Flow Nodes
	case NodeTypeCondition:
		var d ConditionData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeForEach:
		var d ForEachData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeWhileLoop:
		var d WhileLoopData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeFilter:
		var d FilterData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeMap:
		var d MapData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeReduce:
		var d ReduceData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	// Array Processing Nodes
	case NodeTypeSlice:
		var d SliceData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeSort:
		var d SortData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeFind:
		var d FindData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeFlatMap:
		var d FlatMapData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeGroupBy:
		var d GroupByData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeUnique:
		var d UniqueData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeChunk:
		var d ChunkData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeReverse:
		var d ReverseData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypePartition:
		var d PartitionData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeZip:
		var d ZipData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeSample:
		var d SampleData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeRange:
		var d RangeData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeCompact:
		var d CompactData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeTranspose:
		var d TransposeData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	// State & Memory Nodes
	case NodeTypeVariable:
		var d VariableData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeExtract:
		var d ExtractData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeTransform:
		var d TransformData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeAccumulator:
		var d AccumulatorData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeCounter:
		var d CounterData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeParse:
		var d ParseData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeFormat:
		var d FormatData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	// Advanced Control Flow Nodes
	case NodeTypeSwitch:
		var d SwitchData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeParallel:
		var d ParallelData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeJoin:
		var d JoinData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeSplit:
		var d SplitData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeDelay:
		var d DelayData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeCache:
		var d CacheData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	// Error Handling & Resilience Nodes
	case NodeTypeRetry:
		var d RetryData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeTryCatch:
		var d TryCatchData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeTimeout:
		var d TimeoutData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	// Advanced Nodes
	case NodeTypeRateLimiter:
		var d RateLimiterData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeThrottle:
		var d ThrottleData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeSchemaValidator:
		var d SchemaValidatorData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypePaginator:
		var d PaginatorData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	// Context Nodes
	case NodeTypeContextVariable:
		var d ContextVariableData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeContextConstant:
		var d ContextConstantData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	// Visualization Nodes
	case NodeTypeVisualization:
		var d VisualizationData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	case NodeTypeRenderer:
		var d RendererData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, d.Validate()

	default:
		// Unknown node type - could be a custom executor
		// Unmarshal into generic CustomExecutorData
		var d CustomExecutorData
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		// Also store the raw fields for custom executors
		var fields map[string]interface{}
		if err := json.Unmarshal(data, &fields); err != nil {
			return nil, err
		}
		d.Fields = fields
		return d, d.Validate()
	}
}
