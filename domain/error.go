// Package domain provides domain-specific error handling functionality.
// It implements custom error types that satisfy the standard error interface
// and provides additional error comparison capabilities.
package domain

import "fmt"

// Error represents a domain-specific error in the system.
// It's advised that all domain layer errors "inherit" from this type.
type Error struct {
	prevErr error
	msg     string // internal error message
}

func NewError(format string, a ...any) *Error {
	return &Error{
		msg: fmt.Sprintf(format, a...),
	}
}

// NewErrorWithWrap creates a new Error that wraps another error.
func NewErrorWithWrap(err error, format string, a ...any) *Error {
	return &Error{
		prevErr: err,
		msg:     fmt.Sprintf(format, a...),
	}
}

// Error returns the error message, satisfying the error interface.
func (e *Error) Error() string {
	if e.prevErr != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.prevErr)
	}
	return e.msg
}

// Unwrap returns the wrapped error, enabling compatibility with errors.Is and errors.As.
func (e *Error) Unwrap() error {
	return e.prevErr
}
