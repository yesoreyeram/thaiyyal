package executor

import (
	"testing"

	"github.com/yesoreyeram/thaiyyal/backend/pkg/types"
)

func TestRendererExecutor_NodeType(t *testing.T) {
	executor := &RendererExecutor{}
	if executor.NodeType() != types.NodeTypeRenderer {
		t.Errorf("Expected type %v, got %v", types.NodeTypeRenderer, executor.NodeType())
	}
}

func TestRendererExecutor_Validate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Valid node - no configuration needed",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &RendererExecutor{}
			node := types.Node{
				Data: types.NodeData{},
			}
			err := executor.Validate(node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRendererExecutor_Execute(t *testing.T) {
	tests := []struct {
		name    string
		inputs  []interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Pass through string data",
			inputs:  []interface{}{"Hello World"},
			want:    "Hello World",
			wantErr: false,
		},
		{
			name:    "Pass through object data",
			inputs:  []interface{}{map[string]interface{}{"name": "John", "age": 30}},
			want:    map[string]interface{}{"name": "John", "age": 30},
			wantErr: false,
		},
		{
			name: "Pass through array data",
			inputs: []interface{}{[]interface{}{
				map[string]interface{}{"label": "A", "value": 10},
				map[string]interface{}{"label": "B", "value": 20},
			}},
			want: []interface{}{
				map[string]interface{}{"label": "A", "value": 10},
				map[string]interface{}{"label": "B", "value": 20},
			},
			wantErr: false,
		},
		{
			name:    "Pass through number data",
			inputs:  []interface{}{42},
			want:    42,
			wantErr: false,
		},
		{
			name:    "No input returns nil",
			inputs:  []interface{}{},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Multiple inputs - uses first",
			inputs:  []interface{}{"first", "second", "third"},
			want:    "first",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &RendererExecutor{}
			node := types.Node{
				ID:   "renderer1",
				Data: types.NodeData{},
			}

			ctx := &MockExecutionContext{
				inputs: map[string][]interface{}{
					"renderer1": tt.inputs,
				},
			}

			result, err := executor.Execute(ctx, node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Deep comparison for complex types
			if !compareResults(result, tt.want) {
				t.Errorf("Execute() result = %v, want %v", result, tt.want)
			}
		})
	}
}

// Helper function to compare results deeply
func compareResults(got, want interface{}) bool {
	if got == nil && want == nil {
		return true
	}
	if got == nil || want == nil {
		return false
	}

	// For basic types, use direct comparison
	switch v := want.(type) {
	case string, int, float64, bool:
		return got == v
	case map[string]interface{}:
		gotMap, ok := got.(map[string]interface{})
		if !ok {
			return false
		}
		if len(gotMap) != len(v) {
			return false
		}
		for k, wantV := range v {
			gotV, exists := gotMap[k]
			if !exists || !compareResults(gotV, wantV) {
				return false
			}
		}
		return true
	case []interface{}:
		gotSlice, ok := got.([]interface{})
		if !ok {
			return false
		}
		if len(gotSlice) != len(v) {
			return false
		}
		for i, wantV := range v {
			if !compareResults(gotSlice[i], wantV) {
				return false
			}
		}
		return true
	default:
		// For other types, fall back to direct comparison
		return got == want
	}
}
