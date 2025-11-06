package executor

import (
	"fmt"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// RangeExecutor generates an array of numbers
type RangeExecutor struct{}

// Execute generates a range of numbers
func (e *RangeExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
data, err := types.AsRangeData(node.Data)
if err != nil {
return nil, err
}
	// Get start
	start := 0.0
	if startVal, ok := data.Start.(float64); ok {
		start = startVal
	} else if startVal, ok := data.Start.(int); ok {
		start = float64(startVal)
	}

	// Get end
	end := 10.0 // default
	if endVal, ok := data.End.(float64); ok {
		end = endVal
	} else if endVal, ok := data.End.(int); ok {
		end = float64(endVal)
	}

	// Get step
	step := 1.0 // default
	if stepVal, ok := data.Step.(float64); ok {
		step = stepVal
	} else if stepVal, ok := data.Step.(int); ok {
		step = float64(stepVal)
	}

	if step == 0 {
		return nil, fmt.Errorf("range step cannot be 0")
	}

	// Validate range direction
	if step > 0 && start > end {
		return nil, fmt.Errorf("range with positive step requires start <= end, got start=%v end=%v", start, end)
	}
	if step < 0 && start < end {
		return nil, fmt.Errorf("range with negative step requires start >= end, got start=%v end=%v", start, end)
	}

	// Generate range
	var rangeArr []interface{}
	maxItems := 10000 // Safety limit

	if step > 0 {
		for i := start; i <= end && len(rangeArr) < maxItems; i += step {
			rangeArr = append(rangeArr, i)
		}
	} else {
		for i := start; i >= end && len(rangeArr) < maxItems; i += step {
			rangeArr = append(rangeArr, i)
		}
	}

	if len(rangeArr) >= maxItems {
		return nil, fmt.Errorf("range would generate more than %d items, adjust start/end/step", maxItems)
	}

	return map[string]interface{}{
		"range": rangeArr,
		"count": len(rangeArr),
		"start": start,
		"end":   end,
		"step":  step,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *RangeExecutor) NodeType() types.NodeType {
	return types.NodeTypeRange
}

// Validate checks if the node configuration is valid
func (e *RangeExecutor) Validate(node types.Node) error {
data, err := types.AsRangeData(node.Data)
if err != nil {
return err
}
	// Get step for validation
	step := 1.0
	if stepVal, ok := data.Step.(float64); ok {
		step = stepVal
	} else if stepVal, ok := data.Step.(int); ok {
		step = float64(stepVal)
	}

	if step == 0 {
		return fmt.Errorf("range step cannot be 0")
	}

	// Get start and end for validation
	start := 0.0
	if startVal, ok := data.Start.(float64); ok {
		start = startVal
	} else if startVal, ok := data.Start.(int); ok {
		start = float64(startVal)
	}

	end := 10.0
	if endVal, ok := data.End.(float64); ok {
		end = endVal
	} else if endVal, ok := data.End.(int); ok {
		end = float64(endVal)
	}

	// Validate direction
	if step > 0 && start > end {
		return fmt.Errorf("range with positive step requires start <= end")
	}
	if step < 0 && start < end {
		return fmt.Errorf("range with negative step requires start >= end")
	}

	return nil
}
