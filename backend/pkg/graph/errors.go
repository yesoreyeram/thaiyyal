package graph

import "errors"

// Sentinel errors for graph operations
var (
// Graph structure errors
ErrEmptyGraph         = errors.New("graph is empty")
ErrNodeNotFound       = errors.New("node not found in graph")
ErrEdgeNotFound       = errors.New("edge not found in graph")
ErrInvalidEdge        = errors.New("invalid edge")

// Cycle detection errors
ErrCycleDetected      = errors.New("cycle detected in graph")
ErrMultipleCycles     = errors.New("multiple cycles detected in graph")

// Topological sort errors
ErrNotDAG             = errors.New("graph is not a DAG")
ErrCannotSort         = errors.New("cannot perform topological sort")

// Traversal errors
ErrTraversalFailed    = errors.New("graph traversal failed")
ErrInvalidStartNode   = errors.New("invalid start node for traversal")
)
