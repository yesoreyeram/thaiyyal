package workflow

import "fmt"

// ============================================================================
// Node Execution Dispatcher
// ============================================================================
// This file contains the main dispatcher that routes node execution to
// appropriate executor functions based on node type.
// Uses the Strategy Pattern for extensible node type handling.
// ============================================================================

// executeNode dispatches node execution to the appropriate executor based on node type.
// This is the central routing function that implements the Strategy Pattern,
// allowing different execution strategies for different node types.
//
// Supported node types (23 total):
//
// Basic I/O:
//   - number, text_input, visualization
//
// Operations:
//   - operation (arithmetic), text_operation, http
//
// Control Flow:
//   - condition, for_each, while_loop
//
// State & Memory:
//   - variable, extract, transform, accumulator, counter
//
// Advanced Control Flow:
//   - switch, parallel, join, split, delay, cache
//
// Error Handling & Resilience:
//   - retry, try_catch, timeout
//
// Returns:
//   - interface{}: Result of node execution (type depends on node)
//   - error: If node type unknown or execution fails
func (e *Engine) executeNode(node Node) (interface{}, error) {
	// Interpolate templates in node data before execution (except for context nodes)
	if node.Type != NodeTypeContextVariable && node.Type != NodeTypeContextConstant {
		e.interpolateNodeData(&node.Data)
	}

	// Dispatch to appropriate executor based on node type
	// This switch implements the Strategy Pattern
	switch node.Type {
	// Context nodes (executed first to populate context)
	case NodeTypeContextVariable:
		return e.executeContextVariableNode(node)
	case NodeTypeContextConstant:
		return e.executeContextConstantNode(node)
	// Basic I/O nodes
	case NodeTypeNumber:
		return e.executeNumberNode(node)
	case NodeTypeTextInput:
		return e.executeTextInputNode(node)
	case NodeTypeVisualization:
		return e.executeVisualizationNode(node)

	// Operation nodes
	case NodeTypeOperation:
		return e.executeOperationNode(node)
	case NodeTypeTextOperation:
		return e.executeTextOperationNode(node)
	case NodeTypeHTTP:
		return e.executeHTTPNode(node)

	// Control flow nodes
	case NodeTypeCondition:
		return e.executeConditionNode(node)
	case NodeTypeForEach:
		return e.executeForEachNode(node)
	case NodeTypeWhileLoop:
		return e.executeWhileLoopNode(node)

	// State & memory nodes
	case NodeTypeVariable:
		return e.executeVariableNode(node)
	case NodeTypeExtract:
		return e.executeExtractNode(node)
	case NodeTypeTransform:
		return e.executeTransformNode(node)
	case NodeTypeAccumulator:
		return e.executeAccumulatorNode(node)
	case NodeTypeCounter:
		return e.executeCounterNode(node)

	// Advanced control flow nodes
	case NodeTypeSwitch:
		return e.executeSwitchNode(node)
	case NodeTypeParallel:
		return e.executeParallelNode(node)
	case NodeTypeJoin:
		return e.executeJoinNode(node)
	case NodeTypeSplit:
		return e.executeSplitNode(node)
	case NodeTypeDelay:
		return e.executeDelayNode(node)
	case NodeTypeCache:
		return e.executeCacheNode(node)

	// Error handling & resilience nodes
	case NodeTypeRetry:
		inputs := e.getNodeInputs(node.ID)
		return e.executeRetryNode(&node, inputs)
	case NodeTypeTryCatch:
		inputs := e.getNodeInputs(node.ID)
		return e.executeTryCatchNode(&node, inputs)
	case NodeTypeTimeout:
		inputs := e.getNodeInputs(node.ID)
		return e.executeTimeoutNode(&node, inputs)

	default:
		return nil, fmt.Errorf("unknown node type: %s", node.Type)
	}
}

// inferNodeTypes determines node types from data if not explicitly set.
// This allows the frontend to omit node types and have them automatically detected.
//
// Type inference is based on the presence of specific fields in NodeData.
// Some nodes (for_each, while_loop, parallel) require explicit types as they
// have ambiguous fields.
func (e *Engine) inferNodeTypes() {
	for i := range e.nodes {
		if e.nodes[i].Type != "" {
			// Type already set, skip inference
			continue
		}
		
		// Infer type from data fields
		e.nodes[i].Type = inferNodeTypeFromData(e.nodes[i].Data)
	}
}

// inferNodeTypeFromData infers a node's type from its data fields.
// This implements a simple decision tree based on field presence.
//
// Returns:
//   - NodeType: Inferred type, or empty string if cannot infer
func inferNodeTypeFromData(data NodeData) NodeType {
	// Basic I/O nodes (checked first as they're most common)
	if data.Value != nil {
		return NodeTypeNumber
	}
	if data.Text != nil {
		return NodeTypeTextInput
	}
	if data.Mode != nil {
		return NodeTypeVisualization
	}

	// Operation nodes
	if data.Op != nil {
		return NodeTypeOperation
	}
	if data.TextOp != nil {
		return NodeTypeTextOperation
	}
	if data.URL != nil {
		return NodeTypeHTTP
	}

	// Control flow nodes
	if data.Condition != nil {
		return NodeTypeCondition
	}

	// State & memory nodes
	if data.VarName != nil && data.VarOp != nil {
		return NodeTypeVariable
	}
	if data.Field != nil || len(data.Fields) > 0 {
		return NodeTypeExtract
	}
	if data.TransformType != nil {
		return NodeTypeTransform
	}
	if data.AccumOp != nil {
		return NodeTypeAccumulator
	}
	if data.CounterOp != nil {
		return NodeTypeCounter
	}

	// Advanced control flow nodes
	if len(data.Cases) > 0 {
		return NodeTypeSwitch
	}
	if data.JoinStrategy != nil {
		return NodeTypeJoin
	}
	if len(data.Paths) > 0 {
		return NodeTypeSplit
	}
	if data.Duration != nil {
		return NodeTypeDelay
	}
	if data.CacheOp != nil && data.CacheKey != nil {
		return NodeTypeCache
	}

	// Context nodes
	if data.ContextName != nil && data.ContextValue != nil {
		// Default to variable, frontend should specify explicitly
		// This is a best-effort inference
		return NodeTypeContextVariable
	}

	// Error handling & resilience nodes
	if data.MaxAttempts != nil || data.BackoffStrategy != nil {
		return NodeTypeRetry
	}
	if data.FallbackValue != nil || data.ContinueOnError != nil {
		return NodeTypeTryCatch
	}
	if data.Timeout != nil && data.TimeoutAction != nil {
		return NodeTypeTimeout
	}

	// Cannot infer type
	// Note: for_each, while_loop, and parallel require explicit type
	// as they have ambiguous fields
	return ""
}
