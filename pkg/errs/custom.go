package errs

import (
	"fmt"
)

type CustomNotFoundError struct {
	Name string
}

// Deprecated: Use Error instead
type CustomError struct {
	Message string
}

func NewCustomError(message string) *CustomError {
	return &CustomError{
		Message: message,
	}
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("%v", ce.Message)
}
