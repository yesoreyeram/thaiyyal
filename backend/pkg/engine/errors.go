package engine

import "errors"

// Sentinel errors for engine operations
var (
	// Validation errors
	ErrEmptyWorkflow   = errors.New("workflow is empty")
	ErrNoNodes         = errors.New("workflow contains no nodes")
	ErrCycleDetected   = errors.New("cycle detected in workflow graph")
	ErrInvalidNodeType = errors.New("invalid node type")
	ErrMissingNodeID   = errors.New("node ID is required")
	ErrDuplicateNodeID = errors.New("duplicate node ID found")
	ErrInvalidEdge     = errors.New("invalid edge: source or target node not found")

	// Execution errors
	ErrExecutionFailed       = errors.New("workflow execution failed")
	ErrNodeExecutionFailed   = errors.New("node execution failed")
	ErrExecutionTimeout      = errors.New("execution timeout exceeded")
	ErrExecutionCanceled     = errors.New("execution was canceled")
	ErrMaxIterationsExceeded = errors.New("maximum iterations exceeded")

	// Resource errors
	ErrMaxNodesExceeded      = errors.New("maximum number of nodes exceeded")
	ErrMaxEdgesExceeded      = errors.New("maximum number of edges exceeded")
	ErrMaxExecutionsExceeded = errors.New("maximum node executions exceeded")

	// State errors
	ErrStateNotFound = errors.New("state not found")
	ErrInvalidState  = errors.New("invalid state")

	// Custom executor errors
	ErrExecutorNotFound           = errors.New("executor not found for node type")
	ErrExecutorRegistrationFailed = errors.New("failed to register executor")
)
