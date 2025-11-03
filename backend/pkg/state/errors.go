package state

import "errors"

// Sentinel errors for state management operations
var (
// State access errors
ErrKeyNotFound        = errors.New("key not found in state")
ErrInvalidKey         = errors.New("invalid state key")
ErrStateEmpty         = errors.New("state is empty")

// Type errors
ErrTypeMismatch       = errors.New("state value type mismatch")
ErrInvalidValue       = errors.New("invalid state value")

// Transaction errors
ErrTransactionFailed  = errors.New("transaction failed")
ErrTransactionConflict = errors.New("transaction conflict")
ErrTransactionAborted = errors.New("transaction aborted")

// Concurrency errors
ErrConcurrentModification = errors.New("concurrent modification detected")
ErrLockTimeout        = errors.New("lock acquisition timeout")

// Storage errors
ErrStorageFailed      = errors.New("storage operation failed")
ErrPersistenceFailed  = errors.New("persistence failed")
ErrLoadFailed         = errors.New("failed to load state")

// Limit errors
ErrMemoryLimitExceeded = errors.New("state memory limit exceeded")
ErrMaxVariablesExceeded = errors.New("maximum variables exceeded")
ErrTTLExpired         = errors.New("state entry expired")
)
