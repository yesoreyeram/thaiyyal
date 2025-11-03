package observer

import "errors"

// Sentinel errors for observer operations
var (
// Observer errors
ErrObserverPanic      = errors.New("observer panic")
ErrObserverFailed     = errors.New("observer execution failed")
ErrInvalidObserver    = errors.New("invalid observer")

// Registration errors
ErrObserverNotFound   = errors.New("observer not found")
ErrObserverAlreadyRegistered = errors.New("observer already registered")
ErrRegistrationFailed = errors.New("observer registration failed")
)
