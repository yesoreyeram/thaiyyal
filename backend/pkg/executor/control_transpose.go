package executor

import (
	"fmt"
	"log/slog"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

// TransposeExecutor transposes a 2D array (matrix)
type TransposeExecutor struct{}

// Execute transposes a 2D array
func (e *TransposeExecutor) Execute(ctx ExecutionContext, node types.Node) (interface{}, error) {
	inputs := ctx.GetNodeInputs(node.ID)
	if len(inputs) == 0 {
		return nil, fmt.Errorf("transpose node needs at least 1 input")
	}

	input := inputs[0]

	// Check if input is an array
	arr, ok := input.([]interface{})
	if !ok {
		slog.Warn("transpose node received non-array input",
			slog.String("node_id", node.ID),
			slog.String("input_type", fmt.Sprintf("%T", input)),
		)
		return map[string]interface{}{
			"error":         "input is not an array",
			"input":         input,
			"original_type": fmt.Sprintf("%T", input),
		}, nil
	}

	if len(arr) == 0 {
		return map[string]interface{}{
			"transposed": []interface{}{},
			"rows":       0,
			"cols":       0,
		}, nil
	}

	// Check if first element is an array (2D structure)
	firstRow, ok := arr[0].([]interface{})
	if !ok {
		return map[string]interface{}{
			"error":         "input is not a 2D array (array of arrays)",
			"first_element": arr[0],
			"expected":      "array of arrays",
		}, nil
	}

	rows := len(arr)
	cols := len(firstRow)

	// Validate all rows have same length
	for i, row := range arr {
		rowArr, ok := row.([]interface{})
		if !ok {
			return map[string]interface{}{
				"error":     fmt.Sprintf("row %d is not an array", i),
				"row_index": i,
				"row_value": row,
			}, nil
		}
		if len(rowArr) != cols {
			return map[string]interface{}{
				"error":         "all rows must have the same length",
				"expected_cols": cols,
				"row_index":     i,
				"actual_cols":   len(rowArr),
			}, nil
		}
	}

	// Transpose the matrix
	transposed := make([]interface{}, cols)
	for i := 0; i < cols; i++ {
		column := make([]interface{}, rows)
		for j := 0; j < rows; j++ {
			rowArr := arr[j].([]interface{})
			column[j] = rowArr[i]
		}
		transposed[i] = column
	}

	return map[string]interface{}{
		"transposed":      transposed,
		"rows":            rows,
		"cols":            cols,
		"transposed_rows": cols,
		"transposed_cols": rows,
	}, nil
}

// NodeType returns the node type this executor handles
func (e *TransposeExecutor) NodeType() types.NodeType {
	return types.NodeTypeTranspose
}

// Validate checks if the node configuration is valid
func (e *TransposeExecutor) Validate(node types.Node) error {
	// No configuration needed
	return nil
}
