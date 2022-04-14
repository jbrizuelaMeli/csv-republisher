package apierrors

import (
	"fmt"
)

type GenericError struct {
	Message   string
	Code      string
	ErrorCode int
}

type CommunicationError struct {
	GenericError
	StatusCode int
}

func NewCommunicationError(message string, status int) error {
	return CommunicationError{GenericError: GenericError{Message: message}, StatusCode: status}
}

func (e CommunicationError) Error() string {
	return fmt.Sprintf(e.Message)
}
