package types

import "fmt"

// Type assertion helper functions for executors
// These provide a clean API for executors to access their specific data types

// AsNumberData converts NodeDataInterface to NumberData with type checking
func AsNumberData(data NodeDataInterface) (*NumberData, error) {
	if d, ok := data.(NumberData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected NumberData, got %T", data)
}

// AsTextInputData converts NodeDataInterface to TextInputData with type checking
func AsTextInputData(data NodeDataInterface) (*TextInputData, error) {
	if d, ok := data.(TextInputData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected TextInputData, got %T", data)
}

// AsBooleanInputData converts NodeDataInterface to BooleanInputData with type checking
func AsBooleanInputData(data NodeDataInterface) (*BooleanInputData, error) {
	if d, ok := data.(BooleanInputData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected BooleanInputData, got %T", data)
}

// AsDateInputData converts NodeDataInterface to DateInputData with type checking
func AsDateInputData(data NodeDataInterface) (*DateInputData, error) {
	if d, ok := data.(DateInputData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected DateInputData, got %T", data)
}

// AsDateTimeInputData converts NodeDataInterface to DateTimeInputData with type checking
func AsDateTimeInputData(data NodeDataInterface) (*DateTimeInputData, error) {
	if d, ok := data.(DateTimeInputData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected DateTimeInputData, got %T", data)
}

// AsOperationData converts NodeDataInterface to OperationData with type checking
func AsOperationData(data NodeDataInterface) (*OperationData, error) {
	if d, ok := data.(OperationData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected OperationData, got %T", data)
}

// AsTextOperationData converts NodeDataInterface to TextOperationData with type checking
func AsTextOperationData(data NodeDataInterface) (*TextOperationData, error) {
	if d, ok := data.(TextOperationData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected TextOperationData, got %T", data)
}

// AsHTTPData converts NodeDataInterface to HTTPData with type checking
func AsHTTPData(data NodeDataInterface) (*HTTPData, error) {
	if d, ok := data.(HTTPData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected HTTPData, got %T", data)
}

// AsExpressionData converts NodeDataInterface to ExpressionData with type checking
func AsExpressionData(data NodeDataInterface) (*ExpressionData, error) {
	if d, ok := data.(ExpressionData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ExpressionData, got %T", data)
}

// AsConditionData converts NodeDataInterface to ConditionData with type checking
func AsConditionData(data NodeDataInterface) (*ConditionData, error) {
	if d, ok := data.(ConditionData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ConditionData, got %T", data)
}

// AsForEachData converts NodeDataInterface to ForEachData with type checking
func AsForEachData(data NodeDataInterface) (*ForEachData, error) {
	if d, ok := data.(ForEachData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ForEachData, got %T", data)
}

// AsWhileLoopData converts NodeDataInterface to WhileLoopData with type checking
func AsWhileLoopData(data NodeDataInterface) (*WhileLoopData, error) {
	if d, ok := data.(WhileLoopData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected WhileLoopData, got %T", data)
}

// AsFilterData converts NodeDataInterface to FilterData with type checking
func AsFilterData(data NodeDataInterface) (*FilterData, error) {
	if d, ok := data.(FilterData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected FilterData, got %T", data)
}

// AsMapData converts NodeDataInterface to MapData with type checking
func AsMapData(data NodeDataInterface) (*MapData, error) {
	if d, ok := data.(MapData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected MapData, got %T", data)
}

// AsReduceData converts NodeDataInterface to ReduceData with type checking
func AsReduceData(data NodeDataInterface) (*ReduceData, error) {
	if d, ok := data.(ReduceData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ReduceData, got %T", data)
}

// AsSliceData converts NodeDataInterface to SliceData with type checking
func AsSliceData(data NodeDataInterface) (*SliceData, error) {
	if d, ok := data.(SliceData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected SliceData, got %T", data)
}

// AsSortData converts NodeDataInterface to SortData with type checking
func AsSortData(data NodeDataInterface) (*SortData, error) {
	if d, ok := data.(SortData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected SortData, got %T", data)
}

// AsFindData converts NodeDataInterface to FindData with type checking
func AsFindData(data NodeDataInterface) (*FindData, error) {
	if d, ok := data.(FindData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected FindData, got %T", data)
}

// AsFlatMapData converts NodeDataInterface to FlatMapData with type checking
func AsFlatMapData(data NodeDataInterface) (*FlatMapData, error) {
	if d, ok := data.(FlatMapData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected FlatMapData, got %T", data)
}

// AsGroupByData converts NodeDataInterface to GroupByData with type checking
func AsGroupByData(data NodeDataInterface) (*GroupByData, error) {
	if d, ok := data.(GroupByData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected GroupByData, got %T", data)
}

// AsUniqueData converts NodeDataInterface to UniqueData with type checking
func AsUniqueData(data NodeDataInterface) (*UniqueData, error) {
	if d, ok := data.(UniqueData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected UniqueData, got %T", data)
}

// AsChunkData converts NodeDataInterface to ChunkData with type checking
func AsChunkData(data NodeDataInterface) (*ChunkData, error) {
	if d, ok := data.(ChunkData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ChunkData, got %T", data)
}

// AsReverseData converts NodeDataInterface to ReverseData with type checking
func AsReverseData(data NodeDataInterface) (*ReverseData, error) {
	if d, ok := data.(ReverseData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ReverseData, got %T", data)
}

// AsPartitionData converts NodeDataInterface to PartitionData with type checking
func AsPartitionData(data NodeDataInterface) (*PartitionData, error) {
	if d, ok := data.(PartitionData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected PartitionData, got %T", data)
}

// AsZipData converts NodeDataInterface to ZipData with type checking
func AsZipData(data NodeDataInterface) (*ZipData, error) {
	if d, ok := data.(ZipData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ZipData, got %T", data)
}

// AsSampleData converts NodeDataInterface to SampleData with type checking
func AsSampleData(data NodeDataInterface) (*SampleData, error) {
	if d, ok := data.(SampleData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected SampleData, got %T", data)
}

// AsRangeData converts NodeDataInterface to RangeData with type checking
func AsRangeData(data NodeDataInterface) (*RangeData, error) {
	if d, ok := data.(RangeData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected RangeData, got %T", data)
}

// AsCompactData converts NodeDataInterface to CompactData with type checking
func AsCompactData(data NodeDataInterface) (*CompactData, error) {
	if d, ok := data.(CompactData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected CompactData, got %T", data)
}

// AsTransposeData converts NodeDataInterface to TransposeData with type checking
func AsTransposeData(data NodeDataInterface) (*TransposeData, error) {
	if d, ok := data.(TransposeData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected TransposeData, got %T", data)
}

// AsVariableData converts NodeDataInterface to VariableData with type checking
func AsVariableData(data NodeDataInterface) (*VariableData, error) {
	if d, ok := data.(VariableData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected VariableData, got %T", data)
}

// AsExtractData converts NodeDataInterface to ExtractData with type checking
func AsExtractData(data NodeDataInterface) (*ExtractData, error) {
	if d, ok := data.(ExtractData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ExtractData, got %T", data)
}

// AsTransformData converts NodeDataInterface to TransformData with type checking
func AsTransformData(data NodeDataInterface) (*TransformData, error) {
	if d, ok := data.(TransformData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected TransformData, got %T", data)
}

// AsAccumulatorData converts NodeDataInterface to AccumulatorData with type checking
func AsAccumulatorData(data NodeDataInterface) (*AccumulatorData, error) {
	if d, ok := data.(AccumulatorData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected AccumulatorData, got %T", data)
}

// AsCounterData converts NodeDataInterface to CounterData with type checking
func AsCounterData(data NodeDataInterface) (*CounterData, error) {
	if d, ok := data.(CounterData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected CounterData, got %T", data)
}

// AsParseData converts NodeDataInterface to ParseData with type checking
func AsParseData(data NodeDataInterface) (*ParseData, error) {
	if d, ok := data.(ParseData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ParseData, got %T", data)
}

// AsFormatData converts NodeDataInterface to FormatData with type checking
func AsFormatData(data NodeDataInterface) (*FormatData, error) {
	if d, ok := data.(FormatData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected FormatData, got %T", data)
}

// AsSwitchData converts NodeDataInterface to SwitchData with type checking
func AsSwitchData(data NodeDataInterface) (*SwitchData, error) {
	if d, ok := data.(SwitchData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected SwitchData, got %T", data)
}

// AsParallelData converts NodeDataInterface to ParallelData with type checking
func AsParallelData(data NodeDataInterface) (*ParallelData, error) {
	if d, ok := data.(ParallelData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ParallelData, got %T", data)
}

// AsJoinData converts NodeDataInterface to JoinData with type checking
func AsJoinData(data NodeDataInterface) (*JoinData, error) {
	if d, ok := data.(JoinData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected JoinData, got %T", data)
}

// AsSplitData converts NodeDataInterface to SplitData with type checking
func AsSplitData(data NodeDataInterface) (*SplitData, error) {
	if d, ok := data.(SplitData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected SplitData, got %T", data)
}

// AsDelayData converts NodeDataInterface to DelayData with type checking
func AsDelayData(data NodeDataInterface) (*DelayData, error) {
	if d, ok := data.(DelayData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected DelayData, got %T", data)
}

// AsCacheData converts NodeDataInterface to CacheData with type checking
func AsCacheData(data NodeDataInterface) (*CacheData, error) {
	if d, ok := data.(CacheData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected CacheData, got %T", data)
}

// AsRetryData converts NodeDataInterface to RetryData with type checking
func AsRetryData(data NodeDataInterface) (*RetryData, error) {
	if d, ok := data.(RetryData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected RetryData, got %T", data)
}

// AsTryCatchData converts NodeDataInterface to TryCatchData with type checking
func AsTryCatchData(data NodeDataInterface) (*TryCatchData, error) {
	if d, ok := data.(TryCatchData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected TryCatchData, got %T", data)
}

// AsTimeoutData converts NodeDataInterface to TimeoutData with type checking
func AsTimeoutData(data NodeDataInterface) (*TimeoutData, error) {
	if d, ok := data.(TimeoutData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected TimeoutData, got %T", data)
}

// AsRateLimiterData converts NodeDataInterface to RateLimiterData with type checking
func AsRateLimiterData(data NodeDataInterface) (*RateLimiterData, error) {
	if d, ok := data.(RateLimiterData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected RateLimiterData, got %T", data)
}

// AsThrottleData converts NodeDataInterface to ThrottleData with type checking
func AsThrottleData(data NodeDataInterface) (*ThrottleData, error) {
	if d, ok := data.(ThrottleData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ThrottleData, got %T", data)
}

// AsSchemaValidatorData converts NodeDataInterface to SchemaValidatorData with type checking
func AsSchemaValidatorData(data NodeDataInterface) (*SchemaValidatorData, error) {
	if d, ok := data.(SchemaValidatorData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected SchemaValidatorData, got %T", data)
}

// AsPaginatorData converts NodeDataInterface to PaginatorData with type checking
func AsPaginatorData(data NodeDataInterface) (*PaginatorData, error) {
	if d, ok := data.(PaginatorData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected PaginatorData, got %T", data)
}

// AsContextVariableData converts NodeDataInterface to ContextVariableData with type checking
func AsContextVariableData(data NodeDataInterface) (*ContextVariableData, error) {
	if d, ok := data.(ContextVariableData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ContextVariableData, got %T", data)
}

// AsContextConstantData converts NodeDataInterface to ContextConstantData with type checking
func AsContextConstantData(data NodeDataInterface) (*ContextConstantData, error) {
	if d, ok := data.(ContextConstantData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected ContextConstantData, got %T", data)
}

// AsVisualizationData converts NodeDataInterface to VisualizationData with type checking
func AsVisualizationData(data NodeDataInterface) (*VisualizationData, error) {
	if d, ok := data.(VisualizationData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected VisualizationData, got %T", data)
}

// AsRendererData converts NodeDataInterface to RendererData with type checking
func AsRendererData(data NodeDataInterface) (*RendererData, error) {
	if d, ok := data.(RendererData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected RendererData, got %T", data)
}

// AsCustomExecutorData converts NodeDataInterface to CustomExecutorData with type checking
func AsCustomExecutorData(data NodeDataInterface) (*CustomExecutorData, error) {
	if d, ok := data.(CustomExecutorData); ok {
		return &d, nil
	}
	return nil, fmt.Errorf("expected CustomExecutorData, got %T", data)
}
