package reaper

import "fmt"

// Error is a reaper error that can be returned from a command handler to exit with a specific status code
type Error struct {
	Message    string
	ExitStatus int
}

// NewError returns a new error
func NewError(status int, message string) Error {
	return Error{
		ExitStatus: status,
		Message:    message,
	}
}

// NewErrorf returns a new error
func NewErrorf(status int, format string, args ...interface{}) Error {
	return NewError(status, fmt.Sprintf(format, args...))
}

func (e Error) Error() string {
	return e.Message
}
