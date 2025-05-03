package errors

import "fmt"

type Error interface {
	Error() string
	Display() string
}

type BaseError struct {
	message string
}

// Error implements the error interface.
func (e *BaseError) Error() string {
	return e.message
}

func (e *BaseError) Display() string {
	return e.message
}

func New(msg string) Error {
	return &BaseError{
		message: msg,
	}
}

type unexpectedError struct {
	message string
	err     error
}

func IsUnexpectedError(err Error) bool {
	_, ok := err.(*unexpectedError)
	return ok
}

func NewUnexpectedError(err error) *unexpectedError {
	return &unexpectedError{message: "An unexpected error occured", err: err}
}

func (e *unexpectedError) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.err)
}

func (e *unexpectedError) Display() string {
	return "An unexpected error occured!"
}

type SetupError struct {
	message string
}

func NewSetupError(msg string) *SetupError {
	return &SetupError{message: msg}
}

func (e *SetupError) Error() string {
	return e.message
}

func (e *SetupError) Display() string {
	return "An error occured during setup : " + e.message
}
